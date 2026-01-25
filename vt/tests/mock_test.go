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

package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/vt"
)

// TestCursorMove tests several aspects of cursor movement.
func TestCursorMovement(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 5, Y: 3}, vt.MockOptColors(0))
	defer MustClose(t, term)

	MustStart(t, term)

	if size, err := term.WindowSize(); err != nil {
		t.Fatalf("failed getting window size: %v", err)
	} else if size.Height != 3 || size.Width != 5 {
		t.Fatalf("wrong window size X %d Y %d", size.Width, size.Height)
	}
	WriteF(t, term, "\x1b[2;3H")
	CheckPos(t, term, 2, 1)

	WriteF(t, term, "\x1b[20A") // up 20
	CheckPos(t, term, 2, 0)

	WriteF(t, term, "\x1b[20B") // down 20
	CheckPos(t, term, 2, 2)

	WriteF(t, term, "\x1b[A") // up 1
	CheckPos(t, term, 2, 1)

	WriteF(t, term, "\x1b[2C") // right 2
	CheckPos(t, term, 4, 1)

	WriteF(t, term, "\x1b[3D") // left 3
	CheckPos(t, term, 1, 1)

	WriteF(t, term, "\x1b[100D") // left 100
	CheckPos(t, term, 0, 1)

	// Now try the next line and previous line
	WriteF(t, term, "\x1b[2;3H")
	CheckPos(t, term, 2, 1)

	WriteF(t, term, "\x1b[1E")
	CheckPos(t, term, 0, 2)

	WriteF(t, term, "\x1b[2;3H")
	CheckPos(t, term, 2, 1)

	WriteF(t, term, "\x1b[1F")
	CheckPos(t, term, 0, 0)

	WriteF(t, term, "\x1b9")
	CheckPos(t, term, 1, 0)

	WriteF(t, term, "\x1b6")
	CheckPos(t, term, 0, 0)
}

// TestDECALN tests the DEC alignment test (screen filled with 'E').
func TestDECALN(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 5, Y: 3}, vt.MockOptColors(0))
	defer MustClose(t, term)

	MustStart(t, term)

	WriteF(t, term, "\x1b#8")

	for y := range Row(3) {
		for x := range Col(5) {
			CheckAttrs(t, term, x, y, vt.Plain)
			CheckContent(t, term, x, y, "E")
		}
	}

	// clear screen
	WriteF(t, term, "\x1b[H\x1b[J")

	for y := range Row(3) {
		for x := range Col(5) {
			CheckAttrs(t, term, x, y, vt.Plain)
			CheckContent(t, term, x, y, "")
		}
	}

	WriteF(t, term, "\x1b[1m\x1b#8") // bold, DECALN
	for y := range Row(3) {
		for x := range Col(5) {
			CheckAttrs(t, term, x, y, vt.Plain)
			CheckContent(t, term, x, y, "E")
		}
	}
}

// TestBell tests the bell.
func TestBell(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 5, Y: 3}, vt.MockOptColors(0))
	defer MustClose(t, term)
	MustStart(t, term)

	if term.Bells() != 0 {
		t.Errorf("wrong bell count: %d", term.Bells())
	}
	WriteF(t, term, "\x07")
	if term.Bells() != 1 {
		t.Errorf("wrong bell count: %d", term.Bells())
	}
	WriteF(t, term, "\x07")
	if term.Bells() != 2 {
		t.Errorf("wrong bell count: %d", term.Bells())
	}
}

// TestPrimaryDA tests primary device attributes using several mechanisms.
func TestPrimaryDA(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 5, Y: 3}, vt.MockOptColors(0))
	defer MustClose(t, term)

	MustStart(t, term)

	buf := make([]byte, 32)
	WriteF(t, term, "\x1b[c")

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
	WriteF(t, term, "\x1bZ")

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
	term := vt.NewMockTerm(vt.MockOptSize{X: 5, Y: 3}, vt.MockOptColors(0))
	defer MustClose(t, term)

	MustStart(t, term)

	buf := make([]byte, 64)
	WriteF(t, term, "\x1b[>q")

	n, err := term.Read(buf)
	AssertF(t, err == nil, "read failed: %v", err)

	result := string(buf[:n])

	VerifyF(t, strings.HasSuffix(result, "\x1b\\"), "missing suffix ST: %q", result)
	VerifyF(t, strings.HasPrefix(result, "\x1bP>|"), "missing prefix 'ESC P>|': %q", result)
	VerifyF(t, strings.Contains(result, "TCellMock 1.0"), "missing terminal identification: %q", result)
}

