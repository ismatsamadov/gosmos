package main

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"
)

// settle waits, for at most a second, until the number of goroutines has stayed the same for
// 20 checks in a row, so that goroutines an earlier test released have ended before this one
// counts.
func settle() {
	n, same := runtime.NumGoroutine(), 0
	for range 1000 {
		time.Sleep(time.Millisecond)
		m := runtime.NumGoroutine()
		same++
		if m != n {
			n, same = m, 0
		}
		if same == 20 {
			return
		}
	}
}

// waitFor checks every millisecond, for at most a second, until at most n goroutines exist.
func waitFor(n int) {
	for range 1000 {
		if runtime.NumGoroutine() <= n {
			return
		}
		time.Sleep(time.Millisecond)
	}
}

// plant starts n goroutines that wait until the test ends; the test's cleanup releases them and
// waits until they have ended.
func plant(t *testing.T, n int) {
	release := make(chan struct{})
	var wg sync.WaitGroup
	for range n {
		wg.Go(func() { <-release })
	}
	t.Cleanup(func() {
		close(release)
		wg.Wait()
	})
}

func TestGoroutinesThatEnd(t *testing.T) {
	settle()
	got := leftBehind(func() {
		for range 3 {
			go time.Sleep(20 * time.Millisecond)
		}
	})
	if got != 0 {
		t.Errorf("leftBehind(three goroutines that end after 20 ms) = %d, want 0", got)
	}
}

// TestFlownPass flies Furrow-A's pass of 02:22:30 on uplinkd's runPass: the satellite answers,
// runPass returns, and the pass's goroutine ends a moment later, giving its worker back.
func TestFlownPass(t *testing.T) {
	settle()
	u := newUplink("LAN", func(p PassView, lines []string, buf []byte) (Report, error) {
		return Report{Frames: buf}, nil // acknowledged
	})
	pass := PassView{ID: "furrow-a/LAN/22382"}
	var err error
	got := leftBehind(func() {
		err = u.runPass(context.Background(), pass, []string{"PING"})
	})
	if err != nil {
		t.Fatalf("runPass(furrow-a/LAN/22382) = %v, want nil", err)
	}
	if got != 0 {
		t.Errorf("leftBehind(a flown pass on runPass) = %d, want 0", got)
	}
}

func TestPlantedGoroutines(t *testing.T) {
	settle()
	got := leftBehind(func() { plant(t, 3) })
	if got != 3 {
		t.Errorf("leftBehind(three goroutines that wait for the test to end) = %d, want 3", got)
	}
}

// TestEarlierGoroutines runs functions while goroutines started before them are running.
func TestEarlierGoroutines(t *testing.T) {
	settle()
	plant(t, 2)
	if got := leftBehind(func() {}); got != 0 {
		t.Errorf("leftBehind(a function that starts nothing) = %d, with 2 goroutines running from before, want 0", got)
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup
	for range 2 {
		wg.Go(func() { <-stop })
	}
	settle()
	n := runtime.NumGoroutine()
	got := leftBehind(func() {
		close(stop)
		wg.Wait()
		waitFor(n - 2) // until both have ended, and the count is below where it started
	})
	if got != 0 {
		t.Errorf("leftBehind(a function that ends 2 goroutines started before it) = %d, want 0", got)
	}
}
