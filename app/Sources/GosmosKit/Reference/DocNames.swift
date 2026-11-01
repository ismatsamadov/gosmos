import Foundation

/// The standard library's names, for turning code into links to its documentation: every
/// package by import path and by its own name, and every documented name in them, methods as
/// `Type.Method`. Made once from the reference database and kept in memory.
public struct DocNames: Sendable {
    /// Package names by import path.
    private let names: [String: String]
    /// Import paths by package name, shortest first: `rand` is math/rand before
    /// math/rand/v2.
    private let paths: [String: [String]]
    /// Whether each name is deprecated, by "path.Name" or "path.Type.Method".
    private let symbols: [String: Bool]
    /// The types, as "path.Type": a field or interface method has no name of its own, and
    /// links to its type.
    private let types: Set<String>

    /// A documented name: `name` is `Name`, or `Type.Method` for a method.
    public struct Symbol: Sendable {
        public let package: String
        public let name: String
        public let isType: Bool
        public let deprecated: Bool

        public init(package: String, name: String, isType: Bool = false, deprecated: Bool = false) {
            self.package = package
            self.name = name
            self.isType = isType
            self.deprecated = deprecated
        }
    }

    public static let empty = DocNames(packages: [], symbols: [])

    public init(packages: [(path: String, name: String)], symbols: [Symbol]) {
        var names: [String: String] = [:], paths: [String: [String]] = [:]
        for p in packages {
            names[p.path] = p.name
            paths[p.name, default: []].append(p.path)
        }
        self.names = names
        self.paths = paths.mapValues { $0.sorted { ($0.count, $0) < ($1.count, $1) } }
        var deprecated: [String: Bool] = [:], types: Set<String> = []
        for s in symbols {
            deprecated[s.package + "." + s.name] = s.deprecated
            if s.isType { types.insert(s.package + "." + s.name) }
        }
        self.symbols = deprecated
        self.types = types
    }

    /// Whether `name` (`Name` or `Type.Name`) has a link in the package at `path`: a name of
    /// its own, or a field or interface method of a type.
    func has(_ name: String, in path: String) -> Bool {
        if symbols[path + "." + name] != nil { return true }
        guard let dot = name.firstIndex(of: ".") else { return false }
        return types.contains(path + "." + name[..<dot])
    }

    public var isEmpty: Bool { names.isEmpty }

    /// The package's own name, by import path: `rand` for math/rand/v2.
    public func packageName(of path: String) -> String? { names[path] }

    /// The link for code written in text: a package by import path (`fmt`, `net/http`), or a
    /// name qualified by its package's name or path (`fmt.Println`, `*strings.Builder`,
    /// `strings.Builder.WriteString`, `crypto/rand.Read`), with a leading `*` or `&` and a
    /// trailing `()` set aside. A package name several packages share stands for the one with
    /// that name in it, not deprecated there, with the shortest path: `rand.Read` is
    /// crypto/rand's. Anything else, like `wg.Wait()` or `main.go`, has none.
    public func link(forCode text: String) -> DocLink? {
        var code = Substring(text.trimmingCharacters(in: .whitespaces))
        while code.first == "*" || code.first == "&" { code = code.dropFirst() }
        if code.hasSuffix("()") { code = code.dropLast(2) }
        guard !code.isEmpty, code.allSatisfy({ $0.isLetter || $0.isNumber || $0 == "_" || $0 == "." || $0 == "/" }) else { return nil }
        let afterPath = code.lastIndex(of: "/").map { code.index(after: $0) } ?? code.startIndex
        guard let dot = code[afterPath...].firstIndex(of: ".") else {
            return names[String(code)] != nil ? .package(String(code)) : nil
        }
        let qualifier = String(code[..<dot]), name = String(code[code.index(after: dot)...])
        let components = name.split(separator: ".", omittingEmptySubsequences: false)
        guard components.count <= 2, components.allSatisfy({ !$0.isEmpty }) else { return nil }
        let candidates = qualifier.contains("/") ? [qualifier] : paths[qualifier] ?? []
        let found = candidates.filter { has(name, in: $0) }
        guard let path = found.first(where: { symbols[$0 + "." + name] != true }) ?? found.first else { return nil }
        return .symbol(package: path, name: name)
    }

    /// The link for an identifier in a Go file, with the selectors before it, through the
    /// file's imports: renamed (`str.Builder` after `import str "strings"`) or dot imports
    /// included. Only the standard library has links.
    public func link(for parts: [String], imports: [GoSource.Import]) -> DocLink? {
        guard let first = parts.first else { return nil }
        let rest = parts.dropFirst().joined(separator: ".")
        if let imported = imports.first(where: { localName(of: $0) == first }) {
            guard names[imported.path] != nil else { return nil }
            if rest.isEmpty { return .package(imported.path) }
            return has(rest, in: imported.path) ? .symbol(package: imported.path, name: rest) : nil
        }
        // A name a dot import brings in, used without its package's name.
        let name = parts.joined(separator: ".")
        for imported in imports where imported.name == "." && has(name, in: imported.path) {
            return .symbol(package: imported.path, name: name)
        }
        return nil
    }

    /// The link for the identifier at a place in a Go file (GoSource.selector), as
    /// `link(for:imports:)` finds it, except that a lone name standing where Go allows no
    /// package's name has none even when an import has that name: there the code declared
    /// something of its own by it, such as a parameter `time` in `len(time)`.
    public func link(for selector: GoSource.Selector, imports: [GoSource.Import]) -> DocLink? {
        if selector.parts.count == 1, !selector.canBePackage, imports.contains(where: { localName(of: $0) == selector.parts[0] }) {
            return nil
        }
        return link(for: selector.parts, imports: imports)
    }

    /// The name a file uses for an import: the one it gives, else the package's own; a dot
    /// or blank import gives none.
    func localName(of imported: GoSource.Import) -> String? {
        if let name = imported.name { return name == "." || name == "_" ? nil : name }
        if let own = names[imported.path] { return own }
        // Outside the standard library: the path's last element, without a version suffix.
        let elements = imported.path.split(separator: "/")
        if elements.count > 1, let last = elements.last, last.first == "v", last.dropFirst().allSatisfy(\.isNumber) {
            return String(elements[elements.count - 2])
        }
        return elements.last.map(String.init)
    }
}
