package main

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
)

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

// put keeps an uploaded recording under the pass ID in its path, once, after checking its CRC.
func (s *Server) put(req *Request) answer {
	id, ok := strings.CutPrefix(req.Target, "/recordings/")
	if !ok || !validPassID(id) {
		return textAnswer(404, "no such path")
	}
	crc, err := checkCRC(req.Body, req.Header, req.Trailer)
	if err != nil {
		return failure(err)
	}
	if err := s.Store.Put(id, req.Body); err != nil {
		return textAnswer(409, err.Error())
	}
	a := textAnswer(201, fmt.Sprintf("stored %s %d bytes", id, len(req.Body)))
	a.fields = append(a.fields, Field{"Location", req.Target}, Field{"Apogee-Crc32", fmt.Sprintf("%08x", crc)})
	return a
}

// recording answers GET of one recording with its bytes, and HEAD with the same status and
// fields and no body.
func (s *Server) recording(req *Request) answer {
	id := strings.TrimPrefix(req.Target, "/recordings/")
	a := textAnswer(404, "no recording "+id)
	if data, ok := s.Store.Get(id); ok {
		a = answer{status: 200, body: data, fields: []Field{
			{"Content-Type", "application/octet-stream"},
			{"Content-Length", strconv.Itoa(len(data))},
			{"Apogee-Crc32", fmt.Sprintf("%08x", crc32.ChecksumIEEE(data))},
		}}
	}
	return a
}
