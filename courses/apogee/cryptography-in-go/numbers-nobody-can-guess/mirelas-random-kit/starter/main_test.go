package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestNewKey(t *testing.T) {
	a, b := newKey("k-2028-07"), newKey("k-2028-07")
	if a.ID != "k-2028-07" {
		t.Errorf("newKey(%q).ID = %q, want %q", "k-2028-07", a.ID, "k-2028-07")
	}
	if len(a.Secret) != 32 {
		t.Errorf("len(newKey(%q).Secret) = %d, want 32", "k-2028-07", len(a.Secret))
	} else if bytes.Equal(a.Secret, b.Secret) {
		t.Errorf("newKey(%q).Secret = %x… twice in a row, want two different secrets", "k-2028-07", a.Secret[:4])
	}
}

func TestNewNonce(t *testing.T) { checkFresh(t, "newNonce()", newNonce, 12) }

func TestNewSalt(t *testing.T) { checkFresh(t, "newSalt()", newSalt, 16) }

// checkFresh checks that two calls of gen, written call, return size bytes each, and different ones.
func checkFresh(t *testing.T, call string, gen func() []byte, size int) {
	t.Helper()
	a, b := gen(), gen()
	if len(a) != size {
		t.Errorf("len(%s) = %d, want %d", call, len(a), size)
	} else if bytes.Equal(a, b) {
		t.Errorf("%s = %x twice in a row, want two different values", call, a)
	}
}

func TestNewSessionToken(t *testing.T) {
	const base32 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	a, b := newSessionToken(), newSessionToken()
	if len(a) < 26 || strings.Trim(a, base32) != "" {
		t.Errorf("newSessionToken() = %q, want 26 or more characters from A to Z and 2 to 7", a)
	} else if a == b {
		t.Errorf("newSessionToken() = %q twice in a row, want two different tokens", a)
	}
}

// TestOnlyCryptoRand reads the imports of every .go file in the lesson. random.go must take
// its bytes from crypto/rand, and no file may import math/rand or math/rand/v2, whose numbers
// anyone with the seed can make again.
func TestOnlyCryptoRand(t *testing.T) {
	names, err := filepath.Glob("*.go")
	if err != nil || len(names) == 0 {
		t.Fatalf("listing the lesson's .go files: %v", err)
	}
	fromCrypto := false
	for _, name := range names {
		file, err := parser.ParseFile(token.NewFileSet(), name, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("reading %s: %v", name, err)
		}
		for _, imp := range file.Imports {
			path, _ := strconv.Unquote(imp.Path.Value)
			if path == "math/rand" || path == "math/rand/v2" {
				t.Errorf("%s imports %q, want no file to import it", name, path)
			}
			if path == "crypto/rand" && name == "random.go" {
				fromCrypto = true
			}
		}
	}
	if !fromCrypto {
		t.Errorf("random.go imports crypto/rand = false, want true")
	}
}
