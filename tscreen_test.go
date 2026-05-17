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

func TestFiniPreventsSetStyleMutation(t *testing.T) {
	mt := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 2})
	scr, err := NewTerminfoScreenFromTty(mt)
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	bs, ok := scr.(*baseScreen)
	if !ok {
		t.Fatalf("expected *baseScreen, got %T", scr)
	}
	ts, ok := bs.screenImpl.(*tScreen)
	if !ok {
		t.Fatalf("expected *tScreen, got %T", bs.screenImpl)
	}
	if err := scr.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}

	before := StyleDefault.Foreground(ColorRed)
	after := StyleDefault.Foreground(ColorBlue)
	scr.SetStyle(before)
	scr.Fini()
	scr.SetStyle(after)

	ts.Lock()
	defer ts.Unlock()
	if !ts.fini {
		t.Fatal("screen was not marked finished")
	}
	if ts.style != before {
		t.Fatal("SetStyle mutated the screen after Fini")
	}
}

func TestFiniPreventsFurtherTerminalMutation(t *testing.T) {
	tty := &spyTty{MockTerm: vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 2})}
	scr, err := NewTerminfoScreenFromTty(tty)
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	bs, ok := scr.(*baseScreen)
	if !ok {
		t.Fatalf("expected *baseScreen, got %T", scr)
	}
	ts, ok := bs.screenImpl.(*tScreen)
	if !ok {
		t.Fatalf("expected *tScreen, got %T", bs.screenImpl)
	}
	if err := scr.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}

	scr.ShowCursor(1, 1)
	scr.SetCursorStyle(CursorStyleSteadyBlock, ColorRed)
	scr.EnableMouse()
	scr.EnablePaste()
	scr.EnableFocus()
	scr.Fini()

	beforeOutput := tty.Output()
	ts.Lock()
	beforeCursorX, beforeCursorY := ts.cursorx, ts.cursory
	beforeCursorStyle, beforeCursorColor := ts.cursorStyle, ts.cursorColor
	beforeMouseFlags := ts.mouseFlags
	beforePasteEnabled, beforeFocusEnabled := ts.pasteEnabled, ts.focusEnabled
	ts.Unlock()

	scr.ShowCursor(2, 1)
	scr.SetCursorStyle(CursorStyleBlinkingBar, ColorBlue)
	scr.EnableMouse(MouseMotionEvents)
	scr.DisableMouse()
	scr.EnablePaste()
	scr.DisablePaste()
	scr.EnableFocus()
	scr.DisableFocus()
	scr.SetSize(20, 10)
	if err := scr.Suspend(); err != nil {
		t.Fatalf("Suspend after Fini failed: %v", err)
	}
	if err := scr.Resume(); err != nil {
		t.Fatalf("Resume after Fini failed: %v", err)
	}
	if err := scr.Beep(); err != nil {
		t.Fatalf("Beep after Fini failed: %v", err)
	}
	scr.SetTitle("after")
	scr.SetClipboard([]byte("after"))
	scr.GetClipboard()
	scr.ShowNotification("after", "after")

	if got := tty.Output(); got != beforeOutput {
		t.Fatal("terminal output changed after Fini")
	}

	ts.Lock()
	defer ts.Unlock()
	if ts.cursorx != beforeCursorX || ts.cursory != beforeCursorY {
		t.Fatal("cursor position changed after Fini")
	}
	if ts.cursorStyle != beforeCursorStyle || ts.cursorColor != beforeCursorColor {
		t.Fatal("cursor style changed after Fini")
	}
	if ts.mouseFlags != beforeMouseFlags {
		t.Fatal("mouse flags changed after Fini")
	}
	if ts.pasteEnabled != beforePasteEnabled {
		t.Fatal("paste state changed after Fini")
	}
	if ts.focusEnabled != beforeFocusEnabled {
		t.Fatal("focus state changed after Fini")
	}
}

func TestOptAdvancedKeys(t *testing.T) {
	mt := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 2})
	scr, err := NewTerminfoScreenFromTty(mt, OptAdvancedKeys(true))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	bs, ok := scr.(*baseScreen)
	if !ok {
		t.Fatalf("expected *baseScreen, got %T", scr)
	}
	ts, ok := bs.screenImpl.(*tScreen)
	if !ok {
		t.Fatalf("expected *tScreen, got %T", bs.screenImpl)
	}
	if !ts.advancedKeys {
		t.Fatal("advanced keys option was not applied")
	}
	if err := scr.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer scr.Fini()
	if !ts.input.advanced {
		t.Fatal("advanced keys option was not propagated to input parser")
	}
}

