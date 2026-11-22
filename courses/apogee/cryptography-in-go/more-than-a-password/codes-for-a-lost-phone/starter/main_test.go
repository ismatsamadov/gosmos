package main

import (
	"slices"
	"strings"
	"testing"
)

// testRecovery returns a store that holds no recovery codes yet.
func testRecovery() *Recovery {
	return &Recovery{unused: map[string][]string{}}
}

// short cuts a value longer than 50 characters to its first 8 and an ellipsis.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

// TestSheet checks a new sheet: eight different codes from rand.Text, and in the console only
// the digest of each.
func TestSheet(t *testing.T) {
	r := testRecovery()
	sheet := r.Issue("sunniva")
	if len(sheet) != sheetSize {
		t.Fatalf("Issue(%q) returned %d codes, want %d", "sunniva", len(sheet), sheetSize)
	}
	digests := map[string]bool{}
	for _, code := range sheet {
		if len(code) != 26 || strings.Trim(code, "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567") != "" {
			t.Errorf("Issue(%q) gave the code %q, want 26 characters from A to Z and 2 to 7, as rand.Text makes", "sunniva", code)
		}
		if digests[tokenDigest(code)] {
			t.Errorf("Issue(%q) gave the code %q twice, want eight different codes", "sunniva", code)
		}
		digests[tokenDigest(code)] = true
	}
	kept := r.unused["sunniva"]
	if len(kept) != sheetSize {
		t.Fatalf("r.unused[%q] holds %d values after Issue, want %d", "sunniva", len(kept), sheetSize)
	}
	for _, k := range kept {
		if !digests[k] {
			t.Errorf("r.unused[%q] holds %q, want only the tokenDigest of each code on the sheet", "sunniva", short(k))
		}
	}
}

// TestStolenCopy plays whoever copies the console's file of recovery codes, and types each
// value in it as a code.
func TestStolenCopy(t *testing.T) {
	r := testRecovery()
	r.Issue("sunniva")
	for _, kept := range slices.Clone(r.unused["sunniva"]) {
		if r.Use("sunniva", kept) {
			t.Errorf("Use(%q, %q) = true, want false: that value was copied from the console's file", "sunniva", short(kept))
		}
	}
}

// TestOnce uses every code on Sunniva's sheet: each works once, in any order.
func TestOnce(t *testing.T) {
	r := testRecovery()
	sheet := r.Issue("sunniva")
	if !r.Use("sunniva", sheet[3]) {
		t.Fatalf("Use(%q, %q), the fourth code on the sheet, = false, want true", "sunniva", sheet[3])
	}
	if r.Use("sunniva", sheet[3]) {
		t.Errorf("Use(%q, %q) a second time = true, want false: a recovery code works once", "sunniva", sheet[3])
	}
	for i, code := range sheet {
		if i != 3 && !r.Use("sunniva", code) {
			t.Errorf("Use(%q, %q), code %d on the sheet, = false, want true", "sunniva", code, i+1)
		}
	}
}

// TestOtherAccounts checks that a code opens only its own account.
func TestOtherAccounts(t *testing.T) {
	r := testRecovery()
	sunniva := r.Issue("sunniva")
	r.Issue("mirela")
	if r.Use("mirela", sunniva[0]) {
		t.Errorf("Use(%q, %q), Sunniva's code, = true, want false", "mirela", sunniva[0])
	}
	if r.Use("nkem", sunniva[1]) {
		t.Errorf("Use(%q, %q), Sunniva's code for an account with no sheet, = true, want false", "nkem", sunniva[1])
	}
	if !r.Use("sunniva", sunniva[0]) {
		t.Errorf("Use(%q, %q) = false, want true: trying it on Mirela's account must not spend it", "sunniva", sunniva[0])
	}
}

// TestNewSheet prints Sunniva a second sheet: the first one stops working.
func TestNewSheet(t *testing.T) {
	r := testRecovery()
	old := r.Issue("sunniva")
	fresh := r.Issue("sunniva")
	if r.Use("sunniva", old[0]) {
		t.Errorf("Use(%q, %q), from the first sheet, = true, want false after a new sheet", "sunniva", old[0])
	}
	if !r.Use("sunniva", fresh[0]) {
		t.Errorf("Use(%q, %q), from the new sheet, = false, want true", "sunniva", fresh[0])
	}
}
