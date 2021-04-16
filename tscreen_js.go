// Copyright 2021 The TCell Authors
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

// +build js

package tcell

import (
	"errors"
	"os"
	"strings"
	"sync"
	"syscall/js"
)

// engage is used to place the terminal in raw mode and establish screen size, etc.
// Thing of this is as tcell "engaging" the clutch, as it's going to be driving the
// terminal interface.
func (t *tScreen) engage() error {
	t.Lock()
	defer t.Unlock()
	if t.stopQ != nil {
		return errors.New("already engaged")
	}
	if w, h, err := t.getWinSize(); err == nil && w != 0 && h != 0 {
		t.cells.Resize(w, h)
	}
	stopQ := make(chan struct{})
	t.stopQ = stopQ
	t.nonBlocking(false)
	t.enableMouse(t.mouseFlags)
	t.enablePasting(t.pasteEnabled)
	priv := t.privateData.(*jsPrivate)
	if priv.termType == "xterm.js" {
		priv.jsterm.Call("onResize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			t.sigwinch <- nil
			return nil
		}))
	}

	t.wg.Add(2)
	go t.inputLoop(stopQ)
	go t.mainLoop(stopQ)
	return nil
}

// disengage is used to release the terminal back to support from the caller.
// Think of this as tcell disengaging the clutch, so that another application
// can take over the terminal interface.  This restores the TTY mode that was
// present when the application was first started.
func (t *tScreen) disengage() {

	t.Lock()
	t.nonBlocking(true)
	stopQ := t.stopQ
	t.stopQ = nil
	close(stopQ)
	t.Unlock()

	// wait for everything to shut down
	t.wg.Wait()

	priv := t.privateData.(*jsPrivate)
	if priv.termType == "xterm.js" {
		priv.jsterm.Call("onResize", js.Null())
	}

	// put back normal blocking mode
	t.nonBlocking(false)

	// shutdown the screen and disable special modes (e.g. mouse and bracketed paste)
	ti := t.ti
	t.cells.Resize(0, 0)
	t.TPuts(ti.ShowCursor)
	t.TPuts(ti.AttrOff)
	t.TPuts(ti.Clear)
	t.TPuts(ti.ExitCA)
	t.TPuts(ti.ExitKeypad)
	t.enableMouse(0)
	t.enablePasting(false)

}

type jsPrivate struct {
	termType string
	jsterm   js.Value
}

func enosys() js.Error {
	val := js.ValueOf(map[string]interface{}{
		"message": js.ValueOf("not implemented"),
		"code":    js.ValueOf("ENOSYS"),
	})
	return js.Error{Value: val}
}

