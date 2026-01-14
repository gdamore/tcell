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

	mt.KeyEvent(vt.KeyEvent{Code: '1', Down: true})
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.BlinkingBlock)

	mt.KeyEvent(vt.KeyEvent{Code: '2', Down: true})
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.SteadyBlock)

	mt.KeyEvent(vt.KeyEvent{Code: '3', Down: true})
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.BlinkingUnderline)

	mt.KeyEvent(vt.KeyEvent{Code: '4', Base: '4', Down: true})
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.SteadyUnderline)

	mt.KeyEvent(vt.KeyEvent{Code: '5', Base: '5', Down: true})
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.BlinkingBar)

	mt.KeyEvent(vt.KeyEvent{Code: '6', Base: '6', Down: true})
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.SteadyBar)

	mt.KeyEvent(vt.KeyEvent{Code: '0', Base: '0', Down: true})
	time.Sleep(time.Millisecond * 20)
	verifyCursor(t, mt, vt.BlinkingBlock)

	mt.KeyEvent(vt.KeyEvent{Code: 'L', Base: 'L', Mod: vt.ModCtrl, Down: true})
	mt.Drain()
	mt.KeyEvent(vt.KeyEvent{Code: 'Q', Base: 'Q', Mod: vt.ModCtrl, Down: true})
	mt.Backend().GetSize()
	wg.Wait()
}
