//go:build ignore
// +build ignore

// Copyright 2022 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// stress will fill the whole screen with random characters, colors and
// formatting. The frames are pre-generated to draw as fast as possible.
// ESC and Ctrl-C will end the program. Note that resizing isn't supported.
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}

	var frames int

	type cell struct {
		c     rune
		style tcell.Style
	}

	width, height := screen.Size()
	glyphs := []rune{'@', '#', '&', '*', '=', '%', 'Z', 'A'}
	attrs := []tcell.AttrMask{tcell.AttrBold, tcell.AttrReverse, tcell.AttrItalic, tcell.AttrNone}

	// Pre-Generate 100 different frame patterns, so we stress the terminal as
	// much as possible :D
	var patterns [][][]*cell
	for i := 0; i < 100; i++ {
		pattern := make([][]*cell, height)
		for h := 0; h < height; h++ {
			row := make([]*cell, width)
			for w := 0; w < width; w++ {
				rF := int32(rand.Int() % 256)
				gF := int32(rand.Int() % 256)
				bF := int32(rand.Int() % 256)
				rB := int32(rand.Int() % 256)
				gB := int32(rand.Int() % 256)
				bB := int32(rand.Int() % 256)

				row[w] = &cell{
					c: glyphs[rand.Int()%len(glyphs)],
					style: tcell.StyleDefault.
						Attributes(attrs[rand.Int()%len(attrs)]).
						Background(tcell.NewRGBColor(rB, gB, bB)).
						Foreground(tcell.NewRGBColor(rF, gF, bF)),
				}
			}
			pattern[h] = row
		}
		patterns = append(patterns, pattern)
	}

	evCh := make(chan tcell.Event)
	quitCh := make(chan struct{})

	go screen.ChannelEvents(evCh, quitCh)
	startTime := time.Now()
loop:
	for {
		select {
		case event := <-evCh:
			if event, ok := event.(*tcell.EventKey); ok {
				if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyESC {
					close(quitCh)
					break loop
				}
			}
			break
		default:
		}
		pattern := patterns[frames%len(patterns)]
		for h := 0; h < height; h++ {
			for w := 0; w < width; w++ {
				c := pattern[h][w]
				screen.SetContent(w, h, c.c, nil, c.style)
			}
		}
		screen.Show()
		frames++
	}
	duration := time.Since(startTime)
	screen.Fini()
	fps := int(float64(frames) / duration.Seconds())
	fmt.Println("FPS:", fps)
}
