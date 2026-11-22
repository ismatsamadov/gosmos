package main

// drift lists what changes in the code a course gives the learner. The read-only and hidden
// files are the team's codebase as the learner reads it, so a declaration in them changes only
// where a lesson's task changes it, and then the new text appears first in that lesson's
// solution (docs/authoring.md, "Given code changes only for a reason"). drift follows every
// top-level declaration of a given Go file through a course, lesson by lesson, and reports each
// lesson that gives it differently from the last text the learner saw: its last given copy, or
// the solution that changed it. Tests and each lesson's func main are left out.
//
// A change a lesson or the story asks for is fine: the rules allow it, but drift can't know it.
// An accept list names each such change, with its reason, so it stops being reported; an entry
// that no longer matches a change is reported instead, so the list never outlives its changes.

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"maps"
	"os"
	osexec "os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ismatsamadov/gosmos/db"
)

// shownDecl is one top-level declaration as one lesson shows it.
type shownDecl struct {
	lesson int    // the lesson's index in course order
	file   string // the file it is in
	given  bool   // in a read-only or hidden starter file, rather than in the solution
	text   string // as written, its doc comment included
	code   string // printed without comments, so wording and layout don't count
}

// driftChange is one lesson giving a declaration differently from the last text the learner saw.
type driftChange struct {
	lesson, since int
	file, name    string
	codeChanged   bool
	goesBack      bool // the text is one the learner saw before the last
	before, after string
}

// acceptedChange is an entry of an accept list: a change, as drift names it, and why it is fine.
type acceptedChange struct {
	change, reason string
	line           int
	used           bool
}

// readAccepted reads an accept list: for each change, a line that names it as drift does,
// "<course>/<chapter>/<lesson>: <file> gives <declaration>", and an indented line or more
// saying why. Blank lines and lines starting with # are left out.
func readAccepted(path string) ([]*acceptedChange, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var list []*acceptedChange
	for i, line := range strings.Split(string(b), "\n") {
		switch {
		case strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#"):
		case strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t"):
			if len(list) == 0 {
				return nil, fmt.Errorf("%s:%d: a reason comes after the change it is for", path, i+1)
			}
			last := list[len(list)-1]
			last.reason = strings.TrimSpace(last.reason + " " + strings.TrimSpace(line))
		case !strings.Contains(line, " gives "):
			return nil, fmt.Errorf("%s:%d: %q names no change: <lesson>: <file> gives <declaration>", path, i+1, line)
		default:
			list = append(list, &acceptedChange{change: line, line: i + 1})
		}
	}
	for _, a := range list {
		if a.reason == "" {
			return nil, fmt.Errorf("%s:%d: %s has no reason: an indented line after it says why the change is fine", path, a.line, a.change)
		}
	}
	return list, nil
}

