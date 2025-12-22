// Copyright 2025 The TCell Authors
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

// Package mock is a simulated terminal (a terminal emulator if you will!)
// that is intended to be used for testing tcell.  As this package is for
// internal testing of tcell, it carries no stability promise.
package mock

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/vt"
	"github.com/rivo/uniseg"
)

// Cell is a representation of a display cell.
type Cell struct {
	C     []rune // Content, for now only a single rune is supported
	Fg    color.Color
	Bg    color.Color
	Attr  tcell.AttrMask
	Width int // Display width of C.
}

// MockTty is a mock terminal device.
type MockTty struct {
	Cells []Cell // Content of cells
	Rows  vt.Row
	Cols  vt.Col
	Fg    color.Color
	Bg    color.Color
	Attr  tcell.AttrMask
	X     vt.Col // cursor horizontal position
	Y     vt.Row // cursor vertical position
	Bells int    // incremented each time the bell is sounded

	ReadQ  chan any // contents of stdin
	WriteQ chan any // contents of stdout

	// These values can be overridden before Init.

	PrimaryAttributes   string // Primary device attributes, response to CSI-c
	SecondaryAttributes string // Secondary device attributes, response to CSI-c
	ExtendedAttributes  string // Extended attributes (term name, etc.) response to CSI > q
	PrivateModes        map[vt.PrivateMode]vt.ModeStatus

	inited  bool
	started bool
	stopQ   chan struct{}
	resizeQ chan<- bool
	state   int
	waitG   sync.WaitGroup
	escBuf  *bytes.Buffer
	utfBuf  []byte
	utfLen  int
}

// input states for the state machine - this parses data from the app
const (
	stateInit = iota
	stateEsc
	stateCSI
	stateOSC
	stateDCS
	stateSOS
	statePM
	stateAPC
	stateUTF // parsing a Unicode rune
	state3Fp
)

// intParams parses the CSI parameter bytes as a range of at least minNum integers,
// with default values defVal.
func intParams(str string, minNum int, defVal int) []int {
	words := strings.Split(str, ";")
	if len(words) > minNum {
		minNum = len(words)
	}
	result := make([]int, minNum)
	for i := range len(result) {
		result[i] = defVal
	}
	for i, word := range words {
		if v, e := strconv.Atoi(word); e == nil {
			result[i] = v
		}
	}
	return result
}

// moveUp movers the cursor up.
func (mt *MockTty) moveUp(n vt.Row) {
	mt.Y -= n
	if mt.Y < 0 {
		mt.Y = 0
	}
}

// moveDown moves the cursor down.
func (mt *MockTty) moveDown(n vt.Row) {
	mt.Y += n
	if mt.Y > mt.Rows-1 {
		mt.Y = mt.Rows - 1
	}
}

// moveForward moves the cursor to the right.
func (mt *MockTty) moveForward(n vt.Col) {
	mt.X += n
	if mt.X > mt.Cols-1 {
		mt.X = mt.Cols - 1
	}
}

// moveBackward moves the cursor to the left.
func (mt *MockTty) moveBackward(n vt.Col) {
	mt.X -= n
	if mt.X < 0 {
		mt.X = 0
	}
}

// eraseCell erases the cell at the given location.
func (mt *MockTty) eraseCell(x vt.Col, y vt.Row) {
	if x >= 0 && x < mt.Cols && y >= 0 && y < mt.Rows {
		ix := (int(y) * int(mt.Cols)) + int(x)
		mt.Cells[ix] = Cell{C: []rune{' '}, Width: 1, Fg: mt.Fg, Bg: mt.Bg, Attr: mt.Attr}
	}
}

