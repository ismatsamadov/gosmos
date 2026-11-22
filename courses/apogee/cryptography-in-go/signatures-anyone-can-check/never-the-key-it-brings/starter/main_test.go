package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"
)

// notary returns a notary for key id whose seed is the SHA-256 of label, readable on purpose.
func notary(t *testing.T, id, label string) *Notary {
	t.Helper()
	seed := sha256.Sum256([]byte(label))
	n, err := newNotary(id, seed[:])
	if err != nil {
		t.Fatal(err)
	}
	return n
}

// published is what the customer guide lists: the test key r-2028-08's public key, by its ID.
func published(t *testing.T) map[string]ed25519.PublicKey {
	return map[string]ed25519.PublicKey{"r-2028-08": notary(t, "r-2028-08", "receipts r-2028-08").Published()}
}

// sunday is the receipt for Sunday's Furrow-A file.
func sunday() Receipt {
	return Receipt{
		Version: 1, KeyID: "r-2028-08", BookingID: "bk-5274", SatID: "furrow-a", Station: "SQT",
		File: "furrow-a-20280827-1402.tif", Bytes: 10874220000,
		Digest:   "sha256:57f31a45f7f99d89febc9b563983e20056006094293fb00e8d7a03ef78224ce9",
		LandedAt: time.Date(2028, 8, 27, 14, 13, 5, 0, time.UTC),
		SignedAt: time.Date(2028, 9, 1, 9, 0, 0, 0, time.UTC),
	}
}

// wrap signs r with n, as bookingd does.
func wrap(t *testing.T, n *Notary, r Receipt) []byte {
	t.Helper()
	env, err := wrapReceipt(n, r)
	if err != nil {
		t.Fatal(err)
	}
	return env
}

func TestApogeeEnvelope(t *testing.T) {
	env := wrap(t, notary(t, "r-2028-08", "receipts r-2028-08"), sunday())
	r, err := verifyEnvelope(published(t), env)
	if err != nil || r.KeyID != "r-2028-08" || r.BookingID != "bk-5274" || !r.LandedAt.Equal(sunday().LandedAt) {
		t.Errorf("verifyEnvelope(Apogee's envelope) = %s %s landed %v, %v, want r-2028-08 bk-5274 landed %v, nil",
			r.KeyID, r.BookingID, r.LandedAt, err, sunday().LandedAt)
	}
}

// TestKeyItBrings plays Mirela: she signs a receipt landing a day late with a key of her own,
// names r-2028-08 in it, and puts her public key in the envelope as "pub".
func TestKeyItBrings(t *testing.T) {
	mirela := notary(t, "r-2028-08", "a key of Mirela's own")
	late := sunday()
	late.LandedAt = late.LandedAt.Add(24 * time.Hour)
	env := wrap(t, mirela, late)
	withKey := bytes.Replace(env, []byte(`{"key_id"`), []byte(`{"pub":"`+hex.EncodeToString(mirela.Published())+`","key_id"`), 1)
	if r, err := verifyEnvelope(published(t), withKey); err == nil {
		t.Errorf("verifyEnvelope(Mirela's receipt, carrying her key) = landed %v, nil, want an error", r.LandedAt)
	} else if !errors.Is(err, ErrBadEnvelope) {
		t.Errorf("verifyEnvelope(Mirela's receipt, carrying her key) = %v, want ErrBadEnvelope: Envelope has no field pub", err)
	}
	if _, err := verifyEnvelope(published(t), env); !errors.Is(err, ErrBadSignature) {
		t.Errorf("verifyEnvelope(Mirela's receipt, without her key) = %v, want ErrBadSignature", err)
	}
}

func TestUnknownSigner(t *testing.T) {
	r := sunday()
	r.KeyID = "r-2027-01"
	env := wrap(t, notary(t, "r-2027-01", "receipts r-2027-01"), r)
	if _, err := verifyEnvelope(published(t), env); !errors.Is(err, ErrUnknownSigner) {
		t.Errorf("verifyEnvelope(a receipt signed by r-2027-01) = %v, want ErrUnknownSigner", err)
	}
}

// TestInnerKeyID signs with r-2028-08's test key a receipt that names r-2028-10 inside.
func TestInnerKeyID(t *testing.T) {
	r := sunday()
	r.KeyID = "r-2028-10"
	env := wrap(t, notary(t, "r-2028-08", "receipts r-2028-08"), r)
	if _, err := verifyEnvelope(published(t), env); !errors.Is(err, ErrBadEnvelope) {
		t.Errorf("verifyEnvelope(an envelope from r-2028-08 holding a receipt that names r-2028-10) = %v, want ErrBadEnvelope", err)
	}
}
