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
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/vt"
)

// This file implements various tests of the emulator.  Much of these tests
// are "borrowed" (ported from) the tests from Ghostty - https://ghostty.org/docs/vt
// Note that Ghostty's tests assume that STTY modes to expand LF to CF LF are in
// effect (or ANSI mode 20.)  We don't assume that, and add the CR explicitly.

// TestDECSTBMv1 tests full screen scroll region
func TestDECSTBMv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")   //  clear screen
	WriteF(t, term, "ABC\r\n")
	WriteF(t, term, "DEF\r\n")
	WriteF(t, term, "GHI\r\n")
	WriteF(t, term, "\033[r") // scroll region top/bottom
	WriteF(t, term, "\033[T") // scroll down one

	// |c_______|
	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	CheckPos(t, term, 0, 0)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "A")
	CheckContent(t, term, 1, 1, "B")
	CheckContent(t, term, 2, 1, "C")
	CheckContent(t, term, 0, 2, "D")
	CheckContent(t, term, 1, 2, "E")
	CheckContent(t, term, 2, 2, "F")
	CheckContent(t, term, 0, 3, "G")
	CheckContent(t, term, 1, 3, "H")
	CheckContent(t, term, 2, 3, "I")
}

// TestDECSTBMv2 top only
func TestDECSTBMv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")   //  clear screen
	WriteF(t, term, "ABC\r\n")
	WriteF(t, term, "DEF\r\n")
	WriteF(t, term, "GHI\r\n")
	WriteF(t, term, "\033[2;2r") // scroll region top/bottom
	WriteF(t, term, "\033[T")    // scroll down one

	// |________|
	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	CheckPos(t, term, 0, 3) // did not move
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "A")
	CheckContent(t, term, 1, 1, "B")
	CheckContent(t, term, 2, 1, "C")
	CheckContent(t, term, 0, 2, "D")
	CheckContent(t, term, 1, 2, "E")
	CheckContent(t, term, 2, 2, "F")
	CheckContent(t, term, 0, 3, "G")
	CheckContent(t, term, 1, 3, "H")
	CheckContent(t, term, 2, 3, "I")
}

// TestDECSTBMv3 top and bottom
func TestDECSTBMv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")   //  clear screen
	WriteF(t, term, "ABC\r\n")
	WriteF(t, term, "DEF\r\n")
	WriteF(t, term, "GHI\r\n")
	WriteF(t, term, "\033[1;2r") // scroll region top/bottom
	WriteF(t, term, "\033[T")    // scroll down one

	// |________|
	// |ABC_____|
	// |GHI_____|
	// |________|
	CheckPos(t, term, 0, 0)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "A")
	CheckContent(t, term, 1, 1, "B")
	CheckContent(t, term, 2, 1, "C")
	CheckContent(t, term, 0, 2, "G")
	CheckContent(t, term, 1, 2, "H")
	CheckContent(t, term, 2, 2, "I")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestDECSTBMv4 top == bottom
func TestDECSTBMv4(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "ABC\r\n")
	WriteF(t, term, "DEF\r\n")
	WriteF(t, term, "GHI\r\n")
	WriteF(t, term, "\033[2;2r")
	WriteF(t, term, "\033[T")

	// |________|
	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	CheckPos(t, term, 0, 3)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "A")
	CheckContent(t, term, 1, 1, "B")
	CheckContent(t, term, 2, 1, "C")
	CheckContent(t, term, 0, 2, "D")
	CheckContent(t, term, 1, 2, "E")
	CheckContent(t, term, 2, 2, "F")
	CheckContent(t, term, 0, 3, "G")
	CheckContent(t, term, 1, 3, "H")
	CheckContent(t, term, 2, 3, "I")
}

// TestDECSLRMv1 tests full screen right and left margins.
func TestDECSLRMv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[H")
	WriteF(t, term, "\033[J")
	WriteF(t, term, "ABC\n\r")
	WriteF(t, term, "DEF\n\r")
	WriteF(t, term, "GHI\n\r")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[s")
	WriteF(t, term, "\033[X")

	CheckPos(t, term, 0, 0)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 0, 1, "D")
	CheckContent(t, term, 1, 1, "E")
	CheckContent(t, term, 2, 1, "F")
	CheckContent(t, term, 0, 2, "G")
	CheckContent(t, term, 1, 2, "H")
	CheckContent(t, term, 2, 2, "I")
}

