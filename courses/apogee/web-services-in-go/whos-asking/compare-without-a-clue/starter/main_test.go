package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// The beta customers' test API keys, as Yusra gave them out.
const (
	apogeeKey      = "apogee.9e4d2a7c1b0f3e5d8a6c4b2f0e1d3a5c"
	furrowlightKey = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"
	brinewatchKey  = "brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0"
	cirrineKey     = "cirrine.c3a1f7e95d2b4086a7e1c9f3b5d8a024"
)

// keys is today's table: each customer's whole key, by customer ID.
var keys = map[string]string{
	"apogee":      apogeeKey,
	"furrowlight": furrowlightKey,
	"brinewatch":  brinewatchKey,
	"cirrine":     cirrineKey,
}

func TestCheckKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string // "" when the key must be refused
	}{
		{"Apogee's own key", apogeeKey, "apogee"},
		{"Furrowlight's key", furrowlightKey, "furrowlight"},
		{"Brinewatch's key", brinewatchKey, "brinewatch"},
		{"Cirrine's key", cirrineKey, "cirrine"},
		{"a wrong secret", "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6c", ""},
		{"a customer the service doesn't know", "kessick.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b", ""},
		{"Furrowlight's secret with Brinewatch's name", "brinewatch.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b", ""},
		{"an empty key", "", ""},
		{"a key without a dot", "furrowlight5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b", ""},
		{"a key one digit short", "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := checkKey(keys, tt.key)
			if got != tt.want || ok != (tt.want != "") {
				t.Errorf("checkKey(keys, %q) = %q, %v, want %q, %v", tt.key, got, ok, tt.want, tt.want != "")
			}
		})
	}
}

func TestCheckKeyComparesInConstantTime(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "auth.go", nil, 0)
	if err != nil {
		t.Fatalf("parser.ParseFile(auth.go) error = %v, want nil", err)
	}
	calls := 0
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "checkKey" || fn.Body == nil {
			continue
		}
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "subtle" && sel.Sel.Name == "ConstantTimeCompare" {
				calls++
			}
			return true
		})
	}
	if calls == 0 {
		t.Errorf("calls to subtle.ConstantTimeCompare in checkKey = 0, want at least 1: == can stop at the first difference")
	}
}
