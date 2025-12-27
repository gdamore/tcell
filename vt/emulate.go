// Copyright 2025 The TCell Authors
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
	"strconv"
	"strings"

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

	// Drain waits until any queued but not processed input has finished processing.
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

// NewEmulator creates an emulator instance on top of the given backend.
// The input is relative to the emulator, so it receives data from the host,
// whereas the emulator sends data to the application through the output.
func NewEmulator(be Backend) Emulator {
	stopQ := make(chan bool)
	em := &emulator{
		be:         be,
		inBuf:      &bytes.Buffer{},
		writeQ:     make(chan any),
		readQ:      make(chan byte, 128),
		stopQ:      stopQ,
		localModes: map[PrivateMode]ModeStatus{PmAutoMargin: ModeOn},
	}
	close(stopQ)
	em.inb = em.inbInit
	return em
}

// emulator is an implementation of a terminal emulator built on top of
// a Backend.  It implements the common escape sequence handling and high
// level functionality that a real terminal emulator, or a mock, would need.
type emulator struct {
	stopQ     chan bool
	writeQ    chan any  // queues data from application to emulator
	readQ     chan byte // queues data from emulator to application
	wakeQ     chan bool
	be        Backend
	inBuf     *bytes.Buffer // buffer queued for input
	inb       func(byte)    // input byte function (faster than state switch)
	fg        color.Color   // foreground color
	bg        color.Color   // background color
	ul        color.Color   // underline color
	attr      Attr
	utfLen    int
	pos       Coord
	sevenOnly bool   // only allow 7-bit escapes (needed for KOI8, ShiftJIS, etc.)
	name      string // name of this emulator (used for extended attributes)
	vers      string // version string of this emulator (used for extended attributes)
	savedPos  Coord  // saved via DECSC

	localModes map[PrivateMode]ModeStatus // some modes we handle locally
}

func (em *emulator) inbInit(b byte) {
	em.inBuf.Reset()

	// hot path - just doing ASCII directly.
	if b >= ' ' && b < 0x7f {
		// plain ascii
		em.putc(rune(b))
		return
	}

	// For 8-bit encodings, we treat these as Fe sequences.
	// Basically the same as ESC followed by (b - 0x40).
	// TODO: conditionalize this so that we do not do this if
	// the encoding cannot support it (UTF, 8859, and EUC encodings
	// are all fine here, but others like ShiftJIS or KOI8 might not be).
	if b >= 0x80 && b <= 0x9F && !em.sevenOnly {
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
		em.moveLeft()
	case 0x09: // TODO: tab
	case 0x0a: // NL (newline)
		em.nextLine()
	case 0x0b: // VT (vertical tab, treat as LF)
		em.nextLine()
	case 0x0c: // FF (form feed, treat as LF)
		em.nextLine()
	case 0x0d: // CR (carriage return)
		em.setPosition(Coord{0, em.getPosition().Y})
	case 0x0e: // TODO: SO
	case 0x0f: // TODO: SI
	case 0x18: //TODO Cancel (reset parser)
	default:
		// TODO: consider separating Unicode from other 8-bit charsets
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
			em.beep()
		}
	}
}

func (em *emulator) inbEsc(b byte) {
	// By default, reset to init state. Other states will be set explicitly as needed.
	em.inb = em.inbInit

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
	case 'D': // down one line (IND)
		em.moveDown()
	case 'E': // next line (NEL)
		em.nextLine()
	case 'M': // up one line (RI)
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
		em.savedPos = em.getPosition()
	case '8': // restore cursor (DECRC, VT100)
		em.pos = em.savedPos
		em.setPosition(em.pos)
	case '9': // forward index (DECFI, VT420, not widely supported)
		em.moveRight()
	default:
		// ESC-V and ESC-W are for guarded area (TODO)
		// ESC-H horizontal tab set (note that VT52 uses this for home position)
		em.inb = em.inbInit
	}
}

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
		pos := em.getPosition()
		em.setPosition(Coord{0, 0})
		for row := range size.Y {
			for col := range size.X {
				em.be.PutAbs(Coord{X: col, Y: row}, 'E', em.attr)
			}
		}
		em.setPosition(pos)

		// case "%@": // TODO: select 8859-1
		// case "%G": // TODO: select UTF-8
		// case "(A": // TODO: select G0 as UK
		// case "(B": // TODO: select G0 as US
		// case "(C", "(5": // TODO: select G0 as Finnish
		// case "(H", "(7": // TODO: select G0 as Swedish
		// case "(K": // TODO: select G0 as German
		// case "(Q", "(9": // TODO: select G0 as French Candian
		// case "(R", "(f": // TODO: select G0 as French
		// case "(Y": // TODO: select G0 as Italian
	}
}

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

