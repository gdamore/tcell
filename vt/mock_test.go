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
func writeF(t *testing.T, term MockTerm, str string, args ...any) {
	t.Helper()
	b := fmt.Appendf(nil, str, args...)
	for len(b) != 0 {
		if n, err := term.Write(b); err != nil {
			t.Fatalf("Failed to write: %v", err)
		} else {
			b = b[n:]
		}
	}
	if err := term.Drain(); err != nil {
		t.Fatalf("Failed to flush: %v", err)
	}
}

func readF(t *testing.T, term MockTerm) string {
	buf := make([]byte, 128)
	n, err := term.Read(buf)
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

func mustClose(t *testing.T, term MockTerm) {
	t.Helper()
	err := term.Close()
	assertF(t, err == nil, "close failed: %v", err)
}

func mustStart(t *testing.T, term MockTerm) {
	t.Helper()
	err := term.Start()
	assertF(t, err == nil, "start failed: %v", err)
}

func checkPos(t *testing.T, term MockTerm, x Col, y Row) {
	t.Helper()
	verifyF(t, term.Pos().X == x && term.Pos().Y == y,
		"bad position %d,%d (expected %d,%d)", term.Pos().X, term.Pos().Y, x, y)
}

func checkContent(t *testing.T, term MockTerm, x Col, y Row, s string) {
	t.Helper()
	if actual := string(term.GetCell(Coord{X: x, Y: y}).C); actual != s {
		t.Errorf("bad content %d,%d (expected %q got %q)", x, y, s, actual)
	}
}

func checkAttrs(t *testing.T, term MockTerm, x Col, y Row, a Attr) {
	t.Helper()
	if actual := term.GetCell(Coord{X: x, Y: y}).S.Attr(); actual != a {
		t.Errorf("bad attr %d,%d (expected %x got %x)", x, y, a, actual)
	}
}

func checkColors(t *testing.T, term MockTerm, x Col, y Row, fg color.Color, bg color.Color) {
	t.Helper()
	if actual := term.GetCell(Coord{X: x, Y: y}).S.Fg(); actual != fg {
		t.Errorf("bad foreground %d,%d (expected %s got %s)", x, y, fg.String(), actual.String())
	}
	if actual := term.GetCell(Coord{X: x, Y: y}).S.Bg(); actual != bg {
		t.Errorf("bad background %d,%d (expected %s got %s)", x, y, bg.String(), actual.String())
	}
}

// TestCursorMove tests several aspects of cursor movement.
func TestCursorMovement(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, term)

	mustStart(t, term)

	if size, err := term.WindowSize(); err != nil {
		t.Fatalf("failed getting window size: %v", err)
	} else if size.Height != 3 || size.Width != 5 {
		t.Fatalf("wrong window size X %d Y %d", size.Width, size.Height)
	}
	writeF(t, term, "\x1b[2;3H")
	checkPos(t, term, 2, 1)

	writeF(t, term, "\x1b[20A") // up 20
	checkPos(t, term, 2, 0)

	writeF(t, term, "\x1b[20B") // down 20
	checkPos(t, term, 2, 2)

	writeF(t, term, "\x1b[A") // up 1
	checkPos(t, term, 2, 1)

	writeF(t, term, "\x1b[2C") // right 2
	checkPos(t, term, 4, 1)

	writeF(t, term, "\x1b[3D") // left 3
	checkPos(t, term, 1, 1)

	writeF(t, term, "\x1b[100D") // left 100
	checkPos(t, term, 0, 1)

	// Now try the next line and previous line
	writeF(t, term, "\x1b[2;3H")
	checkPos(t, term, 2, 1)

	writeF(t, term, "\x1b[1E")
	checkPos(t, term, 0, 2)

	writeF(t, term, "\x1b[2;3H")
	checkPos(t, term, 2, 1)

	writeF(t, term, "\x1b[1F")
	checkPos(t, term, 0, 0)

	writeF(t, term, "\x1b9")
	checkPos(t, term, 1, 0)

	writeF(t, term, "\x1b6")
	checkPos(t, term, 0, 0)
}

// TestDECALN tests the DEC alignment test (screen filled with 'E').
func TestDECALN(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, term)

	mustStart(t, term)

	writeF(t, term, "\x1b#8")

	for y := range Row(3) {
		for x := range Col(5) {
			checkAttrs(t, term, x, y, Plain)
			checkContent(t, term, x, y, "E")
		}
	}

	// clear screen
	writeF(t, term, "\x1b[H\x1b[J")

	for y := range Row(3) {
		for x := range Col(5) {
			checkAttrs(t, term, x, y, Plain)
			checkContent(t, term, x, y, "")
		}
	}

	writeF(t, term, "\x1b[1m\x1b#8") // bold, DECALN
	for y := range Row(3) {
		for x := range Col(5) {
			checkAttrs(t, term, x, y, Plain)
			checkContent(t, term, x, y, "E")
		}
	}
}

// TestBell tests the bell.
func TestBell(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, term)
	mustStart(t, term)

	if term.Bells() != 0 {
		t.Errorf("wrong bell count: %d", term.Bells())
	}
	writeF(t, term, "\x07")
	if term.Bells() != 1 {
		t.Errorf("wrong bell count: %d", term.Bells())
	}
	writeF(t, term, "\x07")
	if term.Bells() != 2 {
		t.Errorf("wrong bell count: %d", term.Bells())
	}
}

