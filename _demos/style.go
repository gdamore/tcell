//go:build ignore
// +build ignore

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

// unicode just displays a Unicode test on your screen.
// Press ESC to exit the program.
package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

var row = 0
var style = tcell.StyleDefault

func main() {

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

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
	s.PutStrStyled(2, row, "Press ESC to Exit", style)
	row = 4
	s.PutStrStyled(2, row, "Note: Style support is dependent on your terminal.", plain)
	row = 6

	plain = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)

	style = plain
	s.PutStrStyled(2, row, "Plain", style)
	row++

	style = plain.Blink(true)
	s.PutStrStyled(2, row, "Blink", style)
	row++

	style = plain.Reverse(true)
	s.PutStrStyled(2, row, "Reverse", style)
	row++

	style = plain.Dim(true)
	s.PutStrStyled(2, row, "Dim", style)
	row++

	style = plain.Underline(true)
	s.PutStrStyled(2, row, "Underline", style)
	row++

	style = plain.Italic(true)
	s.PutStrStyled(2, row, "Italic", style)
	row++

	style = plain.Bold(true)
	s.PutStrStyled(2, row, "Bold", style)
	row++

	style = plain.Bold(true).Italic(true)
	s.PutStrStyled(2, row, "Bold Italic", style)
	row++

	style = plain.Bold(true).Italic(true).Underline(true)
	s.PutStrStyled(2, row, "Bold Italic Underline", style)
	row++

	style = plain.StrikeThrough(true)
	s.PutStrStyled(2, row, "Strikethrough", style)
	row++

	style = plain.Underline(tcell.UnderlineStyleDouble)
	s.PutStrStyled(2, row, "Double Underline", style)
	row++

	style = plain.Underline(tcell.UnderlineStyleCurly)
	s.PutStrStyled(2, row, "Curly Underline", style)
	row++

	style = plain.Underline(tcell.UnderlineStyleDotted)
	s.PutStrStyled(2, row, "Dotted Underline", style)
	row++

	style = plain.Underline(tcell.UnderlineStyleDashed)
	s.PutStrStyled(2, row, "Dashed Underline", style)
	row++

	style = plain.Underline(true, tcell.ColorBlue)
	s.PutStrStyled(2, row, "Blue Underline", style)
	row++

	style = plain.Underline(tcell.UnderlineStyleSolid, tcell.ColorFireBrick)
	s.PutStrStyled(2, row, "Firebrick Underline", style)
	row++

	style = plain.Underline(tcell.UnderlineStyleCurly, tcell.NewRGBColor(0xc5, 0x8a, 0xf9))
	s.PutStrStyled(2, row, "Pink Curly Underline", style)
	row++

	style = plain.Url("http://github.com/gdamore/tcell")
	s.PutStrStyled(2, row, "HyperLink", style)
	row++

	style = plain.Foreground(tcell.ColorRed)
	s.PutStrStyled(2, row, "Red Foreground", style)
	row++

	style = plain.Background(tcell.ColorRed)
	s.PutStrStyled(2, row, "Red Background", style)
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
