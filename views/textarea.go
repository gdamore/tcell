// Copyright 2016 The Tcell Authors
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
	"strings"
	"sync"

	"github.com/gdamore/tcell"
)

// TextArea is a pannable 2 dimensional text widget. It wraps both
// the view and the model for text in a single, convenient widget.
// Text is provided as an array of strings, each of which represents
// a single line to display.  All text in the TextArea has the same
// style.  An optional soft cursor is available.
type TextArea struct {
	model *linesModel
	once  sync.Once
	CellView
}

type linesModel struct {
	lines  []string
	width  int
	height int
	x      int
	y      int
	hide   bool
	cursor bool
	style  tcell.Style
}

func (m *linesModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	var ch rune
	if x < 0 || y < 0 || y >= len(m.lines) || x >= len(m.lines[y]) {
		return ch, m.style, nil, 1
	}
	// XXX: extend this to support combining and full width chars
	return rune(m.lines[y][x]), m.style, nil, 1
}

func (m *linesModel) GetBounds() (int, int) {
	return m.width, m.height
}

func (m *linesModel) limitCursor() {
	if m.x < 0 {
		m.x = 0
	}
	if m.y < 0 {
		m.y = 0
	}
	if m.x > m.width-1 {
		m.x = m.width - 1
	}
	if m.y > m.height-1 {
		m.y = m.height - 1
	}
}
func (m *linesModel) SetCursor(x, y int) {
	m.x = x
	m.y = y
	m.limitCursor()
}

func (m *linesModel) MoveCursor(x, y int) {
	m.x += x
	m.y += y
	m.limitCursor()
}

func (m *linesModel) GetCursor() (int, int, bool, bool) {
	return m.x, m.y, m.cursor, !m.hide
}

// SetLines sets the content text to display.
func (ta *TextArea) SetLines(lines []string) {
	ta.Init()
	m := ta.model
	m.width = 0
	m.height = len(lines)
	m.lines = append([]string{}, lines...)
	for _, l := range lines {
		if len(l) > m.width {
			m.width = len(l)
		}
	}
	ta.CellView.SetModel(m)
}

func (ta *TextArea) SetStyle(style tcell.Style) {
	ta.model.style = style
	ta.CellView.SetStyle(style)
}

// EnableCursor enables a soft cursor in the TextArea.
func (ta *TextArea) EnableCursor(on bool) {
	ta.Init()
	ta.model.cursor = on
}

// HideCursor hides or shows the cursor in the TextArea.
// If on is true, the cursor is hidden.  Note that a cursor is only
// shown if it is enabled.
func (ta *TextArea) HideCursor(on bool) {
	ta.Init()
	ta.model.hide = on
}

// SetContent is used to set the textual content, passed as a
// single string.  Lines within the string are delimited by newlines.
func (ta *TextArea) SetContent(text string) {
	ta.Init()
	lines := strings.Split(strings.Trim(text, "\n"), "\n")
	ta.SetLines(lines)
}

// Init initializes the TextArea.
func (ta *TextArea) Init() {
	ta.once.Do(func() {
		lm := &linesModel{lines: []string{}, width: 0}
		ta.model = lm
		ta.CellView.Init()
		ta.CellView.SetModel(lm)
	})
}

// NewTextArea creates a blank TextArea.
func NewTextArea() *TextArea {
	ta := &TextArea{}
	ta.Init()
	return ta
}
