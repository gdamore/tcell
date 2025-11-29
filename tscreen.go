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

//go:build !(js && wasm)
// +build !js !wasm

package tcell

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"golang.org/x/text/transform"
)

// NewTerminfoScreen returns a Screen that uses the stock TTY interface
// and POSIX terminal control, combined with a terminfo description taken from
// the $TERM environment variable.  It returns an error if the terminal
// is not supported for any reason.
//
// For terminals that do not support dynamic resize events, the $LINES
// $COLUMNS environment variables can be set to the actual window size,
// otherwise defaults taken from the terminal database are used.
func NewTerminfoScreen() (Screen, error) {
	return NewTerminfoScreenFromTty(nil)
}

// Some terminal escapes that are basically universal.
// We would really like to be able to use private mode queries for some of
// these but generally we've found that support for queries is not always present,
// even when the private modes can be controlled. It appears that *all* terminals
// will happily swallow the escapes that they do not recognize, with the small annoyance
// in "st" where it prints error messages to its stderr (which is usually not visible
// to the user unless they started it from another terminal session).  But apart from
// the complaint to stderr from "st", everything else is fine.
const (
	enableAutoMargin  = "\x1b[?7h" // dec private mode 7
	disableAutoMargin = "\x1b[?7l"
	setCursorPosition = "\x1b[%[1]d;%[2]dH"
	sgr0              = "\x1b[m" // attrOff
	bold              = "\x1b[1m"
	dim               = "\x1b[2m"
	italic            = "\x1b[3m"
	underline         = "\x1b[4m"
	blink             = "\x1b[5m"
	reverse           = "\x1b[7m"
	strikeThrough     = "\x1b[8m"
	showCursor        = "\x1b[?25h"
	hideCursor        = "\x1b[?25l"
	clear             = "\x1b[H\x1b[J" // NB: sun uses \f
	enablePaste       = "\x1b[?2004h"
	disablePaste      = "\x1b[?2004l"
	enableFocus       = "\x1b[?1004h"
	disableFocus      = "\x1b[?1004l"
	startSyncOut      = "\x1b[?2026h"
	endSyncOut        = "\x1b[?2026l"
	doubleUnder       = "\x1b[4:2m"
	curlyUnder        = "\x1b[4:3m"
	dottedUnder       = "\x1b[4:4m"
	dashedUnder       = "\x1b[4:5m"
	underColor        = "\x1b[58:5:%dm"
	underRGB          = "\x1b[58:2::%d:%d:%dm"
	underFg           = "\x1b[59m"
	enableAltChars    = "\x1b(B\x1b)0"                       // set G0 as US-ASCII, G1 as DEC line drawing
	startAltChars     = "\x0e"                               // aka Shift-Out
	endAltChars       = "\x0f"                               // aka Shift-In
	setFg8            = "\x1b[3%dm"                          // for colors less than 8
	setFg256          = "\x1b[38;5;%dm"                      // for colors less than 256
	setFgRgb          = "\x1b[38;2;%d;%d;%dm"                // for RGB
	setBg8            = "\x1b[4%dm"                          // color colors less than 8
	setBg256          = "\x1b[48;5;%dm"                      // for colors less than 256
	setBgRgb          = "\x1b[48;2;%d;%d;%dm"                // for RGB
	setFgBgRgb        = "\x1b[38;2;%d;%d;%d;48;2;%d;%d;%dm"  // for RGB, in one shot
	resetFgBg         = "\x1b[39;49m"                        // ECMA defined
	enterCA           = "\x1b[?1049h"                        // alternate screen
	exitCA            = "\x1b[?1049l"                        // alternate screen
	enterKeypad       = "\x1b[?1h\x1b="                      // Note mode 1 might not be supported everywhere
	exitKeypad        = "\x1b[?1l\x1b>"                      // Also mode 1
)

// NewTerminfoScreenFromTty returns a Screen using a custom Tty implementation.
// If the passed in tty is nil, then a reasonable default (typically /dev/tty)
// is presumed, at least on UNIX hosts. (Windows hosts will typically fail this
// call altogether.)
func NewTerminfoScreenFromTty(tty Tty) (Screen, error) {
	t := &tScreen{tty: tty}

	t.prepareCursorStyles()
	t.prepareExtendedOSC()
	t.buildAcsMap()
	t.resizeQ = make(chan bool, 1)
	t.fallback = make(map[rune]string)
	maps.Copy(t.fallback, RuneFallbacks)

	return &baseScreen{screenImpl: t}, nil
}

