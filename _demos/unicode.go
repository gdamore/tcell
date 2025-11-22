//go:build ignore
// +build ignore

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

// unicode just displays a Unicode test on your screen.
// Press ESC to exit the program.
package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

var row = 0
var style = tcell.StyleDefault

func putln(s tcell.Screen, str string) {

	s.PutStrStyled(1, row, str, style)
	row++
}

func main() {

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	plain := tcell.StyleDefault
	bold := style.Bold(true)

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()

	// we can even try to use unicode window titles!
	s.SetTitle("Unicode Demonstration -- 🤯")

	quit := make(chan struct{})

	style = bold
	putln(s, "Press ESC to Exit")
	putln(s, "Character set: "+s.CharacterSet())
	style = plain

	putln(s, "English:   October")
	putln(s, "Icelandic: október")
	putln(s, "Arabic:    أكتوبر")
	putln(s, "Russian:   октября")
	putln(s, "Greek:     Οκτωβρίου")
	putln(s, "Chinese:   十月 (note, two double wide characters)")
	putln(s, "Combining: A\u030a (should look like Angstrom)")
	putln(s, "Emoticon:  \U0001f618 (blowing a kiss)")
	putln(s, "Airplane:  \u2708 (fly away)")
	putln(s, "Command:   \u2318 (mac clover key)")
	putln(s, "Enclose:   !\u20e3 (should be enclosed exclamation)")
	putln(s, "ZWJ:       \U0001f9db\u200d\u2640 (female vampire)")
	putln(s, "ZWJ:       \U0001f9db\u200d\u2642 (male vampire)")
	putln(s, "Family:    \U0001f469\u200d\U0001f467\u200d\U0001f467 (woman girl girl)\n")
	putln(s, "Region:    \U0001f1fa\U0001f1f8 (USA! USA!)\n")
	putln(s, "")
	putln(s, "Box:")
	putln(s, "┌─┬─┬──┐")
	putln(s, "│·│§│月│ (bullet, lantern, Swiss)")
	putln(s, "├─┼─┼──┤")
	putln(s, "│A│1│😘│ (A, 1, Kiss)")
	putln(s, "├─┼─┼──┤")
	putln(s, "│·│§│🇨🇭│ (bullet, lantern, Swiss)")
	putln(s, "├─┼─┼──┤")
	putln(s, "│◆│↑│  │ (diamond, up arrow, empty)")
	putln(s, "└─┴─┴──┘")

	s.Show()
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	<-quit

	s.Fini()
}
