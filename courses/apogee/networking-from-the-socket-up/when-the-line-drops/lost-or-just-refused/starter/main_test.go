package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestLost(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
		why  string
	}{
		{"the controller hung up", io.EOF, true, "the stream ended"},
		{"cut inside a frame", io.ErrUnexpectedEOF, true, "the stream ended mid-frame"},
		{"no answer in 5s", deadlineError(t), true, "a deadline passed with nothing on the line"},
		{"dishd closed it", net.ErrClosed, true, "the connection is closed"},
		{"a closed pipe", io.ErrClosedPipe, true, "the connection is closed"},
		{"reset by the controller", &net.OpError{Op: "read", Net: "tcp", Err: os.NewSyscallError("read", syscall.ECONNRESET)}, true, "the other end reset the connection"},
		{"a broken pipe", &net.OpError{Op: "write", Net: "tcp", Err: os.NewSyscallError("write", syscall.EPIPE)}, true, "the other end is gone"},
		{"BYE idle", &ByeError{Reason: "idle"}, true, "the controller hung up"},
		{"a damaged length", fmt.Errorf("%w %d", ErrBadLength, 65535), true, "the stream is out of step"},
		{"a damaged answer", ErrBadCRC, true, "the answer can't be trusted"},
		{"NAK wind lock", &NakError{Seq: 121, Code: 3, Reason: "wind lock"}, false, "a NAK is an answer, and the line works"},
		{"the caller stopped", context.Canceled, false, "the caller said stop, not try again"},
		{"the caller's deadline", context.DeadlineExceeded, false, "the caller's time is up, not the line's"},
		{"no error", nil, false, "nothing went wrong"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plain := lost(tt.err)
			if plain != tt.want {
				t.Errorf("lost(%v) = %t, want %t: %s", tt.err, plain, tt.want, tt.why)
			}
			if tt.err == nil {
				return
			}
			wrapped := fmt.Errorf("YST-2: %w", tt.err)
			if got := lost(wrapped); got != tt.want {
				why := tt.why
				if plain == tt.want {
					why = "look through the wrapping with errors.Is or errors.AsType"
				}
				t.Errorf("lost(%v) = %t, want %t: %s", wrapped, got, tt.want, why)
			}
		})
	}
}

// deadlineError returns the error a read gives when its deadline passes on a silent line.
func deadlineError(t *testing.T) error {
	t.Helper()
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()
	a.SetReadDeadline(time.Now())
	_, err := a.Read(make([]byte, 1))
	if err == nil {
		t.Fatal("a read past its deadline worked")
	}
	return err
}