// tScreen represents a screen backed by a terminfo implementation.
type tScreen struct {
	tty          Tty
	h            int
	w            int
	fini         bool
	cells        CellBuffer
	buffering    bool // true if we are collecting writes to buf instead of sending directly to out
	buf          bytes.Buffer
	curstyle     Style
	style        Style
	resizeQ      chan bool
	quit         chan struct{}
	keychan      chan []byte
	cx           int
	cy           int
	clear        bool
	cursorx      int
	cursory      int
	acs          map[rune]string
	charset      string
	encoder      transform.Transformer
	decoder      transform.Transformer
	fallback     map[rune]string
	ncolor       int
	colors       map[Color]Color
	palette      []Color
	truecolor    bool
	legacy       bool
	finiOnce     sync.Once
	enterUrl     string
	exitUrl      string
	setWinSize   string
	cursorStyles map[CursorStyle]string
	cursorStyle  CursorStyle
	cursorColor  Color
	cursorRGB    string
	cursorFg     string
	stopQ        chan struct{}
	eventQ       chan Event
	running      bool
	wg           sync.WaitGroup
	mouseFlags   MouseFlags
	pasteEnabled bool
	focusEnabled bool
	setTitle     string
	saveTitle    string
	restoreTitle string
	title        string
	setClipboard string
	enableCsiU   string
	disableCsiU  string
	input        *inputProcessor

	sync.Mutex
}

func (t *tScreen) Init() error {
	if e := t.initialize(); e != nil {
		return e
	}

	t.keychan = make(chan []byte, 10)

	t.charset = getCharset()
	if enc := GetEncoding(t.charset); enc != nil {
		t.encoder = enc.NewEncoder()
		t.decoder = enc.NewDecoder()
	} else {
		return ErrNoCharset
	}

	// environment overrides
	w := 80
	h := 24
	if i, _ := strconv.Atoi(os.Getenv("LINES")); i != 0 {
		h = i
	}
	if i, _ := strconv.Atoi(os.Getenv("COLUMNS")); i != 0 {
		w = i
	}
	cterm := os.Getenv("COLORTERM")
	nterm := os.Getenv("TERM")

	if t.ncolor == 0 {
		if slices.Contains([]string{"truecolor", "direct", "24bit"}, cterm) || strings.HasSuffix(nterm, "-direct") || strings.HasSuffix(nterm, "-truecolor") {
			t.truecolor = true
			t.ncolor = 256 // base 8-bit palette
		} else if strings.HasSuffix("-256color", nterm) || strings.Contains(cterm, "256") {
			t.ncolor = 256
		} else if strings.HasSuffix("-88color", nterm) {
			t.ncolor = 88
		} else if strings.Contains(nterm, "color") || cterm != "" {
			t.ncolor = 8
		} else if strings.Contains(nterm, "mono") || strings.HasSuffix(nterm, "-m") { // monochrome variants
			t.ncolor = 0
		} else if strings.Contains(nterm, "ansi") || slices.Contains([]string{"dtterm", "xterm", "aixterm", "linux"}, nterm) {
			t.ncolor = 8
		} else if strings.HasPrefix(nterm, "vt") || nterm == "sun" {
			// legacy DEC VT 100/220 etc. family.  (technically the VT525 can do ANSI, but they should set to ansi)
			t.ncolor = 0
		} else {
			// best guess - this covers all the modern variants like ghostty,
			t.ncolor = 256
		}
	}

	if os.Getenv("NO_COLOR") != "" {
		t.truecolor = false
		t.ncolor = 0
	}

	if strings.HasPrefix(nterm, "vt") || strings.Contains(nterm, "ansi") || nterm == "linux" || nterm == "sun" || nterm == "sun-color" {
		// these terminals are "legacy" and not expected to support most OSC functions
		t.legacy = true
	}

	// A user who wants to have his themes honored can
	// set this environment variable.
	if os.Getenv("TCELL_TRUECOLOR") == "disable" {
		t.truecolor = false
	}
	// clip to reasonable limits
	nColors := min(t.ncolor, 256)
	t.colors = make(map[Color]Color, nColors)
	t.palette = make([]Color, nColors)
	for i := range nColors {
		t.palette[i] = Color(i) | ColorValid
		// identity map for our builtin colors
		t.colors[Color(i)|ColorValid] = Color(i) | ColorValid
	}

	t.quit = make(chan struct{})
	t.eventQ = make(chan Event, 256)
	t.input = newInputProcessor(t.eventQ)

	t.Lock()
	t.cx = -1
	t.cy = -1
	t.style = StyleDefault
	t.cells.Resize(w, h)
	t.cursorx = -1
	t.cursory = -1
	t.resize()
	t.Unlock()

	if err := t.engage(); err != nil {
		return err
	}

	return nil
}

