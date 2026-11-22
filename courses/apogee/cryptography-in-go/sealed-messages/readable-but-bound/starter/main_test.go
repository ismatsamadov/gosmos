package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"testing"
)

// cirrineSealHex is Cirrine-1's seal key from the test master k-2028-07 (lesson 5.4).
const cirrineSealHex = "427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171"

// testAEAD returns AES-256-GCM under Cirrine-1's seal key.
func testAEAD(t *testing.T) cipher.AEAD {
	t.Helper()
	key, err := hex.DecodeString(cirrineSealHex)
	if err != nil {
		t.Fatal(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatal(err)
	}
	return aead
}

// header is the readable part of Cirrine-1's first sealed command.
func header() Sealed {
	return Sealed{Version: 2, SatID: "cirrine-1", Seq: 4471, KeyID: "k-2028-07"}
}

func TestRoundTrip(t *testing.T) {
	aead := testAEAD(t)
	s := sealWithHeader(aead, header(), []byte("PING"))
	got, err := openWithHeader(aead, s)
	if err != nil || string(got) != "PING" {
		t.Errorf("openWithHeader(sealWithHeader(PING)) = %q, %v, want %q, nil", got, err, "PING")
	}
}

// TestHeaderBound opens the sealed bytes as Cirrine-1 does: with s.AD() as the associated data.
func TestHeaderBound(t *testing.T) {
	aead := testAEAD(t)
	s := sealWithHeader(aead, header(), []byte("PING"))
	if _, err := aead.Open(nil, s.Nonce, s.Ciphertext, s.AD()); err != nil {
		t.Errorf("Open with s.AD() after sealWithHeader = %v, want nil: seal with s.AD() as the associated data", err)
	}
}

// TestMoved plays Mirela: she keeps the nonce and the sealed bytes and changes the header.
func TestMoved(t *testing.T) {
	aead := testAEAD(t)
	s := sealWithHeader(aead, header(), []byte("PING"))
	tests := []struct {
		name string
		move func(m *Sealed)
	}{
		{"number 4472", func(m *Sealed) { m.Seq = 4472 }},
		{"satellite pipit-2", func(m *Sealed) { m.SatID = "pipit-2" }},
		{"key k-2028-10", func(m *Sealed) { m.KeyID = "k-2028-10" }},
		{"version 3", func(m *Sealed) { m.Version = 3 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moved := s
			tt.move(&moved)
			if got, err := openWithHeader(aead, moved); err == nil {
				t.Errorf("openWithHeader(the PING moved to %s) = %q, nil, want an error", tt.name, got)
			}
		})
	}
}

func TestHeaderReadable(t *testing.T) {
	aead := testAEAD(t)
	s := sealWithHeader(aead, header(), []byte("PING"))
	h := header()
	if s.Version != h.Version || s.SatID != h.SatID || s.Seq != h.Seq || s.KeyID != h.KeyID {
		t.Errorf("sealWithHeader's header = %d %s %d %s, want %d %s %d %s",
			s.Version, s.SatID, s.Seq, s.KeyID, h.Version, h.SatID, h.Seq, h.KeyID)
	}
	if len(s.Nonce) != 12 || len(s.Ciphertext) != 4+16 {
		t.Errorf("sealWithHeader(PING) = a %d-byte nonce and %d sealed bytes, want 12 and 20", len(s.Nonce), len(s.Ciphertext))
	}
}
