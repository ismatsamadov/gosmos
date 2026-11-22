package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"slices"
	"syscall"
	"testing"
	"time"
)

func TestRedial(t *testing.T) {
	t.Run("refused three times, then the line is back", func(t *testing.T) {
		line := newFakeLine(t, 3)
		l := &Link{Dial: line.dial, Wait: line.wait}
		conn, err := l.redial(context.Background())
		if err != nil || conn == nil || conn != line.made(4) {
			t.Fatalf("redial = %s, %v, want the connection of dial 4 and nil", line.which(conn), err)
		}
		if line.dials != 4 {
			t.Errorf("number of dials = %d, want 4", line.dials)
		}
		if want := []time.Duration{time.Second, time.Second, time.Second}; !slices.Equal(line.waits, want) {
			t.Errorf("redial waited %v, want %v: l.Wait(ctx, time.Second) after each refused dial", line.waits, want)
		}
	})
	t.Run("stopped during the second wait", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		line := newFakeLine(t, 100)
		line.stopAt, line.stop = 2, cancel
		l := &Link{Dial: line.dial, Wait: line.wait}
		conn, err := l.redial(ctx)
		if conn != nil || !errors.Is(err, context.Canceled) {
			t.Errorf("redial = %s, %v, want nil and context.Canceled: return the wait's error, so the caller can tell stop from down", line.which(conn), err)
		}
		if line.dials != 2 {
			t.Errorf("number of dials = %d, want 2: no dial after the context ended", line.dials)
		}
	})
	t.Run("the line is up", func(t *testing.T) {
		line := newFakeLine(t, 0)
		l := &Link{Dial: line.dial, Wait: line.wait}
		conn, err := l.redial(context.Background())
		if err != nil || conn == nil || conn != line.made(1) {
			t.Fatalf("redial = %s, %v, want the connection of dial 1 and nil", line.which(conn), err)
		}
		if len(line.waits) != 0 {
			t.Errorf("redial waited %v, want no wait", line.waits)
		}
	})
}

// fakeLine is the island's line in a test. It refuses its first dials, then connects
// each dial through a fresh net.Pipe. Its waits are recorded and return at once, and at
// the wait numbered stopAt it first cancels the caller's context.
type fakeLine struct {
	refusals int
	stopAt   int
	stop     context.CancelFunc
	dials    int
	waits    []time.Duration
	conns    map[int]net.Conn // by dial number, from 1
}

// newFakeLine returns a line that refuses its first refusals dials. When the test ends
// it closes every connection it made.
func newFakeLine(t *testing.T, refusals int) *fakeLine {
	line := &fakeLine{refusals: refusals, conns: map[int]net.Conn{}}
	t.Cleanup(func() {
		for _, c := range line.conns {
			c.Close()
		}
	})
	return line
}

// dial is the link's Dial over this line.
func (fl *fakeLine) dial(ctx context.Context) (net.Conn, error) {
	fl.dials++
	if fl.dials <= fl.refusals {
		return nil, &net.OpError{Op: "dial", Net: "tcp", Err: os.NewSyscallError("connect", syscall.ECONNREFUSED)}
	}
	a, b := net.Pipe()
	b.Close() // nothing answers in this test: the connection itself is what counts
	fl.conns[fl.dials] = a
	return a, nil
}

// wait is the link's Wait: it records d and returns at once, with ctx's error if ctx has
// ended.
func (fl *fakeLine) wait(ctx context.Context, d time.Duration) error {
	fl.waits = append(fl.waits, d)
	if len(fl.waits) == fl.stopAt {
		fl.stop()
	}
	return ctx.Err()
}

// made returns the connection the line made on dial n, or nil.
func (fl *fakeLine) made(n int) net.Conn {
	return fl.conns[n]
}

// which names a connection for a message: the dial that made it, or nil.
func (fl *fakeLine) which(c net.Conn) string {
	if c == nil {
		return "nil"
	}
	for n, made := range fl.conns {
		if made == c {
			return fmt.Sprintf("the connection of dial %d", n)
		}
	}
	return "a connection the line never made"
}
