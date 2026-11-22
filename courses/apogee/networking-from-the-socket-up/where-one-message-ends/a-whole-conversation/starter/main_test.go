package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestExchange(t *testing.T) {
	deadline := time.Now().Add(time.Second) // every read in this test ends by then
	hello := Frame{Type: TypeHello, Body: []byte("dishd")}
	query := Frame{Type: TypeQuery}
	ping := Frame{Type: TypePing}
	point := Frame{Type: TypePoint, Seq: 119, Body: []byte{0x47, 0x95, 0x10, 0x4a}}
	tests := []struct {
		name    string
		sim     func(*Sim)
		frames  []Frame
		want    []byte // the answers' types, in order
		wantErr func(error) bool
		errName string
	}{
		{"the simulator writing one byte at a time", func(s *Sim) { s.Split = 1 },
			[]Frame{hello, query, ping}, []byte{TypeWelcome, TypeStatus, TypePong}, isNil, "nil"},
		{"all three answers glued in one write", func(s *Sim) { s.Glue = 3 },
			[]Frame{hello, query, ping}, []byte{TypeWelcome, TypeStatus, TypePong}, isNil, "nil"},
		{"a BYE idle in place of the second answer", func(s *Sim) { s.MuteAfter, s.IdleLimit = 1, 50*time.Millisecond },
			[]Frame{hello, query, ping}, []byte{TypeWelcome}, isIdleBye, `a *ByeError reading "controller closed: idle"`},
		{"the controller hangs up after two answers", func(s *Sim) { s.DropAfter = 119 },
			[]Frame{hello, query, point}, []byte{TypeWelcome, TypeStatus}, isEOF, "io.EOF"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			tt.sim(sim)
			c, err := newClient(pipeToSim(t, sim, deadline))
			if err != nil {
				t.Fatalf("newClient(simulator) = %v, want a client", err)
			}
			got, err := c.exchange(tt.frames...)
			if names(types(got)) != names(tt.want) || !tt.wantErr(err) {
				t.Errorf("exchange(%s) = %s, %v; want %s, %s", names(types(tt.frames)), names(types(got)), err, names(tt.want), tt.errName)
			}
		})
	}
}

// isNil reports whether err is nil.
func isNil(err error) bool { return err == nil }

// isIdleBye reports whether err is the controller hanging up for being idle.
func isIdleBye(err error) bool {
	bye, ok := errors.AsType[*ByeError](err)
	return ok && bye.Error() == "controller closed: idle"
}

// isEOF reports whether err is the end of the stream.
func isEOF(err error) bool { return errors.Is(err, io.EOF) }

// types lists the frames' types.
func types(frames []Frame) []byte {
	var ts []byte
	for _, f := range frames {
		ts = append(ts, f.Type)
	}
	return ts
}

// names writes frame types as the protocol names them: [WELCOME STATUS PONG].
func names(ts []byte) string {
	known := map[byte]string{TypeHello: "HELLO", TypeQuery: "QUERY", TypePing: "PING", TypePoint: "POINT",
		TypeWelcome: "WELCOME", TypeStatus: "STATUS", TypePong: "PONG", TypeAck: "ACK", TypeNak: "NAK", TypeBye: "BYE"}
	var ns []string
	for _, t := range ts {
		n, ok := known[t]
		if !ok {
			n = fmt.Sprintf("%#02x", t)
		}
		ns = append(ns, n)
	}
	return "[" + strings.Join(ns, " ") + "]"
}

// pipeToSim serves sim on one end of a net.Pipe and returns the other end, with the given
// deadline. The test's end closes, and the simulator stops, when the test ends.
func pipeToSim(t *testing.T, sim *Sim, deadline time.Time) net.Conn {
	t.Helper()
	conn, end := net.Pipe()
	served := make(chan struct{})
	go func() {
		sim.ServeConn(end)
		close(served)
	}()
	conn.SetDeadline(deadline)
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
