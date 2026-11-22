package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Outcome is how a contact ended.
type Outcome struct {
	Kind   string // "flown", "failed" or "skipped"; "stopped" when the rehearsal stopped first
	Reason string // why it failed or was skipped
}

// String is the outcome as the LOS line says it: "flown", "failed no_answer".
func (o Outcome) String() string {
	if o.Reason == "" {
		return o.Kind
	}
	return o.Kind + " " + o.Reason
}

// run is one rehearsal of a planned week: everything its contacts share.
type run struct {
	w      World
	n      *Network
	wk     *Week
	health *Health

	mu       sync.Mutex
	outcomes map[string]Outcome // how each contact ended, by pass ID
}

// fly runs every contact of the week, each on its own clock, and returns once every
// contact is done; when ctx ends first it stops them all and returns once they have.
func (r *run) fly(ctx context.Context) {
	var wg sync.WaitGroup
	for _, c := range r.wk.Contacts {
		wg.Go(func() { r.contact(ctx, c) })
	}
	wg.Wait()
}

// contact flies one contact: it waits for its AOS, flies the pass until its LOS and says
// how it ended, at its LOS or at the moment the rehearsal stops.
func (r *run) contact(ctx context.Context, c Contact) {
	if !waitUntil(ctx, r.w, c.AOS) {
		return
	}
	sayAt(r.w, "AOS %s dish %d", c.ID, c.Dish)
	pctx, cancel := context.WithTimeout(ctx, c.LOS.Sub(r.w.Now()))
	defer cancel()
	p := &pass{r: r, c: c}
	p.fly(pctx)
	<-pctx.Done()
	o := p.outcome()
	if ctx.Err() != nil {
		o = Outcome{Kind: "stopped"}
	}
	sayAt(r.w, "LOS %s %s", c.ID, o)
	r.mu.Lock()
	r.outcomes[c.ID] = o
	r.mu.Unlock()
}

// waitUntil waits until the instant at on w's clock; it reports false when ctx ends first.
func waitUntil(ctx context.Context, w World, at time.Time) bool {
	return sleep(ctx, at.Sub(w.Now()))
}

// sleep waits for d; it reports false when ctx ends first.
func sleep(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return ctx.Err() == nil
	case <-ctx.Done():
		return false
	}
}

// sayAt writes one console line that starts with the instant on w's clock, in one Write:
// "Thu 16:03:46 AOS thistle-2/PEL/02136 dish 1".
func sayAt(w World, format string, args ...any) {
	fmt.Fprintf(w.Out, "%s "+format+"\n", append([]any{w.Now().Format("Mon 15:04:05")}, args...)...)
}

// counts returns how many contacts were flown, failed and skipped.
func (r *run) counts() (flown, failed, skipped int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, o := range r.outcomes {
		switch o.Kind {
		case "flown":
			flown++
		case "failed":
			failed++
		case "skipped":
			skipped++
		}
	}
	return flown, failed, skipped
}
