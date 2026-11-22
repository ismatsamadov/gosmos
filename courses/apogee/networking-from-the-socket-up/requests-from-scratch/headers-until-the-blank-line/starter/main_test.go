package main

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestReadHeader(t *testing.T) {
	tests := []struct {
		name string
		head string
		want Header
	}{
		{"Go's own GET", "Host: 127.0.0.1\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n",
			Header{"host": {"127.0.0.1"}, "user-agent": {"Go-http-client/1.1"}, "accept-encoding": {"gzip"}}},
		{"spaces around a value", "Content-Type:   text/plain  \r\n\r\n", Header{"content-type": {"text/plain"}}},
		{"tabs around a value", "X-Pass:\tfurrow-a/YST/18218\t\r\n\r\n", Header{"x-pass": {"furrow-a/YST/18218"}}},
		{"a field twice", "X-A: 1\r\nx-a: 2\r\n\r\n", Header{"x-a": {"1", "2"}}},
		{"a port after the host", "Host: 127.0.0.1:8080\r\n\r\n", Header{"host": {"127.0.0.1:8080"}}},
		{"an empty value", "X-Empty:\r\n\r\n", Header{"x-empty": {""}}},
		{"no fields at all", "\r\n", Header{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readHeader(bufio.NewReaderSize(strings.NewReader(tt.head), 1024))
			if err != nil || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readHeader(%q) = %v, %v, want %v, nil", tt.head, got, err, tt.want)
			}
		})
	}
	t.Run("Get joins a repeated field", func(t *testing.T) {
		h, err := readHeader(bufio.NewReaderSize(strings.NewReader("X-A: 1\r\nx-a: 2\r\n\r\n"), 1024))
		if got := h.Get("X-A"); got != "1, 2" || err != nil {
			t.Errorf("readHeader(X-A: 1, x-a: 2).Get(\"X-A\") = %q, %v, want \"1, 2\"", got, err)
		}
	})
	t.Run("the body stays in the reader", func(t *testing.T) {
		br := bufio.NewReaderSize(strings.NewReader("Host: 127.0.0.1\r\nContent-Length: 12\r\n\r\nfirst piece\n"), 1024)
		if _, err := readHeader(br); err != nil {
			t.Fatalf("readHeader = %v, want nil", err)
		}
		if rest, _ := io.ReadAll(br); string(rest) != "first piece\n" {
			t.Errorf("after readHeader, the reader holds %q, want %q: stop at the blank line", rest, "first piece\n")
		}
	})
	t.Run("a line with a bare LF", func(t *testing.T) {
		_, err := readHeader(bufio.NewReaderSize(strings.NewReader("Host: 127.0.0.1\n\r\n"), 1024))
		if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != 400 {
			t.Errorf("readHeader(\"Host: 127.0.0.1\\n\") = %v, want readLine's RequestError 400 bad line ending", err)
		}
	})
}
