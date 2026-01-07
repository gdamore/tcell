// Copyright 2026 The TCell Authors
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

package vt

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/color"
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

func readF(t *testing.T, trm MockTerm) string {
	buf := make([]byte, 128)
	n, err := trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
		return ""
	}
	return string(buf[:n])
}

// verifyF validates the condition, printing the message on failure.
func verifyF(t *testing.T, b bool, fmt string, args ...any) {
	t.Helper()
	if !b {
		t.Errorf("validation failure: "+fmt, args...)
	}
}

// assertF validates the condition, and aborts the test if it fails.
func assertF(t *testing.T, b bool, fmt string, args ...any) {
	t.Helper()
	if !b {
		t.Fatalf("validation failure: "+fmt, args...)
	}
}

func mustClose(t *testing.T, trm MockTerm) {
	t.Helper()
	err := trm.Close()
	assertF(t, err == nil, "close failed: %v", err)
}

func mustStart(t *testing.T, trm MockTerm) {
	t.Helper()
	err := trm.Start()
	assertF(t, err == nil, "start failed: %v", err)
}

func checkPos(t *testing.T, trm MockTerm, x Col, y Row) {
	t.Helper()
	verifyF(t, trm.Pos().X == x && trm.Pos().Y == y,
		"bad position %d,%d (expected %d,%d)", trm.Pos().X, trm.Pos().Y, x, y)
}

func checkContent(t *testing.T, trm MockTerm, x Col, y Row, s string) {
	t.Helper()
	if actual := string(trm.GetCell(Coord{X: x, Y: y}).C); actual != s {
		t.Errorf("bad content %d,%d (expected %q got %q)", x, y, s, actual)
	}
}

func checkAttrs(t *testing.T, trm MockTerm, x Col, y Row, a Attr) {
	t.Helper()
	if actual := trm.GetCell(Coord{X: x, Y: y}).S.Attr(); actual != a {
		t.Errorf("bad attr %d,%d (expected %x got %x)", x, y, a, actual)
	}
}

func checkColors(t *testing.T, trm MockTerm, x Col, y Row, fg color.Color, bg color.Color) {
	t.Helper()
	if actual := trm.GetCell(Coord{X: x, Y: y}).S.Fg(); actual != fg {
		t.Errorf("bad foreground %d,%d (expected %s got %s)", x, y, fg.String(), actual.String())
	}
	if actual := trm.GetCell(Coord{X: x, Y: y}).S.Bg(); actual != bg {
		t.Errorf("bad background %d,%d (expected %s got %s)", x, y, bg.String(), actual.String())
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

	writeF(t, trm, "\x1b9")
	checkPos(t, trm, 1, 0)

	writeF(t, trm, "\x1b6")
	checkPos(t, trm, 0, 0)
}

// TestDECALN tests the DEC alignment test (screen filled with 'E').
func TestDECALN(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, trm)

	mustStart(t, trm)

	writeF(t, trm, "\x1b#8")

	for y := range Row(3) {
		for x := range Col(5) {
			checkAttrs(t, trm, x, y, Plain)
			checkContent(t, trm, x, y, "E")
		}
	}

	// clear screen
	writeF(t, trm, "\x1b[H\x1b[J")

	for y := range Row(3) {
		for x := range Col(5) {
			checkAttrs(t, trm, x, y, Plain)
			checkContent(t, trm, x, y, "")
		}
	}

	writeF(t, trm, "\x1b[1m\x1b#8") // bold, DECALN
	for y := range Row(3) {
		for x := range Col(5) {
			checkAttrs(t, trm, x, y, Bold)
			checkContent(t, trm, x, y, "E")
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

// TestExtendedAttr requests the terminal ID using CSI > q.
func TestExtendedAttr(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, trm)

	mustStart(t, trm)

	buf := make([]byte, 64)
	writeF(t, trm, "\x1b[>q")

	n, err := trm.Read(buf)
	assertF(t, err == nil, "read failed: %v", err)

	result := string(buf[:n])

	verifyF(t, strings.HasSuffix(result, "\x1b\\"), "missing suffix ST: %q", result)
	verifyF(t, strings.HasPrefix(result, "\x1bP>|"), "missing prefix 'ESC P>|': %q", result)
	verifyF(t, strings.Contains(result, "TcellMock 1.0"), "missing terminal identification: %q", result)
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
	if s := string(trm.GetCell(Coord{X: 79, Y: 0}).C); s != "A" {
		t.Errorf("last column wrong: %q", s)
	}
	if s := string(trm.GetCell(Coord{X: 0, Y: 1}).C); s != "B" {
		t.Errorf("auto wrap did not work: %q", s)
	}

	// now turn it off
	writeF(t, trm, "\x1b[?7l")

	// mess with 3rd row
	writeF(t, trm, "\x1b[3;80HCD")
	checkPos(t, trm, 79, 2)
	if s := string(trm.GetCell(Coord{X: 79, Y: 2}).C); s != "D" {
		t.Errorf("last column wrong: %q", s)
	}

	// turn it back on
	writeF(t, trm, "\x1b[?7h")

	// demonstrate that writing to the last column does not advance (pending)
	writeF(t, trm, "\x1b[1;80HA")
	checkPos(t, trm, 79, 0)

	// but one more character does advance
	writeF(t, trm, "\x1b[1;80HAB")
	checkPos(t, trm, 1, 1)

	// tab does not advance, but leaves pending state
	writeF(t, trm, "\x1b[1;80HA\t")
	checkPos(t, trm, 79, 0)
	writeF(t, trm, "\x1b[1;80HA\tb")
	checkPos(t, trm, 1, 1)

	// up or down movement resets the pending state
	writeF(t, trm, "\x1b[1;80HA\x1b[AB")
	checkPos(t, trm, 79, 0)
	writeF(t, trm, "\x1b[1;80HA\x1b[BB")
	checkPos(t, trm, 79, 1)

	// forward also resets pending state (which is clipped)
	writeF(t, trm, "\x1b[1;80HA\x1b[CB")
	checkPos(t, trm, 79, 0)
	writeF(t, trm, "\x1b[1;80HA\x1b[CBC")
	checkPos(t, trm, 1, 1)

	// explicit column also resets pending state (which is clipped)
	writeF(t, trm, "\x1b[1;80HA\x1b[80GB")
	checkPos(t, trm, 79, 0)
	writeF(t, trm, "\x1b[1;80HA\x1b[80GBC")
	checkPos(t, trm, 1, 1)

	// newline of course as well (and also VF and FF)
	writeF(t, trm, "\x1b[1;80HA\nB")
	checkPos(t, trm, 79, 1)
	writeF(t, trm, "\x1b[1;80HA\fB")
	checkPos(t, trm, 79, 1)
	writeF(t, trm, "\x1b[1;80HA\vB")
	checkPos(t, trm, 79, 1)

	// and also index
	writeF(t, trm, "\x1b[1;80HA\x1bDB")
	checkPos(t, trm, 79, 1)

	// but not reverse index
	writeF(t, trm, "\x1b[2;80HA\x1bMB")
	checkPos(t, trm, 1, 1)
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
	if s := string(trm.GetCell(Coord{X: 1, Y: 1}).C); s != "å" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(Coord{X: 2, Y: 1}).C); s != "ß" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(Coord{X: 3, Y: 1}).C); s != "c" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(Coord{X: 4, Y: 1}).C); s != "π" {
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
	if s := string(trm.GetCell(Coord{X: 1, Y: 1}).C); s != "å" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(Coord{X: 2, Y: 1}).C); s != "宽" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(Coord{X: 3, Y: 1}).C); s != "" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(Coord{X: 4, Y: 1}).C); s != "c" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(trm.GetCell(Coord{X: 5, Y: 1}).C); s != "π" {
		t.Errorf("decode error wrong: %q", s)
	}
}