// TestCursorReport verifies that cursor position reporting works.
func TestCursorReport(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(0))
	defer MustClose(t, term)

	MustStart(t, term)

	WriteF(t, term, "\x1b[5;10H") // fifth row, tenth column
	CheckPos(t, term, 9, 4)

	WriteF(t, term, "\x1b[6n") // cursor position report

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
	WriteF(t, term, "\b\x1b[6n")
	CheckPos(t, term, 8, 4)
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
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(0))
	defer MustClose(t, term)

	MustStart(t, term)

	WriteF(t, term, "\x1b[20$p") // query for newline mode
	WriteF(t, term, "\x1b[20h")  // turn it on
	WriteF(t, term, "\x1b[20$p") // should read back on
	WriteF(t, term, "\x1b[20l")  // turn it back off
	WriteF(t, term, "\x1b[20$p") // should read back off

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
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(0))
	defer MustClose(t, term)

	MustStart(t, term)

	WriteF(t, term, "\x1b[?7$p")              // query for auto-margin (should start on by default)
	WriteF(t, term, "\x1b[?7l")               // turn it off
	WriteF(t, term, "\x1b[?7$p")              // should read back positive
	WriteF(t, term, "\x1b[?7h")               // put it back on
	WriteF(t, term, "\x1b[?7$p")              // should read back negative
	WriteF(t, term, "\x1b[?1919$p")           // read invalid mode
	WriteF(t, term, "\x1b[?1919h\x1b[?1919l") // togle invalid mode
	WriteF(t, term, "\x1b[?1919$p")           // read invalid mode one more time

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
	WriteF(t, term, "\x1b[?25$p")
	WriteF(t, term, "\x1b[?25l")
	WriteF(t, term, "\x1b[?25$p")
	WriteF(t, term, "\x1b[?25h")
	WriteF(t, term, "\x1b[?25$p")

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
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(0))
	defer MustClose(t, term)
	MustStart(t, term)

	// default is auto-margin is enabled
	WriteF(t, term, "\x1b[2J") // clear the screen
	WriteF(t, term, "\x1b[1;80HAB")
	CheckPos(t, term, 1, 1)
	if s := string(term.GetCell(Coord{X: 79, Y: 0}).C); s != "A" {
		t.Errorf("last column wrong: %q", s)
	}
	if s := string(term.GetCell(Coord{X: 0, Y: 1}).C); s != "B" {
		t.Errorf("auto wrap did not work: %q", s)
	}

	// now turn it off
	WriteF(t, term, "\x1b[?7l")

	// mess with 3rd row
	WriteF(t, term, "\x1b[3;80HCD")
	CheckPos(t, term, 79, 2)
	if s := string(term.GetCell(Coord{X: 79, Y: 2}).C); s != "D" {
		t.Errorf("last column wrong: %q", s)
	}

	// turn it back on
	WriteF(t, term, "\x1b[?7h")

	// demonstrate that writing to the last column does not advance (pending)
	WriteF(t, term, "\x1b[1;80HA")
	CheckPos(t, term, 79, 0)

	// but one more character does advance
	WriteF(t, term, "\x1b[1;80HAB")
	CheckPos(t, term, 1, 1)

	// tab does not advance, but leaves pending state
	WriteF(t, term, "\x1b[1;80HA\t")
	CheckPos(t, term, 79, 0)
	WriteF(t, term, "\x1b[1;80HA\tb")
	CheckPos(t, term, 1, 1)

	// up or down movement resets the pending state
	WriteF(t, term, "\x1b[1;80HA\x1b[AB")
	CheckPos(t, term, 79, 0)
	WriteF(t, term, "\x1b[1;80HA\x1b[BB")
	CheckPos(t, term, 79, 1)

	// forward also resets pending state (which is clipped)
	WriteF(t, term, "\x1b[1;80HA\x1b[CB")
	CheckPos(t, term, 79, 0)
	WriteF(t, term, "\x1b[1;80HA\x1b[CBC")
	CheckPos(t, term, 1, 1)

	// explicit column also resets pending state (which is clipped)
	WriteF(t, term, "\x1b[1;80HA\x1b[80GB")
	CheckPos(t, term, 79, 0)
	WriteF(t, term, "\x1b[1;80HA\x1b[80GBC")
	CheckPos(t, term, 1, 1)

	// newline of course as well (and also VF and FF)
	WriteF(t, term, "\x1b[1;80HA\nB")
	CheckPos(t, term, 79, 1)
	WriteF(t, term, "\x1b[1;80HA\fB")
	CheckPos(t, term, 79, 1)
	WriteF(t, term, "\x1b[1;80HA\vB")
	CheckPos(t, term, 79, 1)

	// and also index
	WriteF(t, term, "\x1b[1;80HA\x1bDB")
	CheckPos(t, term, 79, 1)

	// and also reverse index
	WriteF(t, term, "\x1b[2;80HA\x1bMB")
	CheckPos(t, term, 79, 0)
}

