package main

import (
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestServeOne(t *testing.T) {
	deadline := time.Now().Add(time.Second) // every read in this test ends by then
	tests := []struct {
		name  string
		names []string
		want  string
	}{
		{"both dishes", []string{"YST-1", "YST-2"}, "dishd watching YST-1 YST-2\n"},
		{"Dish 1 alone", []string{"YST-1"}, "dishd watching YST-1\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatalf("listen: %v", err)
			}
			defer ln.Close()
			done := runServeOne(ln, tt.names)

			conn, err := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
			if err != nil {
				t.Fatalf("a watcher could not connect: %v", err)
			}
			defer conn.Close()
			conn.SetDeadline(deadline)
			got, err := readToEnd(conn)
			if got != tt.want || !errors.Is(err, io.EOF) {
				t.Errorf("serveOne(ln, %q): the watcher read %q, then %v; want %q, then EOF", tt.names, got, err, tt.want)
			}
			select {
			case r := <-done:
				if r.panicked != nil {
					t.Errorf("serveOne(ln, %q) panicked: %v", tt.names, r.panicked)
				} else if r.err != nil {
					t.Errorf("serveOne(ln, %q) = %v, want nil", tt.names, r.err)
				}
			case <-time.After(500 * time.Millisecond):
				t.Errorf("serveOne(ln, %q) has not returned after greeting its one watcher", tt.names)
			}
		})
	}
	t.Run("closed listener", func(t *testing.T) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}
		ln.Close()
		select {
		case r := <-runServeOne(ln, []string{"YST-1"}):
			if r.panicked != nil {
				t.Errorf("serveOne(closed listener, [YST-1]) panicked: %v; return Accept's error before using the connection", r.panicked)
			} else if !errors.Is(r.err, net.ErrClosed) {
				t.Errorf("serveOne(closed listener, [YST-1]) = %v, want Accept's error, net.ErrClosed", r.err)
			}
		case <-time.After(500 * time.Millisecond):
			t.Errorf("serveOne(closed listener, [YST-1]) has not returned; return Accept's error")
		}
	})
}

// served is what serveOne returned, or the panic it raised.
type served struct {
	err      error
	panicked any
}

// runServeOne runs serveOne in a goroutine and returns the channel its result arrives on.
func runServeOne(ln net.Listener, names []string) <-chan served {
	done := make(chan served, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				done <- served{panicked: p}
			}
		}()
		done <- served{err: serveOne(ln, names)}
	}()
	return done
}

// readToEnd reads from conn until a read fails, and returns what it read and that failure.
func readToEnd(conn net.Conn) (string, error) {
	var got strings.Builder
	buf := make([]byte, 64)
	for {
		n, err := conn.Read(buf)
		got.Write(buf[:n])
		if err != nil {
			return got.String(), err
		}
	}
}
