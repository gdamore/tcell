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

// Package termbox is a compatibility layer to allow tcells to emulate
// the github.com/nsf/termbox package.
package termbox

import (
	"errors"

	"github.com/gdamore/tcell"
)

var screen tcell.Screen
var outMode OutputMode

func Init() error {
	outMode = OutputNormal
	if s, e := tcell.NewScreen(); e != nil {
		return e
	} else if e = s.Init(); e != nil {
		return e
	} else {
		screen = s
		return nil
	}
}

func Close() {
	screen.Fini()
}

func Flush() error {
	screen.Show()
	return nil
}

func SetCursor(x, y int) {
	screen.ShowCursor(x, y)
}

func HideCursor() {
	SetCursor(-1, -1)
}

func Size() (int, int) {
	return screen.Size()
}

type Attribute uint16

const (
	ColorDefault Attribute = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)
const (
	AttrBold Attribute = 1 << (9 + iota)
	AttrUnderline
	AttrReverse
)

func mkStyle(fg, bg Attribute) tcell.Style {
	st := tcell.StyleDefault

	f := int(fg) & 0x1ff
	b := int(bg) & 0x1ff

	switch outMode {
	case Output256:
		break
	case Output216:
		if f > 216 {
			f = int(ColorDefault)
		} else if f != int(ColorDefault) {
			f += 16
		}
		if b > 216 {
			b = int(ColorDefault)
		} else if b != int(ColorDefault) {
			b += 16
		}
	case OutputGrayscale:
		if f > 24 {
			f = int(ColorDefault)
		} else if f != int(ColorDefault) {
			f += 232
		}
		if b > 24 {
			b = int(ColorDefault)
		} else if b != int(ColorDefault) {
			b += 232
		}
	case OutputNormal:
		f &= 0xf
		b &= 0xf
	}
	st = st.Foreground(tcell.Color(f))
	st = st.Background(tcell.Color(b))
	if (fg&AttrBold != 0) || (bg&AttrBold != 0) {
		st = st.Bold(true)
	}
	if (fg&AttrUnderline != 0) || (bg&AttrUnderline != 0) {
		st = st.Underline(true)
	}
	if (fg&AttrReverse != 0) || (bg&AttrReverse != 0) {
		st = st.Reverse(true)
	}
	return st
}

func Clear(fg, bg Attribute) {
	st := mkStyle(fg, bg)
	w, h := screen.Size()
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			screen.SetContent(col, row, ' ', nil, st)
		}
	}
}

type InputMode int

const (
	InputCurrent InputMode = iota
	InputEsc
	InputAlt
	InputMouse
)

func SetInputMode(mode InputMode) InputMode {
	// We don't do anything else right now
	return InputEsc
}

type OutputMode int

const (
	OutputCurrent OutputMode = iota
	OutputNormal
	Output256
	Output216
	OutputGrayscale
)

func SetOutputMode(mode OutputMode) OutputMode {
	if screen.Colors() < 256 {
		mode = OutputNormal
	}
	switch mode {
	case OutputCurrent:
		return outMode
	case OutputNormal, Output256, Output216, OutputGrayscale:
		outMode = mode
		return mode
	default:
		return outMode
	}
}

func Sync() error {
	screen.Sync()
	return nil
}

func SetCell(x, y int, ch rune, fg, bg Attribute) {
	st := mkStyle(fg, bg)
	screen.SetContent(x, y, ch, nil, st)
}

type EventType uint8
type Modifier tcell.ModMask
type Key tcell.Key

type Event struct {
	Type   EventType
	Mod    Modifier
	Key    Key
	Ch     rune
	Width  int
	Height int
	Err    error
	MouseX int
	MouseY int
	N      int
}

const (
	EventNone EventType = iota
	EventKey
	EventResize
	EventMouse
	EventInterrupt
	EventError
	EventRaw
)

