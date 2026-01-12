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
	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/vt"
)

func TestHello(t *testing.T) {
	// ensure we only use 8 color ANSI for now

	mt := vt.NewMockTerm(vt.MockOptColors(8))
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

	mt.Drain()

	expect := []struct {
		X    vt.Col
		Y    vt.Row
		C    string
		Fg   color.Color
		Bg   color.Color
		Attr vt.Attr
	}{
		{X: 33, Y: 11, C: "H", Fg: color.Teal, Bg: color.Silver},
		{X: 34, Y: 11, C: "e", Fg: color.Teal, Bg: color.Silver},
		{X: 35, Y: 11, C: "l", Fg: color.Teal, Bg: color.Silver},
		{X: 36, Y: 11, C: "l", Fg: color.Teal, Bg: color.Silver},
		{X: 37, Y: 11, C: "o", Fg: color.Teal, Bg: color.Silver},
		{X: 38, Y: 11, C: ",", Fg: color.Teal, Bg: color.Silver},
		{X: 39, Y: 11, C: " ", Fg: color.Teal, Bg: color.Silver},
		{X: 40, Y: 11, C: "W", Fg: color.Teal, Bg: color.Silver},
		{X: 41, Y: 11, C: "o", Fg: color.Teal, Bg: color.Silver},
		{X: 42, Y: 11, C: "r", Fg: color.Teal, Bg: color.Silver},
		{X: 43, Y: 11, C: "l", Fg: color.Teal, Bg: color.Silver},
		{X: 44, Y: 11, C: "d", Fg: color.Teal, Bg: color.Silver},
		{X: 45, Y: 11, C: "!", Fg: color.Teal, Bg: color.Silver},
		{X: 31, Y: 13, C: "P", Fg: color.Silver, Bg: color.Black},
		{X: 32, Y: 13, C: "r", Fg: color.Silver, Bg: color.Black},
		{X: 33, Y: 13, C: "e", Fg: color.Silver, Bg: color.Black},
		{X: 34, Y: 13, C: "s", Fg: color.Silver, Bg: color.Black},
		{X: 35, Y: 13, C: "s", Fg: color.Silver, Bg: color.Black},
		{X: 36, Y: 13, C: " ", Fg: color.Silver, Bg: color.Black},
		{X: 37, Y: 13, C: "E", Fg: color.Silver, Bg: color.Black, Attr: vt.Bold},
		{X: 38, Y: 13, C: "S", Fg: color.Silver, Bg: color.Black, Attr: vt.Bold},
		{X: 39, Y: 13, C: "C", Fg: color.Silver, Bg: color.Black, Attr: vt.Bold},
		{X: 40, Y: 13, C: " ", Fg: color.Silver, Bg: color.Black},
		{X: 41, Y: 13, C: "t", Fg: color.Silver, Bg: color.Black},
		{X: 42, Y: 13, C: "o", Fg: color.Silver, Bg: color.Black},
		{X: 43, Y: 13, C: " ", Fg: color.Silver, Bg: color.Black},
		{X: 44, Y: 13, C: "e", Fg: color.Silver, Bg: color.Black},
		{X: 45, Y: 13, C: "x", Fg: color.Silver, Bg: color.Black},
		{X: 46, Y: 13, C: "i", Fg: color.Silver, Bg: color.Black},
		{X: 47, Y: 13, C: "t", Fg: color.Silver, Bg: color.Black},
		{X: 48, Y: 13, C: ".", Fg: color.Silver, Bg: color.Black},
	}

	for _, v := range expect {
		cell := mt.GetCell(vt.Coord{X: v.X, Y: v.Y})
		if v.C != string(cell.C) {
			t.Errorf("Mismatch string at %d,%d: %q != %q", v.X, v.Y, string(cell.C), v.C)
		}
		if v.Fg != cell.S.Fg() {
			t.Errorf("Mismatch foreground at %d,%d: %s != %s", v.X, v.Y, cell.S.Fg().String(), v.Fg.String())
		}
		if v.Bg != cell.S.Bg() {
			t.Errorf("Mismatch background at %d,%d: %s != %s", v.X, v.Y, cell.S.Bg().String(), v.Bg.String())
		}
		if v.Attr != cell.S.Attr() {
			t.Errorf("Mismatch attr at %d,%d: %x != %x", v.X, v.Y, cell.S.Attr(), v.Attr)
		}
	}

	mt.KeyEvent(vt.KeyEvent{Code: vt.KcEsc, Base: vt.KcEsc, Down: true})
	wg.Wait()
}
