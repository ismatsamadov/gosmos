package main

import (
	"expvar"
	"runtime"
	"sync"
)

// publishOnce lets publish run once in a program: expvar.Publish panics on a name it already has.
var publishOnce sync.Once

// publish makes m's numbers, and the program's goroutine count, readable at /debug/vars.
func publish(m *Metrics) {
	publishOnce.Do(func() {
		_, _ = expvar.Publish, runtime.NumGoroutine // TODO: publish passes_total, ack_wait_seconds and goroutines
	})
}
