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
	"testing"

	"github.com/gdamore/tcell/v3/color"
)

// This file implements various tests of the emulator.  Much of these tests
// are "borrowed" (ported from) the tests from Ghostty - https://ghostty.org/docs/vt
// Note that Ghostty's tests assume that STTY modes to expand LF to CF LF are in
// effect (or ANSI mode 20.)  We don't assume that, and add the CR explicitly.

// TestDECSTBMv1 tests full screen scroll region
func TestDECSTBMv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")   //  clear screen
	writeF(t, term, "ABC\r\n")
	writeF(t, term, "DEF\r\n")
	writeF(t, term, "GHI\r\n")
	writeF(t, term, "\033[r") // scroll region top/bottom
	writeF(t, term, "\033[T") // scroll down one

	// |c_______|
	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	checkPos(t, term, 0, 0)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "A")
	checkContent(t, term, 1, 1, "B")
	checkContent(t, term, 2, 1, "C")
	checkContent(t, term, 0, 2, "D")
	checkContent(t, term, 1, 2, "E")
	checkContent(t, term, 2, 2, "F")
	checkContent(t, term, 0, 3, "G")
	checkContent(t, term, 1, 3, "H")
	checkContent(t, term, 2, 3, "I")
}

// TestDECSTBMv2 top only
func TestDECSTBMv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")   //  clear screen
	writeF(t, term, "ABC\r\n")
	writeF(t, term, "DEF\r\n")
	writeF(t, term, "GHI\r\n")
	writeF(t, term, "\033[2;2r") // scroll region top/bottom
	writeF(t, term, "\033[T")    // scroll down one

	// |________|
	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	checkPos(t, term, 0, 3) // did not move
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "A")
	checkContent(t, term, 1, 1, "B")
	checkContent(t, term, 2, 1, "C")
	checkContent(t, term, 0, 2, "D")
	checkContent(t, term, 1, 2, "E")
	checkContent(t, term, 2, 2, "F")
	checkContent(t, term, 0, 3, "G")
	checkContent(t, term, 1, 3, "H")
	checkContent(t, term, 2, 3, "I")
}

// TestDECSTBMv3 top and bottom
func TestDECSTBMv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")   //  clear screen
	writeF(t, term, "ABC\r\n")
	writeF(t, term, "DEF\r\n")
	writeF(t, term, "GHI\r\n")
	writeF(t, term, "\033[1;2r") // scroll region top/bottom
	writeF(t, term, "\033[T")    // scroll down one

	// |________|
	// |ABC_____|
	// |GHI_____|
	// |________|
	checkPos(t, term, 0, 0)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "A")
	checkContent(t, term, 1, 1, "B")
	checkContent(t, term, 2, 1, "C")
	checkContent(t, term, 0, 2, "G")
	checkContent(t, term, 1, 2, "H")
	checkContent(t, term, 2, 2, "I")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestDECSTBMv4 top == bottom
func TestDECSTBMv4(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "ABC\r\n")
	writeF(t, term, "DEF\r\n")
	writeF(t, term, "GHI\r\n")
	writeF(t, term, "\033[2;2r")
	writeF(t, term, "\033[T")

	// |________|
	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	checkPos(t, term, 0, 3)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "A")
	checkContent(t, term, 1, 1, "B")
	checkContent(t, term, 2, 1, "C")
	checkContent(t, term, 0, 2, "D")
	checkContent(t, term, 1, 2, "E")
	checkContent(t, term, 2, 2, "F")
	checkContent(t, term, 0, 3, "G")
	checkContent(t, term, 1, 3, "H")
	checkContent(t, term, 2, 3, "I")
}

// TestDECSLRMv1 tests full screen right and left margins.
func TestDECSLRMv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[H")
	writeF(t, term, "\033[J")
	writeF(t, term, "ABC\n\r")
	writeF(t, term, "DEF\n\r")
	writeF(t, term, "GHI\n\r")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[s")
	writeF(t, term, "\033[X")

	checkPos(t, term, 0, 0)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 0, 1, "D")
	checkContent(t, term, 1, 1, "E")
	checkContent(t, term, 2, 1, "F")
	checkContent(t, term, 0, 2, "G")
	checkContent(t, term, 1, 2, "H")
	checkContent(t, term, 2, 2, "I")
}

