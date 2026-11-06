import Foundation

/// What a link in a readme or in the documentation does. Gosmos opens nothing outside itself:
/// a link into the Go documentation opens it in the app, a link whose path names a lesson of
/// the readme's library opens that lesson, and a link to a term of its course's field guide
/// shows the term's meaning. Any other address is shown as plain text.
public enum LinkTarget: Hashable, Sendable {
    case doc(DocLink)
    /// A lesson of the library the text belongs to, by uuid.
    case lesson(String)
    /// A term of the field guide of the text's course, as the guide spells it.
    case term(String)

    /// The scheme of the app's own links to lessons: `gosmos-lesson:<uuid>`.
    public static let lessonScheme = "gosmos-lesson"
    /// The scheme of links to the field guide: `gosmos-term:<term>`, the term percent-encoded.
    public static let termScheme = "gosmos-term"

    /// The target of a link's destination, as the Markdown gives it, or as `url` writes it;
    /// nil for one that does nothing. `isLesson` tells whether a uuid names a lesson of the
    /// text's library, and `isTerm` whether a word is a term of its course's field guide.
    public init?(destination: String, isLesson: (String) -> Bool = { _ in false }, isTerm: (String) -> Bool = { _ in false }) {
        let destination = destination.trimmingCharacters(in: .whitespaces)
        if let doc = DocLink(string: destination) {
            self = .doc(doc)
        } else if let id = destination.prefixed(by: LinkTarget.lessonScheme + ":") {
            guard let lesson = id.removingPercentEncoding, isLesson(lesson) else { return nil }
            self = .lesson(lesson)
        } else if let word = destination.prefixed(by: LinkTarget.termScheme + ":") {
            guard let term = word.removingPercentEncoding, isTerm(term) else { return nil }
            self = .term(term)
        } else if let url = URL(string: destination), let lesson = url.pathComponents.last(where: isLesson) {
            self = .lesson(lesson)
        } else {
            return nil
        }
    }

    /// The link as the app writes it into text: the documentation's own (`godoc:`,
    /// `godoc-article:`), or `gosmos-lesson:<uuid>` for a lesson, whatever web address the
    /// text gave for it. So no link the app makes leads to the web, even one that something
    /// other than the app's link handling opened.
    public var url: String {
        switch self {
        case .doc(let link): link.url
        case .lesson(let id): LinkTarget.lessonScheme + ":" + (id.addingPercentEncoding(withAllowedCharacters: LinkTarget.idCharacters) ?? id)
        case .term(let term): LinkTarget.termScheme + ":" + (term.addingPercentEncoding(withAllowedCharacters: LinkTarget.idCharacters) ?? term)
        }
    }

    /// What a lesson's uuid or a term keeps unescaped in its link: either is any text a library
    /// holds.
    private static let idCharacters = CharacterSet(charactersIn: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~")
}

/// Inline code in a lesson's readme that names a standard library package, or a name in one,
/// becomes a link to its documentation: `fmt.Println`, `*strings.Builder`, `net/http`. Code
/// blocks, and code that is a link's text already, stay as they are.
public enum DocLinker {
    public static func link(_ blocks: [MDBlock], names: DocNames) -> [MDBlock] {
        guard !names.isEmpty else { return blocks }
        return blocks.map { link($0, names) }
    }

    static func link(_ block: MDBlock, _ names: DocNames) -> MDBlock {
        switch block {
        case .heading(let level, let tags, let inlines): .heading(level: level, tags: tags, inlines: link(inlines, names))
        case .paragraph(let inlines): .paragraph(link(inlines, names))
        case .quote(let callout, let blocks): .quote(callout: callout, blocks: blocks.map { link($0, names) })
        case .list(let ordered, let start, let tight, let items):
            .list(ordered: ordered, start: start, tight: tight, items: items.map { MDListItem(task: $0.task, blocks: $0.blocks.map { link($0, names) }) })
        case .table(let alignments, let header, let rows):
            .table(alignments: alignments, header: header.map { link($0, names) }, rows: rows.map { $0.map { link($0, names) } })
        case .code, .rule, .html: block
        }
    }

    static func link(_ inlines: [MDInline], _ names: DocNames) -> [MDInline] {
        inlines.map { inline in
            switch inline {
            case .code(let code):
                guard let target = names.link(forCode: code) else { return inline }
                return .link(destination: target.url, title: nil, children: [inline])
            case .emphasis(let c): return .emphasis(link(c, names))
            case .strong(let c): return .strong(link(c, names))
            case .strike(let c): return .strike(link(c, names))
            case .superscript(let c): return .superscript(link(c, names))
            case .text, .link, .image, .lineBreak, .softBreak: return inline
            }
        }
    }
}

/// The first use of each linked term of a course's field guide in a lesson's readme becomes a
/// link to its meaning, so a learner new to the course's world can look any of its words up
/// where they meet it. Terms are matched as the guide spells them, as whole words, also with an
/// ending ("beacons", "stowed", "tracking") and, for a term the guide writes in lowercase, with a
/// capital first letter, as at the start of a sentence ("The Board"); text in code, links and
/// headings stays as it is, and the longer of two terms that start at the same place wins
/// ("slant range" over "range").
public enum TermLinker {
    public static func link(_ blocks: [MDBlock], terms: [CourseStory.Term]) -> [MDBlock] {
        let names = terms.filter(\.linked).map(\.term).filter { !$0.isEmpty }.sorted { $0.count > $1.count }
        guard !names.isEmpty else { return blocks }
        var state = State(names: names, all: Set(names.map { $0.lowercased() }))
        return blocks.map { state.link($0) }
    }

