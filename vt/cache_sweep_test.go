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
	"fmt"
	"testing"
	"unicode/utf8"
)

type benchRuneCache struct {
	entries []benchRuneCacheEntry
	n       int
}

type benchRuneCacheEntry struct {
	r rune
	s string
}

func newBenchRuneCache(size int) *benchRuneCache {
	return &benchRuneCache{entries: make([]benchRuneCacheEntry, size)}
}

func (c *benchRuneCache) stringFor(r rune) string {
	if r < utf8.RuneSelf {
		return asciiRuneStrings[r]
	}

	for i := 0; i < c.n; i++ {
		if c.entries[i].r == r {
			return c.entries[i].s
		}
	}

	s := string(r)
	if c.n < len(c.entries) {
		c.n++
	}
	copy(c.entries[1:c.n], c.entries[:c.n-1])
	c.entries[0] = benchRuneCacheEntry{r: r, s: s}
	return s
}

var sweepMixedRunes32 = []rune{
	'\u0416', '\u0414', '\u042E', '\u042F',
	'\u041F', '\u041B', '\u0424', '\u042B',
	'\u042D', '\u0411', '\u0413', '\u0428',
	'\u00E9', '\u00F6', '\u00FC', '\u00F1',
	'\u00E7', '\u00F8', '\u00E5', '\u00DF',
	'\u0142', '\u0111', '\u0127', '\u0131',
	'\U0001F600', '\U0001F680', '\U0001F9EA', '\U0001F525',
	'\U0001F355', '\U0001F389', '\U0001F4BB', '\U0001F4E6',
}

var sweepMixedRunes64 = append(append([]rune{}, sweepMixedRunes32...),
	'\u0105', '\u0107', '\u0119', '\u0123',
	'\u0135', '\u0137', '\u0144', '\u015B',
	'\u015F', '\u0163', '\u0165', '\u017A',
	'\U0001F31F', '\U0001F4A1', '\U0001F4AF', '\U0001F680',
	'\U0001F9F0', '\U0001F9D1', '\U0001F9EB', '\U0001FAE0',
	'\u250C', '\u2510', '\u2514', '\u2518',
	'\u2500', '\u2502', '\u251C', '\u2524',
	'\u252C', '\u2534', '\u253C', '\u2571',
)

func BenchmarkRuneStringCacheSweep(b *testing.B) {
	for _, tc := range []struct {
		name string
		seq  []rune
	}{
		{name: "mixed32", seq: sweepMixedRunes32},
		{name: "mixed64", seq: sweepMixedRunes64},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for _, size := range []int{8, 16, 32, 64} {
				b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
					cache := newBenchRuneCache(size)
					for _, r := range tc.seq {
						_ = cache.stringFor(r)
					}
					b.ReportAllocs()
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						_ = cache.stringFor(tc.seq[i%len(tc.seq)])
					}
				})
			}
		})
	}
}
