import Foundation
import GRDB

/// Why the Go documentation could not be opened.
public enum ReferenceDatabaseError: Error, LocalizedError, Equatable, Sendable {
    case missing(path: String)
    case notReference(path: String, applicationID: Int)
    case otherVersion(path: String, version: Int)
    case sqlite(String)

    public var errorDescription: String? {
        switch self {
        case .missing(let path): "The Go documentation is missing: there is no file at \(path)."
        case .notReference(let path, let id): "\(path) is not the Go documentation (application_id \(id), not \(ReferenceDatabase.applicationID))."
        case .otherVersion(let path, let version): "\(path) holds the Go documentation in version \(version) of its format; this Gosmos reads version \(ReferenceDatabase.schemaVersion)."
        case .sqlite(let message): "The Go documentation can't be read: \(message)"
        }
    }
}

/// The Go documentation built into the app (db/reference.sql): the standard library, the
/// specification, the memory model, Effective Go, the FAQ and How to Write Go Code, through
/// one read-only GRDB DatabaseQueue. Like the courses built into the app it is opened
/// immutable, so SQLite reads it without any locking.
public final class ReferenceDatabase: Sendable {
    /// 'GSMR'.
    public static let applicationID = 1_196_641_618
    public static let schemaVersion = 1

    let dbQueue: DatabaseQueue

    /// Where the app keeps it: Resources/Reference/go.sqlite.
    public static var bundledURL: URL? {
        guard let url = Bundle.main.resourceURL?.appendingPathComponent("Reference/go.sqlite", isDirectory: false),
              FileManager.default.fileExists(atPath: url.path) else { return nil }
        return url
    }

    /// Opens the file read-only and immutable, after checking it is the Go documentation in
    /// the format this app reads.
    public init(url: URL) throws(ReferenceDatabaseError) {
        let path = url.path
        guard FileManager.default.fileExists(atPath: path) else { throw .missing(path: path) }
        var config = Configuration()
        config.readonly = true
        config.foreignKeysEnabled = false // nothing is written
        config.label = "Go documentation"
        do {
            // A URI, so SQLite takes the immutable flag.
            let queue = try DatabaseQueue(path: URL(fileURLWithPath: path).absoluteString + "?immutable=1", configuration: config)
            let row = try queue.read { db in
                try Row.fetchOne(db, sql: "SELECT (SELECT application_id FROM pragma_application_id) AS a, (SELECT user_version FROM pragma_user_version) AS v")
            }
            let id: Int = row?["a"] ?? 0, version: Int = row?["v"] ?? 0
            if id != ReferenceDatabase.applicationID { throw ReferenceDatabaseError.notReference(path: path, applicationID: id) }
            if version != ReferenceDatabase.schemaVersion { throw ReferenceDatabaseError.otherVersion(path: path, version: version) }
            dbQueue = queue
        } catch let error as ReferenceDatabaseError {
            throw error
        } catch {
            throw .sqlite(AppDatabase.describe(error))
        }
    }

    // MARK: Reads

    /// Every package and documented name, for links in readmes and Quick Help. It waits for
    /// the read, so it is called off the main actor.
    public func names() throws -> DocNames {
        try dbQueue.read { db in try ReferenceDatabase.fetchNames(db) }
    }

    public func outline() async throws -> DocOutline {
        try await dbQueue.read { db in try ReferenceDatabase.fetchOutline(db) }
    }

    /// The page a place is on, its Markdown parsed; nil for a page the documentation lacks.
    public func page(_ link: DocLink) async throws -> DocPage? {
        switch link.page {
        case .package(let path):
            guard let rows = try await dbQueue.read({ db in try ReferenceDatabase.fetchPackageRows(db, path) }) else { return nil }
            return .package(ReferenceDatabase.package(rows))
        case .article(let slug, _):
            guard let rows = try await dbQueue.read({ db in try ReferenceDatabase.fetchArticleRows(db, slug) }) else { return nil }
            return .article(ReferenceDatabase.article(rows))
        case .acknowledgements:
            return .acknowledgements(try await credits())
        case .symbol:
            return nil
        }
    }

    public func credits() async throws -> DocCredits {
        let source = try await dbQueue.read { db -> (credit: String, license: String, goVersion: String) in
            let row = try Row.fetchOne(db, sql: "SELECT credit, license, go_version FROM source")
            return (row?["credit"] ?? "", row?["license"] ?? "", row?["go_version"] ?? "")
        }
        return DocCredits(credit: Readme.parse(source.credit), license: source.license, goVersion: source.goVersion)
    }

