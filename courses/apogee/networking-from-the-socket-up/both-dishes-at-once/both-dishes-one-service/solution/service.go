package main

import (
	"context"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Service is dishd: its dish clients and its watch port.
type Service struct {
	Dishes   []*Client // Dish 1 first
	Watchers *Watchers
	ln       net.Listener
	wg       sync.WaitGroup
	once     sync.Once
}

// Start serves the watch port on ln: every watcher is greeted, added to the set and read from
// until it hangs up.
func (s *Service) Start(ln net.Listener) {
	s.ln = ln
	s.wg.Go(func() {
		if err := serve(ln, &s.wg, s.watch); err != nil {
			log.Printf("dishd: watch port stopped: %v", err)
		}
	})
}

// watch serves one watcher: its greeting, then whatever it sends is read and dropped until it
// hangs up or is closed.
func (s *Service) watch(c net.Conn) {
	names := make([]string, len(s.Dishes))
	for i, d := range s.Dishes {
		names[i] = d.Name
	}
	c.SetWriteDeadline(time.Now().Add(s.Watchers.Timeout))
	if _, err := io.WriteString(c, greeting(names)); err != nil {
		c.Close()
		return
	}
	s.Watchers.Add(c)
	io.Copy(io.Discard, c)
	s.Watchers.Remove(c)
}

// Shutdown stops dishd: it closes the watch port, every watcher and every dish client, then
// waits for every goroutine dishd started. A second call does nothing.
func (s *Service) Shutdown() {
	s.once.Do(func() {
		s.ln.Close()
		s.Watchers.CloseAll()
		for _, d := range s.Dishes {
			d.Close()
		}
		s.wg.Wait()
	})
}

// Run serves the watch port on s.ln and, on every tick, polls every dish and sends each dish's
// line to every watcher. When ctx ends, it shuts dishd down and returns nil.
func (s *Service) Run(ctx context.Context, ticks <-chan time.Time) error {
	s.Start(s.ln)
	for {
		select {
		case <-ctx.Done():
			s.Shutdown()
			return nil
		case <-ticks:
			for _, r := range pollAll(s.Dishes) {
				s.Watchers.Broadcast(watchLine(r))
			}
		}
	}
}
