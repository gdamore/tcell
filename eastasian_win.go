//go:build windows
// +build windows

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

package tcell

import (
	"strings"
	"syscall"
	"unsafe"

	"github.com/rivo/uniseg"
	"golang.org/x/sys/windows"
)

var (
	kernel32            = syscall.NewLazyDLL("kernel32")
	procGetLocaleInfoEx = kernel32.NewProc("GetUserDefaultLocaleName")
)

func init() {

	localeName := make([]uint16, 85) // Windows locale names limited to 85

	r1, _, _ := procGetLocaleInfoEx.Call(uintptr(unsafe.Pointer(&localeName[0])), uintptr(len(localeName)))
	name := windows.UTF16ToString(localeName)

	if r1 == 0 {
		uniseg.EastAsianAmbiguousWidth = 1
	} else if strings.HasPrefix(name, "ko") || strings.HasPrefix(name, "ja") || strings.HasPrefix(name, "cn") {
		if strings.Contains(name, "narrow") || strings.Contains(name, "half") {
			uniseg.EastAsianAmbiguousWidth = 1
		} else {
			uniseg.EastAsianAmbiguousWidth = 2
		}
	} else {
		uniseg.EastAsianAmbiguousWidth = 1
	}
}
