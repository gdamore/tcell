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

import (
	"strings"
	"testing"
)

var mixedCellStream = strings.Repeat("Hello, 世界 👩‍🚀 e\u0301 ", 8)

func BenchmarkCellBufferPutCurrent(b *testing.B) {
	benchCellBufferPut(b, "current", func(cb *CellBuffer, x, y int, str string, style Style) (string, int) {
		return cb.Put(x, y, str, style)
	})
}

func BenchmarkCellBufferPutStreamCurrent(b *testing.B) {
	benchCellBufferPutStream(b, "current", func(cb *CellBuffer, x int, y int, str string, style Style) (string, int) {
		return cb.Put(x, y, str, style)
	})
}

func benchCellBufferPut(b *testing.B, name string, put func(*CellBuffer, int, int, string, Style) (string, int)) {
	cases := []struct {
		name string
		str  string
	}{
		{name: "ascii", str: "Hello, terminal"},
		{name: "combining", str: "e\u0301"},
		{name: "wide", str: "宽"},
		{name: "emoji", str: "👩‍🚀"},
	}

	for _, tc := range cases {
		b.Run(name+"/"+tc.name, func(b *testing.B) {
			cb := &CellBuffer{w: 8, h: 1, cells: make([]cell, 8)}
			style := Style{}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cb.cells[0] = cell{}
				_, _ = put(cb, 0, 0, tc.str, style)
			}
		})
	}
}

func benchCellBufferPutStream(b *testing.B, name string, put func(*CellBuffer, int, int, string, Style) (string, int)) {
	b.Helper()

	b.Run(name+"/mixed-stream", func(b *testing.B) {
		cb := &CellBuffer{w: 128, h: 1, cells: make([]cell, 128)}
		style := Style{}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := range cb.cells {
				cb.cells[j] = cell{}
			}
			x := 0
			rest := mixedCellStream
			for rest != "" && x < cb.w {
				var width int
				rest, width = put(cb, x, 0, rest, style)
				if width == 0 {
					break
				}
				x += width
			}
		}
	})
}
