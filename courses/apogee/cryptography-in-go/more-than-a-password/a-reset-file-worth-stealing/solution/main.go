package main

import (
	"crypto/rand"
	"crypto/subtle"
)

// Start begins a password reset for user and returns the token for the link the console emails.
func (r *Resets) Start(user string) string {
	token := rand.Text()
	r.pending[user] = Reset{Hash: tokenDigest(token), Expires: r.now().Add(resetFor)}
	return token
}

// Redeem returns nil when token is the one in user's pending reset link, and errInvalidReset otherwise.
func (r *Resets) Redeem(user, token string) error {
	p, ok := r.pending[user]
	if !ok || subtle.ConstantTimeCompare([]byte(tokenDigest(token)), []byte(p.Hash)) != 1 {
		return errInvalidReset
	}
	return nil
}