// TODO: DECSLRMv3 left and right this makes use of insert line
// TODO: DECSLRMv4 left and right equal
// TODO: add tests for actual left and right scrolling!

// TestRIv1 top of screen, no scroll
func TestRIv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "A\r\n")
	writeF(t, term, "B\r\n")
	writeF(t, term, "C\r\n")
	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033M")
	writeF(t, term, "X")

	// |Xc______|
	// |A_______|
	// |B_______|
	// |C_______|

	checkPos(t, term, 1, 0)
	checkContent(t, term, 0, 0, "X")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "A")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "B")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "C")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestRIv2 not top of screen, no scroll
func TestRIv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "A\r\n")
	writeF(t, term, "B\r\n")
	writeF(t, term, "C\r\n")
	writeF(t, term, "\033[2;1H")
	writeF(t, term, "\033M")
	writeF(t, term, "X")

	// |Xc______|
	// |B_______|
	// |C_______|
	// |________|

	checkPos(t, term, 1, 0)
	checkContent(t, term, 0, 0, "X")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "B")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "C")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestRIv3 scroll region
func TestRIv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")   //  clear screen
	writeF(t, term, "A\r\n")
	writeF(t, term, "B\r\n")
	writeF(t, term, "C\r\n")
	writeF(t, term, "\033[2;3r")
	writeF(t, term, "\033[2;1H")
	writeF(t, term, "\033M")

	// |A_______|
	// |c_______|
	// |B_______|
	// |________|

	checkPos(t, term, 0, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "B")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestRIv4 outside scroll region - goes to top, does not scroll
func TestRIv4(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")   //  clear screen
	writeF(t, term, "A\r\n")
	writeF(t, term, "B\r\n")
	writeF(t, term, "C\r\n")
	writeF(t, term, "\033[2;3r")
	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033M")

	// |A_______|
	// |B_______|
	// |C_______|
	// |________|

	checkPos(t, term, 0, 0)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "B")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "C")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TODO: RIv5 - left right scroll regions (when we implement left/right regions)
// TODO: RIv6 - outside left/right scroll regions (when we implement left/right regions)

// TestINDv1 no scroll region, top of screen
func TestINDv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")   //  clear screen
	writeF(t, term, "A")
	writeF(t, term, "\033D")
	writeF(t, term, "X")

	// |A_______|
	// |_Xc_____|
	// |________|
	// |________|

	checkPos(t, term, 2, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "X")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestINDv2 no scroll region, bottom of screen
func TestINDv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")   //  clear screen
	writeF(t, term, "\033[4;1H")
	writeF(t, term, "A")
	writeF(t, term, "\033D")
	writeF(t, term, "X")

	// |________|
	// |________|
	// |A_______|
	// |_Xc_____|

	checkPos(t, term, 2, 3)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "A")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "X")
	checkContent(t, term, 2, 3, "")
}

// TestINDv3 inside scroll region
func TestINDv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "A")
	writeF(t, term, "\033D")
	writeF(t, term, "X")

	// |A_______|
	// |_Xc_____|
	// |________|
	// |________|

	checkPos(t, term, 2, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "X")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestINDv4 bottom of scroll region
func TestINDv4(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "\033[4;1H")
	writeF(t, term, "B")
	writeF(t, term, "\033[3;1H")
	writeF(t, term, "A")
	writeF(t, term, "\033D")
	writeF(t, term, "X")

	// |________|
	// |A_______|
	// |_Xc_____|
	// |B_______|

	checkPos(t, term, 2, 2)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "A")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "X")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "B")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestINDv5 bottom of screen with scroll region
func TestINDv5(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "\033[3;1H")
	writeF(t, term, "A")
	writeF(t, term, "\033[4;1H")
	writeF(t, term, "\033D")
	writeF(t, term, "X")

	// |________|
	// |________|
	// |A_______|
	// |________|
	// |Xc______|

	checkPos(t, term, 1, 4)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "A")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
	checkContent(t, term, 0, 4, "X")
	checkContent(t, term, 1, 4, "")
	checkContent(t, term, 2, 4, "")
}

// TestINDv6 tests IND outside of left and right scroll region.
func TestINDv6(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "\033[3;5s")
	writeF(t, term, "\033[3;3H")
	writeF(t, term, "A")
	writeF(t, term, "\033[3;1H")
	writeF(t, term, "\033D")
	writeF(t, term, "X")

	// |________|
	// |________|
	// |XcA_____|

	checkPos(t, term, 1, 2)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "X")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "A")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestINDv7 tests IND inside of left and right scroll region.