// TestKeyEventLegacy tests key events when using the default legacy key protocol.
func TestKeyEventLegacy(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, trm)
	mustStart(t, trm)

	trm.KeyEvent(KeyEvent{Code: 'a', Base: 'a', Down: true})
	trm.KeyEvent(KeyEvent{Code: 'A', Base: 'a', Mod: ModShift, Down: false})
	trm.KeyEvent(KeyEvent{Code: 'B', Base: 'b', Mod: ModShift, Down: true})
	trm.KeyEvent(KeyEvent{Code: 'B', Base: 'b', Mod: ModShift, Down: false})
	trm.KeyEvent(KeyEvent{Code: KcReturn, Down: true})
	trm.KeyEvent(KeyEvent{Code: 'i', Down: true, Mod: ModCtrl})
	trm.KeyEvent(KeyEvent{Code: KcEsc, Down: true})

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
	trm.KeyEvent(KeyEvent{Code: KcF1, Down: true})                                   // SS3 P
	trm.KeyEvent(KeyEvent{Code: KcF1, Down: false})                                  // none
	trm.KeyEvent(KeyEvent{Code: KcF1, Mod: ModShift, Down: true})                    // CSI 1 ; 2 P
	trm.KeyEvent(KeyEvent{Code: KcF2, Mod: ModCtrl, Down: true})                     // CSI 1 ; 5 Q
	trm.KeyEvent(KeyEvent{Code: KcF3, Mod: ModAlt | ModShift | ModCtrl, Down: true}) // ESC CSI 1 ; 6 R
	trm.KeyEvent(KeyEvent{Code: KcF4, Mod: ModAlt | ModCtrl, Down: true})            // ESC CSI 1 ; 5 S
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
	trm.KeyEvent(KeyEvent{Code: KcF5, Down: true})                                   // CSI 15 ~
	trm.KeyEvent(KeyEvent{Code: KcF5, Down: false})                                  // none
	trm.KeyEvent(KeyEvent{Code: KcF6, Mod: ModShift, Down: true})                    // CSI 17 ; 2 ~
	trm.KeyEvent(KeyEvent{Code: KcF7, Mod: ModCtrl, Down: true})                     // CSI 18 ; 5 ~
	trm.KeyEvent(KeyEvent{Code: KcF8, Mod: ModAlt | ModShift | ModCtrl, Down: true}) // ESC CSI 19 ; 6 ~
	trm.KeyEvent(KeyEvent{Code: KcF9, Mod: ModAlt | ModCtrl, Down: true})            // ESC CSI 20 ; 5 ~
	trm.KeyEvent(KeyEvent{Code: KcF20, Mod: ModNone, Down: true})                    // CSI 34 ~
	trm.KeyEvent(KeyEvent{Code: KcHelp, Mod: ModNone, Down: true})                   // CSI 28 ~
	trm.KeyEvent(KeyEvent{Code: KcF15, Mod: ModNone, Down: true})                    // CSI 28 ~
	trm.KeyEvent(KeyEvent{Code: KcMenu, Mod: ModNone, Down: true})                   // CSI 29 ~
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
	assertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	verifyF(t, result == want, "key responses failed: %q != %q", result, want)

	// Misc other keys
	clear(buf)
	trm.KeyEvent(KeyEvent{Code: KcReturn, Down: true})                             // \r
	trm.KeyEvent(KeyEvent{Code: KcTab, Down: true})                                // \t
	trm.KeyEvent(KeyEvent{Code: KcTab, Mod: ModShift, Down: true})                 // CSI Z
	trm.KeyEvent(KeyEvent{Code: 'm', Mod: ModCtrl, Down: true})                    // \r
	trm.KeyEvent(KeyEvent{Code: 'l', Mod: ModCtrl, Down: true})                    // \x0c
	trm.KeyEvent(KeyEvent{Code: KcBackspace, Down: true})                          // \x7f
	trm.KeyEvent(KeyEvent{Code: KcBackspace, Mod: ModCtrl, Down: true})            // \x08
	trm.KeyEvent(KeyEvent{Code: KcBackspace, Mod: ModShift | ModCtrl, Down: true}) // \x08
	trm.KeyEvent(KeyEvent{Code: KcSpace, Mod: ModCtrl, Down: true})                // \x00
	trm.KeyEvent(KeyEvent{Code: KcSpace, Down: true})                              // ' '
	trm.KeyEvent(KeyEvent{Code: 'a', Mod: ModAlt, Down: true})                     // \x1b a
	trm.KeyEvent(KeyEvent{Code: 'a', Mod: ModHyper, Down: true})                   // none
	trm.KeyEvent(KeyEvent{Code: 'a', Mod: ModMeta, Down: true})                    // none
	trm.KeyEvent(KeyEvent{Code: 'j', Mod: ModAlt | ModCtrl, Down: true})           // \x1b\n
	trm.KeyEvent(KeyEvent{Code: 'L', Mod: ModCtrl, Down: true})                    // \x0c
	trm.KeyEvent(KeyEvent{Code: '[', Mod: ModCtrl, Down: true})                    // \x0c

	want = "\r\t\x1b[Z\r\x0c\x7f\x08\x08\x00 \x1ba\x1b\n\x0c\x1b"
	n, err = trm.Read(buf)
	assertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	verifyF(t, result == want, "key responses failed: %q != %q", result, want)

	// Legacy control key mappings (weird ones)
	// 	clear(buf)
	trm.KeyEvent(KeyEvent{Code: '8', Mod: ModCtrl, Down: true}) // \x7F
	trm.KeyEvent(KeyEvent{Code: '4', Mod: ModCtrl, Down: true}) // \x1c
	trm.KeyEvent(KeyEvent{Code: '?', Mod: ModCtrl, Down: true}) // \x1f
	trm.KeyEvent(KeyEvent{Code: '7', Mod: ModCtrl, Down: true}) // \x1f
	trm.KeyEvent(KeyEvent{Code: '7', Mod: ModNone, Down: true}) // 7
	trm.KeyEvent(KeyEvent{Code: '?', Mod: ModNone, Down: true}) // ?
	trm.KeyEvent(KeyEvent{Code: '[', Mod: ModCtrl, Down: true}) // \x1b
	want = "\x7f\x1c\x1f\x1f7?\x1b"
	n, err = trm.Read(buf)
	assertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	verifyF(t, result == want, "key responses failed: %q != %q", result, want)

	// Application cursor keys
	trm.KeyEvent(KeyEvent{Code: KcUp, Down: true})
	writeF(t, trm, "\x1b[?1h")
	trm.KeyEvent(KeyEvent{Code: KcDown, Down: true})
	want = "\x1b[A\x1bOB"
	n, err = trm.Read(buf)
	assertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	verifyF(t, result == want, "key responses failed: %q != %q", result, want)
}

