// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"fmt"
	"testing"

	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/vt"
)

type MockTerm = vt.MockTerm
type Row = vt.Row
type Col = vt.Col
type Coord = vt.Coord
type Attr = vt.Attr

// WriteF writes the string, and ensures it is fully flushed
// before returning.
func WriteF(t *testing.T, term MockTerm, str string, args ...any) {
	t.Helper()
	b := fmt.Appendf(nil, str, args...)
	for len(b) != 0 {
		if n, err := term.Write(b); err != nil {
			t.Fatalf("Failed to write: %v", err)
		} else {
			b = b[n:]
		}
	}
	if err := term.Drain(); err != nil {
		t.Fatalf("Failed to flush: %v", err)
	}
}

// ReadF reads content from the term and returns it as a string.
func ReadF(t *testing.T, term MockTerm) string {
	t.Helper()
	buf := make([]byte, 128)
	n, err := term.Read(buf)
	if err != nil {
		t.Errorf("failed read: %v", err)
		return ""
	}
	return string(buf[:n])
}

// VerifyF validates the condition, printing the message on failure.
func VerifyF(t *testing.T, b bool, fmt string, args ...any) {
	t.Helper()
	if !b {
		t.Errorf("validation failure: "+fmt, args...)
	}
}

// AssertF validates the condition, and aborts the test if it fails.
func AssertF(t *testing.T, b bool, fmt string, args ...any) {
	t.Helper()
	if !b {
		t.Fatalf("validation failure: "+fmt, args...)
	}
}

// MustClose closes or calls Fatalf.
func MustClose(t *testing.T, term MockTerm) {
	t.Helper()
	err := term.Close()
	AssertF(t, err == nil, "close failed: %v", err)
}

// MustStart starts the terminal or calls Fatalf.
func MustStart(t *testing.T, term MockTerm) {
	t.Helper()
	err := term.Start()
	AssertF(t, err == nil, "start failed: %v", err)
}

// CheckPos is verifies the current cursor position of terminal.
func CheckPos(t *testing.T, term MockTerm, x Col, y Row) {
	t.Helper()
	VerifyF(t, term.Pos().X == x && term.Pos().Y == y,
		"bad position %d,%d (expected %d,%d)", term.Pos().X, term.Pos().Y, x, y)
}

// CheckContent verifies the content at a given cell of the terminal.
func CheckContent(t *testing.T, term MockTerm, x Col, y Row, s string) {
	t.Helper()
	if actual := string(term.GetCell(Coord{X: x, Y: y}).C); actual != s {
		t.Errorf("bad content %d,%d (expected %q got %q)", x, y, s, actual)
	}
}

// CheckAttrs verifies the attributes of a given cell.
func CheckAttrs(t *testing.T, term MockTerm, x Col, y Row, a Attr) {
	t.Helper()
	if actual := term.GetCell(Coord{X: x, Y: y}).S.Attr(); actual != a {
		t.Errorf("bad attr %d,%d (expected %x got %x)", x, y, a, actual)
	}
}

// CheckColors verifies the colors of a given cell.
func CheckColors(t *testing.T, term MockTerm, x Col, y Row, fg color.Color, bg color.Color) {
	t.Helper()
	if actual := term.GetCell(Coord{X: x, Y: y}).S.Fg(); actual != fg {
		t.Errorf("bad foreground %d,%d (expected %s got %s)", x, y, fg.String(), actual.String())
	}
	if actual := term.GetCell(Coord{X: x, Y: y}).S.Bg(); actual != bg {
		t.Errorf("bad background %d,%d (expected %s got %s)", x, y, bg.String(), actual.String())
	}
}

// CheckRead verifies that a read matches.
func CheckRead(t *testing.T, term MockTerm, want string) {
	t.Helper()
	result := ReadF(t, term)
	VerifyF(t, want == result, "wrong read %q != %q", result, want)
}
