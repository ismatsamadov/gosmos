package main

import (
	"bufio"
	"net"
	"strings"
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

// send writes a to c within Write, adding Connection: close when it is the last answer, a
// feed, or the server is stopping; a feed goes on streaming until it ends. It reports
// whether c may wait for another request, and if so marks it as waiting.
func (s *Server) send(c net.Conn, a answer, last bool) bool {
	if a.status == 0 {
		return false // nobody to answer
	}
	last = last || a.feed || s.stopping()
	fields := a.fields
	if last {
		fields = append(fields[:len(fields):len(fields)], Field{"Connection", "close"})
	}
	var err error
	if a.feed {
		err = s.stream(c, fields)
	} else {
		limit(c.SetWriteDeadline, s.Limits.Write)
		err = writeResponse(c, s.Now(), a.status, fields, a.body)
	}
	return err == nil && !last && s.rest(c)
}

// stream answers GET /feed: it subscribes to the store's arrivals, writes the head and hands
// the rest to feed, each write within Write.
func (s *Server) stream(c net.Conn, fields []Field) error {
	lines, unsubscribe := s.Store.Subscribe()
	defer unsubscribe()
	bw := bufio.NewWriter(&connWriter{Conn: c, each: s.Limits.Write})
	if err := writeResponse(bw, s.Now(), 200, fields, nil); err != nil {
		return err
	}
	return s.feed(bw, lines)
}

// connWriter writes to a connection, moving the write deadline forward before every write
// when each is set, so a stream lasts as long as each write keeps up.
type connWriter struct {
	net.Conn
	each time.Duration
}

// Write writes to the connection, within each when it is set.
func (w *connWriter) Write(p []byte) (int, error) {
	limit(w.SetWriteDeadline, w.each)
	return w.Conn.Write(p)
}

// rest marks c as waiting between requests and reports true, or reports false, leaving c
// busy, once the server has begun to stop: then c must close instead.
func (s *Server) rest(c net.Conn) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.shutting {
		return false
	}
	s.conns[c] = false
	return true
}

// stopping reports whether the server has begun to stop.
func (s *Server) stopping() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.shutting
}

// closing reports whether req said Connection: close, the token close in any letter case.
func closing(req *Request) bool {
	for token := range strings.SplitSeq(req.Header.Get("connection"), ",") {
		if strings.EqualFold(strings.TrimSpace(token), "close") {
			return true
		}
	}
	return false
}

// closeConns closes the connections waiting between requests, or every connection when
// all is set.
func (s *Server) closeConns(all bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for c, busy := range s.conns {
		if all || !busy {
			c.Close()
		}
	}
}
