package main

import (
	"crypto/hmac"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

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
