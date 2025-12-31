// Copyright 2026 The TCell Authors
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

package mock

import (
	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/vt"
)

// Cell is a representation of a display cell.
type Cell struct {
	C     []rune // Content, for now only a single rune is supported
	Fg    color.Color
	Bg    color.Color
	Attr  vt.Attr
	Width int // Display width of C.
}

// MockBackend provides additional mock-specific capabilities on top of Backend.
// This is meant to facilitate test cases
type MockBackend interface {
	vt.Backend
	vt.Blitter

	// GetCell returns the cell at the given position, or the zero value if the
	// position is out of the bounds of the window.
	GetCell(vt.Coord) Cell

	// Bells counts the number of bells rung.
	Bells() int

	// GetTitle gets the current window title.
	GetTitle() string

	// SetSize is used to resize the window.
	// Newly added cells are empty, and content in old cells that out of range is lost.
	SetSize(vt.Coord)
}

// mockBackend is a mock of a backend device for use with the emulator.
// It implements the following interfaces:
// vt.Backend, vt.Beeper, vt.Colorer, vt.Titler, vt.Resizer, vt.Blitter
type mockBackend struct {
	cells     []Cell // Content of cells
	size      vt.Coord
	pos       vt.Coord
	colors    int
	attr      vt.Attr
	fg        color.Color
	bg        color.Color
	defaultFg color.Color
	defaultBg color.Color
	resizeQ   chan<- bool
	modes     map[vt.PrivateMode]vt.ModeStatus
	bells     int
	errs      int
	title     string
}

func (mb *mockBackend) GetSize() vt.Coord { return mb.size }
func (mb *mockBackend) Beep()             { mb.bells++ }

func (mb *mockBackend) GetPrivateMode(pm vt.PrivateMode) vt.ModeStatus {
	// note default (zero) value is ModeNA
	return mb.modes[pm]
}

func (mb *mockBackend) SetPrivateMode(pm vt.PrivateMode, status vt.ModeStatus) error {
	if old := mb.modes[pm]; old == vt.ModeOn || old == vt.ModeOff {
		if status == vt.ModeOn || status == vt.ModeOff {
			mb.modes[pm] = status
		} else {
			mb.errs++
		}
	} else {
		mb.errs++
	}
	return nil
}

func (mb *mockBackend) PutRune(pos vt.Coord, r rune, width int) {
	if index := mb.index(pos); index >= 0 {
		if r == 0 {
			mb.cells[index].C = nil
			mb.cells[index].Width = 0
		} else {
			mb.cells[index].C = []rune{r}
			mb.cells[index].Width = width
		}
		mb.cells[index].Attr = mb.attr
		mb.cells[index].Fg = mb.fg
		mb.cells[index].Bg = mb.bg
		if width == 2 && pos.X < mb.size.X-1 {
			// wide characters delete the adjacent cell
			index++
			mb.cells[index].C = nil
			mb.cells[index].Width = 0
			mb.cells[index].Attr = mb.attr
			mb.cells[index].Fg = mb.fg
			mb.cells[index].Bg = mb.bg
		}
	} else {
		mb.errs++
	}
}

func (mb *mockBackend) PutGrapheme(pos vt.Coord, grapheme string, width int) {
	if index := mb.index(pos); index >= 0 {
		if grapheme == "" || width == 0 {
			mb.cells[index].C = nil
			mb.cells[index].Width = 0
		} else {
			mb.cells[index].C = []rune(grapheme)
			mb.cells[index].Width = width
		}
		mb.cells[index].Attr = mb.attr
		mb.cells[index].Fg = mb.fg
		mb.cells[index].Bg = mb.bg
		if width == 2 && pos.X < mb.size.X-1 {
			// wide characters delete the adjacent cell
			index++
			mb.cells[index].C = nil
			mb.cells[index].Width = 0
			mb.cells[index].Attr = mb.attr
			mb.cells[index].Fg = mb.fg
			mb.cells[index].Bg = mb.bg
		}
	} else {
		mb.errs++
	}
}

func (mb *mockBackend) isPositionValid(pos vt.Coord) bool {
	return pos.X < mb.size.X && pos.Y < mb.size.Y && pos.X >= 0 && pos.Y >= 0
}

// index calculates the index in the cells array.  If the coordinates are invalid,
// -1 will be returned.
func (mb *mockBackend) index(pos vt.Coord) int {
	if !mb.isPositionValid(pos) {
		return -1
	}
	return int(pos.X) + int(pos.Y)*int(mb.size.X)
}

func (mb *mockBackend) GetCell(pos vt.Coord) Cell {
	if index := mb.index(pos); index >= 0 {
		return mb.cells[index]
	}
	return Cell{}
}

func (mb *mockBackend) Bells() int {
	return mb.bells
}

func (mb *mockBackend) GetPosition() vt.Coord {
	return mb.pos
}

func (mb *mockBackend) SetPosition(pos vt.Coord) {
	pos.X = min(mb.size.X-1, max(0, pos.X))
	pos.Y = min(mb.size.Y-1, max(0, pos.Y))
	mb.pos = pos
}

// setColor is a helper for setting color values.
func (mb *mockBackend) setColor(c color.Color, dst *color.Color, def color.Color) {
	if mb.colors == 0 {
		return
	}
	if c.Valid() {
		if c.IsRGB() {
			if mb.colors > 256 {
				*dst = c
			}
		} else if (int(c) & 255) < mb.colors {
			*dst = c
		}
		return
	} else if c == color.Reset {
		*dst = def
	}
}

