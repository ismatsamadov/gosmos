package main

import (
	"context"
	"errors"
	"math/rand/v2"
	"net"
	"os"
	"slices"
	"sync"
	"syscall"
	"testing"
	"testing/synctest"
	"time"
)

func TestSend(t *testing.T) {
	stow := Frame{Type: TypeStow}
	t.Run("the line drops after the STOW ran", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-2", 2)
			sim.LastSeq, sim.DropAfter = 119, 120
			l := newLine(t, sim).connected(t)
			if err := l.Send(t.Context(), stow); err != nil {
				t.Errorf("Send(STOW 120) = %v, want nil: reconnect, and the WELCOME says it ran", err)
			}
			checkRuns(t, sim, []string{"120 STOW"})
			if n := sim.Conns(); n != 2 {
				t.Errorf("the controller's connection count is %d, want 2: one reconnect", n)
			}
		})
	})
	t.Run("the line drops before the STOW arrived", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-2", 2)
			sim.LastSeq, sim.DropBefore = 119, 120
			l := newLine(t, sim).connected(t)
			if err := l.Send(t.Context(), stow); err != nil {
				t.Errorf("Send(STOW 120) = %v, want nil: reconnect and send it again", err)
			}
			checkRuns(t, sim, []string{"120 STOW"})
			if n := sim.Conns(); n != 2 {
				t.Errorf("the controller's connection count is %d, want 2: one reconnect", n)
			}
		})
	})
	t.Run("Dish 2 is stowed", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-2", 2)
			sim.LastSeq, sim.State, sim.ElDeg = 118, 3, 90
			l := newLine(t, sim).connected(t)
			body, err := pointBody(183.25, 41.70)
			if err != nil {
				t.Fatal(err)
			}
			err = l.Send(t.Context(), Frame{Type: TypePoint, Body: body})
			if nak, ok := errors.AsType[*NakError](err); !ok || nak.Seq != 119 || nak.Code != 4 {
				t.Errorf("Send(POINT 119) = %v, want the NakError seq 119 refused: stowed", err)
			}
			if n := sim.Conns(); n != 1 {
				t.Errorf("the controller's connection count is %d, want 1: a NAK is an answer, not a lost line", n)
			}
		})
	})
	t.Run("the context ends while the line is down", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-2", 2)
			sim.LastSeq, sim.DropBefore = 119, 120
			line := newLine(t, sim)
			l := line.connected(t)
			line.downFor(time.Hour)
			ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
			defer cancel()
			start := time.Now()
			if err := l.Send(ctx, stow); !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("Send(STOW 120) = %v, want context.DeadlineExceeded: the caller stopped waiting for the line", err)
			}
			if took := time.Since(start); took != time.Minute {
				t.Errorf("Send returned after %v, want 1m0s, when its context ended", took)
			}
			checkRuns(t, sim, nil)
		})
	})
}

// checkRuns fails the test unless the controller executed exactly want.
func checkRuns(t *testing.T, sim *Sim, want []string) {
	t.Helper()
	if got := sim.Executed(); !slices.Equal(got, want) {
		t.Errorf("the controller executed %q, want %q: a command sent again keeps its number, and runs once", got, want)
	}
}

// line is the island's line in a test. While it is up, each dial reaches the simulator
// through a fresh net.Pipe; while it is down, each dial is refused. The link's waits
// sleep in the bubble's fake time.
type line struct {
	sim       *Sim
	downUntil time.Time
	ends      []net.Conn
	served    sync.WaitGroup
}

// newLine returns a line to sim. When the test ends it closes every connection it made
// and waits for the simulator to finish with each one.
func newLine(t *testing.T, sim *Sim) *line {
	ln := &line{sim: sim}
	t.Cleanup(func() {
		for _, c := range ln.ends {
			c.Close()
		}
		ln.served.Wait()
	})
	return ln
}

// connected returns a link to the simulator over the line, with dishd's 5 s timeout and
// the jitter of a fixed seed, after its first connect.
func (ln *line) connected(t *testing.T) *Link {
	t.Helper()
	l := &Link{Dial: ln.dial, Wait: wait, Rand: rand.New(rand.NewPCG(1, 2)), Timeout: 5 * time.Second}
	if err := l.connect(t.Context()); err != nil {
		t.Fatalf("connect = %v, want nil", err)
	}
	return l
}

// downFor takes the line down: every dial in the next d is refused.
func (ln *line) downFor(d time.Duration) {
	ln.downUntil = time.Now().Add(d)
}

// dial is the link's Dial over this line.
func (ln *line) dial(ctx context.Context) (net.Conn, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if time.Now().Before(ln.downUntil) {
		return nil, &net.OpError{Op: "dial", Net: "tcp", Err: os.NewSyscallError("connect", syscall.ECONNREFUSED)}
	}
	a, b := net.Pipe()
	ln.ends = append(ln.ends, a, b)
	ln.served.Go(func() { ln.sim.ServeConn(b) })
	return a, nil
}

// wait is the link's Wait: it sleeps d, or ends early with ctx's error.
func wait(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
