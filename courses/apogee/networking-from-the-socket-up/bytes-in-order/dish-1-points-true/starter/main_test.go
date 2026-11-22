package main

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"net"
	"testing"
	"time"
)

func TestCommand(t *testing.T) {
	body, err := pointBody(183.25, 41.70)
	if err != nil {
		t.Fatalf("pointBody(183.25, 41.70) = %v", err)
	}
	point := Frame{Type: TypePoint, Seq: 119, Body: body}
	t.Run("POINT 119 on Dish 1", func(t *testing.T) {
		sim := NewSim("YST-1", 1)
		sim.LastSeq = 118
		c := connect(t, sim)
		if err := c.command(point); err != nil {
			t.Fatalf("command(POINT 119) = %v, want nil", err)
		}
		st, err := c.status()
		if err != nil || math.Abs(st.AzDeg-183.25) > 0.001 || math.Abs(st.ElDeg-41.70) > 0.001 {
			t.Errorf("the status after POINT 119 = az %.2f el %.2f, %v; want az 183.25 el 41.70: read the POINT's answer before anything else", st.AzDeg, st.ElDeg, err)
		}
	})
	refusals := []struct {
		name string
		set  func(*Sim)
		code uint8
		want string
	}{
		{"a dish already stowed", func(s *Sim) { s.State = 3 }, 4, "seq 119 refused: stowed"},
		{"wind at 81 km/h", func(s *Sim) { s.Wind = 81 }, 3, "seq 119 refused: wind lock"},
	}
	for _, r := range refusals {
		t.Run(r.name, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.LastSeq = 118
			r.set(sim)
			err := connect(t, sim).command(point)
			nak, ok := errors.AsType[*NakError](err)
			if !ok || nak.Code != r.code || err.Error() != r.want {
				t.Errorf("command(POINT 119) = %v; want a *NakError with code %d, %q", err, r.code, r.want)
			}
		})
	}
	t.Run("an ACK for seq 118", func(t *testing.T) {
		err := fakeController(t, encodeFrame(Frame{Type: TypeAck, Seq: 118})).command(point)
		want := "answer to seq 119: type 0x90, seq 118"
		if err == nil || err.Error() != want {
			t.Errorf("command(POINT 119) answered by ACK 118 = %v, want %q: that ACK is some other command's", err, want)
		}
	})
	t.Run("a damaged answer", func(t *testing.T) {
		ack := encodeFrame(Frame{Type: TypeAck, Seq: 119})
		ack[len(ack)-1] ^= 0x01
		if err := fakeController(t, ack).command(point); !errors.Is(err, ErrBadCRC) {
			t.Errorf("command(POINT 119) answered by a damaged ACK = %v, want ErrBadCRC", err)
		}
	})
	t.Run("the controller hangs up", func(t *testing.T) {
		bye := encodeFrame(Frame{Type: TypeBye, Body: []byte("replaced")})
		err := fakeController(t, bye).command(point)
		if b, ok := errors.AsType[*ByeError](err); !ok || b.Reason != "replaced" {
			t.Errorf("command(POINT 119) answered by BYE replaced = %v, want the *ByeError recv returns", err)
		}
	})
}

// connect serves sim on a net.Pipe, reads its banner and says HELLO, and returns the client.
func connect(t *testing.T, sim *Sim) *Client {
	t.Helper()
	conn, end := net.Pipe()
	served := make(chan struct{})
	go func() {
		sim.ServeConn(end)
		close(served)
	}()
	t.Cleanup(func() {
		conn.Close()
		<-served
	})
	conn.SetDeadline(time.Now().Add(time.Second))
	c, err := newClient(conn)
	if err != nil {
		t.Fatalf("newClient(simulator) = %v, want a client", err)
	}
	if _, err := c.hello(); err != nil {
		t.Fatalf("hello() = %v, want the WELCOME", err)
	}
	return c
}

// fakeController plays a controller at the far end of a pipe: it sends its banner, reads
// one frame and answers it with the given bytes. It returns dishd's client.
func fakeController(t *testing.T, answer []byte) *Client {
	t.Helper()
	conn, ctrl := net.Pipe()
	deadline := time.Now().Add(time.Second)
	conn.SetDeadline(deadline)
	ctrl.SetDeadline(deadline)
	done := make(chan struct{})
	go func() {
		defer close(done)
		if _, err := io.WriteString(ctrl, "YST-1 READY PROTO 1\r\n"); err != nil {
			return
		}
		var length [2]byte
		if _, err := io.ReadFull(ctrl, length[:]); err != nil {
			return
		}
		if _, err := io.ReadFull(ctrl, make([]byte, binary.BigEndian.Uint16(length[:]))); err != nil {
			return
		}
		ctrl.Write(answer)
	}()
	t.Cleanup(func() {
		conn.Close()
		ctrl.Close()
		<-done
	})
	c, err := newClient(conn)
	if err != nil {
		t.Fatalf("newClient(fake controller) = %v, want a client", err)
	}
	return c
}
