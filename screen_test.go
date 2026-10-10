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
	"strings"
	"sync"
	"testing"
	"unicode/utf8"

	"github.com/gdamore/tcell/v3/vt"
)

func TestSetContent(t *testing.T) {
	_, scr := NewMockScreen(t, vt.MockOptSize{X: 8, Y: 1})
	defer scr.Fini()
	style := StyleDefault.Foreground(ColorRed).Background(ColorBlue)
	for _, tc := range []struct {
		name string
		main rune
		comb []rune
	}{
		{name: "ASCII", main: 'x'},
		{name: "EmptyCombining", main: 'x', comb: []rune{}},
		{name: "Drawing", main: '┼'},
		{name: "Greek", main: 'Ω'},
		{name: "Cyrillic", main: 'Ж'},
		{name: "Arabic", main: 'ش', comb: []rune{'\u064e'}},
		{name: "Hebrew", main: 'ש', comb: []rune{'\u05b8'}},
		{name: "Devanagari", main: 'क', comb: []rune{'\u093f'}},
		{name: "CJK", main: '界'},
		{name: "Emoji", main: '👩', comb: []rune{'\u200d', '🚀'}},
		{name: "Accent", main: 'e', comb: []rune{'\u0301'}},
		{name: "StandaloneCombining", main: '\u0301'},
		{name: "MultipleGraphemes", main: 'e', comb: []rune{'\u0301', 'x'}},
		{name: "LongCombining", main: 'e', comb: []rune(strings.Repeat("\u0301", 100))},
		{name: "NUL", main: 0},
		{name: "DEL", main: 127},
		{name: "RuneSelf", main: utf8.RuneSelf},
		{name: "Negative", main: -1},
		{name: "Surrogate", main: 0xd800},
		{name: "AboveMaxRune", main: utf8.MaxRune + 1},
		{name: "InvalidCombining", main: 'e', comb: []rune{-1, 0xd800}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for range 2 {
				// Compare with the original conversion and public Put path,
				// including segmentation, width, and ColorNone behavior.
				scr.Fill(' ', style)
				input := string(append([]rune{tc.main}, tc.comb...))
				scr.Put(0, 0, input, StyleDefault.Foreground(ColorNone).Background(ColorNone))
				scr.SetContent(4, 0, tc.main, tc.comb, StyleDefault.Foreground(ColorNone).Background(ColorNone))
				want, wantStyle, wantWidth := scr.Get(0, 0)
				got, gotStyle, gotWidth := scr.Get(4, 0)
				if got != want || gotStyle != wantStyle || gotWidth != wantWidth {
					t.Fatalf("got %q/%v/%d, want %q/%v/%d", got, gotStyle, gotWidth, want, wantStyle, wantWidth)
				}
			}
		})
	}
}

func TestSetContentAllocations(t *testing.T) {
	_, scr := NewMockScreen(t, vt.MockOptSize{X: 8, Y: 1})
	defer scr.Fini()
	for _, tc := range []struct {
		name string
		text string
		comb []rune
	}{
		{name: "ASCII", text: "abcdefghijklmnopqrstuvwxyz"},
		{name: "Drawing", text: "─│┌┐└┘├┤┬┴┼═║╔╗╚╝╠╣╦╩╬▀▄█░▒▓"},
		{name: "Latin", text: "éñüçøå"},
		{name: "Greek", text: "Ελληνικά"},
		{name: "Cyrillic", text: "Русский"},
		{name: "CJK", text: "中文日本語한국어"},
		{name: "Emoji", text: "😀😃😄😁"},
		{name: "Accent", text: "eaio", comb: []rune{'\u0301'}},
		{name: "Arabic", text: "شسكت", comb: []rune{'\u064e'}},
		{name: "Devanagari", text: "कखगघ", comb: []rune{'\u093f'}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runes := []rune(tc.text)
			for _, r := range runes {
				scr.SetContent(0, 0, r, tc.comb, StyleDefault)
			}
			i := 0
			allocs := testing.AllocsPerRun(100, func() {
				scr.SetContent(0, 0, runes[i%len(runes)], tc.comb, StyleDefault)
				i++
			})
			if allocs != 0 {
				t.Fatalf("warm SetContent allocated %v times, want 0", allocs)
			}
		})
	}
}

func TestSetContentCacheReuse(t *testing.T) {
	_, scr := NewMockScreen(t, vt.MockOptSize{X: 8, Y: 1})
	defer scr.Fini()
	bs := scr.(*baseScreen)
	scr.SetContent(0, 0, '界', nil, StyleDefault)
	retained, _, _ := scr.Get(0, 0)
	bs.GetCells().SetDirty(0, 0, false)
	// Cycle through a large CJK repertoire to verify earlier glyphs remain
	// available without allocating again after the cache has grown.
	const repertoireSize = 4110
	for i := range repertoireSize {
		scr.SetContent(4, 0, rune(0x4e00+i), nil, StyleDefault)
	}
	if retained != "界" || bs.GetCells().Dirty(0, 0) {
		t.Fatal("cache growth changed retained cell content or dirty state")
	}
	if allocs := testing.AllocsPerRun(5, func() {
		for i := range repertoireSize {
			scr.SetContent(4, 0, rune(0x4e00+i), nil, StyleDefault)
		}
		scr.SetContent(4, 0, '界', nil, StyleDefault)
	}); allocs != 0 {
		t.Fatalf("cached CJK repertoire allocated %v times", allocs)
	}
	if got, _, width := scr.Get(4, 0); got != retained || width != 2 {
		t.Fatalf("reused cell is %q width=%d", got, width)
	}
	// The scratch buffer size does not limit which sequences can be cached.
	long := []rune(strings.Repeat("\u0301", 100))
	scr.SetContent(4, 0, 'e', long, StyleDefault)
	if _, ok := bs.contentCache["e"+string(long)]; !ok {
		t.Fatal("long combining sequence was not cached")
	}
	// Out-of-bounds writes must not fill the cache.
	for _, pos := range [][2]int{{-1, 0}, {0, -1}, {8, 0}, {0, 1}} {
		scr.SetContent(pos[0], pos[1], '☃', nil, StyleDefault)
	}
	if _, ok := bs.contentCache["☃"]; ok {
		t.Fatal("out-of-bounds content was cached")
	}
}

func TestSetContentConcurrent(t *testing.T) {
	_, scr := NewMockScreen(t, vt.MockOptSize{X: 16, Y: 1})
	defer scr.Fini()
	var wg sync.WaitGroup
	for x := range 8 {
		wg.Go(func() {
			for i := range 100 {
				scr.SetContent(x*2, 0, rune(0x4e00+i+x), nil, StyleDefault)
				scr.Get(x*2, 0)
			}
		})
	}
	wg.Wait()
	for x := range 8 {
		if got, _, width := scr.Get(x*2, 0); got != string(rune(0x4e00+99+x)) || width != 2 {
			t.Fatalf("cell %d has %q width=%d", x*2, got, width)
		}
	}
}
