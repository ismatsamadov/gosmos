package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/doc/comment"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// The platform whose build constraints the reference follows: the Macs Gosmos runs on.
const goos, goarch = "darwin", "arm64"

// buildContext is how go/build picks the files of a package of the Go at goroot.
func buildContext(goroot string) *build.Context {
	ctxt := build.Default
	ctxt.GOROOT = goroot
	ctxt.GOPATH = ""
	ctxt.GOOS, ctxt.GOARCH = goos, goarch
	ctxt.CgoEnabled = true
	return &ctxt
}

// public reports whether a standard-library import path is one a program may import: not
// internal, not vendored.
func public(importPath string) bool {
	return !strings.HasPrefix(importPath, "vendor/") && importPath != "internal" &&
		!strings.HasPrefix(importPath, "internal/") && !strings.Contains(importPath, "/internal/") &&
		!strings.HasSuffix(importPath, "/internal")
}

// listPackages finds the public packages of the standard library under goroot as `go list std`
// does, for the reference's platform: every folder of src outside cmd, builtin, testdata and
// folders whose names start with . or _, that holds Go files for the platform.
func listPackages(ctxt *build.Context) ([]*build.Package, error) {
	src := filepath.Join(ctxt.GOROOT, "src")
	var pkgs []*build.Package
	err := filepath.WalkDir(src, func(dir string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !e.IsDir() || dir == src {
			return nil
		}
		rel, err := filepath.Rel(src, dir)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		// builtin only documents the predeclared names, for go doc; go list leaves it out too.
		if name := e.Name(); rel == "cmd" || rel == "builtin" || name == "testdata" || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			return filepath.SkipDir
		}
		if !public(rel) {
			return nil
		}
		p, err := ctxt.ImportDir(dir, 0)
		var noGo *build.NoGoError
		if errors.As(err, &noGo) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("%s: %w", rel, err)
		}
		p.ImportPath = rel
		pkgs = append(pkgs, p)
		return nil
	})
	slices.SortFunc(pkgs, func(a, b *build.Package) int { return strings.Compare(a.ImportPath, b.ImportPath) })
	return pkgs, err
}

// crossCheck compares the packages found with what the Go at goroot lists itself, when it has a
// go command: they must be the same.
func crossCheck(ctxt *build.Context, pkgs []*build.Package) error {
	goCmd := filepath.Join(ctxt.GOROOT, "bin", "go")
	if _, err := os.Stat(goCmd); err != nil {
		return nil // a GOROOT without its tools, such as the tests'
	}
	cmd := exec.Command(goCmd, "list", "std")
	cmd.Dir = ctxt.GOROOT
	cmd.Env = append(os.Environ(), "GOROOT="+ctxt.GOROOT, "GOOS="+goos, "GOARCH="+goarch, "CGO_ENABLED=1", "GOTOOLCHAIN=local", "GOFLAGS=")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("go list std: %w", err)
	}
	var want, got []string
	for p := range strings.FieldsSeq(string(out)) {
		if public(p) {
			want = append(want, p)
		}
	}
	for _, p := range pkgs {
		got = append(got, p.ImportPath)
	}
	slices.Sort(want)
	if !slices.Equal(got, want) {
		return fmt.Errorf("found %d packages, but go list std lists %d public ones:\n  found  %s\n  listed %s",
			len(got), len(want), strings.Join(got, " "), strings.Join(want, " "))
	}
	return nil
}

// A stdPackage is a package's documentation as go/doc reads it, with its examples.
type stdPackage struct {
	path string
	fset *token.FileSet
	doc  *doc.Package
}

