// Copyright 2026 The TCell Authors
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

package vt

import "github.com/gdamore/tcell/v3/color"

// Backend describes the backend of a terminal.
// This can be used to create a real emulator, while allowing the processor
// front end to handle the common details of parsing escape sequences, the state
// machine, and so forth. Backends support a limited set of common functionality,
// including a cursor. They only need to support writing at the cursor.
type Backend interface {

	// GetPrivateMode returns the status of a given private mode.
	GetPrivateMode(PrivateMode) ModeStatus

	// SetPrivateMode sets a private mode to the given status.
	// If either value is invalid, this should simply ignore the operation.
	SetPrivateMode(PrivateMode, ModeStatus) error

	// GetSize returns the size of the terminal in characters.
	// The X and Y are counts, so the bottom right cell should be at coordinate (X-1, Y-1).
	GetSize() Coord

	// SetAttr is used to set the current attributes. It will be used for subsequent writes.
	SetAttr(Attr)

	// Put a single rune at a specific position, using the current attributes and colors.
	// If the rune is 0, then this is an erase and no content should be displayed.
	// The display width is the size in cells for this rune.
	PutRune(Coord, rune, int)

	// Put a grapheme cluster at a specific location, using the current attributes and colors.
	// If the string is empty, then this is an erase and no content should be displayed.
	// The width is supplied as the last parameter, and represents the number character cells expected
	// to be consumed for this grapheme cluster.
	PutGrapheme(Coord, string, int)

	// GetPosition returns the cursor position.
	GetPosition() Coord

	// SetPosition sets the cursor position. If the position is out of bounds,
	// it should be clipped to the window size.
	SetPosition(Coord)

	// Reset resets the terminal to default state.
	Reset()
}

// Beeper can be implemented by a backend to indicate it can ring the bell or beep.
// This is typically done in response to a 0x07 bell.
type Beeper interface {
	Beep()
}

// Colorer can select the colors used.  This interface is stateful, so that
// an implementation needs to remember the values and use them.
type Colorer interface {
	// Colors returns the number of colors this terminal can support.  For direct color,
	// return 1<<24. The XTerm palette is assumed. Monochrome terminals should return 0.
	Colors() int

	// SetFgColor sets the foreground color.
	SetFgColor(color.Color)

	// SetBgColor sets the background color.
	SetBgColor(color.Color)
}

// Resizer adds notifications when the window size changes.
type Resizer interface {
	// NotifyResize registers a channel to be posted to if the window size changes.
	NotifyResize(chan<- bool)
}

// Titler adds support for setting the window title. (Typically this is OSC2.)
// Note that for security reasons we only support setting this.
// We don't bother with icon titles, since few terminal emulators support it, and it
// would be hard for us to do this in any portable fashion.
type Titler interface {
	// SetWindowTitle only changes the window title.
	SetWindowTitle(string)
}

// Blitter implements a cell-level blit, where a rectangular range of cells is copied from one
// location to another.  The source and destination may overlap.  The old locations will remain
// unchanged except of course or cells overwritten by the blit. The content will also be clipped
// to the visible dimensions.
type Blitter interface {
	Blit(src, dst, dim Coord)
}
