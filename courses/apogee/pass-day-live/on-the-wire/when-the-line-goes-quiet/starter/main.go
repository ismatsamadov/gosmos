package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

// main rehearses the week described in the folder named on the command line against the
// real ground segment: the clock, each station's radio over TCP and the customers'
// receivers over HTTPS. Interrupt it to stop the rehearsal.
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: rehearsal <folder>")
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	var d net.Dialer
	w := World{
		Dir: os.Args[1],
		Now: func() time.Time { return time.Now().UTC() },
		Dial: func(ctx context.Context, station string) (net.Conn, error) {
			return d.DialContext(ctx, "tcp", "radio."+strings.ToLower(station)+".apogee.internal:7301")
		},
		HTTP: &http.Client{
			Timeout:       5 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
		Out: &lineWriter{w: os.Stdout},
	}
	if err := rehearse(ctx, w); err != nil {
		fmt.Fprintln(os.Stderr, "rehearsal:", err)
		os.Exit(1)
	}
}

// lineWriter passes each Write on whole, one at a time, so lines from many goroutines
// never mix.
type lineWriter struct {
	mu sync.Mutex
	w  io.Writer
}

// Write writes p in one piece.
func (l *lineWriter) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.w.Write(p)
}
