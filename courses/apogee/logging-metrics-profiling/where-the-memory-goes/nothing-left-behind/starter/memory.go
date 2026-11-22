package main

import (
	"runtime"
	"time"
)

// inUse returns how many bytes of heap the program still uses, counted after a garbage collection.
func inUse() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

// leftBehind runs f and returns how many goroutines it left running, given up to a second to end.
func leftBehind(f func()) int {
	before := runtime.NumGoroutine()
	f()
	left := runtime.NumGoroutine() - before
	deadline := time.Now().Add(time.Second)
	for left > 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
		left = runtime.NumGoroutine() - before
	}
	return max(left, 0)
}
