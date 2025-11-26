// Generated automatically.  DO NOT HAND-EDIT.

package linux

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// Linux console
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:      "linux",
		Colors:    8,
		SetFg:     "\x1b[3%p1%dm",
		SetBg:     "\x1b[4%p1%dm",
		SetFgBg:   "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg: "\x1b[39;49m",
		AltChars:  "++,,--..00``aaffgghhiijjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:  "\x0e",
		ExitAcs:   "\x0f",
		EnableAcs: "\x1b)0",
		Mouse:     "\x1b[M",
	})
}
