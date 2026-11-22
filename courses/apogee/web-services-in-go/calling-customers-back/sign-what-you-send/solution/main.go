package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
)

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
