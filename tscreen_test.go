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

package tcell

import (
	"runtime"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/mock"
)

// This just offers some very basic tests that do not require a full mock.

// drainInput just does a very primitive sleep to allow input to drain.
// We need this because otherwise the application will close too soon before
// consuming characters from input, including sequences that are returned in
// response to queries.
func drainInput() {
	time.Sleep(time.Millisecond * 30)
}

// TestInitScreen just tries to initialze the default screen.
// It requires a working tty.
func TestInitScreen(t *testing.T) {
	s, err := NewTerminfoScreen()
	if err != nil {
		t.Skip("failed to get screen", err)
	}
	if err := s.Init(); err != nil {
		t.Skip("failed to initialize screen", err)
	}
	defer s.Fini()

	if s.CharacterSet() != "UTF-8" {
		t.Fatalf("Character Set (%v) not UTF-8", s.CharacterSet())
	}

	drainInput()
}

func TestNoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	s, err := NewTerminfoScreen()
	if err != nil {
		t.Skip("failed to get screen")
	}
	if err := s.Init(); err != nil {
		t.Skip("failed to initialize screen", err)
	}
	defer s.Fini()

	if s.Colors() != 0 {
		t.Errorf("screen should not have color but had %d", s.Colors())
	}

	drainInput()
}

func NewMockTerm(t *testing.T, opts ...mock.MockOpt) (mock.MockTerm, Screen) {
	t.Helper()
	if runtime.GOOS == "js" {
		t.Skip("not supported on webasm")
		return nil, nil
	}
	mt := mock.NewMockTerm(opts...)
	scr, err := NewTerminfoScreenFromTty(mt)
	if err != nil {
		t.Fatalf("failed to get terminal: %v", err)
		return nil, nil
	}
	if err = scr.Init(); err != nil {
		t.Fatalf("failed to initialize screen: %v", err)
		return nil, nil
	}
	return mt, scr
}