// initialize is used at application startup, and sets up the initial values
// including file descriptors used for terminals and saving the initial state
// so that it can be restored when the application terminates.
// For WASM, the calling thread blocks until "setupTcell" is called from the javascript
// This call passes in the terminal type and the terminal object.
// Then we set up global.fs with the appropriate read and write methods
/*
	var term = new Terminal();
	const fitAddon = new FitAddon.FitAddon();
	term.loadAddon(fitAddon);
	term.open(document.getElementById('terminal'));
	global.term = term;
	fitAddon.fit();
	function resizeFit() {
		fitAddon.fit()
	}
	window.onresize = resizeFit

	const go = new Go();
	go.env.TERM = "xterm-256color";
	WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
		go.run(result.instance);
		setupTcell("xterm.js", term);
	});
*/
func (t *tScreen) initialize() error {
	waitForSetup := make(chan bool)
	js.Global().Set("setupTcell", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Args: type, term, [option]
		// The only supported type is "xterm.js"
		// if option is "NO_FS", fs.read and fs.write are not set.
		//   This is for the case where other file IO is being performed, and we don't want to stomp on filesystem code.
		//   For this case, other code will need to attach stdin and stdout to the terminal.
		if len(args) < 2 {
			result := map[string]interface{}{
				"error": "Invalid no of arguments passed",
			}
			return result
		}
		if t.privateData != nil {
			result := map[string]interface{}{
				"error": "Tcell setup is completed",
			}
			return result
		}
		var priv jsPrivate
		t.privateData = &priv
		priv.termType = args[0].String()
		if priv.termType != "xterm.js" {
			result := map[string]interface{}{
				"error": "Invalid terminal type",
			}
			return result
		}
		priv.jsterm = args[1]
		setFS := true
		if len(args) > 2 && args[2].String() == "NO_FS" {
			setFS = false
		}
		if setFS {
			fs := js.Global().Get("fs")

			// Reading from os.Stdin
			inQ := make(chan []byte, 16)
			var inBuf []byte
			var inM sync.Mutex
			// xterm.js: "onData" callback to put bytes on the queue
			priv.jsterm.Call("onData", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				inQ <- []byte(args[0].String())
				return nil
			}))
			// When reading os.Stdin, pull bytes off the queue.
			fs.Set("read", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				if len(args) != 6 {
					panic("Not enough arguments to read()")
				}
				// args: fd, buffer, offset, length, position, callback
				// this is an asynchronous read.  Launch a goroutine to wait for data.
				go func() {
					fd := args[0].Int()
					// If we aren't reading stdin, we can't
					if fd != 0 {
						args[5].Invoke(enosys(), 0)
						return
					}
					// The default implementation errors if we try to do anything fancy.
					length := args[1].Length()
					if args[2].Int() != 0 || args[3].Int() != length || !args[4].Equal(js.Null()) {
						args[5].Invoke(enosys(), 0)
						return
					}
					// Just in case of concurrent reads
					inM.Lock()
					defer inM.Unlock()
					for len(inQ) > 0 {
						inBuf = append(inBuf, <-inQ...)
					}
					// Block if the buffer is empty
					for len(inBuf) == 0 && length > 0 {
						inBuf = <-inQ
					}
					// Return only what we have
					if len(inBuf) < length {
						length = len(inBuf)
					}
					// Copy into the buffer argument
					js.CopyBytesToJS(args[1], inBuf)
					inBuf = inBuf[length:]
					args[5].Invoke(js.Null(), length)
				}()
				return nil
			}))

			// Writing to os.Stdout and os.Stderr
			var doWrite func(data js.Value) int
			// In case of concurrent writes.
			var outM sync.Mutex
			// xterm.js: Write to term
			doWrite = func(data js.Value) int {
				outM.Lock()
				defer outM.Unlock()
				priv.jsterm.Call("write", data)
				return data.Length()
			}
			// stderr gets sent to the console
			consoleOut := ""
			//var uint8Array = js.Global().Get("Uint8Array")
			doLog := func(data js.Value) int {
				outM.Lock()
				defer outM.Unlock()
				bytes := make([]byte, data.Length())
				js.CopyBytesToGo(bytes, data)
				consoleOut += string(bytes)
				for i := strings.Index(consoleOut, "\n"); i != -1; i = strings.Index(consoleOut, "\n") {
					js.Global().Get("console").Call("log", consoleOut[:i])
					consoleOut = consoleOut[i+1:]
				}
				return data.Length()
			}
			// Write to os.Stdout by calling doWrite
			fs.Set("writeSync", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				if len(args) != 2 {
					panic("Not enough arguments to writesync()")
				}
				// Args: fd, buffer
				fd := args[0].Int()
				// If we aren't writing os.Stdout, write to the console
				switch fd {
				case 1:
					return doWrite(args[1])
				case 2:
					return doLog(args[1])
				}
				return enosys()
			}))
			fs.Set("write", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				if len(args) != 6 {
					panic("Not enough arguments to write()")
				}
				// Args: fd, buf, offset, length, position, callback
				go func() {
					// The default implementation errors if we try to do anything fancy.
					if args[2].Int() != 0 || args[3].Int() != args[1].Length() || !args[4].Equal(js.Null()) {
						args[5].Invoke(enosys(), 0)
						return
					}
					fd := args[0].Int()
					// If we aren't writing os.Stdout, write to the console
					switch fd {
					case 1:
						args[5].Invoke(nil, doWrite(args[1]))
					case 2:
						args[5].Invoke(nil, doLog(args[1]))
					default:
						args[5].Invoke(enosys(), 0)
					}
				}()
				return nil
			}))
		}
		waitForSetup <- true
		return nil
	}))
	<-waitForSetup
	t.out = os.Stdout
	t.in = os.Stdin
	return nil
}

// finalize is used at application shutdown, and restores the terminal
// to it's initial state.  It should not be called more than once.
func (t *tScreen) finalize() {
	t.disengage()
}

// getWinSize is called to obtain the terminal dimensions.
func (t *tScreen) getWinSize() (int, int, error) {
	priv := t.privateData.(*jsPrivate)
	if priv.termType == "xterm.js" {
		return priv.jsterm.Get("cols").Int(), priv.jsterm.Get("rows").Int(), nil
	}
	return 80, 24, nil
}

// Beep emits a beep to the terminal.
func (t *tScreen) Beep() error {
	t.writeString(string(byte(7)))
	return nil
}
