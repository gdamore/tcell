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

// mouse displays a text box and tests mouse interaction.  As you click
// and drag, boxes are displayed on screen.  Other events are reported in
// the box.  Press ESC twice to exit the program.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

var defStyle tcell.Style

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, r rune) {
	rs := string(r)

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.Put(col, y1, string(tcell.RuneHLine), style)
		s.Put(col, y2, string(tcell.RuneHLine), style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.Put(x1, row, string(tcell.RuneVLine), style)
		s.Put(x2, row, string(tcell.RuneVLine), style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.Put(x1, y1, string(tcell.RuneULCorner), style)
		s.Put(x2, y1, string(tcell.RuneURCorner), style)
		s.Put(x1, y2, string(tcell.RuneLLCorner), style)
		s.Put(x2, y2, string(tcell.RuneLRCorner), style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.Put(col, row, rs, style)
		}
	}
}

func drawSelect(s tcell.Screen, x1, y1, x2, y2 int, sel bool) {

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			str, style, width := s.Get(col, row)
			if style == tcell.StyleDefault {
				style = defStyle
			}
			style = style.Reverse(sel)
			s.Put(col, row, str, style)
			col += width - 1 // add an extra column if 2 cells
		}
	}
}

// This program just shows simple mouse and keyboard events.  Press ESC twice to
// exit.
func main() {

	shell := os.Getenv("SHELL")
	if shell == "" {
		if runtime.GOOS == "windows" {
			shell = "CMD.EXE"
		} else {
			shell = "/bin/sh"
		}
	}

	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetTitle("Tcell Mouse Demonstration")
	defStyle = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.EnableFocus()
	s.Clear()

	posfmt := "Mouse: %d, %d  "
	btnfmt := "Buttons: %s"
	keyfmt := "Keys: %s"
	pastefmt := "Paste: [%d] %s"
	focusfmt := "Focus: %s"
	style := tcell.StyleDefault.
		Foreground(tcell.ColorMidnightBlue).Background(tcell.ColorLightCoral)

	mx, my := -1, -1
	ox, oy := -1, -1
	bx, by := -1, -1
	w, h := s.Size()
	lchar := '*'
	bstr := ""
	lks := ""
	pstr := ""
	ecnt := 0
	pasting := false
	focus := true // assume we are focused when we start

	for {
		drawBox(s, 1, 1, 42, 8, style, ' ')
		s.PutStrStyled(2, 2, "Press ESC twice to exit, C to clear.", style)
		s.PutStrStyled(2, 3, fmt.Sprintf(posfmt, mx, my), style)
		s.PutStrStyled(2, 4, fmt.Sprintf(btnfmt, bstr), style)
		s.PutStrStyled(2, 5, fmt.Sprintf(keyfmt, lks), style)

		ps := pstr
		if len(ps) > 26 {
			ps = "..." + ps[len(ps)-24:]
		}
		s.PutStrStyled(2, 6, fmt.Sprintf(pastefmt, len(pstr), ps), style)

		fstr := "false"
		if focus {
			fstr = "true"
		}
		s.PutStrStyled(2, 7, fmt.Sprintf(focusfmt, fstr), style)

		s.Show()
		bstr = ""
		ev := s.PollEvent()
		st := tcell.StyleDefault.Background(tcell.ColorRed)
		up := tcell.StyleDefault.
			Background(tcell.ColorBlue).
			Foreground(tcell.ColorBlack)
		w, h = s.Size()

		// always clear any old selection box
		if ox >= 0 && oy >= 0 && bx >= 0 {
			drawSelect(s, ox, oy, bx, by, false)
		}

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			s.Put(w-1, h-1, "R", st)
		case *tcell.EventKey:
			s.Put(w-2, h-2, string(ev.Rune()), st)
			if pasting {
				s.Put(w-1, h-1, "P", st)
				if ev.Key() == tcell.KeyRune {
					pstr = pstr + string(ev.Rune())
				} else {
					pstr = pstr + "\ufffd" // replacement for now
				}
				lks = ""
				continue
			}
			pstr = ""
			s.Put(w-1, h-1, "K", st)
			if ev.Key() == tcell.KeyEscape {
				ecnt++
				if ecnt > 1 {
					s.Fini()
					os.Exit(0)
				}
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Key() == tcell.KeyCtrlZ {
				// CtrlZ doesn't really suspend the process, but we use it to execute a subshell.
				if err := s.Suspend(); err == nil {
					fmt.Printf("Executing shell (%s -l)...\n", shell)
					fmt.Printf("Exit the shell to return to the demo.\n")
					c := exec.Command(shell, "-l") // NB: -l works for cmd.exe too (ignored)
					c.Stdin = os.Stdin
					c.Stdout = os.Stdout
					c.Stderr = os.Stderr
					c.Run()
					if err := s.Resume(); err != nil {
						panic("failed to resume: " + err.Error())
					}
				}
			} else {
				ecnt = 0
				if ev.Rune() == 'C' || ev.Rune() == 'c' {
					s.Clear()
				}
			}
			lks = ev.Name()
		case *tcell.EventPaste:
			pasting = ev.Start()
			if pasting {
				pstr = ""
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			mods := ev.Modifiers()
			if mods&tcell.ModShift != 0 {
				bstr += " Shift"
			}
			if mods&tcell.ModCtrl != 0 {
				bstr += " Ctrl"
			}
			if mods&tcell.ModAlt != 0 {
				bstr += " Alt"
			}
			if mods&tcell.ModMeta != 0 {
				bstr += " Meta"
			}
			for i := uint(0); i < 8; i++ {
				if int(button)&(1<<i) != 0 {
					bstr += fmt.Sprintf(" Button%d", i+1)
				}
			}
			if button&tcell.WheelUp != 0 {
				bstr += " WheelUp"
			}
			if button&tcell.WheelDown != 0 {
				bstr += " WheelDown"
			}
			if button&tcell.WheelLeft != 0 {
				bstr += " WheelLeft"
			}
			if button&tcell.WheelRight != 0 {
				bstr += " WheelRight"
			}
			// Only buttons, not wheel events
			button &= tcell.ButtonMask(0xff)
			ch := '*'

			if button != tcell.ButtonNone && ox < 0 {
				ox, oy = x, y
			}
			theme := []tcell.Color{
				tcell.ColorGray,
				tcell.ColorRed,
				tcell.ColorLime,
				tcell.ColorYellow,
				tcell.ColorFuchsia,
				tcell.ColorBlue,
				tcell.ColorAqua,
				tcell.ColorSilver,
			}
			switch ev.Buttons() {
			case tcell.ButtonNone:
				if ox >= 0 {
					bg := theme[lchar%8]
					fg := tcell.ColorBlack
					drawBox(s, ox, oy, x, y,
						up.Background(bg).Foreground(fg),
						lchar)
					ox, oy = -1, -1
					bx, by = -1, -1
				}
			case tcell.Button1:
				ch = '1'
			case tcell.Button2:
				ch = '2'
			case tcell.Button3:
				ch = '3'
			case tcell.Button4:
				ch = '4'
			case tcell.Button5:
				ch = '5'
			case tcell.Button6:
				ch = '6'
			case tcell.Button7:
				ch = '7'
			case tcell.Button8:
				ch = '8'
			default:
				ch = '*'

			}
			if button != tcell.ButtonNone {
				bx, by = x, y
			}
			lchar = ch
			s.Put(w-1, h-1, "M", st)
			mx, my = x, y
		case *tcell.EventFocus:
			focus = ev.Focused
		default:
			s.Put(w-1, h-1, "X", st)
		}

		if ox >= 0 && bx >= 0 {
			drawSelect(s, ox, oy, bx, by, true)
		}
	}
}