// TestUnicode tests basic placement of unicode glyphs.
// For now it assumes that the terminal itself supports unicode / latin 1.
func TestUnicode(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(0))
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[2J") // clear the screen
	WriteF(t, term, "\x1b[2;2H")
	CheckPos(t, term, 1, 1)
	WriteF(t, term, "åßcπ")
	CheckPos(t, term, 5, 1)
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
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(0))
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b#8") // fill it with E's (so we can see that wide clears the next cell)
	WriteF(t, term, "\x1b[2;2H")
	CheckPos(t, term, 1, 1)
	WriteF(t, term, "å宽cπ")
	CheckPos(t, term, 6, 1)
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
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(0))
	defer MustClose(t, term)
	MustStart(t, term)

	term.KeyTap(vt.KeyA)
	term.KeyTap(vt.KeyLShift, vt.KeyB)
	term.KeyTap(vt.KeyEnter)
	term.KeyTap(vt.KeyRCtrl, vt.KeyI)
	term.KeyTap(vt.KeyEsc)

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
	term.KeyTap(vt.KeyF1)                                        // SS3 P
	term.KeyTap(vt.KeyLShift, vt.KeyF1)                          // CSI 1 ; 2 P
	term.KeyTap(vt.KeyLCtrl, vt.KeyF2)                           // CSI 1 ; 5 Q
	term.KeyTap(vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyF3) // ESC CSI 1 ; 6 R
	term.KeyTap(vt.KeyRAlt, vt.KeyRCtrl, vt.KeyF4)               // ESC CSI 1 ; 5 S

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
	term.KeyTap(vt.KeyF5)                                        // CSI 15 ~
	term.KeyTap(vt.KeyRShift, vt.KeyF6)                          // CSI 17 ; 2 ~
	term.KeyTap(vt.KeyLCtrl, vt.KeyF7)                           // CSI 18 ; 5 ~
	term.KeyTap(vt.KeyLAlt, vt.KeyRCtrl, vt.KeyLShift, vt.KeyF8) // ESC CSI 19 ; 6 ~
	term.KeyTap(vt.KeyRAlt, vt.KeyLCtrl, vt.KeyF9)               // ESC CSI 20 ; 5 ~
	term.KeyTap(vt.KeyF20)                                       // CSI 34 ~
	term.KeyTap(vt.KeyF15)                                       // CSI 28 ~
	term.KeyTap(vt.KeyMenu)                                      // CSI 29 ~
	want = "\x1b[15~"
	want += "\x1b[17;2~"
	want += "\x1b[18;5~"
	want += "\x1b\x1b[19;6~"
	want += "\x1b\x1b[20;5~"
	want += "\x1b[34~"
	want += "\x1b[28~"
	want += "\x1b[29~"
	n, err = term.Read(buf)
	AssertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	VerifyF(t, result == want, "key responses failed: %q != %q", result, want)

	// Misc other keys
	clear(buf)
	term.KeyTap(vt.KeyEnter)                                // \r
	term.KeyTap(vt.KeyTab)                                  // \t
	term.KeyTap(vt.KeyLShift, vt.KeyTab)                    // CSI Z
	term.KeyTap(vt.KeyLCtrl, vt.KeyM)                       // \r
	term.KeyTap(vt.KeyLCtrl, vt.KeyL)                       // \f
	term.KeyTap(vt.KeyBackspace)                            // \x7f
	term.KeyTap(vt.KeyRCtrl, vt.KeyBackspace)               // \x08
	term.KeyTap(vt.KeyRCtrl, vt.KeyLShift, vt.KeyBackspace) // none
	term.KeyTap(vt.KeyRCtrl, vt.KeySpace)                   // \x00
	term.KeyTap(vt.KeySpace)                                // ' '
	term.KeyTap(vt.KeyRAlt, vt.KeyA)                        // \x1b a
	term.KeyTap(vt.KeyRHyper, vt.KeyA)                      // none
	term.KeyTap(vt.KeyRMeta, vt.KeyA)                       // none
	term.KeyTap(vt.KeyRAlt, vt.KeyLCtrl, vt.KeyJ)           // \x1b\n
	term.KeyTap(vt.KeyRCtrl, vt.KeyL)                       // \x0c
	term.KeyTap(vt.KeyLCtrl, vt.KeyLBrace)                  // \x0c

	want = "\r\t\x1b[Z\r\f\x7f\x08\x00 \x1ba\x1b\n\x0c\x1b"
	n, err = term.Read(buf)
	AssertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	VerifyF(t, result == want, "key responses failed: %q != %q", result, want)

	// Legacy control key mappings (weird ones)
	// Declining a few of the strange ones (control-?)
	clear(buf)
	term.KeyTap(vt.KeyLCtrl, vt.Key8)      // \x7F
	term.KeyTap(vt.KeyLCtrl, vt.Key4)      // \x1c
	term.KeyTap(vt.KeyLCtrl, vt.Key7)      // \x1f
	term.KeyTap(vt.Key7)                   // 7
	term.KeyTap(vt.KeyLShift, vt.KeySlash) // ?
	term.KeyTap(vt.KeyRCtrl, vt.KeyLBrace) // \x1b

	want = "\x7f\x1c\x1f7?\x1b"
	n, err = term.Read(buf)
	AssertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	VerifyF(t, result == want, "key responses failed: %q != %q", result, want)

	// Application cursor keys
	term.KeyTap(vt.KeyUp)
	WriteF(t, term, "\x1b[?1h")
	term.KeyTap(vt.KeyDown)
	want = "\x1b[A\x1bOB"
	n, err = term.Read(buf)
	AssertF(t, err == nil, "failed read: %v", err)

	result = string(buf[:n])
	VerifyF(t, result == want, "key responses failed: %q != %q", result, want)
}

