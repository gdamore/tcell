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

package views

import (
	"github.com/gdamore/tcell"
)

var appWidget Widget
var appScreen tcell.Screen

func SetApplication(app Widget) {
	appWidget = app
}

func AppInit() error {
	if appScreen == nil {
		if scr, e := tcell.NewScreen(); e != nil {
			return e
		} else {
			appScreen = scr
		}
	}
	return nil
}

// AppQuit causes the application to shutdown gracefully.
func AppQuit() {
	ev := &eventAppQuit{}
	ev.SetEventNow()
	if appScreen != nil {
		go func() { appScreen.PostEventWait(ev) }()
	}
}

// AppRedraw causes the application forcibly redraw everything.  Use this
// to clear up screen corruption, etc.
func AppRedraw() {
	ev := &eventAppRedraw{}
	ev.SetEventNow()
	if appScreen != nil {
		go func() { appScreen.PostEventWait(ev) }()
	}
}

// AppDraw asks the application to redraw itself, but only parts of it that
// are changed are updated.
func AppDraw() {
	ev := &eventAppDraw{}
	ev.SetEventNow()
	if appScreen != nil {
		go func() { appScreen.PostEventWait(ev) }()
	}
}

// AppPostFunc posts a function to be executed in the context of the
// application's event loop.  Functions that need to update displayed
// state, etc. can do this to avoid holding locks.
func AppPostFunc(fn func()) {
	ev := &eventAppFunc{fn: fn}
	ev.SetEventNow()
	if appScreen != nil {
		go func() { appScreen.PostEventWait(ev) }()
	}
}

func SetScreen(scr tcell.Screen) {
	appScreen = scr
}

func RunApplication() {

	if appScreen == nil {
		return
	}
	if appWidget == nil {
		return
	}
	scr := appScreen
	scr.Init()

	scr.Clear()
	appWidget.SetView(appScreen)

loop:
	for {
		appWidget.Draw()
		scr.Show()

		ev := scr.PollEvent()
		switch nev := ev.(type) {
		case *eventAppQuit:
			break loop
		case *eventAppDraw:
			scr.Show()
		case *eventAppRedraw:
			scr.Sync()
		case *eventAppFunc:
			nev.fn()
		case *tcell.EventResize:
			scr.Sync()
			appWidget.Resize()
		default:
			appWidget.HandleEvent(ev)
		}
	}
	scr.Fini()
}

type eventAppDraw struct {
	tcell.EventTime
}

type eventAppQuit struct {
	tcell.EventTime
}

type eventAppRedraw struct {
	tcell.EventTime
}

type eventAppFunc struct {
	tcell.EventTime
	fn func()
}
