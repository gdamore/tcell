// Generated automatically.  DO NOT HAND-EDIT.

package dtterm

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// CDE desktop terminal
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:      "dtterm",
		Columns:   80,
		Lines:     24,
		Colors:    8,
		SetFg:     "\x1b[3%p1%dm",
		SetBg:     "\x1b[4%p1%dm",
		SetFgBg:   "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg: "\x1b[39;49m",
		AltChars:  "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:  "\x0e",
		ExitAcs:   "\x0f",
		EnableAcs: "\x1b(B\x1b)0",
	})
}
