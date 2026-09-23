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

package tcell_test

import (
	"testing"

	"github.com/gdamore/tcell/v3"
)

func TestCellBufferResizeInvalidatesBlankCells(t *testing.T) {
	tests := []struct {
		name                string
		oldWidth, oldHeight int
		width, height       int
	}{
		{"initialize", 0, 0, 2, 2},
		{"grow", 3, 2, 4, 3},
		{"shrink", 3, 2, 2, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cb tcell.CellBuffer
			style := tcell.StyleDefault.Bold(true)
			if tt.oldWidth > 0 {
				cb.Resize(tt.oldWidth, tt.oldHeight)
				cb.Put(0, 0, "x", style)
				cb.SetDirty(0, 0, false)
				cb.Put(1, 0, "", tcell.StyleDefault)
			}

			cb.Resize(tt.width, tt.height)
			if w, h := cb.Size(); w != tt.width || h != tt.height {
				t.Fatalf("Size() = (%d, %d), want (%d, %d)", w, h, tt.width, tt.height)
			}
			if tt.oldWidth > 0 {
				if str, gotStyle, width := cb.Get(0, 0); str != "x" || gotStyle != style || width != 1 {
					t.Fatalf("resize did not preserve the nonblank cell: (%q, %v, %d)", str, gotStyle, width)
				}
			}
			for y := 0; y < tt.height; y++ {
				for x := 0; x < tt.width; x++ {
					if !cb.Dirty(x, y) {
						t.Errorf("cell (%d, %d) should be dirty after resize", x, y)
					}
					if tt.oldWidth == 0 || x != 0 || y != 0 {
						if str, _, width := cb.Get(x, y); str != " " || width != 1 {
							t.Errorf("blank cell (%d, %d) = (%q, %d), want a single space", x, y, str, width)
						}
					}
				}
			}
		})
	}
}

func TestCellBufferResizeSameSizeKeepsCellsClean(t *testing.T) {
	var cb tcell.CellBuffer
	cb.Resize(1, 1)
	cb.SetDirty(0, 0, false)
	cb.Resize(1, 1)
	if cb.Dirty(0, 0) {
		t.Fatal("resizing to the same size should leave a clean cell unchanged")
	}
}
