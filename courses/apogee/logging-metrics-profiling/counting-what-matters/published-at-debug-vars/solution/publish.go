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
		expvar.Publish("passes_total", expvar.Func(func() any { return m.Passes.Snapshot() }))
		expvar.Publish("ack_wait_seconds", expvar.Func(func() any { return m.AckWait.Snapshot() }))
		expvar.Publish("goroutines", expvar.Func(func() any { return runtime.NumGoroutine() }))
	})
}