func (t *tScreen) prepareExtendedOSC() {
	if t.legacy {
		return
	}

	// OSC 8 is for enter/exit URL.
	t.enterUrl = "\x1b]8;%[2]s;%[1]s\x1b\\"
	t.exitUrl = "\x1b]8;;\x1b\\"

	// CSI .. t is for window operations.
	t.setWinSize = "\x1b[8;%[2]d;%[1]dt"
	t.saveTitle = "\x1b[22;2t"
	t.restoreTitle = "\x1b[23;2t"
	// this also tries to request that UTF-8 is allowed in the title
	t.setTitle = "\x1b[>2t\x1b]2;%s\x1b\\"

	// OSC 52 is for saving to the clipboard.
	// this string takes a base64 string and sends it to the clipboard.
	// it will also be able to retrieve the clipboard using "?" as the
	// sent string, when we support that.
	t.setClipboard = "\x1b]52;c;%s\x1b\\"

	if t.enableCsiU == "" {
		if runtime.GOOS == "windows" && (os.Getenv("TERM") == "" || os.Getenv("TERM_PROGRAM") == "WezTerm") {
			// on Windows, if we don't have a TERM, use only win32-input-mode
			t.enableCsiU = "\x1b[?9001h"
			t.disableCsiU = "\x1b[?9001l"
		} else if os.Getenv("TERM_PROGRAM") == "WezTerm" {
			// WezTerm is unhappy if we ask for other modes
			t.enableCsiU = "\x1b[>1u"
			t.disableCsiU = "\x1b[<u"
		} else {
			// three advanced keyboard protocols:
			// - xterm modifyOtherKeys (uses CSI 27 ~ )
			// - kitty csi-u (uses CSI u)
			// - win32-input-mode (uses CSI _)
			t.enableCsiU = "\x1b[>4;2m" + "\x1b[>1u" + "\x1b[9001h"
			t.disableCsiU = "\x1b[9001l" + "\x1b[<u" + "\x1b[>4;0m"
		}
	}
}

func (t *tScreen) prepareCursorStyles() {
	if t.legacy {
		return
	}
	t.cursorStyles = map[CursorStyle]string{
		CursorStyleDefault:           "\x1b[0 q",
		CursorStyleBlinkingBlock:     "\x1b[1 q",
		CursorStyleSteadyBlock:       "\x1b[2 q",
		CursorStyleBlinkingUnderline: "\x1b[3 q",
		CursorStyleSteadyUnderline:   "\x1b[4 q",
		CursorStyleBlinkingBar:       "\x1b[5 q",
		CursorStyleSteadyBar:         "\x1b[6 q",
	}
	if t.cursorRGB == "" {
		t.cursorRGB = "\x1b]12;#%02x%02x%02x\007"
		t.cursorFg = "\x1b]112\007"
	}
}

func (t *tScreen) Fini() {
	t.finiOnce.Do(t.finish)
}

func (t *tScreen) finish() {
	close(t.quit)
	t.finalize()
}

func (t *tScreen) SetStyle(style Style) {
	t.Lock()
	if !t.fini {
		t.style = style
	}
	t.Unlock()
}

func (t *tScreen) encodeStr(s string) []byte {

	var dstBuf [128]byte
	var buf []byte
	nb := dstBuf[:]
	dst := 0
	var err error
	if enc := t.encoder; enc != nil {
		enc.Reset()
		dst, _, err = enc.Transform(nb, []byte(s), true)
	}
	if err != nil || dst == 0 || nb[0] == '\x1a' {
		// Combining characters are elided
		r, _ := utf8.DecodeRuneInString(s)
		if len(buf) == 0 {
			if acs, ok := t.acs[r]; ok {
				buf = append(buf, []byte(acs)...)
			} else if fb, ok := t.fallback[r]; ok {
				buf = append(buf, []byte(fb)...)
			} else {
				buf = append(buf, '?')
			}
		}
	} else {
		buf = append(buf, nb[:dst]...)
	}

	return buf
}

