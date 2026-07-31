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

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package tty

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestFileDescriptorPreservesDeadline(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	if _, err = fileDescriptor(r); err != nil {
		t.Fatal(err)
	}
	if err = r.SetReadDeadline(time.Now().Add(10 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	errQ := make(chan error, 1)
	go func() {
		_, err := r.Read(make([]byte, 1))
		errQ <- err
	}()
	select {
	case err = <-errQ:
		if !errors.Is(err, os.ErrDeadlineExceeded) {
			t.Fatalf("got %v, want deadline exceeded", err)
		}
	case <-time.After(time.Second):
		t.Fatal("read did not honor deadline")
	}
}
