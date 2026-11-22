package main

import (
	"bufio"
	"net"
	"time"
)

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