// TODO: DECSLRMv3 left and right this makes use of insert line
// TODO: DECSLRMv4 left and right equal
// TODO: add tests for actual left and right scrolling!

// TestRIv1 top of screen, no scroll
func TestRIv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "A\r\n")
	WriteF(t, term, "B\r\n")
	WriteF(t, term, "C\r\n")
	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033M")
	WriteF(t, term, "X")

	// |Xc______|
	// |A_______|
	// |B_______|
	// |C_______|

	CheckPos(t, term, 1, 0)
	CheckContent(t, term, 0, 0, "X")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "A")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "B")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "C")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestRIv2 not top of screen, no scroll
func TestRIv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "A\r\n")
	WriteF(t, term, "B\r\n")
	WriteF(t, term, "C\r\n")
	WriteF(t, term, "\033[2;1H")
	WriteF(t, term, "\033M")
	WriteF(t, term, "X")

	// |Xc______|
	// |B_______|
	// |C_______|
	// |________|

	CheckPos(t, term, 1, 0)
	CheckContent(t, term, 0, 0, "X")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "B")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "C")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestRIv3 scroll region
func TestRIv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")   //  clear screen
	WriteF(t, term, "A\r\n")
	WriteF(t, term, "B\r\n")
	WriteF(t, term, "C\r\n")
	WriteF(t, term, "\033[2;3r")
	WriteF(t, term, "\033[2;1H")
	WriteF(t, term, "\033M")

	// |A_______|
	// |c_______|
	// |B_______|
	// |________|

	CheckPos(t, term, 0, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "B")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestRIv4 outside scroll region - goes to top, does not scroll
func TestRIv4(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")   //  clear screen
	WriteF(t, term, "A\r\n")
	WriteF(t, term, "B\r\n")
	WriteF(t, term, "C\r\n")
	WriteF(t, term, "\033[2;3r")
	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033M")

	// |A_______|
	// |B_______|
	// |C_______|
	// |________|

	CheckPos(t, term, 0, 0)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "B")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "C")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TODO: RIv5 - left right scroll regions (when we implement left/right regions)
// TODO: RIv6 - outside left/right scroll regions (when we implement left/right regions)

// TestINDv1 no scroll region, top of screen
func TestINDv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")   //  clear screen
	WriteF(t, term, "A")
	WriteF(t, term, "\033D")
	WriteF(t, term, "X")

	// |A_______|
	// |_Xc_____|
	// |________|
	// |________|

	CheckPos(t, term, 2, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "X")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestINDv2 no scroll region, bottom of screen
func TestINDv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")   //  clear screen
	WriteF(t, term, "\033[4;1H")
	WriteF(t, term, "A")
	WriteF(t, term, "\033D")
	WriteF(t, term, "X")

	// |________|
	// |________|
	// |A_______|
	// |_Xc_____|

	CheckPos(t, term, 2, 3)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "A")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "X")
	CheckContent(t, term, 2, 3, "")
}

// TestINDv3 inside scroll region
func TestINDv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "A")
	WriteF(t, term, "\033D")
	WriteF(t, term, "X")

	// |A_______|
	// |_Xc_____|
	// |________|
	// |________|

	CheckPos(t, term, 2, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "X")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestINDv4 bottom of scroll region
func TestINDv4(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "\033[4;1H")
	WriteF(t, term, "B")
	WriteF(t, term, "\033[3;1H")
	WriteF(t, term, "A")
	WriteF(t, term, "\033D")
	WriteF(t, term, "X")

	// |________|
	// |A_______|
	// |_Xc_____|
	// |B_______|

	CheckPos(t, term, 2, 2)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "A")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "X")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "B")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestINDv5 bottom of screen with scroll region
func TestINDv5(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "\033[3;1H")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[4;1H")
	WriteF(t, term, "\033D")
	WriteF(t, term, "X")

	// |________|
	// |________|
	// |A_______|
	// |________|
	// |Xc______|

	CheckPos(t, term, 1, 4)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "A")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
	CheckContent(t, term, 0, 4, "X")
	CheckContent(t, term, 1, 4, "")
	CheckContent(t, term, 2, 4, "")
}

// TestINDv6 tests IND outside of left and right scroll region.
func TestINDv6(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "\033[3;5s")
	WriteF(t, term, "\033[3;3H")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[3;1H")
	WriteF(t, term, "\033D")
	WriteF(t, term, "X")

	// |________|
	// |________|
	// |XcA_____|

	CheckPos(t, term, 1, 2)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "X")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "A")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestINDv7 tests IND inside of left and right scroll region.
