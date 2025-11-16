// This file was originally generated automatically,
// but it is edited to correct for errors in the VT420
// terminfo data.  Additionally we have added extended
// information for the extended F-keys.

package vt420

import "github.com/gdamore/tcell/v2/terminfo"

func init() {

	// DEC VT420
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:              "vt420",
		Columns:           80,
		Lines:             24,
		Bell:              "\a",
		Clear:             "\x1b[H\x1b[2J$<50>",
		ShowCursor:        "\x1b[?25h",
		HideCursor:        "\x1b[?25l",
		AttrOff:           "\x1b[m\x1b(B$<2>",
		Underline:         "\x1b[4m",
		Bold:              "\x1b[1m$<2>",
		Blink:             "\x1b[5m$<2>",
		Reverse:           "\x1b[7m$<2>",
		EnterKeypad:       "\x1b=",
		ExitKeypad:        "\x1b>",
		PadChar:           "\x00",
		AltChars:          "``aaffggjjkkllmmnnooppqqrrssttuuvvwwxxyyzz{{||}}~~",
		EnterAcs:          "\x1b(0$<2>",
		ExitAcs:           "\x1b(B$<4>",
		EnableAcs:         "\x1b)0",
		EnableAutoMargin:  "\x1b[?7h",
		DisableAutoMargin: "\x1b[?7l",
		SetCursor:         "\x1b[%i%p1%d;%p2%dH$<10>",
		CursorBack1:       "\b",
		CursorUp1:         "\x1b[A",
		KeyUp:             "\x1b[A",
		KeyDown:           "\x1b[B",
		KeyRight:          "\x1b[C",
		KeyLeft:           "\x1b[D",
		KeyInsert:         "\x1b[2~",
		KeyDelete:         "\x1b[3~",
		KeyBackspace:      "\b",
		KeyPgUp:           "\x1b[5~",
		KeyPgDn:           "\x1b[6~",
		KeyF1:             "\x1bOP",
		KeyF2:             "\x1bOQ",
		KeyF3:             "\x1bOR",
		KeyF4:             "\x1bOS",
		KeyF6:             "\x1b[17~",
		KeyF7:             "\x1b[18~",
		KeyF8:             "\x1b[19~",
		KeyF9:             "\x1b[20~",
		KeyF10:            "\x1b[21~",
		KeyF11:            "\x1b[23~",
		KeyF12:            "\x1b[24~",
		KeyF13:            "\x1b[25~",
		KeyF14:            "\x1b[26~",
		KeyF15:            "\x1b[28~",
		KeyF16:            "\x1b[29~",
		KeyF17:            "\x1b[31~",
		KeyF18:            "\x1b[32~",
		KeyF19:            "\x1b[33~",
		KeyF20:            "\x1b[34~",
		AutoMargin:        true,
	})
}
