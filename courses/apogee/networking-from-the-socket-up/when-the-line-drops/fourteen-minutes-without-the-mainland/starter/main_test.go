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

func TestDrain(t *testing.T) {
	t.Run("fourteen minutes without the mainland", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-2", 2)
			sim.LastSeq, sim.State, sim.ElDeg = 118, 3, 90 // stowed since the storm
			// The jog runs, and its ACK is lost with the line.
			sim.DropAfter = 121
			line := newLine(t, sim)
			l := line.connected(t)
			line.downFor(14 * time.Minute)
			outcomes, err := l.Drain(t.Context(), calibration(t))
			if len(outcomes) != 4 || slices.ContainsFunc(outcomes, isErr) || err != nil {
				t.Errorf("Drain(UNSTOW, POINT, JOG, STOW) = %v, %v, want 4 nil outcomes and nil", outcomes, err)
			}
			want := []string{"119 UNSTOW", "120 POINT 183.25 41.70", "121 JOG A +0.35", "122 STOW"}
			if got := sim.Executed(); !slices.Equal(got, want) {
				t.Errorf("the controller executed %q, want %q: each command once, in order", got, want)
			}
			var total time.Duration
			for i, d := range line.waits {
				if d > 30*time.Second {
					t.Errorf("wait %d = %v, want at most 30s", i+1, d)
				}
				total += d
			}
			if total < 14*time.Minute {
				t.Errorf("the waits add up to %v, want at least 14m0s: the line was out that long", total)
			}
			if n := sim.Conns(); n != 2 {
				t.Errorf("the controller saw %d connections, want 2: one before the line went, one after", n)
			}
		})
	})
	t.Run("the wind at 81 km/h", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-2", 2)
			sim.LastSeq, sim.Wind = 118, 81 // stowed and locked by the wind
			l := newLine(t, sim).connected(t)
			outcomes, err := l.Drain(t.Context(), calibration(t))
			if len(outcomes) != 4 || err != nil {
				t.Fatalf("Drain(UNSTOW, POINT, JOG, STOW) = %v, %v, want 4 outcomes and nil: a refusal is an outcome, and the queue goes on", outcomes, err)
			}
			for i, seq := range []uint32{119, 120, 121} {
				nak, ok := errors.AsType[*NakError](outcomes[i])
				if !ok || nak.Seq != seq || nak.Code != 3 {
					t.Errorf("outcome %d = %v, want the NakError seq %d refused: wind lock", i+1, outcomes[i], seq)
				}
			}
			if outcomes[3] != nil {
				t.Errorf("outcome 4 = %v, want nil: STOW runs in any wind", outcomes[3])
			}
			if got, want := sim.Executed(), []string{"122 STOW"}; !slices.Equal(got, want) {
				t.Errorf("the controller executed %q, want %q", got, want)
			}
		})
	})
	t.Run("the context ends at minute 5", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim := NewSim("YST-2", 2)
			sim.LastSeq, sim.State, sim.ElDeg = 118, 3, 90
			sim.DropAfter = 121
			line := newLine(t, sim)
			l := line.connected(t)
			line.downFor(14 * time.Minute)
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Minute)
			defer cancel()
			start := time.Now()
			outcomes, err := l.Drain(ctx, calibration(t))
			if len(outcomes) != 2 || slices.ContainsFunc(outcomes, isErr) || !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("Drain(UNSTOW, POINT, JOG, STOW) = %v, %v, want the 2 nil outcomes so far and context.DeadlineExceeded", outcomes, err)
			}
			if took := time.Since(start); took != 5*time.Minute {
				t.Errorf("Drain returned after %v, want 5m0s, when its context ended", took)
			}
			want := []string{"119 UNSTOW", "120 POINT 183.25 41.70", "121 JOG A +0.35"}
			if got := sim.Executed(); !slices.Equal(got, want) {
				t.Errorf("the controller executed %q, want %q", got, want)
			}
		})
	})
}

// calibration is Friday's queue for Dish 2: out of stow, point, Tomás's jog, stow again.
func calibration(t *testing.T) []Frame {
	t.Helper()
	point, err := pointBody(183.25, 41.70)
	if err != nil {
		t.Fatal(err)
	}
	jog, err := jogBody('A', 0.35)
	if err != nil {
		t.Fatal(err)
	}
	return []Frame{{Type: TypeUnstow}, {Type: TypePoint, Body: point}, {Type: TypeJog, Body: jog}, {Type: TypeStow}}
}

// isErr reports whether an outcome is an error.
func isErr(err error) bool { return err != nil }

// line is the island's line in a test. While it is up, each dial reaches the simulator
// through a fresh net.Pipe; while it is down, each dial is refused. It records every
// wait of the link and sleeps it in the bubble's fake time.
type line struct {
	sim       *Sim
	downUntil time.Time
	waits     []time.Duration
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
	l := &Link{Dial: ln.dial, Wait: ln.wait, Rand: rand.New(rand.NewPCG(1, 2)), Timeout: 5 * time.Second}
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

// wait is the link's Wait: it records d and sleeps it, or ends early with ctx.
func (ln *line) wait(ctx context.Context, d time.Duration) error {
	ln.waits = append(ln.waits, d)
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
