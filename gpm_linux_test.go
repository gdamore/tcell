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

//go:build linux

package tcell

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func gpmTestPacket(buttons, mods byte, x, y int, typ int, wdx, wdy int) []byte {
	b := make([]byte, gpmEventSize)
	ne := binary.NativeEndian
	b[0] = buttons
	b[1] = mods
	ne.PutUint16(b[8:], uint16(int16(x)))
	ne.PutUint16(b[10:], uint16(int16(y)))
	ne.PutUint32(b[12:], uint32(typ))
	ne.PutUint16(b[24:], uint16(int16(wdx)))
	ne.PutUint16(b[26:], uint16(int16(wdy)))
	return b
}

func TestParseGPMConsoleTTY(t *testing.T) {
	tests := []struct {
		link string
		want uint32
		ok   bool
	}{
		{link: "/dev/tty1", want: 1, ok: true},
		{link: "/dev/tty63", want: 63, ok: true},
		{link: "/dev/vc/4", want: 4, ok: true},
		{link: "/dev/pts/2", ok: false},
		{link: "/dev/tty", ok: false},
		{link: "/dev/tty0", ok: false},
		{link: "/tmp/tty7", ok: false},
		{link: "/dev/tty2147483648", ok: false},
	}
	for _, tt := range tests {
		got, ok := parseGPMConsoleTTY(tt.link)
		if got != tt.want || ok != tt.ok {
			t.Fatalf("parseGPMConsoleTTY(%q) = %d, %t; want %d, %t", tt.link, got, ok, tt.want, tt.ok)
		}
	}
}

