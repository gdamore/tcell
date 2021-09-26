// +build ignore

// Copyright 2021 The TCell Authors
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

// beep makes a beep every second until you press ESC
package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
)

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault)
	s.Clear()

	s.SetCell(2, 2, tcell.StyleDefault, '0')
	s.SetCursorStyle(tcell.CursorStyleDefault)
	s.ShowCursor(3, 2)
	quit := make(chan struct{})
	style := tcell.StyleDefault
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyRune:
					switch ev.Rune() {
					case '0':
						s.SetContent(2, 2, '0', nil, style)
						s.SetCursorStyle(tcell.CursorStyleDefault)
					case '1':
						s.SetContent(2, 2, '1', nil, style)
						s.SetCursorStyle(tcell.CursorStyleBlinkingBlock)
					case '2':
						s.SetCell(2, 2, tcell.StyleDefault, '2')
						s.SetCursorStyle(tcell.CursorStyleSteadyBlock)
					case '3':
						s.SetCell(2, 2, tcell.StyleDefault, '3')
						s.SetCursorStyle(tcell.CursorStyleBlinkingUnderline)
					case '4':
						s.SetCell(2, 2, tcell.StyleDefault, '4')
						s.SetCursorStyle(tcell.CursorStyleSteadyUnderline)
					case '5':
						s.SetCell(2, 2, tcell.StyleDefault, '5')
						s.SetCursorStyle(tcell.CursorStyleBlinkingBar)
					case '6':
						s.SetCell(2, 2, tcell.StyleDefault, '6')
						s.SetCursorStyle(tcell.CursorStyleSteadyBar)
					}
					s.Show()

				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlC:
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
