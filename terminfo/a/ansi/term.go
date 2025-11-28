// Generated automatically.  DO NOT HAND-EDIT.

package ansi

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// ansi/pc-term compatible with color
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:    "ansi",
		Aliases: []string{"pcansi"},
		Columns: 80,
		Lines:   24,
	})
}
