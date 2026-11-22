package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

const (
	mirelaPassword = "fails closed, always"
	saltHex        = "6d6972656c6120746573742073616c74" // "mirela test salt", the tests' readable salt
	// mirelaRecord is Mirela's password at 600,000 iterations with the test salt.
	mirelaRecord = "pbkdf2-sha256$600000$" + saltHex + "$6d91ceab779ff574b23820b133a740fa76d8b6a0deb5d6f0da2f94a6508dd938"
	// benchRecord is the same password and salt at 1,000 iterations, as the lab bench's
	// account was made.
	benchRecord = "pbkdf2-sha256$1000$" + saltHex + "$5b893367890cac1547d6c07d906a4b46b8a89733bf1f7fb3d7ba24f8e16174b3"
)

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name, owner, record, pw string
		want                    bool
	}{
		{"Mirela's password", "Mirela's record", mirelaRecord, mirelaPassword, true},
		{"a period too many", "Mirela's record", mirelaRecord, mirelaPassword + ".", false},
		{"nothing typed", "Mirela's record", mirelaRecord, "", false},
		{"the bench at 1,000 iterations", "the bench's record", benchRecord, mirelaPassword, true},
		{"the bench and the kiosk's old password", "the bench's record", benchRecord, "fishbowl", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := parseRecord(tt.record)
			if err != nil {
				t.Fatalf("parseRecord(%s) = %v, want nil", tt.owner, err)
			}
			if got, err := checkPassword(r, tt.pw); got != tt.want || err != nil {
				t.Errorf("checkPassword(%s, %q) = %v, %v, want %v, nil", tt.owner, tt.pw, got, err, tt.want)
			}
		})
	}
}

// TestConstantTime reads main.go, helpers included. Comparing the keys with bytes.Equal or ==
// gets every case above right too, but it stops at the first byte that differs, so the time a
// refusal takes says how many bytes of the stored key a guess's key matched.
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
			case "bytes.Equal", "slices.Equal", "reflect.DeepEqual":
				t.Errorf("main.go calls %s at line %d, want subtle.ConstantTimeCompare, which takes as long whichever byte differs", name, fset.Position(n.Pos()).Line)
			}
		case *ast.BinaryExpr:
			if (n.Op == token.EQL || n.Op == token.NEQ) && !fixed(n.X) && !fixed(n.Y) {
				t.Errorf("main.go compares two values with %s at line %d, want the keys compared with subtle.ConstantTimeCompare", n.Op, fset.Position(n.Pos()).Line)
			}
		}
		return true
	})
	if constantTime == 0 {
		t.Errorf("main.go calls subtle.ConstantTimeCompare or hmac.Equal 0 times, want the keys compared with one of them")
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
// about a key, so main.go may compare them with == or !=, as in err != nil.
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
