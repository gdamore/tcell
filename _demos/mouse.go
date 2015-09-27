//+build ignore

// Copyright 2015 The TCell Authors
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

// boxes just displays random colored boxes on your terminal screen.
// Press ESC to exit the program.
package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

// This program just shows simple mouse and keyboard events.  Press ESC to
// exit.
func main() {
	s, e := tcell.NewBufferedScreen()
	if e != nil {
		fmt.Printf("oops: %v", e)
	}
	s.Init()
	s.EnableMouse()
	s.Clear()

	i := 1
	for _, c := range "Press ESC to exit." {
		s.SetCell(i, 1, tcell.StyleDefault, c)
		i++
	}

	for {
		s.Show()
		ev := s.PollEvent()
		st := tcell.StyleDefault.Background(tcell.ColorBrightRed)
		up := tcell.StyleDefault.Background(tcell.ColorBlue)
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			x, y := ev.Size()
			s.SetCell(x-1, y-1, st, 'R')
		case *tcell.EventKey:
			x, y := s.Size()
			s.SetCell(x-2, y-2, st, ev.Rune())
			s.SetCell(x-1, y-1, st, 'K')
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			switch ev.Buttons() {
			case tcell.ButtonNone:
				s.SetCell(x, y, up, '-')
			case tcell.Button1:
				s.SetCell(x, y, st, '1')
			case tcell.Button2:
				s.SetCell(x, y, st, '2')
			case tcell.Button3:
				s.SetCell(x, y, st, '3')
			default:
				s.SetCell(x, y, st, '*')
			}
			x, y = s.Size()
			s.SetCell(x-1, y-1, st, 'M')
		default:
			x, y := s.Size()
			s.SetCell(x-1, y-1, st, 'X')
		}
	}
}
