// Copyright 2026 The TCell Authors
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

package main

import (
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/vt"
)

func TestBeep(t *testing.T) {

	mt := vt.NewMockTerm()
	scr, err := tcell.NewTerminfoScreenFromTty(mt)
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

	mt.Drain()

	// simulate two presses of B
	mt.KeyEvent(vt.KeyEvent{Code: 'b', Down: true})
	mt.KeyEvent(vt.KeyEvent{Code: 'b', Down: false})
	time.Sleep(time.Millisecond * 100)
	mt.KeyEvent(vt.KeyEvent{Code: 'b', Down: true})
	mt.KeyEvent(vt.KeyEvent{Code: 'b', Down: false})

	// control L (forces sync)
	mt.KeyEvent(vt.KeyEvent{Code: 'l', Down: true, Mod: vt.ModCtrl})
	mt.KeyEvent(vt.KeyEvent{Code: 'l', Down: false, Mod: vt.ModCtrl})

	// sleep at least 3 seconds to get the time driven beep
	time.Sleep(time.Millisecond * 3500)

	mt.KeyEvent(vt.KeyEvent{Code: 'q', Down: true, Mod: vt.ModCtrl})
	mt.KeyEvent(vt.KeyEvent{Code: 'q', Down: false, Mod: vt.ModCtrl})

	wg.Wait()

	if cnt := mt.Bells(); cnt != 3 {
		t.Errorf("incorrect bell count %d != 2", cnt)
	}
}
