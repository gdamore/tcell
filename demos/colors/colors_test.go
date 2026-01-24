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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/vt"
)

// TestColors just exercises the code in the colors demo program.
// It does not validate that the content is accurate, that should be done
func TestColors(t *testing.T) {

	interval = time.Microsecond * 10

	for _, colors := range []int{0, 8, 16, 88, 256, 1 << 24} {
		t.Run(fmt.Sprintf("%d_colors)", colors), func(t *testing.T) {
			mt := vt.NewMockTerm(vt.MockOptColors(colors))
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

			time.Sleep(time.Millisecond * 25)
			mt.KeyTap(vt.KeyLCtrl, vt.KeyRShift, vt.KeyL)
			mt.SetSize(vt.Coord{X: 10, Y: 10})
			mt.Drain()
			time.Sleep(time.Millisecond * 25)
			mt.KeyTap(vt.KeyLCtrl, vt.KeyQ)
			mt.Drain()
			wg.Wait()
		})
	}
}
