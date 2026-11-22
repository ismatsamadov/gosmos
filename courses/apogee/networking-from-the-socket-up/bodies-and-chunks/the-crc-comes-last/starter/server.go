package main

import (
	"bufio"
	"errors"
	"net"
	"os"
	"strconv"
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

// answer is what a route decides; serveConn writes it.
type answer struct {
	status int
	fields []Field // Content-Type, Content-Length and the rest, in order
	body   []byte  // nil for HEAD
	feed   bool    // GET /feed: serveConn streams the feed instead, then hangs up
}

// textAnswer is a plain-text answer: msg and a newline, with its type and length.
func textAnswer(status int, msg string) answer {
	body := []byte(msg + "\n")
	return answer{status: status, body: body, fields: []Field{
		{"Content-Type", "text/plain; charset=utf-8"},
		{"Content-Length", strconv.Itoa(len(body))},
	}}
}

// failure is the answer to a request that couldn't be read or kept: a RequestError's own
// status and message, or 408 when a deadline passed. Any other error has no answer, status
// 0, because the client has gone.
func failure(err error) answer {
	if e, ok := errors.AsType[*RequestError](err); ok {
		return textAnswer(e.Status, e.Msg)
	}
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return textAnswer(408, "request timeout")
	}
	return answer{}
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

// read reads the request whose first byte has arrived: its head within Head, then its body,
// each read of which must bring bytes within Read.
func (s *Server) read(cr *connReader, br *bufio.Reader) (*Request, error) {
	limit(cr.SetReadDeadline, s.Limits.Head)
	req, err := readRequest(br, s.Limits.MaxBody)
	if err != nil {
		return nil, err
	}
	limit(cr.SetReadDeadline, 0)
	cr.each = s.Limits.Read
	defer func() { cr.each = 0 }()
	switch {
	case req.Length > 0:
		req.Body, err = readBody(br, req.Length)
	case req.Length < 0:
		req.Body, req.Trailer, err = readChunked(br, s.Limits.MaxBody)
	}
	if err != nil {
		return nil, err
	}
	return req, nil
}

// connReader is what a connection's bufio.Reader reads through. While each is set, it moves
// the read deadline forward before every read, so each read must bring bytes within each.
type connReader struct {
	net.Conn
	each time.Duration
}

// Read reads from the connection, within each when it is set.
func (r *connReader) Read(p []byte) (int, error) {
	if r.each > 0 {
		r.SetReadDeadline(time.Now().Add(r.each))
	}
	return r.Conn.Read(p)
}

// limit sets a deadline d from now with set, or removes it when d is zero.
func limit(set func(time.Time) error, d time.Duration) {
	if d > 0 {
		set(time.Now().Add(d))
	} else {
		set(time.Time{})
	}
}

// track adds c to the server's connections, waiting for its first request, and counts it in
// s.wg. It reports false, having closed c, once the server has begun to stop.
func (s *Server) track(c net.Conn) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.shutting {
		c.Close()
		return false
	}
	if s.conns == nil {
		s.conns = map[net.Conn]bool{}
	}
	s.conns[c] = false
	s.wg.Add(1)
	return true
}

// forget closes c and removes it from the server's connections.
func (s *Server) forget(c net.Conn) {
	s.mu.Lock()
	delete(s.conns, c)
	s.mu.Unlock()
	c.Close()
}

// busy marks c as busy with a request.
func (s *Server) busy(c net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.conns[c]; ok {
		s.conns[c] = true
	}
}
