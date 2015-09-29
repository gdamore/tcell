// +build windows

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
	// "fmt"
	"syscall"
	"unsafe"

	"github.com/mattn/go-runewidth"
)

type cScreen struct {
	in    syscall.Handle
	out   syscall.Handle
	mbtns uint32 // debounce mouse buttons
	evch  chan Event
	quit  chan struct{}

	w int
	h int

	oscreen consoleInfo
	ocursor cursorInfo
	omode   uint32
}

// all Windows systems are little endian
var k32 = syscall.NewLazyDLL("kernel32.dll")

// Note that Windows appends some functions with W to indicate that wide
// characters (Unicode) are in use.  The documentation refers to them
// without this suffix, as the resolution is made via preprocessor.
var (
	procReadConsoleInput             = k32.NewProc("ReadConsoleInputW")
	procGetConsoleCursorInfo         = k32.NewProc("GetConsoleCursorInfo")
	procSetConsoleCursorInfo         = k32.NewProc("SetConsoleCursorInfo")
	procSetConsoleCursorPosition     = k32.NewProc("SetConsoleCursorPosition")
	procSetConsoleMode               = k32.NewProc("SetConsoleMode")
	procGetConsoleMode               = k32.NewProc("GetConsoleMode")
	procWriteConsoleOutput           = k32.NewProc("WriteConsoleOutputW")
	procGetConsoleScreenBufferInfo   = k32.NewProc("GetConsoleScreenBufferInfo")
	procFillConsoleOutputAttribute   = k32.NewProc("FillConsoleOutputAttribute")
	procFillConsoleOutputCharacter   = k32.NewProc("FillConsoleOutputCharacterW")
	procSetConsoleWindowInfo         = k32.NewProc("SetConsoleWindowInfo")
	procSetConsoleScreenBufferSize   = k32.NewProc("SetConsoleScreenBufferSize")
	procSetConsoleActiveScreenBuffer = k32.NewProc("SetConsoleActiveScreenBuffer")
)

// We have to bring in the kernel32.dll directly, so we can get access to some
// system calls that the core Go API lacks.

func NewConsoleScreen() (Screen, error) {
	return &cScreen{}, nil
}

func (s *cScreen) Init() error {

	s.evch = make(chan Event, 2)
	s.quit = make(chan struct{})

	if in, e := syscall.Open("CONIN$", syscall.O_RDWR, 0); e != nil {
		return e
	} else {
		s.in = in
	}
	if out, e := syscall.Open("CONOUT$", syscall.O_RDWR, 0); e != nil {
		syscall.Close(s.in)
		return e
	} else {
		s.out = out
	}

	s.getCursorInfo(&s.ocursor)
	s.getConsoleInfo(&s.oscreen)
	s.getMode(&s.omode)

	if err := s.setMode(modeResizeEn); err != nil {
		syscall.Close(s.in)
		syscall.Close(s.out)
		return err
	}
	s.resize()
	s.Clear()
	s.HideCursor()
	//procSetConsoleActiveScreenBuffer.Call(uintptr(s.out))
	go s.scanInput()

	return nil
}

func (s *cScreen) EnableMouse() {
	s.setMode(modeResizeEn | modeMouseEn)
}

func (s *cScreen) DisableMouse() {
	s.setMode(modeResizeEn)
}

func (s *cScreen) Fini() {
	s.setCursorPos(0, 0)
	s.setCursorInfo(&s.ocursor)
	s.setMode(s.omode)
	s.setBufferSize(int(s.oscreen.size.x), int(s.oscreen.size.y))
	s.Clear()

	close(s.quit)
	syscall.Close(s.in)
	syscall.Close(s.out)
}

func (s *cScreen) PostEvent(ev Event) {
	select {
	case <-s.quit:
	case s.evch <- ev:
	}
}

func (s *cScreen) PollEvent() Event {
	select {
	case <-s.quit:
		return nil
	case ev := <-s.evch:
		return ev
	}
}

type cursorInfo struct {
	size    uint32
	visible uint32
}

type coord struct {
	x int16
	y int16
}

func (c coord) uintptr() uintptr {
	// little endian, put x first
	return uintptr(c.x) | (uintptr(c.y) << 16)
}

type rect struct {
	left   int16
	top    int16
	right  int16
	bottom int16
}

func (s *cScreen) ShowCursor(x, y int) {
	var curinfo cursorInfo

	s.getCursorInfo(&curinfo)
	if x < 0 || y < 0 {
		if curinfo.visible == 0 {
			return
		}
		curinfo.visible = 0
		s.setCursorInfo(&curinfo)
	} else {
		curinfo.visible = 1
		s.setCursorPos(x, y)
		s.setCursorInfo(&curinfo)
	}
}

