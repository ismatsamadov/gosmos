package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"testing"
)

// cirrineSealHex is Cirrine-1's seal key from the test master k-2028-07 (lesson 5.4).
const cirrineSealHex = "427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171"

const (
	known  = "SET_MODE|SAFE"          // a command Mirela knows went up
	secret = "DUMP_FILE|occ-0825.bin" // the command she wants to read
)

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

func TestNewNonceEachSeal(t *testing.T) {
	aead := testAEAD(t)
	n1, s1 := sealCmd(aead, []byte(known))
	n2, s2 := sealCmd(aead, []byte(known))
	if len(n1) != 12 {
		t.Fatalf("sealCmd(%q) nonce = %d bytes, want 12", known, len(n1))
	}
	if bytes.Equal(n1, n2) {
		t.Errorf("sealCmd(%q) twice = nonce %x both times, want a new nonce each time", known, n1)
	}
	if bytes.Equal(s1, s2) {
		t.Errorf("sealCmd(%q) twice = the same sealed bytes, want different ones", known)
	}
}

// TestMirelaReads plays Mirela. She records two sealed commands, knows the first went up as
// SET_MODE|SAFE, and XORs: sealed 1, sealed 2 and the known command.
func TestMirelaReads(t *testing.T) {
	aead := testAEAD(t)
	_, s1 := sealCmd(aead, []byte(known))
	_, s2 := sealCmd(aead, []byte(secret))
	if len(s1) < len(known) || len(s2) < len(known) {
		t.Fatalf("sealCmd = %d and %d sealed bytes, want at least %d each", len(s1), len(s2), len(known))
	}
	read := make([]byte, len(known))
	for i := range read {
		read[i] = s1[i] ^ s2[i] ^ known[i]
	}
	if string(read) == secret[:len(known)] {
		t.Errorf("Mirela's XOR of two sealed commands = %q, want nothing she can read", read)
	}
}

func TestEachOpens(t *testing.T) {
	aead := testAEAD(t)
	for _, cmd := range []string{known, secret} {
		nonce, sealed := sealCmd(aead, []byte(cmd))
		got, err := aead.Open(nil, nonce, sealed, nil)
		if err != nil || string(got) != cmd {
			t.Errorf("Open with sealCmd(%q)'s nonce = %q, %v, want %q, nil", cmd, got, err, cmd)
		}
	}
}
