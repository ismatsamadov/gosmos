package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

// frames are three frames as the bench's radio hands them to uplinkd.
var frames = [][]byte{
	{0x0b, 0x03, 0x37, 0xd1, 0x00, 0x2a},
	{0x0b, 0x03, 0x37, 0xd2, 0x00, 0x2b},
	{0x0b, 0x03, 0x37, 0xd3, 0x00, 0x2c},
}

// hexDump stands in for the costly dump: it counts its calls and writes a frame as hex.
type hexDump struct{ calls int }

// dump counts one call and returns frame as hex.
func (d *hexDump) dump(frame []byte) string {
	d.calls++
	return hex.EncodeToString(frame)
}

// lines splits what a logger wrote into its lines.
func lines(out string) []string {
	if out == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(out, "\n"), "\n")
}

func TestLogFrame(t *testing.T) {
	ctx := context.Background()
	t.Run("at INFO, nothing is built", func(t *testing.T) {
		var buf bytes.Buffer
		l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{ReplaceAttr: dropTime}))
		d := &hexDump{}
		for _, f := range frames {
			logFrame(ctx, l, f, d.dump)
		}
		if d.calls != 0 {
			t.Errorf("three frames at INFO called dump %d times, want 0", d.calls)
		}
		if buf.Len() != 0 {
			t.Errorf("three frames at INFO wrote %q, want nothing", buf.String())
		}
	})
	t.Run("at DEBUG, every frame", func(t *testing.T) {
		var buf bytes.Buffer
		l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug, ReplaceAttr: dropTime}))
		d := &hexDump{}
		for _, f := range frames {
			logFrame(ctx, l, f, d.dump)
		}
		if d.calls != 3 {
			t.Errorf("three frames at DEBUG called dump %d times, want 3", d.calls)
		}
		got := lines(buf.String())
		if len(got) != 3 {
			t.Fatalf("three frames at DEBUG wrote %d lines, want 3", len(got))
		}
		var rec map[string]any
		if err := json.Unmarshal([]byte(got[0]), &rec); err != nil || rec["msg"] != "frame received" || rec["bytes"] != 6.0 || rec["hex"] != "0b0337d1002a" {
			t.Errorf(`the first frame at DEBUG wrote %s, want "frame received" with "bytes":6 and "hex":"0b0337d1002a"`, got[0])
		}
	})
	t.Run("turned up mid-pass", func(t *testing.T) {
		var v slog.LevelVar
		var buf bytes.Buffer
		l := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: &v, ReplaceAttr: dropTime}))
		d := &hexDump{}
		logFrame(ctx, l, frames[0], d.dump)
		logFrame(ctx, l, frames[1], d.dump)
		v.Set(slog.LevelDebug)
		logFrame(ctx, l, frames[2], d.dump)
		if d.calls != 1 {
			t.Errorf("two frames at INFO, then one at DEBUG, called dump %d times, want 1", d.calls)
		}
		if got := lines(buf.String()); len(got) != 1 || !strings.Contains(got[0], `"hex":"0b0337d3002c"`) {
			t.Errorf("two frames at INFO, then one at DEBUG, wrote %q, want the third frame's line only", buf.String())
		}
	})
}
