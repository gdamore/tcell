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

package vt

import (
	"testing"

	"github.com/gdamore/tcell/v3/color"
)

func TestCursorStyleHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		style    CursorStyle
		visible  bool
		blinking bool
	}{
		{name: "steady block", style: SteadyBlock, visible: true, blinking: false},
		{name: "steady bar", style: SteadyBar, visible: true, blinking: false},
		{name: "steady underline", style: SteadyUnderline, visible: true, blinking: false},
		{name: "blinking block", style: BlinkingBlock, visible: true, blinking: true},
		{name: "hidden blink", style: BlinkingBar.Hide(), visible: false, blinking: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.style.IsVisible(); got != tc.visible {
				t.Fatalf("IsVisible() = %v, want %v", got, tc.visible)
			}
			if got := tc.style.IsBlinking(); got != tc.blinking {
				t.Fatalf("IsBlinking() = %v, want %v", got, tc.blinking)
			}

			if got := tc.style.Hide().IsVisible(); got {
				t.Fatalf("Hide() should clear visibility")
			}
			if got := tc.style.Show().IsVisible(); !got {
				t.Fatalf("Show() should set visibility")
			}
			if got := tc.style.Blink().IsBlinking(); !got {
				t.Fatalf("Blink() should set blinking")
			}
			if got := tc.style.Steady().IsBlinking(); got {
				t.Fatalf("Steady() should clear blinking")
			}
		})
	}
}

func TestModeFormatters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		pm          PrivateMode
		enable      string
		disable     string
		query       string
		replyOn     string
		replyOff    string
		ansi        AnsiMode
		ansiEnable  string
		ansiDisable string
		ansiQuery   string
		ansiReply   string
	}{
		{
			name:        "auto margin",
			pm:          PmAutoMargin,
			enable:      "\x1b[?7h",
			disable:     "\x1b[?7l",
			query:       "\x1b[?7$p",
			replyOn:     "\x1b[?7;1$y",
			replyOff:    "\x1b[?7;2$y",
			ansi:        AmNewLineMode,
			ansiEnable:  "\x1b[20h",
			ansiDisable: "\x1b[20l",
			ansiQuery:   "\x1b[20$p",
			ansiReply:   "\x1b[20;1$y",
		},
	}

	for _, tc := range tests {
		if got := tc.pm.Enable(); got != tc.enable {
			t.Fatalf("PrivateMode.Enable() = %q, want %q", got, tc.enable)
		}
		if got := tc.pm.Disable(); got != tc.disable {
			t.Fatalf("PrivateMode.Disable() = %q, want %q", got, tc.disable)
		}
		if got := tc.pm.Query(); got != tc.query {
			t.Fatalf("PrivateMode.Query() = %q, want %q", got, tc.query)
		}
		if got := tc.pm.Reply(ModeOn); got != tc.replyOn {
			t.Fatalf("PrivateMode.Reply(ModeOn) = %q, want %q", got, tc.replyOn)
		}
		if got := tc.pm.Reply(ModeOff); got != tc.replyOff {
			t.Fatalf("PrivateMode.Reply(ModeOff) = %q, want %q", got, tc.replyOff)
		}

		if got := tc.ansi.Enable(); got != tc.ansiEnable {
			t.Fatalf("AnsiMode.Enable() = %q, want %q", got, tc.ansiEnable)
		}
		if got := tc.ansi.Disable(); got != tc.ansiDisable {
			t.Fatalf("AnsiMode.Disable() = %q, want %q", got, tc.ansiDisable)
		}
		if got := tc.ansi.Query(); got != tc.ansiQuery {
			t.Fatalf("AnsiMode.Query() = %q, want %q", got, tc.ansiQuery)
		}
		if got := tc.ansi.Reply(ModeOn); got != tc.ansiReply {
			t.Fatalf("AnsiMode.Reply(ModeOn) = %q, want %q", got, tc.ansiReply)
		}
	}
}

func TestModeStatusChangeable(t *testing.T) {
	t.Parallel()

	cases := []struct {
		status ModeStatus
		want   bool
	}{
		{ModeNA, false},
		{ModeOn, true},
		{ModeOff, true},
		{ModeOnLocked, false},
		{ModeOffLocked, false},
	}

	for _, tc := range cases {
		if got := tc.status.Changeable(); got != tc.want {
			t.Fatalf("Changeable(%v) = %v, want %v", tc.status, got, tc.want)
		}
	}
}

func TestMockBackendResizeAndBells(t *testing.T) {
	t.Parallel()

	mb := NewMockBackend(MockOptSize{X: 2, Y: 2}, MockOptColors(0)).(*mockBackend)

	if got := mb.Bells(); got != 0 {
		t.Fatalf("initial bells = %d, want 0", got)
	}
	mb.Beep()
	mb.Beep()
	if got := mb.Bells(); got != 2 {
		t.Fatalf("bells after beep = %d, want 2", got)
	}

	mb.Put(Coord{X: 0, Y: 0}, Cell{C: "x", S: BaseStyle, W: 1})
	mb.SetPosition(Coord{X: 1, Y: 1})
	mb.SetSize(Coord{X: 3, Y: 1})
	if got := mb.GetSize(); got != (Coord{X: 3, Y: 1}) {
		t.Fatalf("GetSize() = %v, want {3 1}", got)
	}

	if got := mb.GetCell(Coord{X: 0, Y: 0}); got.C != "x" || got.W != 1 {
		t.Fatalf("resized cell preserved incorrectly: %#v", got)
	}
	if got := mb.GetCell(Coord{X: 2, Y: 0}); got.S != BaseStyle {
		t.Fatalf("new cell style = %#v, want BaseStyle", got.S)
	}
}

