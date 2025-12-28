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
	if !strings.HasPrefix(result, "\x1b[63") {
		t.Errorf("Missing prefix 'c': %q", result)
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
	if !strings.HasPrefix(result, "\x1b[63") {
		t.Errorf("Missing prefix 'c': %q", result)
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
	if !strings.HasPrefix(result, "\x1b[P>|") {
		t.Errorf("Missing prefix 'CSI P>|': %q", result)
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
