// Copyright 2025 The TCell Authors
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
	ic "image/color"
	"testing"

	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/vt"
)

// TestColorWrappers just tests the legacy API wrappers.
// The meaty tests are in the color package.
func TestColorWrappers(t *testing.T) {
	p := []color.Color{
		color.XTerm0, color.XTerm1, color.XTerm2, color.XTerm3,
		color.XTerm4, color.XTerm5, color.XTerm6, color.XTerm7,
	}
	for c := ColorBlack; c < Color255; c++ {
		if FindColor(c, p) != color.Find(c, p) {
			t.Errorf("API mismatch for color %s", c.String())
		}
	}
	if GetColor("#112233") != color.GetColor("#112233") {
		t.Errorf("Wrong colors %s != %s", GetColor("#112233").String(), color.GetColor("#112233"))
	}
	red := ic.RGBA{0xFF, 0x00, 0x00, 0x01}
	if FromImageColor(red) != color.FromImageColor(red) {
		t.Errorf("wrong colors %d %d", FromImageColor(red), color.FromImageColor(red))
	}

	if NewHexColor(0x1234) != color.NewHexColor(0x1234) {
		t.Errorf("hex colors don't match: %d %d", NewHexColor(0x1234), color.NewHexColor(0x1234))
	}

	if NewRGBColor(11, 22, 33) != color.NewRGBColor(11, 22, 33) {
		t.Errorf("rgb colors don't match: %d %d", NewRGBColor(11, 22, 33), color.NewRGBColor(11, 22, 33))
	}
}

func TestColorNone(t *testing.T) {
	_, s := NewMockScreen(t, vt.MockOptSize{X: 80, Y: 24})
	defer s.Fini()

	st := StyleDefault.Foreground(ColorBlack).Background(ColorWhite)
	s.Fill(' ', st)
	if _, s1, _ := s.Get(0, 0); s1 != st {
		t.Errorf("Wrong style! fg %s bg %s", s1.fg.String(), s1.bg.String())
	}
	st2 := st.Foreground(ColorNone).Background(ColorNone)
	s.Fill('X', st2)
	if _, s1, _ := s.Get(0, 0); s1 != st {
		t.Errorf("Wrong style! fg %s bg %s", s1.fg.String(), s1.bg.String())
	}
	red := st.Foreground(ColorRed).Background(ColorNone)
	s.Put(1, 0, " ", red)
	if _, s1, _ := s.Get(1, 0); s1 != red.Background(st.bg) {
		t.Errorf("Wrong style! fg %s bg %s", s1.fg.String(), s1.bg.String())
	}
	if _, s1, _ := s.Get(0, 0); s1 != st {
		t.Errorf("Wrong style! fg %s bg %s", s1.fg.String(), s1.bg.String())
	}
	pink := st.Background(ColorPink).Foreground(ColorNone)
	s.Put(1, 0, " ", pink)
	combined := pink.Foreground(ColorRed)

	if _, s1, _ := s.Get(1, 0); s1 != combined {
		t.Errorf("Wrong style! fg %s bg %s", s1.fg.String(), s1.bg.String())
	}
}
