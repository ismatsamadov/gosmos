package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"strings"
	"testing"
	"time"
)

// clock is the console's clock in the tests: a set instant, moved by hand.
type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

// testResets returns an empty reset store whose clock reads Wednesday 25 October 2028, 10:00 UTC.
func testResets() *Resets {
	c := &clock{t: time.Date(2028, 10, 25, 10, 0, 0, 0, time.UTC)}
	return &Resets{pending: map[string]Reset{}, now: c.now}
}

// short cuts a value longer than 50 characters to its first 8 and an ellipsis.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

// TestBackupResetsNobody plays whoever walks off with a copy of the reset file: it holds every
// value the console keeps, and pastes each one into a link as the token.
func TestBackupResetsNobody(t *testing.T) {
	r := testResets()
	tokens := map[string]string{"mirela": r.Start("mirela"), "kiosk": r.Start("kiosk")}
	backup := maps.Clone(r.pending)
	for user, kept := range backup {
		if err := r.Redeem(user, kept.Hash); err == nil {
			t.Errorf("Redeem(%q, %q) = nil, want errInvalidReset: that is the Hash copied from the reset file, and a copy of the file must reset nobody", user, short(kept.Hash))
		}
	}
	for user, token := range tokens {
		if err := r.Redeem(user, token); err != nil {
			t.Errorf("Redeem(%q, %q) = %v, want nil: that is the token from the link", user, token, err)
		}
	}
}

// TestKeptDigest checks what the reset file holds: the token's hex SHA-256, never the token.
func TestKeptDigest(t *testing.T) {
	r := testResets()
	token := r.Start("mirela")
	if len(token) != 26 || strings.Trim(token, "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567") != "" {
		t.Fatalf("Start(%q) = %q, want 26 characters from A to Z and 2 to 7, as rand.Text makes", "mirela", token)
	}
	kept, ok := r.pending["mirela"]
	if !ok {
		t.Fatalf("r.pending[%q] is missing after Start, want the reset kept", "mirela")
	}
	if want := tokenDigest(token); kept.Hash != want {
		t.Errorf("r.pending[%q].Hash = %q, want tokenDigest(token) = %q", "mirela", short(kept.Hash), short(want))
	}
	if want := time.Date(2028, 10, 25, 10, 30, 0, 0, time.UTC); !kept.Expires.Equal(want) {
		t.Errorf("r.pending[%q].Expires = %v, want %v", "mirela", kept.Expires, want)
	}
	if other := r.Start("kiosk"); other == token {
		t.Errorf("Start(%q) and Start(%q) both = %q, want two different tokens", "mirela", "kiosk", token)
	}
}

func TestRedeem(t *testing.T) {
	r := testResets()
	token := r.Start("mirela")
	changed := token[:25] + "A"
	if changed == token {
		changed = token[:25] + "B"
	}
	tests := []struct {
		name, user, token string
	}{
		{"one character changed", "mirela", changed},
		{"another token", "mirela", "AAAAAAAAAAAAAAAAAAAAAAAAAA"},
		{"no token", "mirela", ""},
		{"no reset pending", "nkem", token},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Redeem(tt.user, tt.token); !errors.Is(err, errInvalidReset) {
				t.Errorf("Redeem(%q, %q) = %v, want errInvalidReset", tt.user, tt.token, err)
			}
		})
	}
	if err := r.Redeem("mirela", token); err != nil {
		t.Errorf("Redeem(%q, %q) = %v, want nil: that is the link's token", "mirela", token, err)
	}
}

// TestConstantTime reads main.go, helpers included. Comparing the digests with == or != gets
// every case above right too, but the console compares whatever stands in for a secret with
// subtle.ConstantTimeCompare, which takes as long whichever byte differs.
func TestConstantTime(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", nil, 0)
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}
	constantTime := 0
	ast.Inspect(file, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CallExpr:
			switch name := callName(n); name {
			case "subtle.ConstantTimeCompare", "hmac.Equal":
				constantTime++
			case "bytes.Equal", "slices.Equal", "reflect.DeepEqual", "strings.EqualFold", "strings.Compare":
				t.Errorf("main.go calls %s at line %d, want subtle.ConstantTimeCompare, which takes as long whichever byte differs", name, fset.Position(n.Pos()).Line)
			}
		case *ast.BinaryExpr:
			if (n.Op == token.EQL || n.Op == token.NEQ) && !fixed(n.X) && !fixed(n.Y) {
				t.Errorf("main.go compares two values with %s at line %d, want the digests compared with subtle.ConstantTimeCompare", n.Op, fset.Position(n.Pos()).Line)
			}
		}
		return true
	})
	if constantTime == 0 {
		t.Errorf("main.go calls subtle.ConstantTimeCompare or hmac.Equal 0 times, want the digests compared with one of them")
	}
}

// callName returns "pkg.Func" for a call such as bytes.Equal(a, b), and "" for other calls.
func callName(call *ast.CallExpr) string {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if pkg, ok := sel.X.(*ast.Ident); ok {
			return pkg.Name + "." + sel.Sel.Name
		}
	}
	return ""
}

// fixed reports whether e is nil, true, false, a literal or a length: values that say nothing
// about a secret, so main.go may compare them with == or !=, as in err != nil.
func fixed(e ast.Expr) bool {
	switch e := ast.Unparen(e).(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return e.Name == "nil" || e.Name == "true" || e.Name == "false"
	case *ast.CallExpr:
		fn, ok := e.Fun.(*ast.Ident)
		return ok && fn.Name == "len"
	}
	return false
}
