package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"reflect"
	"sync"
	"testing"
)

// The lab bench key signs the tests' commands; flight keys never appear in code.
var benchKey = []byte("lab bench key, never flown")

// hmacHex returns the hex HMAC-SHA256 of one canonical command alone under key.
func hmacHex(key []byte, canonical string) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(canonical))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestSign(t *testing.T) {
	s := NewSigner(benchKey)
	s.Sign("pipit-2|118|SET_MODE|SAFE")
	const second = "pipit-2|119|PING"
	if got, want := s.Sign(second), hmacHex(benchKey, second); got != want {
		t.Fatalf("second Sign(%q) = %.8s…, want %.8s…, the HMAC of that command alone", second, got, want)
	}

	var wg sync.WaitGroup
	wrong := make(chan string, 400)
	for w := range 8 {
		wg.Go(func() {
			for i := range 50 {
				cmd := fmt.Sprintf("pipit-2|%d|PING", 1000+50*w+i)
				if s.Sign(cmd) != hmacHex(benchKey, cmd) {
					wrong <- cmd
				}
			}
		})
	}
	wg.Wait()
	close(wrong)
	if n := len(wrong); n > 0 {
		t.Errorf("Sign from 8 goroutines at once: %d of 400 signatures wrong, the first for %q, want none", n, <-wrong)
	}
}

// TestNoSharedHash looks inside Signer: a hash.Hash kept there is shared by every message it signs.
func TestNoSharedHash(t *testing.T) {
	hashType := reflect.TypeFor[hash.Hash]()
	signer := reflect.TypeFor[Signer]()
	for i := range signer.NumField() {
		if f := signer.Field(i); f.Type == hashType || f.Type.Implements(hashType) {
			t.Errorf("Signer has field %s of type %s, want a new HMAC made in every Sign", f.Name, f.Type)
		}
	}
}