func TestINDv7(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H") // move to top-left
	writeF(t, term, "\033[0J")
	writeF(t, term, "111111\n\r")
	writeF(t, term, "222222\n\r")
	writeF(t, term, "333333\n\r")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[1;3s")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "\033[3;1H")
	writeF(t, term, "\033D")

	// |222111__|
	// |333222__|
	// |c__333__|

	checkPos(t, term, 0, 2)
	checkContent(t, term, 0, 0, "2")
	checkContent(t, term, 1, 0, "2")
	checkContent(t, term, 2, 0, "2")
	checkContent(t, term, 3, 0, "1")
	checkContent(t, term, 4, 0, "1")
	checkContent(t, term, 5, 0, "1")
	checkContent(t, term, 0, 1, "3")
	checkContent(t, term, 1, 1, "3")
	checkContent(t, term, 2, 1, "3")
	checkContent(t, term, 3, 1, "2")
	checkContent(t, term, 4, 1, "2")
	checkContent(t, term, 5, 1, "2")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 3, 2, "3")
	checkContent(t, term, 4, 2, "3")
	checkContent(t, term, 5, 2, "3")
}

// TestCUDv1 - cursor down
func TestCUDv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "A")
	writeF(t, term, "\033[2B")
	writeF(t, term, "X")

	// |A_______|
	// |________|
	// |_Xc_____|
	// |________|

	checkPos(t, term, 2, 2)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "X")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestCUDv2 - cursor down above bottom margin
func TestCUDv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "\n\n\n\n")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "A")
	writeF(t, term, "\033[5B")
	writeF(t, term, "X")

	// |A_______|
	// |________|
	// |_Xc_____|
	// |________|

	checkPos(t, term, 2, 2)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "X")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestCUDv3 - cursor down below bottom margin
func TestCUDv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "A")
	writeF(t, term, "\033[4;1H")
	writeF(t, term, "\033[5B")
	writeF(t, term, "X")

	// |A_______|
	// |________|
	// |________|
	// |________|
	// |Xc______|

	checkPos(t, term, 1, 4)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
	checkContent(t, term, 0, 4, "X")
	checkContent(t, term, 1, 4, "")
	checkContent(t, term, 2, 4, "")
}

// TestCUUv1 tests cursor up.
func TestCUUv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[3;H")
	writeF(t, term, "A")
	writeF(t, term, "\033[2A")
	writeF(t, term, "X")

	// |_Xc_____|
	// |________|
	// |A_______|

	checkPos(t, term, 2, 0)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "X")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "A")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
}

// TestCUUv2 tests cursor up below the top margin.
func TestCUUv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[2;4r")
	writeF(t, term, "\033[3;1H")
	writeF(t, term, "A")
	writeF(t, term, "\033[5A")
	writeF(t, term, "X")

	// |________|
	// |_Xc_____|
	// |A_______|
	// |________|

	checkPos(t, term, 2, 1)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "X")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "A")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestCUUv3 tests cursor up above the top margin.
func TestCUUv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[3;5r")
	writeF(t, term, "\033[3;1H")
	writeF(t, term, "A")
	writeF(t, term, "\033[2;1H")
	writeF(t, term, "\033[5A")
	writeF(t, term, "X")

	// |Xc______|
	// |________|
	// |A_______|
	// |________|
	// |________|

	checkPos(t, term, 1, 0)
	checkContent(t, term, 0, 0, "X")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "A")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestCNLv1 - cursor next line
func TestCNLv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "A")
	writeF(t, term, "\033[2E")
	writeF(t, term, "X")

	// |A_______|
	// |________|
	// |Xc_____|
	// |________|

	checkPos(t, term, 1, 2)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "X")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestCNLv2 - cursor next line above bottom margin
func TestCNLv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 4})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "\n\n\n\n")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "A")
	writeF(t, term, "\033[5E")
	writeF(t, term, "X")

	// |A_______|
	// |________|
	// |Xc______|
	// |________|

	checkPos(t, term, 1, 2)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "X")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestCNLv3 - cursor next line bottom margin
