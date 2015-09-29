// +build !windows,!nacl,!plan9

// Copyright 2015 The TCell Authors
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
	"os/signal"
	"syscall"
)

// #include <termios.h>
// #include <sys/ioctl.h>
//
// int getwinsize(int fd, int *cols, int *rows) {
// #if defined TIOCGWINSZ
//	struct winsize w;
//	if (ioctl(fd, TIOCGWINSZ, &w) < 0) {
//		return (-1);
//	}
//	*cols = w.ws_col;
//	*rows = w.ws_row;
//	return (0);
// #else
//	return (-1);
// #endif
// }
import "C"

var savedtios map[*tScreen]*C.struct_termios

func init() {
	savedtios = make(map[*tScreen]*C.struct_termios)
}

func (t *tScreen) termioInit() error {
	var e error
	var rv C.int
	var oldtios C.struct_termios
	var newtios C.struct_termios
	var fd C.int

	if t.in, e = os.OpenFile("/dev/tty", os.O_RDONLY, 0); e != nil {
		goto failed
	}
	if t.out, e = os.OpenFile("/dev/tty", os.O_WRONLY, 0); e != nil {
		goto failed
	}

	fd = C.int(t.out.Fd())
	if rv, e = C.tcgetattr(fd, &oldtios); rv != 0 {
		goto failed
	}
	newtios = oldtios
	newtios.c_iflag &^= C.IGNBRK | C.BRKINT | C.PARMRK |
		C.ISTRIP | C.INLCR | C.IGNCR |
		C.ICRNL | C.IXON
	newtios.c_oflag &^= C.OPOST
	newtios.c_lflag &^= C.ECHO | C.ECHONL | C.ICANON |
		C.ISIG | C.IEXTEN
	newtios.c_cflag &^= C.CSIZE | C.PARENB
	newtios.c_cflag |= C.CS8

	// We wake up at the earliest of 100 msec or when data is received.
	// We need to wake up frequently to permit us to exit cleanly and
	// close file descriptors on systems like Darwin, where close does
	// cause a wakeup.  (Probably we could reasonably increase this to
	// something like 1 sec or 500 msec.)
	newtios.c_cc[C.VMIN] = 0
	newtios.c_cc[C.VTIME] = 1

	if rv, e = C.tcsetattr(fd, C.TCSANOW|C.TCSAFLUSH, &newtios); rv != 0 {
		goto failed
	}

	savedtios[t] = &oldtios
	signal.Notify(t.sigwinch, syscall.SIGWINCH)

	if w, h, e := t.getWinSize(); e == nil && w != 0 && h != 0 {
		t.w = w
		t.h = h
	}

	return nil

failed:
	if t.in != nil {
		t.in.Close()
	}
	if t.out != nil {
		t.out.Close()
	}
	return e
}

func (t *tScreen) termioFini() {

	signal.Stop(t.sigwinch)

	if t.out != nil {
		if oldtios, ok := savedtios[t]; ok {
			fd := C.int(t.out.Fd())
			C.tcsetattr(fd, C.TCSANOW, oldtios)
			delete(savedtios, t)
		}
		t.out.Close()
	}
	if t.in != nil {
		t.in.Close()
	}
}

func (t *tScreen) getWinSize() (int, int, error) {
	var cx, cy C.int
	if r, e := C.getwinsize(C.int(t.out.Fd()), &cx, &cy); r == 0 {
		return int(cx), int(cy), nil
	} else {
		return 0, 0, e
	}
}
