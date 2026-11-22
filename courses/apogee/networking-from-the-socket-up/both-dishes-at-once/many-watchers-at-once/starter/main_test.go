package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"testing/synctest"
	"time"
)

const hello = "dishd watching YST-1 YST-2\n"

func TestServe(t *testing.T) {
	t.Run("three watchers at once", func(t *testing.T) {
		// In a synctest bubble, with the watch port in memory: the guards below pass in no
		// real time, and the test knows when every goroutine is waiting.
		synctest.Test(t, func(t *testing.T) {
			ln := newPipeListener()
			release := make(chan struct{})
			letGo := sync.OnceFunc(func() { close(release) })
			var finished atomic.Int32
			handle := func(c net.Conn) { // greets its watcher, then holds it until the test lets go
				defer c.Close()
				c.SetWriteDeadline(time.Now().Add(time.Second))
				io.WriteString(c, hello)
				<-release
				finished.Add(1)
			}
			var wg sync.WaitGroup
			served := make(chan error, 1)
			done := make(chan struct{})
			go func() {
				defer close(done)
				served <- serve(ln, &wg, handle)
			}()
			t.Cleanup(func() {
				letGo()
				ln.Close()
				<-done
				wg.Wait()
			})

			for i := 1; i <= 3; i++ {
				w := ln.dial(t, i)
				if line, err := w.next(); line != hello {
					t.Fatalf("watcher %d read %q, %v; want the greeting %q while the handlers before it still hold their watchers: run each handler in its own goroutine", i, line, err, hello)
				}
			}

			ln.Close()
			select {
			case err := <-served:
				if err != nil {
					t.Errorf("serve(ln, wg, handle) once ln is closed = %v, want nil: a closed listener is the normal end", err)
				}
			case <-time.After(time.Second):
				t.Fatalf("serve(ln, wg, handle) has not returned a second after ln closed: once Accept fails with net.ErrClosed, return nil")
			}

			waited := make(chan struct{})
			go func() {
				wg.Wait()
				close(waited)
			}()
			synctest.Wait() // every goroutine is waiting now: the handlers on the test, Wait on them
			select {
			case <-waited:
				t.Errorf("wg.Wait() returned while the three handlers were still running: start each one with wg.Go")
			default:
			}
			letGo()
			select {
			case <-waited:
			case <-time.After(time.Second):
				t.Fatalf("wg.Wait() has not returned a second after the handlers were let go")
			}
			if n := finished.Load(); n != 3 {
				t.Errorf("wg.Wait() returned after %d of 3 handlers had finished, want all 3", n)
			}
		})
	})
	t.Run("a loopback listener closes", func(t *testing.T) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}
		var wg sync.WaitGroup
		served := make(chan error, 1)
		go func() { served <- serve(ln, &wg, func(c net.Conn) { c.Close() }) }()
		c, err := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
		if err != nil {
			ln.Close()
			t.Fatalf("a watcher could not connect: %v", err)
		}
		c.SetDeadline(time.Now().Add(time.Second))
		io.Copy(io.Discard, c) // the handler hangs up at once
		c.Close()
		ln.Close()
		select {
		case err := <-served:
			if err != nil {
				t.Errorf("serve(a loopback listener) once it is closed = %v, want nil: a closed listener is the normal end", err)
			}
		case <-time.After(time.Second):
			t.Fatalf("serve(a loopback listener) has not returned a second after it closed")
		}
		wg.Wait()
	})
	t.Run("a listener that fails another way", func(t *testing.T) {
		tooMany := &net.OpError{Op: "accept", Net: "tcp", Err: os.NewSyscallError("accept", syscall.EMFILE)}
		var wg sync.WaitGroup
		served := make(chan error, 1)
		go func() { served <- serve(failingListener{tooMany}, &wg, func(c net.Conn) { c.Close() }) }()
		select {
		case err := <-served:
			if err != tooMany {
				t.Errorf("serve(a listener failing with %q) = %v, want that error as it is", tooMany, err)
			}
		case <-time.After(time.Second):
			t.Fatalf("serve(a listener failing with %q) has not returned after a second; want that error", tooMany)
		}
	})
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
		t.Fatalf("the watch port accepted no watcher %d within 1s while the handlers before it still hold their watchers: run each handler in its own goroutine", i)
	}
	return &screen{conn: end, br: bufio.NewReader(end)}
}

// next reads the screen's next line, waiting at most a second.
func (s *screen) next() (string, error) {
	s.conn.SetReadDeadline(time.Now().Add(time.Second))
	return s.br.ReadString('\n')
}

// failingListener is a listener whose Accept always fails with err.
type failingListener struct{ err error }

// Accept fails with the listener's error.
func (l failingListener) Accept() (net.Conn, error) { return nil, l.err }

// Close does nothing.
func (l failingListener) Close() error { return nil }

// Addr returns a loopback address.
func (l failingListener) Addr() net.Addr { return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)} }
