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

//go:build !js && !wasm
// +build !js,!wasm

package tcell

import (
	"fmt"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/vt"
)

// TestInputNullByte tests that null byte (0x00) is correctly handled
// as Ctrl+Space per the fix in the scan() function.
func TestInputNullByte(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)

	// Send null byte
	ip.ScanUTF8([]byte{0x00})

	// Wait briefly for processing
	select {
	case ev := <-evch:
		if kev, ok := ev.(*EventKey); ok {
			if kev.Key() != KeyRune {
				t.Errorf("Expected KeyRune for null byte, got %v", kev.Key())
			}
			if kev.Str() != " " {
				t.Errorf("Expected space character ' ' for null byte, got %q", kev.Str())
			}
			if kev.Modifiers() != ModCtrl {
				t.Errorf("Expected ModCtrl for null byte, got %v", kev.Modifiers())
			}
		} else {
			t.Errorf("Expected EventKey, got %T", ev)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for null byte event")
	}
}

// TestInputControlKeys tests control key handling for bytes 1-31.
// Note: NewEventKey converts control characters to KeyCtrlA-Z for 0x01-0x1A,
// and KeyRune for 0x1C-0x1F with the character and ModCtrl.
func TestInputControlKeys(t *testing.T) {
	tests := []struct {
		name  string
		input byte
		key   Key
		str   string
		mod   ModMask
	}{
		{"Ctrl+A", 0x01, KeyCtrlA, "", ModCtrl},
		{"Ctrl+B", 0x02, KeyCtrlB, "", ModCtrl},
		{"Ctrl+C", 0x03, KeyCtrlC, "", ModCtrl},
		{"Ctrl+D", 0x04, KeyCtrlD, "", ModCtrl},
		{"Ctrl+E", 0x05, KeyCtrlE, "", ModCtrl},
		{"Ctrl+F", 0x06, KeyCtrlF, "", ModCtrl},
		{"Ctrl+G", 0x07, KeyCtrlG, "", ModCtrl},
		{"Ctrl+H", 0x08, KeyBackspace, "", ModNone}, // BS is special
		{"Ctrl+I", 0x09, KeyTab, "", ModNone},       // Tab is special
		{"Ctrl+J", 0x0A, KeyCtrlJ, "", ModCtrl},
		{"Ctrl+K", 0x0B, KeyCtrlK, "", ModCtrl},
		{"Ctrl+L", 0x0C, KeyCtrlL, "", ModCtrl},
		{"Ctrl+M", 0x0D, KeyEnter, "", ModNone}, // CR is Enter
		{"Ctrl+N", 0x0E, KeyCtrlN, "", ModCtrl},
		{"Ctrl+O", 0x0F, KeyCtrlO, "", ModCtrl},
		{"Ctrl+P", 0x10, KeyCtrlP, "", ModCtrl},
		{"Ctrl+Q", 0x11, KeyCtrlQ, "", ModCtrl},
		{"Ctrl+R", 0x12, KeyCtrlR, "", ModCtrl},
		{"Ctrl+S", 0x13, KeyCtrlS, "", ModCtrl},
		{"Ctrl+T", 0x14, KeyCtrlT, "", ModCtrl},
		{"Ctrl+U", 0x15, KeyCtrlU, "", ModCtrl},
		{"Ctrl+V", 0x16, KeyCtrlV, "", ModCtrl},
		{"Ctrl+W", 0x17, KeyCtrlW, "", ModCtrl},
		{"Ctrl+X", 0x18, KeyCtrlX, "", ModCtrl},
		{"Ctrl+Y", 0x19, KeyCtrlY, "", ModCtrl},
		{"Ctrl+Z", 0x1A, KeyCtrlZ, "", ModCtrl},
		{"Ctrl+[", 0x1B, KeyEscape, "", ModNone},  // ESC is special
		{"Ctrl+\\", 0x1C, KeyRune, "\\", ModCtrl}, // becomes KeyRune with string
		{"Ctrl+]", 0x1D, KeyRune, "]", ModCtrl},   // becomes KeyRune with string
		{"Ctrl+^", 0x1E, KeyRune, "^", ModCtrl},   // becomes KeyRune with string
		{"Ctrl+_", 0x1F, KeyRune, "_", ModCtrl},   // becomes KeyRune with string
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputParser(evch)

			// Skip ESC (0x1B) as it has special timeout handling
			if tt.input == 0x1B {
				t.Skip("ESC byte triggers escape sequence state machine with timeout")
				return
			}

			ip.ScanUTF8([]byte{tt.input})

			select {
			case ev := <-evch:
				if kev, ok := ev.(*EventKey); ok {
					if kev.Key() != tt.key {
						t.Errorf("Expected key %v, got %v", tt.key, kev.Key())
					}
					if kev.Str() != tt.str {
						t.Errorf("Expected string %q, got %q", tt.str, kev.Str())
					}
					if kev.Modifiers() != tt.mod {
						t.Errorf("Expected modifiers %v, got %v", tt.mod, kev.Modifiers())
					}
				} else {
					t.Errorf("Expected EventKey, got %T", ev)
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Timeout waiting for control key event")
			}
		})
	}
}

