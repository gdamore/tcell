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
	"fmt"
	"strings"
	"time"
)

// EventKey represents a key press.  Usually this is a key press followed
// by a key release, but since terminal programs don't have a way to report
// key release events, we usually get just one event.  If a key is held down
// then the terminal may synthesize repeated key presses at some predefined
// rate.  We have no control over that, nor visibility into it.
//
// In some cases, we can have a modifier key, such as ModAlt, that can be
// generated with a key press.  (This usually is represented by having the
// high bit set, or in some cases, by sending an ESC prior to the rune.)
//
// If the value of Key() is KeyRune, then the actual key value will be
// available with the Rune() method.  This will be the case for most keys.
// In most situations, the modifiers will not be set.  For example, if the
// rune is 'A', this will be reported without the ModShift bit set, since
// really can't tell if the Shift key was pressed (it might have been CAPSLOCK,
// or a terminal that only can send capitals, or keyboard with separate
// capital letters from lower case letters).
//
// Generally, terminal applications have far less visibility into keyboard
// activity than graphical applications.  Hence, they should avoid depending
// overly much on availability of modifiers, or the availability of any
// specific keys.
type EventKey struct {
	t   time.Time
	mod ModMask
	key Key
	ch  rune
}

// When returns the time when this Event was created, which should closely
// match the time when the key was pressed.
func (*EventKey) When() time.Time {
	return time.Now()
}

// Rune returns the rune corresponding to the key press, if it makes sense.
// The result is only defined if the value of Key() is KeyRune.
func (ev *EventKey) Rune() rune {
	return ev.ch
}

// Key returns a virtual key code.  We use this to identify specific key
// codes, such as KeyEnter, etc.  Most control and function keys are reported
// with unique Key values.  Normal alphanumeric and punctuation keys will
// generally return KeyRune here; the specific key can be further decoded
// using the Rune() function.
func (ev *EventKey) Key() Key {
	return ev.key
}

// ModMask returns the modifiers that were present with the key press.  Note
// that not all platforms and terminals support this equally well, and some
// cases we will not not know for sure.  Hence, applications should avoid
// using this in most circumstances.
func (ev *EventKey) Mod() ModMask {
	return ev.mod
}

// Name returns a printable value or the key stroke.  This can be used
// when printing the event, for example.
func (ev *EventKey) Name() string {
	s := ""
	m := []string{}
	if ev.mod&ModShift != 0 {
		m = append(m, "Shift")
	}
	if ev.mod&ModAlt != 0 {
		m = append(m, "Alt")
	}
	if ev.mod&ModMeta != 0 {
		m = append(m, "Meta")
	}
	if ev.mod&ModCtrl != 0 {
		m = append(m, "Ctrl")
	}

	switch ev.key {
	case KeyRune:
		s = "Rune[" + string(ev.ch) + "]"
	case KeySpace:
		s = "Space"
	case KeyEnter:
		s = "Enter"
	case KeyBackspace:
		s = "Backspace"
	case KeyTab:
		s = "Tab"
	case KeyBacktab:
		s = "Backtab"
	case KeyEsc:
		s = "Esc"
	case KeyBackspace2:
		s = "Backspace2"
	case KeyDelete:
		s = "Delete"
	case KeyInsert:
		s = "Insert"
	case KeyUp:
		s = "Up"
	case KeyDown:
		s = "Down"
	case KeyLeft:
		s = "Left"
	case KeyRight:
		s = "Right"
	case KeyHome:
		s = "Home"
	case KeyEnd:
		s = "End"
	case KeyUpLeft:
		s = "UpLeft"
	case KeyUpRight:
		s = "UpRight"
	case KeyDownLeft:
		s = "DownLeft"
	case KeyDownRight:
		s = "DownRight"
	case KeyCenter:
		s = "Center"
	case KeyPgDn:
		s = "PgDn"
	case KeyPgUp:
		s = "PgUp"
	case KeyClear:
		s = "Clear"
	case KeyExit:
		s = "Exit"
	case KeyCancel:
		s = "Cancel"
	case KeyPause:
		s = "Pause"
	case KeyPrint:
		s = "Print"
	case KeyCtrlSpace:
		s = "Ctrl-Space"
	case KeyCtrlUnderscore:
		s = "Ctrl-_"
	case KeyCtrlRightSq:
		s = "Ctrl-]"
	case KeyCtrlBackslash:
		s = "Ctrl-\\"
	case KeyCtrlCarat:
		s = "Ctrl-^"
	default:
		if ev.key >= KeyF1 && ev.key <= KeyF64 {
			s = fmt.Sprintf("F%d", int(ev.key-KeyF1)+1)
		} else if ev.key >= KeyCtrlA && ev.key <= KeyCtrlZ {
			s = fmt.Sprintf("Ctrl-%c",
				rune(ev.key-KeyCtrlA)+'A')
		} else {
			s = fmt.Sprintf("Key[%d,%d]", ev.key, int(ev.ch))
		}
	}

	if len(m) != 0 {
		if ev.mod&ModCtrl != 0 && strings.HasPrefix(s, "Ctrl-") {
			s = s[5:]
		}
		return fmt.Sprintf("%s+%s", strings.Join(m, "+"), s)
	}
	return s
}

