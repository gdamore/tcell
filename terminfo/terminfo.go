// Copyright 2025 The TCell Authors
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

package terminfo

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// ErrTermNotFound indicates that a suitable terminal entry could
	// not be found.  This can result from either not having TERM set,
	// or from the TERM failing to support certain minimal functionality,
	// in particular absolute cursor addressability (the cup capability)
	// is required.  For example, legacy "adm3" lacks this capability,
	// whereas the slightly newer "adm3a" supports it.  This failure
	// occurs most often with "dumb".
	ErrTermNotFound = errors.New("terminal entry not found")
)

// Terminfo represents a terminfo entry.  Note that we use friendly names
// in Go, but when we write out JSON, we use the same names as terminfo.
// The name, aliases and smous, rmous fields do not come from terminfo directly.
type Terminfo struct {
	Name    string
	Aliases []string
	Columns int // cols
	Lines   int // lines
}

type stack []any

func (st stack) Push(v any) stack {
	if b, ok := v.(bool); ok {
		if b {
			return append(st, 1)
		} else {
			return append(st, 0)
		}
	}
	return append(st, v)
}

func (st stack) PopString() (string, stack) {
	if len(st) > 0 {
		e := st[len(st)-1]
		var s string
		switch v := e.(type) {
		case int:
			s = strconv.Itoa(v)
		case string:
			s = v
		}
		return s, st[:len(st)-1]
	}
	return "", st

}
func (st stack) PopInt() (int, stack) {
	if len(st) > 0 {
		e := st[len(st)-1]
		var i int
		switch v := e.(type) {
		case int:
			i = v
		case string:
			i, _ = strconv.Atoi(v)
		}
		return i, st[:len(st)-1]
	}
	return 0, st
}

// static vars
var svars [26]string

type paramsBuffer struct {
	out bytes.Buffer
	buf bytes.Buffer
}

// Start initializes the params buffer with the initial string data.
// It also locks the paramsBuffer.  The caller must call End() when
// finished.
func (pb *paramsBuffer) Start(s string) {
	pb.out.Reset()
	pb.buf.Reset()
	pb.buf.WriteString(s)
}

// End returns the final output from TParam, but it also releases the lock.
func (pb *paramsBuffer) End() string {
	s := pb.out.String()
	return s
}

// NextCh returns the next input character to the expander.
func (pb *paramsBuffer) NextCh() (byte, error) {
	return pb.buf.ReadByte()
}

// PutCh "emits" (rather schedules for output) a single byte character.
func (pb *paramsBuffer) PutCh(ch byte) {
	pb.out.WriteByte(ch)
}

// PutString schedules a string for output.
func (pb *paramsBuffer) PutString(s string) {
	pb.out.WriteString(s)
}

// TPuts emits the string to the writer, but expands inline padding
// indications (of the form $<[delay]> where [delay] is msec) to
// a suitable time (unless the terminfo string indicates this isn't needed
// by specifying npc - no padding).  All Terminfo based strings should be
// emitted using this function.
func (t *Terminfo) TPuts(w io.Writer, s string) {
	for {
		beg := strings.Index(s, "$<")
		if beg < 0 {
			// Most strings don't need padding, which is good news!
			_, _ = io.WriteString(w, s)
			return
		}
		_, _ = io.WriteString(w, s[:beg])
		s = s[beg+2:]
		end := strings.Index(s, ">")
		if end < 0 {
			// unterminated.. just emit bytes unadulterated
			_, _ = io.WriteString(w, "$<"+s)
			return
		}
		val := s[:end]
		s = s[end+1:]
		padus := 0
		unit := time.Millisecond
		dot := false
	loop:
		for i := range val {
			switch val[i] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				padus *= 10
				padus += int(val[i] - '0')
				if dot {
					unit /= 10
				}
			case '.':
				if !dot {
					dot = true
				} else {
					break loop
				}
			default:
				break loop
			}
		}
	}
}

var (
	dblock    sync.Mutex
	terminfos = make(map[string]*Terminfo)
)

// AddTerminfo can be called to register a new Terminfo entry.
func AddTerminfo(t *Terminfo) {
	dblock.Lock()

	terminfos[t.Name] = t
	for _, x := range t.Aliases {
		terminfos[x] = t
	}
	dblock.Unlock()
}

// LookupTerminfo attempts to find a definition for the named $TERM.
func LookupTerminfo(name string) (*Terminfo, error) {
	if name == "" {
		// else on windows: index out of bounds
		// on the name[0] reference below
		return nil, ErrTermNotFound
	}

	dblock.Lock()
	t := terminfos[name]
	dblock.Unlock()

	if t == nil {
		return nil, ErrTermNotFound
	}

	return t, nil
}

func TerminfoNames() []string {
	res := make([]string, 0, len(terminfos))
	for m := range terminfos {
		res = append(res, m)
	}
	return res
}
