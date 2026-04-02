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

package vt

import (
	"slices"
	"testing"

	"github.com/gdamore/tcell/v3/color"
)

type noMouseBackend struct {
	mb *mockBackend
}

func (n noMouseBackend) GetPrivateMode(pm PrivateMode) ModeStatus { return n.mb.GetPrivateMode(pm) }
func (n noMouseBackend) SetPrivateMode(pm PrivateMode, status ModeStatus) error {
	return n.mb.SetPrivateMode(pm, status)
}
func (n noMouseBackend) GetSize() Coord               { return n.mb.GetSize() }
func (n noMouseBackend) Colors() int                  { return n.mb.Colors() }
func (n noMouseBackend) Put(pos Coord, cell Cell)     { n.mb.Put(pos, cell) }
func (n noMouseBackend) GetPosition() Coord           { return n.mb.GetPosition() }
func (n noMouseBackend) SetPosition(pos Coord)        { n.mb.SetPosition(pos) }
func (n noMouseBackend) Reset()                       { n.mb.Reset() }
func (n noMouseBackend) RaiseResize()                 { n.mb.RaiseResize() }
func (n noMouseBackend) Buffering(enabled bool)       { n.mb.Buffering(enabled) }
func (n noMouseBackend) SetCursor(cs CursorStyle)     { n.mb.SetCursor(cs) }

func TestEmulatorModeHelpers(t *testing.T) {
	mb := NewMockBackend().(*mockBackend)
	em := NewEmulator(mb).(*emulator)

	ansiKeys := em.ansiModeKeys()
	if !slices.Contains(ansiKeys, AmNewLineMode) {
		t.Fatalf("expected AmNewLineMode in ansi mode keys, got %v", ansiKeys)
	}

	privateKeys := em.privateModeKeys()
	if !slices.Contains(privateKeys, PmAutoMargin) {
		t.Fatalf("expected PmAutoMargin in private mode keys, got %v", privateKeys)
	}

	em.setAnsiMode(AmNewLineMode, ModeOn)
	if got := em.getAnsiMode(AmNewLineMode); got != ModeOn {
		t.Fatalf("setAnsiMode did not update changeable mode: got %v", got)
	}

	em.setAnsiMode(AmNewLineMode, ModeNA)
	if got := em.getAnsiMode(AmNewLineMode); got != ModeOn {
		t.Fatalf("setAnsiMode changed mode for non-changeable status: got %v", got)
	}

	em.ansiModes[AmInsertReplace] = ModeOnLocked
	em.setAnsiMode(AmInsertReplace, ModeOff)
	if got := em.getAnsiMode(AmInsertReplace); got != ModeOnLocked {
		t.Fatalf("setAnsiMode changed locked mode: got %v", got)
	}

	em.setAnsiMode(AnsiMode(9999), ModeOn)
	if _, ok := em.ansiModes[AnsiMode(9999)]; ok {
		t.Fatalf("setAnsiMode created an unknown mode entry")
	}
}

func TestEmulatorUpdateMouseReporting(t *testing.T) {
	mb := NewMockBackend().(*mockBackend)
	em := NewEmulator(mb).(*emulator)

	em.localModes[PmMouseButton] = ModeOn
	em.updateMouseReporting()
	if em.mouseReports != MouseButtons {
		t.Fatalf("expected mouse buttons reporting, got %v", em.mouseReports)
	}

	em.localModes[PmMouseDrag] = ModeOn
	em.updateMouseReporting()
	if em.mouseReports != MouseDrag {
		t.Fatalf("expected mouse drag reporting, got %v", em.mouseReports)
	}

	em.localModes[PmMouseMotion] = ModeOn
	em.updateMouseReporting()
	if em.mouseReports != MouseMotion {
		t.Fatalf("expected mouse motion reporting, got %v", em.mouseReports)
	}

	em.localModes[PmMouseButton] = ModeOff
	em.localModes[PmMouseDrag] = ModeOff
	em.localModes[PmMouseMotion] = ModeOff
	em.localModes[PmMouseX10] = ModeOn
	em.updateMouseReporting()
	if em.mouseReports != MouseButtons {
		t.Fatalf("expected mouse X10 reporting to map to buttons, got %v", em.mouseReports)
	}

	em.localModes[PmMouseX10] = ModeOff
	em.updateMouseReporting()
	if em.mouseReports != MouseDisabled {
		t.Fatalf("expected mouse reporting disabled, got %v", em.mouseReports)
	}

	noMouse := NewEmulator(noMouseBackend{mb: NewMockBackend().(*mockBackend)}).(*emulator)
	noMouse.updateMouseReporting()
}

func TestMockBackendHelpers(t *testing.T) {
	mb := NewMockBackend().(*mockBackend)

	style := BaseStyle.WithFg(color.Red).WithBg(color.Blue)
	mb.SetStyle(style)
	if mb.style != style {
		t.Fatalf("SetStyle did not update backend style")
	}

	mb.SetMouse(MouseMotion)
	mb.Buffering(true)
	mb.Buffering(false)

	MockOptNoBlit{}.SetMockOpt(mb)
}
