package main

import (
	"context"
	"errors"
	"fmt"
	"runtime/trace"
)

// ErrNoWorker is runPass's error when no worker came free before the pass's LOS.
var ErrNoWorker = errors.New("no free worker before LOS")

// runPass runs one pass on a free worker and waits for it until LOS.
func (u *Uplink) runPass(ctx context.Context, p PassView, lines []string) error {
	_ = trace.NewTask // TODO: start the pass's task, log its ID, and end the task when runPass returns
	select {
	case u.workers <- struct{}{}:
	case <-ctx.Done():
		return fmt.Errorf("pass %s: %w", p.ID, ErrNoWorker)
	}
	defer func() { <-u.workers }()
	_, err := u.fly(ctx, p, lines, make([]byte, 0, 4<<20)) // ErrNoAck when ctx ends first
	if err != nil {
		return fmt.Errorf("pass %s: %w", p.ID, err)
	}
	return nil
}
