package main

// The simulated world's receivers: each customer's webhook receiver, as Yusra's sample
// receiver behaves, reached through World.HTTP's transport and never through the network.

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// simReceiver is one customer's receiver.
type simReceiver struct {
	ID     string // "kessick"
	Name   string // "Kessick Polytechnic": its satellites' operator
	Host   string // "hooks.kessick.example"
	Secret string // what its events are signed with
}

// simReceivers are the five customers' receivers.
var simReceivers = []simReceiver{
	{"apogee", "Apogee", "hooks.apogee.example", "apogee webhook secret"},
	{"furrowlight", "Furrowlight", "hooks.furrowlight.example", "furrowlight webhook secret"},
	{"brinewatch", "Brinewatch", "hooks.brinewatch.example", "brinewatch webhook secret"},
	{"cirrine", "Cirrine Atmospherics", "hooks.cirrine.example", "cirrine webhook secret"},
	{"kessick", "Kessick Polytechnic", "hooks.kessick.example", "kessick webhook secret"},
}

// simAttempt is one request sent to a receiver.
type simAttempt struct {
	Customer  string    // "cirrine"; "" when the request went to no receiver's host
	At        time.Time // when it was sent
	Down      bool      // the receiver was down: the request failed
	Status    int       // what the receiver answered; 0 when the request failed
	Problem   string    // why the receiver refused it; "" when it took the event
	Duplicate bool      // the receiver had the event already, and dropped it
	Timestamp string    // its Apogee-Timestamp header
	Signature string    // its Apogee-Signature header
	Type      string    // its Content-Type header
	Body      string    // its body, as sent
	Event     Event     // its body, decoded, when it is an event
}

// simHooks is every receiver of one run, and every request they were sent.
type simHooks struct {
	w        *simWorld
	mu       sync.Mutex
	seen     map[string]bool // customer and event ID, of the events taken
	attempts []simAttempt
}

// simNewHooks returns the receivers of a run.
func simNewHooks(w *simWorld) *simHooks {
	return &simHooks{w: w, seen: map[string]bool{}}
}

// RoundTrip hands a request to the receiver its host names. A receiver that is down fails
// the request; any other answers with a status and no body.
func (h *simHooks) RoundTrip(req *http.Request) (*http.Response, error) {
	now := h.w.now()
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(io.LimitReader(req.Body, 64<<10+1))
		req.Body.Close()
	}
	att := simAttempt{At: now, Timestamp: req.Header.Get("Apogee-Timestamp"), Signature: req.Header.Get("Apogee-Signature"),
		Type: req.Header.Get("Content-Type"), Body: string(body)}
	json.Unmarshal(body, &att.Event)
	var rc simReceiver
	for _, r := range simReceivers {
		if r.Host == req.URL.Hostname() {
			rc = r
		}
	}
	att.Customer = rc.ID
	var err error
	switch {
	case rc.ID == "":
		att.Problem = "no such host"
		err = fmt.Errorf("dial %s: no such host", req.URL.Host)
	default:
		att.Status, att.Problem = 0, ""
		for _, o := range h.w.c.receivers {
			if o.Customer == rc.ID && !now.Before(o.From) && now.Before(o.To) {
				att.Status, att.Down = o.Status, o.Status == 0
				att.Problem = "receiver in trouble"
			}
		}
		if att.Down {
			err = errors.New("connection refused")
		} else if att.Status == 0 {
			h.mu.Lock()
			att.Status, att.Problem, att.Duplicate = h.judge(rc, req, body, att.Event, now)
			h.mu.Unlock()
		}
	}
	h.mu.Lock()
	h.attempts = append(h.attempts, att)
	h.mu.Unlock()
	if err != nil {
		return nil, err
	}
	return &http.Response{
		Status:     fmt.Sprintf("%d %s", att.Status, http.StatusText(att.Status)),
		StatusCode: att.Status,
		Proto:      "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1,
		Header:  http.Header{},
		Body:    http.NoBody,
		Request: req,
	}, nil
}

// judge is Yusra's sample receiver: it refuses a request that isn't a POST to /apogee, a
// body over 64 KiB, a timestamp more than 300 seconds from its clock and a signature
// hmac.Equal rejects, and drops an event it already has. The mutex is held.
func (h *simHooks) judge(rc simReceiver, req *http.Request, body []byte, ev Event, now time.Time) (status int, problem string, duplicate bool) {
	if req.Method != http.MethodPost {
		return http.StatusMethodNotAllowed, "not a POST", false
	}
	if req.URL.Scheme != "https" || req.URL.Path != "/apogee" {
		return http.StatusNotFound, "not the receiver's address", false
	}
	if len(body) > 64<<10 {
		return http.StatusRequestEntityTooLarge, "body over 64 KiB", false
	}
	ts, err := strconv.ParseInt(req.Header.Get("Apogee-Timestamp"), 10, 64)
	if err != nil {
		return http.StatusBadRequest, "no timestamp in Unix seconds", false
	}
	if d := now.Unix() - ts; d > 300 || d < -300 {
		return http.StatusBadRequest, "timestamp more than 300 s from the receiver's clock", false
	}
	sig, err := hex.DecodeString(req.Header.Get("Apogee-Signature"))
	mac := hmac.New(sha256.New, []byte(rc.Secret))
	fmt.Fprintf(mac, "%d.", ts)
	mac.Write(body)
	if err != nil || !hmac.Equal(sig, mac.Sum(nil)) {
		return http.StatusUnauthorized, "bad signature", false
	}
	if ev.ID == "" || ev.Type != "data.landed" {
		return http.StatusBadRequest, "not a data.landed event", false
	}
	if h.seen[rc.ID+" "+ev.ID] {
		return http.StatusNoContent, "", true
	}
	h.seen[rc.ID+" "+ev.ID] = true
	return http.StatusNoContent, "", false
}
