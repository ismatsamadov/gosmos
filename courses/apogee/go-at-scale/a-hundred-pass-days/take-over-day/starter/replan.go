package main

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

// Replanner turns news into planning runs and publishes each run's
// schedule, if it is new.
type Replanner struct {
	n     *Network
	pub   *Publisher
	cfg   Config
	now   func() time.Time
	runs  atomic.Uint64 // the number of the last run started
	stows []Stow        // every stow heard of; the Run loop alone changes it
	news  chan struct{} // holds a value once a new schedule is published
	// plan makes a run's schedule: replanFrom within the run limit, or a
	// test's fake.
	plan func(ctx context.Context, cur *Schedule, at time.Time, stows []Stow) (*Schedule, error)
}

// newReplanner returns a Replanner whose runs come after the
// publisher's current schedule.
func newReplanner(n *Network, pub *Publisher, cfg Config, now func() time.Time) *Replanner {
	r := &Replanner{n: n, pub: pub, cfg: cfg, now: now, news: make(chan struct{}, 1)}
	if cur := pub.Current(); cur != nil {
		r.runs.Store(cur.Version)
	}
	r.plan = func(ctx context.Context, cur *Schedule, at time.Time, stows []Stow) (*Schedule, error) {
		return planWithin(ctx, cfg.RunLimit, func(ctx context.Context) (*Schedule, error) {
			return replanFrom(ctx, n, cur, at, stows, cfg)
		})
	}
	return r
}

// publish offers s as run's schedule, made at made, unless its contacts
// are the current schedule's; once s is published, news says so.
func (r *Replanner) publish(s *Schedule, run uint64, made time.Time) {
	if cur := r.pub.Current(); cur != nil && slices.EqualFunc(s.Contacts, cur.Contacts, sameContact) {
		return // nothing new
	}
	s.Version, s.Made = run, made
	if r.pub.Offer(s) {
		select {
		case r.news <- struct{}{}:
		default: // news is already waiting
		}
	}
}

// sameContact reports whether a and b are one pass on one dish.
func sameContact(a, b Contact) bool { return a.ID == b.ID && a.Dish == b.Dish }

// Run plans at every event and every r.cfg.Replan between them, until
// events closes or ctx ends. News cancels the run in progress, and each
// run takes its number as it starts, so the newest news wins.
func (r *Replanner) Run(ctx context.Context, events <-chan Event) error {
	tick := time.NewTicker(r.cfg.Replan)
	defer tick.Stop()
	var runs sync.WaitGroup
	stop := context.CancelFunc(func() {}) // cancels the run in progress
	defer func() {
		runs.Wait()
		stop()
	}()
	for {
		var ev Event
		select {
		case e, ok := <-events:
			if !ok {
				return nil
			}
			ev = e
		case <-tick.C:
			ev = Event{At: r.now()}
		case <-ctx.Done():
			return ctx.Err()
		}
		if ev.Stow != nil {
			r.stows = heard(r.stows, *ev.Stow)
		}
		stop() // the run in progress is out of date
		runCtx, cancel := context.WithCancel(ctx)
		stop = cancel
		run, made, cur, stows := r.runs.Add(1), r.now(), r.pub.Current(), r.stows
		runs.Go(func() {
			if s, err := r.plan(runCtx, cur, ev.At, stows); err == nil {
				r.publish(s, run, made)
			}
		})
	}
}

// heard returns stows with st added: a stow with a zero To begins at
// its station, and one with a To ends the station's stow in progress.
// It leaves stows as it was, since runs in flight still read it.
func heard(stows []Stow, st Stow) []Stow {
	stows = slices.Clone(stows)
	open := func(s Stow) bool { return s.Station == st.Station && s.To.IsZero() }
	if i := slices.IndexFunc(stows, open); i >= 0 && !st.To.IsZero() {
		stows[i].To = st.To
		return stows
	}
	return append(stows, st)
}

// takeOver runs passplan 2 from the cut-over at until ctx ends: run 1
// replans old from at, a Replanner replans on every event and every
// cfg.Replan, and each contact of the newest schedule flies at its AOS.
func takeOver(ctx context.Context, n *Network, old *Schedule, at time.Time, events <-chan Event, cfg Config, now func() time.Time, fly func(Contact)) error {
	return nil
}
