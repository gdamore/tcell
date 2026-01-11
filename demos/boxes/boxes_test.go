package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/vt"
)

func TestBoxes(t *testing.T) {

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
			count = 0
			drawTime = 0
			interval = time.Microsecond * 10

			go func() {
				defer wg.Done()
				main()
			}()

			time.Sleep(time.Millisecond * 25)
			mt.KeyEvent(vt.KeyEvent{Code: 'L', Mod: vt.ModCtrl | vt.ModShift, Down: true})
			mt.SetSize(vt.Coord{X: 10, Y: 10})
			mt.Drain()
			time.Sleep(time.Millisecond * 25)
			mt.KeyEvent(vt.KeyEvent{Code: 'q', Mod: vt.ModCtrl, Down: true})
			mt.Drain()
			wg.Wait()

			if count < 10 { // CI runs *slow*
				t.Errorf("Too few boxes: %d", count)
			}
			if drawTime < time.Microsecond {
				// on windows this can happen because our tick counter is too coarse,
				// and our mock screen is basically limited only by CPU.
				t.Logf("Draw time very short: %s", drawTime)
			}

			// It should not take 10 milliseconds to draw a box,
			// as we generally see values sub millisecond here.
			if drawTime > 10*time.Millisecond*time.Duration(count) {
				t.Errorf("Draw time too long: %s", drawTime)
			}
		})
	}
}