// TestInputNullVsOtherControlChars tests the boundary between
// null byte and other control characters.
func TestInputNullVsOtherControlChars(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)

	// Test null (0x00) - should be KeyRune with Ctrl+Space
	ip.ScanUTF8([]byte{0x00})
	select {
	case ev := <-evch:
		if kev, ok := ev.(*EventKey); ok {
			if kev.Key() != KeyRune {
				t.Errorf("Null byte: expected KeyRune, got %v", kev.Key())
			}
			if kev.Str() != " " {
				t.Errorf("Null byte: expected ' ', got %q", kev.Str())
			}
			if kev.Modifiers() != ModCtrl {
				t.Errorf("Null byte: expected ModCtrl, got %v", kev.Modifiers())
			}
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout for null byte")
	}

	// Test 0x01 (Ctrl+A) - should be KeyCtrlA
	ip.ScanUTF8([]byte{0x01})
	select {
	case ev := <-evch:
		if kev, ok := ev.(*EventKey); ok {
			if kev.Key() != KeyCtrlA {
				t.Errorf("Byte 0x01: expected KeyCtrlA, got %v", kev.Key())
			}
			if kev.Modifiers() != ModCtrl {
				t.Errorf("Byte 0x01: expected ModCtrl, got %v", kev.Modifiers())
			}
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout for 0x01")
	}
}

// TestInputPrintableCharacters tests that printable characters
// are handled correctly without control modifiers.
func TestInputPrintableCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		mod      ModMask
	}{
		{"Space", " ", " ", ModNone},
		{"ExclamationMark", "!", "!", ModNone},
		{"LetterA", "A", "A", ModNone},
		{"LetterZ", "Z", "Z", ModNone},
		{"Lowera", "a", "a", ModNone},
		{"Lowerz", "z", "z", ModNone},
		{"Digit0", "0", "0", ModNone},
		{"Digit9", "9", "9", ModNone},
		{"At", "@", "@", ModNone},
		{"Tilde", "~", "~", ModNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputParser(evch)

			ip.ScanUTF8([]byte(tt.input))

			select {
			case ev := <-evch:
				if kev, ok := ev.(*EventKey); ok {
					if kev.Key() != KeyRune {
						t.Errorf("Expected KeyRune, got %v", kev.Key())
					}
					if kev.Str() != tt.expected {
						t.Errorf("Expected %q, got %q", tt.expected, kev.Str())
					}
					if kev.Modifiers() != tt.mod {
						t.Errorf("Expected modifiers %v, got %v", tt.mod, kev.Modifiers())
					}
				} else {
					t.Errorf("Expected EventKey, got %T", ev)
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Timeout waiting for character event")
			}
		})
	}
}

// TestInputSpecialKeys tests special key handling (tab, backspace, enter).
func TestInputSpecialKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected Key
		mod      ModMask
	}{
		{"Tab", '\t', KeyTab, ModNone},
		{"Backspace_BS", '\b', KeyBackspace, ModNone},
		{"Backspace_DEL", 0x7F, KeyBackspace, ModNone},
		{"Enter_LF", '\n', KeyCtrlJ, ModCtrl},
		{"Enter_CR", '\r', KeyEnter, ModNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputParser(evch)

			ip.ScanUTF8([]byte{tt.input})

			select {
			case ev := <-evch:
				if kev, ok := ev.(*EventKey); ok {
					if kev.Key() != tt.expected {
						t.Errorf("Expected key %v, got %v", tt.expected, kev.Key())
					}
					if kev.Modifiers() != tt.mod {
						t.Errorf("Expected modifiers %v, got %v", tt.mod, kev.Modifiers())
					}
				} else {
					t.Errorf("Expected EventKey, got %T", ev)
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Timeout waiting for special key event")
			}
		})
	}
}

// TestInputSequentialInput tests handling multiple inputs in sequence.
func TestInputSequentialInput(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)

	// Send: null, Ctrl+A, 'B', space
	inputs := []byte{0x00, 0x01, 'B', ' '}
	expected := []struct {
		key Key
		str string
		mod ModMask
	}{
		{KeyRune, " ", ModCtrl}, // null -> Ctrl+Space
		{KeyCtrlA, "", ModCtrl}, // 0x01 -> KeyCtrlA with ModCtrl
		{KeyRune, "B", ModNone}, // 'B'
		{KeyRune, " ", ModNone}, // space
	}

	for i, b := range inputs {
		ip.ScanUTF8([]byte{b})

		select {
		case ev := <-evch:
			if kev, ok := ev.(*EventKey); ok {
				if kev.Key() != expected[i].key {
					t.Errorf("Input %d: expected key %v, got %v", i, expected[i].key, kev.Key())
				}
				if kev.Str() != expected[i].str {
					t.Errorf("Input %d: expected %q, got %q", i, expected[i].str, kev.Str())
				}
				if kev.Modifiers() != expected[i].mod {
					t.Errorf("Input %d: expected modifiers %v, got %v", i, expected[i].mod, kev.Modifiers())
				}
			} else {
				t.Errorf("Input %d: expected EventKey, got %T", i, ev)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Timeout waiting for event %d", i)
		}
	}
}

// TestInputUTF8Characters tests UTF-8 multibyte character handling.
func TestInputUTF8Characters(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"Euro", []byte("€"), "€"},
		{"CJK", []byte("中"), "中"},
		{"Emoji", []byte("😀"), "😀"},
		{"Cyrillic", []byte("Ж"), "Ж"},
		{"Arabic", []byte("ع"), "ع"},
		{"SMP", []byte("🝁"), "🝁"}, // needs full 4 character UTF-8
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputParser(evch)

			ip.ScanUTF8(tt.input)
			time.Sleep(time.Millisecond * 100)
			ip.Scan()

			select {
			case ev := <-evch:
				if kev, ok := ev.(*EventKey); ok {
					if kev.Str() != tt.expected {
						t.Errorf("Expected %q, got %q", tt.expected, kev.Str())
					}
					if kev.Modifiers() != ModNone {
						t.Errorf("Expected ModNone, got %v", kev.Modifiers())
					}
				} else {
					t.Errorf("Expected EventKey, got %T", ev)
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Timeout waiting for UTF-8 character event")
			}
		})
	}
}

