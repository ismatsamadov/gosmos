package main

import (
	"crypto/subtle"
	"time"
)

// Codes checks the six-digit codes from each engineer's authenticator app.
type Codes struct {
	secrets map[string][]byte // user → the 20-byte secret their phone holds too
	last    map[string]uint64 // user → the step of the last code accepted
	now     func() time.Time  // the console's clock
}

// Check reports whether code is user's code at c.now(), a step either side included, for a step after the last one it accepted.
func (c *Codes) Check(user, code string) bool {
	secret, ok := c.secrets[user]
	if !ok {
		return false
	}
	now := timeStep(c.now())
	for step := now - 1; step <= now+1; step++ {
		if subtle.ConstantTimeCompare([]byte(hotp(secret, step)), []byte(code)) == 1 {
			return true
		}
	}
	return false
}
