package main

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"
)

// frames returns a pass's buffer of size bytes, each byte set from its position.
func frames(size int) []byte {
	b := make([]byte, size)
	for i := range b {
		b[i] = byte(i * 7)
	}
	return b
}

// heapInUse returns the heap still in use after a garbage collection.
func heapInUse() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

// mib writes a number of bytes in MiB, to one decimal place.
func mib(n int64) string {
	return fmt.Sprintf("%.1f MiB", float64(n)/(1<<20))
}

func TestTail(t *testing.T) {
	tests := []struct {
		name string
		size int
		n    int
		want int // how many of the last bytes come back
	}{
		{"the last 64 of 4 KiB of frames", 4 << 10, 64, 64},
		{"a buffer of exactly 64", 64, 64, 64},
		{"a buffer shorter than 64, whole", 10, 64, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := frames(tt.size)
			want := bytes.Clone(buf[tt.size-tt.want:])
			got := tail(buf, tt.n)
			if !bytes.Equal(got, want) {
				t.Fatalf("tail(%d bytes, %d) = % x, want % x", tt.size, tt.n, got, want)
			}
			buf[tt.size-1]++ // the buffer is written to after the report took its tail
			if !bytes.Equal(got, want) {
				t.Errorf("tail(%d bytes, %d) changed when its buffer did, want a copy of its own", tt.size, tt.n)
			}
		})
	}
}

// TestTailsLetGo keeps 32 tails of 256 KiB buffers, as a night of pass reports would.
func TestTailsLetGo(t *testing.T) {
	start := heapInUse()
	tails := make([][]byte, 32)
	for i := range tails {
		tails[i] = tail(frames(256<<10), 64)
	}
	held := int64(heapInUse()) - int64(start)
	runtime.KeepAlive(tails)
	if held >= 1<<20 {
		t.Errorf("32 tails of 64 bytes keep %s in use, want under 1 MiB: each window still keeps its 256 KiB buffer", mib(held))
	}
}
