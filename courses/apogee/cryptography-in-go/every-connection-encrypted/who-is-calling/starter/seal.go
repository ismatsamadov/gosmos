package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Sealed is a command sealed for one satellite: only its header is readable.
type Sealed struct {
	Version    int    // 2
	SatID      string // "cirrine-1"
	Seq        uint64 // as Command.Seq: strictly increasing per satellite
	KeyID      string // "k-2028-07": the master version its key comes from
	Nonce      []byte // 12 bytes from crypto/rand, new for every seal
	Ciphertext []byte // AES-256-GCM of "OP|ARG...", its 16-byte tag included
}

// The ways a SEAL line is refused, each wrapped with what was wrong.
var (
	ErrMalformed = errors.New("malformed command line") // "malformed command line: want 7 fields"
	ErrVersion   = errors.New("unsupported version")    // "unsupported version: 3"
	ErrWrongSat  = errors.New("not for this satellite") // `not for this satellite: "pipit-2"`
	// ErrUnknownKey means the ring holds no version with the key ID asked for.
	ErrUnknownKey = errors.New("unknown key")
	ErrRevoked    = errors.New("key revoked") // `key revoked: "k-2028-10"`
	ErrRetired    = errors.New("key retired") // `key retired: "k-2028-07"`
	// ErrReplay refuses a command whose sequence number isn't above the last one accepted.
	ErrReplay = errors.New("replayed sequence number")
	// ErrForged refuses a command whose signature doesn't check out, and never says more.
	ErrForged = errors.New("command failed authentication")
)

// parseSealed reads one SEAL line, refusing anything but the exact format.
func parseSealed(line string) (Sealed, error) {
	f := strings.Split(line, " ")
	if len(f) != 7 {
		return Sealed{}, fmt.Errorf("%w: want 7 fields", ErrMalformed)
	}
	if f[0] == "SEAL" && f[1] != "2" {
		return Sealed{}, fmt.Errorf("%w: %s", ErrVersion, f[1])
	}
	seq, err := strconv.ParseUint(f[3], 10, 64)
	nonce, nonceOK := hexBytes(f[5], 12, 12)
	text, textOK := hexBytes(f[6], 17, 272)
	if f[0] != "SEAL" || !validID(f[2]) || err != nil || seq == 0 || !validID(f[4]) || !nonceOK || !textOK {
		return Sealed{}, fmt.Errorf("%w: want SEAL 2 <sat> <seq> <key id> <nonce> <ciphertext>", ErrMalformed)
	}
	return Sealed{Version: 2, SatID: f[2], Seq: seq, KeyID: f[4], Nonce: nonce, Ciphertext: text}, nil
}

// hexBytes decodes s if it is lowercase hex for least to most bytes.
func hexBytes(s string, least, most int) ([]byte, bool) {
	if len(s)%2 != 0 || len(s) < 2*least || len(s) > 2*most || strings.ToLower(s) != s {
		return nil, false
	}
	b, err := hex.DecodeString(s)
	return b, err == nil
}

// validID reports whether s is a satellite or key ID: lowercase letters, digits and hyphens.
func validID(s string) bool {
	return s != "" && strings.Trim(s, "abcdefghijklmnopqrstuvwxyz0123456789-") == ""
}
