package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"time"
)

// stepSeconds is how long each code lasts: RFC 6238's time step, in seconds.
const stepSeconds = 30

// hotp returns the six-digit code for one counter value, as RFC 4226 makes it.
func hotp(secret []byte, counter uint64) string {
	mac := hmac.New(sha1.New, secret)
	mac.Write(binary.BigEndian.AppendUint64(nil, counter))
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	n := binary.BigEndian.Uint32(sum[offset:]) & 0x7fffffff
	return fmt.Sprintf("%06d", n%1000000)
}

// timeStep returns how many whole 30-second steps have passed since 1970 at t: the counter for t's code.
func timeStep(t time.Time) uint64 {
	return uint64(t.Unix()) / stepSeconds
}

// totp returns the six-digit code an authenticator app shows at t, as RFC 6238 makes it.
func totp(secret []byte, t time.Time) string {
	return hotp(secret, timeStep(t))
}
