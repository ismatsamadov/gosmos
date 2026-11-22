package main

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"crypto/subtle"
)

// checkPassword reports whether pw is the password r was made from.
func checkPassword(r Record, pw string) (bool, error) {
	key, err := pbkdf2.Key(sha256.New, pw, r.Salt, r.Iterations, len(r.Key))
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare(key, r.Key) == 1, nil
}
