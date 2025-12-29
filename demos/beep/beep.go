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

// beep makes a beep for every 3 seconds, or when B is pressed, until you press ESC or CTRL-Q
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v3"
)

func draw(s tcell.Screen, remain int) {
	style := tcell.StyleDefault
	s.Clear()
	s.PutStrStyled(1, 1, fmt.Sprintf("Beep will occur in %d seconds...", remain), style)
	s.PutStrStyled(1, 3, "Press ESC or CTRL-Q to quit, B to beep now.", style.Italic(true))

	s.Show()
}

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
	defer s.Fini()

	s.SetStyle(tcell.StyleDefault)

	remain := 3
	draw(s, remain)
	for {
		select {
		case ev := <-s.EventQ():
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyCtrlQ:
					return
				case tcell.KeyCtrlL:
					draw(s, remain)
					s.Sync()
				case tcell.KeyRune:
					if ev.Str() == "b" || ev.Str() == "B" {
						remain = 3
						s.Beep()
						draw(s, remain)
					}
				}
			case *tcell.EventResize:
				s.Sync()
			}
		case <-time.After(time.Second):
			// imprecise, but good enough for demo
			remain--
			if remain == 0 {
				remain = 3
				s.Beep()
			}
			draw(s, remain)
		}
	}
}
