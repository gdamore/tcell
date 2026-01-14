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

// TestDemo tests the clipboard demo.
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

	mt.KeyEvent(vt.KeyEvent{Code: '1', Base: '1', Down: true})
	mt.Drain()
	time.Sleep(time.Millisecond * 20)

	mt.KeyEvent(vt.KeyEvent{Code: '2', Base: '2', Down: true})
	mt.Drain()
	time.Sleep(time.Millisecond * 20)

	expect := "Enjoy your new clipboard content!"
	if result := string(mt.Backend().GetClipboard()); result != expect {
		t.Errorf("clipboard content did not match! %q != %q", result, expect)
	}

	// a long string
	mt.Backend().SetClipboard([]byte("The quick brown fox jumps over the lazy dog."))
	mt.KeyEvent(vt.KeyEvent{Code: '2', Base: '2', Down: true})
	mt.Drain()
	time.Sleep(time.Millisecond * 20)

	// stick some invalid utf-8
	mt.Backend().SetClipboard([]byte{0xff})
	mt.KeyEvent(vt.KeyEvent{Code: '2', Base: '2', Down: true})
	mt.Drain()
	time.Sleep(time.Millisecond * 20)

	// now nil
	mt.Backend().SetClipboard(nil)
	mt.KeyEvent(vt.KeyEvent{Code: '2', Base: '2', Down: true})
	mt.Drain()
	time.Sleep(time.Millisecond * 20)

	mt.KeyEvent(vt.KeyEvent{Code: 'Q', Base: 'Q', Mod: vt.ModCtrl, Down: true})
	mt.Backend().GetSize()
	wg.Wait()

}
