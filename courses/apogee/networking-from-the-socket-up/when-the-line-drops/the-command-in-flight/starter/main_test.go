package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"slices"
	"testing"
	"time"
)

func TestResume(t *testing.T) {
	jog := Frame{Type: TypeJog, Seq: 122, Body: []byte{'A', 0x00, 0x23}} // A +0.35°
	tests := []struct {
		name     string
		inflight *Frame
		lastSeq  uint32
		want     *Frame
	}{
		{"the jog ran", &jog, 122, nil},
		{"the jog never arrived", &jog, 121, &jog},
		{"later commands ran too", &jog, 150, nil},
		{"nothing in flight", nil, 121, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Link{next: 123, inflight: tt.inflight}
			got := l.resume(Welcome{Dish: 2, FwMajor: 2, FwMinor: 4, LastSeq: tt.lastSeq})
			if !sameFrame(got, tt.want) {
				t.Errorf("resume(last_seq %d) with %s in flight = %s, want %s", tt.lastSeq, show(tt.inflight), show(got), show(tt.want))
			}
			if !sameFrame(l.inflight, tt.want) {
				t.Errorf("after resume(last_seq %d), l.inflight = %s, want %s", tt.lastSeq, show(l.inflight), show(tt.want))
			}
		})
	}
	drops := []struct {
		name string
		drop func(*Sim)
	}{
		{"the line drops after the jog ran", func(s *Sim) { s.DropAfter = 122 }},
		{"the line drops before the jog arrived", func(s *Sim) { s.DropBefore = 122 }},
	}
	for _, d := range drops {
		t.Run(d.name, func(t *testing.T) {
			sim := NewSim("YST-2", 2)
			sim.LastSeq = 121
			d.drop(sim)
			l := &Link{}
			c, w := dialSim(t, sim)
			l.learn(w)
			sent := l.number(Frame{Type: TypeJog, Body: jog.Body})
			l.inflight = &sent
			if err := c.command(sent); !errors.Is(err, io.EOF) {
				t.Fatalf("JOG %d = %v, want EOF: the simulator drops the line", sent.Seq, err)
			}
			c.Close()
			c, w = dialSim(t, sim)
			l.learn(w)
			if f := l.resume(w); f != nil {
				if err := c.command(*f); err != nil {
					t.Errorf("after a WELCOME with last_seq %d, resume gave %s, and sending it = %v, want nil", w.LastSeq, show(f), err)
				}
			}
			want := []string{"122 JOG A +0.35"}
			if got := sim.Executed(); !slices.Equal(got, want) {
				t.Errorf("the controller executed %q, want %q: the jog once, as 122", got, want)
			}
		})
	}
}

// dialSim connects to sim through a fresh net.Pipe, reads the banner and says HELLO, and
// returns the client and the WELCOME. The connection closes when the test ends.
func dialSim(t *testing.T, sim *Sim) (*Client, Welcome) {
	t.Helper()
	a, b := net.Pipe()
	done := make(chan struct{})
	go func() {
		sim.ServeConn(b)
		close(done)
	}()
	t.Cleanup(func() {
		a.Close()
		b.Close()
		<-done
	})
	a.SetDeadline(time.Now().Add(time.Second))
	c, err := newClient(a)
	if err != nil {
		t.Fatalf("newClient: %v", err)
	}
	c.Timeout = time.Second
	w, err := c.hello()
	if err != nil {
		t.Fatalf("HELLO: %v", err)
	}
	return c, w
}

// sameFrame reports whether a and b are both nil, or frames with the same type, seq and
// body.
func sameFrame(a, b *Frame) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Type == b.Type && a.Seq == b.Seq && slices.Equal(a.Body, b.Body)
}

// show describes a frame in flight for a message: "JOG 122" or "nil".
func show(f *Frame) string {
	if f == nil {
		return "nil"
	}
	if f.Type == TypeJog {
		return fmt.Sprintf("JOG %d", f.Seq)
	}
	return fmt.Sprintf("type 0x%02x seq %d", f.Type, f.Seq)
}
