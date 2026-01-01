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
	mt.KeyEvent(vt.KbdEvent{Code: 'b', Down: true})
	mt.KeyEvent(vt.KbdEvent{Code: 'b', Down: false})
	time.Sleep(time.Millisecond * 100)
	mt.KeyEvent(vt.KbdEvent{Code: 'b', Down: true})
	mt.KeyEvent(vt.KbdEvent{Code: 'b', Down: false})

	// control L (forces sync)
	mt.KeyEvent(vt.KbdEvent{Code: 'l', Down: true, Mod: vt.ModCtrl})
	mt.KeyEvent(vt.KbdEvent{Code: 'l', Down: false, Mod: vt.ModCtrl})

	// sleep at least 3 seconds to get the time driven beep
	time.Sleep(time.Millisecond * 3500)

	mt.KeyEvent(vt.KbdEvent{Code: 'q', Down: true, Mod: vt.ModCtrl})
	mt.KeyEvent(vt.KbdEvent{Code: 'q', Down: false, Mod: vt.ModCtrl})

	wg.Wait()

	if cnt := mt.Bells(); cnt != 3 {
		t.Errorf("incorrect bell count %d != 2", cnt)
	}
}
