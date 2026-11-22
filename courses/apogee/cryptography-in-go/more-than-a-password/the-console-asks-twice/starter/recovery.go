package main

import (
	"crypto/rand"
	"crypto/subtle"
	"slices"
)

// sheetSize is how many recovery codes one printed sheet holds.
const sheetSize = 8

// Recovery keeps each engineer's unused recovery codes, as digests only.
type Recovery struct {
	unused map[string][]string // user → the hex SHA-256 of each code not used yet
}

// Issue makes user a new sheet of recovery codes, replacing any older one, and returns the codes to print once.
func (r *Recovery) Issue(user string) []string {
	codes := make([]string, sheetSize)
	digests := make([]string, sheetSize)
	for i := range codes {
		codes[i] = rand.Text()
		digests[i] = tokenDigest(codes[i])
	}
	r.unused[user] = digests
	return codes
}

// Use reports whether code is one of user's unused recovery codes, and spends it if so.
func (r *Recovery) Use(user, code string) bool {
	typed := []byte(tokenDigest(code))
	for i, kept := range r.unused[user] {
		if subtle.ConstantTimeCompare(typed, []byte(kept)) == 1 {
			r.unused[user] = slices.Delete(r.unused[user], i, i+1)
			return true
		}
	}
	return false
}
