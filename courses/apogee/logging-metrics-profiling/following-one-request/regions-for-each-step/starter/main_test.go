package main

import (
	"bytes"
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"runtime/trace"
	"strconv"
	"testing"
	"time"
)

// radio is a fake radio: a satellite in answering acknowledges at once, and any other stays
// silent until the pass's context ends.
type radio struct{ answering map[string]bool }

// fly flies pass p as uplinkd's radio does: it returns at the acknowledgment, or with ErrNoAck
// when ctx ends first.
func (r radio) fly(ctx context.Context, p PassView, lines []string, buf []byte) (Report, error) {
	if r.answering[p.ID] {
		return Report{Frames: buf}, nil
	}
	<-ctx.Done()
	return Report{Frames: buf}, ErrNoAck
}

// replayNight runs three passes of the bench replay, each on an uplinkd of its own, and checks
// what runPass returns: Cirrine-1 never answers, Furrow-B finds all four workers held, and
// Furrow-A flies. Each pass's LOS comes after milliseconds here, not minutes: long enough that
// a free worker is always taken before it, even on a busy Mac.
func replayNight(t *testing.T, while string) {
	t.Helper()
	r := radio{answering: map[string]bool{"furrow-a/LAN/22382": true}}
	for _, p := range []struct {
		id   string
		held int // workers already taken when the pass starts
		los  time.Duration
		want string // "" for a flown pass
	}{
		{"cirrine-1/LAN/17759", 0, 200 * time.Millisecond, "pass cirrine-1/LAN/17759: no ack before LOS"},
		{"furrow-b/LAN/22379", 4, 20 * time.Millisecond, "pass furrow-b/LAN/22379: no free worker before LOS"},
		{"furrow-a/LAN/22382", 0, time.Second, ""},
	} {
		u := newUplink("LAN", r.fly)
		for range p.held {
			u.workers <- struct{}{}
		}
		ctx, cancel := context.WithTimeout(context.Background(), p.los)
		err := u.runPass(ctx, PassView{ID: p.id}, []string{"PING"})
		cancel()
		switch {
		case p.want == "" && err != nil:
			t.Errorf("runPass(%s) %s = %v, want nil", p.id, while, err)
		case p.want != "" && (err == nil || err.Error() != p.want):
			t.Errorf("runPass(%s) %s = %v, want %q", p.id, while, err, p.want)
		}
	}
}

func TestSameResults(t *testing.T) {
	replayNight(t, "with no trace running")
	var buf bytes.Buffer
	if err := trace.Start(&buf); err != nil {
		t.Fatalf("trace.Start = %v, want nil", err)
	}
	replayNight(t, "while a trace runs")
	trace.Stop()
	if head := buf.Bytes()[:min(5, buf.Len())]; !bytes.Equal(head, []byte("go 1.")) {
		t.Errorf("the trace written during the replay starts with %q, want \"go 1.\"", head)
	}
}

// runPassOf parses uplink.go and returns its runPass.
func runPassOf(t *testing.T) *ast.FuncDecl {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), "uplink.go", nil, 0)
	if err != nil {
		t.Fatalf("reading uplink.go: %v", err)
	}
	for _, d := range file.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Name.Name == "runPass" {
			return fn
		}
	}
	t.Fatalf("uplink.go has no runPass")
	return nil
}

// regionCall returns the name a call to trace.StartRegion or trace.WithRegion gives its region.
func regionCall(n ast.Node) (fn, name string, ok bool) {
	call, isCall := n.(*ast.CallExpr)
	if !isCall || len(call.Args) < 2 {
		return "", "", false
	}
	sel, isSel := call.Fun.(*ast.SelectorExpr)
	if !isSel || (sel.Sel.Name != "StartRegion" && sel.Sel.Name != "WithRegion") {
		return "", "", false
	}
	if pkg, isPkg := sel.X.(*ast.Ident); !isPkg || pkg.Name != "trace" {
		return "", "", false
	}
	lit, isLit := call.Args[1].(*ast.BasicLit)
	if !isLit || lit.Kind != token.STRING {
		return sel.Sel.Name, "", true
	}
	name, _ = strconv.Unquote(lit.Value)
	return sel.Sel.Name, name, true
}

// TestRegionsInCode reads uplink.go. Regions change nothing that runPass returns, so the test
// looks for them in the code: one named "wait for worker" and one named "fly", and an End for
// every region StartRegion returns, called in runPass itself rather than in another goroutine.
func TestRegionsInCode(t *testing.T) {
	fn := runPassOf(t)
	names := map[string]bool{}
	var starts []*ast.CallExpr        // every call to trace.StartRegion
	assigned := map[ast.Node]string{} // a StartRegion call: the variable it is assigned to
	chained := map[ast.Node]bool{}    // a StartRegion call ended on the spot: StartRegion(...).End()
	ended := map[string]bool{}        // variables whose End runPass calls itself
	var goroutines []ast.Node
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if fnName, name, ok := regionCall(n); ok {
			names[name] = true
			if fnName == "StartRegion" {
				starts = append(starts, n.(*ast.CallExpr))
			}
		}
		switch n := n.(type) {
		case *ast.GoStmt:
			goroutines = append(goroutines, n)
		case *ast.AssignStmt:
			if id, ok := n.Lhs[0].(*ast.Ident); ok && len(n.Lhs) == 1 && len(n.Rhs) == 1 {
				assigned[n.Rhs[0]] = id.Name
			}
		case *ast.CallExpr:
			sel, ok := n.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "End" {
				return true
			}
			for _, g := range goroutines {
				if g.Pos() <= n.Pos() && n.End() <= g.End() {
					return true // an End in another goroutine doesn't count
				}
			}
			if id, ok := sel.X.(*ast.Ident); ok {
				ended[id.Name] = true
			} else {
				chained[sel.X] = true
			}
		}
		return true
	})
	for _, want := range []string{"wait for worker", "fly"} {
		if !names[want] {
			t.Errorf("uplink.go: runPass has no region named %q, want one", want)
		}
	}
	for _, call := range starts {
		_, name, _ := regionCall(call)
		if v, ok := assigned[call]; chained[call] || (ok && ended[v]) {
			continue
		} else if ok {
			t.Errorf("uplink.go: runPass never calls %s.End() for its region %q, want every StartRegion ended in runPass", v, name)
		} else {
			t.Errorf("uplink.go: runPass never ends its region %q, want every StartRegion ended in runPass", name)
		}
	}
}