    /// Quick Help for a package or a name: its declaration and the first paragraph of its doc.
    public func summary(_ link: DocLink) async throws -> DocSummary? {
        try await dbQueue.read { db in try ReferenceDatabase.fetchSummary(db, link) }
    }

    /// Search while you type: first the names the text is (`Println`, `strings.Builder`,
    /// `rand`), then the names it starts, then what the full-text index finds in names, docs
    /// and the articles, best first.
    public func search(_ text: String, limit: Int = 50) async throws -> [DocSearchHit] {
        try await dbQueue.read { db in try ReferenceDatabase.fetchSearch(db, text, limit: limit) }
    }

    // MARK: Queries shared by the async methods and the tests

    static func fetchNames(_ db: Database) throws -> DocNames {
        let packages = try Row.fetchAll(db, sql: "SELECT import_path, name FROM package").map { (path: $0["import_path"] as String, name: $0["name"] as String) }
        let symbols = try Row.fetchAll(db, sql: """
            SELECT p.import_path, s.receiver, s.name, s.kind, s.deprecated FROM symbol s JOIN package p ON p.id = s.package_id
            """).map { row in
            DocNames.Symbol(package: row["import_path"], name: anchor(row["receiver"], row["name"]), isType: row["kind"] == "type",
                            deprecated: row["deprecated"] as Int != 0)
        }
        return DocNames(packages: packages, symbols: symbols)
    }

    static func fetchOutline(_ db: Database) throws -> DocOutline {
        let articles = try Row.fetchAll(db, sql: "SELECT slug, title FROM article ORDER BY position").map {
            DocOutline.Article(slug: $0["slug"], title: $0["title"])
        }
        var groups: [DocOutline.Group] = []
        for row in try Row.fetchAll(db, sql: "SELECT import_path, synopsis, deprecated FROM package ORDER BY import_path") {
            let package = DocOutline.Package(path: row["import_path"], synopsis: row["synopsis"], deprecated: row["deprecated"] as Int != 0)
            let top = String(package.path.prefix { $0 != "/" })
            if groups.last?.top == top {
                groups[groups.count - 1] = DocOutline.Group(top: top, packages: groups[groups.count - 1].packages + [package])
            } else {
                groups.append(DocOutline.Group(top: top, packages: [package]))
            }
        }
        return DocOutline(articles: articles, groups: groups)
    }

    /// A symbol's anchor on its package's page: `Name`, or `Type.Name` for a method.
    static func anchor(_ receiver: String?, _ name: String) -> String {
        receiver.map { $0 + "." + name } ?? name
    }

    // MARK: Package pages

    /// A package's rows, read in one go; the page is made from them outside the read.
    struct PackageRows: Sendable {
        struct Symbol: Sendable {
            let id: Int64
            let declaration: Int64
            let kind: String
            let receiver: String?
            let name: String
            let deprecated: Bool
        }

        struct Declaration: Sendable {
            let source: String
            let doc: String?
        }

        struct Example: Sendable {
            let id: Int
            let symbol: Int64?
            let suffix: String?
            let doc: String?
            let code: String
            let output: String?
            let unordered: Bool
        }

        let path: String
        let name: String
        let synopsis: String?
        let doc: String?
        let deprecated: Bool
        let symbols: [Symbol]
        /// By id: a const or var block is read once, however many names it has.
        let declarations: [Int64: Declaration]
        let examples: [Example]
    }