func (t *tScreen) sendFgBg(fg Color, bg Color, attr AttrMask) AttrMask {
	if t.Colors() == 0 {
		// foreground vs background, we calculate luminance
		// and possibly do a reverse video
		if !fg.Valid() {
			return attr
		}
		v, ok := t.colors[fg]
		if !ok {
			v = FindColor(fg, []Color{ColorBlack, ColorWhite})
			t.colors[fg] = v
		}
		switch v {
		case ColorWhite:
			return attr
		case ColorBlack:
			return attr ^ AttrReverse
		}
	}

	if fg == ColorReset || bg == ColorReset {
		t.TPuts(resetFgBg)
	}
	if t.truecolor {
		if fg.IsRGB() && bg.IsRGB() {
			r1, g1, b1 := fg.RGB()
			r2, g2, b2 := bg.RGB()
			t.TPuts(fmt.Sprintf(setFgBgRgb, r1, g1, b1, r2, g2, b2))
			return attr
		}

		if fg.IsRGB() {
			r, g, b := fg.RGB()
			t.TPuts(fmt.Sprintf(setFgRgb, r, g, b))
			fg = ColorDefault
		}

		if bg.IsRGB() {
			r, g, b := bg.RGB()
			t.TPuts(fmt.Sprintf(setBgRgb, r, g, b))
			bg = ColorDefault
		}
	}

	if fg.Valid() {
		if v, ok := t.colors[fg]; ok {
			fg = v
		} else {
			v = FindColor(fg, t.palette)
			t.colors[fg] = v
			fg = v
		}
		fgc := fg & 0xffffff
		if fgc < 8 {
			t.TPuts(fmt.Sprintf(setFg8, fgc))
		} else if fgc < 256 {
			t.TPuts(fmt.Sprintf(setFg256, fgc))
		}
	}

	if bg.Valid() {
		if v, ok := t.colors[bg]; ok {
			bg = v
		} else {
			v = FindColor(bg, t.palette)
			t.colors[bg] = v
			bg = v
		}
		bgc := bg & 0xffffff
		if bgc < 8 {
			t.TPuts(fmt.Sprintf(setBg8, bgc))
		} else if bgc < 256 {
			t.TPuts(fmt.Sprintf(setBg256, bgc))
		}

	}

	return attr
}

func (t *tScreen) drawCell(x, y int) int {

	str, style, width := t.cells.Get(x, y)
	if !t.cells.Dirty(x, y) {
		return width
	}

	if t.cy != y || t.cx != x {
		t.TPuts(fmt.Sprintf(setCursorPosition, y+1, x+1))
		t.cx = x
		t.cy = y
	}

	if style == StyleDefault {
		style = t.style
	}
	if style != t.curstyle {
		fg, bg, attrs := style.fg, style.bg, style.attrs

		t.TPuts(sgr0)

		attrs = t.sendFgBg(fg, bg, attrs)
		if attrs&AttrBold != 0 {
			t.TPuts(bold)
		}
		if us, uc := style.ulStyle, style.ulColor; us != UnderlineStyleNone {
			if uc == ColorReset {
				t.TPuts(underFg)
			} else if uc.IsRGB() {
				if v, ok := t.colors[uc]; ok {
					uc = v
				} else {
					v = FindColor(uc, t.palette)
					t.colors[uc] = v
					uc = v
				}
				t.TPuts(fmt.Sprintf(underColor, uc&0xff))
				r, g, b := uc.RGB()
				t.TPuts(fmt.Sprintf(underRGB, r, g, b))
			} else if uc.Valid() {
				t.TPuts(fmt.Sprintf(underColor, uc&0xff))
			}

			t.TPuts(underline) // to ensure everyone gets at least a basic underline
			switch us {
			case UnderlineStyleDouble:
				t.TPuts(doubleUnder)
			case UnderlineStyleCurly:
				t.TPuts(curlyUnder)
			case UnderlineStyleDotted:
				t.TPuts(dottedUnder)
			case UnderlineStyleDashed:
				t.TPuts(dashedUnder)
			}
		}
		if attrs&AttrReverse != 0 {
			t.TPuts(reverse)
		}
		if attrs&AttrBlink != 0 {
			t.TPuts(blink)
		}
		if attrs&AttrDim != 0 {
			t.TPuts(dim)
		}
		if attrs&AttrItalic != 0 {
			t.TPuts(italic)
		}
		if attrs&AttrStrikeThrough != 0 {
			t.TPuts(strikeThrough)
		}

		var newUrl urlInfo
		var oldUrl urlInfo
		if t.curstyle.url != nil {
			oldUrl = *t.curstyle.url
		}
		if style.url != nil {
			newUrl = *style.url
		}
		// URL string can be long, so don't send it unless we really need to
		if t.enterUrl != "" && newUrl != oldUrl {
			if newUrl.url != "" {
				t.TPuts(fmt.Sprintf(t.enterUrl, newUrl.url, newUrl.id))
			} else {
				t.TPuts(t.exitUrl)
			}
		}

		t.curstyle = style
	}

	// now emit runes - taking care to not overrun width with a
	// wide character, and to ensure that we emit exactly one regular
	// character followed up by any residual combing characters

	if width < 1 {
		width = 1
	}

	buf := t.encodeStr(str)
	str = string(buf)

	if width > 1 && str == "?" {
		// No FullWidth character support
		str = "? "
		t.cx = -1
	}

	if x > t.w-width {
		// too wide to fit; emit a single space instead
		width = 1
		str = " "
	}
	t.writeString(str)
	t.cx += width
	t.cells.SetDirty(x, y, false)
	if width > 1 {
		t.cx = -1
	}

	return width
}

