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
	"errors"
	"sync"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

// NewSimulationScreen returns a SimulationScreen.  Note that
// SimulationScreen is also a Screen.
func NewSimulationScreen(charset string) SimulationScreen {
	if charset == "" {
		charset = "UTF-8"
	}
	s := &simscreen{charset: charset}
	return s
}

// SimulationScreen represents a screen simulation.  This is intended to
// be a superset of normal Screens, but also adds some important interfaces
// for testing.
type SimulationScreen interface {
	// InjectKeyBytes injects a stream of bytes corresponding to
	// the native encoding (see charset).  It turns true if the entire
	// set of bytes were processed and delivered as KeyEvents, false
	// if any bytes were not fully understood.  Any bytes that are not
	// fully converted are discarded.
	InjectKeyBytes(buf []byte) bool

	// InjectKey injects a key event.  The rune is a UTF-8 rune, post
	// any translation.
	InjectKey(key Key, r rune, mod ModMask)

	// InjectMouse injects a mouse event.
	InjectMouse(x, y int, buttons ButtonMask, mod ModMask)

	// Resize resizes the underlying physical screen.  It also causes
	// a resize event to be injected during the next Show() or Sync().
	// A new physical contents array will be allocated (with data from
	// the old copied), so any prior value obtained with GetContents
	// won't be used anymore
	Resize(width, height int)

	// GetContents returns screen contents as an array of
	// cells, along with the physical width & height.   Note that the
	// physical contents will be used until the next time Resize()
	// is called.
	GetContents() (cells []SimCell, width int, height int)

	// GetCursor returns the cursor details.
	GetCursor() (x int, y int, visible bool)

	Screen
}

// SimCell represents a simulated screen cell.  The purpose of this
// is to track on screen content.
type SimCell struct {
	// Bytes is the actual character bytes.  Normally this is
	// rune data, but it could be be data in another encoding system.
	Bytes []byte

	// Style is the style used to display the data.
	Style Style

	// Runes is the list of runes, unadulterated, in UTF-8.
	Runes []rune
}

type simscreen struct {
	logw  int
	logh  int
	physw int
	physh int
	style Style
	evch  chan Event
	quit  chan struct{}

	front     []SimCell
	back      []Cell
	clear     bool
	cursorx   int
	cursory   int
	cursorvis bool
	mouse     bool
	charset   string
	encoder   transform.Transformer
	decoder   transform.Transformer
	fillchar  rune
	fillstyle Style

	sync.Mutex
}

func (s *simscreen) Init() error {
	s.evch = make(chan Event, 10)
	s.fillchar = 'X'
	s.fillstyle = StyleDefault
	s.mouse = false
	s.logw = 80
	s.logh = 25
	s.physw = 80
	s.physh = 25
	s.cursorx = -1
	s.cursory = -1
	s.style = StyleDefault

	switch s.charset {
	case "UTF-8", "US-ASCII":
		s.encoder = nil
		s.decoder = nil
	default:
		if enc := GetEncoding(s.charset); enc != nil {
			s.encoder = enc.NewEncoder()
			s.decoder = enc.NewDecoder()
		} else {
			return errors.New("no support for charset " + s.charset)
		}
	}

	s.front = make([]SimCell, s.physw*s.physh)
	s.back = ResizeCells(nil, 0, 0, s.logw, s.logh)

	return nil
}

func (s *simscreen) Fini() {
	if s.quit != nil {
		close(s.quit)
	}
	s.logw = 0
	s.logh = 0
	s.physw = 0
	s.physh = 0
	s.front = nil
	s.back = nil
}

func (s *simscreen) SetStyle(style Style) {
	s.Lock()
	s.style = style
	s.Unlock()
}

func (s *simscreen) Clear() {

	s.Lock()
	ClearCells(s.back, s.style)
	s.Unlock()
}

func (s *simscreen) SetCell(x, y int, style Style, ch ...rune) {

	s.Lock()
	if x < 0 || y < 0 || x >= s.logw || y >= s.logh {
		s.Unlock()
		return
	}
	cell := &s.back[(y*s.logw)+x]
	cell.SetCell(ch, style)
	s.Unlock()
}

