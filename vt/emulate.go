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
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v3/color"
	"github.com/rivo/uniseg"
)

// Emulator is a terminal emulator API. It implements the state machinery
// (escape parsing and so forth) associated with being a terminal emulator.
// The backend handles rendering the content, and some low level details.
//
// NOTE: This is not a committed interface yet, its entirely a work in progress.
type Emulator interface {
	// SetId sets our identity.
	SetId(name string, version string)

	// SendRaw sends raw data to the consumer.  This bypasses the normal encoding,
	// so it should be used with caution.
	SendRaw([]byte)

	// KeyEvent injects a keyboard event into the emulator, which will ultimately
	// result in data being sent via SendRaw.
	KeyEvent(ev KbdEvent)

	// ResizeEvent is called by a backend when the terminal has resized
	// This will send in-band resize notifications if the client has requested them.
	ResizeEvent()

	// Drain waits until any queued but not processed input has finished processing.
	// It also wakes the reader.
	Drain() error

	// Start starts processing.
	Start() error

	// Stop stops processing.
	Stop() error

	// Reader reads data from the emulator.  These are bytes that would be transmitted
	// to a remote party.
	io.Reader

	// Writer writes data to the emulator.  These are commands that the emulator should process.
	io.Writer
}

// Style represents the styling of a cell.
// This is an interface to prevent direct modification.
type Style interface {
	Fg() color.Color              // Fg returns the foreground color.
	Bg() color.Color              // Bg returns the background color.
	Uc() color.Color              // Uc returns the underline color.k
	Attr() Attr                   // Attr returns the associated attributes.
	Url() (string, string)        // Url returns the URL and associated id if one was set.
	WithFg(color.Color) Style     // WithFg creates a new style with the foreground
	WithBg(color.Color) Style     // WithBg creates a new style with the background.
	WithUc(color.Color) Style     // WithUc creates a new style with the underline color
	WithAttr(Attr) Style          // WithAttr creates a new style with the attributes.
	WithUrl(string, string) Style // WithLink creates a new style with the URL and id.
	Equal(Style) bool             // Equal returns true if the styles are the same.
}

// styleStruct implements Style.  Note that it is possible to make this even more
// compact, but we don't think further optimization here on size will justify the
// complexity and runtime performance hit to do so.  We're also already only storing
// a class reference to this per cell.
type styleStruct struct {
	fg   color.Color
	bg   color.Color
	uc   color.Color // underline color
	attr Attr
	url  string // URL
	id   string // Id for link
}

var BaseStyle = &styleStruct{}

func (ss *styleStruct) Fg() color.Color              { return ss.fg }
func (ss *styleStruct) Bg() color.Color              { return ss.bg }
func (ss *styleStruct) Uc() color.Color              { return ss.uc }
func (ss *styleStruct) Attr() Attr                   { return ss.attr }
func (ss *styleStruct) Url() (string, string)        { return ss.url, ss.id }
func (ss *styleStruct) WithFg(fg color.Color) Style  { ns := *ss; ns.fg = fg; return &ns }
func (ss *styleStruct) WithBg(bg color.Color) Style  { ns := *ss; ns.bg = bg; return &ns }
func (ss *styleStruct) WithUc(uc color.Color) Style  { ns := *ss; ns.uc = uc; return &ns }
func (ss *styleStruct) WithAttr(a Attr) Style        { ns := *ss; ns.attr = a; return &ns }
func (ss *styleStruct) WithUrl(url, id string) Style { ns := *ss; ns.url = url; ns.id = id; return &ns }
func (ss *styleStruct) Equal(other Style) bool {
	if s2, ok := other.(*styleStruct); ok {
		return *ss == *s2
	}
	// We have chosen not to support alternative implementations for this compare.
	// We could delegate to the other style, but that could lead to a loop if they
	// do the same.
	return (false)
}

// Cell is a representation of a display cell. Most consumers will not need this.
// Storing the width out of band saves 7 bytes per cell.
type Cell struct {
	C string // Content, it will be a grapheme cluster
	S Style  // Style, a pointer is used efficiency
}

// NewEmulator creates an emulator instance on top of the given backend.
// The input is relative to the emulator, so it receives data from the host,
// whereas the emulator sends data to the application through the output.
func NewEmulator(be Backend) Emulator {
	stopQ := make(chan bool)
	em := &emulator{
		be:     be,
		inBuf:  &bytes.Buffer{},
		writeQ: make(chan any),
		readQ:  make(chan any, 1024),
		stopQ:  stopQ,
		style:  BaseStyle,
		localModes: map[PrivateMode]ModeStatus{
			PmAppCursor:  ModeOff,
			PmAutoMargin: ModeOn,
		},
	}
	if _, ok := be.(Resizer); ok {
		em.localModes[PmResizeReports] = ModeOff
	}
	em.cells = make([]Cell, int(be.GetSize().X)*int(be.GetSize().Y))
	close(stopQ)
	em.inb = em.inbInit
	return em
}

// emulator is an implementation of a terminal emulator built on top of
// a Backend.  It implements the common escape sequence handling and high
// level functionality that a real terminal emulator, or a mock, would need.
type emulator struct {
	stopQ     chan bool
	writeQ    chan any // queues data from application to emulator
	readQ     chan any // queues data from emulator to application
	be        Backend
	inBuf     *bytes.Buffer // buffer queued for input
	inb       func(byte)    // input byte function (faster than state switch)
	style     Style
	utfLen    int
	pos       Coord
	autoWrap  bool        // next character will wrap (auto margin, deferred until char emitted)
	sevenOnly bool        // only allow 7-bit escapes (needed for KOI8, ShiftJIS, etc.)
	name      string      // name of this emulator (used for extended attributes)
	vers      string      // version string of this emulator (used for extended attributes)
	savedPos  Coord       // saved via DECSC
	saved     savedCursor // data saved by save cursor (DECSC)
	sendLock  sync.Mutex  // ensures that send data cannot be intermixed
	tabStops  []Col       // tab stops, ordered. if nil every 8th position is used
	lastIndex int         // index of last cell written + 1 (for grapheme clustering) (zero means none)
	cells     []Cell      // content of cells, we have to maintain our own copy (backend might or might not)

	localModes map[PrivateMode]ModeStatus // some modes we handle locally
}

// savedCursor is the content we save when saving the cursor,
// which is more than just the cursor location itself.
type savedCursor struct {
	pos      Coord
	style    Style
	autoWrap bool
	// We should probably store OSC 8 data here, eventually.
	// TODO: Character sets
	// TODO: Origin mode (DEC Mode 6)
}

