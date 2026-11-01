package main

// The Go these tests read is invented: testdata/goroot holds a tiny standard library of two
// space packages, with packages the reference must leave out beside them, and testdata/website
// the three go.dev articles in miniature.

import (
	"bytes"
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/jobloom/gosmos/db"
)

const testCommit = "0123456789abcdef0123456789abcdef01234567"

func runCmd(t *testing.T, args ...string) (code int, out, errOut string) {
	t.Helper()
	var o, e bytes.Buffer
	code = run(context.Background(), args, &o, &e)
	return code, o.String(), e.String()
}

// buildFixture builds the reference of testdata and returns its path.
func buildFixture(t *testing.T) string {
	t.Helper()
	out := filepath.Join(t.TempDir(), "reference", "go.sqlite")
	os.MkdirAll(filepath.Dir(out), 0o755)
	code, stdout, errOut := runCmd(t, "build", "--goroot", "testdata/goroot", "--website", "testdata/website",
		"--website-commit", testCommit, "--out", out)
	if code != 0 || !regexp.MustCompile(`^built .*go\.sqlite \(\d+\.\d MB\): 2 packages, 8 symbols, 5 examples and 5 articles in 8 sections\n$`).MatchString(stdout) {
		t.Fatalf("build: %d %q %q", code, stdout, errOut)
	}
	return out
}

// rows runs a query and returns its rows, a line each, columns joined by |.
func rows(t *testing.T, d *sql.DB, query string) string {
	t.Helper()
	r, err := d.Query(query)
	if err != nil {
		t.Fatalf("%s: %v", query, err)
	}
	defer r.Close()
	cols, _ := r.Columns()
	var lines []string
	for r.Next() {
		vals := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := r.Scan(ptrs...); err != nil {
			t.Fatal(err)
		}
		var cells []string
		for _, v := range vals {
			if v.Valid {
				cells = append(cells, v.String)
			} else {
				cells = append(cells, "NULL")
			}
		}
		lines = append(lines, strings.Join(cells, "|"))
	}
	return strings.Join(lines, "\n")
}

