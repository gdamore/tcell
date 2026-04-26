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
	"bytes"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/tty"
	"github.com/gdamore/tcell/v3/vt"
)

// This just offers some very basic tests that do not require a full mock.

// drainInput just does a very primitive sleep to allow input to drain.
// We need this because otherwise the application will close too soon before
// consuming characters from input, including sequences that are returned in
// response to queries.
func drainInput() {
	time.Sleep(time.Millisecond * 30)
}

// TestInitScreen just tries to initialize the default screen.
// It requires a working tty.
func TestInitScreen(t *testing.T) {
	s, err := NewTerminfoScreen()
	if err != nil {
		t.Skip("failed to get screen", err)
	}
	if err := s.Init(); err != nil {
		t.Skip("failed to initialize screen", err)
	}
	defer s.Fini()

	if s.CharacterSet() != "UTF-8" {
		t.Fatalf("Character Set (%v) not UTF-8", s.CharacterSet())
	}

	drainInput()
}

type spyTty struct {
	vt.MockTerm
	writes bytes.Buffer
}

func (t *spyTty) Write(b []byte) (int, error) {
	_, _ = t.writes.Write(b)
	return t.MockTerm.Write(b)
}

func (t *spyTty) Output() string {
	return t.writes.String()
}

func TestOptAltScreenDisable(t *testing.T) {
	t.Setenv("TCELL_ALTSCREEN", "")

	tty := &spyTty{MockTerm: vt.NewMockTerm()}
	s, err := NewTerminfoScreenFromTty(tty, OptAltScreen(false))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	s.Fini()

	out := tty.Output()
	if strings.Contains(out, enterCA) {
		t.Fatalf("alternate screen enter escape was emitted")
	}
	if strings.Contains(out, exitCA) {
		t.Fatalf("alternate screen exit escape was emitted")
	}
}

func TestOptAltScreenDefault(t *testing.T) {
	t.Setenv("TCELL_ALTSCREEN", "")

	tty := &spyTty{MockTerm: vt.NewMockTerm()}
	s, err := NewTerminfoScreenFromTty(tty)
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	s.Fini()

	out := tty.Output()
	if !strings.Contains(out, enterCA) {
		t.Fatalf("alternate screen enter escape was not emitted")
	}
	if !strings.Contains(out, exitCA) {
		t.Fatalf("alternate screen exit escape was not emitted")
	}
}

func TestOSC8ControlsAreStrippedFromOutput(t *testing.T) {
	tty := &spyTty{MockTerm: vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})}
	s, err := NewTerminfoScreenFromTty(tty, OptAltScreen(false))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer s.Fini()

	style := StyleDefault.
		Url("http://exa\x07mple.com/\x1b\\path").
		UrlId("id\x00\x1f\x7f\x80\x9fend")

	s.PutStrStyled(0, 0, "X", style)
	s.Show()

	out := tty.Output()
	const prefix = "\x1b]8;id=idend;"
	_, link, ok := strings.Cut(out, prefix)
	if !ok {
		t.Fatalf("missing OSC 8 link open sequence in output: %q", out)
	}
	link, _, ok = strings.Cut(link, "\x1b\\")
	if !ok {
		t.Fatalf("missing OSC 8 terminator in output: %q", out)
	}
	if link != "http://example.com/\\path" {
		t.Fatalf("unexpected emitted URL payload: %q", link)
	}
	for i := 0; i < len(link); i++ {
		c := link[i]
		if c <= 0x1f || c == 0x7f || (c >= 0x80 && c <= 0x9f) {
			t.Fatalf("control characters survived in emitted URL payload: %q", link)
		}
	}
}

