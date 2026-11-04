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

package tcell

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/vt"
)

// Resize notifications must not reconcile the screen before Show or Sync can
// observe a TTY size that changed during Resume and then changed back.
type quietResizeTty struct{ vt.MockTerm }

func (quietResizeTty) NotifyResize(chan<- bool) {}

func TestTerminfoScreenResize(t *testing.T) {
	if runtime.GOOS == "js" {
		t.Skip("not supported on webasm")
	}
	initial := vt.Coord{X: 101, Y: 54}
	cases := []struct {
		name   string
		resume bool
		size   vt.Coord
	}{
		{"transient_width_shrink", true, vt.Coord{X: 100, Y: 54}},
		{"transient_width_grow", true, vt.Coord{X: 102, Y: 54}},
		{"transient_height_shrink", true, vt.Coord{X: 101, Y: 53}},
		{"transient_height_grow", true, vt.Coord{X: 101, Y: 55}},
		{"unchanged_resume", true, initial},
		{"resize_width_shrink", false, vt.Coord{X: 100, Y: 54}},
		{"resize_width_grow", false, vt.Coord{X: 102, Y: 54}},
		{"resize_height_shrink", false, vt.Coord{X: 101, Y: 53}},
		{"resize_height_grow", false, vt.Coord{X: 101, Y: 55}},
		{"unchanged_show", false, initial},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, method := range []string{"Show", "Sync"} {
				t.Run(method, func(t *testing.T) {
					if os.Getenv("TCELL_RESIZE_TEST_CHILD") != "1" {
						// A stalled draw holds the screen lock and prevents Fini.
						// Kill the child instead of leaking a blocked goroutine.
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()
						exe, err := os.Executable()
						if err != nil {
							t.Fatal(err)
						}
						run := "^" + strings.ReplaceAll(t.Name(), "/", "$/^") + "$"
						cmd := exec.CommandContext(ctx, exe, "-test.run="+run, "-test.timeout=5s", "-test.v")
						cmd.Env = append(os.Environ(), "TCELL_RESIZE_TEST_CHILD=1")
						out, err := cmd.CombinedOutput()
						if err != nil {
							t.Fatalf("resize subprocess: %v (context: %v)\n%s", err, ctx.Err(), out)
						}
						return
					}

					t.Setenv("LC_ALL", "C.UTF-8")
					mt := quietResizeTty{MockTerm: vt.NewMockTerm(vt.MockOptSize(initial))}
					s, err := NewTerminfoScreenFromTty(mt, OptNegotiation(false))
					if err != nil {
						t.Fatal(err)
					}
					if err := s.Init(); err != nil {
						t.Fatal(err)
					}
					defer s.Fini()
					s.Show()

					want := tc.size
					if tc.resume {
						if err := s.Suspend(); err != nil {
							t.Fatal(err)
						}
						mt.SetSize(tc.size)
						if err := s.Resume(); err != nil {
							t.Fatal(err)
						}
						want = initial
					}
					mt.SetSize(want)
					draw := s.Show
					if method == "Sync" {
						draw = s.Sync
					}
					draw()
					if w, h := s.Size(); w != int(want.X) || h != int(want.Y) {
						t.Fatalf("screen size = %dx%d, want %dx%d", w, h, want.X, want.Y)
					}
					for _, pos := range []vt.Coord{{}, {X: want.X - 1, Y: want.Y - 1}} {
						if _, _, width := s.Get(int(pos.X), int(pos.Y)); width != 1 {
							t.Fatalf("cell %v has width %d, want 1", pos, width)
						}
					}
					for _, pos := range []vt.Coord{{X: want.X}, {Y: want.Y}} {
						if _, _, width := s.Get(int(pos.X), int(pos.Y)); width != 0 {
							t.Fatalf("out-of-bounds cell %v has width %d, want 0", pos, width)
						}
					}
					pos := vt.Coord{X: want.X - 1, Y: want.Y - 1}
					s.Put(int(pos.X), int(pos.Y), "Z", StyleDefault)
					draw()
					if cell := mt.GetCell(pos); cell.C != "Z" || cell.W != 1 {
						t.Fatalf("rendered bottom-right cell = %q/%d, want Z/1", cell.C, cell.W)
					}
				})
			}
		})
	}
}
