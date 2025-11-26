// Generated automatically.  DO NOT HAND-EDIT.

package vt100

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// DEC VT100 (w/advanced video)
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "vt100",
		Aliases:           []string{"vt100-am", "vt102"},
		Columns:           80,
		Lines:             24,
		Clear:             "\x1b[H\x1b[J",
		AttrOff:           "\x1b[m\x0f",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m",
		Blink:             "\x1b[5m",
		Reverse:           "\x1b[7m",
		EnterKeypad:       "\x1b[?1h\x1b=",
		ExitKeypad:        "\x1b[?1l\x1b>",
		AltChars:          "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x0e",
		ExitAcs:           "\x0f",
		EnableAcs:         "\x1b(B\x1b)0",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH",
	})
}