func (em *emulator) inbOSC(b byte) {
	switch b {
	case 0x9c:
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
	case 0x9c:
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
				em.putc(r)
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
// It ensures a minimum number are present (needed for some safety cases), and ensures
// that default values are filled in, if the individual value is an empty string
func numericParams(str string, minimumLen int, defaultValue int) ([]int, error) {
	ps := strings.Split(str, ";")
	pi := make([]int, max(len(ps), minimumLen))
	for i := range pi {
		pi[i] = defaultValue
	}
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

func (em *emulator) processSgr(str string) {
	words := strings.Split(str, ";")

	// technically parameters for 38 or 48 should be separated by colons, but due to historical
	// accident it is more common to see semicolon separation.  Underline styles are also separated
	// by a colon, if present.
	if len(words) == 0 {
		em.attr = Plain
		return
	}
	for len(words) > 0 {
		// we do this instead of a range so we can lop off
		// multiple words for SGR38 and 48.
		word := words[0]
		words = words[1:]

		if word == "" {
			em.attr = Plain
			continue
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
			em.attr = Plain
		case 1:
			em.attr |= Bold
		case 2:
			em.attr |= Dim
		case 3:
			em.attr |= Italic
		case 4:
			em.attr &^= UnderlineMask
			em.attr |= Underline

			if len(args) > 0 {
				switch args[0] {
				case "2":
					em.attr |= DoubleUnderline
				case "3":
					em.attr |= CurlyUnderline
				case "4":
					em.attr |= DottedUnderline
				case "5":
					em.attr |= DashedUnderline
				}
			}
		case 5, 6:
			em.attr |= Blink // not discriminating between fast and slow blink for now
		case 7:
			em.attr |= Reverse
		case 8: // ignore, its for invisible
		case 9:
			em.attr |= StrikeThrough
		}
	}
}

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

	case "A": // up n times (CUU)
		if pi, err := numericParams(str, 1, 1); err == nil && pi[0] > 0 {
			em.moveUpN(Row(pi[0]))
		}
	case "B": // down n times (CUD)
		if pi, err := numericParams(str, 1, 1); err == nil && pi[0] > 0 {
			em.moveDownN(Row(pi[0]))
		}
	case "C": // forward n times (CUF)
		if pi, err := numericParams(str, 1, 1); err == nil && pi[0] > 0 {
			em.moveRightN(Col(pi[0]))
		}
	case "D": // back n times (CUB)
		if pi, err := numericParams(str, 1, 1); err == nil && pi[0] > 0 {
			em.moveLeftN(Col(pi[0]))
		}
	case "E": // down n times (and reset column) (CNL)
		if pi, err := numericParams(str, 1, 1); err == nil && pi[0] > 0 {
			em.moveDownN(Row(pi[0]))
			pos := em.getPosition()
			pos.X = 0
			em.setPosition(pos)
		}

	case "F": // up n times (and reset column) (CPL)
		if pi, err := numericParams(str, 1, 1); err == nil && pi[0] > 0 {
			em.moveUpN(Row(pi[0]))
			pos := em.getPosition()
			pos.X = 0
			em.setPosition(pos)
		}

	case "G": // cursor column (CHA)
		if pi, err := numericParams(str, 1, 1); err == nil && pi[0] > 0 && pi[0] <= int(em.be.GetSize().X) {
			pos := em.getPosition()
			pos.X = Col(pi[0]) - 1
			em.setPosition(pos)
		}

	case "H", "f": // cursor position (CUP), also (HVP)
		if pi, err := numericParams(str, 2, 1); err == nil {
			pos := em.getPosition()
			wsz := em.be.GetSize()
			row := Row(pi[0])
			col := Col(pi[1])
			row = max(1, min(row, wsz.Y))
			col = max(1, min(col, wsz.X))
			pos.X = col - 1
			pos.Y = row - 1
			em.setPosition(pos)
		}
	case "J": // erase in display (ED)
		if pi, err := numericParams(str, 1, 0); err == nil {
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

	case "K": // erase line (EL)
		if pi, err := numericParams(str, 1, 0); err == nil {
			switch pi[0] {
			case 0:
				em.eraseToLineEnd()
			case 1:
				em.eraseToLineStart()
			case 2:
				em.eraseLine()
			}
		}
	case "c":
		if pi, err := numericParams(str, 1, 0); err == nil && pi[0] == 0 {
			em.sendDA()
		}
	case "m":
		em.processSgr(str)
	case "n":
		em.deviceReport(str)
	case "?h": // DECSET
		if pi, err := numericParams(str, 1, 0); err == nil {
			for _, pm := range pi {
				em.setPrivateMode(PrivateMode(pm), ModeOn)
			}
		}
	case "?l": // DECRST
		if pi, err := numericParams(str, 1, 0); err == nil {
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
		if pi, err := numericParams(str, 1, 0); err == nil && pi[0] == 0 && em.name != "" {
			em.SendRaw(fmt.Appendf(nil, "\x1b[P>|%s %s\x1b\\", em.name, em.vers))
		}
	}
}

func (em *emulator) processOSC() {}

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

func (em *emulator) nextLine() {
	em.moveDown()
	em.pos.X = 0
	em.setPosition(em.pos)
}

func (em *emulator) prevLine() {
	em.moveUp()
	em.pos.X = 0
	em.setPosition(em.pos)
}

func (em *emulator) putc(r rune) {
	if p, ok := em.be.(Positioner); ok {
		p.PutChar(r, em.attr)
	} else {
		autoMargin := em.getPrivateMode(PmAutoMargin) == ModeOn
		old := em.getPosition()
		em.be.PutAbs(em.pos, r, em.attr)
		w := uniseg.StringWidth(string(r))
		em.moveRightN(Col(w))
		if autoMargin && old.X+Col(w) >= em.be.GetSize().X {
			em.nextLine()
		}
	}
}

// eraseCell erases a single cell at the given offset.
// It clears attributes, but leaves the colors intact.
func (em *emulator) eraseCell(c Coord) {
	em.be.PutAbs(c, ' ', Plain)
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

// eraseAll erases the entire screen. It uses the color, but resets all other atributes.
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
	// Reset any modes
	// Select default character sets
	// Set cursor
	// Reset colors
	em.attr = Plain
	em.fg = color.None
	em.bg = color.None
	em.ul = color.None
	if c, ok := em.be.(Colorer); ok {
		c.SetFgColor(em.fg)
		c.SetBgColor(em.bg)
	}
	if c, ok := em.be.(UnderlineColorer); ok {
		c.SetUlColor(em.ul)
	}
	em.setPosition(Coord{0, 0})
	// clear the screen
}

// sendDA ends the primary device attributes.
func (em *emulator) sendDA() {
	buf := &bytes.Buffer{}
	_, _ = fmt.Fprintf(buf, "\x1b[63")
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
func (em *emulator) SendRaw(b []byte) {
	for _, ch := range b {
		select {
		case em.readQ <- ch:
		case <-em.stopQ:
			return
		}
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
	wakeQ := em.wakeQ
	readQ := em.readQ

	n = 0
	if len(data) < 1 {
		return 0, nil
	}
	select {
	case <-stopQ:
		return 0, errors.New("terminal emulator stopped")
	case <-wakeQ:
		return 0, nil
	case ch := <-readQ:
		data[n] = ch
		n++
		for n < len(data) {
			select {
			case ch = <-readQ:
				data[n] = ch
				n++
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