// TestPrimaryDA tests primary device attributes using several mechanisms.
func TestPrimaryDA(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, term)

	mustStart(t, term)

	buf := make([]byte, 32)
	writeF(t, term, "\x1b[c")

	n, err := term.Read(buf)
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
	writeF(t, term, "\x1bZ")

	n, err = term.Read(buf)
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
	term := NewMockTerm(MockOptSize{X: 5, Y: 3}, MockOptColors(0))
	defer mustClose(t, term)

	mustStart(t, term)

	buf := make([]byte, 64)
	writeF(t, term, "\x1b[>q")

	n, err := term.Read(buf)
	assertF(t, err == nil, "read failed: %v", err)

	result := string(buf[:n])

	verifyF(t, strings.HasSuffix(result, "\x1b\\"), "missing suffix ST: %q", result)
	verifyF(t, strings.HasPrefix(result, "\x1bP>|"), "missing prefix 'ESC P>|': %q", result)
	verifyF(t, strings.Contains(result, "TCellMock 1.0"), "missing terminal identification: %q", result)
}

// TestCursorReport verifies that cursor position reporting works.
func TestCursorReport(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, term)

	mustStart(t, term)

	writeF(t, term, "\x1b[5;10H") // fifth row, tenth column
	checkPos(t, term, 9, 4)

	writeF(t, term, "\x1b[6n") // cursor position report

	buf := make([]byte, 32)
	n, err := term.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	if result != "\x1b[5;10R" {
		t.Errorf("wrong report: %q", result)
	}

	buf = make([]byte, 32)
	// move the cursor back one
	writeF(t, term, "\b\x1b[6n")
	checkPos(t, term, 8, 4)
	n, err = term.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result = string(buf[:n])
	if result != "\x1b[5;9R" {
		t.Errorf("wrong report: %q", result)
	}
}

// TestAnsiModes tests the private mode feature.
func TestAnsiModes(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, term)

	mustStart(t, term)

	writeF(t, term, "\x1b[20$p") // query for newline mode
	writeF(t, term, "\x1b[20h")  // turn it on
	writeF(t, term, "\x1b[20$p") // should read back on
	writeF(t, term, "\x1b[20l")  // turn it back off
	writeF(t, term, "\x1b[20$p") // should read back off

	buf := make([]byte, 128)
	n, err := term.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	want := "\x1b[20;2$y" + "\x1b[20;1$y" + "\x1b[20;2$y"
	if result != want {
		t.Errorf("wrong response: %q != %q", result, want)
	}
}

// TestPrivateModes tests the private mode feature.
func TestPrivateModes(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, term)

	mustStart(t, term)

	writeF(t, term, "\x1b[?7$p")              // query for auto-margin (should start on by default)
	writeF(t, term, "\x1b[?7l")               // turn it off
	writeF(t, term, "\x1b[?7$p")              // should read back positive
	writeF(t, term, "\x1b[?7h")               // put it back on
	writeF(t, term, "\x1b[?7$p")              // should read back negative
	writeF(t, term, "\x1b[?1919$p")           // read invalid mode
	writeF(t, term, "\x1b[?1919h\x1b[?1919l") // togle invalid mode
	writeF(t, term, "\x1b[?1919$p")           // read invalid mode one more time

	buf := make([]byte, 128)
	n, err := term.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	want := "\x1b[?7;1$y" + "\x1b[?7;2$y" + "\x1b[?7;1$y" + "\x1b[?1919;0$y" + "\x1b[?1919;0$y"
	if result != want {
		t.Errorf("wrong response: %q != %q", result, want)
	}

	// Lets also test the cursor (show vs hide)
	writeF(t, term, "\x1b[?25$p")
	writeF(t, term, "\x1b[?25l")
	writeF(t, term, "\x1b[?25$p")
	writeF(t, term, "\x1b[?25h")
	writeF(t, term, "\x1b[?25$p")

	buf = make([]byte, 128)
	n, err = term.Read(buf)
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
	term := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, term)
	mustStart(t, term)

	// default is auto-margin is enabled
	writeF(t, term, "\x1b[2J") // clear the screen
	writeF(t, term, "\x1b[1;80HAB")
	checkPos(t, term, 1, 1)
	if s := string(term.GetCell(Coord{X: 79, Y: 0}).C); s != "A" {
		t.Errorf("last column wrong: %q", s)
	}
	if s := string(term.GetCell(Coord{X: 0, Y: 1}).C); s != "B" {
		t.Errorf("auto wrap did not work: %q", s)
	}

	// now turn it off
	writeF(t, term, "\x1b[?7l")

	// mess with 3rd row
	writeF(t, term, "\x1b[3;80HCD")
	checkPos(t, term, 79, 2)
	if s := string(term.GetCell(Coord{X: 79, Y: 2}).C); s != "D" {
		t.Errorf("last column wrong: %q", s)
	}

	// turn it back on
	writeF(t, term, "\x1b[?7h")

	// demonstrate that writing to the last column does not advance (pending)
	writeF(t, term, "\x1b[1;80HA")
	checkPos(t, term, 79, 0)

	// but one more character does advance
	writeF(t, term, "\x1b[1;80HAB")
	checkPos(t, term, 1, 1)

	// tab does not advance, but leaves pending state
	writeF(t, term, "\x1b[1;80HA\t")
	checkPos(t, term, 79, 0)
	writeF(t, term, "\x1b[1;80HA\tb")
	checkPos(t, term, 1, 1)

	// up or down movement resets the pending state
	writeF(t, term, "\x1b[1;80HA\x1b[AB")
	checkPos(t, term, 79, 0)
	writeF(t, term, "\x1b[1;80HA\x1b[BB")
	checkPos(t, term, 79, 1)

	// forward also resets pending state (which is clipped)
	writeF(t, term, "\x1b[1;80HA\x1b[CB")
	checkPos(t, term, 79, 0)
	writeF(t, term, "\x1b[1;80HA\x1b[CBC")
	checkPos(t, term, 1, 1)

	// explicit column also resets pending state (which is clipped)
	writeF(t, term, "\x1b[1;80HA\x1b[80GB")
	checkPos(t, term, 79, 0)
	writeF(t, term, "\x1b[1;80HA\x1b[80GBC")
	checkPos(t, term, 1, 1)

	// newline of course as well (and also VF and FF)
	writeF(t, term, "\x1b[1;80HA\nB")
	checkPos(t, term, 79, 1)
	writeF(t, term, "\x1b[1;80HA\fB")
	checkPos(t, term, 79, 1)
	writeF(t, term, "\x1b[1;80HA\vB")
	checkPos(t, term, 79, 1)

	// and also index
	writeF(t, term, "\x1b[1;80HA\x1bDB")
	checkPos(t, term, 79, 1)

	// and also reverse index
	writeF(t, term, "\x1b[2;80HA\x1bMB")
	checkPos(t, term, 79, 0)
}

