package main

import (
	"bytes"
	"crypto/mlkem"
	"encoding/hex"
	"testing"
)

// short returns the first 8 hex digits of b and an ellipsis.
func short(b []byte) string {
	if len(b) < 4 {
		return hex.EncodeToString(b)
	}
	return hex.EncodeToString(b[:4]) + "…"
}

// receiver makes Ystmark's ML-KEM-768 decapsulation key for one test.
func receiver(t *testing.T) *mlkem.DecapsulationKey768 {
	t.Helper()
	dk, err := mlkem.GenerateKey768()
	if err != nil {
		t.Fatalf("mlkem.GenerateKey768() = %v", err)
	}
	return dk
}

// wrapped runs wrapFor for dk's encapsulation key, stopping the test unless it gives a key and a ciphertext.
func wrapped(t *testing.T, dk *mlkem.DecapsulationKey768) (key, ciphertext []byte) {
	t.Helper()
	ek := dk.EncapsulationKey().Bytes()
	key, ciphertext, err := wrapFor(ek)
	if err != nil || len(key) != 32 || len(ciphertext) != 1088 {
		t.Fatalf("wrapFor(a %d-byte encapsulation key) = %d-byte key, %d-byte ciphertext, %v, want 32, 1088, nil", len(ek), len(key), len(ciphertext), err)
	}
	return key, ciphertext
}

// TestRoundTrip: Port Ellery wraps a key for Ystmark, and Ystmark unwraps the same 32 bytes.
func TestRoundTrip(t *testing.T) {
	dk := receiver(t)
	key, ciphertext := wrapped(t, dk)
	got, err := unwrap(dk, ciphertext)
	if err != nil || !bytes.Equal(got, key) {
		t.Errorf("unwrap(dk, the ciphertext) = %s, %v, want %s, the key wrapFor made", short(got), err, short(key))
	}
}

func TestTwoWrapsDiffer(t *testing.T) {
	dk := receiver(t)
	first, _ := wrapped(t, dk)
	second, _ := wrapped(t, dk)
	if bytes.Equal(first, second) {
		t.Errorf("wrapFor(the same key) = %s twice, want a new key every time", short(first))
	}
}

func TestBadEncapsulationKey(t *testing.T) {
	ek := receiver(t).EncapsulationKey().Bytes()
	if _, _, err := wrapFor(ek[:1183]); err == nil {
		t.Errorf("wrapFor(a 1183-byte encapsulation key) = nil error, want an error")
	}
}

func TestShortCiphertext(t *testing.T) {
	if got, err := unwrap(receiver(t), make([]byte, 100)); err == nil {
		t.Errorf("unwrap(dk, a 100-byte ciphertext) = %d bytes, nil, want an error", len(got))
	}
}

// TestChangedCiphertext changes one byte on the way: no error, but a key the sender never made.
func TestChangedCiphertext(t *testing.T) {
	dk := receiver(t)
	key, ciphertext := wrapped(t, dk)
	ciphertext[500] ^= 0x01
	got, err := unwrap(dk, ciphertext)
	if err != nil || len(got) != 32 || bytes.Equal(got, key) {
		t.Errorf("unwrap(dk, the ciphertext with byte 500 changed) = %s, %v, want 32 bytes other than %s, nil", short(got), err, short(key))
	}
}
