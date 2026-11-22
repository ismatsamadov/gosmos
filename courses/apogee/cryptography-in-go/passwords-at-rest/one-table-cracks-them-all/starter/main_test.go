package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

// accounts are three of the console's accounts as the test knows them: two with passwords
// from the common list, and Mirela's, which isn't on it.
var accounts = []struct{ user, pw string }{
	{"kiosk", "fishbowl"},
	{"visitor", "welcome1"},
	{"mirela", "fails closed, always"},
}

// commonPasswords returns the attacker's list: 1,000 passwords people choose, from 40 words
// and 25 endings.
func commonPasswords() []string {
	words := strings.Fields(`password letmein welcome qwerty dragon monkey sunshine princess
		football baseball shadow master secret admin login iloveyou flower summer winter autumn
		orange banana cheese coffee pepper ginger rocket orbit satellite antenna fishbowl apogee
		pipit ellery harbor seagull lighthouse starlight station uplink`)
	endings := []string{"", "1", "2", "12", "123", "1234", "!", "!!", "01", "07", "11", "13",
		"22", "99", "00", "007", "321", "2025", "2026", "2027", "2028", "?", "#", "1!", "123!"}
	var list []string
	for _, w := range words {
		for _, e := range endings {
			list = append(list, w+e)
		}
	}
	return list
}

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

// TestStolenTable plays whoever has the laptop. Long before Sunday they made one table of
// 1,000 common passwords and their SHA-256 digests, and they look up every field of every
// record in it.
func TestStolenTable(t *testing.T) {
	table := make(map[string]string)
	for _, pw := range commonPasswords() {
		sum := sha256.Sum256([]byte(pw))
		table[hex.EncodeToString(sum[:])] = pw
	}
	for _, a := range accounts {
		t.Run(a.user, func(t *testing.T) {
			record := hashPassword(a.pw)
			for _, field := range strings.Split(record, "$") {
				if pw, found := table[field]; found {
					t.Errorf("hashPassword(%q) = %q, which the table reads as %q at once, want a record no table made in advance can match", a.pw, short(record), pw)
				}
			}
		})
	}
}

func TestSaltedRecord(t *testing.T) {
	const pw = "fishbowl"
	record := hashPassword(pw)
	fields := strings.Split(record, "$")
	if len(fields) != 3 || fields[0] != "salted" || len(fields[1]) != 32 || len(fields[2]) != 64 {
		t.Fatalf("hashPassword(%q) = %q, want salted$<32 hex salt>$<64 hex digest>", pw, short(record))
	}
	salt, err := hex.DecodeString(fields[1])
	if err != nil {
		t.Fatalf("hashPassword(%q) has the salt %q, want 32 hex digits", pw, fields[1])
	}
	sum := sha256.Sum256(append(salt, pw...))
	if want := hex.EncodeToString(sum[:]); fields[2] != want {
		t.Errorf("hashPassword(%q) has the digest %s…, want SHA-256 of the salt, then the password: %s…", pw, fields[2][:8], want[:8])
	}
}

func TestSamePasswordTwice(t *testing.T) {
	const pw = "fishbowl"
	a, b := hashPassword(pw), hashPassword(pw)
	if a == b {
		t.Errorf("hashPassword(%q) = %q twice in a row, want two different records", pw, short(a))
	}
}
