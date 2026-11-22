package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// stored is the tests' maintenance phrase; the real one is never in code.
const stored = "stow before the wind does"

func TestSamePhrase(t *testing.T) {
	tests := []struct {
		name, typed string
		want        bool
	}{
		{"the phrase", stored, true},
		{"one wrong character", "stow before the wind dies", false},
		{"its first two words", "stow before", false},
		{"one word too many", stored + " it", false},
		{"nothing typed", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := samePhrase(tt.typed, stored); got != tt.want {
				t.Errorf("samePhrase(%q, stored) = %v, want %v", tt.typed, got, tt.want)
			}
		})
	}
}

// TestNoEquals reads main.go, samePhrase and any helper beside it. == gives the same answers
// as subtle.ConstantTimeCompare, but it stops at the first difference, so how long it takes
// tells an attacker how much of a guess was right. Comparing the result with a number, as
// in ConstantTimeCompare(...) == 1, is fine.
func TestNoEquals(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", nil, 0)
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}
	calls := 0
	ast.Inspect(file, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CallExpr:
			if sel, ok := n.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "ConstantTimeCompare" {
				if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "subtle" {
					calls++
				}
			}
		case *ast.BinaryExpr:
			if (n.Op == token.EQL || n.Op == token.NEQ) && !isNumber(n.X) && !isNumber(n.Y) {
				t.Errorf("main.go line %d compares with %s, want subtle.ConstantTimeCompare to compare the phrases", fset.Position(n.Pos()).Line, n.Op)
			}
		}
		return true
	})
	if calls == 0 {
		t.Errorf("main.go calls subtle.ConstantTimeCompare 0 times, want 1")
	}
}

// isNumber reports whether e is a whole number written in the code, such as the 1 that
// subtle.ConstantTimeCompare returns for equal slices.
func isNumber(e ast.Expr) bool {
	lit, ok := e.(*ast.BasicLit)
	return ok && lit.Kind == token.INT
}

// TestHashesFirst reads main.go. subtle.ConstantTimeCompare on the phrases themselves gives
// the right answers too, but it answers at once when their lengths differ, so how long it
// takes tells an attacker the phrase's length.
func TestHashesFirst(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, 0)
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}
	hashes := 0
	ast.Inspect(file, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && (sel.Sel.Name == "Sum256" || sel.Sel.Name == "New") {
				if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "sha256" {
					hashes++
				}
			}
		}
		return true
	})
	if hashes == 0 {
		t.Errorf("main.go calls sha256.Sum256 0 times, want both phrases hashed before they are compared")
	}
}
