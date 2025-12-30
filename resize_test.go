// Copyright 2025 The TCell Authors
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

package tcell

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v3/vt"
)

func TestPixelSize(t *testing.T) {
	// simulate 8x16 font, VGA (640x480)
	ev := NewEventResize(80, 25)
	ev.ws.PixelHeight = ev.ws.Height * 16
	ev.ws.PixelWidth = ev.ws.Width * 8
	if w, h := ev.PixelSize(); w != 640 || h != 400 {
		t.Errorf("pixelsize wrong: %d %d", w, h)
	}
}

// TestEventResize creates a screen and does some resizing to prove it works.
func TestEventResize(t *testing.T) {
	mt, ms := NewMockScreen(t)
	defer ms.Fini()

	// drain the eventQ
loop:
	for {
		select {
		case <-ms.EventQ():
		default:
			break loop
		}
	}

	mt.SetSize(vt.Coord{X: 50, Y: 132})

	select {
	case ev := <-ms.EventQ():
		if re, ok := ev.(*EventResize); ok {
			if y, x := re.Size(); y != 50 || x != 132 {
				t.Errorf("wrong size: %d %d", y, x)
			}
		} else {
			t.Errorf("wrong event type %T", ev)
		}
	case <-time.After(time.Millisecond * 100):
		t.Errorf("never got resize signal")
	}
}