const (
	KeyF1         = Key(tcell.KeyF1)
	KeyF2         = Key(tcell.KeyF2)
	KeyF3         = Key(tcell.KeyF3)
	KeyF4         = Key(tcell.KeyF4)
	KeyF5         = Key(tcell.KeyF5)
	KeyF6         = Key(tcell.KeyF6)
	KeyF7         = Key(tcell.KeyF7)
	KeyF8         = Key(tcell.KeyF8)
	KeyF9         = Key(tcell.KeyF9)
	KeyF10        = Key(tcell.KeyF10)
	KeyF11        = Key(tcell.KeyF11)
	KeyF12        = Key(tcell.KeyF12)
	KeyInsert     = Key(tcell.KeyInsert)
	KeyDelete     = Key(tcell.KeyDelete)
	KeyHome       = Key(tcell.KeyHome)
	KeyEnd        = Key(tcell.KeyEnd)
	KeyArrowUp    = Key(tcell.KeyUp)
	KeyArrowDown  = Key(tcell.KeyDown)
	KeyArrowRight = Key(tcell.KeyRight)
	KeyArrowLeft  = Key(tcell.KeyLeft)
	KeyCtrlA      = Key(tcell.KeyCtrlA)
	KeyCtrlB      = Key(tcell.KeyCtrlB)
	KeyCtrlC      = Key(tcell.KeyCtrlC)
	KeyCtrlD      = Key(tcell.KeyCtrlD)
	KeyCtrlE      = Key(tcell.KeyCtrlE)
	KeyCtrlF      = Key(tcell.KeyCtrlF)
	KeyCtrlG      = Key(tcell.KeyCtrlG)
	KeyCtrlH      = Key(tcell.KeyCtrlH)
	KeyCtrlI      = Key(tcell.KeyCtrlI)
	KeyCtrlJ      = Key(tcell.KeyCtrlJ)
	KeyCtrlK      = Key(tcell.KeyCtrlK)
	KeyCtrlL      = Key(tcell.KeyCtrlL)
	KeyCtrlM      = Key(tcell.KeyCtrlM)
	KeyCtrlN      = Key(tcell.KeyCtrlN)
	KeyCtrlO      = Key(tcell.KeyCtrlO)
	KeyCtrlP      = Key(tcell.KeyCtrlP)
	KeyCtrlQ      = Key(tcell.KeyCtrlQ)
	KeyCtrlR      = Key(tcell.KeyCtrlR)
	KeyCtrlS      = Key(tcell.KeyCtrlS)
	KeyCtrlT      = Key(tcell.KeyCtrlT)
	KeyCtrlU      = Key(tcell.KeyCtrlU)
	KeyCtrlV      = Key(tcell.KeyCtrlV)
	KeyCtrlW      = Key(tcell.KeyCtrlW)
	KeyCtrlX      = Key(tcell.KeyCtrlX)
	KeyCtrlY      = Key(tcell.KeyCtrlY)
	KeyCtrlZ      = Key(tcell.KeyCtrlZ)
	KeyBackspace  = Key(tcell.KeyBackspace)
	KeyBackspace2 = Key(tcell.KeyBackspace2)
	KeyTab        = Key(tcell.KeyTab)
	KeyEnter      = Key(tcell.KeyEnter)
	KeySpace      = Key(tcell.KeySpace)
	KeyEsc        = Key(tcell.KeyEscape)
	KeyPgdn       = Key(tcell.KeyPgDn)
	KeyPgup       = Key(tcell.KeyPgUp)
	MouseLeft     = Key(tcell.KeyF63) // arbitrary assignments
	MouseRight    = Key(tcell.KeyF62)
	MouseMiddle   = Key(tcell.KeyF61)
)

const (
	ModAlt = Modifier(tcell.ModAlt)
)

func makeEvent(tev tcell.Event) Event {
	switch tev := tev.(type) {
	case *tcell.EventInterrupt:
		return Event{Type: EventInterrupt}
	case *tcell.EventResize:
		w, h := tev.Size()
		return Event{Type: EventResize, Width: w, Height: h}
	case *tcell.EventKey:
		k := tev.Key()
		ch := rune(0)
		if k == tcell.KeyRune {
			ch = tev.Rune()
		}
		mod := tev.Mod()
		return Event{
			Type: EventKey,
			Key:  Key(k),
			Ch:   ch,
			Mod:  Modifier(mod),
		}
	default:
		return Event{Type: EventNone}
	}
}

func ParseEvent(data []byte) Event {
	// Not supported
	return Event{Type: EventError, Err: errors.New("no raw events")}
}

func PollRawEvent(data []byte) Event {
	// Not supported
	return Event{Type: EventError, Err: errors.New("no raw events")}
}

func PollEvent() Event {
	ev := screen.PollEvent()
	return makeEvent(ev)
}

func Interrupt() {
	screen.PostEvent(tcell.NewEventInterrupt(nil))
}

type Cell struct {
	Ch rune
	Fg Attribute
	Bg Attribute
}