    struct State {
        /// The terms not linked yet, longest first.
        var names: [String]
        /// Every linked term, lowercased, linked yet or not.
        let all: Set<String>

        mutating func link(_ block: MDBlock) -> MDBlock {
            switch block {
            case .paragraph(let inlines): return .paragraph(link(inlines))
            case .quote(let callout, let blocks): return .quote(callout: callout, blocks: blocks.map { link($0) })
            case .list(let ordered, let start, let tight, let items):
                return .list(ordered: ordered, start: start, tight: tight,
                             items: items.map { item in MDListItem(task: item.task, blocks: item.blocks.map { link($0) }) })
            case .table(let alignments, let header, let rows):
                let header = header.map { link($0) }
                return .table(alignments: alignments, header: header, rows: rows.map { $0.map { link($0) } })
            case .heading, .code, .rule, .html: return block
            }
        }

        mutating func link(_ inlines: [MDInline]) -> [MDInline] {
            var out: [MDInline] = []
            for inline in inlines {
                switch inline {
                case .text(let text): out += link(text)
                case .emphasis(let c): out.append(.emphasis(link(c)))
                case .strong(let c): out.append(.strong(link(c)))
                case .strike(let c): out.append(.strike(link(c)))
                case .superscript(let c): out.append(.superscript(link(c)))
                case .code, .link, .image, .lineBreak, .softBreak: out.append(inline)
                }
            }
            return out
        }

        /// The text, its first match of a term not linked yet made a link, and so on in the rest.
        mutating func link(_ text: String) -> [MDInline] {
            guard !names.isEmpty, let (name, range) = firstMatch(in: text) else { return [.text(text)] }
            names.removeAll { $0 == name }
            var out: [MDInline] = []
            if range.lowerBound > text.startIndex { out.append(.text(String(text[..<range.lowerBound]))) }
            out.append(.link(destination: LinkTarget.term(name).url, title: nil, children: [.text(String(text[range]))]))
            if range.upperBound < text.endIndex { out += link(String(text[range.upperBound...])) }
            return out
        }

        /// The earliest whole-word use of a term in `text`, with its ending if it has one.
        func firstMatch(in text: String) -> (String, Range<String.Index>)? {
            var best: (String, Range<String.Index>)?
            for name in names {
                let found = TermLinker.spellings(of: name).compactMap { use(of: $0, as: name, in: text) }.min { $0.lowerBound < $1.lowerBound }
                guard let range = found else { continue }
                // Longest first, so a later name only wins by starting sooner.
                if best == nil || range.lowerBound < best!.1.lowerBound { best = (name, range) }
            }
            return best
        }

        /// The first whole-word use of `spelling` in `text`, passing over a word that with its
        /// ending is another term of the guide: "uplinkd" is never "uplink" with a "d".
        func use(of spelling: String, as name: String, in text: String) -> Range<String.Index>? {
            var from = text.startIndex
            while let range = TermLinker.wholeWord(spelling, in: text, from: from) {
                let word = text[range].lowercased()
                if word == name.lowercased() || !all.contains(word) { return range }
                from = range.upperBound
            }
            return nil
        }
    }

    /// A term as a readme may write it: as the guide spells it and, when the guide writes it in
    /// lowercase, with a capital first letter.
    static func spellings(of term: String) -> [String] {
        guard let first = term.first, first.isLowercase else { return [term] }
        return [term, first.uppercased() + term.dropFirst()]
    }

    /// The endings a term may take in a readme: plurals, and the forms of a verb.
    static let endings = ["ing", "es", "ed", "s", "d"]

    /// The first use of `word` in `text` that is a whole word, also with one of its endings,
    /// from `start` on.
    static func wholeWord(_ word: String, in text: String, from start: String.Index? = nil) -> Range<String.Index>? {
        var from = start ?? text.startIndex
        while let found = text.range(of: word, range: from..<text.endIndex) {
            let before = found.lowerBound == text.startIndex ? nil : text[text.index(before: found.lowerBound)]
            var end = found.upperBound
            for ending in endings where text[end...].hasPrefix(ending) {
                let after = text.index(end, offsetBy: ending.count)
                if after == text.endIndex || !isWordCharacter(text[after]) { end = after; break }
            }
            let after = end == text.endIndex ? nil : text[end]
            if !(before.map(isWordCharacter) ?? false) && !(after.map(isWordCharacter) ?? false) {
                return found.lowerBound..<end
            }
            from = text.index(after: found.lowerBound)
        }
        return nil
    }

    static func isWordCharacter(_ c: Character) -> Bool { c.isLetter || c.isNumber || c == "_" }
}
