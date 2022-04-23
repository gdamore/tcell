//go:build ignore
// +build ignore

// Copyright 2022 The TCell Authors
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

package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"

	"github.com/mattn/go-runewidth"
)

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func displayDemo(s tcell.Screen) {
	w, h := s.Size()
	s.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorCadetBlue.TrueColor()).Background(tcell.ColorWhite)
	sizeStr := fmt.Sprintf("%d x %d", w, h)
	helpStr := "Use cursors to resize, ESC to exit."
	emitStr(s, (w-len(sizeStr))/2, h/2, style, sizeStr)
	emitStr(s, (w-len(helpStr))/2, h/2+1, tcell.StyleDefault, helpStr)
	s.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	displayDemo(s)

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			displayDemo(s)
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				s.Fini()
				os.Exit(0)
			case tcell.KeyRight:
				w, h := s.Size()
				s.SetSize(w+1, h)
			case tcell.KeyLeft:
				w, h := s.Size()
				s.SetSize(w-1, h)
			case tcell.KeyUp:
				w, h := s.Size()
				s.SetSize(w, h-1)
			case tcell.KeyDown:
				w, h := s.Size()
				s.SetSize(w, h+1)
			}
		}
	}
}
