// Copyright 2024 The TCell Authors
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
	"testing"

	"github.com/gdamore/tcell/v3/color"
)

func TestStyle(t *testing.T) {
	_, s := NewMockScreen(t)
	defer s.Fini()

	style := StyleDefault
	fg, bg, attr := style.fg, style.bg, style.attrs

	if fg != ColorDefault || bg != ColorDefault || attr != AttrNone {
		t.Errorf("Bad default style (%v, %v, %v)", fg, bg, attr)
	}

	s2 := style.
		Background(ColorRed).
		Foreground(ColorBlue).
		Blink(true)

	fg, bg, attr = s2.fg, s2.bg, s2.attrs
	if fg != ColorBlue || bg != ColorRed || attr != AttrBlink {
		t.Errorf("Bad custom style (%v, %v, %v)", fg, bg, attr)
	}

	if !s2.HasBlink() {
		t.Errorf("blink is false")
	}
	if s2.HasItalic() {
		t.Errorf("italic is true")
	}
	if s2.HasDim() {
		t.Errorf("dim is true")
	}
	if s2.HasBold() {
		t.Errorf("bold is true")
	}
	if s2.HasUnderline() {
		t.Errorf("underline is true")
	}
	if s2.HasReverse() {
		t.Errorf("reverse is true")
	}
	if s2.HasStrikeThrough() {
		t.Errorf("strike-through is true")
	}
	if id, url := s2.GetUrl(); id != "" || url != "" {
		t.Errorf("url not empty: %q %q", id, url)
	}
	if s2.GetAttributes() != AttrBlink {
		t.Errorf("wrong attrs %x", s2.GetAttributes())
	}
	if s2.GetBackground() != color.Red || s2.GetForeground() != color.Blue {
		t.Errorf("wrong colors %s %s", s2.GetForeground().String(), s2.GetBackground().String())
	}

	us := s2.Url("http://example.com")
	if id, url := us.GetUrl(); id != "" || url != "http://example.com" {
		t.Errorf("url wrong: %q %q", id, url)
	}
	us = us.Url("http://example.com").UrlId("someId")
	if id, url := us.GetUrl(); id != "someId" || url != "http://example.com" {
		t.Errorf("url wrong: %q %q", id, url)
	}
	us = us.Underline(UnderlineStyleDotted, color.Pink)
	if us.GetUnderlineColor() != color.Pink {
		t.Errorf("underline color wrong: %q", us.GetUnderlineColor().String())
	}
	if us.GetUnderlineStyle() != UnderlineStyleDotted {
		t.Errorf("wrong underline style: %v", us.GetUnderlineStyle())
	}

	us = us.Normal().Reverse(true).Italic(false)
	if us.GetAttributes() != AttrReverse {
		t.Errorf("wrong attributes: %v", us.GetAttributes())
	}
}
