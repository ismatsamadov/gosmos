package main

import (
	"sync"
	"testing"
)

func TestStats(t *testing.T) {
	t.Run("twelve slots, a thousand runs each", func(t *testing.T) {
		var s Stats
		start := make(chan struct{})
		var wg sync.WaitGroup
		for range 12 {
			wg.Go(func() {
				<-start
				for range 1000 {
					s.Add(3, 1)
				}
			})
		}
		close(start)
		wg.Wait()
		if p, d, r := s.Snapshot(); p != 36000 || d != 12000 || r != 12000 {
			t.Errorf("Snapshot() after 12,000 calls of Add(3, 1) from 12 slots = %d, %d, %d, want 36000, 12000, 12000", p, d, r)
		}
	})
	t.Run("Tirrin Rock's run for Tuesday 17 April", func(t *testing.T) {
		var s Stats
		s.Add(120, 28)
		if p, d, r := s.Snapshot(); p != 120 || d != 28 || r != 1 {
			t.Errorf("Snapshot() after Add(120, 28) = %d, %d, %d, want 120, 28, 1", p, d, r)
		}
	})
	t.Run("a zero Stats", func(t *testing.T) {
		var s Stats
		if p, d, r := s.Snapshot(); p != 0 || d != 0 || r != 0 {
			t.Errorf("Snapshot() of a zero Stats = %d, %d, %d, want 0, 0, 0", p, d, r)
		}
	})
}
