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

//go:build !js && !wasm
// +build !js,!wasm

package mock

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/vt"
)

// TestMockStart just verifies that we can start the mock tty.
func TestMockStart(t *testing.T) {
	mt := &MockTty{}
	mt.Reset()
	screen, err := tcell.NewTerminfoScreenFromTty(mt)
	if err != nil {
		t.Fatalf("cannot get terminfo screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		t.Fatalf("failed to init screen: %v", err)
	}

	w, h := screen.Size()
	if w != 80 || h != 24 {
		t.Errorf("bad window size %d x %d", w, h)
	}
	time.Sleep(time.Millisecond * 100)
	screen.Fini()
}

func TestMockDECALN(t *testing.T) {
	mt := &MockTty{Rows: 3, Cols: 5, Fg: tcell.ColorWhite, Bg: tcell.ColorBlack}
	mt.Reset()
	if err := mt.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	mt.Write([]byte("\x1b#8"))
	mt.Drain()
	mt.Stop()
	if len(mt.Cells) != 15 {
		t.Fatalf("Wrong cell length: %d", len(mt.Cells))
	}
	for i := range mt.Cells {
		if string(mt.Cells[i].C) != "E" {
			t.Errorf("wrong value at %d: %s", i, string(mt.Cells[i].C))
			break
		}
		if mt.Cells[i].Attr != tcell.AttrNone {
			t.Errorf("wrong attr at %d: %x", i, mt.Cells[i].Attr)
		}
		if mt.Cells[i].Fg != tcell.ColorWhite {
			t.Errorf("wrong fg at %d: %s", i, mt.Cells[i].Fg.String())
		}
		if mt.Cells[i].Bg != tcell.ColorBlack {
			t.Errorf("wrong bg at %d: %s", i, mt.Cells[i].Bg.String())
		}
		if mt.Cells[i].Width != 1 {
			t.Errorf("wrong width at %d: %d", i, mt.Cells[i].Width)
		}

	}
	if err := mt.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	mt.Fg = tcell.ColorRed
	mt.Attr = tcell.AttrBold
	mt.Write([]byte("\x1b#8"))
	mt.Drain()
	for i := range mt.Cells {
		if string(mt.Cells[i].C) != "E" {
			t.Errorf("wrong value at %d: %s", i, string(mt.Cells[i].C))
			break
		}
		if mt.Cells[i].Attr != tcell.AttrBold {
			t.Errorf("wrong attr at %d: %x", i, mt.Cells[i].Attr)
		}
		if mt.Cells[i].Fg != tcell.ColorRed {
			t.Errorf("wrong fg at %d: %s", i, mt.Cells[i].Fg.String())
		}
		if mt.Cells[i].Bg != tcell.ColorBlack {
			t.Errorf("wrong bg at %d: %s", i, mt.Cells[i].Bg.String())
		}
	}
	mt.Stop()

	mt.Close()
}

func TestMockCursorMovement(t *testing.T) {
	mt := &MockTty{Rows: 3, Cols: 5, Fg: tcell.ColorWhite, Bg: tcell.ColorBlack}
	mt.Reset()
	if err := mt.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	checkPos := func(x vt.Col, y vt.Row) {
		mt.Drain()
		if mt.X != x || mt.Y != y {
			t.Errorf("bad position %d,%d (expected %d,%d)", mt.X, mt.Y, x, y)
		}
	}
	mt.Write([]byte("\x1b[2;3H"))
	checkPos(2, 1)

	mt.Write([]byte("\x1b[20A"))
	checkPos(2, 0)

	mt.Write([]byte("\x1b[20B"))
	checkPos(2, 2)

	mt.Write([]byte("\x1b[A"))
	checkPos(2, 1)

	mt.Write([]byte("\x1b[2C"))
	checkPos(4, 1)

	mt.Write([]byte("\x1b[3D"))
	checkPos(1, 1)

	mt.Write([]byte("\x1b[100D"))
	checkPos(0, 1)

	// Now try the next line and previous line
	mt.Write([]byte("\x1b[2;3H"))
	checkPos(2, 1)
	mt.Write([]byte("\x1b[1E"))
	checkPos(0, 2)

	mt.Write([]byte("\x1b[2;3H"))
	checkPos(2, 1)
	mt.Write([]byte("\x1b[1F"))
	checkPos(0, 0)

	mt.Stop()
	mt.Close()
}
