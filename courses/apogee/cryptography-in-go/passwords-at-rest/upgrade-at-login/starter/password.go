package main

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// iterations is how many rounds of work derive spends on every password.
const iterations = 600000

// Record is a parsed password record.
type Record struct {
	Iterations int    // 600000
	Salt       []byte // 16 bytes
	Key        []byte // 32 bytes
}

// derive returns the 32-byte key the console stores for pw with the given salt.
func derive(pw string, salt []byte) ([]byte, error) {
	return pbkdf2.Key(sha256.New, pw, salt, iterations, 32)
}

// formatRecord writes r as the ops console stores it.
func formatRecord(r Record) string {
	return fmt.Sprintf("pbkdf2-sha256$%d$%x$%x", r.Iterations, r.Salt, r.Key)
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

// checkPassword reports whether pw is the password r was made from.
func checkPassword(r Record, pw string) (bool, error) {
	key, err := pbkdf2.Key(sha256.New, pw, r.Salt, r.Iterations, len(r.Key))
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare(key, r.Key) == 1, nil
}

// newRecord makes a record of pw with a fresh salt and today's count.
func newRecord(pw string) (string, error) {
	salt := newSalt()
	key, err := derive(pw, salt)
	if err != nil {
		return "", err
	}
	return formatRecord(Record{Iterations: iterations, Salt: salt, Key: key}), nil
}
