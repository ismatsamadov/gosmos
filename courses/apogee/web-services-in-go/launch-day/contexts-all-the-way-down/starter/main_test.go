package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// slowDisk is the rehearsal's slow disk: a save waits until release is closed or its context ends, and,
// like every store call, gives up with the context's error once the context has ended.
type slowDisk struct {
	entered  chan struct{}
	release  chan struct{}
	canceled chan struct{} // closed when a save sees its context end
	once     sync.Once
}

// newSlowDisk returns a slow disk, released when the test ends at the latest.
func newSlowDisk(t *testing.T) *slowDisk {
	d := &slowDisk{entered: make(chan struct{}, 1), release: make(chan struct{}), canceled: make(chan struct{})}
	t.Cleanup(func() {
		select {
		case <-d.release:
		default:
			close(d.release)
		}
	})
	return d
}

// save saves a booking: it waits for release, or for ctx to end.
func (d *slowDisk) save(ctx context.Context) error {
	d.entered <- struct{}{}
	select {
	case <-d.release:
	case <-ctx.Done():
		d.once.Do(func() { close(d.canceled) })
	}
	return ctx.Err()
}

// app returns the service as the tests run it: POST /bookings saves to the slow disk.
func (d *slowDisk) app() *App {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /bookings", func(w http.ResponseWriter, r *http.Request) {
		if err := d.save(r.Context()); err != nil {
			http.Error(w, "request canceled", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"id":"bk-4851"}`)
	})
	return &App{Service: &Service{}, Handler: mux}
}

// start runs app on a new 127.0.0.1 listener and returns its address and run's result.
func start(t *testing.T, ctx context.Context, app *App, limits Timeouts) (string, <-chan error) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	done := make(chan error, 1)
	go func() { done <- run(ctx, ln, app, limits) }()
	return ln.Addr().String(), done
}

// inFlight sends Furrowlight's POST /bookings to addr in the background and passes on its status, or 0 if no
// answer came.
func inFlight(addr string) <-chan int {
	status := make(chan int, 1)
	go func() {
		booking := strings.NewReader(`{"sat_id":"furrow-a","station":"SQT","aos":"2028-02-01T14:02:10Z","los":"2028-02-01T14:11:52Z"}`)
		resp, err := (&http.Client{Timeout: 2 * time.Second}).Post("http://"+addr+"/bookings", "application/json", booking)
		if err != nil {
			status <- 0
			return
		}
		resp.Body.Close()
		status <- resp.StatusCode
	}()
	return status
}

// within waits up to a second for something on ch, and fails the test if nothing comes.
func within[T any](t *testing.T, ch <-chan T, what string) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(time.Second):
		t.Fatalf("waited a second for %s", what)
		var zero T
		return zero
	}
}

func TestStuckRequestStops(t *testing.T) {
	disk := newSlowDisk(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr, done := start(t, ctx, disk.app(), Timeouts{Shutdown: 100 * time.Millisecond})
	inFlight(addr)
	within(t, disk.entered, "POST /bookings to reach the disk")
	cancel()
	err := within(t, done, "run to return after the 100ms shutdown limit")
	if !errors.Is(err, context.DeadlineExceeded) || err.Error() != "shutdown: context deadline exceeded" {
		t.Errorf("run with a request stuck past the limit = %v, want %q, which errors.Is finds as context.DeadlineExceeded", err, "shutdown: context deadline exceeded")
	}
	select {
	case <-disk.canceled:
	case <-time.After(time.Second):
		t.Error("the stuck save never saw its context canceled, want it stopped once the shutdown gave up")
	}
}

func TestRequestsThatFinishStillFinish(t *testing.T) {
	disk := newSlowDisk(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr, done := start(t, ctx, disk.app(), Timeouts{Shutdown: time.Second})
	status := inFlight(addr)
	within(t, disk.entered, "POST /bookings to reach the disk")
	cancel()
	for start := time.Now(); ; time.Sleep(time.Millisecond) {
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err != nil {
			break // the shutdown has begun
		}
		conn.Close()
		if time.Since(start) > time.Second {
			t.Fatal("new connections were still taken a second after the stop, want the shutdown begun")
		}
	}
	close(disk.release)
	if got := within(t, status, "the request in flight to finish"); got != http.StatusCreated {
		t.Errorf("POST /bookings in flight during the shutdown = %d, want 201: a request that finishes in time keeps its context", got)
	}
	if err := within(t, done, "run to return"); err != nil {
		t.Errorf("run = %v, want nil: nothing outlived the limit", err)
	}
}

func TestNothingSlow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	_, done := start(t, ctx, newSlowDisk(t).app(), Timeouts{Shutdown: time.Second})
	cancel()
	if err := within(t, done, "run to return"); err != nil {
		t.Errorf("run with no request in flight = %v, want nil", err)
	}
}