func TestNewGPMConsoleMouseWithMockServer(t *testing.T) {
	oldCtlPath := gpmCtlPath
	oldFDLink := gpmFDLink
	defer func() {
		gpmCtlPath = oldCtlPath
		gpmFDLink = oldFDLink
	}()

	gpmFDLink = func(fd string) (string, error) {
		if fd == "2" {
			return "/dev/tty7", nil
		}
		return "", os.ErrNotExist
	}

	gpmCtlPath = filepath.Join(t.TempDir(), "gpmctl")
	listener, err := net.Listen("unix", gpmCtlPath)
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer listener.Close()

	serverErr := make(chan error, 1)
	connectPacket := make(chan []byte, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			serverErr <- err
			return
		}
		defer conn.Close()

		packet := make([]byte, gpmConnectSize)
		if _, err = io.ReadFull(conn, packet); err != nil {
			serverErr <- err
			return
		}
		connectPacket <- packet

		if _, err = conn.Write(gpmTestPacket(0, 0, 5, 6, gpmEventMove, 0, 0)); err != nil {
			if !errors.Is(err, net.ErrClosed) {
				serverErr <- err
			}
			return
		}
		serverErr <- nil
	}()

	eventQ := make(chan Event, 1)
	console, err := newGPMConsoleMouse(eventQ)
	if err != nil {
		t.Fatalf("newGPMConsoleMouse failed: %v", err)
	}
	cm := console.(*gpmConsoleMouse)
	cm.SetSize(80, 25)

	select {
	case packet := <-connectPacket:
		ne := binary.NativeEndian
		if ne.Uint16(packet[0:]) != 0xffff || ne.Uint16(packet[2:]) != 0xffff {
			t.Fatalf("unexpected masks in connect packet: % x", packet)
		}
		if vc := ne.Uint32(packet[12:]); vc != 7 {
			t.Fatalf("connect packet vc = %d, want 7", vc)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for GPM connect packet")
	}

	if err = cm.Start(MouseMotionEvents); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer cm.Stop()

	select {
	case ev := <-eventQ:
		mouse, ok := ev.(*EventMouse)
		if !ok {
			t.Fatalf("event type = %T, want *EventMouse", ev)
		}
		if mouse.Buttons() != ButtonNone {
			t.Fatalf("buttons = %x, want none", mouse.Buttons())
		}
		if x, y := mouse.Position(); x != 4 || y != 5 {
			t.Fatalf("position = %d,%d, want 4,5", x, y)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for mock GPM event")
	}

	select {
	case err := <-serverErr:
		if err != nil {
			t.Fatalf("mock server failed: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for mock server")
	}
}

func TestGPMButtonAndModifierMapping(t *testing.T) {
	if got := gpmButtonMask(gpmButtonLeft|gpmButtonRight|gpmButtonMiddle|gpmButtonFourth, 0, 0); got != Button1|Button2|Button3|Button4 {
		t.Fatalf("button mapping = %x", got)
	}
	if got := gpmButtonMask(gpmButtonUp|gpmButtonDown, -1, 1); got != WheelUp|WheelDown|WheelLeft {
		t.Fatalf("wheel mapping = %x", got)
	}
	if got := gpmButtonMask(0, 1, -1); got != WheelRight|WheelDown {
		t.Fatalf("wheel delta mapping = %x", got)
	}
	if got := gpmModMask(gpmModShift | gpmModAltGr | gpmModCtrl | gpmModAlt); got != ModShift|ModAlt|ModCtrl {
		t.Fatalf("modifier mapping = %x", got)
	}
}

func TestGPMEventFromPacketButtonsAndClipping(t *testing.T) {
	cm := &gpmConsoleMouse{eventQ: make(chan Event, 1), flags: MouseButtonEvents}
	cm.SetSize(10, 5)

	ev, ok := cm.eventFromPacket(gpmTestPacket(gpmButtonLeft, gpmModShift, 3, 4, gpmEventDown|gpmEventSingle, 0, 0))
	if !ok {
		t.Fatal("button down was ignored")
	}
	if ev.Buttons() != Button1 || ev.Modifiers() != ModShift {
		t.Fatalf("down event = buttons %x modifiers %x", ev.Buttons(), ev.Modifiers())
	}
	if x, y := ev.Position(); x != 2 || y != 3 {
		t.Fatalf("down position = %d,%d", x, y)
	}

	ev, ok = cm.eventFromPacket(gpmTestPacket(gpmButtonRight, 0, 99, -4, gpmEventDown, 0, 0))
	if !ok {
		t.Fatal("button chord was ignored")
	}
	if ev.Buttons() != Button1|Button2 {
		t.Fatalf("chord buttons = %x", ev.Buttons())
	}
	if x, y := ev.Position(); x != 9 || y != 0 {
		t.Fatalf("clipped position = %d,%d", x, y)
	}

	ev, ok = cm.eventFromPacket(gpmTestPacket(gpmButtonLeft, 0, 3, 4, gpmEventUp, 0, 0))
	if !ok {
		t.Fatal("button up was ignored")
	}
	if ev.Buttons() != Button2 {
		t.Fatalf("buttons after release = %x", ev.Buttons())
	}
}

func TestGPMEventFromPacketRespectsMouseFlags(t *testing.T) {
	cm := &gpmConsoleMouse{eventQ: make(chan Event, 1), flags: MouseButtonEvents}
	cm.SetSize(10, 5)
	if _, ok := cm.eventFromPacket(gpmTestPacket(0, 0, 3, 4, gpmEventMove, 0, 0)); ok {
		t.Fatal("motion event passed with MouseButtonEvents only")
	}
	if _, ok := cm.eventFromPacket(gpmTestPacket(gpmButtonLeft, 0, 3, 4, gpmEventDrag, 0, 0)); ok {
		t.Fatal("drag event passed with MouseButtonEvents only")
	}

	cm.flags = MouseDragEvents
	if ev, ok := cm.eventFromPacket(gpmTestPacket(gpmButtonLeft, 0, 3, 4, gpmEventDrag, 0, 0)); !ok || ev.Buttons() != Button1 {
		t.Fatalf("drag event = %#v ok %t", ev, ok)
	}

	cm.flags = MouseMotionEvents
	if ev, ok := cm.eventFromPacket(gpmTestPacket(0, 0, 3, 4, gpmEventMove, 0, 0)); !ok || ev.Buttons() != ButtonNone {
		t.Fatalf("motion event = %#v ok %t", ev, ok)
	}
}

func TestGPMEventFromPacketWheel(t *testing.T) {
	cm := &gpmConsoleMouse{eventQ: make(chan Event, 1), flags: MouseButtonEvents}
	cm.SetSize(10, 5)

	ev, ok := cm.eventFromPacket(gpmTestPacket(gpmButtonUp, gpmModCtrl, 3, 4, gpmEventMove, 0, 0))
	if !ok {
		t.Fatal("wheel up was ignored")
	}
	if ev.Buttons() != WheelUp || ev.Modifiers() != ModCtrl {
		t.Fatalf("wheel up event = buttons %x modifiers %x", ev.Buttons(), ev.Modifiers())
	}

	ev, ok = cm.eventFromPacket(gpmTestPacket(0, 0, 3, 4, 0, 1, -1))
	if !ok {
		t.Fatal("wheel deltas were ignored")
	}
	if ev.Buttons() != WheelRight|WheelDown {
		t.Fatalf("wheel delta event = buttons %x", ev.Buttons())
	}
}

func TestGPMEventFromShortPacketIgnored(t *testing.T) {
	cm := &gpmConsoleMouse{eventQ: make(chan Event, 1), flags: MouseButtonEvents}
	if _, ok := cm.eventFromPacket(make([]byte, gpmEventSize-1)); ok {
		t.Fatal("short packet produced event")
	}
}

func TestGPMStartStopWithPartialFrame(t *testing.T) {
	reader, writer := net.Pipe()
	cm := &gpmConsoleMouse{sock: reader, eventQ: make(chan Event, 1)}
	if err := cm.Start(MouseButtonEvents); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if _, err := writer.Write(make([]byte, gpmEventSize/2)); err != nil {
		t.Fatalf("partial write failed: %v", err)
	}
	_ = writer.Close()

	select {
	case <-cm.done:
	case <-time.After(time.Second):
		t.Fatal("GPM loop did not exit after partial frame EOF")
	}
	cm.Stop()
}
