//go:build ignore
// +build ignore

// Copyright 2024 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This command is used to generate suitable configuration files in either
// go syntax or in JSON.  It defaults to JSON output on stdout.  If no
// term values are specified on the command line, then $TERM is used.
//
// Usage is like this:
//
// mkinfo [-go file.go] [-quiet] [-nofatal] [-I <import>] [-P <pkg}] [<term>...]
//
// -go       specifies Go output into the named file.  Use - for stdout.
// -nofatal  indicates that errors loading definitions should not be fatal
// -P pkg    use the supplied package name
// -I import use the named import instead of github.com/gdamore/tcell/v2/terminfo
//

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2/terminfo"
)

type termcap struct {
	name    string
	desc    string
	aliases []string
	bools   map[string]bool
	nums    map[string]int
	strs    map[string]string
}

func (tc *termcap) getnum(s string) int {
	return (tc.nums[s])
}

func (tc *termcap) getflag(s string) bool {
	return (tc.bools[s])
}

func (tc *termcap) getstr(s string) string {
	return (tc.strs[s])
}

const (
	NONE = iota
	CTRL
	ESC
)

var notAddressable = errors.New("terminal not cursor addressable")

func unescape(s string) string {
	// Various escapes are in \x format.  Control codes are
	// encoded as ^M (carat followed by ASCII equivalent).
	// Escapes are: \e, \E - escape
	//  \0 NULL, \n \l \r \t \b \f \s for equivalent C escape.
	buf := &bytes.Buffer{}
	esc := NONE

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch esc {
		case NONE:
			switch c {
			case '\\':
				esc = ESC
			case '^':
				esc = CTRL
			default:
				buf.WriteByte(c)
			}
		case CTRL:
			buf.WriteByte(c ^ 1<<6)
			esc = NONE
		case ESC:
			switch c {
			case 'E', 'e':
				buf.WriteByte(0x1b)
			case '0', '1', '2', '3', '4', '5', '6', '7':
				if i+2 < len(s) && s[i+1] >= '0' && s[i+1] <= '7' && s[i+2] >= '0' && s[i+2] <= '7' {
					buf.WriteByte(((c - '0') * 64) + ((s[i+1] - '0') * 8) + (s[i+2] - '0'))
					i = i + 2
				} else if c == '0' {
					buf.WriteByte(0)
				}
			case 'n':
				buf.WriteByte('\n')
			case 'r':
				buf.WriteByte('\r')
			case 't':
				buf.WriteByte('\t')
			case 'b':
				buf.WriteByte('\b')
			case 'f':
				buf.WriteByte('\f')
			case 's':
				buf.WriteByte(' ')
			case 'l':
				panic("WTF: weird format: " + s)
			default:
				buf.WriteByte(c)
			}
			esc = NONE
		}
	}
	return (buf.String())
}

func (tc *termcap) setupterm(name string) error {
	cmd := exec.Command("infocmp", "-x", "-1", name)
	output := &bytes.Buffer{}
	cmd.Stdout = output

	tc.strs = make(map[string]string)
	tc.bools = make(map[string]bool)
	tc.nums = make(map[string]int)

	err := cmd.Run()
	if err != nil {
		return err
	}

	// Now parse the output.
	// We get comment lines (starting with "#"), followed by
	// a header line that looks like "<name>|<alias>|...|<desc>"
	// then capabilities, one per line, starting with a tab and ending
	// with a comma and newline.
	lines := strings.Split(output.String(), "\n")
	for len(lines) > 0 && strings.HasPrefix(lines[0], "#") {
		lines = lines[1:]
	}

	// Ditch trailing empty last line
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	header := lines[0]
	header = strings.TrimSuffix(header, ",")
	names := strings.Split(header, "|")
	tc.name = names[0]
	names = names[1:]
	if len(names) > 0 {
		tc.desc = names[len(names)-1]
		names = names[:len(names)-1]
	}
	tc.aliases = names
	for _, val := range lines[1:] {
		if (!strings.HasPrefix(val, "\t")) ||
			(!strings.HasSuffix(val, ",")) {
			return (errors.New("malformed infocmp: " + val))
		}

		val = val[1:]
		val = val[:len(val)-1]

		if k := strings.SplitN(val, "=", 2); len(k) == 2 {
			tc.strs[k[0]] = unescape(k[1])
		} else if k := strings.SplitN(val, "#", 2); len(k) == 2 {
			if u, err := strconv.ParseUint(k[1], 0, 0); err != nil {
				return (err)
			} else {
				tc.nums[k[0]] = int(u)
			}
		} else {
			tc.bools[val] = true
		}
	}
	return nil
}

