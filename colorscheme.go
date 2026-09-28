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

package tcell

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v3/color"
)

// EventColorScheme reports whether the terminal uses a dark or a light
// color scheme. It is sent after EnableColorScheme, if the terminal supports
// color scheme reports (DEC private mode 2031): once for the current scheme,
// and again whenever it changes, for example when the system switches
// between light and dark.
type EventColorScheme struct {
	EventTime

	// Dark is true for a dark color scheme, false for a light one.
	Dark bool
}

// NewEventColorScheme returns an EventColorScheme.
func NewEventColorScheme(dark bool) *EventColorScheme {
	ev := &EventColorScheme{Dark: dark}
	ev.SetEventNow()
	return ev
}

// EventBackgroundColor reports the terminal's default background color, in
// reply to GetBackgroundColor.
type EventBackgroundColor struct {
	EventTime

	// Color is the background color, as an RGB color.
	Color color.Color
}

// NewEventBackgroundColor returns an EventBackgroundColor.
func NewEventBackgroundColor(c color.Color) *EventBackgroundColor {
	ev := &EventBackgroundColor{Color: c}
	ev.SetEventNow()
	return ev
}

// Dark reports whether the background is dark: whether its luma, the
// Rec. 709 weighted sum of its sRGB components, is at most one half. This
// puts the boundary at mid-gray, as perceived. (Weighting linearized
// components, WCAG's relative luminance, would put it near #bcbcbc and call
// light-gray backgrounds dark.)
func (ev *EventBackgroundColor) Dark() bool {
	r, g, b := ev.Color.RGB()
	return 0.2126*float64(r)+0.7152*float64(g)+0.0722*float64(b) <= 0.5*255
}

// parseXColor parses a color as terminals report it in reply to OSC color
// queries: "rgb:RRRR/GGGG/BBBB", with one to four hex digits for each
// component, or "rgba:RRRR/GGGG/BBBB/AAAA", whose alpha component must be
// valid too but is ignored.
func parseXColor(s string) (color.Color, bool) {
	var parts []string
	want := 0
	if rest, ok := strings.CutPrefix(s, "rgb:"); ok {
		parts, want = strings.Split(rest, "/"), 3
	} else if rest, ok := strings.CutPrefix(s, "rgba:"); ok {
		parts, want = strings.Split(rest, "/"), 4
	}
	if want == 0 || len(parts) != want {
		return color.Default, false
	}
	var rgb [4]int32
	for i, p := range parts {
		if len(p) < 1 || len(p) > 4 {
			return color.Default, false
		}
		v, err := strconv.ParseUint(p, 16, 16)
		if err != nil {
			return color.Default, false
		}
		// Scale to 8 bits: "f", "ff", "fff" and "ffff" are all full.
		rgb[i] = int32(v * 255 / (1<<(4*len(p)) - 1))
	}
	return color.NewRGBColor(rgb[0], rgb[1], rgb[2]), true
}
