package main

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestPipePair(t *testing.T) {
	var ends []net.Conn
	t.Run("the pair carries bytes", func(t *testing.T) {
		a, b := pipePair(t)
		if a == nil || b == nil {
			t.Fatalf("pipePair(t) = %v, %v, want two connected ends", a, b)
		}
		ends = []net.Conn{a, b}
		deadline := time.Now().Add(time.Second)
		a.SetDeadline(deadline)
		b.SetDeadline(deadline)
		checks := []struct {
			name     string
			from, to net.Conn
			msg      string
		}{
			{"first end to second", a, b, "PING"},
			{"second end to first", b, a, "PONG"},
			{"the watch port's greeting", a, b, greeting([]string{"YST-1", "YST-2"})},
		}
		for _, c := range checks {
			sent := make(chan error, 1)
			go func() { // a pipe has no buffer: this Write waits for the Read below
				_, err := io.WriteString(c.from, c.msg)
				sent <- err
			}()
			buf := make([]byte, 64)
			n, err := c.to.Read(buf)
			werr := <-sent
			if errors.Is(err, io.ErrClosedPipe) || errors.Is(werr, io.ErrClosedPipe) {
				t.Fatalf("%s: Read = %v, Write = %v, want both to work: the ends were closed before the test used them, as a defer in pipePair would; t.Cleanup waits for the test to end", c.name, err, werr)
			}
			if got := string(buf[:n]); got != c.msg || err != nil {
				t.Errorf("%s: one Read = %q, %v, want %q, nil", c.name, got, err, c.msg)
			}
			if werr != nil {
				t.Errorf("%s: Write = %v, want nil", c.name, werr)
			}
		}
	})
	if ends == nil {
		return
	}
	for i, end := range ends {
		end.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
		if _, err := end.Read(make([]byte, 1)); !errors.Is(err, io.ErrClosedPipe) {
			t.Errorf("Read on end %d after the test that made the pair ended = %v, want io.ErrClosedPipe: close both ends with t.Cleanup", i+1, err)
		}
	}
}
