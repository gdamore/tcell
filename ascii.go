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

package tcell

import (
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

type ascii struct{ transform.NopResetter }
type asciiDecoder struct{ transform.NopResetter }
type asciiEncoder struct{ transform.NopResetter }

// ASCII represents an basic 7-bit ASCII scheme.  It decodes directly to UTF-8
// without change, as all ASCII values are legal UTF-8.  It encodes any UTF-8
// runes outside of ASCII to 0x1A, the ASCII substitution character.
var ASCII encoding.Encoding = ascii{}

func (ascii) NewDecoder() transform.Transformer {
	return asciiDecoder{}
}

func (ascii) NewEncoder() transform.Transformer {
	return asciiEncoder{}
}

func (asciiDecoder) Transform(dst, src []byte, atEOF bool) (int, int, error) {
	var e error
	var ndst, nsrc int
	for _, c := range src {
		if ndst >= len(dst) {
			e = transform.ErrShortDst
			break
		}
		dst[ndst] = c
		ndst++
		nsrc++
	}
	return ndst, nsrc, e
}

func (asciiEncoder) Transform(dst, src []byte, atEOF bool) (int, int, error) {
	var e error
	var ndst, nsrc int
	var sz int
	for nsrc < len(src) {
		if ndst >= len(dst) {
			e = transform.ErrShortDst
			break
		}
		r := rune(src[nsrc])
		if r < utf8.RuneSelf {
			dst[ndst] = uint8(r)
			nsrc++
			ndst++
			continue
		}

		// No valid runes beyond ASCII.  However, we need to consume
		// the full rune, and report incomplete runes properly.

		// Attempt to decode a multibyte rune
		r, sz = utf8.DecodeRune(src[nsrc:])
		if sz == 1 {
			// If its inconclusive due to insufficient data in
			// in the source, report it
			if !atEOF && !utf8.FullRune(src[nsrc:]) {
				e = transform.ErrShortSrc
				break
			}
		}
		nsrc += sz
		dst[ndst] = encoding.ASCIISub
		ndst++
	}
	return ndst, nsrc, e
}
