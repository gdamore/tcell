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
	"testing"
	"unicode/utf8"

	"github.com/gdamore/tcell/v3/vt"
)

func TestSetContent(t *testing.T) {
	_, scr := NewMockScreen(t, vt.MockOptSize{X: 80, Y: 24})
	defer scr.Fini()

	style := StyleDefault.Foreground(ColorRed).Background(ColorBlue)

	t.Run("ascii", func(t *testing.T) {
		scr.SetContent(0, 0, 'H', nil, style)
		got, gotStyle, gotWidth := scr.Get(0, 0)
		if got != "H" {
			t.Errorf("got %q, want %q", got, "H")
		}
		if gotStyle != style {
			t.Errorf("got style %v, want %v", gotStyle, style)
		}
		if gotWidth != 1 {
			t.Errorf("got width %d, want 1", gotWidth)
		}
	})

	t.Run("ascii_empty_combining_slice", func(t *testing.T) {
		scr.SetContent(1, 0, 'i', []rune{}, style)
		got, _, gotWidth := scr.Get(1, 0)
		if got != "i" {
			t.Errorf("got %q, want %q", got, "i")
		}
		if gotWidth != 1 {
			t.Errorf("got width %d, want 1", gotWidth)
		}
	})

	t.Run("non_ascii", func(t *testing.T) {
		scr.SetContent(2, 0, '世', nil, style)
		got, _, gotWidth := scr.Get(2, 0)
		if got != "世" {
			t.Errorf("got %q, want %q", got, "世")
		}
		if gotWidth != 2 {
			t.Errorf("got width %d, want 2", gotWidth)
		}
	})

	t.Run("combining", func(t *testing.T) {
		scr.SetContent(4, 0, 'e', []rune{'\u0301'}, style)
		got, _, gotWidth := scr.Get(4, 0)
		if got != "e\u0301" {
			t.Errorf("got %q, want %q", got, "e\u0301")
		}
		if gotWidth != 1 {
			t.Errorf("got width %d, want 1", gotWidth)
		}
	})

	t.Run("out_of_bounds", func(t *testing.T) {
		// Should not panic on negative or beyond screen bounds
		scr.SetContent(-1, 0, 'x', nil, style)
		scr.SetContent(0, -1, 'x', nil, style)
		scr.SetContent(100, 100, 'x', nil, style)
	})

	t.Run("negative_rune", func(t *testing.T) {
		scr.SetContent(5, 0, rune(-1), nil, style)
		got, _, _ := scr.Get(5, 0)
		if got != string(utf8.RuneError) {
			t.Errorf("got %q, want RuneError", got)
		}
	})
}

func TestSetContentASCIIAllocations(t *testing.T) {
	_, scr := NewMockScreen(t, vt.MockOptSize{X: 80, Y: 24})
	defer scr.Fini()

	style := StyleDefault

	t.Run("redraw_zero_allocs", func(t *testing.T) {
		scr.SetContent(0, 0, 'A', nil, style)
		allocs := testing.AllocsPerRun(100, func() {
			scr.SetContent(0, 0, 'A', nil, style)
		})
		if allocs != 0 {
			t.Errorf("SetContent identical redraw allocated %v times, want 0", allocs)
		}
	})

	t.Run("changed_ascii_zero_allocs", func(t *testing.T) {
		var toggle bool
		allocs := testing.AllocsPerRun(100, func() {
			if toggle {
				scr.SetContent(0, 0, 'A', nil, style)
			} else {
				scr.SetContent(0, 0, 'B', nil, style)
			}
			toggle = !toggle
		})
		if allocs != 0 {
			t.Errorf("SetContent changed ASCII allocated %v times, want 0", allocs)
		}
	})

	t.Run("empty_slice_zero_allocs", func(t *testing.T) {
		emptyComb := []rune{}
		allocs := testing.AllocsPerRun(100, func() {
			scr.SetContent(0, 0, 'C', emptyComb, style)
		})
		if allocs != 0 {
			t.Errorf("SetContent empty combining slice allocated %v times, want 0", allocs)
		}
	})

	t.Run("all_printable_ascii_zero_allocs", func(t *testing.T) {
		r := rune(' ')
		allocs := testing.AllocsPerRun(100, func() {
			scr.SetContent(0, 0, r, nil, style)
			r++
			if r > '~' {
				r = ' '
			}
		})
		if allocs != 0 {
			t.Errorf("SetContent printable ASCII range allocated %v times, want 0", allocs)
		}
	})
}
