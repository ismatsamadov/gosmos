import Foundation

/// A place in the Go documentation, as links name it: `godoc:strings` (a package),
/// `godoc:strings.Builder` or `godoc:io.Reader.Read` (a name in one), and
/// `godoc-article:spec#Slice_types` (an article, at a section).
public enum DocLink: Hashable, Sendable {
    case package(String)
    /// A name in a package, as the package's page anchors it: `Builder`, `Reader.Read`.
    case symbol(package: String, name: String)
    case article(slug: String, anchor: String?)
    /// Where the documentation comes from, and its licenses.
    case acknowledgements

    public static let scheme = "godoc"
    public static let articleScheme = "godoc-article"

    public init?(string: String) {
        if let rest = string.prefixed(by: DocLink.scheme + ":") {
            // Import paths have no dots, so the first dot after the last slash starts the name.
            let afterPath = rest.lastIndex(of: "/").map { rest.index(after: $0) } ?? rest.startIndex
            if let dot = rest[afterPath...].firstIndex(of: ".") {
                let package = String(rest[..<dot]), name = String(rest[rest.index(after: dot)...])
                guard !package.isEmpty, !name.isEmpty else { return nil }
                self = .symbol(package: package, name: name)
            } else {
                guard !rest.isEmpty else { return nil }
                self = .package(String(rest))
            }
        } else if let rest = string.prefixed(by: DocLink.articleScheme + ":") {
            let parts = rest.split(separator: "#", maxSplits: 1, omittingEmptySubsequences: false)
            guard let slug = parts.first, !slug.isEmpty else { return nil }
            self = .article(slug: String(slug), anchor: parts.count > 1 && !parts[1].isEmpty ? String(parts[1]) : nil)
        } else {
            return nil
        }
    }

    public init?(_ url: URL) {
        self.init(string: url.absoluteString)
    }

    /// The link as the documentation writes it; the acknowledgements have none.
    public var url: String {
        switch self {
        case .package(let path): DocLink.scheme + ":" + path
        case .symbol(let package, let name): DocLink.scheme + ":" + package + "." + name
        case .article(let slug, let anchor): DocLink.articleScheme + ":" + slug + (anchor.map { "#" + $0 } ?? "")
        case .acknowledgements: ""
        }
    }

    /// The page the place is on: a package's, an article, or the acknowledgements.
    public var page: DocLink {
        switch self {
        case .symbol(let package, _): .package(package)
        case .article(let slug, _): .article(slug: slug, anchor: nil)
        case .package, .acknowledgements: self
        }
    }

    /// Where on its page the place is: a name, or an article's section.
    public var anchor: String? {
        switch self {
        case .symbol(_, let name): name
        case .article(_, let anchor): anchor
        case .package, .acknowledgements: nil
        }
    }
}

extension String {
    /// The rest of the string after `prefix`, or nil when it does not start with it.
    func prefixed(by prefix: String) -> Substring? {
        hasPrefix(prefix) ? dropFirst(prefix.count) : nil
    }
}
