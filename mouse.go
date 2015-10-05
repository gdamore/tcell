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
	"time"
)

// EventMouse is a mouse event.  It is sent on either mouse up or mouse down
// events.  It is also sent on mouse motion events - if the terminal supports
// it.  We make every effort to ensure that mouse release events are delivered.
// Hence, click drag can be identified by a motion event with the mouse down,
// without any intervening button release.
//
// Mouse wheel events, when reported, may appear on their own as individual
// impulses; that is, there will normally not be a release event delivered
// for mouse wheel movements.
//
// Most terminals cannot report the state of more than one button at a time --
// and many cannot report motion events.  (Windows consoles, modern XTerm, and
// modern emulators like iTerm2, are known to support this well, though.)
//
// Applications can inspect the time between events to figure out double clicks
// and such.
type EventMouse struct {
	t   time.Time
	btn ButtonMask
	mod ModMask
	x   int
	y   int
}

func (ev *EventMouse) When() time.Time {
	return ev.t
}

// ButtonMask returns the list of buttons that were pressed.
func (ev *EventMouse) Buttons() ButtonMask {
	return ev.btn
}

// Modifiers returns a list of keyboard modifiers that were pressed
// with the mouse button(s).
func (ev *EventMouse) Modifiers() ModMask {
	return ev.mod
}

// Position returns the mouse position in character cells.  The origin
// 0, 0 is at the upper left corner.
func (ev *EventMouse) Position() (int, int) {
	return ev.x, ev.y
}

// NewEventMouse is used to create a new mouse event.  Applications
// shouldn't need to use this; its mostly for screen implementors.
func NewEventMouse(x, y int, btn ButtonMask, mod ModMask) *EventMouse {
	return &EventMouse{t: time.Now(), x: x, y: y, btn: btn, mod: mod}
}

// BtnMask is a mask of mouse buttons.
type ButtonMask int16

const (
	// Button1 is usually the left mouse button.
	Button1 ButtonMask = 1 << iota
	// Button2 is usually the middle mouse button, for three button mice.
	Button2
	// Button3 is usually the right mouse button on 2 or 3 button mice.
	Button3
	Button4
	Button5
	Button6
	Button7
	Button8
	// WheelUp indicates the wheel being moved up, away from the user.
	WheelUp
	// WheelDown indicates the wheel being moved down, towards the user.
	WheelDown
	WheelLeft
	WheelRight
)
const ButtonNone ButtonMask = 0
