package main

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

const (
	saltHex = "6d6972656c6120746573742073616c74" // "mirela test salt", the tests' readable salt
	keyHex  = "6d91ceab779ff574b23820b133a740fa76d8b6a0deb5d6f0da2f94a6508dd938"
	// mirelaRecord is Mirela's record: her password at 600,000 iterations with the test salt.
	mirelaRecord = "pbkdf2-sha256$600000$" + saltHex + "$" + keyHex
	// kioskOld is the kiosk's record from the stolen file: the unsalted SHA-256 of "fishbowl".
	kioskOld = "sha256$1c4bb463c52d458ae6e9225430c45c753c5ae417bb4e7e4530773aaa0cda201b"
)

// short cuts each long field of a record to its first 8 characters and an ellipsis.
func short(record string) string {
	if len(record) <= 50 {
		return record
	}
	fields := strings.Split(record, "$")
	for i, f := range fields {
		if len(f) > 16 {
			fields[i] = f[:8] + "…"
		}
	}
	return strings.Join(fields, "$")
}

// show writes a Record in one short line.
func show(r Record) string {
	return fmt.Sprintf("{Iterations: %d, Salt: %d bytes, Key: %d bytes}", r.Iterations, len(r.Salt), len(r.Key))
}

func TestFormatRecord(t *testing.T) {
	key, _ := hex.DecodeString(keyHex)
	r := Record{Iterations: 600000, Salt: []byte("mirela test salt"), Key: key}
	if got := formatRecord(r); got != mirelaRecord {
		t.Errorf("formatRecord(Mirela's record) = %q, want %q", short(got), short(mirelaRecord))
	}
}

func TestParseRecord(t *testing.T) {
	tests := []struct {
		name, record string
		iterations   int
	}{
		{"Mirela's record", mirelaRecord, 600000},
		{"the lowest count, 1,000", "pbkdf2-sha256$1000$" + saltHex + "$" + keyHex, 1000},
		{"the highest count, 10,000,000", "pbkdf2-sha256$10000000$" + saltHex + "$" + keyHex, 10000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := parseRecord(tt.record)
			if err != nil || r.Iterations != tt.iterations || string(r.Salt) != "mirela test salt" || hex.EncodeToString(r.Key) != keyHex {
				t.Errorf("parseRecord(%q) = %s, %v, want {Iterations: %d, Salt: 16 bytes, Key: 32 bytes} holding the record's salt and key, nil", short(tt.record), show(r), err, tt.iterations)
			}
		})
	}
}

func TestParseRecordRefuses(t *testing.T) {
	tests := []struct{ name, record string }{
		{"the kiosk's old record", kioskOld},
		{"a count of 0", "pbkdf2-sha256$0$" + saltHex + "$" + keyHex},
		{"a count of 999", "pbkdf2-sha256$999$" + saltHex + "$" + keyHex},
		{"a count of 10,000,001", "pbkdf2-sha256$10000001$" + saltHex + "$" + keyHex},
		{"a 30-digit salt", "pbkdf2-sha256$600000$" + saltHex[:30] + "$" + keyHex},
		{"a 62-digit key", "pbkdf2-sha256$600000$" + saltHex + "$" + keyHex[:62]},
		{"a missing field", "pbkdf2-sha256$600000$" + keyHex},
		{"a fifth field", mirelaRecord + "$admin"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if r, err := parseRecord(tt.record); err == nil {
				t.Errorf("parseRecord(%q) = %s, nil, want an error", short(tt.record), show(r))
			}
		})
	}
}
