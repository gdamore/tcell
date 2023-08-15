// Copyright 2022 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package encoding

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestGBK(t *testing.T) {
	enc := tcell.GetEncoding("GBK")
	if enc == nil {
		t.Fatal("NULL encoding for GBK")
	}
	glyph, _ := enc.NewDecoder().Bytes([]byte{0x82, 0x74})
	if string(glyph) != "倀" {
		t.Errorf("failed to match: %s != 倀", string(glyph))
	}
}

func TestAscii(t *testing.T) {
	encodings := []string{
		"ASCII",
		"ISO-8859-1",
		"KOI8-R",
		"KOI8-U",
		"SJIS",
		"Big5",
		"GB2312",
		"GB18030",
		"EUC-JP",
		"EUCKR",
	}

	for _, name := range encodings {
		t.Run(name, func(t *testing.T) {
			enc := tcell.GetEncoding(name)
			if enc == nil {
				t.Errorf("Failed getting encoding for %s", name)
				return
			}
			encoder := enc.NewEncoder()
			decoder := enc.NewDecoder()
			// Ensure that all US-ASCII (lower 7 bit values) encode and decode identically
			for i := byte(0); i < 126; i++ { // well, KOI8-R has some problem with "~"
				s := string([]byte{i})
				if x, err := encoder.String(s); err != nil || x != s {
					t.Errorf("failed encoding for character: %d, err %v expect %s got %s", i, err, s, x)
				}
				if x, err := decoder.String(s); err != nil || x != s {
					t.Errorf("failed decoding for character: %d, err %v expect %s got %s", i, err, s, x)
				}
			}
		})
	}
}