func TestOptKeyboardProtocol(t *testing.T) {
	mt := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 2})
	scr, err := NewTerminfoScreenFromTty(mt, OptKeyboardProtocol(KittyKeyboard))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	bs, ok := scr.(*baseScreen)
	if !ok {
		t.Fatalf("expected *baseScreen, got %T", scr)
	}
	ts, ok := bs.screenImpl.(*tScreen)
	if !ok {
		t.Fatalf("expected *tScreen, got %T", bs.screenImpl)
	}
	if !ts.forceKbd || ts.forcedKbd != KittyKeyboard {
		t.Fatalf("forced keyboard protocol = (%v, %v), want (true, %v)", ts.forceKbd, ts.forcedKbd, KittyKeyboard)
	}
}

func TestOptNegotiation(t *testing.T) {
	mt := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 2})
	scr, err := NewTerminfoScreenFromTty(mt, OptNegotiation(false))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	bs, ok := scr.(*baseScreen)
	if !ok {
		t.Fatalf("expected *baseScreen, got %T", scr)
	}
	ts, ok := bs.screenImpl.(*tScreen)
	if !ok {
		t.Fatalf("expected *tScreen, got %T", bs.screenImpl)
	}
	if ts.negotiate {
		t.Fatal("negotiation option was not applied")
	}
}

func TestSetSizeTakesScreenLock(t *testing.T) {
	mt := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 2})
	ts := &tScreen{tty: mt, w: 8, h: 2}

	done := make(chan struct{})
	ts.Lock()
	go func() {
		ts.SetSize(8, 2)
		close(done)
	}()

	select {
	case <-done:
		ts.Unlock()
		t.Fatal("SetSize returned while the screen lock was held")
	case <-time.After(10 * time.Millisecond):
	}

	ts.Unlock()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("SetSize did not return after the screen lock was released")
	}
}

func TestOptControlStringLimit(t *testing.T) {
	mt := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 2})
	scr, err := NewTerminfoScreenFromTty(mt, OptControlStringLimit(4096))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	bs, ok := scr.(*baseScreen)
	if !ok {
		t.Fatalf("expected *baseScreen, got %T", scr)
	}
	ts, ok := bs.screenImpl.(*tScreen)
	if !ok {
		t.Fatalf("expected *tScreen, got %T", bs.screenImpl)
	}
	if ts.controlStringLimit != 4096 {
		t.Fatalf("control string limit = %d, want %d", ts.controlStringLimit, 4096)
	}
	if err := scr.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer scr.Fini()
	if ts.input.controlStringMax != 4096 {
		t.Fatalf("input parser control string limit = %d, want %d", ts.input.controlStringMax, 4096)
	}
}

func TestOptControlStringLimitUnlimited(t *testing.T) {
	mt := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 2})
	scr, err := NewTerminfoScreenFromTty(mt, OptControlStringLimit(0))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	bs, ok := scr.(*baseScreen)
	if !ok {
		t.Fatalf("expected *baseScreen, got %T", scr)
	}
	ts, ok := bs.screenImpl.(*tScreen)
	if !ok {
		t.Fatalf("expected *tScreen, got %T", bs.screenImpl)
	}
	if err := scr.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer scr.Fini()

	payload := bytes.Repeat([]byte{'x'}, defaultControlStringLimit+1)
	ts.input.ScanUTF8(append([]byte("\x1b]"), payload...))

	if ts.input.state != istOsc {
		t.Fatalf("parser state = %v, want %v", ts.input.state, istOsc)
	}
	if !bytes.Equal(ts.input.strBuf, payload) {
		t.Fatalf("string buffer length = %d, want %d", len(ts.input.strBuf), len(payload))
	}
}