// TestSgrAttr tests a variety of combinations of Sgr settings.
func TestSgrAttr(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[H")
	writeF(t, trm, "\x1b[1mA") // bold
	checkAttrs(t, trm, 0, 0, Bold)
	checkContent(t, trm, 0, 0, "A")

	writeF(t, trm, "\x1b[2mB") // dim
	checkAttrs(t, trm, 1, 0, Dim)
	checkContent(t, trm, 1, 0, "B")

	writeF(t, trm, "\x1b[22mC") // clear both
	checkAttrs(t, trm, 2, 0, Plain)
	checkContent(t, trm, 2, 0, "C")

	writeF(t, trm, "\x1b[3;2mD") // italic, dim
	checkAttrs(t, trm, 3, 0, Italic|Dim)
	checkContent(t, trm, 3, 0, "D")

	writeF(t, trm, "\x1b[22mE") // remove dim, should leave italic
	checkAttrs(t, trm, 4, 0, Italic)
	checkContent(t, trm, 4, 0, "E")

	writeF(t, trm, "\x1b[23mF") // clear italic
	checkAttrs(t, trm, 5, 0, Plain)
	checkContent(t, trm, 5, 0, "F")

	writeF(t, trm, "\x1b[3;4mG") // simple underline
	checkAttrs(t, trm, 6, 0, Italic|Underline)
	checkContent(t, trm, 6, 0, "G")

	writeF(t, trm, "\x1b[21mH") // double underline (ECMA)
	checkAttrs(t, trm, 7, 0, Italic|DoubleUnderline)
	checkContent(t, trm, 7, 0, "H")

	writeF(t, trm, "\x1b[4mI") // simple underline
	checkAttrs(t, trm, 8, 0, Italic|Underline)
	checkContent(t, trm, 8, 0, "I")

	writeF(t, trm, "\x1b[4:2mJ") // double underline
	checkAttrs(t, trm, 9, 0, Italic|DoubleUnderline)
	checkContent(t, trm, 9, 0, "J")

	writeF(t, trm, "\x1b[4:3mK") // curly underline
	checkAttrs(t, trm, 10, 0, Italic|CurlyUnderline)
	checkContent(t, trm, 10, 0, "K")

	writeF(t, trm, "\x1b[4:4mL") // dotted underline
	checkAttrs(t, trm, 11, 0, Italic|DottedUnderline)
	checkContent(t, trm, 11, 0, "L")

	writeF(t, trm, "\x1b[4:5mM") // dashed underline
	checkAttrs(t, trm, 12, 0, Italic|DashedUnderline)
	checkContent(t, trm, 12, 0, "M")

	writeF(t, trm, "\x1b[4:9mN") // junk treats as plain underline
	checkAttrs(t, trm, 13, 0, Italic|Underline)
	checkContent(t, trm, 13, 0, "N")

	writeF(t, trm, "\x1b[4:5;24mO") // clustering, clear it
	checkAttrs(t, trm, 14, 0, Italic)
	checkContent(t, trm, 14, 0, "O")

	writeF(t, trm, "\x1b[0;9;7;53mP") // clear, strike-through, reverse, over-lined
	checkAttrs(t, trm, 15, 0, StrikeThrough|Reverse|Overline)
	checkContent(t, trm, 15, 0, "P")

	writeF(t, trm, "\x1b[5;27;29;55mQ")
	checkAttrs(t, trm, 16, 0, Blink)
	checkContent(t, trm, 16, 0, "Q")

	writeF(t, trm, "\x1b[25mR")
	checkAttrs(t, trm, 17, 0, Plain)
	checkContent(t, trm, 17, 0, "R")
}

// TestSgrColor8 tests simple ECMA 48 ANSI color (only 8 possible color values.)
func TestSgrColor8(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[36;42m\x1b#8")
	checkColors(t, trm, 0, 0, color.Teal, color.Green)

	writeF(t, trm, "\x1b[H\x1b[39mA")
	checkColors(t, trm, 0, 0, color.Silver, color.Green)

	writeF(t, trm, "\x1b[49mA")
	checkColors(t, trm, 1, 0, color.Silver, color.Black)

	// verify zero clears colors, first write some non zero colors
	writeF(t, trm, "\x1b[36;42mD")
	checkColors(t, trm, 2, 0, color.Teal, color.Green)

	// then send zero, should go to default colors
	writeF(t, trm, "\x1b[0mA")
	checkColors(t, trm, 3, 0, color.Silver, color.Black)
}

// TestSgrColor256 tests simple ECMA 48 ANSI color (256 possible color values.)
func TestSgrColor256(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(256))
	defer mustClose(t, trm)
	mustStart(t, trm)

	// foreground
	writeF(t, trm, "\x1b[38:5:6m\x1b[42m\x1b#8")
	checkColors(t, trm, 0, 0, color.Teal, color.Green)
	writeF(t, trm, "\x1b[38;5;5m\x1b[42mA")
	checkColors(t, trm, 0, 0, color.XTerm5, color.Green)

	writeF(t, trm, "\x1b[38;5;212;1mB")
	checkColors(t, trm, 1, 0, color.XTerm212, color.Green)
	checkAttrs(t, trm, 1, 0, Bold)

	// background
	writeF(t, trm, "\x1b[48;5;2mA")
	checkColors(t, trm, 2, 0, color.XTerm212, color.XTerm2)
	checkAttrs(t, trm, 2, 0, Bold)

	writeF(t, trm, "\x1b[48:5:134;2mC")
	checkColors(t, trm, 3, 0, color.XTerm212, color.XTerm134)
	checkAttrs(t, trm, 3, 0, Dim)

	// mix background and foreground using colons
	writeF(t, trm, "\x1b[48:5:135;38:5:22mC")
	checkColors(t, trm, 4, 0, color.XTerm22, color.XTerm135)

	// and using semicolons
	writeF(t, trm, "\x1b[48;5;136;38;5;23mC")
	checkColors(t, trm, 5, 0, color.XTerm23, color.XTerm136)

	// underline colors - it uses the same parser so we won't check
	// all the variations
	writeF(t, trm, "\x1b[58;5;21;4mC")
	verifyF(t, trm.GetCell(Coord{X: 6, Y: 0}).S.Uc() == color.XTerm21, "underline color is wrong")

	writeF(t, trm, "\x1b[H\x1b[39;49mA")
	checkColors(t, trm, 0, 0, color.Silver, color.Black)

	// verify zero clears colors, first write some non zero colors
	writeF(t, trm, "\x1b[H\x1b[36;42mA")
	checkColors(t, trm, 0, 0, color.Teal, color.Green)

	// then send zero, should go to default colors
	writeF(t, trm, "\x1b[H\x1b[0mA")
	checkColors(t, trm, 0, 0, color.Silver, color.Black)

	// fuzz some things
	writeF(t, trm, "\x1b[m\x1b[H")
	writeF(t, trm, "\x1b[38:3m")
	writeF(t, trm, "\x1b[38;3m")
	writeF(t, trm, "\x1b[38;2m")
	writeF(t, trm, "\x1b[38;2:300m")
	writeF(t, trm, "\x1b[38;5m")
	writeF(t, trm, "\x1b[38;5:300m")
	writeF(t, trm, "\x1b[38:2m")
	writeF(t, trm, "\x1b[38:5m")
	writeF(t, trm, "\x1b[38:2;1;1;1m")
	writeF(t, trm, "\x1b[38:2:1:1m")
	writeF(t, trm, "A")
	checkColors(t, trm, 0, 0, color.Silver, color.Black)
}

// TestSgrColorRGB tests simple ECMA 48 ANSI color (full 24-bit color.)
func TestSgrColorRGB(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(1<<24))
	defer mustClose(t, trm)
	mustStart(t, trm)

	// foreground
	writeF(t, trm, "\x1b[m\x1b[H\x1b[38:2:255:0:0mA")
	checkColors(t, trm, 0, 0, color.NewRGBColor(255, 0, 0), color.Black)

	writeF(t, trm, "\x1b[m\x1b[H\x1b[38;2;2;0;0mA")
	checkColors(t, trm, 0, 0, color.NewRGBColor(2, 0, 0), color.Black)

	// background
	writeF(t, trm, "\x1b[m\x1b[H\x1b[48:2:1:2:3mA")
	checkColors(t, trm, 0, 0, color.Silver, color.NewRGBColor(1, 2, 3))

	writeF(t, trm, "\x1b[m\x1b[H\x1b[48;2;4;5;6mA")
	checkColors(t, trm, 0, 0, color.Silver, color.NewRGBColor(4, 5, 6))

	// full colors
	writeF(t, trm, "\x1b[m\x1b[H\x1b[38;2;99;88;77;48;2;4;5;6;58;2;99;98;91;1mA")
	checkColors(t, trm, 0, 0, color.NewRGBColor(99, 88, 77), color.NewRGBColor(4, 5, 6))
	verifyF(t, trm.GetCell(Coord{X: 0, Y: 0}).S.Uc() == color.NewRGBColor(99, 98, 91), "underline color is wrong")
	checkAttrs(t, trm, 0, 0, Bold)
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

	trm.SetSize(Coord{X: 132, Y: 24})
	if sz, err := trm.WindowSize(); err != nil || sz.Height != 24 || sz.Width != 132 {
		t.Errorf("resize did not occur: %v %d %d", err, sz.Height, sz.Width)
	}
	for y := range Row(24) {
		for x := range Col(80) {
			if s := string(trm.GetCell(Coord{X: x, Y: y}).C); s != "E" {
				t.Errorf("resize content at %d,%d wrong: %q", x, y, s)
			}
		}
	}
	for y := range Row(24) {
		for x := Col(80); x < 132; x++ {
			if s := string(trm.GetCell(Coord{X: x, Y: y}).C); s != "" {
				t.Errorf("resize content at %d,%d wrong: %q", x, y, s)
			}
		}
	}
	select {
	case <-resizeQ:
	case <-time.After(time.Millisecond * 100):
		t.Errorf("resize signal failure")
	}
	// this forces a flush of the write queue
	trm.Write([]byte{})

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

// TestTabs tests tab stop functionality.
func TestTabs(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "a\tC")
	if s := string(trm.GetCell(Coord{X: 8, Y: 0}).C); s != "C" {
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
	// delete this one (do it twice to exercise the does not exist flow)
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

func TestVerticalPos(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[2;2H")
	writeF(t, trm, "\x1b[10d")
	if pos := trm.Pos(); pos.X != 1 || pos.Y != 9 {
		t.Errorf("wrong position: %d %d", pos.X, pos.Y)
	}
	writeF(t, trm, "\x1b[2e")
	if pos := trm.Pos(); pos.X != 1 || pos.Y != 11 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, trm, "\x1b[e")
	if pos := trm.Pos(); pos.X != 1 || pos.Y != 12 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, trm, "\x1b[0e")
	if pos := trm.Pos(); pos.X != 1 || pos.Y != 13 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, trm, "\x1b[50d")
	if pos := trm.Pos(); pos.X != 1 || pos.Y != 23 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, trm, "\x1b[50e")
	if pos := trm.Pos(); pos.X != 1 || pos.Y != 23 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, trm, "\x1b[0d")
	if pos := trm.Pos(); pos.X != 1 || pos.Y != 0 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, trm, "\x1b[10d")
	if pos := trm.Pos(); pos.X != 1 || pos.Y != 9 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, trm, "\x1b[1d")
	if pos := trm.Pos(); pos.X != 1 || pos.Y != 0 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
}

