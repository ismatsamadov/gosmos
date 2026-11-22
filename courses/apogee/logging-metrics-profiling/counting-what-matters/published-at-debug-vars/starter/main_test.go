package main

import (
	"encoding/json"
	"expvar"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

// expvar complains about a name published twice through the log package, which would print
// the Mac's clock in front of it.
func init() { log.SetFlags(0) }

// station is the one Metrics this test binary publishes, made at first use: expvar keeps a
// name for as long as the program runs, so every test, each time it runs, reads these numbers.
var station = sync.OnceValue(func() *Metrics {
	m := &Metrics{AckWait: NewHistogram(0.5, 1, 2, 5, 10, 30, 60, 120, 300)}
	m.Passes.With("brine-3", "failed").Inc()
	m.Passes.With("brine-3", "flown").Add(23)
	m.Passes.With("pipit-1", "flown").Add(12)
	for _, w := range []struct {
		n int
		d time.Duration
	}{
		{3, 900 * time.Millisecond}, {40, 1500 * time.Millisecond}, {15, 3 * time.Second},
		{2, 7 * time.Second}, {1, 21510 * time.Millisecond},
	} {
		for range w.n {
			m.AckWait.Observe(w.d)
		}
	}
	publish(m)
	return m
})

// published returns the station's Metrics, published at first use, and stops the test when
// publish panicked, as it does when two of its variables get the same name.
func published(t *testing.T) *Metrics {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("publish(m) panicked (%s), want three variables, each published once", strings.TrimSpace(fmt.Sprint(r)))
		}
	}()
	return station()
}

// vars answers GET /debug/vars through expvar.Handler, as the debug port does, and splits the
// answer into its variables.
func vars(t *testing.T) map[string]json.RawMessage {
	t.Helper()
	published(t)
	rec := httptest.NewRecorder()
	expvar.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/debug/vars", nil))
	var all map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &all); err != nil {
		t.Fatalf("GET /debug/vars = not JSON (%v), want one JSON object", err)
	}
	return all
}

// passesTotal decodes passes_total from GET /debug/vars.
func passesTotal(t *testing.T) map[string]map[string]uint64 {
	t.Helper()
	raw, ok := vars(t)["passes_total"]
	if !ok {
		t.Fatalf(`GET /debug/vars has no "passes_total", want the passes by satellite and outcome`)
	}
	var got map[string]map[string]uint64
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("passes_total = %s, want counts by satellite and then outcome (%v)", raw, err)
	}
	return got
}

func TestPassesTotal(t *testing.T) {
	got := passesTotal(t)
	for _, c := range []struct {
		sat, outcome string
		want         uint64
	}{
		{"brine-3", "failed", 1}, {"brine-3", "flown", 23}, {"pipit-1", "flown", 12},
	} {
		if n := got[c.sat][c.outcome]; n != c.want {
			t.Errorf("passes_total[%q][%q] = %d, want %d", c.sat, c.outcome, n, c.want)
		}
	}
}

func TestAckWaitSeconds(t *testing.T) {
	raw, ok := vars(t)["ack_wait_seconds"]
	if !ok {
		t.Fatalf(`GET /debug/vars has no "ack_wait_seconds", want the histogram`)
	}
	want := HistogramSnapshot{
		Bounds: []float64{0.5, 1, 2, 5, 10, 30, 60, 120, 300},
		Counts: []uint64{0, 3, 40, 15, 2, 1, 0, 0, 0, 0},
		Count:  61,
		SumMs:  143210,
	}
	var got HistogramSnapshot
	err := json.Unmarshal(raw, &got)
	if err != nil || !slices.Equal(got.Bounds, want.Bounds) || !slices.Equal(got.Counts, want.Counts) ||
		got.Count != want.Count || got.SumMs != want.SumMs {
		b, _ := json.Marshal(want)
		t.Errorf("ack_wait_seconds = %s, want %s", raw, b)
	}
}

func TestGoroutines(t *testing.T) {
	raw, ok := vars(t)["goroutines"]
	if !ok {
		t.Fatalf(`GET /debug/vars has no "goroutines", want the program's goroutine count`)
	}
	var n float64
	if err := json.Unmarshal(raw, &n); err != nil || n < 1 || n != math.Trunc(n) {
		t.Errorf("goroutines = %s, want a whole number of goroutines, at least 1", raw)
	}
}

// TestReadWhenAsked counts one more pass between two reads: a value read at request time
// shows it, and one worked out once, at publish, never would.
func TestReadWhenAsked(t *testing.T) {
	before := passesTotal(t)["cirrine-1"]["flown"]
	published(t).Passes.With("cirrine-1", "flown").Inc()
	if after := passesTotal(t)["cirrine-1"]["flown"]; after != before+1 {
		t.Errorf(`passes_total["cirrine-1"]["flown"] after one more Inc() = %d, want %d`, after, before+1)
	}
}

func TestExpvarsOwn(t *testing.T) {
	all := vars(t)
	for _, name := range []string{"cmdline", "memstats"} {
		if _, ok := all[name]; !ok {
			t.Errorf("GET /debug/vars has no %q, want expvar's own variables kept", name)
		}
	}
}

func TestPublishTwice(t *testing.T) {
	published(t)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("publish(m) a second time panicked (%s), want it harmless", strings.TrimSpace(fmt.Sprint(r)))
		}
	}()
	publish(&Metrics{AckWait: NewHistogram(1, 2)})
}
