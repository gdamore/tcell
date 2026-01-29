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

package tcell

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/vt"
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

func getMouseEvent(t *testing.T, s Screen) (*EventMouse, bool) {
	t.Helper()
	for {
		select {
		case ev := <-s.EventQ():
			if me, ok := ev.(*EventMouse); ok {
				return me, true
			}
			t.Logf("Got different event: %T", ev)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for mouse event")
			return nil, false
		}
	}
}

func getFocusEvent(t *testing.T, s Screen) (*EventFocus, bool) {
	t.Helper()
	for {
		select {
		case ev := <-s.EventQ():
			if fe, ok := ev.(*EventFocus); ok {
				return fe, true
			}
			t.Logf("Got different event: %T", ev)
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for focus event")
			return nil, false
		}
	}
}

func TestMouseEvents(t *testing.T) {
	term, s := NewMockScreen(t)
	defer s.Fini()

	s.EnableMouse(MouseMotionEvents)
	time.Sleep(time.Millisecond * 10)

	// simple mouse click
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button1,
		Down:     true,
		Motion:   false,
		Mod:      vt.ModNone,
	})
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button1,
		Down:     false,
		Motion:   false,
		Mod:      vt.ModNone,
	})

	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != Button1 {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	}

	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != ButtonNone {
		t.Errorf("still got buttons %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	}

	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.WheelUp,
		Down:     true,
		Mod:      vt.ModLAlt,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != WheelUp {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModAlt {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModAlt)
	}

	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.WheelDown,
		Down:     true,
		Mod:      vt.ModRMeta | vt.ModRCtrl,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != WheelDown {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModCtrl|ModAlt { // NB: ModAlt used for Meta in mouse events
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModAlt|ModCtrl)
	}

	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.WheelLeft,
		Down:     true,
		Mod:      vt.ModRShift,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != WheelLeft {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModShift)
	}

	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.WheelRight,
		Down:     true,
		Mod:      vt.ModLShift,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != WheelRight {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModShift)
	}

	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button2,
		Down:     true,
		Mod:      vt.ModLShift,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != Button2 {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModShift)
	}

	// this will be a chord of buttons 2 and 3
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button3,
		Down:     true,
		Mod:      vt.ModLShift,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != Button3|Button2 {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModShift)
	}

	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button3,
		Down:     false,
		Mod:      vt.ModLShift,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != Button2 {
		t.Errorf("wrong buttons: %x", me.Buttons())
	}
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button2,
		Down:     false,
		Mod:      vt.ModLShift,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != ButtonNone {
		t.Errorf("wrong buttons: %x", me.Buttons())
	}

	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.NoButton,
		Down:     true,
		Mod:      vt.ModLShift,
		Motion:   true,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != ButtonNone {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModShift)
	}

	// this simulates a certain kind of broken report
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button1,
		Down:     false,
		Mod:      vt.ModLShift,
		Motion:   true,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != ButtonNone {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModShift)
	}

	// send malformed mouse events (fuzz testing)
	term.SendRaw([]byte("\x1b[<3M"))
	term.SendRaw([]byte("\x1b[<3;2m"))

	select {
	case ev := <-s.EventQ():
		t.Fatalf("Got unexpected event: %T", ev)
	case <-time.After(time.Millisecond * 50):
	}

	// Now try the upper buttons (uncommon)
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button4,
		Down:     true,
		Mod:      vt.ModShift,
		Motion:   false,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != Button4 {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModShift)
	}
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button4,
		Down:     false,
		Mod:      vt.ModShift,
		Motion:   false,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != ButtonNone {
		t.Errorf("wrong buttons: %x", me.Buttons())
	}

	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button5,
		Down:     true,
		Mod:      vt.ModShift,
		Motion:   false,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != Button5 {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModShift)
	}

	// chord with 5
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button6,
		Down:     true,
		Mod:      vt.ModShift,
		Motion:   false,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != Button6|Button5 {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModShift)
	}

	// and now adding the final chord with button 7
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button7,
		Down:     true,
		Mod:      vt.ModShift,
		Motion:   false,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != Button5|Button6|Button7 {
		t.Errorf("wrong buttons: %x", me.Buttons())
	} else if x, y := me.Position(); x != 2 || y != 3 {
		t.Errorf("wrong position %d,%d != 2,3", x, y)
	} else if me.Modifiers() != ModShift {
		t.Errorf("wrong modifiers %x != %x", me.Modifiers(), ModShift)
	}

	// and release them out of order
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button6,
		Down:     false,
		Mod:      vt.ModRShift,
		Motion:   false,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != Button5|Button7 {
		t.Errorf("wrong buttons: %x", me.Buttons())
	}
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button5,
		Down:     false,
		Mod:      vt.ModRShift,
		Motion:   false,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != Button7 {
		t.Errorf("wrong buttons: %x", me.Buttons())
	}
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button7,
		Down:     false,
		Mod:      vt.ModRShift,
		Motion:   false,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != ButtonNone {
		t.Errorf("wrong buttons: %x", me.Buttons())
	}

	// Spurious release event
	term.MouseEvent(vt.MouseEvent{
		Position: vt.Coord{X: 2, Y: 3},
		Button:   vt.Button7,
		Down:     false,
		Mod:      vt.ModRShift,
		Motion:   false,
	})
	if me, ok := getMouseEvent(t, s); !ok {
		return
	} else if me.Buttons() != ButtonNone {
		t.Errorf("wrong buttons: %x", me.Buttons())
	}
}

func TestFocusEvents(t *testing.T) {
	term, s := NewMockScreen(t)
	defer s.Fini()

	s.EnableFocus()
	time.Sleep(time.Millisecond * 10)

	term.FocusEvent(true)
	if ev, ok := getFocusEvent(t, s); !ok {
		return
	} else if !ev.Focused {
		t.Error("focused was not set")
	}

	term.FocusEvent(false)
	if ev, ok := getFocusEvent(t, s); !ok {
		return
	} else if ev.Focused {
		t.Error("focused was set")
	}

	term.FocusEvent(true)
	if ev, ok := getFocusEvent(t, s); !ok {
		return
	} else if !ev.Focused {
		t.Error("focused was not set")
	}
}
