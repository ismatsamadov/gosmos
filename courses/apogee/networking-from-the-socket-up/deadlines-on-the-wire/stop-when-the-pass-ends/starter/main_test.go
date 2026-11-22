package main

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

var hello = Frame{Type: TypeHello, Body: []byte("dishd")}

// Every test runs in a synctest bubble, whose clock is fake: dishd's 5s timeout keeps its real
// length and passes in no real time.

func TestRecvCtx(t *testing.T) {
	t.Run("the pass ends during a silent read", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.Silent = true // reads every frame and answers none
			c := connect(t, sim)
			send(t, c, hello)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			defer cancel()
			start := time.Now()
			r := recvCtxWithin(ctx, c, time.Minute)
			took := time.Since(start)
			if !r.done || !errors.Is(r.err, context.DeadlineExceeded) || took != 20*time.Millisecond {
				t.Errorf("recvCtx(a context ending after 20ms, a silent controller) = %v after %v, want context.DeadlineExceeded after 20ms, not c.Timeout's 5s", r.err, took)
			}
		})
	})
	t.Run("a pass that has already ended", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			sim.Silent = true
			c, line := connectGated(t, sim) // holds back a deadline cut, so only checking first ends the call at once
			send(t, c, hello)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			got := make(chan received, 1)
			go func() {
				f, err := recvCtx(ctx, c)
				got <- received{f, err, true}
			}()
			synctest.Wait()
			select {
			case r := <-got:
				if !errors.Is(r.err, context.Canceled) {
					t.Errorf("recvCtx(a canceled context) = %v, want context.Canceled", r.err)
				}
			default:
				line.open()
				r := <-got
				t.Errorf("recvCtx(a canceled context) = %v, but only after it had started to read, want context.Canceled at once: check ctx.Err() before anything else", r.err)
			}
		})
	})
	t.Run("an answer in time, then a later read", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := connect(t, NewSim("YST-1", 1))
			send(t, c, hello)
			pass, endPass := context.WithCancel(context.Background())
			defer endPass()
			r := recvCtxWithin(pass, c, time.Minute)
			if !r.done || r.err != nil || r.f.Type != TypeWelcome {
				t.Fatalf("recvCtx(a pass still on) after HELLO = type 0x%02x, %v, want WELCOME (0x81), nil", r.f.Type, r.err)
			}
			later := make(chan received, 1)
			go func() {
				f, err := c.recv()
				later <- received{f, err, true}
			}()
			synctest.Wait() // the later read is waiting for an answer
			endPass()
			synctest.Wait()
			send(t, c, Frame{Type: TypePing})
			if r := <-later; r.err != nil || r.f.Type != TypePong {
				t.Errorf("a later recv(), waiting when the first pass's context ended = type 0x%02x, %v, want PONG (0x82), nil: call stop once the read is done, or the old context still cuts reads", r.f.Type, r.err)
			}
		})
	})
	t.Run("the pass ends as the answer arrives", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-1", 1)
			answer := make(chan struct{})
			sim.Before = func(req string) {
				if req == "HELLO" {
					<-answer // the WELCOME waits until the test lets it go
				}
			}
			c, line := connectGated(t, sim)
			var letGo sync.Once
			release := func() { letGo.Do(func() { close(answer) }) }
			t.Cleanup(release) // runs before the line closes, so the simulator can finish
			send(t, c, hello)
			pass, endPass := context.WithCancel(context.Background())
			defer endPass()
			got := make(chan received, 1)
			go func() {
				f, err := recvCtx(pass, c)
				got <- received{f, err, true}
			}()
			synctest.Wait() // recvCtx is waiting for the WELCOME
			endPass()       // LOS: the function given to AfterFunc starts, and its deadline is held back
			synctest.Wait()
			if line.held() == 0 {
				line.open()
				t.Fatalf("when the pass's context ended, nothing set c.conn's read deadline: register context.AfterFunc to set it to now")
			}
			release() // the WELCOME arrives all the same, and the read returns it
			synctest.Wait()
			select {
			case r := <-got:
				line.open()
				t.Fatalf("recvCtx(a pass ending as the WELCOME arrives) returned %v while the deadline set at LOS was still on its way, want it to wait: stop reports false once that function has started, and its late deadline would cut the next read", r.err)
			default:
			}
			line.open() // the held deadline lands now
			r := <-got
			if !(r.err == nil && r.f.Type == TypeWelcome) && !errors.Is(r.err, context.Canceled) {
				t.Errorf("recvCtx(a pass ending as the WELCOME arrives) = type 0x%02x, %v, want WELCOME (0x81), nil, or context.Canceled", r.f.Type, r.err)
			}
			send(t, c, Frame{Type: TypePing})
			if f, err := c.recv(); err != nil || f.Type != TypePong {
				t.Errorf("the next recv() after that pass = type 0x%02x, %v, want PONG (0x82), nil", f.Type, err)
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

// connectGated is connect with the client's end of the line wrapped in a gatedLine, which holds
// back every read deadline that has already passed until the test opens it.
func connectGated(t *testing.T, sim *Sim) (*Client, *gatedLine) {
	t.Helper()
	end, simEnd := net.Pipe()
	line := &gatedLine{Conn: end, gate: make(chan struct{})}
	var wg sync.WaitGroup
	wg.Go(func() { sim.ServeConn(simEnd) })
	t.Cleanup(func() {
		line.open()
		end.Close()
		simEnd.Close()
		wg.Wait()
	})
	c, err := newClient(line)
	if err != nil {
		t.Fatalf("newClient(simulator) = %v, want a client", err)
	}
	return c, line
}

// gatedLine is a connection that holds back deadlines that have already passed: such a
// deadline takes effect only once the test opens the gate, as if the goroutine setting it
// were slow to run. Deadlines still ahead pass straight through.
type gatedLine struct {
	net.Conn
	gate     chan struct{}
	openOnce sync.Once
	mu       sync.Mutex
	waiting  int
}

// SetReadDeadline sets the read deadline, after the gate opens if t has already passed.
func (g *gatedLine) SetReadDeadline(t time.Time) error {
	g.hold(t)
	return g.Conn.SetReadDeadline(t)
}

// SetDeadline sets both deadlines, after the gate opens if t has already passed.
func (g *gatedLine) SetDeadline(t time.Time) error {
	g.hold(t)
	return g.Conn.SetDeadline(t)
}

// hold waits for the gate if t has already passed.
func (g *gatedLine) hold(t time.Time) {
	if t.IsZero() || t.After(time.Now()) {
		return
	}
	g.mu.Lock()
	g.waiting++
	g.mu.Unlock()
	<-g.gate
}

// held reports how many deadlines have been held back so far.
func (g *gatedLine) held() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.waiting
}

// open lets every held deadline through, and every later one at once.
func (g *gatedLine) open() {
	g.openOnce.Do(func() { close(g.gate) })
}

// send sends f through c, failing the test if it can't.
func send(t *testing.T, c *Client, f Frame) {
	t.Helper()
	if err := c.send(f); err != nil {
		t.Fatalf("send(type 0x%02x) = %v, want nil", f.Type, err)
	}
}

// received is what one read returned; done is false when it was still waiting.
type received struct {
	f    Frame
	err  error
	done bool
}

// recvCtxWithin calls recvCtx and waits for it at most limit. If the call is still waiting
// then, it hangs up on the controller so the call ends, and reports done as false.
func recvCtxWithin(ctx context.Context, c *Client, limit time.Duration) received {
	got := make(chan received, 1)
	go func() {
		f, err := recvCtx(ctx, c)
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
