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

package vt

import (
	"testing"

	"github.com/gdamore/tcell/v3/color"
)

func BenchmarkPutRuneCurrent(b *testing.B) {
	benchPutRune(b, "current", (*emulator).putRune)
}

func benchPutRune(b *testing.B, name string, put func(*emulator, rune)) {
	cases := []struct {
		name string
		r    rune
		seq  []rune
	}{
		{name: "ascii", r: 'a'},
		{name: "width", r: 'π'},
		{name: "wide", r: '宽'},
		{name: "combining", r: '\u0301'},
		{name: "line", seq: []rune{
			'\u250C', '\u2510', '\u2514', '\u2518',
			'\u2500', '\u2502', '\u251C', '\u2524',
			'\u252C', '\u2534', '\u253C', '\u256D',
			'\u256E', '\u256F', '\u2570', '\u2571',
		}},
		{name: "mixed32", seq: []rune{
			'\u0416', '\u0414', '\u042E', '\u042F',
			'\u041F', '\u041B', '\u0424', '\u042B',
			'\u042D', '\u0411', '\u0413', '\u0428',
			'\u00E9', '\u00F6', '\u00FC', '\u00F1',
			'\u00E7', '\u00F8', '\u00E5', '\u00DF',
			'\u0142', '\u0111', '\u0127', '\u0131',
			'\U0001F600', '\U0001F680', '\U0001F9EA', '\U0001F525',
			'\U0001F355', '\U0001F389', '\U0001F4BB', '\U0001F4E6',
		}},
		{name: "mixed64", seq: sweepMixedRunes64},
	}

	for _, tc := range cases {
		b.Run(name+"/"+tc.name, func(b *testing.B) {
			em := NewEmulator(NewMockBackend(MockOptSize{X: 8, Y: 1}, MockOptColors(0))).(*emulator)
			em.style = BaseStyle.WithFg(color.White).WithBg(color.Black)
			em.defaultStyle = em.style
			em.localModes[PmGraphemeClusters] = ModeOn
			em.localModes[PmAutoMargin] = ModeOn
			em.cells[0].S = em.style
			em.cells[0].W = 1
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r := tc.r
				if len(tc.seq) != 0 {
					r = tc.seq[i%len(tc.seq)]
				}
				em.pos = Coord{X: 1, Y: 0}
				em.lastIndex = 1
				em.autoWrap = false
				em.cells[0].C = "e"
				em.cells[0].S = em.style
				em.cells[0].W = 1
				put(em, r)
			}
		})
	}
}