// TestSgrAttr tests a variety of combinations of Sgr settings.
func TestSgrAttr(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[H")
	WriteF(t, term, "\x1b[1mA") // bold
	CheckAttrs(t, term, 0, 0, vt.Bold)
	CheckContent(t, term, 0, 0, "A")

	WriteF(t, term, "\x1b[2mB") // dim
	CheckAttrs(t, term, 1, 0, vt.Dim)
	CheckContent(t, term, 1, 0, "B")

	WriteF(t, term, "\x1b[22mC") // clear both
	CheckAttrs(t, term, 2, 0, vt.Plain)
	CheckContent(t, term, 2, 0, "C")

	WriteF(t, term, "\x1b[3;2mD") // italic, dim
	CheckAttrs(t, term, 3, 0, vt.Italic|vt.Dim)
	CheckContent(t, term, 3, 0, "D")

	WriteF(t, term, "\x1b[22mE") // remove dim, should leave italic
	CheckAttrs(t, term, 4, 0, vt.Italic)
	CheckContent(t, term, 4, 0, "E")

	WriteF(t, term, "\x1b[23mF") // clear italic
	CheckAttrs(t, term, 5, 0, vt.Plain)
	CheckContent(t, term, 5, 0, "F")

	WriteF(t, term, "\x1b[3;4mG") // simple underline
	CheckAttrs(t, term, 6, 0, vt.Italic|vt.Underline)
	CheckContent(t, term, 6, 0, "G")

	WriteF(t, term, "\x1b[21mH") // double underline (ECMA)
	CheckAttrs(t, term, 7, 0, vt.Italic|vt.DoubleUnderline)
	CheckContent(t, term, 7, 0, "H")

	WriteF(t, term, "\x1b[4mI") // simple underline
	CheckAttrs(t, term, 8, 0, vt.Italic|vt.Underline)
	CheckContent(t, term, 8, 0, "I")

	WriteF(t, term, "\x1b[4:2mJ") // double underline
	CheckAttrs(t, term, 9, 0, vt.Italic|vt.DoubleUnderline)
	CheckContent(t, term, 9, 0, "J")

	WriteF(t, term, "\x1b[4:3mK") // curly underline
	CheckAttrs(t, term, 10, 0, vt.Italic|vt.CurlyUnderline)
	CheckContent(t, term, 10, 0, "K")

	WriteF(t, term, "\x1b[4:4mL") // dotted underline
	CheckAttrs(t, term, 11, 0, vt.Italic|vt.DottedUnderline)
	CheckContent(t, term, 11, 0, "L")

	WriteF(t, term, "\x1b[4:5mM") // dashed underline
	CheckAttrs(t, term, 12, 0, vt.Italic|vt.DashedUnderline)
	CheckContent(t, term, 12, 0, "M")

	WriteF(t, term, "\x1b[4:9mN") // junk treats as plain underline
	CheckAttrs(t, term, 13, 0, vt.Italic|vt.Underline)
	CheckContent(t, term, 13, 0, "N")

	WriteF(t, term, "\x1b[4:5;24mO") // clustering, clear it
	CheckAttrs(t, term, 14, 0, vt.Italic)
	CheckContent(t, term, 14, 0, "O")

	WriteF(t, term, "\x1b[0;9;7;53mP") // clear, strike-through, reverse, over-lined
	CheckAttrs(t, term, 15, 0, vt.StrikeThrough|vt.Reverse|vt.Overline)
	CheckContent(t, term, 15, 0, "P")

	WriteF(t, term, "\x1b[5;27;29;55mQ")
	CheckAttrs(t, term, 16, 0, vt.Blink)
	CheckContent(t, term, 16, 0, "Q")

	WriteF(t, term, "\x1b[25mR")
	CheckAttrs(t, term, 17, 0, vt.Plain)
	CheckContent(t, term, 17, 0, "R")
}

// TestSgrColor8 tests simple ECMA 48 ANSI color (only 8 possible color values.)
func TestSgrColor8(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[36;42m\x1b#8")
	CheckColors(t, term, 0, 0, color.Teal, color.Green)

	WriteF(t, term, "\x1b[H\x1b[39mA")
	CheckColors(t, term, 0, 0, color.Silver, color.Green)

	WriteF(t, term, "\x1b[49mA")
	CheckColors(t, term, 1, 0, color.Silver, color.Black)

	// verify zero clears colors, first write some non zero colors
	WriteF(t, term, "\x1b[36;42mD")
	CheckColors(t, term, 2, 0, color.Teal, color.Green)

	// then send zero, should go to default colors
	WriteF(t, term, "\x1b[0mA")
	CheckColors(t, term, 3, 0, color.Silver, color.Black)
}

// TestSgrColor256 tests simple ECMA 48 ANSI color (256 possible color values.)
func TestSgrColor256(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(256))
	defer MustClose(t, term)
	MustStart(t, term)

	// foreground
	WriteF(t, term, "\x1b[38:5:6m\x1b[42m\x1b#8")
	CheckColors(t, term, 0, 0, color.Teal, color.Green)
	WriteF(t, term, "\x1b[38;5;5m\x1b[42mA")
	CheckColors(t, term, 0, 0, color.XTerm5, color.Green)

	WriteF(t, term, "\x1b[38;5;212;1mB")
	CheckColors(t, term, 1, 0, color.XTerm212, color.Green)
	CheckAttrs(t, term, 1, 0, vt.Bold)

	// background
	WriteF(t, term, "\x1b[48;5;2mA")
	CheckColors(t, term, 2, 0, color.XTerm212, color.XTerm2)
	CheckAttrs(t, term, 2, 0, vt.Bold)

	WriteF(t, term, "\x1b[48:5:134;2mC")
	CheckColors(t, term, 3, 0, color.XTerm212, color.XTerm134)
	CheckAttrs(t, term, 3, 0, vt.Dim)

	// mix background and foreground using colons
	WriteF(t, term, "\x1b[48:5:135;38:5:22mC")
	CheckColors(t, term, 4, 0, color.XTerm22, color.XTerm135)

	// and using semicolons
	WriteF(t, term, "\x1b[48;5;136;38;5;23mC")
	CheckColors(t, term, 5, 0, color.XTerm23, color.XTerm136)

	// underline colors - it uses the same parser so we won't check
	// all the variations
	WriteF(t, term, "\x1b[58;5;21;4mC")
	VerifyF(t, term.GetCell(Coord{X: 6, Y: 0}).S.Uc() == color.XTerm21, "underline color is wrong")

	WriteF(t, term, "\x1b[H\x1b[39;49mA")
	CheckColors(t, term, 0, 0, color.Silver, color.Black)

	// verify zero clears colors, first write some non zero colors
	WriteF(t, term, "\x1b[H\x1b[36;42mA")
	CheckColors(t, term, 0, 0, color.Teal, color.Green)

	// then send zero, should go to default colors
	WriteF(t, term, "\x1b[H\x1b[0mA")
	CheckColors(t, term, 0, 0, color.Silver, color.Black)

	// fuzz some things
	WriteF(t, term, "\x1b[m\x1b[H")
	WriteF(t, term, "\x1b[38:3m")
	WriteF(t, term, "\x1b[38;3m")
	WriteF(t, term, "\x1b[38;2m")
	WriteF(t, term, "\x1b[38;2:300m")
	WriteF(t, term, "\x1b[38;5m")
	WriteF(t, term, "\x1b[38;5:300m")
	WriteF(t, term, "\x1b[38:2m")
	WriteF(t, term, "\x1b[38:5m")
	WriteF(t, term, "\x1b[38:2;1;1;1m")
	WriteF(t, term, "\x1b[38:2:1:1m")
	WriteF(t, term, "A")
	CheckColors(t, term, 0, 0, color.Silver, color.Black)
}

// TestSgrColorRGB tests simple ECMA 48 ANSI color (full 24-bit color.)
func TestSgrColorRGB(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24}, vt.MockOptColors(1<<24))
	defer MustClose(t, term)
	MustStart(t, term)

	// foreground
	WriteF(t, term, "\x1b[m\x1b[H\x1b[38:2:255:0:0mA")
	CheckColors(t, term, 0, 0, color.NewRGBColor(255, 0, 0), color.Black)

	WriteF(t, term, "\x1b[m\x1b[H\x1b[38;2;2;0;0mA")
	CheckColors(t, term, 0, 0, color.NewRGBColor(2, 0, 0), color.Black)

	// background
	WriteF(t, term, "\x1b[m\x1b[H\x1b[48:2:1:2:3mA")
	CheckColors(t, term, 0, 0, color.Silver, color.NewRGBColor(1, 2, 3))

	WriteF(t, term, "\x1b[m\x1b[H\x1b[48;2;4;5;6mA")
	CheckColors(t, term, 0, 0, color.Silver, color.NewRGBColor(4, 5, 6))

	// full colors
	WriteF(t, term, "\x1b[m\x1b[H\x1b[38;2;99;88;77;48;2;4;5;6;58;2;99;98;91;1mA")
	CheckColors(t, term, 0, 0, color.NewRGBColor(99, 88, 77), color.NewRGBColor(4, 5, 6))
	VerifyF(t, term.GetCell(Coord{X: 0, Y: 0}).S.Uc() == color.NewRGBColor(99, 98, 91), "underline color is wrong")
	CheckAttrs(t, term, 0, 0, vt.Bold)
}