// TestInputEdgeCases tests edge cases and boundary conditions.
func TestInputEdgeCases(t *testing.T) {
	t.Run("EmptyInput", func(t *testing.T) {
		evch := make(chan Event, 10)
		ip := newInputParser(evch)

		ip.ScanUTF8([]byte{})
		time.Sleep(time.Millisecond * 100)
		ip.Scan()

		select {
		case ev := <-evch:
			t.Errorf("Expected no event for empty input, got %T", ev)
		case <-time.After(50 * time.Millisecond):
			// Success - no event generated
		}
	})

	t.Run("MultipleNullBytes", func(t *testing.T) {
		evch := make(chan Event, 10)
		ip := newInputParser(evch)

		ip.ScanUTF8([]byte{0x00, 0x00, 0x00})

		for i := range 3 {
			select {
			case ev := <-evch:
				if kev, ok := ev.(*EventKey); ok {
					if kev.Key() != KeyRune || kev.Str() != " " || kev.Modifiers() != ModCtrl {
						t.Errorf("Null byte %d: expected KeyRune with Ctrl+Space, got key=%v str=%q mod=%v",
							i, kev.Key(), kev.Str(), kev.Modifiers())
					}
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("Timeout waiting for null byte %d", i)
			}
		}
	})

	t.Run("BoundaryByte0x1F", func(t *testing.T) {
		evch := make(chan Event, 10)
		ip := newInputParser(evch)

		// 0x1F is the last control character before space (0x20)
		// NewEventKey transforms it to KeyRune with "_" and ModCtrl
		ip.ScanUTF8([]byte{0x1F})

		select {
		case ev := <-evch:
			if kev, ok := ev.(*EventKey); ok {
				if kev.Key() != KeyRune {
					t.Errorf("Expected KeyRune for 0x1F, got %v", kev.Key())
				}
				if kev.Str() != "_" {
					t.Errorf("Expected '_' for 0x1F, got %q", kev.Str())
				}
				if kev.Modifiers() != ModCtrl {
					t.Errorf("Expected ModCtrl for 0x1F, got %v", kev.Modifiers())
				}
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timeout for 0x1F")
		}
	})

	t.Run("BoundaryByte0x20", func(t *testing.T) {
		evch := make(chan Event, 10)
		ip := newInputParser(evch)

		// 0x20 is space - first printable character
		ip.ScanUTF8([]byte{0x20})

		select {
		case ev := <-evch:
			if kev, ok := ev.(*EventKey); ok {
				if kev.Key() != KeyRune {
					t.Errorf("Expected KeyRune for 0x20, got %v", kev.Key())
				}
				if kev.Str() != " " {
					t.Errorf("Expected ' ' for 0x20, got %q", kev.Str())
				}
				if kev.Modifiers() != ModNone {
					t.Errorf("Expected ModNone for 0x20, got %v", kev.Modifiers())
				}
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timeout for 0x20")
		}
	})
}

// TestInputConcurrentAccess tests that the parser is safe
// for concurrent access (it uses a mutex internally).
func TestInputConcurrentAccess(t *testing.T) {
	evch := make(chan Event, 100)
	ip := newInputParser(evch)

	done := make(chan bool)

	// Goroutine 1: Send null bytes
	go func() {
		for range 10 {
			ip.ScanUTF8([]byte{0x00})
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Goroutine 2: Send regular characters
	go func() {
		for range 10 {
			ip.ScanUTF8([]byte{'A'})
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Verify we got events (exact order not guaranteed)
	eventCount := 0
	timeout := time.After(200 * time.Millisecond)
	for eventCount < 20 {
		select {
		case <-evch:
			eventCount++
		case <-timeout:
			t.Fatalf("Expected 20 events, got %d", eventCount)
		}
	}
}

// TestInputStateTransitions tests that the null byte fix doesn't
// interfere with state machine transitions.
func TestInputStateTransitions(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)

	// Send null in initial state
	ip.ScanUTF8([]byte{0x00})

	select {
	case ev := <-evch:
		if kev, ok := ev.(*EventKey); ok {
			if kev.Key() != KeyRune || kev.Str() != " " || kev.Modifiers() != ModCtrl {
				t.Errorf("Expected KeyRune with Ctrl+Space, got key=%v str=%q mod=%v",
					kev.Key(), kev.Str(), kev.Modifiers())
			}
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout")
	}

	// Verify state machine returns to initial state by sending a normal char
	ip.ScanUTF8([]byte{'X'})

	select {
	case ev := <-evch:
		if kev, ok := ev.(*EventKey); ok {
			if kev.Key() != KeyRune || kev.Str() != "X" || kev.Modifiers() != ModNone {
				t.Errorf("Expected KeyRune 'X' with ModNone after null, got key=%v str=%q mod=%v",
					kev.Key(), kev.Str(), kev.Modifiers())
			}
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout")
	}
}

// TestSpecialKeys tests that special keys (F-keys, home, del, etc. work as expected)
func TestSpecialKeys(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectedKey Key
		expectedMod ModMask
		expectedStr string
	}{
		{"Esc", []byte{'\x1b'}, KeyEscape, ModNone, ""},
		{"Esc-Esc", []byte{'\x1b', '\x1b'}, KeyEscape, ModAlt, ""},
		{"Esc-Y", []byte{'\x1b', 'Y'}, KeyRune, ModAlt, "Y"},
		{"Esc-Ctrl-B", []byte{'\x1b', '\x02'}, KeyRune, ModAlt | ModCtrl, "b"},
		{"Esc-[", []byte{'\x1b', '['}, KeyRune, ModAlt, "["},
		{"Tab", []byte{'\t'}, KeyTab, ModNone, ""},
		{"NL", []byte{'\n'}, KeyCtrlJ, ModCtrl, ""},
		{"CR", []byte{'\r'}, KeyEnter, ModNone, ""},
		{"Backspace", []byte{'\b'}, KeyBackspace, ModNone, ""},
		{"Delete", []byte{'\x7f'}, KeyBackspace, ModNone, ""},
		{"CSI-A", []byte{'\x1b', '[', 'A'}, KeyUp, ModNone, ""},
		{"CSI-B", []byte{'\x1b', '[', 'B'}, KeyDown, ModNone, ""},
		{"CSI-C", []byte{'\x1b', '[', 'C'}, KeyRight, ModNone, ""},
		{"CSI-D", []byte{'\x1b', '[', 'D'}, KeyLeft, ModNone, ""},
		{"CSI-E", []byte{'\x1b', '[', 'E'}, KeyClear, ModNone, ""},
		{"CSI-F", []byte{'\x1b', '[', 'F'}, KeyEnd, ModNone, ""},
		{"CSI-H", []byte{'\x1b', '[', 'H'}, KeyHome, ModNone, ""},
		{"CSI-L", []byte{'\x1b', '[', 'L'}, KeyInsert, ModNone, ""},
		{"CSI-P", []byte{'\x1b', '[', 'P'}, KeyF1, ModNone, ""},
		{"CSI-Q", []byte{'\x1b', '[', 'Q'}, KeyF2, ModNone, ""},
		{"CSI-S", []byte{'\x1b', '[', 'S'}, KeyF4, ModNone, ""},
		{"CSI-Z", []byte{'\x1b', '[', 'Z'}, KeyBacktab, ModNone, ""},
		{"CSI-a", []byte{'\x1b', '[', 'a'}, KeyUp, ModShift, ""},
		{"CSI-b", []byte{'\x1b', '[', 'b'}, KeyDown, ModShift, ""},
		{"CSI-c", []byte{'\x1b', '[', 'c'}, KeyRight, ModShift, ""},
		{"CSI-d", []byte{'\x1b', '[', 'd'}, KeyLeft, ModShift, ""},
		{"CSI-15~", []byte{'\x1b', '[', '1', '5', '~'}, KeyF5, ModNone, ""},
		{"CSI-17~", []byte{'\x1b', '[', '1', '7', '~'}, KeyF6, ModNone, ""},
		{"CSI-18~", []byte{'\x1b', '[', '1', '8', '~'}, KeyF7, ModNone, ""},
		{"CSI-19~", []byte{'\x1b', '[', '1', '9', '~'}, KeyF8, ModNone, ""},
		{"CSI-20~", []byte{'\x1b', '[', '2', '0', '~'}, KeyF9, ModNone, ""},
		{"CSI-21~", []byte{'\x1b', '[', '2', '1', '~'}, KeyF10, ModNone, ""},
		{"CSI-23~", []byte{'\x1b', '[', '2', '3', '~'}, KeyF11, ModNone, ""},
		{"CSI-24~", []byte{'\x1b', '[', '2', '4', '~'}, KeyF12, ModNone, ""},
		{"CSI-1-$", []byte{'\x1b', '[', '1', '$'}, KeyHome, ModShift, ""}, // rxvt bs
		{"SS3-F1", []byte{'\x1b', 'O', 'P'}, KeyF1, ModNone, ""},
		{"SS3-F2", []byte{'\x1b', 'O', 'Q'}, KeyF2, ModNone, ""},
		{"SS3-F3", []byte{'\x1b', 'O', 'R'}, KeyF3, ModNone, ""},
		{"SS3-F4", []byte{'\x1b', 'O', 'S'}, KeyF4, ModNone, ""},
		{"SS3-F4-Shift", []byte{'\x1b', 'O', '1', ';', '2', 'S'}, KeyF4, ModShift, ""},
		{"SS3-F4-Ctrl", []byte{'\x1b', 'O', '1', ';', '5', 'S'}, KeyF4, ModCtrl, ""},
		{"SS3-F4-Ctrl-Short", []byte{'\x1b', 'O', '5', 'S'}, KeyF4, ModCtrl, ""},
		{"SS3-F4-Ctrl-Shift", []byte{'\x1b', 'O', ';', '6', 'S'}, KeyF4, ModCtrl | ModShift, ""},
		{"SS3-F2-Meta", []byte{'\x1b', 'O', ';', '9', 'Q'}, KeyF2, ModMeta, ""},
		{"CSI-F2-Alt", []byte{'\x1b', '[', '1', ';', '3', 'Q'}, KeyF2, ModAlt, ""},
		{"CSI-F2-Hyper", []byte{'\x1b', '[', '1', ';', '1', '7', 'Q'}, KeyF2, ModHyper, ""},
		{"CSI-F2-Super", []byte{'\x1b', '[', '1', ';', '3', '3', 'Q'}, KeyF2, ModMeta, ""},
		{"Ctrl-Home", []byte{'\x1b', '[', '1', ';', '5', '~'}, KeyHome, ModCtrl, ""},
		{"SS3-Home", []byte{'\x1b', 'O', 'H'}, KeyHome, ModNone, ""},
		{"SS3-Clear", []byte{'\x1b', 'O', 'E'}, KeyClear, ModNone, ""},
		{"ESC-Tab", []byte{'\x1b', '\t'}, KeyBacktab, ModNone, ""}, // linux console special
		{"Linux-F1", []byte{'\x1b', '[', '[', 'A'}, KeyF1, ModNone, ""},
		{"Linux-F2", []byte{'\x1b', '[', '[', 'B'}, KeyF2, ModNone, ""},
		{"Linux-F3", []byte{'\x1b', '[', '[', 'C'}, KeyF3, ModNone, ""},
		{"Linux-F4", []byte{'\x1b', '[', '[', 'D'}, KeyF4, ModNone, ""},
		{"Linux-F5", []byte{'\x1b', '[', '[', 'E'}, KeyF5, ModNone, ""},
		{"XTerm-Alt-Tab", []byte{'\x1b', '[', '2', '7', ';', '3', ';', '9', '~'}, KeyTab, ModAlt, ""}, // modifyOtherKeys == 1
		{"Alt-F7", []byte{'\x1b', '[', '1', '8', ';', '3', '~'}, KeyF7, ModAlt, ""},
		{"XTerm-Shift-Tab", []byte{'\x1b', '[', '2', '7', ';', '2', ';', '9', '~'}, KeyBacktab, ModNone, ""}, // modifyOtherKeys == 2
		{"XTerm-Space", []byte{'\x1b', '[', '2', '7', ';', '1', ';', '3', '2', '~'}, KeyRune, ModNone, " "},  // modifyOtherKeys == 3
		{"Kitty-Esc", []byte{'\x1b', '[', '2', '7', 'u'}, KeyEsc, ModNone, ""},
		{"Kitty-Control-I", []byte{'\x1b', '[', '1', '0', '5', ';', '5', 'u'}, 'I', ModCtrl, ""},
		{"Win-Shift-A", []byte{'\x1b', '[', '6', '5', ';', '0', ';', '6', '5', ';', '1', ';', '1', '6', '_'}, KeyRune, ModNone, "A"},
		{"Win-Ctrl-1", []byte{'\x1b', '[', '4', '9', ';', '0', ';', '4', '9', ';', '1', ';', '8', '_'}, KeyRune, ModCtrl, "1"},
		{"Win-Ctrl-A", []byte{'\x1b', '[', '6', '5', ';', '0', ';', '1', ';', '1', ';', '8', '_'}, KeyCtrlA, ModCtrl, ""},
		{"Win-Ctrl-Up", []byte{'\x1b', '[', '3', '8', ';', '0', ';', '0', ';', '1', ';', '8', '_'}, KeyUp, ModCtrl, ""},
		{"Win-Ctrl-Up-2", []byte{'\x1b', '[', '3', '8', ';', '0', ';', '0', ';', '1', ';', '4', '_'}, KeyUp, ModCtrl, ""},
		{"Win-Alt-F1", []byte{'\x1b', '[', '1', '1', '2', ';', '0', ';', '0', ';', '1', ';', '1', '_'}, KeyF1, ModAlt, ""},
		{"Win-Alt-F1-2", []byte{'\x1b', '[', '1', '1', '2', ';', '0', ';', '0', ';', '1', ';', '2', '_'}, KeyF1, ModAlt, ""},
		{"Win-AltGr-E", []byte{'\x1b', '[', '6', '9', ';', '0', ';', '6', '9', ';', '1', ';', '5', '_'}, KeyRune, ModNone, "E"},
		{"Win-Ignore-Release", []byte{'\x1b', '[', '6', '5', ';', '0', ';', '6', '5', ';', '0', ';', '1', '6', '_', 'C'}, KeyRune, ModNone, "C"},
		{"Win-Mod-Ignore-Shift", []byte{'\x1b', '[', '1', '6', ';', '0', ';', '1', '1', ';', '1', ';', '1', '6', '_', 'C'}, KeyRune, ModNone, "C"},
		{"Win-Mod-Ignore-Ctrl", []byte{'\x1b', '[', '1', '7', ';', '0', ';', '1', '3', ';', '1', ';', '1', '6', '_', 'C'}, KeyRune, ModNone, "C"},
		{"Win-Mod-Ignore-Alt", []byte{'\x1b', '[', '1', '8', ';', '0', ';', '1', '4', ';', '1', ';', '1', '6', '_', 'C'}, KeyRune, ModNone, "C"},
		{"Win-Surrogates", []byte{
			'\x1b', '[', '0', ';', '0', ';', '5', '5', '3', '5', '6', ';', '1', ';', '0', ';', '0', '_',
			'\x1b', '[', '0', ';', '0', ';', '5', '7', '2', '5', '6', ';', '1', ';', '0', ';', '0', '_',
		}, KeyRune, ModNone, "🎨"},
		{"Win-Nested-Shift-B", []byte{
			'\x1b', '[', '0', ';', '0', ';', '2', '7', ';', '1', ';', '0', ';', '0', '_', // ESC
			'\x1b', '[', '0', ';', '0', ';', '9', '1', ';', '1', ';', '0', ';', '0', '_', // [
			'\x1b', '[', '0', ';', '0', ';', '5', '4', ';', '1', ';', '0', ';', '0', '_', // 6
			'\x1b', '[', '0', ';', '0', ';', '5', '4', ';', '1', ';', '0', ';', '0', '_', // 6
			'\x1b', '[', '0', ';', '0', ';', '5', '9', ';', '1', ';', '0', ';', '0', '_', // ;
			'\x1b', '[', '0', ';', '0', ';', '4', '8', ';', '1', ';', '0', ';', '0', '_', // 0
			'\x1b', '[', '0', ';', '0', ';', '5', '9', ';', '1', ';', '0', ';', '0', '_', // ;
			'\x1b', '[', '0', ';', '0', ';', '5', '4', ';', '1', ';', '0', ';', '0', '_', // 6
			'\x1b', '[', '0', ';', '0', ';', '5', '4', ';', '1', ';', '0', ';', '0', '_', // 6
			'\x1b', '[', '0', ';', '0', ';', '5', '9', ';', '1', ';', '0', ';', '0', '_', // ;
			'\x1b', '[', '0', ';', '0', ';', '4', '9', ';', '1', ';', '0', ';', '0', '_', // 1
			'\x1b', '[', '0', ';', '0', ';', '5', '9', ';', '1', ';', '0', ';', '0', '_', // ;
			'\x1b', '[', '0', ';', '0', ';', '4', '9', ';', '1', ';', '0', ';', '0', '_', // 1
			'\x1b', '[', '0', ';', '0', ';', '5', '4', ';', '1', ';', '0', ';', '0', '_', // 6
			'\x1b', '[', '0', ';', '0', ';', '5', '9', ';', '1', ';', '0', ';', '0', '_', // ;
			'\x1b', '[', '0', ';', '0', ';', '4', '9', ';', '1', ';', '0', ';', '0', '_', // 1
			'\x1b', '[', '0', ';', '0', ';', '9', '5', ';', '1', ';', '0', ';', '0', '_', // _
		}, KeyRune, ModNone, "B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputParser(evch)

			ip.ScanUTF8(tt.input)

			expected := NewEventKey(tt.expectedKey, "", tt.expectedMod)

			select {
			case ev := <-evch:
				if kev, ok := ev.(*EventKey); ok {
					if kev.Key() != tt.expectedKey || kev.Modifiers() != tt.expectedMod || kev.Str() != string(tt.expectedStr) {
						t.Errorf("Expected %q, got %q (rune expected %q got %q)", expected.Name(), kev.Name(), string(tt.expectedStr), kev.Str())
					} else {
						t.Logf("Key %s ok", kev.Name())
					}
				} else {
					t.Errorf("Expected EventKey, got %T", ev)
				}

			case <-time.After(100 * time.Millisecond):
				ip.Scan()
				select {
				case ev := <-evch:
					if kev, ok := ev.(*EventKey); ok {
						if kev.Key() != tt.expectedKey || kev.Modifiers() != tt.expectedMod {
							t.Errorf("Expected %q, got %q", expected.Name(), kev.Name())
						}
					} else {
						t.Errorf("Expected EventKey, got %T", ev)
					}

				default:
					t.Fatalf("Timeout waiting for key event")
				}
			}
		})
	}
}

// TestDecPrivateModeResponse tests responses to various DEC private mode queries
func TestDecPrivateModeResponse(t *testing.T) {
	tests := []struct {
		bytes  string
		result eventPrivateMode
		usable bool
	}{
		{"\x1b[?1001;0$y", eventPrivateMode{Mode: 1001, Status: vt.ModeNA}, false},
		{"\x1b[?1004;2$y", eventPrivateMode{Mode: 1004, Status: vt.ModeOff}, true},
		{"\x1b[?7;1$y", eventPrivateMode{Mode: vt.PmAutoMargin, Status: vt.ModeOn}, true},
		{"\x1b[?25;3$y", eventPrivateMode{Mode: vt.PmShowCursor, Status: vt.ModeOnLocked}, false},
		{"\x1b[?12;4$y", eventPrivateMode{Mode: vt.PmBlinkCursor, Status: vt.ModeOffLocked}, false},
		{"\x1b[?990;$y", eventPrivateMode{Mode: 990, Status: vt.ModeNA}, false},
		{"\x1b[?991$y", eventPrivateMode{Mode: 991, Status: vt.ModeNA}, false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputParser(evch)

			// Send null in initial state
			ip.ScanUTF8([]byte(tt.bytes))

			select {
			case ev := <-evch:
				if pev, ok := ev.(*eventPrivateMode); ok {
					if pev.Status.Changeable() != tt.usable {
						t.Errorf("Private mode usability wrong %v != %v", pev.Status.Changeable(), tt.usable)
					}
					if *pev != tt.result {
						t.Errorf("Private mode mismatch for %d", tt.result.Mode)
					}
				} else {
					t.Errorf("Got unexpected event %T", ev)
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Timeout")
			}
		})
	}
}

func TestIgnoredSequences(t *testing.T) {
	tests := []struct {
		name  string
		bytes string
	}{
		{"LoneST", "\x9c"}, // 7 bit version would be confused with Alt-\
		{"SoS", "\x1bXdata\x1b\\"},
		{"SoS-Bell", "\x1bXdata\x07"},
		{"SoS-Embed-ESC", "\x1bXab\x1bcde\x1b\\"},
		{"PM", "\x1b^data\x07"},
		{"PM8", "\x9edata\x07"},
		{"PM-Bell", "\x1b^data\x07"},
		{"APC", "\x1b_data\x07"},
		{"APC8", "\x9fdata\x07"},
		{"APC-Bell", "\x1b_data\x07"},
		{"OSC", "\x1b]junk\x1b\\"},
		{"OSC8", "\x9djunk\x1b\\"},
		{"OSC-Bell", "\x1b]junk\x07"},
		{"DCS", "\x1bPjunk\x1b\\"},
		{"DCS8", "\x90junk\x1b\\"},
		{"DCS-Bell", "\x1bPjunk\x07"},
		{"SS2", "\x1bN1"},
		{"SS28", "\x8e1"},
		{"BadCSI", "\x1b[\x07"},
		{"BadUTF8", "\xe0\xff"},
		{"Win32Shift", "\x1b[16;0;0;1;1;1_"},
		{"Win32Ctrl", "\x1b[17;0;0;1;1;1_"},
		{"Win32Alt", "\x1b[18;0;0;1;1;1_"},
		{"Win32CapsLock", "\x1b[20;0;0;1;1;1_"},
		{"Win32KeyUp", "\x1b[13;0;13;0;1;1_"},
		{"RuntDA1", "\x1b[?c"},
		{"RuntWindowNotice", "\x1b[t"},
		{"OtherIntermediates", "\x1b[1 ~"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputParser(evch)

			// send event plus DECID so we get a final event
			ip.ScanUTF8(append([]byte(test.bytes), '\x1B', '[', '?', '6', 'c'))
			select {
			case ev := <-evch:
				if _, ok := ev.(*eventPrimaryAttributes); ok {
					return
				} else {
					t.Errorf("Got unexpected event %T", ev)
					if ev, ok := ev.(*EventKey); ok {
						t.Logf("Key %s", ev.Name())
					}
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Timeout")
			}
		})
	}
}

// TestKeyboardMode tests the responses to various keyboard modes
func TestKeyboardMode(t *testing.T) {
	tests := []struct {
		bytes  string
		result Event
	}{
		{"\x1b[?9001;0$y", &eventPrivateMode{Mode: 9001, Status: vt.ModeNA}},
		{"\x1b[?9001;2$y", &eventPrivateMode{Mode: 9001, Status: vt.ModeOff}},
		{"\x1b[>4;0m", &eventXTermKbdMode{Mode: XtermKbdModeOff}},
		{"\x1b[?0u", &eventKittyKbdMode{Mode: KittyKbdModeOff}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputParser(evch)

			// Send null in initial state
			ip.ScanUTF8([]byte(tt.bytes))

			select {
			case ev := <-evch:
				switch expect := tt.result.(type) {
				case *eventPrivateMode:
					if actual, ok := ev.(*eventPrivateMode); ok {
						if actual.Mode != expect.Mode || actual.Status != expect.Status {
							t.Errorf("Wrong mode or status: %v %v != %v %v", actual.Mode, actual.Status, expect.Mode, expect.Status)
						}
						return
					}
				case *eventXTermKbdMode:
					if actual, ok := ev.(*eventXTermKbdMode); ok {
						if actual.Mode != expect.Mode {
							t.Errorf("Wrong mode: %v != %v", actual.Mode, expect.Mode)
						}
						return
					}
				case *eventKittyKbdMode:
					if actual, ok := ev.(*eventKittyKbdMode); ok {
						if actual.Mode != expect.Mode {
							t.Errorf("Wrong mode: %v != %v", actual.Mode, expect.Mode)
						}
						return
					}
				}
				t.Errorf("Wrong type, expected %T got %T", tt.result, ev)
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Timeout")
			}
		})
	}
}

// firstKey drains the event channel and returns the first EventKey received,
// or nil if none arrives within 100ms.
func firstKey(evch chan Event) *EventKey {
	var got *EventKey
	for {
		select {
		case ev := <-evch:
			if got == nil {
				if kev, ok := ev.(*EventKey); ok {
					got = kev
				}
			}
			continue
		case <-time.After(100 * time.Millisecond):
		}
		break
	}
	return got
}

func nextKey(evch chan Event) *EventKey {
	for {
		select {
		case ev := <-evch:
			if kev, ok := ev.(*EventKey); ok {
				return kev
			}
		case <-time.After(100 * time.Millisecond):
			return nil
		}
	}
}

func TestAdvancedControlKeys(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)
	ip.advanced = true

	ip.ScanUTF8([]byte{0x01, '\t'})

	got := nextKey(evch)
	if got == nil {
		t.Fatal("expected Ctrl-A event")
	}
	if got.Key() != KeyRune || got.Str() != "a" || got.Modifiers() != ModCtrl || got.Physical() != Key('a') {
		t.Fatalf("expected Ctrl-A as KeyRune a + ModCtrl, got key=%v str=%q mod=%v physical=%v", got.Key(), got.Str(), got.Modifiers(), got.Physical())
	}

	got = nextKey(evch)
	if got == nil {
		t.Fatal("expected Tab event")
	}
	if got.Key() != KeyTab || got.Modifiers() != ModNone {
		t.Fatalf("expected ambiguous Ctrl-I/Tab to remain Tab, got key=%v str=%q mod=%v", got.Key(), got.Str(), got.Modifiers())
	}
}

func TestAdvancedKittyKeyMetadata(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)
	ip.advanced = true

	ip.ScanUTF8([]byte("\x1b[65:97;2u\x1b[65:97;2:3u\x1b[57448;1:3u"))

	got := nextKey(evch)
	if got == nil {
		t.Fatal("expected shifted A press")
	}
	if got.Key() != KeyRune || got.Str() != "A" || got.Modifiers() != ModShift || got.Physical() != Key('a') || !got.Pressed() || got.Repeat() != 1 {
		t.Fatalf("unexpected shifted A press: key=%v str=%q mod=%v physical=%v pressed=%v repeat=%v", got.Key(), got.Str(), got.Modifiers(), got.Physical(), got.Pressed(), got.Repeat())
	}

	got = nextKey(evch)
	if got == nil {
		t.Fatal("expected shifted A release")
	}
	if got.Key() != KeyRune || got.Str() != "A" || got.Modifiers() != ModShift || got.Physical() != Key('a') || got.Pressed() {
		t.Fatalf("unexpected shifted A release: key=%v str=%q mod=%v physical=%v pressed=%v", got.Key(), got.Str(), got.Modifiers(), got.Physical(), got.Pressed())
	}

	got = nextKey(evch)
	if got == nil {
		t.Fatal("expected right Ctrl release")
	}
	if got.Key() != KeyCtrl || got.Modifiers()&ModRCtrl != ModRCtrl || got.Modifiers()&ModCtrl == 0 || got.Pressed() {
		t.Fatalf("unexpected right Ctrl release: key=%v mod=%v pressed=%v", got.Key(), got.Modifiers(), got.Pressed())
	}
}

func TestAdvancedKittyParameterParseErrorKeepsSubparametersAligned(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)
	ip.advanced = true

	ip.ScanUTF8([]byte("\x1b[65;999999999999999999999999999999999:3u"))

	got := nextKey(evch)
	if got == nil {
		t.Fatal("expected key event")
	}
	if got.Key() != KeyRune || got.Str() != "A" || got.Pressed() {
		t.Fatalf("expected shifted A release despite malformed modifier parameter, got key=%v str=%q pressed=%v", got.Key(), got.Str(), got.Pressed())
	}
}

func TestKeyConversionBounds(t *testing.T) {
	if key, ok := keyFromInt(32767); !ok || key != Key(32767) {
		t.Fatalf("keyFromInt max = %v, %v", key, ok)
	}
	if _, ok := keyFromInt(32768); ok {
		t.Fatal("keyFromInt accepted overflow")
	}
	if _, ok := keyFromInt(-1); ok {
		t.Fatal("keyFromInt accepted negative")
	}
	if key, ok := keyFromRune(32767); !ok || key != Key(32767) {
		t.Fatalf("keyFromRune max = %v, %v", key, ok)
	}
	if _, ok := keyFromRune(32768); ok {
		t.Fatal("keyFromRune accepted overflow")
	}
	if b, ok := asciiByteFromInt(0x7f); !ok || b != 0x7f {
		t.Fatalf("asciiByteFromInt max = %v, %v", b, ok)
	}
	if _, ok := asciiByteFromInt(0); ok {
		t.Fatal("asciiByteFromInt accepted zero")
	}
	if _, ok := asciiByteFromInt(0x80); ok {
		t.Fatal("asciiByteFromInt accepted non-ASCII")
	}
}

func TestAdvancedPhysicalKeyOverflowIsUnknown(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)
	ip.advanced = true

	ip.ScanUTF8([]byte("\x1b[65:40000;2u"))

	got := nextKey(evch)
	if got == nil {
		t.Fatal("expected key event")
	}
	if got.Key() != KeyRune || got.Str() != "A" || got.Physical() != 0 {
		t.Fatalf("expected overflow physical key to be unknown, got key=%v str=%q physical=%v", got.Key(), got.Str(), got.Physical())
	}
}

func TestLargeUnicodePhysicalKeyIsUnknown(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)
	ip.advanced = true

	ip.ScanUTF8([]byte("😀"))

	got := nextKey(evch)
	if got == nil {
		t.Fatal("expected key event")
	}
	if got.Key() != KeyRune || got.Str() != "😀" || got.Physical() != 0 {
		t.Fatalf("expected large Unicode physical key to be unknown, got key=%v str=%q physical=%v", got.Key(), got.Str(), got.Physical())
	}
}

