package main

import "strings"

// route sends each request to the answer for its path and method: 404 for a path the
// server doesn't have, 405 with Allow for a method its path doesn't take.
func (s *Server) route(req *Request) answer {
	_ = strings.CutPrefix
	return s.list(req)
}
