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

import "time"

// EventPaste represents a bracketed paste event.
type EventPaste struct {
	t    time.Time
	text string
}

// When returns the time when this Event was created, which should closely
// match the time when the paste was made.
func (e *EventPaste) When() time.Time {
	return e.t
}

// Text returns the text that was pasted
func (e *EventPaste) Text() string {
	return e.text
}

// NewEventPaste creates a new paste event from the given text
func NewEventPaste(text string) *EventPaste {
	return &EventPaste{
		t:    time.Now(),
		text: text,
	}
}