func drift(ctx context.Context, path, course, acceptPath string, diffs bool, out io.Writer) (bool, error) {
	var accepted []*acceptedChange
	if acceptPath != "" {
		var err error
		if accepted, err = readAccepted(acceptPath); err != nil {
			return false, err
		}
	}
	d, recovered, err := db.Inspect(path)
	if err != nil {
		return false, err
	}
	defer d.Close()
	if recovered {
		fmt.Fprintln(out, db.RecoveredMessage(path))
	}
	var version int
	if err := d.QueryRowContext(ctx, `PRAGMA user_version`).Scan(&version); err != nil {
		return false, fmt.Errorf("%s: %w", path, err)
	}
	if version != db.Version {
		return false, fmt.Errorf("%s is of format %d, not %d; make it again with this version's tools", path, version, db.Version)
	}
	ls, err := loadCodeLessons(ctx, d, course)
	if err != nil {
		return false, err
	}
	var changes []driftChange
	for _, c := range driftChanges(ls) {
		key := fmt.Sprintf("%s: %s gives %s", ls[c.lesson].name, c.file, c.name)
		if i := slices.IndexFunc(accepted, func(a *acceptedChange) bool { return a.change == key }); i >= 0 {
			accepted[i].used = true
			continue
		}
		changes = append(changes, c)
	}
	var stale []*acceptedChange
	for _, a := range accepted {
		// A course left out was not looked at.
		if lesson, _, _ := strings.Cut(a.change, ": "); !a.used && (course == "" || strings.HasPrefix(lesson, course+"/")) {
			stale = append(stale, a)
		}
	}
	for _, a := range stale {
		fmt.Fprintf(out, "%s:%d: %s: no such change; take it out of the list\n", acceptPath, a.line, a.change)
	}
	for _, c := range changes {
		what := "only its comments or layout changed"
		if c.codeChanged {
			what = "its code changed"
		}
		if c.goesBack {
			what += ", back to an older text,"
		}
		fmt.Fprintf(out, "%s: %s gives %s: %s since %s\n", ls[c.lesson].name, c.file, c.name, what, ls[c.since].name)
		if diffs {
			showDiff(out, c.before, c.after)
		}
	}
	if len(changes) > 0 || len(stale) > 0 {
		fmt.Fprintf(out, "%s in the given code of %s\n", plural(len(changes), "change"), path)
		return false, nil
	}
	used := 0
	for _, a := range accepted {
		if a.used {
			used++
		}
	}
	fmt.Fprintf(out, "ok: %s: the given code changes only where a solution changed it first (%s)", path, plural(len(ls), "code lesson"))
	if used > 0 {
		fmt.Fprintf(out, ", or where %s accepts it (%s)", acceptPath, plural(used, "change"))
	}
	fmt.Fprintln(out)
	return true, nil
}

// driftChanges follows each declaration through its course, in lesson order: within a lesson
// its solution comes first, so a task that changes a given declaration and then gives it
// changed is fine.
func driftChanges(ls []*codeLesson) []driftChange {
	shown := map[string][]shownDecl{} // by course, then the declaration's name
	for i, l := range ls {
		courseSlug, _, _ := strings.Cut(l.name, "/")
		for _, name := range slices.Sorted(maps.Keys(l.solution)) {
			for decl, s := range declarations(name, l.solution[name]) {
				s.lesson = i
				shown[courseSlug+"\x00"+decl] = append(shown[courseSlug+"\x00"+decl], s)
			}
		}
		for _, f := range l.starter {
			if f.editable() {
				continue
			}
			for decl, s := range declarations(f.name, f.content) {
				s.lesson, s.given = i, true
				shown[courseSlug+"\x00"+decl] = append(shown[courseSlug+"\x00"+decl], s)
			}
		}
	}
	var changes []driftChange
	for key, uses := range shown {
		_, name, _ := strings.Cut(key, "\x00")
		seen := map[string]bool{}
		for i, u := range uses {
			if i > 0 && u.given && u.text != uses[i-1].text {
				last := uses[i-1]
				changes = append(changes, driftChange{lesson: u.lesson, since: last.lesson, file: u.file, name: name,
					codeChanged: u.code != last.code, goesBack: seen[u.text], before: last.text, after: u.text})
			}
			seen[u.text] = true
		}
	}
	slices.SortFunc(changes, func(a, b driftChange) int {
		if a.lesson != b.lesson {
			return a.lesson - b.lesson
		}
		return strings.Compare(a.name, b.name)
	})
	return changes
}