func TestSaveCursorPosition(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[1;1H\x1b[J")
	writeF(t, trm, "\x1b[1;5H")
	writeF(t, trm, "A")
	writeF(t, trm, "\x1b7")
	writeF(t, trm, "\x1b[1;1H")
	writeF(t, trm, "B")
	writeF(t, trm, "\x1b8")
	writeF(t, trm, "X")

	checkContent(t, trm, 0, 0, "B")
	checkContent(t, trm, 4, 0, "A")
	checkContent(t, trm, 5, 0, "X")
}

func TestSaveCursorWrapState(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[1;1H\x1b[J")
	writeF(t, trm, "\x1b[80G")
	writeF(t, trm, "A")
	writeF(t, trm, "\x1b7") // save cursor
	writeF(t, trm, "\x1b[1;1H")
	writeF(t, trm, "B")
	writeF(t, trm, "\x1b8") // restore cursor
	writeF(t, trm, "X")
	checkContent(t, trm, 0, 0, "B")
	checkContent(t, trm, 79, 0, "A")
	checkContent(t, trm, 0, 1, "X")
}

func TestSaveCursorSgr(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)
	writeF(t, trm, "\x1b[1;1H\x1b[0J")
	writeF(t, trm, "\x1b[1;4;33;44m")
	writeF(t, trm, "A")
	checkPos(t, trm, 1, 0)
	writeF(t, trm, "\x1b7")
	checkPos(t, trm, 1, 0)
	writeF(t, trm, "\x1b[0m")
	checkPos(t, trm, 1, 0)
	writeF(t, trm, "BE")
	checkPos(t, trm, 3, 0)
	writeF(t, trm, "\x1b8")
	checkPos(t, trm, 1, 0)
	writeF(t, trm, "X")
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "X")
	checkContent(t, trm, 2, 0, "E")
	checkAttrs(t, trm, 0, 0, Bold|Underline)
	checkAttrs(t, trm, 1, 0, Bold|Underline)
	checkAttrs(t, trm, 2, 0, Plain)
	checkColors(t, trm, 0, 0, color.XTerm3, color.XTerm4)
	checkColors(t, trm, 1, 0, color.XTerm3, color.XTerm4)
	checkColors(t, trm, 2, 0, color.Silver, color.Black)
}

func TestReset(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, trm)
	mustStart(t, trm)

	// write a bunch of stuff to create state (so we can verify it gets reset)
	writeF(t, trm, "\x1b[1;4;33;44m")
	writeF(t, trm, "\x1b#8")
	writeF(t, trm, "\x1b[1;80HX")
	writeF(t, trm, "\x1b7")    // save cursor
	writeF(t, trm, "\x1b[?7l") // disable automargin
	writeF(t, trm, "\x1bc")
	checkPos(t, trm, 0, 0)
	for row := range Row(24) {
		for col := range Col(80) {
			checkAttrs(t, trm, col, row, Plain)
			checkContent(t, trm, col, row, "")
			checkColors(t, trm, col, row, color.Silver, color.Black)
		}
	}
	writeF(t, trm, "\x1b8") // restore cursor
	checkPos(t, trm, 0, 0)
	writeF(t, trm, "X")
	checkAttrs(t, trm, 0, 0, Plain)
	checkContent(t, trm, 0, 0, "X")
	checkColors(t, trm, 0, 0, color.Silver, color.Black)
	writeF(t, trm, "\x1b[?7$p")

	// verify mode reset
	want := "\x1b[?7;1$y"
	buf := make([]byte, 128)
	n, err := trm.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	verifyF(t, result == want, "wrong mode: %q != %q", result, want)
}

// backendBox makes a backend box, filled with increasing letters (modulo 16)
func backendBox(t *testing.T, mb MockBackend, tl Coord, br Coord, attr Attr) {
	t.Helper()
	mb.SetStyle(BaseStyle.WithAttr(attr))
	hex := []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P'}
	assertF(t, len(hex) == 16, "wrong hex string")
	i := 0
	for row := tl.Y; row <= br.Y; row++ {
		for col := tl.X; col <= br.X; col++ {
			mb.PutRune(Coord{Y: row, X: col}, hex[i%16], 1)
			i++
		}
	}
}

func TestBackendBlit(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, trm)
	mustStart(t, trm)

	mb := trm.Backend()
	assertF(t, mb != nil, "backend is nil")

	backendBox(t, mb, Coord{X: 0, Y: 0}, Coord{X: 1, Y: 1}, Bold)
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "B")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "C")
	checkContent(t, trm, 1, 1, "D")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 1, 2, "")
	checkContent(t, trm, 2, 2, "")
	checkAttrs(t, trm, 2, 2, Plain)

	// blit the entire box down and right 1
	mb.(Blitter).Blit(Coord{X: 0, Y: 0}, Coord{X: 1, Y: 1}, Coord{X: 2, Y: 2})
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "B")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "C")
	checkContent(t, trm, 1, 1, "A")
	checkContent(t, trm, 2, 1, "B")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 1, 2, "C")
	checkContent(t, trm, 2, 2, "D")

	// spot check attributes
	checkAttrs(t, trm, 2, 2, Bold)
}

