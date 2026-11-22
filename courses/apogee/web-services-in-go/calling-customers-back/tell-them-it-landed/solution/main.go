package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// post sends ev to url as a JSON POST, and returns an error unless the receiver answers 2xx.
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
