package main

import (
	"errors"
	"net"
	"os"
	"testing"
	"testing/synctest"
	"time"
)

var (
	hello = Frame{Type: TypeHello, Body: []byte("dishd")}
	ping  = Frame{Type: TypePing}
)

// Every test runs in a synctest bubble, whose clock is fake: the controller and dishd keep
// their real limits (5s for dishd's reads, 30s of silence before the controller hangs up),
// and those seconds pass in no real time.

func TestRecv(t *testing.T) {
	t.Run("a controller that never answers", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.Silent = true // reads every frame and answers none, like Tuesday's dead line
			c := connect(t, sim)
			send(t, c, hello)
			start := time.Now()
			r := recvWithin(c, time.Minute)
			took := time.Since(start)
			if !r.done {
				t.Fatalf("recv() from a controller that never answers is still waiting after %v, want an error after c.Timeout, 5s: set the read deadline before reading", took)
			}
			if !errors.Is(r.err, os.ErrDeadlineExceeded) || took != 5*time.Second {
				t.Errorf("recv() from a controller that never answers = %v after %v, want os.ErrDeadlineExceeded after c.Timeout, 5s", r.err, took)
			}
		})
	})
	t.Run("a controller that answers", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := connect(t, NewSim("YST-1", 1))
			send(t, c, hello)
			r := recvWithin(c, time.Minute)
			if !r.done || r.err != nil || r.f.Type != TypeWelcome {
				t.Errorf("recv() after HELLO = type 0x%02x, %v, want WELCOME (0x81), nil", r.f.Type, r.err)
			}
		})
	})
	t.Run("five exchanges 2s apart", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := connect(t, NewSim("YST-1", 1))
			frames := []Frame{hello, ping, ping, ping, ping}
			want := []byte{TypeWelcome, TypePong, TypePong, TypePong, TypePong}
			for i, f := range frames {
				if i > 0 {
					time.Sleep(2 * time.Second)
				}
				send(t, c, f)
				r := recvWithin(c, time.Minute)
				if !r.done || r.err != nil || r.f.Type != want[i] {
					t.Fatalf("exchange %d, %ds after the first: recv() = type 0x%02x, %v, want type 0x%02x, nil: a deadline is an instant, so give every read a new one", i+1, 2*i, r.f.Type, r.err, want[i])
				}
			}
		})
	})
}

// connect serves sim on one end of a pipe and returns a client on the other end, its banner
// read. When the test ends, both ends close and the simulator stops.
func connect(t *testing.T, sim *Sim) *Client {
	t.Helper()
	end, simEnd := net.Pipe()
	stopped := make(chan struct{})
	go func() {
		sim.ServeConn(simEnd)
		close(stopped)
	}()
	t.Cleanup(func() {
		end.Close()
		simEnd.Close()
		<-stopped
	})
	c, err := newClient(end)
	if err != nil {
		t.Fatalf("newClient(simulator) = %v, want a client", err)
	}
	return c
}

// send sends f through c, failing the test if it can't.
func send(t *testing.T, c *Client, f Frame) {
	t.Helper()
	if err := c.send(f); err != nil {
		t.Fatalf("send(type 0x%02x) = %v, want nil", f.Type, err)
	}
}

// received is what one call to recv returned; done is false when it was still waiting.
type received struct {
	f    Frame
	err  error
	done bool
}

// recvWithin calls c.recv and waits for it at most limit. If the call is still waiting then,
// it hangs up on the controller so the call ends, and reports done as false.
func recvWithin(c *Client, limit time.Duration) received {
	got := make(chan received, 1)
	go func() {
		f, err := c.recv()
		got <- received{f, err, true}
	}()
	select {
	case r := <-got:
		return r
	case <-time.After(limit):
		c.conn.Close()
		<-got
		return received{}
	}
}
