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

//go:build !js && !wasm
// +build !js,!wasm

package mock

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
)

// TestMockStart just verifies that we can start the mock tty.
func TestMockStart(t *testing.T) {
	mt := &MockTty{}
	mt.Reset()
	screen, err := tcell.NewTerminfoScreenFromTty(mt)
	if err != nil {
		t.Fatalf("cannot get terminfo screen: %v", err)
	}
	if err = screen.Init(); err != nil {
		t.Fatalf("failed to init screen: %v", err)
	}

	w, h := screen.Size()
	if w != 80 || h != 24 {
		t.Errorf("bad window size %d x %d", w, h)
	}
	time.Sleep(time.Millisecond * 100)
	screen.Fini()
}
