package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

// Two passes whose records overlapped at Lantern Atoll on 8 November: Cirrine-1's began at
// 01:53:09, while Brine-3's still had a minute and a half to go.
var (
	brine3 = PassView{
		ID: "brine-3/LAN/14383",
		Pass: Pass{
			SatID:      "brine-3",
			Station:    "LAN",
			AOS:        time.Date(2028, 11, 8, 1, 48, 20, 0, time.UTC),
			LOS:        time.Date(2028, 11, 8, 1, 54, 41, 0, time.UTC),
			MaxElevDeg: 19.1,
		},
	}
	cirrine = PassView{
		ID: "cirrine-1/LAN/16402",
		Pass: Pass{
			SatID:      "cirrine-1",
			Station:    "LAN",
			AOS:        time.Date(2028, 11, 8, 1, 55, 9, 0, time.UTC),
			LOS:        time.Date(2028, 11, 8, 2, 2, 39, 0, time.UTC),
			MaxElevDeg: 22.6,
		},
	}
)

// bench flies p through l on the lab bench, whose radio always answers: the six records
// uplinkd's runPass writes for a pass that flies.
func bench(l *slog.Logger, p PassView) {
	l.Info("pass started")
	l.Debug("worker taken")
	l.Debug("lines sent", "lines", 1)
	l.Debug("frame received", "bytes", 1024)
	l.Debug("ack received", "ack_ms", 1500)
	l.Info("pass flown", "max_elev_deg", p.MaxElevDeg, "wait_ms", 131500)
}

// debugBase returns a logger at DEBUG that writes JSON lines to buf, without the time.
func debugBase(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug, ReplaceAttr: dropTime}))
}

// records decodes every line in buf.
func records(t *testing.T, buf *bytes.Buffer) []map[string]any {
	t.Helper()
	var recs []map[string]any
	for _, line := range strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n") {
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("the bench wrote %q, want JSON lines", line)
		}
		recs = append(recs, rec)
	}
	return recs
}

// show writes a decoded value as JSON writes it, or says it is missing.
func show(v any) string {
	if v == nil {
		return "missing"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

// carries reports whether rec carries pass p's ID and group.
func carries(rec map[string]any, p PassView) bool {
	group, _ := rec["pass"].(map[string]any)
	return rec["pass_id"] == p.ID && group["sat_id"] == p.SatID && group["station"] == p.Station &&
		group["aos"] == p.AOS.Format(time.RFC3339) && group["los"] == p.LOS.Format(time.RFC3339)
}

// checkPass reports the records that don't carry pass p's ID and group: how many, and the first.
func checkPass(t *testing.T, recs []map[string]any, p PassView) {
	t.Helper()
	var bad []map[string]any
	for _, rec := range recs {
		if !carries(rec, p) {
			bad = append(bad, rec)
		}
	}
	if len(bad) > 0 {
		t.Errorf(`passLogger(l, %s): %d of %d records lack the pass; the first, %q, has "pass_id" = %s and "pass" = %s`,
			p.ID, len(bad), len(recs), bad[0]["msg"], show(bad[0]["pass_id"]), show(bad[0]["pass"]))
	}
}

func TestPassLogger(t *testing.T) {
	t.Run("all six records carry the pass", func(t *testing.T) {
		var buf bytes.Buffer
		bench(passLogger(debugBase(&buf), cirrine), cirrine)
		recs := records(t, &buf)
		if len(recs) != 6 {
			t.Fatalf("one pass on passLogger(l, p) wrote %d lines, want 6", len(recs))
		}
		checkPass(t, recs, cirrine)
	})
	t.Run("two passes at once keep their own", func(t *testing.T) {
		var buf bytes.Buffer
		base := debugBase(&buf)
		var wg sync.WaitGroup
		for _, p := range []PassView{brine3, cirrine} {
			wg.Go(func() { bench(passLogger(base, p), p) })
		}
		wg.Wait()
		byPass := map[string][]map[string]any{}
		for _, rec := range records(t, &buf) {
			id, _ := rec["pass_id"].(string)
			byPass[id] = append(byPass[id], rec)
		}
		if len(byPass[brine3.ID]) != 6 || len(byPass[cirrine.ID]) != 6 {
			t.Fatalf("two passes at once wrote %d lines with %s's ID and %d with %s's, want 6 each",
				len(byPass[brine3.ID]), brine3.ID, len(byPass[cirrine.ID]), cirrine.ID)
		}
		checkPass(t, byPass[brine3.ID], brine3)
		checkPass(t, byPass[cirrine.ID], cirrine)
	})
	t.Run("the parent logger is unchanged", func(t *testing.T) {
		var buf bytes.Buffer
		base := debugBase(&buf)
		passLogger(base, cirrine)
		base.Info("lines queued", "lines", 1)
		raw := strings.TrimSpace(buf.String())
		if strings.Contains(raw, `"pass_id"`) || strings.Contains(raw, `"pass"`) {
			t.Errorf("after passLogger(l, p), l wrote %s, want no pass", raw)
		}
	})
}
