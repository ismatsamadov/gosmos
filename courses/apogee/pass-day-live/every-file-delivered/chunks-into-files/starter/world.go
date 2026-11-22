package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"
)

// World is everything outside the rehearsal: main builds the real one, and the tests
// build a simulated one.
type World struct {
	Dir  string                                                      // the rehearsal's folder: its inputs, and the plan and reports it writes
	Now  func() time.Time                                            // the clock: every instant comes from it
	Dial func(ctx context.Context, station string) (net.Conn, error) // a station's radio, for the pass line (chapter 2)
	HTTP *http.Client                                                // customers' webhook receivers (chapter 5)
	Out  io.Writer                                                   // the console: one line per event, safe for many goroutines
}
