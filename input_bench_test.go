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
)

func BenchmarkInputParser(b *testing.B) {
	cases := []struct {
		name string
		data string
	}{
		{"ascii", strings.Repeat("hello world, typing some text. ", 32)},
		{"paste", "\x1b[200~" + strings.Repeat("func main() { fmt.Println(\"x\") }\n", 2000) + "\x1b[201~"},
		{"arrows", strings.Repeat("\x1b[A\x1b[B\x1b[C\x1b[D\x1b[1;5C\x1bOP", 100)},
		{"key", "a"},
	}
	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			evch := make(chan Event, 1<<17)
			ip := newInputParser(evch)
			data := []byte(tc.data)
			b.SetBytes(int64(len(data)))
			b.ReportAllocs()
			for b.Loop() {
				ip.ScanUTF8(data)
				for len(evch) > 0 {
					<-evch
				}
			}
		})
	}
}