// NewEventKey attempts to create a suitable event.  It parses the various
// ASCII control sequences if KeyRune is passed for Key, but if the caller
// has more precise information it should set that specifically.  Callers
// that aren't sure about modifier state (most) should just pass ModNone.
func NewEventKey(k Key, ch rune, mod ModMask) *EventKey {
	if k == KeyRune && (ch <= ' ' || ch == 0x7f) {
		// Turn specials into proper key codes.  This is for
		// control characters and the DEL.
		k = Key(ch)
		if mod == ModNone && ch < ' ' {
			switch Key(ch) {
			case KeyBackspace, KeyTab, KeyEsc, KeyEnter:
				// these keys are directly typeable without CTRL
			default:
				// most likely entered with a CTRL keypress
				mod = ModCtrl
			}
		}
	}
	return &EventKey{t: time.Now(), key: k, ch: ch, mod: mod}
}

// ModMask is a mask of modifier keys.  Note that it will not always be
// possible to report modifier keys.
type ModMask int16

const (
	ModShift ModMask = 1 << iota
	ModCtrl
	ModAlt
	ModMeta
)
const ModNone ModMask = 0

// Key is a generic value for representing keys, and especially special
// keys (function keys, cursor movement keys, etc.)  For normal keys, like
// ASCII letters, we use KeyRune, and then expect the application to
// inspect the Rune() member of the EventKey.
type Key int16

const (
	KeyRune Key = iota + 256
	KeyUp
	KeyDown
	KeyRight
	KeyLeft
	KeyUpLeft
	KeyUpRight
	KeyDownLeft
	KeyDownRight
	KeyCenter
	KeyPgUp
	KeyPgDn
	KeyHome
	KeyEnd
	KeyInsert
	KeyDelete
	KeyHelp
	KeyExit
	KeyClear
	KeyCancel
	KeyPrint
	KeyPause
	KeyBacktab
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyF26
	KeyF27
	KeyF28
	KeyF29
	KeyF30
	KeyF31
	KeyF32
	KeyF33
	KeyF34
	KeyF35
	KeyF36
	KeyF37
	KeyF38
	KeyF39
	KeyF40
	KeyF41
	KeyF42
	KeyF43
	KeyF44
	KeyF45
	KeyF46
	KeyF47
	KeyF48
	KeyF49
	KeyF50
	KeyF51
	KeyF52
	KeyF53
	KeyF54
	KeyF55
	KeyF56
	KeyF57
	KeyF58
	KeyF59
	KeyF60
	KeyF61
	KeyF62
	KeyF63
	KeyF64
)

const (
	KeyCtrlSpace Key = iota
	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlI
	KeyCtrlJ
	KeyCtrlK
	KeyCtrlL
	KeyCtrlM
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ
	KeyCtrlLeftSq // Escape
	KeyCtrlBackslash
	KeyCtrlRightSq
	KeyCtrlCarat
	KeyCtrlUnderscore
)

// Special values - these are fixed in an attempt to make it more likely
// that aliases will encode the same way.

// These are the defined ASCII values for key codes.  They generally match
// with KeyCtrl values.
const (
	KeyNUL Key = iota
	KeySOH
	KeySTX
	KeyETX
	KeyEOT
	KeyENQ
	KeyACK
	KeyBEL
	KeyBS
	KeyTAB
	KeyLF
	KeyVT
	KeyFF
	KeyCR
	KeySO
	KeySI
	KeyDLE
	KeyDC1
	KeyDC2
	KeyDC3
	KeyDC4
	KeyNAK
	KeySYN
	KeyETB
	KeyCAN
	KeyEM
	KeySUB
	KeyESC
	KeyFS
	KeyGS
	KeyRS
	KeyUS
	KeySP
	KeyDEL Key = 0x7F
)

const (
	KeyBackspace  = KeyBS
	KeyTab        = KeyTAB
	KeyEsc        = KeyESC
	KeyEscape     = KeyESC
	KeyEnter      = KeyCR
	KeySpace      = KeySP
	KeyBackspace2 = KeyDEL // This is delete back, not forward
)
