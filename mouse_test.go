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

package tcell

import (
	"testing"
	"time"
)

func TestMouseEventFields(t *testing.T) {
	now := time.Now()
	ev := NewEventMouse(4, 5, ButtonMiddle, ModShift)
	if ev.when.Before(now) {
		t.Errorf("time went backwards? %s", ev.when)
	}
	if time.Now().Before(ev.when) {
		t.Errorf("time also went backwards? %s", ev.when)
	}
	if ev.Buttons() != Button3 {
		t.Errorf("wrong buttons %v", ev.Buttons())
	}
	if ev.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %v", ev.Modifiers())
	}
	if x, y := ev.Position(); x != 4 || y != 5 {
		t.Errorf("wrong position %d %d", x, y)
	}
}
