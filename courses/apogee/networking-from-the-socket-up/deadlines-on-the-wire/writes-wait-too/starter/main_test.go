package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

var (
	hello = Frame{Type: TypeHello, Body: []byte("dishd")}
	ping  = Frame{Type: TypePing}
	point = Frame{Type: TypePoint, Seq: 119, Body: []byte{0x47, 0x95, 0x10, 0x4a}} // 183.25°, 41.70°
)

// Every test runs in a synctest bubble, whose clock is fake: dishd's 5s timeout keeps its real
// length and passes in no real time.

func TestSend(t *testing.T) {
	t.Run("a controller that never reads", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.NoRead = true // sends its banner, then never reads a byte
			c := connect(t, sim)
			start := time.Now()
			err, done := sendWithin(c, point, time.Minute)
			took := time.Since(start)
			if !done {
				t.Fatalf("send(POINT 119) to a controller that never reads is still writing after %v, want an error after c.Timeout, 5s: set a write deadline before writing", took)
			}
			if !errors.Is(err, os.ErrDeadlineExceeded) || took != 5*time.Second {
				t.Errorf("send(POINT 119) to a controller that never reads = %v after %v, want os.ErrDeadlineExceeded after 5s", err, took)
			}
		})
	})
	t.Run("a POINT arrives byte for byte", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			end, line := net.Pipe()
			defer line.Close()
			defer end.Close()
			line.SetDeadline(time.Now().Add(time.Minute))
			go line.Write([]byte("YST-1 READY PROTO 1\r\n"))
			c, err := newClient(end)
			if err != nil {
				t.Fatalf("newClient(a banner) = %v, want a client", err)
			}
			got := make(chan []byte, 1)
			go func() {
				b := make([]byte, 15)
				n, _ := io.ReadFull(line, b)
				got <- b[:n]
			}()
			if err, done := sendWithin(c, point, time.Minute); !done || err != nil {
				t.Fatalf("send(POINT 119) to a controller that reads = %v, want nil", err)
			}
			want := "00 0d 10 00 00 00 77 47 95 10 4a 20 eb ca 40"
			if b := fmt.Sprintf("% x", <-got); b != want {
				t.Errorf("send(POINT 119) wrote %s, want %s", b, want)
			}
		})
	})
	t.Run("five sends 2s apart", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.WriteLimit = time.Minute // waits for dishd to read its answers, which it never does here
			c := connect(t, sim)
			frames := []Frame{hello, ping, ping, ping, ping}
			for i, f := range frames {
				if i > 0 {
					time.Sleep(2 * time.Second)
				}
				if err, done := sendWithin(c, f, time.Minute); !done || err != nil {
					t.Fatalf("send %d, %ds after the first: send(type 0x%02x) = %v, want nil: a deadline is an instant, so give every write a new one", i+1, 2*i, f.Type, err)
				}
			}
			synctest.Wait()
			want := []string{"HELLO dishd", "PING", "PING", "PING", "PING"}
			if got := sim.Received(); !slices.Equal(got, want) {
				t.Errorf("after five sends, the simulator received %q, want %q", got, want)
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

// sendWithin calls c.send(f) and waits for it at most limit. If the call is still writing
// then, it hangs up on the controller so the call ends, and reports done as false.
func sendWithin(c *Client, f Frame, limit time.Duration) (err error, done bool) {
	sent := make(chan error, 1)
	go func() { sent <- c.send(f) }()
	select {
	case err := <-sent:
		return err, true
	case <-time.After(limit):
		c.conn.Close()
		<-sent
		return nil, false
	}
}
