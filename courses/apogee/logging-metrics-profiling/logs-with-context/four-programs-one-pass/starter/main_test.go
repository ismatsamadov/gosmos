package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// cirrine is the pass ID of Cirrine-1's pass of 8 November.
const cirrine = "cirrine-1/LAN/16402"

// fakeClock is a clock the test sets by hand.
type fakeClock struct{ now time.Time }

// Now returns the clock's time.
func (c *fakeClock) Now() time.Time { return c.now }

// timeOf writes a line's time as it has it, unless it is the Mac's own clock, which no message prints.
func timeOf(rec map[string]any) string {
	s, ok := rec["time"].(string)
	switch {
	case !ok:
		return "missing"
	case strings.HasPrefix(s, "2028-"):
		return fmt.Sprintf("%q", s)
	default:
		return "the Mac's own clock"
	}
}

func TestProgramLog(t *testing.T) {
	clock := &fakeClock{now: time.Date(2028, 11, 8, 1, 48, 30, 0, time.UTC)}
	var buf bytes.Buffer
	l := programLog(&buf, "uplinkd", "LAN", clock.Now)
	l.InfoContext(withPass(context.Background(), cirrine), "lines queued", "lines", 1)
	const call = `programLog(w, "uplinkd", "LAN", clock.Now).InfoContext(ctx, "lines queued", "lines", 1)`
	raw := strings.TrimSuffix(buf.String(), "\n")
	var rec map[string]any
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		t.Fatalf("%s wrote %q, want one JSON line", call, buf.String())
	}
	if rec["time"] != "2028-11-08T01:48:30Z" {
		t.Errorf(`%s: "time" = %s, want "2028-11-08T01:48:30Z", the clock's`, call, timeOf(rec))
	}
	for _, f := range []struct{ key, want string }{
		{"program", "uplinkd"},
		{"site", "LAN"},
		{"pass_id", cirrine},
	} {
		if rec[f.key] != f.want {
			t.Errorf("%s: %q = %s, want %q", call, f.key, field(rec, f.key), f.want)
		}
	}
	if p, s := strings.Index(raw, `"program":`), strings.Index(raw, `"site":`); p >= 0 && s >= 0 && s < p {
		t.Errorf(`%s wrote "site" before "program", want program, then site`, call)
	}
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

// canonical is Lantern Atoll's failure line of 8 November, as the night shift's search must read it.
const canonical = `{"time":"2028-11-08T02:02:39Z","level":"ERROR","msg":"pass failed","program":"uplinkd","site":"LAN","pass":{"sat_id":"cirrine-1","station":"LAN","aos":"2028-11-08T01:55:09Z","los":"2028-11-08T02:02:39Z"},"max_elev_deg":22.6,"wait_ms":570000,"err":"pass cirrine-1/LAN/16402: context deadline exceeded","pass_id":"cirrine-1/LAN/16402"}`

// TestReplay follows Cirrine-1's pass through the four programs' logs, as the night shift's search does.
func TestReplay(t *testing.T) {
	var bookingd, passplan, hq, lan bytes.Buffer
	replayNovember8(&bookingd, &passplan, &hq, &lan)
	var logs [][]Line
	for _, log := range []struct {
		name string
		buf  *bytes.Buffer
	}{{"bookingd", &bookingd}, {"passplan", &passplan}, {"uplinkd HQ", &hq}, {"uplinkd LAN", &lan}} {
		lines, broken, err := follow(bytes.NewReader(log.buf.Bytes()), cirrine)
		if err != nil || broken != 0 {
			t.Fatalf("follow(%s's log) = %d broken lines and the error %v, want JSON lines only", log.name, broken, err)
		}
		logs = append(logs, lines)
	}
	got := merge(logs...)
	want := []struct{ time, program, site, msg string }{
		{"2028-11-06T10:14:03Z", "bookingd", "HQ", "booking created"},
		{"2028-11-07T00:00:41Z", "passplan", "HQ", "pass planned"},
		{"2028-11-08T01:44:40Z", "uplinkd", "HQ", "command sealed"},
		{"2028-11-08T01:48:30Z", "uplinkd", "HQ", "lines delivered"},
		{"2028-11-08T01:48:30Z", "uplinkd", "LAN", "lines queued"},
		{"2028-11-08T01:53:09Z", "uplinkd", "LAN", "pass started"},
		{"2028-11-08T02:02:39Z", "uplinkd", "LAN", "pass failed"},
	}
	if len(got) != len(want) {
		t.Fatalf("follow and merge found %d lines of %s in the replay, want %d, one per step from booking to failure", len(got), cirrine, len(want))
	}
	for i, w := range want {
		l := got[i]
		at := l.Time.UTC().Format(time.RFC3339)
		if !strings.HasPrefix(at, "2028-") {
			at = "the Mac's own clock"
		}
		if !l.Time.Equal(mustTime(w.time)) || l.Program != w.program || l.Site != w.site || l.Msg != w.msg {
			t.Fatalf("line %d = %s %s %s %q, want %s %s %s %q", i+1, at, l.Program, l.Site, l.Msg, w.time, w.program, w.site, w.msg)
		}
		if n := strings.Count(l.Raw, `"pass_id":`); n != 1 {
			t.Errorf("line %d, %s's %q, has %d pass_id keys, want 1: %s", i+1, w.program, w.msg, n, l.Raw)
		}
	}
	if last := got[len(got)-1].Raw; last != canonical {
		t.Errorf("Lantern Atoll's failure line is\n%s\nwant\n%s", last, canonical)
	}
}

// mustTime reads one of the test's RFC 3339 times.
func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
