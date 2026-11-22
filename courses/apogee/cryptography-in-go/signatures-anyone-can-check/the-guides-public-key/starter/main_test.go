package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

// testKey returns the test receipt key r-2028-08, made from a readable seed.
func testKey() ed25519.PrivateKey {
	seed := sha256.Sum256([]byte("receipts r-2028-08"))
	return ed25519.NewKeyFromSeed(seed[:])
}

func TestRoundTrip(t *testing.T) {
	pub := testKey().Public().(ed25519.PublicKey)
	text, err := publishKey(pub)
	if err != nil {
		t.Fatalf("publishKey(r-2028-08) = %v, want nil", err)
	}
	if block, _ := pem.Decode(text); block == nil || block.Type != "PUBLIC KEY" {
		t.Fatalf("publishKey(r-2028-08) = %q, want a PUBLIC KEY block", text)
	}
	got, err := readKey(text)
	if err != nil || !bytes.Equal(got, pub) {
		t.Errorf("readKey(publishKey(r-2028-08)) = %x, %v, want %x, nil", got, err, pub)
	}
}

func TestRefusals(t *testing.T) {
	good, err := publishKey(testKey().Public().(ed25519.PublicKey))
	if err != nil {
		t.Fatalf("publishKey(r-2028-08) = %v, want nil", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(testKey())
	if err != nil {
		t.Fatal(err)
	}
	private := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	tests := []struct {
		name string
		text []byte
	}{
		{"a PRIVATE KEY block", private},
		{"text after the block", append(bytes.Clone(good), "r-2028-10 coming soon\n"...)},
		{"two PUBLIC KEY blocks", append(bytes.Clone(good), good...)},
		{"nothing at all", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := readKey(tt.text); err == nil {
				t.Errorf("readKey(%s) = %x, nil, want an error", tt.name, got)
			}
		})
	}
}

// TestNotEd25519 gives readKey a PUBLIC KEY block that holds another kind of key, ECDSA P-256.
func TestNotEd25519(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	text := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
	if _, err := readKey(text); err == nil || err.Error() != "not an Ed25519 key" {
		t.Errorf("readKey(a P-256 PUBLIC KEY block) = %v, want the error %q", err, "not an Ed25519 key")
	}
}
