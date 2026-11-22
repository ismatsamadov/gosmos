package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// fakeClock is a clock the test sets by hand.
type fakeClock struct{ now time.Time }

// Now returns the clock's time.
func (c *fakeClock) Now() time.Time { return c.now }

// Add moves the clock forward by d.
func (c *fakeClock) Add(d time.Duration) { c.now = c.now.Add(d) }

// cirrine is the pass that failed at 02:02:39 on 8 November, Cirrine-1's sealed PING still queued.
var cirrine = PassView{
	ID: "cirrine-1/LAN/16402",
	Pass: Pass{
		SatID:      "cirrine-1",
		Station:    "LAN",
		AOS:        time.Date(2028, 11, 8, 1, 55, 9, 0, time.UTC),
		LOS:        time.Date(2028, 11, 8, 2, 2, 39, 0, time.UTC),
		MaxElevDeg: 22.6,
	},
}

// brine3 is the first pass Lantern Atoll flew after the restart at 03:14.
var brine3 = PassView{
	ID: "brine-3/LAN/14384",
	Pass: Pass{
		SatID:      "brine-3",
		Station:    "LAN",
		AOS:        time.Date(2028, 11, 8, 3, 26, 17, 0, time.UTC),
		LOS:        time.Date(2028, 11, 8, 3, 32, 10, 0, time.UTC),
		MaxElevDeg: 16.6,
	},
}

// failure is the error cirrine's pass ended with when its time ran out at LOS.
var failure = fmt.Errorf("pass %s: %w", cirrine.ID, context.DeadlineExceeded)

// want is what one line must hold: its level, message, pass, wait and error ("" for none).
type want struct {
	level, msg, passID string
	elevDeg, waitMs    float64
	err                string
}

// decode splits out into lines and decodes each as a JSON object; it stops the test,
// naming call, at a line that isn't one.
func decode(t *testing.T, call, out string) []map[string]any {
	t.Helper()
	var recs []map[string]any
	for _, line := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		if line == "" {
			continue
		}
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("%s wrote %q, want a JSON object (%v)", call, line, err)
		}
		recs = append(recs, rec)
	}
	return recs
}