func TestOptSanitizeContent(t *testing.T) {
	t.Run("disabled by default", func(t *testing.T) {
		mt := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 2})
		scr, err := NewTerminfoScreenFromTty(mt)
		if err != nil {
			t.Fatalf("failed to get screen: %v", err)
		}
		if err := scr.Init(); err != nil {
			t.Fatalf("failed to initialize screen: %v", err)
		}
		defer scr.Fini()

		scr.PutStr(0, 0, "\x1bA\x07B")
		if got, _, _ := scr.Get(0, 0); !strings.Contains(got, "\x1b") {
			t.Fatalf("expected control bytes to remain when sanitizer is disabled, got %q", got)
		}
		if got, _, _ := scr.Get(1, 0); !strings.Contains(got, "\x07") {
			t.Fatalf("expected control bytes to remain when sanitizer is disabled, got %q", got)
		}
	})

	t.Run("enabled", func(t *testing.T) {
		mt := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 2})
		scr, err := NewTerminfoScreenFromTty(mt, OptSanitizeContent(true))
		if err != nil {
			t.Fatalf("failed to get screen: %v", err)
		}
		if err := scr.Init(); err != nil {
			t.Fatalf("failed to initialize screen: %v", err)
		}
		defer scr.Fini()

		scr.PutStr(0, 0, "\x1bA\x07B")
		if got, _, _ := scr.Get(0, 0); got != "A" {
			t.Fatalf("unexpected sanitized cell content at 0,0: %q", got)
		}
		if got, _, _ := scr.Get(1, 0); got != "B" {
			t.Fatalf("unexpected sanitized cell content at 1,0: %q", got)
		}
	})
}

func TestNewScreenSanitizeContentOption(t *testing.T) {
	scr, err := NewScreen(OptSanitizeContent(true))
	if err != nil {
		t.Skipf("failed to get screen: %v", err)
	}
	if err := scr.Init(); err != nil {
		t.Skipf("failed to initialize screen: %v", err)
	}
	defer scr.Fini()

	scr.PutStr(0, 0, "\x1bA\x07B")
	if got, _, _ := scr.Get(0, 0); got != "A" {
		t.Fatalf("unexpected sanitized cell content at 0,0: %q", got)
	}
	if got, _, _ := scr.Get(1, 0); got != "B" {
		t.Fatalf("unexpected sanitized cell content at 1,0: %q", got)
	}
}

func TestNewScreenShimScreen(t *testing.T) {
	_, scr := NewMockScreen(t)
	ShimScreen(scr)

	got, err := NewScreen()
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if got != scr {
		t.Fatalf("unexpected shimmed screen: got %T, want %T", got, scr)
	}
}

func TestUnlockRegionRedrawsUntouchedBlankCell(t *testing.T) {
	mt, scr := NewMockScreen(t, vt.MockOptSize{X: 1, Y: 1})
	defer scr.Fini()

	mt.Backend().Put(vt.Coord{X: 0, Y: 0}, vt.Cell{C: "X", S: vt.BaseStyle, W: 1})
	scr.LockRegion(0, 0, 1, 1, true)
	scr.LockRegion(0, 0, 1, 1, false)
	scr.Show()

	if got := mt.GetCell(vt.Coord{X: 0, Y: 0}); got.C != " " || got.W != 1 {
		t.Fatalf("unlocking an untouched blank cell should repaint a space, got %#v", got)
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
	if _, afterLink, ok := strings.Cut(out, "X"); !ok || !strings.Contains(afterLink, "\x1b]8;;\x1b\\") {
		t.Fatalf("missing OSC 8 link close sequence after linked content: %q", out)
	}
	for i := 0; i < len(link); i++ {
		c := link[i]
		if c <= 0x1f || c == 0x7f || (c >= 0x80 && c <= 0x9f) {
			t.Fatalf("control characters survived in emitted URL payload: %q", link)
		}
	}
}

func TestOSC8IdWithoutUrlDoesNotEmitClose(t *testing.T) {
	tty := &spyTty{MockTerm: vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})}
	s, err := NewTerminfoScreenFromTty(tty, OptAltScreen(false))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer s.Fini()

	tty.writes.Reset()
	s.PutStrStyled(0, 0, "X", StyleDefault.UrlId("orphan"))
	s.Show()

	if out := tty.Output(); strings.Contains(out, "\x1b]8;;\x1b\\") {
		t.Fatalf("unexpected OSC 8 close sequence for id-only style: %q", out)
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

func TestApplyKnownTerminalProfile(t *testing.T) {
	tests := []struct {
		name         string
		goos         string
		termProgram  string
		wantKnown    bool
		wantInitted  bool
		wantMouse    bool
		wantMouseSgr bool
		wantKittyKbd bool
		wantWin32Kbd bool
	}{
		{
			name:         "AppleTerminal",
			goos:         "darwin",
			termProgram:  "Apple_Terminal",
			wantKnown:    true,
			wantMouse:    true,
			wantMouseSgr: true,
		},
		{
			name:         "LocalWindowsWezTerm",
			goos:         "windows",
			termProgram:  "WezTerm",
			wantKnown:    true,
			wantInitted:  true,
			wantMouse:    true,
			wantMouseSgr: true,
			wantWin32Kbd: true,
		},
		{
			name:         "LocalUnixWezTerm",
			goos:         "linux",
			termProgram:  "WezTerm",
			wantKnown:    true,
			wantInitted:  true,
			wantMouse:    true,
			wantMouseSgr: true,
			wantKittyKbd: true,
		},
		{
			name:        "UnknownTerminal",
			goos:        "linux",
			termProgram: "Other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &tScreen{}
			if got := s.applyKnownTerminalProfile(tt.goos, tt.termProgram); got != tt.wantKnown {
				t.Fatalf("applyKnownTerminalProfile() = %v, want %v", got, tt.wantKnown)
			}
			if s.initted != tt.wantInitted {
				t.Fatalf("initted = %v, want %v", s.initted, tt.wantInitted)
			}
			if s.haveMouse != tt.wantMouse {
				t.Fatalf("haveMouse = %v, want %v", s.haveMouse, tt.wantMouse)
			}
			if s.haveMouseSgr != tt.wantMouseSgr {
				t.Fatalf("haveMouseSgr = %v, want %v", s.haveMouseSgr, tt.wantMouseSgr)
			}
			if s.haveKittyKbd != tt.wantKittyKbd {
				t.Fatalf("haveKittyKbd = %v, want %v", s.haveKittyKbd, tt.wantKittyKbd)
			}
			if s.haveWin32Kbd != tt.wantWin32Kbd {
				t.Fatalf("haveWin32Kbd = %v, want %v", s.haveWin32Kbd, tt.wantWin32Kbd)
			}
		})
	}
}