// TestUnicode tests basic placement of unicode glyphs.
// For now it assumes that the terminal itself supports unicode / latin 1.
func TestUnicode(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b[2J") // clear the screen
	writeF(t, term, "\x1b[2;2H")
	checkPos(t, term, 1, 1)
	writeF(t, term, "åßcπ")
	checkPos(t, term, 5, 1)
	if s := string(term.GetCell(Coord{X: 1, Y: 1}).C); s != "å" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(term.GetCell(Coord{X: 2, Y: 1}).C); s != "ß" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(term.GetCell(Coord{X: 3, Y: 1}).C); s != "c" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(term.GetCell(Coord{X: 4, Y: 1}).C); s != "π" {
		t.Errorf("decode error wrong: %q", s)
	}
}

// TestUnicodeWide tests a wide unicode glyph.
func TestUnicodeWide(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b#8") // fill it with E's (so we can see that wide clears the next cell)
	writeF(t, term, "\x1b[2;2H")
	checkPos(t, term, 1, 1)
	writeF(t, term, "å宽cπ")
	checkPos(t, term, 6, 1)
	if s := string(term.GetCell(Coord{X: 1, Y: 1}).C); s != "å" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(term.GetCell(Coord{X: 2, Y: 1}).C); s != "宽" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(term.GetCell(Coord{X: 3, Y: 1}).C); s != "" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(term.GetCell(Coord{X: 4, Y: 1}).C); s != "c" {
		t.Errorf("decode error wrong: %q", s)
	}
	if s := string(term.GetCell(Coord{X: 5, Y: 1}).C); s != "π" {
		t.Errorf("decode error wrong: %q", s)
	}
}

// TestKeyEventLegacy tests key events when using the default legacy key protocol.
func TestKeyEventLegacy(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(0))
	defer mustClose(t, term)
	mustStart(t, term)

	term.KeyTap(KeyA)
	term.KeyTap(KeyLShift, KeyB)
	term.KeyTap(KeyEnter)
	term.KeyTap(KeyRCtrl, KeyI)
	term.KeyTap(KeyEsc)

	buf := make([]byte, 256)
	n, err := term.Read(buf)
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
	term.KeyTap(KeyF1)                               // SS3 P
	term.KeyTap(KeyLShift, KeyF1)                    // CSI 1 ; 2 P
	term.KeyTap(KeyLCtrl, KeyF2)                     // CSI 1 ; 5 Q
	term.KeyTap(KeyLAlt, KeyLCtrl, KeyLShift, KeyF3) // ESC CSI 1 ; 6 R
	term.KeyTap(KeyRAlt, KeyRCtrl, KeyF4)            // ESC CSI 1 ; 5 S

	want = "\x1bOP"
	want += "\x1b[1;2P"
	want += "\x1b[1;5Q"
	want += "\x1b\x1b[1;6R"
	want += "\x1b\x1b[1;5S"
	n, err = term.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result = string(buf[:n])
	if result != want {
		t.Errorf("key responses failed: %q != %q", result, want)
	}

	// CSI based F-keys
	buf = make([]byte, 256)
	term.KeyTap(KeyF5)                               // CSI 15 ~
	term.KeyTap(KeyRShift, KeyF6)                    // CSI 17 ; 2 ~
	term.KeyTap(KeyLCtrl, KeyF7)                     // CSI 18 ; 5 ~
	term.KeyTap(KeyLAlt, KeyRCtrl, KeyLShift, KeyF8) // ESC CSI 19 ; 6 ~
	term.KeyTap(KeyRAlt, KeyLCtrl, KeyF9)            // ESC CSI 20 ; 5 ~
	term.KeyTap(KeyF20)                              // CSI 34 ~
	term.KeyTap(KeyF15)                              // CSI 28 ~
	term.KeyTap(KeyMenu)                             // CSI 29 ~
	want = "\x1b[15~"
	want += "\x1b[17;2~"
	want += "\x1b[18;5~"
	want += "\x1b\x1b[19;6~"
	want += "\x1b\x1b[20;5~"
	want += "\x1b[34~"
	want += "\x1b[28~"
	want += "\x1b[29~"
	n, err = term.Read(buf)
	assertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	verifyF(t, result == want, "key responses failed: %q != %q", result, want)

	// Misc other keys
	clear(buf)
	term.KeyTap(KeyEnter)                          // \r
	term.KeyTap(KeyTab)                            // \t
	term.KeyTap(KeyLShift, KeyTab)                 // CSI Z
	term.KeyTap(KeyLCtrl, KeyM)                    // \r
	term.KeyTap(KeyLCtrl, KeyL)                    // \f
	term.KeyTap(KeyBackspace)                      // \x7f
	term.KeyTap(KeyRCtrl, KeyBackspace)            // \x08
	term.KeyTap(KeyRCtrl, KeyLShift, KeyBackspace) // none
	term.KeyTap(KeyRCtrl, KeySpace)                // \x00
	term.KeyTap(KeySpace)                          // ' '
	term.KeyTap(KeyRAlt, KeyA)                     // \x1b a
	term.KeyTap(KeyRHyper, KeyA)                   // none
	term.KeyTap(KeyRMeta, KeyA)                    // none
	term.KeyTap(KeyRAlt, KeyLCtrl, KeyJ)           // \x1b\n
	term.KeyTap(KeyRCtrl, KeyL)                    // \x0c
	term.KeyTap(KeyLCtrl, KeyLBrace)               // \x0c

	want = "\r\t\x1b[Z\r\f\x7f\x08\x00 \x1ba\x1b\n\x0c\x1b"
	n, err = term.Read(buf)
	assertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	verifyF(t, result == want, "key responses failed: %q != %q", result, want)

	// Legacy control key mappings (weird ones)
	// Declining a few of the strange ones (control-?)
	clear(buf)
	term.KeyTap(KeyLCtrl, Key8)      // \x7F
	term.KeyTap(KeyLCtrl, Key4)      // \x1c
	term.KeyTap(KeyLCtrl, Key7)      // \x1f
	term.KeyTap(Key7)                // 7
	term.KeyTap(KeyLShift, KeySlash) // ?
	term.KeyTap(KeyRCtrl, KeyLBrace) // \x1b

	want = "\x7f\x1c\x1f7?\x1b"
	n, err = term.Read(buf)
	assertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	verifyF(t, result == want, "key responses failed: %q != %q", result, want)

	// Application cursor keys
	term.KeyTap(KeyUp)
	writeF(t, term, "\x1b[?1h")
	term.KeyTap(KeyDown)
	want = "\x1b[A\x1bOB"
	n, err = term.Read(buf)
	assertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	verifyF(t, result == want, "key responses failed: %q != %q", result, want)
}

