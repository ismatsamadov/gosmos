package main

import (
	"bytes"
	"context"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// baseline is a healthy uplinkd's own goroutines: the tests count as Lantern Atoll's gauge does,
// those 14 and whatever the passes leave behind.
const baseline = 14

// lantern holds twenty of Lantern Atoll's passes of 5 and 6 February, from the orbit model, each
// replayed here as a pass whose satellite never answered.
var lantern = []string{
	"cirrine-1/LAN/17729", "furrow-a/LAN/22351", "furrow-b/LAN/22349", "furrow-a/LAN/22352",
	"brine-3/LAN/15748", "pipit-1/LAN/14629", "cirrine-1/LAN/17736", "furrow-a/LAN/22358",
	"furrow-b/LAN/22356", "furrow-a/LAN/22359", "brine-3/LAN/15755", "pipit-1/LAN/14637",
	"cirrine-1/LAN/17744", "furrow-a/LAN/22366", "furrow-b/LAN/22364", "furrow-a/LAN/22367",
	"brine-3/LAN/15763", "brine-3/LAN/15764", "pipit-1/LAN/14644", "pipit-1/LAN/14652",
}

// radio is a fake radio. A satellite in answering acknowledges at once; any other stays silent,
// and the radio gives up on it only when the test says, long after runPass has given up at LOS.
type radio struct {
	answering map[string]bool
	mu        sync.Mutex
	gaveUp    map[string]chan struct{} // closed when the radio gives up on a pass
}

// newRadio returns a fake radio whose satellites answer on the passes named.
func newRadio(answering ...string) *radio {
	r := &radio{answering: map[string]bool{}, gaveUp: map[string]chan struct{}{}}
	for _, id := range answering {
		r.answering[id] = true
	}
	return r
}

// silence returns the channel that closes when the radio gives up on pass id.
func (r *radio) silence(id string) chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.gaveUp[id] == nil {
		r.gaveUp[id] = make(chan struct{})
	}
	return r.gaveUp[id]
}

// fly flies pass p as uplinkd's radio does: it returns at the acknowledgment, or with ErrNoAck
// once it gives up.
func (r *radio) fly(p PassView, lines []string, buf []byte) (Report, error) {
	if r.answering[p.ID] {
		return Report{Frames: buf}, nil
	}
	<-r.silence(p.ID)
	return Report{Frames: buf}, ErrNoAck
}

// giveUp makes the radio give up on pass id, as it does 15 minutes after the pass's last line.
func (r *radio) giveUp(id string) { close(r.silence(id)) }

// afterLOS returns a pass's context once its LOS has passed: already over.
func afterLOS() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	cancel() // the context is over already, and its reason stays the deadline
	return ctx
}

// leaked returns how many goroutines the leak profile reports stuck in a function named name.
func leaked(name string) int {
	var buf bytes.Buffer
	if err := pprof.Lookup("goroutineleak").WriteTo(&buf, 1); err != nil {
		return 0
	}
	total, count := 0, 0
	for _, line := range strings.Split(buf.String(), "\n") {
		if n, _, ok := strings.Cut(line, " @ "); ok {
			count, _ = strconv.Atoi(n)
		} else if strings.HasPrefix(line, "#") && strings.Contains(line, name) {
			total += count
			count = 0 // a stack counts once, however many of its lines match
		}
	}
	return total
}

// TestTwentyFailedPasses fails twenty passes in a row. Each satellite stays silent, runPass
// gives up at LOS, and only then does the radio give up too.
func TestTwentyFailedPasses(t *testing.T) {
	r := newRadio()
	u := newUplink("LAN", r.fly)
	var errs []error
	left := leftBehind(func() {
		for _, id := range lantern {
			errs = append(errs, u.runPass(afterLOS(), PassView{ID: id}, []string{"PING"}))
			r.giveUp(id)
		}
	})
	for i, err := range errs {
		if want := "pass " + lantern[i] + ": context deadline exceeded"; err == nil || err.Error() != want {
			t.Errorf("runPass(%s, after LOS) = %v, want %q", lantern[i], err, want)
		}
	}
	if left != 0 {
		t.Errorf("goroutines after 20 failed passes = %d, want %d", baseline+left, baseline)
	}
	if held := len(u.workers); held != 0 {
		t.Errorf("workers held after 20 failed passes = %d, want 0", held)
	}
	if n := leaked("runPass"); n != 0 {
		t.Errorf("the leak profile after 20 failed passes shows %d goroutines stuck in runPass, want 0", n)
	}
}

// TestFlownPass flies Furrow-A's pass of 02:22:30, which its satellite acknowledged.
func TestFlownPass(t *testing.T) {
	u := newUplink("LAN", newRadio("furrow-a/LAN/22382").fly)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // LOS, as far as the test waits
	defer cancel()
	if err := u.runPass(ctx, PassView{ID: "furrow-a/LAN/22382"}, []string{"PING"}); err != nil {
		t.Errorf("runPass(furrow-a/LAN/22382, acknowledged) = %v, want nil", err)
	}
}
