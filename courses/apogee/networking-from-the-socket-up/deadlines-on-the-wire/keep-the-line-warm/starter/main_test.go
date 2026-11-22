package main

import (
	"errors"
	"net"
	"os"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// Every test runs in a synctest bubble, whose clock is fake: ticks 10s apart, dishd's 5s
// timeout and the controller's 30s idle limit keep their real lengths and pass in no real time.

func TestHeartbeat(t *testing.T) {
	t.Run("three ticks", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			c := connect(t, sim)
			exchange(t, c, Frame{Type: TypeHello, Body: []byte("dishd")}, TypeWelcome)
			err, done := heartbeatWithin(c, ticks(t, 3), 2*time.Minute)
			if !done || err != nil {
				t.Fatalf("heartbeat(three ticks 10s apart) = %v, want nil once the ticks end (a controller whose answers go unread hangs up after 10s)", err)
			}
			if got, want := sim.Received(), []string{"HELLO dishd", "PING", "PING", "PING"}; !slices.Equal(got, want) {
				t.Errorf("after heartbeat(three ticks), the simulator received %q, want %q", got, want)
			}
			if err := c.send(Frame{Type: TypeQuery}); err != nil {
				t.Fatalf("send(QUERY) after the heartbeat = %v, want nil", err)
			}
			if f, err := c.recv(); err != nil || f.Type != TypeStatus {
				t.Errorf("after heartbeat(three ticks), recv() for a QUERY = type 0x%02x, %v, want STATUS (0x83): read the PONG of every PING, or it waits for the next command", f.Type, err)
			}
		})
	})
	t.Run("the controller goes quiet after two answers", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.MuteAfter = 2 // answers HELLO and the first PING, then nothing
			c := connect(t, sim)
			exchange(t, c, Frame{Type: TypeHello, Body: []byte("dishd")}, TypeWelcome)
			start := time.Now()
			err, done := heartbeatWithin(c, ticks(t, 3), 2*time.Minute)
			took := time.Since(start)
			if !done || !errors.Is(err, os.ErrDeadlineExceeded) || took != 25*time.Second {
				t.Errorf("heartbeat(three ticks 10s apart) on a line that goes quiet after the first PONG = %v after %v, want os.ErrDeadlineExceeded after 25s: read each PONG, and return the first error", err, took)
			}
			if got, want := sim.Received(), []string{"HELLO dishd", "PING", "PING"}; !slices.Equal(got, want) {
				t.Errorf("after the second tick went unanswered, the simulator had received %q, want %q", got, want)
			}
		})
	})
	t.Run("an answer that isn't a PONG", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := connect(t, NewSim("YST-1", 1)) // no HELLO yet, so PING gets NAK "hello first"
			err, done := heartbeatWithin(c, ticks(t, 1), 2*time.Minute)
			if !done || err == nil {
				t.Errorf("heartbeat(one tick) answered by a NAK = %v, want an error: only a PONG proves the line", err)
			}
		})
	})
	t.Run("no ticks", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			c := connect(t, sim)
			exchange(t, c, Frame{Type: TypeHello, Body: []byte("dishd")}, TypeWelcome)
			none := make(chan time.Time)
			close(none)
			if err, done := heartbeatWithin(c, none, 2*time.Minute); !done || err != nil {
				t.Fatalf("heartbeat(no ticks) = %v, want nil at once", err)
			}
			c.Timeout = time.Minute // long enough to hear the controller give up
			start := time.Now()
			_, err := c.recv()
			var bye *ByeError
			if !errors.As(err, &bye) || bye.Reason != "idle" || time.Since(start) != 30*time.Second {
				t.Errorf("recv() with no heartbeat = %v after %v, want controller closed: idle after the controller's 30s idle limit", err, time.Since(start))
			}
		})
	})
}

// ticks returns a channel that ticks n times, 10s apart, as dishd's heartbeat ticker would,
// then closes. The ticking stops when the test ends.
func ticks(t *testing.T, n int) <-chan time.Time {
	ch := make(chan time.Time, n)
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Go(func() {
		defer close(ch)
		for range n {
			select {
			case <-stop:
				return
			case now := <-time.After(10 * time.Second):
				ch <- now
			}
		}
	})
	t.Cleanup(func() {
		close(stop)
		wg.Wait()
	})
	return ch
}

// heartbeatWithin calls heartbeat and waits for it at most limit. If it is still running then,
// it hangs up on the controller so the call ends, and reports done as false.
func heartbeatWithin(c *Client, ticks <-chan time.Time, limit time.Duration) (err error, done bool) {
	result := make(chan error, 1)
	go func() { result <- heartbeat(c, ticks) }()
	select {
	case err := <-result:
		return err, true
	case <-time.After(limit):
		c.conn.Close()
		<-result
		return nil, false
	}
}

// exchange sends f and requires an answer of type want.
func exchange(t *testing.T, c *Client, f Frame, want byte) {
	t.Helper()
	if err := c.send(f); err != nil {
		t.Fatalf("send(type 0x%02x) = %v, want nil", f.Type, err)
	}
	if got, err := c.recv(); err != nil || got.Type != want {
		t.Fatalf("recv() after type 0x%02x = type 0x%02x, %v, want type 0x%02x", f.Type, got.Type, err, want)
	}
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