// TestTitles tests that we can set a window title.
func TestTitles(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	defer MustClose(t, term)
	MustStart(t, term)
	WriteF(t, term, "\x1b]2;Test Application\x1b\\")
	if s := term.GetTitle(); s != "Test Application" {
		t.Errorf("wrong title: %q", s)
	}

	// test ST termination using legacy bell character
	WriteF(t, term, "\x1b]2;Bell Ring\x07")
	if s := term.GetTitle(); s != "Bell Ring" {
		t.Errorf("wrong title: %q", s)
	}

	// try using 8-bit sequence
	WriteF(t, term, "\x9d2;Eight Bits\x9c")
	if s := term.GetTitle(); s != "Eight Bits" {
		t.Errorf("wrong title: %q", s)
	}
}

// TestResize tests resizing the terminal
func TestResize(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	defer MustClose(t, term)
	MustStart(t, term)

	// with E, and enable notifications
	WriteF(t, term, "\x1b#8\x1b[?2048h")
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
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "a\tC")
	if s := string(term.GetCell(Coord{X: 8, Y: 0}).C); s != "C" {
		t.Errorf("tab did not work: %q", s)
	}
	WriteF(t, term, "\x1b[3I")
	if x := term.Pos().X; x != 32 {
		t.Errorf("wrong position %d", x)
	}

	WriteF(t, term, "\x1b[2Z")
	if x := term.Pos().X; x != 16 {
		t.Errorf("wrong position %d", x)
	}

	WriteF(t, term, "\x1b[3g") // clear all tabs
	WriteF(t, term, "\x1b[I")  // one tab, should go to right margin
	if x := term.Pos().X; x != 79 {
		t.Errorf("wrong position: %d", x)
	}

	WriteF(t, term, "\x1b[Z")
	if x := term.Pos().X; x != 0 {
		t.Errorf("wrong position: %d", x)
	}

	// reset tabs
	WriteF(t, term, "\x1b[?5W")

	WriteF(t, term, "\t")
	if x := term.Pos().X; x != 8 {
		t.Errorf("wrong position: %d", x)
	}
	// clear this position, advance one
	WriteF(t, term, "\x1b[gA")
	if x := term.Pos().X; x != 9 {
		t.Errorf("wrong position: %d", x)
	}
	WriteF(t, term, "\x1bH")
	WriteF(t, term, "\t")
	if x := term.Pos().X; x != 16 {
		t.Errorf("wrong position: %d", x)
	}
	WriteF(t, term, "\x1b[Z")
	if x := term.Pos().X; x != 9 {
		t.Errorf("wrong position: %d", x)
	}
	WriteF(t, term, "\x1b[Z")
	if x := term.Pos().X; x != 0 {
		t.Errorf("wrong position: %d", x)
	}
	WriteF(t, term, "\x1b[1;10H") // goto position 9
	if x := term.Pos().X; x != 9 {
		t.Errorf("wrong position: %d", x)
	}
	// delete this one (do it twice to exercise the does not exist flow)
	WriteF(t, term, "\x1b[0g")
	WriteF(t, term, "\x1b[0g")

	// advance to next tab, then back, we should go to 0
	WriteF(t, term, "\t")
	if x := term.Pos().X; x != 16 {
		t.Errorf("wrong position: %d", x)
	}
	WriteF(t, term, "\x1b[Z")
	if x := term.Pos().X; x != 0 {
		t.Errorf("wrong position: %d", x)
	}
	WriteF(t, term, "\x1b[20I")
	WriteF(t, term, "\t\t")
	if pos := term.Pos(); pos.X != 79 || pos.Y != 0 {
		t.Errorf("wrong position: %d %d", pos.X, pos.Y)
	}

	// now backwards
	WriteF(t, term, "\x1b[20Z")
	WriteF(t, term, "\x1b[Z")
	if pos := term.Pos(); pos.X != 0 || pos.Y != 0 {
		t.Errorf("wrong position: %d %d", pos.X, pos.Y)
	}
}