func TestStringCacheHelpers(t *testing.T) {
	t.Parallel()

	mb := NewMockBackend(MockOptSize{X: 4, Y: 1}, MockOptColors(0)).(*mockBackend)
	em := NewEmulator(mb).(*emulator)

	if got := em.runeString('a'); got != "a" {
		t.Fatalf("runeString(ascii) = %q, want %q", got, "a")
	}
	if got := em.runeString('\u03c0'); got != "π" {
		t.Fatalf("runeString(non-ascii) = %q, want %q", got, "π")
	}
	if got := em.runeString('\u03c0'); got != "π" {
		t.Fatalf("runeString cache hit returned %q, want %q", got, "π")
	}

	if got := em.clusterString([]byte("e\u0301")); got != "e\u0301" {
		t.Fatalf("clusterString(combining) = %q, want %q", got, "e\u0301")
	}
	if got := em.clusterString([]byte("e\u0301")); got != "e\u0301" {
		t.Fatalf("clusterString cache hit returned %q, want %q", got, "e\u0301")
	}
}

func TestEmulatorBellDispatch(t *testing.T) {
	t.Parallel()

	mb := NewMockBackend(MockOptSize{X: 4, Y: 1}, MockOptColors(0)).(*mockBackend)
	em := NewEmulator(mb).(*emulator)
	em.style = BaseStyle.WithFg(color.White).WithBg(color.Black)
	em.defaultStyle = em.style

	em.beep()
	if got := mb.Bells(); got != 1 {
		t.Fatalf("bell count = %d, want 1", got)
	}
}

func TestShouldCheckGrapheme(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		prev byte
		r    rune
		want bool
	}{
		{name: "crlf", prev: '\r', r: '\n', want: true},
		{name: "ascii false", prev: 'a', r: 'b', want: false},
		{name: "mark", prev: 'e', r: '\u0301', want: true},
		{name: "zwj", prev: 'x', r: '\u200d', want: true},
		{name: "vs16", prev: 'x', r: '\uFE0F', want: true},
		{name: "supplementary vs", prev: 'x', r: '\U000E0100', want: true},
		{name: "plain non-ascii", prev: 'x', r: 'π', want: false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := shouldCheckGrapheme(tc.prev, tc.r); got != tc.want {
				t.Fatalf("shouldCheckGrapheme(%q, %q) = %v, want %v", tc.prev, tc.r, got, tc.want)
			}
		})
	}
}

func TestPutRunePaths(t *testing.T) {
	newEm := func(width int) *emulator {
		em := NewEmulator(NewMockBackend(MockOptSize{X: Col(width), Y: 1}, MockOptColors(0))).(*emulator)
		em.style = BaseStyle.WithFg(color.White).WithBg(color.Black)
		em.defaultStyle = em.style
		em.localModes[PmGraphemeClusters] = ModeOn
		em.localModes[PmAutoMargin] = ModeOn
		em.cells[0].S = em.style
		em.cells[0].W = 1
		return em
	}

	t.Run("ascii falls through", func(t *testing.T) {
		em := newEm(4)
		em.cells[0].C = "a"
		em.lastIndex = 1
		em.setPosition(Coord{X: 1, Y: 0})

		em.putRune('b')

		if got := em.cells[1].C; got != "b" {
			t.Fatalf("cell[1].C = %q, want %q", got, "b")
		}
		if got := em.cells[0].C; got != "a" {
			t.Fatalf("cell[0].C = %q, want %q", got, "a")
		}
	})

	t.Run("combining mark extends cluster", func(t *testing.T) {
		em := newEm(4)
		em.cells[0].C = "e"
		em.lastIndex = 1
		em.setPosition(Coord{X: 1, Y: 0})

		em.putRune('\u0301')

		if got := em.cells[0].C; got != "e\u0301" {
			t.Fatalf("cell[0].C = %q, want %q", got, "e\u0301")
		}
		if got := em.cells[0].W; got != 1 {
			t.Fatalf("cell[0].W = %d, want 1", got)
		}
	})

	t.Run("wide grapheme clears next cell", func(t *testing.T) {
		em := newEm(4)
		em.cells[0].C = "\u2764"
		em.cells[1].C = "x"
		em.cells[1].S = em.style
		em.cells[1].W = 1
		em.lastIndex = 1
		em.setPosition(Coord{X: 1, Y: 0})

		em.putRune('\uFE0F')

		if got := em.cells[0].W; got != 2 {
			t.Fatalf("cell[0].W = %d, want 2", got)
		}
		if got := em.cells[1].C; got != "" || em.cells[1].W != 0 {
			t.Fatalf("cell[1] not cleared: %#v", em.cells[1])
		}
	})

	t.Run("auto wrap when width reaches margin", func(t *testing.T) {
		em := newEm(2)
		em.cells[0].C = "a"
		em.lastIndex = 1
		em.setPosition(Coord{X: 1, Y: 0})

		em.putRune('宽')

		if !em.autoWrap {
			t.Fatalf("autoWrap = false, want true")
		}
	})

	t.Run("grapheme extension at margin preserves wrap", func(t *testing.T) {
		em := newEm(2)
		em.cells[1].C = "\u2764"
		em.cells[1].S = em.style
		em.cells[1].W = 1
		em.lastIndex = 2
		em.setPosition(Coord{X: 1, Y: 0})

		em.putRune('\uFE0F')

		if !em.autoWrap {
			t.Fatalf("autoWrap = false, want true")
		}
		if got := em.getPosition(); got.X != 1 {
			t.Fatalf("cursor X = %d, want 1", got.X)
		}
	})
}