func TestCNLv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "A")
	writeF(t, term, "\033[4;3H")
	writeF(t, term, "\033[5E")
	writeF(t, term, "X")

	// |A_______|
	// |________|
	// |________|
	// |________|
	// |Xc______|

	checkPos(t, term, 1, 4)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
	checkContent(t, term, 0, 4, "X")
	checkContent(t, term, 1, 4, "")
	checkContent(t, term, 2, 4, "")
}

// TestCPLv1 tests cursor previous line.
func TestCPLv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[3;H")
	writeF(t, term, "A")
	writeF(t, term, "\033[2F")
	writeF(t, term, "X")

	// |Xc______|
	// |________|
	// |A_______|

	checkPos(t, term, 1, 0)
	checkContent(t, term, 0, 0, "X")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "A")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
}

// TestCPLv2 tests cursor previous line below the top margin.
func TestCPLv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[2;4r")
	writeF(t, term, "\033[3;1H")
	writeF(t, term, "A")
	writeF(t, term, "\033[5F")
	writeF(t, term, "X")

	// |________|
	// |Xc______|
	// |A_______|
	// |________|

	checkPos(t, term, 1, 1)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "X")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "A")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestCPLv3 tests cursor previous line above the top margin.
func TestCPLv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H")
	writeF(t, term, "\033[0J")
	writeF(t, term, "\033[3;5r")
	writeF(t, term, "\033[3;1H")
	writeF(t, term, "A")
	writeF(t, term, "\033[2;2H")
	writeF(t, term, "\033[5F")
	writeF(t, term, "X")

	// |Xc______|
	// |________|
	// |A_______|
	// |________|
	// |________|

	checkPos(t, term, 1, 0)
	checkContent(t, term, 0, 0, "X")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 0, 2, "A")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
}

// TestILv1 tests inserting a line.
func TestILv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC\n\r")
	writeF(t, term, "DEF\n\r")
	writeF(t, term, "GHI\n\r")
	writeF(t, term, "\033[2;2H")
	writeF(t, term, "\033[L")

	// |ABC_____|
	// |c_______|
	// |DEF_____|
	// |GHI_____|
	// |________|

	checkPos(t, term, 0, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 3, 1, "")
	checkContent(t, term, 0, 2, "D")
	checkContent(t, term, 1, 2, "E")
	checkContent(t, term, 2, 2, "F")
	checkContent(t, term, 3, 2, "")
	checkContent(t, term, 0, 3, "G")
	checkContent(t, term, 1, 3, "H")
	checkContent(t, term, 2, 3, "I")
	checkContent(t, term, 3, 3, "")
}

// TestILv2 tests insert line outside of the scroll region.
func TestILv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC\n\r")
	writeF(t, term, "DEF\n\r")
	writeF(t, term, "GHI\n\r")
	writeF(t, term, "\033[3;4r")
	writeF(t, term, "\033[2;2H")
	writeF(t, term, "\033[L")

	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	// |________|
	// |________|

	checkPos(t, term, 1, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 0, 1, "D")
	checkContent(t, term, 1, 1, "E")
	checkContent(t, term, 2, 1, "F")
	checkContent(t, term, 3, 1, "")
	checkContent(t, term, 0, 2, "G")
	checkContent(t, term, 1, 2, "H")
	checkContent(t, term, 2, 2, "I")
	checkContent(t, term, 3, 2, "")
}

// TestILv3 tests insert line with top and bottom scroll regions.
func TestILv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC\n\r")
	writeF(t, term, "DEF\n\r")
	writeF(t, term, "GHI\n\r")
	writeF(t, term, "123\n\r")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "\033[2;2H")
	writeF(t, term, "\033[L")

	// |ABC_____|
	// |c_______|
	// |DEF_____|
	// |123_____|

	checkPos(t, term, 0, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 0, 1, "")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 3, 1, "")
	checkContent(t, term, 0, 2, "D")
	checkContent(t, term, 1, 2, "E")
	checkContent(t, term, 2, 2, "F")
	checkContent(t, term, 3, 2, "")
	checkContent(t, term, 0, 3, "1")
	checkContent(t, term, 1, 3, "2")
	checkContent(t, term, 2, 3, "3")
	checkContent(t, term, 3, 3, "")
}

