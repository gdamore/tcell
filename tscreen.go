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
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"sync"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

// NewTerminfoScreen returns a Screen that uses the stock TTY interface
// and POSIX termios, combined with a terminfo description taken from
// the $TERM environment variable.  It returns an error if the terminal
// is not supported for any reason.
//
// For terminals that do not support dynamic resize events, the $LINES
// $COLUMNS environment variables can be set to the actual window size,
// otherwise defaults taken from the terminal database are used.
func NewTerminfoScreen() (Screen, error) {
	ti, e := LookupTerminfo(os.Getenv("TERM"))
	if e != nil {
		return nil, e
	}
	t := &tScreen{ti: ti}

	t.keys = make(map[Key][]byte)
	if len(ti.Mouse) > 0 {
		t.mouse = []byte(ti.Mouse)
	}
	t.prepareKeys()
	t.buildAcsMap()
	t.sigwinch = make(chan os.Signal, 10)

	return t, nil
}

// tScreen represents a screen backed by a terminfo implementation.
type tScreen struct {
	ti       *Terminfo
	h        int
	w        int
	fini     bool
	cells    CellBuffer
	in       *os.File
	out      *os.File
	curstyle Style
	style    Style
	evch     chan Event
	sigwinch chan os.Signal
	quit     chan struct{}
	indoneq  chan struct{}
	keys     map[Key][]byte
	cx       int
	cy       int
	mouse    []byte
	clear    bool
	cursorx  int
	cursory  int
	tiosp    *termiosPrivate
	baud     int
	wasbtn   bool
	acs      map[rune]string
	charset  string
	encoder  transform.Transformer
	decoder  transform.Transformer

	sync.Mutex
}

func (t *tScreen) Init() error {
	t.evch = make(chan Event, 10)
	t.indoneq = make(chan struct{})
	t.charset = "UTF-8"

	t.charset = t.getCharset()
	switch t.charset {
	case "UTF-8", "US-ASCII":
		t.encoder = nil
		t.decoder = nil
	default:
		if enc := GetEncoding(t.charset); enc != nil {
			t.encoder = enc.NewEncoder()
			t.decoder = enc.NewDecoder()
		} else {
			return errors.New("no support for charset " + t.charset)
		}
	}
	ti := t.ti

	// environment overrides
	w := ti.Columns
	h := ti.Lines
	if i, _ := strconv.Atoi(os.Getenv("LINES")); i != 0 {
		h = i
	}
	if i, _ := strconv.Atoi(os.Getenv("COLUMNS")); i != 0 {
		w = i
	}
	if e := t.termioInit(); e != nil {
		return e
	}

	t.TPuts(ti.EnterCA)
	t.TPuts(ti.EnterKeypad)
	t.TPuts(ti.HideCursor)
	t.TPuts(ti.Clear)

	t.quit = make(chan struct{})

	t.Lock()
	t.cx = -1
	t.cy = -1
	t.style = StyleDefault
	t.cells.Resize(w, h)
	t.cursorx = -1
	t.cursory = -1
	t.resize()
	t.Unlock()

	go t.inputLoop()

	return nil
}

func (t *tScreen) prepareKey(key Key, val string) {
	if val != "" {
		t.keys[key] = []byte(val)
	}
}