func TestINDv7(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H") // move to top-left
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "111111\n\r")
	WriteF(t, term, "222222\n\r")
	WriteF(t, term, "333333\n\r")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[1;3s")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "\033[3;1H")
	WriteF(t, term, "\033D")

	// |222111__|
	// |333222__|
	// |c__333__|

	CheckPos(t, term, 0, 2)
	CheckContent(t, term, 0, 0, "2")
	CheckContent(t, term, 1, 0, "2")
	CheckContent(t, term, 2, 0, "2")
	CheckContent(t, term, 3, 0, "1")
	CheckContent(t, term, 4, 0, "1")
	CheckContent(t, term, 5, 0, "1")
	CheckContent(t, term, 0, 1, "3")
	CheckContent(t, term, 1, 1, "3")
	CheckContent(t, term, 2, 1, "3")
	CheckContent(t, term, 3, 1, "2")
	CheckContent(t, term, 4, 1, "2")
	CheckContent(t, term, 5, 1, "2")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 3, 2, "3")
	CheckContent(t, term, 4, 2, "3")
	CheckContent(t, term, 5, 2, "3")
}

// TestCUDv1 - cursor down
func TestCUDv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "A")
	WriteF(t, term, "\033[2B")
	WriteF(t, term, "X")

	// |A_______|
	// |________|
	// |_Xc_____|
	// |________|

	CheckPos(t, term, 2, 2)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "X")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestCUDv2 - cursor down above bottom margin
func TestCUDv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\n\n\n\n")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[5B")
	WriteF(t, term, "X")

	// |A_______|
	// |________|
	// |_Xc_____|
	// |________|

	CheckPos(t, term, 2, 2)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "X")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestCUDv3 - cursor down below bottom margin
func TestCUDv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[4;1H")
	WriteF(t, term, "\033[5B")
	WriteF(t, term, "X")

	// |A_______|
	// |________|
	// |________|
	// |________|
	// |Xc______|

	CheckPos(t, term, 1, 4)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
	CheckContent(t, term, 0, 4, "X")
	CheckContent(t, term, 1, 4, "")
	CheckContent(t, term, 2, 4, "")
}

// TestCUUv1 tests cursor up.
func TestCUUv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[3;H")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[2A")
	WriteF(t, term, "X")

	// |_Xc_____|
	// |________|
	// |A_______|

	CheckPos(t, term, 2, 0)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "X")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "A")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
}

// TestCUUv2 tests cursor up below the top margin.
func TestCUUv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[2;4r")
	WriteF(t, term, "\033[3;1H")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[5A")
	WriteF(t, term, "X")

	// |________|
	// |_Xc_____|
	// |A_______|
	// |________|

	CheckPos(t, term, 2, 1)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "X")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "A")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestCUUv3 tests cursor up above the top margin.
func TestCUUv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[3;5r")
	WriteF(t, term, "\033[3;1H")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[2;1H")
	WriteF(t, term, "\033[5A")
	WriteF(t, term, "X")

	// |Xc______|
	// |________|
	// |A_______|
	// |________|
	// |________|

	CheckPos(t, term, 1, 0)
	CheckContent(t, term, 0, 0, "X")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "A")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestCNLv1 - cursor next line
func TestCNLv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "A")
	WriteF(t, term, "\033[2E")
	WriteF(t, term, "X")

	// |A_______|
	// |________|
	// |Xc_____|
	// |________|

	CheckPos(t, term, 1, 2)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "X")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestCNLv2 - cursor next line above bottom margin
func TestCNLv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\n\n\n\n")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[5E")
	WriteF(t, term, "X")

	// |A_______|
	// |________|
	// |Xc______|
	// |________|

	CheckPos(t, term, 1, 2)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "X")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestCNLv3 - cursor next line bottom margin
func TestCNLv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[4;3H")
	WriteF(t, term, "\033[5E")
	WriteF(t, term, "X")

	// |A_______|
	// |________|
	// |________|
	// |________|
	// |Xc______|

	CheckPos(t, term, 1, 4)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
	CheckContent(t, term, 0, 4, "X")
	CheckContent(t, term, 1, 4, "")
	CheckContent(t, term, 2, 4, "")
}