func (em *emulator) saveCursor() {
	em.saved.pos = em.getPosition()
	em.saved.style = em.style
	em.saved.autoWrap = em.autoWrap
}

func (em *emulator) restoreCursor() {
	em.setPosition(em.saved.pos)
	em.autoWrap = em.saved.autoWrap
	em.style = em.saved.style
	em.be.SetAttr(em.style.Attr())
	if c, ok := em.be.(Colorer); ok {
		c.SetFgColor(em.style.Fg())
		c.SetBgColor(em.style.Bg())
	}
}

// inbInit processes bytes received in the "default" state. Most often these are just
// text characters to display on screen, but if ESC is seen then additional processing will result.
func (em *emulator) inbInit(b byte) {
	em.inBuf.Reset()

	// hot path - just doing ASCII directly.
	if b >= ' ' && b < 0x7f {
		// plain ascii
		em.putRune(rune(b))
		return
	}

	// For 8-bit encodings, we treat these as Fe sequences.
	// Basically the same as ESC followed by (b - 0x40).
	// TODO: condition this so that we do not do this if
	// the encoding cannot support it (UTF, 8859, and EUC encodings
	// are all fine here, but others like ShiftJIS or KOI8 might not be).
	if b >= 0x80 && b <= 0x9F && !em.sevenOnly {
		em.lastIndex = 0
		em.inbEsc(b - 0x40)
		return
	}

	// TODO: To support non-UTF-8 locales, include a check here for > 0x7F.  Those locales
	// might preclude 8-bit control sequences - 8859 character sets are fine, but e.g. KOI8,
	// and ShiftJIS use values in those ranges.

	switch b {
	case 0x1b: // ESC (escape)
		em.inb = em.inbEsc
	case 0x07: // BEL (bell)
		em.beep()
	case 0x08: // BS (backspace)
		em.lastIndex = 0
		em.moveLeft()
	case 0x09: // horizontal tab
		em.lastIndex = 0
		em.nextTab()
	case 0x0a: // NL (newline)
		em.lastIndex = 0
		em.nextLine()
	case 0x0b: // VT (vertical tab, treat as LF)
		em.lastIndex = 0
		em.nextLine()
	case 0x0c: // FF (form feed, treat as LF)
		em.lastIndex = 0
		em.nextLine()
	case 0x0d: // CR (carriage return)
		em.lastIndex = 0
		em.setPosition(Coord{0, em.getPosition().Y})
	case 0x0e: // TODO: SO
		em.lastIndex = 0
	case 0x0f: // TODO: SI
		em.lastIndex = 0
	case 0x18: //TODO Cancel (reset parser)
		em.lastIndex = 0
	default:
		// TODO: consider separating Unicode from other 8-bit character sets
		if b&0xE0 == 0xC0 {
			em.utfLen = 2
			em.inb = em.inbUTF
			em.inBuf.WriteByte(b)
		} else if b&0xF0 == 0xE0 {
			em.utfLen = 3
			em.inb = em.inbUTF
			em.inBuf.WriteByte(b)
		} else if b&0xF8 == 0xF0 {
			em.utfLen = 4
			em.inb = em.inbUTF
			em.inBuf.WriteByte(b)
		} else {
			em.lastIndex = 0
			em.beep()
		}
	}
}

// inbEsc processes the next byte after an escape character is seen.
func (em *emulator) inbEsc(b byte) {
	// By default, reset to init state. Other states will be set explicitly as needed.
	em.inb = em.inbInit
	em.lastIndex = 0

	switch b {
	case '[':
		em.inb = em.inbCSI
	case ']':
		em.inb = em.inbOSC
	case ' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/':
		// 0x20 - 0x2F -- usually followed by just one terminating character, but could include others
		em.inb = em.inbNF
		em.inBuf.WriteByte(b)
	case '^': // privacy message (PM)
		em.inb = em.inbStr
	case '_': // application program command (APC)
		em.inb = em.inbStr
	case 'D': // down one line (IND) - note does not reset auto wrap
		// TODO: should scroll if needed
		em.moveDown()
	case 'E': // next line (NEL)
		em.nextLine()
	case 'H': // set tab stop (HTS) - VT52 is go home, but we do not support VT52
		em.setTabStop(em.getPosition().X)
	case 'M': // up one line (RI) - note does not reset autoWrap
		em.moveUp()
	case 'N': // single shift two (SS2) (TODO)
	case 'O': // single shift three (SS3) (TODO)
	case 'P':
		em.inb = em.inbStr // device control string (DCS) (TODO)
	case 'X': // start of string (SOS)
		em.inb = em.inbStr
	case 'Z': // DECID, obsolete form to get primary DA
		em.sendDA()
	case 'c': // RIS, soft reset
		em.softReset()
	case '6': // back index (DECBI, VT420, not widely supported)
		em.moveLeft()
	case '7': // save cursor (DECSC, VT100)
		em.saveCursor()
	case '8': // restore cursor (DECRC, VT100)
		em.restoreCursor()
	case '9': // forward index (DECFI, VT420, not widely supported)
		em.moveRight()
	default:
		// ESC-V and ESC-W are for guarded area (TODO)
		em.inb = em.inbInit
	}
}

// inbNF processes bytes that are part of an "nF" sequence (see ECMA-48).
func (em *emulator) inbNF(b byte) {
	if b >= 0x20 && b <= 0x2F {
		em.inBuf.WriteByte(b)
		return
	}
	if b < 0x20 || b > 0x7E { // not a valid sequence
		em.beep()
		em.inb = em.inbInit
		return
	}
	em.inBuf.WriteByte(b)
	em.inb = em.inbInit
	switch em.inBuf.String() {
	case "#8": // DECALN - fill screen with 'E'
		size := em.be.GetSize()
		em.autoWrap = false
		em.setPosition(Coord{0, 0})
		em.be.SetAttr(em.style.Attr())
		for row := range size.Y {
			for col := range size.X {
				em.be.PutRune(Coord{X: col, Y: row}, 'E', 1)
			}
		}
		// most implementations leave the cursor at home for this
		em.setPosition(Coord{0, 0})

		// case "%@": // TODO: select 8859-1
		// case "%G": // TODO: select UTF-8
		// case "(A": // TODO: select G0 as UK
		// case "(B": // TODO: select G0 as US
		// case "(C", "(5": // TODO: select G0 as Finnish
		// case "(H", "(7": // TODO: select G0 as Swedish
		// case "(K": // TODO: select G0 as German
		// case "(Q", "(9": // TODO: select G0 as French Canadian
		// case "(R", "(f": // TODO: select G0 as French
		// case "(Y": // TODO: select G0 as Italian
	}
}

