package main

import (
	"errors"
	"testing"
)

func TestContentLength(t *testing.T) {
	tests := []struct {
		name   string
		method string
		values []string // the Content-Length lines, nil for none
		want   int64
		status int
		msg    string
	}{
		{"an upload of 4096 bytes", "PUT", []string{"4096"}, 4096, 0, ""},
		{"leading zeros", "PUT", []string{"0004096"}, 4096, 0, ""},
		{"exactly the limit", "PUT", []string{"65536"}, 65536, 0, ""},
		{"an empty upload", "PUT", []string{"0"}, 0, 0, ""},
		{"a letter", "PUT", []string{"12a"}, 0, 400, `bad content-length "12a"`},
		{"a minus sign", "PUT", []string{"-1"}, 0, 400, `bad content-length "-1"`},
		{"a plus sign", "PUT", []string{"+5"}, 0, 400, `bad content-length "+5"`},
		{"a space inside", "PUT", []string{"4 096"}, 0, 400, `bad content-length "4 096"`},
		{"an empty value", "PUT", []string{""}, 0, 400, `bad content-length ""`},
		{"a list in one field", "PUT", []string{"4096, 4096"}, 0, 400, `bad content-length "4096, 4096"`},
		{"one over the limit", "PUT", []string{"65537"}, 0, 413, "body over 65536 bytes"},
		{"20 digits", "PUT", []string{"12345678901234567890"}, 0, 413, "body over 65536 bytes"},
		{"two Content-Length lines", "PUT", []string{"4096", "4096"}, 0, 400, "more than one content-length"},
		{"GET with a body", "GET", []string{"5"}, 0, 400, "GET has no body"},
		{"HEAD with a body", "HEAD", []string{"5"}, 0, 400, "HEAD has no body"},
		{"GET with 0", "GET", []string{"0"}, 0, 0, ""},
		{"GET with none", "GET", nil, 0, 0, ""},
		{"an upload with none", "PUT", nil, 0, 411, "length required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Header{}
			if tt.values != nil {
				h["content-length"] = tt.values
			}
			got, err := contentLength(tt.method, h, 65536)
			if tt.status == 0 {
				if got != tt.want || err != nil {
					t.Errorf("contentLength(%s, %q, 65536) = %d, %v, want %d, nil", tt.method, tt.values, got, err, tt.want)
				}
				return
			}
			if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != tt.status || e.Msg != tt.msg {
				t.Errorf("contentLength(%s, %q, 65536) = %d, %v, want the RequestError %d %s", tt.method, tt.values, got, err, tt.status, tt.msg)
			}
		})
	}
	t.Run("a max of 0 means no limit", func(t *testing.T) {
		h := Header{"content-length": {"17179869184"}} // 16 GiB, the production limit
		if got, err := contentLength("PUT", h, 0); got != 17179869184 || err != nil {
			t.Errorf("contentLength(PUT, [\"17179869184\"], 0) = %d, %v, want 17179869184, nil", got, err)
		}
	})
}