    static func fetchPackageRows(_ db: Database, _ path: String) throws -> PackageRows? {
        guard let package = try Row.fetchOne(db, sql: "SELECT id, import_path, name, synopsis, doc, deprecated FROM package WHERE import_path = ?",
                                             arguments: [path]) else { return nil }
        let id: Int64 = package["id"]
        let symbols = try Row.fetchAll(db, sql: """
            SELECT id, declaration_id, kind, receiver, name, deprecated FROM symbol WHERE package_id = ? ORDER BY position
            """, arguments: [id]).map {
            PackageRows.Symbol(id: $0["id"], declaration: $0["declaration_id"], kind: $0["kind"], receiver: $0["receiver"], name: $0["name"],
                               deprecated: $0["deprecated"] as Int != 0)
        }
        var declarations: [Int64: PackageRows.Declaration] = [:]
        for row in try Row.fetchAll(db, sql: """
            SELECT id, source, doc FROM declaration WHERE id IN (SELECT declaration_id FROM symbol WHERE package_id = ?)
            """, arguments: [id]) {
            declarations[row["id"]] = PackageRows.Declaration(source: row["source"], doc: row["doc"])
        }
        let examples = try Row.fetchAll(db, sql: """
            SELECT id, symbol_id, suffix, doc, code, output, unordered FROM example WHERE package_id = ? ORDER BY position
            """, arguments: [id]).map {
            PackageRows.Example(id: $0["id"], symbol: $0["symbol_id"], suffix: $0["suffix"], doc: $0["doc"], code: $0["code"],
                                output: $0["output"], unordered: $0["unordered"] as Int != 0)
        }
        return PackageRows(path: package["import_path"], name: package["name"], synopsis: package["synopsis"], doc: package["doc"],
                           deprecated: package["deprecated"] as Int != 0, symbols: symbols, declarations: declarations, examples: examples)
    }

    /// Lays a package's page out as go doc does: consts, vars and funcs, then each type with
    /// the consts, vars and funcs grouped with it and its methods. A const or var block shows
    /// once, at its first name.
    static func package(_ rows: PackageRows) -> DocPackage {
        var examples: [Int64?: [DocExample]] = [:]
        for row in rows.examples {
            let example = DocExample(id: row.id, title: row.suffix.map { "Example (\($0))" } ?? "Example", doc: Readme.parse(row.doc ?? ""),
                                     code: row.code, output: row.output, unordered: row.unordered)
            examples[row.symbol, default: []].append(example)
        }

        struct Pending {
            var first: PackageRows.Symbol
            var names: [String]
            var examples: [DocExample]
            var deprecated: Bool
            var isMember: Bool
        }
        var order: [Int64] = [], pending: [Int64: Pending] = [:]
        var inType = false
        for row in rows.symbols {
            if row.kind == "type" { inType = true }
            if pending[row.declaration] == nil {
                order.append(row.declaration)
                pending[row.declaration] = Pending(first: row, names: [], examples: [], deprecated: true, isMember: inType && row.kind != "type")
            }
            pending[row.declaration]?.names.append(anchor(row.receiver, row.name))
            pending[row.declaration]?.examples += examples[row.id] ?? []
            // A block is deprecated when all its names are.
            if !row.deprecated { pending[row.declaration]?.deprecated = false }
        }

        var constants: [DocDeclaration] = [], variables: [DocDeclaration] = [], functions: [DocDeclaration] = [], types: [DocDeclaration] = []
        var index: [DocIndexEntry] = []
        for id in order {
            guard let p = pending[id], let kind = DocKind(rawValue: p.first.kind), let first = p.names.first,
                  let row = rows.declarations[id] else { continue }
            let source = row.source
            let declaration = DocDeclaration(id: first, kind: kind, names: p.names, title: title(kind, name: p.first.name, source: source),
                                             source: source, doc: Readme.parse(row.doc ?? ""), deprecated: p.deprecated,
                                             isMember: p.isMember, examples: p.examples)
            switch kind {
            case .const where !p.isMember: constants.append(declaration)
            case .var where !p.isMember: variables.append(declaration)
            case .func where !p.isMember: functions.append(declaration)
            default: types.append(declaration)
            }
            if kind == .func || kind == .method || kind == .type {
                index.append(DocIndexEntry(anchor: first, text: kind == .type ? "type " + p.first.name : signature(source),
                                           isMember: p.isMember, deprecated: p.deprecated))
            }
        }
        // The index starts with the package's constants and variables, as on pkg.go.dev.
        var head: [DocIndexEntry] = []
        if let c = constants.first { head.append(DocIndexEntry(anchor: c.id, text: "Constants", isMember: false, deprecated: false)) }
        if let v = variables.first { head.append(DocIndexEntry(anchor: v.id, text: "Variables", isMember: false, deprecated: false)) }

        let sections = [("Constants", constants), ("Variables", variables), ("Functions", functions), ("Types", types)]
            .filter { !$0.1.isEmpty }.map { DocPackage.DocSection(title: $0.0, declarations: $0.1) }
        return DocPackage(path: rows.path, name: rows.name, synopsis: rows.synopsis, doc: Readme.parse(rows.doc ?? ""), deprecated: rows.deprecated,
                          examples: examples[nil] ?? [], index: head + index, sections: sections)
    }

