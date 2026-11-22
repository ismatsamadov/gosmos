package main

import (
	"bytes"
	"testing"
)

// TestRestart starts two Noncers, as two runs of uplinkd would, and compares their first nonces.
func TestRestart(t *testing.T) {
	var first, second Noncer
	a, b := first.Nonce(), second.Nonce()
	if bytes.Equal(a, b) {
		t.Errorf("Nonce() = %x on the first call of each run, want two different nonces", a)
	}
}

// TestThousand takes 1,000 nonces from one Noncer: each 12 bytes, and none twice.
func TestThousand(t *testing.T) {
	var n Noncer
	seen := make(map[string]int, 1000)
	for call := 1; call <= 1000; call++ {
		nonce := n.Nonce()
		if len(nonce) != 12 {
			t.Fatalf("len(Nonce()) = %d on call %d, want 12", len(nonce), call)
		}
		if first, ok := seen[string(nonce)]; ok {
			t.Fatalf("Nonce() = %x on calls %d and %d, want 1,000 different nonces", nonce, first, call)
		}
		seen[string(nonce)] = call
	}
}
