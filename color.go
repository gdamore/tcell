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

// Color represents a color.  The low numeric values are the same as used
// by ECMA-48, and beyond that XTerm.  A 24-bit RGB value may be used by
// adding in the ColorIsRGB flag.  For Color names we use the W3C approved
// color names.
//
// Note that on various terminals colors may be approximated however, or
// not supported at all.  If no suitable representation for a color is known,
// the library will simply not set any color, deferring to whatever default
// attributes the terminal uses.
type Color int32

const (
	// ColorDefault is used to leave the Color unchanged from whatever
	// system or teminal default may exist.
	ColorDefault Color = -1

	// ColorIsRGB is used to indicate that the numeric value is not
	// a known color constant, but rather an RGB value.  The lower
	// order 3 bytes are RGB.
	ColorIsRGB Color = 1 << 24
)

// Note that the order of these options is important -- it follows the
// definitions used by ECMA and XTerm.  Hence any further named colors
// must begin at a value not less than 256.
const (
	ColorBlack Color = iota
	ColorMaroon
	ColorGreen
	ColorOlive
	ColorNavy
	ColorPurple
	ColorTeal
	ColorSilver
	ColorGrey
	ColorRed
	ColorLime
	ColorYellow
	ColorBlue
	ColorFuchsia
	ColorAqua
	ColorWhite
)
