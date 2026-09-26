// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

import "testing"

func TestSetContentPreservesGraphemes(t *testing.T) {
	tests := []struct {
		name      string
		primary   rune
		combining []rune
		want      string
		width     int
	}{
		{"printable ASCII", 'x', nil, "x", 1},
		{"combining", 'e', []rune{'\u0301'}, "e\u0301", 1},
		{"wide", '宽', nil, "宽", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &baseScreen{screenImpl: &tScreen{cells: CellBuffer{w: 8, h: 1, cells: make([]cell, 8)}}}
			s.SetContent(0, 0, tt.primary, tt.combining, StyleDefault)
			got, _, width := s.Get(0, 0)
			if got != tt.want || width != tt.width {
				t.Fatalf("Get() = %q, width %d; want %q, width %d", got, width, tt.want, tt.width)
			}
		})
	}
}

func TestSetContentASCIIRepaintDoesNotAllocate(t *testing.T) {
	s := &baseScreen{screenImpl: &tScreen{cells: CellBuffer{w: 8, h: 1, cells: make([]cell, 8)}}}
	s.SetContent(0, 0, 'x', nil, StyleDefault)
	if allocs := testing.AllocsPerRun(1000, func() {
		s.SetContent(0, 0, 'x', nil, StyleDefault)
	}); allocs != 0 {
		t.Fatalf("unchanged printable ASCII cell redraw allocated %g times per call", allocs)
	}
}