// handleCsi process a fully complete CSI sequence.
func (mt *MockTty) handleCsi(final byte) {
	// NB: Technically the spec allows characters 0x3A through 0x3F to be
	// intermixed with digits, but all of the escape sequences basically have
	// at most one of these as a prefix, which when combined with the middle
	// initial, dictates the final function.
	str := mt.escBuf.String()
	mt.escBuf.Reset()
	funcId := ""
	if len(str) > 0 && str[0] >= 0x3A && str[0] <= 0x3F {
		funcId = string([]byte{str[0], final})
		str = str[1:]
	} else {
		funcId = string([]byte{final})
	}
	// insert the last intermediate byte as well -- right now we will only
	// have at most one of these (no other variants are known, though they
	// could legally exist)  This will look like e.g. "?$p" ($ is the intermediate)
	// for a private mode query.
	if len(str) > 0 && str[len(str)-1] >= 0x20 && str[len(str)-1] <= 0x2F {
		if len(funcId) == 2 {
			funcId = string([]byte{funcId[0], str[len(str)-1], final})
		} else {
			funcId = string([]byte{final, str[len(str)-1]})
		}
		str = str[0 : len(str)-1]
	}

	switch funcId {
	case "A": // up n times (CUU)
		if y := vt.Row(intParams(str, 1, 1)[0]); y >= 1 {
			mt.moveUp(y)
		}
	case "B": // down n times (CUD)
		mt.moveDown(vt.Row(intParams(str, 1, 1)[0]))

	case "C": // forward n times (CUF)
		mt.moveForward(vt.Col(intParams(str, 1, 1)[0]))

	case "D": // back n times (CUB)
		mt.moveBackward(vt.Col(intParams(str, 1, 1)[0]))

	case "E": // down n times (and reset column) (CNL)
		mt.moveDown(vt.Row(intParams(str, 1, 1)[0]))
		mt.X = 0

	case "F": // up n times (and reset column) (CPL)
		mt.moveUp(vt.Row(intParams(str, 1, 1)[0]))
		mt.X = 0

	case "G": // cursor column (CHA)
		if x := vt.Col(intParams(str, 1, 1)[0]); x > 0 && x <= mt.Cols {
			mt.X = x - 1
		}
	case "H", "f": // cursor position (CUP), also (HVP)
		if pos := vt.Row(intParams(str, 2, 1)[0]); pos > 0 && pos <= mt.Rows {
			mt.Y = pos - 1
		}
		if pos := vt.Col(intParams(str, 2, 1)[1]); pos >= 1 && pos <= mt.Cols {
			mt.X = pos - 1
		}
	case "I": // TODO: advance to next tab stop
	case "J": // erase in display (ED)
		switch intParams(str, 1, 0)[0] {
		case 0: // erase below
			for y := mt.Y + 1; y < mt.Rows; y++ {
				for x := vt.Col(0); x < mt.Cols; x++ {
					mt.eraseCell(x, y)
				}
			}
		case 1: // erase above
			for y := vt.Row(0); y < mt.Y; y++ {
				for x := vt.Col(0); x < mt.Cols; x++ {
					mt.eraseCell(x, y)
				}
			}
			// erase preceding on the same line
			for x := vt.Col(0); x < mt.X; x++ {
				mt.eraseCell(x, mt.Y)
			}
		case 2: // erase all
			for y := range mt.Rows {
				for x := range mt.Cols {
					mt.eraseCell(x, y)
				}
			}
			// others not supported (3 is erase saved lines)
		}
	case "K":
		switch intParams(str, 1, 0)[0] {
		case 0:
			for x := mt.X; x < mt.Cols; x++ {
				mt.eraseCell(x, mt.Y)
			}
		case 1:
			for x := vt.Col(0); x <= mt.X; x++ {
				mt.eraseCell(x, mt.Y)
			}
		case 2:
			for x := vt.Col(0); x < mt.Cols; x++ {
				mt.eraseCell(x, mt.Y)
			}
		}
	case "L": // TODO: insert line (IL)
	case "M": // TODO: delete line (DL)
	case "P": // TODO: delete characters
	case "c": // send primary DA
		if intParams(str, 1, 0)[0] == 0 {
			mt.reply(mt.PrimaryAttributes)
		}
	case "n":
		switch intParams(str, 1, 0)[0] {
		case 5:
			// device status ("OK")
			mt.reply("\x1b[0n")
		case 6: // cursor position report (CPR)
			mt.reply(fmt.Sprintf("\x1b[%d;%dR", mt.Y+1, mt.X+1))
		}
	case ">c":
		if intParams(str, 1, 0)[0] == 0 {
			mt.reply(mt.SecondaryAttributes)
		}
	case ">q":
		if intParams(str, 1, 0)[0] == 0 {
			mt.reply(mt.ExtendedAttributes)
		}
	case "?h": // set private mode
		pm := vt.PrivateMode(intParams(str, 1, 0)[0])
		if mt.PrivateModes[pm].Changeable() {
			mt.PrivateModes[pm] = vt.ModeOn
		}
	case "?l": // reset private mode
		pm := vt.PrivateMode(intParams(str, 1, 0)[0])
		if mt.PrivateModes[pm].Changeable() {
			mt.PrivateModes[pm] = vt.ModeOff
		}
	case "?$p": // private mode query
		pm := vt.PrivateMode(intParams(str, 1, 0)[0])
		mt.reply(pm.Reply(mt.PrivateModes[pm]))
	}
}

