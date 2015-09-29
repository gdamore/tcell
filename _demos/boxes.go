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
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

func makebox(s tcell.BufferedScreen) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	lx := rand.Int() % w
	ly := rand.Int() % h
	lw := rand.Int() % (w-lx)
	lh := rand.Int() % (h-ly)
	st := tcell.StyleDefault.Background(tcell.Color(rand.Int() % s.Colors()))

	for row := 0; row < lh; row++ {
		for col := 0; col < lw; col++ {
			s.SetCell(lx + col, ly + row, st, ' ')
		}
	}
	s.Show()
}

func main() {
	s, e := tcell.NewBufferedScreen()
	if e != nil {
		panic(e.Error())
	}
	if e = s.Init(); e != nil {
		panic(e.Error())
	}
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					s.Fini()
					os.Exit(0)
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	for {
		makebox(s)
		time.Sleep(time.Millisecond*10)
	}
}