// TestILv4 tests insert line with left and right margins.
func TestILv4(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC123\n\r")
	writeF(t, term, "DEF456\n\r")
	writeF(t, term, "GHI789\n\r")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[2;4s")
	writeF(t, term, "\033[2;2H")
	writeF(t, term, "\033[L")

	// |ABC123__|
	// |Dc__56__|
	// |GEF489__|
	// |_HI7____|

	checkPos(t, term, 0, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "1")
	checkContent(t, term, 4, 0, "2")
	checkContent(t, term, 5, 0, "3")
	checkContent(t, term, 0, 1, "D")
	checkContent(t, term, 1, 1, "")
	checkContent(t, term, 2, 1, "")
	checkContent(t, term, 3, 1, "")
	checkContent(t, term, 4, 1, "5")
	checkContent(t, term, 5, 1, "6")
	checkContent(t, term, 0, 2, "G")
	checkContent(t, term, 1, 2, "E")
	checkContent(t, term, 2, 2, "F")
	checkContent(t, term, 3, 2, "4")
	checkContent(t, term, 4, 2, "8")
	checkContent(t, term, 5, 2, "9")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "H")
	checkContent(t, term, 2, 3, "I")
	checkContent(t, term, 3, 3, "7")
	checkContent(t, term, 4, 3, "")
	checkContent(t, term, 5, 3, "")
}

// TestDLv1 tests a simple delete line.
func TestDLv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC\n\r")
	writeF(t, term, "DEF\n\r")
	writeF(t, term, "GHI\n\r")
	writeF(t, term, "\033[2;2H")
	writeF(t, term, "\033[M")

	// |ABC_____|
	// |GHI_____|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 0, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 0, 1, "G")
	checkContent(t, term, 1, 1, "H")
	checkContent(t, term, 2, 1, "I")
	checkContent(t, term, 3, 1, "")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 2, 3, "")
}

// TestDLv2 delete line outside of the scroll region.
func TestDLv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC\n\r")
	writeF(t, term, "DEF\n\r")
	writeF(t, term, "GHI\n\r")
	writeF(t, term, "\033[3;4r")
	writeF(t, term, "\033[2;2H")
	writeF(t, term, "\033[M")

	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	// |________|
	// |________|

	checkPos(t, term, 1, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 0, 1, "D")
	checkContent(t, term, 1, 1, "E")
	checkContent(t, term, 2, 1, "F")
	checkContent(t, term, 3, 1, "")
	checkContent(t, term, 0, 2, "G")
	checkContent(t, term, 1, 2, "H")
	checkContent(t, term, 2, 2, "I")
	checkContent(t, term, 3, 2, "")
}

// TestDLv3 delete line with top and bottom scroll regions.
func TestDLv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC\n\r")
	writeF(t, term, "DEF\n\r")
	writeF(t, term, "GHI\n\r")
	writeF(t, term, "123\n\r")
	writeF(t, term, "\033[1;3r")
	writeF(t, term, "\033[2;2H")
	writeF(t, term, "\033[M")

	// |ABC_____|
	// |GHI_____|
	// |________|
	// |123_____|

	checkPos(t, term, 0, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 0, 1, "G")
	checkContent(t, term, 1, 1, "H")
	checkContent(t, term, 2, 1, "I")
	checkContent(t, term, 3, 1, "")
	checkContent(t, term, 0, 2, "")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 3, 2, "")
	checkContent(t, term, 0, 3, "1")
	checkContent(t, term, 1, 3, "2")
	checkContent(t, term, 2, 3, "3")
	checkContent(t, term, 3, 3, "")
}

// TestDLv4 delete line with left and right margins.
func TestDLv4(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC123\n\r")
	writeF(t, term, "DEF456\n\r")
	writeF(t, term, "GHI789\n\r")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[2;4s")
	writeF(t, term, "\033[2;2H")
	writeF(t, term, "\033[M")

	// |ABC123__|
	// |DHI756__|
	// |G___89__|
	// |________|

	checkPos(t, term, 0, 1)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "1")
	checkContent(t, term, 4, 0, "2")
	checkContent(t, term, 5, 0, "3")
	checkContent(t, term, 0, 1, "D")
	checkContent(t, term, 1, 1, "H")
	checkContent(t, term, 2, 1, "I")
	checkContent(t, term, 3, 1, "7")
	checkContent(t, term, 4, 1, "5")
	checkContent(t, term, 5, 1, "6")
	checkContent(t, term, 0, 2, "G")
	checkContent(t, term, 1, 2, "")
	checkContent(t, term, 2, 2, "")
	checkContent(t, term, 3, 2, "")
	checkContent(t, term, 4, 2, "8")
	checkContent(t, term, 5, 2, "9")
	checkContent(t, term, 0, 3, "")
	checkContent(t, term, 1, 3, "")
	checkContent(t, term, 2, 3, "")
	checkContent(t, term, 3, 3, "")
	checkContent(t, term, 4, 3, "")
	checkContent(t, term, 5, 3, "")
}

