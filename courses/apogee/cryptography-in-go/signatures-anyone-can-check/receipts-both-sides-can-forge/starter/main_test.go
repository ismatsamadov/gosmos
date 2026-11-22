package main

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// r202808Public is the public key of the test receipt key r-2028-08, whose seed is the SHA-256
// of "receipts r-2028-08".
const r202808Public = "ba6c6226d2c3c056684e9d9a5e489db44a148e27e3fcb8f2ba16a6bddb0a86b2"

// The receipt for Sunday's Furrow-A file, and the one Furrowlight would like to show.
const (
	receipt = "bk-5274 furrow-a-20280827-1402.tif landed 2028-08-27T14:13:05Z"
	dayLate = "bk-5274 furrow-a-20280827-1402.tif landed 2028-08-28T14:13:05Z"
)

// seedOf returns the SHA-256 of label: a test seed, readable on purpose.
func seedOf(label string) []byte {
	sum := sha256.Sum256([]byte(label))
	return sum[:]
}

// short returns s, or its first 8 characters and an ellipsis when it is longer than 50.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

// testNotary returns bookingd's notary for the test key r-2028-08.
func testNotary(t *testing.T) *Notary {
	t.Helper()
	n, err := newNotary("r-2028-08", seedOf("receipts r-2028-08"))
	if err != nil {
		t.Fatalf("newNotary(r-2028-08, a 32-byte seed) = %v, want nil", err)
	}
	return n
}

func TestApogeeSigns(t *testing.T) {
	n := testNotary(t)
	if got := hex.EncodeToString(n.Published()); got != r202808Public {
		t.Errorf("Published() = %s, want r-2028-08's public key %s", short(got), short(r202808Public))
	}
	if !checkReceipt(n.Published(), []byte(receipt), n.Sign([]byte(receipt))) {
		t.Errorf("checkReceipt(Published(), receipt, Sign(receipt)) = false, want true")
	}
}

// TestFurrowlightForges plays Furrowlight, which holds everything a customer is given: what
// Published returns. It signs the day-late receipt every way those bytes could allow.
func TestFurrowlightForges(t *testing.T) {
	published := testNotary(t).Published()
	mac := hmac.New(sha256.New, published)
	mac.Write([]byte(dayLate))
	if checkReceipt(published, []byte(dayLate), mac.Sum(nil)) {
		t.Errorf("checkReceipt(the day-late receipt, an HMAC under Published()) = true, want false")
	}
	if len(published) == ed25519.SeedSize {
		forger := ed25519.NewKeyFromSeed(published)
		if checkReceipt(published, []byte(dayLate), ed25519.Sign(forger, []byte(dayLate))) {
			t.Errorf("checkReceipt(the day-late receipt, signed with Published() as a seed) = true, want false")
		}
	}
}

func TestChangedReceipt(t *testing.T) {
	n := testNotary(t)
	if checkReceipt(n.Published(), []byte(dayLate), n.Sign([]byte(receipt))) {
		t.Errorf("checkReceipt(the day-late receipt, Apogee's signature of the real one) = true, want false")
	}
}

func TestShortKey(t *testing.T) {
	n := testNotary(t)
	sig := n.Sign([]byte(receipt))
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("checkReceipt(a 31-byte key, ...) panicked: %v, want false", r)
		}
	}()
	if checkReceipt(n.Published()[:31], []byte(receipt), sig) {
		t.Errorf("checkReceipt(a 31-byte key, ...) = true, want false")
	}
}

func TestShortSeed(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("newNotary(r-2028-08, a 31-byte seed) panicked: %v, want an error", r)
		}
	}()
	if _, err := newNotary("r-2028-08", seedOf("receipts r-2028-08")[:31]); err == nil {
		t.Errorf("newNotary(r-2028-08, a 31-byte seed) = nil error, want an error")
	}
}
