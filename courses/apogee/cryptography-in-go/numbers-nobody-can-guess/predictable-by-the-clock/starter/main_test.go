package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"math/rand/v2"
	"strconv"
	"testing"
	"time"
)

// clockNonce makes the nonce the draft helper made in the given second, as Mirela would.
func clockNonce(sec int64) []byte {
	r := rand.New(rand.NewPCG(uint64(sec), 0))
	nonce := make([]byte, 12)
	for i := range nonce {
		nonce[i] = byte(r.Uint32())
	}
	return nonce
}

// TestAttackerPredicts plays Mirela. She knows roughly when newNonce ran, so she makes the
// draft helper's nonce for every second from 2 s before the call to 2 s after it.
func TestAttackerPredicts(t *testing.T) {
	from := time.Now().Unix() - 2
	nonce := newNonce()
	to := time.Now().Unix() + 2
	for sec := from; sec <= to; sec++ {
		if bytes.Equal(nonce, clockNonce(sec)) {
			t.Errorf("newNonce() = %x, the nonce of the clock's second %d, want bytes no clock can make again", nonce, sec)
			return
		}
	}
}

func TestLength(t *testing.T) {
	if got := len(newNonce()); got != 12 {
		t.Errorf("len(newNonce()) = %d, want 12", got)
	}
}

func TestTwoDiffer(t *testing.T) {
	a, b := newNonce(), newNonce()
	if bytes.Equal(a, b) {
		t.Errorf("newNonce() = %x twice in a row, want two different nonces", a)
	}
}

// TestNoMathRand reads main.go's imports. math/rand and math/rand/v2 make numbers that anyone
// with the seed can make again, so no secret may come from them.
func TestNoMathRand(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}
	for _, imp := range file.Imports {
		if path, _ := strconv.Unquote(imp.Path.Value); path == "math/rand" || path == "math/rand/v2" {
			t.Errorf("main.go imports %q, want crypto/rand in its place", path)
		}
	}
}
