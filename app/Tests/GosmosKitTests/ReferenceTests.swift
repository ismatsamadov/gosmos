import Foundation
import GRDB
import Testing
@testable import GosmosKit

/// A Go reference made from the repository's db/reference.sql in $TMPDIR, with a few packages,
/// names, examples and articles written for these tests in the shape of the real ones.
enum ReferenceFixture {
    static let schemaURL = Fixtures.repository.appendingPathComponent("db/reference.sql")
    /// The full reference the app ships, for the time budgets.
    static let fullURL = Fixtures.repository.appendingPathComponent("reference/go.sqlite")

    static func database() throws -> (ReferenceDatabase, URL) {
        let dir = try Fixtures.tempDir("reference")
        let url = dir.appendingPathComponent("go.sqlite")
        let queue = try DatabaseQueue(path: url.path)
        try queue.writeWithoutTransaction { db in
            try db.execute(sql: String(contentsOf: schemaURL, encoding: .utf8))
            try db.execute(sql: "PRAGMA journal_mode = DELETE")
        }
        try queue.write(fill)
        try queue.close()
        return (try ReferenceDatabase(url: url), dir)
    }

    /// Packages with their names in go doc's order: (kind, receiver, name, deprecated, source,
    /// doc), a block's names sharing the first one's declaration when `block` names it.
    static func fill(_ db: Database) throws {
        try db.execute(sql: """
            INSERT INTO source(singleton, go_version, platform, website_commit, generated_at, credit, license)
            VALUES (1, 'go1.27.3', 'darwin/arm64', ?, '2026-11-01T00:00:00Z', ?, ?)
            """, arguments: [String(repeating: "ab", count: 20),
                             "The Go documentation is by The Go Authors, under CC BY 4.0 (https://creativecommons.org/licenses/by/4.0/).",
                             "Copyright 2009 The Go Authors.\n\nRedistribution and use in source and binary forms…"])
        var packageID = 0, declarationID = 0, symbolID = 0, entryID = 0
        func entry(_ kind: String, package: Int? = nil, symbol: Int? = nil, slug: String? = nil, position: Int? = nil, title: String, body: String?) throws {
            entryID += 1
            try db.execute(sql: """
                INSERT INTO search_entry(id, kind, package_id, symbol_id, article_slug, section_position, title, body) VALUES (?,?,?,?,?,?,?,?)
                """, arguments: [entryID, kind, package, symbol, slug, position, title, body])
        }
        func package(_ path: String, _ name: String, doc: String, deprecated: Bool = false,
                     _ symbols: [(kind: String, receiver: String?, name: String, deprecated: Bool, source: String, doc: String?)],
                     blocks: [String: String] = [:], examples: [(symbol: String?, suffix: String?, code: String, output: String?)] = []) throws {
            packageID += 1
            let synopsis = String(doc.prefix { $0 != "." }) + "."
            try db.execute(sql: "INSERT INTO package(id, import_path, name, synopsis, doc, deprecated) VALUES (?,?,?,?,?,?)",
                           arguments: [packageID, path, name, synopsis, doc, deprecated ? 1 : 0])
            try entry("package", package: packageID, title: path, body: synopsis)
            var declarations: [String: Int] = [:], ids: [String: Int] = [:]
            for (position, s) in symbols.enumerated() {
                let anchor = s.receiver.map { $0 + "." + s.name } ?? s.name
                let declaration: Int
                if let first = blocks[s.name], let shared = declarations[first] {
                    declaration = shared
                } else {
                    declarationID += 1
                    declaration = declarationID
                    try db.execute(sql: "INSERT INTO declaration(id, package_id, kind, source, doc) VALUES (?,?,?,?,?)",
                                   arguments: [declaration, packageID, s.kind, s.source, s.doc])
                }
                declarations[s.name] = declaration
                symbolID += 1
                ids[anchor] = symbolID
                try db.execute(sql: """
                    INSERT INTO symbol(id, package_id, declaration_id, kind, receiver, name, synopsis, deprecated, position) VALUES (?,?,?,?,?,?,?,?,?)
                    """, arguments: [symbolID, packageID, declaration, s.kind, s.receiver, s.name, s.doc.map { String($0.prefix { $0 != "." }) + "." },
                                     s.deprecated ? 1 : 0, position])
                try entry("symbol", symbol: symbolID, title: path + "." + anchor, body: s.doc)
            }
            for (position, e) in examples.enumerated() {
                try db.execute(sql: """
                    INSERT INTO example(package_id, symbol_id, suffix, position, doc, code, output, unordered, whole_file) VALUES (?,?,?,?,?,?,?,0,0)
                    """, arguments: [packageID, e.symbol.flatMap { ids[$0] }, e.suffix, position, nil, e.code, e.output])
            }
        }

        try package("fmt", "fmt", doc: "Package fmt implements formatted I/O with functions analogous to C's printf and scanf.", [
            ("func", nil, "Printf", false, "func Printf(format string, a ...any) (n int, err error)", "Printf formats according to a format specifier and writes to standard output."),
            ("func", nil, "Println", false, "func Println(a ...any) (n int, err error)",
             "Println formats using the default formats for its operands and writes to standard output. It returns the number of bytes written and any write error encountered.\n\nSpaces are always added between operands, as by [io.Writer]."),
        ], examples: [("Println", nil, "fmt.Println(\"Kim is\", 22, \"years old.\")", "Kim is 22 years old.\n")])
        try package("log", "log", doc: "Package log implements a simple logging package.", [
            ("func", nil, "Println", false, "func Println(v ...any)", "Println calls Output to print to the standard logger."),
            ("type", nil, "Logger", false, "type Logger struct {\n\t// contains filtered or unexported fields\n}", "A Logger represents an active logging object."),
            ("method", "Logger", "Println", false, "func (l *Logger) Println(v ...any)", "Println calls l.Output to print to the logger."),
        ])
        try package("archive/builder", "builder", doc: "Package builder builds archives, one Builder at a time; a builder is not a strings builder.", [
            ("type", nil, "Builder", false, "type Builder struct{}", "A Builder builds an archive."),
        ])
        try package("strings", "strings", doc: "Package strings implements simple functions to manipulate UTF-8 encoded strings.\n\nFor more, see https://blog.golang.org/strings.", [
            ("func", nil, "Clone", false, "func Clone(s string) string", "Clone returns a fresh copy of s."),
            ("type", nil, "Builder", false, "type Builder struct {\n\t// contains filtered or unexported fields\n}",
             "A Builder is used to efficiently build a string using [Builder.Write] methods. It minimizes memory copying."),
            ("method", "Builder", "Grow", false, "func (b *Builder) Grow(n int)", "Grow grows b's capacity, if necessary, to guarantee space for another n bytes."),
            ("method", "Builder", "WriteString", false, "func (b *Builder) WriteString(s string) (int, error)", "WriteString appends the contents of s to b's buffer."),
            ("type", nil, "Reader", false, "type Reader struct {\n\t// contains filtered or unexported fields\n}", "A Reader implements the [io.Reader] interface by reading from a string."),
            ("func", nil, "NewReader", false, "func NewReader(s string) *Reader", "NewReader returns a new [Reader] reading from s."),
            ("method", "Reader", "Read", false, "func (r *Reader) Read(b []byte) (n int, err error)", "Read implements the [io.Reader] interface."),
        ], examples: [("Builder", nil, "var b strings.Builder\nb.WriteString(\"ignition\")\nfmt.Println(b.String())", "ignition\n"),
                      (nil, "fields", "fmt.Println(strings.Fields(\"  a b  \"))", nil)])
        try package("io", "io", doc: "Package io provides basic interfaces to I/O primitives.", [
            ("const", nil, "SeekStart", false, "const (\n\tSeekStart   = 0 // seek relative to the origin of the file\n\tSeekCurrent = 1 // seek relative to the current offset\n)",
             "Seek whence values."),
            ("const", nil, "SeekCurrent", false, "", nil),
            ("type", nil, "Reader", false, "type Reader interface {\n\tRead(p []byte) (n int, err error)\n}", "Reader is the interface that wraps the basic Read method."),
            ("type", nil, "Writer", false, "type Writer interface {\n\tWrite(p []byte) (n int, err error)\n}", "Writer is the interface that wraps the basic Write method."),
        ], blocks: ["SeekStart": "SeekStart", "SeekCurrent": "SeekStart"])
        try package("crypto/rand", "rand", doc: "Package rand implements a cryptographically secure random number generator.", [
            ("func", nil, "Read", false, "func Read(b []byte) (n int, err error)", "Read fills b with cryptographically secure random bytes."),
        ])
        try package("math/rand", "rand", doc: "Package rand implements pseudo-random number generators.", [
            ("func", nil, "Intn", false, "func Intn(n int) int", "Intn returns a non-negative pseudo-random number in [0,n)."),
            ("func", nil, "Read", true, "func Read(p []byte) (n int, err error)", "Read generates len(p) random bytes.\n\nDeprecated: Use [crypto/rand.Read]."),
        ])
        try package("time", "time", doc: "Package time provides functionality for measuring and displaying time.", [
            ("type", nil, "Duration", false, "type Duration int64", "A Duration represents the elapsed time between two instants."),
            ("const", nil, "Nanosecond", false, "const (\n\tNanosecond Duration = 1\n\tSecond = 1000000000 * Nanosecond\n)", "Common durations."),
            ("const", nil, "Second", false, "", nil),
        ], blocks: ["Nanosecond": "Nanosecond", "Second": "Nanosecond"])
        try package("net/http", "http", doc: "Package http provides HTTP client and server implementations.", [
            ("func", nil, "Get", false, "func Get(url string) (resp *Response, err error)", "Get issues a GET to the specified URL."),
        ])

        try db.execute(sql: "INSERT INTO article(slug, title, subtitle, lead, position) VALUES ('spec', 'The Go Programming Language Specification', 'Language version go1.27', NULL, 0)")
        try db.execute(sql: "INSERT INTO article(slug, title, subtitle, lead, position) VALUES ('mem', 'The Go Memory Model', NULL, 'Go programs share memory.', 1)")
        for (position, s) in [("Introduction", 2, "Introduction", "This is the reference manual for the Go programming language."),
                              ("Types", 2, "Types", "A type determines a set of values."),
                              ("Slice_types", 3, "Slice types", "A slice is a descriptor for a contiguous segment of an underlying array."),
                              ("Struct_types", 3, "Struct types", "A struct is a sequence of named elements, called fields. See [Slice types](godoc-article:spec#Slice_types).")].enumerated() {
            try db.execute(sql: "INSERT INTO article_section(article_slug, position, anchor, level, title, body) VALUES ('spec',?,?,?,?,?)",
                           arguments: [position, s.0, s.1, s.2, s.3])
            try entry("section", slug: "spec", position: position, title: s.2, body: s.3)
        }
    }

