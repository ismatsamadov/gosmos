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
	defer trace.Stop()
	replayNight(t, "while a trace runs")
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

// traceCall returns the arguments of a call to trace.<name>, or false when call is another call.
func traceCall(call *ast.CallExpr, name string) ([]ast.Expr, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != name {
		return nil, false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return call.Args, ok && pkg.Name == "trace"
}

// isString reports whether e is the string literal s.
func isString(e ast.Expr, s string) bool {
	lit, ok := e.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	v, err := strconv.Unquote(lit.Value)
	return err == nil && v == s
}

// isIdent reports whether e is the identifier name.
func isIdent(e ast.Expr, name string) bool {
	id, ok := e.(*ast.Ident)
	return ok && name != "" && id.Name == name
}

// TestTaskInCode reads uplink.go. A task changes nothing that runPass returns, so the test looks
// for it in the code: trace.NewTask(ctx, "pass") at the start, its End deferred, the pass ID
// logged with trace.Log (or trace.Logf), and the task's context handed to fly.
func TestTaskInCode(t *testing.T) {
	fn := runPassOf(t)
	var taskCtx, task string
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || len(as.Lhs) != 2 || len(as.Rhs) != 1 {
			return true
		}
		call, ok := as.Rhs[0].(*ast.CallExpr)
		c, okc := as.Lhs[0].(*ast.Ident)
		k, okk := as.Lhs[1].(*ast.Ident)
		if !ok || !okc || !okk {
			return true
		}
		if args, ok := traceCall(call, "NewTask"); ok && len(args) == 2 && isString(args[1], "pass") {
			taskCtx, task = c.Name, k.Name
		}
		return true
	})
	if task == "" {
		t.Fatalf(`uplink.go: runPass never starts trace.NewTask(ctx, "pass"), want one task per pass`)
	}
	var ended, logged, flies bool
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.DeferStmt:
			if sel, ok := n.Call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "End" && isIdent(sel.X, task) {
				ended = true
			}
		case *ast.CallExpr:
			for _, name := range []string{"Log", "Logf"} {
				if args, ok := traceCall(n, name); ok && len(args) >= 3 && isIdent(args[0], taskCtx) && isString(args[1], "pass") {
					logged = true
				}
			}
			if sel, ok := n.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "fly" && len(n.Args) > 0 && isIdent(n.Args[0], taskCtx) {
				flies = true
			}
		}
		return true
	})
	if !ended {
		t.Errorf("uplink.go: runPass has no defer %s.End(), want the task ended with defer, however runPass returns", task)
	}
	if !logged {
		t.Errorf(`uplink.go: runPass has no trace.Log(%s, "pass", ...), want the pass ID logged on the task`, taskCtx)
	}
	if !flies {
		t.Errorf("uplink.go: runPass never calls u.fly(%s, ...), want the task's context handed to fly", taskCtx)
	}
}
