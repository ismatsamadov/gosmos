package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// Event is what the service tells a customer when something happens to a booking.
type Event struct {
	ID        string    `json:"id"`   // "evt-000118": the same on every attempt
	Type      string    `json:"type"` // "data.landed", the only type so far
	BookingID string    `json:"booking_id"`
	SatID     string    `json:"sat_id"`
	Station   string    `json:"station"`
	File      string    `json:"file"`
	Bytes     int64     `json:"bytes"`
	LandedAt  time.Time `json:"landed_at"`
}

// signature returns the hex HMAC-SHA256, under secret, of the timestamp, a dot, then the body.
func signature(secret []byte, ts int64, body []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(strconv.FormatInt(ts, 10) + "."))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
