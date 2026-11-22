package main

import (
	"bufio"
	"context"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

const (
	tracking = "YST-1 TRACKING az=183.25 el=41.70 wind=64"
	idle     = "YST-2 IDLE az=180.00 el=45.00 wind=64"
)

// Every test runs in a synctest bubble, whose clock is fake: ticks, timeouts and every wait
// for a line keep their real lengths and pass in no real time. The watch port is in memory,
// each watcher on one end of a net.Pipe.

func TestRun(t *testing.T) {
	t.Run("both dishes, two watchers, three ticks", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim1, sim2 := yst1(), yst2()
			d := startDishd(t, sim1, sim2)
			watchers := d.watch(t, 2)
			for tick := 1; tick <= 3; tick++ {
				d.tick(t)
				for i, w := range watchers {
					for _, want := range []string{tracking, idle} {
						if got := w.next(); got != want+"\n" {
							t.Fatalf("tick %d: watcher %d read %q, want %q: on each tick, broadcast every dish's watchLine, Dish 1 first", tick, i+1, got, want+"\n")
						}
					}
				}
			}
			d.stop(t, watchers)
			for _, s := range []*Sim{sim1, sim2} {
				if n := d.closed(s); n != 1 {
					t.Errorf("after Run returned, a simulator counted %d connections closed by dishd, want 1: shut down before returning", n)
				}
			}
		})
	})
	t.Run("Dish 2's controller goes down after the first tick", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			sim1, sim2 := yst1(), yst2()
			d := startDishd(t, sim1, sim2)
			d.hangUpOnSecondQuery(sim2)
			watchers := d.watch(t, 2)
			for tick, lines := range [][]string{{tracking, idle}, {tracking, "YST-2 DOWN closed"}} {
				d.tick(t)
				for i, w := range watchers {
					for _, want := range lines {
						if got := w.next(); got != want+"\n" {
							t.Fatalf("tick %d: watcher %d read %q, want %q: one dish down must not stop the other's line", tick+1, i+1, got, want+"\n")
						}
					}
				}
			}
			d.stop(t, watchers)
			if n := d.closed(sim1); n != 1 {
				t.Errorf("after Run returned, YST-1's simulator counted %d connections closed by dishd, want 1: shut down before returning", n)
			}
		})
	})
}

// yst1 is a simulator for Dish 1, tracking at 183.25° and 41.70°.
func yst1() *Sim {
	sim := NewSim("YST-1", 1)
	sim.State, sim.AzDeg, sim.ElDeg, sim.Wind, sim.LastSeq = 2, 183.25, 41.70, 64, 118
	return sim
}

// yst2 is a simulator for Dish 2, idle at 180.00° and 45.00° in the same wind.
func yst2() *Sim {
	sim := NewSim("YST-2", 2)
	sim.Wind = 64
	return sim
}

// dishd is one run of the service under test.
type dishd struct {
	s       *Service
	ln      *pipeListener
	ticks   chan time.Time
	cancel  context.CancelFunc
	result  chan error // what Run returned
	ends    map[*Sim]net.Conn
	stopped map[*Sim]chan struct{}
}

// startDishd connects a client to each simulator, makes a Service with an in-memory watch
// port and starts Run. When the test ends, everything it opened is closed.
func startDishd(t *testing.T, sims ...*Sim) *dishd {
	t.Helper()
	d := &dishd{ln: newPipeListener(), ticks: make(chan time.Time), result: make(chan error, 1),
		ends: map[*Sim]net.Conn{}, stopped: map[*Sim]chan struct{}{}}
	var clients []*Client
	for _, sim := range sims {
		end, simEnd := net.Pipe()
		stopped := make(chan struct{})
		d.ends[sim], d.stopped[sim] = simEnd, stopped
		go func() {
			sim.ServeConn(simEnd)
			close(stopped)
		}()
		t.Cleanup(func() {
			end.Close()
			simEnd.Close()
			<-stopped
		})
		c, err := newClient(end)
		if err != nil {
			t.Fatalf("newClient(simulator) = %v, want a client", err)
		}
		if _, err := c.hello(); err != nil {
			t.Fatalf("hello() = %v, want a WELCOME", err)
		}
		clients = append(clients, c)
	}
	t.Cleanup(func() { d.ln.Close() })
	d.s = &Service{Dishes: clients, Watchers: &Watchers{Timeout: 2 * time.Second}, ln: d.ln}
	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel
	t.Cleanup(cancel)
	go func() { d.result <- d.s.Run(ctx, d.ticks) }()
	return d
}

// hangUpOnSecondQuery makes sim's controller hang up when it is asked for its status the
// second time, as a rebooting controller would.
func (d *dishd) hangUpOnSecondQuery(sim *Sim) {
	var mu sync.Mutex
	queries := 0
	end := d.ends[sim]
	sim.Before = func(req string) {
		if req != "QUERY" {
			return
		}
		mu.Lock()
		queries++
		second := queries == 2
		mu.Unlock()
		if second {
			end.Close()
		}
	}
}

