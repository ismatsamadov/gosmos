package main

import (
	"context"
	"math/rand/v2"
	"regexp"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeClock is a clock the test sets by hand. Safe for many goroutines.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

// Now returns the clock's time.
func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Add moves the clock forward by d.
func (c *fakeClock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// prep is the pass prep of Cirrine-1's pass on 7 February, two minutes before its AOS.
var prep = time.Date(2029, 2, 7, 1, 16, 8, 0, time.UTC)

// newTracer returns a tracer on a fake clock at prep, drawing IDs from rand.NewPCG(1, 2).
func newTracer() (*Tracer, *fakeClock) {
	clock := &fakeClock{now: prep}
	return &Tracer{Now: clock.Now, Rand: rand.New(rand.NewPCG(1, 2))}, clock
}

var (
	traceID = regexp.MustCompile(`^[0-9a-f]{16}$`)
	spanID  = regexp.MustCompile(`^[0-9a-f]{8}$`)
)

// watched is a source of random numbers that counts the draws made while nobody held mu.
type watched struct {
	src  rand.Source
	mu   *sync.Mutex
	free atomic.Int32
}

// Uint64 returns src's next number, counting it when mu was free.
func (w *watched) Uint64() uint64 {
	if w.mu.TryLock() {
		w.mu.Unlock()
		w.free.Add(1)
	}
	return w.src.Uint64()
}

// TestFirstPass starts Cirrine-1's pass span and its worker span, as uplinkd does at pass prep.
func TestFirstPass(t *testing.T) {
	tr, _ := newTracer()
	src := &watched{src: rand.NewPCG(1, 2), mu: &tr.mu}
	tr.Rand = rand.New(src) // the same numbers, watched
	ctx, pass := tr.Start(withPass(context.Background(), "cirrine-1/LAN/17759"), "pass")
	if pass.TraceID != "c4f5a58656eef510" || pass.SpanID != "9dcec3ad" || pass.ParentID != "" {
		t.Errorf(`Start(ctx with no span, "pass") = trace %q, span %q, parent %q, want trace "c4f5a58656eef510", span "9dcec3ad", parent "": the trace ID drawn first`,
			pass.TraceID, pass.SpanID, pass.ParentID)
	}
	if pass.Name != "pass" || pass.PassID != "cirrine-1/LAN/17759" || !pass.Start.Equal(prep) {
		t.Errorf(`Start(ctx with pass cirrine-1/LAN/17759, "pass") = name %q, pass %q, start %v, want "pass", "cirrine-1/LAN/17759", %v`,
			pass.Name, pass.PassID, pass.Start, prep)
	}
	got, ok := spanFrom(ctx)
	if !ok || got != pass {
		t.Fatalf(`spanFrom(the context Start returned) = %v, %v, want the pass span, true`, got, ok)
	}
	_, worker := tr.Start(ctx, "worker")
	if worker.TraceID != pass.TraceID || worker.ParentID != pass.SpanID || worker.SpanID != "c8d04605" {
		t.Errorf(`Start(the pass's context, "worker") = trace %q, span %q, parent %q, want trace %q, span "c8d04605", parent %q`,
			worker.TraceID, worker.SpanID, worker.ParentID, pass.TraceID, pass.SpanID)
	}
	if worker.PassID != "cirrine-1/LAN/17759" {
		t.Errorf(`Start(the pass's context, "worker").PassID = %q, want "cirrine-1/LAN/17759"`, worker.PassID)
	}
	if n := src.free.Load(); n != 0 {
		t.Errorf("two Start calls drew %d times from Rand without holding t.mu, want every draw made holding it", n)
	}
}

// TestNoPass starts a span from a context with neither a span nor a pass ID.
func TestNoPass(t *testing.T) {
	tr, _ := newTracer()
	_, s := tr.Start(context.Background(), "deliver")
	if !traceID.MatchString(s.TraceID) || !spanID.MatchString(s.SpanID) || s.ParentID != "" || s.PassID != "" {
		t.Errorf(`Start(context.Background(), "deliver") = trace %q, span %q, parent %q, pass %q, want 16 and 8 lowercase hex digits, no parent and no pass`,
			s.TraceID, s.SpanID, s.ParentID, s.PassID)
	}
}

// TestEnd ends Cirrine-1's worker span: it took the station's last worker at once.
func TestEnd(t *testing.T) {
	tr, clock := newTracer()
	ctx, pass := tr.Start(withPass(context.Background(), "cirrine-1/LAN/17759"), "pass")
	_, worker := tr.Start(ctx, "worker")
	tr.End(worker)
	clock.Add(11*time.Minute + 5*time.Second) // LOS, 01:27:13
	tr.End(pass)
	pass.Name = "changed after End"
	done := tr.Spans()
	if len(done) != 2 || done[0].Name != "worker" || done[1].Name != "pass" {
		t.Fatalf("Spans() after two End calls = %v, want the worker span, then the pass span", done)
	}
	if d := done[1].End.Sub(done[1].Start); d != 11*time.Minute+5*time.Second {
		t.Errorf("the pass span lasted %v, want 11m5s, from pass prep to LOS", d)
	}
}

// TestTwoGoroutines starts the spans of two passes at once, as a station flying two passes does.
func TestTwoGoroutines(t *testing.T) {
	tr, _ := newTracer()
	ids := []string{"furrow-b/LAN/22379", "furrow-a/LAN/22382"}
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Go(func() {
			for range 50 {
				ctx, pass := tr.Start(withPass(context.Background(), id), "pass")
				_, worker := tr.Start(ctx, "worker")
				if worker.ParentID != pass.SpanID || worker.TraceID != pass.TraceID || worker.PassID != id {
					t.Errorf("the worker span of %s = trace %q, parent %q, pass %q, want trace %q, parent %q, pass %q",
						id, worker.TraceID, worker.ParentID, worker.PassID, pass.TraceID, pass.SpanID, id)
					return
				}
			}
		})
	}
	wg.Wait()
}
