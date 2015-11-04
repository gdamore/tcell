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
	"strings"

	"github.com/gdamore/tcell"
)

// TextArea is a pannable 2 dimensional text widget. It wraps both
// the view and the model for text in a single, convenient widget.
// Text is provided as an array of strings, each of which represents
// a single line to display.  All text in the TextArea has the same
// style.  An optional soft cursor is available.
type TextArea struct {
	view  *CellView
	model *linesModel
}

type linesModel struct {
	lines  []string
	width  int
	height int
	x      int
	y      int
	hide   bool
	cursor bool
}

func (m *linesModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	var ch rune
	if x < 0 || y < 0 || y >= len(m.lines) || x >= len(m.lines[y]) {
		return ch, tcell.StyleDefault, nil, 1
	}
	// XXX: extend this to support combining and full width chars
	return rune(m.lines[y][x]), tcell.StyleDefault, nil, 1
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
	m := ta.model
	m.width = 0
	m.height = len(lines)
	m.lines = append([]string{}, lines...)
	for _, l := range lines {
		if len(l) > m.width {
			m.width = len(l)
		}
	}
	ta.view.SetModel(m)
}

// EnableCursor enables a soft cursor in the TextArea.
func (ta *TextArea) EnableCursor(on bool) {
	ta.model.cursor = on
}

// HideCursor hides or shows the cursor in the TextArea.
// If on is true, the cursor is hidden.  Note that a cursor is only
// shown if it is enabled.
func (ta *TextArea) HideCursor(on bool) {
	ta.model.hide = on
}

// Draw draws the TextArea.
func (ta *TextArea) Draw() {
	ta.view.Draw()
}

// HandleEvent handles any events.
func (ta *TextArea) HandleEvent(ev tcell.Event) bool {
	return ta.view.HandleEvent(ev)
}

// Resize is called when the drawing context (View) changes size.
func (ta *TextArea) Resize() {
	ta.view.Resize()
}

// SetView sets the drawing context.
func (ta *TextArea) SetView(view View) {
	ta.view.SetView(view)
}

// SetContent is used to set the textual content, passed as a
// single string.  Lines within the string are delimited by newlines.
func (ta *TextArea) SetContent(text string) {
	lines := strings.Split(strings.Trim(text, "\n"), "\n")
	ta.SetLines(lines)
}

// NewTextArea creates a blank TextArea.
func NewTextArea() *TextArea {
	lm := &linesModel{lines: []string{}, width: 0}
	ta := &TextArea{model: lm}
	ta.view = NewCellView()
	ta.view.SetModel(lm)
	return ta
}
