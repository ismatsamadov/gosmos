package main

import "strings"

// route sends each request to its answer; a request it has no route for gets the list of
// recordings.
func (s *Server) route(req *Request) answer {
	switch {
	case req.Method == "PUT":
		return s.put(req)
	case strings.HasPrefix(req.Target, "/recordings/"):
		return s.recording(req)
	case req.Target == "/feed":
		return feedAnswer()
	}
	return s.list(req)
}
