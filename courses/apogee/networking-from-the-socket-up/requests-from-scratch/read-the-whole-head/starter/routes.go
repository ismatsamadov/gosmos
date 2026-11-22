package main

import "strconv"

// route answers every good request with the list of recordings, for now.
func (s *Server) route(req *Request) answer {
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
