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

package tty

import (
	"reflect"
	"testing"
)

func TestDecodeUTF16Rune(t *testing.T) {
	tests := []struct {
		name string
		in   []rune
		want []rune
	}{
		{
			name: "bmp rune",
			in:   []rune{'A'},
			want: []rune{'A'},
		},
		{
			name: "valid surrogate pair",
			in:   []rune{0xD83D, 0xDE00},
			want: []rune{'😀'},
		},
		{
			name: "orphaned high before bmp",
			in:   []rune{0xD83D, ' '},
			want: []rune{'�', ' '},
		},
		{
			name: "orphaned high before high",
			in:   []rune{0xD83D, 0xD83D, 0xDE00},
			want: []rune{'�', '😀'},
		},
		{
			name: "orphaned low",
			in:   []rune{0xDE00},
			want: []rune{'�'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var surrogate rune
			var got []rune
			for _, wc := range tt.in {
				got = append(got, decodeUTF16Rune(&surrogate, wc)...)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("decodeUTF16Rune(%U) = %U, want %U", tt.in, got, tt.want)
			}
		})
	}
}
