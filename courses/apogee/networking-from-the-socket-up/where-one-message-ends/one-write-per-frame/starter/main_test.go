package main

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
)

func TestWriteFrame(t *testing.T) {
	t.Run("a PING", func(t *testing.T) {
		var r recorder
		err := writeFrame(&r, Frame{Type: TypePing})
		want := "00 09 02 00 00 00 00 bc e2 a4 7d"
		if err != nil || len(r.writes) != 1 || fmt.Sprintf("% x", r.writes[0]) != want {
			t.Errorf("writeFrame(PING) made %d Writes: %s, %v; want one Write, %s, nil", len(r.writes), r.String(), err, want)
		}
	})
	t.Run("HELLO, QUERY and PING", func(t *testing.T) {
		var r recorder
		frames := []Frame{{Type: TypeHello, Body: []byte("dishd")}, {Type: TypeQuery}, {Type: TypePing}}
		for _, f := range frames {
			if err := writeFrame(&r, f); err != nil {
				t.Fatalf("writeFrame(type %#02x) = %v, want nil", f.Type, err)
			}
		}
		if len(r.writes) != 3 {
			t.Fatalf("three writeFrame calls made %d Writes: %s; want 3, one per frame", len(r.writes), r.String())
		}
		for i, f := range frames {
			if want := encodeFrame(f); !bytes.Equal(r.writes[i], want) {
				t.Errorf("Write %d = % x, want the whole frame % x", i+1, r.writes[i], want)
			}
		}
	})
	t.Run("a line that has gone", func(t *testing.T) {
		if err := writeFrame(gone{}, Frame{Type: TypePing}); !errors.Is(err, net.ErrClosed) {
			t.Errorf("writeFrame(a closed line, PING) = %v, want the Write's error, net.ErrClosed", err)
		}
	})
}

// recorder keeps a copy of every Write it is given.
type recorder struct{ writes [][]byte }

// Write records a copy of p.
func (r *recorder) Write(p []byte) (int, error) {
	r.writes = append(r.writes, bytes.Clone(p))
	return len(p), nil
}

// String shows each Write's bytes in brackets.
func (r *recorder) String() string {
	var parts []string
	for _, w := range r.writes {
		parts = append(parts, fmt.Sprintf("[% x]", w))
	}
	return strings.Join(parts, " ")
}

// gone is a connection that has closed: every Write fails.
type gone struct{}

// Write fails as a closed connection's does.
func (gone) Write(p []byte) (int, error) { return 0, net.ErrClosed }
