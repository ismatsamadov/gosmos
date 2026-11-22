package main

import (
	"context"
	"math/rand/v2"
	"net"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestRedial(t *testing.T) {
	t.Run("refused six times, then the line is back", func(t *testing.T) {
		line := newFakeLine(t, 6)
		l := &Link{Dial: line.dial, Wait: line.wait, Rand: rand.New(rand.NewPCG(1, 2))}
		if conn, err := l.redial(context.Background()); conn == nil || err != nil {
			t.Fatalf("redial = %v, %v, want the seventh dial's connection and nil", conn, err)
		}
		checkWaits(t, line.waits, rand.New(rand.NewPCG(1, 2)), 1, 6)
	})
	t.Run("the next drop starts from the first delay", func(t *testing.T) {
		line := newFakeLine(t, 6)
		l := &Link{Dial: line.dial, Wait: line.wait, Rand: rand.New(rand.NewPCG(1, 2))}
		if conn, err := l.redial(context.Background()); conn == nil || err != nil {
			t.Fatalf("redial = %v, %v, want the seventh dial's connection and nil", conn, err)
		}
		r := rand.New(rand.NewPCG(1, 2))
		for range 6 {
			r.Float64() // the six draws of the first drop
		}
		line.refuse(1)
		line.waits = nil
		if conn, err := l.redial(context.Background()); conn == nil || err != nil {
			t.Fatalf("redial after a new drop = %v, %v, want a connection and nil", conn, err)
		}
		if len(line.waits) != 1 {
			t.Fatalf("after one more refused dial, redial waited %v, want one wait", line.waits)
		}
		if want := redialDelay(1, r); line.waits[0] != want {
			t.Errorf("after a dial worked, the next drop's first wait = %v, want %v (between 500ms and 1s): set l.attempt back to 0 when a dial works", line.waits[0], want)
		}
	})
	t.Run("sixty refusals in a row", func(t *testing.T) {
		line := newFakeLine(t, 60)
		l := &Link{Dial: line.dial, Wait: line.wait, Rand: rand.New(rand.NewPCG(1, 2))}
		if conn, err := l.redial(context.Background()); conn == nil || err != nil {
			t.Fatalf("redial = %v, %v, want the 61st dial's connection and nil", conn, err)
		}
		for i, d := range line.waits {
			if d > 30*time.Second {
				t.Errorf("wait %d = %v, want at most 30s", i+1, d)
			}
		}
		checkWaits(t, line.waits, rand.New(rand.NewPCG(1, 2)), 1, 60)
	})
}

// checkWaits fails the test, at the first wait that differs, unless waits are
// redialDelay's for attempts from, from+1 ... up to to, drawn in order from r.
func checkWaits(t *testing.T, waits []time.Duration, r *rand.Rand, from, to int) {
	t.Helper()
	if len(waits) != to-from+1 {
		t.Fatalf("redial waited %d times, want %d: once after each refused dial", len(waits), to-from+1)
	}
	for i, got := range waits {
		attempt := from + i
		lo, hi := window(attempt)
		if want := redialDelay(attempt, r); got != want {
			t.Errorf("wait %d = %v, want %v (between %v and %v): redialDelay(%d, l.Rand)", attempt, got, want, lo, hi, attempt)
			return
		}
	}
}

// window is the range of the wait after the attempt-th failed dial in a row: a second
// doubled for each earlier failure, at most 30 s, times a jitter from 0.5 to 1.
func window(attempt int) (lo, hi time.Duration) {
	hi = time.Second
	for i := 1; i < attempt && hi < 30*time.Second; i++ {
		hi *= 2
	}
	hi = min(hi, 30*time.Second)
	return hi / 2, hi
}

// fakeLine is the island's line in a test: it refuses dials until it has refused
// enough, then connects each dial through a fresh net.Pipe. Its waits are recorded and
// return at once.
type fakeLine struct {
	refuseUntil int // dials up to this number are refused
	dials       int
	waits       []time.Duration
	conns       []net.Conn
}

// newFakeLine returns a line that refuses its first refusals dials. When the test ends
// it closes every connection it made.
func newFakeLine(t *testing.T, refusals int) *fakeLine {
	line := &fakeLine{refuseUntil: refusals}
	t.Cleanup(func() {
		for _, c := range line.conns {
			c.Close()
		}
	})
	return line
}

// refuse drops the line again: it refuses the next n dials.
func (fl *fakeLine) refuse(n int) {
	fl.refuseUntil = fl.dials + n
}

// dial is the link's Dial over this line.
func (fl *fakeLine) dial(ctx context.Context) (net.Conn, error) {
	fl.dials++
	if fl.dials <= fl.refuseUntil {
		return nil, &net.OpError{Op: "dial", Net: "tcp", Err: os.NewSyscallError("connect", syscall.ECONNREFUSED)}
	}
	a, b := net.Pipe()
	b.Close() // nothing answers in this test: the connection itself is what counts
	fl.conns = append(fl.conns, a)
	return a, nil
}

// wait is the link's Wait: it records d and returns at once, with ctx's error if ctx has
// ended.
func (fl *fakeLine) wait(ctx context.Context, d time.Duration) error {
	fl.waits = append(fl.waits, d)
	return ctx.Err()
}