// TestCPLv1 tests cursor previous line.
func TestCPLv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[3;H")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[2F")
	WriteF(t, term, "X")

	// |Xc______|
	// |________|
	// |A_______|

	CheckPos(t, term, 1, 0)
	CheckContent(t, term, 0, 0, "X")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "A")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
}

// TestCPLv2 tests cursor previous line below the top margin.
func TestCPLv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[2;4r")
	WriteF(t, term, "\033[3;1H")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[5F")
	WriteF(t, term, "X")

	// |________|
	// |Xc______|
	// |A_______|
	// |________|

	CheckPos(t, term, 1, 1)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "X")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "A")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestCPLv3 tests cursor previous line above the top margin.
func TestCPLv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H")
	WriteF(t, term, "\033[0J")
	WriteF(t, term, "\033[3;5r")
	WriteF(t, term, "\033[3;1H")
	WriteF(t, term, "A")
	WriteF(t, term, "\033[2;2H")
	WriteF(t, term, "\033[5F")
	WriteF(t, term, "X")

	// |Xc______|
	// |________|
	// |A_______|
	// |________|
	// |________|

	CheckPos(t, term, 1, 0)
	CheckContent(t, term, 0, 0, "X")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 0, 2, "A")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
}

// TestILv1 tests inserting a line.
func TestILv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC\n\r")
	WriteF(t, term, "DEF\n\r")
	WriteF(t, term, "GHI\n\r")
	WriteF(t, term, "\033[2;2H")
	WriteF(t, term, "\033[L")

	// |ABC_____|
	// |c_______|
	// |DEF_____|
	// |GHI_____|
	// |________|

	CheckPos(t, term, 0, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 3, 1, "")
	CheckContent(t, term, 0, 2, "D")
	CheckContent(t, term, 1, 2, "E")
	CheckContent(t, term, 2, 2, "F")
	CheckContent(t, term, 3, 2, "")
	CheckContent(t, term, 0, 3, "G")
	CheckContent(t, term, 1, 3, "H")
	CheckContent(t, term, 2, 3, "I")
	CheckContent(t, term, 3, 3, "")
}

// TestILv2 tests insert line outside of the scroll region.
func TestILv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC\n\r")
	WriteF(t, term, "DEF\n\r")
	WriteF(t, term, "GHI\n\r")
	WriteF(t, term, "\033[3;4r")
	WriteF(t, term, "\033[2;2H")
	WriteF(t, term, "\033[L")

	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	// |________|
	// |________|

	CheckPos(t, term, 1, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 0, 1, "D")
	CheckContent(t, term, 1, 1, "E")
	CheckContent(t, term, 2, 1, "F")
	CheckContent(t, term, 3, 1, "")
	CheckContent(t, term, 0, 2, "G")
	CheckContent(t, term, 1, 2, "H")
	CheckContent(t, term, 2, 2, "I")
	CheckContent(t, term, 3, 2, "")
}

// TestILv3 tests insert line with top and bottom scroll regions.
func TestILv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC\n\r")
	WriteF(t, term, "DEF\n\r")
	WriteF(t, term, "GHI\n\r")
	WriteF(t, term, "123\n\r")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "\033[2;2H")
	WriteF(t, term, "\033[L")

	// |ABC_____|
	// |c_______|
	// |DEF_____|
	// |123_____|

	CheckPos(t, term, 0, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 0, 1, "")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 3, 1, "")
	CheckContent(t, term, 0, 2, "D")
	CheckContent(t, term, 1, 2, "E")
	CheckContent(t, term, 2, 2, "F")
	CheckContent(t, term, 3, 2, "")
	CheckContent(t, term, 0, 3, "1")
	CheckContent(t, term, 1, 3, "2")
	CheckContent(t, term, 2, 3, "3")
	CheckContent(t, term, 3, 3, "")
}