// TestDCv1 tests a simple delete character.
func TestDCv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC123")
	writeF(t, term, "\033[3G")
	writeF(t, term, "\033[2P")

	// |AB23____|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 2, 0)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "2")
	checkContent(t, term, 3, 0, "3")
	checkContent(t, term, 4, 0, "")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 6, 0, "")
}

// TestDCv2 tests delete character SGR state.
func TestDCv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC123")
	writeF(t, term, "\033[3G")
	writeF(t, term, "\033[41m")
	writeF(t, term, "\x1b[1m")
	writeF(t, term, "\033[2P")

	// |AB23____|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 2, 0)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "2")
	checkContent(t, term, 3, 0, "3")
	checkContent(t, term, 4, 0, "")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 6, 0, "")
	checkAttrs(t, term, 6, 0, Plain)
	checkAttrs(t, term, 7, 0, Plain)
	checkColors(t, term, 6, 0, color.Silver, color.Maroon)
	checkColors(t, term, 7, 0, color.Silver, color.Maroon)
}

// TestDCv3 tests delete outside the left/right scroll region.
func TestDCv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC123")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[3;5s")
	writeF(t, term, "\x1b[2G")
	writeF(t, term, "\033[P")

	// |ABC123__|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 1, 0)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "1")
	checkContent(t, term, 4, 0, "2")
	checkContent(t, term, 5, 0, "3")
	checkContent(t, term, 6, 0, "")
}

// TestDCv4 tests delete inside the left/right scroll region.
func TestDCv4(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC123")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[3;5s")
	writeF(t, term, "\x1b[4G")
	writeF(t, term, "\033[P")

	// |ABC2_3__|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 3, 0)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "B")
	checkContent(t, term, 2, 0, "C")
	checkContent(t, term, 3, 0, "2")
	checkContent(t, term, 4, 0, "")
	checkContent(t, term, 5, 0, "3")
	checkContent(t, term, 6, 0, "")
}

// TestDCv5 tests delete character splitting a wide character.
func TestDCv5(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "A橋123")
	writeF(t, term, "\x1b[3G")
	writeF(t, term, "\033[P")

	// |A_123___|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 2, 0)
	checkContent(t, term, 0, 0, "A")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "1")
	checkContent(t, term, 3, 0, "2")
	checkContent(t, term, 4, 0, "3")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 6, 0, "")
}

// TestICHv1 tests insert character without a scroll region, fitting on screen.
func TestICHv1(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "ABC")
	writeF(t, term, "\033[1G")
	writeF(t, term, "\x1b[2@")
	writeF(t, term, "X")

	// |XcABC___|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 1, 0)
	checkContent(t, term, 0, 0, "X")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "A")
	checkContent(t, term, 3, 0, "B")
	checkContent(t, term, 4, 0, "C")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 6, 0, "")
}

// TestICHv2 tests insert character SGR state.
func TestICHv2(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "ABC")
	writeF(t, term, "\033[1G")
	writeF(t, term, "\033[41m\033[7m")
	writeF(t, term, "\x1b[2@")
	writeF(t, term, "X")

	// |XcABC___|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 1, 0)
	checkContent(t, term, 0, 0, "X")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "A")
	checkContent(t, term, 3, 0, "B")
	checkContent(t, term, 4, 0, "C")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 6, 0, "")
	checkAttrs(t, term, 0, 0, Reverse)
	checkAttrs(t, term, 1, 0, Reverse)
	checkAttrs(t, term, 2, 0, Plain)
	checkAttrs(t, term, 3, 0, Plain)
	checkAttrs(t, term, 4, 0, Plain)
	checkColors(t, term, 0, 0, color.Silver, color.Maroon)
	checkColors(t, term, 1, 0, color.Silver, color.Maroon)
	checkColors(t, term, 2, 0, color.Silver, color.Black)
	checkColors(t, term, 3, 0, color.Silver, color.Black)
	checkColors(t, term, 4, 0, color.Silver, color.Black)
}

