package main

import (
	"runtime"
	"runtime/debug"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestGetScratchEmpty(t *testing.T) {
	_, cwr := stationDay("CWR", june28)
	for range 20 {
		p := getScratch()
		if len(*p) != 0 {
			t.Fatalf("len(*getScratch()) = %d after a scratch held Cape Wren's %d went back, want 0", len(*p), len(cwr))
		}
		*p = append(*p, cwr...)
		putScratch(p)
	}
}

func TestPutScratchSizes(t *testing.T) {
	t.Run("a scratch of 65,536 passes is kept", func(t *testing.T) {
		if raceOn {
			t.Skip("under -race, Put drops one value in four on purpose")
		}
		// A pool keeps a part for each CPU it runs on, and the
		// collector may empty it at any time: until this check
		// ends, one CPU and no collection.
		defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
		defer debug.SetGCPercent(debug.SetGCPercent(-1))
		getScratch() // the scratch the pool would hand out first
		kept := make([]PassView, 0, maxScratch)
		putScratch(&kept)
		if getScratch() != &kept {
			t.Errorf("getScratch() after putScratch of a scratch of capacity %d = another scratch, want that one: %d or less goes back in the pool", maxScratch, maxScratch)
		}
	})
	t.Run("a scratch of 65,537 passes is dropped", func(t *testing.T) {
		big := make([]PassView, 0, maxScratch+1)
		for range 20 {
			putScratch(&big)
			if p := getScratch(); cap(*p) > maxScratch {
				t.Fatalf("cap(*getScratch()) = %d after putScratch of that size, want at most %d: a bigger scratch isn't kept", cap(*p), maxScratch)
			}
		}
	})
}

func TestScratchAllocs(t *testing.T) {
	if raceOn {
		t.Skip("allocation counts differ under -race")
	}
	_, cwr := stationDay("CWR", june28)
	// The collector may empty the pool at any time: none while
	// counting.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))
	allocs := testing.AllocsPerRun(100, func() {
		p := getScratch()
		*p = append(*p, cwr...)
		putScratch(p)
	})
	if allocs >= 1 {
		t.Errorf("allocations in a round of getScratch, Cape Wren's %d candidates and putScratch = %.0f, want 0 on average: the pool hands the same scratch out again", len(cwr), allocs)
	}
}

func TestScratchTwelveSlots(t *testing.T) {
	sites := newNetwork().Sites[:12]
	days := make([][]PassView, len(sites))
	for i, s := range sites {
		_, days[i] = stationDay(s.Code, june28)
	}
	var ready, done sync.WaitGroup
	ready.Add(len(sites))
	done.Add(len(sites))
	wrong := make(chan string, len(sites))
	for i, s := range sites {
		go func() {
			defer done.Done()
			p := getScratch()
			*p = append(*p, days[i]...)
			ready.Done()
			ready.Wait() // all twelve hold a scratch
			if !slices.Equal(*p, days[i]) {
				wrong <- s.Code
			}
			putScratch(p)
		}()
	}
	finished := make(chan struct{})
	go func() {
		done.Wait()
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatalf("12 slots still waiting for getScratch after 1s, want each served at once")
	}
	close(wrong)
	for code := range wrong {
		t.Errorf("%s's scratch, with 12 slots at work, held other passes than %s's, want each slot's own", code, code)
	}
}