// This program is used to collect data from the system's terminfo library,
// and write it into Go source code.  That is, we maintain our terminfo
// capabilities encoded in the program.  It should never need to be run by
// an end user, but developers can use this to add codes for additional
// terminal types.
func getinfo(name string) (*terminfo.Terminfo, string, error) {
	var tc termcap
	if err := tc.setupterm(name); err != nil {
		return nil, "", err
	}
	t := &terminfo.Terminfo{}
	// If this is an alias record, then just emit the alias
	t.Name = tc.name
	if t.Name != name {
		return t, "", nil
	}
	t.Aliases = tc.aliases
	t.Colors = tc.getnum("colors")
	t.Columns = tc.getnum("cols")
	t.Lines = tc.getnum("lines")
	t.Clear = tc.getstr("clear")
	t.EnterCA = tc.getstr("smcup")
	t.ExitCA = tc.getstr("rmcup")
	t.ShowCursor = tc.getstr("cnorm")
	t.HideCursor = tc.getstr("civis")
	t.AttrOff = tc.getstr("sgr0")
	t.Underline = tc.getstr("smul")
	t.Bold = tc.getstr("bold")
	t.Blink = tc.getstr("blink")
	t.Dim = tc.getstr("dim")
	t.Italic = tc.getstr("sitm")
	t.Reverse = tc.getstr("rev")
	t.EnterKeypad = tc.getstr("smkx")
	t.ExitKeypad = tc.getstr("rmkx")
	t.SetFg = tc.getstr("setaf")
	t.SetBg = tc.getstr("setab")
	t.ResetFgBg = tc.getstr("op")
	t.SetCursor = tc.getstr("cup")
	t.InsertChar = tc.getstr("ich1")
	t.AutoMargin = tc.getflag("am")
	t.AltChars = tc.getstr("acsc")
	t.EnterAcs = tc.getstr("smacs")
	t.ExitAcs = tc.getstr("rmacs")
	t.EnableAcs = tc.getstr("enacs")
	t.StrikeThrough = tc.getstr("smxx")
	t.Mouse = tc.getstr("kmous")
	t.EnableAutoMargin = tc.getstr("smam")
	t.DisableAutoMargin = tc.getstr("rmam")

	// Technically the RGB flag that is provided for xterm-direct is not
	// quite right.  The problem is that the -direct flag that was introduced
	// with ncurses 6.1 requires a parsing for the parameters that we lack.
	// For this case we'll just assume it's XTerm compatible.  Someday this
	// may be incorrect, but right now it is correct, and nobody uses it
	// anyway.
	if tc.getflag("Tc") {
		// This presumes XTerm 24-bit true color.
		t.TrueColor = true
	} else if tc.getflag("RGB") {
		// This is for xterm-direct, which uses a different scheme entirely.
		// (ncurses went a very different direction from everyone else, and
		// so it's unlikely anything is using this definition.)
		t.TrueColor = true
		t.SetBg = "\x1b[%?%p1%{8}%<%t4%p1%d%e%p1%{16}%<%t10%p1%{8}%-%d%e48;5;%p1%d%;m"
		t.SetFg = "\x1b[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m"
	}

	// We only support colors in ANSI 8 or 256 color mode.
	if t.Colors < 8 || t.SetFg == "" {
		t.Colors = 0
	}
	if t.SetCursor == "" {
		return nil, "", notAddressable
	}

	// For padding, we lookup the pad char.  If that isn't present,
	// and npc is *not* set, then we assume a null byte.
	t.PadChar = tc.getstr("pad")
	if t.PadChar == "" {
		if !tc.getflag("npc") {
			t.PadChar = "\u0000"
		}
	}

	// For terminals that use "standard" SGR sequences, lets combine the
	// foreground and background together.
	if strings.HasPrefix(t.SetFg, "\x1b[") &&
		strings.HasPrefix(t.SetBg, "\x1b[") &&
		strings.HasSuffix(t.SetFg, "m") &&
		strings.HasSuffix(t.SetBg, "m") {
		fg := t.SetFg[:len(t.SetFg)-1]
		r := regexp.MustCompile("%p1")
		bg := r.ReplaceAllString(t.SetBg[2:], "%p2")
		t.SetFgBg = fg + ";" + bg
	}

	if tc.getflag("XT") {
		t.XTermLike = true
	}
	return t, tc.desc, nil
}

