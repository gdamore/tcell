// Copyright 2019 The TCell Authors
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

// EventBundle wraps a slice of Events, allowing multiple events to
// be posted "atomically" to the event-handling goroutine.
type EventBundle struct {
	t      time.Time
	events []Event
}

// When returns the time when this event was created.
func (ev *EventBundle) When() time.Time {
	return ev.t
}

// Events returns the bundled events as a slice.
func (ev *EventBundle) Events() []Event {
	return ev.events
}

// NewEventBundle creates an EventBundle containing the provided
// events.
func NewEventBundle(events []Event) *EventBundle {
	return &EventBundle{t: time.Now(), events: events}
}
