package main

import (
	"bytes"
	"fmt"
	"net"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestEncodeFrame(t *testing.T) {
	tests := []struct {
		name string
		f    Frame
		want string
	}{
		{"PING", Frame{Type: TypePing}, "00 09 02 00 00 00 00 bc e2 a4 7d"},
		{"POINT seq 119 to 183.25, 41.70", Frame{Type: TypePoint, Seq: 119, Body: []byte{0x47, 0x95, 0x10, 0x4a}},
			"00 0d 10 00 00 00 77 47 95 10 4a 20 eb ca 40"},
		{"HELLO from dishd", Frame{Type: TypeHello, Body: []byte("dishd")},
			"00 0e 01 00 00 00 00 64 69 73 68 64 97 98 f8 ab"},
		{"STOW seq 120", Frame{Type: TypeStow, Seq: 120}, "00 09 11 00 00 00 78 c5 7c b0 21"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf("% x", encodeFrame(tt.f))
			if got != tt.want {
				t.Errorf("encodeFrame(%s) = %s, want %s%s", tt.name, got, tt.want, mistake(got))
			}
		})
	}
	t.Run("a nil body and an empty one", func(t *testing.T) {
		if a, b := encodeFrame(Frame{Type: TypePing}), encodeFrame(Frame{Type: TypePing, Body: []byte{}}); !bytes.Equal(a, b) {
			t.Errorf("encodeFrame(PING with a nil body) = % x, but with an empty body % x; want the same frame", a, b)
		}
	})
	t.Run("the simulator ACKs the POINT", func(t *testing.T) {
		sim := NewSim("YST-1", 1)
		sim.LastSeq = 118
		conn, end := net.Pipe()
		served := make(chan struct{})
		go func() {
			sim.ServeConn(end)
			close(served)
		}()
		defer func() {
			conn.Close()
			<-served
		}()
		conn.SetDeadline(time.Now().Add(time.Second))
		c, err := newClient(conn)
		if err != nil {
			t.Fatalf("newClient(simulator) = %v, want a client", err)
		}
		if _, err := c.hello(); err != nil {
			t.Fatalf("hello() = %v, want the WELCOME: the controller checks every frame's CRC", err)
		}
		if err := c.send(Frame{Type: TypePoint, Seq: 119, Body: []byte{0x47, 0x95, 0x10, 0x4a}}); err != nil {
			t.Fatalf("send(POINT 119) = %v, want nil", err)
		}
		ans, err := c.recv()
		if err != nil || ans.Type != TypeAck || ans.Seq != 119 {
			t.Fatalf("the answer to POINT 119 = type %#02x seq %d, %v; want ACK 119 (type 0x90)", ans.Type, ans.Seq, err)
		}
		if got, want := sim.Executed(), []string{"119 POINT 183.25 41.70"}; !slices.Equal(got, want) {
			t.Errorf("the simulator executed %q, want %q", got, want)
		}
	})
}

// mistake names the two CRC mistakes the protocol's table shows for a PING.
func mistake(got string) string {
	switch {
	case strings.HasSuffix(got, "c0 a3 dd d6"):
		return ": that CRC covers the length too, and should start after it"
	case strings.HasSuffix(got, "7d a4 e2 bc"):
		return ": that CRC is little-endian, and should be big-endian"
	}
	return ""
}
