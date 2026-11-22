package main

import (
	"bytes"
	"testing"
)

// The forged field report: Mirela's copy and the original share the CRC-32 f6ced6c0.
var (
	report  = []byte("NDVI field 12: 0.71\n")
	forgery = []byte("NDVI field 12: 0.17\n\x1a\x15\xbc\xd1")
)

const (
	reportSHA256 = "72198855a7d55dd3cbd5fff4d06168339fa8d32086a6cd18899400210b610aae"
	emptySHA256  = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
)

// short returns s, or its first 8 characters and an ellipsis when it is longer than 50.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

func TestReport(t *testing.T) {
	if got := fingerprint(report); got != reportSHA256 {
		t.Errorf("fingerprint(report) = %q, want %q", short(got), short(reportSHA256))
	}
}

func TestForgery(t *testing.T) {
	if got := fingerprint(forgery); got == fingerprint(report) {
		t.Errorf("fingerprint(forgery) = %q, the report's own, want a fingerprint of its own", short(got))
	}
}

func TestOneByteChanged(t *testing.T) {
	file := bytes.Repeat([]byte("Furrow-A field 12 "), 57)[:1024]
	before := fingerprint(file)
	for at := range file {
		changed := bytes.Clone(file)
		changed[at] ^= 0x20
		if got := fingerprint(changed); got == before {
			t.Errorf("fingerprint(1 KiB file with byte %d changed) = %q, as before, want a new fingerprint", at, short(got))
			return
		}
	}
}

func TestEmptyFile(t *testing.T) {
	if got := fingerprint(nil); got != emptySHA256 {
		t.Errorf("fingerprint(empty file) = %q, want %q", short(got), short(emptySHA256))
	}
}
