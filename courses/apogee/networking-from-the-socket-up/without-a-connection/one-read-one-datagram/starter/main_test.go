package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"
)

// Frames from Monday's passes, as the recorder sends them: 70, 72 and 72 bytes.
const (
	pipit2Frame  = "pipit-2 30412 2028-06-12T06:37:15Z batt_v=8.04 temp_c=9.5 mode=NOMINAL"
	furrowAFrame = "furrow-a 52817 2028-06-12T07:38:02Z batt_v=7.86 temp_c=12.5 mode=NOMINAL"
	nextFrame    = "furrow-a 52818 2028-06-12T07:38:03Z batt_v=7.86 temp_c=12.5 mode=NOMINAL"
)

func TestReadDatagram(t *testing.T) {
	t.Run("three frames, three datagrams", func(t *testing.T) {
		view, rec := sockets(t)
		frames := []string{pipit2Frame, furrowAFrame, nextFrame}
		for _, f := range frames {
			send(t, rec, view, []byte(f))
		}
		var got [][]byte
		for range frames {
			d, from, err := readDatagram(view)
			if err != nil {
				t.Fatalf("readDatagram() after three frames = %v, want each frame and nil", err)
			}
			if from == nil || from.String() != rec.LocalAddr().String() {
				t.Errorf("readDatagram() sender = %v, want the recorder's %v", from, rec.LocalAddr())
			}
			got = append(got, d)
		}
		for i, f := range frames {
			if string(got[i]) != f {
				t.Errorf("readDatagram() number %d = %q (%d bytes), want %q (%d bytes)", i+1, got[i], len(got[i]), f, len(f))
			}
		}
	})
	t.Run("512 bytes, the most a datagram may hold", func(t *testing.T) {
		view, rec := sockets(t)
		want := bytes.Repeat([]byte("x"), 512)
		send(t, rec, view, want)
		d, _, err := readDatagram(view)
		if err != nil || !bytes.Equal(d, want) {
			t.Errorf("readDatagram() of 512 bytes = %d bytes, %v; want all 512 and nil", len(d), err)
		}
	})
	for _, size := range []int{513, 600} {
		t.Run(fmt.Sprintf("%d bytes, then a frame", size), func(t *testing.T) {
			view, rec := sockets(t)
			send(t, rec, view, bytes.Repeat([]byte("y"), size))
			send(t, rec, view, []byte(nextFrame))
			d, from, err := readDatagram(view)
			if !errors.Is(err, ErrTooBig) {
				t.Fatalf("readDatagram() of %d bytes = %d bytes, %v; want nil, ErrTooBig", size, len(d), err)
			}
			if d != nil {
				t.Errorf("readDatagram() of %d bytes returned %d bytes with ErrTooBig, want nil", size, len(d))
			}
			if from == nil || from.String() != rec.LocalAddr().String() {
				t.Errorf("readDatagram() of %d bytes: sender = %v, want the recorder's %v", size, from, rec.LocalAddr())
			}
			d, _, err = readDatagram(view)
			if err != nil || string(d) != nextFrame {
				t.Errorf("readDatagram() after the %d-byte datagram = %q, %v; want the next frame whole, %q, nil: the rest of a cut datagram is gone", size, d, err, nextFrame)
			}
		})
	}
	t.Run("an empty datagram", func(t *testing.T) {
		view, rec := sockets(t)
		send(t, rec, view, nil)
		d, from, err := readDatagram(view)
		if err != nil || len(d) != 0 || from == nil {
			t.Errorf("readDatagram() of an empty datagram = %q, %v, %v; want 0 bytes, the recorder's address, nil", d, from, err)
		}
	})
}

// sockets opens the live view's socket and the recorder's on loopback, closes both when the
// test ends, and limits the live view's reads to a second, so a wrong read can't hang the test.
func sockets(t *testing.T) (view, rec net.PacketConn) {
	t.Helper()
	view, rec = listen(t), listen(t)
	view.SetReadDeadline(time.Now().Add(time.Second))
	return view, rec
}

// listen opens a UDP socket on 127.0.0.1, on a port the system picks, until the test ends.
func listen(t *testing.T) net.PacketConn {
	t.Helper()
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { pc.Close() })
	return pc
}

// send sends d as one datagram from one socket to another.
func send(t *testing.T, from, to net.PacketConn, d []byte) {
	t.Helper()
	if _, err := from.WriteTo(d, to.LocalAddr()); err != nil {
		t.Fatalf("sending %d bytes: %v", len(d), err)
	}
}