func (t *tScreen) prepareKeys() {
	ti := t.ti
	t.prepareKey(KeyBackspace, ti.KeyBackspace)
	t.prepareKey(KeyF1, ti.KeyF1)
	t.prepareKey(KeyF2, ti.KeyF2)
	t.prepareKey(KeyF3, ti.KeyF3)
	t.prepareKey(KeyF4, ti.KeyF4)
	t.prepareKey(KeyF5, ti.KeyF5)
	t.prepareKey(KeyF6, ti.KeyF6)
	t.prepareKey(KeyF7, ti.KeyF7)
	t.prepareKey(KeyF8, ti.KeyF8)
	t.prepareKey(KeyF9, ti.KeyF9)
	t.prepareKey(KeyF10, ti.KeyF10)
	t.prepareKey(KeyF11, ti.KeyF11)
	t.prepareKey(KeyF12, ti.KeyF12)
	t.prepareKey(KeyF13, ti.KeyF13)
	t.prepareKey(KeyF14, ti.KeyF14)
	t.prepareKey(KeyF15, ti.KeyF15)
	t.prepareKey(KeyF16, ti.KeyF16)
	t.prepareKey(KeyF17, ti.KeyF17)
	t.prepareKey(KeyF18, ti.KeyF18)
	t.prepareKey(KeyF19, ti.KeyF19)
	t.prepareKey(KeyF20, ti.KeyF20)
	t.prepareKey(KeyF21, ti.KeyF21)
	t.prepareKey(KeyF22, ti.KeyF22)
	t.prepareKey(KeyF23, ti.KeyF23)
	t.prepareKey(KeyF24, ti.KeyF24)
	t.prepareKey(KeyF25, ti.KeyF25)
	t.prepareKey(KeyF26, ti.KeyF26)
	t.prepareKey(KeyF27, ti.KeyF27)
	t.prepareKey(KeyF28, ti.KeyF28)
	t.prepareKey(KeyF29, ti.KeyF29)
	t.prepareKey(KeyF30, ti.KeyF30)
	t.prepareKey(KeyF31, ti.KeyF31)
	t.prepareKey(KeyF32, ti.KeyF32)
	t.prepareKey(KeyF33, ti.KeyF33)
	t.prepareKey(KeyF34, ti.KeyF34)
	t.prepareKey(KeyF35, ti.KeyF35)
	t.prepareKey(KeyF36, ti.KeyF36)
	t.prepareKey(KeyF37, ti.KeyF37)
	t.prepareKey(KeyF38, ti.KeyF38)
	t.prepareKey(KeyF39, ti.KeyF39)
	t.prepareKey(KeyF40, ti.KeyF40)
	t.prepareKey(KeyF41, ti.KeyF41)
	t.prepareKey(KeyF42, ti.KeyF42)
	t.prepareKey(KeyF43, ti.KeyF43)
	t.prepareKey(KeyF44, ti.KeyF44)
	t.prepareKey(KeyF45, ti.KeyF45)
	t.prepareKey(KeyF46, ti.KeyF46)
	t.prepareKey(KeyF47, ti.KeyF47)
	t.prepareKey(KeyF48, ti.KeyF48)
	t.prepareKey(KeyF49, ti.KeyF49)
	t.prepareKey(KeyF50, ti.KeyF50)
	t.prepareKey(KeyF51, ti.KeyF51)
	t.prepareKey(KeyF52, ti.KeyF52)
	t.prepareKey(KeyF53, ti.KeyF53)
	t.prepareKey(KeyF54, ti.KeyF54)
	t.prepareKey(KeyF55, ti.KeyF55)
	t.prepareKey(KeyF56, ti.KeyF56)
	t.prepareKey(KeyF57, ti.KeyF57)
	t.prepareKey(KeyF58, ti.KeyF58)
	t.prepareKey(KeyF59, ti.KeyF59)
	t.prepareKey(KeyF60, ti.KeyF60)
	t.prepareKey(KeyF61, ti.KeyF61)
	t.prepareKey(KeyF62, ti.KeyF62)
	t.prepareKey(KeyF63, ti.KeyF63)
	t.prepareKey(KeyF64, ti.KeyF64)
	t.prepareKey(KeyInsert, ti.KeyInsert)
	t.prepareKey(KeyDelete, ti.KeyDelete)
	t.prepareKey(KeyHome, ti.KeyHome)
	t.prepareKey(KeyEnd, ti.KeyEnd)
	t.prepareKey(KeyUp, ti.KeyUp)
	t.prepareKey(KeyDown, ti.KeyDown)
	t.prepareKey(KeyLeft, ti.KeyLeft)
	t.prepareKey(KeyRight, ti.KeyRight)
	t.prepareKey(KeyPgUp, ti.KeyPgUp)
	t.prepareKey(KeyPgDn, ti.KeyPgDn)
	t.prepareKey(KeyHelp, ti.KeyHelp)
	t.prepareKey(KeyPrint, ti.KeyPrint)
	t.prepareKey(KeyCancel, ti.KeyCancel)
	t.prepareKey(KeyExit, ti.KeyExit)
	t.prepareKey(KeyBacktab, ti.KeyBacktab)
}