// TestSgrAttr tests a variety of combinations of Sgr settings.
func TestSgrAttr(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b[H")
	writeF(t, term, "\x1b[1mA") // bold
	checkAttrs(t, term, 0, 0, Bold)
	checkContent(t, term, 0, 0, "A")

	writeF(t, term, "\x1b[2mB") // dim
	checkAttrs(t, term, 1, 0, Dim)
	checkContent(t, term, 1, 0, "B")

	writeF(t, term, "\x1b[22mC") // clear both
	checkAttrs(t, term, 2, 0, Plain)
	checkContent(t, term, 2, 0, "C")

	writeF(t, term, "\x1b[3;2mD") // italic, dim
	checkAttrs(t, term, 3, 0, Italic|Dim)
	checkContent(t, term, 3, 0, "D")

	writeF(t, term, "\x1b[22mE") // remove dim, should leave italic
	checkAttrs(t, term, 4, 0, Italic)
	checkContent(t, term, 4, 0, "E")

	writeF(t, term, "\x1b[23mF") // clear italic
	checkAttrs(t, term, 5, 0, Plain)
	checkContent(t, term, 5, 0, "F")

	writeF(t, term, "\x1b[3;4mG") // simple underline
	checkAttrs(t, term, 6, 0, Italic|Underline)
	checkContent(t, term, 6, 0, "G")

	writeF(t, term, "\x1b[21mH") // double underline (ECMA)
	checkAttrs(t, term, 7, 0, Italic|DoubleUnderline)
	checkContent(t, term, 7, 0, "H")

	writeF(t, term, "\x1b[4mI") // simple underline
	checkAttrs(t, term, 8, 0, Italic|Underline)
	checkContent(t, term, 8, 0, "I")

	writeF(t, term, "\x1b[4:2mJ") // double underline
	checkAttrs(t, term, 9, 0, Italic|DoubleUnderline)
	checkContent(t, term, 9, 0, "J")

	writeF(t, term, "\x1b[4:3mK") // curly underline
	checkAttrs(t, term, 10, 0, Italic|CurlyUnderline)
	checkContent(t, term, 10, 0, "K")

	writeF(t, term, "\x1b[4:4mL") // dotted underline
	checkAttrs(t, term, 11, 0, Italic|DottedUnderline)
	checkContent(t, term, 11, 0, "L")

	writeF(t, term, "\x1b[4:5mM") // dashed underline
	checkAttrs(t, term, 12, 0, Italic|DashedUnderline)
	checkContent(t, term, 12, 0, "M")

	writeF(t, term, "\x1b[4:9mN") // junk treats as plain underline
	checkAttrs(t, term, 13, 0, Italic|Underline)
	checkContent(t, term, 13, 0, "N")

	writeF(t, term, "\x1b[4:5;24mO") // clustering, clear it
	checkAttrs(t, term, 14, 0, Italic)
	checkContent(t, term, 14, 0, "O")

	writeF(t, term, "\x1b[0;9;7;53mP") // clear, strike-through, reverse, over-lined
	checkAttrs(t, term, 15, 0, StrikeThrough|Reverse|Overline)
	checkContent(t, term, 15, 0, "P")

	writeF(t, term, "\x1b[5;27;29;55mQ")
	checkAttrs(t, term, 16, 0, Blink)
	checkContent(t, term, 16, 0, "Q")

	writeF(t, term, "\x1b[25mR")
	checkAttrs(t, term, 17, 0, Plain)
	checkContent(t, term, 17, 0, "R")
}

// TestSgrColor8 tests simple ECMA 48 ANSI color (only 8 possible color values.)
func TestSgrColor8(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b[36;42m\x1b#8")
	checkColors(t, term, 0, 0, color.Teal, color.Green)

	writeF(t, term, "\x1b[H\x1b[39mA")
	checkColors(t, term, 0, 0, color.Silver, color.Green)

	writeF(t, term, "\x1b[49mA")
	checkColors(t, term, 1, 0, color.Silver, color.Black)

	// verify zero clears colors, first write some non zero colors
	writeF(t, term, "\x1b[36;42mD")
	checkColors(t, term, 2, 0, color.Teal, color.Green)

	// then send zero, should go to default colors
	writeF(t, term, "\x1b[0mA")
	checkColors(t, term, 3, 0, color.Silver, color.Black)
}

// TestSgrColor256 tests simple ECMA 48 ANSI color (256 possible color values.)
func TestSgrColor256(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(256))
	defer mustClose(t, term)
	mustStart(t, term)

	// foreground
	writeF(t, term, "\x1b[38:5:6m\x1b[42m\x1b#8")
	checkColors(t, term, 0, 0, color.Teal, color.Green)
	writeF(t, term, "\x1b[38;5;5m\x1b[42mA")
	checkColors(t, term, 0, 0, color.XTerm5, color.Green)

	writeF(t, term, "\x1b[38;5;212;1mB")
	checkColors(t, term, 1, 0, color.XTerm212, color.Green)
	checkAttrs(t, term, 1, 0, Bold)

	// background
	writeF(t, term, "\x1b[48;5;2mA")
	checkColors(t, term, 2, 0, color.XTerm212, color.XTerm2)
	checkAttrs(t, term, 2, 0, Bold)

	writeF(t, term, "\x1b[48:5:134;2mC")
	checkColors(t, term, 3, 0, color.XTerm212, color.XTerm134)
	checkAttrs(t, term, 3, 0, Dim)

	// mix background and foreground using colons
	writeF(t, term, "\x1b[48:5:135;38:5:22mC")
	checkColors(t, term, 4, 0, color.XTerm22, color.XTerm135)

	// and using semicolons
	writeF(t, term, "\x1b[48;5;136;38;5;23mC")
	checkColors(t, term, 5, 0, color.XTerm23, color.XTerm136)

	// underline colors - it uses the same parser so we won't check
	// all the variations
	writeF(t, term, "\x1b[58;5;21;4mC")
	verifyF(t, term.GetCell(Coord{X: 6, Y: 0}).S.Uc() == color.XTerm21, "underline color is wrong")

	writeF(t, term, "\x1b[H\x1b[39;49mA")
	checkColors(t, term, 0, 0, color.Silver, color.Black)

	// verify zero clears colors, first write some non zero colors
	writeF(t, term, "\x1b[H\x1b[36;42mA")
	checkColors(t, term, 0, 0, color.Teal, color.Green)

	// then send zero, should go to default colors
	writeF(t, term, "\x1b[H\x1b[0mA")
	checkColors(t, term, 0, 0, color.Silver, color.Black)

	// fuzz some things
	writeF(t, term, "\x1b[m\x1b[H")
	writeF(t, term, "\x1b[38:3m")
	writeF(t, term, "\x1b[38;3m")
	writeF(t, term, "\x1b[38;2m")
	writeF(t, term, "\x1b[38;2:300m")
	writeF(t, term, "\x1b[38;5m")
	writeF(t, term, "\x1b[38;5:300m")
	writeF(t, term, "\x1b[38:2m")
	writeF(t, term, "\x1b[38:5m")
	writeF(t, term, "\x1b[38:2;1;1;1m")
	writeF(t, term, "\x1b[38:2:1:1m")
	writeF(t, term, "A")
	checkColors(t, term, 0, 0, color.Silver, color.Black)
}

// TestSgrColorRGB tests simple ECMA 48 ANSI color (full 24-bit color.)
func TestSgrColorRGB(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24}, MockOptColors(1<<24))
	defer mustClose(t, term)
	mustStart(t, term)

	// foreground
	writeF(t, term, "\x1b[m\x1b[H\x1b[38:2:255:0:0mA")
	checkColors(t, term, 0, 0, color.NewRGBColor(255, 0, 0), color.Black)

	writeF(t, term, "\x1b[m\x1b[H\x1b[38;2;2;0;0mA")
	checkColors(t, term, 0, 0, color.NewRGBColor(2, 0, 0), color.Black)

	// background
	writeF(t, term, "\x1b[m\x1b[H\x1b[48:2:1:2:3mA")
	checkColors(t, term, 0, 0, color.Silver, color.NewRGBColor(1, 2, 3))

	writeF(t, term, "\x1b[m\x1b[H\x1b[48;2;4;5;6mA")
	checkColors(t, term, 0, 0, color.Silver, color.NewRGBColor(4, 5, 6))

	// full colors
	writeF(t, term, "\x1b[m\x1b[H\x1b[38;2;99;88;77;48;2;4;5;6;58;2;99;98;91;1mA")
	checkColors(t, term, 0, 0, color.NewRGBColor(99, 88, 77), color.NewRGBColor(4, 5, 6))
	verifyF(t, term.GetCell(Coord{X: 0, Y: 0}).S.Uc() == color.NewRGBColor(99, 98, 91), "underline color is wrong")
	checkAttrs(t, term, 0, 0, Bold)
}