    /// The fixture's names, without a database, for the tests of resolution.
    static let names = DocNames(
        packages: [("fmt", "fmt"), ("strings", "strings"), ("io", "io"), ("crypto/rand", "rand"), ("math/rand", "rand"), ("math/rand/v2", "rand"),
                   ("time", "time"), ("net/http", "http"), ("testing", "testing")],
        symbols: [.init(package: "fmt", name: "Println"), .init(package: "strings", name: "Builder", isType: true),
                  .init(package: "strings", name: "Builder.WriteString"), .init(package: "io", name: "Reader", isType: true),
                  .init(package: "crypto/rand", name: "Read"), .init(package: "math/rand", name: "Read", deprecated: true),
                  .init(package: "math/rand/v2", name: "IntN"), .init(package: "time", name: "Duration", isType: true),
                  .init(package: "time", name: "Second"), .init(package: "net/http", name: "Get"),
                  .init(package: "testing", name: "T", isType: true)])
}

/// Links into the documentation, from text and from Go code.
@Suite struct DocResolutionTests {
    @Test func linksParseAsTheDocumentationWritesThem() {
        #expect(DocLink(string: "godoc:strings") == .package("strings"))
        #expect(DocLink(string: "godoc:strings.Builder") == .symbol(package: "strings", name: "Builder"))
        #expect(DocLink(string: "godoc:io.Reader.Read") == .symbol(package: "io", name: "Reader.Read"))
        #expect(DocLink(string: "godoc:math/rand/v2.IntN") == .symbol(package: "math/rand/v2", name: "IntN"))
        #expect(DocLink(string: "godoc:net/http") == .package("net/http"))
        #expect(DocLink(string: "godoc-article:spec#Struct_types") == .article(slug: "spec", anchor: "Struct_types"))
        #expect(DocLink(string: "godoc-article:spec") == .article(slug: "spec", anchor: nil))
        for bad in ["godoc:", "https://pkg.go.dev/strings", "godoc-article:", "strings.Builder"] { #expect(DocLink(string: bad) == nil, "\(bad)") }
        for link in [DocLink.package("net/http"), .symbol(package: "io", name: "Reader.Read"), .article(slug: "spec", anchor: "Types")] {
            #expect(DocLink(string: link.url) == link)
        }
        #expect(DocLink.symbol(package: "io", name: "Reader.Read").page == .package("io"))
        #expect(DocLink.symbol(package: "io", name: "Reader.Read").anchor == "Reader.Read")
    }

    /// Code in a readme's text: a package by path, or a name qualified by its package.
    @Test func codeInTextResolves() {
        let names = ReferenceFixture.names
        #expect(names.link(forCode: "fmt.Println") == .symbol(package: "fmt", name: "Println"))
        #expect(names.link(forCode: "fmt.Println()") == .symbol(package: "fmt", name: "Println"))
        #expect(names.link(forCode: "*strings.Builder") == .symbol(package: "strings", name: "Builder"))
        #expect(names.link(forCode: "strings.Builder.WriteString") == .symbol(package: "strings", name: "Builder.WriteString"))
        #expect(names.link(forCode: "time.Duration") == .symbol(package: "time", name: "Duration"))
        #expect(names.link(forCode: "*testing.T") == .symbol(package: "testing", name: "T"))
        #expect(names.link(forCode: "fmt") == .package("fmt"))
        #expect(names.link(forCode: "net/http") == .package("net/http"))
        #expect(names.link(forCode: "http.Get") == .symbol(package: "net/http", name: "Get"))
        // A name in several packages of one name: the one not deprecated.
        #expect(names.link(forCode: "rand.Read") == .symbol(package: "crypto/rand", name: "Read"))
        #expect(names.link(forCode: "rand.IntN") == .symbol(package: "math/rand/v2", name: "IntN"))
        // An interface method has no name of its own: it links to its type.
        #expect(names.link(forCode: "io.Reader.Read") == .symbol(package: "io", name: "Reader.Read"))
        for text in ["main.go", "wg.Wait()", "t.Errorf", "go.mod", "json", "fmt.Sprintf(\"%d\", n)", "strings.Missing", "Println", "nil"] {
            #expect(names.link(forCode: text) == nil, "\(text)")
        }
    }

    @Test func importsOfAGoFile() {
        let source = """
        // Package main says hi.
        package main

        import "fmt" // the one everyone has
        import (
        \tstr "strings"
        \t. "math"
        \t_ "embed"
        \t/* a comment */ "net/http"
        \t`time`
        )

        func main() {}
        """
        #expect(GoSource.imports(in: source) == [
            .init(path: "fmt"), .init(path: "strings", name: "str"), .init(path: "math", name: "."), .init(path: "embed", name: "_"),
            .init(path: "net/http"), .init(path: "time"),
        ])
        #expect(GoSource.imports(in: "package main\n\nfunc main() {}\n").isEmpty)
        #expect(GoSource.imports(in: "not go").isEmpty)
    }

    /// Option-click and Show Documentation: the identifier at a place, with what qualifies it,
    /// resolved through the file's imports.
    @Test func identifiersInTheEditorResolveThroughImports() throws {
        let source = """
        package main

        import (
        \t"fmt"
        \tstr "strings"
        \t"net/http"
        )

        // fmt.Println in a comment is not code.
        func main() {
        \tvar b str.Builder
        \tb.WriteString("go")
        \tfmt.Println(b.String(), "fmt.Println")
        \t_, _ = http.Get("http://localhost")
        \tvar d time.Duration
        }
        """
        let ns = source as NSString
        let imports = GoSource.imports(in: source)
        func link(at needle: String, offset: Int = 0, occurrence: Int = 0) -> DocLink? {
            var range = ns.range(of: needle)
            for _ in 0..<occurrence { range = ns.range(of: needle, range: NSRange(location: NSMaxRange(range), length: ns.length - NSMaxRange(range))) }
            guard let selector = GoSource.selector(at: range.location + offset, in: source) else { return nil }
            return ReferenceFixture.names.link(for: selector.parts, imports: imports)
        }
        // On the name and on the package's name, before and after the dot.
        #expect(link(at: "Builder") == .symbol(package: "strings", name: "Builder"))
        #expect(link(at: "Builder", offset: 7) == .symbol(package: "strings", name: "Builder")) // just past its end
        #expect(link(at: "str.Builder") == .package("strings"))
        #expect(link(at: "Println", occurrence: 1) == .symbol(package: "fmt", name: "Println"))
        #expect(link(at: "Get") == .symbol(package: "net/http", name: "Get"))
        // A method through a variable, a package not imported, a comment and a string: none.
        #expect(link(at: "WriteString") == nil)
        #expect(link(at: "Duration") == nil)
        #expect(link(at: "fmt.Println in") == nil)
        #expect(link(at: "fmt.Println\"", offset: 5) == nil)
        let selector = try #require(GoSource.selector(at: ns.range(of: "Builder").location + 2, in: source))
        #expect(selector.parts == ["str", "Builder"] && selector.range == ns.range(of: "str.Builder"))
        // A dot import brings names in without their package's name.
        let dotted = "package main\n\nimport . \"fmt\"\n\nfunc main() { Println() }\n"
        let at = (dotted as NSString).range(of: "Println").location
        #expect(GoSource.selector(at: at, in: dotted).flatMap { ReferenceFixture.names.link(for: $0.parts, imports: GoSource.imports(in: dotted)) }
            == .symbol(package: "fmt", name: "Println"))
    }
}