// declarations returns the top-level declarations of a Go file, by a name such as
// "func (Clock).Now", "type Pass" or "const a, b", leaving out tests, func main and imports.
// Each name of a group counts on its own, so moving a declaration into or out of a group, or
// rewording one member, still shows. A file that doesn't parse has none.
func declarations(name, src string) map[string]shownDecl {
	if filepath.Ext(name) != ".go" || strings.HasSuffix(name, "_test.go") {
		return nil
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, src, parser.ParseComments)
	if err != nil {
		return nil
	}
	decls := map[string]shownDecl{}
	add := func(key string, doc *ast.CommentGroup, node ast.Node, end token.Pos, bare ast.Decl) {
		start := node.Pos()
		if doc != nil {
			start = doc.Pos()
		}
		text := src[fset.Position(start).Offset:fset.Position(end).Offset]
		var code bytes.Buffer
		if err := printer.Fprint(&code, token.NewFileSet(), withoutComments(bare)); err != nil {
			return
		}
		decls[key] = shownDecl{file: filepath.Base(name), text: text, code: code.String()}
	}
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil && d.Name.Name == "main" {
				continue // each lesson's main is its own program
			}
			key := "func " + d.Name.Name
			if d.Recv != nil && len(d.Recv.List) > 0 {
				t := d.Recv.List[0].Type
				if s, ok := t.(*ast.StarExpr); ok {
					t = s.X
				}
				if id, ok := t.(*ast.Ident); ok {
					key = "func (" + id.Name + ")." + d.Name.Name
				}
			}
			add(key, d.Doc, d, d.End(), d)
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				continue
			}
			for _, sp := range d.Specs {
				var doc *ast.CommentGroup
				end := sp.End()
				var names []string
				switch sp := sp.(type) {
				case *ast.TypeSpec:
					doc, names = sp.Doc, []string{sp.Name.Name}
					if sp.Comment != nil && sp.Comment.End() > end {
						end = sp.Comment.End()
					}
				case *ast.ValueSpec:
					doc = sp.Doc
					for _, n := range sp.Names {
						names = append(names, n.Name)
					}
					if sp.Comment != nil && sp.Comment.End() > end {
						end = sp.Comment.End()
					}
				}
				bare := &ast.GenDecl{Tok: d.Tok, Specs: []ast.Spec{sp}}
				if doc == nil && !d.Lparen.IsValid() {
					doc, bare.Doc = d.Doc, d.Doc
				}
				add(d.Tok.String()+" "+strings.Join(names, ", "), doc, sp, end, bare)
			}
		}
	}
	return decls
}

// withoutComments returns d with its doc comments and line comments taken out, so printing it
// shows only its code; a directive such as //go:embed stays, since it is code. It changes d
// itself; declarations reads each text before it calls this.
func withoutComments(d ast.Decl) ast.Decl {
	ast.Inspect(d, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncDecl:
			n.Doc = directives(n.Doc)
		case *ast.GenDecl:
			n.Doc = directives(n.Doc)
		case *ast.Field:
			n.Doc, n.Comment = nil, nil
		case *ast.ValueSpec:
			n.Doc, n.Comment = directives(n.Doc), nil
		case *ast.TypeSpec:
			n.Doc, n.Comment = directives(n.Doc), nil
		}
		return true
	})
	return d
}

// directives returns the //go: lines of a doc comment, or nil if it has none.
func directives(g *ast.CommentGroup) *ast.CommentGroup {
	if g == nil {
		return nil
	}
	var kept []*ast.Comment
	for _, c := range g.List {
		if strings.HasPrefix(c.Text, "//go:") {
			kept = append(kept, c)
		}
	}
	if kept == nil {
		return nil
	}
	return &ast.CommentGroup{List: kept}
}

// showDiff prints a unified diff from before to after, indented, with diff(1).
func showDiff(out io.Writer, before, after string) {
	dir, err := os.MkdirTemp("", "gosmos-drift-")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)
	a, b := filepath.Join(dir, "before"), filepath.Join(dir, "after")
	if os.WriteFile(a, []byte(before+"\n"), 0o600) != nil || os.WriteFile(b, []byte(after+"\n"), 0o600) != nil {
		return
	}
	text, _ := osexec.Command("diff", "-u", a, b).Output() // diff exits 1 when the files differ
	for i, line := range strings.Split(string(text), "\n") {
		if i >= 2 && line != "" { // the two header lines name temporary files
			fmt.Fprintln(out, "        "+line)
		}
	}
}
