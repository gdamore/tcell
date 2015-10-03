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

package tcell

import (
	"github.com/mattn/go-runewidth"
)

// Cell represents a single character cell.  This is primarily intended for
// use by Screen implementors.
type Cell struct {
	Ch    []rune
	Dirty bool
	Width uint8
	Style Style
}

// ClearCells clears the entire set of cells, making them all whitespace with
// the provided attribute.
func ClearCells(c []Cell, style Style) {
	for i := range c {
		c[i].Ch = nil
		c[i].Style = style
		c[i].Width = 1
		c[i].Dirty = true
	}
}

// InvalidateCells marks all cells in the array dirty.
func InvalidateCells(c []Cell) {
	for i := range c {
		c[i].Dirty = true
	}
}

// ResizeCells is used to create a new cells array, with different dimensions,
// while preserving the original contents.  The returned array may be the same
// as the original, if we can reuse it.  Hence, the old array should no longer
// be used by the caller after this call.  The cells will be marked dirty so
// that they can be redrawn.
func ResizeCells(oldc []Cell, oldw, oldh, neww, newh int) []Cell {

	if oldh == newh && oldw == neww {
		return oldc
	}
	newc := oldc

	// Probably are other conditions where we could reuse, but if there is
	// any doubt at all, its easier & safest to just realloc the window.
	if newh > oldh || neww > oldw {
		newc = make([]Cell, neww*newh)
	}
	for row := 0; row < newh && row < oldh; row++ {
		for col := 0; col < oldw && col < neww; col++ {
			newc[(row*neww)+col] = oldc[(row*oldw)+col]
			newc[(row*neww)+col].Dirty = true
		}
	}
	return newc
}

// SetCell writes the contents into the cell.  It ensures that at most one
// nonzero width rune is present in the Ch array (and if any zero width runes
// are present without a non-zero one, then a space is inserted), and updates
// the Dirty bit if the contents are different than they were.
func (c *Cell) SetCell(ch []rune, style Style) {

	c.PutChars(ch)
	c.PutStyle(style)
/*
	var mainc rune
	var width uint8
	var compc []rune

	width = 1
	mainc = ' '
	for _, r := range ch {
		if r < ' ' {
			// skip over non-printable control characters
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 1:
			mainc = r
			width = 1
		case 2:
			mainc = r
			width = 2
		case 0:
			compc = append(compc, r)
		}
	}

	newch := append([]rune{mainc}, compc...)
	if len(newch) != len(c.Ch) || style != c.Style || c.Dirty {
		c.Dirty = true
	} else {
		for i := range newch {
			if newch[i] != c.Ch[i] {
				c.Dirty = true
			}
		}
	}
	c.Ch = newch
	c.Style = style
	c.Width = width
*/
}

func (c *Cell) PutChars(ch []rune) {

	var mainc rune
	var width uint8
	var compc []rune

	width = 1
	mainc = ' '
	for _, r := range ch {
		if r < ' ' {
			// skip over non-printable control characters
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 1:
			mainc = r
			width = 1
		case 2:
			mainc = r
			width = 2
		case 0:
			compc = append(compc, r)
		}
	}

	newch := append([]rune{mainc}, compc...)
	if len(newch) != len(c.Ch) {
		c.Dirty = true
	} else {
		for i := range newch {
			if newch[i] != c.Ch[i] {
				c.Dirty = true
			}
		}
	}
	c.Ch = newch
	c.Width = width
}

func (c *Cell) PutChar(ch rune) {
	c.PutChars([]rune{ch})
}

func (c *Cell) PutStyle(style Style) {
	if c.Style != style {
		c.Style = style
		c.Dirty = true
	}
}
