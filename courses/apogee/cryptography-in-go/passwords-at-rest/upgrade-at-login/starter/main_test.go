package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

const (
	mirelaPassword = "fails closed, always"
	saltHex        = "6d6972656c6120746573742073616c74" // "mirela test salt", the tests' readable salt
	// mirelaRecord is Mirela's password at 600,000 iterations: a record as strong as new ones.
	mirelaRecord = "pbkdf2-sha256$600000$" + saltHex + "$6d91ceab779ff574b23820b133a740fa76d8b6a0deb5d6f0da2f94a6508dd938"
	// benchRecord is the same password at 1,000 iterations, as the lab bench's account was made.
	benchRecord = "pbkdf2-sha256$1000$" + saltHex + "$5b893367890cac1547d6c07d906a4b46b8a89733bf1f7fb3d7ba24f8e16174b3"
	// strongerRecord is the same password at 600,001 iterations: one more than new records use.
	strongerRecord = "pbkdf2-sha256$600001$" + saltHex + "$d6f613b555d5c2d60ad22f13f06cf3c6e5756f38dd0058b22d77df3854831823"
	// kioskOld is the kiosk's record from the stolen file: the unsalted SHA-256 of "fishbowl".
	kioskOld = "sha256$1c4bb463c52d458ae6e9225430c45c753c5ae417bb4e7e4530773aaa0cda201b"
)

// short cuts each long field of a record to its first 8 characters and an ellipsis.
func short(record string) string {
	if len(record) <= 50 {
		return record
	}
	fields := strings.Split(record, "$")
	for i, f := range fields {
		if len(f) > 16 {
			fields[i] = f[:8] + "…"
		}
	}
	return strings.Join(fields, "$")
}

// TestUpgradeOld logs in with the kiosk's old record, then with the record that replaces it,
// which must accept the same password and need no upgrade of its own.
func TestUpgradeOld(t *testing.T) {
	ok, upgraded, err := login(kioskOld, "fishbowl")
	if !ok || upgraded == "" || err != nil {
		t.Fatalf("login(kioskOld, %q) = %v, %q, %v, want true, a new record, nil", "fishbowl", ok, short(upgraded), err)
	}
	if r, err := parseRecord(upgraded); err != nil || r.Iterations != 600000 {
		t.Fatalf("login(kioskOld, %q) returned %q, want a pbkdf2-sha256 record at 600000 iterations", "fishbowl", short(upgraded))
	}
	ok, again, err := login(upgraded, "fishbowl")
	if !ok || again != "" || err != nil {
		t.Errorf("login(the new record, %q) = %v, %q, %v, want true, \"\", nil", "fishbowl", ok, short(again), err)
	}
}

func TestUpgradeFewIterations(t *testing.T) {
	ok, upgraded, err := login(benchRecord, mirelaPassword)
	if !ok || err != nil {
		t.Fatalf("login(benchRecord, %q) = %v, %q, %v, want true, a new record, nil", mirelaPassword, ok, short(upgraded), err)
	}
	r, err := parseRecord(upgraded)
	if err != nil || r.Iterations != 600000 || string(r.Salt) == "mirela test salt" {
		t.Errorf("login(benchRecord, %q) returned %q, want a pbkdf2-sha256 record at 600000 iterations with a new salt", mirelaPassword, short(upgraded))
	}
}

func TestCurrentRecord(t *testing.T) {
	ok, upgraded, err := login(mirelaRecord, mirelaPassword)
	if !ok || upgraded != "" || err != nil {
		t.Errorf("login(mirelaRecord, %q) = %v, %q, %v, want true, \"\", nil", mirelaPassword, ok, short(upgraded), err)
	}
}

// TestStrongerRecord logs in with a record made with more iterations than new records use. It
// isn't weak, so replacing it with a 600,000-iteration record would make it weaker.
func TestStrongerRecord(t *testing.T) {
	ok, upgraded, err := login(strongerRecord, mirelaPassword)
	if !ok || upgraded != "" || err != nil {
		t.Errorf("login(strongerRecord, %q) = %v, %q, %v, want true, \"\", nil", mirelaPassword, ok, short(upgraded), err)
	}
}

func TestWrongPassword(t *testing.T) {
	tests := []struct{ name, owner, record, pw string }{
		{"the kiosk's old record", "kioskOld", kioskOld, "fishbowl!"},
		{"the bench's 1,000-iteration record", "benchRecord", benchRecord, "fishbowl"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok, upgraded, err := login(tt.record, tt.pw); ok || upgraded != "" || err != nil {
				t.Errorf("login(%s, %q) = %v, %q, %v, want false, \"\", nil", tt.owner, tt.pw, ok, short(upgraded), err)
			}
		})
	}
}

func TestUnreadableRecord(t *testing.T) {
	for _, record := range []string{
		"md5$" + saltHex,
		"pbkdf2-sha256$999$" + saltHex + "$5b893367890cac1547d6c07d906a4b46b8a89733bf1f7fb3d7ba24f8e16174b3",
	} {
		if ok, upgraded, err := login(record, mirelaPassword); ok || upgraded != "" || err == nil {
			t.Errorf("login(%q, %q) = %v, %q, %v, want false, \"\" and an error", short(record), mirelaPassword, ok, short(upgraded), err)
		}
	}
}

// TestConstantTime reads main.go, helpers included. Comparing an old record's digest with ==
// gets every case above right too, but it stops at the first character that differs, so the
// time a refusal takes says how much of the stored digest a guess's digest matched.
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
				t.Errorf("main.go compares two values with %s at line %d, want the old record compared with subtle.ConstantTimeCompare", n.Op, fset.Position(n.Pos()).Line)
			}
		}
		return true
	})
	if constantTime == 0 {
		t.Errorf("main.go calls subtle.ConstantTimeCompare or hmac.Equal 0 times, want an old record compared with one of them")
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

// fixed reports whether e is nil, true, false, a literal, a length or an iteration count:
// values that say nothing about a password, so main.go may compare them with == or !=.
func fixed(e ast.Expr) bool {
	switch e := ast.Unparen(e).(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return e.Name == "nil" || e.Name == "true" || e.Name == "false" || e.Name == "iterations"
	case *ast.SelectorExpr:
		return e.Sel.Name == "Iterations"
	case *ast.CallExpr:
		fn, ok := e.Fun.(*ast.Ident)
		return ok && fn.Name == "len"
	}
	return false
}
