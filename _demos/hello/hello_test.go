package main

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/gdamore/tcell/v3/mock"
	"github.com/gdamore/tcell/v3/vt"
)

func TestHello(t *testing.T) {
	// ensure we only use 8 color ANSI for now
	os.Setenv("TERM", "ansi")
	os.Setenv("COLORTERM", "")

	mt := mock.NewMockTerm(mock.MockOptColors(8))
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

	expect := []struct {
		X    vt.Col
		Y    vt.Row
		C    string
		Fg   color.Color
		Bg   color.Color
		Attr vt.Attr
	}{
		{X: 33, Y: 12, C: "H", Fg: color.Teal, Bg: color.Silver},
		{X: 34, Y: 12, C: "e", Fg: color.Teal, Bg: color.Silver},
		{X: 35, Y: 12, C: "l", Fg: color.Teal, Bg: color.Silver},
		{X: 36, Y: 12, C: "l", Fg: color.Teal, Bg: color.Silver},
		{X: 37, Y: 12, C: "o", Fg: color.Teal, Bg: color.Silver},
		{X: 38, Y: 12, C: ",", Fg: color.Teal, Bg: color.Silver},
		{X: 39, Y: 12, C: " ", Fg: color.Teal, Bg: color.Silver},
		{X: 40, Y: 12, C: "W", Fg: color.Teal, Bg: color.Silver},
		{X: 41, Y: 12, C: "o", Fg: color.Teal, Bg: color.Silver},
		{X: 42, Y: 12, C: "r", Fg: color.Teal, Bg: color.Silver},
		{X: 43, Y: 12, C: "l", Fg: color.Teal, Bg: color.Silver},
		{X: 44, Y: 12, C: "d", Fg: color.Teal, Bg: color.Silver},
		{X: 45, Y: 12, C: "!", Fg: color.Teal, Bg: color.Silver},
	}

	for _, v := range expect {
		cell := mt.GetCell(vt.Coord{X: v.X, Y: v.Y})
		if v.C != string(cell.C) {
			t.Errorf("Mismatch string at %d,%d: %q != %q", v.X, v.Y, string(cell.C), v.C)
		}
		if v.Fg != cell.Fg {
			t.Errorf("Mismatch foreground: %s != %s", cell.Fg.String(), v.Fg.String())
		}
		if v.Bg != cell.Bg {
			t.Errorf("Mismatch foreground: %s != %s", cell.Fg.String(), v.Fg.String())
		}
	}

	mt.KeyEvent(vt.KbdEvent{Code: vt.KcEsc, Base: vt.KcEsc, Down: true})
	wg.Wait()
}
