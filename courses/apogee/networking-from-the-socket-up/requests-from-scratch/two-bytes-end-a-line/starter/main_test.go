package main

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	tests := []struct {
		name    string
		stream  string
		want    string
		wantErr error // nil, errLineTooLong, io.EOF or io.ErrUnexpectedEOF
		status  int   // a RequestError's status, when there is one
	}{
		{"a request line", "GET / HTTP/1.1\r\n", "GET / HTTP/1.1", nil, 0},
		{"a blank line", "\r\n", "", nil, 0},
		{"1022 characters and CR LF", strings.Repeat("a", 1022) + "\r\n", strings.Repeat("a", 1022), nil, 0},
		{"1023 characters and CR LF", strings.Repeat("a", 1023) + "\r\n", "", errLineTooLong, 0},
		{"a bare LF", "Host: x\n", "", nil, 400},
		{"a lone CR inside the line", "Host: a\rb\r\n", "", nil, 400},
		{"two CRs before the LF", "Host: x\r\r\n", "", nil, 400},
		{"a CR alone at the end", "Host: x\r", "", io.ErrUnexpectedEOF, 0},
		{"a stream that ends mid-line", "GET / HT", "", io.ErrUnexpectedEOF, 0},
		{"an empty stream", "", "", io.EOF, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLine(bufio.NewReaderSize(strings.NewReader(tt.stream), 1024))
			switch {
			case tt.status != 0:
				if e, ok := errors.AsType[*RequestError](err); !ok || e.Status != tt.status || e.Msg != "bad line ending" {
					t.Errorf("readLine(%.40q) = %.40q, %v, want the RequestError 400 bad line ending", tt.stream, got, err)
				}
			case tt.wantErr != nil:
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("readLine(%.40q) = %.40q, %v, want %v", tt.stream, got, err, tt.wantErr)
				}
			case got != tt.want || err != nil:
				t.Errorf("readLine(%.40q) = %.40q, %v, want %.40q, nil", tt.stream, got, err, tt.want)
			}
		})
	}
	t.Run("64 KiB and no line end", func(t *testing.T) {
		src := &countingReader{r: strings.NewReader(strings.Repeat("a", 65536))}
		_, err := readLine(bufio.NewReaderSize(src, 1024))
		if !errors.Is(err, errLineTooLong) || src.n > 1024 {
			t.Errorf("readLine(65536 bytes with no line end) = %v after reading %d bytes, want errLineTooLong after at most 1024", err, src.n)
		}
	})
	t.Run("two lines, one at a time", func(t *testing.T) {
		br := bufio.NewReaderSize(strings.NewReader("Host: 127.0.0.1\r\nUser-Agent: Go-http-client/1.1\r\n"), 1024)
		first, err1 := readLine(br)
		second, err2 := readLine(br)
		if first != "Host: 127.0.0.1" || second != "User-Agent: Go-http-client/1.1" || err1 != nil || err2 != nil {
			t.Errorf("two readLine calls = %q, %v and %q, %v, want each line in turn: copy the line before the next read reuses the buffer", first, err1, second, err2)
		}
	})
}

// countingReader counts the bytes read through it.
type countingReader struct {
	r io.Reader
	n int
}

// Read reads from r and counts what it got.
func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += n
	return n, err
}
