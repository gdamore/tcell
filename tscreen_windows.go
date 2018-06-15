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
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// we don't have a signal type in windows to send to the 'sigwinch' handler, so here's one
type mySignal struct{}

func (mySignal) String() string { return "" }
func (mySignal) Signal()        {}

// info pack returned from `stty size speed` command
type sttyInfo struct {
	w, h, speed int
}

// private data
type termiosPrivate struct {

	// used to kill a monitor goroutines when this tScreen implementation is finalized
	termSignal chan bool

	// acks the termSignal
	termSignalAck chan bool

	// whether we successfully initialized (we'll have some methods called, even if we didn't init correctly?)
	initialized bool

	// name of cygwin/msys pty we're attached to
	ptyname string

	// last cached info, so we have a uniform sense of the terminal dimensions
	info sttyInfo

	// memory of tty state; restore on exit
	stty_blob string
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

func (t *tScreen) termioInit() error {

	// General plan: detect whether we're in a cygwin-type environment (or msys) that we can control with stty
	// return ErrNoScreen if we aren't (and windows console IO will take over)

	var err error

	// First, see if we're in cygwin and what pty we're on
	// we need to know that because a normal windows program (as a Go program is) isn't linked against any
	// cygwin parts that give it wisdom about the cygwin environment.
	// However, we have some wisdom, and can see if the stdio we were provided are a pipe controlled by cygwin
	// If we know what pty to address, the cygwin tools know how to take it from there
	// (all the cygwin processes coordinate with a bunch of global objects)
	// We can be pretty sure stty is accessible, because we've inherited a cygwin $PATH with a /usr/bin/stty in it
	ptynum := GetMSYSTerminal(os.Stdout) //stdin breaks things for some reason, in GetFileInformationByHandleEx
	if ptynum == -1 {
		return ErrNoScreen
	}

	// we're going to need ptyname and generally depend on building ourselves up as we go, so prep our private data store
	t.tiosp = &termiosPrivate{
		termSignal:    make(chan bool, 1),
		termSignalAck: make(chan bool, 1),
		ptyname:       fmt.Sprintf("/dev/pty%d", ptynum),
	}

	// try getting dump of current tty state.
	// if this fails, stty isn't available at all, and we can't operate a terminal in this way
	var stty_blob strings.Builder
	cmdGrabBlob := exec.Command("stty", "-F", t.tiosp.ptyname, "-g")
	cmdGrabBlob.Stdout = &stty_blob
	err = cmdGrabBlob.Run()
	if err != nil {
		return ErrNoScreen
	}

	// tcsetattr-like operation to get a more raw mode in the tty
	cmdSetRaw := exec.Command("stty", "-F", t.tiosp.ptyname, "raw", "-echo")
	cmdSetRaw.Stdout = os.Stdout
	err = cmdSetRaw.Run()
	if err != nil {
		return ErrNoScreen
	}

	// fetch initial parameters. in case stty is malfunctioning and returning garbage, this may return an error
	tmp, err := t.sttyPoll()
	if err != nil {
		return ErrNoScreen
	}

	// apply initial size and speed param
	t.cells.Resize(tmp.w, tmp.h)
	t.baud = tmp.speed

	// things are looking good. finish up
	t.tiosp.stty_blob = stty_blob.String()
	t.tiosp.initialized = true

	// start the monitor for terminal size changes
	go t.monitorSize()

	// the finale -- due to our precautions, sending terminal data through stdin/stdout will now work
	// however, we need to duplicate the handles, because tscreen is going to expect the handles to be closed later
	// (else its input loop never terminates)
	// The consequence of not duplicating is your 2nd tscreen will get its input mixed up with the disappeared 1st tscreen
	t.in = os.NewFile(duplicateHandle(os.Stdin.Fd()), "stdin")
	t.out = os.NewFile(duplicateHandle(os.Stdout.Fd()), "stdout")

	return nil
}

func (t *tScreen) doUpdateSize() {

	// poll information through stty; send to tcell's resize signal handler

	var err error

	lastInfo := t.tiosp.info

	t.tiosp.info, err = t.sttyPoll()
	if err != nil {
		panic(err)
	}

	if lastInfo != t.tiosp.info {
		t.sigwinch <- mySignal{}
	}
}

func (t *tScreen) monitorSize() {

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

func (t *tScreen) sttyPoll() (ret sttyInfo, err error) {

	// return nice errors from here, because stty could malfunction, i guess

	var outbuf strings.Builder
	cmd := exec.Command("stty", "-F", t.tiosp.ptyname, "size", "speed")
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

	if !t.tiosp.initialized {
		return
	}

	// tell monitor thread to stop and wait for ack
	// this is so we make sure we arent sending resizes while trying to shut down
	t.tiosp.termSignal <- true
	<-t.tiosp.termSignalAck
}

func (t *tScreen) termioFini() {

	if !t.tiosp.initialized {
		return
	}

	// this waits for the tscreen main loop to exit...
	<-t.indoneq

	// restore previous terminal state
	cmd := exec.Command("stty", "-F", t.tiosp.ptyname, t.tiosp.stty_blob)
	cmd.Stdout = os.Stdout
	cmd.Run()

	// do something to close duplicated handles?
	// but there's a problem... with this, something seems to wait for a \n before proceeding
	// skip it for now. input handler threads will die after they inadvertently eat one character
	// t.in.Close()
	// t.out.Close()
	// syscall.CloseHandle(syscall.Handle(t.in.Fd()))
	// syscall.CloseHandle(syscall.Handle(t.out.Fd()))

	t.tiosp.initialized = false
}
