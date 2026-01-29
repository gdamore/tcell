// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"testing"

	"github.com/gdamore/tcell/v3/vt"
)

func TestMouse1006(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(0))
	defer MustClose(t, term)

	MustStart(t, term)

	WriteF(t, term, "\x1b[?1006$p")
	result := ReadF(t, term)
	want := "\x1b[?1006;2$y"
	VerifyF(t, result == want, "wrong mode report: %q != %q", result, want)

	// turn on mouse reports
	WriteF(t, term, "\x1b[?1000h\x1b[?1002h\x1b[?1003h\x1b[?1006h")

	// simple mouse click
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button1,
		Down:     true,
		Motion:   false,
		Mod:      vt.ModNone,
	})
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button1,
		Down:     false,
		Motion:   false,
		Mod:      vt.ModNone,
	})

	result = ReadF(t, term)
	want = "\x1b[<0;3;4M\x1b[<0;3;4m"
	VerifyF(t, result == want, "wrong mouse event: %q != %q", result, want)

	// modified drag (10x10 square)
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button2,
		Down:     true,
		Motion:   false,
		Mod:      vt.ModLShift,
	})
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 12, Y: 13},
		Button:   vt.Button2,
		Down:     false,
		Motion:   true,
		Mod:      vt.ModRShift,
	})

	result = ReadF(t, term)
	want = "\x1b[<6;3;4M\x1b[<38;13;14m"
	VerifyF(t, result == want, "wrong mouse event: %q != %q", result, want)
}

func TestMouseX10(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(0))
	defer MustClose(t, term)

	MustStart(t, term)

	WriteF(t, term, "\x1b[?9$p")
	result := ReadF(t, term)
	want := "\x1b[?9;2$y"
	VerifyF(t, result == want, "wrong mode report: %q != %q", result, want)

	// turn on mouse reports
	WriteF(t, term, "\x1b[?9h")

	// simple mouse click
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button1,
		Down:     true,
		Motion:   false,
		Mod:      vt.ModNone,
	})
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button1,
		Down:     false,
		Motion:   false,
		Mod:      vt.ModNone,
	})

	result = ReadF(t, term)
	want = "\x1b[M #$"
	VerifyF(t, result == want, "wrong mouse event: %q != %q", result, want)

	// button 3 mouse click
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button3,
		Down:     true,
		Motion:   false,
		Mod:      vt.ModNone,
	})
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button3,
		Down:     false,
		Motion:   false,
		Mod:      vt.ModNone,
	})

	result = ReadF(t, term)
	want = "\x1b[M!#$"
	VerifyF(t, result == want, "wrong mouse event: %q != %q", result, want)

	// modified drag (10x10 square) - we only see the start
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button2,
		Down:     true,
		Motion:   false,
		Mod:      vt.ModLShift,
	})
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 12, Y: 13},
		Button:   vt.Button2,
		Down:     false,
		Motion:   true,
		Mod:      vt.ModLShift,
	})

	result = ReadF(t, term)
	want = "\x1b[M\"#$"
	VerifyF(t, result == want, "wrong mouse event: %q != %q", result, want)
}

// TestMouse1000 tests basic mouse mode (VT200, no tracking)
func TestMouse1000(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(0))
	defer MustClose(t, term)

	MustStart(t, term)

	WriteF(t, term, "\x1b[?1000$p")
	result := ReadF(t, term)
	want := "\x1b[?1000;2$y"
	VerifyF(t, result == want, "wrong mode report: %q != %q", result, want)

	// turn on mouse reports
	WriteF(t, term, "\x1b[?1000h")

	// simple mouse click
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button1,
		Down:     true,
		Motion:   false,
		Mod:      vt.ModNone,
	})
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button1,
		Down:     false,
		Motion:   false,
		Mod:      vt.ModNone,
	})

	result = ReadF(t, term)
	want = "\x1b[M #$\x1b[M##$"
	VerifyF(t, result == want, "wrong mouse event: %q != %q", result, want)

	// modified drag (10x10 square) - we only see the start
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 2, Y: 3},
		Button:   vt.Button2,
		Down:     true,
		Motion:   false,
		Mod:      vt.ModLShift,
	})
	term.MouseEvent(vt.MouseEvent{ // should be suppressed
		Position: Coord{X: 11, Y: 12},
		Button:   vt.NoButton,
		Down:     false,
		Motion:   true,
		Mod:      vt.ModLShift,
	})
	term.MouseEvent(vt.MouseEvent{
		Position: Coord{X: 12, Y: 13},
		Button:   vt.Button2,
		Down:     false,
		Motion:   true,
		Mod:      vt.ModRShift,
	})

	result = ReadF(t, term)
	want = "\x1b[M&#$\x1b[M'-."
	VerifyF(t, result == want, "wrong mouse event: %q != %q", result, want)
}
