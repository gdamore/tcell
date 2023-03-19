//go:build ignore
// +build ignore

// Copyright 2023 The TCell Authors
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
//
// sixel displays a sixel and demonstrates how to use direct drawing
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"

	"github.com/mattn/go-runewidth"
	"github.com/mattn/go-sixel"
)

type imageData struct {
	width  int // width in pixels
	height int // height in pixels
	data   *bytes.Buffer
}

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

func displayHelloWorld(s tcell.Screen) {
	w, h := s.Size()
	s.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorCadetBlue.TrueColor()).Background(tcell.ColorWhite)
	emitStr(s, w/2-7, h/2, style, "Hello, World!")
	emitStr(s, w/2-9, h/2+1, tcell.StyleDefault, "Press ESC to exit.")
	emitStr(s, w/2-18, h/2+2, tcell.StyleDefault, "Press Enter to toggle sixel lock.")
	s.Show()
}

func displaySixel(s tcell.Screen, img *imageData, lock bool) {
	tty, ok := s.Tty()
	if !ok {
		s.Fini()
		log.Fatal("not a terminal")
	}
	ws, err := tty.WindowSize()
	if err != nil {
		s.Fini()
		log.Fatal(err)
	}
	// Get the dimensions of a single cell
	cw, ch := ws.CellDimensions()
	if cw == 0 || ch == 0 {
		s.Fini()
		log.Fatal("terminal does not support sixel graphics")
		return
	}

	// Calculate the image dimensions in cells. We round up to prevent
	// drawing on a partially filled cell
	sixelWidth := int(math.Ceil(float64(img.width) / float64(cw)))
	sixelHeight := int(math.Ceil(float64(img.height) / float64(ch)))

	sixelX := ws.Width/2 - (sixelWidth / 2) // Center the image horizontally
	sixelY := ws.Height/2 - sixelHeight - 2
	if sixelY < 0 {
		sixelY = 0
	}
	// Lock the region where we will draw the sixel, this prevents tcell
	// from drawing over this area
	s.LockRegion(sixelX, sixelY, sixelWidth, sixelHeight, lock)

	// Get the terminfo for our current terminal
	ti, err := tcell.LookupTerminfo(os.Getenv("TERM"))
	if err != nil {
		s.Fini()
		log.Fatal(err)
	}

	emitStr(s, sixelX, sixelY, tcell.StyleDefault, "This text is behind")
	emitStr(s, sixelX, sixelY+1, tcell.StyleDefault, "     the sixel")

	// Move the cursor to our draw position
	ti.TPuts(tty, ti.TGoto(sixelX, sixelY))
	// Draw the sixel data
	ti.TPuts(tty, img.data.String())

	s.Show()
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return png.Decode(f)
}

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

	raw, err := loadImage("./logos/tcell.png")
	if err != nil {
		s.Fini()
		log.Println("couldn't load image. try running from the root directory")
		log.Fatalf("        go run ./_demos/sixel.go")
	}

	img := &imageData{
		width:  raw.Bounds().Dx(),
		height: raw.Bounds().Dy(),
		data:   bytes.NewBuffer(nil),
	}
	enc := sixel.NewEncoder(img.data)
	if err := enc.Encode(raw); err != nil {
		s.Fini()
		log.Fatal(err)
	}

	lock := true
	displayHelloWorld(s)
	displaySixel(s, img, lock)

	for {
		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			displayHelloWorld(s)
			displaySixel(s, img, lock)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEnter {
				lock = !lock
				displaySixel(s, img, lock)
			}
		}
	}
}
