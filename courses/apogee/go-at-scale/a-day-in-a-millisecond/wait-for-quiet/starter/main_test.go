package main

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// burst is Ostrel Island's eleven stow and unstow messages, 400 ms apart, from at.
func burst(at time.Duration) []time.Duration {
	var times []time.Duration
	for i := range 11 {
		times = append(times, at+time.Duration(i)*400*time.Millisecond)
	}
	return times
}

// when writes times after the bubble's start, the first three and a count of the rest, or never.
func when(times []time.Duration) string {
	if len(times) == 0 {
		return "never"
	}
	var s []string
	for _, d := range times[:min(len(times), 3)] {
		s = append(s, d.String())
	}
	if len(times) > 3 {
		return fmt.Sprintf("%s and %d more", strings.Join(s, ", "), len(times)-3)
	}
	return strings.Join(s, ", ")
}

// try calls f, ending the test with name's panic instead of ending the test binary.
func try(t *testing.T, name string, f func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("%s panicked: %v", name, r)
		}
	}()
	f()
}

func TestDebouncer(t *testing.T) {
	tests := []struct {
		name   string
		notify []time.Duration // when Notify is called
		stop   time.Duration   // when Stop is called; 0 for never
		fires  []time.Duration // when fire should run
	}{
		{"eleven messages in four seconds", burst(0), 0, []time.Duration{9 * time.Second}},
		{"a second burst 10 s after the first", append(burst(0), burst(14*time.Second)...), 0, []time.Duration{9 * time.Second, 23 * time.Second}},
		{"Stop 2 s after the last message", burst(0), 6 * time.Second, nil},
		{"Stop before any message", nil, time.Second, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				start := time.Now()
				var mu sync.Mutex // fire runs on the timer's goroutine
				var fires []time.Duration
				d := newDebouncer(5*time.Second, func() {
					mu.Lock()
					defer mu.Unlock()
					fires = append(fires, time.Since(start))
				})
				for _, at := range tt.notify {
					time.Sleep(at - time.Since(start))
					try(t, "Notify()", d.Notify)
				}
				if tt.stop > 0 {
					time.Sleep(tt.stop - time.Since(start))
					try(t, "Stop()", d.Stop)
				}
				time.Sleep(time.Minute - time.Since(start))
				synctest.Wait()
				mu.Lock()
				defer mu.Unlock()
				if !slices.Equal(fires, tt.fires) {
					t.Errorf("newDebouncer(5s, fire): fire times = %s, want %s", when(fires), when(tt.fires))
				}
			})
		})
	}
}
