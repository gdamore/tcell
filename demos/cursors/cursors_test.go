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

package main

import (
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/vt"
)

func verifyCursor(t *testing.T, term vt.MockTerm, expected vt.CursorStyle) {
	t.Helper()
	if actual := term.Backend().GetCursor(); actual != expected {
		t.Errorf("wrong cursor style %x != %x", actual, expected)
	}
}

// TestDemo tests the cursors demo.
func TestDemo(t *testing.T) {

	mt := vt.NewMockTerm()
	scr, err := tcell.NewTerminfoScreenFromTty(mt, tcell.OptColors(8), tcell.OptTerm("ansi"))
	if err != nil {
		t.Fatalf("failed to create screen: %v", err)
	}
	tcell.ShimScreen(scr)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		main()
	}()

	// give enough time for the screen to draw
	// this is needed because main is running asynchronously.
	time.Sleep(time.Millisecond * 50)

	mt.KeyTap(vt.Key1)
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.BlinkingBlock)

	mt.KeyTap(vt.Key2)
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.SteadyBlock)

	mt.KeyTap(vt.Key3)
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.BlinkingUnderline)

	mt.KeyTap(vt.Key4)
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.SteadyUnderline)

	mt.KeyTap(vt.Key5)
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.BlinkingBar)

	mt.KeyTap(vt.Key6)
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.SteadyBar)

	mt.KeyTap(vt.Key0)
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.BlinkingBlock)

	mt.KeyTap(vt.KeyRCtrl, vt.KeyL)
	mt.Drain()
	mt.KeyTap(vt.KeyRCtrl, vt.KeyQ)
	mt.Backend().GetSize()
	wg.Wait()
}
