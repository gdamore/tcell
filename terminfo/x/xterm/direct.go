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

// This terminal definition is derived from the xterm-256color definition, but
// makes use of the RGB property these terminals have to support direct color.
// The terminfo entry for this uses a new format for the color handling introduced
// by ncurses 6.1 (and used by nobody else), so this override ensures we get
// good handling even in the face of this.

package xterm

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// derived from xterm-256color, but adds full RGB support
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:    "xterm-direct",
		Aliases: []string{"alacritty-direct", "xterm-truecolor"},
		Columns: 80,
		Lines:   24,
	})
}
