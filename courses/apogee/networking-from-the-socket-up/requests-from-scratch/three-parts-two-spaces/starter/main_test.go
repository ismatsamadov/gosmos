package main

import (
	"errors"
	"testing"
)

func TestParseRequestLine(t *testing.T) {
	tests := []struct {
		name           string
		line           string
		method, target string
		status         int
		msg            string
	}{
		{"the recorder's list request", "GET /recordings HTTP/1.1", "GET", "/recordings", 0, ""},
		{"an upload", "PUT /recordings/furrow-a/YST/18218 HTTP/1.1", "PUT", "/recordings/furrow-a/YST/18218", 0, ""},
		{"a method in lower case, checked later", "get /recordings HTTP/1.1", "get", "/recordings", 0, ""},
		{"two spaces", "GET  /recordings HTTP/1.1", "", "", 400, "bad request line"},
		{"a tab", "GET\t/recordings HTTP/1.1", "", "", 400, "bad request line"},
		{"a space at the end", "GET /recordings HTTP/1.1 ", "", "", 400, "bad request line"},
		{"no version", "GET /recordings", "", "", 400, "bad request line"},
		{"a target without its slash", "GET recordings HTTP/1.1", "", "", 400, "bad request line"},
		{"a version cut short", "GET /recordings HTTP/1", "", "", 400, "bad request line"},
		{"an empty line", "", "", "", 400, "bad request line"},
		{"HTTP/1.0", "GET /recordings HTTP/1.0", "", "", 505, `version "HTTP/1.0" not supported`},
		{"HTTP/2.0", "GET /recordings HTTP/2.0", "", "", 505, `version "HTTP/2.0" not supported`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method, target, err := parseRequestLine(tt.line)
			if tt.status == 0 {
				if method != tt.method || target != tt.target || err != nil {
					t.Errorf("parseRequestLine(%q) = %q, %q, %v, want %q, %q, nil", tt.line, method, target, err, tt.method, tt.target)
				}
				return
			}
			if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != tt.status || e.Msg != tt.msg {
				t.Errorf("parseRequestLine(%q) = %q, %q, %v, want the RequestError %d %s", tt.line, method, target, err, tt.status, tt.msg)
			}
		})
	}
}
