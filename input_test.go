package tcell

import (
	"testing"
	"time"
)

// firstKey drains the event channel and returns the first EventKey received,
// or nil if none arrives within 100ms.
func firstKey(evch chan Event) *EventKey {
	var got *EventKey
	for {
		select {
		case ev := <-evch:
			if got == nil {
				if kev, ok := ev.(*EventKey); ok {
					got = kev
				}
			}
			continue
		case <-time.After(100 * time.Millisecond):
		}
		break
	}
	return got
}

// TestEscDuringCsiResetsParser verifies that an ESC byte received while
// the parser is in CSI state correctly transitions to escape state per
// ECMA-48 §5.3.1, rather than being swallowed as a "bad parse".
func TestEscDuringCsiResetsParser(t *testing.T) {
	evch := make(chan Event, 10)
	ip := NewInputProcessor(evch)

	// Feed a partial CSI (ESC [) followed immediately by a new escape
	// sequence for Down arrow (ESC [ B). Without the fix, the second
	// ESC would be swallowed and 'B' would be emitted as a literal key.
	ip.ScanUTF8([]byte("\x1b[\x1b[B"))

	got := firstKey(evch)
	if got == nil {
		t.Fatal("expected a key event, got none")
	}
	if got.Key() != KeyDown {
		t.Errorf("expected KeyDown, got key=%v rune=%v mod=%v", got.Key(), got.Rune(), got.Modifiers())
	}
}

// TestEscDuringSs3ResetsParser verifies that an ESC byte received while
// the parser is in SS3 state correctly transitions to escape state per
// ECMA-48 §5.3.1.
func TestEscDuringSs3ResetsParser(t *testing.T) {
	evch := make(chan Event, 10)
	ip := NewInputProcessor(evch)

	// Feed ESC O (enters SS3) then ESC [ B (Down arrow). Without the
	// fix, the ESC is lost inside SS3 handling.
	ip.ScanUTF8([]byte("\x1bO\x1b[B"))

	got := firstKey(evch)
	if got == nil {
		t.Fatal("expected a key event, got none")
	}
	if got.Key() != KeyDown {
		t.Errorf("expected KeyDown, got key=%v rune=%v mod=%v", got.Key(), got.Rune(), got.Modifiers())
	}
}
