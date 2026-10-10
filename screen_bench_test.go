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

	"github.com/gdamore/tcell/v3/vt"
)

func BenchmarkSetContent(b *testing.B) {
	for _, tc := range []struct {
		name  string
		runes [2]rune
		comb  []rune
	}{
		{name: "ASCII", runes: [2]rune{'x', 'y'}},
		{name: "Unicode", runes: [2]rune{'é', 'ñ'}},
		{name: "Drawing", runes: [2]rune{'┼', '─'}},
		{name: "Cyrillic", runes: [2]rune{'Ж', 'Я'}},
		{name: "Wide", runes: [2]rune{'界', '世'}},
		{name: "Emoji", runes: [2]rune{'😀', '😃'}},
		{name: "Combining", runes: [2]rune{'e', 'a'}, comb: []rune{'\u0301'}},
		{name: "Arabic", runes: [2]rune{'ش', 'س'}, comb: []rune{'\u064e'}},
		{name: "Devanagari", runes: [2]rune{'क', 'ख'}, comb: []rune{'\u093f'}},
	} {
		for _, changed := range []bool{false, true} {
			name := tc.name + "/Redraw"
			if changed {
				name = tc.name + "/Changed"
			}
			b.Run(name, func(b *testing.B) {
				scr, err := NewTerminfoScreenFromTty(vt.NewMockTerm(vt.MockOptSize{X: 4, Y: 1}))
				if err != nil {
					b.Fatal(err)
				}
				if err := scr.Init(); err != nil {
					b.Fatal(err)
				}
				b.Cleanup(scr.Fini)
				for _, r := range tc.runes {
					scr.SetContent(0, 0, r, tc.comb, StyleDefault)
				}
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					r := tc.runes[0]
					if changed {
						r = tc.runes[i%2]
					}
					scr.SetContent(0, 0, r, tc.comb, StyleDefault)
				}
			})
		}
	}
}
