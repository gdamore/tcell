// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

import "testing"

func TestEventKeyNameSidedModifiers(t *testing.T) {
	tests := []struct {
		key  Key
		mod  ModMask
		want string
	}{
		{KeyShift, ModLShift, "LeftShift"},
		{KeyShift, ModRShift, "RightShift"},
		{KeyCtrl, ModLCtrl, "LeftCtrl"},
		{KeyCtrl, ModRCtrl, "RightCtrl"},
		{KeyAlt, ModLAlt, "LeftAlt"},
		{KeyAlt, ModRAlt, "RightAlt"},
		{KeyMeta, ModLMeta, "LeftMeta"},
		{KeyMeta, ModRMeta, "RightMeta"},
		{KeyHyper, ModLHyper, "LeftHyper"},
		{KeyHyper, ModRHyper, "RightHyper"},
		{KeyShift, ModShift, "Shift"},
		{KeyCtrl, ModCtrl, "Ctrl"},
		{KeyAlt, ModAlt, "Alt"},
		{KeyMeta, ModMeta, "Meta"},
		{KeyHyper, ModHyper, "Hyper"},
	}

	for _, tt := range tests {
		ev := NewEventKeyEx(tt.key, "", tt.mod, true, tt.key, 1)
		if got := ev.Name(); got != tt.want {
			t.Errorf("Name() = %q, want %q", got, tt.want)
		}
	}
}

func TestEventKeyAdvancedNormalization(t *testing.T) {
	tests := []struct {
		name         string
		ev           *EventKey
		wantKey      Key
		wantStr      string
		wantMod      ModMask
		wantPressed  bool
		wantPhysical Key
		wantRepeat   int
	}{
		{
			name:         "advanced ctrl letter",
			ev:           NewEventKeyEx(KeyRune, "a", ModCtrl, false, KeyA, 0),
			wantKey:      KeyRune,
			wantStr:      "a",
			wantMod:      ModCtrl,
			wantPressed:  false,
			wantPhysical: KeyA,
			wantRepeat:   1,
		},
		{
			name:         "advanced shift tab",
			ev:           NewEventKeyEx(KeyTab, "", ModShift, true, KeyTab, 2),
			wantKey:      KeyTab,
			wantStr:      "",
			wantMod:      ModShift,
			wantPressed:  true,
			wantPhysical: KeyTab,
			wantRepeat:   2,
		},
		{
			name:         "backspace alias",
			ev:           NewEventKeyEx(KeyBackspace2, "", ModNone, true, KeyBackspace2, 1),
			wantKey:      KeyBackspace,
			wantStr:      "",
			wantMod:      ModNone,
			wantPressed:  true,
			wantPhysical: KeyBackspace2,
			wantRepeat:   1,
		},
		{
			name:         "legacy shift printable normalizes",
			ev:           NewEventKey(KeyRune, "A", ModShift),
			wantKey:      KeyRune,
			wantStr:      "A",
			wantMod:      ModNone,
			wantPressed:  true,
			wantPhysical: 0,
			wantRepeat:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.ev.Key() != tt.wantKey || tt.ev.Str() != tt.wantStr ||
				tt.ev.Modifiers() != tt.wantMod || tt.ev.Pressed() != tt.wantPressed ||
				tt.ev.Physical() != tt.wantPhysical || tt.ev.Repeat() != tt.wantRepeat {
				t.Fatalf("unexpected event: key=%v str=%q mod=%v pressed=%v physical=%v repeat=%v",
					tt.ev.Key(), tt.ev.Str(), tt.ev.Modifiers(), tt.ev.Pressed(), tt.ev.Physical(), tt.ev.Repeat())
			}
		})
	}
}

func TestPhysicalKeyAliases(t *testing.T) {
	if KeyA != Key('a') || KeyZ != Key('z') || Key0 != Key('0') || Key9 != Key('9') || KeySpace != Key(' ') {
		t.Fatalf("physical key aliases do not match their rune values")
	}
	tests := []struct {
		name string
		key  Key
		want Key
	}{
		{"grave", KeyGrave, Key('`')},
		{"minus", KeyMinus, Key('-')},
		{"equal", KeyEqual, Key('=')},
		{"left bracket", KeyLBracket, Key('[')},
		{"right bracket", KeyRBracket, Key(']')},
		{"backslash", KeyBackslash, Key('\\')},
		{"semicolon", KeySemicolon, Key(';')},
		{"quote", KeyQuote, Key('\'')},
		{"comma", KeyComma, Key(',')},
		{"period", KeyPeriod, Key('.')},
		{"slash", KeySlash, Key('/')},
	}
	for _, tt := range tests {
		if tt.key != tt.want {
			t.Fatalf("%s key alias = %v, want %v", tt.name, tt.key, tt.want)
		}
	}
}

func TestPhysicalKeySlashLogicalText(t *testing.T) {
	unshifted := NewEventKeyEx(KeyRune, "/", ModNone, true, KeySlash, 1)
	if unshifted.Str() != "/" || unshifted.Physical() != KeySlash {
		t.Fatalf("unshifted slash: str=%q physical=%v", unshifted.Str(), unshifted.Physical())
	}

	shifted := NewEventKeyEx(KeyRune, "?", ModShift, true, KeySlash, 1)
	if shifted.Str() != "?" || shifted.Physical() != KeySlash {
		t.Fatalf("shifted slash: str=%q physical=%v", shifted.Str(), shifted.Physical())
	}
}
