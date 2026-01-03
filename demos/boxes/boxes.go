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

// boxes just displays random colored boxes on your terminal screen.
// Press escape or control-q to exit the program.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func makeBox(s tcell.Screen) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	glyphs := []string{"@", "#", "&", "*", "=", "%", "Z", "A"}

	lx := rand.Int() % w
	ly := rand.Int() % h
	lw := rand.Int() % (w - lx)
	lh := rand.Int() % (h - ly)
	st := tcell.StyleDefault
	gl := " "
	if s.Colors() > 256 {
		rgb := tcell.NewHexColor(int32(rand.Int() & 0xffffff))
		st = st.Background(rgb)
	} else if s.Colors() > 1 {
		st = st.Background(color.Color(rand.Int()%s.Colors()) | color.IsValid)
	} else {
		st = st.Reverse(rand.Int()%2 == 0)
		gl = glyphs[rand.Int()%len(glyphs)]
	}

	for row := range lh {
		for col := range lw {
			s.PutStrStyled(lx+col, ly+row, gl, st)
		}
	}
	s.Show()
}

var (
	count    = 0
	interval = time.Millisecond * 5
	drawTime time.Duration
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

	s.SetStyle(tcell.StyleDefault.
		Foreground(color.Black).
		Background(color.White))
	s.Clear()

loop:
	for {
		select {
		case ev := <-s.EventQ():
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyCtrlQ:
					break loop
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		case <-time.After(interval):
		}
		start := time.Now()
		makeBox(s)
		count++
		drawTime += time.Since(start)
	}

	s.Fini()
	fmt.Printf("Finished %d boxes in %s (drawing time only)\n", count, drawTime)
	if count > 0 {
		fmt.Printf("Average is %0.3f ms / box\n", (float64(drawTime)/float64(count))/1000000.0)
	}
}