// inbCSI handles bytes that are part of a CSI based sequence.
func (em *emulator) inbCSI(b byte) {
	if (b >= 0x30) && (b <= 0x3F) {
		em.inBuf.WriteByte(b) // parameter bytes
	} else if (b >= 0x20) && (b <= 0x2F) {
		em.inBuf.WriteByte(b) // intermediate bytes
	} else if b >= 0x40 && (b <= 0x7F) {
		em.inb = em.inbInit
		em.processCsi(b)
	} else {
		// error state
		em.beep()
		em.inb = em.inbInit
	}
}

// inbOSC handles bytes that are part of on OSC sequences (operating system command).
func (em *emulator) inbOSC(b byte) {
	switch b {
	case 0x9c, 0x07:
		em.inb = em.inbInit
		em.processOSC()
	case '\\':
		if buf := em.inBuf.Bytes(); len(buf) > 0 && buf[len(buf)-1] == 0x1b {
			em.inb = em.inbInit
			em.inBuf.Truncate(em.inBuf.Len() - 1)
			em.processOSC()
		} else {
			em.inBuf.WriteByte(b)
		}
	default:
		em.inBuf.WriteByte(b)
	}

}

// inbStr handles PM, SOS, and any other string we want to consume and discard.
func (em *emulator) inbStr(b byte) {
	switch b {
	case 0x9c, 0x07:
		em.inb = em.inbInit
	case '\\':
		if buf := em.inBuf.Bytes(); len(buf) > 0 && buf[len(buf)-1] == 0x1b {
			em.inb = em.inbInit
			em.inBuf.Truncate(em.inBuf.Len() - 1)
		} else {
			em.inBuf.WriteByte(b)
		}
	default:
		em.inBuf.WriteByte(b)
	}
}

// inbUTF handles continuation bytes for UTF-8 sequences.
func (em *emulator) inbUTF(b byte) {
	if b&0xC0 == 0x80 {
		// good continuation byte
		em.inBuf.WriteByte(b)
		if em.inBuf.Len() == em.utfLen {
			em.inb = em.inbInit
			r, _, err := em.inBuf.ReadRune()
			if err != nil {
				em.beep()
			} else {
				em.putRune(r)
			}
		}
	} else {
		em.beep()
		em.inb = em.inbInit
	}
}

func (em *emulator) beep() {
	if beeper, ok := em.be.(Beeper); ok {
		beeper.Beep()
	}
}

// numericParams splits the string consisting of numeric parameters into integers.
// It ensures a minimum number are present (needed for some safety cases).
// Empty strings default to zero.
func numericParams(str string, minimumLen int) ([]int, error) {
	ps := strings.Split(str, ";")
	pi := make([]int, max(len(ps), minimumLen))
	for i, str := range ps {
		if str != "" {
			if iv, err := strconv.Atoi(str); err != nil {
				return nil, err
			} else {
				pi[i] = iv
			}
		}
	}
	return pi, nil
}

// splitSgrArgs grabs either 2 arguments, or 4 arguments for palette or rgb values
// used with SGR 38 and 48.
func splitSgrArgs(args []string, words []string) ([]string, []string) {
	if len(args) > 0 {
		return args, words
	}
	if len(words) == 0 {
		return nil, nil
	}
	switch i, _ := strconv.Atoi(words[0]); i {
	case 2: // RGB value follows, 3 parameters
		if len(words) < 4 {
			return nil, nil
		}
		return words[:4], words[4:]
	case 5: // single palette index follows
		if len(words) < 2 {
			return nil, nil
		}
		return words[:2], words[2:]
	}
	return nil, nil
}