func TestBuild(t *testing.T) {
	path := buildFixture(t)
	d, err := db.OpenReference(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	for _, tc := range []struct{ name, query, want string }{
		{"the source", `SELECT go_version, platform, website_commit, generated_at GLOB '20*Z', license LIKE 'Copyright 2009 The Go Authors.%', credit LIKE '%Go go1.99.0, for darwin/arm64%' || '%' || website_commit || '%Creative Commons Attribution 4.0%' FROM source`,
			"go1.99.0|darwin/arm64|" + testCommit + "|1|1|1"},
		{"only public packages for darwin/arm64, without internal, vendor, cmd, builtin and testdata", `SELECT import_path, name, synopsis, deprecated FROM package ORDER BY id`,
			"space/orbit|orbit|Package orbit plans the passes of a craft over a ground station.|0\nspace/relay|relay|Package relay passes signals between stations.|0"},
		{"the package's doc", `SELECT doc FROM package WHERE import_path = 'space/orbit'`, "Package orbit plans the passes of a craft over a ground station.\n\n" +
			"### Passes\n\n" +
			"A pass starts when the craft rises above the horizon; see [Craft.Raise](godoc:space/orbit.Craft.Raise), [relay.Signal](godoc:space/relay.Signal) and strings.Builder. " +
			"The [Go spec](godoc-article:spec#Types) has the rules, and https://example.com/orbits explains more.\n\n" +
			"```go\nc := orbit.NewCraft(\"Wren\")\nc.Raise(2)\n```"},
		{"the symbols, in go doc's order", `SELECT p.import_path, s.position, s.kind, s.receiver, s.name, s.declaration_id, s.synopsis, s.deprecated FROM symbol s JOIN package p ON p.id = s.package_id ORDER BY s.id`,
			"space/orbit|0|const|NULL|Low|1|low orbit|0\n" +
				"space/orbit|1|const|NULL|High|1|Altitudes a craft can hold, in km.|0\n" +
				"space/orbit|2|var|NULL|Drift|2|Drift is how far a craft wanders each pass.|0\n" +
				"space/orbit|3|func|NULL|Plan|3|Plan plans a pass.|1\n" +
				"space/orbit|4|type|NULL|Craft|4|A Craft flies over the station.|0\n" +
				"space/orbit|5|func|NULL|NewCraft|5|NewCraft returns a craft called name.|0\n" +
				"space/orbit|6|method|Craft|Raise|6|Raise lifts the craft by km; see relay.Signal.Send.|0\n" +
				"space/relay|0|type|NULL|Signal|7|A Signal travels between stations.|0"},
		{"a const block, kept whole once", `SELECT kind, source, doc FROM declaration WHERE id = 1`,
			"const|const (\n\tLow  = 400 // low orbit\n\tHigh = 1200\n)|Altitudes a craft can hold, in km."},
		{"a type, its unexported field left out, its doc's links and characters", `SELECT source, doc FROM declaration WHERE id = 4`,
			"type Craft struct {\n\tName string // what the station calls it\n\t// contains filtered or unexported fields\n}|" +
				`A Craft flies over the station. Its [Craft.Name](godoc:space/orbit.Craft.Name) is shown \& logged with \~care.`},
		{"a method, linking to an interface method", `SELECT source, doc FROM declaration WHERE id = 6`,
			"func (c *Craft) Raise(km int)|Raise lifts the craft by km; see [relay.Signal.Send](godoc:space/relay.Signal.Send)."},
		{"a deprecated func", `SELECT doc FROM declaration WHERE id = 3`,
			"Plan plans a pass.\n\nDeprecated: Use [Craft.Raise](godoc:space/orbit.Craft.Raise) instead."},
		{"the examples", `SELECT e.package_id, e.position, s.name, e.suffix, e.doc, e.code, e.output, e.unordered, e.whole_file FROM example e LEFT JOIN symbol s ON s.id = e.symbol_id WHERE e.package_id = 1 ORDER BY e.position`,
			"1|0|NULL|NULL|Plan a pass.|fmt.Println(orbit.Low)|400\n|0|0\n" +
				"1|1|NewCraft|NULL|NULL|orbit.NewCraft(\"Ystmark\")|NULL|0|0\n" +
				"1|2|Raise|NULL|NULL|c := orbit.NewCraft(\"Wren\")\nc.Raise(2)\nfmt.Println(c.Name)|Wren\n|0|0\n" +
				"1|3|Raise|twice|Raise it twice; the lines come in any order.|fmt.Println(\"up\")\nfmt.Println(\"up again\")|up again\nup\n|1|0"},
		{"a whole-file example", `SELECT s.name, e.code LIKE 'package relay_test%type beacon struct{}%func ExampleSignal() {%', e.output, e.whole_file FROM example e JOIN symbol s ON s.id = e.symbol_id WHERE e.package_id = 2`,
			"Signal|1|<nil>\n|1"},
		{"the articles", `SELECT slug, title, subtitle, lead, position FROM article ORDER BY position`,
			"spec|The Go Programming Language Specification|Language version go1.99 (Nov 1, 2026)|NULL|0\n" +
				"mem|The Go Memory Model|Version of November 1, 2026|NULL|1\n" +
				"effective-go|Effective Go|NULL|> [!NOTE]\n> This document is old. See the release notes.|2\n" +
				"faq|Frequently Asked Questions (FAQ)|NULL|NULL|3\n" +
				"code|How to Write Go Code|NULL|NULL|4"},
		{"the spec's sections", `SELECT position, anchor, level, title, body FROM article_section WHERE article_slug = 'spec' ORDER BY position`,
			"0|Introduction|2|Introduction|This is the reference manual; see [the memory model](godoc-article:mem), [types](godoc-article:spec#Types) and [`Raise`](godoc:space/orbit.Craft.Raise). Programs are built from *packages* \\& files.\n" +
				"1|Types|2|Types|```\nType = TypeName [ TypeArgs ] .\n```\n\n- a₁, 2^32^\n- second\n" +
				"2|Struct_types|3|Struct types|```go\nstruct {\n\tx, y int\n}\n```\n\nSee an essay."},
		{"the memory model, its rules set apart", `SELECT body FROM article_section WHERE article_slug = 'mem'`,
			"> A send happens before the receive.\n\nRead the [spec](godoc-article:spec#Struct_types) and [`relay.Signal`](godoc:space/relay.Signal)."},
		{"Effective Go, its template run", `SELECT body FROM article_section WHERE article_slug = 'effective-go'`,
			"```go\nfunc main() {\n    fmt.Println(\"orbit\")\n}\n```\n\n- **Indentation**\n\n  Use tabs.\n\n- **Line length**\n\n  None.\n\nTemplate text such as `{{.}}` stays."},
		{"the FAQ, from Markdown", `SELECT anchor, level, title, body FROM article_section WHERE article_slug = 'faq' ORDER BY position`,
			"Origins|2|Origins|NULL\nwhy|3|Why a new language?|Because. See the [spec's types](godoc-article:spec#Types) and [Effective Go](godoc-article:effective-go#introduction).\n\n```go\ndone <- true\n```\n\n- one\n- two"},
		{"How to Write Go Code, its shell session", `SELECT body FROM article_section WHERE article_slug = 'code'`,
			"Use the [orbit package](godoc:space/orbit).\n\n```\n$ go install example/user/hello\n$ export PATH=$PATH:$(dirname $(go list -f '{{.Target}}' .))\n```"},
		{"search entries, one for each package, symbol and section", `SELECT kind, count(*) FROM search_entry GROUP BY kind ORDER BY kind`, "package|2\nsection|8\nsymbol|8"},
		{"a symbol found by its qualified name", `SELECT e.title FROM search JOIN search_entry e ON e.id = search.rowid WHERE search MATCH '"orbit" "craft" "rai"*'`, "space/orbit.Craft.Raise"},
		{"a section found by its text", `SELECT e.kind, e.title, e.body FROM search JOIN search_entry e ON e.id = search.rowid WHERE search MATCH 'essay'`, "section|Struct types|struct {\nx, y int\n}\nSee an essay."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := rows(t, d, tc.query); got != tc.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}

func TestBuildReplacesOnlyAWholeFile(t *testing.T) {
	path := buildFixture(t)
	before, _ := os.ReadFile(path)
	code, _, errOut := runCmd(t, "build", "--goroot", "testdata/goroot", "--website", t.TempDir(), "--website-commit", testCommit, "--out", path)
	if code != 1 || !strings.Contains(errOut, "effective-go: ") {
		t.Errorf("build without the website's files: %d %q", code, errOut)
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("a failed build changed the reference")
	}
	if _, err := os.Stat(path + ".new"); err == nil {
		t.Error("a failed build left its new file behind")
	}
}

func TestCheck(t *testing.T) {
	path := buildFixture(t)
	before, _ := os.ReadFile(path)
	code, out, errOut := runCmd(t, "check", path)
	want := "ok: " + path + " is the Go reference for go1.99.0 on darwin/arm64 (golang/website 0123456789ab), with 2 packages, 8 symbols, 5 examples and 5 articles in 8 sections\n"
	if code != 0 || out != want {
		t.Errorf("check: %d %q %q, want %q", code, out, errOut, want)
	}
	if after, _ := os.ReadFile(path); !bytes.Equal(before, after) {
		t.Error("check changed the file")
	}

	d, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{
		`UPDATE package SET doc = doc || ' See [maps](godoc-article:spec#Map_types).' WHERE import_path = 'space/relay'`,
		`DELETE FROM search_entry WHERE kind = 'section' AND article_slug = 'code'`,
	} {
		if _, err := d.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
	d.Close()
	code, out, _ = runCmd(t, "check", path)
	want = "search_entry holds 17 rows, want one for each package, symbol and section: 18\n" +
		"package space/relay links to godoc-article:spec#Map_types, which is not in the reference\n"
	if code != 1 || out != want {
		t.Errorf("check of a broken reference: %d %q, want %q", code, out, want)
	}

	library := filepath.Join(t.TempDir(), "library.sqlite")
	l, err := db.Create(library, db.Info{Title: "Apogee", Version: "1"})
	if err != nil {
		t.Fatal(err)
	}
	l.Close()
	if code, _, errOut := runCmd(t, "check", library); code != 1 || !strings.Contains(errOut, "is not a Gosmos Go reference") {
		t.Errorf("check of a course library: %d %q", code, errOut)
	}
}

func TestUsage(t *testing.T) {
	for _, args := range [][]string{
		nil,
		{"publish"},
		{"check"},
		{"check", "a.sqlite", "b.sqlite"},
		{"build", "--goroot", "testdata/goroot", "--out", "x.sqlite"},
		{"build", "--goroot", "testdata/goroot", "--website", "testdata/website", "--website-commit", testCommit, "--out", "x.sqlite", "extra"},
	} {
		if code, _, errOut := runCmd(t, args...); code != 2 || !strings.Contains(errOut, "usage:") {
			t.Errorf("%q: %d %q", args, code, errOut)
		}
	}
	if code, _, errOut := runCmd(t, "build", "--goroot", "testdata/goroot", "--website", "testdata/website", "--website-commit", "b8d34bb", "--out", filepath.Join(t.TempDir(), "x.sqlite")); code != 1 || !strings.Contains(errOut, "is not a full commit hash") {
		t.Errorf("a short commit: %d %q", code, errOut)
	}
}
