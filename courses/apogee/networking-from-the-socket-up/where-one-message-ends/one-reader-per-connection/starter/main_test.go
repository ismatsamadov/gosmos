package main

import (
	"io"
	"net"
	"testing"
	"time"
)

// Frames dishd sends, as they go on the wire.
const (
	helloDishd = "\x00\x0e\x01\x00\x00\x00\x00dishd\x97\x98\xf8\xab"
	query      = "\x00\x09\x03\x00\x00\x00\x00\x81\x82\x8d\xcd"
	ping       = "\x00\x09\x02\x00\x00\x00\x00\xbc\xe2\xa4\x7d"
)

func TestRecv(t *testing.T) {
	tests := []struct {
		name       string
		holdBanner bool
		split      int
	}{
		{"the banner, then WELCOME and STATUS glued", false, 0},
		{"banner, WELCOME and STATUS in one write", true, 0},
		{"the same, one byte at a time", true, 1},
	}
	for _, tt := range tests {
		ok := t.Run(tt.name, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.Wind = 81 // Wednesday's wind: the dish stows itself
			sim.Glue, sim.HoldBanner, sim.Split = 2, tt.holdBanner, tt.split
			conn := pipeToSim(t, sim)
			sendBoth := func() {
				if _, err := io.WriteString(conn, helloDishd+query); err != nil {
					t.Fatalf("sending HELLO and QUERY: %v", err)
				}
			}
			if tt.holdBanner {
				sendBoth()
			}
			c, err := newClient(conn)
			if err != nil || c.Name != "YST-1" {
				t.Fatalf("newClient(simulator) = %v, want the client of YST-1", err)
			}
			if !tt.holdBanner {
				sendBoth()
			}
			answers := []struct {
				typ  byte
				name string
			}{{TypeWelcome, "WELCOME"}, {TypeStatus, "STATUS"}}
			for i, want := range answers {
				f, err := c.recv()
				if err != nil || f.Type != want.typ {
					t.Fatalf("recv() #%d = type %#02x, %v; want the %s (type %#02x): does every read go through c.br?", i+1, f.Type, err, want.name, want.typ)
				}
				if want.typ == TypeStatus && (len(f.Body) != 16 || f.Body[1] != 3) {
					t.Errorf("the STATUS body is % x, want state 3, STOWED, at index 1", f.Body)
				}
			}
			if _, err := io.WriteString(conn, ping); err != nil {
				t.Fatalf("sending PING: %v", err)
			}
			if f, err := c.recv(); err != nil || f.Type != TypePong {
				t.Errorf("recv() after a PING = type %#02x, %v; want the PONG (type 0x82)", f.Type, err)
			}
		})
		if !ok {
			return // the cases after it would fail the same way
		}
	}
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