// processSgr processes SGR commands (things that change how characters are displayed).
func (em *emulator) processSgr(str string) {
	words := strings.Split(str, ";")

	// technically parameters for 38 or 48 should be separated by colons, but due to historical
	// accident it is more common to see semicolon separation.  Underline styles are also separated
	// by a colon, if present.
	if len(words) == 0 {
		words = []string{"0"}
	}
	for len(words) > 0 {
		// we do this instead of a range so we can lop off
		// multiple words for SGR38 and 48.
		word := words[0]
		words = words[1:]

		if word == "" {
			word = "0"
		}
		args := []string(nil)
		if strings.Contains(word, ":") {
			args = strings.Split(word, ":")
			word = args[0]
			args = args[1:]
		}

		v, err := strconv.Atoi(word)
		if err != nil {
			// just swallow it for now
			return
		}
		switch v {
		case 0:
			em.style = BaseStyle
			em.be.SetAttr(Plain)
			if c, ok := em.be.(Colorer); ok && c.Colors() > 0 {
				c.SetFgColor(color.Reset)
				c.SetBgColor(color.Reset)
			}
		case 1:
			em.style = em.style.WithAttr((em.style.Attr() &^ Dim) | Bold)
			em.be.SetAttr(em.style.Attr())
		case 2:
			em.style = em.style.WithAttr((em.style.Attr() &^ Bold) | Dim)
			em.be.SetAttr(em.style.Attr())
		case 3:
			em.style = em.style.WithAttr(em.style.Attr() | Italic)
			em.be.SetAttr(em.style.Attr())
		case 4:
			em.style = em.style.WithAttr((em.style.Attr() &^ UnderlineMask) | Underline)

			if len(args) > 0 {
				switch args[0] {
				case "2":
					em.style = em.style.WithAttr(em.style.Attr() | DoubleUnderline)
				case "3":
					em.style = em.style.WithAttr(em.style.Attr() | CurlyUnderline)
				case "4":
					em.style = em.style.WithAttr(em.style.Attr() | DottedUnderline)
				case "5":
					em.style = em.style.WithAttr(em.style.Attr() | DashedUnderline)
				}
			}
			em.be.SetAttr(em.style.Attr())
		case 5, 6:
			em.style = em.style.WithAttr(em.style.Attr() | Blink)
			em.be.SetAttr(em.style.Attr())
		case 7:
			em.style = em.style.WithAttr(em.style.Attr() | Reverse)
			em.be.SetAttr(em.style.Attr())
		case 8: // ignore, its for invisible
		case 9:
			em.style = em.style.WithAttr(em.style.Attr() | StrikeThrough)
			em.be.SetAttr(em.style.Attr())
		case 21: // Doubly underlined, per ECMA
			em.style = em.style.WithAttr((em.style.Attr() &^ UnderlineMask) | DoubleUnderline)
			em.be.SetAttr(em.style.Attr())
		case 22:
			em.style = em.style.WithAttr(em.style.Attr() &^ (Bold | Dim))
			em.be.SetAttr(em.style.Attr())
		case 23:
			em.style = em.style.WithAttr(em.style.Attr() &^ Italic)
			em.be.SetAttr(em.style.Attr())
		case 24:
			em.style = em.style.WithAttr(em.style.Attr() &^ UnderlineMask)
			em.be.SetAttr(em.style.Attr())
		case 25:
			em.style = em.style.WithAttr(em.style.Attr() &^ Blink)
			em.be.SetAttr(em.style.Attr())
		case 27:
			em.style = em.style.WithAttr(em.style.Attr() &^ Reverse)
			em.be.SetAttr(em.style.Attr())
		case 29:
			em.style = em.style.WithAttr(em.style.Attr() &^ StrikeThrough)
			em.be.SetAttr(em.style.Attr())

		case 30, 31, 32, 33, 34, 35, 36, 37: // simple foreground colors
			if c, ok := em.be.(Colorer); ok && c.Colors() > 0 {
				em.style = em.style.WithFg(color.Black + color.Color(v-30))
				c.SetFgColor(em.style.Fg())
			}
		case 38:
			args, words = splitSgrArgs(args, words)
		case 39:
			if c, ok := em.be.(Colorer); ok && c.Colors() > 0 {
				em.style = em.style.WithFg(color.Reset)
				c.SetFgColor(em.style.Fg())
			}
		case 40, 41, 42, 43, 44, 45, 46, 47: // simple background colors
			if c, ok := em.be.(Colorer); ok && c.Colors() > 0 {
				em.style = em.style.WithBg(color.Black + color.Color(v-40))
				c.SetBgColor(em.style.Bg())
			}
		case 48: // TODO:
			args, words = splitSgrArgs(args, words)
		case 49:
			if c, ok := em.be.(Colorer); ok && c.Colors() > 0 {
				em.style = em.style.WithBg(color.Reset)
				c.SetBgColor(em.style.Bg())
			}
		case 53:
			em.style = em.style.WithAttr(em.style.Attr() | Overline)
			em.be.SetAttr(em.style.Attr())
		case 55:
			em.style = em.style.WithAttr(em.style.Attr() &^ Overline)
			em.be.SetAttr(em.style.Attr())
		}
	}
}

// processCursorUp implements CUU.
func (em *emulator) processCursorUp(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveUpN(Row(max(1, pi[0])))
	}
}

// processCursorDown implements CUD.
func (em *emulator) processCursorDown(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveDownN(Row(max(1, pi[0])))
	}
}

// processCursorForward implements CUF.
func (em *emulator) processCursorForward(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveRightN(Col(max(1, pi[0])))
	}
}

// processCursorBackward implements CUB.
func (em *emulator) processCursorBackward(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveLeftN(Col(max(1, pi[0])))
	}
}

// processCursorNextLine implements CNL.
func (em *emulator) processCursorNextLine(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveDownN(Row(max(1, pi[0])))
		pos := em.getPosition()
		pos.X = 0
		em.setPosition(pos)
	}
}

// processCursorPreviousLine implements CPL.
func (em *emulator) processCursorPreviousLine(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		em.moveUpN(Row(max(1, pi[0])))
		pos := em.getPosition()
		pos.X = 0
		em.setPosition(pos)
	}
}

// processCursorColumn implements CHA.
func (em *emulator) processCursorColumn(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		pos := em.getPosition()
		pos.X = min(Col(max(1, pi[0])), em.be.GetSize().X) - 1
		em.setPosition(pos)
	}
}

// processCursorPosition implements CUP, and also HVP.
func (em *emulator) processCursorPosition(str string) {
	if pi, err := numericParams(str, 2); err == nil {
		em.autoWrap = false
		pos := em.getPosition()
		wsz := em.be.GetSize()
		row := Row(max(1, pi[0]))
		col := Col(max(1, pi[1]))
		row = max(1, min(row, wsz.Y))
		col = max(1, min(col, wsz.X))
		pos.X = col - 1
		pos.Y = row - 1
		em.setPosition(pos)
	}
}

// processCursorTab implements CHT.
func (em *emulator) processCursorTab(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		// Note: tab does not clear this field.
		for range max(1, pi[0]) {
			em.nextTab()
		}
	}
}

// processCursorBackTab implements CBT.
func (em *emulator) processCursorBackTab(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		em.autoWrap = false
		for range max(1, pi[0]) {
			em.prevTab()
		}
	}
}

// processEraseDisplay implements ED.
func (em *emulator) processEraseDisplay(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		switch pi[0] {
		case 0: // erase below
			em.eraseBelow()
		case 1: // erase above
			em.eraseAbove()
		case 2: // erase all
			em.eraseAll()
			// others not supported (3 is erase saved lines)
		}
	}
}

// processEraseLine implements EL.
func (em *emulator) processEraseLine(str string) {
	if pi, err := numericParams(str, 1); err == nil {
		switch pi[0] {
		case 0:
			em.eraseToLineEnd()
		case 1:
			em.eraseToLineStart()
		case 2:
			em.eraseLine()
		}
	}
}

