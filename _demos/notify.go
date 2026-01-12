//go:build ignore
// +build ignore

// Copyright 2025 The TCell Authors
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
	"os"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func displayHelloWorld(s tcell.Screen, secs int) {
	w, h := s.Size()
	s.Clear()
	style := tcell.StyleDefault.Foreground(color.CadetBlue.TrueColor()).Background(color.White)
	msg := "Notification Demo"
	s.PutStrStyled((w-len(msg))/2, h/2-1, msg, style)
	msg = "(Minimize This Window)"
	s.PutStrStyled((w-len(msg))/2, h/2+1, msg, style)
	if secs > 0 {
		msg = fmt.Sprintf("Incoming in %d Seconds", secs)
	} else {
		msg = "Notification Sent!"
	}
	s.PutStr((w-len(msg))/2, h/2+3, msg)
	msg = "Press ESC to exit, ENTER to restart."
	s.PutStr((w-len(msg))/2, h/2+5, msg)
	s.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(color.Black).
		Foreground(color.White)
	s.SetStyle(defStyle)

	when := 10
	displayHelloWorld(s, when)

	ticker := time.NewTicker(time.Second)
	var ev tcell.Event
	for {
		select {
		case ev = <-s.EventQ():
		case <-ticker.C:
			if when > 0 {
				when--
				if when == 0 {
					s.ShowNotification("Ding Dong!", "The wicked witch is dead.")
				}
			}
			displayHelloWorld(s, when)
			continue
		}
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter:
				when = 10
			case tcell.KeyEscape:
				s.Fini()
				os.Exit(0)
			}
		}
		displayHelloWorld(s, when)
	}
}
