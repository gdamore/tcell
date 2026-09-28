// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
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
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/color"
)

// waitEvent returns the next event of type T, skipping others.
func waitEvent[T Event](t *testing.T, s Screen) T {
	t.Helper()
	for {
		select {
		case ev := <-s.EventQ():
			if e, ok := ev.(T); ok {
				return e
			}
		case <-time.After(time.Second):
			var zero T
			t.Fatalf("timeout waiting for %T", zero)
			return zero
		}
	}
}

// noEvent fails if an event of type T arrives within a short time.
func noEvent[T Event](t *testing.T, s Screen) {
	t.Helper()
	deadline := time.After(100 * time.Millisecond)
	for {
		select {
		case ev := <-s.EventQ():
			if _, ok := ev.(T); ok {
				t.Fatalf("unexpected %T", ev)
			}
		case <-deadline:
			return
		}
	}
}

func TestColorSchemeEvents(t *testing.T) {
	term, s := NewMockScreen(t)
	defer s.Fini()

	// No reports until they are enabled.
	term.SetColorScheme(false, color.White)
	noEvent[*EventColorScheme](t, s)

	// Enabling reports the current scheme.
	s.EnableColorScheme()
	if ev := waitEvent[*EventColorScheme](t, s); ev.Dark {
		t.Error("current scheme: dark, want light")
	}

	// Then each change.
	term.SetColorScheme(true, color.Black)
	if ev := waitEvent[*EventColorScheme](t, s); !ev.Dark {
		t.Error("after switching to dark: light")
	}
	term.SetColorScheme(false, color.White)
	if ev := waitEvent[*EventColorScheme](t, s); ev.Dark {
		t.Error("after switching to light: dark")
	}

	s.DisableColorScheme()
	time.Sleep(10 * time.Millisecond)
	term.SetColorScheme(true, color.Black)
	noEvent[*EventColorScheme](t, s)
}

func TestColorSchemeAfterResume(t *testing.T) {
	term, s := NewMockScreen(t)
	defer s.Fini()

	s.EnableColorScheme()
	waitEvent[*EventColorScheme](t, s)
	if err := s.Suspend(); err != nil {
		t.Fatal(err)
	}
	if err := s.Resume(); err != nil {
		t.Fatal(err)
	}
	// Resuming turns reports on again and asks for the scheme again.
	waitEvent[*EventColorScheme](t, s)
	term.SetColorScheme(false, color.White)
	if ev := waitEvent[*EventColorScheme](t, s); ev.Dark {
		t.Error("after resume and switching to light: dark")
	}
}

func TestBackgroundColor(t *testing.T) {
	term, s := NewMockScreen(t)
	defer s.Fini()

	term.SetColorScheme(true, color.NewRGBColor(0x1e, 0x1e, 0x2e))
	s.GetBackgroundColor()
	ev := waitEvent[*EventBackgroundColor](t, s)
	if r, g, b := ev.Color.RGB(); r != 0x1e || g != 0x1e || b != 0x2e {
		t.Errorf("background = %02x%02x%02x, want 1e1e2e", r, g, b)
	}
	if !ev.Dark() {
		t.Error("1e1e2e should be dark")
	}

	term.SetColorScheme(false, color.NewRGBColor(0xfa, 0xfa, 0xfa))
	s.GetBackgroundColor()
	if ev := waitEvent[*EventBackgroundColor](t, s); ev.Dark() {
		t.Error("fafafa should be light")
	}
}

// Replies as real terminals send them, including ones tcell ignores.
func TestColorSchemeReplies(t *testing.T) {
	term, s := NewMockScreen(t)
	defer s.Fini()

	term.SendRaw([]byte("\x1b[?997;1n"))
	if ev := waitEvent[*EventColorScheme](t, s); !ev.Dark {
		t.Error("997;1 is dark")
	}
	term.SendRaw([]byte("\x1b[?997;2n"))
	if ev := waitEvent[*EventColorScheme](t, s); ev.Dark {
		t.Error("997;2 is light")
	}

	// BEL-terminated, and with two digits per component.
	term.SendRaw([]byte("\x1b]11;rgb:ff/80/00\x07"))
	ev := waitEvent[*EventBackgroundColor](t, s)
	if r, g, b := ev.Color.RGB(); r != 0xff || g != 0x80 || b != 0 {
		t.Errorf("rgb:ff/80/00 = %02x%02x%02x", r, g, b)
	}

	// Malformed replies post nothing, and don't disturb what follows.
	term.SendRaw([]byte("\x1b[?997;7n\x1b]11;rgb:zz/00/00\x1b\\\x1b[?997;1n"))
	if ev := waitEvent[*EventColorScheme](t, s); !ev.Dark {
		t.Error("the report after malformed ones was lost")
	}
}

func TestParseXColor(t *testing.T) {
	cases := []struct {
		in      string
		r, g, b int32
		ok      bool
	}{
		{"rgb:0000/0000/0000", 0, 0, 0, true},
		{"rgb:ffff/ffff/ffff", 255, 255, 255, true},
		{"rgb:1e1e/1e1e/2e2e", 0x1e, 0x1e, 0x2e, true},
		{"rgb:f/8/0", 255, 0x88, 0, true},
		{"rgb:fff/800/000", 255, 0x7f, 0, true},
		{"rgba:ffff/0000/0000/ffff", 255, 0, 0, true},
		{"rgb:ff/ff", 0, 0, 0, false},
		{"rgb:fffff/0/0", 0, 0, 0, false},
		{"rgb://", 0, 0, 0, false},
		{"#ffffff", 0, 0, 0, false},
		{"", 0, 0, 0, false},
	}
	for _, c := range cases {
		col, ok := parseXColor(c.in)
		if ok != c.ok {
			t.Errorf("%q: ok = %v", c.in, ok)
			continue
		}
		if !ok {
			continue
		}
		if r, g, b := col.RGB(); r != c.r || g != c.g || b != c.b {
			t.Errorf("%q = %02x%02x%02x, want %02x%02x%02x", c.in, r, g, b, c.r, c.g, c.b)
		}
	}
}
