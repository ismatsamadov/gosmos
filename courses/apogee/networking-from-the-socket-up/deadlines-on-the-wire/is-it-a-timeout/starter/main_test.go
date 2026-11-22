package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"testing/iotest"
	"time"
)

func TestDescribe(t *testing.T) {
	quiet := deadlineError(t)
	tests := []struct {
		name    string
		dish    string
		timeout time.Duration
		err     error
		want    string
	}{
		{"the line goes quiet", "YST-1", 50 * time.Millisecond, quiet, "YST-1: no answer in 50ms"},
		{"the same, wrapped by recv", "YST-1", 50 * time.Millisecond, fmt.Errorf("recv: %w", quiet), "YST-1: no answer in 50ms"},
		{"Dish 2 quiet for dishd's 5s", "YST-2", 5 * time.Second, fmt.Errorf("recv: %w", quiet), "YST-2: no answer in 5s"},
		{"the controller hangs up between frames", "YST-1", 50 * time.Millisecond, io.EOF, "YST-1: closed"},
		{"the controller hangs up inside a frame", "YST-1", 50 * time.Millisecond, fmt.Errorf("reading STATUS: %w", io.ErrUnexpectedEOF), "YST-1: closed"},
		{"a damaged frame", "YST-1", 50 * time.Millisecond, ErrBadCRC, "YST-1: bad crc"},
		{"the controller says why it hangs up", "YST-1", 50 * time.Millisecond, &ByeError{Reason: "idle"}, "YST-1: controller closed: idle"},
		{"a context's deadline, not the line's", "YST-1", 50 * time.Millisecond, context.DeadlineExceeded, "YST-1: context deadline exceeded"},
		{"iotest's pretend timeout", "YST-1", 50 * time.Millisecond, iotest.ErrTimeout, "YST-1: timeout"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := describe(tt.dish, tt.timeout, tt.err); got != tt.want {
				t.Errorf("describe(%q, %v, %v) = %q, want %q", tt.dish, tt.timeout, tt.err, got, tt.want)
			}
		})
	}
}

// deadlineError returns the error of a real read that ran out of time: a read on one end of a
// pipe whose other end never writes, with a deadline 50ms away.
func deadlineError(t *testing.T) error {
	t.Helper()
	end, other := net.Pipe()
	defer end.Close()
	defer other.Close()
	end.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	_, err := end.Read(make([]byte, 16))
	if err == nil {
		t.Fatalf("a read with nothing to read returned no error; want a deadline error")
	}
	return err
}