/// Links in a readme: inline code naming the standard library links to its documentation, and
/// nothing leads out of the app.
@Suite struct ReadmeLinkTests {
    /// Every link's destination in blocks, in order.
    static func destinations(_ blocks: [MDBlock]) -> [String] {
        func inlines(_ items: [MDInline]) -> [String] {
            items.flatMap { item -> [String] in
                switch item {
                case .link(let destination, _, let children): [destination] + inlines(children)
                case .emphasis(let c), .strong(let c), .strike(let c), .superscript(let c): inlines(c)
                default: []
                }
            }
        }
        return blocks.flatMap { block -> [String] in
            switch block {
            case .heading(_, _, let i), .paragraph(let i): inlines(i)
            case .quote(_, let b): destinations(b)
            case .list(_, _, _, let items): items.flatMap { destinations($0.blocks) }
            case .table(_, let header, let rows): (header + rows.flatMap { $0 }).flatMap(inlines)
            case .code, .rule, .html: []
            }
        }
    }

    @Test func inlineCodeLinksButCodeBlocksDoNot() {
        let readme = """
        Print it with `fmt.Println`, or build it in a `*strings.Builder` first; `main.go` holds it.

        - Wait for `time.Second`.

        ```go
        fmt.Println("hi")
        ```

        [Read `fmt.Println` there](godoc:fmt.Println).
        """
        let linked = DocLinker.link(Readme.parse(readme), names: ReferenceFixture.names)
        #expect(linked[0] == .paragraph([
            .text("Print it with "), .link(destination: "godoc:fmt.Println", title: nil, children: [.code("fmt.Println")]),
            .text(", or build it in a "), .link(destination: "godoc:strings.Builder", title: nil, children: [.code("*strings.Builder")]),
            .text(" first; "), .code("main.go"), .text(" holds it."),
        ]))
        #expect(linked[1] == .list(ordered: false, start: 1, tight: true, items: [MDListItem(task: nil, blocks: [.paragraph([
            .text("Wait for "), .link(destination: "godoc:time.Second", title: nil, children: [.code("time.Second")]), .text("."),
        ])])]))
        #expect(linked[2] == .code(language: "go", text: "fmt.Println(\"hi\")"))
        // Code that is a link's text already stays as it is.
        #expect(linked[3] == Readme.parse(readme)[3])
        // Without the documentation, nothing changes.
        #expect(DocLinker.link(Readme.parse(readme), names: .empty) == Readme.parse(readme))
    }

    /// Gosmos opens nothing outside itself: a web address in a readme, linked or bare, does
    /// nothing; documentation links and links to the library's lessons act.
    @Test func webLinksInAReadmeOpenNothing() {
        // Written as a link, bare (the parser autolinks it, as on GitHub), and as an autolink.
        let blocks = Readme.parse("See [the tour](https://go.dev/tour), https://pkg.go.dev/fmt or <mailto:gopher@example.com>.")
        let destinations = Self.destinations(blocks)
        #expect(destinations == ["https://go.dev/tour", "https://pkg.go.dev/fmt", "mailto:gopher@example.com"])
        let lesson: String = Fixtures.codeLesson
        for destination in destinations {
            #expect(LinkTarget(destination: destination, isLesson: { $0 == lesson }) == nil, "\(destination)")
        }
        #expect(LinkTarget(destination: "godoc:fmt.Println") == .doc(.symbol(package: "fmt", name: "Println")))
        #expect(LinkTarget(destination: "godoc-article:spec#Slice_types") == .doc(.article(slug: "spec", anchor: "Slice_types")))
        #expect(LinkTarget(destination: "https://example.com/lessons/\(lesson)", isLesson: { $0 == lesson }) == .lesson(lesson))
        #expect(LinkTarget(destination: "../\(lesson)", isLesson: { $0 == lesson }) == .lesson(lesson))
    }
}

