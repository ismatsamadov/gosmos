package main

import (
	"errors"
	"testing"
)

// Payloads from the protocol's table: each frame's bytes after its length.
const (
	// mondayStatus is Dish 1's STATUS at 09:14 on Monday 27 March.
	mondayStatus = "\x83\x00\x00\x00\x00\x01\x02\x47\x95\x10\x4a\x02\x80\xff\x83\x00\x01\x00\x00\x00\x76\x62\x9d\xe0\xf1"
	// nakWindLock is the NAK refusing seq 121 for the wind lock.
	nakWindLock = "\x91\x00\x00\x00\x79\x03wind lock\x7e\x15\x64\xe0"
)

func TestDecodeFrame(t *testing.T) {
	good := []struct {
		name, payload string
		typ           byte
		seq           uint32
		body          string
	}{
		{"Monday's STATUS", mondayStatus, 0x83, 0, mondayStatus[5:21]},
		{"a NAK of seq 121", nakWindLock, 0x91, 121, "\x03wind lock"},
	}
	for _, g := range good {
		t.Run(g.name, func(t *testing.T) {
			f, err := callDecodeFrame(t, []byte(g.payload))
			if err != nil || f.Type != g.typ || f.Seq != g.seq || string(f.Body) != g.body {
				t.Errorf("decodeFrame(%s) = type %#02x seq %d body % x, %v; want type %#02x seq %d body % x, nil",
					g.name, f.Type, f.Seq, f.Body, err, g.typ, g.seq, g.body)
			}
		})
	}
	t.Run("every single bit of the STATUS flipped in turn", func(t *testing.T) {
		accepted, first := 0, -1
		var firstFrame Frame
		var firstErr error
		for bit := range 8 * len(mondayStatus) {
			damaged := []byte(mondayStatus)
			damaged[bit/8] ^= 1 << (7 - bit%8)
			f, err := callDecodeFrame(t, damaged)
			if !errors.Is(err, ErrBadCRC) || !empty(f) {
				accepted++
				if first < 0 {
					first, firstFrame, firstErr = bit, f, err
				}
			}
		}
		if accepted > 0 {
			t.Errorf("decodeFrame missed %d of %d damaged STATUS payloads; with bit %d flipped it returned type %#02x seq %d body % x, %v; want an empty Frame and ErrBadCRC",
				accepted, 8*len(mondayStatus), first, firstFrame.Type, firstFrame.Seq, firstFrame.Body, firstErr)
		}
	})
	short := []struct {
		name, payload, want string
	}{
		{"8 bytes, one too few", mondayStatus[:8], "bad length 8"},
		{"nothing at all", "", "bad length 0"},
	}
	for _, s := range short {
		t.Run(s.name, func(t *testing.T) {
			f, err := callDecodeFrame(t, []byte(s.payload))
			if !errors.Is(err, ErrBadLength) || err.Error() != s.want || !empty(f) {
				t.Errorf("decodeFrame(%d bytes) = type %#02x, %v; want an empty Frame and ErrBadLength as %q", len(s.payload), f.Type, err, s.want)
			}
		})
	}
}

// callDecodeFrame calls decodeFrame, turning a panic into a failure that says where it came from.
func callDecodeFrame(t *testing.T, payload []byte) (Frame, error) {
	t.Helper()
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("decodeFrame(%d bytes) panicked: %v; check the length before reading any field", len(payload), p)
		}
	}()
	return decodeFrame(payload)
}

// empty reports whether f is the zero Frame, with no field filled in.
func empty(f Frame) bool {
	return f.Type == 0 && f.Seq == 0 && f.Body == nil
}
