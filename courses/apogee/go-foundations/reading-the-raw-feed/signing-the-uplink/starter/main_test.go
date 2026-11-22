package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// benchKey signs commands on the lab bench at Port Ellery; flight keys never live in code.
var benchKey = []byte("lab bench key, never flown")

const (
	benchCmd = "pipit-2|118|SET_MODE|SAFE"
	// benchSig is the HMAC-SHA256 of benchCmd under benchKey, in hex.
	benchSig = "512545bcb43a5bb6ccb9f9c8511f4cc8c5fedeb38c226376af5b5aaae60e3b4d"
	// lastDigitChanged is benchSig with its last hex digit changed from d to e.
	lastDigitChanged = "512545bcb43a5bb6ccb9f9c8511f4cc8c5fedeb38c226376af5b5aaae60e3b4e"
)

func TestSign(t *testing.T) {
	if got := sign(benchCmd, benchKey); got != benchSig {
		t.Errorf("sign(%q, benchKey) = %q, want %q", benchCmd, got, benchSig)
	}
}

func TestVerify(t *testing.T) {
	cases := []struct {
		name, cmd, sig string
		key            []byte
		want           bool
	}{
		{"signed command", benchCmd, benchSig, benchKey, true},
		{"changed command", "pipit-2|118|SET_MODE|NOMINAL", benchSig, benchKey, false},
		{"changed signature", benchCmd, lastDigitChanged, benchKey, false},
		{"wrong key", benchCmd, benchSig, []byte("someone else's key"), false},
		{"empty signature", benchCmd, "", benchKey, false},
		{"not hex", benchCmd, "sig=" + benchSig, benchKey, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := verify(c.cmd, c.sig, c.key); got != c.want {
				t.Errorf("verify(%q, %q, %q) = %v, want %v", c.cmd, c.sig, c.key, got, c.want)
			}
		})
	}
}

// TestVerifyUsesHMACEqual reads main.go. A comparison with == gets every case above right, but
// it can stop at the first difference, so its timing tells an attacker how much of a forgery
// is right.
func TestVerifyUsesHMACEqual(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, 0)
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}
	uses := 0
	ast.Inspect(file, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel.Name == "Equal" {
			if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "hmac" {
				uses++
			}
		}
		return true
	})
	if uses == 0 {
		t.Errorf("main.go calls hmac.Equal 0 times, want verify to compare the signatures with it")
	}
}