func (t *tScreen) Fini() {
	ti := t.ti
	t.Lock()
	t.cells.Resize(0, 0)
	t.TPuts(ti.ShowCursor)
	t.TPuts(ti.AttrOff)
	t.TPuts(ti.Clear)
	t.TPuts(ti.ExitCA)
	t.TPuts(ti.ExitKeypad)
	t.TPuts(ti.TParm(ti.MouseMode, 0))
	t.curstyle = Style(-1)
	t.clear = false
	t.fini = true
	t.Unlock()

	if t.quit != nil {
		close(t.quit)
	}
	t.termioFini()
}

func (t *tScreen) SetStyle(style Style) {
	t.Lock()
	if !t.fini {
		t.style = style
	}
	t.Unlock()
}

func (t *tScreen) Clear() {
	t.Lock()
	if !t.fini {
		t.cells.Fill(' ', t.style)
	}
	t.Unlock()
}

func (t *tScreen) SetContent(x, y int, mainc rune, combc []rune, style Style) {
	t.Lock()
	if !t.fini {
		t.cells.SetContent(x, y, mainc, combc, style)
	}
	t.Unlock()
}

func (t *tScreen) GetContent(x, y int) (rune, []rune, Style, int) {
	t.Lock()
	mainc, combc, style, width := t.cells.GetContent(x, y)
	t.Unlock()
	return mainc, combc, style, width
}

func (t *tScreen) SetCell(x, y int, style Style, ch ...rune) {
	if len(ch) > 0 {
		t.SetContent(x, y, ch[0], ch[1:], style)
	} else {
		t.SetContent(x, y, ' ', nil, style)
	}
}

func (t *tScreen) encodeRune(r rune, buf []byte) []byte {

	// all the character sets we care about are ASCII supersets
	if r < 0x80 {
		buf = append(buf, byte(r))
		return buf
	}

	enc := t.encoder
	if enc == nil {
		// This is probably ASCII.  Only append a filler character
		// for the main rune, ignoring combining runes.
		if len(buf) == 0 {
			buf = append(buf, '?')
		}
		return buf
	}

	nb := make([]byte, 6)
	ob := make([]byte, 6)
	num := utf8.EncodeRune(ob, r)
	ob = ob[:num]
	t.encoder.Reset()
	dst, _, err := t.encoder.Transform(nb, ob, true)
	if err == nil {
		nb = nb[:dst]
		if len(nb) == 0 || nb[0] == '\x1a' {
			// special handling for replacement chars,
			// we can't display them, so use "?" for primary char,
			// and elide combining
			if len(buf) == 0 {
				if acs, ok := t.acs[r]; ok {
					buf = append(buf, []byte(acs)...)
				} else {
					buf = append(buf, '?')
				}
			}
		} else {
			buf = append(buf, nb...)
		}
	} else if len(buf) == 0 {
		buf = append(buf, '?')
	}
	return buf
}

