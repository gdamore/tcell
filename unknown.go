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

// EventUnknown is sent when the terminal sends an escape sequence that
// tcell does not interpret itself. This lets applications write their own
// queries to the Tty and receive the replies.
type EventUnknown struct {
	EventTime

	// The sequence, starting with the escape character. Control strings
	// (such as OSC, DCS, and APC) always end with ST (ESC \), even if the
	// terminal ended them with BEL.
	Sequence string
}

// NewEventUnknown creates an EventUnknown for the given sequence.
func NewEventUnknown(seq string) *EventUnknown {
	ev := &EventUnknown{Sequence: seq}
	ev.SetEventNow()
	return ev
}
