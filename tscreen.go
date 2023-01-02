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

package tcell

import (
	// import the stock terminals
	_ "github.com/gdamore/tcell/v2/terminfo/base"
)

// tKeyCode represents a combination of a key code and modifiers.
type tKeyCode struct {
	key Key
	mod ModMask
}

func (t *tScreen) SetStyle(style Style) {
	t.Lock()
	if !t.fini {
		t.style = style
	}
	t.Unlock()
}

func (t *tScreen) Clear() {
	t.Fill(' ', t.style)
	t.Lock()
	t.clear = true
	w, h := t.cells.Size()
	// because we are going to clear (see t.clear) in the next cycle,
	// let's also unmark the dirty bit so that we don't waste cycles
	// drawing things that are already dealt with via the clear escape sequence.
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			t.cells.SetDirty(col, row, false)
		}
	}
	t.Unlock()
}

func (t *tScreen) Fill(r rune, style Style) {
	t.Lock()
	if !t.fini {
		t.cells.Fill(r, style)
	}
	t.Unlock()
}

func (t *tScreen) EnableMouse(flags ...MouseFlags) {
	var f MouseFlags
	flagsPresent := false
	for _, flag := range flags {
		f |= flag
		flagsPresent = true
	}
	if !flagsPresent {
		f = MouseMotionEvents | MouseDragEvents | MouseButtonEvents
	}

	t.Lock()
	t.mouseFlags = f
	t.enableMouse(f)
	t.Unlock()
}

func (t *tScreen) DisableMouse() {
	t.Lock()
	t.mouseFlags = 0
	t.enableMouse(0)
	t.Unlock()
}

func (t *tScreen) SetContent(x, y int, mainc rune, combc []rune, style Style) {
	t.Lock()
	if !t.fini {
		t.cells.SetContent(x, y, mainc, combc, style)
	}
	t.Unlock()
}

func (t *tScreen) GetContent(x, y int) (rune, []rune, Style, int) {
	t.Lock()
	mainc, combc, style, width := t.cells.GetContent(x, y)
	t.Unlock()
	return mainc, combc, style, width
}

func (t *tScreen) SetCell(x, y int, style Style, ch ...rune) {
	if len(ch) > 0 {
		t.SetContent(x, y, ch[0], ch[1:], style)
	} else {
		t.SetContent(x, y, ' ', nil, style)
	}
}

func (t *tScreen) HideCursor() {
	t.ShowCursor(-1, -1)
}

func (t *tScreen) Show() {
	t.Lock()
	if !t.fini {
		t.resize()
		t.draw()
	}
	t.Unlock()
}

func (t *tScreen) EnablePaste() {
	t.Lock()
	t.pasteEnabled = true
	t.enablePasting(true)
	t.Unlock()
}

func (t *tScreen) DisablePaste() {
	t.Lock()
	t.pasteEnabled = false
	t.enablePasting(false)
	t.Unlock()
}

func (t *tScreen) Size() (int, int) {
	t.Lock()
	w, h := t.w, t.h
	t.Unlock()
	return w, h
}

func (t *tScreen) ChannelEvents(ch chan<- Event, quit <-chan struct{}) {
	defer close(ch)
	for {
		select {
		case <-quit:
			return
		case <-t.quit:
			return
		case ev := <-t.evch:
			select {
			case <-quit:
				return
			case <-t.quit:
				return
			case ch <- ev:
			}
		}
	}
}

func (t *tScreen) PollEvent() Event {
	select {
	case <-t.quit:
		return nil
	case ev := <-t.evch:
		return ev
	}
}

func (t *tScreen) HasPendingEvent() bool {
	return len(t.evch) > 0
}

func (t *tScreen) PostEventWait(ev Event) {
	t.evch <- ev
}

func (t *tScreen) PostEvent(ev Event) error {
	select {
	case t.evch <- ev:
		return nil
	default:
		return ErrEventQFull
	}
}

func (t *tScreen) clip(x, y int) (int, int) {
	w, h := t.cells.Size()
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x > w-1 {
		x = w - 1
	}
	if y > h-1 {
		y = h - 1
	}
	return x, y
}

func (t *tScreen) Sync() {
	t.Lock()
	t.cx = -1
	t.cy = -1
	if !t.fini {
		t.resize()
		t.clear = true
		t.cells.Invalidate()
		t.draw()
	}
	t.Unlock()
}

func (t *tScreen) RegisterRuneFallback(orig rune, fallback string) {
	t.Lock()
	t.fallback[orig] = fallback
	t.Unlock()
}

func (t *tScreen) UnregisterRuneFallback(orig rune) {
	t.Lock()
	delete(t.fallback, orig)
	t.Unlock()
}
