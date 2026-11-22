package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// The test vector: the shares of "ephemeral PEL", which dials, and "ephemeral YST", which
// answers; the secret they agree on; and the session key made from it.
var (
	pelShare = fromHex("a10c608ed294ad1d437b121ce63f146ee081f200b8b71126bc8f5642e69aa93e")
	ystShare = fromHex("b6f8a5a974e292828807a52f32334ac3925ea872ca82de318f53e971fafc5a66")
	agreed   = fromHex("c50d64c87ddd744d6c24842274ac98d74c33be2bdadd89463eae753cae91cc7a")
)

const want = "24e57dc80a28df831da79d5b65636ba94ed13a4274908279576a15428a06958e"

// fromHex decodes the test vector's hex, which is always well formed.
func fromHex(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

// short returns b in hex, cut to its first 8 digits and an ellipsis when longer than 50.
func short(b []byte) string {
	if s := hex.EncodeToString(b); len(s) <= 50 {
		return s
	}
	return hex.EncodeToString(b[:4]) + "…"
}

func TestVector(t *testing.T) {
	got, err := sessionKey(agreed, pelShare, ystShare)
	if err != nil || hex.EncodeToString(got) != want {
		t.Errorf("sessionKey(c50d64c8…, PEL's share, YST's share) = %s, %v, want %s…, nil", short(got), err, want[:8])
	}
}

// TestNotTheRawSecret: the session key is never the X25519 secret itself.
func TestNotTheRawSecret(t *testing.T) {
	got, err := sessionKey(agreed, pelShare, ystShare)
	if err != nil || bytes.Equal(got, agreed) {
		t.Errorf("sessionKey(c50d64c8…, PEL's share, YST's share) = %s, %v, want a key other than the secret c50d64c8…", short(got), err)
	}
}

// TestOrder swaps the two shares, as if each end thought it had dialed: the keys must differ.
func TestOrder(t *testing.T) {
	right, err1 := sessionKey(agreed, pelShare, ystShare)
	swapped, err2 := sessionKey(agreed, ystShare, pelShare)
	if err1 != nil || err2 != nil || bytes.Equal(right, swapped) {
		t.Errorf("sessionKey(…, YST's share, PEL's share) = %s, %v, want a key other than %s", short(swapped), err2, short(right))
	}
}

// TestSecretSize sends secrets that aren't 32 bytes: each is an error, never a key.
func TestSecretSize(t *testing.T) {
	for _, n := range []int{0, 31, 33} {
		if got, err := sessionKey(make([]byte, n), pelShare, ystShare); err == nil {
			t.Errorf("sessionKey(a %d-byte secret, …) = %d bytes, nil, want an error", n, len(got))
		}
	}
}
