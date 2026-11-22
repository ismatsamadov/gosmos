package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

// key is the lab bench's copy of k-2028-10, its test secret readable on purpose.
var key = Key{ID: "k-2028-10", Secret: []byte("uplink next master key, 32 bytes")}

// leaks returns each form of key's secret that out holds: its text, its hex or its base64.
func leaks(out string) []string {
	var found []string
	for _, form := range []struct{ what, s string }{
		{"the secret", string(key.Secret)},
		{"the secret in hex", hex.EncodeToString(key.Secret)},
		{"the secret in base64", base64.StdEncoding.EncodeToString(key.Secret)},
	} {
		if strings.Contains(out, form.s) {
			found = append(found, form.what)
		}
	}
	return found
}

func TestKeyInJSON(t *testing.T) {
	tests := []struct {
		name  string
		write func(l *slog.Logger)
		path  []string // where the key sits in the decoded line
	}{
		{"the sealing line", func(l *slog.Logger) {
			l.Info("command sealed", "seq", 4688, "key", key)
		}, []string{"key"}},
		{"inside a group", func(l *slog.Logger) {
			l.Info("command sealed", slog.Group("command", "seq", 4688, "key", key))
		}, []string{"command", "key"}},
		{"through With", func(l *slog.Logger) {
			l.With("key", key).Info("command sealed", "seq", 4688)
		}, []string{"key"}},
		{"by pointer", func(l *slog.Logger) {
			l.Info("command sealed", "seq", 4688, "key", &key)
		}, []string{"key"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.write(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime})))
			line := strings.TrimSuffix(buf.String(), "\n")
			for _, what := range leaks(line) {
				t.Errorf("the JSON line holds %s: %s", what, line)
			}
			var v any
			if err := json.Unmarshal([]byte(line), &v); err != nil {
				t.Fatalf("the line %q isn't JSON (%v)", line, err)
			}
			for _, k := range tt.path {
				m, _ := v.(map[string]any)
				v = m[k]
			}
			if v != "k-2028-10 (secret hidden)" {
				b, _ := json.Marshal(v)
				t.Errorf("%q = %s, want \"k-2028-10 (secret hidden)\"", strings.Join(tt.path, "."), b)
			}
		})
	}
}

func TestKeyInText(t *testing.T) {
	var buf bytes.Buffer
	slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime})).Info("command sealed", "seq", 4688, "key", key)
	line := strings.TrimSuffix(buf.String(), "\n")
	for _, what := range leaks(line) {
		t.Errorf("the text line holds %s: %s", what, line)
	}
	if want := `key="k-2028-10 (secret hidden)"`; !strings.Contains(line, want) {
		t.Errorf("the text line = %s, want it to hold %s", line, want)
	}
}

func TestKeyThroughFmt(t *testing.T) {
	for _, verb := range []string{"%v", "%+v", "%s"} {
		if got := fmt.Sprintf(verb, key); got != "k-2028-10 (secret hidden)" {
			t.Errorf("fmt.Sprintf(%q, key) = %q, want \"k-2028-10 (secret hidden)\", as String prints it", verb, got)
		}
	}
}
