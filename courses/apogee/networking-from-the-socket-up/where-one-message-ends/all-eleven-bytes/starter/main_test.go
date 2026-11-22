package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

// Frames from the protocol's table, as they go on the wire.
const (
	helloDishd = "\x00\x0e\x01\x00\x00\x00\x00dishd\x97\x98\xf8\xab"
	welcome    = "\x00\x10\x81\x00\x00\x00\x00\x01\x02\x04\x00\x00\x00\x00\x9e\x88\xe9\x9f"
	ping       = "\x00\x09\x02\x00\x00\x00\x00\xbc\xe2\xa4\x7d"
	pong       = "\x00\x09\x82\x00\x00\x00\x00\x0d\x04\x14\xef"
)

func TestReadExactly(t *testing.T) {
	tests := []struct {
		name string
		r    io.Reader
		want string
		err  error
	}{
		{"a PONG, one byte at a time", iotest.OneByteReader(strings.NewReader(pong)), pong, nil},
		{"a PONG, half of the rest at a time", iotest.HalfReader(strings.NewReader(pong)), pong, nil},
		{"a stream that ends after 5 of 11 bytes", strings.NewReader(pong[:5]), "", io.ErrUnexpectedEOF},
		{"an empty stream", strings.NewReader(""), "", io.EOF},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readExactly(tt.r, 11)
			if tt.err == nil && (!bytes.Equal(got, []byte(tt.want)) || err != nil) {
				t.Errorf("readExactly(%s, 11) = % x, %v; want % x, nil", tt.name, got, err, tt.want)
			}
			if tt.err != nil && !errors.Is(err, tt.err) {
				t.Errorf("readExactly(%s, 11) error = %v, want %v", tt.name, err, tt.err)
			}
		})
	}
	t.Run("the simulator writing 3 bytes at a time", func(t *testing.T) {
		sim := NewSim("YST-1", 1)
		sim.Split = 3
		conn := pipeToSim(t, sim)
		steps := []struct {
			send string // what dishd sends before reading, if anything
			n    int
			want string
		}{
			{"", 21, "YST-1 READY PROTO 1\r\n"},
			{helloDishd, 18, welcome},
			{ping, 11, pong},
		}
		for _, s := range steps {
			if s.send != "" {
				if _, err := io.WriteString(conn, s.send); err != nil {
					t.Fatalf("sending % x: %v", s.send, err)
				}
			}
			got, err := readExactly(conn, s.n)
			if !bytes.Equal(got, []byte(s.want)) || err != nil {
				t.Fatalf("readExactly(simulator, %d) = % x, %v; want % x, nil", s.n, got, err, s.want)
			}
		}
	})
}

// pipeToSim serves sim on one end of a net.Pipe and returns the other end, with a deadline
// a second away. The test's end closes, and the simulator stops, when the test ends.
func pipeToSim(t *testing.T, sim *Sim) net.Conn {
	t.Helper()
	conn, end := net.Pipe()
	served := make(chan struct{})
	go func() {
		sim.ServeConn(end)
		close(served)
	}()
	conn.SetDeadline(time.Now().Add(time.Second))
	t.Cleanup(func() {
		conn.Close()
		select {
		case <-served:
		case <-time.After(time.Second):
			t.Error("the simulator did not stop within a second of its connection closing")
		}
	})
	return conn
}