func TestVerticalPos(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[2;2H")
	WriteF(t, term, "\x1b[10d")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 9 {
		t.Errorf("wrong position: %d %d", pos.X, pos.Y)
	}
	WriteF(t, term, "\x1b[2e")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 11 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	WriteF(t, term, "\x1b[e")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 12 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	WriteF(t, term, "\x1b[0e")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 13 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	WriteF(t, term, "\x1b[50d")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 23 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	WriteF(t, term, "\x1b[50e")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 23 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	WriteF(t, term, "\x1b[0d")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 0 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	WriteF(t, term, "\x1b[10d")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 9 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
	WriteF(t, term, "\x1b[1d")
	if pos := term.Pos(); pos.X != 1 || pos.Y != 0 {
		t.Errorf("wrong position; %d %d", pos.X, pos.Y)
	}
}

func TestSaveCursorPosition(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[1;1H\x1b[J")
	WriteF(t, term, "\x1b[1;5H")
	WriteF(t, term, "A")
	WriteF(t, term, "\x1b7")
	WriteF(t, term, "\x1b[1;1H")
	WriteF(t, term, "B")
	WriteF(t, term, "\x1b8")
	WriteF(t, term, "X")

	CheckContent(t, term, 0, 0, "B")
	CheckContent(t, term, 4, 0, "A")
	CheckContent(t, term, 5, 0, "X")
}

func TestSaveCursorWrapState(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[1;1H\x1b[J")
	WriteF(t, term, "\x1b[80G")
	WriteF(t, term, "A")
	WriteF(t, term, "\x1b7") // save cursor
	WriteF(t, term, "\x1b[1;1H")
	WriteF(t, term, "B")
	WriteF(t, term, "\x1b8") // restore cursor
	WriteF(t, term, "X")
	CheckContent(t, term, 0, 0, "B")
	CheckContent(t, term, 79, 0, "A")
	CheckContent(t, term, 0, 1, "X")
}

func TestSaveCursorSgr(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	defer MustClose(t, term)
	MustStart(t, term)
	WriteF(t, term, "\x1b[1;1H\x1b[0J")
	WriteF(t, term, "\x1b[1;4;33;44m")
	WriteF(t, term, "A")
	CheckPos(t, term, 1, 0)
	WriteF(t, term, "\x1b7")
	CheckPos(t, term, 1, 0)
	WriteF(t, term, "\x1b[0m")
	CheckPos(t, term, 1, 0)
	WriteF(t, term, "BE")
	CheckPos(t, term, 3, 0)
	WriteF(t, term, "\x1b8")
	CheckPos(t, term, 1, 0)
	WriteF(t, term, "X")
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "X")
	CheckContent(t, term, 2, 0, "E")
	CheckAttrs(t, term, 0, 0, vt.Bold|vt.Underline)
	CheckAttrs(t, term, 1, 0, vt.Bold|vt.Underline)
	CheckAttrs(t, term, 2, 0, vt.Plain)
	CheckColors(t, term, 0, 0, color.XTerm3, color.XTerm4)
	CheckColors(t, term, 1, 0, color.XTerm3, color.XTerm4)
	CheckColors(t, term, 2, 0, color.Silver, color.Black)
}

func TestReset(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 80, Y: 24})
	defer MustClose(t, term)
	MustStart(t, term)

	// write a bunch of stuff to create state (so we can verify it gets reset)
	WriteF(t, term, "\x1b[1;4;33;44m")
	WriteF(t, term, "\x1b#8")
	WriteF(t, term, "\x1b[1;80HX")
	WriteF(t, term, "\x1b7")    // save cursor
	WriteF(t, term, "\x1b[?7l") // disable automargin
	WriteF(t, term, "\x1bc")
	CheckPos(t, term, 0, 0)
	for row := range Row(24) {
		for col := range Col(80) {
			CheckAttrs(t, term, col, row, vt.Plain)
			CheckContent(t, term, col, row, "")
			CheckColors(t, term, col, row, color.Silver, color.Black)
		}
	}
	WriteF(t, term, "\x1b8") // restore cursor
	CheckPos(t, term, 0, 0)
	WriteF(t, term, "X")
	CheckAttrs(t, term, 0, 0, vt.Plain)
	CheckContent(t, term, 0, 0, "X")
	CheckColors(t, term, 0, 0, color.Silver, color.Black)
	WriteF(t, term, "\x1b[?7$p")

	// verify mode reset
	want := "\x1b[?7;1$y"
	buf := make([]byte, 128)
	n, err := term.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
	}
	result := string(buf[:n])
	VerifyF(t, result == want, "wrong mode: %q != %q", result, want)
}

// backendBox makes a backend box, filled with increasing letters (modulo 16)
func backendBox(t *testing.T, mb vt.MockBackend, tl Coord, br Coord, attr Attr) {
	t.Helper()
	style := vt.BaseStyle.WithAttr(attr)
	hex := []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P'}
	AssertF(t, len(hex) == 16, "wrong hex string")
	i := 0
	for row := tl.Y; row <= br.Y; row++ {
		for col := tl.X; col <= br.X; col++ {
			mb.Put(Coord{Y: row, X: col}, vt.Cell{C: string(hex[i%16]), W: 1, S: style})
			i++
		}
	}
}

func TestBackendBlit(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 10, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	mb := term.Backend()
	AssertF(t, mb != nil, "backend is nil")

	backendBox(t, mb, Coord{X: 0, Y: 0}, Coord{X: 1, Y: 1}, vt.Bold)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "C")
	CheckContent(t, term, 1, 1, "D")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckAttrs(t, term, 2, 2, vt.Plain)

	// blit the entire box down and right 1
	mb.(vt.Blitter).Blit(Coord{X: 0, Y: 0}, Coord{X: 1, Y: 1}, Coord{X: 2, Y: 2})
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "C")
	CheckContent(t, term, 1, 1, "A")
	CheckContent(t, term, 2, 1, "B")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "C")
	CheckContent(t, term, 2, 2, "D")

	// spot check attributes
	CheckAttrs(t, term, 2, 2, vt.Bold)
}

