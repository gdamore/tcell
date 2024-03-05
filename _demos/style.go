//go:build ignore
// +build ignore

// Copyright 2024 The TCell Authors
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

// unicode just displays a Unicode test on your screen.
// Press ESC to exit the program.
package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	runewidth "github.com/mattn/go-runewidth"
)

var row = 0
var style = tcell.StyleDefault

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}

func main() {

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	plain := tcell.StyleDefault
	bold := style.Bold(true)

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()

	quit := make(chan struct{})

	style = bold.Foreground(tcell.ColorBlue).Background(tcell.ColorSilver)

	row = 2
	puts(s, style, 2, row, "Press ESC to Exit")
	row = 4
	puts(s, plain, 2, row, "Note: Style support is dependent on your terminal.")
	row = 6

	plain = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	style = plain
	puts(s, style, 2, row, "Plain")
	row++

	style = plain.Blink(true)
	puts(s, style, 2, row, "Blink")
	row++

	style = plain.Reverse(true)
	puts(s, style, 2, row, "Reverse")
	row++

	style = plain.Dim(true)
	puts(s, style, 2, row, "Dim")
	row++

	style = plain.Underline(true)
	puts(s, style, 2, row, "Underline")
	row++

	style = plain.Italic(true)
	puts(s, style, 2, row, "Italic")
	row++

	style = plain.Bold(true)
	puts(s, style, 2, row, "Bold")
	row++

	style = plain.Bold(true).Italic(true)
	puts(s, style, 2, row, "Bold Italic")
	row++

	style = plain.Bold(true).Italic(true).Underline(true)
	puts(s, style, 2, row, "Bold Italic Underline")
	row++

	style = plain.StrikeThrough(true)
	puts(s, style, 2, row, "Strikethrough")
	row++

	style = plain.Underline(tcell.UnderlineStyleDouble)
	puts(s, style, 2, row, "Double Underline")
	row++

	style = plain.Underline(tcell.UnderlineStyleCurly)
	puts(s, style, 2, row, "Curly Underline")
	row++

	style = plain.Underline(tcell.UnderlineStyleDotted)
	puts(s, style, 2, row, "Dotted Underline")
	row++

	style = plain.Underline(tcell.UnderlineStyleDashed)
	puts(s, style, 2, row, "Dashed Underline")
	row++

	style = plain.Underline(true, tcell.ColorBlue)
	puts(s, style, 2, row, "Blue Underline")
	row++

	style = plain.Underline(tcell.UnderlineStyleSolid, tcell.ColorFireBrick)
	puts(s, style, 2, row, "Firebrick Underline")
	row++

	style = plain.Underline(tcell.UnderlineStyleCurly, tcell.NewRGBColor(0xc5, 0x8a, 0xf9))
	puts(s, style, 2, row, "Pink Curly Underline")
	row++

	style = plain.Url("http://github.com/gdamore/tcell")
	puts(s, style, 2, row, "HyperLink")
	row++

	style = plain.Foreground(tcell.ColorRed)
	puts(s, style, 2, row, "Red Foreground")
	row++

	style = plain.Background(tcell.ColorRed)
	puts(s, style, 2, row, "Red Background")
	row++

	s.Show()
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	<-quit

	s.Fini()
}
