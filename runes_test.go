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

func TestCanDisplay(t *testing.T) {
	Convey("With a UTF-8 screen", t,
		WithScreen(t, "UTF-8", func(s SimulationScreen) {
			So(s.CharacterSet(), ShouldEqual, "UTF-8")
			So(s.CanDisplay('a', true), ShouldBeTrue)
			So(s.CanDisplay(RuneHLine, true), ShouldBeTrue)
			So(s.CanDisplay(RuneHLine, false), ShouldBeTrue)
			So(s.CanDisplay('⌀', false), ShouldBeTrue)
		}))

	Convey("With an ASCII screen", t,
		WithScreen(t, "US-ASCII", func(s SimulationScreen) {
			So(s.CharacterSet(), ShouldEqual, "US-ASCII")
			So(s.CanDisplay('a', true), ShouldBeTrue)
			So(s.CanDisplay(RuneHLine, true), ShouldBeTrue)
			So(s.CanDisplay(RunePi, false), ShouldBeFalse)
			So(s.CanDisplay('⌀', false), ShouldBeFalse)
		}))

}

func TestRegisterFallback(t *testing.T) {
	Convey("With an ASCII screen", t,
		WithScreen(t, "US-ASCII", func(s SimulationScreen) {
			So(s.CharacterSet(), ShouldEqual, "US-ASCII")
			s.RegisterRuneFallback('⌀', "o")
			So(s.CanDisplay('⌀', false), ShouldBeFalse)
			So(s.CanDisplay('⌀', true), ShouldBeTrue)
			s.UnregisterRuneFallback('⌀')
			So(s.CanDisplay('⌀', false), ShouldBeFalse)
			So(s.CanDisplay('⌀', true), ShouldBeFalse)
		}))

}

func TestUnregisterFallback(t *testing.T) {
	Convey("With an ASCII screen (HLine)", t,
		WithScreen(t, "US-ASCII", func(s SimulationScreen) {
			So(s.CharacterSet(), ShouldEqual, "US-ASCII")
			So(s.CanDisplay(RuneHLine, true), ShouldBeTrue)
			s.UnregisterRuneFallback(RuneHLine)
			So(s.CanDisplay(RuneHLine, true), ShouldBeFalse)
		}))
}