func (c *cScreen) HideCursor() {
	c.ShowCursor(10, 5)
}

type charInfo struct {
	ch   uint16
	attr uint16
}

type inputRecord struct {
	typ  uint16
	_    uint16
	data [16]byte
}

const (
	keyEvent    uint16 = 1
	mouseEvent  uint16 = 2
	resizeEvent uint16 = 4
	menuEvent   uint16 = 8  // don't use
	focusEvent  uint16 = 16 // don't use
)

type mouseRecord struct {
	x     int16
	y     int16
	btns  uint32
	mod   uint32
	flags uint32
}

type resizeRecord struct {
	x int16
	y int16
}

type keyRecord struct {
	isdown int32
	repeat uint16
	kcode  uint16
	scode  uint16
	ch     uint16
	mod    uint32
}

const (
	// Constants per Microsoft.  We don't put the modifiers
	// here.
	vkCancel = 0x03
	vkBack   = 0x08 // Backspace
	vkTab    = 0x09
	vkClear  = 0x0c
	vkReturn = 0x0d
	vkEscape = 0x1b
	vkSpace  = 0x20
	vkPrior  = 0x21 // PgUp
	vkNext   = 0x22 // PgDn
	vkEnd    = 0x23
	vkHome   = 0x24
	vkLeft   = 0x25
	vkUp     = 0x26
	vkRight  = 0x27
	vkDown   = 0x28
	vkInsert = 0x2d
	vkDelete = 0x2e
	vkHelp   = 0x2f
	vkF1     = 0x70
	vkF2     = 0x71
	vkF3     = 0x72
	vkF4     = 0x73
	vkF5     = 0x74
	vkF6     = 0x75
	vkF7     = 0x76
	vkF8     = 0x77
	vkF9     = 0x78
	vkF10    = 0x79
	vkF11    = 0x7a
	vkF12    = 0x7b
	vkF13    = 0x7c
	vkF14    = 0x7d
	vkF15    = 0x7e
	vkF16    = 0x7f
	vkF17    = 0x80
	vkF18    = 0x81
	vkF19    = 0x82
	vkF20    = 0x83
	vkF21    = 0x84
	vkF22    = 0x85
	vkF23    = 0x86
	vkF24    = 0x87
)

// NB: All Windows platforms are little endian.  We assume this
// never, ever change.  The following code is endian safe. and does
// not use unsafe pointers.
func getu32(v []byte) uint32 {
	return uint32(v[0]) + (uint32(v[1]) << 8) + (uint32(v[2]) << 16) + (uint32(v[3]) << 24)
}
func geti32(v []byte) int32 {
	return int32(getu32(v))
}
func getu16(v []byte) uint16 {
	return uint16(v[0]) + (uint16(v[1]) << 8)
}
func geti16(v []byte) int16 {
	return int16(getu16(v))
}

// Convert windows dwControlKeyState to modifier mask
func mod2mask(cks uint32) ModMask {
	mm := ModNone
	// Left or right control
	if (cks & (0x0008 | 0x0004)) != 0 {
		mm |= ModCtrl
	}
	// Left or right alt
	if (cks & (0x0002 | 0x0001)) != 0 {
		mm |= ModAlt
	}
	// Any shift
	if (cks & 0x0010) != 0 {
		mm |= ModShift
	}
	return mm
}

