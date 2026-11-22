package main

import (
	"bytes"
	"context"
	"errors"
	"net"
	"os"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

var (
	hello = Frame{Type: TypeHello, Body: []byte("dishd")}
	query = Frame{Type: TypeQuery}
)

// Sunniva's replay of Tuesday's storm. Every line runs in a synctest bubble, whose clock is
// fake: a byte every 0.4s, dishd's 5s timeout and a pass ending at 2s keep their real lengths
// and pass in no real time.

func TestDo(t *testing.T) {
	t.Run("a calm line", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := connect(t, NewSim("YST-1", 1))
			r := doWithin(context.Background(), c, hello)
			if r.err != nil || r.f.Type != TypeWelcome || r.took != 0 {
				t.Errorf("Do(ctx, HELLO) on a calm line = type 0x%02x, %v after %v, want WELCOME (0x81), nil at once", r.f.Type, r.err, r.took)
			}
		})
	})
	t.Run("a slow line, a byte every 0.4s", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.Split, sim.Gap = 1, 400*time.Millisecond // every message a byte at a time
			c := connect(t, sim)
			if r := doWithin(context.Background(), c, hello); r.err != nil || r.f.Type != TypeWelcome {
				t.Fatalf("Do(ctx, HELLO) on a slow line = type 0x%02x, %v after %v, want WELCOME (0x81), nil", r.f.Type, r.err, r.took)
			}
			r := doWithin(context.Background(), c, query)
			want := []byte{0x01, 0x00, 0x46, 0x50, 0x11, 0x94, 0x00, 0xc8, 0xff, 0x83, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}
			if r.err != nil || r.f.Type != TypeStatus || !bytes.Equal(r.f.Body, want) || r.took != 10400*time.Millisecond {
				t.Errorf("Do(ctx, QUERY) on a slow line = type 0x%02x, body % x, %v after %v, want STATUS (0x83), body % x, nil after 10.4s", r.f.Type, r.f.Body, r.err, r.took, want)
			}
		})
	})
	t.Run("a silent line", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.Silent = true // reads every frame and answers none
			c := connect(t, sim)
			r := doWithin(context.Background(), c, hello)
			if r.err == nil || r.took != 5*time.Second {
				t.Fatalf("Do(ctx, HELLO) on a silent line = %v after %v, want an error after 5s", r.err, r.took)
			}
			if got, want := describe(c.Name, c.Timeout, r.err), "YST-1: no answer in 5s"; got != want {
				t.Errorf("describe(%q, %v, Do's error %v) = %q, want %q", c.Name, c.Timeout, r.err, got, want)
			}
		})
	})
	t.Run("a controller that stops reading", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.NoRead = true // sends its banner, then never reads a byte
			c := connect(t, sim)
			r := doWithin(context.Background(), c, hello)
			var op *net.OpError
			if !errors.Is(r.err, os.ErrDeadlineExceeded) || !errors.As(r.err, &op) || op.Op != "write" || r.took != 5*time.Second {
				t.Errorf("Do(ctx, HELLO) to a controller that never reads = %v after %v, want the write's own deadline error after 5s: when the send fails, return its error without reading", r.err, r.took)
			}
		})
	})
	t.Run("the pass ends after 2s", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.Silent = true
			c := connect(t, sim)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			r := doWithin(ctx, c, hello)
			if !errors.Is(r.err, context.DeadlineExceeded) || r.took != 2*time.Second {
				t.Errorf("Do(a context ending after 2s, HELLO) on a silent line = %v after %v, want context.DeadlineExceeded after 2s: read the answer with recvCtx", r.err, r.took)
			}
		})
	})
}

// connect serves sim on one end of a pipe and returns a client on the other end, its banner
// read. When the test ends, both ends close and the simulator stops.
func connect(t *testing.T, sim *Sim) *Client {
	t.Helper()
	end, simEnd := net.Pipe()
	var wg sync.WaitGroup
	wg.Go(func() { sim.ServeConn(simEnd) })
	t.Cleanup(func() {
		end.Close()
		simEnd.Close()
		wg.Wait()
	})
	c, err := newClient(end)
	if err != nil {
		t.Fatalf("newClient(simulator) = %v, want a client", err)
	}
	return c
}

// done is what one call to Do returned, and how long it took.
type done struct {
	f    Frame
	err  error
	took time.Duration
}

// doWithin calls c.Do(ctx, f) and waits for it at most two minutes. If the call is still
// running then, it hangs up on the controller so the call ends.
func doWithin(ctx context.Context, c *Client, f Frame) done {
	start := time.Now()
	got := make(chan done, 1)
	go func() {
		ans, err := c.Do(ctx, f)
		got <- done{ans, err, time.Since(start)}
	}()
	select {
	case r := <-got:
		return r
	case <-time.After(2 * time.Minute):
		c.conn.Close()
		r := <-got
		r.err = errors.New("still waiting after 2m0s")
		return r
	}
}
