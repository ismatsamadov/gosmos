package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

// minute is 01:58 on 4 December, in the middle of furrow-a/LAN/21393 and of a downpour.
var minute = time.Date(2028, 12, 4, 1, 58, 0, 0, time.UTC)

// line is one record the sampler passed, as the JSON handler wrote it.
type line struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Msg     string    `json:"msg"`
	Dropped *int      `json:"dropped"` // nil when the line has no "dropped"
}

// newSampler returns a Sampler on a JSON handler that writes every level to buf.
func newSampler(buf *bytes.Buffer) *Sampler {
	return NewSampler(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

// send hands h one record of level and msg at time at, as a logger would.
func send(t *testing.T, h slog.Handler, at time.Time, level slog.Level, msg string) {
	t.Helper()
	r := slog.NewRecord(at, level, msg, 0)
	if msg == "carrier lost" {
		r.AddAttrs(slog.Float64("margin_db", -1.5))
	}
	if err := h.Handle(context.Background(), r); err != nil {
		t.Fatalf("Handle returned %v, want nil", err)
	}
}

// warnings hands h n carrier lost warnings, ten a second from start.
func warnings(t *testing.T, h slog.Handler, start time.Time, n int) {
	t.Helper()
	for i := range n {
		send(t, h, start.Add(time.Duration(i)*100*time.Millisecond), slog.LevelWarn, "carrier lost")
	}
}

// decode returns the lines in buf, stopping the test at one that isn't JSON.
func decode(t *testing.T, buf *bytes.Buffer) []line {
	t.Helper()
	var lines []line
	for _, s := range strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n") {
		if s == "" {
			continue
		}
		var l line
		if err := json.Unmarshal([]byte(s), &l); err != nil {
			t.Fatalf("the sampler passed on a line that isn't JSON: %q (%v)", s, err)
		}
		lines = append(lines, l)
	}
	return lines
}

// show writes a line's dropped, or "none".
func show(d *int) string {
	if d == nil {
		return "none"
	}
	b, _ := json.Marshal(*d)
	return string(b)
}

func TestSampleAMinute(t *testing.T) {
	var buf bytes.Buffer
	s := newSampler(&buf)
	warnings(t, s, minute, 250)
	send(t, s, minute.Add(time.Minute), slog.LevelWarn, "carrier lost")
	lines := decode(t, &buf)
	if len(lines) != 13 {
		t.Fatalf("250 warnings in one minute and one in the next wrote %d lines, want 12 and 1", len(lines))
	}
	for i, want := range []struct {
		at      time.Duration // after 01:58:00
		dropped int           // 0 for none
	}{
		{0, 0}, {100 * time.Millisecond, 0}, {200 * time.Millisecond, 0}, {300 * time.Millisecond, 0},
		{400 * time.Millisecond, 0}, {500 * time.Millisecond, 0}, {600 * time.Millisecond, 0},
		{700 * time.Millisecond, 0}, {800 * time.Millisecond, 0}, {900 * time.Millisecond, 0},
		{10900 * time.Millisecond, 99}, // the 110th
		{20900 * time.Millisecond, 99}, // the 210th
		{time.Minute, 40},              // 01:59:00, the next minute's first
	} {
		l := lines[i]
		if !l.Time.Equal(minute.Add(want.at)) {
			t.Errorf("line %d is the warning of %s, want the one of %s", i+1, l.Time.Format("15:04:05.0"), minute.Add(want.at).Format("15:04:05.0"))
		}
		switch {
		case want.dropped == 0 && l.Dropped != nil:
			t.Errorf("line %d: \"dropped\" = %s, want no \"dropped\": nothing was dropped before it", i+1, show(l.Dropped))
		case want.dropped > 0 && (l.Dropped == nil || *l.Dropped != want.dropped):
			t.Errorf("line %d: \"dropped\" = %s, want %d", i+1, show(l.Dropped), want.dropped)
		}
	}
	if got := s.Dropped(); got != 238 {
		t.Errorf("after 251 warnings and 13 lines, s.Dropped() = %d, want 238", got)
	}
}

func TestSampleNeverDropsErrors(t *testing.T) {
	var buf bytes.Buffer
	s := newSampler(&buf)
	for i := range 30 {
		send(t, s, minute.Add(time.Duration(i)*time.Second), slog.LevelError, "pass failed")
	}
	lines := decode(t, &buf)
	if len(lines) != 30 {
		t.Errorf("30 ERROR records in one minute wrote %d lines, want all 30", len(lines))
	}
	for i, l := range lines {
		if l.Dropped != nil {
			t.Errorf("ERROR line %d: \"dropped\" = %s, want none: no ERROR is ever dropped", i+1, show(l.Dropped))
		}
	}
}

func TestSampleCountsMessagesApart(t *testing.T) {
	var buf bytes.Buffer
	s := newSampler(&buf)
	for i := range 15 {
		at := minute.Add(time.Duration(i) * 100 * time.Millisecond)
		send(t, s, at, slog.LevelWarn, "carrier lost")
		send(t, s, at, slog.LevelDebug, "frame received")
	}
	n := map[string]int{}
	for _, l := range decode(t, &buf) {
		n[l.Msg]++
	}
	if n["carrier lost"] != 10 || n["frame received"] != 10 {
		t.Errorf("15 warnings and 15 frame records in one minute wrote %d and %d lines, want 10 of each: each message counts on its own", n["carrier lost"], n["frame received"])
	}
}

func TestSampleSharedByWith(t *testing.T) {
	var buf bytes.Buffer
	s := newSampler(&buf)
	pass := slog.New(s).With("pass_id", "furrow-a/LAN/21393").Handler()
	start := minute.Add(30 * time.Second)
	warnings(t, s, start, 6)
	warnings(t, pass, start.Add(time.Second), 6)
	if lines := decode(t, &buf); len(lines) != 10 {
		t.Errorf("6 warnings through the sampler and 6 through a logger made with With wrote %d lines, want 10: they share one count", len(lines))
	}
	if got := s.Dropped(); got != 2 {
		t.Errorf("s.Dropped() = %d, want 2", got)
	}
}

func TestSampleManyGoroutines(t *testing.T) {
	var buf bytes.Buffer
	s := newSampler(&buf)
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for i := range 100 {
				r := slog.NewRecord(minute.Add(time.Duration(i)*100*time.Millisecond), slog.LevelWarn, "carrier lost", 0)
				s.Handle(context.Background(), r)
			}
		})
	}
	wg.Wait()
	if lines := decode(t, &buf); len(lines) != 17 {
		t.Errorf("800 warnings from 8 goroutines in one minute wrote %d lines, want 17: the first 10, then the 110th to the 710th", len(lines))
	}
	if got := s.Dropped(); got != 783 {
		t.Errorf("s.Dropped() = %d, want 783", got)
	}
}
