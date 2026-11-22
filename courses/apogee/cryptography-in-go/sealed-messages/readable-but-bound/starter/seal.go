package main

import "fmt"

// Sealed is a command sealed for one satellite: only its header is readable.
type Sealed struct {
	Version    int    // 2
	SatID      string // "cirrine-1"
	Seq        uint64 // as Command.Seq: strictly increasing per satellite
	KeyID      string // "k-2028-07": the master version its key comes from
	Nonce      []byte // 12 bytes from crypto/rand, new for every seal
	Ciphertext []byte // AES-256-GCM of "OP|ARG...", its 16-byte tag included
}

// AD returns s's header as the associated data its seal binds: SEAL|2|cirrine-1|4471|k-2028-07.
func (s Sealed) AD() []byte {
	return fmt.Appendf(nil, "SEAL|%d|%s|%d|%s", s.Version, s.SatID, s.Seq, s.KeyID)
}
