package main

import (
	"runtime"
	"time"
)

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
