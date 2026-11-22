package main

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// ephemeral returns the test X25519 key whose private bytes are the SHA-256 of label.
func ephemeral(t *testing.T, label string) *ecdh.PrivateKey {
	t.Helper()
	seed := sha256.Sum256([]byte(label))
	key, err := ecdh.X25519().NewPrivateKey(seed[:])
	if err != nil {
		t.Fatalf("ecdh.X25519().NewPrivateKey(SHA-256 of %q) = %v", label, err)
	}
	return key
}

// short returns b in hex, cut to its first 8 digits and an ellipsis when longer than 50.
func short(b []byte) string {
	if s := hex.EncodeToString(b); len(s) <= 50 {
		return s
	}
	return hex.EncodeToString(b[:4]) + "…"
}

// TestBothWays makes two fresh key pairs, as Port Ellery and Ystmark do, and agrees from each end.
func TestBothWays(t *testing.T) {
	pel, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("ecdh.X25519().GenerateKey = %v", err)
	}
	yst, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("ecdh.X25519().GenerateKey = %v", err)
	}
	atPEL, err := agree(pel, yst.PublicKey().Bytes())
	if err != nil || len(atPEL) != 32 {
		t.Fatalf("agree(PEL's key, YST's share) = %d bytes, %v, want 32 bytes, nil", len(atPEL), err)
	}
	atYST, err := agree(yst, pel.PublicKey().Bytes())
	if err != nil || !bytes.Equal(atYST, atPEL) {
		t.Errorf("agree(YST's key, PEL's share) = %s, %v, want %s, the secret PEL arrived at", short(atYST), err, short(atPEL))
	}
}

// TestFixedKeys checks the keys this week's tests share: "ephemeral PEL" with Ystmark's share
// arrives at c50d64c8…cc7a.
func TestFixedKeys(t *testing.T) {
	const want = "c50d64c87ddd744d6c24842274ac98d74c33be2bdadd89463eae753cae91cc7a"
	pel, yst := ephemeral(t, "ephemeral PEL"), ephemeral(t, "ephemeral YST")
	got, err := agree(pel, yst.PublicKey().Bytes())
	if err != nil || hex.EncodeToString(got) != want {
		t.Errorf("agree(ephemeral PEL, YST's share) = %s, %v, want %s…, nil", short(got), err, want[:8])
	}
}

// TestThirdKey gives Ystmark's share to a third key, "ephemeral MID": it must arrive somewhere else.
func TestThirdKey(t *testing.T) {
	pel, yst, mid := ephemeral(t, "ephemeral PEL"), ephemeral(t, "ephemeral YST"), ephemeral(t, "ephemeral MID")
	atPEL, err := agree(pel, yst.PublicKey().Bytes())
	if err != nil || len(atPEL) != 32 {
		t.Fatalf("agree(ephemeral PEL, YST's share) = %d bytes, %v, want 32 bytes, nil", len(atPEL), err)
	}
	atMID, err := agree(mid, yst.PublicKey().Bytes())
	if err != nil || len(atMID) != 32 || bytes.Equal(atMID, atPEL) {
		t.Errorf("agree(ephemeral MID, YST's share) = %s, %v, want 32 bytes other than PEL's %s", short(atMID), err, short(atPEL))
	}
}
