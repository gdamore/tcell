// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/vt"
)

// TestScanCodes tests that the keys all have correct scan codes.
func TestScanCodes(t *testing.T) {
	cases := []struct {
		name string
		key  vt.Key
		code vt.ScanCode
	}{
		{"A", vt.KeyA, 0x1e},
		{"B", vt.KeyB, 0x30},
		{"C", vt.KeyC, 0x2e},
		{"D", vt.KeyD, 0x20},
		{"E", vt.KeyE, 0x12},
		{"F", vt.KeyF, 0x21},
		{"G", vt.KeyG, 0x22},
		{"H", vt.KeyH, 0x23},
		{"I", vt.KeyI, 0x17},
		{"J", vt.KeyJ, 0x24},
		{"K", vt.KeyK, 0x25},
		{"L", vt.KeyL, 0x26},
		{"M", vt.KeyM, 0x32},
		{"N", vt.KeyN, 0x31},
		{"O", vt.KeyO, 0x18},
		{"P", vt.KeyP, 0x19},
		{"Q", vt.KeyQ, 0x10},
		{"R", vt.KeyR, 0x13},
		{"S", vt.KeyS, 0x1f},
		{"T", vt.KeyT, 0x14},
		{"U", vt.KeyU, 0x16},
		{"V", vt.KeyV, 0x2f},
		{"W", vt.KeyW, 0x11},
		{"X", vt.KeyX, 0x2d},
		{"Y", vt.KeyY, 0x15},
		{"Z", vt.KeyZ, 0x2c},
		{"1", vt.Key1, 0x02},
		{"2", vt.Key2, 0x03},
		{"3", vt.Key3, 0x04},
		{"4", vt.Key4, 0x05},
		{"5", vt.Key5, 0x06},
		{"6", vt.Key6, 0x07},
		{"7", vt.Key7, 0x08},
		{"8", vt.Key8, 0x09},
		{"9", vt.Key9, 0x0a},
		{"0", vt.Key0, 0x0b},
		{"enter", vt.KeyEnter, 0x1c},
		{"escape", vt.KeyEsc, 0x01},
		{"bs", vt.KeyBackspace, 0x0e},
		{"tab", vt.KeyTab, 0x0f},
		{"space", vt.KeySpace, 0x39},
		{"minus", vt.KeyMinus, 0x0c},
		{"equals", vt.KeyEqual, 0x0d},
		{"lbrace", vt.KeyLBrace, 0x1a},
		{"rbrace", vt.KeyRBrace, 0x1b},
		{"backslash", vt.KeyBackslash, 0x2b},
		{"semicolon", vt.KeySemi, 0x27},
		{"quote", vt.KeyQuote, 0x28},
		{"grave", vt.KeyGrave, 0x29},
		{"comma", vt.KeyComma, 0x33},
		{"period", vt.KeyPeriod, 0x34},
		{"slash", vt.KeySlash, 0x35},
		{"capsLock", vt.KeyCapsLock, 0x3a},
		{"f1", vt.KeyF1, 0x3b},
		{"f2", vt.KeyF2, 0x3c},
		{"f3", vt.KeyF3, 0x3d},
		{"f4", vt.KeyF4, 0x3e},
		{"f5", vt.KeyF5, 0x3f},
		{"f6", vt.KeyF6, 0x40},
		{"f7", vt.KeyF7, 0x41},
		{"f8", vt.KeyF8, 0x42},
		{"f9", vt.KeyF9, 0x43},
		{"f10", vt.KeyF10, 0x44},
		{"f11", vt.KeyF11, 0x57},
		{"f12", vt.KeyF12, 0x58},
		{"prtScr", vt.KeyPrtScr, 0x54},
		{"scrLock", vt.KeyScrLock, 0x46},
		{"pause", vt.KeyPause, 0xe046},
		{"ins", vt.KeyInsert, 0xe052},
		{"home", vt.KeyHome, 0xe047},
		{"pgUp", vt.KeyPgUp, 0xe049},
		{"del", vt.KeyDelete, 0xe053},
		{"end", vt.KeyEnd, 0xe04f},
		{"pgDn", vt.KeyPgDn, 0xe051},
		{"right", vt.KeyRight, 0xe04d},
		{"left", vt.KeyLeft, 0xe04b},
		{"down", vt.KeyDown, 0xe050},
		{"up", vt.KeyUp, 0xe048},
		{"numLock", vt.KeyNumLock, 0x45},
		{"div", vt.KeyPadDiv, 0xe035},
		{"mul", vt.KeyPadMul, 0x37},
		{"sub", vt.KeyPadSub, 0x4a},
		{"add", vt.KeyPadAdd, 0x4e},
		{"enter", vt.KeyPadEnter, 0xe01c},
		{"pad1", vt.KeyPad1, 0x4f},
		{"pad2", vt.KeyPad2, 0x50},
		{"pad3", vt.KeyPad3, 0x51},
		{"pad4", vt.KeyPad4, 0x4b},
		{"pad5", vt.KeyPad5, 0x4c},
		{"pad6", vt.KeyPad6, 0x4d},
		{"pad7", vt.KeyPad7, 0x47},
		{"pad8", vt.KeyPad8, 0x48},
		{"pad9", vt.KeyPad9, 0x49},
		{"pad0", vt.KeyPad0, 0x52},
		{"padDecimal", vt.KeyPadDec, 0x53},
		{"isoBackSlash", vt.KeyIsoBackSlash, 0x56},
		{"padEqual", vt.KeyPadEqual, 0x59},
		{"f13", vt.KeyF13, 0x64},
		{"f14", vt.KeyF14, 0x65},
		{"f15", vt.KeyF15, 0x66},
		{"f16", vt.KeyF16, 0x67},
		{"f17", vt.KeyF17, 0x68},
		{"f18", vt.KeyF18, 0x69},
		{"f19", vt.KeyF19, 0x6a},
		{"f20", vt.KeyF20, 0x6b},
		{"f21", vt.KeyF21, 0x6c},
		{"f22", vt.KeyF22, 0x6d},
		{"f23", vt.KeyF23, 0x6e},
		{"f24", vt.KeyF24, 0x76},
		{"padComma", vt.KeyPadComma, 0x7e},
		{"padIsoSlash", vt.KeyIsoSlash, 0x73},
		{"hiragana", vt.KeyHiragana, 0x70},
		{"yen", vt.KeyYen, 0x7d},
		{"convert", vt.KeyConvert, 0x79},
		{"noConvert", vt.KeyNonConvert, 0x7b},
		{"leftCtrl", vt.KeyLCtrl, 0x1d},
		{"leftShift", vt.KeyLShift, 0x2a},
		{"leftAlt", vt.KeyLAlt, 0x38},
		{"leftMeta", vt.KeyLMeta, 0xe05b},
		{"rightCtrl", vt.KeyRCtrl, 0xe01d},
		{"rightShift", vt.KeyRShift, 0x36},
		{"rightAlt", vt.KeyRAlt, 0xe038},
		{"rightMeta", vt.KeyRMeta, 0xe05c},
		{"menu", vt.KeyMenu, 0xe05d},
	}
	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			VerifyF(t, cases[i].key.ScanCode() == cases[i].code, "Scan code %x did not match %x", cases[i].key.ScanCode(), cases[i].code)
		})
	}

	// Let's also make sure that we have no duplicate scan codes.
	seen := make(map[vt.ScanCode]vt.Key)
	for i := range cases {
		sc := cases[i].key.ScanCode()
		if other, ok := seen[sc]; ok {
			t.Errorf("Duplicate scan %x code for key %x and %x", sc, cases[i].key, other)
		}
		seen[sc] = cases[i].key
	}
}

