package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// landed is the event for Furrow-A's image from Monday's pass over Salar Quieto.
var landed = Event{
	ID:        "evt-000001",
	Type:      "data.landed",
	BookingID: "bk-4718",
	SatID:     "furrow-a",
	Station:   "SQT",
	File:      "furrow-a-20280117-1416.tif",
	Bytes:     10875000000,
	LandedAt:  time.Date(2028, 1, 17, 14, 27, 52, 0, time.UTC),
}

// secret is Furrowlight's test webhook secret, and ts the time of the attempt, 14:27:53 UTC.
var (
	secret = []byte("furrowlight webhook secret")
	ts     = int64(1831732073)
)

// arrival is what a receiver saw of one request.
type arrival struct {
	method, contentType, body string
}

// newReceiver starts a receiver that answers status to every request and passes on what it saw.
func newReceiver(t *testing.T, status int) (*httptest.Server, <-chan arrival) {
	seen := make(chan arrival, 10)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		seen <- arrival{r.Method, r.Header.Get("Content-Type"), string(body)}
		w.WriteHeader(status)
		io.WriteString(w, http.StatusText(status))
	}))
	t.Cleanup(srv.Close)
	return srv, seen
}

func TestPostSendsTheEvent(t *testing.T) {
	srv, seen := newReceiver(t, http.StatusNoContent)
	if err := post(context.Background(), srv.Client(), srv.URL+"/hooks/apogee", secret, ts, landed); err != nil {
		t.Fatalf("post(evt-000001) = %v, want nil", err)
	}
	var got arrival
	select {
	case got = <-seen:
	default:
		t.Fatal("post(evt-000001) returned nil, but the receiver got no request")
	}
	want, _ := json.Marshal(landed)
	if got.method != http.MethodPost {
		t.Errorf("post(evt-000001) sent method %q, want %q", got.method, http.MethodPost)
	}
	if got.contentType != "application/json" {
		t.Errorf("post(evt-000001) sent Content-Type %q, want %q", got.contentType, "application/json")
	}
	if got.body != string(want) {
		t.Errorf("post(evt-000001) sent body %q, want %q", got.body, want)
	}
}

func TestPostCountsEvery2xx(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"Furrowlight's receiver says OK", http.StatusOK},
		{"a receiver that queues events", http.StatusAccepted},
		{"a receiver that says nothing", http.StatusNoContent},
		{"an address nobody serves", http.StatusNotFound},
		{"a receiver that crashed", http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, _ := newReceiver(t, tt.status)
			err := post(context.Background(), srv.Client(), srv.URL, secret, ts, landed)
			call := "post(answer " + strconv.Itoa(tt.status) + ")"
			if tt.status <= 299 {
				if err != nil {
					t.Errorf("%s = %v, want nil: any 2xx means the event arrived", call, err)
				}
				return
			}
			answer, ok := errors.AsType[*AnswerError](err)
			if !ok {
				t.Fatalf("%s = %v, want an *AnswerError", call, err)
			}
			if answer.StatusCode != tt.status {
				t.Errorf("%s gave StatusCode %d, want %d", call, answer.StatusCode, tt.status)
			}
			if want := strconv.Itoa(tt.status) + " " + http.StatusText(tt.status); !strings.Contains(err.Error(), want) {
				t.Errorf("%s = %q, want an error naming %q", call, err.Error(), want)
			}
		})
	}
}

// closeCounter wraps a transport and counts the answer bodies it hands out and the ones closed.
type closeCounter struct {
	inner          http.RoundTripper
	mu             sync.Mutex
	opened, closed int
}

// RoundTrip passes req on and wraps the answer's body to count its Close.
func (c *closeCounter) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.inner.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	c.opened++
	c.mu.Unlock()
	resp.Body = &countedBody{ReadCloser: resp.Body, counter: c}
	return resp, nil
}

// counts returns how many answer bodies were handed out and how many were closed.
func (c *closeCounter) counts() (opened, closed int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.opened, c.closed
}

// countedBody is an answer body that tells its closeCounter when it is closed.
type countedBody struct {
	io.ReadCloser
	counter *closeCounter
	once    sync.Once
}

// Close closes the body and counts it, once.
func (b *countedBody) Close() error {
	b.once.Do(func() {
		b.counter.mu.Lock()
		b.counter.closed++
		b.counter.mu.Unlock()
	})
	return b.ReadCloser.Close()
}

func TestPostClosesEveryBody(t *testing.T) {
	for _, status := range []int{http.StatusNoContent, http.StatusInternalServerError} {
		srv, _ := newReceiver(t, status)
		counter := &closeCounter{inner: srv.Client().Transport}
		c := &http.Client{Transport: counter, Timeout: 2 * time.Second}
		post(context.Background(), c, srv.URL, secret, ts, landed)
		if opened, closed := counter.counts(); opened != 1 || closed != 1 {
			t.Errorf("post(answer %d) closed %d of %d answer bodies, want 1 of 1", status, closed, opened)
		}
	}
}

func TestPostWithNoReceiver(t *testing.T) {
	srv, _ := newReceiver(t, http.StatusNoContent)
	srv.Close()
	err := post(context.Background(), &http.Client{Timeout: 2 * time.Second}, srv.URL, secret, ts, landed)
	if err == nil {
		t.Fatal("post(to a closed receiver) = nil, want an error")
	}
	if _, ok := errors.AsType[*AnswerError](err); ok {
		t.Errorf("post(to a closed receiver) = %v, want a connection error, not an *AnswerError", err)
	}
}

func TestPostKeepsTheContext(t *testing.T) {
	srv, seen := newReceiver(t, http.StatusNoContent)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := post(ctx, srv.Client(), srv.URL, secret, ts, landed); !errors.Is(err, context.Canceled) {
		t.Errorf("post(a canceled context) = %v, want context.Canceled", err)
	}
	select {
	case <-seen:
		t.Error("post(a canceled context) reached the receiver, want nothing sent")
	default:
	}
}