func TestAdvancedWinKeyMetadata(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)
	ip.advanced = true

	ip.ScanUTF8([]byte("\x1b[65;0;65;1;16;3_\x1b[65;0;65;0;16;1_\x1b[17;0;0;1;8;1_\x1b[65;0;1;1;4;1_\x1b[160;0;0;1;20;1_"))

	got := nextKey(evch)
	if got == nil {
		t.Fatal("expected Win32 shifted A press")
	}
	if got.Key() != KeyRune || got.Str() != "A" || got.Modifiers() != ModShift || got.Physical() != Key('a') || !got.Pressed() || got.Repeat() != 3 {
		t.Fatalf("unexpected Win32 shifted A press: key=%v str=%q mod=%v physical=%v pressed=%v repeat=%v", got.Key(), got.Str(), got.Modifiers(), got.Physical(), got.Pressed(), got.Repeat())
	}

	got = nextKey(evch)
	if got == nil {
		t.Fatal("expected Win32 shifted A release")
	}
	if got.Key() != KeyRune || got.Str() != "A" || got.Modifiers() != ModShift || got.Physical() != Key('a') || got.Pressed() {
		t.Fatalf("unexpected Win32 shifted A release: key=%v str=%q mod=%v physical=%v pressed=%v", got.Key(), got.Str(), got.Modifiers(), got.Physical(), got.Pressed())
	}

	got = nextKey(evch)
	if got == nil {
		t.Fatal("expected Win32 Ctrl modifier press")
	}
	if got.Key() != KeyCtrl || got.Modifiers()&ModLCtrl != ModLCtrl || got.Modifiers()&ModCtrl == 0 || !got.Pressed() {
		t.Fatalf("unexpected Win32 Ctrl press: key=%v mod=%v pressed=%v", got.Key(), got.Modifiers(), got.Pressed())
	}

	got = nextKey(evch)
	if got == nil {
		t.Fatal("expected Win32 right-Ctrl A press")
	}
	if got.Key() != KeyRune || got.Str() != "a" || got.Modifiers()&ModRCtrl != ModRCtrl || got.Modifiers()&ModCtrl == 0 {
		t.Fatalf("unexpected Win32 right-Ctrl A press: key=%v str=%q mod=%v", got.Key(), got.Str(), got.Modifiers())
	}

	got = nextKey(evch)
	if got == nil {
		t.Fatal("expected Win32 left-Shift press with right-Ctrl held")
	}
	if got.Key() != KeyShift || got.Modifiers()&ModLShift != ModLShift ||
		got.Modifiers()&ModRCtrl != ModRCtrl || got.Modifiers()&ModCtrl == 0 {
		t.Fatalf("unexpected Win32 left-Shift press with right-Ctrl held: key=%v mod=%v", got.Key(), got.Modifiers())
	}
}

