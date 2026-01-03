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
			interval = time.Microsecond * 100

			go func() {
				defer wg.Done()
				main()
			}()

			time.Sleep(time.Millisecond * 25)
			mt.KeyEvent(vt.KbdEvent{Code: 'L', Mod: vt.ModCtrl | vt.ModShift, Down: true})
			mt.SetSize(vt.Coord{X: 10, Y: 10})
			mt.Drain()
			time.Sleep(time.Millisecond * 25)
			mt.KeyEvent(vt.KbdEvent{Code: 'q', Mod: vt.ModCtrl, Down: true})
			mt.Drain()
			wg.Wait()

			if count < 10 { // CI runs *slow*
				t.Errorf("Too few boxes: %d", count)
			}
			if drawTime < time.Microsecond {
				t.Errorf("Interval too short: %s", drawTime)
			}
			if drawTime > 100*time.Millisecond { // longer because CI/CD timekeeping is awful
				t.Errorf("Interval too long: %s", drawTime)
			}
		})
	}
}
