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

package vt

import (
	"time"

	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/tty"
)

// mockTerm implements MockTerm.
type mockTerm struct {
	mb MockBackend
	em Emulator
}

// Stop the terminal.
func (mt *mockTerm) Stop() error {
	return mt.em.Stop()
}

// Start the terminal.
func (mt *mockTerm) Start() error {
	return mt.em.Start()
}

// Drain all output from the terminal, ensuring
// any queued commands are processed.
func (mt *mockTerm) Drain() error {
	return mt.em.Drain()
}

// Read data from the terminal. This is called by a terminal
// application (e.g. via tcell Tty.)  Read data will include
// key strokes, mouse events, and responses to terminal queries.
func (mt *mockTerm) Read(data []byte) (int, error) {
	return mt.em.Read(data)
}

// Write data to the terminal, typically either commands or data
// that should be displayed on the virtual screen.
func (mt *mockTerm) Write(b []byte) (n int, err error) {
	return mt.em.Write(b)
}

// WindowSize obtains the dimensions of the window.
func (mt *mockTerm) WindowSize() (tty.WindowSize, error) {
	sz := mt.mb.GetSize()
	// No pixel sizes for now
	return tty.WindowSize{Width: int(sz.X), Height: int(sz.Y)}, nil
}

// NotifyResize registers a channel to be signaled when a resize has occurred.
// In real terminal emulators this would be posted (non-blocking) by a signal handler.
func (mt *mockTerm) NotifyResize(resizeq chan<- bool) {
	if rs, ok := mt.mb.(Resizer); ok {
		rs.NotifyResize(resizeq)
	}
}

// Close closes the terminal, after which it should no longer be used. Stop is implied.
func (mt *mockTerm) Close() error {
	return mt.Stop()
}

// Pos returns the cursor position.
func (mt *mockTerm) Pos() Coord {
	return mt.mb.GetPosition()
}

// GetCell returns the contents of the cell at the given coordinates, or a zero value
// if the coordinates are out of range.
func (mt *mockTerm) GetCell(pos Coord) MockCell {
	return mt.mb.GetCell(pos)
}

// Bells counts the number of times the bell has rung.
func (mt *mockTerm) Bells() int {
	return mt.mb.Bells()
}

func (mt *mockTerm) KeyEvent(ev KbdEvent) {
	mt.em.KeyEvent(ev)
	if ev.Code == KcEsc {
		// Inject a delay to simulate human typing.
		// Necessary to disambiguate Escape from other sequences.
		time.Sleep(time.Millisecond * 150)
	}
}

// GetTitle returns the current window title.
func (mt *mockTerm) GetTitle() string {
	return mt.mb.GetTitle()
}

// SetSize is used to change the terminal size.
func (mt *mockTerm) SetSize(size Coord) {
	mt.mb.SetSize(size)
	mt.em.ResizeEvent()
}

// Backend returns the backend for testing.
func (mt *mockTerm) Backend() MockBackend {
	return mt.mb
}

// MockTerm is a mock terminal (emulator).  It can be used to
// test the emulator itself, or to test applications (or tcell) that
// uses the terminal.  It also implements the Tty interface used
// by tcell itself.
type MockTerm interface {
	tty.Tty

	// Pos reports the current cursor position.
	Pos() Coord

	// GetCell returns the current cell.
	GetCell(Coord) MockCell

	// Bells returns the number of times the bell has been rung.
	Bells() int

	// Inject a keyboard event.
	KeyEvent(KbdEvent)

	// GetTitle obtains the current window title.
	GetTitle() string

	// SetSize is used to resize the terminal.
	SetSize(Coord)

	// Backend returns the backend (used for testing).
	Backend() MockBackend
}

// NewMockTerm gives a mock terminal emulator.
func NewMockTerm(opts ...MockOpt) MockTerm {
	mt := &mockTerm{}
	mt.mb = NewMockBackend(opts...)
	mt.em = NewEmulator(mt.mb)
	mt.em.SetId("TcellMock", "1.0")
	return mt
}

// MockCell is a representation of a display cell.
// It only adds a width so we can use that for verification.
type MockCell struct {
	Cell
	Width int // Display width of C.
}

// MockBackend provides additional mock-specific capabilities on top of Backend.
// This is meant to facilitate test cases
type MockBackend interface {
	Backend
	Blitter

	// GetCell returns the cell at the given position, or the zero value if the
	// position is out of the bounds of the window.
	GetCell(Coord) MockCell

	// Bells counts the number of bells rung.
	Bells() int

	// GetTitle gets the current window title.
	GetTitle() string

	// SetSize is used to resize the window.
	// Newly added cells are empty, and content in old cells that out of range is lost.
	SetSize(Coord)
}

// mockBackend is a mock of a backend device for use with the emulator.
// It implements the following interfaces:
// vt.Backend, vt.Beeper, vt.Colorer, vt.Titler, vt.Resizer, vt.Blitter
type mockBackend struct {
	cells        []MockCell // Content of cells
	size         Coord
	pos          Coord
	colors       int
	style        Style
	defaultStyle Style
	resizeQ      chan<- bool
	modes        map[PrivateMode]ModeStatus
	bells        int
	errs         int
	title        string
}

