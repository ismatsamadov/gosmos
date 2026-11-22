package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

// payloadText returns a raw payload's bytes as text, and prints any other payload.
func payloadText(payload any) string {
	if raw, ok := payload.(json.RawMessage); ok {
		return string(raw)
	}
	return fmt.Sprint(payload)
}

func TestDecodeEnvelope(t *testing.T) {
	for _, tc := range []struct {
		name, report, kind string
		id                 int64
		payload            string
	}{
		{"image report",
			`{"version":2,"kind":"image","id":9007199254741003,"payload":{"sat_id": "furrow-a", "station": "SQT", "chunks": 2000, "missing": 9}}`,
			"image", 9007199254741003, `{"sat_id": "furrow-a", "station": "SQT", "chunks": 2000, "missing": 9}`},
		{"station report",
			`{"version":2,"kind":"station","id":9007199254741005,"payload":{"station":"SQT","day":"2027-08-24","passes":6,"missing":31}}`,
			"station", 9007199254741005, `{"station":"SQT","day":"2027-08-24","passes":6,"missing":31}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			env, err := decodeEnvelope([]byte(tc.report))
			if err != nil {
				t.Fatalf("decodeEnvelope(%s) error = %v, want nil", tc.name, err)
			}
			if env.Version != 2 || env.Kind != tc.kind || env.ID != tc.id {
				t.Errorf("decodeEnvelope(%s) = version %d, kind %q, ID %d, want 2, %q, %d", tc.name, env.Version, env.Kind, env.ID, tc.kind, tc.id)
			}
			if got := payloadText(env.Payload); got != tc.payload {
				t.Errorf("decodeEnvelope(%s).Payload = %q, want %q, as it arrived", tc.name, got, tc.payload)
			}
		})
	}
}

func TestDecodeEnvelopeRefuses(t *testing.T) {
	for _, tc := range []struct{ name, report string }{
		{"cut off", `{"version":2,"kind":"image","id":9007199254741003,"payload":{"sat_id":"fur`},
		{"version as text", `{"version":"2","kind":"image","id":9007199254741003,"payload":{}}`},
		{"no payload", `{"version":2,"kind":"station","id":9007199254741005}`},
		{"null payload", `{"version":2,"kind":"station","id":9007199254741005,"payload":null}`},
		{"a list", `[2,"image",9007199254741003]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := decodeEnvelope([]byte(tc.report)); err == nil {
				t.Errorf("decodeEnvelope(%s) error = nil, want an error", tc.name)
			}
		})
	}
}
