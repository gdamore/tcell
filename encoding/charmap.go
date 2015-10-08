// Copyright 2015 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package encoding

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// suitable for 8-bit encodings
type cmap struct {
	transform.NopResetter
	bytes map[rune]byte
	runes [256]rune // offset by 128, as all values are identical
	ascii bool
}

type cmapDecoder struct {
	transform.NopResetter
	cmap *cmap
}
type cmapEncoder struct {
	transform.NopResetter
	cmap *cmap
}

func (c *cmap) Init() {
	c.bytes = make(map[rune]byte)
	for i := 0; i < 256; i++ {
		c.bytes[rune(i)] = byte(i)
		c.runes[i] = rune(i)
		c.ascii = true
	}
}

func (c *cmap) Map(b byte, r rune) {
	if b < 128 {
		c.ascii = false
	}

	// delete the old self-mapping
	delete(c.bytes, rune(b))

	// and add the new one
	c.bytes[r] = b
	c.runes[int(b)] = r
}

func (c *cmap) NewDecoder() transform.Transformer {
	return cmapDecoder{cmap: c}
}

func (c *cmap) NewEncoder() transform.Transformer {
	return cmapEncoder{cmap: c}
}

func (d cmapDecoder) Transform(dst, src []byte, atEOF bool) (int, int, error) {
	var e error
	var ndst, nsrc int

	for _, c := range src {
		if d.cmap.ascii && c < utf8.RuneSelf {
			if ndst >= len(dst) {
				e = transform.ErrShortDst
				break
			}
			dst[ndst] = c
			ndst++
			nsrc++
			continue
		}

		r := d.cmap.runes[c]
		l := utf8.RuneLen(r)

		// l will be a positive number, because we never inject invalid
		// runes into the rune map.

		if ndst+l > len(dst) {
			e = transform.ErrShortDst
			break
		}
		utf8.EncodeRune(dst[ndst:], r)
		ndst += l
		nsrc++
	}
	return ndst, nsrc, e
}

func (d cmapEncoder) Transform(dst, src []byte, atEOF bool) (int, int, error) {
	var e error
	var ndst, nsrc int
	for nsrc < len(src) {
		if ndst >= len(dst) {
			e = transform.ErrShortDst
			break
		}
		ch := src[nsrc]
		if d.cmap.ascii && ch < utf8.RuneSelf {
			dst[ndst] = ch
			nsrc++
			ndst++
			continue
		}

		// No valid runes beyond 0xFF.  However, we need to consume
		// the full rune, and report incomplete runes properly.

		// Attempt to decode a multibyte rune
		r, sz := utf8.DecodeRune(src[nsrc:])
		if r == utf8.RuneError && sz == 1 {
			// If its inconclusive due to insufficient data in
			// in the source, report it
			if !atEOF && !utf8.FullRune(src[nsrc:]) {
				e = transform.ErrShortSrc
				break
			}
		}

		if c, ok := d.cmap.bytes[r]; ok {
			dst[ndst] = c
		} else {
			dst[ndst] = encoding.ASCIISub
		}
		nsrc += sz
		ndst++
	}

	return ndst, nsrc, e
}
