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

package mock

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/vt"
)

// mockTerm implements MockTerm.
type mockTerm struct {
	mb MockBackend
	em vt.Emulator
}

// Stop the terminal.
func (mt *mockTerm) Stop() error {
	return mt.em.Stop()
}

// Start the terminal.
func (mt *mockTerm) Start() error {
	return mt.em.Start()
}

// Drain all output from the terminal, ensuring
// any queued commands are processed.
func (mt *mockTerm) Drain() error {
	return mt.em.Drain()
}

// Read data from the terminal. This is called by a terminal
// application (e.g. via tcell Tty.)  Read data will include
// key strokes, mouse events, and responses to terminal queries.
func (mt *mockTerm) Read(data []byte) (int, error) {
	return mt.em.Read(data)
}

// Write data to the terminal, typically either commands or data
// that should be displayed on the virtual screen.
func (mt *mockTerm) Write(b []byte) (n int, err error) {
	return mt.em.Write(b)
}

// WindowSize obtains the dimensions of the window.
func (mt *mockTerm) WindowSize() (tcell.WindowSize, error) {
	sz := mt.mb.GetSize()
	// No pixel sizes for now
	return tcell.WindowSize{Width: int(sz.X), Height: int(sz.Y)}, nil
}

// NotifyResize registers a channel to be signaled when a resize has occurred.
// In real terminal emulators this would be posted (non-blocking) by a signal handler.
func (mt *mockTerm) NotifyResize(resizeq chan<- bool) {
	if rs, ok := mt.mb.(vt.Resizer); ok {
		rs.NotifyResize(resizeq)
	}
}

// Close closes the terminal, after which it should no longer be used. Stop is implied.
func (mt *mockTerm) Close() error {
	return mt.Stop()
}

// Pos returns the cursor position.
func (mt *mockTerm) Pos() vt.Coord {
	return mt.mb.GetPosition()
}

// GetCell returns the contents of the cell at the given coordinates, or a zero value
// if the coordinates are out of range.
func (mt *mockTerm) GetCell(pos vt.Coord) Cell {
	return mt.mb.GetCell(pos)
}

// Bells counts the number of times the bell has rung.
func (mt *mockTerm) Bells() int {
	return mt.mb.Bells()
}

// MockTerm is a mock terminal (emulator).  It can be used to
// test the emulator itself, or to test applications (or tcell) that
// uses the terminal.  It also implements the Tty interface used
// by tcell itself.
type MockTerm interface {
	tcell.Tty

	// Pos reports the current cursor position.
	Pos() vt.Coord

	// GetCell returns the current cell.
	GetCell(vt.Coord) Cell

	// Bells returns the number of times the bell has been rung.
	Bells() int
}

// NewMockTerm gives a mock terminal emulator.
func NewMockTerm(opts ...MockOpt) MockTerm {
	mt := &mockTerm{}
	mt.mb = NewMockBackend(opts...)
	mt.em = vt.NewEmulator(mt.mb)
	mt.em.SetId("TcellMock", "1.0")
	return mt
}
