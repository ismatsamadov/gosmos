import Foundation

/// What a name in the standard library is.
public enum DocKind: String, Sendable, CaseIterable {
    case const, `var`, `func`, type, method

    /// The keyword its declaration starts with.
    public var keyword: String { self == .method ? "func" : rawValue }
}

/// The documentation's contents, for the browser's list: the articles about the language, and
/// the standard library's packages by the first element of their paths.
public struct DocOutline: Sendable {
    public struct Article: Identifiable, Hashable, Sendable {
        public let slug: String
        public let title: String
        public var id: String { slug }
        public var link: DocLink { .article(slug: slug, anchor: nil) }
    }

    public struct Package: Identifiable, Hashable, Sendable {
        public let path: String
        public let synopsis: String?
        public let deprecated: Bool
        public var id: String { path }
        public var link: DocLink { .package(path) }
    }

    /// The packages under one first path element: crypto, crypto/aes, crypto/cipher and so on.
    public struct Group: Identifiable, Hashable, Sendable {
        public let top: String
        public let packages: [Package]
        public var id: String { top }
    }

    public let articles: [Article]
    public let groups: [Group]

    public static let empty = DocOutline(articles: [], groups: [])

    public init(articles: [Article], groups: [Group]) {
        self.articles = articles
        self.groups = groups
    }
}

/// A page of the documentation: a package's, an article, or the acknowledgements.
public enum DocPage: Sendable {
    case package(DocPackage)
    case article(DocArticle)
    case acknowledgements(DocCredits)

    public var link: DocLink {
        switch self {
        case .package(let p): .package(p.path)
        case .article(let a): .article(slug: a.slug, anchor: nil)
        case .acknowledgements: .acknowledgements
        }
    }

    public var title: String {
        switch self {
        case .package(let p): p.path
        case .article(let a): a.title
        case .acknowledgements: "Acknowledgements"
        }
    }
}

/// A package's page: its doc, an index of its names, and each declaration with its doc and
/// examples, in go doc's order.
public struct DocPackage: Sendable {
    public let path: String
    public let name: String
    /// The doc's first sentence.
    public let synopsis: String?
    public let doc: [MDBlock]
    public let deprecated: Bool
    /// The package's own examples.
    public let examples: [DocExample]
    public let index: [DocIndexEntry]
    /// Constants, Variables, Functions and Types, those it has.
    public let sections: [DocSection]
    /// The declaration showing each name, by anchor (`Builder.Grow`): a const or var block
    /// shows all its names.
    let declarations: [String: String]

    public init(path: String, name: String, synopsis: String?, doc: [MDBlock], deprecated: Bool, examples: [DocExample],
                index: [DocIndexEntry], sections: [DocSection]) {
        self.path = path
        self.name = name
        self.synopsis = synopsis
        self.doc = doc
        self.deprecated = deprecated
        self.examples = examples
        self.index = index
        self.sections = sections
        var declarations: [String: String] = [:]
        for section in sections {
            for d in section.declarations {
                for name in d.names where declarations[name] == nil { declarations[name] = d.id }
            }
        }
        self.declarations = declarations
    }

    /// The declaration on the page for a name: `Builder.Grow` is the method's, and a field or
    /// interface method such as `Reader.Read` is its type's. Nil for a name not on the page.
    public func declaration(for anchor: String) -> String? {
        if let id = declarations[anchor] { return id }
        if let dot = anchor.firstIndex(of: ".") { return declarations[String(anchor[..<dot])] }
        return nil
    }

    /// A group of declarations: Constants, Variables, Functions or Types.
    public struct DocSection: Identifiable, Sendable {
        public let title: String
        public let declarations: [DocDeclaration]
        public var id: String { title }
    }
}

/// One declaration as go doc shows it: a func, method or type, or a whole const or var block.
public struct DocDeclaration: Identifiable, Sendable {
    /// The anchor of its first name.
    public let id: String
    public let kind: DocKind
    /// The anchors of the names it declares: `Println`, `Builder.Grow`, every name of a block.
    public let names: [String]
    /// Its heading: `func Println`, `type Builder`, `func (*Builder) Grow`; nil for a const or
    /// var block, which shows under its section's heading.
    public let title: String?
    /// The Go source, without its doc comment or body.
    public let source: String
    public let doc: [MDBlock]
    public let deprecated: Bool
    /// Grouped with a type: the type's funcs and methods, and consts and vars of its type.
    public let isMember: Bool
    public let examples: [DocExample]