func (t *tScreen) ShowCursor(x, y int) {
	t.Lock()
	t.cursorx = x
	t.cursory = y
	t.Unlock()
}

func (t *tScreen) SetCursor(cs CursorStyle, cc Color) {
	t.Lock()
	t.cursorStyle = cs
	t.cursorColor = cc
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
	t.TPuts(fmt.Sprintf(setCursorPosition, y+1, x+1))
	t.TPuts(showCursor)
	if t.cursorStyles != nil {
		if esc, ok := t.cursorStyles[t.cursorStyle]; ok {
			t.TPuts(esc)
		}
	}
	if t.cursorRGB != "" {
		if t.cursorColor == ColorReset {
			t.TPuts(t.cursorFg)
		} else if t.cursorColor.Valid() {
			r, g, b := t.cursorColor.RGB()
			t.TPuts(fmt.Sprintf(t.cursorRGB, r, g, b))
		}
	}
	t.cx = x
	t.cy = y
}

// writeString sends a string to the terminal. The string is sent as-is and
// this function does not expand inline padding indications (of the form
// $<[delay]> where [delay] is msec). In order to have these expanded, use
// TPuts. If the screen is "buffering", the string is collected in a buffer,
// with the intention that the entire buffer be sent to the terminal in one
// write operation at some point later.
func (t *tScreen) writeString(s string) {
	if t.buffering {
		_, _ = io.WriteString(&t.buf, s)
	} else {
		_, _ = io.WriteString(t.tty, s)
	}
}

func (t *tScreen) TPuts(s string) {
	t.writeString(s)
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
	t.TPuts(sgr0)
	t.TPuts(t.exitUrl)
	_ = t.sendFgBg(t.style.fg, t.style.bg, AttrNone)

	t.TPuts(clear)

	t.clear = false
}

func (t *tScreen) startBuffering() {
	t.TPuts(startSyncOut)
}

func (t *tScreen) endBuffering() {
	t.TPuts(endSyncOut)
}

func (t *tScreen) hideCursor() {
	// just in case we cannot hide it, move it to the end
	t.cx, t.cy = t.cells.Size()
	t.TPuts(fmt.Sprintf(setCursorPosition, t.cy+1, t.cx+1))
	// then hide it
	t.TPuts(hideCursor)
}

func (t *tScreen) draw() {
	// clobber cursor position, because we're going to change it all
	t.cx = -1
	t.cy = -1
	// make no style assumptions
	t.curstyle = styleInvalid

	t.buf.Reset()
	t.buffering = true
	t.startBuffering()
	defer func() {
		t.buffering = false
		t.endBuffering()
	}()

	// hide the cursor while we move stuff around
	t.hideCursor()

	if t.clear {
		t.clearScreen()
	}

	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			width := t.drawCell(x, y)
			if width > 1 {
				if x+1 < t.w {
					// this is necessary so that if we ever
					// go back to drawing that cell, we
					// actually will *draw* it.
					t.cells.SetDirty(x+1, y, true)
				}
			}
			x += width - 1
		}
	}

	// restore the cursor
	t.showCursor()

	_, _ = t.buf.WriteTo(t.tty)
}

