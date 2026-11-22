package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// testNotary returns bookingd's notary for the test receipt key r-2028-08, whose seed is the
// SHA-256 of a label, readable on purpose.
func testNotary(t *testing.T) *Notary {
	t.Helper()
	seed := sha256.Sum256([]byte("receipts r-2028-08"))
	n, err := newNotary("r-2028-08", seed[:])
	if err != nil {
		t.Fatal(err)
	}
	return n
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

// apogeeEnvelope returns the envelope bookingd sends with sunday's receipt.
func apogeeEnvelope(t *testing.T, n *Notary) []byte {
	t.Helper()
	env, err := wrapReceipt(n, sunday())
	if err != nil {
		t.Fatal(err)
	}
	return env
}

func TestApogeeEnvelope(t *testing.T) {
	n := testNotary(t)
	r, err := openEnvelope(n.Published(), apogeeEnvelope(t, n))
	want := sunday()
	if err != nil || r.BookingID != want.BookingID || r.File != want.File || r.Bytes != want.Bytes ||
		r.Digest != want.Digest || !r.LandedAt.Equal(want.LandedAt) || !r.SignedAt.Equal(want.SignedAt) {
		t.Errorf("openEnvelope(Apogee's envelope) = %+v, %v, want %+v, nil", r, err, want)
	}
}

// TestRefund plays Mirela: she adds "refund":true inside Apogee's receipt and passes it on.
func TestRefund(t *testing.T) {
	n := testNotary(t)
	refund := bytes.Replace(apogeeEnvelope(t, n), []byte(`"v":1,`), []byte(`"v":1,"refund":true,`), 1)
	_, err := openEnvelope(n.Published(), refund)
	if err == nil {
		t.Errorf("openEnvelope(Apogee's receipt with \"refund\":true added) = nil error, want an error")
	} else if !errors.Is(err, ErrBadSignature) {
		t.Errorf("openEnvelope(Apogee's receipt with \"refund\":true added) = %v, want ErrBadSignature: check the signature before decoding", err)
	}
}

// TestRespaced passes Apogee's receipt on with a space after every comma: the same fields,
// other bytes.
func TestRespaced(t *testing.T) {
	n := testNotary(t)
	respaced := bytes.ReplaceAll(apogeeEnvelope(t, n), []byte(`,"`), []byte(`, "`))
	if _, err := openEnvelope(n.Published(), respaced); !errors.Is(err, ErrBadSignature) {
		t.Errorf("openEnvelope(Apogee's receipt, respaced) = %v, want ErrBadSignature: Apogee signed other bytes", err)
	}
}

func TestShortSig(t *testing.T) {
	n := testNotary(t)
	var e Envelope
	if err := json.Unmarshal(apogeeEnvelope(t, n), &e); err != nil {
		t.Fatal(err)
	}
	e.Sig = e.Sig[:127]
	env, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := openEnvelope(n.Published(), env); !errors.Is(err, ErrBadSignature) {
		t.Errorf("openEnvelope(a sig of 127 hex digits) = %v, want ErrBadSignature", err)
	}
}

// TestUnknownField signs, with r-2028-08's test key, a receipt holding a field Receipt doesn't
// have: the signature is good, and strict decoding must still refuse it.
func TestUnknownField(t *testing.T) {
	n := testNotary(t)
	b, err := json.Marshal(sunday())
	if err != nil {
		t.Fatal(err)
	}
	extra := bytes.Replace(b, []byte(`"v":1,`), []byte(`"v":1,"refund":true,`), 1)
	env, err := json.Marshal(Envelope{KeyID: "r-2028-08", Receipt: extra, Sig: hex.EncodeToString(n.Sign(extra))})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := openEnvelope(n.Published(), env); !errors.Is(err, ErrBadEnvelope) {
		t.Errorf("openEnvelope(a signed receipt with a refund field) = %v, want ErrBadEnvelope: decode strictly", err)
	}
}