func TestBackendBlitReverse(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, trm)
	mustStart(t, trm)

	mb := trm.Backend()
	assertF(t, mb != nil, "backend is nil")

	backendBox(t, mb, Coord{X: 1, Y: 1}, Coord{X: 2, Y: 2}, Bold)
	checkContent(t, trm, 0, 0, "")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 1, 1, "A")
	checkContent(t, trm, 2, 1, "B")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 1, 2, "C")
	checkContent(t, trm, 2, 2, "D")
	checkAttrs(t, trm, 0, 0, Plain)

	// blit the entire box down and right 1
	mb.(Blitter).Blit(Coord{X: 1, Y: 1}, Coord{X: 0, Y: 0}, Coord{X: 2, Y: 2})
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "B")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "C")
	checkContent(t, trm, 1, 1, "D")
	checkContent(t, trm, 2, 1, "B")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 1, 2, "C")
	checkContent(t, trm, 2, 2, "D")

	// spot check attributes
	checkAttrs(t, trm, 2, 2, Bold)
}

func TestGraphemeCluster(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[H")
	// grapheme clustering is off
	writeF(t, trm, "🇨🇭") // flag + regional indicator C + regional indicator H
	writeF(t, trm, "A")

	// we advanced by four columns (two wide emoji), and a single character
	checkPos(t, trm, 5, 0)
	checkContent(t, trm, 0, 0, "\U0001f1e8") // regional indicator C
	checkContent(t, trm, 1, 0, "")           // empty
	checkContent(t, trm, 2, 0, "\U0001f1ed") // regional indicator H
	checkContent(t, trm, 3, 0, "")           // empty
	checkContent(t, trm, 4, 0, "A")          // empty

	// now turn on grapheme clustering
	writeF(t, trm, "\x1b[?2027h")
	writeF(t, trm, "\x1b[H\x1b[J")
	checkPos(t, trm, 0, 0)
	writeF(t, trm, "\U0001f1e8\U0001f1ed")
	checkPos(t, trm, 2, 0)
	writeF(t, trm, "A")
	checkPos(t, trm, 3, 0)
	checkContent(t, trm, 0, 0, "🇨🇭")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "A")

	// lets also verify it works with automargin
	// RECALL: Maximum width is 10
	writeF(t, trm, "\x1b[7h")    // should already be on
	writeF(t, trm, "\x1b[1;10H") // last position in first row
	writeF(t, trm, "🇨🇭A")        // flag + regional indicator C + regional indicator H
	checkPos(t, trm, 1, 1)
	checkContent(t, trm, 9, 0, "🇨🇭")
	checkContent(t, trm, 0, 1, "A")
}