// processCsi processes CSI sequences.
func (em *emulator) processCsi(final byte) {
	// CSI sequences are supported in several different possible ways:
	// parameters may have a prefix character that is not numeric, typically
	// indicating a whole different mode of operation than the final byte.
	// There may also be intermediate bytes, but we only look for one, because
	// the use cases we have this are that only a single intermediate byte is
	// sometimes used to affect function.  (E.g. $ in some cases.)
	cmd := ""
	if em.inBuf.Len() > 0 {
		if b := em.inBuf.Bytes()[0]; b > '9' && b <= '?' {
			cmd += string(b)
			em.inBuf.ReadByte()
		}
	}
	if l := em.inBuf.Len(); l > 0 {
		if b := em.inBuf.Bytes()[l-1]; b >= 0x20 && b <= 0x2F {
			cmd += string(b)
			em.inBuf.Truncate(l - 1)
		}
	}
	cmd += string(final)

	str := em.inBuf.String()
	switch cmd {

	case "A":
		em.processCursorUp(str)
	case "B":
		em.processCursorDown(str)
	case "C":
		em.processCursorForward(str)
	case "D":
		em.processCursorBackward(str)
	case "E":
		em.processCursorNextLine(str)
	case "F":
		em.processCursorPreviousLine(str)
	case "G":
		em.processCursorColumn(str)
	case "H", "f":
		em.processCursorPosition(str)
	case "I":
		em.processCursorTab(str)
	case "J":
		em.processEraseDisplay(str)
	case "K":
		em.processEraseLine(str)
	case "Z":
		em.processCursorBackTab(str)

	case "c":
		if pi, err := numericParams(str, 1); err == nil && pi[0] == 0 {
			em.sendDA()
		}
	case "d": // move to specific row (VPA)
		if pi, err := numericParams(str, 1); err == nil {
			pos := em.getPosition()
			pos.Y = min(Row(max(1, pi[0])), em.be.GetSize().Y) - 1
			em.setPosition(pos)
		}
	case "e": // advance by rows (VPR)
		if pi, err := numericParams(str, 1); err == nil {
			pos := em.getPosition()
			pos.Y = min(pos.Y+Row(max(1, pi[0])), em.be.GetSize().Y-1)
			em.setPosition(pos)
		}
	case "g": // tab clear (TBC)
		if pi, err := numericParams(str, 1); err == nil {
			switch pi[0] {
			case 0: // clear stop at current column
				em.clrTabStop(em.getPosition().X)
			case 3: // clear all columns
				em.tabStops = []Col{} // this is distinct from nil
			}
		}
	case "m":
		em.processSgr(str)
	case "n":
		em.deviceReport(str)
	case "?W":
		if pi, err := numericParams(str, 1); err == nil && pi[0] == 5 {
			// DECST8C - reset tab stops to default (VT510)
			em.tabStops = nil
		}
	case "?h": // DECSET
		if pi, err := numericParams(str, 1); err == nil {
			for _, pm := range pi {
				em.setPrivateMode(PrivateMode(pm), ModeOn)
			}
		}
	case "?l": // DECRST
		if pi, err := numericParams(str, 1); err == nil {
			for _, pm := range pi {
				em.setPrivateMode(PrivateMode(pm), ModeOff)
			}
		}
	case "?$p": // DECRQM - only a single numeric parameter (mode number) can be supplied (VT300+)
		if pm, err := strconv.Atoi(str); err == nil {
			status := em.getPrivateMode(PrivateMode(pm))
			em.SendRaw(fmt.Appendf(nil, "\x1b[?%d;%d$y", pm, status))
		}
	case ">q":
		if pi, err := numericParams(str, 1); err == nil && pi[0] == 0 && em.name != "" {
			em.SendRaw(fmt.Appendf(nil, "\x1bP>|%s %s\x1b\\", em.name, em.vers))
		}
	}
}

// processOSC processes an operating system command.
// TODO: add support for these - e.g. OSC 8 for hyperlinks, OSC 52 for clipboard access, etc.
func (em *emulator) processOSC() {
	// Every OSC we support has a number, semicolon, then string.
	ns, str, ok := strings.Cut(em.inBuf.String(), ";")
	if !ok {
		return
	}
	if num, err := strconv.Atoi(ns); err != nil {
		return
	} else {
		switch num {
		case 2: // Set window title
			if t, ok := em.be.(Titler); ok {
				// TODO: possibly validate the UTF-8 content?
				t.SetWindowTitle(str)
			}
		}
	}
}

func (em *emulator) getPosition() Coord {
	pos := em.be.GetPosition()
	em.pos = pos
	return em.pos
}

func (em *emulator) setPosition(pos Coord) {
	em.pos = pos
	em.be.SetPosition(pos)
}

func (em *emulator) deviceReport(s string) {
	switch s {
	case "5":
		em.SendRaw([]byte("\x1b[0n"))
	case "6":
		pos := em.getPosition()
		em.SendRaw(fmt.Appendf(nil, "\x1b[%d;%dR", pos.Y+1, pos.X+1))
	default: // ignore
	}
}

func (em *emulator) moveUpN(count Row) {
	for range count {
		em.moveUp()
	}
}

func (em *emulator) moveDownN(count Row) {
	for range count {
		em.moveDown()
	}
}

func (em *emulator) moveLeftN(count Col) {
	for range count {
		em.moveLeft()
	}
}

func (em *emulator) moveRightN(count Col) {
	for range count {
		em.moveRight()
	}
}

func (em *emulator) moveDown() {
	pos := em.getPosition()
	win := em.be.GetSize()
	if pos.Y == win.Y-1 {
		// TODO: scroll
	} else {
		pos.Y++
		em.setPosition(pos)
	}
}

func (em *emulator) moveUp() {
	pos := em.getPosition()
	if pos.Y == 0 {
		// TODO: scroll
	} else {
		pos.Y--
		em.setPosition(pos)
	}
}

func (em *emulator) moveLeft() {
	em.autoWrap = false
	pos := em.getPosition()
	if pos.X > 0 {
		pos.X--
		em.setPosition(pos)
	}
}

func (em *emulator) moveRight() {
	pos := em.getPosition()
	win := em.be.GetSize()
	if pos.X < win.X-1 {
		pos.X++
		em.setPosition(pos)
	}
}

// nextLine is like CNL with 1, but it optionally also scrolls.
func (em *emulator) nextLine() {
	em.autoWrap = false
	em.moveDown()
	em.pos.X = 0
	em.setPosition(em.pos)
}

// nextTab advances to the next tab stop, or the end of
// the line if there is no further tab.
func (em *emulator) nextTab() {
	maxX := em.be.GetSize().X - 1
	curX := em.getPosition().X
	if curX == maxX { // already at end
		return
	}
	nextX := maxX
	if em.tabStops == nil {
		// just advance to the next one
		nextX = min((curX+8)&^7, maxX)
	} else {
		for _, p := range em.tabStops {
			if p > curX {
				nextX = p
				break
			}
		}
	}
	em.setPosition(Coord{X: nextX, Y: em.pos.Y})
}

func (em *emulator) prevTab() {
	curX := em.getPosition().X
	if curX == 0 {
		return
	}
	nextX := curX - 1
	if em.tabStops == nil {
		nextX &^= 7
	} else if i, exist := slices.BinarySearch(em.tabStops, nextX); exist {
		nextX = em.tabStops[i]
	} else if i > 0 {
		nextX = em.tabStops[i-1]
	} else {
		nextX = 0
	}
	em.setPosition(Coord{X: nextX, Y: em.pos.Y})
}