// TestTitles tests that we can set a window title.
func TestTitles(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, term)
	mustStart(t, term)
	writeF(t, term, "\x1b]2;Test Application\x1b\\")
	if s := term.GetTitle(); s != "Test Application" {
		t.Errorf("wrong title: %q", s)
	}

	// test ST termination using legacy bell character
	writeF(t, term, "\x1b]2;Bell Ring\x07")
	if s := term.GetTitle(); s != "Bell Ring" {
		t.Errorf("wrong title: %q", s)
	}

	// try using 8-bit sequence
	writeF(t, term, "\x9d2;Eight Bits\x9c")
	if s := term.GetTitle(); s != "Eight Bits" {
		t.Errorf("wrong title: %q", s)
	}
}

// TestResize tests resizing the terminal
func TestResize(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, term)
	mustStart(t, term)

	// with E, and enable notifications
	writeF(t, term, "\x1b#8\x1b[?2048h")
	resizeQ := make(chan bool, 1)
	term.NotifyResize(resizeQ)

	term.SetSize(Coord{X: 132, Y: 24})
	if sz, err := term.WindowSize(); err != nil || sz.Height != 24 || sz.Width != 132 {
		t.Errorf("resize did not occur: %v %d %d", err, sz.Height, sz.Width)
	}
	for y := range Row(24) {
		for x := range Col(80) {
			if s := string(term.GetCell(Coord{X: x, Y: y}).C); s != "E" {
				t.Errorf("resize content at %d,%d wrong: %q", x, y, s)
			}
		}
	}
	for y := range Row(24) {
		for x := Col(80); x < 132; x++ {
			if s := string(term.GetCell(Coord{X: x, Y: y}).C); s != "" {
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
	term.Write([]byte{})

	buf := make([]byte, 128)
	n, err := term.Read(buf)
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
	term := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "a\tC")
	if s := string(term.GetCell(Coord{X: 8, Y: 0}).C); s != "C" {
		t.Errorf("tab did not work: %q", s)
	}
	writeF(t, term, "\x1b[3I")
	if x := term.Pos().X; x != 32 {
		t.Errorf("wrong position %d", x)
	}

	writeF(t, term, "\x1b[2Z")
	if x := term.Pos().X; x != 16 {
		t.Errorf("wrong position %d", x)
	}

	writeF(t, term, "\x1b[3g") // clear all tabs
	writeF(t, term, "\x1b[I")  // one tab, should go to right margin
	if x := term.Pos().X; x != 79 {
		t.Errorf("wrong position: %d", x)
	}

	writeF(t, term, "\x1b[Z")
	if x := term.Pos().X; x != 0 {
		t.Errorf("wrong position: %d", x)
	}

	// reset tabs
	writeF(t, term, "\x1b[?5W")

	writeF(t, term, "\t")
	if x := term.Pos().X; x != 8 {
		t.Errorf("wrong position: %d", x)
	}
	// clear this position, advance one
	writeF(t, term, "\x1b[gA")
	if x := term.Pos().X; x != 9 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, term, "\x1bH")
	writeF(t, term, "\t")
	if x := term.Pos().X; x != 16 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, term, "\x1b[Z")
	if x := term.Pos().X; x != 9 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, term, "\x1b[Z")
	if x := term.Pos().X; x != 0 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, term, "\x1b[1;10H") // goto position 9
	if x := term.Pos().X; x != 9 {
		t.Errorf("wrong position: %d", x)
	}
	// delete this one (do it twice to exercise the does not exist flow)
	writeF(t, term, "\x1b[0g")
	writeF(t, term, "\x1b[0g")

	// advance to next tab, then back, we should go to 0
	writeF(t, term, "\t")
	if x := term.Pos().X; x != 16 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, term, "\x1b[Z")
	if x := term.Pos().X; x != 0 {
		t.Errorf("wrong position: %d", x)
	}
	writeF(t, term, "\x1b[20I")
	writeF(t, term, "\t\t")
	if pos := term.Pos(); pos.X != 79 || pos.Y != 0 {
		t.Errorf("wrong position: %d %d", pos.X, pos.Y)
	}

	// now backwards
	writeF(t, term, "\x1b[20Z")
	writeF(t, term, "\x1b[Z")
	if pos := term.Pos(); pos.X != 0 || pos.Y != 0 {
		t.Errorf("wrong position: %d %d", pos.X, pos.Y)
	}
}

