// Generated automatically.  DO NOT HAND-EDIT.

package vt100

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// DEC VT100 (w/advanced video)
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:        "vt100",
		Aliases:     []string{"vt100-am", "vt102"},
		Columns:     80,
		Lines:       24,
		EnterKeypad: "\x1b[?1h\x1b=",
		ExitKeypad:  "\x1b[?1l\x1b>",
		AltChars:    "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:    "\x0e",
		ExitAcs:     "\x0f",
		EnableAcs:   "\x1b(B\x1b)0",
	})
}