func TestAppleTerminalProfileSkipsStartupQueries(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "Apple_Terminal")
	t.Setenv("TERM_PROGRAM_VERSION", "999")

	tty := &spyTty{MockTerm: vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})}
	s, err := NewTerminfoScreenFromTty(tty, OptAltScreen(false))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer s.Fini()

	out := tty.Output()
	for _, seq := range []string{
		vt.PmResizeReports.Query(),
		vt.PmMouseButton.Query(),
		vt.PmMouseSgr.Query(),
		vt.PmWin32Input.Query(),
		queryKittyKbd,
		queryXTermKbd,
		requestExtAttr,
	} {
		if strings.Contains(out, seq) {
			t.Fatalf("Apple Terminal profile emitted startup query %q", seq)
		}
	}
	if !strings.Contains(out, requestPrimaryDA) {
		t.Fatal("Apple Terminal profile did not emit primary DA")
	}
	name, version := s.Terminal()
	if name != "Terminal.app" || version != "999" {
		t.Fatalf("Terminal() = %q, %q, want %q, %q", name, version, "Terminal.app", "999")
	}
}

func TestKeyboardProbePolicy(t *testing.T) {
	tests := []struct {
		name                string
		goos                string
		wantWindowSizeQuery bool
		wantXTermQuery      bool
	}{
		{
			name:                "Unix",
			goos:                "linux",
			wantWindowSizeQuery: true,
			wantXTermQuery:      true,
		},
		{
			name: "Windows",
			goos: "windows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := useVTWindowSizeQuery(tt.goos); got != tt.wantWindowSizeQuery {
				t.Fatalf("useVTWindowSizeQuery() = %v, want %v", got, tt.wantWindowSizeQuery)
			}
			if got := useXTermKeyboardQuery(tt.goos); got != tt.wantXTermQuery {
				t.Fatalf("useXTermKeyboardQuery() = %v, want %v", got, tt.wantXTermQuery)
			}
		})
	}
}

