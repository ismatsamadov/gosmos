package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"strings"
	"testing"
)

// request is the line Brinewatch's security office signs for the handover, naming its new
// receiving key by the key's whole fingerprint.
const request = "brinewatch handover sha256:d2499efc1976897aa14ea5d3948c4524bb826ad42b6ac1b43c5772f6b2ab8c1e"

// brinewatchSigns signs text as Brinewatch's own tool does: RSA-PSS over SHA-256, its random
// part as long as the hash.
func brinewatchSigns(t *testing.T, priv *rsa.PrivateKey, text string) []byte {
	t.Helper()
	digest := sha256.Sum256([]byte(text))
	sig, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA256, digest[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	if err != nil {
		t.Fatal(err)
	}
	return sig
}

// TestSignRequest checks signRequest's signature the way Brinewatch's tools check one.
func TestSignRequest(t *testing.T) {
	sign := rsaKey(t, signPEM)
	sig, err := signRequest(sign, []byte(request))
	if err != nil || len(sig) != 384 {
		t.Fatalf("signRequest(the request) = %d bytes, %v, want 384 bytes, nil", len(sig), err)
	}
	digest := sha256.Sum256([]byte(request))
	if err := rsa.VerifyPSS(&sign.PublicKey, crypto.SHA256, digest[:], sig, nil); err != nil {
		t.Errorf("rsa.VerifyPSS(signRequest(the request)) = %v, want nil", err)
	}
}

// TestNeverAlike signs the request twice: the random bytes PSS mixes in make the two signatures
// differ, and checkRequest accepts both.
func TestNeverAlike(t *testing.T) {
	sign := rsaKey(t, signPEM)
	a, errA := signRequest(sign, []byte(request))
	b, errB := signRequest(sign, []byte(request))
	if errA != nil || errB != nil || len(a) == 0 {
		t.Fatalf("signRequest(the request) twice = %d bytes, %v and %v, want 384 bytes, nil", len(a), errA, errB)
	}
	if bytes.Equal(a, b) {
		t.Errorf("signRequest(the request) twice = the same signature, want two different ones")
	}
	for i, sig := range [][]byte{a, b} {
		if err := checkRequest(&sign.PublicKey, []byte(request), sig); err != nil {
			t.Errorf("checkRequest(the request, signature %d of 2) = %v, want nil", i+1, err)
		}
	}
}

// TestBrinewatchsSignature checks a signature made by Brinewatch's own tool.
func TestBrinewatchsSignature(t *testing.T) {
	sign := rsaKey(t, signPEM)
	sig := brinewatchSigns(t, sign, request)
	if err := checkRequest(&sign.PublicKey, []byte(request), sig); err != nil {
		t.Errorf("checkRequest(the request, Brinewatch's signature) = %v, want nil", err)
	}
}

// TestRefused gives checkRequest signatures that must not count.
func TestRefused(t *testing.T) {
	sign, recv := rsaKey(t, signPEM), rsaKey(t, receivePEM)
	good := brinewatchSigns(t, sign, request)
	digest := sha256.Sum256([]byte(request))
	older, err := rsa.SignPKCS1v15(rand.Reader, sign, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	changed := strings.Replace(request, "8c1e", "8c1f", 1)
	tests := []struct {
		name string
		req  string
		sig  []byte
	}{
		{"a request changed after it was signed", changed, good},
		{"a signature by the receiving key", request, brinewatchSigns(t, recv, request)},
		{"an older kind of RSA signature, PKCS #1 v1.5", request, older},
		{"a signature cut to 383 bytes", request, good[:383]},
		{"no signature at all", request, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkRequest(&sign.PublicKey, []byte(tt.req), tt.sig); err == nil {
				t.Errorf("checkRequest(%s) = nil, want an error", tt.name)
			}
		})
	}
}
