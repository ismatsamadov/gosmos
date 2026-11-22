package main

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

// checkDelivery checks a delivered file against the envelope of its receipt, with the receipt
// keys Apogee publishes, and returns the receipt.
func checkDelivery(keys []ReceiptKey, env []byte, file io.Reader) (Receipt, error) {
	e, err := decodeEnvelope(env)
	if err != nil {
		return Receipt{}, err
	}
	k, ok := findKey(keys, e.KeyID)
	if !ok {
		return Receipt{}, fmt.Errorf("%w: %q", ErrUnknownSigner, e.KeyID)
	}
	if k.Revoked {
		return Receipt{}, fmt.Errorf("%w: %q", ErrRevoked, k.ID)
	}
	r, err := openReceipt(k.Public, e)
	if err != nil {
		return Receipt{}, err
	}
	if r.KeyID != e.KeyID {
		return Receipt{}, fmt.Errorf("%w: receipt names key %q", ErrBadEnvelope, r.KeyID)
	}
	if r.SignedAt.Before(k.Activate) || !k.Retire.IsZero() && !r.SignedAt.Before(k.Retire) {
		return Receipt{}, fmt.Errorf("%w: %s", ErrOutsideWindow, r.SignedAt.Format(time.RFC3339))
	}
	if err := checkFile(r, file); err != nil {
		return Receipt{}, err
	}
	return r, nil
}

// checkFile hashes file with the receipt's own hash and returns ErrWrongFile unless its size
// and digest are the receipt's.
func checkFile(r Receipt, file io.Reader) error {
	alg, sum, err := parseDigest(r.Digest)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBadEnvelope, err)
	}
	h, err := newHash(alg)
	if err != nil {
		return err
	}
	n, err := io.Copy(h, file)
	if err != nil {
		return err
	}
	if n != r.Bytes || !bytes.Equal(h.Sum(nil), sum) {
		return ErrWrongFile
	}
	return nil
}
