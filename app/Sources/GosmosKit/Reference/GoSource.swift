import Foundation

/// Go source as Quick Help reads it: a file's imports, and the qualified identifier at a place
/// in the code, such as `strings.Builder` under the pointer.
public enum GoSource {
    /// An import of a Go file. `name` is what the file calls the package (`str "strings"`),
    /// "." for a dot import or "_" for a blank one; nil when the package's own name applies.
    public struct Import: Hashable, Sendable {
        public let path: String
        public let name: String?

        public init(path: String, name: String? = nil) {
            self.path = path
            self.name = name
        }
    }

    /// An identifier with the selectors before it: `["strings", "Builder"]` for Builder in
    /// `strings.Builder`. `range` covers them all, in UTF-16.
    public struct Selector: Hashable, Sendable {
        public let parts: [String]
        public let range: NSRange
        /// Whether the identifier stands where Go allows a package's name: before a dot
        /// (`strings` in `strings.Builder`) or before an import path (`str` in `str "strings"`).
        /// Anywhere else, as `strings` in `len(strings)`, the name is something the code
        /// declared in the package's stead, such as a variable.
        public let canBePackage: Bool

        public init(parts: [String], range: NSRange, canBePackage: Bool = true) {
            self.parts = parts
            self.range = range
            self.canBePackage = canBePackage
        }
    }

    /// The file's imports: the import declarations after its package clause, single or in
    /// parentheses. Code that stops making sense ends the list there.
    public static func imports(in source: String) -> [Import] {
        var tokens = Tokens(Array(source.utf16))
        guard tokens.next() == .word("package"), case .word? = tokens.next() else { return [] }
        var imports: [Import] = []
        while let token = tokens.next() {
            switch token {
            case .punct(0x3B): // ;
                continue
            case .word("import"):
                if tokens.peek() == .punct(0x28) { // (
                    tokens.skip()
                    while let next = tokens.peek(), next != .punct(0x29) { // )
                        if next == .punct(0x3B) {
                            tokens.skip()
                        } else if let spec = spec(&tokens) {
                            imports.append(spec)
                        } else {
                            return imports
                        }
                    }
                    tokens.skip()
                } else if let spec = spec(&tokens) {
                    imports.append(spec)
                } else {
                    return imports
                }
            default:
                return imports
            }
        }
        return imports
    }

    /// One import spec: an optional name, then the path.
    private static func spec(_ tokens: inout Tokens) -> Import? {
        var name: String?
        switch tokens.next() {
        case .word(let word)?: name = word
        case .punct(0x2E)?: name = "." // .
        case .string(let path)?: return Import(path: path)
        default: return nil
        }
        guard case .string(let path)? = tokens.next() else { return nil }
        return Import(path: path, name: name)
    }

    /// The identifier at `offset` (UTF-16), or just before it, with up to two selectors before
    /// it: `pkg.Name` and `pkg.Type.Method`. Nil in a comment or a string, or away from any
    /// identifier.
    public static func selector(at offset: Int, in source: String) -> Selector? {
        let s = Array(source.utf16)
        guard offset >= 0, offset <= s.count else { return nil }
        var at = offset
        if at == s.count || !isIdentifier(s[at]) {
            guard at > 0, isIdentifier(s[at - 1]) else { return nil }
            at -= 1
        }
        var start = at, end = at
        while start > 0, isIdentifier(s[start - 1]) { start -= 1 }
        while end < s.count, isIdentifier(s[end]) { end += 1 }
        guard !isDigit(s[start]) else { return nil }
        let inLiteral = GoLexer.tokens(source).contains { token in
            (token.kind == .comment || token.kind == .string) && NSLocationInRange(start, token.range)
        }
        guard !inLiteral else { return nil }
        func word(_ from: Int, _ to: Int) -> String { String(decoding: s[from..<to], as: UTF16.self) }
        var parts = [word(start, end)]
        var from = start
        while parts.count < 3, from >= 2, s[from - 1] == 0x2E, isIdentifier(s[from - 2]) { // .
            var first = from - 2
            while first > 0, isIdentifier(s[first - 1]) { first -= 1 }
            guard !isDigit(s[first]) else { break }
            parts.insert(word(first, from - 1), at: 0)
            from = first
        }
        // After spaces or tabs: a dot, but not the ... of a variadic argument, or an import path.
        var next = end
        while next < s.count, s[next] == 0x20 || s[next] == 0x09 { next += 1 }
        let dot = next < s.count && s[next] == 0x2E && !(next + 1 < s.count && s[next + 1] == 0x2E)
        let path = next < s.count && (s[next] == 0x22 || s[next] == 0x60) // " or `
        return Selector(parts: parts, range: NSRange(location: from, length: end - from), canBePackage: dot || path)
    }

    static func isDigit(_ c: UInt16) -> Bool { c >= 0x30 && c <= 0x39 }

    /// A letter, a digit or an underscore; anything beyond ASCII counts as a letter, as in the
    /// lexer.
    static func isIdentifier(_ c: UInt16) -> Bool {
        isDigit(c) || (c >= 0x41 && c <= 0x5A) || (c >= 0x61 && c <= 0x7A) || c == 0x5F || c >= 0x80
    }

    /// Go's tokens as far as import declarations need them: words, string literals and
    /// punctuation, without white space and comments.
    struct Tokens {
        enum Token: Equatable {
            case word(String)
            case string(String)
            case punct(UInt16)
        }

        let s: [UInt16]
        var i = 0
        var peeked: Token??

        init(_ s: [UInt16]) { self.s = s }

        mutating func peek() -> Token? {
            if let peeked { return peeked }
            let token = read()
            peeked = .some(token)
            return token
        }

        mutating func skip() { _ = next() }

        mutating func next() -> Token? {
            if let peeked {
                self.peeked = nil
                return peeked
            }
            return read()
        }

        private mutating func read() -> Token? {
            while i < s.count {
                let c = s[i]
                if c == 0x20 || c == 0x09 || c == 0x0A || c == 0x0D {
                    i += 1
                } else if c == 0x2F, i + 1 < s.count, s[i + 1] == 0x2F { // line comment
                    while i < s.count, s[i] != 0x0A { i += 1 }
                } else if c == 0x2F, i + 1 < s.count, s[i + 1] == 0x2A { // general comment
                    i += 2
                    while i < s.count, !(s[i] == 0x2A && i + 1 < s.count && s[i + 1] == 0x2F) { i += 1 }
                    i = min(s.count, i + 2)
                } else {
                    break
                }
            }
            guard i < s.count else { return nil }
            let c = s[i], start = i
            if GoSource.isIdentifier(c) {
                while i < s.count, GoSource.isIdentifier(s[i]) { i += 1 }
                return .word(String(decoding: s[start..<i], as: UTF16.self))
            }
            if c == 0x22 || c == 0x60 { // " or `
                i += 1
                while i < s.count, s[i] != c, c == 0x60 || s[i] != 0x0A {
                    i += c == 0x22 && s[i] == 0x5C ? 2 : 1
                }
                let end = min(i, s.count)
                i = min(s.count, i + 1)
                return .string(String(decoding: s[(start + 1)..<end], as: UTF16.self))
            }
            i += 1
            return .punct(c)
        }
    }
}
