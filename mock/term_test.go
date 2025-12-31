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

package mock

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/vt"
)

// writeF writes the string, and ensures it is fully flushed
// before returning.
func writeF(t *testing.T, trm MockTerm, str string, args ...any) {
	t.Helper()
	b := fmt.Appendf(nil, str, args...)
	for len(b) != 0 {
		if n, err := trm.Write(b); err != nil {
			t.Fatalf("Failed to write: %v", err)
		} else {
			b = b[n:]
		}
	}
	if err := trm.Drain(); err != nil {
		t.Fatalf("Failed to flush: %v", err)
	}
}

func mustClose(t *testing.T, trm MockTerm) {
	t.Helper()
	if err := trm.Close(); err != nil {
		t.Errorf("close failed: %v", err)
	}
}

func mustStart(t *testing.T, trm MockTerm) {
	t.Helper()
	if err := trm.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}
}

func checkPos(t *testing.T, trm MockTerm, x vt.Col, y vt.Row) {
	t.Helper()
	if trm.Pos().X != x || trm.Pos().Y != y {
		t.Errorf("bad position %d,%d (expected %d,%d)", trm.Pos().X, trm.Pos().Y, x, y)
	}
}

// TestCursorMove tests several aspects of cursor movement.
func TestCursorMovement(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, trm)

	mustStart(t, trm)

	if size, err := trm.WindowSize(); err != nil {
		t.Fatalf("failed getting window size: %v", err)
	} else if size.Height != 3 || size.Width != 5 {
		t.Fatalf("wrong window size X %d Y %d", size.Width, size.Height)
	}
	writeF(t, trm, "\x1b[2;3H")
	checkPos(t, trm, 2, 1)

	writeF(t, trm, "\x1b[20A") // up 20
	checkPos(t, trm, 2, 0)

	writeF(t, trm, "\x1b[20B") // down 20
	checkPos(t, trm, 2, 2)

	writeF(t, trm, "\x1b[A") // up 1
	checkPos(t, trm, 2, 1)

	writeF(t, trm, "\x1b[2C") // right 2
	checkPos(t, trm, 4, 1)

	writeF(t, trm, "\x1b[3D") // left 3
	checkPos(t, trm, 1, 1)

	writeF(t, trm, "\x1b[100D") // left 100
	checkPos(t, trm, 0, 1)

	// Now try the next line and previous line
	writeF(t, trm, "\x1b[2;3H")
	checkPos(t, trm, 2, 1)

	writeF(t, trm, "\x1b[1E")
	checkPos(t, trm, 0, 2)

	writeF(t, trm, "\x1b[2;3H")
	checkPos(t, trm, 2, 1)

	writeF(t, trm, "\x1b[1F")
	checkPos(t, trm, 0, 0)
}

// TestDECALN tests the DEC alignment test (screen filled with 'E').
func TestDECALN(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, trm)

	mustStart(t, trm)

	writeF(t, trm, "\x1b#8")

	for y := range vt.Row(3) {
		for x := range vt.Col(5) {
			cell := trm.GetCell(vt.Coord{X: x, Y: y})
			if cell.Attr != vt.Plain {
				t.Errorf("wrong attr at %d,%d: %x", x, y, cell.Attr)
			}
			if string(cell.C) != "E" {
				t.Errorf("wrong content at %d,%d: %q", x, y, string(cell.C))
			}
		}
	}

	// clear screen
	writeF(t, trm, "\x1b[H\x1b[J")

	for y := range vt.Row(3) {
		for x := range vt.Col(5) {
			cell := trm.GetCell(vt.Coord{X: x, Y: y})
			if cell.Attr != vt.Plain {
				t.Errorf("wrong attr at %d,%d: %x", x, y, cell.Attr)
			}
			if string(cell.C) != "" {
				t.Errorf("wrong content at %d,%d: %q", x, y, string(cell.C))
			}
		}
	}

	writeF(t, trm, "\x1b[1m\x1b#8") // bold, DECALN
	for y := range vt.Row(3) {
		for x := range vt.Col(5) {
			cell := trm.GetCell(vt.Coord{X: x, Y: y})
			if cell.Attr != vt.Bold {
				t.Errorf("wrong attr at %d,%d: %x", x, y, cell.Attr)
			}
			if string(cell.C) != "E" {
				t.Errorf("wrong content at %d,%d: %q", x, y, string(cell.C))
			}
		}
	}
}

