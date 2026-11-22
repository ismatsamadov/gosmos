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
	_ = bytes.NewReader
	_ = body // TODO: POST body to url through c, close the answer's body, and refuse any answer outside 2xx
	return nil
}
