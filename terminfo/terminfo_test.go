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

package terminfo

import (
	"testing"
)

// This terminfo entry is a stripped down version from
// xterm-256color, but I've added some of my own entries.
var testTerminfo = &Terminfo{
	Name:      "simulation_test",
	Columns:   80,
	Lines:     24,
	XTermLike: true,
}

func TestTerminfoExpansion(t *testing.T) {
	ti := testTerminfo

	// This tests some conditionals
	if ti.TParm("A[%p1%2.2X]B", 47) != "A[2F]B" {
		t.Error("TParm conditionals failed")
	}

	// Color tests.

	type testCase struct {
		expect string
		format string
		params []any
	}

	cases := []testCase{
		{expect: "0a", format: "%p1%02x", params: []any{10}},
		{expect: "0A", format: "%p1%02X", params: []any{10}},
		{expect: "A", format: "%p1%c", params: []any{65}},
		{expect: "A", format: "%'A'%c", params: []any{}},
		{expect: "65", format: "%'A'%d", params: []any{}},
		{expect: "7", format: "%i%p1%p2%+%d", params: []any{2, 3}},
		{expect: "abc", format: "%p1%s", params: []any{"abc"}},
		{expect: "1%d", format: "1%%d", params: []any{}},
		{expect: "abc", format: "%p1%s%", params: []any{"abc"}}, // unterminated %
		{expect: "  abc", format: "%p1%5s", params: []any{"abc"}},
		{expect: "abc  ", format: "%p1%:-5s", params: []any{"abc"}},
		{expect: "15", format: "%{3}%p1%*%d", params: []any{5}},
		{expect: " A", format: "%p1%2c", params: []any{65}},
		{expect: "4", format: "%p1%l%d", params: []any{"four"}},
		{expect: "0", format: "%pA%d", params: []any{}}, // missing/invalid parameter
		{expect: "5", format: "%p1%p2%/%d", params: []any{15, 3}},
		{expect: "0", format: "%p1%p2%/%d", params: []any{3, 15}},
		{expect: "0", format: "%p1%p2%/%d", params: []any{3, 0}},
		{expect: "3", format: "%p1%p2%m%d", params: []any{15, 4}},
		{expect: "0", format: "%p1%p2%m%d", params: []any{3, 0}},
		{expect: "2", format: "%p1%Pa%{4}%{3}%ga%d", params: []any{2}},
		{expect: "2", format: "%p1%PA%{4}%{3}%gA%d", params: []any{2}},
		{expect: "0", format: "%p1%PA%{4}%{3}%ga%d", params: []any{2}},
		{expect: "0", format: "%p1%Pz%{4}%{3}%gZ%d", params: []any{2}},
		{expect: "0", format: "%d", params: []any{}}, // underflow
		{expect: "", format: "%s", params: []any{}},  // underflow
		{expect: "1", format: "%p1%p2%=%d", params: []any{3, 3}},
		{expect: "0", format: "%p1%p2%=%d", params: []any{3, 4}},
		{expect: "1", format: "%p1%p2%=%!%d", params: []any{3, 4}},
		{expect: "1", format: "%p1%p2%>%d", params: []any{4, 3}},
		{expect: "3", format: "%p1%p2%|%d", params: []any{1, 2}},
		{expect: "2", format: "%p1%p2%&%d", params: []any{2, 3}},
		{expect: "1", format: "%p1%p2%^%d", params: []any{2, 3}},
		{expect: "f", format: "%p1%~%{255}%&%x", params: []any{0xf0}},
		{expect: "%Z", format: "%Z", params: []any{2, 3}}, // unknown sequence
	}

	for i := range cases {
		if res := ti.TParm(cases[i].format, cases[i].params...); res != cases[i].expect {
			t.Errorf("Format case %d failed: Format %q got %q", i, cases[i].format, res)
		}
	}
	t.Logf("Tested %d cases", len(cases))
}

func TestStringParameter(t *testing.T) {
	ti := testTerminfo
	enterUrl := "\x1b]8;%p2%s;%p1%s\x1b\\"

	s := ti.TParm(enterUrl, "https://example.org/test")
	if s != "\x1b]8;;https://example.org/test\x1b\\" {
		t.Errorf("Result string failed: %s", s)
	}
	s = ti.TParm(enterUrl, "https://example.org/test", "id=1234")
	if s != "\x1b]8;id=1234;https://example.org/test\x1b\\" {
		t.Errorf("Result string failed: %s", s)
	}
}
