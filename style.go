// Copyright 2024 The TCell Authors
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

// Style represents a complete text style, including both foreground color,
// background color, and additional attributes such as "bold" or "underline".
//
// Note that not all terminals can display all colors or attributes, and
// many might have specific incompatibilities between specific attributes
// and color combinations.
//
// To use Style, just declare a variable of its type.
type Style struct {
	fg    Color
	bg    Color
	under Color
	attrs AttrMask
	url   string
	urlId string
}

// StyleDefault represents a default style, based upon the context.
// It is the zero value.
var StyleDefault Style

// styleInvalid is just an arbitrary invalid style used internally.
var styleInvalid = Style{attrs: AttrInvalid}

// Foreground returns a new style based on s, with the foreground color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Foreground(c Color) Style {
	s2 := s
	s2.fg = c
	return s2
}

// Background returns a new style based on s, with the background color set
// as requested.  ColorDefault can be used to select the global default.
func (s Style) Background(c Color) Style {
	s2 := s
	s2.bg = c
	return s2
}

// Decompose breaks a style up, returning the foreground, background,
// and other attributes.  The URL if set is not included.
// Deprecated: Applications should not attempt to decompose style,
// as this content is not sufficient to describe the actual style.
func (s Style) Decompose() (fg Color, bg Color, attr AttrMask) {
	return s.fg, s.bg, s.attrs
}

func (s Style) setAttrs(attrs AttrMask, on bool) Style {
	s2 := s
	if on {
		s2.attrs |= attrs
	} else {
		s2.attrs &^= attrs
	}
	return s2
}

// Normal returns the style with all attributes disabled.
func (s Style) Normal() Style {
	return Style{
		fg: s.fg,
		bg: s.bg,
	}
}

// Bold returns a new style based on s, with the bold attribute set
// as requested.
func (s Style) Bold(on bool) Style {
	return s.setAttrs(AttrBold, on)
}

// Blink returns a new style based on s, with the blink attribute set
// as requested.
func (s Style) Blink(on bool) Style {
	return s.setAttrs(AttrBlink, on)
}

// Dim returns a new style based on s, with the dim attribute set
// as requested.
func (s Style) Dim(on bool) Style {
	return s.setAttrs(AttrDim, on)
}

// Italic returns a new style based on s, with the italic attribute set
// as requested.
func (s Style) Italic(on bool) Style {
	return s.setAttrs(AttrItalic, on)
}

// Reverse returns a new style based on s, with the reverse attribute set
// as requested.  (Reverse usually changes the foreground and background
// colors.)
func (s Style) Reverse(on bool) Style {
	return s.setAttrs(AttrReverse, on)
}

// Underline returns a new style based on s, with the underline attribute set
// as requested.
func (s Style) Underline(on bool) Style {
	return s.setAttrs(AttrUnderline, on)
}

// StrikeThrough sets strikethrough mode.
func (s Style) StrikeThrough(on bool) Style {
	return s.setAttrs(AttrStrikeThrough, on)
}

func (s Style) DoubleUnderline(on bool) Style {
	return s.setAttrs(AttrDoubleUnderline, on)
}

func (s Style) CurlyUnderline(on bool) Style {
	return s.setAttrs(AttrCurlyUnderline, on)
}

func (s Style) DottedUnderline(on bool) Style {
	return s.setAttrs(AttrDottedUnderline, on)
}

func (s Style) DashedUnderline(on bool) Style {
	return s.setAttrs(AttrDashedUnderline, on)
}

func (s Style) UnderlineColor(c Color) Style {
	s2 := s
	s2.under = c
	return s2
}

// Attributes returns a new style based on s, with its attributes set as
// specified.
func (s Style) Attributes(attrs AttrMask) Style {
	s2 := s
	s2.attrs = attrs
	return s2
}

// Url returns a style with the Url set.  If the provided Url is not empty,
// and the terminal supports it, text will typically be marked up as a clickable
// link to that Url.  If the Url is empty, then this mode is turned off.
func (s Style) Url(url string) Style {
	s2 := s
	s2.url = url
	return s2
}

// UrlId returns a style with the UrlId set. If the provided UrlId is not empty,
// any marked up Url with this style will be given the UrlId also. If the
// terminal supports it, any text with the same UrlId will be grouped as if it
// were one Url, even if it spans multiple lines.
func (s Style) UrlId(id string) Style {
	s2 := s
	s2.urlId = "id=" + id
	return s2
}
