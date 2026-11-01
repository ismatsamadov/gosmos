# The Go reference

Gosmos carries the documentation of the Go it runs with, so that looking something up never
needs the network: `reference/go.sqlite`, one SQLite file the app opens read-only. It holds the
standard library, with every exported name's declaration, doc and examples, and five articles:
The Go Programming Language Specification, The Go Memory Model, Effective Go, the FAQ and How
to Write Go Code. Like a course library, the data is the file itself, committed to the
repository; nothing else is kept beside it.

## The file

- SQLite 3, rollback journal (`journal_mode` DELETE).
- `PRAGMA application_id = 1196641618` ('GSMR'), `PRAGMA user_version = 1`. The app refuses
  any other file.
- The tables are exactly those in `db/reference.sql`, which documents every column:
  - `source`: one row, saying which Go (`go_version`, such as go1.27.3) and which platform's
    build constraints (`darwin/arm64`) the standard library comes from, the golang/website
    commit of three of the articles, when the file was made, and the credit and license the
    app shows (see below).
  - `package`: every public package of `go list std`: not internal, not vendored. Its
    synopsis is the doc's first sentence, as plain text; `deprecated` is 1 when its doc has a
    paragraph that starts with `Deprecated:`.
  - `declaration` and `symbol`: every exported const, var, func, type and method, as go doc
    shows them. A declaration is Go source without its doc comment or body, with unexported
    fields left out; a const or var block is one declaration, kept whole, that each of its
    names points to. Symbols keep go doc's order (`position`): the package's consts, vars and
    funcs, then each type followed by the consts, vars and funcs grouped with it and its
    methods. A method's `receiver` is its type's name. Methods a type gets by embedding an
    exported type are documented on that type, as in go doc.
  - `example`: the examples in the package's tests, each of the package or of one of its
    symbols. `code` is the example function's body, unindented and without its output
    comment, or the whole test file when go doc shows the whole file (`whole_file`). `output`
    is what it must print: NULL when it has no output comment and is only compiled, '' when it
    must print nothing; `unordered` is 1 for `Unordered output:`.
  - `article` and `article_section`: each article is cut at its headings. A section is a
    heading (`level`, `title` as plain text, `anchor` as on go.dev) and its body up to the next
    heading of any level; `article.lead` is what comes before the first heading.
  - `search_entry` and the full-text index `search` over it: one row for each package (found
    by its import path and synopsis), symbol (by its import path and name, such as
    `crypto/sha256.New` or `strings.Builder.Grow`, and its synopsis) and section (by its title
    and text).

### Markdown and links

Docs, examples' docs, leads and section bodies are Markdown of the kind the app renders for
lesson readmes: paragraphs, headings, fenced code, lists, block quotes, inline code,
emphasis, `^superscript^`, and the `> [!NOTE]` callout. There are no pictures and no raw HTML.

Every link leads inside the reference:

| Link | Leads to |
|---|---|
| `godoc:strings` | a package |
| `godoc:strings.Builder` | a const, var, func or type |
| `godoc:strings.Builder.Grow` | a method |
| `godoc:io.Reader.Read` | a field or interface method, which has no row of its own: its type |
| `godoc-article:spec` | an article |
| `godoc-article:spec#Struct_types` | a section of one |

The import path is everything up to the first dot after the last slash (the standard
library's paths have none). Any other link in the sources is kept as its text, since nothing
in Gosmos leads to the web.

### Search

Query `search` with FTS5 and join `search_entry` on its rowid. The tokenizer is unicode61, so
`.`, `/` and `_` separate words: quote each word of what the user typed and make the last a
prefix.

```sql
SELECT e.kind, e.title, e.package_id, e.symbol_id, e.article_slug, e.section_position
FROM search JOIN search_entry e ON e.id = search.rowid
WHERE search MATCH '"strings" "buil"*'
ORDER BY bm25(search, 10.0, 1.0)
LIMIT 50;
```

## How it is made

`make reference` remakes the file, and `make check` checks it with `go run ./cmd/reference
check reference/go.sqlite`. The tool uses no network: the Makefile fetches the three files it
needs from golang/website (and the five small programs Effective Go quotes) once, at the
commit `WEBSITE_COMMIT` names, from raw.githubusercontent.com into `$TMPDIR/godocs/website`.

```
go run ./cmd/reference build --goroot <dir> --website <dir> --website-commit <sha> --out <file>
```

- **Packages**: the tool walks `$GOROOT/src` as `go list std` does, for darwin/arm64 with cgo:
  every folder outside `cmd` and `builtin` that holds Go files for the platform, without
  internal and vendored packages. When the Go has its go command, the list must match `go
  list std`'s, or the build stops. go/parser reads each package's files and tests, and go/doc
  documents them (exported names only, examples from `_test.go` files).
- **Doc comments** become Markdown with go/doc/comment's printer: headings at level 3, code
  blocks as ```go fences, doc links such as `[strings.Builder]` as `godoc:` links. A link to a
  go.dev article or package page becomes the reference's own link; any other URL, linked or
  plain, is left as text.
- **Articles**: the spec and the memory model come from `$GOROOT/doc`, so they match the Go;
  Effective Go (`_content/doc/effective_go.html`), the FAQ (`_content/doc/faq.md`; the FAQ is
  Markdown on the website now) and How to Write Go Code (`_content/doc/code.html`) from the
  website. The tool reads each as go.dev builds it: its metadata, its template when it asks
  for one (`{{raw}}`, `{{code}}`), and goldmark for Markdown. The HTML then becomes Markdown:
  - each heading starts a section, whose anchor is the heading's id;
  - `<pre>` becomes a ```go fence, or a plain fence for the spec's grammar (`ebnf` and
    `grammar` blocks) and shell sessions;
  - lists, inline code, emphasis and strong carry over; definition lists become lists whose
    items start with the term in bold; Effective Go's note becomes a `> [!NOTE]` callout and
    the memory model's rules block quotes;
  - subscripts become Unicode subscripts where one exists (`a₁`, `pᵢ`), superscripts `^32^`;
  - links to a section of the same or another article become `godoc-article:` links (an id
    inside a section leads to that section), links to pkg.go.dev or golang.org/pkg
    `godoc:` links, and every other link its text;
  - pictures, scripts, styles and comments are dropped;
  - a table would become a GitHub table, which the app draws, but none of the five has one.
- The new file is written beside `--out`, checked, and only then put in its place.

`reference check` reports one line for each problem: a file that is not a Gosmos Go reference
of version 1; `quick_check` or `foreign_key_check` failing; a `source` other than one row; a
method of a type its package does not have; a declaration without a symbol; a package, symbol
or section that search does not find; and any link that does not lead to something in the
file.

## Credit

The reference is the work of The Go Authors. `source.credit` says so, and the app shows it
with `source.license`:

- The standard library's documentation and examples come from Go's source code: Copyright The
  Go Authors, used under Go's BSD-3-Clause license, which `source.license` holds as Go ships it.
- The spec, the memory model, Effective Go, the FAQ and How to Write Go Code are by The Go
  Authors, published on go.dev under the Creative Commons Attribution 4.0 International License
  (creativecommons.org/licenses/by/4.0), with their code under the BSD license. Gosmos changes
  them only to Markdown, without pictures, and with outside links kept as their text.
