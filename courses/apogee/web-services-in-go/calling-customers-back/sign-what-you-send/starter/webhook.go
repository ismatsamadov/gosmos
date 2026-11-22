package main

import "time"

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

// AnswerError is a receiver's answer outside 200 to 299: the event did not arrive.
type AnswerError struct {
	StatusCode int    // 503
	Status     string // "503 Service Unavailable"
}

// Error says what the receiver answered.
func (e *AnswerError) Error() string { return "receiver answered " + e.Status }