func TestKeyboardProtocolHelpers(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		want      KeyProtocol
		wantValid bool
	}{
		{name: "legacy", value: "legacy", want: LegacyKeyboard, wantValid: true},
		{name: "kitty", value: "kitty", want: KittyKeyboard, wantValid: true},
		{name: "win32", value: "win32", want: Win32Keyboard, wantValid: true},
		{name: "xterm", value: "xterm", want: XTermKeyboard, wantValid: true},
		{name: "invalid", value: "bogus"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseKeyboardProtocol(tt.value)
			if got != tt.want || ok != tt.wantValid {
				t.Fatalf("parseKeyboardProtocol(%q) = (%v, %v), want (%v, %v)", tt.value, got, ok, tt.want, tt.wantValid)
			}
		})
	}

	if validKeyboardProtocol(KeyProtocol(99)) {
		t.Fatal("invalid keyboard protocol was accepted")
	}
}

func TestForceKeyboardProtocol(t *testing.T) {
	s := &tScreen{
		haveKittyKbd: true,
		haveWin32Kbd: true,
		haveXTermKbd: true,
	}
	if s.forceKeyboardProtocol(KeyProtocol(99)) {
		t.Fatal("invalid keyboard protocol was forced")
	}
	if s.forceKbd {
		t.Fatal("invalid keyboard protocol changed override state")
	}
	if !s.forceKeyboardProtocol(XTermKeyboard) {
		t.Fatal("valid keyboard protocol was not forced")
	}
	s.applyKeyboardProtocolOverride()
	if s.haveKittyKbd || s.haveWin32Kbd || !s.haveXTermKbd {
		t.Fatalf("forced XTerm protocol state = kitty:%v win32:%v xterm:%v", s.haveKittyKbd, s.haveWin32Kbd, s.haveXTermKbd)
	}
}

func TestApplyEnvironmentOverrides(t *testing.T) {
	tests := []struct {
		name             string
		keyboardProtocol string
		negotiate        string
		mouse            string
		startForceKbd    bool
		startForcedKbd   KeyProtocol
		startNegotiate   bool
		wantForceKbd     bool
		wantForcedKbd    KeyProtocol
		wantNegotiate    bool
		wantMouseOff     bool
	}{
		{
			name:           "defaults",
			startNegotiate: true,
			wantNegotiate:  true,
			wantForcedKbd:  LegacyKeyboard,
		},
		{
			name:             "force all",
			keyboardProtocol: "win32",
			negotiate:        "disable",
			mouse:            "disable",
			startNegotiate:   true,
			wantForceKbd:     true,
			wantForcedKbd:    Win32Keyboard,
			wantMouseOff:     true,
		},
		{
			name:             "auto resets options",
			keyboardProtocol: "auto",
			negotiate:        "auto",
			startForceKbd:    true,
			startForcedKbd:   KittyKeyboard,
			wantNegotiate:    true,
			wantForcedKbd:    KittyKeyboard,
		},
		{
			name:             "invalid leaves options",
			keyboardProtocol: "bogus",
			negotiate:        "bogus",
			startForceKbd:    true,
			startForcedKbd:   XTermKeyboard,
			startNegotiate:   false,
			wantForceKbd:     true,
			wantForcedKbd:    XTermKeyboard,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TCELL_KEYBOARD_PROTOCOL", tt.keyboardProtocol)
			t.Setenv("TCELL_NEGOTIATE", tt.negotiate)
			t.Setenv("TCELL_MOUSE", tt.mouse)
			s := &tScreen{
				forceKbd:  tt.startForceKbd,
				forcedKbd: tt.startForcedKbd,
				negotiate: tt.startNegotiate,
			}
			s.applyEnvironmentOverrides()
			if s.forceKbd != tt.wantForceKbd {
				t.Fatalf("forceKbd = %v, want %v", s.forceKbd, tt.wantForceKbd)
			}
			if s.forcedKbd != tt.wantForcedKbd {
				t.Fatalf("forcedKbd = %v, want %v", s.forcedKbd, tt.wantForcedKbd)
			}
			if s.negotiate != tt.wantNegotiate {
				t.Fatalf("negotiate = %v, want %v", s.negotiate, tt.wantNegotiate)
			}
			if s.mouseDisabled != tt.wantMouseOff {
				t.Fatalf("mouseDisabled = %v, want %v", s.mouseDisabled, tt.wantMouseOff)
			}
		})
	}
}

