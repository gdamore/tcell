//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

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

package tcell

import (
	"os"
	"strings"

	"github.com/rivo/uniseg"
)

func init() {

	name := os.Getenv("LANG")
	if name == "" {
		name = os.Getenv("LC_ALL")
	}

	if strings.HasPrefix(name, "ko") || strings.HasPrefix(name, "ja") || strings.HasPrefix(name, "cn") {
		if strings.Contains(name, "narrow") || strings.Contains(name, "half") {
			uniseg.EastAsianAmbiguousWidth = 1
		} else {
			uniseg.EastAsianAmbiguousWidth = 2
		}
	} else {
		uniseg.EastAsianAmbiguousWidth = 1
	}
}
