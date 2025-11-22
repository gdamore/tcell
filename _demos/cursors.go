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

	text := "This demonstrates cursor styles.  Press 0 through 6 to change the style."
	s.PutStr(1, 1, text)

	s.Put(2, 2, "0", tcell.StyleDefault)
	s.SetCursorStyle(tcell.CursorStyleDefault)
	s.ShowCursor(3, 2)
	quit := make(chan struct{})
	go func() {
		for {
			s.Show()
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyRune:
					switch ev.Rune() {
					case '0':
						s.Put(2, 2, "0", tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleDefault, tcell.ColorReset)
					case '1':
						s.Put(2, 2, "1", tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleBlinkingBlock, tcell.ColorGreen)
					case '2':
						s.Put(2, 2, "2", tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleSteadyBlock, tcell.ColorBlue)
					case '3':
						s.Put(2, 2, "3", tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleBlinkingUnderline, tcell.ColorRed)
					case '4':
						s.Put(2, 2, "4", tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleSteadyUnderline, tcell.ColorOrange)
					case '5':
						s.Put(2, 2, "5", tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleBlinkingBar, tcell.ColorYellow)
					case '6':
						s.Put(2, 2, "6", tcell.StyleDefault)
						s.SetCursorStyle(tcell.CursorStyleSteadyBar, tcell.ColorPink)
					}

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
