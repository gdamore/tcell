// Copyright 2015 The Tcell Authors
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

package views

import (
	"github.com/gdamore/tcell"
)

// SimpleStyledText is a form of Text that offers highlighting of the text
// using simple in-line markup.  Its intention is to make it easier to mark
// up hot // keys for menubars, etc.
type SimpleStyledText struct {
	styleN tcell.Style
	styleA tcell.Style
	styleB tcell.Style
	styleS tcell.Style
	styleU tcell.Style
	markup []rune
	init   bool
	Text
}

// SetMarkup sets the text used for the string.  It applies markup as follows
// (modeled on tcsh style prompt markup):
//
// * %% - emit a single % in current style
// * %N - normal style
// * %A - alternate style
// * %S - start standout (reverse) style
// * %B - start bold style
// * %U - start underline style
//
// Note that for simplicity, combining styles is not supported.  By default
// the alternate style is the same as standout (reverse) mode.
//
// Arguably we could have used Markdown syntax instead, but properly doing all
// of Markdown is not trivial, and these escape sequences make it clearer that
// we are not even attempting to do that.
func (t *SimpleStyledText) SetMarkup(s string) {

	markup := []rune(s)
	styl := make([]tcell.Style, 0, len(markup))
	text := make([]rune, 0, len(markup))

	style := t.StyleN()
	esc := false
	for _, r := range markup {
		if esc {
			esc = false
			switch r {
			case '%':
				text = append(text, '%')
				styl = append(styl, style)
			case 'N':
				style = t.StyleN()
			case 'A':
				style = t.StyleA()
			case 'B':
				style = t.StyleB()
			case 'S':
				style = t.StyleS()
			case 'U':
				style = t.StyleU()
			default:
				text = append(append(text, '%'), r)
				styl = append(append(styl, style), style)
			}
			continue
		}
		switch r {
		case '%':
			esc = true
			continue
		default:
			text = append(text, r)
			styl = append(styl, style)
		}
	}

	t.Text.SetText(string(text))
	for i, s := range styl {
		t.SetStyleAt(i, s)
	}
	t.markup = markup
}

// Markup returns the text that was set, including markup.
func (t *SimpleStyledText) Markup() string {
	return string(t.markup)
}

// SetStyleN sets the style used for N (normal).
func (t *SimpleStyledText) SetStyleN(style tcell.Style) {
	t.styleN = style
	t.Text.SetStyle(style)
}

// SetStyleA sets the style used for A (alternate).  Existing text is not
// changed, so call this before doing SetText.
func (t *SimpleStyledText) SetStyleA(style tcell.Style) {
	t.styleA = style
}

// SetStyleB sets the style used for B (bold).  Existing text is not
// changed, so call this before doing SetText.
func (t *SimpleStyledText) SetStyleB(style tcell.Style) {
	t.styleB = style
}

// SetStyleU sets the style used for U (underline).  Existing text is not
// changed, so call this before doing SetText.
func (t *SimpleStyledText) SetStyleU(style tcell.Style) {
	t.styleU = style
}

// SetStyleS sets the style used for S (standout).  Existing text is not
// changed, so call this before doing SetText.
func (t *SimpleStyledText) SetStyleS(style tcell.Style) {
	t.styleS = style
}

// StyleN returns the previously set N (normal) style.
func (t *SimpleStyledText) StyleN() tcell.Style {
	return t.styleN
}

// StyleA returns the previously set A (alternate) style.
func (t *SimpleStyledText) StyleA() tcell.Style {
	if t.styleA == tcell.StyleDefault {
		return t.styleN.Reverse(true).Bold(true)
	}
	return t.styleA
}

// StyleB returns the previously set B (bold) style.
func (t *SimpleStyledText) StyleB() tcell.Style {
	if t.styleB == tcell.StyleDefault {
		return t.styleN.Bold(true)
	}
	return t.styleB
}

// StyleS returns the previously set S (standout) style.
func (t *SimpleStyledText) StyleS() tcell.Style {
	if t.styleS == tcell.StyleDefault {
		return t.styleN.Reverse(true)
	}
	return t.styleS
}

// StyleU returns the previously set U (underline) style.
func (t *SimpleStyledText) StyleU() tcell.Style {
	if t.styleU == tcell.StyleDefault {
		return t.styleN.Underline(true)
	}
	return t.styleU
}

// NewSimpleStyledText creates an empty Text.
func NewSimpleStyledText() *SimpleStyledText {
	return &SimpleStyledText{}
}
