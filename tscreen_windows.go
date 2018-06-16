// +build windows

// Copyright 2015 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// we don't have a signal type in windows to send to the 'sigwinch' handler, so here's one
type mySignal struct{}

func (mySignal) String() string { return "" }
func (mySignal) Signal()        {}

// info pack returned from `stty size speed` command
// also reused by the VT console mode just for storing the dimensions
type ttyinfo struct {
	w, h, speed int
}

type tscreenMode int

const (
	tscreenModeNone = iota

	// use stty and msys pipes back to the shell
	tscreenModeMsys

	// use windows >= 10.0.10586 native console virtual terminal mode (added for WSL most likely)
	tscreenModeVTConsole
)

type termiosPrivateMsys struct {
	// name of cygwin/msys pty we're attached to
	ptyname string

	// memory of tty state; restore on exit
	stty_blob string
}

type termiosPrivateVT struct {
	// stdin and stdout, used by tscreenModeVTConsole
	conin_file, conout_file *os.File

	// console handles, needed for tscreenModeVTConsole
	conin, conout syscall.Handle
}

// private data
type termiosPrivate struct {
	// which of the major operating modes we're in
	mode tscreenMode

	// data for msys mode
	msys termiosPrivateMsys

	// data for VT Console mode
	vt termiosPrivateVT

	// used to kill a monitor goroutines when this tScreen implementation is finalized
	termSignal chan bool

	// acks the termSignal
	termSignalAck chan bool

	// whether we successfully initialized (we'll have some methods called, even if we didn't init correctly?)
	initialized bool

	// last cached info, so we have a uniform sense of the terminal dimensions
	info ttyinfo
}

// Just a wrapper for the DuplicateHandle syscall
func duplicateHandle(handle uintptr) uintptr {
	var dup uintptr
	p, _ := syscall.GetCurrentProcess()
	err := syscall.DuplicateHandle(p, syscall.Handle(handle), p, (*syscall.Handle)(&dup), 0, true, syscall.DUPLICATE_SAME_ACCESS)
	if err != nil {
		panic("Couldn't duplicate stdio handle")
	}
	return dup
}

func (t *tScreen) termioPreInit() (termvar string, err error) {

	// We'd have a terminfo for cygwin, but based on my study,
	// $TERM=cygwin is only used by the cygwin bash which is using the windows console
	// On pre-WSL systems, we'll have to handle these with the console screen

	termvar = os.Getenv("TERM")

	// start building private data to save work later
	t.tiosp = &termiosPrivate{}
	priv := t.tiosp

	//check if we have msys pipes to
	ptynum := GetMSYSTerminal(os.Stdout) //stdin breaks things for some reason, in GetFileInformationByHandleEx
	if ptynum != -1 {
		t.tiosp.mode = tscreenModeMsys
		t.tiosp.msys.ptyname = fmt.Sprintf("/dev/pty%d", ptynum)
		return termvar, nil
	}

	// ---------- WSL ANALYSIS ----------
	// in WSL bash I found:
	//  inputmode =  $3D8 - ENABLE_VIRTUAL_TERMINAL_INPUT | ENABLE_AUTO_POSITION  | ENABLE_EXTENDED_FLAGS | ENABLE_QUICK_EDIT_MODE  | ENABLE_MOUSE_INPUT | ENABLE_WINDOW_INPUT
	//  outputmode = $007 - ENABLE_PROCESSED_OUTPUT | ENABLE_WRAP_AT_EOL_OUTPUT | ENABLE_VIRTUAL_TERMINAL_PROCESSING
	// in windows 10 cmd.exe I found:
	//  inputmode =  $1F7 - ENABLE_AUTO_POSITION | ENABLE_EXTENDED_FLAGS | ENABLE_QUICK_EDIT_MODE | ENABLE_INSERT_MODE | ENABLE_MOUSE_INPUT | ENABLE_ECHO_INPUT | ENABLE_LINE_INPUT | ENABLE_PROCESSED_INPUT
	//  outputmode = $003 - ENABLE_PROCESSED_OUTPUT | ENABLE_WRAP_AT_EOL_OUTPUT
	// in windows 7 while debugging in vscode, I found:
	//  inputmode =  $1B7 - ENABLE_AUTO_POSITION | ENABLE_EXTENDED_FLAGS | ENABLE_INSERT_MODE | ENABLE_MOUSE_INPUT | ENABLE_ECHO_INPUT | ENABLE_LINE_INPUT | ENABLE_PROCESSED_INPUT
	//  outputmode = $003 - ENABLE_PROCESSED_OUTPUT | ENABLE_WRAP_AT_EOL_OUTPUT
	// in windows 7 cmd.exe, I found:
	//  inputmode =  $1A7 - ENABLE_AUTO_POSITION | ENABLE_EXTENDED_FLAGS | ENABLE_INSERT_MODE | ENABLE_ECHO_INPUT | ENABLE_LINE_INPUT | ENABLE_PROCESSED_INPUT
	//  outputmode = $003 - ENABLE_PROCESSED_OUTPUT | ENABLE_WRAP_AT_EOL_OUTPUT
	// I really don't understand why the difference between all these, but it seems clear that:
	// 1. From evidence, we can use VIRTUAL_TERMINAL_PROCESSING and/or ENABLE_VIRTUAL_TERMINAL_INPUT to know we're launched from bash
	// 2. From logic, if the console is telling us it's a virtual terminal, why not use it that way?
	// see here for some related https://stackoverflow.com/questions/46030331/enable-virtual-terminal-processing-and-disable-newline-auto-return-failing
	// ----------------------------------

	//try opening the consoles and check the input mode. that's how we know whether we can activate tscreenModeVTConsole
	var inmode uint32
	priv.vt.conin, err = syscall.Open("CONIN$", syscall.O_RDWR, 0)
	priv.vt.conout, err = syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	_, _, err = procGetConsoleMode.Call(uintptr(priv.vt.conin), uintptr(unsafe.Pointer(&inmode)))
	if inmode&0x0200 != 0 { //ENABLE_VIRTUAL_TERMINAL_INPUT
		priv.mode = tscreenModeVTConsole

		//the docs say something like it's based on xterm, but since this is what bash says, we'll go with this
		//(might just be xterm though, I don't know)
		return "xterm-256color", nil
	}

	//never mind, we don't want to deal with these here
	syscall.Close(priv.vt.conin)
	syscall.Close(priv.vt.conout)

	// sometimes cygwin will be running on old windows; we need to handle that with the console.
	// since we do have a cygwin entry in our terminfo, we return "" here instead to kick ourselves to the console screen
	if termvar == "cygwin" {
		return "", nil
	}

	// we don't know what to do with any other termvars yet. give up, I guess

	return "", nil
}

