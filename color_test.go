// Copyright 2023 The TCell Authors
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
	ic "image/color"
	"testing"
)

func TestColorValues(t *testing.T) {
	var values = []struct {
		color Color
		hex   int32
	}{
		{ColorRed, 0xFF0000},
		{ColorGreen, 0x008000},
		{ColorLime, 0x00FF00},
		{ColorBlue, 0x0000FF},
		{ColorBlack, 0x000000},
		{ColorWhite, 0xFFFFFF},
		{ColorSilver, 0xC0C0C0},
		{ColorNavy, 0x000080},
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
	for i := 0; i < 255; i++ {
		pal = append(pal, PaletteColor(i))
	}

	// Exact color fitting on ANSI colors
	for i := 0; i < 7; i++ {
		if FindColor(PaletteColor(i), pal[:8]) != PaletteColor(i) {
			t.Errorf("Color ANSI fit fail at %d", i)
		}
	}
	// Grey is closest to Silver
	if FindColor(PaletteColor(8), pal[:8]) != PaletteColor(7) {
		t.Errorf("Grey does not fit to silver")
	}
	// Color fitting of upper 8 colors.
	for i := 9; i < 16; i++ {
		if FindColor(PaletteColor(i), pal[:8]) != PaletteColor(i%8) {
			t.Errorf("Color fit fail at %d", i)
		}
	}
	// Imperfect fit
	if FindColor(ColorOrangeRed, pal[:16]) != ColorRed ||
		FindColor(ColorAliceBlue, pal[:16]) != ColorWhite ||
		FindColor(ColorPink, pal) != Color217 ||
		FindColor(ColorSienna, pal) != Color173 ||
		FindColor(GetColor("#00FD00"), pal) != ColorLime {
		t.Errorf("Imperfect color fit")
	}

}

func TestColorNameLookup(t *testing.T) {
	var values = []struct {
		name  string
		color Color
		rgb   bool
	}{
		{"#FF0000", ColorRed, true},
		{"black", ColorBlack, false},
		{"orange", ColorOrange, false},
		{"door", ColorDefault, false},
	}
	for _, v := range values {
		c := GetColor(v.name)
		if c.Hex() != v.color.Hex() {
			t.Errorf("Wrong color for %v: %v", v.name, c.Hex())
		}
		if v.rgb {
			if c&ColorIsRGB == 0 {
				t.Errorf("Color should have RGB")
			}
		} else {
			if c&ColorIsRGB != 0 {
				t.Errorf("Named color should not be RGB")
			}
		}

		if c.TrueColor().Hex() != v.color.Hex() {
			t.Errorf("TrueColor did not match")
		}
	}

	// these colors only have strings (for debugging), you cannot use them to create a color
	if ColorNone.String() != "none" {
		t.Errorf("ColorNone did not match")
	}
	if ColorReset.String() != "reset" {
		t.Errorf("ColorReset did not match")
	}
	if ColorDefault.String() != "default" {
		t.Errorf("ColorDefault did not match")
	}
}

func TestColorRGB(t *testing.T) {
	r, g, b := GetColor("#112233").RGB()
	if r != 0x11 || g != 0x22 || b != 0x33 {
		t.Errorf("RGB wrong (%x, %x, %x)", r, g, b)
	}
}

func TestFromImageColor(t *testing.T) {
	red := ic.RGBA{0xFF, 0x00, 0x00, 0x00}
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

func TestColorNone(t *testing.T) {
	s := mkTestScreen(t, "")
	s.Init()
	s.SetSize(80, 24)
	st := StyleDefault.Foreground(ColorBlack).Background(ColorWhite)
	s.Fill(' ', st)
	if _, _, s1, _ := s.GetContent(0, 0); s1 != st {
		t.Errorf("Wrong style! fg %s bg %s", s1.fg.String(), s1.bg.String())
	}
	st2 := st.Foreground(ColorNone).Background(ColorNone)
	s.Fill('X', st2)
	if _, _, s1, _ := s.GetContent(0, 0); s1 != st {
		t.Errorf("Wrong style! fg %s bg %s", s1.fg.String(), s1.bg.String())
	}
	red := st.Foreground(ColorRed).Background(ColorNone)
	s.SetContent(1, 0, ' ', nil, red)
	if _, _, s1, _ := s.GetContent(1, 0); s1 != red.Background(st.bg) {
		t.Errorf("Wrong style! fg %s bg %s", s1.fg.String(), s1.bg.String())
	}
	if _, _, s1, _ := s.GetContent(0, 0); s1 != st {
		t.Errorf("Wrong style! fg %s bg %s", s1.fg.String(), s1.bg.String())
	}
	pink := st.Background(ColorPink).Foreground(ColorNone)
	s.SetContent(1, 0, ' ', nil, pink)
	combined := pink.Foreground(ColorRed)

	if _, _, s1, _ := s.GetContent(1, 0); s1 != combined {
		t.Errorf("Wrong style! fg %s bg %s", s1.fg.String(), s1.bg.String())
	}
}
