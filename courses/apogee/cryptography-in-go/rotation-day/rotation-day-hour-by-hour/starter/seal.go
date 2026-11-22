package main

import (
	"crypto/aes"
	"crypto/cipher"
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

// AD returns s's header as the associated data its seal binds: SEAL|2|cirrine-1|4471|k-2028-07.
func (s Sealed) AD() []byte {
	return fmt.Appendf(nil, "SEAL|%d|%s|%d|%s", s.Version, s.SatID, s.Seq, s.KeyID)
}

// formatSealed writes s as one SEAL line: SEAL 2 <sat> <seq> <key id> <nonce> <ciphertext>.
func formatSealed(s Sealed) string {
	return fmt.Sprintf("SEAL %d %s %d %s %x %x", s.Version, s.SatID, s.Seq, s.KeyID, s.Nonce, s.Ciphertext)
}

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

// newAEAD returns AES-256-GCM under a 32-byte seal key.
func newAEAD(key []byte) (cipher.AEAD, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("seal key has %d bytes, want 32", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

// openCommand opens s with a satellite's seal key and returns the command sealed in it.
func openCommand(key []byte, s Sealed) (Command, error) {
	aead, err := newAEAD(key)
	if err != nil {
		return Command{}, err
	}
	if len(s.Nonce) != aead.NonceSize() {
		return Command{}, ErrForged
	}
	plaintext, err := aead.Open(nil, s.Nonce, s.Ciphertext, s.AD())
	if err != nil {
		return Command{}, ErrForged
	}
	fields := strings.Split(string(plaintext), "|")
	return Command{SatID: s.SatID, Seq: s.Seq, Op: fields[0], Args: fields[1:]}, nil
}

// sealCommand seals c for its satellite under master version keyID and returns its SEAL line.
func sealCommand(ring *KeyRing, keyID string, c Command) (string, error) {
	if _, err := canonical(c); err != nil {
		return "", err
	}
	key, err := ring.SatKey(keyID, "seal", c.SatID)
	if err != nil {
		return "", err
	}
	aead, err := newAEAD(key)
	if err != nil {
		return "", err
	}
	s := Sealed{Version: 2, SatID: c.SatID, Seq: c.Seq, KeyID: keyID, Nonce: newNonce()}
	plaintext := strings.Join(append([]string{c.Op}, c.Args...), "|")
	s.Ciphertext = aead.Seal(nil, s.Nonce, []byte(plaintext), s.AD())
	return formatSealed(s), nil
}
