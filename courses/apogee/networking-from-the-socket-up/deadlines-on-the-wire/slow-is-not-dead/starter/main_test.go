package main

import (
	"bytes"
	"errors"
	"net"
	"os"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// status is a whole STATUS frame from Dish 1, 27 bytes: tracking at 183.25° and 41.70°.
var status = []byte{
	0x00, 0x19, 0x83, 0x00, 0x00, 0x00, 0x00,
	0x01, 0x02, 0x47, 0x95, 0x10, 0x4a, 0x02, 0x80, 0xff, 0x83, 0x00, 0x01, 0x00, 0x00, 0x00, 0x76,
	0x62, 0x9d, 0xe0, 0xf1,
}

// Every test runs in a synctest bubble, whose clock is fake: a byte every 0.4s and dishd's
// 5s timeout keep their real lengths and pass in no real time.

func TestIdleReader(t *testing.T) {
	t.Run("a STATUS at a byte every 0.4s", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := slowLine(t, status, 400*time.Millisecond)
			start := time.Now()
			r := recvWithin(c, time.Minute)
			took := time.Since(start)
			if r.err != nil || !r.done {
				t.Fatalf("recv() of a 27-byte STATUS arriving a byte every 0.4s = %v after %v, want the frame after 10.4s: give every Read its own deadline, not the whole frame one", r.err, took)
			}
			if r.f.Type != TypeStatus || !bytes.Equal(r.f.Body, status[7:23]) || took != 10400*time.Millisecond {
				t.Errorf("recv() of a STATUS arriving a byte every 0.4s = type 0x%02x, body % x after %v, want type 0x83, body % x after 10.4s", r.f.Type, r.f.Body, took, status[7:23])
			}
		})
	})
	t.Run("a byte every 6s", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := slowLine(t, status, 6*time.Second)
			start := time.Now()
			r := recvWithin(c, time.Minute)
			took := time.Since(start)
			if !errors.Is(r.err, os.ErrDeadlineExceeded) || took != 5*time.Second {
				t.Errorf("recv() of a STATUS arriving a byte every 6s = %v after %v, want os.ErrDeadlineExceeded after 5s: a read that brings nothing in c.Timeout means the line is gone", r.err, took)
			}
		})
	})
	t.Run("a silent simulator, asked 20s after its banner", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.Silent = true // reads every frame and answers none
			c := connect(t, sim)
			time.Sleep(20 * time.Second)
			if err := c.send(Frame{Type: TypeHello, Body: []byte("dishd")}); err != nil {
				t.Fatalf("send(HELLO) = %v, want nil", err)
			}
			start := time.Now()
			r := recvWithin(c, time.Minute)
			took := time.Since(start)
			if !errors.Is(r.err, os.ErrDeadlineExceeded) || took != 5*time.Second {
				t.Errorf("recv() from a silent simulator = %v after %v, want os.ErrDeadlineExceeded 5s after the read began: set the deadline before each Read, not after", r.err, took)
			}
		})
	})
}

// slowLine plays a controller on one end of a pipe: it sends its banner whole, then data a
// byte at a time, gap apart, the way Tuesday's backhaul did. It returns a client on the other
// end. When the test ends, both ends close and the controller stops.
func slowLine(t *testing.T, data []byte, gap time.Duration) *Client {
	t.Helper()
	end, line := net.Pipe()
	line.SetDeadline(time.Now().Add(time.Minute))
	var wg sync.WaitGroup
	wg.Go(func() {
		if _, err := line.Write([]byte("YST-1 READY PROTO 1\r\n")); err != nil {
			return
		}
		for i := range data {
			if i > 0 {
				time.Sleep(gap)
			}
			if _, err := line.Write(data[i : i+1]); err != nil {
				return
			}
		}
	})
	t.Cleanup(func() {
		end.Close()
		line.Close()
		wg.Wait()
	})
	c, err := newClient(end)
	if err != nil {
		t.Fatalf("newClient(a banner sent whole) = %v, want a client", err)
	}
	return c
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
