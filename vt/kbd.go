// Copyright 2025 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vt

type KbdEvent struct {
	Down   bool     // true if event is for key down event
	Repeat int      // if > 1, a repeat count
	Code   KeyCode  // logical key code (X11 key symbol, e.g. 'A')
	Base   KeyCode  // base key code (physical key, e.g 'a'), may be zero if same as code
	Mod    Modifier // modifiers
	Utf    string   // if non-empty, the unicode content for this
}

type Modifier int

const (
	ModNone  = Modifier(0)
	ModShift = Modifier(1 << iota)
	ModCtrl
	ModAlt
	ModMeta
	ModHyper
	ModCapsLock
	ModNumLock
)

// KeyCode represents either a physical key, or logical key.
// When it is the "code", it will represent the logical key such as 'A'.
// But (on English keyboards), 'A' becomes 'a' (with the ModShift) modifier,
// as a base.  (Sometimes we don't care about the logical keyboard, but only
// about a physical key, and using the base lets us deal with that.)
//
// Normally KeyCodes are numerically the value of the key they represent.  E.g
// 65 for 'A' or 97 for 'a'.
//
// Following the Kitty protocol, we use the Private Use Area 57344 - 63743 for keys
// that do not correspond to a Unicode value (such as function keys).  However, we
// do not use the same values as we need to represent more keys here than kitty does.
//
// This should be sufficient to support all known keyboard protocols.
//
// Note that Windows defines key codes for mice and game pads. Those should be treated
// as different events entirely.
type KeyCode rune

const (
	// definitions in ASCII space
	KcBackspace = KeyCode(0x08)
	KcTab       = KeyCode(0x09)
	KcReturn    = KeyCode(0x0d)
	KcEsc       = KeyCode(0x1B)
	KcSpace     = KeyCode(0x20)
	KcDelete    = KeyCode(0x7F) // not KcDel, but legacy Unix DELETE

	// misc function keys
	KcCancel = KeyCode(0xf000) + iota
	KcPause
	KcMenu
	KcSelect
	KcPrint
	KcExecute
	KcPrtScr
	KcHelp
	KcSleep
	KcApp // application key

	// navigation keys
	KcUp
	KcDown
	KcLeft
	KcRight
	KcCenter // center key in cursor keypad
	KcUpLeft
	KcUpRight
	KcDownLeft
	KcDownRight
	KcPgUp
	KcPgDn
	KcHome
	KcEnd
	KcIns
	KcDel

	// keypad keys
	KcPad0
	KcPad1
	KcPad2
	KcPad3
	KcPad4
	KcPad5
	KcPad6
	KcPad7
	KcPad8
	KcPad9
	KcPadDec // decimal (.)
	KcPadDiv // divide (/)
	KcPadMul // multiply (*)
	KcPadAdd // add (+)
	KcPadSub // subtract (-)
	KcPadEq  // equals (=)
	KcPadSep // separator (comma)
	KcPadEnter
	KcPadUp
	KcPadDown
	KcPadLeft
	KcPadRight
	KcPadPgUp
	KcPadPgDn
	KcPadHome
	KcPadEnd
	KcPadBegin
	KcPadIns
	KcPadDel

	// media keys
	KcMediaPlay
	KcMediaStop
	KcMediaPause // usually play is just modal
	KcMediaRewind
	KcMediaForward
	KcMediaNext
	KcMediaPrev
	KcMediaVolUp
	KcMediaVolDn
	KcMediaMute
	KcMediaRecord

	// browser keys (rarely supported)
	KcBrowserBack
	KcBrowserForward
	KcBrowserStop
	KcBrowserRefresh
	KcBrowserSearch
	KcBrowserFavorite
	KcBrowserHome

	// application keys
	KcApp1
	KcApp2
	KcApp3
	KcApp4
	KcApp5
	KcApp6
	KcApp7
	KcApp8

	// modifier keys (can be presented separately)
	KcRShift
	KcLShift
	KcLCtrl
	KcRCtrl
	KcLSuper
	KcRSuper
	KcLMeta
	KcRMeta
	KcLHyper
	KcRHyper
	KcIsoL3 // ISO level 3 shift - usually AltGr
	KcIsoL5 // ISO level 5 shift - custom usually

	// lock keys
	KcCapsLock
	KcScrLock
	KcNumLock

	// IME related
	KcKana
	KcHangul
	KcJunja
	KcKanji
	KcImeOn
	KcImeOff
	KcConvert    // IME convert
	KcNonConvert // IME non-convert
	KcAccept     // IME accept

	// virtual unicode - use to inject raw UTF-8
	KcUTF

	// up to 256 F-keys
	KcF1 = KeyCode(0xf100) + iota
	KcF2
	KcF3
	KcF4
	KcF5
	KcF6
	KcF7
	KcF8
	KcF9
	KcF10
	KcF11
	KcF12
	KcF13
	KcF14
	KcF15
	KcF16
	KcF17
	KcF18
	KcF19
	KcF20
	KcF21
	KcF22
	KcF23
	KcF24
	KcF25
	KcF26
	KcF27
	KcF28
	KcF29
	KcF30
	KcF31
	KcF32
	KcF33
	KcF34
)