func (s *cScreen) getConsoleInput() error {
	rec := &inputRecord{}
	var nrec int32
	rv, _, er := procReadConsoleInput.Call(
		uintptr(s.in),
		uintptr(unsafe.Pointer(rec)),
		uintptr(1),
		uintptr(unsafe.Pointer(&nrec)))
	if rv == 0 {
		return er
	}
	if nrec != 1 {
		return nil
	}
	switch rec.typ {
	case keyEvent:
		krec := &keyRecord{}
		krec.isdown = geti32(rec.data[0:])
		krec.repeat = getu16(rec.data[4:])
		krec.kcode = getu16(rec.data[6:])
		krec.scode = getu16(rec.data[8:])
		krec.ch = getu16(rec.data[10:])
		krec.mod = getu32(rec.data[12:])

		if krec.isdown == 0 || krec.repeat < 1 {
			// its a key release event, ignore it
			return nil
		}
		if krec.ch != 0 {
			// synthesized key code
			for krec.repeat > 0 {
				s.PostEvent(NewEventKey(KeyRune, rune(krec.ch), mod2mask(krec.mod)))
				krec.repeat--
			}
			return nil
		}
		key := KeyNUL // impossible on Windows
		switch krec.kcode {
		case vkBack:
			key = KeyBackspace
		case vkTab:
			key = KeyTab
		case vkPrior:
			key = KeyPgUp
		case vkNext:
			key = KeyPgDn
		case vkReturn:
			key = KeyEnter
		case vkEnd:
			key = KeyEnd
		case vkHome:
			key = KeyHome
		case vkLeft:
			key = KeyLeft
		case vkUp:
			key = KeyUp
		case vkRight:
			key = KeyRight
		case vkDown:
			key = KeyDown
		case vkInsert:
			key = KeyInsert
		case vkDelete:
			key = KeyDelete
		case vkHelp:
			key = KeyHelp
		case vkF1:
			key = KeyF1
		case vkF2:
			key = KeyF2
		case vkF3:
			key = KeyF3
		case vkF4:
			key = KeyF4
		case vkF5:
			key = KeyF5
		case vkF6:
			key = KeyF6
		case vkF7:
			key = KeyF7
		case vkF8:
			key = KeyF8
		case vkF9:
			key = KeyF9
		case vkF10:
			key = KeyF10
		case vkF11:
			key = KeyF11
		case vkF12:
			key = KeyF12
		case vkF13:
			key = KeyF13
		case vkF14:
			key = KeyF14
		case vkF15:
			key = KeyF15
		case vkF16:
			key = KeyF16
		case vkF17:
			key = KeyF17
		case vkF18:
			key = KeyF18
		case vkF19:
			key = KeyF19
		case vkF20:
			key = KeyF20
		case vkF21:
			key = KeyF21
		case vkF22:
			key = KeyF22
		case vkF23:
			key = KeyF23
		case vkF24:
			key = KeyF24
		default:
			return nil
		}
		for krec.repeat > 0 {
			s.PostEvent(NewEventKey(key, rune(krec.ch), mod2mask(krec.mod)))
			krec.repeat--
		}

	case mouseEvent:
		var mrec mouseRecord
		mrec.x = geti16(rec.data[0:])
		mrec.y = geti16(rec.data[2:])
		mrec.btns = getu32(rec.data[4:])
		mrec.mod = getu32(rec.data[8:])
		mrec.flags = getu32(rec.data[12:]) // not using yet
		btns := ButtonNone

		if mrec.btns == s.mbtns {
			// If the buttons have not changed,
			// then don't report the event.  We aren't
			// reporting motion events at this time.
			return nil
		}

		s.mbtns = mrec.btns
		if mrec.btns&0x1 != 0 {
			btns |= Button1
		}
		if mrec.btns&0x2 != 0 {
			btns |= Button2
		}
		if mrec.btns&0x4 != 0 {
			btns |= Button3
		}
		if mrec.btns&0x8 != 0 {
			btns |= Button4
		}
		if mrec.btns&0x10 != 0 {
			btns |= Button5
		}

		s.PostEvent(NewEventMouse(int(mrec.x), int(mrec.y), btns, mod2mask(mrec.mod)))

	case resizeEvent:
		var rrec resizeRecord
		rrec.x = geti16(rec.data[0:])
		rrec.y = geti16(rec.data[2:])
		s.PostEvent(NewEventResize(int(rrec.x), int(rrec.y)))

	default:
	}
	return nil
}

func (s *cScreen) scanInput() {
	for {
		if e := s.getConsoleInput(); e != nil {
			return
		}
	}
}

// Windows console can display 8 characters, in either low or high intensity
func (s *cScreen) Colors() int {
	return 16
}

// Windows uses RGB signals
func mapColor2RGB(c Color) uint16 {
	switch c {
	case ColorBlack:
		return 0
		// primaries
	case ColorRed:
		return 0x4
	case ColorGreen:
		return 0x2
	case ColorBlue:
		return 0x1
	case ColorYellow:
		return 0x6
	case ColorMagenta:
		return 0x5
	case ColorCyan:
		return 0x3
	case ColorWhite:
		return 0x7
	// bright variants
	case ColorGrey:
		return 0x8
	case ColorBrightRed:
		return 0xc
	case ColorBrightGreen:
		return 0xa
	case ColorBrightBlue:
		return 0x9
	case ColorBrightYellow:
		return 0xe
	case ColorBrightMagenta:
		return 0xd
	case ColorBrightCyan:
		return 0xb
	case ColorBrightWhite:
		return 0xf
	}
	return 0
}