    /// A declaration's heading: `func Println`, `func (*Builder) Grow`, `type Builder`; none for
    /// a const or var block.
    static func title(_ kind: DocKind, name: String, source: String) -> String? {
        switch kind {
        case .const, .var: return nil
        case .type: return "type " + name
        case .func: return "func " + name
        case .method:
            // The receiver's type as the source writes it, without the receiver's name.
            guard let open = source.firstIndex(of: "("), let close = source[open...].firstIndex(of: ")") else { return "func " + name }
            let receiver = source[source.index(after: open)..<close].split(separator: " ").last.map(String.init) ?? ""
            return "func (\(receiver)) \(name)"
        }
    }

    /// A declaration on one line, for the index.
    static func signature(_ source: String) -> String {
        source.split(whereSeparator: \.isNewline).map { $0.trimmingCharacters(in: .whitespaces) }.joined(separator: " ")
            .replacingOccurrences(of: "( ", with: "(").replacingOccurrences(of: ", )", with: ")")
    }

    // MARK: Articles

    struct ArticleRows: Sendable {
        struct Section: Sendable {
            let anchor: String
            let level: Int
            let title: String
            let body: String?
        }

        let slug: String
        let title: String
        let subtitle: String?
        let lead: String?
        let sections: [Section]
    }

    static func fetchArticleRows(_ db: Database, _ slug: String) throws -> ArticleRows? {
        guard let article = try Row.fetchOne(db, sql: "SELECT slug, title, subtitle, lead FROM article WHERE slug = ?", arguments: [slug]) else { return nil }
        let sections = try Row.fetchAll(db, sql: "SELECT anchor, level, title, body FROM article_section WHERE article_slug = ? ORDER BY position",
                                        arguments: [slug]).map {
            ArticleRows.Section(anchor: $0["anchor"], level: $0["level"], title: $0["title"], body: $0["body"])
        }
        return ArticleRows(slug: article["slug"], title: article["title"], subtitle: article["subtitle"], lead: article["lead"], sections: sections)
    }

    static func article(_ rows: ArticleRows) -> DocArticle {
        DocArticle(slug: rows.slug, title: rows.title, subtitle: rows.subtitle, lead: Readme.parse(rows.lead ?? ""),
                   sections: rows.sections.map { DocArticle.Section(anchor: $0.anchor, level: $0.level, title: $0.title, body: Readme.parse($0.body ?? "")) })
    }

    // MARK: Quick Help

    static func fetchSummary(_ db: Database, _ link: DocLink) throws -> DocSummary? {
        switch link {
        case .package(let path):
            guard let row = try Row.fetchOne(db, sql: "SELECT name, doc, deprecated FROM package WHERE import_path = ?", arguments: [path]) else { return nil }
            return DocSummary(link: link, title: path, declaration: "package " + (row["name"] as String),
                              paragraph: firstParagraph(row["doc"]), deprecated: row["deprecated"] as Int != 0)
        case .symbol(let path, let name):
            let parts = name.split(separator: ".").map(String.init)
            let sql = """
                SELECT s.kind, s.name, s.deprecated, d.source, d.doc
                FROM symbol s JOIN package p ON p.id = s.package_id JOIN declaration d ON d.id = s.declaration_id
                WHERE p.import_path = ? AND s.name = ? AND s.receiver IS ?
                """
            var row: Row?
            if parts.count == 2 {
                row = try Row.fetchOne(db, sql: sql, arguments: [path, parts[1], parts[0]])
                // A field or an interface method: its type says what it is.
                if row == nil { row = try Row.fetchOne(db, sql: sql, arguments: [path, parts[0], nil]) }
            } else {
                row = try Row.fetchOne(db, sql: sql, arguments: [path, name, nil])
            }
            guard let row, let kind = DocKind(rawValue: row["kind"]) else { return nil }
            return DocSummary(link: link, title: path + "." + name, declaration: declaration(kind, name: row["name"], source: row["source"]),
                              paragraph: firstParagraph(row["doc"]), deprecated: row["deprecated"] as Int != 0)
        case .article, .acknowledgements:
            return nil
        }
    }