func TestWinKeyASCIIFallbackRangeCheck(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)

	ip.ScanUTF8([]byte("\x1b[0;0;65;1;0;1_"))

	got := nextKey(evch)
	if got == nil {
		t.Fatal("expected Win32 ASCII fallback key")
	}
	if got.Key() != KeyRune || got.Str() != "A" {
		t.Fatalf("unexpected Win32 ASCII fallback key: key=%v str=%q", got.Key(), got.Str())
	}

	ip.ScanUTF8([]byte("\x1b[0;0;0;1;0;1_"))

	if got = nextKey(evch); got != nil {
		t.Fatalf("unexpected Win32 NUL fallback key: key=%v str=%q", got.Key(), got.Str())
	}
}

func TestAdvancedModifierHelpers(t *testing.T) {
	winCases := []struct {
		vk      int
		wantKey Key
		wantMod ModMask
	}{
		{0x10, KeyShift, ModShift},
		{0xa0, KeyShift, ModLShift},
		{0xa1, KeyShift, ModRShift},
		{0x11, KeyCtrl, ModCtrl},
		{0xa2, KeyCtrl, ModLCtrl},
		{0xa3, KeyCtrl, ModRCtrl},
		{0x12, KeyAlt, ModAlt},
		{0xa4, KeyAlt, ModLAlt},
		{0xa5, KeyAlt, ModRAlt},
		{0x5b, KeyMeta, ModLMeta},
		{0x5c, KeyMeta, ModRMeta},
		{0x14, KeyCapsLock, ModNone},
	}
	for _, tt := range winCases {
		key, mod, ok := winModifierKey(tt.vk)
		if !ok || key != tt.wantKey || mod != tt.wantMod {
			t.Fatalf("winModifierKey(%#x) = %v, %v, %v; want %v, %v, true", tt.vk, key, mod, ok, tt.wantKey, tt.wantMod)
		}
	}
	if _, _, ok := winModifierKey(0xff); ok {
		t.Fatal("winModifierKey accepted unknown key")
	}

	kittyCases := map[int]ModMask{
		57441: ModLShift,
		57447: ModRShift,
		57442: ModLCtrl,
		57448: ModRCtrl,
		57443: ModLAlt,
		57449: ModRAlt,
		57444: ModLMeta,
		57450: ModRMeta,
	}
	for code, want := range kittyCases {
		if got := kittyModifierKey(code); got != want {
			t.Fatalf("kittyModifierKey(%d) = %v, want %v", code, got, want)
		}
	}
	if got := kittyModifierKey(12345); got != ModNone {
		t.Fatalf("kittyModifierKey unknown = %v, want ModNone", got)
	}
}

