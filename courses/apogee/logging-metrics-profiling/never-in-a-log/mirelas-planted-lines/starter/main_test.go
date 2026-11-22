package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

// The test secrets of the booking service, the uplink keys and the ops console, readable on purpose.
const (
	furrowKey   = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"
	brineKey    = "brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0"
	cirrineKey  = "cirrine.c3a1f7e95d2b4086a7e1c9f3b5d8a024"
	furrowToken = "furrowlight.1830682800.767de44da59ef9b9e3afb9354270dd7cc56e26ceb00650b8a615afe1cd10b616"
	password    = "fails closed, always"
	hookSecret  = "cirrine webhook secret"
)

// key is the lab bench's copy of k-2028-10.
var key = Key{ID: "k-2028-10", Secret: []byte("uplink next master key, 32 bytes")}

// plantedSecrets is everything that must appear nowhere in the output.
var plantedSecrets = []struct{ what, s string }{
	{"Furrowlight's key", furrowKey[len("furrowlight."):]},
	{"Brinewatch's key", brineKey[len("brinewatch."):]},
	{"Cirrine's key", cirrineKey[len("cirrine."):]},
	{"Furrowlight's token", furrowToken[len("furrowlight.1830682800."):]},
	{"k-2028-10's secret", string(key.Secret)},
	{"k-2028-10's secret in hex", hex.EncodeToString(key.Secret)},
	{"k-2028-10's secret in base64", base64.StdEncoding.EncodeToString(key.Secret)},
	{"Mirela's password", password},
	{"Cirrine's webhook secret", hookSecret},
	{"Hollis's full address", "hollis.vey@"},
	{"Cirrine's full address", "uplink@cirrine.example"},
}

// plant writes one record, at a fixed time, through h.
func plant(h slog.Handler, at time.Time, level slog.Level, msg string, attrs ...slog.Attr) error {
	r := slog.NewRecord(at, level, msg, 0)
	r.AddAttrs(attrs...)
	return h.Handle(context.Background(), r)
}

// line is one planted line: how Mirela writes it, and the fields its JSON must hold, by path.
type line struct {
	name  string
	write func(l *slog.Logger, at time.Time) error
	want  map[string]any
}

// mirelasLines are her twelve lines of Friday 1 December.
func mirelasLines() []line {
	errInvalid := errors.New("invalid token")
	return []line{
		{"a header dump", func(l *slog.Logger, at time.Time) error {
			h := http.Header{"Authorization": {"Key " + furrowKey}, "Content-Type": {"application/json"}}
			return plant(l.Handler(), at, slog.LevelDebug, "request headers", slog.Any("header", h))
		}, map[string]any{"header.Authorization": []any{"[redacted]"}, "header.Content-Type": []any{"application/json"}}},
		{"a key inside an error", func(l *slog.Logger, at time.Time) error {
			return plant(l.Handler(), at, slog.LevelError, "request failed",
				slog.Any("err", fmt.Errorf("verify token %q: %w", furrowKey, errInvalid)))
		}, map[string]any{"err": `verify token "furrowlight.[redacted]": invalid token`}},
		{"a Key", func(l *slog.Logger, at time.Time) error {
			return plant(l.Handler(), at, slog.LevelInfo, "command sealed", slog.Int("seq", 4688), slog.Any("key", key))
		}, map[string]any{"seq": 4688.0, "key": "k-2028-10 (secret hidden)"}},
		{"a password in a group", func(l *slog.Logger, at time.Time) error {
			return plant(l.Handler(), at, slog.LevelInfo, "booking created", slog.String("booking_id", "bk-5611"),
				slog.Group("login", slog.String("user", "mirela"), slog.String("password", password)))
		}, map[string]any{"booking_id": "bk-5611", "login.user": "mirela", "login.password": "[redacted]"}},
		{"a token in the message", func(l *slog.Logger, at time.Time) error {
			return plant(l.Handler(), at, slog.LevelError, "token "+furrowToken+" refused", slog.Int("status", 401))
		}, map[string]any{"msg": "token furrowlight.[redacted] refused", "status": 401.0}},
		{"Hollis's address under email", func(l *slog.Logger, at time.Time) error {
			return plant(l.Handler(), at, slog.LevelError, "request failed", slog.String("email", "hollis.vey@furrowlight.example"),
				slog.Any("err", errors.New("deliver evt-000118: 503 Service Unavailable")))
		}, map[string]any{"email": "h***@furrowlight.example", "err": "deliver evt-000118: 503 Service Unavailable"}},
		{"an APIKey under a key of its own", func(l *slog.Logger, at time.Time) error {
			return plant(l.Handler(), at, slog.LevelInfo, "booking created", slog.String("booking_id", "bk-5611"),
				slog.Any("customer_key", APIKey(cirrineKey)))
		}, map[string]any{"booking_id": "bk-5611", "customer_key": "cirrine.[redacted]"}},
		{"headers dumped as attributes", func(l *slog.Logger, at time.Time) error {
			return plant(l.Handler(), at, slog.LevelDebug, "request headers", slog.Group("header",
				slog.String("authorization", "Bearer "+furrowToken), slog.String("Cookie", "token="+furrowToken)))
		}, map[string]any{"header.authorization": "[redacted]", "header.Cookie": "token=furrowlight.[redacted]"}},
		{"a key under api_key", func(l *slog.Logger, at time.Time) error {
			return plant(l.Handler(), at, slog.LevelError, "request failed", slog.String("api_key", brineKey),
				slog.Any("err", errors.New("invalid api key")))
		}, map[string]any{"api_key": "[redacted]", "err": "invalid api key"}},
		{"a webhook's secret and a contact", func(l *slog.Logger, at time.Time) error {
			return plant(l.Handler(), at, slog.LevelInfo, "booking created",
				slog.Group("webhook", slog.String("secret", hookSecret), slog.Int("attempts", 8)),
				slog.Group("contact", slog.String("email", "uplink@cirrine.example")))
		}, map[string]any{"webhook.secret": "[redacted]", "webhook.attempts": 8.0, "contact.email": "u***@cirrine.example"}},
		{"the 8 November failure, with nothing to hide", func(l *slog.Logger, at time.Time) error {
			return plant(l.Handler(), at, slog.LevelError, "pass failed", slog.String("pass_id", "cirrine-1/LAN/16402"),
				slog.Float64("max_elev_deg", 22.6), slog.Int64("wait_ms", 570000),
				slog.Any("err", fmt.Errorf("pass cirrine-1/LAN/16402: %w", context.DeadlineExceeded)))
		}, map[string]any{"level": "ERROR", "msg": "pass failed", "pass_id": "cirrine-1/LAN/16402", "max_elev_deg": 22.6,
			"wait_ms": 570000.0, "err": "pass cirrine-1/LAN/16402: context deadline exceeded"}},
		{"a token given to With", func(l *slog.Logger, at time.Time) error {
			return plant(l.With("program", "bookingd", "Token", furrowToken).Handler(), at, slog.LevelInfo, "booking created",
				slog.String("booking_id", "bk-5613"))
		}, map[string]any{"program": "bookingd", "Token": "[redacted]", "booking_id": "bk-5613"}},
	}
}