// watcher is a watcher on dishd's watch port: a reader of its own takes every line it is sent.
type watcher struct {
	lines chan string   // every line read
	gone  chan struct{} // closed when dishd hangs up on the watcher
}

// next returns the watcher's next line, or "" if none comes within a second.
func (w *watcher) next() string {
	select {
	case line := <-w.lines:
		return line
	case <-time.After(time.Second):
		return ""
	}
}

// watch connects n watchers to the watch port and checks each one's greeting. Once every
// goroutine in the bubble is waiting again, dishd has added them all to its set.
func (d *dishd) watch(t *testing.T, n int) []*watcher {
	t.Helper()
	var ws []*watcher
	for i := 1; i <= n; i++ {
		conn := d.ln.dial(t, i)
		w := &watcher{lines: make(chan string, 16), gone: make(chan struct{})}
		go func() {
			defer close(w.gone)
			br := bufio.NewReader(conn)
			for {
				line, err := br.ReadString('\n')
				if err != nil {
					return
				}
				w.lines <- line
			}
		}()
		if got, want := w.next(), "dishd watching YST-1 YST-2\n"; got != want {
			t.Fatalf("watcher %d read %q, want the greeting %q", i, got, want)
		}
		ws = append(ws, w)
	}
	synctest.Wait()
	return ws
}

// tick sends Run one tick, failing the test if Run doesn't take it within a second.
func (d *dishd) tick(t *testing.T) {
	t.Helper()
	select {
	case d.ticks <- time.Now():
	case <-time.After(time.Second):
		t.Fatalf("Run took no tick within 1s: loop, and wait for the next tick with select")
	}
}

// stop cancels Run's context and checks that Run returns nil and every watcher is hung up on.
func (d *dishd) stop(t *testing.T, watchers []*watcher) {
	t.Helper()
	d.cancel()
	select {
	case err := <-d.result:
		if err != nil {
			t.Errorf("Run(ctx, ticks) after ctx was canceled = %v, want nil", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("Run(ctx, ticks) has not returned a second after ctx was canceled")
	}
	for i, w := range watchers {
		select {
		case <-w.gone:
		case <-time.After(time.Second):
			t.Errorf("a second after Run returned, watcher %d was still connected, want it hung up on: call s.Shutdown() before returning", i+1)
		}
	}
}

// closed waits for sim to stop serving, a second at most, and returns how many connections
// it counted closed by dishd.
func (d *dishd) closed(sim *Sim) int {
	select {
	case <-d.stopped[sim]:
	case <-time.After(time.Second):
	}
	return sim.ClosedByClient()
}

// pipeListener is a watch port in memory, for a synctest bubble: each dial hands Accept one
// end of a net.Pipe.
type pipeListener struct {
	conns  chan net.Conn
	done   chan struct{}
	once   sync.Once
	failed atomic.Bool
}

// newPipeListener returns a listener that accepts nothing until a watcher dials.
func newPipeListener() *pipeListener {
	return &pipeListener{conns: make(chan net.Conn), done: make(chan struct{})}
}

// Accept waits for the next watcher. Once the listener is closed it fails with net.ErrClosed;
// a call after that ends the goroutine making it, so a loop that never stops fails the test
// instead of spinning.
func (l *pipeListener) Accept() (net.Conn, error) {
	select {
	case c := <-l.conns:
		return c, nil
	case <-l.done:
	}
	if l.failed.Swap(true) {
		runtime.Goexit()
	}
	return nil, &net.OpError{Op: "accept", Net: "pipe", Err: net.ErrClosed}
}

// Close stops the listener; Accept then fails.
func (l *pipeListener) Close() error {
	l.once.Do(func() { close(l.done) })
	return nil
}

// Addr returns the listener's address.
func (l *pipeListener) Addr() net.Addr { return pipeAddr{} }

// pipeAddr is the address of a pipeListener.
type pipeAddr struct{}

// Network returns "pipe".
func (pipeAddr) Network() string { return "pipe" }

// String returns "watch port".
func (pipeAddr) String() string { return "watch port" }

// dial connects watcher number i, failing the test if nobody accepts it within a second, and
// returns the watcher's end of the connection, which closes when the test ends.
func (l *pipeListener) dial(t *testing.T, i int) net.Conn {
	t.Helper()
	dishd, end := net.Pipe()
	t.Cleanup(func() {
		end.Close()
		dishd.Close()
	})
	select {
	case l.conns <- dishd:
	case <-time.After(time.Second):
		t.Fatalf("the watch port accepted no watcher %d within 1s: start it with s.Start(s.ln)", i)
	}
	return end
}
