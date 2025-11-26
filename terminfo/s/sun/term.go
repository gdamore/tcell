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

// This terminal definition is hand-coded, as the default terminfo for
// this terminal is busted with respect to color.  Unlike pretty much every
// other ANSI compliant terminal, this terminal cannot combine foreground and
// background escapes.  The default terminfo also only provides escapes for
// 16-bit color. We also added support for disabling auto margins, which was
// added to illumos back in 2021.

package sun

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// Sun Microsystems Inc. workstation console
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:    "sun",
		Columns: 80,
		Lines:   34,
	})

	// Sun Microsystems Workstation console with color support (IA systems)
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:      "sun-color",
		Columns:   80,
		Lines:     34,
		Colors:    256,
		SetFg:     "\x1b[38;5;%p1%dm",
		SetBg:     "\x1b[48;5;%p1%dm",
		ResetFgBg: "\x1b[0m",
	})
}
