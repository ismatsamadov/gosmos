package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// cirrine is the pass that failed at 02:02:39 on 8 November.
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

// atOffset returns p with its times at an offset of -07:00: the same moments, written another way.
func atOffset(p PassView) PassView {
	zone := time.FixedZone("", -7*60*60)
	p.AOS, p.LOS = p.AOS.In(zone), p.LOS.In(zone)
	return p
}

// show writes a decoded value as JSON writes it, or says it is missing.
func show(v any) string {
	if v == nil {
		return "missing"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func TestPassGroupJSON(t *testing.T) {
	want := `"pass":{"sat_id":"cirrine-1","station":"LAN","aos":"2028-11-08T01:55:09Z","los":"2028-11-08T02:02:39Z"}`
	for _, tt := range []struct {
		name string
		p    PassView
	}{
		{"Cirrine-1 on 8 November", cirrine},
		{"the same pass, its times at -07:00", atOffset(cirrine)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime}))
			l.Info("pass started", passGroup(tt.p))
			raw := strings.TrimSuffix(buf.String(), "\n")
			var rec map[string]any
			if err := json.Unmarshal([]byte(raw), &rec); err != nil {
				t.Fatalf("passGroup(p) wrote %q, want one JSON line", raw)
			}
			group, ok := rec["pass"].(map[string]any)
			if !ok {
				t.Fatalf(`passGroup(p) wrote "pass":%s, want an object`, show(rec["pass"]))
			}
			right := true // whether every value is right, so only the keys' order can still be wrong
			for _, f := range []struct{ key, want string }{
				{"sat_id", "cirrine-1"},
				{"station", "LAN"},
				{"aos", "2028-11-08T01:55:09Z"},
				{"los", "2028-11-08T02:02:39Z"},
			} {
				if group[f.key] != f.want {
					right = false
					t.Errorf("passGroup(p): pass.%s = %s, want %q", f.key, show(group[f.key]), f.want)
				}
			}
			if len(group) != 4 {
				t.Errorf("passGroup(p) wrote %d keys in the group, %s, want exactly 4", len(group), show(group))
			} else if right && !strings.Contains(raw, want) {
				t.Errorf("passGroup(p) wrote %s, want the keys in this order: %s", raw, want)
			}
		})
	}
}

func TestPassGroupText(t *testing.T) {
	var buf bytes.Buffer
	l := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime}))
	l.Info("pass started", passGroup(cirrine))
	want := "pass.sat_id=cirrine-1 pass.station=LAN pass.aos=2028-11-08T01:55:09.000Z pass.los=2028-11-08T02:02:39.000Z"
	if got := strings.TrimSpace(buf.String()); !strings.HasSuffix(got, want) {
		why := ""
		if strings.Contains(got, "pass.aos=2028-11-08T01:55:09Z") {
			why = ": aos and los must be times, not text, so the text handler writes their milliseconds"
		}
		t.Errorf("passGroup(p) through the text handler wrote %q, want it to end %q%s", got, want, why)
	}
}
