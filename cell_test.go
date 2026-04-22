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

import "testing"

func TestCellBufferPutAndSanitize(t *testing.T) {
	t.Run("unsanitized", func(t *testing.T) {
		cb := &CellBuffer{w: 4, h: 1, cells: make([]cell, 4)}
		rest, width := cb.Put(0, 0, "A\x1bB", StyleDefault)
		if rest != "\x1bB" {
			t.Fatalf("unexpected remainder: %q", rest)
		}
		if width != 1 {
			t.Fatalf("unexpected width: %d", width)
		}
		if got, _, gotWidth := cb.Get(0, 0); got != "A" || gotWidth != 1 {
			t.Fatalf("unexpected cell content: %q width=%d", got, gotWidth)
		}
	})

	t.Run("sanitized", func(t *testing.T) {
		cb := &CellBuffer{w: 4, h: 1, cells: make([]cell, 4), sanitizeContent: true}
		rest, width := cb.Put(0, 0, "A\x1bB", StyleDefault)
		if rest != "B" {
			t.Fatalf("unexpected sanitized remainder: %q", rest)
		}
		if width != 1 {
			t.Fatalf("unexpected width: %d", width)
		}
		if got, _, gotWidth := cb.Get(0, 0); got != "A" || gotWidth != 1 {
			t.Fatalf("unexpected sanitized cell content: %q width=%d", got, gotWidth)
		}
	})
}

func TestCellBufferWideAndColorNone(t *testing.T) {
	cb := &CellBuffer{w: 3, h: 1, cells: make([]cell, 3)}

	base := StyleDefault.Foreground(ColorRed).Background(ColorBlue)
	cb.Put(0, 0, "a", base)
	cb.SetDirty(0, 0, false)
	cb.Put(1, 0, "z", base)
	cb.SetDirty(1, 0, false)

	reused := Style{fg: ColorNone, bg: ColorNone}
	rest, width := cb.Put(0, 0, "你", reused)
	if rest != "" {
		t.Fatalf("unexpected remainder for wide rune: %q", rest)
	}
	if width != 2 {
		t.Fatalf("unexpected width for wide rune: %d", width)
	}
	if !cb.Dirty(0, 0) {
		t.Fatalf("wide rune should dirty the base cell")
	}
	if !cb.Dirty(1, 0) {
		t.Fatalf("wide rune should dirty the second cell")
	}

	got, gotStyle, gotWidth := cb.Get(0, 0)
	if got != "你" || gotWidth != 2 {
		t.Fatalf("unexpected wide cell content: %q width=%d", got, gotWidth)
	}
	if gotStyle.GetForeground() != ColorRed || gotStyle.GetBackground() != ColorBlue {
		t.Fatalf("ColorNone should preserve existing colors, got fg=%v bg=%v", gotStyle.GetForeground(), gotStyle.GetBackground())
	}
}

func TestCellBufferDirtyState(t *testing.T) {
	cb := &CellBuffer{w: 2, h: 1, cells: make([]cell, 2)}

	cb.cells[0].setDirty(false)
	if got := cb.cells[0].lastStr; got != " " {
		t.Fatalf("setDirty(false) should normalize empty content to space, got %q", got)
	}
	if got := cb.cells[0].lastStyle; got != cb.cells[0].currStyle {
		t.Fatalf("setDirty(false) should copy style")
	}

	cb.cells[0].setDirty(true)
	if got := cb.cells[0].lastStr; got != "" {
		t.Fatalf("setDirty(true) should clear lastStr, got %q", got)
	}

	cb.Put(0, 0, "x", StyleDefault)
	cb.SetDirty(0, 0, false)
	if cb.Dirty(0, 0) {
		t.Fatalf("clean cell should not be dirty")
	}

	cb.cells[0].currStyle = cb.cells[0].currStyle.Bold(true)
	if !cb.Dirty(0, 0) {
		t.Fatalf("style change should make cell dirty")
	}
	cb.SetDirty(0, 0, false)
	cb.cells[0].currStyle = cb.cells[0].lastStyle
	cb.cells[0].currStr = "y"
	if !cb.Dirty(0, 0) {
		t.Fatalf("content change should make cell dirty")
	}
}

func TestCellBufferLockResizeFill(t *testing.T) {
	cb := &CellBuffer{w: 2, h: 2, cells: make([]cell, 4)}
	cb.Fill('x', StyleDefault)
	if got, _, _ := cb.Get(1, 1); got != "x" {
		t.Fatalf("unexpected fill content: %q", got)
	}

	cb.LockCell(1, 1)
	if cb.Dirty(1, 1) {
		t.Fatalf("locked cell should not be dirty")
	}
	cb.UnlockCell(1, 1)
	if !cb.Dirty(1, 1) {
		t.Fatalf("unlocked cell should be dirty")
	}

	cb.Put(0, 0, "ab", StyleDefault)
	cb.Resize(3, 1)
	if got, _, _ := cb.Get(0, 0); got != "a" {
		t.Fatalf("resize should preserve content, got %q", got)
	}
	if _, _, width := cb.Get(0, 0); width != 1 {
		t.Fatalf("unexpected width after resize: %d", width)
	}

	cb.Invalidate()
	if !cb.Dirty(0, 0) {
		t.Fatalf("invalidate should mark content dirty")
	}
}

func TestCellBufferOutOfRangeAndResizeNoop(t *testing.T) {
	cb := &CellBuffer{w: 2, h: 1, cells: make([]cell, 2)}

	if rest, width := cb.Put(-1, 0, "x", StyleDefault); rest != "x" || width != 0 {
		t.Fatalf("out-of-range put should be a no-op, got rest=%q width=%d", rest, width)
	}

	cb.Put(0, 0, "x", StyleDefault.Foreground(ColorRed).Background(ColorBlue))
	cb.SetDirty(0, 0, false)

	if rest, width := cb.Put(0, 0, "", StyleDefault.Foreground(ColorGreen)); rest != "" || width != 0 {
		t.Fatalf("empty put should be a no-op, got rest=%q width=%d", rest, width)
	}
	afterStr, afterStyle, afterWidth := cb.Get(0, 0)
	if afterStr != " " || afterStyle.GetForeground() != ColorGreen || afterWidth != 1 {
		t.Fatalf("empty put should clear the cell and apply the new style: got=%q/%v/%d", afterStr, afterStyle, afterWidth)
	}
	if got := cb.Dirty(0, 0); !got {
		t.Fatalf("empty put should mark the cell dirty")
	}

	cb.LockCell(-1, 0)
	cb.LockCell(2, 0)
	cb.UnlockCell(-1, 0)
	cb.UnlockCell(2, 0)

	cb.Resize(2, 1)
	if w, h := cb.Size(); w != 2 || h != 1 {
		t.Fatalf("resize no-op changed size to %dx%d", w, h)
	}
}
