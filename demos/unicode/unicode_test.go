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

// TestUnicode just exercises the code in the unicode demo program.
// It does not validate that the content is accurate, that should be done
// manually by running the program with a real terminal.
func TestUnicode(t *testing.T) {

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

	// control L (forces sync)
	mt.KeyTap(vt.KeyLCtrl, vt.KeyL)

	attrs := 0
	var lastAttr vt.Attr
	for y := range vt.Row(24) {
		for x := range vt.Col(80) {
			if attr := mt.GetCell(vt.Coord{X: vt.Col(x), Y: vt.Row(y)}).S.Attr(); attr != lastAttr {
				attrs++
				lastAttr = attr
			}
		}
	}
	time.Sleep(time.Millisecond * 10)
	mt.KeyTap(vt.KeyLCtrl, vt.KeyQ)

	wg.Wait()
}
