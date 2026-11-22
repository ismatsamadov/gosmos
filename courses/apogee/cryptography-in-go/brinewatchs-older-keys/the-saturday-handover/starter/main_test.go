package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"strings"
	"testing"
)

// Brinewatch's two new secrets (test values).
const (
	newAPIKey  = "brinewatch.20281017a1b2c3d4e5f60718293a4b5c"
	newWebhook = "brinewatch webhook secret 2028-10"
)

// line returns the request line naming pub by its whole fingerprint, worked out by the test.
func line(t *testing.T, pub *rsa.PublicKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(der)
	return "brinewatch handover sha256:" + hex.EncodeToString(sum[:])
}

// request returns what Brinewatch's security office sends: text, signed by signer as its own tool
// signs, with keyPEM beside it.
func request(t *testing.T, text string, signer *rsa.PrivateKey, keyPEM []byte) Request {
	t.Helper()
	digest := sha256.Sum256([]byte(text))
	sig, err := rsa.SignPSS(rand.Reader, signer, crypto.SHA256, digest[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	if err != nil {
		t.Fatal(err)
	}
	return Request{Text: []byte(text), Sig: sig, Key: keyPEM}
}

// handover runs Saturday's real handover and returns its two lines.
func handover(t *testing.T) []string {
	t.Helper()
	recv, sign := rsaKey(t, receivePEM), rsaKey(t, signPEM)
	r := request(t, line(t, &recv.PublicKey), sign, publicPEM(t, &recv.PublicKey))
	lines, err := handOver(&sign.PublicKey, r, []byte(newAPIKey), []byte(newWebhook))
	if err != nil || len(lines) != 2 {
		t.Fatalf("handOver(Brinewatch's request) = %d lines, %v, want 2 lines, nil", len(lines), err)
	}
	return lines
}

// lockedIn returns the bytes that line i carries after its name.
func lockedIn(t *testing.T, lines []string, i int, name string) []byte {
	t.Helper()
	got, hexText, _ := strings.Cut(lines[i], " ")
	locked, err := hex.DecodeString(hexText)
	if got != name || err != nil || len(locked) != 384 {
		shown := lines[i]
		if len(shown) > 50 {
			shown = shown[:24] + "…"
		}
		t.Fatalf("handOver line %d = %q, want %q, a space and 384 bytes in hex", i+1, shown, name)
	}
	return locked
}

// TestHandover opens both lines as Brinewatch will, with its receiving key and each line's label.
func TestHandover(t *testing.T) {
	recv := rsaKey(t, receivePEM)
	lines := handover(t)
	for i, want := range []struct{ name, label, secret string }{
		{"api-key", "apogee api key", newAPIKey},
		{"webhook-secret", "apogee webhook secret", newWebhook},
	} {
		locked := lockedIn(t, lines, i, want.name)
		got, err := rsa.DecryptOAEP(sha256.New(), nil, recv, locked, []byte(want.label))
		if err != nil || string(got) != want.secret {
			t.Errorf("Brinewatch opening line %d under %q = %q, %v, want %q, nil", i+1, want.label, got, err, want.secret)
		}
	}
}

// TestLabelsBind swaps what the two lines carry, as someone on the way might: neither opens as
// the other secret.
func TestLabelsBind(t *testing.T) {
	recv := rsaKey(t, receivePEM)
	lines := handover(t)
	apiKey, webhook := lockedIn(t, lines, 0, "api-key"), lockedIn(t, lines, 1, "webhook-secret")
	if _, err := rsa.DecryptOAEP(sha256.New(), nil, recv, webhook, []byte("apogee api key")); err == nil {
		t.Errorf("Brinewatch opening the webhook-secret line under %q = nil, want an error", "apogee api key")
	}
	if _, err := rsa.DecryptOAEP(sha256.New(), nil, recv, apiKey, []byte("apogee webhook secret")); err == nil {
		t.Errorf("Brinewatch opening the api-key line under %q = nil, want an error", "apogee webhook secret")
	}
}

// TestRefused gives handOver requests that must get no secret at all.
func TestRefused(t *testing.T) {
	recv, sign := rsaKey(t, receivePEM), rsaKey(t, signPEM)
	impostor, old := rsaKey(t, impostorPEM), rsaKey(t, oldPEM)
	recvPub, impostorPub := publicPEM(t, &recv.PublicKey), publicPEM(t, &impostor.PublicKey)
	changed := request(t, line(t, &recv.PublicKey), sign, impostorPub)
	changed.Text = []byte(line(t, &impostor.PublicKey))
	tests := []struct {
		name string
		r    Request
	}{
		{"signed by the impostor", request(t, line(t, &recv.PublicKey), impostor, recvPub)},
		{"signed with the receiving key", request(t, line(t, &recv.PublicKey), recv, recvPub)},
		{"changed on the way to name the impostor's key", changed},
		{"with the impostor's key swapped in beside it", request(t, line(t, &recv.PublicKey), sign, impostorPub)},
		{"naming the key by six hex digits", request(t, line(t, &recv.PublicKey)[:33], sign, recvPub)},
		{"for Brinewatch's 2021 key", request(t, line(t, &old.PublicKey), sign, publicPEM(t, &old.PublicKey))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, err := handOver(&sign.PublicKey, tt.r, []byte(newAPIKey), []byte(newWebhook))
			if err == nil || len(lines) != 0 {
				t.Errorf("handOver(a request %s) = %d lines, %v, want no lines and an error", tt.name, len(lines), err)
			}
		})
	}
}
