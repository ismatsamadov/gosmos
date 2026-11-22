package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"strings"
	"testing"
	"time"
)

// answer is what the service sent back to one request.
type answer struct {
	Status int
	Header http.Header
	Body   []byte
}

// errorCode returns the "error" of a JSON ErrorBody answer, or "" when the answer is not one.
func (a answer) errorCode() string {
	mediaType, _, _ := mime.ParseMediaType(a.Header.Get("Content-Type"))
	var body ErrorBody
	if mediaType != "application/json" || json.Unmarshal(a.Body, &body) != nil {
		return ""
	}
	return body.Code
}

// request sends method and path to base with token as a Bearer token, body as JSON and header as name, value pairs.
func request(t *testing.T, base, method, path, token string, body any, header ...string) answer {
	t.Helper()
	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		r = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, base+path, r)
	if err != nil {
		t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for i := 0; i+1 < len(header); i += 2 {
		req.Header.Set(header[i], header[i+1])
	}
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%s %s: reading the answer: %v", method, path, err)
	}
	return answer{resp.StatusCode, resp.Header, data}
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