// TestBaseKeys tests that all the "base" keys have correct values and correct shifted values.
// These are the values used for the Kitty protocol.
func TestBaseKeys(t *testing.T) {
	cases := []struct {
		name    string
		key     vt.Key
		base    vt.BaseKey
		shifted rune
	}{
		{"A", vt.KeyA, 'a', 'A'},
		{"B", vt.KeyB, 'b', 'B'},
		{"C", vt.KeyC, 'c', 'C'},
		{"D", vt.KeyD, 'd', 'D'},
		{"E", vt.KeyE, 'e', 'E'},
		{"F", vt.KeyF, 'f', 'F'},
		{"G", vt.KeyG, 'g', 'G'},
		{"H", vt.KeyH, 'h', 'H'},
		{"I", vt.KeyI, 'i', 'I'},
		{"J", vt.KeyJ, 'j', 'J'},
		{"K", vt.KeyK, 'k', 'K'},
		{"L", vt.KeyL, 'l', 'L'},
		{"M", vt.KeyM, 'm', 'M'},
		{"N", vt.KeyN, 'n', 'N'},
		{"O", vt.KeyO, 'o', 'O'},
		{"P", vt.KeyP, 'p', 'P'},
		{"Q", vt.KeyQ, 'q', 'Q'},
		{"R", vt.KeyR, 'r', 'R'},
		{"S", vt.KeyS, 's', 'S'},
		{"T", vt.KeyT, 't', 'T'},
		{"U", vt.KeyU, 'u', 'U'},
		{"V", vt.KeyV, 'v', 'V'},
		{"W", vt.KeyW, 'w', 'W'},
		{"X", vt.KeyX, 'x', 'X'},
		{"Y", vt.KeyY, 'y', 'Y'},
		{"Z", vt.KeyZ, 'z', 'Z'},
		{"1", vt.Key1, '1', '!'},
		{"2", vt.Key2, '2', '@'},
		{"3", vt.Key3, '3', '#'},
		{"4", vt.Key4, '4', '$'},
		{"5", vt.Key5, '5', '%'},
		{"6", vt.Key6, '6', '^'},
		{"7", vt.Key7, '7', '&'},
		{"8", vt.Key8, '8', '*'},
		{"9", vt.Key9, '9', '('},
		{"0", vt.Key0, '0', ')'},
		{"minus", vt.KeyMinus, '-', '_'},
		{"equal", vt.KeyEqual, '=', '+'},
		{"lbrace", vt.KeyLBrace, '[', '{'},
		{"backslash", vt.KeyBackslash, '\\', '|'},
		{"rbrace", vt.KeyRBrace, ']', '}'},
		{"semi", vt.KeySemi, ';', ':'},
		{"backslash-iso", vt.KeyIsoBackSlash, '\\', '|'},
		{"comma", vt.KeyComma, ',', '<'},
		{"period", vt.KeyPeriod, '.', '>'},
		{"slash", vt.KeySlash, '/', '?'},
		{"yen", vt.KeyYen, '¥', '|'},
	}
	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			base := cases[i].key.KittyBase()
			shift := base.Shifted()
			VerifyF(t, base == cases[i].base, "Base code %x did not match %x", base, cases[i].base)
			VerifyF(t, shift == cases[i].shifted, "Base code %q %x did not match %q %x", string(shift), shift, string(cases[i].shifted), cases[i].shifted)
		})
	}
}

