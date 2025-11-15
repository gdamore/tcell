package extended

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2/terminfo"
)

func TestAmbiguousKeys(t *testing.T) {
	keys := make(map[string]string)
	otrm := make(map[string][]string)

	addSeq := func(t *testing.T, seq string, name string, key string) {
		if seq == "" {
			return
		}
		switch keys[seq] {
		case key:
			otrm[seq] = append(otrm[seq], name)
		case "":
			keys[seq] = key
			otrm[seq] = append(otrm[seq], name)
		default:
			if (key == "Help" && keys[seq] == "F15") || (key == "F15" && keys[seq] == "Help") {
				// VT220 derived terminals co-mapped F15 and Help
				keys[seq] = "F15"
				otrm[seq] = append(otrm[seq], name)
			} else {
				t.Errorf("Conflicting map for %s %s with %s %s", name, key, strings.Join(otrm[seq], ","), keys[seq])
			}
		}
	}

	for _, name := range terminfo.TerminfoNames() {
		term, _ := terminfo.LookupTerminfo(name)
		println("STARTING", name)
		addSeq(t, term.KeyF1, name, "F1")
		addSeq(t, term.KeyF2, name, "F2")
		addSeq(t, term.KeyF3, name, "F3")
		addSeq(t, term.KeyF4, name, "F4")
		addSeq(t, term.KeyF5, name, "F5")
		addSeq(t, term.KeyF6, name, "F6")
		addSeq(t, term.KeyF7, name, "F7")
		addSeq(t, term.KeyF8, name, "F8")
		addSeq(t, term.KeyF9, name, "F9")
		addSeq(t, term.KeyF10, name, "F10")
		addSeq(t, term.KeyF11, name, "F11")
		addSeq(t, term.KeyF12, name, "F12")
		addSeq(t, term.KeyF13, name, "F13")
		addSeq(t, term.KeyF14, name, "F14")
		addSeq(t, term.KeyF15, name, "F15")
		addSeq(t, term.KeyF16, name, "F16")
		addSeq(t, term.KeyF17, name, "F17")
		addSeq(t, term.KeyF18, name, "F18")
		addSeq(t, term.KeyF19, name, "F19")
		addSeq(t, term.KeyF20, name, "F20")
		addSeq(t, term.KeyF21, name, "F21")
		addSeq(t, term.KeyF22, name, "F22")
		addSeq(t, term.KeyF23, name, "F23")
		addSeq(t, term.KeyF24, name, "F24")
		addSeq(t, term.KeyF25, name, "F25")
		addSeq(t, term.KeyF26, name, "F26")
		addSeq(t, term.KeyF27, name, "F27")
		addSeq(t, term.KeyF28, name, "F28")
		addSeq(t, term.KeyF29, name, "F29")
		addSeq(t, term.KeyF30, name, "F30")
		addSeq(t, term.KeyF31, name, "F31")
		addSeq(t, term.KeyF32, name, "F32")
		addSeq(t, term.KeyF33, name, "F33")
		addSeq(t, term.KeyF34, name, "F34")
		addSeq(t, term.KeyF35, name, "F35")

		addSeq(t, term.KeyInsert, name, "Insert")
		addSeq(t, term.KeyDelete, name, "Delete")
		addSeq(t, term.KeyHome, name, "Home")
		addSeq(t, term.KeyEnd, name, "End")
		addSeq(t, term.KeyHelp, name, "Help")
		addSeq(t, term.KeyPgDn, name, "PgDn")
		addSeq(t, term.KeyPgUp, name, "PgUp")
		addSeq(t, term.KeyUp, name, "Up")
		addSeq(t, term.KeyDown, name, "Down")
		addSeq(t, term.KeyLeft, name, "Left")
		addSeq(t, term.KeyRight, name, "Right")
		addSeq(t, term.KeyBacktab, name, "Backtab")
		addSeq(t, term.KeyExit, name, "Exit")
		addSeq(t, term.KeyClear, name, "Clear")
		addSeq(t, term.KeyPrint, name, "Print")
		addSeq(t, term.KeyCancel, name, "Cancel")
		addSeq(t, term.KeyBackspace, name, "Backspace")
	}

	for seq1 := range keys {
		for seq2 := range keys {
			if seq1 == seq2 {
				continue
			}
			if strings.HasPrefix(seq1, seq2) || strings.HasPrefix(seq2, seq1) {
				t.Errorf("Key %s from %s conflicts with %s from %s", keys[seq1], otrm[seq1], keys[seq2], otrm[seq2])
			}
		}
	}
}
