package main

import (
	"runtime"
	"time"
)

// leftBehind runs f and returns how many goroutines it left running, given up to a second to end.
func leftBehind(f func()) int {
	_, _ = runtime.NumGoroutine, time.Sleep // TODO: count before and after f, waiting for the count to come back
	f()
	return 0
}