// TestBell tests the bell.
func TestBell(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, trm)
	mustStart(t, trm)

	if trm.Bells() != 0 {
		t.Errorf("wrong bell count: %d", trm.Bells())
	}
	writeF(t, trm, "\x07")
	if trm.Bells() != 1 {
		t.Errorf("wrong bell count: %d", trm.Bells())
	}
	writeF(t, trm, "\x07")
	if trm.Bells() != 2 {
		t.Errorf("wrong bell count: %d", trm.Bells())
	}
}

// TestPrimaryDA tests primary device attributes using several mechanisms.
func TestPrimaryDA(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, trm)

	mustStart(t, trm)

	buf := make([]byte, 32)
	writeF(t, trm, "\x1b[c")

	n, err := trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	if !strings.HasSuffix(result, "c") {
		t.Errorf("Missing suffix 'c': %q", result)
	}
	if !strings.HasPrefix(result, "\x1b[?63") {
		t.Errorf("Missing prefix '\x1b[?63': %q", result)
	}

	// Legacy version
	buf = make([]byte, 32) // to start over
	writeF(t, trm, "\x1bZ")

	n, err = trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result = string(buf[:n])
	if !strings.HasSuffix(result, "c") {
		t.Errorf("Missing suffix 'c': %q", result)
	}
	if !strings.HasPrefix(result, "\x1b[?63") {
		t.Errorf("Missing prefix '\x1b[?63': %q", result)
	}
}

// TestExtendedAttr requests ther terminal ID using CSI > q.
func TestExtendedAttr(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, trm)

	mustStart(t, trm)

	buf := make([]byte, 64)
	writeF(t, trm, "\x1b[>q")

	n, err := trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}

	result := string(buf[:n])

	if !strings.HasSuffix(result, "\x1b\\") {
		t.Errorf("Missing suffix ST: %q", result)
	}
	if !strings.HasPrefix(result, "\x1bP>|") {
		t.Errorf("Missing prefix 'ESC P>|': %q", result)
	}
	if !strings.Contains(result, "TcellMock 1.0") {
		t.Errorf("Missing terminal identification")
	}
}

// TestCursorReport verifies that cursor position reporting works.
func TestCursorReport(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, trm)

	mustStart(t, trm)

	writeF(t, trm, "\x1b[5;10H") // fifth row, tenth column
	checkPos(t, trm, 9, 4)

	writeF(t, trm, "\x1b[6n") // cursor position report

	buf := make([]byte, 32)
	n, err := trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	if result != "\x1b[5;10R" {
		t.Errorf("wrong report: %q", result)
	}

	buf = make([]byte, 32)
	// move the cursor back one
	writeF(t, trm, "\b\x1b[6n")
	checkPos(t, trm, 8, 4)
	n, err = trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result = string(buf[:n])
	if result != "\x1b[5;9R" {
		t.Errorf("wrong report: %q", result)
	}
}

// TestPrivateModes tests the private mode feature.
func TestPrivateModes(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, trm)

	mustStart(t, trm)

	writeF(t, trm, "\x1b[?7$p")              // query for automargin (should start on by default)
	writeF(t, trm, "\x1b[?7l")               // turn it off
	writeF(t, trm, "\x1b[?7$p")              // should readback positive
	writeF(t, trm, "\x1b[?7h")               // put it back on
	writeF(t, trm, "\x1b[?7$p")              // should read back negative
	writeF(t, trm, "\x1b[?1919$p")           // read invalid mode
	writeF(t, trm, "\x1b[?1919h\x1b[?1919l") // toogle invalid mode
	writeF(t, trm, "\x1b[?1919$p")           // read invalid mode one more time

	buf := make([]byte, 128)
	n, err := trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	want := "\x1b[?7;1$y" + "\x1b[?7;2$y" + "\x1b[?7;1$y" + "\x1b[?1919;0$y" + "\x1b[?1919;0$y"
	if result != want {
		t.Errorf("wrong response: %q != %q", result, want)
	}

	// Lets also test the cursor (show vs hide)
	writeF(t, trm, "\x1b[?25$p")
	writeF(t, trm, "\x1b[?25l")
	writeF(t, trm, "\x1b[?25$p")
	writeF(t, trm, "\x1b[?25h")
	writeF(t, trm, "\x1b[?25$p")

	buf = make([]byte, 128)
	n, err = trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result = string(buf[:n])

	want = "\x1b[?25;1$y" + "\x1b[?25;2$y" + "\x1b[?25;1$y"
	if result != want {
		t.Errorf("wrong response: %q != %q", result, want)
	}
}

