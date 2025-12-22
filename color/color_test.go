// Copyright 2025 The TCell Authors
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

package color

import (
	ic "image/color"
	"testing"
)

func TestColorValues(t *testing.T) {
	var values = []struct {
		color Color
		hex   int32
	}{
		{Red, 0xFF0000},
		{Green, 0x008000},
		{Lime, 0x00FF00},
		{Blue, 0x0000FF},
		{Black, 0x000000},
		{White, 0xFFFFFF},
		{Silver, 0xC0C0C0},
		{Navy, 0x000080},
	}

	for _, tc := range values {
		if tc.color.Hex() != tc.hex {
			t.Errorf("Color: %x != %x", tc.color.Hex(), tc.hex)
		}

		if tc.color.TrueColor().Hex() != tc.hex {
			t.Errorf("TrueColor %x != %x", tc.color.TrueColor().Hex(), tc.hex)
		}
	}
}

func TestColorFitting(t *testing.T) {
	var pal []Color
	for i := range 255 {
		pal = append(pal, PaletteColor(i))
	}

	// Exact color fitting on ANSI colors
	for i := range 7 {
		if Find(PaletteColor(i), pal[:8]) != PaletteColor(i) {
			t.Errorf("Color ANSI fit fail at %d", i)
		}
	}
	// Grey is closest to Silver
	if Find(PaletteColor(8), pal[:8]) != PaletteColor(7) {
		t.Errorf("Grey does not fit to silver")
	}
	// Color fitting of upper 8 colors.
	for i := 9; i < 16; i++ {
		if Find(PaletteColor(i), pal[:8]) != PaletteColor(i%8) {
			t.Errorf("Color fit fail at %d", i)
		}
	}
	// Imperfect fit
	if Find(OrangeRed, pal[:16]) != Red ||
		Find(AliceBlue, pal[:16]) != White ||
		Find(Pink, pal) != XTerm217 ||
		Find(Sienna, pal) != XTerm173 ||
		Find(GetColor("#00FD00"), pal) != Lime {
		t.Errorf("Imperfect color fit")
	}

}

func TestColorNameLookup(t *testing.T) {
	var values = []struct {
		name  string
		color Color
		rgb   bool
	}{
		{"#FF0000", Red, true},
		{"black", Black, false},
		{"orange", Orange, true},
		{"door", Default, false},
	}
	for _, v := range values {
		c := GetColor(v.name)
		if c.Hex() != v.color.Hex() {
			t.Errorf("Wrong color for %v: %v", v.name, c.Hex())
		}
		if v.rgb {
			if !c.IsRGB() {
				t.Errorf("Color should have RGB: %v", v.name)
			}
		} else {
			if c.IsRGB() {
				t.Errorf("Named color should not be RGB: %v", v.name)
			}
		}

		if c.TrueColor().Hex() != v.color.Hex() {
			t.Errorf("TrueColor did not match")
		}
	}

	// these colors only have strings (for debugging), you cannot use them to create a color
	if None.String() != "none" {
		t.Errorf("ColorNone did not match")
	}
	if Reset.String() != "reset" {
		t.Errorf("Reset did not match")
	}
	if Default.String() != "default" {
		t.Errorf("Default did not match")
	}
}

func TestColorRGB(t *testing.T) {
	r, g, b := GetColor("#112233").RGB()
	if r != 0x11 || g != 0x22 || b != 0x33 {
		t.Errorf("RGB wrong (%x, %x, %x)", r, g, b)
	}
}

func TestFromImageColor(t *testing.T) {
	red := ic.RGBA{0xFF, 0x00, 0x00, 0x01}
	white := ic.Gray{0xFF}
	cyan := ic.CMYK{0xFF, 0x00, 0x00, 0x00}

	if hex := FromImageColor(red).Hex(); hex != 0xFF0000 {
		t.Errorf("%v is not 0xFF0000", hex)
	}
	if hex := FromImageColor(white).Hex(); hex != 0xFFFFFF {
		t.Errorf("%v is not 0xFFFFFF", hex)
	}
	if hex := FromImageColor(cyan).Hex(); hex != 0x00FFFF {
		t.Errorf("%v is not 0x00FFFF", hex)
	}
}

func TestColorRGBA(t *testing.T) {
	r, g, b, a := Red.RGBA()
	if r != 0xffff || g != 0 || b != 0 || a != 0xffff {
		t.Errorf("Wrong RGBA for red: %x %x %x %x", r, g, b, a)
	}
	r, g, b, a = Red.TrueColor().RGBA()
	if r != 0xffff || g != 0 || b != 0 || a != 0xffff {
		t.Errorf("Wrong RGBA for red.TrueColor: %x %x %x %x", r, g, b, a)
	}

	r, g, b, a = Default.RGBA()
	if r != 0 || g != 0 || b != 0 || a != 0 {
		t.Errorf("Non-zero RGBA for default")
	}

	r, g, b, a = None.RGBA()
	if r != 0 || g != 0 || b != 0 || a != 0 {
		t.Errorf("Non-zero RGBA for none")
	}
	r, g, b, a = Reset.RGBA()
	if r != 0 || g != 0 || b != 0 || a != 0 {
		t.Errorf("Non-zero RGBA for reset")
	}
}