// TestILv4 tests insert line with left and right margins.
func TestILv4(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC123\n\r")
	WriteF(t, term, "DEF456\n\r")
	WriteF(t, term, "GHI789\n\r")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[2;4s")
	WriteF(t, term, "\033[2;2H")
	WriteF(t, term, "\033[L")

	// |ABC123__|
	// |Dc__56__|
	// |GEF489__|
	// |_HI7____|

	CheckPos(t, term, 0, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "1")
	CheckContent(t, term, 4, 0, "2")
	CheckContent(t, term, 5, 0, "3")
	CheckContent(t, term, 0, 1, "D")
	CheckContent(t, term, 1, 1, "")
	CheckContent(t, term, 2, 1, "")
	CheckContent(t, term, 3, 1, "")
	CheckContent(t, term, 4, 1, "5")
	CheckContent(t, term, 5, 1, "6")
	CheckContent(t, term, 0, 2, "G")
	CheckContent(t, term, 1, 2, "E")
	CheckContent(t, term, 2, 2, "F")
	CheckContent(t, term, 3, 2, "4")
	CheckContent(t, term, 4, 2, "8")
	CheckContent(t, term, 5, 2, "9")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "H")
	CheckContent(t, term, 2, 3, "I")
	CheckContent(t, term, 3, 3, "7")
	CheckContent(t, term, 4, 3, "")
	CheckContent(t, term, 5, 3, "")
}

// TestDLv1 tests a simple delete line.
func TestDLv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC\n\r")
	WriteF(t, term, "DEF\n\r")
	WriteF(t, term, "GHI\n\r")
	WriteF(t, term, "\033[2;2H")
	WriteF(t, term, "\033[M")

	// |ABC_____|
	// |GHI_____|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 0, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 0, 1, "G")
	CheckContent(t, term, 1, 1, "H")
	CheckContent(t, term, 2, 1, "I")
	CheckContent(t, term, 3, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 2, 3, "")
}

// TestDLv2 delete line outside of the scroll region.
func TestDLv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC\n\r")
	WriteF(t, term, "DEF\n\r")
	WriteF(t, term, "GHI\n\r")
	WriteF(t, term, "\033[3;4r")
	WriteF(t, term, "\033[2;2H")
	WriteF(t, term, "\033[M")

	// |ABC_____|
	// |DEF_____|
	// |GHI_____|
	// |________|
	// |________|

	CheckPos(t, term, 1, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 0, 1, "D")
	CheckContent(t, term, 1, 1, "E")
	CheckContent(t, term, 2, 1, "F")
	CheckContent(t, term, 3, 1, "")
	CheckContent(t, term, 0, 2, "G")
	CheckContent(t, term, 1, 2, "H")
	CheckContent(t, term, 2, 2, "I")
	CheckContent(t, term, 3, 2, "")
}

// TestDLv3 delete line with top and bottom scroll regions.
func TestDLv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC\n\r")
	WriteF(t, term, "DEF\n\r")
	WriteF(t, term, "GHI\n\r")
	WriteF(t, term, "123\n\r")
	WriteF(t, term, "\033[1;3r")
	WriteF(t, term, "\033[2;2H")
	WriteF(t, term, "\033[M")

	// |ABC_____|
	// |GHI_____|
	// |________|
	// |123_____|

	CheckPos(t, term, 0, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 0, 1, "G")
	CheckContent(t, term, 1, 1, "H")
	CheckContent(t, term, 2, 1, "I")
	CheckContent(t, term, 3, 1, "")
	CheckContent(t, term, 0, 2, "")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 3, 2, "")
	CheckContent(t, term, 0, 3, "1")
	CheckContent(t, term, 1, 3, "2")
	CheckContent(t, term, 2, 3, "3")
	CheckContent(t, term, 3, 3, "")
}

// TestDLv4 delete line with left and right margins.
func TestDLv4(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC123\n\r")
	WriteF(t, term, "DEF456\n\r")
	WriteF(t, term, "GHI789\n\r")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[2;4s")
	WriteF(t, term, "\033[2;2H")
	WriteF(t, term, "\033[M")

	// |ABC123__|
	// |DHI756__|
	// |G___89__|
	// |________|

	CheckPos(t, term, 0, 1)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "1")
	CheckContent(t, term, 4, 0, "2")
	CheckContent(t, term, 5, 0, "3")
	CheckContent(t, term, 0, 1, "D")
	CheckContent(t, term, 1, 1, "H")
	CheckContent(t, term, 2, 1, "I")
	CheckContent(t, term, 3, 1, "7")
	CheckContent(t, term, 4, 1, "5")
	CheckContent(t, term, 5, 1, "6")
	CheckContent(t, term, 0, 2, "G")
	CheckContent(t, term, 1, 2, "")
	CheckContent(t, term, 2, 2, "")
	CheckContent(t, term, 3, 2, "")
	CheckContent(t, term, 4, 2, "8")
	CheckContent(t, term, 5, 2, "9")
	CheckContent(t, term, 0, 3, "")
	CheckContent(t, term, 1, 3, "")
	CheckContent(t, term, 2, 3, "")
	CheckContent(t, term, 3, 3, "")
	CheckContent(t, term, 4, 3, "")
	CheckContent(t, term, 5, 3, "")
}

