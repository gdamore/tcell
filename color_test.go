// Copyright 2018 The TCell Authors
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
	"testing"
)

func TestColorValues(t *testing.T) {
	var values = []struct {
		color Color
		hex   int32
	}{
		{ColorRed, 0x00FF0000},
		{ColorGreen, 0x00008000},
		{ColorLime, 0x0000FF00},
		{ColorBlue, 0x000000FF},
		{ColorBlack, 0x00000000},
		{ColorWhite, 0x00FFFFFF},
		{ColorSilver, 0x00C0C0C0},
	}

	for _, tc := range values {
		if tc.color.Hex() != tc.hex {
			t.Errorf("Color: %x != %x", tc.color.Hex(), tc.hex)
		}
	}
}

func TestColorFitting(t *testing.T) {
	pal := []Color{}
	for i := 0; i < 255; i++ {
		pal = append(pal, Color(i))
	}

	// Exact color fitting on ANSI colors
	for i := 0; i < 7; i++ {
		if FindColor(Color(i), pal[:8]) != Color(i) {
			t.Errorf("Color ANSI fit fail at %d", i)
		}
	}
	// Grey is closest to Silver
	if FindColor(Color(8), pal[:8]) != Color(7) {
		t.Errorf("Grey does not fit to silver")
	}
	// Color fitting of upper 8 colors.
	for i := 9; i < 16; i++ {
		if FindColor(Color(i), pal[:8]) != Color(i%8) {
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
	}{
		{"#FF0000", ColorRed},
		{"black", ColorBlack},
		{"orange", ColorOrange},
	}
	for _, v := range values {
		c := GetColor(v.name)
		if c.Hex() != v.color.Hex() {
			t.Errorf("Wrong color for %v: %v", v.name, c.Hex())
		}
	}
}

func TestColorRGB(t *testing.T) {
	r, g, b := GetColor("#112233").RGB()
	if r != 0x11 || g != 0x22 || b != 0x33 {
		t.Errorf("RGB wrong (%x, %x, %x)", r, g, b)
	}
}
