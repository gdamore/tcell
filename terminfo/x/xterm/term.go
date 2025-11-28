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

package xterm

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// xterm terminal emulator (X Window System)
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:        "xterm",
		Aliases:     []string{"xterm-debian"},
		Columns:     80,
		Lines:       24,
		EnterKeypad: "\x1b[?1h\x1b=",
		ExitKeypad:  "\x1b[?1l\x1b>",
		Mouse:       "\x1b[<",
		XTermLike:   true,
	})

	// xterm with 256 colors
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:        "xterm-256color",
		Aliases:     []string{"alacritty", "ghostty", "rio", "st", "xterm-ghostty", "xterm-kitty"},
		Columns:     80,
		Lines:       24,
		EnterKeypad: "\x1b[?1h\x1b=",
		ExitKeypad:  "\x1b[?1l\x1b>",
		Mouse:       "\x1b[<",
		XTermLike:   true,
	})
}
