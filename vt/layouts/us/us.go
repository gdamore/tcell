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

package us

import "github.com/gdamore/tcell/v3/vt"

const (
	UsInternationalLayout = "US International" // name of the US International
)
const (
	deadCarat = vt.DeadRune + '^'
	deadAcute = vt.DeadRune + '\''
	deadQuote = vt.DeadRune + '"'
	deadGrave = vt.DeadRune + '`'
	deadTilde = vt.DeadRune + '~'
)

var usInternationalUpper = map[vt.Key]rune{
	vt.KeyA:     'Á',
	vt.KeyD:     'Ð',
	vt.KeyE:     'É',
	vt.KeyI:     'Í',
	vt.KeyL:     'Ø',
	vt.KeyN:     'Ñ',
	vt.KeyO:     'Ó',
	vt.KeyP:     'Ö',
	vt.KeyQ:     'Ä',
	vt.KeyT:     'Þ',
	vt.KeyU:     'Ú',
	vt.KeyW:     'Å',
	vt.KeyY:     'Ü',
	vt.KeyZ:     'Æ',
	vt.KeyComma: 'Ç',
}

var usInternationalLower = map[vt.Key]rune{
	vt.KeyA:     'á',
	vt.KeyD:     'ð',
	vt.KeyE:     'é',
	vt.KeyI:     'í',
	vt.KeyL:     'ø',
	vt.KeyN:     'ñ',
	vt.KeyO:     'ó',
	vt.KeyP:     'ö',
	vt.KeyQ:     'ä',
	vt.KeyT:     'þ',
	vt.KeyU:     'ú',
	vt.KeyW:     'å',
	vt.KeyY:     'ü',
	vt.KeyZ:     'æ',
	vt.KeyComma: 'ç',
}