func (t *tScreen) drawCell(x, y int) int {

	ti := t.ti

	mainc, combc, style, width := t.cells.GetContent(x, y)

	if !t.cells.Dirty(x, y) {
		return width
	}

	if t.cy != y || t.cx != x {
		t.TPuts(ti.TGoto(x, y))
		t.cx = x
		t.cy = y
	}

	if style == StyleDefault {
		style = t.style
	}
	if style != t.curstyle {
		fg, bg, attrs := style.Decompose()

		t.TPuts(ti.AttrOff)
		if attrs&AttrBold != 0 {
			t.TPuts(ti.Bold)
		}
		if attrs&AttrUnderline != 0 {
			t.TPuts(ti.Underline)
		}
		if attrs&AttrReverse != 0 {
			t.TPuts(ti.Reverse)
		}
		if attrs&AttrBlink != 0 {
			t.TPuts(ti.Blink)
		}
		if attrs&AttrDim != 0 {
			t.TPuts(ti.Dim)
		}
		if fg != ColorDefault {
			c := int(fg) - 1
			t.TPuts(ti.TParm(ti.SetFg, c))
		}
		if bg != ColorDefault {
			c := int(bg) - 1
			t.TPuts(ti.TParm(ti.SetBg, c))
		}
		t.curstyle = style
	}
	// now emit runes - taking care to not overrun width with a
	// wide character, and to ensure that we emit exactly one regular
	// character followed up by any residual combing characters

	if width < 1 {
		width = 1
	}

	var str string

	switch t.charset {
	case "UTF-8":
		str = string(mainc)
		if combc != nil {
			str += string(combc)
		}
	default:
		// Non-Unicode systems.  Make do.
		buf := make([]byte, 0, 6)

		buf = t.encodeRune(mainc, buf)
		for _, r := range combc {
			buf = t.encodeRune(r, buf)
		}

		str = string(buf)
		if width > 1 && str == "?" {
			// No FullWidth character support
			str = "? "
		}
	}

	// XXX: check for hazeltine not being able to display ~

	if x > t.w-width {
		// too wide to fit; emit a single space instead
		width = 1
		str = " "
	}
	io.WriteString(t.out, str)
	t.cx += width
	t.cells.SetDirty(x, y, false)

	return width
}

func (t *tScreen) ShowCursor(x, y int) {
	t.Lock()
	t.cursorx = x
	t.cursory = y
	t.Unlock()
}

func (t *tScreen) HideCursor() {
	t.ShowCursor(-1, -1)
}

func (t *tScreen) showCursor() {

	x, y := t.cursorx, t.cursory
	w, h := t.cells.Size()
	if x < 0 || y < 0 || x >= w || y >= h {
		t.hideCursor()
		return
	}
	if t.cx != x || t.cy != y {
		t.TPuts(t.ti.TGoto(x, y))
	}
	t.TPuts(t.ti.ShowCursor)
	t.cx = x
	t.cy = y
}

func (t *tScreen) TPuts(s string) {
	t.ti.TPuts(t.out, s, t.baud)
}

func (t *tScreen) Show() {
	t.Lock()
	if !t.fini {
		t.resize()
		t.draw()
	}
	t.Unlock()
}

func (t *tScreen) clearScreen() {
	t.TPuts(t.ti.Clear)
	t.clear = false
}

func (t *tScreen) hideCursor() {
	// does not update cursor position
	if t.ti.HideCursor != "" {
		t.TPuts(t.ti.HideCursor)
	} else {
		// No way to hide cursor, stick it
		// at bottom right of screen
		t.cx, t.cy = t.cells.Size()
		t.TPuts(t.ti.TGoto(t.cx, t.cy))
	}
}

func (t *tScreen) draw() {
	// clobber cursor position, because we're gonna change it all
	t.cx = -1
	t.cy = -1

	// hide the cursor while we move stuff around
	t.hideCursor()

	if t.clear {
		t.clearScreen()
	}

	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			width := t.drawCell(x, y)
			x += width - 1
		}
	}

	// restore the cursor
	t.showCursor()
}

func (t *tScreen) EnableMouse() {
	if len(t.mouse) != 0 {
		t.TPuts(t.ti.TParm(t.ti.MouseMode, 1))
	}
}

func (t *tScreen) DisableMouse() {
	if len(t.mouse) != 0 {
		t.TPuts(t.ti.TParm(t.ti.MouseMode, 0))
	}
}

func (t *tScreen) Size() (int, int) {
	t.Lock()
	w, h := t.w, t.h
	t.Unlock()
	return w, h
}

func (t *tScreen) resize() {
	if w, h, e := t.getWinSize(); e == nil {
		if w != t.w || h != t.h {
			t.cx = -1
			t.cy = -1

			t.cells.Resize(w, h)
			t.cells.Invalidate()
			t.h = h
			t.w = w
			ev := NewEventResize(w, h)
			t.PostEvent(ev)
		}
	}
}

func (t *tScreen) Colors() int {
	// this doesn't change, no need for lock
	return t.ti.Colors
}

