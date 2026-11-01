package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jobloom/gosmos/db"
)

// counts is what a reference holds.
type counts struct {
	packages, symbols, examples, articles, sections int
}

func (c counts) String() string {
	return fmt.Sprintf("%s, %s, %s and %s in %s", plural(c.packages, "package"), plural(c.symbols, "symbol"),
		plural(c.examples, "example"), plural(c.articles, "article"), plural(c.sections, "section"))
}

func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

var reCommit = regexp.MustCompile(`^[0-9a-f]{40}$`)

// buildReference makes the reference anew at out: it writes a new file beside it, checks it, and only
// then puts it in out's place.
func buildReference(ctx context.Context, goroot, website, commit, out string, now time.Time) (counts, error) {
	if !reCommit.MatchString(commit) {
		return counts{}, fmt.Errorf("--website-commit %q is not a full commit hash", commit)
	}
	version, err := os.ReadFile(filepath.Join(goroot, "VERSION"))
	if err != nil {
		return counts{}, fmt.Errorf("the Go at %s: %w", goroot, err)
	}
	goVersion, _, _ := strings.Cut(strings.TrimSpace(string(version)), "\n")
	// Homebrew keeps Go's LICENSE beside its GOROOT, libexec, rather than in it.
	license, err := os.ReadFile(filepath.Join(goroot, "LICENSE"))
	if errors.Is(err, fs.ErrNotExist) {
		license, err = os.ReadFile(filepath.Join(goroot, "..", "LICENSE"))
	}
	if err != nil {
		return counts{}, fmt.Errorf("the Go at %s: %w", goroot, err)
	}

	ctxt := buildContext(goroot)
	found, err := listPackages(ctxt)
	if err != nil {
		return counts{}, err
	}
	if err := crossCheck(ctxt, found); err != nil {
		return counts{}, err
	}
	l := &linker{symbols: map[string]map[string]string{}, articles: map[string]*article{}}
	var pkgs []*stdPackage
	for _, bp := range found {
		p, err := loadPackage(bp)
		if err != nil {
			return counts{}, err
		}
		pkgs = append(pkgs, p)
		l.symbols[p.path] = p.names()
	}
	articles, err := loadArticles(goroot, website)
	if err != nil {
		return counts{}, err
	}
	for _, a := range articles {
		l.articles[a.slug] = a
	}

	tmp := out + ".new"
	if err := os.Remove(tmp); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return counts{}, err
	}
	d, err := db.CreateReference(tmp)
	if err != nil {
		return counts{}, err
	}
	c, err := write(d, l, pkgs, articles, source{goVersion, commit, now, string(license)})
	if err == nil {
		_, err = d.Exec(`INSERT INTO search(search) VALUES ('optimize')`)
	}
	if err == nil {
		_, err = d.Exec(`VACUUM`)
	}
	if err := errors.Join(err, d.Close()); err != nil {
		os.Remove(tmp)
		return counts{}, err
	}
	if err := checkFile(ctx, tmp); err != nil {
		os.Remove(tmp)
		return counts{}, err
	}
	return c, os.Rename(tmp, out)
}

// checkFile runs the reference's checks on the file at path.
func checkFile(ctx context.Context, path string) error {
	d, err := db.OpenReference(path)
	if err != nil {
		return err
	}
	defer d.Close()
	problems, err := db.ValidateReference(ctx, d)
	if err != nil {
		return err
	}
	if len(problems) > 0 {
		return errors.New("the reference failed its checks:\n  " + strings.Join(problems, "\n  "))
	}
	return nil
}

type source struct {
	goVersion, commit string
	at                time.Time
	license           string
}

// credit says whose work the reference is, and under which licenses, as the app shows it.
func (s source) credit() string {
	return fmt.Sprintf(`Gosmos's Go reference is the work of The Go Authors, converted to Markdown for Gosmos: pictures are left out, and links that lead outside the reference are kept as their text.

- The standard library's documentation and examples come from the source code of Go %s, for %s/%s. Copyright The Go Authors; used under Go's BSD-3-Clause license, below.
- The Go Programming Language Specification and The Go Memory Model come from the same Go, and Effective Go, Frequently Asked Questions (FAQ) and How to Write Go Code from the golang/website repository at commit %s. By The Go Authors, published on go.dev under the Creative Commons Attribution 4.0 International License (creativecommons.org/licenses/by/4.0); the code in them is under the BSD-3-Clause license, below.`,
		s.goVersion, goos, goarch, s.commit)
}