func TestVerticalPos(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b[2;2H")
	writeF(t, term, "\x1b[10d")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 9 {
		t.Errorf("wrong position: %d %d", pos.X, pos.Y)
	}
	writeF(t, term, "\x1b[2e")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 11 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, term, "\x1b[e")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 12 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, term, "\x1b[0e")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 13 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, term, "\x1b[50d")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 23 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, term, "\x1b[50e")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 23 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, term, "\x1b[0d")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 0 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, term, "\x1b[10d")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 9 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	writeF(t, term, "\x1b[1d")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 0 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
}

func TestSaveCursorPosition(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b[1;1H\x1b[J")
	writeF(t, term, "\x1b[1;5H")
	writeF(t, term, "A")
	writeF(t, term, "\x1b7")
	writeF(t, term, "\x1b[1;1H")
	writeF(t, term, "B")
	writeF(t, term, "\x1b8")
	writeF(t, term, "X")

	checkContent(t, term, 0, 0, "B")
	checkContent(t, term, 4, 0, "A")
	checkContent(t, term, 5, 0, "X")
}

func TestSaveCursorWrapState(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b[1;1H\x1b[J")
	writeF(t, term, "\x1b[80G")
	writeF(t, term, "A")
	writeF(t, term, "\x1b7") // save cursor
	writeF(t, term, "\x1b[1;1H")
	writeF(t, term, "B")
	writeF(t, term, "\x1b8") // restore cursor
	writeF(t, term, "X")
	checkContent(t, term, 0, 0, "B")
	checkContent(t, term, 79, 0, "A")
	checkContent(t, term, 0, 1, "X")
}

func TestSaveCursorSgr(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, term)
	mustStart(t, term)
	writeF(t, term, "\x1b[1;1H\x1b[0J")
	writeF(t, term, "\x1b[1;4;33;44m")
	writeF(t, term, "A")
	checkPos(t, term, 1, 0)
	writeF(t, term, "\x1b7")
	checkPos(t, term, 1, 0)
	writeF(t, term, "\x1b[0m")
	checkPos(t, term, 1, 0)
	writeF(t, term, "BE")
	checkPos(t, term, 3, 0)
	writeF(t, term, "\x1b8")
	checkPos(t, term, 1, 0)
	writeF(t, term, "X")
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "X")
	checkContent(t, term, 2, 0, "E")
	checkAttrs(t, term, 0, 0, Bold|Underline)
	checkAttrs(t, term, 1, 0, Bold|Underline)
	checkAttrs(t, term, 2, 0, Plain)
	checkColors(t, term, 0, 0, color.XTerm3, color.XTerm4)
	checkColors(t, term, 1, 0, color.XTerm3, color.XTerm4)
	checkColors(t, term, 2, 0, color.Silver, color.Black)
}

func TestReset(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 80, Y: 24})
	defer mustClose(t, term)
	mustStart(t, term)

	// write a bunch of stuff to create state (so we can verify it gets reset)
	writeF(t, term, "\x1b[1;4;33;44m")
	writeF(t, term, "\x1b#8")
	writeF(t, term, "\x1b[1;80HX")
	writeF(t, term, "\x1b7")    // save cursor
	writeF(t, term, "\x1b[?7l") // disable automargin
	writeF(t, term, "\x1bc")
	checkPos(t, term, 0, 0)
	for row := range Row(24) {
		for col := range Col(80) {
			checkAttrs(t, term, col, row, Plain)
			checkContent(t, term, col, row, "")
			checkColors(t, term, col, row, color.Silver, color.Black)
		}
	}
	writeF(t, term, "\x1b8") // restore cursor
	checkPos(t, term, 0, 0)
	writeF(t, term, "X")
	checkAttrs(t, term, 0, 0, Plain)
	checkContent(t, term, 0, 0, "X")
	checkColors(t, term, 0, 0, color.Silver, color.Black)
	writeF(t, term, "\x1b[?7$p")

	// verify mode reset
	want := "\x1b[?7;1$y"
	buf := make([]byte, 128)
	n, err := term.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	verifyF(t, result == want, "wrong mode: %q != %q", result, want)
}