// handleOSC processes a fully complete OSC command.
func (mt *MockTty) handleOSC(_ string) {
	mt.state = stateInit
}

// handle3Fp handles private 3Fp cases (esc-#)
func (mt *MockTty) handle3Fp(final byte) {
	mt.state = stateInit
	switch final {
	case '8': // DECALN
		if mt.escBuf.Len() == 0 {
			for i := range len(mt.Cells) {
				mt.Cells[i] = Cell{
					C:     []rune{'E'},
					Fg:    mt.Fg,
					Bg:    mt.Bg,
					Attr:  mt.Attr,
					Width: 1,
				}
			}
		}
	}
}

// reply sends the given string to the application's stdin.
func (mt *MockTty) reply(s string) {
	for _, b := range []byte(s) {
		select {
		case mt.ReadQ <- b:
		case <-mt.stopQ:
		}
	}
}

// run the input processor.
func (mt *MockTty) run() {
	defer mt.waitG.Done()
	mt.state = stateInit
	var ch byte
	for {
		select {
		case item := <-mt.WriteQ: // read from application output
			switch v := item.(type) {
			case byte:
				ch = v
			case chan struct{}: // this was a drain checkpoint
				close(v)
				continue
			}
		case <-mt.stopQ:
			return
		}

		// calculate the position in the cells because we will use it a bit
		ix := int(mt.X) + int(mt.Y)*int(mt.Cols)

		switch mt.state {
		case stateInit:
			switch {
			case ch >= ' ' && ch < '\x7E':
				// normal ASCII, this is the easiest case
				mt.Cells[ix] = Cell{C: []rune{rune(ch)}, Fg: mt.Fg, Bg: mt.Bg, Attr: mt.Attr, Width: 1}
				if mt.X < mt.Cols-1 {
					mt.X++
				}
			case ch == '\x1B':
				mt.state = stateEsc
			case (ch & 0xE0) == 0xC0:
				mt.utfBuf = []byte{ch}
				mt.utfLen = 2
				mt.state = stateUTF
			case (ch & 0xF0) == 0xE0:
				mt.utfBuf = []byte{ch}
				mt.utfLen = 3
				mt.state = stateUTF
			case (ch & 0xF8) == 0xF0:
				mt.utfBuf = []byte{ch}
				mt.utfLen = 4
				mt.state = stateUTF
			// C0 cases
			case ch == '\x07':
				mt.Bells++
			case ch == '\x08':
				if mt.X > 0 {
					mt.X--
				}
			case ch == '\x0e': // TODO: SO
			case ch == '\x0f': // TODO: SI
			// C1 cases - note that these are all < 0xC0 so do not overlap with UTF
			case ch == 0x9b: // 8-bit CSI, not recommended but real terminals grok
				mt.state = stateCSI
				mt.escBuf.Reset()
			case ch == 0x9d:
				mt.state = stateOSC
				mt.escBuf.Reset()
			case ch == 0x9e:
				mt.state = statePM
				mt.escBuf.Reset()
			case ch == 0x9f:
				mt.state = stateAPC
				mt.escBuf.Reset()
			case ch == 0x98:
				mt.state = stateSOS
				mt.escBuf.Reset()

			// case ch == 0x84: // TODO: IND
			// case ch == 0x85: // TODO: NEL
			// case ch == 0x88: // TODO: HTS
			// case ch == 0x8e: // TODO: SSG2
			// case ch == 0x8f: // TODO: SSG3
			// case ch == 0x9a: // TODO: DECID
			default:
				// all other unknown cases, just ignore it (we don't do 8-bit decodes right now, as we are unicode only)
			}
		case stateUTF:
			if (ch & 0xC0) != 0x80 {
				// miscoded, discard this and reset the state
				mt.state = stateInit
				continue
			}
			mt.utfBuf = append(mt.utfBuf, ch)
			if len(mt.utfBuf) < mt.utfLen {
				continue
			}

			mt.state = stateInit
			// should be a full rune
			// TODO: possibly consider this as appending to the grapheme cluster?
			if r, l := utf8.DecodeRune(mt.utfBuf); l == mt.utfLen && r != utf8.RuneError {
				mt.Cells[ix] = Cell{C: []rune{r}, Fg: mt.Fg, Bg: mt.Bg, Attr: mt.Attr}
				if mt.X < mt.Cols-1 {
					mt.X++
				}
				if w := uniseg.StringWidth(string(r)); w > 1 && mt.X < mt.Cols {
					mt.Cells[ix].Width = w
					mt.Cells[ix+1] = Cell{C: nil, Fg: mt.Fg, Bg: mt.Bg, Attr: mt.Attr, Width: 0}
					if mt.X < mt.Cols-1 {
						mt.X++
					}
				}
			}
		case stateEsc:
			switch ch {
			case '[':
				mt.state = stateCSI
				mt.escBuf.Reset()
			case ']':
				mt.state = stateOSC
				mt.escBuf.Reset()
			case 'D': // down one line
				if mt.Y < mt.Rows-1 {
					mt.Y++
				}
			case 'E':
				if mt.Y < mt.Rows-1 {
					mt.Y++
				}
				mt.X = 0
			case 'H': // TODO: set tab stop
			case 'M':
				if mt.Y > 0 {
					mt.Y--
				}
			case 'N': // TODO: SS2: set G2 (not doing)
			case 'O': // TODO: SS3: set G3
			case 'P':
				mt.state = stateDCS
				mt.escBuf.Reset()
			case 'X':
				mt.state = stateSOS
				mt.escBuf.Reset()
			case 'Z':
				mt.reply(mt.PrimaryAttributes)
			case '\\': // TODO: end of string
			case '^':
				mt.state = statePM
				mt.escBuf.Reset()
			case '_':
				mt.state = stateAPC
				mt.escBuf.Reset()
			case '6': // TODO: move left one, possibly moving data
			case '7': // TODO: save state
			case '8': // TODO: restore state
			case '9': // TODO: move right one, possibly moving data
			case '=': // TODO: application keypad mode
			case '>': // TODO: normal keypad mode
			case 'c': // TODO: terminal reset
			case 'g':
				mt.Bells++
			case 'n': // TODO: LS2
			case 'o': // TODO: LS3
			case '|': // TODO LS3R
			case '}': // TODO LS2R
			case '#':
				mt.state = state3Fp
				mt.escBuf.Reset()
			case '$', '(', ')', '*', '+': // TODO: select character set
			}

		case stateCSI:
			if (ch >= 0x30) && (ch <= 0x3F) {
				mt.escBuf.WriteByte(ch) // parameter bytes
			} else if (ch >= 0x20) && (ch <= 0x2F) {
				mt.escBuf.WriteByte(ch) // intermediate bytes
			} else if ch >= 0x40 && (ch <= 0x7F) {
				mt.state = stateInit
				mt.handleCsi(ch)
			}

		case stateOSC:
			if ch == '\x9c' {
				mt.handleOSC(mt.escBuf.String())
			} else if buf := mt.escBuf.Bytes(); len(buf) > 0 && buf[len(buf)-1] == '\x1b' && ch == '\\' {
				buf = buf[:len(buf)-1]
				mt.handleOSC(string(buf))
			} else {
				mt.escBuf.WriteByte(ch)
			}

		case statePM, stateSOS, stateAPC, stateDCS: // none of these have handling for now
			if ch == '\x9c' {
				mt.state = stateInit
			} else if buf := mt.escBuf.Bytes(); len(buf) > 0 && buf[len(buf)-1] == '\x1b' && ch == '\\' {
				mt.state = stateInit
			} else {
				mt.escBuf.WriteByte(ch)
			}
		case state3Fp:
			if ch >= 0x20 && ch <= 0x2F {
				mt.escBuf.WriteByte(ch)
			} else if ch >= 0x30 && ch <= 0x3F {
				mt.handle3Fp(ch)
			} else {
				// unexpected
				mt.state = stateInit
			}
		}
	}
}

