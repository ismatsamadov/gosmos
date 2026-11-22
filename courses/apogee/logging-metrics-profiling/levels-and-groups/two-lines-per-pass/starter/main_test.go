package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

func TestLevelFor(t *testing.T) {
	tests := []struct {
		event string
		want  slog.Level
	}{
		{"pass started", slog.LevelInfo},
		{"worker taken", slog.LevelDebug},
		{"lines sent", slog.LevelDebug},
		{"frame received", slog.LevelDebug},
		{"ack received", slog.LevelDebug},
		{"carrier lost", slog.LevelWarn},
		{"pass flown", slog.LevelInfo},
		{"pass failed", slog.LevelError},
		{"pass skipped", slog.LevelWarn},
		{"lines queued", slog.LevelInfo},
		{"command sealed", slog.LevelInfo},
		{"lines delivered", slog.LevelInfo},
		{"radio restarted", slog.LevelInfo}, // an event nobody listed
	}
	for _, tt := range tests {
		t.Run(tt.event, func(t *testing.T) {
			if got := levelFor(tt.event); got != tt.want {
				t.Errorf("levelFor(%q) = %v, want %v", tt.event, got, tt.want)
			}
		})
	}
}

// night is what runPass wrote for each of Lantern Atoll's four passes from 01:46 to 03:28 on 8 November, event by event.
var night = []struct{ pass, event string }{
	{"brine-3/LAN/14383", "pass started"},
	{"brine-3/LAN/14383", "worker taken"},
	{"brine-3/LAN/14383", "lines sent"},
	{"brine-3/LAN/14383", "frame received"},
	{"brine-3/LAN/14383", "carrier lost"},
	{"brine-3/LAN/14383", "carrier lost"},
	{"cirrine-1/LAN/16402", "pass started"},
	{"brine-3/LAN/14383", "pass failed"},
	{"cirrine-1/LAN/16402", "pass failed"},
	{"furrow-a/LAN/20998", "pass started"},
	{"furrow-a/LAN/20998", "pass failed"},
	{"brine-3/LAN/14384", "pass started"},
	{"brine-3/LAN/14384", "worker taken"},
	{"brine-3/LAN/14384", "lines sent"},
	{"brine-3/LAN/14384", "frame received"},
	{"brine-3/LAN/14384", "frame received"},
	{"brine-3/LAN/14384", "ack received"},
	{"brine-3/LAN/14384", "pass flown"},
}

// replay writes the night through a logger at level, each event at levelFor's level, and returns
// the lines it wrote as "<pass> <message>".
func replay(t *testing.T, level slog.Level) []string {
	t.Helper()
	var buf bytes.Buffer
	l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: level, ReplaceAttr: dropTime}))
	for _, e := range night {
		l.Log(context.Background(), levelFor(e.event), e.event, "pass_id", e.pass)
	}
	var got []string
	for _, line := range strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n") {
		var rec struct {
			Msg    string `json:"msg"`
			PassID string `json:"pass_id"`
		}
		if line != "" && json.Unmarshal([]byte(line), &rec) == nil {
			got = append(got, rec.PassID+" "+rec.Msg)
		}
	}
	return got
}

func TestTheNightAtInfo(t *testing.T) {
	want := []string{
		"brine-3/LAN/14383 pass started",
		"brine-3/LAN/14383 carrier lost",
		"brine-3/LAN/14383 carrier lost",
		"cirrine-1/LAN/16402 pass started",
		"brine-3/LAN/14383 pass failed",
		"cirrine-1/LAN/16402 pass failed",
		"furrow-a/LAN/20998 pass started",
		"furrow-a/LAN/20998 pass failed",
		"brine-3/LAN/14384 pass started",
		"brine-3/LAN/14384 pass flown",
	}
	got := replay(t, slog.LevelInfo)
	for i := range max(len(got), len(want)) {
		switch {
		case i >= len(got):
			t.Fatalf("at INFO the night wrote %d lines, want %d: line %d, %q, is missing", len(got), len(want), i+1, want[i])
		case i >= len(want):
			t.Fatalf("at INFO the night wrote %d lines, want %d: each pass's start and end, and 2 warnings; line %d is %q", len(got), len(want), i+1, got[i])
		case got[i] != want[i]:
			t.Fatalf("at INFO the night wrote %d lines, want %d: line %d = %q, want %q", len(got), len(want), i+1, got[i], want[i])
		}
	}
}

func TestTheNightAtDebug(t *testing.T) {
	got := replay(t, slog.LevelDebug)
	if len(got) != len(night) {
		t.Fatalf("at DEBUG the night wrote %d lines, want all %d events", len(got), len(night))
	}
	for i, e := range night {
		if want := fmt.Sprintf("%s %s", e.pass, e.event); got[i] != want {
			t.Fatalf("at DEBUG, line %d = %q, want %q", i+1, got[i], want)
		}
	}
}