func (s *simscreen) PutCell(x, y int, cell *Cell) {
	s.Lock()
	if x < 0 || y < 0 || x >= s.logw || y >= s.logh {
		s.Unlock()
		return
	}
	cp := &s.back[(y*s.logw)+x]
	cp.PutStyle(cell.Style)
	cp.PutChars(cell.Ch)
	s.Unlock()
}

func (s *simscreen) GetCell(x, y int) *Cell {
	s.Lock()
	if x < 0 || y < 0 || x >= s.logw || y >= s.logh {
		s.Unlock()
		return nil
	}
	cell := s.back[(y*s.logw)+x]
	s.Unlock()
	return &cell
}

func (s *simscreen) drawCell(x, y int, cell *Cell) {
	if x >= s.physw || y >= s.physh || x < 0 || y < 0 {
		return
	}
	simc := &s.front[(y*s.physw)+x]
	if cell.Style == StyleDefault {
		simc.Style = s.style
	} else {
		simc.Style = cell.Style
	}
	simc.Runes = nil
	simc.Runes = append(simc.Runes, cell.Ch...)

	// now emit runes - taking care to not overrun width with a
	// wide character, and to ensure that we emit exactly one regular
	// character followed up by any residual combing characters

	width := int(cell.Width)
	simc.Bytes = nil

	if len(cell.Ch) == 0 {
		simc.Bytes = []byte{' '}
		return
	}
	if width > 1 && x >= s.physw-1 {
		simc.Runes = []rune{' '}
		simc.Bytes = []byte{' '}
		return
	}

	enc := s.encoder
	ubuf := make([]byte, 12)
	lbuf := make([]byte, 12)

	for _, r := range simc.Runes {
		l := utf8.EncodeRune(ubuf, r)
		if enc == nil {
			simc.Bytes = append(simc.Bytes, ubuf[:l]...)
			if s.charset == "US-ASCII" {
				return
			}
			continue
		}
		nout, _, _ := enc.Transform(lbuf, ubuf[:l], true)
		if nout == 1 && lbuf[0] == '\x1a' {
			// replacement character
			if simc.Bytes == nil {
				simc.Bytes = append(simc.Bytes, '?')
			}
		} else if nout > 0 {
			simc.Bytes = append(simc.Bytes, lbuf[:nout]...)
		}
	}
}

func (s *simscreen) ShowCursor(x, y int) {
	s.Lock()
	s.cursorx, s.cursory = x, y
	s.showCursor()
	s.Unlock()
}

func (s *simscreen) HideCursor() {
	s.ShowCursor(-1, -1)
}

func (s *simscreen) showCursor() {

	x, y := s.cursorx, s.cursory
	if x < 0 || y < 0 || x >= s.physw || y >= s.physh {
		s.cursorvis = false
	} else {
		s.cursorvis = true
	}
}

func (s *simscreen) hideCursor() {
	// does not update cursor position
	s.cursorvis = false
}

func (s *simscreen) Show() {
	s.Lock()
	s.resize()
	s.draw()
	s.Unlock()
}

func (s *simscreen) clearScreen() {
	// We emulate a hardware clear by filling with a specific pattern
	for i := range s.front {
		s.front[i].Style = s.fillstyle
		s.front[i].Runes = []rune{s.fillchar}
		s.front[i].Bytes = []byte{byte(s.fillchar)}
	}
	s.clear = false
}

func (s *simscreen) draw() {
	// hide the cursor while we move stuff around
	s.hideCursor()

	if s.clear {
		s.clearScreen()
	}

	for row := 0; row < s.logh; row++ {
		for col := 0; col < s.logw; col++ {
			cell := &s.back[(row*s.logw)+col]
			if !cell.Dirty {
				continue
			}
			s.drawCell(col, row, cell)
			if cell.Width > 1 {
				col++
			}
			cell.Dirty = false
		}
	}

	// restore the cursor
	s.showCursor()
}

func (s *simscreen) EnableMouse() {
	s.mouse = true
}

func (s *simscreen) DisableMouse() {
	s.mouse = false
}

func (s *simscreen) Size() (int, int) {
	s.Lock()
	w, h := s.logw, s.logh
	s.Unlock()
	return w, h
}

