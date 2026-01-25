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

	"github.com/gdamore/tcell/v3/vt"
	"github.com/gdamore/tcell/v3/vt/layouts/us"
)

func TestGetLayout(t *testing.T) {
	l := vt.GetLayout("US International")
	if l == nil {
		t.Fatalf("no keyboard layout found")
	}
	if l != us.UsInternational {
		t.Errorf("got the wrong layout: %s", l.Name)
	}
}

func TestUsAltGr(t *testing.T) {

	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	term.SetLayout(vt.GetLayout("US International"))
	MustStart(t, term)

	cases := []struct {
		name string
		keys []vt.Key
		want string
	}{
		{"AltGr-0", []vt.Key{vt.KeyRAlt, vt.Key0}, "’"},
		{"AltGr-1", []vt.Key{vt.KeyRAlt, vt.Key1}, "¡"},
		{"AltGr-2", []vt.Key{vt.KeyRAlt, vt.Key2}, "²"},
		{"AltGr-3", []vt.Key{vt.KeyRAlt, vt.Key3}, "³"},
		{"AltGr-4", []vt.Key{vt.KeyRAlt, vt.Key4}, "¤"},
		{"AltGr-5", []vt.Key{vt.KeyRAlt, vt.Key5}, "€"},
		{"AltGr-6", []vt.Key{vt.KeyRAlt, vt.Key6}, "¼"},
		{"AltGr-7", []vt.Key{vt.KeyRAlt, vt.Key7}, "½"},
		{"AltGr-8", []vt.Key{vt.KeyRAlt, vt.Key8}, "¾"},
		{"AltGr-9", []vt.Key{vt.KeyRAlt, vt.Key9}, "‘"},
		{"AltGr-A", []vt.Key{vt.KeyRAlt, vt.KeyA}, "á"},
		{"AltGr-C", []vt.Key{vt.KeyRAlt, vt.KeyC}, "©"},
		{"AltGr-D", []vt.Key{vt.KeyRAlt, vt.KeyD}, "ð"},
		{"AltGr-E", []vt.Key{vt.KeyRAlt, vt.KeyE}, "é"},
		{"AltGr-I", []vt.Key{vt.KeyRAlt, vt.KeyI}, "í"},
		{"AltGr-L", []vt.Key{vt.KeyRAlt, vt.KeyL}, "ø"},
		{"AltGr-M", []vt.Key{vt.KeyRAlt, vt.KeyM}, "µ"},
		{"AltGr-N", []vt.Key{vt.KeyRAlt, vt.KeyN}, "ñ"},
		{"AltGr-O", []vt.Key{vt.KeyRAlt, vt.KeyO}, "ó"},
		{"AltGr-P", []vt.Key{vt.KeyRAlt, vt.KeyP}, "ö"},
		{"AltGr-Q", []vt.Key{vt.KeyRAlt, vt.KeyQ}, "ä"},
		{"AltGr-R", []vt.Key{vt.KeyRAlt, vt.KeyR}, "®"},
		{"AltGr-S", []vt.Key{vt.KeyRAlt, vt.KeyS}, "ß"},
		{"AltGr-T", []vt.Key{vt.KeyRAlt, vt.KeyT}, "þ"},
		{"AltGr-U", []vt.Key{vt.KeyRAlt, vt.KeyU}, "ú"},
		{"AltGr-W", []vt.Key{vt.KeyRAlt, vt.KeyW}, "å"},
		{"AltGr-Y", []vt.Key{vt.KeyRAlt, vt.KeyY}, "ü"},
		{"AltGr-Z", []vt.Key{vt.KeyRAlt, vt.KeyZ}, "æ"},
		{"AltGr-Semi", []vt.Key{vt.KeyRAlt, vt.KeySemi}, "¶"},
		{"AltGr-Equal", []vt.Key{vt.KeyRAlt, vt.KeyEqual}, "×"},
		{"AltGr-Comma", []vt.Key{vt.KeyRAlt, vt.KeyComma}, "ç"},
		{"AltGr-Minus", []vt.Key{vt.KeyRAlt, vt.KeyMinus}, "¥"},
		{"AltGr-Slash", []vt.Key{vt.KeyRAlt, vt.KeySlash}, "¿"},
		{"AltGr-LBrace", []vt.Key{vt.KeyRAlt, vt.KeyLBrace}, "«"},
		{"AltGr-Backslash", []vt.Key{vt.KeyRAlt, vt.KeyBackslash}, "¬"},
		{"AltGr-RBrace", []vt.Key{vt.KeyRAlt, vt.KeyRBrace}, "»"},
		{"AltGr-Quote", []vt.Key{vt.KeyRAlt, vt.KeyQuote}, "´"},

		{"S-AltGr-0", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.Key0}, ""},
		{"S-AltGr-1", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.Key1}, "¹"},
		{"S-AltGr-2", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.Key2}, ""},
		{"S-AltGr-3", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.Key3}, ""},
		{"S-AltGr-4", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.Key4}, "£"},
		{"S-AltGr-5", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.Key5}, ""},
		{"S-AltGr-6", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.Key6}, ""},
		{"S-AltGr-7", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.Key7}, ""},
		{"S-AltGr-8", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.Key8}, ""},
		{"S-AltGr-9", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.Key9}, ""},
		{"S-AltGr-A", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyA}, "Á"},
		{"S-AltGr-B", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyB}, ""},
		{"S-AltGr-C", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyC}, "¢"},
		{"S-AltGr-D", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyD}, "Ð"},
		{"S-AltGr-E", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyE}, "É"},
		{"S-AltGr-F", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyF}, ""},
		{"S-AltGr-G", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyG}, ""},
		{"S-AltGr-H", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyH}, ""},
		{"S-AltGr-I", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyI}, "Í"},
		{"S-AltGr-J", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyJ}, ""},
		{"S-AltGr-K", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyK}, ""},
		{"S-AltGr-L", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyL}, "Ø"},
		{"S-AltGr-M", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyM}, ""},
		{"S-AltGr-N", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyN}, "Ñ"},
		{"S-AltGr-O", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyO}, "Ó"},
		{"S-AltGr-P", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyP}, "Ö"},
		{"S-AltGr-Q", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyQ}, "Ä"},
		{"S-AltGr-R", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyR}, ""},
		{"S-AltGr-S", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyS}, "§"},
		{"S-AltGr-T", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyT}, "Þ"},
		{"S-AltGr-U", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyU}, "Ú"},
		{"S-AltGr-V", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyV}, ""},
		{"S-AltGr-W", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyW}, "Å"},
		{"S-AltGr-X", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyX}, ""},
		{"S-AltGr-Y", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyY}, "Ü"},
		{"S-AltGr-Z", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyZ}, "Æ"},
		{"S-AltGr-Semi", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeySemi}, "°"},
		{"S-AltGr-Equal", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyEqual}, "÷"},
		{"S-AltGr-Comma", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyComma}, "Ç"},
		{"S-AltGr-Minus", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyMinus}, ""},
		{"S-AltGr-Slash", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeySlash}, ""},
		{"S-AltGr-LBrace", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyLBrace}, ""},
		{"S-AltGr-RBrace", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyRBrace}, ""},
		{"S-AltGr-BackSlash", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyBackslash}, "¦"},
		{"S-AltGr-Quote", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyQuote}, "¨"},
		{"S-AltGr-IsoBackSlash", []vt.Key{vt.KeyRAlt, vt.KeyLShift, vt.KeyIsoBackSlash}, ""},

		{"Alt-Ctrl-0", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.Key0}, "’"},
		{"Alt-Ctrl-1", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.Key1}, "¡"},
		{"Alt-Ctrl-2", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.Key2}, "²"},
		{"Alt-Ctrl-3", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.Key3}, "³"},
		{"Alt-Ctrl-4", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.Key4}, "¤"},
		{"Alt-Ctrl-5", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.Key5}, "€"},
		{"Alt-Ctrl-6", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.Key6}, "¼"},
		{"Alt-Ctrl-7", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.Key7}, "½"},
		{"Alt-Ctrl-8", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.Key8}, "¾"},
		{"Alt-Ctrl-9", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.Key9}, "‘"},
		{"Alt-Ctrl-A", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyA}, "á"},
		{"Alt-Ctrl-C", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyC}, "©"},
		{"Alt-Ctrl-D", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyD}, "ð"},
		{"Alt-Ctrl-E", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyE}, "é"},
		{"Alt-Ctrl-I", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyI}, "í"},
		{"Alt-Ctrl-L", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyL}, "ø"},
		{"Alt-Ctrl-M", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyM}, "µ"},
		{"Alt-Ctrl-N", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyN}, "ñ"},
		{"Alt-Ctrl-O", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyO}, "ó"},
		{"Alt-Ctrl-P", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyP}, "ö"},
		{"Alt-Ctrl-Q", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyQ}, "ä"},
		{"Alt-Ctrl-R", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyR}, "®"},
		{"Alt-Ctrl-S", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyS}, "ß"},
		{"Alt-Ctrl-T", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyT}, "þ"},
		{"Alt-Ctrl-U", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyU}, "ú"},
		{"Alt-Ctrl-W", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyW}, "å"},
		{"Alt-Ctrl-Y", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyY}, "ü"},
		{"Alt-Ctrl-Z", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyZ}, "æ"},
		{"Alt-Ctrl-Semi", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeySemi}, "¶"},
		{"Alt-Ctrl-Equal", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyEqual}, "×"},
		{"Alt-Ctrl-Comma", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyComma}, "ç"},
		{"Alt-Ctrl-Minus", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyMinus}, "¥"},
		{"Alt-Ctrl-Slash", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeySlash}, "¿"},
		{"Alt-Ctrl-LBrace", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLBrace}, "«"},
		{"Alt-Ctrl-Backslash", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyBackslash}, "¬"},
		{"Alt-Ctrl-RBrace", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyRBrace}, "»"},
		{"Alt-Ctrl-Quote", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyQuote}, "´"},

		{"S-Alt-Ctrl-0", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.Key0}, ""},
		{"S-Alt-Ctrl-1", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.Key1}, "¹"},
		{"S-Alt-Ctrl-2", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.Key2}, ""},
		{"S-Alt-Ctrl-3", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.Key3}, ""},
		{"S-Alt-Ctrl-4", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.Key4}, "£"},
		{"S-Alt-Ctrl-5", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.Key5}, ""},
		{"S-Alt-Ctrl-6", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.Key6}, ""},
		{"S-Alt-Ctrl-7", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.Key7}, ""},
		{"S-Alt-Ctrl-8", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.Key8}, ""},
		{"S-Alt-Ctrl-9", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.Key9}, ""},
		{"S-Alt-Ctrl-A", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyA}, "Á"},
		{"S-Alt-Ctrl-B", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyB}, ""},
		{"S-Alt-Ctrl-C", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyC}, "¢"},
		{"S-Alt-Ctrl-D", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyD}, "Ð"},
		{"S-Alt-Ctrl-E", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyE}, "É"},
		{"S-Alt-Ctrl-F", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyF}, ""},
		{"S-Alt-Ctrl-G", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyG}, ""},
		{"S-Alt-Ctrl-H", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyH}, ""},
		{"S-Alt-Ctrl-I", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyI}, "Í"},
		{"S-Alt-Ctrl-J", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyJ}, ""},
		{"S-Alt-Ctrl-K", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyK}, ""},
		{"S-Alt-Ctrl-L", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyL}, "Ø"},
		{"S-Alt-Ctrl-M", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyM}, ""},
		{"S-Alt-Ctrl-N", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyN}, "Ñ"},
		{"S-Alt-Ctrl-O", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyO}, "Ó"},
		{"S-Alt-Ctrl-P", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyP}, "Ö"},
		{"S-Alt-Ctrl-Q", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyQ}, "Ä"},
		{"S-Alt-Ctrl-R", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyR}, ""},
		{"S-Alt-Ctrl-S", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyS}, "§"},
		{"S-Alt-Ctrl-T", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyT}, "Þ"},
		{"S-Alt-Ctrl-U", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyU}, "Ú"},
		{"S-Alt-Ctrl-V", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyV}, ""},
		{"S-Alt-Ctrl-W", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyW}, "Å"},
		{"S-Alt-Ctrl-X", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyX}, ""},
		{"S-Alt-Ctrl-Y", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyY}, "Ü"},
		{"S-Alt-Ctrl-Z", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyZ}, "Æ"},
		{"S-Alt-Ctrl-Semi", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeySemi}, "°"},
		{"S-Alt-Ctrl-Equal", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyEqual}, "÷"},
		{"S-Alt-Ctrl-Comma", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyComma}, "Ç"},
		{"S-Alt-Ctrl-Minus", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyMinus}, ""},
		{"S-Alt-Ctrl-Slash", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeySlash}, ""},
		{"S-Alt-Ctrl-LBrace", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyLBrace}, ""},
		{"S-Alt-Ctrl-RBrace", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyRBrace}, ""},
		{"S-Alt-Ctrl-BackSlash", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyBackslash}, "¦"},
		{"S-Alt-Ctrl-Quote", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyQuote}, "¨"},
		{"S-Alt-Ctrl-IsoBackSlash", []vt.Key{vt.KeyLAlt, vt.KeyLCtrl, vt.KeyLShift, vt.KeyIsoBackSlash}, ""},
	}
	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			term.KeyTap(cases[i].keys...)
			term.KeyTap(vt.KeySpace) // add a vanilla space to force output if empty
			CheckRead(t, term, cases[i].want+" ")
		})
	}
}