// backendBox makes a backend box, filled with increasing letters (modulo 16)
func backendBox(t *testing.T, mb MockBackend, tl Coord, br Coord, attr Attr) {
	t.Helper()
	style := BaseStyle.WithAttr(attr)
	hex := []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P'}
	assertF(t, len(hex) == 16, "wrong hex string")
	i := 0
	for row := tl.Y; row <= br.Y; row++ {
		for col := tl.X; col <= br.X; col++ {
			mb.Put(Coord{Y: row, X: col}, Cell{C: string(hex[i%16]), W: 1, S: style})
			i++
		}
	}
}

func TestBackendBlit(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	mb := term.Backend()
	assertF(t, mb != nil, "backend is nil")

	backendBox(t, mb, Coord{X: 0, Y: 0}, Coord{X: 1, Y: 1}, Bold)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "C")
	checkContent(t, term, 1, 1, "D")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkAttrs(t, term, 2, 2, Plain)

	// blit the entire box down and right 1
	mb.(Blitter).Blit(Coord{X: 0, Y: 0}, Coord{X: 1, Y: 1}, Coord{X: 2, Y: 2})
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "C")
	checkContent(t, term, 1, 1, "A")
	checkContent(t, term, 2, 1, "B")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "C")
	checkContent(t, term, 2, 2, "D")

	// spot check attributes
	checkAttrs(t, term, 2, 2, Bold)
}

func TestBackendBlitReverse(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	mb := term.Backend()
	assertF(t, mb != nil, "backend is nil")

	backendBox(t, mb, Coord{X: 1, Y: 1}, Coord{X: 2, Y: 2}, Bold)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "A")
	checkContent(t, term, 2, 1, "B")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "C")
	checkContent(t, term, 2, 2, "D")
	checkAttrs(t, term, 0, 0, Plain)

	// blit the entire box down and right 1
	mb.(Blitter).Blit(Coord{X: 1, Y: 1}, Coord{X: 0, Y: 0}, Coord{X: 2, Y: 2})
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "C")
	checkContent(t, term, 1, 1, "D")
	checkContent(t, term, 2, 1, "B")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "C")
	checkContent(t, term, 2, 2, "D")

	// spot check attributes
	checkAttrs(t, term, 2, 2, Bold)
}

func TestGraphemeCluster(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b[H")
	// grapheme clustering is off
	writeF(t, term, "🇨🇭") // flag + regional indicator C + regional indicator H
	writeF(t, term, "A")

	// we advanced by four columns (two wide emoji), and a single character
	checkPos(t, term, 5, 0)
	checkContent(t, term, 0, 0, "\U0001f1e8") // regional indicator C
	checkContent(t, term, 1, 0, "")           // empty
	checkContent(t, term, 2, 0, "\U0001f1ed") // regional indicator H
	checkContent(t, term, 3, 0, "")           // empty
	checkContent(t, term, 4, 0, "A")          // empty

	// now turn on grapheme clustering
	writeF(t, term, "\x1b[?2027h")
	writeF(t, term, "\x1b[H\x1b[J")
	checkPos(t, term, 0, 0)
	writeF(t, term, "\U0001f1e8\U0001f1ed")
	checkPos(t, term, 2, 0)
	writeF(t, term, "A")
	checkPos(t, term, 3, 0)
	checkContent(t, term, 0, 0, "🇨🇭")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "A")

	// lets also verify it works with automargin
	// RECALL: Maximum width is 10
	writeF(t, term, "\x1b[7h")    // should already be on
	writeF(t, term, "\x1b[1;10H") // last position in first row
	writeF(t, term, "🇨🇭A")        // flag + regional indicator C + regional indicator H
	checkPos(t, term, 1, 1)
	checkContent(t, term, 9, 0, "🇨🇭")
	checkContent(t, term, 0, 1, "A")
}