func TestForcedKeyboardProtocolSkipsKeyboardQueries(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "")
	tty := &spyTty{MockTerm: vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})}
	s, err := NewTerminfoScreenFromTty(tty, OptKeyboardProtocol(KittyKeyboard), OptAdvancedKeys(true), OptAltScreen(false))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer s.Fini()

	out := tty.Output()
	for _, seq := range []string{vt.PmWin32Input.Query(), queryKittyKbd, queryXTermKbd} {
		if strings.Contains(out, seq) {
			t.Fatalf("forced keyboard protocol emitted competing query %q", seq)
		}
	}
	for _, seq := range []string{vt.PmResizeReports.Query(), vt.PmMouseButton.Query(), vt.PmMouseSgr.Query(), requestExtAttr} {
		if !strings.Contains(out, seq) {
			t.Fatalf("forced keyboard protocol suppressed unrelated query %q", seq)
		}
	}
	if !strings.Contains(out, enableKittyKbdAdv) {
		t.Fatalf("forced Kitty protocol did not emit advanced enable sequence")
	}
}

func TestNegotiationDisabledSkipsStartupQueries(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "")
	tty := &spyTty{MockTerm: vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})}
	s, err := NewTerminfoScreenFromTty(tty, OptNegotiation(false), OptAltScreen(false))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer s.Fini()

	out := tty.Output()
	for _, seq := range []string{
		requestWindowSize,
		vt.PmResizeReports.Query(),
		vt.PmMouseButton.Query(),
		vt.PmMouseSgr.Query(),
		vt.PmWin32Input.Query(),
		queryKittyKbd,
		queryXTermKbd,
		requestExtAttr,
		requestPrimaryDA,
	} {
		if strings.Contains(out, seq) {
			t.Fatalf("disabled negotiation emitted startup query %q", seq)
		}
	}
}

func TestMouseDisabledPreventsEnablement(t *testing.T) {
	t.Setenv("TCELL_MOUSE", "disable")
	tty := &spyTty{MockTerm: vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})}
	s, err := NewTerminfoScreenFromTty(tty, OptNegotiation(false), OptAltScreen(false))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer s.Fini()

	s.EnableMouse()
	out := tty.Output()
	for _, seq := range []string{vt.PmMouseButton.Enable(), vt.PmMouseDrag.Enable(), vt.PmMouseMotion.Enable(), vt.PmMouseSgr.Enable()} {
		if strings.Contains(out, seq) {
			t.Fatalf("disabled mouse emitted enable sequence %q", seq)
		}
	}
}

func TestMouseWithoutSgrPreventsEnablement(t *testing.T) {
	tty := &spyTty{MockTerm: vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})}
	s, err := NewTerminfoScreenFromTty(tty, OptNegotiation(false), OptAltScreen(false))
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if err := s.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer s.Fini()

	ts := s.(*baseScreen).screenImpl.(*tScreen)
	ts.haveMouse = true
	ts.haveMouseSgr = false
	s.EnableMouse()

	out := tty.Output()
	for _, seq := range []string{vt.PmMouseButton.Enable(), vt.PmMouseDrag.Enable(), vt.PmMouseMotion.Enable(), vt.PmMouseSgr.Enable()} {
		if strings.Contains(out, seq) {
			t.Fatalf("mouse without SGR support emitted enable sequence %q", seq)
		}
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

func TestSetTitleStripsOSCControls(t *testing.T) {
	mt := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	scr, err := NewTerminfoScreenFromTty(mt)
	if err != nil {
		t.Fatalf("failed to get terminal: %v", err)
	}

	scr.SetTitle("good\x07title\x1b\\end")
	if err := scr.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer scr.Fini()

	if got := mt.GetTitle(); got != "goodtitle\\end" {
		t.Fatalf("title not sanitized: %q", got)
	}
}

func TestShowNotificationStripsOSCControls(t *testing.T) {
	mt := &spyTty{MockTerm: vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})}
	scr, err := NewTerminfoScreenFromTty(mt)
	if err != nil {
		t.Fatalf("failed to get terminal: %v", err)
	}
	if err := scr.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}
	defer scr.Fini()

	before := mt.Output()
	scr.ShowNotification("tit\x07le", "bo\x1b\\dy")
	delta := mt.Output()[len(before):]

	if strings.Contains(delta, "tit\x07le") || strings.Contains(delta, "bo\x1b\\dy") {
		t.Fatalf("notification payload still contains control characters: %q", delta)
	}
	if !strings.Contains(delta, "title") || !strings.Contains(delta, "bo\\dy") {
		t.Fatalf("notification payload missing sanitized strings: %q", delta)
	}
}
