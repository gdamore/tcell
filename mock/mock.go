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
	"github.com/gdamore/tcell/v3"
)

type Cell struct {
	C     []rune // Content, for now only a single rune is supported
	Fg    tcell.Color
	Bg    tcell.Color
	Attr  tcell.AttrMask
	Width int // Display width of C.
}

type MockTty struct {
	Cells []Cell // Content of cells
	Rows  int
	Cols  int
	Fg    tcell.Color
	Bg    tcell.Color
	attr  tcell.AttrMask
	X     int // cursor position
	Y     int // cursor position

	ReadQ  chan byte // contents of stdin
	WriteQ chan byte // contents of stdout

	// These values can be overridden before Init.

	PrimaryAttributes  string // Primary device attributes, response to CSI-c
	ExtendedAttributes string // Extended attributes (term name, etc.) response to CSI > q

	inited  bool
	started bool
	stopQ   chan struct{}
	resizeQ chan<- bool
}

func (mt *MockTty) Start() error {
	mt.started = true
	mt.stopQ = make(chan struct{})
	return nil
}

func (mt *MockTty) Stop() error {
	select {
	case <-mt.stopQ:
	default:
		close(mt.stopQ)
	}
	mt.started = false
	return nil
}

func (mt *MockTty) Read(b []byte) (int, error) {

	for n := range len(b) {
		select {
		case b[n] = <-mt.ReadQ:
		case <-mt.stopQ:
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
	close(mt.stopQ)
loop:
	for {
		select {
		case <-mt.ReadQ:
		default:
			break loop
		}
	}
	return nil
}

func (mt *MockTty) NotifyResize(rq chan<- bool) {
	mt.resizeQ = rq
}

func (mt *MockTty) WindowSize() (tcell.WindowSize, error) {
	return tcell.WindowSize{Height: mt.Rows, Width: mt.Cols}, nil
}

// Reset is not part of the TTY interface, and is for testing.
func (mt *MockTty) Reset() {
	if !mt.inited {
		mt.inited = true
		mt.Rows = 24
		mt.Cols = 80
		mt.Cells = make([]Cell, mt.Cols*mt.Rows)
		mt.Fg = tcell.ColorWhite
		mt.Bg = tcell.ColorBlack
		mt.attr = tcell.AttrNone
		mt.PrimaryAttributes = "\x1b[?62;1;7;22c"
		mt.ExtendedAttributes = "\x1b[P>|simtty 0.1.2\x1b\\"
		mt.WriteQ = make(chan byte, 128)
		mt.ReadQ = make(chan byte, 128)
	}
}