// at returns the time of Mirela's nth line, a second apart from 16:00.
func at(n int) time.Time {
	return time.Date(2028, 12, 1, 16, 0, n, 0, time.UTC)
}

// field returns the value at a dotted path in a decoded line, and whether it is there.
func field(rec map[string]any, path string) (any, bool) {
	var v any = rec
	for _, k := range strings.Split(path, ".") {
		m, ok := v.(map[string]any)
		if !ok {
			return nil, false
		}
		if v, ok = m[k]; !ok {
			return nil, false
		}
	}
	return v, true
}

// show writes v as JSON writes it.
func show(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func TestSafeLogger(t *testing.T) {
	var buf bytes.Buffer
	l := safeLogger(&buf)
	lines := mirelasLines()
	for i, ln := range lines {
		if err := ln.write(l, at(i)); err != nil {
			t.Fatalf("line %d (%s): writing it failed: %v", i+1, ln.name, err)
		}
	}
	out := buf.String()
	for _, sc := range plantedSecrets {
		if strings.Contains(out, sc.s) {
			t.Errorf("the output holds %s, want it kept out", sc.what)
		}
	}
	got := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(got) != len(lines) {
		t.Fatalf("twelve planted records wrote %d lines, want %d", len(got), len(lines))
	}
	for i, ln := range lines {
		var rec map[string]any
		if err := json.Unmarshal([]byte(got[i]), &rec); err != nil {
			t.Errorf("line %d (%s) isn't a JSON object (%v): %s", i+1, ln.name, err, got[i])
			continue
		}
		if want := at(i).Format(time.RFC3339); rec["time"] != want {
			t.Errorf("line %d (%s): \"time\" = %s, want %q", i+1, ln.name, show(rec["time"]), want)
		}
		for _, path := range slices.Sorted(maps.Keys(ln.want)) {
			want := ln.want[path]
			v, ok := field(rec, path)
			switch {
			case !ok:
				t.Errorf("line %d (%s): %q is missing, want %s", i+1, ln.name, path, show(want))
			case !reflect.DeepEqual(v, want):
				t.Errorf("line %d (%s): %q = %s, want %s", i+1, ln.name, path, show(v), show(want))
			}
		}
	}
}
