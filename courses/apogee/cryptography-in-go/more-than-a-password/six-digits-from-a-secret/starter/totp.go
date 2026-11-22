package main

import "time"

// stepSeconds is how long each code lasts: RFC 6238's time step, in seconds.
const stepSeconds = 30

// timeStep returns how many whole 30-second steps have passed since 1970 at t: the counter for t's code.
func timeStep(t time.Time) uint64 {
	return uint64(t.Unix()) / stepSeconds
}

// totp returns the six-digit code an authenticator app shows at t, as RFC 6238 makes it.
func totp(secret []byte, t time.Time) string {
	return hotp(secret, timeStep(t))
}
