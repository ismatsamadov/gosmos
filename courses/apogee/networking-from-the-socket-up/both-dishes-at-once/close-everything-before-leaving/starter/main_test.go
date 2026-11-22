package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

const hello = "dishd watching YST-1 YST-2\n"

// The first two tests run in a synctest bubble, with the watch port in memory: their guards
// pass in no real time, and the test knows when every goroutine is waiting.

func TestShutdown(t *testing.T) {
	t.Run("two dishes and three watchers", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim1, sim2 := NewSim("YST-1", 1), NewSim("YST-2", 2)
			dish1, stopped1 := dish(t, sim1)
			dish2, stopped2 := dish(t, sim2)
			ln := newPipeListener()
			t.Cleanup(func() { ln.Close() })
			s := &Service{Dishes: []*Client{dish1, dish2}, Watchers: &Watchers{Timeout: 2 * time.Second}}
			s.Start(ln)
			var screens []*screen
			for i := 1; i <= 3; i++ {
				w := ln.dial(t, i)
				if line, err := w.next(); line != hello {
					t.Fatalf("watcher %d read %q, %v, want the greeting %q", i, line, err, hello)
				}
				screens = append(screens, w)
			}

			if !within(s.Shutdown, time.Second) {
				t.Fatalf("Shutdown() has not returned after 1s, with three watchers connected: close the listener, the watchers and the dishes before waiting, or the watchers' goroutines never end")
			}
			for i, w := range screens {
				if line, err := w.next(); !errors.Is(err, io.EOF) {
					t.Errorf("after Shutdown(), watcher %d's next read = %q, %v, want io.EOF: close every watcher", i+1, line, err)
				}
			}
			for _, d := range []struct {
				name    string
				sim     *Sim
				stopped <-chan struct{}
			}{{"YST-1", sim1, stopped1}, {"YST-2", sim2, stopped2}} {
				select {
				case <-d.stopped:
					if n := d.sim.ClosedByClient(); n != 1 {
						t.Errorf("after Shutdown(), %s's simulator counted %d connections closed by dishd, want 1", d.name, n)
					}
				case <-time.After(time.Second):
					t.Errorf("a second after Shutdown(), %s's simulator still had dishd's connection open: close every dish client", d.name)
				}
			}
			if !ln.isClosed() {
				t.Errorf("after Shutdown(), the watch port was still open: close the listener")
			}
			second := make(chan struct{})
			go func() {
				s.Shutdown()
				close(second)
			}()
			synctest.Wait()
			select {
			case <-second:
			default:
				t.Errorf("a second Shutdown() has not returned at once, want it to do nothing")
			}
		})
	})
	t.Run("a watcher who joins as dishd stops", func(t *testing.T) {
		// The watcher joins at exactly the worst moment: accepted before Shutdown, greeted only
		// after the watchers have been closed.
		synctest.Test(t, func(t *testing.T) {
			dish1, _ := dish(t, NewSim("YST-1", 1))
			dish2, _ := dish(t, NewSim("YST-2", 2))
			ln := newPipeListener()
			t.Cleanup(func() { ln.Close() })
			s := &Service{Dishes: []*Client{dish1, dish2}, Watchers: &Watchers{Timeout: 2 * time.Second}}
			s.Start(ln)
			sunniva := ln.dial(t, 1)
			if line, err := sunniva.next(); line != hello {
				t.Fatalf("Sunniva's terminal read %q, %v, want the greeting %q", line, err, hello)
			}
			ingrid := ln.dial(t, 2) // accepted, and her greeting waits until she reads it

			stopped := make(chan struct{})
			go func() {
				s.Shutdown()
				close(stopped)
			}()
			if line, err := sunniva.next(); !errors.Is(err, io.EOF) {
				t.Fatalf("after Shutdown() began, Sunniva's terminal read %q, %v, want io.EOF: close every watcher before waiting", line, err)
			}
			synctest.Wait()
			select {
			case <-stopped:
				t.Errorf("Shutdown() returned while the goroutine greeting Ingrid's screen was still running: wait for s.wg after closing everything")
			default:
			}
			if line, err := ingrid.next(); line != hello {
				t.Fatalf("Ingrid's screen, accepted before Shutdown(), read %q, %v, want the greeting %q", line, err, hello)
			}
			if line, err := ingrid.next(); !errors.Is(err, io.EOF) {
				t.Errorf("Ingrid's screen, greeted after the watchers were closed, read %q, %v, want io.EOF: a watcher that joins after CloseAll must be hung up on too", line, err)
			}
			select {
			case <-stopped:
			case <-time.After(time.Minute):
				t.Fatalf("Shutdown() has not returned a minute after it began, while a watcher who joined as it ran was still connected")
			}
		})
	})
	t.Run("the watch port on loopback", func(t *testing.T) {
		dish1, _ := dish(t, NewSim("YST-1", 1))
		dish2, _ := dish(t, NewSim("YST-2", 2))
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}
		t.Cleanup(func() { ln.Close() })
		s := &Service{Dishes: []*Client{dish1, dish2}, Watchers: &Watchers{Timeout: time.Second}}
		s.Start(ln)
		t.Cleanup(func() { within(s.wg.Wait, time.Second) }) // after the watcher hangs up
		c, err := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
		if err != nil {
			t.Fatalf("a watcher could not connect: %v", err)
		}
		t.Cleanup(func() { c.Close() })
		c.SetDeadline(time.Now().Add(time.Second))
		if line, err := bufio.NewReader(c).ReadString('\n'); line != hello {
			t.Fatalf("a watcher read %q, %v, want the greeting %q", line, err, hello)
		}
		if !within(s.Shutdown, time.Second) {
			t.Fatalf("Shutdown() has not returned after 1s, with a watcher connected")
		}
		if c, err := net.DialTimeout("tcp", ln.Addr().String(), time.Second); err == nil {
			c.Close()
			t.Errorf("after Shutdown(), a new watcher could still connect to the watch port: close the listener")
		}
	})
}