func dotGoAddInt(w io.Writer, n string, i int) {
	if i == 0 {
		// initialized to 0, ignore
		return
	}
	fmt.Fprintf(w, "\t\t%-13s %d,\n", n+":", i)
}
func dotGoAddStr(w io.Writer, n string, s string) {
	if s == "" {
		return
	}
	fmt.Fprintf(w, "\t\t%-13s %q,\n", n+":", s)
}
func dotGoAddFlag(w io.Writer, n string, b bool) {
	if !b {
		// initialized to 0, ignore
		return
	}
	fmt.Fprintf(w, "\t\t%-13s true,\n", n+":")
}

func dotGoAddArr(w io.Writer, n string, a []string) {
	if len(a) == 0 {
		return
	}
	fmt.Fprintf(w, "\t\t%-13s []string{", n+":")
	did := false
	for _, b := range a {
		if did {
			fmt.Fprint(w, ", ")
		}
		did = true
		fmt.Fprintf(w, "%q", b)
	}
	fmt.Fprintln(w, "},")
}

func dotGoHeader(w io.Writer, packname, tipackname string) {
	fmt.Fprintln(w, "// Generated automatically.  DO NOT HAND-EDIT.")
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "package %s\n", packname)
	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "import \"%s\"\n", tipackname)
	fmt.Fprintln(w, "")
}

func dotGoTrailer(w io.Writer) {
}

