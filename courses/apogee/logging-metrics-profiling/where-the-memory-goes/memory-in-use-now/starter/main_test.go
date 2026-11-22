package main

import (
	"fmt"
	"runtime"
	"testing"
)

// mib writes a number of bytes in MiB, to one decimal place.
func mib(n int64) string {
	return fmt.Sprintf("%.1f MiB", float64(n)/(1<<20))
}

// TestInUse holds 16 buffers of 1 MiB, as four stuck passes hold their frames, then lets them go.
func TestInUse(t *testing.T) {
	start := inUse()
	buffers := make([][]byte, 16)
	for i := range buffers {
		buffers[i] = make([]byte, 1<<20)
	}
	holding := inUse()
	runtime.KeepAlive(buffers)
	buffers = nil
	after := inUse()

	if rise := int64(holding) - int64(start); rise < 12<<20 {
		t.Errorf("inUse() while 16 MiB is held = %s above the start, want at least 12 MiB more", mib(rise))
	}
	if rest := int64(after) - int64(start); rest >= 4<<20 {
		t.Errorf("inUse() after the 16 MiB is let go = %s above the start, want under 4 MiB: TotalAlloc never falls, and HeapAlloc falls only after runtime.GC()", mib(rest))
	}
}