// dish serves sim on one end of a pipe and returns a client on the other end that has said
// HELLO, and a channel closed when the simulator stops. Both ends close when the test ends.
func dish(t *testing.T, sim *Sim) (*Client, <-chan struct{}) {
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
	if _, err := c.hello(); err != nil {
		t.Fatalf("hello() = %v, want a WELCOME", err)
	}
	return c, stopped
}

// within runs f in a goroutine and reports whether it returned within limit.
func within(f func(), limit time.Duration) bool {
	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(limit):
		return false
	}
}

// pipeListener is a watch port in memory, for a synctest bubble: each dial hands Accept one
// end of a net.Pipe.
type pipeListener struct {
	conns  chan net.Conn
	done   chan struct{}
	once   sync.Once
	failed atomic.Bool
}

// newPipeListener returns a listener that accepts nothing until a watcher dials.
func newPipeListener() *pipeListener {
	return &pipeListener{conns: make(chan net.Conn), done: make(chan struct{})}
}

// Accept waits for the next watcher. Once the listener is closed it fails with net.ErrClosed;
// a call after that ends the goroutine making it, so a loop that never stops fails the test
// instead of spinning.
func (l *pipeListener) Accept() (net.Conn, error) {
	select {
	case c := <-l.conns:
		return c, nil
	case <-l.done:
	}
	if l.failed.Swap(true) {
		runtime.Goexit()
	}
	return nil, &net.OpError{Op: "accept", Net: "pipe", Err: net.ErrClosed}
}

// Close stops the listener; Accept then fails.
func (l *pipeListener) Close() error {
	l.once.Do(func() { close(l.done) })
	return nil
}

// isClosed reports whether the listener has been closed.
func (l *pipeListener) isClosed() bool {
	select {
	case <-l.done:
		return true
	default:
		return false
	}
}

// Addr returns the listener's address.
func (l *pipeListener) Addr() net.Addr { return pipeAddr{} }

// pipeAddr is the address of a pipeListener.
type pipeAddr struct{}

// Network returns "pipe".
func (pipeAddr) Network() string { return "pipe" }

// String returns "watch port".
func (pipeAddr) String() string { return "watch port" }

// screen is a watcher on a pipeListener: the test's end of its connection.
type screen struct {
	conn net.Conn
	br   *bufio.Reader
}

// dial connects watcher number i, failing the test if nobody accepts it within a second. Its
// connection closes when the test ends.
func (l *pipeListener) dial(t *testing.T, i int) *screen {
	t.Helper()
	dishd, end := net.Pipe()
	t.Cleanup(func() {
		end.Close()
		dishd.Close()
	})
	select {
	case l.conns <- dishd:
	case <-time.After(time.Second):
		t.Fatalf("the watch port accepted no watcher %d within 1s", i)
	}
	return &screen{conn: end, br: bufio.NewReader(end)}
}

// next reads the screen's next line, waiting at most a second.
func (s *screen) next() (string, error) {
	s.conn.SetReadDeadline(time.Now().Add(time.Second))
	return s.br.ReadString('\n')
}
