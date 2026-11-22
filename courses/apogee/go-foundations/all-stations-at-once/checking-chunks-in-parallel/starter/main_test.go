package main

import (
	"fmt"
	"hash/crc32"
	"runtime"
	"runtime/metrics"
	"slices"
	"testing"
	"time"
)

// image returns the 2,000 chunks of a Furrow-A image and their CRC-32 sums, with the chunks
// numbered in bad corrupted on the way down.
func image(bad ...int) ([][]byte, []uint32) {
	chunks := make([][]byte, 2000)
	sums := make([]uint32, 2000)
	for i := range chunks {
		chunks[i] = []byte(fmt.Sprintf("furrow-a img-0623 chunk %04d", i))
		sums[i] = crc32.ChecksumIEEE(chunks[i])
	}
	for _, i := range bad {
		chunks[i][0] ^= 0xff
	}
	return chunks, sums
}

// goroutinesStarted returns how many goroutines the program has started so far, and false if
// this version of Go can't tell.
func goroutinesStarted() (uint64, bool) {
	s := []metrics.Sample{{Name: "/sched/goroutines-created:goroutines"}}
	metrics.Read(s)
	if s[0].Value.Kind() != metrics.KindUint64 {
		return 0, false
	}
	return s[0].Value.Uint64(), true
}

// checkWithin calls checkAll and stops the test if it hasn't returned within 1 s. It also
// returns how many goroutines checkAll started, and whether that could be counted.
func checkWithin(t *testing.T, chunks [][]byte, sums []uint32, workers int) ([]bool, uint64, bool) {
	t.Helper()
	type result struct {
		ok      []bool
		started uint64
		counted bool
	}
	done := make(chan result, 1)
	runtime.GC() // the collector starts its own goroutines now, not while checkAll runs
	go func() {
		before, counted := goroutinesStarted()
		ok := checkAll(chunks, sums, workers)
		after, _ := goroutinesStarted()
		done <- result{ok, after - before, counted}
	}()
	select {
	case r := <-done:
		return r.ok, r.started, r.counted
	case <-time.After(time.Second):
		t.Fatalf("checkAll(chunks, sums, %d) has not returned after 1s: close the index channel after the last send, then wait", workers)
		return nil, 0, false
	}
}

func TestCheckAll(t *testing.T) {
	bad := []int{17, 1024, 1998}
	cases := []struct {
		name    string
		workers int
	}{
		{"one worker", 1},
		{"eight workers", 8},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			chunks, sums := image(bad...)
			ok, started, counted := checkWithin(t, chunks, sums, c.workers)
			if len(ok) != len(chunks) {
				t.Fatalf("checkAll(chunks, sums, %d) returned %d results, want %d", c.workers, len(ok), len(chunks))
			}
			var failed []int
			for i, good := range ok {
				if !good {
					failed = append(failed, i)
				}
			}
			if len(failed) > 10 {
				t.Errorf("checkAll(chunks, sums, %d) failed %d chunks, want only chunks %v", c.workers, len(failed), bad)
			} else if !slices.Equal(failed, bad) {
				t.Errorf("checkAll(chunks, sums, %d) failed chunks %v, want %v", c.workers, failed, bad)
			}
			// A pool has one goroutine per worker; a little more allows a goroutine that sends.
			if counted && (started < uint64(c.workers) || started > uint64(2*c.workers+2)) {
				t.Errorf("checkAll(chunks, sums, %d) started %d goroutines, want %d, one per worker", c.workers, started, c.workers)
			}
		})
	}
}