// write fills a new reference in one transaction.
func write(d *sql.DB, l *linker, pkgs []*stdPackage, articles []*article, src source) (counts, error) {
	var c counts
	tx, err := d.Begin()
	if err != nil {
		return c, err
	}
	defer tx.Rollback()
	var stmtErr error
	prepare := func(query string) *sql.Stmt {
		s, err := tx.Prepare(query)
		stmtErr = errors.Join(stmtErr, err)
		return s
	}
	insPackage := prepare(`INSERT INTO package (id, import_path, name, synopsis, doc, deprecated) VALUES (?, ?, ?, ?, ?, ?)`)
	insDecl := prepare(`INSERT INTO declaration (id, package_id, kind, source, doc) VALUES (?, ?, ?, ?, ?)`)
	insSymbol := prepare(`INSERT INTO symbol (id, package_id, declaration_id, kind, receiver, name, synopsis, deprecated, position) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	insExample := prepare(`INSERT INTO example (id, package_id, symbol_id, suffix, position, doc, code, output, unordered, whole_file) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	insArticle := prepare(`INSERT INTO article (slug, title, subtitle, lead, position) VALUES (?, ?, ?, ?, ?)`)
	insSection := prepare(`INSERT INTO article_section (article_slug, position, anchor, level, title, body) VALUES (?, ?, ?, ?, ?, ?)`)
	insSearch := prepare(`INSERT INTO search_entry (kind, package_id, symbol_id, article_slug, section_position, title, body) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if stmtErr != nil {
		return c, stmtErr
	}
	if _, err := tx.Exec(`INSERT INTO source (singleton, go_version, platform, website_commit, generated_at, credit, license) VALUES (1, ?, ?, ?, ?, ?, ?)`,
		src.goVersion, goos+"/"+goarch, src.commit, src.at.UTC().Format(time.RFC3339), src.credit(), src.license); err != nil {
		return c, fmt.Errorf("source: %w", err)
	}

	var declID, symbolID, exampleID int
	for i, p := range pkgs {
		pkgID := i + 1
		pd := p.doc
		synopsis := pd.Synopsis(pd.Doc)
		if _, err := insPackage.Exec(pkgID, p.path, pd.Name, null(synopsis), null(l.markdown(p, pd.Doc)), deprecated(pd.Doc)); err != nil {
			return c, fmt.Errorf("package %s: %w", p.path, err)
		}
		if _, err := insSearch.Exec("package", pkgID, nil, nil, nil, p.path, null(synopsis)); err != nil {
			return c, fmt.Errorf("package %s: %w", p.path, err)
		}
		decls, syms, examples, err := p.rows(l)
		if err != nil {
			return c, err
		}
		firstDecl, firstSymbol := declID+1, symbolID+1
		for _, x := range decls {
			declID++
			if _, err := insDecl.Exec(declID, pkgID, x.kind, x.source, null(x.doc)); err != nil {
				return c, fmt.Errorf("a declaration of %s: %w", p.path, err)
			}
		}
		for pos, s := range syms {
			symbolID++
			qualified := p.path + "." + s.name
			if s.receiver != "" {
				qualified = p.path + "." + s.receiver + "." + s.name
			}
			if _, err := insSymbol.Exec(symbolID, pkgID, firstDecl+s.decl, s.kind, null(s.receiver), s.name, null(s.synopsis), s.deprecated, pos); err != nil {
				return c, fmt.Errorf("%s: %w", qualified, err)
			}
			if _, err := insSearch.Exec("symbol", nil, symbolID, nil, nil, qualified, null(s.synopsis)); err != nil {
				return c, fmt.Errorf("%s: %w", qualified, err)
			}
		}
		for pos, x := range examples {
			exampleID++
			var symbol any
			if x.symbol >= 0 {
				symbol = firstSymbol + x.symbol
			}
			var output any
			if x.output != nil {
				output = *x.output
			}
			if _, err := insExample.Exec(exampleID, pkgID, symbol, null(x.suffix), pos, null(x.doc), x.code, output, x.unordered, x.wholeFile); err != nil {
				return c, fmt.Errorf("an example of %s: %w", p.path, err)
			}
		}
		c.packages++
		c.symbols += len(syms)
		c.examples += len(examples)
	}

	for pos, a := range articles {
		conv := &converter{link: func(href string) string { return l.target(a.base, href) }}
		if _, err := insArticle.Exec(a.slug, a.title, null(a.subtitle), null(conv.markdown(a.lead)), pos); err != nil {
			return c, fmt.Errorf("article %s: %w", a.slug, err)
		}
		for i, s := range a.sections {
			if _, err := insSection.Exec(a.slug, i, s.anchor, s.level, s.title, null(conv.markdown(s.nodes))); err != nil {
				return c, fmt.Errorf("section %s of article %s: %w", s.anchor, a.slug, err)
			}
			if _, err := insSearch.Exec("section", nil, nil, a.slug, i, s.title, null(plainText(s.nodes))); err != nil {
				return c, fmt.Errorf("section %s of article %s: %w", s.anchor, a.slug, err)
			}
		}
		c.articles++
		c.sections += len(a.sections)
	}
	return c, tx.Commit()
}

// null is s, or NULL for an empty string.
func null(s string) any {
	if s == "" {
		return nil
	}
	return s
}
