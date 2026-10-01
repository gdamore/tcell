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

// The EventUnknown events are sent when the terminal sends an escape
// sequence that tcell does not interpret itself, one type per kind of
// sequence. This lets applications write their own queries to the Tty and
// receive the replies.
//
// Sequence holds the whole sequence, starting with the escape character.
// Control strings (OSC, DCS, APC, PM, and SOS) always end with ST (ESC \),
// even if the terminal ended them with BEL.

// EventUnknownCSI is sent for a CSI sequence that tcell does not interpret.
type EventUnknownCSI struct {
	EventTime
	Sequence string
}

// EventUnknownOSC is sent for an OSC string that tcell does not interpret.
type EventUnknownOSC struct {
	EventTime
	Sequence string
}

// EventUnknownDCS is sent for a DCS string that tcell does not interpret.
type EventUnknownDCS struct {
	EventTime
	Sequence string
}

// EventUnknownAPC is sent for an APC string that tcell does not interpret.
type EventUnknownAPC struct {
	EventTime
	Sequence string
}

// EventUnknownPM is sent for a PM string that tcell does not interpret.
type EventUnknownPM struct {
	EventTime
	Sequence string
}

// EventUnknownSOS is sent for an SOS string that tcell does not interpret.
type EventUnknownSOS struct {
	EventTime
	Sequence string
}

// EventUnknownSS2 is sent for an SS2 sequence that tcell does not interpret.
type EventUnknownSS2 struct {
	EventTime
	Sequence string
}

// EventUnknownSS3 is sent for an SS3 sequence that tcell does not interpret.
type EventUnknownSS3 struct {
	EventTime
	Sequence string
}