/// The reference database, as the browser, search and Quick Help read it.
@Suite struct ReferenceDatabaseTests {
    @Test func searchPutsExactNamesFirst() async throws {
        let (reference, dir) = try ReferenceFixture.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        #expect(try await reference.search("Println").map(\.title).prefix(3) == ["fmt.Println", "log.Println", "log.Logger.Println"])
        #expect(try await reference.search("println").first?.title == "fmt.Println")
        let builder = try await reference.search("Builder")
        #expect(builder.first?.link == .symbol(package: "strings", name: "Builder"))
        #expect(builder.map(\.title).prefix(2) == ["strings.Builder", "archive/builder.Builder"])
        #expect(try await reference.search("strings.Builder.Gr").first?.title == "strings.Builder.Grow")
        #expect(try await reference.search("Builder.Grow").first?.title == "strings.Builder.Grow")
        #expect(try await reference.search("rand").prefix(2).map(\.title) == ["math/rand", "crypto/rand"])
        // Prefixes, then the full text: sections of the articles too.
        #expect(try await reference.search("Prin").map(\.title).prefix(3) == ["fmt.Printf", "fmt.Println", "log.Println"])
        let slices = try await reference.search("slice")
        #expect(slices.contains { $0.link == .article(slug: "spec", anchor: "Slice_types") && $0.kind == .section && $0.detail == "The Go Programming Language Specification" })
        #expect(try await reference.search("  ").isEmpty)
        #expect(try await reference.search("\"*(").isEmpty)
    }