// TestDCv1 tests a simple delete character.
func TestDCv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC123")
	WriteF(t, term, "\033[3G")
	WriteF(t, term, "\033[2P")

	// |AB23____|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 2, 0)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "2")
	CheckContent(t, term, 3, 0, "3")
	CheckContent(t, term, 4, 0, "")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 6, 0, "")
}

// TestDCv2 tests delete character SGR state.
func TestDCv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC123")
	WriteF(t, term, "\033[3G")
	WriteF(t, term, "\033[41m")
	WriteF(t, term, "\x1b[1m")
	WriteF(t, term, "\033[2P")

	// |AB23____|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 2, 0)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "2")
	CheckContent(t, term, 3, 0, "3")
	CheckContent(t, term, 4, 0, "")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 6, 0, "")
	CheckAttrs(t, term, 6, 0, vt.Plain)
	CheckAttrs(t, term, 7, 0, vt.Plain)
	CheckColors(t, term, 6, 0, color.Silver, color.Maroon)
	CheckColors(t, term, 7, 0, color.Silver, color.Maroon)
}

// TestDCv3 tests delete outside the left/right scroll region.
func TestDCv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC123")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[3;5s")
	WriteF(t, term, "\x1b[2G")
	WriteF(t, term, "\033[P")

	// |ABC123__|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 1, 0)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "1")
	CheckContent(t, term, 4, 0, "2")
	CheckContent(t, term, 5, 0, "3")
	CheckContent(t, term, 6, 0, "")
}

// TestDCv4 tests delete inside the left/right scroll region.
func TestDCv4(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC123")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[3;5s")
	WriteF(t, term, "\x1b[4G")
	WriteF(t, term, "\033[P")

	// |ABC2_3__|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 3, 0)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "B")
	CheckContent(t, term, 2, 0, "C")
	CheckContent(t, term, 3, 0, "2")
	CheckContent(t, term, 4, 0, "")
	CheckContent(t, term, 5, 0, "3")
	CheckContent(t, term, 6, 0, "")
}

// TestDCv5 tests delete character splitting a wide character.
func TestDCv5(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "A橋123")
	WriteF(t, term, "\x1b[3G")
	WriteF(t, term, "\033[P")

	// |A_123___|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 2, 0)
	CheckContent(t, term, 0, 0, "A")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "1")
	CheckContent(t, term, 3, 0, "2")
	CheckContent(t, term, 4, 0, "3")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 6, 0, "")
}

// TestICHv1 tests insert character without a scroll region, fitting on screen.
func TestICHv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "ABC")
	WriteF(t, term, "\033[1G")
	WriteF(t, term, "\x1b[2@")
	WriteF(t, term, "X")

	// |XcABC___|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 1, 0)
	CheckContent(t, term, 0, 0, "X")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "A")
	CheckContent(t, term, 3, 0, "B")
	CheckContent(t, term, 4, 0, "C")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 6, 0, "")
}

// TestICHv2 tests insert character SGR state.
func TestICHv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "ABC")
	WriteF(t, term, "\033[1G")
	WriteF(t, term, "\033[41m\033[7m")
	WriteF(t, term, "\x1b[2@")
	WriteF(t, term, "X")

	// |XcABC___|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 1, 0)
	CheckContent(t, term, 0, 0, "X")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "A")
	CheckContent(t, term, 3, 0, "B")
	CheckContent(t, term, 4, 0, "C")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 6, 0, "")
	CheckAttrs(t, term, 0, 0, vt.Reverse)
	CheckAttrs(t, term, 1, 0, vt.Reverse)
	CheckAttrs(t, term, 2, 0, vt.Plain)
	CheckAttrs(t, term, 3, 0, vt.Plain)
	CheckAttrs(t, term, 4, 0, vt.Plain)
	CheckColors(t, term, 0, 0, color.Silver, color.Maroon)
	CheckColors(t, term, 1, 0, color.Silver, color.Maroon)
	CheckColors(t, term, 2, 0, color.Silver, color.Black)
	CheckColors(t, term, 3, 0, color.Silver, color.Black)
	CheckColors(t, term, 4, 0, color.Silver, color.Black)
}