func (t *tScreen) EnableMouse(flags ...MouseFlags) {
	var f MouseFlags
	flagsPresent := false
	for _, flag := range flags {
		f |= flag
		flagsPresent = true
	}
	if !flagsPresent {
		f = MouseMotionEvents | MouseDragEvents | MouseButtonEvents
	}

	t.Lock()
	t.mouseFlags = f
	t.enableMouse(f)
	t.Unlock()
}

func (t *tScreen) enableMouse(f MouseFlags) {
	// Rather than using terminfo to find mouse escape sequences, we rely on the fact that
	// pretty much *every* terminal that supports mouse tracking follows the
	// XTerm standards (the modern ones).  It is expected that all terminals understand
	// the same DEC private modes.  Note that the SGR mode is required for the mouse sequences
	// to be understood.

	// start by disabling all tracking.
	t.TPuts("\x1b[?1000l\x1b[?1002l\x1b[?1003l\x1b[?1006l")
	if f&MouseButtonEvents != 0 {
		t.TPuts("\x1b[?1000h")
	}
	if f&MouseDragEvents != 0 {
		t.TPuts("\x1b[?1002h")
	}
	if f&MouseMotionEvents != 0 {
		t.TPuts("\x1b[?1003h")
	}
	if f&(MouseButtonEvents|MouseDragEvents|MouseMotionEvents) != 0 {
		t.TPuts("\x1b[?1006h")
	}
}

func (t *tScreen) DisableMouse() {
	t.Lock()
	t.mouseFlags = 0
	t.enableMouse(0)
	t.Unlock()
}

func (t *tScreen) EnablePaste() {
	t.Lock()
	t.pasteEnabled = true
	t.enablePasting(true)
	t.Unlock()
}

func (t *tScreen) DisablePaste() {
	t.Lock()
	t.pasteEnabled = false
	t.enablePasting(false)
	t.Unlock()
}

func (t *tScreen) enablePasting(on bool) {
	var s string
	if on {
		s = enablePaste
	} else {
		s = disablePaste
	}
	if s != "" {
		t.TPuts(s)
	}
}

func (t *tScreen) EnableFocus() {
	t.Lock()
	t.focusEnabled = true
	t.enableFocusReporting()
	t.Unlock()
}

func (t *tScreen) DisableFocus() {
	t.Lock()
	t.focusEnabled = false
	t.disableFocusReporting()
	t.Unlock()
}

func (t *tScreen) enableFocusReporting() {
	t.TPuts(enableFocus)
}

func (t *tScreen) disableFocusReporting() {
	t.TPuts(disableFocus)
}

func (t *tScreen) Size() (int, int) {
	t.Lock()
	w, h := t.w, t.h
	t.Unlock()
	return w, h
}

func (t *tScreen) resize() {
	ws, err := t.tty.WindowSize()
	if err != nil {
		return
	}
	if ws.Width == t.w && ws.Height == t.h {
		return
	}
	t.cx = -1
	t.cy = -1

	t.cells.Resize(ws.Width, ws.Height)
	t.cells.Invalidate()
	t.h = ws.Height
	t.w = ws.Width
	t.input.SetSize(ws.Width, ws.Height)
}

func (t *tScreen) Colors() int {
	// this doesn't change, no need for lock
	if t.truecolor {
		return 1 << 24
	}
	return t.ncolor
}

// nColors returns the size of the built-in palette.
// This is distinct from Colors(), as it will generally
// always be a small number. (<= 256)
func (t *tScreen) nColors() int {
	return t.ncolor
}