func TestEraseAbove(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b#8")
	writeF(t, term, "\x1b[3;5H")

	writeF(t, term, "\x1b[31;1;42m") // set some colors and bold
	writeF(t, term, "\x1b[1J")
	// cursor is at 4,2
	for row := range Row(5) {
		for col := range Col(10) {
			if row < 2 || row == 2 && col < 5 {
				checkContent(t, term, col, row, "")
				checkAttrs(t, term, col, row, Plain)
				checkColors(t, term, col, row, color.XTerm1, color.XTerm2)
			} else {
				checkContent(t, term, col, row, "E")
				checkAttrs(t, term, col, row, Plain)
				checkColors(t, term, col, row, color.Silver, color.Black)
			}
		}
	}
}

func TestEraseLine(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b#8")

	// erase to end
	writeF(t, term, "\x1b[31;1;42m") // set some colors and bold
	writeF(t, term, "\x1b[1;5H")
	writeF(t, term, "\x1b[0K")
	// cursor is at 4,0
	checkPos(t, term, 4, 0)
	for col := range Col(10) {
		row := Row(0)
		if col < 4 {
			checkContent(t, term, col, row, "E")
			checkAttrs(t, term, col, row, Plain)
			checkColors(t, term, col, row, color.Silver, color.Black)
		} else {
			checkContent(t, term, col, row, "")
			checkAttrs(t, term, col, row, Plain)
			checkColors(t, term, col, row, color.XTerm1, color.XTerm2)
		}
	}
	checkPos(t, term, 4, 0)

	// erase to beginning
	writeF(t, term, "\x1b[2;5H")
	writeF(t, term, "\x1b[1K")
	// cursor is at 4,1
	checkPos(t, term, 4, 1)
	for col := range Col(10) {
		row := Row(1)
		if col > 4 {
			checkContent(t, term, col, row, "E")
			checkAttrs(t, term, col, row, Plain)
			checkColors(t, term, col, row, color.Silver, color.Black)
		} else {
			checkContent(t, term, col, row, "")
			checkAttrs(t, term, col, row, Plain)
			checkColors(t, term, col, row, color.XTerm1, color.XTerm2)
		}
	}
	checkPos(t, term, 4, 1)

	// erase entire line
	writeF(t, term, "\x1b[3;5H")
	writeF(t, term, "\x1b[2K")
	// cursor is at 4,2
	checkPos(t, term, 4, 2)
	for col := range Col(10) {
		row := Row(2)
		checkContent(t, term, col, row, "")
		checkAttrs(t, term, col, row, Plain)
		checkColors(t, term, col, row, color.XTerm1, color.XTerm2)
	}
	checkPos(t, term, 4, 2)
}

// TestNewLineScroll tests scrolling with a new line.
// This is one of the most fundamental operations for a terminal.
func TestNewLineScroll(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 10, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b[H\x1b[J") // home and clear
	writeF(t, term, "\x1b[1;1HA")
	writeF(t, term, "\x1b[2;1HB")
	writeF(t, term, "\x1b[5;1HC") // first column on last row
	checkContent(t, term, 0, 4, "C")
	writeF(t, term, "\n") // new line should scroll
	checkPos(t, term, 1, 4)
	checkContent(t, term, 0, 0, "B")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 0, 4, "")
	checkContent(t, term, 0, 3, "C")
}

// TestNewLineScrollNoBlitter tests scrolling with a new line,
// using the fallback copy for backends without Blit support.
func TestNewLineScrollNoBlitter(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 10, Y: 5}, MockOptNoBlit{})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\x1b[H\x1b[J") // home and clear
	writeF(t, term, "\x1b[1;1HA")
	writeF(t, term, "\x1b[2;1HB")
	writeF(t, term, "\x1b[5;1HC") // first column on last row
	checkContent(t, term, 0, 4, "C")
	writeF(t, term, "\n") // new line should scroll
	checkPos(t, term, 1, 4)
	checkContent(t, term, 0, 0, "B")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 0, 4, "")
	checkContent(t, term, 0, 3, "C")
}

func TestNewLineModes(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 10, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)
	writeF(t, term, "\x1b[H\x1b[J")
	writeF(t, term, "ABC\n")
	checkPos(t, term, 3, 1)
	term.KeyTap(KeyEnter)
	writeF(t, term, "\x1b[20h")
	term.KeyTap(KeyEnter)
	writeF(t, term, "DEF")
	checkPos(t, term, 6, 1)
	writeF(t, term, "\n")
	checkPos(t, term, 0, 2)
	writeF(t, term, "GHI")
	checkPos(t, term, 3, 2)
	writeF(t, term, "\x1b[20l")
	writeF(t, term, "\n")
	checkPos(t, term, 3, 3)
	term.KeyTap(KeyEnter)

	// |ABC_____|
	// |___DEF__|
	// |GHI_____|
	// |___c____|
	//
	// input stream contains \r\r\n\r

	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 4, 0, "")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 3, 1, "D")
	checkContent(t, term, 4, 1, "E")
	checkContent(t, term, 5, 1, "F")
	checkContent(t, term, 0, 2, "G")
	checkContent(t, term, 1, 2, "H")
	checkContent(t, term, 2, 2, "I")
	checkContent(t, term, 3, 2, "")
	checkContent(t, term, 4, 2, "")
	checkContent(t, term, 5, 2, "")

	result := readF(t, term)
	want := "\r\r\n\r"
	verifyF(t, result == want, "response incorrect: %q != %q", result, want)
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
