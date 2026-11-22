package main

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Sampler is a handler that passes, of each level and message, the first 10 records of a UTC minute and then one in a hundred; ERROR records always pass.
type Sampler struct {
	inner  slog.Handler
	shared *sampling // the counts, shared with every Sampler that With makes from this one
}

// sampling holds the counts of a Sampler and of every Sampler made from it.
type sampling struct {
	mu      sync.Mutex
	tallies map[tallyKey]*tally
	total   uint64 // records dropped so far, for log_lines_dropped_total
}

// tallyKey is what records are counted by: their level and their message.
type tallyKey struct {
	level slog.Level
	msg   string
}

// tally counts the records of one level and message in their current UTC minute.
type tally struct {
	minute  time.Time // the UTC minute being counted
	seen    int       // records of that minute so far
	dropped int       // records dropped since one last passed, carried into the next minute
}

// NewSampler returns a Sampler that hands the records it passes to inner.
func NewSampler(inner slog.Handler) *Sampler {
	return &Sampler{inner: inner, shared: &sampling{tallies: map[tallyKey]*tally{}}}
}

// Enabled reports whether inner handles records at level.
func (s *Sampler) Enabled(ctx context.Context, level slog.Level) bool {
	return s.inner.Enabled(ctx, level)
}

// WithAttrs returns a Sampler for inner with attrs added, sharing this one's counts.
func (s *Sampler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Sampler{inner: s.inner.WithAttrs(attrs), shared: s.shared}
}

// WithGroup returns a Sampler for inner with the group opened, sharing this one's counts.
func (s *Sampler) WithGroup(name string) slog.Handler {
	return &Sampler{inner: s.inner.WithGroup(name), shared: s.shared}
}

// Dropped returns how many records the Sampler and the Samplers made from it have dropped.
func (s *Sampler) Dropped() uint64 {
	s.shared.mu.Lock()
	defer s.shared.mu.Unlock()
	return s.shared.total
}

// tallyOf returns the tally of r's level and message, its count started again when r is in a new UTC minute; the caller holds mu.
func (sh *sampling) tallyOf(r slog.Record) *tally {
	k := tallyKey{r.Level, r.Message}
	t := sh.tallies[k]
	if t == nil {
		t = &tally{}
		sh.tallies[k] = t
	}
	if m := r.Time.UTC().Truncate(time.Minute); !m.Equal(t.minute) {
		t.minute, t.seen = m, 0
	}
	return t
}

// Handle counts r and passes it on, or drops it, as the sampling rule says.
func (s *Sampler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level < slog.LevelError {
		sh := s.shared
		sh.mu.Lock()
		t := sh.tallyOf(r)
		t.seen++
		dropped, keep := t.dropped, t.seen <= 10 || t.seen%100 == 10
		if keep {
			t.dropped = 0
		} else {
			t.dropped++
			sh.total++
		}
		sh.mu.Unlock()
		if !keep {
			return nil
		}
		if dropped > 0 {
			r.AddAttrs(slog.Int("dropped", dropped))
		}
	}
	return s.inner.Handle(ctx, r)
}
