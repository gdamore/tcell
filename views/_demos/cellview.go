// Copyright 2015 The Tcell Authors
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

package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type model struct {
	x    int
	y    int
	endx int
	endy int
	hide bool
	enab bool
	loc  string
}

func (m *model) GetBounds() (int, int) {
	return m.endx, m.endy
}

func (m *model) MoveCursor(offx, offy int) {
	m.x += offx
	m.y += offy
	m.limitCursor()
}

func (m *model) limitCursor() {
	if m.x < 0 {
		m.x = 0
	}
	if m.x > m.endx-1 {
		m.x = m.endx - 1
	}
	if m.y < 0 {
		m.y = 0
	}
	if m.y > m.endy-1 {
		m.y = m.endy - 1
	}
	m.loc = fmt.Sprintf("Cursor is %d,%d", m.x, m.y)
}

func (m *model) GetCursor() (int, int, bool, bool) {
	return m.x, m.y, m.enab, !m.hide
}

func (m *model) SetCursor(x int, y int) {
	m.x = x
	m.y = y

	m.limitCursor()
}

func (m *model) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	dig := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	var ch rune
	style := tcell.StyleDefault
	if x >= 60 || y >= 15 {
		return ch, style, nil, 1
	}
	colors := []tcell.Color{
		tcell.ColorWhite,
		tcell.ColorGreen,
		tcell.ColorMaroon,
		tcell.ColorNavy,
		tcell.ColorOlive,
	}
	if y == 0 && x < len(m.loc) {
		style = style.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorLime)
		ch = rune(m.loc[x])
	} else {
		ch = dig[(x)%len(dig)]
		style = style.
			Foreground(colors[(y)%len(colors)]).
			Background(tcell.ColorBlack)
	}
	return ch, style, nil, 1
}

type MyApp struct {
	view   views.View
	main   *views.CellView
	keybar *views.SimpleStyledText
	status *views.SimpleStyledText
	model  *model

	views.Panel
}

func (a *MyApp) HandleEvent(ev tcell.Event) bool {

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlL:
			views.AppRedraw()
			return true
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'Q', 'q':
				views.AppQuit()
				return true
			case 'S', 's':
				a.model.hide = false
				a.updateKeys()
				return true
			case 'H', 'h':
				a.model.hide = true
				a.updateKeys()
				return true
			case 'E', 'e':
				a.model.enab = true
				a.updateKeys()
				return true
			case 'D', 'd':
				a.model.enab = false
				a.updateKeys()
				return true
			}
		}
	}
	return a.Panel.HandleEvent(ev)
}

func (a *MyApp) Draw() {
	a.status.SetMarkup(a.model.loc)
	a.Panel.Draw()
}

func (a *MyApp) updateKeys() {
	m := a.model
	w := "[%AQ%N] Quit"
	if !m.enab {
		w += "  [%AE%N] Enable cursor"
	} else {
		w += "  [%AD%N] Disable cursor"
		if !m.hide {
			w += "  [%AH%N] Hide cursor"
		} else {
			w += "  [%AS%N] Show cursor"
		}
	}
	a.keybar.SetMarkup(w)
	views.AppDraw()
}

func main() {

	app := &MyApp{}

	app.model = &model{endx: 60, endy: 15}

	if e := views.AppInit(); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(1)
	}

	title := views.NewTextBar()
	title.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorTeal).
		Foreground(tcell.ColorWhite))
	title.SetCenter("CellView Test", tcell.StyleDefault)
	title.SetRight("Example v1.0", tcell.StyleDefault)

	app.keybar = views.NewSimpleStyledText()
	app.keybar.SetStyleN(tcell.StyleDefault.
		Background(tcell.ColorSilver).
		Foreground(tcell.ColorBlack))
	app.keybar.SetStyleA(tcell.StyleDefault.
		Background(tcell.ColorSilver).
		Foreground(tcell.ColorRed))
	app.keybar.SetMarkup("[%AQ%N] Quit")

	app.status = views.NewSimpleStyledText()
	app.status.SetStyleN(tcell.StyleDefault.
		Background(tcell.ColorYellow).
		Foreground(tcell.ColorBlack))
	app.status.SetMarkup("My status is here.")

	app.main = views.NewCellView()
	app.main.SetModel(app.model)
	app.main.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorBlack))

	app.SetMenu(app.keybar)
	app.SetTitle(title)
	app.SetContent(app.main)
	app.SetStatus(app.status)

	app.updateKeys()
	views.AppSetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	views.SetApplication(app)
	views.RunApplication()
}