// TestICHv3 tests insert character shifting off screen
func TestICHv3(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "\033[8G")
	writeF(t, term, "\033[2D")
	writeF(t, term, "ABC")
	writeF(t, term, "\033[2D")
	writeF(t, term, "\033[2@")
	writeF(t, term, "X")

	// |_____XcA|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 6, 0)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 4, 0, "")
	checkContent(t, term, 5, 0, "X")
	checkContent(t, term, 6, 0, "")
	checkContent(t, term, 7, 0, "A")
}

// TestICHv4 tests insert character inside left and right scroll regions.
func TestICHv4(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[3;5s")
	writeF(t, term, "\033[3G")
	writeF(t, term, "ABC")
	writeF(t, term, "\033[3G")
	writeF(t, term, "\033[2@")
	writeF(t, term, "X")

	// |__XcA___|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 3, 0)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "X")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 4, 0, "A")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 6, 0, "")
	checkContent(t, term, 7, 0, "")
}

// TestICHv5 tests insert character outside left and right scroll regions.
func TestICHv5(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[3;5s")
	writeF(t, term, "\033[3G")
	writeF(t, term, "ABC")
	writeF(t, term, "\033[1G")
	writeF(t, term, "\033[2@")
	writeF(t, term, "X")

	// |XcABC___|
	// |________|
	// |________|
	// |________|

	checkPos(t, term, 1, 0)
	checkContent(t, term, 0, 0, "X")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "A")
	checkContent(t, term, 3, 0, "B")
	checkContent(t, term, 4, 0, "C")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 6, 0, "")
	checkContent(t, term, 7, 0, "")
}

// TestICHv5a tests insert character outside left and right scroll regions resets auto-wrap.
func TestICHv5a(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[1;1H\033[0J")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[3;5s")
	writeF(t, term, "\033[3G")
	writeF(t, term, "ABC")
	writeF(t, term, "\033[8G")
	writeF(t, term, "X")
	writeF(t, term, "\033[2@")
	writeF(t, term, "Y")
	writeF(t, term, "Z")

	// |__ABC__Y|
	// |Zc______|
	// |________|
	// |________|

	checkPos(t, term, 1, 1)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "A")
	checkContent(t, term, 3, 0, "B")
	checkContent(t, term, 4, 0, "C")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 6, 0, "")
	checkContent(t, term, 7, 0, "Y")
	checkContent(t, term, 0, 1, "Z")
}

// TestICHv6 tests insert character splitting a wide character on the right of the screen.
func TestICHv6(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[8G")
	writeF(t, term, "\033[1D")
	writeF(t, term, "橋")
	writeF(t, term, "\033[2D")
	writeF(t, term, "\033[@")
	writeF(t, term, "X")

	// |_____Xc_|

	checkPos(t, term, 6, 0)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 4, 0, "")
	checkContent(t, term, 5, 0, "X")
	checkContent(t, term, 6, 0, "")
	checkContent(t, term, 7, 0, "")
}

// TestICHv6a tests insert character splitting a wide character at the right margin.
func TestICHv6a(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[4G")
	writeF(t, term, "橋")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[3;5s")
	writeF(t, term, "\033[3G")
	writeF(t, term, "\033[@")
	writeF(t, term, "X")

	// |__Xc____|

	checkPos(t, term, 3, 0)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "X")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 4, 0, "")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 6, 0, "")
}

// TestICHv6b tests insert character splitting a wide character at the left margin.
func TestICHv6b(t *testing.T) {
	term := NewMockTerm(MockOptSize{X: 8, Y: 5})
	defer mustClose(t, term)
	mustStart(t, term)

	writeF(t, term, "\033[2G")
	writeF(t, term, "橋")
	writeF(t, term, "\033[?69h")
	writeF(t, term, "\033[3;5s")
	writeF(t, term, "\033[3G")
	writeF(t, term, "\033[@")
	writeF(t, term, "X")

	// |__Xc____|

	checkPos(t, term, 3, 0)
	checkContent(t, term, 0, 0, "")
	checkContent(t, term, 1, 0, "")
	checkContent(t, term, 2, 0, "X")
	checkContent(t, term, 3, 0, "")
	checkContent(t, term, 4, 0, "")
	checkContent(t, term, 5, 0, "")
	checkContent(t, term, 6, 0, "")
}
