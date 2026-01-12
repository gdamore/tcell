// Copyright 2025 The TCell Authors
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

package vt

import "testing"

func TestPmStrings(t *testing.T) {
	pm := PmAutoMargin
	if s := pm.Enable(); s != "\x1b[?7h" {
		t.Errorf("enable string wrong: %s", s)
	}
	if s := pm.Disable(); s != "\x1b[?7l" {
		t.Errorf("disable string wrong: %s", s)
	}
	if s := pm.Query(); s != "\x1b[?7$p" {
		t.Errorf("query string wrong: %s", s)
	}
	if s := pm.Reply(ModeNA); s != "\x1b[?7;0$y" {
		t.Errorf("reply(NA) string wrong: %s", s)
	}
	if s := pm.Reply(ModeOn); s != "\x1b[?7;1$y" {
		t.Errorf("reply(On) string wrong: %s", s)
	}
	if s := pm.Reply(ModeOff); s != "\x1b[?7;2$y" {
		t.Errorf("reply(Off) string wrong: %s", s)
	}
	if s := pm.Reply(ModeOnLocked); s != "\x1b[?7;3$y" {
		t.Errorf("reply(OnLocked) string wrong: %s", s)
	}
	if s := pm.Reply(ModeOffLocked); s != "\x1b[?7;4$y" {
		t.Errorf("reply(OffLocked) string wrong: %s", s)
	}
}
