// Copyright 2015 The Tcell Authors
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

package views

import (
	"github.com/gdamore/tcell"
)

// BoxLayout is a container Widget that lays out its child widgets in
// either a horizontal row or a vertical column.
type BoxLayout struct {
	view    View
	orient  Orientation
	style   tcell.Style // backing style
	cells   []*boxLayoutCell
	width   int
	height  int
	changed bool

	WidgetWatchers
}

type boxLayoutCell struct {
	widget Widget
	fill   float64 // fill factor - 0.0 means no expansion
	pad    int     // count of padding spaces (stretch)
	frac   float64 // calculated residual spacing, used internally
	view   *ViewPort
}

func (b *BoxLayout) layout() {
	if b.view == nil {
		return
	}
	w, h := b.view.Size()

	minx, miny, totx, toty := 0, 0, 0, 0
	totf := 0.0
	for _, c := range b.cells {
		x, y := c.widget.Size()
		totx += x
		toty += y
		totf += c.fill
		if x > minx {
			minx = x
		}
		if y > miny {
			miny = y
		}
	}

	extra := 0
	if b.orient == Horizontal {
		extra = w - totx
		b.width = totx
		b.height = miny
	} else {
		extra = h - toty
		b.width = minx
		b.height = toty
	}
	if extra < 0 {
		extra = 0
	}
	resid := extra
	if totf == 0 {
		resid = 0
	}
	for _, c := range b.cells {
		if c.fill > 0 {
			c.frac = float64(extra) * c.fill / totf
			c.pad = int(c.frac)
			c.frac -= float64(c.pad)
			resid -= c.pad
		} else {
			c.pad = 0
			c.frac = 0
		}
	}

	// Distribute any left over padding.  We try to give it to the
	// the cells with the highest residual fraction.  It should be
	// the case that no single cell gets more than one more cell.
	for resid > 0 {
		var best *boxLayoutCell = nil
		for _, c := range b.cells {
			if c.fill == 0 {
				continue
			}
			if best == nil || c.frac > best.frac {
				best = c
			}
		}
		best.pad++
		best.frac = 0
		resid--
	}

	x, y, xinc, yinc := 0, 0, 0, 0
	for _, c := range b.cells {
		cw, ch := c.widget.Size()

		switch b.orient {
		case Horizontal:
			xinc = cw + c.pad
			cw += c.pad
			ch = h

		case Vertical:
			yinc = ch + c.pad
			ch += c.pad
			cw = w

		default:
			panic("Bad orientation")

		}
		c.view.Resize(x, y, cw, ch)
		x += xinc
		y += yinc
	}
	b.changed = false
}

func (b *BoxLayout) Resize() {
	b.layout()

	// Now also let the children know we resized.
	for i := range b.cells {
		b.cells[i].widget.Resize()
	}
	b.PostEventWidgetResize(b)
}

// Draw is called to update the displayed content.
func (b *BoxLayout) Draw() {

	if b.view == nil {
		return
	}
	if b.changed {
		b.layout()
	}
	b.view.Clear()
	w, h := b.view.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			b.view.SetContent(x, y, ' ', nil, b.style)
		}
	}
	for i := range b.cells {
		b.cells[i].widget.Draw()
	}
}

// Size returns the preferred size in character cells (width, height).
func (b *BoxLayout) Size() (int, int) {
	return b.width, b.height
}

// SetView sets the View object used for the text bar.
func (b *BoxLayout) SetView(view View) {
	b.changed = true
	b.view = view
	for _, c := range b.cells {
		c.view.SetView(view)
	}
}

// HandleEvent implements a tcell.EventHandler.  The only events
// we care about are Widget change events from our children. We
// watch for those so that if the child changes, we can arrange
// to update our layout.
func (b *BoxLayout) HandleEvent(ev tcell.Event) bool {
	switch ev.(type) {
	case *EventWidgetContent:
		b.changed = true
		return true
	}
	return false
}

// Add() adds a widget to the end of the BoxLayout.
func (b *BoxLayout) AddWidget(widget Widget, fill float64) {
	c := &boxLayoutCell{
		widget: widget,
		fill:   fill,
		view:   NewViewPort(b.view, 0, 0, 0, 0),
	}
	c.widget.SetView(c.view)
	b.cells = append(b.cells, c)
	b.changed = true
	widget.Watch(b)
	b.PostEventWidgetContent(b)
}

func (b *BoxLayout) RemoveWidget(widget Widget) {
	for i := 0; i < len(b.cells); i++ {
		if b.cells[i].widget == widget {
			b.cells = append(b.cells[:i-1], b.cells[i+1:]...)
			return
		}
	}
	b.changed = true
	widget.Unwatch(b)
	b.PostEventWidgetContent(b)
}

func (b *BoxLayout) SetOrientation(orient Orientation) {
	b.orient = orient
	b.changed = true
	b.PostEventWidgetContent(b)
}

// SetStyle sets the style used.
func (b *BoxLayout) SetStyle(style tcell.Style) {
	b.style = style
	b.PostEventWidgetContent(b)
}

// NewBoxLayout creates an empty BoxLayout.
func NewBoxLayout(orient Orientation) *BoxLayout {
	return &BoxLayout{orient: orient}
}