func (t *tScreen) PollEvent() Event {
	select {
	case <-t.quit:
		return nil
	case ev := <-t.evch:
		return ev
	}
}

// bulidAcsMap builds a map of characters that we translate from Unicode to
// alternate character encodings.  To do this, we use the standard VT100 ACS
// maps.  This is only done if the terminal lacks support for Unicode; we
// always prefer to emit Unicode glyphs when we are able.

type acsMap struct {
	utf   rune   // UTF-8 glyph
	vt100 rune   // VT100 name
	ascii string // ASCII default
}

func (t *tScreen) buildAcsMap() {
	acsstr := t.ti.AltChars
	vtNames := []acsMap{
		{RuneSterling, '}', "f"},
		{RuneDArrow, '.', "v"},
		{RuneLArrow, ',', "<"},
		{RuneRArrow, '+', ">"},
		{RuneUArrow, '-', "^"},
		{RuneBullet, '~', "o"},
		{RuneBoard, 'h', "#"},
		{RuneCkBoard, 'a', ":"},
		{RuneDegree, 'f', "\\"},
		{RuneDiamond, '`', "+"},
		{RuneGEqual, 'z', ">"},
		{RunePi, '{', "*"},
		{RuneHLine, 'q', "-"},
		{RuneLantern, 'i', "#"},
		{RunePlus, 'n', "+"},
		{RuneLEqual, 'y', "<"},
		{RuneLLCorner, 'm', "+"},
		{RuneLRCorner, 'j', "+"},
		{RuneNEqual, '|', "!"},
		{RunePlMinus, 'g', "#"},
		{RuneS1, 'o', "~"},
		{RuneS3, 'p', "-"},
		{RuneS7, 'r', "-"},
		{RuneS9, 's', "_"},
		{RuneBlock, '0', "#"},
		{RuneTTee, 'w', "+"},
		{RuneRTee, 'u', "+"},
		{RuneLTee, 't', "+"},
		{RuneBTee, 'v', "+"},
		{RuneULCorner, 'l', "+"},
		{RuneURCorner, 'k', "+"},
		{RuneVLine, 'x', "|"},
	}
	t.acs = make(map[rune]string)
	for i := range vtNames {
		// prefill defaults
		t.acs[vtNames[i].utf] = vtNames[i].ascii
	}
	for len(acsstr) > 2 {
		srcv := rune(acsstr[0])
		dstv := string(acsstr[1])
		// O(n*2), but n is pretty small (30)
		for i := range vtNames {
			if srcv == vtNames[i].vt100 {
				t.acs[vtNames[i].utf] =
					t.ti.EnterAcs + dstv + t.ti.ExitAcs
				break
			}
		}
		acsstr = acsstr[2:]
	}
}

func (t *tScreen) PostEvent(ev Event) {
	select {
	case t.evch <- ev:
	default:
		// drop the event on the floor
	}
}

func (t *tScreen) postMouseEvent(x, y, btn int) {

	// XTerm mouse events only report at most one button at a time,
	// which may include a wheel button.  Wheel motion events are
	// reported as single impulses, while other button events are reported
	// as separate press & release events.

	button := ButtonNone
	mod := ModNone

	// Mouse wheel has bit 6 set, no release events.  It should be noted
	// that wheel events are sometimes misdelivered as mouse button events
	// during a click-drag, so we debounce these, considering them to be
	// button press events unless we see an intervening release event.
	switch btn & 0x43 {
	case 0:
		button = Button1
		t.wasbtn = true
	case 1:
		button = Button2
		t.wasbtn = true
	case 2:
		button = Button3
		t.wasbtn = true
	case 3:
		button = ButtonNone
		t.wasbtn = false
	case 0x40:
		if !t.wasbtn {
			button = WheelUp
		} else {
			button = Button1
		}
	case 0x41:
		if !t.wasbtn {
			button = WheelDown
		} else {
			button = Button2
		}
	}

	if btn&0x4 != 0 {
		mod |= ModShift
	}
	if btn&0x8 != 0 {
		mod |= ModMeta
	}
	if btn&0x10 != 0 {
		mod |= ModCtrl
	}

	// Some terminals will report mouse coordinates outside the
	// screen, especially with click-drag events.  Clip the coordinates
	// to the screen in that case.
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	w, h := t.cells.Size()
	if x > w-1 {
		x = w - 1
	}
	if y > h-1 {
		y = h - 1
	}
	ev := NewEventMouse(x, y, button, mod)
	t.PostEvent(ev)
}