// Map a tcell style to Windows attributes
func mapStyle(style Style) uint16 {
	f, b, a := style.Decompose()
	if f == ColorDefault {
		f = ColorWhite
	}
	if b == ColorDefault {
		b = ColorBlack
	}
	attr := mapColor2RGB(f)
	attr |= (mapColor2RGB(b) << 4)
	if a&AttrBold != 0 {
		attr |= 0x8
	}
	if a&AttrDim != 0 {
		attr &^= 0x8
	}
	if a&AttrUnderline != 0 {
		attr |= 0x8000
	}
	if a&AttrReverse != 0 {
		attr |= 0x4000
	}
	// Blink is unsupported
	return attr
}

func (s *cScreen) SetCell(x, y int, style Style, ch ...rune) {
	r := ' '
	w := 0
	for i := range ch {
		if w = runewidth.RuneWidth(ch[i]); w != 0 {
			r = ch[i]
			break
		}
	}

	// Windows console lacks support for combining chars
	if w == 0 {
		r = ' '
		w = 1
	}

	rec := rect{left: int16(x), right: int16(x), top: int16(y), bottom: int16(y)}
	pos := coord{x: int16(0), y: int16(0)}
	siz := coord{x: 1, y: 1}
	dat := charInfo{ch: uint16(r), attr: mapStyle(style)}

	procWriteConsoleOutput.Call(
		uintptr(s.out),
		uintptr(unsafe.Pointer(&dat)),
		siz.uintptr(),
		pos.uintptr(),
		uintptr(unsafe.Pointer(&rec)))
}

type consoleInfo struct {
	size  coord
	pos   coord
	attrs uint16
	win   rect
	maxsz coord
}

func (s *cScreen) getConsoleInfo(info *consoleInfo) {
	procGetConsoleScreenBufferInfo.Call(
		uintptr(s.out),
		uintptr(unsafe.Pointer(info)))
}

func (s *cScreen) getCursorInfo(info *cursorInfo) {
	procGetConsoleCursorInfo.Call(
		uintptr(s.out),
		uintptr(unsafe.Pointer(info)))
}

func (s *cScreen) setCursorInfo(info *cursorInfo) {
	procSetConsoleCursorInfo.Call(
		uintptr(s.out),
		uintptr(unsafe.Pointer(info)))
}

func (s *cScreen) setCursorPos(x, y int) {
	procSetConsoleCursorPosition.Call(
		uintptr(s.out),
		coord{int16(x), int16(y)}.uintptr())
}

func (s *cScreen) setBufferSize(x, y int) {
	procSetConsoleScreenBufferSize.Call(
		uintptr(s.out),
		coord{int16(x), int16(y)}.uintptr())
}

func (s *cScreen) Size() (int, int) {

	info := consoleInfo{}
	s.getConsoleInfo(&info)
	w := int((info.win.right - info.win.left) + 1)
	h := int((info.win.bottom - info.win.top) + 1)

	return w, h
}

func (s *cScreen) resize() {

	info := consoleInfo{}
	s.getConsoleInfo(&info)

	w := int((info.win.right - info.win.left) + 1)
	h := int((info.win.bottom - info.win.top) + 1)
	if s.w == w && s.h == h {
		return
	}
	r := rect{0, 0, int16(w - 1), int16(h - 1)}
	procSetConsoleWindowInfo.Call(
		uintptr(s.out),
		uintptr(1),
		uintptr(unsafe.Pointer(&r)))

	s.setBufferSize(w, h)
}

func (s *cScreen) Clear() {
	pos := coord{0, 0}
	attr := uint16(0x7) // default white fg, black bg)
	x, y := s.Size()
	scratch := uint32(0)
	count := uint32(x * y)

	procFillConsoleOutputAttribute.Call(
		uintptr(s.out),
		uintptr(attr),
		uintptr(count),
		pos.uintptr(),
		uintptr(unsafe.Pointer(&scratch)))
	procFillConsoleOutputCharacter.Call(
		uintptr(s.out),
		uintptr(' '),
		uintptr(count),
		pos.uintptr(),
		uintptr(unsafe.Pointer(&scratch)))
}

const (
	modeMouseEn  uint32 = 0x0010
	modeResizeEn uint32 = 0x0008
)

func (s *cScreen) setMode(mode uint32) error {
	rv, _, err := procSetConsoleMode.Call(
		uintptr(s.in),
		uintptr(mode))
	if rv == 0 {
		return err
	}
	return nil
}

func (s *cScreen) getMode(v *uint32) {
	procGetConsoleMode.Call(
		uintptr(s.in),
		uintptr(unsafe.Pointer(v)))
}
