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

// Color represents a color.  We encode it in a 64-bit int for efficiency.
// The coding is (MSB): <16b rsvd><16b attr><16b fgcolor><16b bgcolor>.
// This gives 16bit color options, if it ever becomes truly necessary.
// I can't see a need to get beyond 256 colors.
type Style int64

func NewStyle() Style {
	return Style(0)
}

const StyleDefault Style = 0

func (s Style) Foreground(c Color) Style {
	return (s &^ Style(0xffff0000)) | (Style(c) << 16)
}

func (s Style) Background(c Color) Style {
	return (s &^ (0xffff)) | Style(c)
}

func (s Style) Decompose() (Color, Color, AttrMask) {
	return Color((s >> 16) & 0xffff),
		Color(s & 0xfffff),
		AttrMask((s >> 32) & 0xffff)
}

func (s Style) setAttrs(attrs Style, on bool) Style {
	if on {
		return s | (attrs << 32)
	} else {
		return s &^ (attrs << 32)
	}
}

// Normal returns the style with all attributes disabled.
func (s Style) Normal() Style {
	return s &^ (Style(0xfffff) << 32)
}

func (s Style) Bold(on bool) Style {
	return s.setAttrs(Style(AttrBold), on)
}

func (s Style) Blink(on bool) Style {
	return s.setAttrs(Style(AttrBlink), on)
}

func (s Style) Dim(on bool) Style {
	return s.setAttrs(Style(AttrDim), on)
}

func (s Style) Reverse(on bool) Style {
	return s.setAttrs(Style(AttrReverse), on)
}

func (s Style) Underline(on bool) Style {
	return s.setAttrs(Style(AttrUnderline), on)
}
