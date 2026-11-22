package main

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

// The HELLO lines of the test keys: "identity PEL" signing the share of "ephemeral PEL", and the
// same for YST.
const (
	pelHello = "HELLO PEL a10c608ed294ad1d437b121ce63f146ee081f200b8b71126bc8f5642e69aa93e " +
		"6da1b3b87f974f60b3e11e8420d9b8085d7d554e06823cde432819d6ec74abf1" +
		"53200bbeedf7d649e40fa42cb9cd344dba8712b25e87da0a25dc04628336390a"
	ystHello = "HELLO YST b6f8a5a974e292828807a52f32334ac3925ea872ca82de318f53e971fafc5a66 " +
		"ebe7d1730794166876e9d88347b47cc8f643b233f882acd87288e37725945ab2" +
		"41ab348be4a78717f9d955c3fc6ae79bc1f20e162b0df76c4b3a79234ab96203"
)

// seed returns the SHA-256 of label: every test key here is made from one.
func seed(label string) []byte {
	sum := sha256.Sum256([]byte(label))
	return sum[:]
}

// identity returns the test identity of the station code, made from the seed "identity <code>".
func identity(code string) Identity {
	return Identity{Code: code, Key: ed25519.NewKeyFromSeed(seed("identity " + code))}
}

// knows returns what a station holds in advance: the public identity key of each code.
func knows(codes ...string) map[string]ed25519.PublicKey {
	known := make(map[string]ed25519.PublicKey)
	for _, code := range codes {
		known[code] = identity(code).Key.Public().(ed25519.PublicKey)
	}
	return known
}

// ephemeral returns the test X25519 key of code, made from the seed "ephemeral <code>".
func ephemeral(t *testing.T, code string) *ecdh.PrivateKey {
	t.Helper()
	key, err := ecdh.X25519().NewPrivateKey(seed("ephemeral " + code))
	if err != nil {
		t.Fatalf("ecdh.X25519().NewPrivateKey(SHA-256 of %q) = %v", "ephemeral "+code, err)
	}
	return key
}

// brief returns line with each field longer than 50 characters cut to 8 and an ellipsis.
func brief(line string) string {
	fields := strings.Split(line, " ")
	for i, f := range fields {
		if len(f) > 50 {
			fields[i] = f[:8] + "…"
		}
	}
	return strings.Join(fields, " ")
}

// TestHello checks both stations' HELLO lines byte for byte: Ed25519 signs the same bytes the same way.
func TestHello(t *testing.T) {
	for _, tt := range []struct{ code, want string }{{"PEL", pelHello}, {"YST", ystHello}} {
		if got := hello(identity(tt.code), ephemeral(t, tt.code)); got != tt.want {
			t.Errorf("hello(%s) = %q, want %q", tt.code, brief(got), brief(tt.want))
		}
	}
}

// TestCheckHello: each end reads the other's HELLO, holding the other's identity key in advance.
func TestCheckHello(t *testing.T) {
	tests := []struct {
		line, code string
		known      map[string]ed25519.PublicKey
	}{
		{pelHello, "PEL", knows("PEL")}, // at Ystmark
		{ystHello, "YST", knows("YST")}, // at Port Ellery
	}
	for _, tt := range tests {
		code, share, err := checkHello(tt.line, tt.known)
		if err != nil || code != tt.code || share == nil || !share.Equal(ephemeral(t, tt.code).PublicKey()) {
			t.Errorf("checkHello(%q) = %q, %v, want %q, its share, nil", brief(tt.line), code, err, tt.code)
		}
	}
}

// TestMiddle plays the company that runs the backhaul. It answers PEL with its own share under
// YST's name, and YST with the same share under PEL's, signed with the only identity key it has,
// "identity MID". An end that accepts shares a secret with the middle.
func TestMiddle(t *testing.T) {
	mid := ephemeral(t, "MID")
	swaps := []struct {
		end, claims string
		known       map[string]ed25519.PublicKey
	}{
		{"PEL", "YST", knows("YST")},
		{"YST", "PEL", knows("PEL")},
	}
	for _, s := range swaps {
		line := hello(Identity{Code: s.claims, Key: identity("MID").Key}, mid)
		code, share, err := checkHello(line, s.known)
		if errors.Is(err, ErrBadHello) {
			continue
		}
		stolen := "nothing yet"
		if share != nil {
			atEnd, _ := ephemeral(t, s.end).ECDH(share)
			atMid, _ := mid.ECDH(ephemeral(t, s.end).PublicKey())
			if len(atEnd) == 32 && bytes.Equal(atEnd, atMid) {
				stolen = s.end + "'s secret " + hex.EncodeToString(atEnd[:4]) + "…"
			}
		}
		t.Errorf("checkHello(the middle's HELLO to %s, as %s) = %q, %v, want ErrBadHello: the middle holds %s", s.end, s.claims, code, err, stolen)
	}
}

// TestRefusals sends lines each end must refuse, every one with an error wrapping ErrBadHello.
func TestRefusals(t *testing.T) {
	sig := ystHello[strings.LastIndex(ystHello, " ")+1:]
	share := strings.Fields(ystHello)[2]
	tests := []struct{ name, line string }{
		{"YST's HELLO with one share digit changed", strings.Replace(ystHello, "YST b6f8", "YST a6f8", 1)},
		{"the middle's own HELLO, as MID", hello(identity("MID"), ephemeral(t, "MID"))},
		{"YST's HELLO without its signature", "HELLO YST " + share},
		{"PEL's own HELLO, sent back to PEL", pelHello},
		{"YST's HELLO with a fifth field", ystHello + " again"},
		{"YST's HELLO with its share in capitals", "HELLO YST " + strings.ToUpper(share) + " " + sig},
		{"an empty line", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Port Ellery's end: it knows YST alone.
			if code, _, err := checkHello(tt.line, knows("YST")); !errors.Is(err, ErrBadHello) {
				t.Errorf("checkHello(%q) = %q, %v, want ErrBadHello", brief(tt.line), code, err)
			}
		})
	}
}

// TestMessages checks the two refusals the log will show most often.
func TestMessages(t *testing.T) {
	mid := ephemeral(t, "MID")
	tests := []struct{ name, line, want string }{
		{"the middle as MID", hello(identity("MID"), mid), `bad hello: unknown station "MID"`},
		{"the middle as YST", hello(Identity{Code: "YST", Key: identity("MID").Key}, mid), "bad hello: signature does not verify"},
	}
	for _, tt := range tests {
		if _, _, err := checkHello(tt.line, knows("YST")); err == nil || err.Error() != tt.want {
			t.Errorf("checkHello(%s) = %v, want %q", tt.name, err, tt.want)
		}
	}
}