func (t *tScreen) initModeMsys() error {
	// General plan: detect whether we're in a cygwin-type environment (or msys) that we can control with stty
	// return ErrNoScreen if we aren't (and windows console IO will take over)

	var err error
	priv := t.tiosp

	priv.termSignal = make(chan bool, 1)
	priv.termSignalAck = make(chan bool, 1)

	// try getting dump of current tty state.
	// if this fails, stty isn't available at all, and we can't operate a terminal in this way
	var stty_blob strings.Builder
	cmdGrabBlob := exec.Command("stty", "-F", priv.msys.ptyname, "-g")
	cmdGrabBlob.Stdout = &stty_blob
	err = cmdGrabBlob.Run()
	if err != nil {
		return ErrNoScreen
	}

	// tcsetattr-like operation to get a more raw mode in the tty
	cmdSetRaw := exec.Command("stty", "-F", priv.msys.ptyname, "raw", "-echo")
	cmdSetRaw.Stdout = os.Stdout
	err = cmdSetRaw.Run()
	if err != nil {
		return ErrNoScreen
	}

	// fetch initial parameters. in case stty is malfunctioning and returning garbage, this may return an error
	priv.info, err = t.sttyPoll()
	if err != nil {
		return ErrNoScreen
	}

	// apply initial size and speed param
	t.cells.Resize(priv.info.w, priv.info.h)
	t.baud = priv.info.speed

	// the finale -- due to our precautions, sending terminal data through stdin/stdout will now work
	// however, we need to duplicate the handles, because tscreen is going to expect the handles to be closed later
	// (else its input loop never terminates)
	// The consequence of not duplicating is your 2nd tscreen will get its input mixed up with the disappeared 1st tscreen
	t.in = os.NewFile(duplicateHandle(os.Stdin.Fd()), "stdin")
	t.out = os.NewFile(duplicateHandle(os.Stdout.Fd()), "stdout")

	// start the monitor for terminal size changes
	go t.goMonitorSize()

	// things are looking good. finish up
	t.tiosp.msys.stty_blob = stty_blob.String()
	t.tiosp.initialized = true

	return nil
}

func (t *tScreen) initModeVTConsole() error {

	priv := t.tiosp

	// just use conin and conout
	// we can't use the handles we already have because we need a file
	// (maybe we could build a file around the handle? not sure)
	in, e := os.OpenFile("CONIN$", syscall.O_RDONLY, 0)
	if e != nil {
		return e
	}
	out, e := os.OpenFile("CONOUT$", syscall.O_WRONLY, 0)
	if e != nil {
		in.Close()
		return e
	}

	priv.vt.conin_file = in
	priv.vt.conout_file = out
	t.in = in
	t.out = out

	priv.info, e = t.vtPoll()
	if e != nil {
		return e
	}
	t.cells.Resize(priv.info.w, priv.info.h)
	t.baud = 38400 //that's what stty from WSL bash said

	t.tiosp.initialized = true

	// start the monitor for terminal size changes
	go t.goMonitorSize()

	return nil
}

