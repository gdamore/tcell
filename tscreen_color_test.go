package tcell

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v3/vt"
)

func TestDefaultColors(t *testing.T) {
	tty := &spyTty{MockTerm: vt.NewMockTerm()}
	scr, err := NewTerminfoScreenFromTty(tty)
	if err != nil {
		t.Fatalf("failed to get screen: %v", err)
	}
	if err := scr.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
	}

	const (
		setFg   = "\x1b]10;#c0caf5\x07"
		setBg   = "\x1b]11;#1a1b26\x07"
		resetFg = "\x1b]110\x07"
		resetBg = "\x1b]111\x07"
	)
	// written returns what step wrote to the terminal.
	written := func(step func()) string {
		before := len(tty.Output())
		step()
		return tty.Output()[before:]
	}
	expect := func(t *testing.T, got string, want ...string) {
		t.Helper()
		for _, w := range want {
			if !strings.Contains(got, w) {
				t.Fatalf("%q not written, got %q", w, got)
			}
		}
	}

	t.Run("set", func(t *testing.T) {
		expect(t, written(func() {
			scr.SetForegroundColor(NewHexColor(0xc0caf5))
			scr.SetBackgroundColor(NewHexColor(0x1a1b26))
		}), setFg, setBg)
	})
	t.Run("restored while suspended", func(t *testing.T) {
		expect(t, written(func() { _ = scr.Suspend() }), resetFg, resetBg)
		expect(t, written(func() { _ = scr.Resume() }), setFg, setBg)
	})
	t.Run("reset", func(t *testing.T) {
		expect(t, written(func() { scr.SetForegroundColor(ColorDefault) }), resetFg)
		if got := written(scr.Fini); strings.Contains(got, resetFg) || !strings.Contains(got, resetBg) {
			t.Fatalf("Fini wrote %q, want only the background restored", got)
		}
	})
}
