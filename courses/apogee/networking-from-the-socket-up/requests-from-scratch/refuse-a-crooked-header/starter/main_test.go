package main

import (
	"bufio"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestReadHeader(t *testing.T) {
	tests := []struct {
		name   string
		head   string
		status int
		msg    string
	}{
		{"no colon", "Host 127.0.0.1\r\n\r\n", 400, "bad header line"},
		{"a space before the colon", "Host : 127.0.0.1\r\n\r\n", 400, "bad header line"},
		{"a space inside the name", "Ho st: 127.0.0.1\r\n\r\n", 400, "bad header line"},
		{"a tab inside the name", "Ho\tst: 127.0.0.1\r\n\r\n", 400, "bad header line"},
		{"an empty name", ": 127.0.0.1\r\n\r\n", 400, "bad header line"},
		{"a slash in the name", "Pass/ID: furrow-a/YST/18218\r\n\r\n", 400, "bad header line"},
		{"a line folded with a space", "X-Note: first\r\n second\r\n\r\n", 400, "bad header line"},
		{"a line folded with a tab", "X-Note: first\r\n\tsecond\r\n\r\n", 400, "bad header line"},
		{"a NUL in the value", "X-Note: a\x00b\r\n\r\n", 400, "bad header line"},
		{"a bell in the value", "X-Note: a\x07b\r\n\r\n", 400, "bad header line"},
		{"a DEL in the value", "X-Note: a\x7fb\r\n\r\n", 400, "bad header line"},
		{"a header line over 1024 bytes", "X-Long: " + strings.Repeat("a", 1100) + "\r\n\r\n", 431, "header line over 1024 bytes"},
		{"33 fields", fields(33) + "\r\n", 431, "more than 32 header fields"},
		{"33 lines of one field", strings.Repeat("X-A: 1\r\n", 33) + "\r\n", 431, "more than 32 header fields"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := readHeader(bufio.NewReaderSize(strings.NewReader(tt.head), 1024))
			if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != tt.status || e.Msg != tt.msg {
				t.Errorf("readHeader(%.40q) = %v, %v, want the RequestError %d %s", tt.head, h, err, tt.status, tt.msg)
			}
		})
	}
	good := []struct {
		name string
		head string
		want Header
	}{
		{"32 fields", fields(32) + "\r\n", nil},
		{"a port after the host", "Host: 127.0.0.1:8080\r\n\r\n", Header{"host": {"127.0.0.1:8080"}}},
		{"every token character", "X-a1!#$%&'*+.^_`|~: ok\r\n\r\n", Header{"x-a1!#$%&'*+.^_`|~": {"ok"}}},
		{"an empty value", "X-Empty:\r\n\r\n", Header{"x-empty": {""}}},
		{"a tab inside the value", "X-Note: a\tb\r\n\r\n", Header{"x-note": {"a\tb"}}},
		{"bytes from 0x80 up", "X-Azimuth: 183.25°\r\n\r\n", Header{"x-azimuth": {"183.25°"}}},
	}
	for _, tt := range good {
		t.Run(tt.name, func(t *testing.T) {
			h, err := readHeader(bufio.NewReaderSize(strings.NewReader(tt.head), 1024))
			if err != nil {
				t.Fatalf("readHeader(%.40q) = %v, want the fields and nil", tt.head, err)
			}
			if tt.want == nil && len(h) != 32 {
				t.Errorf("readHeader(32 fields) gave %d fields, want 32", len(h))
			}
			if tt.want != nil && !reflect.DeepEqual(h, tt.want) {
				t.Errorf("readHeader(%.40q) = %v, want %v", tt.head, h, tt.want)
			}
		})
	}
}

// fields returns n header lines, X-Field-1 to X-Field-n, each ending in CR LF.
func fields(n int) string {
	var b strings.Builder
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, "X-Field-%d: %d\r\n", i, i)
	}
	return b.String()
}
