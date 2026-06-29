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

//go:build windows
// +build windows

package tty

import (
	"encoding/binary"
	"testing"
)

func makeWinKeyRecord(keyDown bool, repeat, virtualKey, scanCode, unicodeChar uint16, controlState uint32) [16]byte {
	var data [16]byte
	if keyDown {
		binary.LittleEndian.PutUint32(data[0:], 1)
	}
	binary.LittleEndian.PutUint16(data[4:], repeat)
	binary.LittleEndian.PutUint16(data[6:], virtualKey)
	binary.LittleEndian.PutUint16(data[8:], scanCode)
	binary.LittleEndian.PutUint16(data[10:], unicodeChar)
	binary.LittleEndian.PutUint32(data[12:], controlState)
	return data
}

func TestEncodeWinKeyRecordPreservesKeyState(t *testing.T) {
	tests := []struct {
		name string
		data [16]byte
		want string
	}{
		{
			name: "press",
			data: makeWinKeyRecord(true, 3, 65, 30, 'A', 0x10),
			want: "\x1b[65;30;65;1;16;3_",
		},
		{
			name: "release",
			data: makeWinKeyRecord(false, 1, 65, 30, 'A', 0x10),
			want: "\x1b[65;30;65;0;16;1_",
		},
		{
			name: "zero repeat becomes one",
			data: makeWinKeyRecord(true, 0, 65, 30, 'A', 0),
			want: "\x1b[65;30;65;1;0;1_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var surrogate rune
			if got := string(encodeWinKeyRecord(tt.data, &surrogate)); got != tt.want {
				t.Fatalf("encodeWinKeyRecord() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEncodeWinKeyRecordFallbackIgnoresRelease(t *testing.T) {
	var surrogate rune
	data := makeWinKeyRecord(false, 1, 0, 0, 'A', 0)
	if got := encodeWinKeyRecord(data, &surrogate); len(got) != 0 {
		t.Fatalf("encodeWinKeyRecord() = %q, want no bytes", string(got))
	}

	data = makeWinKeyRecord(true, 1, 0, 0, 'A', 0)
	if got := string(encodeWinKeyRecord(data, &surrogate)); got != "A" {
		t.Fatalf("encodeWinKeyRecord() = %q, want %q", got, "A")
	}

	data = makeWinKeyRecord(true, 3, 0, 0, 'A', 0)
	if got := string(encodeWinKeyRecord(data, &surrogate)); got != "AAA" {
		t.Fatalf("encodeWinKeyRecord() = %q, want %q", got, "AAA")
	}
}