// initTabStops initializes the tab stops assuming every 8th column
// is a tab stop.  This should only be called if the user is intentionally
// changing the tab stops, because it will no longer support expanding
// tab stops on resizing.
func (em *emulator) initTabStops() {
	if em.tabStops == nil {
		// no tab stop at offset 0 since that would be pointless
		em.tabStops = make([]Col, 0, int(em.be.GetSize().X/8)+1)
		for col := Col(8); col < em.be.GetSize().X; col += 8 {
			em.tabStops = append(em.tabStops, col)
		}
	}
}

// setTabStop sets a tab stop at the given location.
// This calls  initTabStops - please see the description of that function for ramifications.
func (em *emulator) setTabStop(ts Col) {
	em.initTabStops()
	if index, exist := slices.BinarySearch(em.tabStops, ts); !exist {
		em.tabStops = slices.Insert(em.tabStops, index, ts)
	}
}

// clrTabStop clears the tab stop at the given column.  This calls
// initTabStops - please see the description of that function for ramifications.
func (em *emulator) clrTabStop(ts Col) {
	em.initTabStops()
	em.tabStops = slices.DeleteFunc(em.tabStops, func(x Col) bool { return x == ts })
}

// index obtains the index in the cells slice for the given coordinates,
// which must be within the bounds of the display size.
func (em *emulator) index(c Coord) int {
	dim := em.be.GetSize()
	return int(c.Y)*int(dim.X) + int(c.X)
}

// putRune puts out a single rune.  This might be a subsequent part of a grapheme cluster, in
// which case it will be emitted together with the preceding base character.
func (em *emulator) putRune(r rune) {
	dim := em.be.GetSize()
	em.be.SetAttr(em.style.Attr())

	if lastIdx := em.lastIndex; lastIdx != 0 {
		lastIdx--
		if pm := em.getPrivateMode(PmGraphemeClusters); pm == ModeOn || pm == ModeOnLocked {
			// maybe we need to update the last index
			str := em.cells[lastIdx].C + string(r)
			if cs, rest, width, _ := uniseg.FirstGraphemeClusterInString(str, -1); rest == "" {
				// we are adding to a cluster
				em.cells[lastIdx].C = cs
				col := Col(lastIdx) % dim.X
				row := Row(lastIdx / int(dim.X))
				// we may have to move position if this switches to wide, so recalculate expected end
				end := (col + Col(width)) % dim.X
				if em.getPrivateMode(PmAutoMargin) == ModeOn && end >= dim.X {
					em.autoWrap = true
				}
				if width == 2 && col < dim.X-1 {
					// erase the next cell before putting down a character
					em.cells[lastIdx+1].C = ""
					em.cells[lastIdx+1].S = em.cells[lastIdx].S
				}
				// we leave the em.lastIndex for now, we might keep extending this cluster
				em.be.PutGrapheme(Coord{X: col, Y: row}, cs, width)
				em.setPosition(Coord{X: end, Y: row})
				return
			}
		}
	}

	if em.autoWrap {
		em.nextLine()
	}

	autoMargin := em.getPrivateMode(PmAutoMargin) == ModeOn

	pos := em.getPosition()
	w := uniseg.StringWidth(string(r))
	if autoMargin && pos.X+Col(w) >= dim.X {
		em.autoWrap = true
	}
	index := em.index(pos)
	em.cells[index].C = string(r)
	em.cells[index].S = em.style
	em.be.PutRune(em.pos, r, w)
	em.lastIndex = index + 1

	if w == 2 && pos.X < dim.X-1 {
		index++
		em.cells[index].C = ""
		em.cells[index].S = em.style
	}
	// Advance the cursor. This will stop at the margin.
	// Note that if auto margin is enabled, we will have set
	// autoWrap above if we were at the margin already.
	em.moveRightN(Col(w))
}

// eraseCell erases a single cell at the given offset.
// It clears attributes, but leaves the colors intact.
func (em *emulator) eraseCell(c Coord) {
	em.be.PutRune(c, 0, 0)
}

// eraseBelow erases from (and including) the current cursor position to the end of the window.
func (em *emulator) eraseBelow() {
	size := em.be.GetSize()
	pos := em.getPosition()
	for x := pos.X; x < size.X; x++ {
		em.eraseCell(Coord{X: x, Y: pos.Y})
	}
	for y := pos.Y + 1; y < size.Y; y++ {
		for x := Col(0); x < size.X; x++ {
			em.eraseCell(Coord{X: x, Y: y})
		}
	}
	em.setPosition(pos)
}

// eraseAbove erases from the origin to (and including) the current cursor position.
func (em *emulator) eraseAbove() {
	size := em.be.GetSize()
	pos := em.getPosition()
	for y := Row(0); y < pos.Y; y++ {
		for x := Col(0); x < size.X; x++ {
			em.eraseCell(Coord{X: x, Y: y})
		}
	}
	for x := Col(0); x <= pos.X; x++ {
		em.eraseCell(Coord{X: x, Y: pos.Y})
	}
	em.setPosition(pos)
}

// eraseAll erases the entire screen. It uses the color, but resets all other attributes.
func (em *emulator) eraseAll() {
	size := em.be.GetSize()
	pos := em.getPosition()
	for y := Row(0); y < size.Y; y++ {
		for x := Col(0); x < size.X; x++ {
			em.eraseCell(Coord{X: x, Y: y})
		}
	}
	em.setPosition(pos)
}

// eraseToLineEnd erases to the end of the line, including the cursor position.
func (em *emulator) eraseToLineEnd() {
	size := em.be.GetSize()
	pos := em.getPosition()
	for x := pos.X; x < size.X; x++ {
		em.eraseCell(Coord{x, pos.Y})
	}
	em.setPosition(pos)
}

// eraseToLineStart erases to the start of the line, including the cursor position.
func (em *emulator) eraseToLineStart() {
	pos := em.getPosition()
	for x := Col(0); x <= pos.X; x++ {
		em.eraseCell(Coord{x, pos.Y})
	}
	em.setPosition(pos)
}

// eraseLine erases the entire line.
func (em *emulator) eraseLine() {
	size := em.be.GetSize()
	pos := em.getPosition()
	for x := range size.X {
		em.eraseCell(Coord{x, pos.Y})
	}
	em.setPosition(pos)
}

