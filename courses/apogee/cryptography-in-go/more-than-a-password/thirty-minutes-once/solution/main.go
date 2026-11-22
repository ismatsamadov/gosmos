package main

import "crypto/subtle"

// Redeem spends user's pending reset when token is its token and it hasn't expired, and returns errInvalidReset otherwise.
func (r *Resets) Redeem(user, token string) error {
	p, ok := r.pending[user]
	if !ok || !r.now().Before(p.Expires) || subtle.ConstantTimeCompare([]byte(tokenDigest(token)), []byte(p.Hash)) != 1 {
		return errInvalidReset
	}
	delete(r.pending, user)
	return nil
}