// UsInternational keyboard layout, includes Alt-Gr support (as Alt-Ctrl).
// Keys that are not overridden here fallback to the standard US Layout.
var UsInternational = &vt.Layout{
	Name:    UsInternationalLayout,
	Base:    vt.KeyboardANSI,
	Locking: vt.KeyboardANSI.Locking,
	Modifiers: map[vt.Key]vt.Modifier{
		vt.KeyLShift: vt.ModShift,
		vt.KeyRShift: vt.ModShift,
		vt.KeyLCtrl:  vt.ModCtrl,
		vt.KeyRCtrl:  vt.ModCtrl,
		vt.KeyRAlt:   vt.ModAlt | vt.ModCtrl, // acts as AltGr
		vt.KeyLAlt:   vt.ModAlt,
		vt.KeyRMeta:  vt.ModMeta,
		vt.KeyLMeta:  vt.ModMeta,
		vt.KeyRHyper: vt.ModHyper,
		vt.KeyLHyper: vt.ModHyper,
	},
	Maps: []vt.ModifierMap{
		{
			Mod:  vt.ModNone,
			Mask: vt.ModShift | vt.ModCtrl | vt.ModAlt,
			Map: map[vt.Key]rune{
				vt.KeyGrave: deadGrave,
				vt.KeyQuote: deadAcute,
			},
		},
		{
			Mod:  vt.ModShift,
			Mask: vt.ModShift | vt.ModCtrl | vt.ModAlt,
			Map: map[vt.Key]rune{
				vt.Key6:     deadCarat,
				vt.KeyGrave: deadTilde,
				vt.KeyQuote: deadQuote,
			},
		},
		// Letters
		{
			Mod:  vt.ModCtrl | vt.ModAlt,
			Mask: vt.ModCtrl | vt.ModAlt | vt.ModShift | vt.ModCapsLock,
			Map:  usInternationalLower,
		},
		{
			Mod:  vt.ModCtrl | vt.ModAlt | vt.ModShift,
			Mask: vt.ModCtrl | vt.ModAlt | vt.ModShift | vt.ModCapsLock,
			Map:  usInternationalUpper,
		},
		{
			Mod:  vt.ModCtrl | vt.ModAlt | vt.ModCapsLock,
			Mask: vt.ModCtrl | vt.ModAlt | vt.ModShift | vt.ModCapsLock,
			Map:  usInternationalUpper,
		},
		{
			Mod:  vt.ModCtrl | vt.ModAlt | vt.ModShift | vt.ModCapsLock,
			Mask: vt.ModCtrl | vt.ModAlt | vt.ModShift | vt.ModCapsLock,
			Map:  usInternationalLower,
		},
		// Extended non-letter forms
		{
			Mod:  vt.ModCtrl | vt.ModAlt,
			Mask: vt.ModCtrl | vt.ModAlt | vt.ModShift,
			Map: map[vt.Key]rune{
				vt.Key0:         '’',
				vt.Key1:         '¡',
				vt.Key2:         '²',
				vt.Key3:         '³',
				vt.Key4:         '¤',
				vt.Key5:         '€',
				vt.Key6:         '¼',
				vt.Key7:         '½',
				vt.Key8:         '¾',
				vt.Key9:         '‘',
				vt.KeyC:         '©',
				vt.KeyM:         'µ',
				vt.KeyR:         '®',
				vt.KeyS:         'ß', // technically also a letter, but no capital form
				vt.KeySemi:      '¶',
				vt.KeyEqual:     '×',
				vt.KeyMinus:     '¥',
				vt.KeySlash:     '¿',
				vt.KeyLBrace:    '«',
				vt.KeyBackslash: '¬',
				vt.KeyRBrace:    '»',
				vt.KeyQuote:     '´',
			},
		},
		// Even more extended non-letter forms
		{
			Mod:  vt.ModCtrl | vt.ModAlt | vt.ModShift,
			Mask: vt.ModCtrl | vt.ModAlt | vt.ModShift,
			Map: map[vt.Key]rune{
				vt.Key1:         '¹',
				vt.Key4:         '£',
				vt.KeyC:         '¢',
				vt.KeyS:         '§',
				vt.KeySemi:      '°',
				vt.KeyEqual:     '÷',
				vt.KeyBackslash: '¦',
				vt.KeyQuote:     '¨',
			},
		},
	},
	DeadKeys: map[rune]vt.DeadKey{
		deadCarat: {
			Next: map[rune]vt.DeadKey{
				' ': {U: '^'},
				'a': {U: 'â'},
				'A': {U: 'Â'},
				'e': {U: 'ê'},
				'E': {U: 'Ê'},
				'i': {U: 'î'},
				'I': {U: 'Î'},
				'o': {U: 'ô'},
				'O': {U: 'Ô'},
				'u': {U: 'û'},
				'U': {U: 'Û'},
			},
		},
		deadAcute: {
			Next: map[rune]vt.DeadKey{
				' ': {U: '\''},
				'a': {U: 'á'},
				'A': {U: 'Á'},
				'c': {U: 'ç'},
				'C': {U: 'Ç'},
				'e': {U: 'é'},
				'E': {U: 'É'},
				'i': {U: 'í'},
				'I': {U: 'Í'},
				'o': {U: 'ó'},
				'O': {U: 'Ó'},
				'u': {U: 'ú'},
				'U': {U: 'Ú'},
				'y': {U: 'ý'},
				'Y': {U: 'Ý'},
			},
		},
		deadQuote: {
			Next: map[rune]vt.DeadKey{
				' ': {U: '"'},
				'a': {U: 'ä'},
				'A': {U: 'Ä'},
				'e': {U: 'ë'},
				'E': {U: 'Ë'},
				'i': {U: 'ï'},
				'I': {U: 'Ï'},
				'o': {U: 'ö'},
				'O': {U: 'Ö'},
				'u': {U: 'ü'},
				'U': {U: 'Ü'},
				'y': {U: 'ÿ'},
				'Y': {U: 'Ÿ'},
			},
		},
		deadGrave: {
			Next: map[rune]vt.DeadKey{
				' ': {U: '`'},
				'a': {U: 'à'},
				'A': {U: 'À'},
				'e': {U: 'è'},
				'E': {U: 'È'},
				'i': {U: 'ì'},
				'I': {U: 'Ì'},
				'o': {U: 'ò'},
				'O': {U: 'Ò'},
				'u': {U: 'ù'},
				'U': {U: 'Ù'},
			},
		},
		deadTilde: {
			Next: map[rune]vt.DeadKey{
				' ': {U: '~'},
				'a': {U: 'ã'},
				'A': {U: 'Ã'},
				'n': {U: 'ñ'},
				'N': {U: 'Ñ'},
				'o': {U: 'õ'},
				'O': {U: 'Õ'},
			},
		},
	},
}

func init() {
	vt.RegisterLayout(UsInternational)
}