func (mb *mockBackend) GetSize() Coord { return mb.size }
func (mb *mockBackend) Beep()          { mb.bells++ }

func (mb *mockBackend) GetPrivateMode(pm PrivateMode) ModeStatus {
	// note default (zero) value is ModeNA
	return mb.modes[pm]
}

func (mb *mockBackend) SetPrivateMode(pm PrivateMode, status ModeStatus) error {
	if old := mb.modes[pm]; old == ModeOn || old == ModeOff {
		if status == ModeOn || status == ModeOff {
			mb.modes[pm] = status
		} else {
			mb.errs++
		}
	} else {
		mb.errs++
	}
	return nil
}

func (mb *mockBackend) PutRune(pos Coord, r rune, width int) {
	if index := mb.index(pos); index >= 0 {
		if r == 0 {
			mb.cells[index].C = ""
			mb.cells[index].Width = 0
		} else {
			mb.cells[index].C = string(r)
			mb.cells[index].Width = width
		}
		mb.cells[index].S = mb.style
		if width == 2 && pos.X < mb.size.X-1 {
			// wide characters delete the adjacent cell
			index++
			mb.cells[index].C = ""
			mb.cells[index].Width = 0
			mb.cells[index].S = mb.style
		}
	} else {
		mb.errs++
	}
}

func (mb *mockBackend) PutGrapheme(pos Coord, grapheme string, width int) {
	if index := mb.index(pos); index >= 0 {
		if grapheme == "" || width == 0 {
			mb.cells[index].C = ""
			mb.cells[index].Width = 0
		} else {
			mb.cells[index].C = grapheme
			mb.cells[index].Width = width
		}
		mb.cells[index].S = mb.style
		if width == 2 && pos.X < mb.size.X-1 {
			// wide characters delete the adjacent cell
			index++
			mb.cells[index].C = ""
			mb.cells[index].Width = 0
			mb.cells[index].S = mb.style
		}
	} else {
		mb.errs++
	}
}

func (mb *mockBackend) isPositionValid(pos Coord) bool {
	return pos.X < mb.size.X && pos.Y < mb.size.Y && pos.X >= 0 && pos.Y >= 0
}

// index calculates the index in the cells array.  If the coordinates are invalid,
// -1 will be returned.
func (mb *mockBackend) index(pos Coord) int {
	if !mb.isPositionValid(pos) {
		return -1
	}
	return int(pos.X) + int(pos.Y)*int(mb.size.X)
}

func (mb *mockBackend) GetCell(pos Coord) MockCell {
	if index := mb.index(pos); index >= 0 {
		return mb.cells[index]
	}
	return MockCell{}
}

func (mb *mockBackend) Bells() int {
	return mb.bells
}

func (mb *mockBackend) GetPosition() Coord {
	return mb.pos
}

func (mb *mockBackend) SetPosition(pos Coord) {
	pos.X = min(mb.size.X-1, max(0, pos.X))
	pos.Y = min(mb.size.Y-1, max(0, pos.Y))
	mb.pos = pos
}

func (mb *mockBackend) Colors() int {
	return mb.colors
}

func (mb *mockBackend) SetStyle(style Style) {
	mb.style = style
}

func (mb *mockBackend) GetStyle() Style {
	return mb.style
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
func (mb *mockBackend) SetSize(size Coord) {
	old := mb.cells
	ox := int(mb.size.X)
	oy := int(mb.size.Y)
	nx := int(size.X)
	ny := int(size.Y)
	cells := make([]MockCell, int(size.Y)*int(size.X))
	for i := range cells {
		cells[i].S = BaseStyle
	}
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
	mb.style = mb.defaultStyle

	mb.title = ""
	mb.errs = 0
	mb.bells = 0
	mb.pos = Coord{X: 0, Y: 0}
	mb.modes[PmShowCursor] = ModeOn
	mb.modes[PmGraphemeClusters] = ModeOff
}

func (mb *mockBackend) Blit(src, dst, dim Coord) {
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
type MockOptSize Coord

func (o MockOptSize) SetMockOpt(mb *mockBackend) { mb.size = Coord(o) }

// MockOptColors changes the number of colors the terminal supports.
type MockOptColors int

func (o MockOptColors) SetMockOpt(mb *mockBackend) { mb.colors = int(o) }

// NewMockBackend returns a MockBackend modified by the given options.
// The default is a fully featured 256-color backend with initial size 80x24.
func NewMockBackend(options ...MockOpt) MockBackend {
	mb := &mockBackend{
		size:         Coord{X: 80, Y: 24},
		colors:       256,
		style:        BaseStyle,
		defaultStyle: BaseStyle.WithFg(color.Silver).WithBg(color.Black),
	}

	for _, opt := range options {
		opt.SetMockOpt(mb)
		// TODO: possibly be could be "filtered" for some options (e.g. to hide colorer API, etc.)
	}

	if mb.colors > 0 {
		mb.style = mb.defaultStyle
	}
	mb.cells = make([]MockCell, int(mb.size.X)*int(mb.size.Y))
	for i := range mb.cells {
		mb.cells[i].S = BaseStyle
	}

	mb.modes = make(map[PrivateMode]ModeStatus)
	mb.modes[PmShowCursor] = ModeOn
	mb.modes[PmGraphemeClusters] = ModeOff
	return mb
}
