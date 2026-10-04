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

// shown marks every cell as drawn, the way Show leaves the buffer.
func shown(cb *CellBuffer) {
	w, h := cb.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			cb.SetDirty(x, y, false)
		}
	}
}

func TestCellBufferRedrawAfterFillIsClean(t *testing.T) {
	cb := &CellBuffer{}
	cb.Resize(4, 1)
	cb.SetContent(0, 0, 'x', nil, StyleDefault)
	shown(cb)

	// Clear and draw the same frame again, as an application that redraws
	// everything each frame does.
	cb.Fill(' ', StyleDefault)
	cb.SetContent(0, 0, 'x', nil, StyleDefault)
	for x := 0; x < 4; x++ {
		if cb.Dirty(x, 0) {
			t.Errorf("cell %d is dirty after redrawing the same content", x)
		}
	}

	cb.SetContent(1, 0, 'y', nil, StyleDefault)
	if !cb.Dirty(1, 0) {
		t.Error("a cell whose content changed is not dirty")
	}
}

func TestCellBufferWidePutMarksNextCellDirty(t *testing.T) {
	cb := &CellBuffer{}
	cb.Resize(4, 1)
	shown(cb)

	cb.SetContent(0, 0, '世', nil, StyleDefault)
	if !cb.Dirty(0, 0) || !cb.Dirty(1, 0) {
		t.Errorf("wide character: dirty = %v, %v, want both cells dirty", cb.Dirty(0, 0), cb.Dirty(1, 0))
	}
}
