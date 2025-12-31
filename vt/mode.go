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

// Package vt provides common definitions for VT derived terminals and applications.
// This includes the venerable VT100, XTerm, and newer emulators such as Kitty and
// the Windows Terminal.
package vt

import "fmt"

// PrivateMode describes a DEC Private Mode.
type PrivateMode int

const (
	PmAppCursor      PrivateMode = 1 // application cursor keys
	PmAutoMargin     PrivateMode = 7
	PmAltScreen      PrivateMode = 1049 // 47 and 1047 are alternates, but we use 1049
	PmBlinkCursor    PrivateMode = 12
	PmShowCursor     PrivateMode = 25
	PmMouseButton    PrivateMode = 1000
	PmMouseDrag      PrivateMode = 1002
	PmMouseMotion    PrivateMode = 1003
	PmFocusReports   PrivateMode = 1004
	PmMouseSgr       PrivateMode = 1006
	PmMouseSgrPixel  PrivateMode = 1016
	PmBracketedPaste PrivateMode = 2004
	PmSyncOutput     PrivateMode = 2026
	PmResizeReports  PrivateMode = 2048 // send in-band resize reports
)

// Enable returns the string used to enable this private mode.
func (pm PrivateMode) Enable() string {
	return fmt.Sprintf("\x1b[?%dh", pm)
}

// Disable returns the string used to disable this private mode.
func (pm PrivateMode) Disable() string {
	return fmt.Sprintf("\x1b[?%dl", pm)
}

// Query returns the string used to query the state of this private mode.
func (pm PrivateMode) Query() string {
	return fmt.Sprintf("\x1b[?%d$p", pm)
}

// Reply returns a string representing a query reply for the given mode and status.
func (pm PrivateMode) Reply(status ModeStatus) string {
	return fmt.Sprintf("\x1b[?%d;%d$y", pm, status)
}

// ModeStatus represents the status of the mode.
type ModeStatus int

const (
	ModeNA        ModeStatus = 0 // Mode is not supported (or unknown)
	ModeOn        ModeStatus = 1 // Mode is on (e.g. via CSI-h)
	ModeOff       ModeStatus = 2 // Mode is off (e.g. via CSI-l)
	ModeOnLocked  ModeStatus = 3 // Mode is hardwired on
	ModeOffLocked ModeStatus = 4 // Mode is hardwired off
)

// Changeable indicates that the mode may be changed.
func (ms ModeStatus) Changeable() bool {
	return ms == ModeOn || ms == ModeOff
}
