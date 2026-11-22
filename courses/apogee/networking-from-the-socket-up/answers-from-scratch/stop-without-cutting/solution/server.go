package main

import (
	"bufio"
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

// Limits are the station server's limits. Production values in comments; zero means none.
type Limits struct {
	MaxBody int64         // 16 GiB (17179869184): a Furrow-A pass is about 11 GB. Tests: 65536
	Head    time.Duration // 5s: the whole head, from its first byte
	Read    time.Duration // 5s: each read of a body must bring bytes
	Idle    time.Duration // 60s: a kept-alive connection waits this long for the next request
	Write   time.Duration // 30s: each answer
}

// Server is the station server.
type Server struct {
	Store  *Store
	Limits Limits
	Now    func() time.Time

	mu       sync.Mutex
	ln       net.Listener
	conns    map[net.Conn]bool // every open connection: true while it is busy with a request
	wg       sync.WaitGroup    // counts the connections being served
	shutting bool              // the server has begun to stop
}

// Serve accepts connections on ln and serves each one in its own goroutine. Once ln is
// closed, it waits for every connection to end and returns nil.
func (s *Server) Serve(ln net.Listener) error {
	s.mu.Lock()
	s.ln = ln
	if s.shutting {
		ln.Close()
	}
	s.mu.Unlock()
	defer s.wg.Wait()
	for {
		c, err := ln.Accept()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		if err != nil {
			return err
		}
		if s.track(c) {
			go func() {
				defer s.wg.Done()
				s.serveConn(c)
			}()
		}
	}
}

// serveConn answers the requests that arrive on c, one after another, and closes c after
// the last.
func (s *Server) serveConn(c net.Conn) {
	defer s.forget(c)
	cr := &connReader{Conn: c}
	br := bufio.NewReaderSize(cr, 1024)
	for {
		limit(c.SetReadDeadline, s.Limits.Idle)
		if _, err := br.Peek(1); err != nil {
			return // closed, or quiet for Idle: hang up without a word
		}
		s.busy(c)
		req, err := s.read(cr, br)
		if err != nil {
			s.send(c, failure(err), true)
			return
		}
		a := s.route(req)
		if !s.send(c, a, closing(req) || a.status/100 != 2) {
			return
		}
	}
}

// Shutdown stops the server without cutting a request in flight: no new connections, the
// idle ones closed, every feed ended, and each request in flight answered before its
// connection closes. It returns nil once all are done, or ctx's error, having closed
// everything, if ctx ends first.
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.shutting = true
	if s.ln != nil {
		s.ln.Close()
	}
	s.mu.Unlock()
	s.closeConns(false)
	s.Store.Close()
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.closeConns(true)
		<-done
		return ctx.Err()
	}
}
