package main

import (
	"context"
	"net"
	"time"
)

// Link keeps one controller connected through drops and numbers its commands.
type Link struct {
	Dial func(ctx context.Context) (net.Conn, error)      // in tests, a fake line
	Wait func(ctx context.Context, d time.Duration) error // sleeps, or ends with ctx
}

// redial dials the controller until a dial works, waiting a second after each failed
// dial, and returns the wait's error when ctx ends first.
func (l *Link) redial(ctx context.Context) (net.Conn, error) {
	for {
		conn, err := l.Dial(ctx)
		if err == nil {
			return conn, nil
		}
		if err := l.Wait(ctx, time.Second); err != nil {
			return nil, err
		}
	}
}
