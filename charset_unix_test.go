//go:build unix
// +build unix

package tcell

import "testing"

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

func TestGetCharset(t *testing.T) {
	cases := []struct {
		name    string
		locale  string
		charset string
	}{
		{"POSIX", "POSIX", "US-ASCII"},
		{"C", "C", "US-ASCII"},
		{"C Unicode", "C.UTF-8", "UTF-8"},
		{"Old UK euro", "UK.ISO-8859-1@euro", "ISO-8859-1"},
		{"Unset", "", "UTF-8"},
	}

	for _, tc := range cases {
		t.Setenv("LANG", "")
		t.Setenv("LC_ALL", "")
		t.Setenv("LC_CTYPE", "")
		for _, env := range []string{"LANG", "LC_ALL", "LC_CTYPE"} {
			t.Run(tc.name+" "+env, func(t *testing.T) {
				t.Setenv(env, tc.locale)
				if actual := getCharset(); actual != tc.charset {
					t.Errorf("charset %q did not match %q", actual, tc.charset)
				}
			})
		}
	}
}
