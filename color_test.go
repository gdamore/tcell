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

func TestColor(t *testing.T) {
	pal := []Color{}
	for i := 0; i < 255; i++ {
		pal = append(pal, Color(i))
	}
	Convey("Color values are correct", t, func() {
		So(ColorRed.Hex(), ShouldEqual, 0x00FF0000)
		So(ColorGreen.Hex(), ShouldEqual, 0x00008000)
		So(ColorLime.Hex(), ShouldEqual, 0x0000FF00)
		So(ColorBlue.Hex(), ShouldEqual, 0x000000FF)
		So(ColorBlack.Hex(), ShouldEqual, 0x00000000)
		So(ColorWhite.Hex(), ShouldEqual, 0x00FFFFFF)
		So(ColorSilver.Hex(), ShouldEqual, 0x00C0C0C0)
	})

	Convey("Color fitting from 16 colors to 8 colors works", t, func() {
		for i := 0; i < 7; i++ {
			So(FindColor(Color(i), pal[:8]), ShouldEqual, Color(i))
		}
		// Grey is closest to Silver
		So(FindColor(Color(8), pal[:8]), ShouldEqual, Color(7))

		for i := 9; i < 16; i++ {
			So(FindColor(Color(i), pal[:8]), ShouldEqual, Color(i%8))
		}
	})

	Convey("Color lookups by name work", t, func() {
		So(GetColor("red"), ShouldEqual, ColorRed)
		So(GetColor("#FF0000").Hex(), ShouldEqual, ColorRed.Hex())
		So(GetColor("black"), ShouldEqual, ColorBlack)
		So(GetColor("orange"), ShouldEqual, ColorOrange)
	})

	Convey("Color imperfect fit works", t, func() {
		So(FindColor(ColorOrangeRed, pal[:16]), ShouldEqual, ColorRed)
		So(FindColor(ColorAliceBlue, pal[:16]), ShouldEqual, ColorWhite)
		So(FindColor(ColorPink, pal), ShouldEqual, Color217)
		So(FindColor(ColorSienna, pal), ShouldEqual, Color173)
		So(FindColor(GetColor("#00FD00"), pal), ShouldEqual, ColorLime)
	})

	Convey("Color RGB breakdown works", t, func() {
		r, g, b := GetColor("#112233").RGB()
		So(r, ShouldEqual, 0x11)
		So(g, ShouldEqual, 0x22)
		So(b, ShouldEqual, 0x33)
	})
}
