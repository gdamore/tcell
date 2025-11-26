// Generated automatically.  DO NOT HAND-EDIT.

package emacs

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// GNU Emacs term.el terminal emulation
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:    "eterm",
		Columns: 80,
		Lines:   24,
		EnterCA: "\x1b7\x1b[?47h",
		ExitCA:  "\x1b[2J\x1b[?47l\x1b8",
	})

	// Emacs term.el terminal emulator term-protocol-version 0.96
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:      "eterm-color",
		Columns:   80,
		Lines:     24,
		Colors:    8,
		EnterCA:   "\x1b7\x1b[?47h",
		ExitCA:    "\x1b[2J\x1b[?47l\x1b8",
		SetFg:     "\x1b[%p1%{30}%+%dm",
		SetBg:     "\x1b[%p1%'('%+%dm",
		SetFgBg:   "\x1b[%p1%{30}%+%d;%p2%'('%+%dm",
		ResetFgBg: "\x1b[39;49m",
	})
}
