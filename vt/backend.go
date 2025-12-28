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

	// Put a single character at a specific position with the given attributes.
	// The cursor should be moved to this location.  It may advance.  If the value
	// of the rune is zero, then the cell is being erased, and no content should be
	// displayed.
	PutAbs(Coord, rune, Attr)

	// GetPosition returns the cursor position.
	GetPosition() Coord

	// SetPosition sets the cursor position. If the position is out of bounds,
	// it should be clipped to the window size.
	SetPosition(Coord)
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

// UnderlineColorer adds underline color management to the Colorer.
type UnderlineColorer interface {
	Colorer

	// SetUlColor sets the underline color.
	// If color.None is chosen, then default foreground color is used.
	SetUlColor(color.Color)
}

// Positioner tracks the cursor. Backends that implement this should automatically
// advance the position after writing content.
type Positioner interface {

	// Put the rune at the cursor position with the given attributes. The cursor is expected
	// to advance by the rune width.
	PutChar(rune, Attr)
}

// GraphemePositioner adds the ability to write a grapheme cluster to Positioner.
// If a terminal does not implement this, then it will not support grapheme clusters.
type GraphemePositioner interface {
	Positioner

	// Put the grapheme cluster at the cursor position. This will occupy a single
	// narrow or wide terminal cell.
	PutGrapheme([]rune)
}

// Resizer adds notifications when the window size changes.
type Resizer interface {
	// NotifyResize registers a channel to be posted to if the window size changes.
	NotifyResize(chan<- bool)
}
