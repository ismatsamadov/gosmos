package main

import (
	"crypto/subtle"
	"errors"
)

// Console is the ops console's login.
type Console struct {
	records map[string]string                                      // user → record text
	kdf     func(pw string, salt []byte, iter int) ([]byte, error) // PBKDF2-HMAC-SHA256; tests count calls
}

// errInvalidLogin is the console's one answer to a failed login.
var errInvalidLogin = errors.New("invalid user or password")

// dummyRecord is a valid 600,000-iteration record whose password nobody knows. An unknown
// user is checked against it, so a failed login costs the same work whoever tries.
const dummyRecord = "pbkdf2-sha256$600000$64756d6d79207265636f72642c203136$6e6f2070617373776f726420646572697665732074686973206b65792c203332"

// matches reports whether pw opens a stored record, deriving its key through c.kdf.
func (c *Console) matches(stored, pw string) (bool, error) {
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
