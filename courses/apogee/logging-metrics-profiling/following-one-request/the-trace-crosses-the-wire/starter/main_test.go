package main

import (
	"context"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
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

// monday is Monday 19 February 2029, the day spans reach every uplinkd.
var monday = time.Date(2029, 2, 19, 0, 0, 0, 0, time.UTC)

var traceID = regexp.MustCompile(`^[0-9a-f]{16}$`)

// arrival is what Lantern Atoll's handler for POST /uplink saw of one request.
type arrival struct {
	header string // the Apogee-Trace header as it arrived
	span   *Span  // the span the handler's context carried; nil for none
}

// station starts Lantern Atoll's side of POST /uplink: a handler, inside continueTrace's "queue"
// span, that notes what each request brought and queues its lines. Its tracer draws its own IDs.
func station(t *testing.T) (*httptest.Server, *Tracer, func() []arrival) {
	t.Helper()
	tr := &Tracer{Now: (&fakeClock{now: monday}).Now, Rand: rand.New(rand.NewPCG(3, 4))}
	var mu sync.Mutex
	var seen []arrival
	queue := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, _ := spanFrom(r.Context())
		mu.Lock()
		seen = append(seen, arrival{r.Header.Get("Apogee-Trace"), s})
		mu.Unlock()
		w.WriteHeader(http.StatusAccepted) // lines queued
	})
	srv := httptest.NewServer(continueTrace(tr, "queue", queue))
	t.Cleanup(srv.Close)
	return srv, tr, func() []arrival {
		mu.Lock()
		defer mu.Unlock()
		return seen
	}
}

// deliver sends req through client, and stops the test unless the station accepted it.
func deliver(t *testing.T, client *http.Client, req *http.Request) {
	t.Helper()
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST /uplink: %v", err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("POST /uplink = %d, want 202: a span never refuses a request", res.StatusCode)
	}
}

// TestAcrossTheWire sends Cirrine-1's lines from headquarters, inside its "deliver" span.
func TestAcrossTheWire(t *testing.T) {
	srv, tr, seen := station(t)
	hq := &Tracer{Now: (&fakeClock{now: monday}).Now, Rand: rand.New(rand.NewPCG(1, 2))}
	ctx, span := hq.Start(withPass(context.Background(), "cirrine-1/LAN/17759"), "deliver")
	client := &http.Client{Transport: traceTransport{next: srv.Client().Transport}, Timeout: 2 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, srv.URL+"/uplink", strings.NewReader("PING\n"))
	if err != nil {
		t.Fatal(err)
	}
	deliver(t, client, req)

	if got := req.Header.Get("Apogee-Trace"); got != "" {
		t.Errorf("headquarters' own request after Do has Apogee-Trace %q, want none: set it on a clone", got)
	}
	got := seen()
	if len(got) != 1 {
		t.Fatalf("the station's handler saw %d requests, want 1", len(got))
	}
	const want = "c4f5a58656eef510-9dcec3ad"
	if got[0].header != want {
		t.Errorf("Apogee-Trace at the station = %q, want %q: the deliver span's trace ID and span ID", got[0].header, want)
	}
	s := got[0].span
	if s == nil {
		t.Fatalf("the station's handler found no span in its context, want the queue span")
	}
	if s.Name != "queue" || s.TraceID != span.TraceID || s.ParentID != span.SpanID {
		t.Errorf("the station's span = %q, trace %q, parent %q, want \"queue\", trace %q, parent %q",
			s.Name, s.TraceID, s.ParentID, span.TraceID, span.SpanID)
	}
	if done := tr.Spans(); len(done) != 1 || done[0].Name != "queue" {
		t.Errorf("the station's finished spans = %v, want the queue span, ended once the request was served", done)
	}
}

// TestNoSpan sends lines from a context that carries no span: no header goes out, and the station
// starts a trace of its own.
func TestNoSpan(t *testing.T) {
	srv, _, seen := station(t)
	client := &http.Client{Transport: traceTransport{next: srv.Client().Transport}, Timeout: 2 * time.Second}
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/uplink", strings.NewReader("PING\n"))
	if err != nil {
		t.Fatal(err)
	}
	deliver(t, client, req)
	got := seen()
	if len(got) != 1 {
		t.Fatalf("the station's handler saw %d requests, want 1", len(got))
	}
	if got[0].header != "" {
		t.Errorf("Apogee-Trace from a context with no span = %q, want no header", got[0].header)
	}
	if s := got[0].span; s == nil || !traceID.MatchString(s.TraceID) || s.ParentID != "" {
		t.Errorf("the station's span for a request with no header = %+v, want a new trace with no parent", s)
	}
}

// TestBadHeaders sends Apogee-Trace headers that name no usable span.
func TestBadHeaders(t *testing.T) {
	for _, h := range []string{
		"c4f5a58656eef510",
		"C4F5A58656EEF510-9DCEC3AD",
		"c4f5a58656eef510-9dcec3ad-01",
		"c4f5a586-9dcec3ad",
		"c4f5a58656eef510-",
	} {
		t.Run(h, func(t *testing.T) {
			srv, _, seen := station(t)
			req, err := http.NewRequest(http.MethodPost, srv.URL+"/uplink", strings.NewReader("PING\n"))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Apogee-Trace", h)
			deliver(t, srv.Client(), req)
			got := seen()
			if len(got) != 1 {
				t.Fatalf("the station's handler saw %d requests, want 1", len(got))
			}
			if s := got[0].span; s == nil || !traceID.MatchString(s.TraceID) || s.ParentID != "" {
				t.Errorf("the station's span for Apogee-Trace %q = %+v, want a new trace with no parent", h, s)
			}
		})
	}
}
