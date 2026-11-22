package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"
)

// Console is the ops console's login.
type Console struct {
	records  map[string]string                                      // user → record text
	kdf      func(pw string, salt []byte, iter int) ([]byte, error) // PBKDF2-HMAC-SHA256; tests count calls
	save     func(user, record string) error                        // stores an upgraded record
	sessions map[string]string                                      // hex SHA-256 of a token → user
}

// errInvalidLogin is the console's one answer to a failed login.
var errInvalidLogin = errors.New("invalid user or password")

// dummyRecord is a valid 600,000-iteration record whose password nobody knows. An unknown
// user is checked against it, so a failed login costs the same work whoever tries.
const dummyRecord = "pbkdf2-sha256$600000$64756d6d79207265636f72642c203136$6e6f2070617373776f726420646572697665732074686973206b65792c203332"

// Check returns nil when pw is user's password, and an error otherwise.
func (c *Console) Check(user, pw string) error {
	stored, known := c.records[user]
	if !known {
		stored = dummyRecord
	}
	ok, err := c.matches(stored, pw)
	if err != nil || !ok || !known {
		return errInvalidLogin
	}
	return nil
}

// matches reports whether pw opens a stored record, old sha256$ ones included, deriving one key through c.kdf for either kind.
func (c *Console) matches(stored, pw string) (bool, error) {
	if old, found := strings.CutPrefix(stored, "sha256$"); found {
		if _, err := c.matches(dummyRecord, pw); err != nil { // the work a new record costs
			return false, err
		}
		sum := sha256.Sum256([]byte(pw))
		typed := hex.EncodeToString(sum[:])
		return subtle.ConstantTimeCompare([]byte(typed), []byte(old)) == 1, nil
	}
	r, err := parseRecord(stored)
	if err != nil {
		return false, err
	}
	key, err := c.kdf(pw, r.Salt, r.Iterations)
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare(key, r.Key) == 1, nil
}

// needsUpgrade reports whether a stored record is an old sha256$ one or has fewer iterations than new ones.
func needsUpgrade(stored string) bool {
	if strings.HasPrefix(stored, "sha256$") {
		return true
	}
	r, err := parseRecord(stored)
	return err == nil && r.Iterations < iterations
}

// newRecord makes a record of pw with a fresh salt and today's count, deriving through c.kdf.
func (c *Console) newRecord(pw string) (string, error) {
	salt := newSalt()
	key, err := c.kdf(pw, salt, iterations)
	if err != nil {
		return "", err
	}
	return formatRecord(Record{Iterations: iterations, Salt: salt, Key: key}), nil
}