// loadPackage parses a package's files for the platform, and its tests for their examples.
func loadPackage(p *build.Package) (*stdPackage, error) {
	fset := token.NewFileSet()
	var files []*ast.File
	for _, name := range slices.Concat(p.GoFiles, p.CgoFiles, p.TestGoFiles, p.XTestGoFiles) {
		f, err := parser.ParseFile(fset, filepath.Join(p.Dir, name), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	d, err := doc.NewFromFiles(fset, files, p.ImportPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", p.ImportPath, err)
	}
	return &stdPackage{path: p.ImportPath, fset: fset, doc: d}, nil
}

// The rows a package becomes.
type (
	declRow struct {
		kind, source, doc string
	}
	symbolRow struct {
		kind, receiver, name, synopsis string
		decl                           int // index into the package's declarations
		deprecated                     bool
	}
	exampleRow struct {
		symbol               int // index into the package's symbols; -1 for the package
		suffix, doc, code    string
		output               *string
		unordered, wholeFile bool
	}
)

// each calls fn for every symbol of the package in go doc's order, with its declaration.
func (p *stdPackage) each(fn func(kind, recv, name string, decl ast.Node, docText string, spec *ast.ValueSpec, examples []*doc.Example)) {
	values := func(vs []*doc.Value, kind string) {
		for _, v := range vs {
			for _, s := range v.Decl.Specs {
				spec := s.(*ast.ValueSpec)
				for _, n := range spec.Names {
					if n.IsExported() {
						fn(kind, "", n.Name, v.Decl, v.Doc, spec, nil)
					}
				}
			}
		}
	}
	funcs := func(fs []*doc.Func, kind, recv string) {
		for _, f := range fs {
			fn(kind, recv, f.Name, f.Decl, f.Doc, nil, f.Examples)
		}
	}
	values(p.doc.Consts, "const")
	values(p.doc.Vars, "var")
	funcs(p.doc.Funcs, "func", "")
	for _, t := range p.doc.Types {
		fn("type", "", t.Name, t.Decl, t.Doc, nil, t.Examples)
		values(t.Consts, "const")
		values(t.Vars, "var")
		funcs(t.Funcs, "func", "")
		funcs(t.Methods, "method", t.Name)
	}
}

// names are the package's symbols, Name or Type.Name, with their kinds.
func (p *stdPackage) names() map[string]string {
	m := map[string]string{}
	p.each(func(kind, recv, name string, _ ast.Node, _ string, _ *ast.ValueSpec, _ []*doc.Example) {
		if recv != "" {
			name = recv + "." + name
		}
		m[name] = kind
	})
	return m
}

// rows turns the package into its declarations, symbols and examples.
func (p *stdPackage) rows(l *linker) (decls []declRow, syms []symbolRow, examples []exampleRow, err error) {
	declIndex := map[ast.Node]int{}
	for _, ex := range p.doc.Examples {
		examples = append(examples, p.example(l, ex, -1))
	}
	p.each(func(kind, recv, name string, decl ast.Node, docText string, spec *ast.ValueSpec, exs []*doc.Example) {
		if err != nil {
			return
		}
		i, ok := declIndex[decl]
		if !ok {
			var b bytes.Buffer
			if err = format.Node(&b, p.fset, decl); err != nil {
				err = fmt.Errorf("%s.%s: %w", p.path, name, err)
				return
			}
			i = len(decls)
			declIndex[decl] = i
			decls = append(decls, declRow{kind: kind, source: b.String(), doc: l.markdown(p, docText)})
		}
		own := docText // a const or var of a block is summed up by its own comment, if it has one
		if spec != nil {
			if spec.Doc != nil {
				own = spec.Doc.Text()
			} else if spec.Comment != nil {
				own = spec.Comment.Text()
			}
		}
		syms = append(syms, symbolRow{kind: kind, receiver: recv, name: name, decl: i,
			synopsis: p.doc.Synopsis(own), deprecated: deprecated(docText) || deprecated(own)})
		for _, ex := range exs {
			examples = append(examples, p.example(l, ex, len(syms)-1))
		}
	})
	return decls, syms, examples, err
}

// reOutput finds the output comment that ends an example: its code stops there.
var reOutput = regexp.MustCompile(`(?i)//[[:space:]]*(unordered )?output:`)

// example is one example: a function's body, unindented, as go doc shows it, or the whole file
// when go/doc took it whole.
func (p *stdPackage) example(l *linker, ex *doc.Example, symbol int) exampleRow {
	var b bytes.Buffer
	format.Node(&b, p.fset, &printer.CommentedNode{Node: ex.Code, Comments: ex.Comments})
	code := b.String()
	_, whole := ex.Code.(*ast.File)
	if !whole {
		code = strings.TrimSuffix(strings.TrimPrefix(code, "{"), "}")
		lines := strings.Split(code, "\n")
		for i, line := range lines {
			lines[i] = strings.TrimPrefix(line, "\t")
		}
		code = strings.Join(lines, "\n")
		if loc := reOutput.FindStringIndex(code); loc != nil {
			code = code[:loc[0]]
		}
	}
	r := exampleRow{symbol: symbol, suffix: ex.Suffix, doc: l.markdown(p, ex.Doc),
		code: strings.TrimRight(strings.Trim(code, "\n"), " \t\n"), unordered: ex.Unordered, wholeFile: whole}
	if ex.Output != "" || ex.EmptyOutput {
		r.output = &ex.Output
	}
	return r
}

// deprecated reports whether a doc comment has a paragraph that starts with "Deprecated:".
func deprecated(text string) bool {
	for p := range strings.SplitSeq(text, "\n\n") {
		if strings.HasPrefix(strings.TrimSpace(p), "Deprecated:") {
			return true
		}
	}
	return false
}

// markdown writes a doc comment of the package as Markdown: headings at level 3, code blocks as
// ```go fences, doc links ([strings.Builder]) as godoc: links, links to go.dev's articles and
// package pages as the reference's own, and every other link as its text.
func (l *linker) markdown(p *stdPackage, text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	d := p.doc.Parser().Parse(text)
	pr := &comment.Printer{
		HeadingLevel: 3,
		HeadingID:    func(*comment.Heading) string { return "" }, // the app's Markdown has no {#id}
		DocLinkURL: func(link *comment.DocLink) string {
			path := link.ImportPath
			if path == "" {
				path = p.path
			}
			return l.symbolURL(path, link.Recv, link.Name)
		},
	}
	var blocks []string
	for _, blk := range d.Content {
		switch x := blk.(type) {
		case *comment.Code:
			blocks = append(blocks, fence("go", strings.TrimRight(x.Text, "\n")))
			continue
		case *comment.Paragraph:
			x.Text = l.texts(p, x.Text)
		case *comment.Heading:
			x.Text = l.texts(p, x.Text)
		case *comment.List:
			for _, item := range x.Items {
				for _, c := range item.Content {
					if para, ok := c.(*comment.Paragraph); ok {
						para.Text = l.texts(p, para.Text)
					}
				}
			}
		}
		md := strings.TrimRight(string(pr.Markdown(&comment.Doc{Content: []comment.Block{blk}})), "\n")
		md = mdLeftovers.Replace(md) // what comment's printer leaves unescaped: entities, strikethrough
		if strings.HasPrefix(md, ">") {
			md = `\` + md
		}
		blocks = append(blocks, md)
	}
	return strings.Join(blocks, "\n\n")
}

var mdLeftovers = strings.NewReplacer("&", `\&`, "~", `\~`)

// texts rewrites the links in a doc comment's text: a doc link to a package the reference does
// not have, and any URL that does not lead into the reference, become their text.
func (l *linker) texts(p *stdPackage, in []comment.Text) []comment.Text {
	var out []comment.Text
	for _, t := range in {
		switch t := t.(type) {
		case *comment.Link:
			if u := l.target(goDev, t.URL); u != "" {
				out = append(out, &comment.Link{Text: t.Text, URL: u})
			} else {
				out = append(out, t.Text...)
			}
		case *comment.DocLink:
			path := t.ImportPath
			if path == "" {
				path = p.path
			}
			if l.symbols[path] == nil {
				out = append(out, t.Text...)
			} else {
				out = append(out, t)
			}
		default:
			out = append(out, t)
		}
	}
	return out
}