func TestBackendBlitReverse(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 10, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	mb := term.Backend()
	AssertF(t, mb != nil, "backend is nil")

	backendBox(t, mb, Coord{X: 1, Y: 1}, Coord{X: 2, Y: 2}, vt.Bold)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "A")
	CheckContent(t, term, 2, 1, "B")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "C")
	CheckContent(t, term, 2, 2, "D")
	CheckAttrs(t, term, 0, 0, vt.Plain)

	// blit the entire box down and right 1
	mb.(vt.Blitter).Blit(Coord{X: 1, Y: 1}, Coord{X: 0, Y: 0}, Coord{X: 2, Y: 2})
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "C")
	CheckContent(t, term, 1, 1, "D")
	CheckContent(t, term, 2, 1, "B")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "C")
	CheckContent(t, term, 2, 2, "D")

	// spot check attributes
	CheckAttrs(t, term, 2, 2, vt.Bold)
}

func TestGraphemeCluster(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 10, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[H")
	// grapheme clustering is off
	WriteF(t, term, "🇨🇭") // flag + regional indicator C + regional indicator H
	WriteF(t, term, "A")

	// we advanced by four columns (two wide emoji), and a single character
	CheckPos(t, term, 5, 0)
	CheckContent(t, term, 0, 0, "\U0001f1e8") // regional indicator C
	CheckContent(t, term, 1, 0, "")           // empty
	CheckContent(t, term, 2, 0, "\U0001f1ed") // regional indicator H
	CheckContent(t, term, 3, 0, "")           // empty
	CheckContent(t, term, 4, 0, "A")          // empty

	// now turn on grapheme clustering
	WriteF(t, term, "\x1b[?2027h")
	WriteF(t, term, "\x1b[H\x1b[J")
	CheckPos(t, term, 0, 0)
	WriteF(t, term, "\U0001f1e8\U0001f1ed")
	CheckPos(t, term, 2, 0)
	WriteF(t, term, "A")
	CheckPos(t, term, 3, 0)
	CheckContent(t, term, 0, 0, "🇨🇭")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "A")

	// lets also verify it works with automargin
	// RECALL: Maximum width is 10
	WriteF(t, term, "\x1b[7h")    // should already be on
	WriteF(t, term, "\x1b[1;10H") // last position in first row
	WriteF(t, term, "🇨🇭A")        // flag + regional indicator C + regional indicator H
	CheckPos(t, term, 1, 1)
	CheckContent(t, term, 9, 0, "🇨🇭")
	CheckContent(t, term, 0, 1, "A")
}

func TestEraseAbove(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 10, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b#8")
	WriteF(t, term, "\x1b[3;5H")

	WriteF(t, term, "\x1b[31;1;42m") // set some colors and bold
	WriteF(t, term, "\x1b[1J")
	// cursor is at 4,2
	for row := range Row(5) {
		for col := range Col(10) {
			if row < 2 || row == 2 && col < 5 {
				CheckContent(t, term, col, row, "")
				CheckAttrs(t, term, col, row, vt.Plain)
				CheckColors(t, term, col, row, color.XTerm1, color.XTerm2)
			} else {
				CheckContent(t, term, col, row, "E")
				CheckAttrs(t, term, col, row, vt.Plain)
				CheckColors(t, term, col, row, color.Silver, color.Black)
			}
		}
	}
}

func TestEraseLine(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 10, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b#8")

	// erase to end
	WriteF(t, term, "\x1b[31;1;42m") // set some colors and bold
	WriteF(t, term, "\x1b[1;5H")
	WriteF(t, term, "\x1b[0K")
	// cursor is at 4,0
	CheckPos(t, term, 4, 0)
	for col := range Col(10) {
		row := Row(0)
		if col < 4 {
			CheckContent(t, term, col, row, "E")
			CheckAttrs(t, term, col, row, vt.Plain)
			CheckColors(t, term, col, row, color.Silver, color.Black)
		} else {
			CheckContent(t, term, col, row, "")
			CheckAttrs(t, term, col, row, vt.Plain)
			CheckColors(t, term, col, row, color.XTerm1, color.XTerm2)
		}
	}
	CheckPos(t, term, 4, 0)

	// erase to beginning
	WriteF(t, term, "\x1b[2;5H")
	WriteF(t, term, "\x1b[1K")
	// cursor is at 4,1
	CheckPos(t, term, 4, 1)
	for col := range Col(10) {
		row := Row(1)
		if col > 4 {
			CheckContent(t, term, col, row, "E")
			CheckAttrs(t, term, col, row, vt.Plain)
			CheckColors(t, term, col, row, color.Silver, color.Black)
		} else {
			CheckContent(t, term, col, row, "")
			CheckAttrs(t, term, col, row, vt.Plain)
			CheckColors(t, term, col, row, color.XTerm1, color.XTerm2)
		}
	}
	CheckPos(t, term, 4, 1)

	// erase entire line
	WriteF(t, term, "\x1b[3;5H")
	WriteF(t, term, "\x1b[2K")
	// cursor is at 4,2
	CheckPos(t, term, 4, 2)
	for col := range Col(10) {
		row := Row(2)
		CheckContent(t, term, col, row, "")
		CheckAttrs(t, term, col, row, vt.Plain)
		CheckColors(t, term, col, row, color.XTerm1, color.XTerm2)
	}
	CheckPos(t, term, 4, 2)
}

