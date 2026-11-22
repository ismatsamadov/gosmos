package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand/v2"
	"sync"
	"testing"
	"time"
)

// run returns the schedule of the planning run numbered v.
func run(v uint64) *Schedule { return &Schedule{Version: v} }

// swaps reports whether publish.go calls CompareAndSwap anywhere.
func swaps(t *testing.T) bool {
	f, err := parser.ParseFile(token.NewFileSet(), "publish.go", nil, 0)
	if err != nil {
		t.Fatalf("reading publish.go: %v", err)
	}
	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel.Name == "CompareAndSwap" {
			found = true
		}
		return !found
	})
	return found
}

func TestOffer(t *testing.T) {
	t.Run("run 1's plan, the first", func(t *testing.T) {
		var p Publisher
		if !p.Offer(run(1)) {
			t.Errorf("Offer(version 1) to an empty publisher = false, want true")
		}
		if cur := p.Current(); cur == nil || cur.Version != 1 {
			t.Errorf("Current() after Offer(version 1) is not version 1")
		}
	})
	t.Run("run 2 finishes before run 1", func(t *testing.T) {
		var p Publisher
		p.Offer(run(2))
		if p.Offer(run(1)) {
			t.Errorf("Offer(version 1) after version 2 = true, want false")
		}
		if cur := p.Current(); cur == nil || cur.Version != 2 {
			t.Errorf("Current() after versions 2 and 1 is not version 2, want the newer plan to stay")
		}
	})
	t.Run("run 2 offered twice", func(t *testing.T) {
		var p Publisher
		first, again := run(2), run(2)
		p.Offer(first)
		if p.Offer(again) {
			t.Errorf("Offer(version 2) when version 2 is current = true, want false")
		}
		if p.Current() != first {
			t.Errorf("Current() after offering version 2 twice is the second, want the first to stay")
		}
	})
	t.Run("a hundred runs finish at once, in any order", func(t *testing.T) {
		r := rand.New(rand.NewPCG(1, 2))
		for round := range 50 {
			var p Publisher
			order := r.Perm(100)
			start := make(chan struct{})
			var wg sync.WaitGroup
			for _, i := range order {
				wg.Go(func() {
					<-start
					p.Offer(run(uint64(i + 1)))
				})
			}
			close(start)
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(time.Second): // never wait for an offer that keeps swapping against an old load
				t.Fatalf("round %d: offers of versions 1 to 100 still running after a second, want each to return: when the swap fails, load again", round+1)
			}
			if cur := p.Current(); cur == nil || cur.Version != 100 {
				v := uint64(0)
				if cur != nil {
					v = cur.Version
				}
				t.Fatalf("round %d: Current() after offers of versions 1 to 100 = version %d, want 100", round+1, v)
			}
		}
	})
	t.Run("Offer swaps, never only stores", func(t *testing.T) {
		if !swaps(t) {
			t.Errorf("publish.go never calls CompareAndSwap: between Offer's check and its Store, another offer can get in")
		}
	})
}
