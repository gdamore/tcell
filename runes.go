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

// The names of these constants are chosen to match Terminfo names,
// modulo case, and changing the prefix from ACS_ to Rune.  These are
// the runes we provide extra special handling for, with ASCII fallbacks
// for terminals that lack them.

const (
	RuneSterling = '£'
	RuneDArrow   = '↓'
	RuneLArrow   = '←'
	RuneRArrow   = '→'
	RuneUArrow   = '↑'
	RuneBullet   = '·'
	RuneBoard    = '░'
	RuneCkBoard  = '▒'
	RuneDegree   = '°'
	RuneDiamond  = '◆'
	RuneGEqual   = '≥'
	RunePi       = 'π'
	RuneHLine    = '─'
	RuneLantern  = '§'
	RunePlus     = '┼'
	RuneLEqual   = '≤'
	RuneLLCorner = '└'
	RuneLRCorner = '┘'
	RuneNEqual   = '≠'
	RunePlMinus  = '±'
	RuneS1       = '⎺'
	RuneS3       = '⎻'
	RuneS7       = '⎼'
	RuneS9       = '⎽'
	RuneBlock    = '█'
	RuneTTee     = '┬'
	RuneRTee     = '┤'
	RuneLTee     = '├'
	RuneBTee     = '┴'
	RuneULCorner = '┌'
	RuneURCorner = '┐'
	RuneVLine    = '│'
)
