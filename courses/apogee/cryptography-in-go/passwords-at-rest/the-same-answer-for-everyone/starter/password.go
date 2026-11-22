package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Record is a parsed password record.
type Record struct {
	Iterations int    // 600000
	Salt       []byte // 16 bytes
	Key        []byte // 32 bytes
}

// parseRecord reads a stored record, refusing all but pbkdf2-sha256 with a count from 1,000 to 10,000,000, a 16-byte salt and a 32-byte key.
func parseRecord(s string) (Record, error) {
	f := strings.Split(s, "$")
	if len(f) != 4 || f[0] != "pbkdf2-sha256" {
		return Record{}, errors.New("record is not pbkdf2-sha256$<count>$<salt>$<key>")
	}
	count, err := strconv.Atoi(f[1])
	if err != nil || count < 1000 || count > 10000000 {
		return Record{}, fmt.Errorf("count %q is not from 1000 to 10000000", f[1])
	}
	salt, err := hex.DecodeString(f[2])
	if err != nil || len(salt) != 16 {
		return Record{}, errors.New("salt is not 32 hex digits")
	}
	key, err := hex.DecodeString(f[3])
	if err != nil || len(key) != 32 {
		return Record{}, errors.New("key is not 64 hex digits")
	}
	return Record{Iterations: count, Salt: salt, Key: key}, nil
}
