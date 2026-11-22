package main

import (
	"fmt"
	"strconv"
	"strings"
)

// route sends each request to its answer; a request it has no route for gets the list of
// recordings.
func (s *Server) route(req *Request) answer {
	switch {
	case req.Method == "PUT":
		return s.put(req)
	}
	return s.list(req)
}

// list answers with one line per kept recording; HEAD gets the same fields and no body.
func (s *Server) list(req *Request) answer {
	var body []byte
	for _, line := range s.Store.List() {
		body = append(body, line+"\n"...)
	}
	a := answer{status: 200, body: body, fields: []Field{
		{"Content-Type", "text/plain; charset=utf-8"},
		{"Content-Length", strconv.Itoa(len(body))},
	}}
	if req.Method == "HEAD" {
		a.body = nil
	}
	return a
}

// put checks a body's CRC, then keeps it under the pass ID in its path and answers 201.
func (s *Server) put(req *Request) answer {
	if _, err := checkCRC(req.Body, req.Header, req.Trailer); err != nil {
		return failure(err)
	}
	id := strings.TrimPrefix(req.Target, "/recordings/")
	if err := s.Store.Put(id, req.Body); err != nil {
		return textAnswer(409, err.Error())
	}
	return textAnswer(201, fmt.Sprintf("stored %s %d bytes", id, len(req.Body)))
}
