package main

import (
	"bytes"
	"crypto/hkdf"
	"crypto/sha256"
	"encoding/hex"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// testMaster is the tests' stand-in for k-2028-07: 32 readable bytes.
var testMaster = Key{ID: "k-2028-07", Secret: []byte("uplink test master key, 32 bytes")}

// shortHex returns a key's first 8 hex digits and an ellipsis, or "nothing" for no key.
func shortHex(key []byte) string {
	if len(key) == 0 {
		return "nothing"
	}
	return hex.EncodeToString(key[:min(len(key), 4)]) + "…"
}

func TestFleetKeys(t *testing.T) {
	keys, err := fleetKeys(testMaster, fleet)
	if err != nil {
		t.Fatalf("fleetKeys(testMaster, fleet) = %v, want 14 keys", err)
	}
	if len(keys) != 14 {
		t.Errorf("len(fleetKeys(testMaster, fleet)) = %d, want 14: a sign and a seal key for each of 7 satellites", len(keys))
	}
	const sealCirrine = "427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171"
	if got := hex.EncodeToString(keys["seal|cirrine-1"]); got != sealCirrine {
		t.Errorf("fleetKeys(testMaster, fleet)[%q] = %s, want 427a8ae7…", "seal|cirrine-1", shortHex(keys["seal|cirrine-1"]))
	}
	for _, sat := range fleet {
		for _, purpose := range []string{"sign", "seal"} {
			info := purpose + "|" + sat
			want, _ := hkdf.Key(sha256.New, testMaster.Secret, salt, info, 32)
			if got := keys[info]; !bytes.Equal(got, want) {
				t.Errorf("fleetKeys(testMaster, fleet)[%q] = %s, want %s, the key hkdf.Key derives with the same salt and info", info, shortHex(got), shortHex(want))
				return // the first wrong key says enough
			}
		}
	}
}

func TestEmptyList(t *testing.T) {
	if keys, err := fleetKeys(testMaster, nil); err != nil || len(keys) != 0 {
		t.Errorf("fleetKeys(testMaster, nil) = %d keys, %v, want 0 keys, nil", len(keys), err)
	}
}

// TestExtractOnce reads main.go. Calling hkdf.Key for every key gets every case above right
// too, but it runs the extract step fourteen times; so does hkdf.Extract inside the loop.
func TestExtractOnce(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", nil, 0)
	if err != nil {
		t.Fatalf("reading main.go: %v", err)
	}
	calls := map[string]int{}
	var stack []ast.Node
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
			return true
		}
		stack = append(stack, n)
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if pkg, ok := sel.X.(*ast.Ident); !ok || pkg.Name != "hkdf" {
			return true
		}
		calls[sel.Sel.Name]++
		if sel.Sel.Name == "Extract" {
			for _, outer := range stack {
				switch outer.(type) {
				case *ast.ForStmt, *ast.RangeStmt:
					t.Errorf("main.go calls hkdf.Extract inside a loop at line %d, want it once, before the loop", fset.Position(call.Pos()).Line)
				}
			}
		}
		return true
	})
	if calls["Extract"] != 1 {
		t.Errorf("main.go calls hkdf.Extract %d times, want 1", calls["Extract"])
	}
	if calls["Expand"] == 0 {
		t.Errorf("main.go calls hkdf.Expand 0 times, want a call for every key")
	}
	if calls["Key"] != 0 {
		t.Errorf("main.go calls hkdf.Key %d times, want 0: extract once, then expand", calls["Key"])
	}
}