func TestEraseAbove(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b#8")
	writeF(t, trm, "\x1b[3;5H")

	writeF(t, trm, "\x1b[31;1;42m") // set some colors and bold
	writeF(t, trm, "\x1b[1J")
	// cursor is at 4,2
	for row := range Row(5) {
		for col := range Col(10) {
			if row < 2 || row == 2 && col < 5 {
				checkContent(t, trm, col, row, "")
				checkAttrs(t, trm, col, row, Bold)
				checkColors(t, trm, col, row, color.XTerm1, color.XTerm2)
			} else {
				checkContent(t, trm, col, row, "E")
				checkAttrs(t, trm, col, row, Plain)
				checkColors(t, trm, col, row, color.Silver, color.Black)
			}
		}
	}
}

func TestEraseLine(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b#8")

	// erase to end
	writeF(t, trm, "\x1b[31;1;42m") // set some colors and bold
	writeF(t, trm, "\x1b[1;5H")
	writeF(t, trm, "\x1b[0K")
	// cursor is at 4,0
	checkPos(t, trm, 4, 0)
	for col := range Col(10) {
		row := Row(0)
		if col < 4 {
			checkContent(t, trm, col, row, "E")
			checkAttrs(t, trm, col, row, Plain)
			checkColors(t, trm, col, row, color.Silver, color.Black)
		} else {
			checkContent(t, trm, col, row, "")
			checkAttrs(t, trm, col, row, Bold)
			checkColors(t, trm, col, row, color.XTerm1, color.XTerm2)
		}
	}
	checkPos(t, trm, 4, 0)

	// erase to beginning
	writeF(t, trm, "\x1b[2;5H")
	writeF(t, trm, "\x1b[1K")
	// cursor is at 4,1
	checkPos(t, trm, 4, 1)
	for col := range Col(10) {
		row := Row(1)
		if col > 4 {
			checkContent(t, trm, col, row, "E")
			checkAttrs(t, trm, col, row, Plain)
			checkColors(t, trm, col, row, color.Silver, color.Black)
		} else {
			checkContent(t, trm, col, row, "")
			checkAttrs(t, trm, col, row, Bold)
			checkColors(t, trm, col, row, color.XTerm1, color.XTerm2)
		}
	}
	checkPos(t, trm, 4, 1)

	// erase entire line
	writeF(t, trm, "\x1b[3;5H")
	writeF(t, trm, "\x1b[2K")
	// cursor is at 4,2
	checkPos(t, trm, 4, 2)
	for col := range Col(10) {
		row := Row(2)
		checkContent(t, trm, col, row, "")
		checkAttrs(t, trm, col, row, Bold)
		checkColors(t, trm, col, row, color.XTerm1, color.XTerm2)
	}
	checkPos(t, trm, 4, 2)
}

// TestNewLineScroll tests scrolling with a new line.
// This is one of the most fundamental operations for a terminal.
func TestNewLineScroll(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[H\x1b[J") // home and clear
	writeF(t, trm, "\x1b[1;1HA")
	writeF(t, trm, "\x1b[2;1HB")
	writeF(t, trm, "\x1b[5;1HC") // first column on last row
	checkContent(t, trm, 0, 4, "C")
	writeF(t, trm, "\n") // new line should scroll
	checkPos(t, trm, 1, 4)
	checkContent(t, trm, 0, 0, "B")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 0, 4, "")
	checkContent(t, trm, 0, 3, "C")
}

// TestNewLineScrollNoBlitter tests scrolling with a new line,
// using the fallback copy for backends without Blit support.
func TestNewLineScrollNoBlitter(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 10, Y: 5}, MockOptNoBlit{})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[H\x1b[J") // home and clear
	writeF(t, trm, "\x1b[1;1HA")
	writeF(t, trm, "\x1b[2;1HB")
	writeF(t, trm, "\x1b[5;1HC") // first column on last row
	checkContent(t, trm, 0, 4, "C")
	writeF(t, trm, "\n") // new line should scroll
	checkPos(t, trm, 1, 4)
	checkContent(t, trm, 0, 0, "B")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 0, 4, "")
	checkContent(t, trm, 0, 3, "C")
}

// TestScrollUp tests scrolling up. The cursor position is retained.
func TestScrollUp(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[H\x1b[J") // home and clear
	writeF(t, trm, "\x1b[1;1HA")
	writeF(t, trm, "\x1b[2;1HB")
	writeF(t, trm, "\x1b[5;1HC") // first column on last row
	checkContent(t, trm, 0, 4, "C")
	writeF(t, trm, "\x1b[5;5H") // fifth column on last row
	checkPos(t, trm, 4, 4)
	writeF(t, trm, "\x1bD")
	checkPos(t, trm, 4, 4)
	checkContent(t, trm, 0, 0, "B")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 0, 4, "")
	checkContent(t, trm, 0, 3, "C")
	writeF(t, trm, "\x1bE") // this is like a newline
	checkPos(t, trm, 0, 4)
	checkContent(t, trm, 0, 2, "C")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 0, 4, "")

	writeF(t, trm, "\x1b[H\x1bJ")
	writeF(t, trm, "\x1b[1;1HA")
	writeF(t, trm, "\x1b[2;1HB")
	writeF(t, trm, "\x1b[3;1HC")
	writeF(t, trm, "\x1b[4;1HD")
	writeF(t, trm, "\x1b[5;1HE")
	writeF(t, trm, "\x1b[3;3H")
	writeF(t, trm, "\x1b[3S") // scroll in place, leaves cursor where it is
	checkContent(t, trm, 0, 0, "D")
	checkContent(t, trm, 0, 1, "E")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 0, 4, "")
	checkPos(t, trm, 2, 2)
}