// vtACSNames is a map of bytes defined by terminfo that are used in
// the terminals Alternate Character Set to represent other glyphs.
// For example, the upper left corner of the box drawing set can be
// displayed by printing "l" while in the alternate character set.
// It's not quite that simple, since the "l" is the terminfo name,
// and it may be necessary to use a different character based on
// the terminal implementation (or the terminal may lack support for
// this altogether).  These values are from the DEC VT100, and all
// modern terminal emulators support this as charset 0.
var vtACSNames = map[byte]rune{
	'`': RuneDiamond,
	'a': RuneCkBoard,
	'f': RuneDegree,
	'g': RunePlMinus,
	'h': RuneBoard,
	'i': RuneLantern,
	'j': RuneLRCorner,
	'k': RuneURCorner,
	'l': RuneULCorner,
	'm': RuneLLCorner,
	'n': RunePlus,
	'o': RuneS1,
	'p': RuneS3,
	'q': RuneHLine,
	'r': RuneS7,
	's': RuneS9,
	't': RuneLTee,
	'u': RuneRTee,
	'v': RuneBTee,
	'w': RuneTTee,
	'x': RuneVLine,
	'y': RuneLEqual,
	'z': RuneGEqual,
	'{': RunePi,
	'|': RuneNEqual,
	'}': RuneSterling,
	'~': RuneBullet,
}

// buildAcsMap builds a map of characters that we translate from Unicode to
// alternate character encodings.  To do this, we use the standard VT100 ACS
// maps.  This is only done if the terminal lacks support for Unicode; we
// always prefer to emit Unicode glyphs when we are able.
func (t *tScreen) buildAcsMap() {
	const acsstr = "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~"

	t.acs = make(map[rune]string)
	for b, r := range vtACSNames {
		t.acs[r] = startAltChars + string(b) + endAltChars
	}
}

func (t *tScreen) scanInput(buf *bytes.Buffer) {
	for buf.Len() > 0 {
		utf := make([]byte, min(8, max(buf.Len()*2, 128)))
		nOut, nIn, e := t.decoder.Transform(utf, buf.Bytes(), true)
		_ = buf.Next(nIn)
		t.input.ScanUTF8(utf[:nOut])
		if e == transform.ErrShortSrc {
			return
		}
	}
}

func (t *tScreen) mainLoop(stopQ chan struct{}) {
	defer t.wg.Done()
	buf := &bytes.Buffer{}
	for {
		select {
		case <-stopQ:
			return
		case <-t.quit:
			return
		case <-t.resizeQ:
			t.Lock()
			t.cx = -1
			t.cy = -1
			t.resize()
			t.cells.Invalidate()
			t.draw()
			t.Unlock()
			continue
		case chunk := <-t.keychan:
			buf.Write(chunk)
			t.scanInput(buf)
		}
	}
}