// TestICHv3 tests insert character shifting off screen
func TestICHv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "\033[8G")
	WriteF(t, term, "\033[2D")
	WriteF(t, term, "ABC")
	WriteF(t, term, "\033[2D")
	WriteF(t, term, "\033[2@")
	WriteF(t, term, "X")

	// |_____XcA|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 6, 0)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 4, 0, "")
	CheckContent(t, term, 5, 0, "X")
	CheckContent(t, term, 6, 0, "")
	CheckContent(t, term, 7, 0, "A")
}

// TestICHv4 tests insert character inside left and right scroll regions.
func TestICHv4(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[3;5s")
	WriteF(t, term, "\033[3G")
	WriteF(t, term, "ABC")
	WriteF(t, term, "\033[3G")
	WriteF(t, term, "\033[2@")
	WriteF(t, term, "X")

	// |__XcA___|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 3, 0)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "X")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 4, 0, "A")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 6, 0, "")
	CheckContent(t, term, 7, 0, "")
}

// TestICHv5 tests insert character outside left and right scroll regions.
func TestICHv5(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[3;5s")
	WriteF(t, term, "\033[3G")
	WriteF(t, term, "ABC")
	WriteF(t, term, "\033[1G")
	WriteF(t, term, "\033[2@")
	WriteF(t, term, "X")

	// |XcABC___|
	// |________|
	// |________|
	// |________|

	CheckPos(t, term, 1, 0)
	CheckContent(t, term, 0, 0, "X")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "A")
	CheckContent(t, term, 3, 0, "B")
	CheckContent(t, term, 4, 0, "C")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 6, 0, "")
	CheckContent(t, term, 7, 0, "")
}

// TestICHv5a tests insert character outside left and right scroll regions resets auto-wrap.
func TestICHv5a(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[1;1H\033[0J")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[3;5s")
	WriteF(t, term, "\033[3G")
	WriteF(t, term, "ABC")
	WriteF(t, term, "\033[8G")
	WriteF(t, term, "X")
	WriteF(t, term, "\033[2@")
	WriteF(t, term, "Y")
	WriteF(t, term, "Z")

	// |__ABC__Y|
	// |Zc______|
	// |________|
	// |________|

	CheckPos(t, term, 1, 1)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "A")
	CheckContent(t, term, 3, 0, "B")
	CheckContent(t, term, 4, 0, "C")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 6, 0, "")
	CheckContent(t, term, 7, 0, "Y")
	CheckContent(t, term, 0, 1, "Z")
}

// TestICHv6 tests insert character splitting a wide character on the right of the screen.
func TestICHv6(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[8G")
	WriteF(t, term, "\033[1D")
	WriteF(t, term, "橋")
	WriteF(t, term, "\033[2D")
	WriteF(t, term, "\033[@")
	WriteF(t, term, "X")

	// |_____Xc_|

	CheckPos(t, term, 6, 0)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 4, 0, "")
	CheckContent(t, term, 5, 0, "X")
	CheckContent(t, term, 6, 0, "")
	CheckContent(t, term, 7, 0, "")
}

// TestICHv6a tests insert character splitting a wide character at the right margin.
func TestICHv6a(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[4G")
	WriteF(t, term, "橋")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[3;5s")
	WriteF(t, term, "\033[3G")
	WriteF(t, term, "\033[@")
	WriteF(t, term, "X")

	// |__Xc____|

	CheckPos(t, term, 3, 0)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "X")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 4, 0, "")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 6, 0, "")
}

// TestICHv6b tests insert character splitting a wide character at the left margin.
func TestICHv6b(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[2G")
	WriteF(t, term, "橋")
	WriteF(t, term, "\033[?69h")
	WriteF(t, term, "\033[3;5s")
	WriteF(t, term, "\033[3G")
	WriteF(t, term, "\033[@")
	WriteF(t, term, "X")

	// |__Xc____|

	CheckPos(t, term, 3, 0)
	CheckContent(t, term, 0, 0, "")
	CheckContent(t, term, 1, 0, "")
	CheckContent(t, term, 2, 0, "X")
	CheckContent(t, term, 3, 0, "")
	CheckContent(t, term, 4, 0, "")
	CheckContent(t, term, 5, 0, "")
	CheckContent(t, term, 6, 0, "")
}

// TestDECSCUSRv1 tests setting the cursor.
func TestDECSCUSRv1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[ q")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingBlock, "")
	WriteF(t, term, "\033[0 q")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingBlock, "")
	WriteF(t, term, "\033[1 q")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingBlock, "")
	WriteF(t, term, "\033[2 q")
	VerifyF(t, term.Backend().GetCursor() == vt.SteadyBlock, "")
	WriteF(t, term, "\033[3 q")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingUnderline, "")
	WriteF(t, term, "\033[4 q")
	VerifyF(t, term.Backend().GetCursor() == vt.SteadyUnderline, "")
	WriteF(t, term, "\033[5 q")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingBar, "")
	WriteF(t, term, "\033[6 q")
	VerifyF(t, term.Backend().GetCursor() == vt.SteadyBar, "")
}

