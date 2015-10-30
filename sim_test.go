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

func WithScreen(t *testing.T, charset string, fn func(s SimulationScreen)) func() {
	return func() {
		s := NewSimulationScreen(charset)
		So(s, ShouldNotBeNil)
		e := s.Init()
		So(e, ShouldBeNil)
		Reset(func() {
			s.Fini()
		})
		fn(s)
	}
}

func TestInitScreen(t *testing.T) {

	Convey("Init a screen", t, WithScreen(t, "", func(s SimulationScreen) {

		Convey("Size should be valid", func() {
			x, y := s.Size()
			So(x, ShouldEqual, 80)
			So(y, ShouldEqual, 25)
		})

		Convey("Default charset is UTF-8", func() {
			So(s.CharacterSet(), ShouldEqual, "UTF-8")
		})

		Convey("Backing size is correct", func() {
			b, x, y := s.GetContents()
			So(b, ShouldNotBeNil)
			So(x, ShouldEqual, 80)
			So(y, ShouldEqual, 25)
			So(len(b), ShouldEqual, x*y)
		})
	}))

}

func TestClearScreen(t *testing.T) {
	Convey("Clear screen", t, WithScreen(t, "", func(s SimulationScreen) {

		s.Clear()
		b, x, y := s.GetContents()
		So(b, ShouldNotBeNil)
		So(x, ShouldEqual, 80)
		So(y, ShouldEqual, 25)
		So(len(b), ShouldEqual, x*y)
		s.Sync()

		nmatch := 0
		for i := 0; i < x*y; i++ {
			if len(b[i].Runes) == 1 && b[i].Runes[0] == ' ' {
				nmatch++
			}
		}
		So(nmatch, ShouldEqual, x*y)

		nmatch = 0
		for i := 0; i < x*y; i++ {
			if len(b[i].Bytes) == 1 && b[i].Bytes[0] == ' ' {
				nmatch++
			}
		}
		So(nmatch, ShouldEqual, x*y)

		nmatch = 0
		for i := 0; i < x*y; i++ {
			if b[i].Style == StyleDefault {
				nmatch++
			}
		}
		So(nmatch, ShouldEqual, x*y)
	}))
}

func TestSetCell(t *testing.T) {
	st := StyleDefault.Background(ColorRed).Blink(true)
	Convey("Set contents", t, WithScreen(t, "", func(s SimulationScreen) {
		s.SetCell(2, 5, st, '@')
		b, x, y := s.GetContents()
		So(len(b), ShouldEqual, x*y)
		So(x, ShouldEqual, 80)
		So(y, ShouldEqual, 25)
		s.Show()

		sc := &b[5*80+2]
		So(len(sc.Runes), ShouldEqual, 1)
		So(len(sc.Bytes), ShouldEqual, 1)
		So(sc.Bytes[0], ShouldEqual, '@')
		So(sc.Runes[0], ShouldEqual, '@')
		So(sc.Style, ShouldEqual, st)
	}))
}

func TestResize(t *testing.T) {
	st := StyleDefault.Background(ColorYellow).Underline(true)
	Convey("Resize", t, WithScreen(t, "", func(s SimulationScreen) {
		s.SetCell(2, 5, st, '&')
		b, x, y := s.GetContents()
		So(len(b), ShouldEqual, x*y)
		So(x, ShouldEqual, 80)
		So(y, ShouldEqual, 25)
		s.Show()

		sc := &b[5*80+2]
		So(len(sc.Runes), ShouldEqual, 1)
		So(len(sc.Bytes), ShouldEqual, 1)
		So(sc.Bytes[0], ShouldEqual, '&')
		So(sc.Runes[0], ShouldEqual, '&')
		So(sc.Style, ShouldEqual, st)

		Convey("Do resize", func() {
			s.SetSize(30, 10)
			s.Show()
			b2, x2, y2 := s.GetContents()
			So(b2, ShouldNotEqual, b)
			So(x2, ShouldEqual, 30)
			So(y2, ShouldEqual, 10)

			sc2 := &b[5*80+2]
			So(len(sc2.Runes), ShouldEqual, 1)
			So(len(sc2.Bytes), ShouldEqual, 1)
			So(sc2.Bytes[0], ShouldEqual, '&')
			So(sc2.Runes[0], ShouldEqual, '&')
			So(sc2.Style, ShouldEqual, st)
		})
	}))
}