func (mb *mockBackend) Colors() int {
	return mb.colors
}

func (mb *mockBackend) SetFgColor(c color.Color) {
	mb.setColor(c, &mb.fg, mb.defaultFg)
}

func (mb *mockBackend) SetBgColor(c color.Color) {
	mb.setColor(c, &mb.bg, mb.defaultBg)
}

func (mb *mockBackend) SetAttr(attr vt.Attr) {
	mb.attr = attr
}

// SetWindowTitle implements the Titler interface.
func (mb *mockBackend) SetWindowTitle(title string) {
	mb.title = title
}

// GetTitle allows test code to observe what was set with SetWindowTitle.
func (mb *mockBackend) GetTitle() string {
	return mb.title
}

// NotifyResize registers a channel to be written to (non-blocking) if the
// backend changes size.
func (mb *mockBackend) NotifyResize(rq chan<- bool) {
	mb.resizeQ = rq
}

// SetSize is used to change the size of the virtual terminal. Cells that are
// added are treated as empty, while cells that are removed are just lost.
// (Note that at least one other emulator erases content on a resize.  There is
// standard for what to do here.)
func (mb *mockBackend) SetSize(size vt.Coord) {
	old := mb.cells
	ox := int(mb.size.X)
	oy := int(mb.size.Y)
	nx := int(size.X)
	ny := int(size.Y)
	cells := make([]Cell, int(size.Y)*int(size.X))
	for y := range min(ny, oy) {
		for x := range min(nx, ox) {
			cells[y*nx+x] = old[y*ox+x]
		}
	}
	mb.cells = cells
	mb.size = size
	mb.pos.X = min(mb.pos.X, size.X-1)
	mb.pos.Y = min(mb.pos.Y, size.Y-1)
	if rq := mb.resizeQ; rq != nil {
		select {
		case rq <- true:
		default:
		}
	}
}

// Reset the terminal to startup defaults.
func (mb *mockBackend) Reset() {
	mb.fg = mb.defaultFg
	mb.bg = mb.defaultBg
	mb.attr = vt.Plain
	mb.title = ""
	mb.errs = 0
	mb.bells = 0
	mb.pos = vt.Coord{X: 0, Y: 0}
	mb.modes[vt.PmShowCursor] = vt.ModeOn
	mb.modes[vt.PmGraphemeClusters] = vt.ModeOff
}

func (mb *mockBackend) Blit(src, dst, dim vt.Coord) {
	// clip to visible source
	if dim.X+src.X > mb.size.X {
		dim.X = mb.size.X - src.X
	}
	if dim.Y+src.Y > mb.size.Y {
		dim.Y = mb.size.Y - src.Y
	}
	// and clip to final destination
	if dim.X+dst.X > mb.size.X {
		dim.X = mb.size.X - dst.X
	}
	if dim.Y+dst.Y > mb.size.Y {
		dim.Y = mb.size.Y - dst.Y
	}

	// gap represents decrement when shifting to the next row --
	// skipping over the irrelevant cells. (The increment in the
	// index when going from last cell of row to first cell of next row,
	// or vice versa.)
	gap := int(mb.size.X - dim.X)

	// the following logic is carefully constructed to avoid expensive
	// operations in the loops (only addition or subtraction)
	if mb.index(src) > mb.index(dst) { // source appears later, so we can forward copy
		si := mb.index(src)
		di := mb.index(dst)
		for range dim.Y {
			for range dim.X {
				mb.cells[di] = mb.cells[si]
				di++
				si++
			}
			// advance to next row
			si += gap
			di += gap
		}
	} else { // source appears later, so we have to reverse copy
		src.Y += dim.Y - 1
		dst.Y += dim.Y - 1
		src.X += dim.X - 1
		dst.X += dim.X - 1
		si := mb.index(src)
		di := mb.index(dst)

		for range dim.Y {
			for range dim.X {
				mb.cells[di] = mb.cells[si]
				si--
				di--
			}
			si -= gap
			di -= gap
		}
	}
}

// MockOpt is an interface by which options can change the behavior of the mocked terminal.
// This is intended to permit easier testing.
type MockOpt interface{ SetMockOpt(mb *mockBackend) }

// MockOptSize changes the default terminal size, which is normally 80x24.
type MockOptSize vt.Coord

func (o MockOptSize) SetMockOpt(mb *mockBackend) { mb.size = vt.Coord(o) }

// MockOptColors changes the number of colors the terminal supports.
type MockOptColors int

func (o MockOptColors) SetMockOpt(mb *mockBackend) { mb.colors = int(o) }

// NewMockBackend returns a MockBackend modified by the given options.
// The default is a fully featured 256-color backend with initial size 80x24.
func NewMockBackend(options ...MockOpt) MockBackend {
	mb := &mockBackend{
		size:      vt.Coord{X: 80, Y: 24},
		colors:    256,
		defaultFg: color.Silver,
		defaultBg: color.Black,
	}

	for _, opt := range options {
		opt.SetMockOpt(mb)
		// TODO: possibly be could be "filtered" for some options (e.g. to hide colorer API, etc.)
	}

	if mb.colors > 0 {
		mb.fg = mb.defaultFg
		mb.bg = mb.defaultBg
	}
	mb.cells = make([]Cell, int(mb.size.X)*int(mb.size.Y))
	mb.modes = make(map[vt.PrivateMode]vt.ModeStatus)
	mb.modes[vt.PmShowCursor] = vt.ModeOn
	mb.modes[vt.PmGraphemeClusters] = vt.ModeOff
	return mb
}