// parseSgrMouse attempts to locate an SGR mouse record at the start of the
// buffer.  It returns true, true if it found one, and the associated bytes
// be removed from the buffer.  It returns true, false if the buffer might
// contain such an event, but more bytes are necessary (partial match), and
// false, false if the content is definitely *not* an SGR mouse record.
func (t *tScreen) parseSgrMouse(buf *bytes.Buffer) (bool, bool) {

	b := buf.Bytes()

	var x, y, btn, state int
	dig := false
	neg := false
	i := 0
	val := 0

	for i = range b {
		switch b[i] {
		case '\x1b':
			if state != 0 {
				return false, false
			}
			state = 1

		case '\x9b':
			if state != 0 {
				return false, false
			}
			state = 2

		case '[':
			if state != 1 {
				return false, false
			}
			state = 2

		case '<':
			if state != 2 {
				return false, false
			}
			val = 0
			dig = false
			neg = false
			state = 3

		case '-':
			if state != 3 || state != 4 || state != 5 {
				return false, false
			}
			if dig || neg {
				return false, false
			}
			neg = true // stay in state

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if state != 3 && state != 4 && state != 5 {
				return false, false
			}
			val *= 10
			val += int(b[i] - '0')
			dig = true // stay in state

		case ';':
			if neg {
				val = -val
			}
			switch state {
			case 3:
				btn, val = val, 0
				neg, dig, state = false, false, 4
			case 4:
				x, val = val, 0
				neg, dig, state = false, false, 5
			default:
				return false, false
			}

		case 'm', 'M':
			if state != 5 {
				return false, false
			}
			if neg {
				val = -val
			}
			y = val

			// We don't care about the motion bit
			btn &^= 32
			if b[i] == 'm' {
				// mouse release, clear all buttons
				btn |= 3
				btn &^= 0x40
			}
			// consume the event bytes
			for i >= 0 {
				buf.ReadByte()
				i--
			}
			t.postMouseEvent(x, y, btn)
			return true, true
		}
	}

	// incomplete & inconclusve at this point
	return true, false
}

// parseXtermMouse is like parseSgrMouse, but it parses a legacy
// X11 mouse record.
func (t *tScreen) parseXtermMouse(buf *bytes.Buffer) (bool, bool) {

	b := buf.Bytes()

	state := 0
	btn := 0
	x := 0
	y := 0

	for i := range b {
		switch state {
		case 0:
			switch b[i] {
			case '\x1b':
				state = 1
			case '\x9b':
				state = 2
			default:
				return false, false
			}
		case 1:
			if b[i] != '[' {
				return false, false
			}
			state = 2
		case 2:
			if b[i] != 'M' {
				return false, false
			}
			state++
		case 3:
			btn = int(b[i])
			state++
		case 4:
			x = int(b[i]) - 32 - 1
			state++
		case 5:
			y = int(b[i]) - 32 - 1
			for i >= 0 {
				buf.ReadByte()
				i--
			}
			t.postMouseEvent(x, y, btn)
			return true, true
		}
	}
	return true, false
}

func (t *tScreen) parseFunctionKey(buf *bytes.Buffer) (bool, bool) {
	b := buf.Bytes()
	partial := false
	for k, esc := range t.keys {
		if bytes.HasPrefix(b, esc) {
			// matched
			var r rune
			if len(esc) == 1 {
				r = rune(b[0])
			}
			ev := NewEventKey(k, r, ModNone)
			t.PostEvent(ev)
			for i := 0; i < len(esc); i++ {
				buf.ReadByte()
			}
			return true, true
		}
		if bytes.HasPrefix(esc, b) {
			partial = true
		}
	}
	return partial, false
}