func TestKeyboardProtocol(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*tScreen)
		want  KeyProtocol
	}{
		{
			name: "Legacy",
			want: LegacyKeyboard,
		},
		{
			name: "Xterm",
			setup: func(s *tScreen) {
				s.haveXTermKbd = true
			},
			want: XTermKeyboard,
		},
		{
			name: "Kitty",
			setup: func(s *tScreen) {
				s.haveKittyKbd = true
			},
			want: KittyKeyboard,
		},
		{
			name: "Win32",
			setup: func(s *tScreen) {
				s.haveWin32Kbd = true
			},
			want: Win32Keyboard,
		},
		{
			name: "KittyBeforeXterm",
			setup: func(s *tScreen) {
				s.haveKittyKbd = true
				s.haveXTermKbd = true
			},
			want: KittyKeyboard,
		},
		{
			name: "Win32BeforeKittyAndXterm",
			setup: func(s *tScreen) {
				s.haveWin32Kbd = true
				s.haveKittyKbd = true
				s.haveXTermKbd = true
			},
			want: Win32Keyboard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &tScreen{}
			if tt.setup != nil {
				tt.setup(s)
			}
			if got := s.KeyboardProtocol(); got != tt.want {
				t.Fatalf("KeyboardProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessInitQKeyboardProtocol(t *testing.T) {
	tests := []struct {
		name string
		evs  []Event
		want KeyProtocol
	}{
		{
			name: "Legacy",
			want: LegacyKeyboard,
		},
		{
			name: "Xterm",
			evs:  []Event{&eventXTermKbdMode{Mode: XtermKbdModeExt}},
			want: XTermKeyboard,
		},
		{
			name: "Kitty",
			evs:  []Event{&eventKittyKbdMode{Mode: KittyKbdModeBase}},
			want: KittyKeyboard,
		},
		{
			name: "Win32",
			evs:  []Event{&eventPrivateMode{Mode: vt.PmWin32Input, Status: vt.ModeOff}},
			want: Win32Keyboard,
		},
		{
			name: "KittyPreferredOverXterm",
			evs: []Event{
				&eventXTermKbdMode{Mode: XtermKbdModeExt},
				&eventKittyKbdMode{Mode: KittyKbdModeBase},
			},
			want: KittyKeyboard,
		},
		{
			name: "Win32PreferredOverKittyAndXterm",
			evs: []Event{
				&eventXTermKbdMode{Mode: XtermKbdModeExt},
				&eventKittyKbdMode{Mode: KittyKbdModeBase},
				&eventPrivateMode{Mode: vt.PmWin32Input, Status: vt.ModeOff},
			},
			want: Win32Keyboard,
		},
		{
			name: "UnchangeableWin32IsIgnored",
			evs:  []Event{&eventPrivateMode{Mode: vt.PmWin32Input, Status: vt.ModeNA}},
			want: LegacyKeyboard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &tScreen{
				initQ: make(chan Event, len(tt.evs)+1),
			}
			for _, ev := range tt.evs {
				s.initQ <- ev
			}
			s.initQ <- &eventPrimaryAttributes{}
			s.Lock()
			s.processInitQ()
			s.Unlock()

			if got := s.KeyboardProtocol(); got != tt.want {
				t.Fatalf("KeyboardProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInitScreenStdio just tries to initialize the default screen using standard I/O.
// It requires a working tty.
func TestInitScreenStdio(t *testing.T) {
	tty, err := tty.NewStdIoTty()
	if err != nil {
		t.Skip("maybe stdin is not a tty?")
		return
	}
	s, err := NewTerminfoScreenFromTty(tty)
	if err := s.Init(); err != nil {
		t.Skip("failed to initialize screen", err)
		tty.Close()
		return
	}
	defer s.Fini()

	if s.CharacterSet() != "UTF-8" {
		t.Fatalf("Character Set (%v) not UTF-8", s.CharacterSet())
	}
	drainInput()
}

func TestNotDevNull(t *testing.T) {
	tty, err := tty.NewDevTtyFromDev("/dev/null")
	if err == nil {
		tty.Close()
		t.Error("open /dev/null as tty should not have passed")
	}
}

func TestNoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	s, err := NewTerminfoScreen()
	if err != nil {
		t.Skip("failed to get screen")
	}
	if err := s.Init(); err != nil {
		t.Skip("failed to initialize screen", err)
	}
	defer s.Fini()

	if s.Colors() != 0 {
		t.Errorf("screen should not have color but had %d", s.Colors())
	}

	drainInput()
}

func NewMockScreen(t *testing.T, opts ...vt.MockOpt) (vt.MockTerm, Screen) {
	t.Helper()
	if runtime.GOOS == "js" {
		t.Skip("not supported on webasm")
		return nil, nil
	}
	mt := vt.NewMockTerm(opts...)
	scr, err := NewTerminfoScreenFromTty(mt)
	if err != nil {
		t.Fatalf("failed to get terminal: %v", err)
		return nil, nil
	}
	if err = scr.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
		return nil, nil
	}
	return mt, scr
}