// TestScrollDown tests scrolling down. The cursor position is retained.
func TestScrollDown(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\x1b[H\x1b[J") // home and clear
	writeF(t, trm, "\x1b[1;1HA")
	writeF(t, trm, "\x1b[2;1HB")
	writeF(t, trm, "\x1b[4;1HC") // first column on penultimate row
	writeF(t, trm, "\x1b[5;1HD") // first column on last row
	checkContent(t, trm, 0, 3, "C")
	checkContent(t, trm, 0, 4, "D")
	writeF(t, trm, "\x1b[1;5H") // fifth column on first row
	checkPos(t, trm, 4, 0)
	writeF(t, trm, "\x1bM")
	checkPos(t, trm, 4, 0)
	checkContent(t, trm, 0, 0, "")
	checkContent(t, trm, 0, 1, "A")
	checkContent(t, trm, 0, 2, "B")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 0, 4, "C")

	writeF(t, trm, "\x1b[H\x1bJ")
	writeF(t, trm, "\x1b[1;1HA")
	writeF(t, trm, "\x1b[2;1HB")
	writeF(t, trm, "\x1b[3;1HC")
	writeF(t, trm, "\x1b[4;1HD")
	writeF(t, trm, "\x1b[5;1HE")
	writeF(t, trm, "\x1b[3;3H")
	writeF(t, trm, "\x1b[3T") // scroll in place, leaves cursor where it is
	checkContent(t, trm, 0, 0, "")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 0, 3, "A")
	checkContent(t, trm, 0, 4, "B")
	checkPos(t, trm, 2, 2)
}

// TestDECSTBMv1 tests full screen scroll region (test courtesy of ghostty)
func TestDECSTBMv1(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H") // move to top-left
	writeF(t, trm, "\033[0J")   //  clear screen
	writeF(t, trm, "ABC\r\n")
	writeF(t, trm, "DEF\r\n")
	writeF(t, trm, "GHI\r\n")
	writeF(t, trm, "\033[r") // scroll region top/bottom
	writeF(t, trm, "\033[T") // scroll down one

	// |c_______|
	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	checkPos(t, trm, 0, 0)
	checkContent(t, trm, 0, 0, "")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "A")
	checkContent(t, trm, 1, 1, "B")
	checkContent(t, trm, 2, 1, "C")
	checkContent(t, trm, 0, 2, "D")
	checkContent(t, trm, 1, 2, "E")
	checkContent(t, trm, 2, 2, "F")
	checkContent(t, trm, 0, 3, "G")
	checkContent(t, trm, 1, 3, "H")
	checkContent(t, trm, 2, 3, "I")
}

// TestDECSTBMv2 top only
func TestDECSTBMv2(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H") // move to top-left
	writeF(t, trm, "\033[0J")   //  clear screen
	writeF(t, trm, "ABC\r\n")
	writeF(t, trm, "DEF\r\n")
	writeF(t, trm, "GHI\r\n")
	writeF(t, trm, "\033[2;2r") // scroll region top/bottom
	writeF(t, trm, "\033[T")    // scroll down one

	// |________|
	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	checkPos(t, trm, 0, 3) // did not move
	checkContent(t, trm, 0, 0, "")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "A")
	checkContent(t, trm, 1, 1, "B")
	checkContent(t, trm, 2, 1, "C")
	checkContent(t, trm, 0, 2, "D")
	checkContent(t, trm, 1, 2, "E")
	checkContent(t, trm, 2, 2, "F")
	checkContent(t, trm, 0, 3, "G")
	checkContent(t, trm, 1, 3, "H")
	checkContent(t, trm, 2, 3, "I")
}

// TestDECSTBMv3 top and bottom
func TestDECSTBMv3(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H") // move to top-left
	writeF(t, trm, "\033[0J")   //  clear screen
	writeF(t, trm, "ABC\r\n")
	writeF(t, trm, "DEF\r\n")
	writeF(t, trm, "GHI\r\n")
	writeF(t, trm, "\033[1;2r") // scroll region top/bottom
	writeF(t, trm, "\033[T")    // scroll down one

	// |________|
	// |ABC_____|
	// |GHI_____|
	// |________|
	checkPos(t, trm, 0, 0)
	checkContent(t, trm, 0, 0, "")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "A")
	checkContent(t, trm, 1, 1, "B")
	checkContent(t, trm, 2, 1, "C")
	checkContent(t, trm, 0, 2, "G")
	checkContent(t, trm, 1, 2, "H")
	checkContent(t, trm, 2, 2, "I")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
}

// TestDECSTBMv4 top == bottom
func TestDECSTBMv4(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H")
	writeF(t, trm, "\033[0J")
	writeF(t, trm, "ABC\r\n")
	writeF(t, trm, "DEF\r\n")
	writeF(t, trm, "GHI\r\n")
	writeF(t, trm, "\033[2;2r")
	writeF(t, trm, "\033[T")

	// |________|
	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	checkPos(t, trm, 0, 3)
	checkContent(t, trm, 0, 0, "")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "A")
	checkContent(t, trm, 1, 1, "B")
	checkContent(t, trm, 2, 1, "C")
	checkContent(t, trm, 0, 2, "D")
	checkContent(t, trm, 1, 2, "E")
	checkContent(t, trm, 2, 2, "F")
	checkContent(t, trm, 0, 3, "G")
	checkContent(t, trm, 1, 3, "H")
	checkContent(t, trm, 2, 3, "I")
}

// TestRIv1 top of screen, no scroll
func TestRIv1(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H")
	writeF(t, trm, "\033[0J")
	writeF(t, trm, "A\r\n")
	writeF(t, trm, "B\r\n")
	writeF(t, trm, "C\r\n")
	writeF(t, trm, "\033[1;1H")
	writeF(t, trm, "\033M")
	writeF(t, trm, "X")

	// |Xc______|
	// |A_______|
	// |B_______|
	// |C_______|

	checkPos(t, trm, 1, 0)
	checkContent(t, trm, 0, 0, "X")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "A")
	checkContent(t, trm, 1, 1, "")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "B")
	checkContent(t, trm, 1, 2, "")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "C")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
}

// TestRIv2 not top of screen, no scroll
func TestRIv2(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H")
	writeF(t, trm, "\033[0J")
	writeF(t, trm, "A\r\n")
	writeF(t, trm, "B\r\n")
	writeF(t, trm, "C\r\n")
	writeF(t, trm, "\033[2;1H")
	writeF(t, trm, "\033M")
	writeF(t, trm, "X")

	// |Xc______|
	// |B_______|
	// |C_______|
	// |________|

	checkPos(t, trm, 1, 0)
	checkContent(t, trm, 0, 0, "X")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "B")
	checkContent(t, trm, 1, 1, "")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "C")
	checkContent(t, trm, 1, 2, "")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
}

