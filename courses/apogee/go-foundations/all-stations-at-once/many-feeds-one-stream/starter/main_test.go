package main

import (
	"slices"
	"testing"
	"time"
)

// heard returns a closed channel holding seqs: a station feed that heard those beacons.
func heard(seqs ...int) <-chan int {
	ch := make(chan int, len(seqs))
	for _, seq := range seqs {
		ch <- seq
	}
	close(ch)
	return ch
}

// mergeAll merges feeds and returns everything merge sends, sorted. It stops the test if the
// merged channel is still open after 1 s.
func mergeAll(t *testing.T, feeds ...<-chan int) []int {
	t.Helper()
	result := make(chan []int, 1)
	go func() {
		got := []int{}
		for seq := range merge(feeds...) {
			got = append(got, seq)
		}
		result <- got
	}()
	select {
	case got := <-result:
		slices.Sort(got)
		return got
	case <-time.After(time.Second):
		t.Fatal("merge(...) is still open after 1s: close its channel once every feed is done")
		return nil
	}
}

func TestMerge(t *testing.T) {
	t.Parallel()
	got := mergeAll(t, heard(9001, 9002), heard(9003, 9004, 9005), heard())
	want := []int{9001, 9002, 9003, 9004, 9005}
	if !slices.Equal(got, want) {
		t.Errorf("merge(heard(9001, 9002), heard(9003, 9004, 9005), heard()) sent %v (sorted), want %v", got, want)
	}
}

func TestMergeHearsEveryFeedAtOnce(t *testing.T) {
	t.Parallel()
	lan := make(chan int) // Lantern Atoll is still listening and has heard nothing yet
	defer close(lan)
	first := make(chan int, 1)
	go func() {
		out := merge(lan, heard(9002))
		first <- <-out
	}()
	select {
	case got := <-first:
		if got != 9002 {
			t.Errorf("first value from merge(silent lan, heard(9002)) = %d, want 9002", got)
		}
	case <-time.After(time.Second):
		t.Fatal("merge(silent lan, heard(9002)) sent nothing within 1s: give each feed its own goroutine")
	}
}