// softReset performs a soft reset.
func (em *emulator) softReset() {
	// TODO:
	// Select default character sets
	em.tabStops = nil
	em.autoWrap = false
	em.style = BaseStyle
	em.saved = savedCursor{style: BaseStyle}
	em.be.Reset()
	// start by resetting all modes
	for pm := range em.localModes {
		em.setPrivateMode(pm, ModeOff)
	}
	// and set any that should reset on (auto-margin)
	em.setPrivateMode(PmAutoMargin, ModeOn)
	em.setPrivateMode(PmShowCursor, ModeOn)
	em.setPosition(Coord{0, 0})
	em.eraseAll()
}

// sendDA ends the primary device attributes.
func (em *emulator) sendDA() {
	buf := &bytes.Buffer{}
	_, _ = fmt.Fprintf(buf, "\x1b[?63")
	if _, ok := em.be.(Colorer); ok {
		fmt.Fprintf(buf, ";22")
	}
	// 9 for NRC?
	// 15 for graphics?
	// 52 for clipboard access?
	buf.WriteRune('c')
	em.SendRaw(buf.Bytes())
}

// getPrivateMode returns the value of a DEC private mode.
func (em *emulator) getPrivateMode(pm PrivateMode) ModeStatus {
	if ms, ok := em.localModes[pm]; ok {
		return ms
	}
	return em.be.GetPrivateMode(pm)
}

// setPrivateMode sets the DEC private mode.
func (em *emulator) setPrivateMode(pm PrivateMode, ms ModeStatus) {
	if ms != ModeOn && ms != ModeOff {
		return
	}
	if old, ok := em.localModes[pm]; ok && (old == ModeOn || old == ModeOff) {
		em.localModes[pm] = ms
	} else {
		_ = em.be.SetPrivateMode(pm, ms)
	}
}

// SendRaw allows raw data to be sent to the application.
// This is done in a thread-safe way, so that content is not intermingled.
func (em *emulator) SendRaw(b []byte) {
	em.sendLock.Lock()
	defer em.sendLock.Unlock()

	// Do not attempt to send *anything* if we are stopped.
	select {
	case <-em.stopQ:
		return
	default:
	}

	// Try to write to the readQ, but if we cannot, then wait until
	// either we can, or the stopQ is fired.  This ensures that we avoid
	// breaking up content if at all possible.
	for _, ch := range b {
		select {
		case em.readQ <- ch:
		default:
			select {
			case em.readQ <- ch:
			case <-em.stopQ:
				return
			}
		}
	}
}

// KbdEvent injects a keyboard event into the emulator
func (em *emulator) KeyEvent(ev KbdEvent) {
	// TODO: more add support for other keyboard protocols, right now we only do legacy
	em.keyLegacy(ev)
}

// ResizeEvent is called by the backend when a resize occurs.  A real backend with a child
// process (essentially a "real emulator") should probably also fire SIGWINCH if appropriate.
// That would be the job of something other than this code.
func (em *emulator) ResizeEvent() {
	// reload our position it may have changed
	em.pos = em.getPosition()
	if em.getPrivateMode(PmResizeReports) == ModeOn { // NB: we never support "ModeOnLocked"
		newSz := em.be.GetSize()
		// NB: for now we do not support pixel sizes
		em.SendRaw(fmt.Appendf(nil, "\x1b[48;%d;%d;0;0t", newSz.Y, newSz.X))
	}
}

var legacyKeys = map[KeyCode]struct {
	K  string // unmodified key
	A  string // unmodified in application cursor mode (smkx)
	S  string // with shift (if empty use regular modifier)
	C  string // with control (if empty use regular modifier)
	CS string // with ctrl-shift
}{
	KcF1:        {K: "\x1bOP"}, // SS3 P
	KcF2:        {K: "\x1bOQ"}, // SS3 Q
	KcF3:        {K: "\x1bOR"}, // SS3 R
	KcF4:        {K: "\x1bOS"}, // SS3 S
	KcF5:        {K: "\x1b[15~"},
	KcF6:        {K: "\x1b[17~"},
	KcF7:        {K: "\x1b[18~"},
	KcF8:        {K: "\x1b[19~"},
	KcF9:        {K: "\x1b[20~"},
	KcF10:       {K: "\x1b[21~"},
	KcF11:       {K: "\x1b[23~"},
	KcF12:       {K: "\x1b[24~"},
	KcF13:       {K: "\x1b[25~"},
	KcF14:       {K: "\x1b[26~"},
	KcF15:       {K: "\x1b[28~"},
	KcF16:       {K: "\x1b[29~"},
	KcF17:       {K: "\x1b[31~"},
	KcF18:       {K: "\x1b[32~"},
	KcF19:       {K: "\x1b[33~"},
	KcF20:       {K: "\x1b[34~"},
	KcUp:        {K: "\x1b[A", A: "\x1bOA"},
	KcDown:      {K: "\x1b[B", A: "\x1bOB"},
	KcRight:     {K: "\x1b[C", A: "\x1bOC"},
	KcLeft:      {K: "\x1b[D", A: "\x1bOD"},
	KcHome:      {K: "\x1b[H", A: "\x1bOH"},
	KcEnd:       {K: "\x1b[F", A: "\x1bOF"},
	KcPgUp:      {K: "\x1b[5~"},
	KcPgDn:      {K: "\x1b[6~"},
	KcDel:       {K: "\x1b[3~"},
	KcIns:       {K: "\x1b[2~"},
	KcHelp:      {K: "\x1b[28~"}, // also F15
	KcMenu:      {K: "\x1b[29~"}, // also F16
	KcTab:       {K: "\t", S: "\x1b[Z", CS: "\x1b[Z"},
	KcBackspace: {K: "\x7f", S: "\x7f", C: "\x08", CS: "\x08"},
	KcDelete:    {K: "\x08", S: "\x08", C: "\x7f", CS: "\x7f"},
	KcSpace:     {K: " ", S: " ", C: "\x00", CS: "\x00"},
	KcReturn:    {K: "\r", S: "\r", CS: "\r"},
	KcEsc:       {K: "\x1b", S: "\x1b", C: "\x1b"},

	// These ones are weird legacy control sequences that we mostly
	// do not care about.  We don't include shifted variants.
	KeyCode('2'): {K: "2", C: "\x00"},
	KeyCode('3'): {K: "3", C: "\x1b"},
	KeyCode('4'): {K: "4", C: "\x1c"},
	KeyCode('5'): {K: "5", C: "\x1d"},
	KeyCode('6'): {K: "6", C: "\x1e"},
	KeyCode('7'): {K: "7", C: "\x1f"},
	KeyCode('8'): {K: "8", C: "\x7f"},
	KeyCode('['): {K: "[", C: "\x1b"},
	KeyCode('/'): {K: "/", C: "\x1c"},
	KeyCode(']'): {K: "]", C: "\x1d"},
	KeyCode('~'): {K: "~", C: "\x1e"},
	KeyCode('?'): {K: "?", C: "\x1f"},
}

