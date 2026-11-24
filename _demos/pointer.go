//go:build ignore
// +build ignore

// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
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

	"github.com/gdamore/tcell/v3"
)

// The names of the CSS cursor property.  Terminals that support OSC 22 differ
// in which of these they understand.
var shapes = []string{
	"default", "pointer", "text", "vertical-text", "crosshair", "help",
	"wait", "progress", "move", "grab", "grabbing", "not-allowed", "no-drop",
	"copy", "alias", "cell", "context-menu", "zoom-in", "zoom-out",
	"all-scroll", "col-resize", "row-resize", "ew-resize", "ns-resize",
	"nesw-resize", "nwse-resize", "n-resize", "e-resize", "s-resize",
	"w-resize", "ne-resize", "nw-resize", "se-resize", "sw-resize",
}

const (
	columnWidth = 16
	firstRow    = 2
)

func columns(s tcell.Screen) int {
	w, _ := s.Size()
	return max(1, w/columnWidth)
}

// shapeAt returns the name drawn at the given position, or "" if there is none.
func shapeAt(s tcell.Screen, x, y int) string {
	if y < firstRow || (y-firstRow)%2 != 0 || x < 0 {
		return ""
	}
	cols := columns(s)
	col := x / columnWidth
	if col >= cols {
		return ""
	}
	i := (y-firstRow)/2*cols + col
	if i >= len(shapes) || x%columnWidth >= len(shapes[i]) {
		return ""
	}
	return shapes[i]
}

func draw(s tcell.Screen, current string) {
	s.Clear()
	s.PutStr(0, 0, "Hover over a name to see that pointer shape. Press ESC to exit.")
	cols := columns(s)
	for i, name := range shapes {
		style := tcell.StyleDefault
		if name == current {
			style = style.Reverse(true)
		}
		s.PutStrStyled(i%cols*columnWidth, firstRow+i/cols*2, name, style)
	}
	s.Show()
}

// This program shows the mouse pointer in the shape whose name it hovers over.
func main() {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	s.EnableMouse()

	current := ""
	draw(s, current)
	for ev := range s.EventQ() {
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventMouse:
			x, y := ev.Position()
			if shape := shapeAt(s, x, y); shape != current {
				current = shape
				s.SetPointerShape(current)
			}
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
		}
		draw(s, current)
	}
}
