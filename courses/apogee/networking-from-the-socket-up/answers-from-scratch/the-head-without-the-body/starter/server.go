package main

import (
	"bufio"
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

// serveConn reads one request from c, with its body, answers it and closes c: one request
// per connection, so every answer says Connection: close.
func (s *Server) serveConn(c net.Conn) {
	defer s.forget(c)
	cr := &connReader{Conn: c}
	br := bufio.NewReaderSize(cr, 1024)
	if _, err := br.Peek(1); err != nil {
		return // the connection ended before a request began
	}
	s.busy(c)
	var a answer
	if req, err := s.read(cr, br); err != nil {
		a = failure(err)
	} else {
		a = s.route(req)
	}
	if a.status == 0 {
		return
	}
	limit(c.SetWriteDeadline, s.Limits.Write)
	writeResponse(c, s.Now(), a.status, append(a.fields, Field{"Connection", "close"}), a.body)
}
