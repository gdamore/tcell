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

// Style represents a complete text style, including both foreground
// and background color.  We encode it in a 64-bit int for efficiency.
// The coding is (MSB): <16b rsvd><16b attr><16b fgcolor><16b bgcolor>.
// This gives 16bit color options, if it ever becomes truly necessary.
// However, applications must not rely on this encoding.
//
// Note that not all terminals can display all colors or attributes, and
// many might have specific incompatibilities between specific attributes
// and color combinations.
//
// To use Style, just declare a variable of its type.
type Style int64

// StyleDefault represents a default style, based upon the context.
// It is the zero value.
const StyleDefault Style = 0

// Foreground returns a new style based on s, with the foreground color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Foreground(c Color) Style {
	return (s &^ Style(0xffff0000)) | (Style(c) << 16)
}

// Background returns a new style based on s, with the background color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Background(c Color) Style {
	return (s &^ (0xffff)) | Style(c)
}

// Decompose breaks a style up, returning the foreground, background,
// and other attributes.
func (s Style) Decompose() (fg Color, bg Color, attr AttrMask) {
	return Color((s >> 16) & 0xffff),
		Color(s & 0xfffff),
		AttrMask((s >> 32) & 0xffff)
}

func (s Style) setAttrs(attrs Style, on bool) Style {
	if on {
		return s | (attrs << 32)
	}
	return s &^ (attrs << 32)
}

// Normal returns the style with all attributes disabled.
func (s Style) Normal() Style {
	return s &^ (Style(0xfffff) << 32)
}

// Bold returns a new style based on s, with the bold attribute set
// as requested.
func (s Style) Bold(on bool) Style {
	return s.setAttrs(Style(AttrBold), on)
}

// Blink returns a new style based on s, with the blink attribute set
// as requested.
func (s Style) Blink(on bool) Style {
	return s.setAttrs(Style(AttrBlink), on)
}

// Dim returns a new style based on s, with the dim attribute set
// as requested.
func (s Style) Dim(on bool) Style {
	return s.setAttrs(Style(AttrDim), on)
}

// Reverse returns a new style based on s, with the reverse attribute set
// as requested.  (Reverse usually changes the foreground and background
// colors.)
func (s Style) Reverse(on bool) Style {
	return s.setAttrs(Style(AttrReverse), on)
}

// Underline returns a new style based on s, with the underline attribute set
// as requested.
func (s Style) Underline(on bool) Style {
	return s.setAttrs(Style(AttrUnderline), on)
}