// TestNewLineScroll tests scrolling with a new line.
// This is one of the most fundamental operations for a terminal.
func TestNewLineScroll(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 10, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[H\x1b[J") // home and clear
	WriteF(t, term, "\x1b[1;1HA")
	WriteF(t, term, "\x1b[2;1HB")
	WriteF(t, term, "\x1b[5;1HC") // first column on last row
	CheckContent(t, term, 0, 4, "C")
	WriteF(t, term, "\n") // new line should scroll
	CheckPos(t, term, 1, 4)
	CheckContent(t, term, 0, 0, "B")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 0, 4, "")
	CheckContent(t, term, 0, 3, "C")
}

// TestNewLineScrollNoBlitter tests scrolling with a new line,
// using the fallback copy for backends without Blit support.
func TestNewLineScrollNoBlitter(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 10, Y: 5}, vt.MockOptNoBlit{})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[H\x1b[J") // home and clear
	WriteF(t, term, "\x1b[1;1HA")
	WriteF(t, term, "\x1b[2;1HB")
	WriteF(t, term, "\x1b[5;1HC") // first column on last row
	CheckContent(t, term, 0, 4, "C")
	WriteF(t, term, "\n") // new line should scroll
	CheckPos(t, term, 1, 4)
	CheckContent(t, term, 0, 0, "B")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 0, 4, "")
	CheckContent(t, term, 0, 3, "C")
}

func TestNewLineModes(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 10, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)
	WriteF(t, term, "\x1b[H\x1b[J")
	WriteF(t, term, "ABC\n")
	CheckPos(t, term, 3, 1)
	term.KeyTap(vt.KeyEnter)
	WriteF(t, term, "\x1b[20h")
	term.KeyTap(vt.KeyEnter)
	WriteF(t, term, "DEF")
	CheckPos(t, term, 6, 1)
	WriteF(t, term, "\n")
	CheckPos(t, term, 0, 2)
	WriteF(t, term, "GHI")
	CheckPos(t, term, 3, 2)
	WriteF(t, term, "\x1b[20l")
	WriteF(t, term, "\n")
	CheckPos(t, term, 3, 3)
	term.KeyTap(vt.KeyEnter)

	// |ABC_____|
	// |___DEF__|
	// |GHI_____|
	// |___c____|
	//
	// input stream contains \r\r\n\r

	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 4, 0, "")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 3, 1, "D")
	CheckContent(t, term, 4, 1, "E")
	CheckContent(t, term, 5, 1, "F")
	CheckContent(t, term, 0, 2, "G")
	CheckContent(t, term, 1, 2, "H")
	CheckContent(t, term, 2, 2, "I")
	CheckContent(t, term, 3, 2, "")
	CheckContent(t, term, 4, 2, "")
	CheckContent(t, term, 5, 2, "")

	result := ReadF(t, term)
	want := "\r\r\n\r"
	VerifyF(t, result == want, "response incorrect: %q != %q", result, want)
}

// TestScrollUp tests scrolling up. The cursor position is retained.
func TestScrollUp(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 10, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[H\x1b[J") // home and clear
	WriteF(t, term, "\x1b[1;1HA")
	WriteF(t, term, "\x1b[2;1HB")
	WriteF(t, term, "\x1b[5;1HC") // first column on last row
	CheckContent(t, term, 0, 4, "C")
	WriteF(t, term, "\x1b[5;5H") // fifth column on last row
	CheckPos(t, term, 4, 4)
	WriteF(t, term, "\x1bD")
	CheckPos(t, term, 4, 4)
	CheckContent(t, term, 0, 0, "B")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 0, 4, "")
	CheckContent(t, term, 0, 3, "C")
	WriteF(t, term, "\x1bE") // this is like a newline
	CheckPos(t, term, 0, 4)
	CheckContent(t, term, 0, 2, "C")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 0, 4, "")

	WriteF(t, term, "\x1b[H\x1bJ")
	WriteF(t, term, "\x1b[1;1HA")
	WriteF(t, term, "\x1b[2;1HB")
	WriteF(t, term, "\x1b[3;1HC")
	WriteF(t, term, "\x1b[4;1HD")
	WriteF(t, term, "\x1b[5;1HE")
	WriteF(t, term, "\x1b[3;3H")
	WriteF(t, term, "\x1b[3S") // scroll in place, leaves cursor where it is
	CheckContent(t, term, 0, 0, "D")
	CheckContent(t, term, 0, 1, "E")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 0, 4, "")
	CheckPos(t, term, 2, 2)
}

// TestScrollDown tests scrolling down. The cursor position is retained.
func TestScrollDown(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 10, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\x1b[H\x1b[J") // home and clear
	WriteF(t, term, "\x1b[1;1HA")
	WriteF(t, term, "\x1b[2;1HB")
	WriteF(t, term, "\x1b[4;1HC") // first column on penultimate row
	WriteF(t, term, "\x1b[5;1HD") // first column on last row
	CheckContent(t, term, 0, 3, "C")
	CheckContent(t, term, 0, 4, "D")
	WriteF(t, term, "\x1b[1;5H") // fifth column on first row
	CheckPos(t, term, 4, 0)
	WriteF(t, term, "\x1bM")
	CheckPos(t, term, 4, 0)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 0, 1, "A")
	CheckContent(t, term, 0, 2, "B")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 0, 4, "C")

	WriteF(t, term, "\x1b[H\x1bJ")
	WriteF(t, term, "\x1b[1;1HA")
	WriteF(t, term, "\x1b[2;1HB")
	WriteF(t, term, "\x1b[3;1HC")
	WriteF(t, term, "\x1b[4;1HD")
	WriteF(t, term, "\x1b[5;1HE")
	WriteF(t, term, "\x1b[3;3H")
	WriteF(t, term, "\x1b[3T") // scroll in place, leaves cursor where it is
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 0, 3, "A")
	CheckContent(t, term, 0, 4, "B")
	CheckPos(t, term, 2, 2)
}