    /// A declaration for Quick Help: for a name in a const or var block, its own line of the
    /// block, after the keyword.
    static func declaration(_ kind: DocKind, name: String, source: String) -> String {
        guard kind == .const || kind == .var, source.contains("\n") else { return source }
        let line = source.split(whereSeparator: \.isNewline).map { $0.trimmingCharacters(in: .whitespaces) }.first { line in
            line.hasPrefix(name) && !(line.dropFirst(name.count).first.map { $0.isLetter || $0.isNumber || $0 == "_" } ?? false)
        }
        return line.map { kind.rawValue + " " + $0 } ?? source
    }

    static func firstParagraph(_ markdown: String?) -> [MDBlock] {
        guard let markdown else { return [] }
        return Readme.parse(markdown).first { if case .paragraph = $0 { true } else { false } }.map { [$0] } ?? []
    }

    // MARK: Search

    static func fetchSearch(_ db: Database, _ text: String, limit: Int) throws -> [DocSearchHit] {
        let query = text.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !query.isEmpty, limit > 0 else { return [] }
        var hits: [DocSearchHit] = [], seen: Set<DocLink> = []
        func add(_ hit: DocSearchHit) {
            if hits.count < limit, seen.insert(hit.link).inserted { hits.append(hit) }
        }
        if let name = NameQuery(query) {
            for hit in try fetchPackages(db, name, prefix: false, limit: limit) { add(hit) }
            for hit in try fetchSymbols(db, name, prefix: false, limit: limit) { add(hit) }
            for hit in try fetchPackages(db, name, prefix: true, limit: limit) { add(hit) }
            for hit in try fetchSymbols(db, name, prefix: true, limit: limit) { add(hit) }
        }
        if hits.count < limit, let match = ftsQuery(query) {
            for hit in try fetchFullText(db, match, limit: limit + hits.count) { add(hit) }
        }
        return hits
    }

    /// The text as a Go name: a package's path or name (`rand`, `net/http`), or a name with
    /// what qualifies it: `Println`, `fmt.Println`, `Builder.Grow`, `strings.Builder.Gr`.
    struct NameQuery {
        /// A package's name or path, or a lone word that may be one.
        var package: String?
        var receiver: String?
        var name: String

        init?(_ text: String) {
            guard text.allSatisfy({ $0.isLetter || $0.isNumber || $0 == "_" || $0 == "." || $0 == "/" }) else { return nil }
            let afterPath = text.lastIndex(of: "/").map { text.index(after: $0) } ?? text.startIndex
            var parts = [String]()
            if let dot = text[afterPath...].firstIndex(of: ".") {
                parts = [String(text[..<dot])] + text[text.index(after: dot)...].split(separator: ".", omittingEmptySubsequences: false).map(String.init)
            } else {
                parts = [text]
            }
            guard parts.count <= 3, parts.allSatisfy({ !$0.isEmpty }) || parts.last == "" else { return nil }
            switch parts.count {
            case 1:
                name = parts[0]
            case 2:
                // Packages are lower case and exported names upper case: Builder.Grow is a method.
                if parts[0].first?.isUppercase == true {
                    receiver = parts[0]
                } else {
                    package = parts[0]
                }
                name = parts[1]
            default:
                package = parts[0]
                receiver = parts[1]
                name = parts[2]
            }
        }
    }

    /// Packages whose path or name is the query, or starts with it, in the same case: package
    /// names are lower case, so `Builder` means a name, not the package builder.
    static func fetchPackages(_ db: Database, _ q: NameQuery, prefix: Bool, limit: Int) throws -> [DocSearchHit] {
        guard q.package == nil, q.receiver == nil, !q.name.isEmpty else { return [] }
        // The query has letters, digits, _ and / only (NameQuery), none special to GLOB.
        let condition = prefix ? "(import_path GLOB :q || '*' OR name GLOB :q || '*')" : "(import_path = :q OR name = :q)"
        return try Row.fetchAll(db, sql: """
            SELECT import_path, synopsis, deprecated FROM package WHERE \(condition)
            ORDER BY length(import_path), import_path LIMIT :limit
            """, arguments: ["q": q.name, "limit": limit]).map {
            DocSearchHit(link: .package($0["import_path"]), kind: .package, title: $0["import_path"], detail: $0["synopsis"],
                         deprecated: $0["deprecated"] as Int != 0)
        }
    }

