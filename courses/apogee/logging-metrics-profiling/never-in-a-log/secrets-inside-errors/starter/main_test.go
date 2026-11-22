package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// The booking service's test key and token for Furrowlight, readable on purpose.
const (
	furrowKey   = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"
	furrowToken = "furrowlight.1830682800.767de44da59ef9b9e3afb9354270dd7cc56e26ceb00650b8a615afe1cd10b616"
)

func TestScrub(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"the 14th line's error", `verify token "` + furrowKey + `": invalid token`,
			`verify token "furrowlight.[redacted]": invalid token`},
		{"a key and a token in one header dump", "map[Authorization:[Key " + furrowKey + " Bearer " + furrowToken + "]]",
			"map[Authorization:[Key furrowlight.[redacted] Bearer furrowlight.[redacted]]]"},
		{"Brinewatch's key in a sentence", "key brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0 refused",
			"key brinewatch.[redacted] refused"},
		{"a near miss of 31 hex digits", "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6",
			"furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6"},
		{"a pass ID", "pass cirrine-1/LAN/16402: context deadline exceeded",
			"pass cirrine-1/LAN/16402: context deadline exceeded"},
		{"an address, which is maskEmail's work", "hollis.vey@furrowlight.example",
			"hollis.vey@furrowlight.example"},
		{"nothing at all", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := scrub(tt.s); got != tt.want {
				t.Errorf("scrub(%q) = %q, want %q", tt.s, got, tt.want)
			}
		})
	}
}

// compiles reports whether n calls regexp.Compile, regexp.MustCompile or one of their POSIX forms.
func compiles(n ast.Node) bool {
	found := false
	ast.Inspect(n, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "regexp" && strings.Contains(sel.Sel.Name, "Compile") {
			found = true
		}
		return true
	})
	return found
}

// TestScrubCompilesOnce reads redact.go: right and wrong code scrub alike, but only one of them
// compiles the pattern a single time.
func TestScrubCompilesOnce(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "redact.go", nil, 0)
	if err != nil {
		t.Fatalf("reading redact.go: %v", err)
	}
	once := false
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil && d.Name.Name == "scrub" && compiles(d.Body) {
				t.Errorf("redact.go: scrub calls regexp.MustCompile or Compile, so its pattern is compiled again for every string; want it compiled once, outside scrub")
			}
			if d.Recv == nil && d.Name.Name == "init" && compiles(d.Body) {
				once = true
			}
		case *ast.GenDecl:
			if d.Tok == token.VAR && compiles(d) {
				once = true
			}
		}
	}
	if !once {
		t.Errorf("redact.go has no variable set by regexp.MustCompile outside every function, want one that holds secretPattern, compiled once")
	}
}
