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

package tcell

import (
	"strings"
	"testing"

	"github.com/clipperhouse/displaywidth"
)

var (
	mixedWidthText = strings.Repeat("Hello, 世界 👩‍🚀 e\u0301 — π «ambiguous» ", 8)
	ansiWidthText  = strings.Repeat("\x1b[31mHello\x1b[0m, 世界 👩‍🚀 e\u0301 \x1b]0;title\x07 ", 8)
)

func BenchmarkStringWidthMixedCurrent(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = currentMixedWidth(mixedWidthText)
	}
}

func BenchmarkStringWidthANSIControl(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = currentANSIWidth(ansiWidthText)
	}
}

func currentMixedWidth(s string) int {
	return displaywidth.Options{EastAsianWidth: true}.String(s)
}

func currentANSIWidth(s string) int {
	return displaywidth.Options{
		EastAsianWidth:       true,
		ControlSequences:     true,
		ControlSequences8Bit: true,
	}.String(s)
}