// TestKeyRepeat tests simple key repeat
func TestKeyRepeat(t *testing.T) {
	term := vt.NewMockTerm()
	defer MustClose(t, term)

	MustStart(t, term)

	// these are unreasonable repeat rates, but its somewhere to start
	term.SetRepeat(time.Millisecond*50, time.Millisecond*25)

	term.KeyPress(vt.KeyX)
	time.Sleep(100 * time.Millisecond)
	term.KeyPress(vt.KeyX)
	term.KeyRelease(vt.KeyX)
	term.KeyRelease(vt.KeyCapsLock)

	CheckRead(t, term, "xxxx") // 0 ms, 50 ms, 75 ms, 100 ms
}

// TestKeyRepeatCapsLock ensures that caps lock does not repeat
func TestKeyRepeatCapsLock(t *testing.T) {
	term := vt.NewMockTerm()
	defer MustClose(t, term)

	MustStart(t, term)

	// these are unreasonable repeat rates, but its somewhere to start
	term.SetRepeat(time.Millisecond*50, time.Millisecond*25)

	term.KeyPress(vt.KeyCapsLock)
	term.KeyPress(vt.KeyZ)
	time.Sleep(100 * time.Millisecond)
	term.KeyPress(vt.KeyZ)
	term.KeyRelease(vt.KeyZ)
	term.KeyRelease(vt.KeyCapsLock)

	CheckRead(t, term, "ZZZZ") // 0 ms, 50 ms, 75 ms, 100 ms
}