func dotGoInfo(w io.Writer, terms []*TData) {

	fmt.Fprintln(w, "func init() {")
	for _, t := range terms {
		fmt.Fprintf(w, "\n\t// %s\n", t.Desc)
		fmt.Fprintln(w, "\tterminfo.AddTerminfo(&terminfo.Terminfo{")
		dotGoAddStr(w, "Name", t.Name)
		dotGoAddArr(w, "Aliases", t.Aliases)
		dotGoAddInt(w, "Columns", t.Columns)
		dotGoAddInt(w, "Lines", t.Lines)
		dotGoAddInt(w, "Colors", t.Colors)
		dotGoAddStr(w, "Clear", t.Clear)
		dotGoAddStr(w, "EnterCA", t.EnterCA)
		dotGoAddStr(w, "ExitCA", t.ExitCA)
		dotGoAddStr(w, "ShowCursor", t.ShowCursor)
		dotGoAddStr(w, "HideCursor", t.HideCursor)
		dotGoAddStr(w, "AttrOff", t.AttrOff)
		dotGoAddStr(w, "Underline", t.Underline)
		dotGoAddStr(w, "Bold", t.Bold)
		dotGoAddStr(w, "Dim", t.Dim)
		dotGoAddStr(w, "Italic", t.Italic)
		dotGoAddStr(w, "Blink", t.Blink)
		dotGoAddStr(w, "Reverse", t.Reverse)
		dotGoAddStr(w, "EnterKeypad", t.EnterKeypad)
		dotGoAddStr(w, "ExitKeypad", t.ExitKeypad)
		dotGoAddStr(w, "SetFg", t.SetFg)
		dotGoAddStr(w, "SetBg", t.SetBg)
		dotGoAddStr(w, "SetFgBg", t.SetFgBg)
		dotGoAddStr(w, "ResetFgBg", t.ResetFgBg)
		dotGoAddStr(w, "PadChar", t.PadChar)
		dotGoAddStr(w, "AltChars", t.AltChars)
		dotGoAddStr(w, "EnterAcs", t.EnterAcs)
		dotGoAddStr(w, "ExitAcs", t.ExitAcs)
		dotGoAddStr(w, "EnableAcs", t.EnableAcs)
		dotGoAddStr(w, "EnableAutoMargin", t.EnableAutoMargin)
		dotGoAddStr(w, "DisableAutoMargin", t.DisableAutoMargin)
		dotGoAddStr(w, "SetFgRGB", t.SetFgRGB)
		dotGoAddStr(w, "SetBgRGB", t.SetBgRGB)
		dotGoAddStr(w, "SetFgBgRGB", t.SetFgBgRGB)
		dotGoAddStr(w, "StrikeThrough", t.StrikeThrough)
		dotGoAddStr(w, "Mouse", t.Mouse)
		dotGoAddStr(w, "SetCursor", t.SetCursor)
		dotGoAddFlag(w, "TrueColor", t.TrueColor)
		dotGoAddFlag(w, "AutoMargin", t.AutoMargin)
		dotGoAddStr(w, "InsertChar", t.InsertChar)
		dotGoAddFlag(w, "XTermLike", t.XTermLike)
		fmt.Fprintln(w, "\t})")
	}
	fmt.Fprintln(w, "}")
}

var packname = ""
var tipackname = "github.com/gdamore/tcell/v2/terminfo"

func dotGoFile(fname string, terms []*TData) error {
	w := os.Stdout
	var e error
	if fname != "-" && fname != "" {
		if w, e = os.Create(fname); e != nil {
			return e
		}
	}
	if packname == "" {
		packname = strings.Replace(terms[0].Name, "-", "_", -1)
	}
	dotGoHeader(w, packname, tipackname)
	dotGoInfo(w, terms)
	dotGoTrailer(w)
	if w != os.Stdout {
		w.Close()
	}
	cmd := exec.Command("go", "fmt", fname)
	cmd.Run()
	return nil
}

type TData struct {
	Desc string

	terminfo.Terminfo
}

func main() {
	gofile := ""
	nofatal := false
	quiet := false
	all := false

	flag.StringVar(&gofile, "go", "", "generate go source in named file")
	flag.StringVar(&tipackname, "I", tipackname, "import package path")
	flag.StringVar(&packname, "P", packname, "package name (go source)")
	flag.BoolVar(&nofatal, "nofatal", false, "errors are not fatal")
	flag.BoolVar(&quiet, "quiet", false, "suppress error messages")
	flag.BoolVar(&all, "all", false, "load all terminals from terminfo")
	flag.Parse()
	var e error

	args := flag.Args()
	if len(args) == 0 {
		args = []string{os.Getenv("TERM")}
	}

	tdata := make([]*TData, 0)

	for _, term := range args {
		if t, desc, e := getinfo(term); e != nil {
			if all && e == notAddressable {
				continue
			}
			if !quiet {
				fmt.Fprintf(os.Stderr,
					"Failed loading %s: %v\n", term, e)
			}
			if !nofatal {
				os.Exit(1)
			}
		} else {
			tdata = append(tdata, &TData{
				Desc:     desc,
				Terminfo: *t,
			})
		}
	}

	if len(tdata) == 0 {
		// No data.
		os.Exit(0)
	}

	e = dotGoFile(gofile, tdata)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Failed %s: %v", gofile, e)
		os.Exit(1)
	}
}
