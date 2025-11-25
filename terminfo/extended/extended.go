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

// Package extended contains an extended set of terminal descriptions.
// Applications desiring to have a better chance of Just Working by
// default should include this package.  This will significantly increase
// the size of the program.
package extended

import (
	// The following imports just register themselves --
	// these are the terminal types we aggregate in this package.
	_ "github.com/gdamore/tcell/v3/terminfo/a/aixterm"
	_ "github.com/gdamore/tcell/v3/terminfo/a/ansi"
	_ "github.com/gdamore/tcell/v3/terminfo/d/dtterm"
	_ "github.com/gdamore/tcell/v3/terminfo/e/emacs"
	_ "github.com/gdamore/tcell/v3/terminfo/l/linux"
	_ "github.com/gdamore/tcell/v3/terminfo/r/rxvt"
	_ "github.com/gdamore/tcell/v3/terminfo/s/screen"
	_ "github.com/gdamore/tcell/v3/terminfo/s/simpleterm"
	_ "github.com/gdamore/tcell/v3/terminfo/s/sun"
	_ "github.com/gdamore/tcell/v3/terminfo/t/tmux"
	_ "github.com/gdamore/tcell/v3/terminfo/v/vt100"
	_ "github.com/gdamore/tcell/v3/terminfo/v/vt102"
	_ "github.com/gdamore/tcell/v3/terminfo/v/vt220"
	_ "github.com/gdamore/tcell/v3/terminfo/x/xterm"
)
