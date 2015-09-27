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

// Screen represents the physical (or emulated) screen, without any buffering.
// This can be a terminal window or a physical console.  Platforms implement
// this differerently.  Applications are unlikely to interface directly with
// this.
type Screen interface {
	// Init initializes the screen for use.
	Init() error

	// Fini finalizes the screen also releasing resources.
	Fini()

	// Clear erases the screen.
	Clear()

	// SetCell sets the cell at the given location.
	// The ch array contains at most one rune of width > 0, and the
	// runes with zero width (combining marks) must follow the first
	// non-zero width character.  (If only combining marks are present,
	// a space character will be filled in.)
	//
	// Note that double wide runes occupy two cells, and attempts to
	// place a character at the immediately adjacent cell will have
	// undefined effects.  Double wide runes that are printed in the
	// last column will be replaced with a single width space on output.
	//
	// Note that unlike the higher level interfaces, this operates
	// immediately, without any buffering.
	//
	// SetCell may change the cursor location.  Callers should explictly
	// save and restore cursor state if neccesary.  The cursor visibility
	// is not affected, so callers probably should hide the cursor when
	// calling this.
	SetCell(x int, y int, style Style, ch ...rune)

	// ShowCursor is used to display the cursor at a given location.
	// If the coordinates -1, -1 are given or are otherwise outside the
	// dimensions of the screen, the cursor will be hidden.
	ShowCursor(x int, y int)

	// HideCursor is used to hide the cursor.  Its an alias for
	// ShowCursor(-1, -1).
	HideCursor()

	// Size returns the screen size as width, height.  This changes in
	// response to a call to Clear or Flush.
	Size() (int, int)

	// PollEvent waits for events to arrive.  Main application loops
	// must spin on this to prevent the application from stalling.
	// Furthermore, this will return nil if the Screen is finalized.
	PollEvent() Event

	// PostEvent posts an event into the event stream.
	PostEvent(Event)

	// EnableMouse enables the mouse.  (If your terminal supports it.)
	EnableMouse()

	// DisableMouse disables the mouse.
	DisableMouse()

	// Colors returns the number of colors.  All colors are assumed to
	// use the ANSI color map.
	Colors() int
}

func NewScreen() (Screen, error) {
	// First we attempt to obtain a terminfo screen.  This should work
	// in most places if $TERM is set.
	if s, e := NewTerminfoScreen(); s != nil {
		return s, nil

	} else if s, _ := NewConsoleScreen(); s != nil {
		return s, nil

	} else {
		return nil, e
	}
}
