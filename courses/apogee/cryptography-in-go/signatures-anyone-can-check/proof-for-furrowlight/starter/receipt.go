package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

// Receipt is Apogee's signed word that a booking's file was delivered.
type Receipt struct {
	Version   int       `json:"v"`          // 1
	KeyID     string    `json:"key_id"`     // "r-2028-08": the key that signed it
	BookingID string    `json:"booking_id"` // "bk-5274"
	SatID     string    `json:"sat_id"`
	Station   string    `json:"station"`
	File      string    `json:"file"` // "furrow-a-20280827-1402.tif"
	Bytes     int64     `json:"bytes"`
	Digest    string    `json:"digest"`    // "sha256:<64 hex>", or "sha3-256:" for Cirrine
	LandedAt  time.Time `json:"landed_at"` // from the Landing (Season 3)
	SignedAt  time.Time `json:"signed_at"` // bookingd's clock, truncated to the second
}

// Envelope carries a receipt and its signature.
type Envelope struct {
	KeyID   string          `json:"key_id"`  // the same as the receipt's
	Receipt json.RawMessage `json:"receipt"` // the exact bytes that were signed
	Sig     string          `json:"sig"`     // 128 hex digits: Ed25519 over Receipt's bytes
}

// ReceiptKey is one published receipt key and its signing window.
type ReceiptKey struct {
	ID       string            // "r-2028-08"
	Public   ed25519.PublicKey // 32 bytes
	Activate time.Time         // receipts signed from here...
	Retire   time.Time         // ...until here (zero: still signing)
	Revoked  bool
}

// The ways a receipt is refused.
var (
	ErrBadEnvelope   = errors.New("malformed receipt")
	ErrUnknownSigner = errors.New("unknown receipt key")
	ErrBadSignature  = errors.New("receipt signature does not verify")
	ErrOutsideWindow = errors.New("receipt signed outside its key's window")
	ErrWrongFile     = errors.New("file does not match receipt")
	ErrRevoked       = errors.New("key revoked") // as for sealed commands
)

// Notary is bookingd's receipt signer.
type Notary struct {
	KeyID string
	priv  ed25519.PrivateKey
}

// newNotary returns the notary for receipt key keyID, made from its 32-byte seed.
func newNotary(keyID string, seed []byte) (*Notary, error) {
	if len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("seed has %d bytes, want %d", len(seed), ed25519.SeedSize)
	}
	return &Notary{KeyID: keyID, priv: ed25519.NewKeyFromSeed(seed)}, nil
}

// Sign returns the notary's signature over msg.
func (n *Notary) Sign(msg []byte) []byte {
	return ed25519.Sign(n.priv, msg)
}

// Published returns what customers are given to check receipts with: the 32-byte public key.
func (n *Notary) Published() []byte {
	return n.priv.Public().(ed25519.PublicKey)
}

// wrapReceipt signs r's JSON with n and returns the envelope that carries both.
func wrapReceipt(n *Notary, r Receipt) ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return json.Marshal(Envelope{KeyID: n.KeyID, Receipt: b, Sig: hex.EncodeToString(n.Sign(b))})
}

// checkSig returns ErrBadSignature unless sigHex is 128 lowercase hex digits of pub's signature
// over msg.
func checkSig(pub ed25519.PublicKey, msg []byte, sigHex string) error {
	if len(sigHex) != 2*ed25519.SignatureSize || strings.ToLower(sigHex) != sigHex || len(pub) != ed25519.PublicKeySize {
		return ErrBadSignature
	}
	sig, err := hex.DecodeString(sigHex)
	if err != nil || !ed25519.Verify(pub, msg, sig) {
		return ErrBadSignature
	}
	return nil
}

// decodeStrict decodes the one JSON value in data into v, refusing fields v doesn't have and
// anything after the value.
func decodeStrict(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return errors.New("more data after the value")
	}
	return nil
}

// decodeEnvelope reads an envelope strictly: 64 KiB at most, and no field Envelope doesn't have.
func decodeEnvelope(env []byte) (Envelope, error) {
	if len(env) > 64<<10 {
		return Envelope{}, fmt.Errorf("%w: over 64 KiB", ErrBadEnvelope)
	}
	var e Envelope
	if err := decodeStrict(env, &e); err != nil {
		return Envelope{}, fmt.Errorf("%w: %v", ErrBadEnvelope, err)
	}
	return e, nil
}

// openReceipt checks e's signature with pub over the receipt's exact bytes, then decodes them
// strictly.
func openReceipt(pub ed25519.PublicKey, e Envelope) (Receipt, error) {
	if err := checkSig(pub, e.Receipt, e.Sig); err != nil {
		return Receipt{}, err
	}
	var r Receipt
	if err := decodeStrict(e.Receipt, &r); err != nil {
		return Receipt{}, fmt.Errorf("%w: %v", ErrBadEnvelope, err)
	}
	return r, nil
}

// findKey returns the published key named id, if keys has one.
func findKey(keys []ReceiptKey, id string) (ReceiptKey, bool) {
	for _, k := range keys {
		if k.ID == id {
			return k, true
		}
	}
	return ReceiptKey{}, false
}
