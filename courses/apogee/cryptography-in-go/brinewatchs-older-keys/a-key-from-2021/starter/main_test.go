package main

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/pem"
	"fmt"
	"math/big"
	"testing"
)

// keyText says what readRSA returned, for messages.
func keyText(pub *rsa.PublicKey) string {
	if pub == nil {
		return "no key"
	}
	return fmt.Sprintf("a %d-bit key", pub.N.BitLen())
}

func TestRSAKeys(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"Brinewatch's receiving key, 3,072 bits", receivePEM},
		{"Brinewatch's signing key, 3,072 bits", signPEM},
		{"the impostor's key, 2,048 bits exactly", impostorPEM},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priv := rsaKey(t, tt.key)
			got, err := readRSA(publicPEM(t, &priv.PublicKey))
			if err != nil || got == nil || !got.Equal(&priv.PublicKey) {
				t.Errorf("readRSA(%s) = %s, %v, want that key and nil", tt.name, keyText(got), err)
			}
		})
	}
}

// TestTooShort reads the key Brinewatch has used since 2021, a valid RSA key of 1,024 bits, and
// a public key one bit short of 2,048, whose number is the impostor's halved: readRSA only
// measures it.
func TestTooShort(t *testing.T) {
	n := new(big.Int).Rsh(rsaKey(t, impostorPEM).N, 1)
	n.SetBit(n, 0, 1)
	tests := []struct {
		name string
		pub  *rsa.PublicKey
		want string
	}{
		{"Brinewatch's 2021 key", &rsaKey(t, oldPEM).PublicKey, "want at least 2048 bits, got 1024"},
		{"a key of 2,047 bits", &rsa.PublicKey{N: n, E: 65537}, "want at least 2048 bits, got 2047"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readRSA(publicPEM(t, tt.pub))
			if got != nil || err == nil || err.Error() != tt.want {
				t.Errorf("readRSA(%s) = %s, %v, want no key and the error %q", tt.name, keyText(got), err, tt.want)
			}
		})
	}
}

// TestOtherKinds gives readRSA PUBLIC KEY blocks that hold other kinds of key.
func TestOtherKinds(t *testing.T) {
	receipt := sha256.Sum256([]byte("receipts r-2028-08"))
	share := sha256.Sum256([]byte("ephemeral PEL"))
	x, err := ecdh.X25519().NewPrivateKey(share[:])
	if err != nil {
		t.Fatal(err)
	}
	p256, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		pub  any
		want string
	}{
		{"Apogee's receipt key r-2028-08", ed25519.NewKeyFromSeed(receipt[:]).Public(), "want an RSA key, got Ed25519"},
		{"Port Ellery's backhaul share", x.PublicKey(), "want an RSA key, got X25519"},
		{"an ECDSA P-256 key", &p256.PublicKey, "want an RSA key, got *ecdsa.PublicKey"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readRSA(publicPEM(t, tt.pub))
			if got != nil || err == nil || err.Error() != tt.want {
				t.Errorf("readRSA(%s) = %s, %v, want no key and the error %q", tt.name, keyText(got), err, tt.want)
			}
		})
	}
}

// TestNotOneKey gives readRSA input that isn't one whole PUBLIC KEY block.
func TestNotOneKey(t *testing.T) {
	good := publicPEM(t, &rsaKey(t, receivePEM).PublicKey)
	block, _ := pem.Decode(good)
	cut := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: block.Bytes[:200]})
	tests := []struct {
		name string
		text []byte
	}{
		{"a PRIVATE KEY block", []byte(receivePEM)},
		{"two PUBLIC KEY blocks", append(bytes.Clone(good), good...)},
		{"a key cut short", cut},
		{"nothing at all", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := readRSA(tt.text); got != nil || err == nil {
				t.Errorf("readRSA(%s) = %s, %v, want no key and an error", tt.name, keyText(got), err)
			}
		})
	}
}
