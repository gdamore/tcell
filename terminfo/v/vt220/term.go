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

package vt220

import "github.com/gdamore/tcell/v3/terminfo"

func init() {

	// DEC VT220 and later (through VT420)
	terminfo.AddTerminfo(&terminfo.Terminfo{
		Name:    "vt220",
		Aliases: []string{"vt200", "vt320", "vt400", "vt400-24", "dec-vt400", "vt420"},
		Columns: 80,
		Lines:   24,
	})
}
