package main

import (
	"fmt"
	"log/slog"
	"runtime/debug"
)

// guardPass prepares pass id with work. A panic in work is recovered, so the station flies on,
// and ends the pass: guardPass writes its failure line.
func guardPass(l *slog.Logger, id string, work func()) {
	defer func() {
		if v := recover(); v != nil {
			err := fmt.Errorf("pass %s: panic: %v", id, v)
			l.Error("pass failed", "pass_id", id, "err", err, "stack", string(debug.Stack()))
		}
	}()
	work()
}
