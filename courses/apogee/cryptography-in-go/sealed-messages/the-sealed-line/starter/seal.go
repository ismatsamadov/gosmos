package main

import (
	"encoding/hex"
	"errors"
	"fmt"
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

// The ways a SEAL line is refused when it is read, each wrapped with what was wrong.
var (
	ErrMalformed = errors.New("malformed command line") // "malformed command line: want 7 fields"
	ErrVersion   = errors.New("unsupported version")    // "unsupported version: 3"
)

// AD returns s's header as the associated data its seal binds: SEAL|2|cirrine-1|4471|k-2028-07.
func (s Sealed) AD() []byte {
	return fmt.Appendf(nil, "SEAL|%d|%s|%d|%s", s.Version, s.SatID, s.Seq, s.KeyID)
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