// TestRIv3 scroll region
func TestRIv3(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H") // move to top-left
	writeF(t, trm, "\033[0J")   //  clear screen
	writeF(t, trm, "A\r\n")
	writeF(t, trm, "B\r\n")
	writeF(t, trm, "C\r\n")
	writeF(t, trm, "\033[2;3r")
	writeF(t, trm, "\033[2;1H")
	writeF(t, trm, "\033M")

	// |A_______|
	// |c_______|
	// |B_______|
	// |________|

	checkPos(t, trm, 0, 1)
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 1, 1, "")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "B")
	checkContent(t, trm, 1, 2, "")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
}

// TestRIv4 outside scroll region - goes to top, does not scroll
func TestRIv4(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H") // move to top-left
	writeF(t, trm, "\033[0J")   //  clear screen
	writeF(t, trm, "A\r\n")
	writeF(t, trm, "B\r\n")
	writeF(t, trm, "C\r\n")
	writeF(t, trm, "\033[2;3r")
	writeF(t, trm, "\033[1;1H")
	writeF(t, trm, "\033M")

	// |A_______|
	// |B_______|
	// |C_______|
	// |________|

	checkPos(t, trm, 0, 0)
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "B")
	checkContent(t, trm, 1, 1, "")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "C")
	checkContent(t, trm, 1, 2, "")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
}

// TODO: RIv5 - left right scroll regions (when we implement left/right regions)
// TODO: RIv6 - outside left/right scroll regions (when we implement left/right regions)

// TestINDv1 no scroll region, top of screen
func TestINDv1(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H") // move to top-left
	writeF(t, trm, "\033[0J")   //  clear screen
	writeF(t, trm, "A")
	writeF(t, trm, "\033D")
	writeF(t, trm, "X")

	// |A_______|
	// |_Xc_____|
	// |________|
	// |________|

	checkPos(t, trm, 2, 1)
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 1, 1, "X")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 1, 2, "")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
}

// TestINDv2 no scroll region, bottom of screen
func TestINDv2(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H") // move to top-left
	writeF(t, trm, "\033[0J")   //  clear screen
	writeF(t, trm, "\033[4;1H")
	writeF(t, trm, "A")
	writeF(t, trm, "\033D")
	writeF(t, trm, "X")

	// |________|
	// |________|
	// |A_______|
	// |_Xc_____|

	checkPos(t, trm, 2, 3)
	checkContent(t, trm, 0, 0, "")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 1, 1, "")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "A")
	checkContent(t, trm, 1, 2, "")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "X")
	checkContent(t, trm, 2, 3, "")
}

// TestINDv3 inside scroll region
func TestINDv3(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H") // move to top-left
	writeF(t, trm, "\033[0J")
	writeF(t, trm, "\033[1;3r")
	writeF(t, trm, "A")
	writeF(t, trm, "\033D")
	writeF(t, trm, "X")

	// |A_______|
	// |_Xc_____|
	// |________|
	// |________|

	checkPos(t, trm, 2, 1)
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 1, 1, "X")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 1, 2, "")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
}

// TestINDv4 bottom of scroll region
func TestINDv4(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H") // move to top-left
	writeF(t, trm, "\033[0J")
	writeF(t, trm, "\033[1;3r")
	writeF(t, trm, "\033[4;1H")
	writeF(t, trm, "B")
	writeF(t, trm, "\033[3;1H")
	writeF(t, trm, "A")
	writeF(t, trm, "\033D")
	writeF(t, trm, "X")

	// |________|
	// |A_______|
	// |_Xc_____|
	// |B_______|

	checkPos(t, trm, 2, 2)
	checkContent(t, trm, 0, 0, "")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "A")
	checkContent(t, trm, 1, 1, "")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 1, 2, "X")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "B")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
}

// TestINDv5 bottom of screen with scroll region
func TestINDv5(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H") // move to top-left
	writeF(t, trm, "\033[0J")
	writeF(t, trm, "\033[1;3r")
	writeF(t, trm, "\033[3;1H")
	writeF(t, trm, "A")
	writeF(t, trm, "\033[4;1H")
	writeF(t, trm, "\033D")
	writeF(t, trm, "X")

	// |________|
	// |________|
	// |A_______|
	// |________|
	// |Xc______|

	checkPos(t, trm, 1, 4)
	checkContent(t, trm, 0, 0, "")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 1, 1, "")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "A")
	checkContent(t, trm, 1, 2, "")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
	checkContent(t, trm, 0, 4, "X")
	checkContent(t, trm, 1, 4, "")
	checkContent(t, trm, 2, 4, "")
}

// TODO: INDv6 - outside of left/right scroll region (when we have them)
// TODO: INDv7 - inside of left/right scroll region (when we have them)

// TestCUDv1 - cursor down
func TestCUDv1(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "A")
	writeF(t, trm, "\033[2B")
	writeF(t, trm, "X")

	// |A_______|
	// |________|
	// |_Xc_____|
	// |________|

	checkPos(t, trm, 2, 2)
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 1, 1, "")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 1, 2, "X")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
}

// TestCUDv2 - cursor down above bottom margin
func TestCUDv2(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H")
	writeF(t, trm, "\033[0J")
	writeF(t, trm, "\n\n\n\n")
	writeF(t, trm, "\033[1;3r")
	writeF(t, trm, "A")
	writeF(t, trm, "\033[5B")
	writeF(t, trm, "X")

	// |A_______|
	// |________|
	// |_Xc_____|
	// |________|

	checkPos(t, trm, 2, 2)
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 1, 1, "")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 1, 2, "X")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
}

// TestCUDv3 - cursor down below bottom margin
func TestCUDv3(t *testing.T) {
	trm := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, trm)
	mustStart(t, trm)

	writeF(t, trm, "\033[1;1H")
	writeF(t, trm, "\033[0J")
	writeF(t, trm, "\033[1;3r")
	writeF(t, trm, "A")
	writeF(t, trm, "\033[4;1H")
	writeF(t, trm, "\033[5B")
	writeF(t, trm, "X")

	// |A_______|
	// |________|
	// |________|
	// |________|
	// |Xc______|

	checkPos(t, trm, 1, 4)
	checkContent(t, trm, 0, 0, "A")
	checkContent(t, trm, 1, 0, "")
	checkContent(t, trm, 2, 0, "")
	checkContent(t, trm, 0, 1, "")
	checkContent(t, trm, 1, 1, "")
	checkContent(t, trm, 2, 1, "")
	checkContent(t, trm, 0, 2, "")
	checkContent(t, trm, 1, 2, "")
	checkContent(t, trm, 2, 2, "")
	checkContent(t, trm, 0, 3, "")
	checkContent(t, trm, 1, 3, "")
	checkContent(t, trm, 2, 3, "")
	checkContent(t, trm, 0, 4, "X")
	checkContent(t, trm, 1, 4, "")
	checkContent(t, trm, 2, 4, "")
}

// TODO: Test cases for CUU, CNL, CPL.
