package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"sync"
	"testing"
)

// TestIncFromManyGoroutines counts 10,000 passes from 50 goroutines at once.
func TestIncFromManyGoroutines(t *testing.T) {
	var c Counter
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() {
			for range 200 {
				c.Inc()
			}
		})
	}
	wg.Wait()
	if got := c.Value(); got != 10000 {
		t.Errorf("Value() after 50 goroutines × 200 Inc() = %d, want 10000", got)
	}
}

func TestAdd(t *testing.T) {
	var c Counter
	c.Add(5)
	c.Inc()
	if got := c.Value(); got != 6 {
		t.Errorf("Value() after Add(5) and Inc() = %d, want 6", got)
	}
}

func TestNewCounter(t *testing.T) {
	var c Counter
	if got := c.Value(); got != 0 {
		t.Errorf("Value() of a new Counter = %d, want 0", got)
	}
}

// TestAtomicField reads metrics.go. A plain uint64 with ++ gives the right count most of the
// time, so only the code shows whether Counter keeps its count in an atomic.Uint64.
func TestAtomicField(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "metrics.go", nil, 0)
	if err != nil {
		t.Fatalf("reading metrics.go: %v", err)
	}
	var fields []string
	ast.Inspect(file, func(n ast.Node) bool {
		spec, ok := n.(*ast.TypeSpec)
		if !ok || spec.Name.Name != "Counter" {
			return true
		}
		if st, ok := spec.Type.(*ast.StructType); ok {
			for _, f := range st.Fields.List {
				fields = append(fields, types.ExprString(f.Type))
			}
		}
		return false
	})
	for _, f := range fields {
		if f == "atomic.Uint64" {
			return
		}
	}
	t.Errorf("Counter's field types in metrics.go = %v, want [atomic.Uint64]", fields)
}