    public init(id: String, kind: DocKind, names: [String], title: String?, source: String, doc: [MDBlock], deprecated: Bool,
                isMember: Bool, examples: [DocExample]) {
        self.id = id
        self.kind = kind
        self.names = names
        self.title = title
        self.source = source
        self.doc = doc
        self.deprecated = deprecated
        self.isMember = isMember
        self.examples = examples
    }
}

/// An example from a package's tests: its code, and what it prints when that is checked.
public struct DocExample: Identifiable, Sendable {
    public let id: Int
    /// "Example", or "Example (Grow)" for one with a suffix.
    public let title: String
    public let doc: [MDBlock]
    public let code: String
    /// Nil when the example is only compiled; empty when it prints nothing.
    public let output: String?
    /// Its output's lines may come in any order.
    public let unordered: Bool

    public init(id: Int, title: String, doc: [MDBlock], code: String, output: String?, unordered: Bool) {
        self.id = id
        self.title = title
        self.doc = doc
        self.code = code
        self.output = output
        self.unordered = unordered
    }
}

/// A line of a package's index: a const or var section, a func or type, or a type's member,
/// with the anchor it leads to.
public struct DocIndexEntry: Identifiable, Hashable, Sendable {
    public let anchor: String
    /// The declaration's signature on one line, or "Constants" and "Variables".
    public let text: String
    public let isMember: Bool
    public let deprecated: Bool
    public var id: String { anchor }

    public init(anchor: String, text: String, isMember: Bool, deprecated: Bool) {
        self.anchor = anchor
        self.text = text
        self.isMember = isMember
        self.deprecated = deprecated
    }
}

/// An article from go.dev, with its sections in order.
public struct DocArticle: Sendable {
    public let slug: String
    public let title: String
    /// Such as the spec's language version.
    public let subtitle: String?
    /// What comes before the first heading.
    public let lead: [MDBlock]
    public let sections: [Section]

    /// A heading and what follows it up to the next heading.
    public struct Section: Identifiable, Sendable {
        public let anchor: String
        public let level: Int
        public let title: String
        public let body: [MDBlock]
        public var id: String { anchor }

        public init(anchor: String, level: Int, title: String, body: [MDBlock]) {
            self.anchor = anchor
            self.level = level
            self.title = title
            self.body = body
        }
    }

    public init(slug: String, title: String, subtitle: String?, lead: [MDBlock], sections: [Section]) {
        self.slug = slug
        self.title = title
        self.subtitle = subtitle
        self.lead = lead
        self.sections = sections
    }

    /// The table of contents: the sections two levels below the title.
    public var contents: [Section] { sections.filter { $0.level <= 3 } }
}

/// Where the documentation comes from, and its licenses, as the reference database says.
public struct DocCredits: Sendable {
    /// Whose work each part is, and its license.
    public let credit: [MDBlock]
    /// Go's LICENSE file, as it is.
    public let license: String
    /// Such as go1.27.3.
    public let goVersion: String

    public init(credit: [MDBlock], license: String, goVersion: String) {
        self.credit = credit
        self.license = license
        self.goVersion = goVersion
    }
}

/// What search found in the documentation: a package, a name in one, or a section of an article.
public struct DocSearchHit: Identifiable, Hashable, Sendable {
    public enum Kind: Hashable, Sendable {
        case package
        case symbol(DocKind)
        case section
    }

    public let link: DocLink
    public let kind: Kind
    /// The package's path, the qualified name (`strings.Builder.Grow`), or the section's title.
    public let title: String
    /// The synopsis, or the article a section is in.
    public let detail: String?
    public let deprecated: Bool
    public var id: String { link.url }

    public init(link: DocLink, kind: Kind, title: String, detail: String?, deprecated: Bool) {
        self.link = link
        self.kind = kind
        self.title = title
        self.detail = detail
        self.deprecated = deprecated
    }
}

/// Quick Help: a name's declaration and the first paragraph of its doc.
public struct DocSummary: Sendable {
    public let link: DocLink
    /// The qualified name, such as `strings.Builder`.
    public let title: String
    /// Its declaration, or the line of it for a name in a const or var block; a package's
    /// clause.
    public let declaration: String
    public let paragraph: [MDBlock]
    public let deprecated: Bool

    public init(link: DocLink, title: String, declaration: String, paragraph: [MDBlock], deprecated: Bool) {
        self.link = link
        self.title = title
        self.declaration = declaration
        self.paragraph = paragraph
        self.deprecated = deprecated
    }
}
