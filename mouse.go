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
// events.  (Eventually we can also arrange for mouse motion events, but only
// with genuine xterm -- other emulators lack support for tracking this.)
// Most terminals cannot report the state of more than one button at a time --
// that is buttons are seen together.
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

func (ev *EventMouse) Buttons() ButtonMask {
	return ev.btn
}

func (ev *EventMouse) Modifiers() ModMask {
	return ev.mod
}

func (ev *EventMouse) Position() (int, int) {
	return ev.x, ev.y
}

func NewEventMouse(x, y int, btn ButtonMask, mod ModMask) *EventMouse {
	return &EventMouse{t: time.Now(), x: x, y: y, btn: btn, mod: mod}
}

// BtnMask is a mask of mouse buttons.
type ButtonMask int16

const (
	Button1 ButtonMask = 1 << iota
	Button2
	Button3
	Button4
	Button5
	Button6
	Button7
	Button8
)
const ButtonNone ButtonMask = 0