    @Test func packagePagesFollowGoDoc() async throws {
        let (reference, dir) = try ReferenceFixture.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        guard case .package(let strings)? = try await reference.page(.symbol(package: "strings", name: "Builder.Grow")) else {
            Issue.record("no strings page")
            return
        }
        #expect(strings.name == "strings" && strings.synopsis == "Package strings implements simple functions to manipulate UTF-8 encoded strings.")
        #expect(strings.sections.map(\.title) == ["Functions", "Types"])
        #expect(strings.sections[1].declarations.map(\.id) == ["Builder", "Builder.Grow", "Builder.WriteString", "Reader", "NewReader", "Reader.Read"])
        #expect(strings.sections[1].declarations.map(\.title) == ["type Builder", "func (*Builder) Grow", "func (*Builder) WriteString", "type Reader",
                                                                  "func NewReader", "func (*Reader) Read"])
        // NewReader comes after Reader: it is grouped with its type.
        #expect(strings.index.map(\.text) == ["func Clone(s string) string", "type Builder", "func (b *Builder) Grow(n int)",
                                              "func (b *Builder) WriteString(s string) (int, error)", "type Reader",
                                              "func NewReader(s string) *Reader", "func (r *Reader) Read(b []byte) (n int, err error)"])
        #expect(strings.index.filter(\.isMember).count == 4)
        #expect(strings.sections[1].declarations[0].examples.map(\.output) == ["ignition\n"])
        #expect(strings.examples.map(\.title) == ["Example (fields)"] && strings.examples[0].output == nil)
        // The package doc's bare web address is autolinked by the parser, and does nothing.
        let addresses = ReadmeLinkTests.destinations(strings.doc)
        #expect(addresses == ["https://blog.golang.org/strings"])
        #expect(addresses.allSatisfy { LinkTarget(destination: $0) == nil })

        guard case .package(let io)? = try await reference.page(.package("io")) else { Issue.record("no io page"); return }
        // A const block shows once, for all its names; an interface method shows its type.
        #expect(io.sections.map(\.title) == ["Constants", "Types"])
        #expect(io.sections[0].declarations.count == 1 && io.sections[0].declarations[0].names == ["SeekStart", "SeekCurrent"])
        #expect(io.declaration(for: "SeekCurrent") == "SeekStart")
        #expect(io.declaration(for: "Reader.Read") == "Reader")
        #expect(io.declaration(for: "Missing") == nil)
        #expect(io.index.first?.text == "Constants")

        guard case .package(let time)? = try await reference.page(.package("time")) else { Issue.record("no time page"); return }
        #expect(time.sections.map(\.title) == ["Types"] && time.sections[0].declarations.map(\.isMember) == [false, true])

        guard case .package(let rand)? = try await reference.page(.package("math/rand")) else { Issue.record("no math/rand page"); return }
        #expect(rand.sections[0].declarations.map(\.deprecated) == [false, true])
        #expect(try await reference.page(.package("missing")) == nil)
    }

    @Test func articlesQuickHelpAndCredits() async throws {
        let (reference, dir) = try ReferenceFixture.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        let outline = try await reference.outline()
        #expect(outline.articles.map(\.slug) == ["spec", "mem"])
        #expect(outline.groups.map(\.top) == ["archive", "crypto", "fmt", "io", "log", "math", "net", "strings", "time"])
        #expect(outline.groups.first { $0.top == "math" }?.packages.map(\.path) == ["math/rand"])

        guard case .article(let spec)? = try await reference.page(.article(slug: "spec", anchor: "Slice_types")) else { Issue.record("no spec"); return }
        #expect(spec.subtitle == "Language version go1.27" && spec.sections.count == 4)
        #expect(spec.contents.map(\.anchor) == ["Introduction", "Types", "Slice_types", "Struct_types"])
        // Links between sections stay links: they act in the app.
        #expect(spec.sections[3].body == [.paragraph([.text("A struct is a sequence of named elements, called fields. See "),
                                                      .link(destination: "godoc-article:spec#Slice_types", title: nil, children: [.text("Slice types")]),
                                                      .text(".")])])

        let println = try #require(try await reference.summary(.symbol(package: "fmt", name: "Println")))
        #expect(println.title == "fmt.Println" && println.declaration == "func Println(a ...any) (n int, err error)")
        #expect(Readme.plainText(println.paragraph).hasPrefix("Println formats using the default formats") && println.paragraph.count == 1)
        // A name in a block: its own line.
        #expect(try await reference.summary(.symbol(package: "io", name: "SeekCurrent"))?.declaration == "const SeekCurrent = 1 // seek relative to the current offset")
        #expect(try await reference.summary(.symbol(package: "io", name: "Reader.Read"))?.declaration.hasPrefix("type Reader interface") == true)
        #expect(try await reference.summary(.package("strings"))?.declaration == "package strings")
        #expect(try await reference.summary(.symbol(package: "strings", name: "Missing")) == nil)

        let credits = try await reference.credits()
        #expect(credits.goVersion == "go1.27.3" && credits.license.hasPrefix("Copyright 2009 The Go Authors."))
        #expect(Readme.plainText(credits.credit).contains("CC BY 4.0"))

        let names = try reference.names()
        #expect(names.link(forCode: "strings.Builder") == .symbol(package: "strings", name: "Builder"))
        #expect(names.link(forCode: "io.SeekCurrent") == .symbol(package: "io", name: "SeekCurrent"))
        #expect(names.link(forCode: "rand.Read") == .symbol(package: "crypto/rand", name: "Read"))
    }

    @Test func refusesOtherFiles() throws {
        let (_, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        let library = dir.appendingPathComponent(Fixtures.libraryName)
        #expect(throws: ReferenceDatabaseError.notReference(path: library.path, applicationID: AppDatabase.applicationID)) {
            try ReferenceDatabase(url: library)
        }
        let missing = dir.appendingPathComponent("go.sqlite")
        #expect(throws: ReferenceDatabaseError.missing(path: missing.path)) { try ReferenceDatabase(url: missing) }
    }
}
