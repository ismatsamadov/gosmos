package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"time"
)

// resetFor is how long a reset link works after the console sends it.
const resetFor = 30 * time.Minute

// resetSent is the reset form's one answer, whatever name it is given.
const resetSent = "If that account exists, a reset link is on its way."

// errInvalidReset is the console's one answer to a reset link it won't take.
var errInvalidReset = errors.New("invalid or expired reset link")

// Reset is one pending password reset, as the console keeps it.
type Reset struct {
	Hash    string    // hex SHA-256 of the link's token
	Expires time.Time // the link is refused from this instant
}

// Resets holds the ops console's pending password resets, one per account.
type Resets struct {
	pending map[string]Reset         // user → their pending reset
	now     func() time.Time         // the console's clock
	known   func(user string) bool   // whether user has an account
	send    func(user, token string) // queues an email to user with the reset link
}

// tokenDigest returns the hex SHA-256 of a token: what the console keeps instead of it.
func tokenDigest(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// Start begins a password reset for user and returns the token for the link the console emails.
func (r *Resets) Start(user string) string {
	token := rand.Text()
	r.pending[user] = Reset{Hash: tokenDigest(token), Expires: r.now().Add(resetFor)}
	return token
}

// Redeem spends user's pending reset when token is its token and it hasn't expired, and returns errInvalidReset otherwise.
func (r *Resets) Redeem(user, token string) error {
	p, ok := r.pending[user]
	if !ok || !r.now().Before(p.Expires) || subtle.ConstantTimeCompare([]byte(tokenDigest(token)), []byte(p.Hash)) != 1 {
		return errInvalidReset
	}
	delete(r.pending, user)
	return nil
}