func TestWinModifierState(t *testing.T) {
	if got := calcWinModifier(0x1f, true); got != ModShift|ModLCtrl|ModRCtrl|ModLAlt|ModRAlt {
		t.Fatalf("advanced win modifier state = %v", got)
	}
	if got := calcWinModifier(0x1f, false); got != ModShift|ModCtrl|ModAlt {
		t.Fatalf("legacy win modifier state = %v", got)
	}
}

// TestEscDuringCsiResetsParser verifies that an ESC byte received while
// the parser is in CSI state correctly transitions to escape state per
// ECMA-48 §5.3.1, rather than being swallowed as a "bad parse".
func TestEscDuringCsiResetsParser(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)

	// Feed a partial CSI (ESC [) followed immediately by a new escape
	// sequence for Down arrow (ESC [ B). Without the fix, the second
	// ESC would be swallowed and 'B' would be emitted as a literal key.
	ip.ScanUTF8([]byte("\x1b[\x1b[B"))

	got := firstKey(evch)
	if got == nil {
		t.Fatal("expected a key event, got none")
	}
	if got.Key() != KeyDown {
		t.Errorf("expected KeyDown, got key=%v str=%q mod=%v", got.Key(), got.Str(), got.Modifiers())
	}
}

// TestEscDuringSs3ResetsParser verifies that an ESC byte received while
// the parser is accumulating SS3 parameters correctly transitions to
// escape state per ECMA-48 §5.3.1.
func TestEscDuringSs3ResetsParser(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputParser(evch)

	// Feed a partial SS3 with parameter (ESC O 1) followed immediately
	// by a new escape sequence for Down arrow (ESC [ B). Without the
	// fix, the ESC would be lost inside the SS3 parameter accumulation.
	ip.ScanUTF8([]byte("\x1bO1\x1b[B"))

	got := firstKey(evch)
	if got == nil {
		t.Fatal("expected a key event, got none")
	}
	if got.Key() != KeyDown {
		t.Errorf("expected KeyDown, got key=%v str=%q mod=%v", got.Key(), got.Str(), got.Modifiers())
	}
}