func (mt *MockTty) Start() error {
	mt.started = true
	mt.stopQ = make(chan struct{})
	mt.waitG.Add(1)
	go mt.run()
	return nil
}

func (mt *MockTty) Stop() error {
	select {
	case <-mt.stopQ:
	default:
		close(mt.stopQ)
	}
	mt.waitG.Wait()
	mt.started = false
	return nil
}

func (mt *MockTty) Read(b []byte) (int, error) {

	// This implements VMIN 1, VTIME 0...
	// I.e. it will wait for at least one character,
	// then it reads as many as it can (up to the requested amount).
	// It returns once at least one character is available, unless
	// it is aborted.
	if len(b) == 0 {
		return 0, nil
	}
	gotOne := false
	for !gotOne {
		select {
		// case b[0] = <-mt.ReadQ:
		case v := <-mt.ReadQ:
			switch v := v.(type) {
			case byte:
				b[0] = v
				gotOne = true
			case chan struct{}:
				close(v)
				return 0, nil
			case bool:
				return 0, nil
			}
		case <-mt.stopQ:
			return 0, nil
		}
	}

	n := 1

	for n < len(b) {
		select {
		case v := <-mt.ReadQ:
			switch v := v.(type) {
			case byte:
				b[n] = v
				n++
			case chan struct{}:
				close(v)
			}
		default:
			return n, nil
		}
	}
	return len(b), nil
}

