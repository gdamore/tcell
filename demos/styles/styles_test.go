package main

import (
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/mock"
	"github.com/gdamore/tcell/v3/vt"
)

// TestStyles just exercises the code in the styles demo program.
// It does not validate that the content is accurate, that should be done
// manually by running the program with a real terminal.
func TestStyles(t *testing.T) {

	mt := mock.NewMockTerm()
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
	mt.KeyEvent(vt.KbdEvent{Code: 'l', Down: true, Mod: vt.ModCtrl})
	mt.KeyEvent(vt.KbdEvent{Code: 'l', Down: false, Mod: vt.ModCtrl})

	attrs := 0
	var lastAttr vt.Attr
	for y := range vt.Row(24) {
		for x := range vt.Col(80) {
			if attr := mt.GetCell(vt.Coord{X: vt.Col(x), Y: vt.Row(y)}).Attr; attr != lastAttr {
				attrs++
				lastAttr = attr
			}
		}
	}
	mt.KeyEvent(vt.KbdEvent{Code: 'q', Down: true, Mod: vt.ModCtrl})
	mt.KeyEvent(vt.KbdEvent{Code: 'q', Down: false, Mod: vt.ModCtrl})

	wg.Wait()
	if attrs < 8 {
		t.Errorf("Not enough attribute changes changes: %d", attrs)
	}
}