// TestAutoMargin tests the auto margin feature.
func TestAutoMargin(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, trm)
	mustStart(t, trm)

	// default is automargin is enabled
	writeF(t, trm, "\x1b[2J") // clear the screen
	writeF(t, trm, "\x1b[1;80HAB")
	checkPos(t, trm, 1, 1)
	if s := string(trm.GetCell(vt.Coord{X: 79, Y: 0}).C); s != "A" {
		t.Errorf("last column wrong: %q", s)
	}
	if s := string(trm.GetCell(vt.Coord{X: 0, Y: 1}).C); s != "B" {
		t.Errorf("auto wrap did not work: %q", s)
	}

	// now turn it off
	writeF(t, trm, "\x1b[?7l")

	// mess with 3rd row
	writeF(t, trm, "\x1b[3;80HCD")
	checkPos(t, trm, 79, 2)
	if s := string(trm.GetCell(vt.Coord{X: 79, Y: 2}).C); s != "D" {
		t.Errorf("last column wrong: %q", s)
	}
}

// TestUnicode tests basic placement of unicode glyphs.
// For now it assumes that the terminal itself supports unicode / latin 1.
func TestUnicode(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[2J") // clear the screen
	writeF(t, trm, "\x1b[2;2H")
	checkPos(t, trm, 1, 1)
	writeF(t, trm, "åßcπ")
	checkPos(t, trm, 5, 1)
	if s := string(trm.GetCell(vt.Coord{X: 1, Y: 1}).C); s != "å" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(vt.Coord{X: 2, Y: 1}).C); s != "ß" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(vt.Coord{X: 3, Y: 1}).C); s != "c" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(vt.Coord{X: 4, Y: 1}).C); s != "π" {
		t.Errorf("decode error wrong: %q", s)
	}
}

// TestUnicodeWide tests a wide unicode glyph.
func TestUnicodeWide(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b#8") // fill it with E's (so we can see that wide clears the next cell)
	writeF(t, trm, "\x1b[2;2H")
	checkPos(t, trm, 1, 1)
	writeF(t, trm, "å宽cπ")
	checkPos(t, trm, 6, 1)
	if s := string(trm.GetCell(vt.Coord{X: 1, Y: 1}).C); s != "å" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(vt.Coord{X: 2, Y: 1}).C); s != "宽" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(vt.Coord{X: 3, Y: 1}).C); s != "" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(vt.Coord{X: 4, Y: 1}).C); s != "c" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(vt.Coord{X: 5, Y: 1}).C); s != "π" {
		t.Errorf("decode error wrong: %q", s)
	}
}