// toASCIIUpper returns the equivalent upper case ASCII (and true),
// if the input is an ASCII letter.  Otherwise it returns 0, false.
func toASCIIUpper(r rune) (rune, bool) {
	if r >= 'a' && r <= 'z' {
		return (r - 32), true
	} else if r >= 'A' && r <= 'Z' {
		return r, true
	}
	return 0, false
}

// keyLegacy handles a keyboard event when in legacy vt220 style mode.
func (em *emulator) keyLegacy(ev KbdEvent) {
	if !ev.Down { // legacy protocol does not support key release
		return
	}
	if ev.Mod&(ModHyper|ModMeta) != 0 { // legacy protocol does not support these
		return
	}

	if v, ok := legacyKeys[ev.Code]; ok {
		str := ""
		match := false
		switch ev.Mod & (ModShift | ModCtrl) {
		case ModNone:
			if em.getPrivateMode(PmAppCursor) == ModeOn && v.A != "" {
				str = v.A
			} else {
				str = v.K
			}
			match = true
		case ModShift:
			if str = v.S; str != "" {
				match = true
			}
		case ModCtrl:
			if str = v.C; str != "" {
				match = true
			}
		case ModCtrl | ModShift:
			if str = v.CS; str != "" {
				match = true
			}
		}
		if !match {
			// No specific modifiers present, lets add them. There are two cases,
			// one for SS3 based keys and another for CSI based keys.  SS3 based
			// keys are converted to CSI - 1 ; mod ; final
			// Note: legacy encoding does not use modifiers for alt or super - alt will be
			// determined by sending an escape prefix.
			mod := 0
			if ev.Mod&ModShift != 0 {
				mod |= 1
			}
			if ev.Mod&ModCtrl != 0 {
				mod |= 4
			}
			if strings.HasPrefix(v.K, "\x1bO") {
				str = fmt.Sprintf("\x1b[1;%d%c", mod+1, v.K[len(v.K)-1])
			} else {
				str = fmt.Sprintf("%s;%d%c", v.K[:len(v.K)-1], mod+1, v.K[len(v.K)-1])
			}
		}
		if ev.Mod&ModAlt != 0 {
			em.SendRaw(append([]byte{'\x1b'}, []byte(str)...)) // alt sends leading escape
		} else {
			em.SendRaw([]byte(str))
		}
		return
	}

	// fallback control key handling
	if u, ok := toASCIIUpper(rune(ev.Code)); ok && ev.Mod&ModCtrl != 0 {
		b := byte(u) - 'A' + 1
		if ev.Mod&ModAlt != 0 {
			em.SendRaw([]byte{'\x1b', b})
		} else {
			em.SendRaw([]byte{b})
		}
		return
	}

	if ev.Code > KcSpace && ev.Code < 0x7F && ev.Mod&ModCtrl == ModNone {
		if ev.Mod&ModAlt != 0 {
			em.SendRaw([]byte{'\x1b'})
		}
		em.SendRaw([]byte{byte(ev.Code)})
		return
	}
}

// SetId sets the terminal name and version.
func (em *emulator) SetId(name string, version string) {
	em.name = name
	em.vers = version
}

// Start the terminal emulator.
func (em *emulator) Start() error {
	select {
	case <-em.stopQ:
	default:
		// already running
		return errors.New("terminal already started")
	}
	stopQ := make(chan bool)
	em.stopQ = stopQ
	go em.run(stopQ)
	return nil
}

// Stop the terminal emulator.  This also wakes any blocked
// Read or Write calls, which will return an error.
func (em *emulator) Stop() error {
	select {
	case <-em.stopQ:
	default:
		close(em.stopQ)
	}
	return nil
}

// Drain pending output to the terminal emulator.
func (em *emulator) Drain() error {
	q := make(chan bool)
	select {
	case em.writeQ <- q:
	case <-em.stopQ:
	}
	select {
	case <-q:
	case <-em.stopQ:
	}
	// make sure to wake the reader
	select {
	case em.readQ <- true:
	default:
	}
	return nil
}

// Write data to the emulator (commands).
func (em *emulator) Write(data []byte) (n int, err error) {
	stopQ := em.stopQ
	writeQ := em.writeQ
	drainQ := make(chan bool)
	select {
	case writeQ <- data:
		// we add the drainQ for synchronization, so that we only
		// return after the the emulator has processed this.
		select {
		case <-stopQ:
			return 0, errors.New("terminal emulator stopped")
		case writeQ <- drainQ:
		}
		select {
		case <-stopQ:
			return 0, errors.New("terminal emulator stopped")
		case <-drainQ:
			return len(data), nil
		}
	case <-stopQ:
		return 0, errors.New("terminal emulator stopped")
	}
}

// Read data (key events, etc.) from the emulator.
func (em *emulator) Read(data []byte) (n int, err error) {
	stopQ := em.stopQ
	readQ := em.readQ

	n = 0
	if len(data) < 1 {
		return 0, nil
	}
	select {
	case <-stopQ:
		return 0, errors.New("terminal emulator stopped")
	case v := <-readQ:
		// The data arriving in the channel may be a byte, or it might be a bool
		// trying to force a wake up.  Note that the bool may be intermingled with other
		// bytes, so we check it. Also data may have arrived since the bool was posted,
		// so make sure we don't terminate until we have collected all the relevant data
		// that we can (up to the limit of what was requested.)
		if ch, ok := v.(byte); ok {
			data[n] = ch
			n++
		}
		for n < len(data) {
			select {
			case v = <-readQ:
				if ch, ok := v.(byte); ok {
					data[n] = ch
					n++
				}
			default:
				return n, nil
			}
		}
		return n, nil
	}
}

func (em *emulator) run(stopQ <-chan bool) {
	for {
		select {
		case item := <-em.writeQ:
			switch d := item.(type) {
			case byte:
				em.inb(d)
			case []byte:
				for _, ch := range d {
					em.inb(ch)
				}
			case chan bool:
				close(d)
			}
		case <-stopQ:
			return
		}
	}
}
