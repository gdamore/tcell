// Copyright 2018 The TCell Authors
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

package terminfo

import (
	"bytes"
	"os"
	"testing"
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

func TestTerminfoExpansion(t *testing.T) {

	ti := testTerminfo

	// Tests %i and basic parameter strings too
	if ti.TGoto(7, 9) != "\x1b[10;8H" {
		t.Error("TGoto expansion failed")
	}

	// This tests some conditionals
	if ti.TParm("A[%p1%2.2X]B", 47) != "A[2F]B" {
		t.Error("TParm conditionals failed")
	}

	// Color tests.
	if ti.TParm(ti.SetFg, 7) != "\x1b[37m" {
		t.Error("SetFg(7) failed")
	}
	if ti.TParm(ti.SetFg, 15) != "\x1b[97m" {
		t.Error("SetFg(15) failed")
	}
	if ti.TParm(ti.SetFg, 200) != "\x1b[38;5;200m" {
		t.Error("SetFg(200) failed")
	}

	if ti.TParm(ti.MouseMode, 1) != "\x1b[?1000h\x1b[?1003h\x1b[?1006h" {
		t.Error("Enable mouse mode failed")
	}
	if ti.TParm(ti.MouseMode, 0) != "\x1b[?1000l\x1b[?1003l\x1b[?1006l" {
		t.Error("Disable mouse mode failed")
	}
}

func TestTerminfoBaud19200(t *testing.T) {
	ti := testTerminfo
	buf := bytes.NewBuffer(nil)
	ti.TPuts(buf, ti.Blink, 19200)
	s := string(buf.Bytes())
	if s != "\x1b2ms\x00\x00\x00\x00" {
		t.Error("1920 baud failed")
	}
}
func TestTerminfoBaud50(t *testing.T) {
	ti := testTerminfo
	buf := bytes.NewBuffer(nil)
	ti.TPuts(buf, ti.Blink, 50)
	s := string(buf.Bytes())
	if s != "\x1b2ms" {
		t.Error("50 baud failed")
	}
}

func TestTerminfoBasic(t *testing.T) {

	os.Setenv("TCELLDB", "testdata/test1")
	ti, err := LookupTerminfo("test1")
	if ti == nil || err != nil || ti.Columns != 80 {
		t.Errorf("Failed test1 lookup: %v", err)
	}

	ti, err = LookupTerminfo("alias1")
	if ti == nil || err != nil || ti.Columns != 80 {
		t.Errorf("Failed alias1 lookup: %v", err)
	}

	os.Setenv("TCELLDB", "testdata")
	ti, err = LookupTerminfo("test2")
	if ti == nil || err != nil || ti.Columns != 80 {
		t.Errorf("Failed test2 lookup: %v", err)
	}
	if len(ti.Aliases) != 1 || ti.Aliases[0] != "alias2" {
		t.Errorf("Alias for test2 wrong")
	}
}

func TestTerminfoBadName(t *testing.T) {

	os.Setenv("TCELLDB", "testdata")
	if ti, err := LookupTerminfo("test3"); err == nil || ti != nil {
		t.Error("Bad name should not have resolved")
	}
}

func TestTerminfoLoop(t *testing.T) {
	os.Setenv("TCELLDB", "testdata")
	if ti, err := LookupTerminfo("loop1"); ti != nil || err == nil {
		t.Error("Loop loop1 should not have resolved")
	}
}

func TestTerminfoGzip(t *testing.T) {

	os.Setenv("TCELLDB", "testdata")
	if ti, err := LookupTerminfo("test-gzip"); ti == nil || err != nil ||
		ti.Columns != 80 {
		t.Error("test-gzip filed")
	}

	if ti, err := LookupTerminfo("alias-gzip"); ti == nil || err != nil ||
		ti.Columns != 80 {
		t.Error("alias-gzip failed")
	}
}

func TestTerminfoBadAlias(t *testing.T) {

	os.Setenv("TCELLDB", "testdata")
	if ti, err := LookupTerminfo("alias-none"); err == nil || ti != nil {
		t.Errorf("Bad alias should not have worked")
	}
}

func TestTerminfoCombined(t *testing.T) {

	os.Setenv("TCELLDB", "testdata/combined")

	var values = []struct {
		name  string
		lines int
	}{
		{"combined2", 102},
		{"alias-comb1", 101},
		{"combined3", 103},
		{"combined1", 101},
	}
	for _, v := range values {
		if ti, e := LookupTerminfo(v.name); e != nil || ti == nil ||
			ti.Lines != v.lines {
			t.Errorf("Combined terminal for %s wrong", v.name)
		}
	}
}

func BenchmarkSetFgBg(b *testing.B) {
	ti := testTerminfo

	for i := 0; i < b.N; i++ {
		ti.TParm(ti.SetFg, 100, 200)
		ti.TParm(ti.SetBg, 100, 200)
	}
}
