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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStyle(t *testing.T) {
	Convey("Style checks", t, WithScreen(t, "", func(s SimulationScreen) {

		style := StyleDefault
		fg, bg, attr := style.Decompose()

		So(fg, ShouldEqual, ColorDefault)
		So(bg, ShouldEqual, ColorDefault)
		So(attr, ShouldEqual, AttrNone)

		s2 := style.
			Background(ColorRed).
			Foreground(ColorBlue).
			Blink(true)

		fg, bg, attr = s2.Decompose()
		So(fg, ShouldEqual, ColorBlue)
		So(bg, ShouldEqual, ColorRed)
		So(attr, ShouldEqual, AttrBlink)
	}))
}