func (t *tScreen) inputLoop(stopQ chan struct{}) {

	defer t.wg.Done()
	for {
		select {
		case <-stopQ:
			return
		default:
		}
		chunk := make([]byte, 128)
		n, e := t.tty.Read(chunk)
		switch e {
		case nil:
		default:
			t.Lock()
			running := t.running
			t.Unlock()
			if running {
				select {
				case t.eventQ <- NewEventError(e):
				case <-t.quit:
				}
			}
			return
		}
		if n > 0 {
			t.keychan <- chunk[:n]
		}
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

func (t *tScreen) RegisterRuneFallback(orig rune, fallback string) {
	t.Lock()
	t.fallback[orig] = fallback
	t.Unlock()
}

func (t *tScreen) UnregisterRuneFallback(orig rune) {
	t.Lock()
	delete(t.fallback, orig)
	t.Unlock()
}

func (t *tScreen) SetSize(w, h int) {
	if t.setWinSize != "" {
		t.TPuts(fmt.Sprintf(t.setWinSize, w, h))
	}
	t.cells.Invalidate()
	t.resize()
}

func (t *tScreen) Resize(int, int, int, int) {}

func (t *tScreen) Suspend() error {
	t.disengage()
	return nil
}

func (t *tScreen) Resume() error {
	return t.engage()
}

func (t *tScreen) Tty() (Tty, bool) {
	return t.tty, true
}

// engage is used to place the terminal in raw mode and establish screen size, etc.
// Think of this is as tcell "engaging" the clutch, as it's going to be driving the
// terminal interface.
func (t *tScreen) engage() error {
	t.Lock()
	defer t.Unlock()
	if t.tty == nil {
		return ErrNoScreen
	}
	t.tty.NotifyResize(func() {
		select {
		case t.resizeQ <- true:
		default:
		}
	})
	if t.running {
		return errors.New("already engaged")
	}
	if err := t.tty.Start(); err != nil {
		return err
	}
	t.running = true
	if ws, err := t.tty.WindowSize(); err == nil && ws.Width != 0 && ws.Height != 0 {
		t.cells.Resize(ws.Width, ws.Height)
	}
	stopQ := make(chan struct{})
	t.stopQ = stopQ
	t.enableMouse(t.mouseFlags)
	t.enablePasting(t.pasteEnabled)
	if t.focusEnabled {
		t.enableFocusReporting()
	}
	if os.Getenv("TCELL_ALTSCREEN") != "disable" {
		// Technically this may not be right, but every terminal we know about
		// (even Wyse 60) uses this to enter the alternate screen buffer, and
		// possibly save and restore the window title and/or icon.
		// (In theory there could be terminals that don't support X,Y cursor
		// positions without a setup command, but we don't support them.)
		t.TPuts(enterCA)
		t.TPuts(t.saveTitle)
	}
	t.TPuts(enterKeypad)
	t.TPuts(hideCursor)
	t.TPuts(enableAltChars)
	t.TPuts(disableAutoMargin)
	t.TPuts(clear)
	if t.title != "" && t.setTitle != "" {
		t.TPuts(fmt.Sprintf(t.setTitle, t.title))
	}
	t.TPuts(t.enableCsiU)

	t.wg.Add(2)
	go t.inputLoop(stopQ)
	go t.mainLoop(stopQ)
	return nil
}

// disengage is used to release the terminal back to support from the caller.
// Think of this as tcell disengaging the clutch, so that another application
// can take over the terminal interface.  This restores the TTY mode that was
// present when the application was first started.
func (t *tScreen) disengage() {

	t.Lock()
	if !t.running {
		t.Unlock()
		return
	}

	t.running = false
	stopQ := t.stopQ
	close(stopQ)
	_ = t.tty.Drain()
	t.Unlock()

	t.tty.NotifyResize(nil)
	// wait for everything to shut down
	t.wg.Wait()

	// shutdown the screen and disable special modes (e.g. mouse and bracketed paste)
	t.cells.Resize(0, 0)
	t.TPuts(showCursor)
	if t.cursorStyles != nil && t.cursorStyle != CursorStyleDefault {
		t.TPuts(t.cursorStyles[CursorStyleDefault])
	}
	if t.cursorFg != "" && t.cursorColor.Valid() {
		t.TPuts(t.cursorFg)
	}
	t.TPuts(exitKeypad)
	t.TPuts(resetFgBg)
	t.TPuts(sgr0)
	t.TPuts(enableAutoMargin)
	t.TPuts(t.disableCsiU)
	if os.Getenv("TCELL_ALTSCREEN") != "disable" {
		t.TPuts(t.restoreTitle)
		t.TPuts(clear)
		t.TPuts(exitCA)
	}
	t.enableMouse(0)
	t.enablePasting(false)
	t.disableFocusReporting()

	_ = t.tty.Stop()
}

// Beep emits a beep to the terminal.
func (t *tScreen) Beep() error {
	t.writeString(string(byte(7)))
	return nil
}

// finalize is used to at application shutdown, and restores the terminal
// to it's initial state.  It should not be called more than once.
func (t *tScreen) finalize() {
	t.disengage()
	_ = t.tty.Close()
}

func (t *tScreen) StopQ() <-chan struct{} {
	return t.quit
}

func (t *tScreen) EventQ() chan Event {
	return t.eventQ
}

func (t *tScreen) GetCells() *CellBuffer {
	return &t.cells
}

func (t *tScreen) SetTitle(title string) {
	t.Lock()
	t.title = title
	if t.setTitle != "" && t.running {
		t.TPuts(fmt.Sprintf(t.setTitle, title))
	}
	t.Unlock()
}

func (t *tScreen) SetClipboard(data []byte) {
	// Post binary data to the system clipboard.  It might be UTF-8, it might not be.
	t.Lock()
	if t.setClipboard != "" {
		encoded := base64.StdEncoding.EncodeToString(data)
		t.TPuts(fmt.Sprintf(t.setClipboard, encoded))
	}
	t.Unlock()
}

func (t *tScreen) GetClipboard() {
	t.Lock()
	if t.setClipboard != "" {
		t.TPuts(fmt.Sprintf(t.setClipboard, "?"))
	}
	t.Unlock()
}
