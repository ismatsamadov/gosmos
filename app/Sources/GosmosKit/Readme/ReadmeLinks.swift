import Foundation

/// What a link in a readme or in the documentation does. Gosmos opens nothing outside itself:
/// a link into the Go documentation opens it in the app, and a link whose path names a lesson
/// of the readme's library opens that lesson. Any other address is shown as plain text.
public enum LinkTarget: Hashable, Sendable {
    case doc(DocLink)
    /// A lesson of the library the text belongs to, by uuid.
    case lesson(String)

    /// The target of a link's destination, as the Markdown gives it; nil for one that does
    /// nothing. `isLesson` tells whether a uuid names a lesson of the text's library.
    public init?(destination: String, isLesson: (String) -> Bool = { _ in false }) {
        let destination = destination.trimmingCharacters(in: .whitespaces)
        if let doc = DocLink(string: destination) {
            self = .doc(doc)
        } else if let url = URL(string: destination), let lesson = url.pathComponents.last(where: isLesson) {
            self = .lesson(lesson)
        } else {
            return nil
        }
    }
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
