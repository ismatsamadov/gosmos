package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
