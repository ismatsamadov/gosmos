package main

import (
	"crypto/sha256"
	"time"
)

// claim decides what a request with key id and this body gets: a refusal, a kept answer to replay, or nil, nil to run it.
func (s *keyStore) claim(id keyID, body []byte) (*kept, error) {
	sum := sha256.Sum256(body)
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	if k, ok := s.keys[id]; ok && now.Sub(k.At) < 24*time.Hour {
		switch {
		case k.Sum != sum:
			return nil, differentRequest(id.Key)
		case !k.Done:
			return nil, inProgress(id.Key)
		}
		return k, nil
	}
	s.keys[id] = &kept{Sum: sum, At: now}
	return nil, nil
}
