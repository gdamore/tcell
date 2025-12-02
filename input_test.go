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

package tcell

import (
	"testing"
	"time"
)

// TestInputProcessorNullByte tests that null byte (0x00) is correctly handled
// as Ctrl+Space per the fix in the scan() function.
func TestInputProcessorNullByte(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputProcessor(evch)

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

// TestInputProcessorControlKeys tests control key handling for bytes 1-31.
// Note: NewEventKey converts control characters to KeyCtrlA-Z for 0x01-0x1A,
// and KeyRune for 0x1C-0x1F with the character and ModCtrl.
func TestInputProcessorControlKeys(t *testing.T) {
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
		{"Ctrl+J", 0x0A, KeyEnter, "", ModNone},     // LF is Enter
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
		{"Ctrl+[", 0x1B, KeyEscape, "", ModNone},    // ESC is special
		{"Ctrl+\\", 0x1C, KeyRune, "\\", ModCtrl},   // becomes KeyRune with string
		{"Ctrl+]", 0x1D, KeyRune, "]", ModCtrl},     // becomes KeyRune with string
		{"Ctrl+^", 0x1E, KeyRune, "^", ModCtrl},     // becomes KeyRune with string
		{"Ctrl+_", 0x1F, KeyRune, "_", ModCtrl},     // becomes KeyRune with string
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputProcessor(evch)

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

// TestInputProcessorNullVsOtherControlChars tests the boundary between
// null byte and other control characters.
func TestInputProcessorNullVsOtherControlChars(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputProcessor(evch)

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

// TestInputProcessorPrintableCharacters tests that printable characters
// are handled correctly without control modifiers.
func TestInputProcessorPrintableCharacters(t *testing.T) {
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
			ip := newInputProcessor(evch)

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

// TestInputProcessorSpecialKeys tests special key handling (tab, backspace, enter).
func TestInputProcessorSpecialKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected Key
		mod      ModMask
	}{
		{"Tab", '\t', KeyTab, ModNone},
		{"Backspace_BS", '\b', KeyBackspace, ModNone},
		{"Backspace_DEL", 0x7F, KeyBackspace, ModNone},
		{"Enter_LF", '\n', KeyEnter, ModNone},
		{"Enter_CR", '\r', KeyEnter, ModNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputProcessor(evch)

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

// TestInputProcessorSequentialInput tests handling multiple inputs in sequence.
func TestInputProcessorSequentialInput(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputProcessor(evch)

	// Send: null, Ctrl+A, 'B', space
	inputs := []byte{0x00, 0x01, 'B', ' '}
	expected := []struct {
		key Key
		str string
		mod ModMask
	}{
		{KeyRune, " ", ModCtrl},  // null -> Ctrl+Space
		{KeyCtrlA, "", ModCtrl},  // 0x01 -> KeyCtrlA with ModCtrl
		{KeyRune, "B", ModNone},  // 'B'
		{KeyRune, " ", ModNone},  // space
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

// TestInputProcessorUTF8Characters tests UTF-8 multibyte character handling.
func TestInputProcessorUTF8Characters(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evch := make(chan Event, 10)
			ip := newInputProcessor(evch)

			ip.ScanUTF8(tt.input)

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

// TestInputProcessorEdgeCases tests edge cases and boundary conditions.
func TestInputProcessorEdgeCases(t *testing.T) {
	t.Run("EmptyInput", func(t *testing.T) {
		evch := make(chan Event, 10)
		ip := newInputProcessor(evch)

		ip.ScanUTF8([]byte{})

		select {
		case ev := <-evch:
			t.Errorf("Expected no event for empty input, got %T", ev)
		case <-time.After(50 * time.Millisecond):
			// Success - no event generated
		}
	})

	t.Run("MultipleNullBytes", func(t *testing.T) {
		evch := make(chan Event, 10)
		ip := newInputProcessor(evch)

		ip.ScanUTF8([]byte{0x00, 0x00, 0x00})

		for i := 0; i < 3; i++ {
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
		ip := newInputProcessor(evch)

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
		ip := newInputProcessor(evch)

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

// TestInputProcessorConcurrentAccess tests that the inputProcessor is safe
// for concurrent access (it uses a mutex internally).
func TestInputProcessorConcurrentAccess(t *testing.T) {
	evch := make(chan Event, 100)
	ip := newInputProcessor(evch)

	done := make(chan bool)

	// Goroutine 1: Send null bytes
	go func() {
		for i := 0; i < 10; i++ {
			ip.ScanUTF8([]byte{0x00})
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Goroutine 2: Send regular characters
	go func() {
		for i := 0; i < 10; i++ {
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

// TestInputProcessorStateTransitions tests that the null byte fix doesn't
// interfere with state machine transitions.
func TestInputProcessorStateTransitions(t *testing.T) {
	evch := make(chan Event, 10)
	ip := newInputProcessor(evch)

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