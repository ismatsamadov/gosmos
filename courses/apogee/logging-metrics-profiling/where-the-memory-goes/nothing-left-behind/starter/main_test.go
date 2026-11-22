package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

// baseline is a healthy uplinkd's own goroutines: the tests count as Lantern Atoll's gauge does,
// those 14 and whatever the passes leave behind.
const baseline = 14

// radio is a fake radio. A satellite in answering acknowledges at once; any other stays silent
// until the pass's context ends, or until the radio gives up on it when the test ends.
type radio struct {
	answering map[string]bool
	gaveUp    chan struct{} // closed when the test ends, long after anyone's LOS
}

// newRadio returns a fake radio whose satellites answer on the passes named.
func newRadio(t *testing.T, answering ...string) *radio {
	r := &radio{answering: map[string]bool{}, gaveUp: make(chan struct{})}
	for _, id := range answering {
		r.answering[id] = true
	}
	t.Cleanup(func() { close(r.gaveUp) })
	return r
}

// fly flies pass p as uplinkd's radio does now: it returns at the acknowledgment, or with
// ErrNoAck when ctx ends first.
func (r *radio) fly(ctx context.Context, p PassView, lines []string, buf []byte) (Report, error) {
	if r.answering[p.ID] {
		return Report{Frames: buf}, nil
	}
	select {
	case <-ctx.Done():
	case <-r.gaveUp:
	}
	return Report{Frames: buf}, ErrNoAck
}

// hold takes n of u's workers, as passes still flying hold them, and gives them back when the
// test ends.
func hold(t *testing.T, u *Uplink, n int) {
	for range n {
		u.workers <- struct{}{}
	}
	t.Cleanup(func() {
		for range n {
			<-u.workers
		}
	})
}

// mib writes a number of bytes in MiB, to one decimal place.
func mib(n int64) string {
	return fmt.Sprintf("%.1f MiB", float64(n)/(1<<20))
}

// night is Lantern Atoll's night of 6–7 February: Cirrine-1 never answered, and both Furrows
// would have flown on any worker. Each pass's LOS comes after milliseconds here, not minutes:
// long enough that a free worker is always taken before it, even on a busy Mac.
var night = []struct {
	id   string
	los  time.Duration // from the start of runPass
	want string        // runPass's error; "" for a flown pass
}{
	{"cirrine-1/LAN/17759", 200 * time.Millisecond, "pass cirrine-1/LAN/17759: no ack before LOS"},
	{"furrow-b/LAN/22379", 200 * time.Millisecond, ""},
	{"furrow-a/LAN/22382", 200 * time.Millisecond, ""},
}

// TestReplay replays the night as it began, three of the four workers held since the natural
// losses of 19 and 30 January and 4 February.
func TestReplay(t *testing.T) {
	u := newUplink("LAN", newRadio(t, "furrow-b/LAN/22379", "furrow-a/LAN/22382").fly)
	hold(t, u, 3)
	start := inUse()
	errs := make([]error, len(night))
	left := leftBehind(func() {
		for i, p := range night {
			ctx, cancel := context.WithTimeout(context.Background(), p.los)
			errs[i] = u.runPass(ctx, PassView{ID: p.id}, []string{"PING"})
			cancel()
		}
		for i, p := range night {
			err := errs[i]
			switch {
			case p.want == "" && err != nil:
				t.Errorf("runPass(%s, after the failed pass) = %v, want nil: a worker, and the pass flown", p.id, err)
			case p.want != "" && (err == nil || err.Error() != p.want):
				t.Errorf("runPass(%s) = %v, want %q", p.id, err, p.want)
			case p.want != "" && !errors.Is(err, ErrNoAck):
				t.Errorf("errors.Is(runPass(%s), ErrNoAck) = false, want true", p.id)
			}
		}
		if t.Failed() {
			// A wrong runPass can leave goroutines that end only with the test, and
			// leftBehind would wait its whole second for them: stop here.
			t.FailNow()
		}
	})
	grew := int64(inUse()) - int64(start)

	if left != 0 {
		t.Errorf("goroutines after the replay = %d, want %d", baseline+left, baseline)
	}
	if grew >= 4<<20 {
		t.Errorf("heap in use after the replay = %s above the start, want under 4 MiB", mib(grew))
	}
}

// TestNoWorker runs Furrow-B's pass of 01:36:52 into its LOS with all four workers held, as four
// passes still flying would hold them.
func TestNoWorker(t *testing.T) {
	u := newUplink("LAN", newRadio(t, "furrow-b/LAN/22379").fly)
	hold(t, u, 4)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err := u.runPass(ctx, PassView{ID: "furrow-b/LAN/22379"}, []string{"PING"})
	const want = "pass furrow-b/LAN/22379: no free worker before LOS"
	if err == nil || err.Error() != want {
		t.Fatalf("runPass(furrow-b/LAN/22379, four workers held) = %v, want %q", err, want)
	}
	if !errors.Is(err, ErrNoWorker) {
		t.Errorf("errors.Is(runPass(furrow-b/LAN/22379), ErrNoWorker) = false, want true")
	}
}
