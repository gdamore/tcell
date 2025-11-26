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
		Name:              "xterm",
		Aliases:           []string{"xterm-debian"},
		Columns:           80,
		Lines:             24,
		Colors:            8,
		Clear:             "\x1b[H\x1b[2J",
		EnterCA:           "\x1b[?1049h\x1b[22;0;0t",
		ExitCA:            "\x1b[?1049l\x1b[23;0;0t",
		ShowCursor:        "\x1b[?12l\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b(B\x1b[m",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m",
		Dim:               "\x1b[2m",
		Italic:            "\x1b[3m",
		Blink:             "\x1b[5m",
		Reverse:           "\x1b[7m",
		EnterKeypad:       "\x1b[?1h\x1b=",
		ExitKeypad:        "\x1b[?1l\x1b>",
		SetFg:             "\x1b[3%p1%dm",
		SetBg:             "\x1b[4%p1%dm",
		SetFgBg:           "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg:         "\x1b[39;49m",
		AltChars:          "``aaffggiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x1b(0",
		ExitAcs:           "\x1b(B",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		StrikeThrough:     "\x1b[9m",
		Mouse:             "\x1b[<",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH",
		XTermLike:         true,
	})

	// xterm with 88 colors
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "xterm-88color",
		Columns:           80,
		Lines:             24,
		Colors:            88,
		Clear:             "\x1b[H\x1b[2J",
		EnterCA:           "\x1b[?1049h\x1b[22;0;0t",
		ExitCA:            "\x1b[?1049l\x1b[23;0;0t",
		ShowCursor:        "\x1b[?12l\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b(B\x1b[m",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m",
		Dim:               "\x1b[2m",
		Italic:            "\x1b[3m",
		Blink:             "\x1b[5m",
		Reverse:           "\x1b[7m",
		EnterKeypad:       "\x1b[?1h\x1b=",
		ExitKeypad:        "\x1b[?1l\x1b>",
		SetFg:             "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m",
		SetBg:             "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m",
		SetFgBg:           "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;;%?%p2%{8}%<%t4%p2%d%e%p2%{16}%<%t10%p2%{8}%-%d%e48;5;%p2%d%;m",
		ResetFgBg:         "\x1b[39;49m",
		AltChars:          "``aaffggiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x1b(0",
		ExitAcs:           "\x1b(B",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		StrikeThrough:     "\x1b[9m",
		Mouse:             "\x1b[<",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH",
		XTermLike:         true,
	})

	// xterm with 256 colors
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "xterm-256color",
		Aliases:           []string{"alacritty", "ghostty", "rio", "st", "xterm-ghostty", "xterm-kitty"},
		Columns:           80,
		Lines:             24,
		Colors:            256,
		Clear:             "\x1b[H\x1b[2J",
		EnterCA:           "\x1b[?1049h\x1b[22;0;0t",
		ExitCA:            "\x1b[?1049l\x1b[23;0;0t",
		ShowCursor:        "\x1b[?12l\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b(B\x1b[m",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m",
		Dim:               "\x1b[2m",
		Italic:            "\x1b[3m",
		Blink:             "\x1b[5m",
		Reverse:           "\x1b[7m",
		EnterKeypad:       "\x1b[?1h\x1b=",
		ExitKeypad:        "\x1b[?1l\x1b>",
		SetFg:             "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m",
		SetBg:             "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m",
		SetFgBg:           "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;;%?%p2%{8}%<%t4%p2%d%e%p2%{16}%<%t10%p2%{8}%-%d%e48;5;%p2%d%;m",
		ResetFgBg:         "\x1b[39;49m",
		AltChars:          "``aaffggiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x1b(0",
		ExitAcs:           "\x1b(B",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		StrikeThrough:     "\x1b[9m",
		Mouse:             "\x1b[<",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH",
		XTermLike:         true,
		TrueColor:         true,
	})
}