func TestDeadKeys(t *testing.T) {

	term := vt.NewMockTerm(vt.MockOptSize{X: 8, Y: 4})
	defer MustClose(t, term)
	term.SetLayout(vt.GetLayout("US International"))
	MustStart(t, term)

	term.KeyTap(vt.KeyGrave)
	term.KeyTap(vt.KeyA)
	want := "à"
	result := ReadF(t, term)

	VerifyF(t, want == result, "wrong read %q != %q", result, want)

	cases := []struct {
		name string
		keys []vt.Key
		key  []vt.Key
		want string
	}{
		{"Carat", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeySpace}, "^"},
		{"Carat-a", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeyA}, "â"},
		{"Carat-e", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeyE}, "ê"},
		{"Carat-i", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeyI}, "î"},
		{"Carat-o", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeyO}, "ô"},
		{"Carat-u", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeyU}, "û"},
		{"Carat-A", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeyLShift, vt.KeyA}, "Â"},
		{"Carat-E", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeyLShift, vt.KeyE}, "Ê"},
		{"Carat-I", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeyLShift, vt.KeyI}, "Î"},
		{"Carat-O", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeyLShift, vt.KeyO}, "Ô"},
		{"Carat-U", []vt.Key{vt.KeyLShift, vt.Key6}, []vt.Key{vt.KeyLShift, vt.KeyU}, "Û"},

		{"Acute", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeySpace}, "'"},
		{"Acute-a", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyA}, "á"},
		{"Acute-c", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyC}, "ç"},
		{"Acute-e", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyE}, "é"},
		{"Acute-i", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyI}, "í"},
		{"Acute-o", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyO}, "ó"},
		{"Acute-u", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyU}, "ú"},
		{"Acute-y", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyY}, "ý"},
		{"Acute-A", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyA}, "Á"},
		{"Acute-C", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyC}, "Ç"},
		{"Acute-E", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyE}, "É"},
		{"Acute-I", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyI}, "Í"},
		{"Acute-O", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyO}, "Ó"},
		{"Acute-U", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyU}, "Ú"},
		{"Acute-Y", []vt.Key{vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyY}, "Ý"},

		{"Quotes", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeySpace}, "\""},
		{"Quotes-a", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyA}, "ä"},
		{"Quotes-e", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyE}, "ë"},
		{"Quotes-i", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyI}, "ï"},
		{"Quotes-o", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyO}, "ö"},
		{"Quotes-u", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyU}, "ü"},
		{"Quotes-y", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyY}, "ÿ"},
		{"Quotes-A", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyA}, "Ä"},
		{"Quotes-E", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyE}, "Ë"},
		{"Quotes-I", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyI}, "Ï"},
		{"Quotes-O", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyO}, "Ö"},
		{"Quotes-U", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyU}, "Ü"},
		{"Quotes-Y", []vt.Key{vt.KeyLShift, vt.KeyQuote}, []vt.Key{vt.KeyLShift, vt.KeyY}, "Ÿ"},

		{"Grave", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeySpace}, "`"},
		{"Grave-a", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeyA}, "à"},
		{"Grave-e", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeyE}, "è"},
		{"Grave-i", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeyI}, "ì"},
		{"Grave-o", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeyO}, "ò"},
		{"Grave-u", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeyU}, "ù"},
		{"Grave-A", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeyLShift, vt.KeyA}, "À"},
		{"Grave-E", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeyLShift, vt.KeyE}, "È"},
		{"Grave-I", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeyLShift, vt.KeyI}, "Ì"},
		{"Grave-O", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeyLShift, vt.KeyO}, "Ò"},
		{"Grave-U", []vt.Key{vt.KeyGrave}, []vt.Key{vt.KeyLShift, vt.KeyU}, "Ù"},

		{"Tilde-a", []vt.Key{vt.KeyLShift, vt.KeyGrave}, []vt.Key{vt.KeyA}, "ã"},
		{"Tilde-n", []vt.Key{vt.KeyLShift, vt.KeyGrave}, []vt.Key{vt.KeyN}, "ñ"},
		{"Tilde-o", []vt.Key{vt.KeyLShift, vt.KeyGrave}, []vt.Key{vt.KeyO}, "õ"},
		{"Tilde-A", []vt.Key{vt.KeyLShift, vt.KeyGrave}, []vt.Key{vt.KeyLShift, vt.KeyA}, "Ã"},
		{"Tilde-N", []vt.Key{vt.KeyLShift, vt.KeyGrave}, []vt.Key{vt.KeyLShift, vt.KeyN}, "Ñ"},
		{"Tilde-O", []vt.Key{vt.KeyLShift, vt.KeyGrave}, []vt.Key{vt.KeyLShift, vt.KeyO}, "Õ"},
	}
	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) {
			term.KeyTap(cases[i].keys...)
			term.KeyTap(cases[i].key...)
			term.KeyTap(vt.KeySpace) // add a vanilla space to force output if empty
			CheckRead(t, term, cases[i].want+" ")
		})
	}

}
