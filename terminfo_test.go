// Copyright 2016 The TCell Authors
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
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// This terminfo entry is a stripped down version from
// xterm-256color, but I've added some of my own entries.
var testTerminfo = &Terminfo{
	Name:      "simulation_test",
	Columns:   80,
	Lines:     24,
	Colors:    256,
	Bell:      "\a",
	Blink:     "\x1b2ms$<2>",
	Reverse:   "\x1b[7m",
	SetFg:     "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m",
	SetBg:     "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m",
	AltChars:  "``aaffggiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
	Mouse:     "\x1b[M",
	MouseMode: "%?%p1%{1}%=%t%'h'%Pa%e%'l'%Pa%;\x1b[?1000%ga%c\x1b[?1003%ga%c\x1b[?1006%ga%c",
	SetCursor: "\x1b[%i%p1%d;%p2%dH",
	PadChar:   "\x00",
}

func TestTerminfo(t *testing.T) {

	ti := testTerminfo

	Convey("Terminfo parameter processing", t, func() {
		// This tests %i, and basic parameter strings too
		Convey("TGoto works", func() {
			s := ti.TGoto(7, 9)
			So(s, ShouldEqual, "\x1b[10;8H")
		})

		// This tests some conditionals
		Convey("TParm extended formats work", func() {
			s := ti.TParm("A[%p1%2.2X]B", 47)
			So(s, ShouldEqual, "A[2F]B")
		})

		// This tests some conditionals
		Convey("TParm colors work", func() {
			s := ti.TParm(ti.SetFg, 7)
			So(s, ShouldEqual, "\x1b[37m")

			s = ti.TParm(ti.SetFg, 15)
			So(s, ShouldEqual, "\x1b[97m")

			s = ti.TParm(ti.SetFg, 200)
			So(s, ShouldEqual, "\x1b[38;5;200m")
		})

		// This tests variables
		Convey("TParm mouse mode works", func() {
			s := ti.TParm(ti.MouseMode, 1)
			So(s, ShouldEqual, "\x1b[?1000h\x1b[?1003h\x1b[?1006h")
			s = ti.TParm(ti.MouseMode, 0)
			So(s, ShouldEqual, "\x1b[?1000l\x1b[?1003l\x1b[?1006l")
		})

	})

	Convey("Terminfo delay handling", t, func() {

		Convey("19200 baud", func() {
			buf := bytes.NewBuffer(nil)
			ti.TPuts(buf, ti.Blink, 19200)
			s := string(buf.Bytes())
			So(s, ShouldEqual, "\x1b2ms\x00\x00\x00\x00")
		})

		Convey("50 baud", func() {
			buf := bytes.NewBuffer(nil)
			ti.TPuts(buf, ti.Blink, 50)
			s := string(buf.Bytes())
			So(s, ShouldEqual, "\x1b2ms")
		})
	})
}

func BenchmarkSetFgBg(b *testing.B) {
	ti := testTerminfo

	for i := 0; i < b.N; i++ {
		ti.TParm(ti.SetFg, 100, 200)
		ti.TParm(ti.SetBg, 100, 200)
	}
}
