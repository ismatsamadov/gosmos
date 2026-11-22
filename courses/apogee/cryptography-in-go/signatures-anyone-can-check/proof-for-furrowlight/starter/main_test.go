package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha3"
	"encoding/hex"
	"errors"
	"testing"
	"time"
)

// The receipt keys' windows, as the customer guide lists them.
var (
	augustActivates = time.Date(2028, 9, 1, 0, 0, 0, 0, time.UTC)
	augustRetires   = time.Date(2028, 10, 17, 0, 0, 0, 0, time.UTC)
	fridayMorning   = time.Date(2028, 9, 1, 9, 0, 0, 0, time.UTC)
)

// notary returns bookingd's notary for receipt key id, whose test seed is the SHA-256 of
// "receipts <id>", readable on purpose.
func notary(t *testing.T, id string) *Notary {
	t.Helper()
	seed := sha256.Sum256([]byte("receipts " + id))
	n, err := newNotary(id, seed[:])
	if err != nil {
		t.Fatal(err)
	}
	return n
}

// published is receipt-keys.txt as the guide lists it in September: r-2028-08 alone, still signing.
func published(t *testing.T) []ReceiptKey {
	return []ReceiptKey{{ID: "r-2028-08", Public: notary(t, "r-2028-08").Published(), Activate: augustActivates}}
}

// standIn is the tests' 64 KiB stand-in for Sunday's 10.9 GB Furrow-A image.
func standIn() []byte {
	return bytes.Repeat([]byte("field 7 "), 8192)
}

// digestOf returns file's digest named for its hash, sha256 or sha3-256.
func digestOf(alg string, file []byte) string {
	var sum [32]byte
	if alg == "sha3-256" {
		sum = sha3.Sum256(file)
	} else {
		sum = sha256.Sum256(file)
	}
	return alg + ":" + hex.EncodeToString(sum[:])
}

// envelope signs, with key id, the receipt for the stand-in file under bk-5274.
func envelope(t *testing.T, id, digest string, signedAt time.Time) []byte {
	t.Helper()
	r := Receipt{
		Version: 1, KeyID: id, BookingID: "bk-5274", SatID: "furrow-a", Station: "SQT",
		File: "furrow-a-20280827-1402.tif", Bytes: int64(len(standIn())), Digest: digest,
		LandedAt: time.Date(2028, 8, 27, 14, 13, 5, 0, time.UTC), SignedAt: signedAt,
	}
	env, err := wrapReceipt(notary(t, id), r)
	if err != nil {
		t.Fatal(err)
	}
	return env
}

func TestSundaysFile(t *testing.T) {
	env := envelope(t, "r-2028-08", digestOf("sha256", standIn()), fridayMorning)
	r, err := checkDelivery(published(t), env, bytes.NewReader(standIn()))
	if err != nil || r.BookingID != "bk-5274" || r.Bytes != 65536 {
		t.Errorf("checkDelivery(Sunday's receipt, its file) = %s %d bytes, %v, want bk-5274 65536 bytes, nil", r.BookingID, r.Bytes, err)
	}
}

// TestSHA3 checks a receipt fingerprinted with SHA3-256, as Cirrine's are.
func TestSHA3(t *testing.T) {
	env := envelope(t, "r-2028-08", digestOf("sha3-256", standIn()), fridayMorning)
	if _, err := checkDelivery(published(t), env, bytes.NewReader(standIn())); err != nil {
		t.Errorf("checkDelivery(a sha3-256 receipt, its file) = %v, want nil", err)
	}
}

func TestWrongFile(t *testing.T) {
	changed := standIn()
	changed[40000] ^= 0x01
	tests := []struct {
		name string
		file []byte
	}{
		{"one byte changed", changed},
		{"one byte short", standIn()[:65535]},
		{"one byte too many", append(standIn(), '\n')},
	}
	env := envelope(t, "r-2028-08", digestOf("sha256", standIn()), fridayMorning)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := checkDelivery(published(t), env, bytes.NewReader(tt.file)); !errors.Is(err, ErrWrongFile) {
				t.Errorf("checkDelivery(Sunday's receipt, the file with %s) = %v, want ErrWrongFile", tt.name, err)
			}
		})
	}
}

// TestWindow lists r-2028-08 as the guide will from 9 October, with its retire time.
func TestWindow(t *testing.T) {
	keys := published(t)
	keys[0].Retire = augustRetires
	tests := []struct {
		name     string
		signedAt time.Time
		want     error
	}{
		{"at its activate time", augustActivates, nil},
		{"one second before it activates", augustActivates.Add(-time.Second), ErrOutsideWindow},
		{"one second before it retires", augustRetires.Add(-time.Second), nil},
		{"at its retire time", augustRetires, ErrOutsideWindow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := envelope(t, "r-2028-08", digestOf("sha256", standIn()), tt.signedAt)
			_, err := checkDelivery(keys, env, bytes.NewReader(standIn()))
			if tt.want == nil && err != nil {
				t.Errorf("checkDelivery(a receipt signed %s) = %v, want nil", tt.name, err)
			}
			if tt.want != nil && !errors.Is(err, tt.want) {
				t.Errorf("checkDelivery(a receipt signed %s) = %v, want ErrOutsideWindow", tt.name, err)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	revoked := published(t)
	revoked[0].Revoked = true
	digest := digestOf("sha256", standIn())
	forged := bytes.Replace(envelope(t, "r-2028-10", digest, fridayMorning), []byte(`"r-2028-10"`), []byte(`"r-2028-08"`), 2)
	tests := []struct {
		name     string
		keys     []ReceiptKey
		env      []byte
		want     error
		wantName string
	}{
		{"signed by r-2028-08, now revoked", revoked, envelope(t, "r-2028-08", digest, fridayMorning), ErrRevoked, "ErrRevoked"},
		{"signed by r-2028-10, not yet published", published(t), envelope(t, "r-2028-10", digest, fridayMorning), ErrUnknownSigner, "ErrUnknownSigner"},
		{"signed by r-2028-10, naming r-2028-08", published(t), forged, ErrBadSignature, "ErrBadSignature"},
		{"forged and revoked: revoked comes first", revoked, forged, ErrRevoked, "ErrRevoked"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := checkDelivery(tt.keys, tt.env, bytes.NewReader(standIn())); !errors.Is(err, tt.want) {
				t.Errorf("checkDelivery(a receipt %s) = %v, want %s", tt.name, err, tt.wantName)
			}
		})
	}
}

// TestOrder gives a receipt signed after its key retired, for a file that has changed: the
// window is checked before the file is read.
func TestOrder(t *testing.T) {
	keys := published(t)
	keys[0].Retire = augustRetires
	env := envelope(t, "r-2028-08", digestOf("sha256", standIn()), augustRetires)
	if _, err := checkDelivery(keys, env, bytes.NewReader([]byte("not the file"))); !errors.Is(err, ErrOutsideWindow) {
		t.Errorf("checkDelivery(a receipt signed at its key's retire time, another file) = %v, want ErrOutsideWindow", err)
	}
}
