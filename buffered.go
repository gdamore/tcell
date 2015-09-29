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

import (
	"sync"

	"github.com/mattn/go-runewidth"
)

// BufferedScreen is like the base screen, but is buffered.  Each screen
// represents the root window that applications interface with.
type BufferedScreen interface {
	// Init initializes the screen for use.
	Init() error

	// Fini finazlizes the screen also releasing resources.
	Fini()

	// Clear erases the screen.
	Clear()

	// SetCell sets the cell at the given location.
	// The ch array contains at most one rune of width > 0, and the
	// runes with zero width (combining marks) must follow the first
	// non-zero width character.  (If only combining marks are present,
	// a space character will be filled in.)
	//
	// Note that double wide runes occupy two cells, and attempts to
	// place a character at the immediately adjacent cell will have
	// undefined effects.  Double wide runes that are printed in the
	// last column will be replaced with a single width space on output.
	SetCell(x int, y int, style Style, ch ...rune)

	// ShowCursor is used to display the cursor at a given location.
	ShowCursor(x int, y int)

	// HideCursor is used to hide the cursor.
	HideCursor()

	// Size returns the screen size as width, height.  This changes in
	// response to a call to Clear or Flush.
	Size() (int, int)

	// PollEvent waits for events to arrive.  Main application loops
	// can generally spin on this.
	PollEvent() Event

	// PostEvent posts an event into the event stream.
	PostEvent(Event)

	// Colors returns the number of colors.  All colors are assumed to
	// use the ANSI color map.
	Colors() int

	// EnableMouse enables mouse events, if your terminal has support
	// for them.
	EnableMouse()

	// DisableMouse disables mouse events.
	DisableMouse()

	// Sync synchronizes the buffered content with the screen, without
	// making any assumptions about the content that is displayed.
	// This is most often useful when some other program has altered the
	// screen state.  Because this is a full redraw, it can be visually
	// jarring & expensive, and should only be done when truly needed.
	Sync()

	// Show writes the contents of the buffer to the physical screen.
	// Only the contents that have changed will be written.  This is what
	// applications call to redraw the screen.
	Show()
}

type cell struct {
	ch    []rune
	dirty bool
	width uint8
	style Style
}

type bScreen struct {
	s       Screen
	cells   []cell
	w       int
	h       int
	cursorx int
	cursory int
	clear   bool

	sync.Mutex
}

func (b *bScreen) Init() error {
	if e := b.s.Init(); e != nil {
		return e
	}

	// Allocate cells
	b.w, b.h = b.s.Size()
	b.cells = make([]cell, b.w*b.h)
	b.cursorx = -1
	b.cursory = -1
	return nil
}

func (b *bScreen) Fini() {
	b.Lock()
	b.cells = nil
	b.w = 0
	b.h = 0
	b.clear = false
	b.cursorx = -1
	b.cursory = -1
	b.Unlock()
	b.s.Fini()
}

func (b *bScreen) Colors() int {
	return b.s.Colors()
}

func (b *bScreen) Size() (int, int) {
	b.Lock()
	w, h := b.w, b.h
	b.Unlock()
	return w, h
}

func (b *bScreen) PollEvent() Event {
	ev := b.s.PollEvent()

	// We need to capture resize events from the bottom screen, so that
	// we can readjust our buffer. This is important since the we need to
	// do any adjustment before the application starts drawing in response,
	// or coordinates may be erroneously believed out of range and results
	// discarded.
	if _, ok := ev.(*EventResize); ok {
		b.Lock()
		b.checkResize()
		b.Unlock()
	}
	return ev
}

func (b *bScreen) PostEvent(ev Event) {
	// See comment in PollEvent for why we do this.  Note that normally
	// events are posted directly into the screen below, so these are only
	// application supplied events.
	if _, ok := ev.(*EventResize); ok {
		b.Lock()
		b.checkResize()
		b.Unlock()
	}
	b.s.PostEvent(ev)
}

func (b *bScreen) Clear() {
	b.Lock()
	for i := range b.cells {
		b.cells[i].dirty = true
		b.cells[i].style = StyleDefault
		b.cells[i].ch = nil
	}
	b.Unlock()
}

func (b *bScreen) HideCursor() {
	b.Lock()
	b.cursorx = -1
	b.cursory = -1
	b.Unlock()
}

func (b *bScreen) ShowCursor(x, y int) {
	b.Lock()
	b.cursorx = x
	b.cursory = y
	b.Unlock()
}

func (b *bScreen) checkResize() {
	// Must be called with lock held!
	w, h := b.s.Size()
	if w == b.w && h == b.h {
		return
	}
	// We could reuse the cells, if we knew that both the row size did not
	// increase, and the total size did not increase.  For now we just
	// take the lazy approach and grow a new cells structure.  It would be
	// bad if the window size changes very frequently, but that shouldn't
	// happen.
	newc := make([]cell, w*h)
	for row := 0; row < h && row < b.h; row++ {
		for col := 0; col < w && col < b.w; col++ {
			newc[(row*w)+col] = b.cells[(row*b.w)+col]
			newc[(row*w)+col].dirty = true
		}
	}
	b.w = w
	b.h = h
	b.cells = newc
	// force a full screen redraw - just to be sure
}

func (b *bScreen) SetCell(x int, y int, style Style, ch ...rune) {
	// compare ch, compare style
	b.Lock()
	if x < 0 || y < 0 || x >= b.w || y >= b.h {
		b.Unlock()
		return
	}
	cell := &b.cells[(y*b.w)+x]

	// check to see if its the same value, if it is, don't mark it dirty
	match := false
	if len(ch) == len(cell.ch) && style == cell.style {
		match = true
		for i, r := range cell.ch {
			if ch[i] != r {
				match = false
				break
			}
		}
	}
	if !match {
		cell.dirty = true
		cell.ch = ch
		cell.style = style
		cell.width = 1
		for i := range cell.ch {
			if runewidth.RuneWidth(ch[i]) == 2 {
				cell.width = 2
			}
		}
	}
	b.Unlock()
}

func (b *bScreen) Show() {
	b.Lock()

	b.checkResize()

	b.s.HideCursor()

	if b.clear {
		b.s.Clear()
		b.clear = false
	}

	for row := 0; row < b.h; row++ {
		for col := 0; col < b.w; col++ {
			c := &b.cells[(row*b.w)+col]
			if !c.dirty {
				continue
			}
			b.s.SetCell(col, row, c.style, c.ch...)
			c.dirty = false
			if c.width == 2 {
				col++
			}
		}
	}

	b.s.ShowCursor(b.cursorx, b.cursory)
	b.Unlock()
}

func (b *bScreen) Sync() {
	b.Lock()
	// do a clear screen and also mark everything dirty
	b.clear = true
	for i := range b.cells {
		b.cells[i].dirty = true
	}
	b.Unlock()
	b.Show()
}

func (b *bScreen) EnableMouse() {
	b.Lock()
	b.s.EnableMouse()
	b.Unlock()
}

func (b *bScreen) DisableMouse() {
	b.Lock()
	b.s.DisableMouse()
	b.Unlock()
}

func NewBufferedScreen() (BufferedScreen, error) {
	s, e := NewScreen()
	if e != nil {
		return nil, e
	}
	return MakeBufferedScreen(s), nil
}

func MakeBufferedScreen(s Screen) BufferedScreen {
	return &bScreen{s: s}
}