// TestDECSCUSRv2 tests setting the cursor while hidden.
func TestDECSCUSRv2(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[?25l")
	WriteF(t, term, "\033[ q")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingBlock.Hide(), "")
	WriteF(t, term, "\033[0 q")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingBlock.Hide(), "")
	WriteF(t, term, "\033[1 q")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingBlock.Hide(), "")
	WriteF(t, term, "\033[2 q")
	VerifyF(t, term.Backend().GetCursor() == vt.SteadyBlock.Hide(), "")
	WriteF(t, term, "\033[3 q")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingUnderline.Hide(), "")
	WriteF(t, term, "\033[4 q")
	VerifyF(t, term.Backend().GetCursor() == vt.SteadyUnderline.Hide(), "")
	WriteF(t, term, "\033[5 q")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingBar.Hide(), "")
	WriteF(t, term, "\033[6 q")
	VerifyF(t, term.Backend().GetCursor() == vt.SteadyBar.Hide(), "")
}

// TestDECSCUSRv3 tests setting the blink private mode.
func TestDECSCUSRv3(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033[0 q")
	WriteF(t, term, "\033[?12l")
	VerifyF(t, term.Backend().GetCursor() == vt.SteadyBlock, "")
	WriteF(t, term, "\033[?12h")
	VerifyF(t, term.Backend().GetCursor() == vt.BlinkingBlock, "")
}

// TestOSC8v1 tests hyperlinks.
func TestOSC8v1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	WriteF(t, term, "\033]8;bob=junk:id=123:junk;https://example.com\033\\")
	WriteF(t, term, "Hello")
	WriteF(t, term, "\033]8;;\033\\")
	WriteF(t, term, " World!")

	term.Drain()
	for col := range Col(10) {
		cell := term.GetCell(Coord{X: col, Y: 0})
		if col < 5 {
			if u, i := cell.S.Url(); u != "https://example.com" || i != "123" {
				t.Errorf("wrong data at column %d: url %q id %q", col, u, i)
			}
		} else {
			if u, i := cell.S.Url(); u != "" || i != "" {
				t.Errorf("wrong data at column %d: url %q id %q", col, u, i)
			}
		}
	}
}

// TestOSC52v1 tests pasting and retrieving the clipboard.
func TestOSC52v1(t *testing.T) {
	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 5})
	defer MustClose(t, term)
	MustStart(t, term)

	str := "Hello, World."
	enc := make([]byte, base64.StdEncoding.EncodedLen(len(str)))
	base64.StdEncoding.Encode(enc, []byte(str))

	WriteF(t, term, "\033]52;c;%s\033\\", enc)
	VerifyF(t, string(term.Backend().GetClipboard()) == str, "clipboard contents did not match (%s != %s)",
		string(term.Backend().GetClipboard()), str)
	WriteF(t, term, "\033]52;c;?\033\\")

	response := ReadF(t, term)
	expect := fmt.Sprintf("\033]52;c;%s\033\\", enc)
	VerifyF(t, response == expect, "retrieved value did not match (%q != %q)", response, expect)

	// any invalid base64 clears the clipboard
	WriteF(t, term, "\033]52;c;junk1\033\\")
	VerifyF(t, string(term.Backend().GetClipboard()) == "", "")

	// empty clears the clipboard
	term.Backend().SetClipboard([]byte("something"))
	WriteF(t, term, "\033]52;c;\033\\")
	VerifyF(t, string(term.Backend().GetClipboard()) == "", "")

	// but a bogus sequence without two fields does nothing
	term.Backend().SetClipboard([]byte("something"))
	WriteF(t, term, "\033]52;x\033\\")
	VerifyF(t, string(term.Backend().GetClipboard()) == "something", "")
}