// TestKbdEventLegacy tests key events when using the default legacy key protocol.
func TestKbdEventLegacy(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, trm)
	mustStart(t, trm)

	trm.KeyEvent(vt.KbdEvent{Code: 'a', Base: 'a', Down: true})
	trm.KeyEvent(vt.KbdEvent{Code: 'A', Base: 'a', Mod: vt.ModShift, Down: false})
	trm.KeyEvent(vt.KbdEvent{Code: 'B', Base: 'b', Mod: vt.ModShift, Down: true})
	trm.KeyEvent(vt.KbdEvent{Code: 'B', Base: 'b', Mod: vt.ModShift, Down: false})
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcReturn, Down: true})
	trm.KeyEvent(vt.KbdEvent{Code: 'i', Down: true, Mod: vt.ModCtrl})
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcEsc, Down: true})

	buf := make([]byte, 256)
	n, err := trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	want := "aB\r\x09\x1B"
	if result != want {
		t.Errorf("key responses failed: %q != %q", result, want)
	}

	// SS3 based F-keys
	clear(buf)
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF1, Down: true})                                            // SS3 P
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF1, Down: false})                                           // none
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF1, Mod: vt.ModShift, Down: true})                          // CSI 1 ; 2 P
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF2, Mod: vt.ModCtrl, Down: true})                           // CSI 1 ; 5 Q
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF3, Mod: vt.ModAlt | vt.ModShift | vt.ModCtrl, Down: true}) // ESC CSI 1 ; 6 R
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF4, Mod: vt.ModAlt | vt.ModCtrl, Down: true})               // ESC CSI 1 ; 5 S
	want = "\x1bOP"
	want += "\x1b[1;2P"
	want += "\x1b[1;5Q"
	want += "\x1b\x1b[1;6R"
	want += "\x1b\x1b[1;5S"
	n, err = trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result = string(buf[:n])
	if result != want {
		t.Errorf("key responses failed: %q != %q", result, want)
	}

	// CSI based F-keys
	buf = make([]byte, 256)
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF5, Down: true})                                            // CSI 15 ~
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF5, Down: false})                                           // none
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF6, Mod: vt.ModShift, Down: true})                          // CSI 17 ; 2 ~
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF7, Mod: vt.ModCtrl, Down: true})                           // CSI 18 ; 5 ~
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF8, Mod: vt.ModAlt | vt.ModShift | vt.ModCtrl, Down: true}) // ESC CSI 19 ; 6 ~
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF9, Mod: vt.ModAlt | vt.ModCtrl, Down: true})               // ESC CSI 20 ; 5 ~
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF20, Mod: vt.ModNone, Down: true})                          // CSI 34 ~
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcHelp, Mod: vt.ModNone, Down: true})                         // CSI 28 ~
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcF15, Mod: vt.ModNone, Down: true})                          // CSI 28 ~
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcMenu, Mod: vt.ModNone, Down: true})                         // CSI 29 ~
	want = "\x1b[15~"
	want += "\x1b[17;2~"
	want += "\x1b[18;5~"
	want += "\x1b\x1b[19;6~"
	want += "\x1b\x1b[20;5~"
	want += "\x1b[34~"
	want += "\x1b[28~"
	want += "\x1b[28~"
	want += "\x1b[29~"
	n, err = trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result = string(buf[:n])
	if result != want {
		t.Errorf("key responses failed: %q != %q", result, want)
	}

	// Misc other keys
	clear(buf)
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcReturn, Down: true})                                   // \r
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcTab, Down: true})                                      // \t
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcTab, Mod: vt.ModShift, Down: true})                    // CSI Z
	trm.KeyEvent(vt.KbdEvent{Code: 'm', Mod: vt.ModCtrl, Down: true})                          // \r
	trm.KeyEvent(vt.KbdEvent{Code: 'l', Mod: vt.ModCtrl, Down: true})                          // \x0c
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcBackspace, Down: true})                                // \x7f
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcBackspace, Mod: vt.ModCtrl, Down: true})               // \x08
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcBackspace, Mod: vt.ModShift | vt.ModCtrl, Down: true}) // \x08
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcSpace, Mod: vt.ModCtrl, Down: true})                   // \x00
	trm.KeyEvent(vt.KbdEvent{Code: vt.KcSpace, Down: true})                                    // ' '
	trm.KeyEvent(vt.KbdEvent{Code: 'a', Mod: vt.ModAlt, Down: true})                           // \x1b a
	trm.KeyEvent(vt.KbdEvent{Code: 'a', Mod: vt.ModHyper, Down: true})                         // none
	trm.KeyEvent(vt.KbdEvent{Code: 'a', Mod: vt.ModMeta, Down: true})                          // none
	trm.KeyEvent(vt.KbdEvent{Code: 'j', Mod: vt.ModAlt | vt.ModCtrl, Down: true})              // \x1b\n
	trm.KeyEvent(vt.KbdEvent{Code: 'L', Mod: vt.ModCtrl, Down: true})                          // \x0c
	trm.KeyEvent(vt.KbdEvent{Code: '[', Mod: vt.ModCtrl, Down: true})                          // \x0c

	want = "\r\t\x1b[Z\r\x0c\x7f\x08\x08\x00 \x1ba\x1b\n\x0c\x1b"
	n, err = trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result = string(buf[:n])
	if result != want {
		t.Errorf("key responses failed: %q != %q", result, want)
	}

	// Legacy control key mappings (weird ones)
	// 	clear(buf)
	trm.KeyEvent(vt.KbdEvent{Code: '8', Mod: vt.ModCtrl, Down: true}) // \x7F
	trm.KeyEvent(vt.KbdEvent{Code: '4', Mod: vt.ModCtrl, Down: true}) // \x1c
	trm.KeyEvent(vt.KbdEvent{Code: '?', Mod: vt.ModCtrl, Down: true}) // \x1f
	trm.KeyEvent(vt.KbdEvent{Code: '7', Mod: vt.ModCtrl, Down: true}) // \x1f
	trm.KeyEvent(vt.KbdEvent{Code: '7', Mod: vt.ModNone, Down: true}) // 7
	trm.KeyEvent(vt.KbdEvent{Code: '?', Mod: vt.ModNone, Down: true}) // ?
	trm.KeyEvent(vt.KbdEvent{Code: '[', Mod: vt.ModCtrl, Down: true}) // \x1b
	want = "\x7f\x1c\x1f\x1f7?\x1b"
	n, err = trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result = string(buf[:n])
	if result != want {
		t.Errorf("key responses failed: %q != %q", result, want)
	}
}

