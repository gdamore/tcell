//go:build wasip1 || wasip2

// Copyright 2026 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

import "errors"

// initialize on WASI: there is no /dev/tty and no termios control, so a
// terminal screen can only work over an explicitly provided Tty
// implementation; without one, initialization fails.
func (t *tScreen) initialize() error {
	if t.tty == nil {
		return errors.New("tcell: no tty available")
	}
	return nil
}