// field writes the value of rec's key as JSON writes it, or says it is missing.
func field(rec map[string]any, key string) string {
	v, ok := rec[key]
	if !ok {
		return "missing"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// timeOf writes rec's time as the line has it, unless it is the Mac's own clock, which no
// message prints.
func timeOf(rec map[string]any) string {
	s, ok := rec["time"].(string)
	switch {
	case !ok:
		return field(rec, "time")
	case strings.HasPrefix(s, "2028-"):
		return fmt.Sprintf("%q", s)
	default:
		return "the Mac's own clock"
	}
}

// check compares one decoded line with w, naming call in every message.
func check(t *testing.T, call string, rec map[string]any, w want) {
	t.Helper()
	if rec["level"] != w.level || rec["msg"] != w.msg {
		t.Errorf(`%s: "level" = %s and "msg" = %s, want %q and %q`, call, field(rec, "level"), field(rec, "msg"), w.level, w.msg)
	}
	if rec["pass_id"] != w.passID {
		t.Errorf(`%s: "pass_id" = %s, want %q`, call, field(rec, "pass_id"), w.passID)
	}
	if rec["max_elev_deg"] != w.elevDeg {
		t.Errorf(`%s: "max_elev_deg" = %s, want %v, a number`, call, field(rec, "max_elev_deg"), w.elevDeg)
	}
	if rec["wait_ms"] != w.waitMs {
		t.Errorf(`%s: "wait_ms" = %s, want %v, a number`, call, field(rec, "wait_ms"), w.waitMs)
	}
	_, hasErr := rec["err"]
	switch {
	case w.err == "" && hasErr:
		t.Errorf(`%s wrote "err":%s, want no "err" key`, call, field(rec, "err"))
	case w.err != "" && rec["err"] != w.err:
		t.Errorf(`%s: "err" = %s, want %q, the whole text`, call, field(rec, "err"), w.err)
	}
}

func TestNewPassLog(t *testing.T) {
	clock := &fakeClock{now: time.Date(2028, 11, 8, 1, 53, 9, 0, time.UTC)}
	var buf bytes.Buffer
	l := newPassLog(&buf, clock.Now)
	l.Info("pass started")
	clock.Add(9*time.Minute + 30*time.Second)
	l.Error("pass failed")
	const call = "newPassLog(w, clock.Now)"
	recs := decode(t, call, buf.String())
	if len(recs) != 2 {
		t.Fatalf("two records on %s wrote %d lines, want 2", call, len(recs))
	}
	for i, w := range []struct{ time, level, msg string }{
		{"2028-11-08T01:53:09Z", "INFO", "pass started"},
		{"2028-11-08T02:02:39Z", "ERROR", "pass failed"},
	} {
		if rec := recs[i]; rec["time"] != w.time {
			t.Errorf(`%s, line %d: "time" = %s, want %q, the clock's`, call, i+1, timeOf(rec), w.time)
		} else if rec["level"] != w.level || rec["msg"] != w.msg {
			t.Errorf(`%s, line %d: "level" = %s and "msg" = %s, want %q and %q`, call, i+1, field(rec, "level"), field(rec, "msg"), w.level, w.msg)
		}
	}
}

func TestPassLine(t *testing.T) {
	tests := []struct {
		name string
		call string
		p    PassView
		wait time.Duration
		err  error
		want want
	}{
		{"Cirrine-1 at 02:02:39", "passLine(l, cirrine-1/LAN/16402, 9m30s, err)", cirrine, 9*time.Minute + 30*time.Second, failure,
			want{"ERROR", "pass failed", cirrine.ID, 22.6, 570000, failure.Error()}},
		{"Brine-3 after the restart", "passLine(l, brine-3/LAN/14384, 3m21s, nil)", brine3, 3*time.Minute + 21*time.Second, nil,
			want{"INFO", "pass flown", brine3.ID, 16.6, 201000, ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			passLine(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime})), tt.p, tt.wait, tt.err)
			recs := decode(t, tt.call, buf.String())
			if len(recs) != 1 {
				t.Fatalf("%s wrote %d lines, want 1", tt.call, len(recs))
			}
			check(t, tt.call, recs[0], tt.want)
		})
	}
}

// TestReplay plays the night's failure and the next flown pass through uplinkd's new logger,
// as Nkem does on the lab bench.
func TestReplay(t *testing.T) {
	clock := &fakeClock{now: time.Date(2028, 11, 8, 2, 2, 39, 0, time.UTC)}
	var buf bytes.Buffer
	l := newPassLog(&buf, clock.Now)
	passLine(l, cirrine, 9*time.Minute+30*time.Second, failure)
	clock.Add(time.Hour + 24*time.Minute + 59*time.Second) // Brine-3's acknowledgment, 03:27:38
	passLine(l, brine3, 3*time.Minute+21*time.Second, nil)
	const call = "the replay"
	recs := decode(t, call, buf.String())
	if len(recs) != 2 {
		t.Fatalf("two passLine calls on newPassLog(w, clock.Now) wrote %d lines, want 2", len(recs))
	}
	for i, w := range []struct {
		time string
		want want
	}{
		{"2028-11-08T02:02:39Z", want{"ERROR", "pass failed", cirrine.ID, 22.6, 570000, failure.Error()}},
		{"2028-11-08T03:27:38Z", want{"INFO", "pass flown", brine3.ID, 16.6, 201000, ""}},
	} {
		line := fmt.Sprintf("%s, line %d", call, i+1)
		if recs[i]["time"] != w.time {
			t.Errorf(`%s: "time" = %s, want %q, the clock's`, line, timeOf(recs[i]), w.time)
		}
		check(t, line, recs[i], w.want)
	}
}