// TestSgrAttr tests a variety of combinations of Sgr settings.
func TestSgrAttr(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[H")
	writeF(t, trm, "\x1b[1mA") // bold
	if attr := trm.GetCell(vt.Coord{X: 0, Y: 0}).Attr; attr != vt.Bold {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[2mB")                                          // dim
	if attr := trm.GetCell(vt.Coord{X: 1, Y: 0}).Attr; attr != vt.Dim { // dim is exclusive of bold
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[22mC") // clear both
	if attr := trm.GetCell(vt.Coord{X: 2, Y: 0}).Attr; attr != vt.Plain {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[3;2mD") // italic, dim
	if attr := trm.GetCell(vt.Coord{X: 3, Y: 0}).Attr; attr != vt.Italic|vt.Dim {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[22mE") // dim, should leave italic
	if attr := trm.GetCell(vt.Coord{X: 4, Y: 0}).Attr; attr != vt.Italic {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[23mF") // clear italic
	if attr := trm.GetCell(vt.Coord{X: 5, Y: 0}).Attr; attr != vt.Plain {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[3;4mG") // simple underline
	if attr := trm.GetCell(vt.Coord{X: 6, Y: 0}).Attr; attr != vt.Italic|vt.Underline {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[21mH") // double underline (ECMA)
	if attr := trm.GetCell(vt.Coord{X: 7, Y: 0}).Attr; attr != vt.Italic|vt.DoubleUnderline {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[4mI") // simple underline
	if attr := trm.GetCell(vt.Coord{X: 8, Y: 0}).Attr; attr != vt.Italic|vt.Underline {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[4:2mJ") // simple underline
	if attr := trm.GetCell(vt.Coord{X: 9, Y: 0}).Attr; attr != vt.Italic|vt.DoubleUnderline {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[4:3mI") // curly underline
	if attr := trm.GetCell(vt.Coord{X: 10, Y: 0}).Attr; attr != vt.Italic|vt.CurlyUnderline {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[4:4mJ") // dotted underline
	if attr := trm.GetCell(vt.Coord{X: 11, Y: 0}).Attr; attr != vt.Italic|vt.DottedUnderline {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[4:5mK") // dashed underline
	if attr := trm.GetCell(vt.Coord{X: 12, Y: 0}).Attr; attr != vt.Italic|vt.DashedUnderline {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[4:9mL") // junk treats as plain
	if attr := trm.GetCell(vt.Coord{X: 13, Y: 0}).Attr; attr != vt.Italic|vt.Underline {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[4:5;24mM") // clustering, clear it
	if attr := trm.GetCell(vt.Coord{X: 14, Y: 0}).Attr; attr != vt.Italic {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[0;9;7;53mN") // clear, strikethrough, reverse, overlined
	if attr := trm.GetCell(vt.Coord{X: 15, Y: 0}).Attr; attr != vt.StrikeThrough|vt.Reverse|vt.Overline {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[5;27;29;55mO")
	if attr := trm.GetCell(vt.Coord{X: 16, Y: 0}).Attr; attr != vt.Blink {
		t.Errorf("wrong attr: %x", attr)
	}
	writeF(t, trm, "\x1b[25mP")
	if attr := trm.GetCell(vt.Coord{X: 17, Y: 0}).Attr; attr != vt.Plain {
		t.Errorf("wrong attr: %x", attr)
	}
}

// TestSgrColor8 tests simple ECMA 48 ANSI color (only 8 possible color values.)
func TestSgrColor8(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[36;42m\x1b#8")
	if fg := trm.GetCell(vt.Coord{X: 0, Y: 0}).Fg; fg != color.Teal {
		t.Errorf("wrong fg: %s\n", fg.String())
	}
	if bg := trm.GetCell(vt.Coord{X: 0, Y: 0}).Bg; bg != color.Green {
		t.Errorf("wrong bg: %s\n", bg.String())
	}
	writeF(t, trm, "\x1b[H\x1b[39mA")
	if fg := trm.GetCell(vt.Coord{X: 0, Y: 0}).Fg; fg != color.Silver {
		t.Errorf("wrong fg: %s\n", fg.String())
	}
	if bg := trm.GetCell(vt.Coord{X: 0, Y: 0}).Bg; bg != color.Green {
		t.Errorf("wrong bg: %s\n", bg.String())
	}
	writeF(t, trm, "\x1b[49mA")
	if fg := trm.GetCell(vt.Coord{X: 1, Y: 0}).Fg; fg != color.Silver {
		t.Errorf("wrong fg: %s\n", fg.String())
	}
	if bg := trm.GetCell(vt.Coord{X: 1, Y: 0}).Bg; bg != color.Black {
		t.Errorf("wrong bg: %s\n", bg.String())
	}

	// verify zero clears colors, first write some non zero colors
	writeF(t, trm, "\x1b[36;42mD")
	if fg := trm.GetCell(vt.Coord{X: 2, Y: 0}).Fg; fg != color.Teal {
		t.Errorf("wrong fg: %s\n", fg.String())
	}
	if bg := trm.GetCell(vt.Coord{X: 2, Y: 0}).Bg; bg != color.Green {
		t.Errorf("wrong bg: %s\n", bg.String())
	}
	// then send zero, should go to default colors
	writeF(t, trm, "\x1b[0mA")
	if fg := trm.GetCell(vt.Coord{X: 3, Y: 0}).Fg; fg != color.Silver {
		t.Errorf("wrong fg: %s\n", fg.String())
	}
	if bg := trm.GetCell(vt.Coord{X: 3, Y: 0}).Bg; bg != color.Black {
		t.Errorf("wrong bg: %s\n", bg.String())
	}
}

// TestTitler tests that we can set a window title.
func TestTitler(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)
	writeF(t, trm, "\x1b]2;Test Application\x1b\\")
	if s := trm.GetTitle(); s != "Test Application" {
		t.Errorf("wrong title: %q", s)
	}

	// test ST termination using legacy bell character
	writeF(t, trm, "\x1b]2;Bell Ring\x07")
	if s := trm.GetTitle(); s != "Bell Ring" {
		t.Errorf("wrong title: %q", s)
	}

	// try using 8-bit sequence
	writeF(t, trm, "\x9d2;Eight Bits\x9c")
	if s := trm.GetTitle(); s != "Eight Bits" {
		t.Errorf("wrong title: %q", s)
	}
}

// TestResize tests resizing the terminal
func TestResize(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	// with E, and enable notifications
	writeF(t, trm, "\x1b#8\x1b[?2048h")
	resizeQ := make(chan bool, 1)
	trm.NotifyResize(resizeQ)

	trm.SetSize(vt.Coord{X: 132, Y: 24})
	if sz, err := trm.WindowSize(); err != nil || sz.Height != 24 || sz.Width != 132 {
		t.Errorf("resize did not occur: %v %d %d", err, sz.Height, sz.Width)
	}
	for y := range vt.Row(24) {
		for x := range vt.Col(80) {
			if s := string(trm.GetCell(vt.Coord{X: x, Y: y}).C); s != "E" {
				t.Errorf("resize content at %d,%d wrong: %q", x, y, s)
			}
		}
	}
	for y := range vt.Row(24) {
		for x := vt.Col(80); x < 132; x++ {
			if s := string(trm.GetCell(vt.Coord{X: x, Y: y}).C); s != "" {
				t.Errorf("resize content at %d,%d wrong: %q", x, y, s)
			}
		}
	}
	select {
	case <-resizeQ:
	case <-time.After(time.Millisecond * 100):
		t.Errorf("resize signal failure")
	}

	buf := make([]byte, 128)
	n, err := trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	if result != "\x1b[48;24;132;0;0t" && result != "\x1b[48;24;132t" {
		t.Errorf("key responses failed: %q", result)
	}
}

// TestTabs tests tab stop funtionality.
func TestTabs(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "a\tC")
	if s := string(trm.GetCell(vt.Coord{X: 8, Y: 0}).C); s != "C" {
		t.Errorf("tab did not work: %q", s)
	}
	writeF(t, trm, "\x1b[3I")
	if x := trm.Pos().X; x != 32 {
		t.Errorf("wrong position %d", x)
	}

	writeF(t, trm, "\x1b[2Z")
	if x := trm.Pos().X; x != 16 {
		t.Errorf("wrong position %d", x)
	}

	writeF(t, trm, "\x1b[3g") // clear all tabs
	writeF(t, trm, "\x1b[I")  // one tab, should go to right margin
	if x := trm.Pos().X; x != 79 {
		t.Errorf("wrong position: %d", x)
	}

	writeF(t, trm, "\x1b[Z")
	if x := trm.Pos().X; x != 0 {
		t.Errorf("wrong position: %d", x)
	}

	// reset tabs
	writeF(t, trm, "\x1b[?5W")

	writeF(t, trm, "\t")
	if x := trm.Pos().X; x != 8 {
		t.Errorf("wrong position: %d", x)
	}
	// clear this position, advance one
	writeF(t, trm, "\x1b[gA")
	if x := trm.Pos().X; x != 9 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, trm, "\x1bH")
	writeF(t, trm, "\t")
	if x := trm.Pos().X; x != 16 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, trm, "\x1b[Z")
	if x := trm.Pos().X; x != 9 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, trm, "\x1b[Z")
	if x := trm.Pos().X; x != 0 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, trm, "\x1b[1;10H") // goto position 9
	if x := trm.Pos().X; x != 9 {
		t.Errorf("wrong position: %d", x)
	}
	// delete this one (do it twice to exericse the does not exist flow)
	writeF(t, trm, "\x1b[0g")
	writeF(t, trm, "\x1b[0g")

	// advance to next tab, then back, we should go to 0
	writeF(t, trm, "\t")
	if x := trm.Pos().X; x != 16 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, trm, "\x1b[Z")
	if x := trm.Pos().X; x != 0 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, trm, "\x1b[20I")
	writeF(t, trm, "\t\t")
	if pos := trm.Pos(); pos.X != 79 || pos.Y != 0 {
		t.Errorf("wrong position: %d %d", pos.X, pos.Y)
	}

	// now backwards
	writeF(t, trm, "\x1b[20Z")
	writeF(t, trm, "\x1b[Z")
	if pos := trm.Pos(); pos.X != 0 || pos.Y != 0 {
		t.Errorf("wrong position: %d %d", pos.X, pos.Y)
	}
}
