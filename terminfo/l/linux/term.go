// Generated automatically.  DO NOT HAND-EDIT.

package linux

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// Linux console
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:   "linux",
		Colors: 8,
		Mouse:  "\x1b[M",
	})
}
