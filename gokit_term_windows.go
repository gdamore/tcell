// based on https://github.com/go-kit -- modified to indicate which pty we're on

// // Based on ssh/terminal:
// // Copyright 2011 The Go Authors. All rights reserved.
// // Use of this source code is governed by a BSD-style
// // license that can be found in the LICENSE file.

// +build windows

package tcell

import (
	"encoding/binary"
	"io"
	"regexp"
	"strconv"
	"syscall"
	"unsafe"
)

type fder interface {
	Fd() uintptr
}

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procGetFileInformationByHandleEx = kernel32.NewProc("GetFileInformationByHandleEx")

	// Originally this was \d? but that doesn't support pty10, etc.
	// moreover, pty0 appears (at least in my testing) rather than just pty
	// TODO - notify origin they should change regex to support >= pty10
	msysPipeNameRegex = regexp.MustCompile(`\\(cygwin|msys)-\w+-pty(\d+)-(to|from)-master`)
)

const (
	fileNameInfo = 0x02
)

// GetMSYSTerminal returns the pty# if w writes to a MSYS/MSYS2 terminal; otherwise, -1
func GetMSYSTerminal(w io.Writer) int {
	var handle syscall.Handle

	if fw, ok := w.(fder); ok {
		handle = syscall.Handle(fw.Fd())
	} else {
		// The writer has no file-descriptor and so can't be a terminal.
		return -1
	}

	// MSYS(2) terminal reports as a pipe for STDIN/STDOUT/STDERR. If it isn't
	// a pipe, it can't be a MSYS(2) terminal.
	filetype, err := syscall.GetFileType(handle)

	if filetype != syscall.FILE_TYPE_PIPE || err != nil {
		return -1
	}

	// MSYS2/Cygwin terminal's name looks like: \msys-dd50a72ab4668b33-pty2-to-master
	data := make([]byte, 256, 256)

	r, _, e := syscall.Syscall6(
		procGetFileInformationByHandleEx.Addr(),
		4,
		uintptr(handle),
		uintptr(fileNameInfo),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		0,
		0,
	)

	if r != 0 && e == 0 {
		// The first 4 bytes of the buffer are the size of the UTF16 name, in bytes.
		unameLen := binary.LittleEndian.Uint32(data[:4]) / 2
		uname := make([]uint16, unameLen, unameLen)

		for i := uint32(0); i < unameLen; i++ {
			uname[i] = binary.LittleEndian.Uint16(data[i*2+4 : i*2+2+4])
		}

		name := syscall.UTF16ToString(uname)

		found := msysPipeNameRegex.FindStringSubmatch(name)
		if len(found) != 4 {
			return -1
		}

		ptynum, err := strconv.Atoi(found[2])
		if err != nil {
			return -1
		}

		return ptynum

	}

	return -1
}
