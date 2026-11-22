package main

import (
	"crypto/subtle"
	"time"
)

// Codes checks the six-digit codes from each engineer's authenticator app.
type Codes struct {
	secrets map[string][]byte // user → the 20-byte secret their phone holds too
	now     func() time.Time  // the console's clock
}

// Check reports whether code is user's code at c.now(), the step before and the step after included.
func (c *Codes) Check(user, code string) bool {
	secret, ok := c.secrets[user]
	if !ok {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(totp(secret, c.now())), []byte(code)) == 1
}
