package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"
)

// Brinewatch's two new secrets (test values) and the labels they are locked under.
const (
	newAPIKey    = "brinewatch.20281017a1b2c3d4e5f60718293a4b5c"
	newWebhook   = "brinewatch webhook secret 2028-10"
	apiKeyLabel  = "apogee api key"
	webhookLabel = "apogee webhook secret"
)

// newSecrets lists each secret with its label.
func newSecrets() []struct{ name, secret, label string } {
	return []struct{ name, secret, label string }{
		{"the new API key", newAPIKey, apiKeyLabel},
		{"the new webhook secret", newWebhook, webhookLabel},
	}
}

// TestBrinewatchOpens locks each secret with encryptTo and opens it as Brinewatch's own tools
// will: rsa.DecryptOAEP with SHA-256, the receiving key's private half and the same label.
func TestBrinewatchOpens(t *testing.T) {
	recv := rsaKey(t, receivePEM)
	for _, s := range newSecrets() {
		locked, err := encryptTo(&recv.PublicKey, []byte(s.secret), s.label)
		if err != nil || len(locked) != 384 {
			t.Errorf("encryptTo(%s) = %d bytes, %v, want 384 bytes, nil", s.name, len(locked), err)
			continue
		}
		if bytes.Contains(locked, []byte(s.secret)) {
			t.Errorf("encryptTo(%s) = bytes holding the secret itself, want it locked", s.name)
		}
		got, err := rsa.DecryptOAEP(sha256.New(), nil, recv, locked, []byte(s.label))
		if err != nil || string(got) != s.secret {
			t.Errorf("Brinewatch opening encryptTo(%s) = %q, %v, want %q, nil", s.name, got, err, s.secret)
		}
	}
}

// TestDecryptWith opens secrets the test locked the way Brinewatch's tools lock them.
func TestDecryptWith(t *testing.T) {
	recv := rsaKey(t, receivePEM)
	for _, s := range newSecrets() {
		locked, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &recv.PublicKey, []byte(s.secret), []byte(s.label))
		if err != nil {
			t.Fatal(err)
		}
		got, err := decryptWith(recv, locked, s.label)
		if err != nil || string(got) != s.secret {
			t.Errorf("decryptWith(%s, locked by the test) = %q, %v, want %q, nil", s.name, got, err, s.secret)
		}
	}
}

// TestNeverTheSameTwice locks the API key twice: the random bytes OAEP mixes in make the two
// ciphertexts differ, and each still opens.
func TestNeverTheSameTwice(t *testing.T) {
	recv := rsaKey(t, receivePEM)
	a, errA := encryptTo(&recv.PublicKey, []byte(newAPIKey), apiKeyLabel)
	b, errB := encryptTo(&recv.PublicKey, []byte(newAPIKey), apiKeyLabel)
	if errA != nil || errB != nil || len(a) == 0 {
		t.Fatalf("encryptTo(the new API key) twice = %d bytes, %v and %v, want 384 bytes, nil", len(a), errA, errB)
	}
	if bytes.Equal(a, b) {
		t.Errorf("encryptTo(the new API key) twice = the same bytes, want two different ciphertexts")
	}
	for _, locked := range [][]byte{a, b} {
		if got, err := decryptWith(recv, locked, apiKeyLabel); err != nil || string(got) != newAPIKey {
			t.Errorf("decryptWith(one of the two) = %q, %v, want %q, nil", got, err, newAPIKey)
		}
	}
}

// TestLabelBinds opens the webhook secret's ciphertext under the API key's label: it must fail,
// so that one secret's ciphertext can never pass for the other's.
func TestLabelBinds(t *testing.T) {
	recv := rsaKey(t, receivePEM)
	locked, err := encryptTo(&recv.PublicKey, []byte(newWebhook), webhookLabel)
	if err != nil || len(locked) == 0 {
		t.Fatalf("encryptTo(the new webhook secret) = %d bytes, %v, want 384 bytes, nil", len(locked), err)
	}
	if got, err := decryptWith(recv, locked, apiKeyLabel); err == nil {
		t.Errorf("decryptWith(the webhook secret's ciphertext, %q) = %q, nil, want an error", apiKeyLabel, got)
	}
	if _, err := rsa.DecryptOAEP(sha256.New(), nil, recv, locked, []byte(apiKeyLabel)); err == nil {
		t.Errorf("Brinewatch opening the webhook secret's ciphertext under %q = nil, want an error: encryptTo must lock with its label", apiKeyLabel)
	}
}

// TestDecryptRefuses gives decryptWith the API key's ciphertext in ways that must not open.
func TestDecryptRefuses(t *testing.T) {
	recv, sign := rsaKey(t, receivePEM), rsaKey(t, signPEM)
	locked, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &recv.PublicKey, []byte(newAPIKey), []byte(apiKeyLabel))
	if err != nil {
		t.Fatal(err)
	}
	changed := bytes.Clone(locked)
	changed[100] ^= 0x01
	tests := []struct {
		name   string
		priv   *rsa.PrivateKey
		locked []byte
	}{
		{"with Brinewatch's signing key", sign, locked},
		{"with one byte changed", recv, changed},
		{"cut to 383 bytes", recv, locked[:383]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := decryptWith(tt.priv, tt.locked, apiKeyLabel); err == nil {
				t.Errorf("decryptWith(the API key's ciphertext, %s) = %q, nil, want an error", tt.name, got)
			}
		})
	}
}

// TestTooLong locks the most a 3,072-bit key can take, 318 bytes, and one byte more.
func TestTooLong(t *testing.T) {
	recv := rsaKey(t, receivePEM)
	if got, err := encryptTo(&recv.PublicKey, make([]byte, 318), apiKeyLabel); err != nil || len(got) != 384 {
		t.Errorf("encryptTo(318 bytes) = %d bytes, %v, want 384 bytes, nil: 318 is the most a 3,072-bit key locks", len(got), err)
	}
	if _, err := encryptTo(&recv.PublicKey, make([]byte, 319), apiKeyLabel); err == nil {
		t.Errorf("encryptTo(319 bytes) = nil, want an error: 318 is the most a 3,072-bit key locks")
	}
}
