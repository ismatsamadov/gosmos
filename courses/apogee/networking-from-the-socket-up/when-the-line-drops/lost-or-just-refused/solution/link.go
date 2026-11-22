package main

import (
	"context"
	"errors"
	"io"
	"math/rand/v2"
	"net"
	"os"
	"syscall"
	"time"
)

// Link keeps one controller connected through drops and numbers its commands.
type Link struct {
	Dial    func(ctx context.Context) (net.Conn, error)      // in tests, a fake line
	Wait    func(ctx context.Context, d time.Duration) error // sleeps, or ends with ctx
	Rand    *rand.Rand                                       // jitter (math/rand/v2)
	attempt int                                              // failed dials in a row
}

// redial dials the controller until a dial works, waiting redialDelay after each failed
// dial, and returns the wait's error when ctx ends first.
func (l *Link) redial(ctx context.Context) (net.Conn, error) {
	for {
		conn, err := l.Dial(ctx)
		if err == nil {
			l.attempt = 0
			return conn, nil
		}
		l.attempt++
		if err := l.Wait(ctx, redialDelay(l.attempt, l.Rand)); err != nil {
			return nil, err
		}
	}
}

// lost reports whether err means the line is lost, so dishd hangs up and reconnects.
func lost(err error) bool {
	if _, ok := errors.AsType[*ByeError](err); ok {
		return true
	}
	for _, target := range []error{io.EOF, io.ErrUnexpectedEOF, os.ErrDeadlineExceeded,
		net.ErrClosed, io.ErrClosedPipe, syscall.ECONNRESET, syscall.EPIPE,
		ErrBadLength, ErrBadCRC} {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
