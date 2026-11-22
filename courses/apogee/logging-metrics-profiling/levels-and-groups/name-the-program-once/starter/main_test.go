package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

// jsonBase returns the base logger every program starts from: JSON lines on buf, without the time.
func jsonBase(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{ReplaceAttr: dropTime}))
}

// lastLine decodes the last line in buf; it stops the test, naming call, when there is none.
func lastLine(t *testing.T, call string, buf *bytes.Buffer) (map[string]any, string) {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
	raw := lines[len(lines)-1]
	var rec map[string]any
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		t.Fatalf("%s wrote %q, want one JSON line", call, raw)
	}
	return rec, raw
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

func TestProgramLogger(t *testing.T) {
	t.Run("every record names its program and site", func(t *testing.T) {
		var buf bytes.Buffer
		l := programLogger(jsonBase(&buf), "uplinkd", "LAN")
		for _, tt := range []struct {
			call, own string
			write     func()
		}{
			{`Warn("carrier lost", "margin_db", -1.5)`, "margin_db", func() { l.Warn("carrier lost", "margin_db", -1.5) }},
			{`Info("lines queued", "lines", 1)`, "lines", func() { l.Info("lines queued", "lines", 1) }},
		} {
			tt.write()
			call := `programLogger(base, "uplinkd", "LAN").` + tt.call
			rec, raw := lastLine(t, call, &buf)
			if rec["program"] != "uplinkd" || rec["site"] != "LAN" {
				t.Errorf(`%s: "program" = %s and "site" = %s, want "uplinkd" and "LAN"`, call, field(rec, "program"), field(rec, "site"))
				continue
			}
			p, s, own := strings.Index(raw, `"program":`), strings.Index(raw, `"site":`), strings.Index(raw, `"`+tt.own+`":`)
			if !(p < s && s < own) {
				t.Errorf("%s wrote %s, want program, then site, then the record's own %s", call, raw, tt.own)
			}
		}
	})
	t.Run("the base logger is unchanged", func(t *testing.T) {
		var buf bytes.Buffer
		base := jsonBase(&buf)
		programLogger(base, "uplinkd", "LAN")
		base.Info("lines queued", "lines", 1)
		rec, raw := lastLine(t, `base.Info("lines queued", "lines", 1)`, &buf)
		if _, ok := rec["program"]; ok {
			t.Errorf("after programLogger(base, ...), base wrote %s, want no program", raw)
		}
		if _, ok := rec["site"]; ok {
			t.Errorf("after programLogger(base, ...), base wrote %s, want no site", raw)
		}
	})
	t.Run("Lantern Atoll and Ystmark don't mix", func(t *testing.T) {
		var buf bytes.Buffer
		base := jsonBase(&buf)
		lan := programLogger(base, "uplinkd", "LAN")
		yst := programLogger(base, "uplinkd", "YST")
		for _, tt := range []struct {
			l    *slog.Logger
			site string
		}{{lan, "LAN"}, {yst, "YST"}, {lan, "LAN"}} {
			tt.l.Info("pass started")
			call := `the ` + tt.site + ` logger's Info("pass started")`
			rec, raw := lastLine(t, call, &buf)
			if rec["site"] != tt.site || strings.Count(raw, `"site":`) != 1 {
				t.Errorf("%s wrote %s, want one site, %q", call, raw, tt.site)
			}
		}
	})
}