    /// Names that are the query, or start with it: funcs and types before consts and vars and
    /// those before methods; then by their package's path, shortest first. The same case as
    /// typed comes first.
    static func fetchSymbols(_ db: Database, _ q: NameQuery, prefix: Bool, limit: Int) throws -> [DocSearchHit] {
        guard !q.name.isEmpty || q.package != nil || q.receiver != nil else { return [] }
        var conditions = [prefix ? "s.name LIKE :p ESCAPE '\\'" : "s.name = :n COLLATE NOCASE"]
        if q.name.isEmpty { conditions = ["1"] }
        if q.package != nil { conditions.append("(p.name = :pkg OR p.import_path = :pkg)") }
        if q.receiver != nil { conditions.append("s.receiver = :r COLLATE NOCASE") }
        return try Row.fetchAll(db, sql: """
            SELECT p.import_path, s.kind, s.receiver, s.name, s.synopsis, s.deprecated
            FROM symbol s JOIN package p ON p.id = s.package_id
            WHERE \(conditions.joined(separator: " AND "))
            ORDER BY s.name <> :n, \(prefix ? "length(s.name), " : "")
                     CASE s.kind WHEN 'func' THEN 0 WHEN 'type' THEN 0 WHEN 'method' THEN 2 ELSE 1 END,
                     length(p.import_path), p.import_path, s.receiver
            LIMIT :limit
            """, arguments: ["n": q.name, "p": likePrefix(q.name), "pkg": q.package, "r": q.receiver, "limit": limit]).compactMap(symbolHit)
    }

    static func symbolHit(_ row: Row) -> DocSearchHit? {
        guard let kind = DocKind(rawValue: row["kind"]) else { return nil }
        let path: String = row["import_path"], name = anchor(row["receiver"], row["name"])
        return DocSearchHit(link: .symbol(package: path, name: name), kind: .symbol(kind), title: path + "." + name,
                            detail: row["synopsis"], deprecated: row["deprecated"] as Int != 0)
    }

    /// What the full-text index finds, names weighing ten times the text.
    static func fetchFullText(_ db: Database, _ match: String, limit: Int) throws -> [DocSearchHit] {
        try Row.fetchAll(db, sql: """
            WITH found AS (
                SELECT rowid AS id, bm25(search, 10.0, 1.0) AS score FROM search WHERE search MATCH ? ORDER BY score LIMIT ?
            )
            SELECT e.kind AS entry, e.title, p.import_path AS package_path, p.synopsis AS package_synopsis, p.deprecated AS package_deprecated,
                   sp.import_path, s.kind, s.receiver, s.name, s.synopsis, s.deprecated,
                   sec.anchor, sec.article_slug, a.title AS article_title
            FROM found JOIN search_entry e ON e.id = found.id
            LEFT JOIN package p ON p.id = e.package_id
            LEFT JOIN symbol s ON s.id = e.symbol_id
            LEFT JOIN package sp ON sp.id = s.package_id
            LEFT JOIN article_section sec ON sec.article_slug = e.article_slug AND sec.position = e.section_position
            LEFT JOIN article a ON a.slug = e.article_slug
            ORDER BY found.score
            """, arguments: [match, limit]).compactMap { row in
            switch row["entry"] as String {
            case "package":
                return DocSearchHit(link: .package(row["package_path"]), kind: .package, title: row["package_path"], detail: row["package_synopsis"],
                                    deprecated: row["package_deprecated"] as Int != 0)
            case "symbol":
                return symbolHit(row)
            default:
                return DocSearchHit(link: .article(slug: row["article_slug"], anchor: row["anchor"]), kind: .section, title: row["title"],
                                    detail: row["article_title"], deprecated: false)
            }
        }
    }

    /// The text as an FTS5 query: every word must appear, the last as a prefix; quotes keep
    /// FTS5's syntax characters literal. Words split where the index splits them, so
    /// `strings.Buil` is "strings" "buil"*.
    static func ftsQuery(_ text: String) -> String? {
        let words = text.split(whereSeparator: { !$0.isLetter && !$0.isNumber }).map(String.init)
        guard !words.isEmpty else { return nil }
        return words.enumerated().map { i, w in
            "\"" + w.replacingOccurrences(of: "\"", with: "\"\"") + "\"" + (i == words.count - 1 ? "*" : "")
        }.joined(separator: " ")
    }

    /// A LIKE pattern for names starting with `text`.
    static func likePrefix(_ text: String) -> String {
        text.replacingOccurrences(of: "\\", with: "\\\\").replacingOccurrences(of: "%", with: "\\%").replacingOccurrences(of: "_", with: "\\_") + "%"
    }
}
