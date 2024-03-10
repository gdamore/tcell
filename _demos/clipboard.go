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

package main

import (
	"fmt"
	"os"
	"unicode/utf8"

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

var clipboard []byte

func displayHelloWorld(s tcell.Screen) {
	w, h := s.Size()
	s.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorCadetBlue.TrueColor()).Background(tcell.ColorWhite)
	emitStr(s, w/2-14, h/2, style, "Press 1 to set clipboard")
	emitStr(s, w/2-14, h/2+1, style, "Press 2 to get clipboard")

	msg := ""
	if utf8.Valid(clipboard) {
		cp := string(clipboard)
		if len(cp) >= w-25 {
			cp = cp[:21] + " ..."
		}
		msg = fmt.Sprintf("Clipboard (%d bytes): %s", len(clipboard), cp)
	} else if clipboard != nil {
		msg = fmt.Sprintf("Clipboard (%d bytes) Not Valid UTF-8", len(clipboard))
	} else {
		msg = "No clipboard data"
	}
	emitStr(s, (w-len(msg))/2, h/2+3, tcell.StyleDefault, msg)
	emitStr(s, w/2-9, h/2+5, tcell.StyleDefault, "Press ESC to exit.")
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

	displayHelloWorld(s)

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			displayHelloWorld(s)
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch ev.Rune() {
				case '1':
					s.SetClipboard([]byte("Enjoy your new clipboard content!"))
				case '2':
					s.GetClipboard()
				}
			case tcell.KeyEscape:
				s.Fini()
				os.Exit(0)
			}
		case *tcell.EventClipboard:
			clipboard = ev.Data()
			displayHelloWorld(s)
		}
	}
}
