package main

import "strings"

// route sends each request to the answer for its path and method: 404 for a path the
// server doesn't have, 405 with Allow for a method its path doesn't take.
func (s *Server) route(req *Request) answer {
	id, isRecording := strings.CutPrefix(req.Target, "/recordings/")
	isRecording = isRecording && validPassID(id)
	switch {
	case req.Target == "/recordings" && req.Method != "PUT":
		return s.list(req)
	case req.Target == "/recordings":
		return refuse(req, 405, "method not allowed", Field{"Allow", "GET, HEAD"})
	case isRecording && req.Method == "PUT":
		return s.put(req)
	case isRecording:
		return s.recording(req)
	case req.Target == "/feed" && req.Method == "GET":
		return feedAnswer()
	case req.Target == "/feed":
		return refuse(req, 405, "method not allowed", Field{"Allow", "GET"})
	}
	return refuse(req, 404, "no such path")
}
