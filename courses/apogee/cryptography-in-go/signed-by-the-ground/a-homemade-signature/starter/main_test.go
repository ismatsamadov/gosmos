package main

import (
	"crypto/sha256"
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"testing"
)

// The tests' keys are readable on purpose; flight keys never appear in code.
var (
	benchKey  = []byte("lab bench key, never flown")       // the lab bench's since 2027
	pipit3Key = []byte("pipit-3 test key, 32 bytes long!") // Mirela knows only its length
)

// extend plays Mirela. From the tag of msg and the key's length alone, it makes the message
// msg, SHA-256's padding, then extra, and the tag that sha256(key + message) gives it.
func extend(t *testing.T, msgTag []byte, msg string, keyLen int, extra string) (string, []byte) {
	t.Helper()
	n := keyLen + len(msg)
	pad := []byte{0x80}
	for (n+len(pad))%64 != 56 {
		pad = append(pad, 0)
	}
	pad = binary.BigEndian.AppendUint64(pad, uint64(n)*8)
	state := append([]byte("sha\x03"), msgTag...)
	state = append(state, make([]byte, 64)...)
	state = binary.BigEndian.AppendUint64(state, uint64(n+len(pad)))
	h := sha256.New()
	if err := h.(encoding.BinaryUnmarshaler).UnmarshalBinary(state); err != nil {
		t.Fatalf("rebuilding SHA-256 from a %d-byte tag: %v", len(msgTag), err)
	}
	h.Write([]byte(extra))
	return msg + string(pad) + extra, h.Sum(nil)
}

func TestLengthExtension(t *testing.T) {
	const msg = "pipit-3|212|SET_MODE|SAFE"
	msgTag := tag(pipit3Key, msg)
	if len(msgTag) != sha256.Size {
		t.Fatalf("len(tag(key, %q)) = %d, want %d", msg, len(msgTag), sha256.Size)
	}
	forged, forgedTag := extend(t, msgTag, msg, len(pipit3Key), "|REBOOT")
	if valid(pipit3Key, forged, forgedTag) {
		t.Errorf("valid(key, %q + SHA-256's padding + %q, Mirela's tag) = true, want false", msg, "|REBOOT")
	}
	if !valid(pipit3Key, msg, msgTag) {
		t.Errorf("valid(key, %q, tag(key, %q)) = false, want true", msg, msg)
	}
}

func TestTag(t *testing.T) {
	const want = "512545bcb43a5bb6ccb9f9c8511f4cc8c5fedeb38c226376af5b5aaae60e3b4d"
	if got := hex.EncodeToString(tag(benchKey, "pipit-2|118|SET_MODE|SAFE")); got != want {
		t.Errorf("tag(benchKey, %q) = %q, want %q, the lab bench's signature since 2027", "pipit-2|118|SET_MODE|SAFE", short(got), short(want))
	}
}

// short returns s, or its first 8 characters and an ellipsis when it is longer than 50.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}
