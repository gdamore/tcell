// Generated automatically.  DO NOT HAND-EDIT.

package aixterm

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// IBM Aixterm Terminal Emulator
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:      "aixterm",
		Columns:   80,
		Lines:     25,
		Colors:    8,
		SetFg:     "\x1b[3%p1%dm",
		SetBg:     "\x1b[4%p1%dm",
		SetFgBg:   "\x1b[3%p1%d;4%p2%dm",
		ResetFgBg: "\x1b[32m\x1b[40m",
	})
}