func (t *tScreen) termioInit() error {
	switch t.tiosp.mode {
	case tscreenModeMsys:
		if err := t.initModeMsys(); err != nil {
			return err
		}
	case tscreenModeVTConsole:
		if err := t.initModeVTConsole(); err != nil {
			return err
		}
	default:
		return errors.New("tscreen tried to init with no good mode")
	}
	return nil
}

func (t *tScreen) doUpdateSize() {
	var err error
	priv := t.tiosp
	lastInfo := priv.info

	if priv.mode == tscreenModeMsys {
		// poll information through stty; send to tcell's resize signal handler
		priv.info, err = t.sttyPoll()
	} else {
		priv.info, err = t.vtPoll()
	}

	if err != nil {
		panic(err)
	}

	if lastInfo != priv.info {
		t.sigwinch <- mySignal{}
	}
}

func (t *tScreen) goMonitorSize() {

	// It's uncertain what's the ideal value here
	// I just picked a value that looked OK to a human but was as long as possible
	ticker := time.NewTicker(time.Millisecond * 250)
	defer ticker.Stop()

LOOP:
	for {
		select {
		case <-ticker.C:
			t.doUpdateSize()

		case <-t.tiosp.termSignal:
			break LOOP
		}
	}

	t.tiosp.termSignalAck <- true
}

func (t *tScreen) vtPoll() (ret ttyinfo, err error) {
	// not an accident, we do reuse the ttyinfo struct
	priv := t.tiosp
	info := consoleInfo{}
	procGetConsoleScreenBufferInfo.Call(uintptr(priv.vt.conout), uintptr(unsafe.Pointer(&info)))
	ret.w = int((info.win.right - info.win.left) + 1)
	ret.h = int((info.win.bottom - info.win.top) + 1)
	return
}

func (t *tScreen) sttyPoll() (ret ttyinfo, err error) {

	priv := t.tiosp

	// return nice errors from here, because stty could malfunction, i guess

	var outbuf strings.Builder
	cmd := exec.Command("stty", "-F", priv.msys.ptyname, "size", "speed")
	cmd.Stdout = &outbuf
	err = cmd.Run()
	if err != nil {
		return
	}
	result := outbuf.String()

	// parse what stty returns in this format:
	// rows cols
	// baud

	// 1. split lines
	lines := strings.Split(result, "\n")
	if len(lines) != 3 { // (yes, we get an empty line)
		return ret, fmt.Errorf("stty returned invalid format")
	}

	// 2. split the size line
	fields := strings.Fields(lines[0])
	if err != nil || len(fields) != 2 {
		return ret, fmt.Errorf("stty returned invalid format")
	}

	// strings to integers
	var errw, errh, errs error
	ret.w, errh = strconv.Atoi(fields[1]) //not a mistake, we get rows and cols (not width and height) from stty
	ret.h, errw = strconv.Atoi(fields[0])
	ret.speed, errs = strconv.Atoi(lines[1])
	if errw != nil || errh != nil || errs != nil {
		return ret, fmt.Errorf("stty returned invalid format")
	}

	return
}

func (t *tScreen) getWinSize() (int, int, error) {

	priv := t.tiosp
	if !priv.initialized {
		return 0, 0, fmt.Errorf("calling getWinSize on uninitialized tscreen_windows part")
	}

	return priv.info.w, priv.info.h, nil
}

func (t *tScreen) termioPreFini() {

	priv := t.tiosp
	if !priv.initialized {
		return
	}

	if priv.mode == tscreenModeMsys {
		// tell monitor thread to stop and wait for ack
		// this is so we make sure we arent sending resizes while trying to shut down
		t.tiosp.termSignal <- true
		<-t.tiosp.termSignalAck
	}
}

func (t *tScreen) termioFini() {

	priv := t.tiosp
	if !priv.initialized {
		return
	}

	if priv.mode == tscreenModeMsys {

		// this waits for the tscreen main loop to exit...
		<-t.indoneq

		// restore previous terminal state
		cmd := exec.Command("stty", "-F", priv.msys.ptyname, priv.msys.stty_blob)
		cmd.Stdout = os.Stdout
		cmd.Run()

		// do something to close duplicated handles?
		// but there's a problem... with this, something seems to wait for a \n before proceeding
		// skip it for now. input handler threads will die after they inadvertently eat one character
		// t.in.Close()
		// t.out.Close()
		// syscall.CloseHandle(syscall.Handle(t.in.Fd()))
		// syscall.CloseHandle(syscall.Handle(t.out.Fd()))
	}

	if priv.mode == tscreenModeVTConsole {
		//see "race condition" discussion in console_win.go. Is this another solution?
		//Without this, the outstanding read on conin will hang even after this handle is closed
		syscall.CancelIoEx(syscall.Handle(priv.vt.conin_file.Fd()), nil)

		syscall.CloseHandle(priv.vt.conout)
		syscall.CloseHandle(priv.vt.conin)

		priv.vt.conout_file.Close()
		priv.vt.conin_file.Close()
	}

	priv.initialized = false
}
