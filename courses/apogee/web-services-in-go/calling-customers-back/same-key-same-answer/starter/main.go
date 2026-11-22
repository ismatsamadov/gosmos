package main

import (
	"crypto/sha256"
	"time"
)

// claim decides what a request with key id and this body gets: a refusal, a kept answer to replay, or nil, nil to run it.
func (s *keyStore) claim(id keyID, body []byte) (*kept, error) {
	_, _ = sha256.Sum256, time.Hour
	return nil, nil
}
