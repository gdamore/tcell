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

// Screen represents the physical (or emulated) screen.
// This can be a terminal window or a physical console.  Platforms implement
// this differerently.
type Screen interface {
	// Init initializes the screen for use.
	Init() error

	// Fini finalizes the screen also releasing resources.
	Fini()

	// Clear erases the screen.  The contents of any screen buffers
	// will also be cleared.  This has the logical effect of
	// filling the screen with spaces, using the global default style.
	Clear()

	// SetCell sets the cell at the given location.
	// The ch list contains at most one rune of width > 0, and the
	// runes with zero width (combining marks) must follow the first
	// non-zero width character.  (If only combining marks are present,
	// a space character will be filled in.)
	//
	// Note that double wide runes occupy two cells, and attempts to
	// place a character at the immediately adjacent cell will have
	// undefined effects.  Double wide runes that are printed in the
	// last column will be replaced with a single width space on output.
	//
	// SetCell may change the cursor location.  Callers should explictly
	// save and restore cursor state if neccesary.  The cursor visibility
	// is not affected, so callers probably should hide the cursor when
	// calling this.
	//
	// Note that the results will not be visible until either Show() or
	// Sync() are called.
	SetCell(x int, y int, style Style, ch ...rune)

	// GetContent returns the contents at the given location.  If the
	// coordinates are out of range, then the values will be 0, nil,
	// StyleDefault.  Note that the contents returned are logical contents
	// and may not actually be what is displayed, but rather are what will
	// be displayed if Show() or Sync() is called.  The width is the width
	// in screen cells - this should either be 1 or 2.
	GetContent(x, y int) (mainc rune, combc []rune, style Style, width int)

	// SetContent sets the contents of the given cell location.  If
	// the coordinates are out of range, then the operation is ignored.
	// The first rune is the primary non-zero width rune.  The array
	// that follows is a possible list of combining characters to append.
	// The results are not displayd until Show() or Sync() is called.
	SetContent(x int, y int, mainc rune, combc []rune, style Style)

	// SetStyle sets the default style to use when clearing the screen
	// or when StyleDefault is specified.  If it is also StyleDefault,
	// then whatever system/terminal default is relevant will be used.
	SetStyle(style Style)

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
	// use the ANSI color map.  If a terminal is monochrome, it will
	// return 0.
	Colors() int

	// Show takes any output that was deferred due to buffering, and
	// flushes it to the physical display.  It does so in the most
	// efficient and least visually disruptive manner possible.
	Show()

	// Sync works like Show(), but it updates every visible cell on the
	// physical display, assuming that it is not synchronized with any
	// internal model.  This may be both expensive and visually jarring,
	// so it should only be used when believed to actually be necessary.
	// Typically this is called as a result of a user-requested redraw
	// (e.g. to clear up on screen corruption caused by some other program),
	// or during a resize event.
	Sync()

	// CharacterSet() returns information about the character set.
	// This isn't the full locale, but it does give us the input/ouput
	// character set.  Note that this is just for diagnostic purposes,
	// we normally translate input/output to/from UTF-8, regardless of
	// what the user's environment is.
	CharacterSet() string
}

// NewScreen returns a default Screen suitable for the user's terminal
// environment.
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