func (s *simscreen) resize() {
	var ev Event
	w, h := s.physw, s.physh
	if w != s.logw || h != s.logh {
		ev = NewEventResize(w, h)
		s.back = ResizeCells(s.back, s.logw, s.logh, w, h)
		s.logw = w
		s.logh = h
	}
	if ev != nil {
		s.PostEvent(ev)
	}
}

func (s *simscreen) Colors() int {
	return 256
}

func (s *simscreen) PollEvent() Event {
	select {
	case <-s.quit:
		return nil
	case ev := <-s.evch:
		return ev
	}
}

func (s *simscreen) PostEvent(ev Event) {
	select {
	case s.evch <- ev:
	default:
		// drop the event on the floor
	}
}

func (s *simscreen) InjectMouse(x, y int, buttons ButtonMask, mod ModMask) {
	ev := NewEventMouse(x, y, buttons, mod)
	s.PostEvent(ev)
}

func (s *simscreen) InjectKey(key Key, r rune, mod ModMask) {
	ev := NewEventKey(KeyRune, r, ModNone)
	s.PostEvent(ev)
}

func (s *simscreen) InjectKeyBytes(b []byte) bool {
	failed := false

outer:
	for len(b) > 0 {
		if b[0] >= ' ' && b[0] <= 0x7F {
			// printable ASCII easy to deal with -- no encodings
			ev := NewEventKey(KeyRune, rune(b[0]), ModNone)
			s.PostEvent(ev)
			b = b[1:]
			continue
		}

		if b[0] < 0x80 {
			mod := ModNone
			// No encodings start with low numbered values
			if Key(b[0]) >= KeyCtrlA && Key(b[0]) <= KeyCtrlZ {
				mod = ModCtrl
			}
			ev := NewEventKey(Key(b[0]), 0, mod)
			s.PostEvent(ev)
			continue
		}

		switch s.charset {
		case "UTF-8":
			r, l := utf8.DecodeRune(b)
			if r == utf8.RuneError && (l == 0 || l == 1) {
				failed = true
				// yank off one byte
				b = b[1:]
			} else {
				b = b[l:]
				ev := NewEventKey(KeyRune, r, ModNone)
				s.PostEvent(ev)
				continue
			}

		case "US-ASCII":
			// ASCII cannot generate this, so most likely it was
			// entered as an Alt sequence
			ev := NewEventKey(KeyRune, rune(b[0]-128), ModAlt)
			s.PostEvent(ev)
			b = b[1:]
			continue

		default:
			utfb := make([]byte, len(b)*4) // worst case
			dec := s.decoder
			if dec == nil {
				failed = true
				b = b[1:]
				continue
			}

			// take care to consume at *most* a single rune
			for l := 1; l < len(b); l++ {
				dec.Reset()
				nout, nin, _ := dec.Transform(utfb, b[:l], true)

				if nout != 0 {
					r, _ := utf8.DecodeRune(utfb[:nout])
					ev := NewEventKey(KeyRune, r, ModNone)
					s.PostEvent(ev)
					b = b[nin:]
					continue outer
				}
			}
			failed = true
			b = b[1:]
			continue
		}
	}

	return failed == false
}

func (s *simscreen) Sync() {
	s.Lock()
	s.clear = true
	s.resize()
	InvalidateCells(s.back)
	s.draw()
	s.Unlock()
}

func (s *simscreen) CharacterSet() string {
	return s.charset
}

func (s *simscreen) Resize(w, h int) {
	s.Lock()
	newc := make([]SimCell, w*h)
	for row := 0; row < h && row < s.physh; row++ {
		for col := 0; col < w && col < s.physw; col++ {
			newc[(row*w)+col] = s.front[(row*s.physw)+col]
		}
	}
	s.physw = w
	s.physh = h
	s.Unlock()
}

func (s *simscreen) GetContents() ([]SimCell, int, int) {
	s.Lock()
	cells, w, h := s.front, s.physw, s.physh
	s.Unlock()
	return cells, w, h
}

func (s *simscreen) GetCursor() (int, int, bool) {
	s.Lock()
	x, y, vis := s.cursorx, s.cursory, s.cursorvis
	s.Unlock()
	return x, y, vis
}