func (mt *MockTty) Write(b []byte) (int, error) {
	for n := range b {
		select {
		case mt.WriteQ <- b[n]:
		case <-mt.stopQ:
			return n, nil
		}
	}
	return len(b), nil
}

func (mt *MockTty) Close() error { return nil }

func (mt *MockTty) Drain() error {

	// inject a checkpoint to make sure the entire
	// write q to this point is done - but we have a
	wcq := make(chan struct{})
	mt.WriteQ <- wcq
	select {
	case <-wcq:
	case <-time.After(time.Millisecond * 200):
	case <-mt.stopQ:
	}

	// readq side, we don't need to wait, as the
	// read call will be interrupted which is all that is needed
	select {
	case mt.ReadQ <- true:
	case <-mt.stopQ:
	}

	return nil
}

func (mt *MockTty) NotifyResize(rq chan<- bool) {
	mt.resizeQ = rq
}

func (mt *MockTty) WindowSize() (tcell.WindowSize, error) {
	return tcell.WindowSize{Height: int(mt.Rows), Width: int(mt.Cols)}, nil
}

// Reset is not part of the TTY interface, and is for testing.
func (mt *MockTty) Reset() {
	if !mt.inited {
		mt.inited = true
		if mt.Rows == 0 {
			mt.Rows = 24
		}
		if mt.Cols == 0 {
			mt.Cols = 80
		}
		mt.Cells = make([]Cell, int(mt.Cols)*int(mt.Rows))
		mt.Fg = color.White
		mt.Bg = color.Black
		mt.Attr = tcell.AttrNone
		mt.PrimaryAttributes = "\x1b[?62;1;22;52c"
		mt.SecondaryAttributes = "\x1b[>1;10c"
		mt.ExtendedAttributes = "\x1b[P>|simtty 0.1.2\x1b\\"
		mt.WriteQ = make(chan any, 128)
		mt.ReadQ = make(chan any, 128)
		mt.state = stateInit
		mt.escBuf = &bytes.Buffer{}
		mt.PrivateModes = map[vt.PrivateMode]vt.ModeStatus{
			vt.PmAutoMargin:    vt.ModeOffLocked, // forced auto margin
			vt.PmFocusReports:  vt.ModeOff,       // focus mode
			vt.PmResizeReports: vt.ModeOff,       // resize reports
		}
	}
}
