// Copyright 2025 The Tcell Authors
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

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/uniseg"
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
	strings [][]string
	widths  [][]int
	width   int
	height  int
	x       int
	y       int
	hide    bool
	cursor  bool
	style   tcell.Style
}

func (m *linesModel) Get(x, y int) (string, tcell.Style, int) {
	if x < 0 || y < 0 || y >= m.height || y >= len(m.strings) || x >= len(m.strings[y]) {
		return "", m.style, 1
	}
	return m.strings[y][x], m.style, m.widths[y][x]
}

func (m *linesModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	if x < 0 || y < 0 || y >= m.height || y >= len(m.strings) || x >= len(m.strings[y]) {
		return 0, m.style, nil, 1
	}
	var runes = []rune(m.strings[y][x])
	switch len(runes) {
	case 0:
		return 0, m.style, nil, m.widths[y][x]
	case 1:
		return runes[0], m.style, nil, m.widths[y][x]
	default:
		return runes[0], m.style, runes[1:], m.widths[y][x]
	}
}

func (m *linesModel) GetBounds() (int, int) {
	return m.width, m.height
}

func (m *linesModel) limitCursor() {
	if m.x > m.width-1 {
		m.x = m.width - 1
	}
	if m.y > m.height-1 {
		m.y = m.height - 1
	}
	if m.x < 0 {
		m.x = 0
	}
	if m.y < 0 {
		m.y = 0
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

	// extend slice before using m.runes[row] to avoid panic
	slice := make([][]string, len(lines))
	m.strings = slice
	m.widths = make([][]int, len(lines))

	for row, line := range lines {
		state := -1
		var cl string
		var w int
		rw := 0
		for line != "" {
			cl, line, w, state = uniseg.FirstGraphemeClusterInString(line, state)
			m.strings[row] = append(m.strings[row], cl)
			m.widths[row] = append(m.widths[row], w)
			rw += w
			m.width = max(rw, m.width)
		}
	}

	m.height = len(m.strings)

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
		lm := &linesModel{strings: [][]string{}, widths: [][]int{}}
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