func (t *tScreen) parseRune(buf *bytes.Buffer) (bool, bool) {
	b := buf.Bytes()
	if b[0] >= ' ' && b[0] <= 0x7F {
		// printable ASCII easy to deal with -- no encodings
		ev := NewEventKey(KeyRune, rune(b[0]), ModNone)
		t.PostEvent(ev)
		buf.ReadByte()
		return true, true
	}

	if b[0] < 0x80 {
		// No encodings start with low numbered values
		return false, false
	}

	switch t.charset {
	case "UTF-8":
		if utf8.FullRune(b) {
			r, _, e := buf.ReadRune()
			if e == nil {
				ev := NewEventKey(KeyRune, r, ModNone)
				t.PostEvent(ev)
				return true, true
			}
		}
	case "US-ASCII":
		// ASCII cannot generate this, so most likely it was
		// entered as an Alt sequence
		ev := NewEventKey(KeyRune, rune(b[0]-128), ModAlt)
		t.PostEvent(ev)
		buf.ReadByte()
		return true, true

	default:
		utfb := make([]byte, 12)
		for l := 1; l <= len(b); l++ {
			t.decoder.Reset()
			nout, nin, _ := t.decoder.Transform(utfb, b[:l], true)
			if nout != 0 {
				if r, _ := utf8.DecodeRune(utfb[:nout]); r != utf8.RuneError {
					ev := NewEventKey(KeyRune, r, ModNone)
					t.PostEvent(ev)
				}
				for eat := 0; eat < nin; eat++ {
					buf.ReadByte()
				}
				return true, true
			}
		}
	}
	// Looks like potential escape
	return true, false
}

func (t *tScreen) scanInput(buf *bytes.Buffer, expire bool) {

	for {
		b := buf.Bytes()
		if len(b) == 0 {
			buf.Reset()
			return
		}

		partials := 0

		if part, comp := t.parseRune(buf); comp {
			continue
		} else if part {
			partials++
		}

		if part, comp := t.parseFunctionKey(buf); comp {
			continue
		} else if part {
			partials++
		}

		// Only parse mouse records if this term claims to have
		// mouse support

		if t.ti.Mouse != "" {
			if part, comp := t.parseXtermMouse(buf); comp {
				continue
			} else if part {
				partials++
			}

			if part, comp := t.parseSgrMouse(buf); comp {
				continue
			} else if part {
				partials++
			}
		}

		if partials == 0 || expire {
			// Nothing was going to match, or we timed out
			// waiting for more data -- just deliver the characters
			// to the app & let them sort it out.  Possibly we should only
			// do this for control characters such like ESC.
			by, _ := buf.ReadByte()
			ev := NewEventKey(KeyRune, rune(by), ModNone)
			t.PostEvent(ev)
			continue
		}

		// well we have some partial data, wait until we get
		// some more
		break
	}
}

func (t *tScreen) inputLoop() {
	buf := &bytes.Buffer{}

	chunk := make([]byte, 128)
	for {
		select {
		case <-t.quit:
			close(t.indoneq)
			return
		case <-t.sigwinch:
			t.Lock()
			t.cx = -1
			t.cy = -1
			t.resize()
			t.cells.Invalidate()
			t.draw()
			t.Unlock()
			continue
		default:
		}
		n, e := t.in.Read(chunk)
		switch e {
		case io.EOF:
			// If we timeout waiting for more bytes, then it's
			// time to give up on it.  Even at 300 baud it takes
			// less than 0.5 ms to transmit a whole byte.
			if buf.Len() > 0 {
				t.scanInput(buf, true)
			}
			continue
		case nil:
		default:
			close(t.indoneq)
			return
		}
		buf.Write(chunk[:n])
		// Now we need to parse the input buffer for events
		t.scanInput(buf, false)
	}
}

func (t *tScreen) Sync() {
	t.Lock()
	t.cx = -1
	t.cy = -1
	if !t.fini {
		t.resize()
		t.clear = true
		t.cells.Invalidate()
		t.draw()
	}
	t.Unlock()
}

func (t *tScreen) CharacterSet() string {
	return t.charset
}
