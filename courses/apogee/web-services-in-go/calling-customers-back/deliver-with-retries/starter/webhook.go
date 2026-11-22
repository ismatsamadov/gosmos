package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
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

// Delivery is one event on its way to one customer.
type Delivery struct {
	Customer string    `json:"customer"` // "furrowlight"
	Event    Event     `json:"event"`
	Attempts int       `json:"attempts"` // made so far
	NextAt   time.Time `json:"next_at"`  // when the next attempt is due
}

// ErrUndelivered means an event did not reach its customer and will not be tried again.
var ErrUndelivered = errors.New("undelivered")

// AnswerError is a receiver's answer outside 200 to 299: the event did not arrive.
type AnswerError struct {
	StatusCode int    // 503
	Status     string // "503 Service Unavailable"
}

// Error says what the receiver answered.
func (e *AnswerError) Error() string { return "receiver answered " + e.Status }

// post sends ev to url as a signed JSON POST, and returns an error unless the receiver answers 2xx.
func post(ctx context.Context, c *http.Client, url string, secret []byte, ts int64, ev Event) error {
	body, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Apogee-Timestamp", strconv.FormatInt(ts, 10))
	req.Header.Set("Apogee-Signature", signature(secret, ts, body))
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return &AnswerError{StatusCode: resp.StatusCode, Status: resp.Status}
	}
	return nil
}

// signature returns the hex HMAC-SHA256, under secret, of the timestamp, a dot, then the body.
func signature(secret []byte, ts int64, body []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(strconv.FormatInt(ts, 10) + "."))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// verifyEvent checks a signed request's timestamp and signature, and only then decodes the event in it.
func verifyEvent(r *http.Request, secret []byte, now time.Time) (Event, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 64<<10))
	if err != nil {
		return Event{}, err
	}
	ts, err := strconv.ParseInt(r.Header.Get("Apogee-Timestamp"), 10, 64)
	if err != nil {
		return Event{}, errors.New("timestamp missing or not a number")
	}
	if age := now.Unix() - ts; age < -300 || age > 300 {
		return Event{}, errors.New("timestamp more than 300 seconds from now")
	}
	if !hmac.Equal([]byte(r.Header.Get("Apogee-Signature")), []byte(signature(secret, ts, body))) {
		return Event{}, errors.New("signature does not match")
	}
	var ev Event
	if err := json.Unmarshal(body, &ev); err != nil {
		return Event{}, err
	}
	return ev, nil
}

// backoff returns the wait after failed attempt number attempt: from 2s, doubling up to 5 minutes, with jitter from r.
func backoff(attempt int, r *rand.Rand) time.Duration {
	wait := 2 * time.Second
	for i := 1; i < attempt && wait < 5*time.Minute; i++ {
		wait *= 2
	}
	wait = min(wait, 5*time.Minute)
	return time.Duration(float64(wait) * (0.5 + r.Float64()/2))
}

// retryable reports whether trying again can help, from the receiver's status, or from err when no answer came.
func retryable(status int, err error) bool {
	switch {
	case errors.Is(err, context.Canceled):
		return false
	case status == 0:
		return err != nil
	case status == http.StatusRequestTimeout, status == http.StatusTooManyRequests:
		return true
	default:
		return status >= 500 && status <= 599
	}
}

// statusOf returns the status of the receiver's answer that err carries, or 0 when no answer came.
func statusOf(err error) int {
	if answer, ok := errors.AsType[*AnswerError](err); ok {
		return answer.StatusCode
	}
	return 0
}

// Sender delivers events to customers' webhooks.
type Sender struct {
	Client    *http.Client                                     // Timeout 5s; never follows a redirect
	Customers map[string]Customer                              // by ID: each one's WebhookURL and WebhookKey
	Now       func() time.Time                                 // each attempt's timestamp
	Rand      *rand.Rand                                       // jitter; used only under mu
	Wait      func(ctx context.Context, d time.Duration) error // sleeps in production; tests record d
	mu        sync.Mutex
}

// schedule works out the wait after d's latest failed attempt and sets d.NextAt to when it ends.
func (s *Sender) schedule(d *Delivery) time.Duration {
	s.mu.Lock()
	wait := backoff(d.Attempts, s.Rand)
	s.mu.Unlock()
	d.NextAt = s.Now().Add(wait)
	return wait
}