// TestKeyRepeatNoAlt ensures that alt keys do not repeat
func TestKeyRepeatNoAlt(t *testing.T) {
	term := vt.NewMockTerm()
	defer MustClose(t, term)

	MustStart(t, term)

	// these are unreasonable repeat rates, but its somewhere to start
	term.SetRepeat(time.Millisecond*50, time.Millisecond*25)

	term.KeyPress(vt.KeyLAlt)
	term.KeyPress(vt.KeyZ)
	time.Sleep(100 * time.Millisecond)
	term.KeyPress(vt.KeyZ)
	term.KeyRelease(vt.KeyZ)
	term.KeyRelease(vt.KeyLAlt)

	CheckRead(t, term, "\x1bz") // 0 ms, 50 ms, 75 ms, 100 ms
}

// TestKeyRepeatShift ensures that shifted keys still work as long as repeat is held down.
func TestKeyRepeatShift(t *testing.T) {
	term := vt.NewMockTerm()
	defer MustClose(t, term)

	MustStart(t, term)

	// these are unreasonable repeat rates, but its somewhere to start
	term.SetRepeat(time.Millisecond*50, time.Millisecond*25)

	term.KeyPress(vt.KeyLShift)
	term.KeyPress(vt.Key1)
	time.Sleep(100 * time.Millisecond)
	term.KeyPress(vt.Key1)
	term.KeyRelease(vt.Key1)
	term.KeyRelease(vt.KeyLShift)

	CheckRead(t, term, "!!!!") // 0 ms, 50 ms, 75 ms, 100 ms
}

// TestKeyRepeatShiftRelease ensures that releasing shift breaks repeat.
func TestKeyRepeatShiftRelease(t *testing.T) {
	term := vt.NewMockTerm()
	defer MustClose(t, term)

	MustStart(t, term)

	// these are unreasonable repeat rates, but its somewhere to start
	term.SetRepeat(time.Millisecond*50, time.Millisecond*25)

	term.KeyPress(vt.KeyLShift)
	term.KeyPress(vt.Key2)
	term.KeyRelease(vt.KeyLShift)
	time.Sleep(100 * time.Millisecond)
	term.KeyPress(vt.Key2)
	term.KeyRelease(vt.Key2)

	CheckRead(t, term, "@222") // 0 ms, 50 ms, 75 ms, 100 ms
}

func TestKeyRepeatCursor(t *testing.T) {
	term := vt.NewMockTerm()
	defer MustClose(t, term)

	MustStart(t, term)

	// these are unreasonable repeat rates, but its somewhere to start
	term.SetRepeat(time.Millisecond*50, time.Millisecond*25)

	term.KeyPress(vt.KeyRight)
	time.Sleep(100 * time.Millisecond)
	term.KeyPress(vt.KeyRight)
	term.KeyRelease(vt.KeyRight)
	term.KeyRelease(vt.KeyRight)

	CheckRead(t, term, "\x1b[C\x1b[C\x1b[C\x1b[C") // 0 ms, 50 ms, 75 ms, 100 ms
}
