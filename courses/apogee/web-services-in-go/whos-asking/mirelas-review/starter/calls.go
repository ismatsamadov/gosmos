package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// The beta customers' test API keys, as Yusra gave them out.
const (
	furrowlightKey = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"
	brinewatchKey  = "brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0"
	cirrineKey     = "cirrine.c3a1f7e95d2b4086a7e1c9f3b5d8a024"
)

// privateRoutes are the service's routes that need a token, with paths that name a real pass and booking.
var privateRoutes = []string{
	"GET /passes",
	"GET /passes/furrow-a/SQT/16351",
	"POST /bookings",
	"GET /bookings",
	"GET /bookings/bk-4501",
	"DELETE /bookings/bk-4501",
}

// send makes one request to the service at base and returns the answer's status and body. route is a
// method and a path, such as "GET /bookings", and auth the Authorization header ("" sends none).
func send(t *testing.T, base, route, auth string) (int, []byte) {
	t.Helper()
	method, path, _ := strings.Cut(route, " ")
	req, err := http.NewRequest(method, base+path, nil)
	if err != nil {
		t.Fatalf("%s: %v", route, err)
	}
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	client := &http.Client{Timeout: 2 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("%s: %v", route, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%s: reading the answer: %v", route, err)
	}
	return res.StatusCode, body
}

// call sends route to the service at base with auth, and returns the status and, for an error, its message.
func call(t *testing.T, base, route, auth string) (status int, message string) {
	t.Helper()
	status, body := send(t, base, route, auth)
	var e ErrorBody
	if status >= 400 && json.Unmarshal(body, &e) == nil {
		message = e.Message
	}
	return status, message
}

// token trades key for a token at POST /tokens, and stops the review if the service refuses.
func token(t *testing.T, base, key string) string {
	t.Helper()
	status, body := send(t, base, "POST /tokens", "Key "+key)
	var tr TokenResponse
	if err := json.Unmarshal(body, &tr); status != http.StatusOK || err != nil {
		t.Fatalf("POST /tokens with a key = %d %s, want 200 and a token", status, body)
	}
	return tr.Token
}

// bookingsOf returns the bookings GET /bookings lists for the token tok, and stops the review if it can't.
func bookingsOf(t *testing.T, base, tok string) []Booking {
	t.Helper()
	status, body := send(t, base, "GET /bookings", "Bearer "+tok)
	var list []Booking
	if err := json.Unmarshal(body, &list); status != http.StatusOK || err != nil {
		t.Fatalf("GET /bookings = %d %s, want 200 and a list", status, body)
	}
	return list
}
