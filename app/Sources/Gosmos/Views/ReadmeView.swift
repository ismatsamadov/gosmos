import AppKit
import CodeEditor
import GosmosKit
import SwiftUI

/// A lesson's readme as native text, in a readable column centered in its pane, then the
/// lesson's help (hints and its reference solution), then anything the lesson type adds
/// (the quiz).
struct ReadmeColumn<Footer: View>: View {
    let lesson: LessonModel
    @ViewBuilder var footer: Footer
    @Environment(AppModel.self) private var app
    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    var body: some View {
        ScrollViewReader { proxy in
            ScrollView {
                VStack(alignment: .leading, spacing: 0) {
                    LessonHeader(lesson: lesson)
                        .padding(.bottom, 24)
                    // Not once the lesson is marked not done again: its XP is gone then.
                    if lesson.justPassed && (lesson.isReview || lesson.isDone) {
                        MissionCompletePanel(lesson: lesson)
                            .padding(.bottom, 24)
                    }
                    // A readme that is only its title shows no reading block. Its inline code that
                    // names the standard library links to the documentation.
                    if !lesson.parsed.readme.isEmpty {
                        MarkdownBlocks(blocks: DocLinker.link(lesson.parsed.readme, names: app.docNames))
                    }
                    if !lesson.hints.hints.isEmpty {
                        HintsPanel(state: lesson.hints, lockedReason: lockedReason) { lesson.revealHint() }
                            .padding(.top, 24)
                    }
                    if lesson.hasSolution {
                        SolutionRow(lesson: lesson)
                            .padding(.top, lesson.hints.hints.isEmpty ? 24 : 10)
                    }
                    footer
                        .padding(.top, 24)
                        .id(ReadmeColumnAnchor.footer)
                }
                .frame(maxWidth: Theme.readableWidth, alignment: .leading)
                .padding(.horizontal, 32)
                .padding(.top, 26)
                .padding(.bottom, 40)
                .frame(maxWidth: .infinity)
            }
            // Text scrolled up under the toolbar stops at a hard edge, so it never runs into
            // the window's title.
            .scrollEdgeEffectStyle(.hard, for: .top)
            .onChange(of: lesson.footerRequests) {
                withAnimation(reduceMotion ? nil : .smooth) { proxy.scrollTo(ReadmeColumnAnchor.footer, anchor: .top) }
            }
        }
        .textSelection(.enabled)
        .background(Theme.surface)
        // Links to lessons of this library act here too; the window follows them (AppModel.follow).
        .environment(\.linkPolicy, LinkPolicy(catalog: app.library(lesson.detail.library)?.catalog ?? .empty))
    }

    /// Why the last hint waits: for a failed Submit.
    private var lockedReason: String { "The last hint opens after a Submit that fails." }
}

enum ReadmeColumnAnchor: Hashable { case footer }

/// The lesson's reference solution: a button that shows it beside your code, and what a look
/// does. Before you pass, looking is fine but brings the lesson back for review; after, it
/// changes nothing.
struct SolutionRow: View {
    let lesson: LessonModel

    var body: some View {
        HStack(alignment: .center, spacing: 12) {
            Image(systemName: "doc.text.magnifyingglass")
                .symbolRenderingMode(.hierarchical)
                .foregroundStyle(.secondary)
                .accessibilityHidden(true)
            Text(caption)
                .font(.system(size: 13))
                .foregroundStyle(Theme.secondaryText)
                .fixedSize(horizontal: false, vertical: true)
                .frame(maxWidth: .infinity, alignment: .leading)
            Button("Show Solution") { lesson.showSolution() }
                .secondaryGlassButton()
                .controlSize(.small)
                .help("Show the solution beside your code (⌥⌘S)")
        }
        .padding(.horizontal, 14)
        .padding(.vertical, 12)
        .missionCard(radius: 14)
        .accessibilityElement(children: .contain)
    }

    private var caption: String {
        if lesson.isDone { return "Compare your code with the lesson's solution." }
        if lesson.solutionWasViewed { return "You looked at the solution, so the lesson comes back for review a day after you pass it." }
        return "Stuck? Looking at the solution is fine; the lesson then comes back for review a day after you pass it."
    }
}

/// A compact, open heading keeps the lesson itself above the fold.
struct LessonHeader: View {
    let lesson: LessonModel
    @Environment(AppModel.self) private var app

    var body: some View {
        let detail = lesson.detail
        VStack(alignment: .leading, spacing: 15) {
            HStack(alignment: .firstTextBaseline, spacing: 8) {
                Text(eyebrow)
                    .font(.system(size: 10, weight: .medium))
                    .foregroundStyle(Theme.secondaryText)
                    .fixedSize(horizontal: false, vertical: true)
                Spacer(minLength: 0)
                if !lesson.isReview {
                    Label("\(LearningProgress.lessonXP) XP", systemImage: "bolt.fill")
                        .font(.system(size: 10, weight: .medium)).foregroundStyle(Theme.gold)
                        .fixedSize()
                }
            }
            HStack(alignment: .firstTextBaseline, spacing: 10) {
                Text(detail.title)
                    .font(.system(size: Theme.headingSize(1), weight: .semibold))
                    .tracking(-0.7)
                    .fixedSize(horizontal: false, vertical: true)
                    .accessibilityAddTraits(.isHeader)
                let done = app.isDone(detail.ref)
                Image(systemName: "checkmark.circle.fill")
                    .font(.system(size: 20))
                    .foregroundStyle(.tint)
                    .opacity(done ? 1 : 0)
                    .symbolEffect(.bounce, value: done)
                    .accessibilityLabel("Done")
                    .accessibilityHidden(!done)
            }
            HStack(spacing: 10) {
                // The type always shows. A done lesson adds that it is completed; a redo is
                // practice, whatever the lesson's own done mark says.
                let completed = !lesson.isReview && lesson.isDone
                MissionPill(text: lesson.isReview ? "\(detail.type.label) · Practice" : completed ? "\(detail.type.label) · Completed" : detail.type.label,
                            symbol: completed ? "checkmark" : detail.type.symbol,
                            color: completed ? Theme.success : Theme.accent)
                if !detail.topics.isEmpty {
                    Text(detail.topics.prefix(2).joined(separator: " · "))
                        .font(.system(size: 11)).foregroundStyle(Theme.secondaryText).lineLimit(2)
                }
            }
            if let chapter = app.library(detail.library)?.catalog.chapter(of: detail.id) {
                let done = app.doneCount(chapter, in: detail.library)
                VStack(alignment: .leading, spacing: 8) {
                    MissionProgress(fraction: Double(done) / Double(max(1, chapter.lessons.count)), height: 3)
                        .accessibilityLabel("Chapter progress")
                    Text("\(done) of \(chapter.lessons.count) lessons complete in this chapter")
                        .font(.system(size: 10)).foregroundStyle(Theme.secondaryText)
                }
                .padding(.top, 3)
            }
        }
        .padding(.bottom, 6)
        .frame(maxWidth: .infinity, alignment: .leading)
    }

    private var eyebrow: String {
        let detail = lesson.detail
        var parts = lesson.isReview ? ["Redo for review"] : [String]()
        if let catalog = app.library(detail.library)?.catalog, let loc = catalog.location(of: detail.id) {
            let chapter = catalog.courses[loc.course].chapters[loc.chapter]
            parts.append("Chapter \(loc.chapter + 1) · Lesson \(loc.lesson + 1) of \(chapter.lessons.count)")
        }
        if !detail.summary.required { parts.append("Optional") }
        return parts.joined(separator: "  ·  ")
    }
}

struct MissionCompletePanel: View {
    let lesson: LessonModel
    @Environment(AppModel.self) private var app

    var body: some View {
        // A redo graded by a Submit that failed can pass later, as practice; the review still
        // brings it back, so the panel doesn't say the knowledge held.
        let practiced = lesson.isReview && lesson.code?.reviewPassed == false
        VStack(alignment: .leading, spacing: 14) {
            HStack(spacing: 12) {
                MissionIcon(symbol: lesson.isReview ? "sparkles" : "checkmark.seal.fill", color: Theme.success)
                VStack(alignment: .leading, spacing: 4) {
                    Text(!lesson.isReview ? "Mission accomplished!" : practiced ? "Solved." : "Skills refreshed.")
                        .font(.system(size: 17, weight: .semibold))
                    // True on a first pass and on a pass of a lesson already done, which adds no XP.
                    Text(!lesson.isReview ? "This lesson’s \(LearningProgress.lessonXP) XP is part of your journey. Keep exploring."
                         : practiced ? "It comes back for review in a day, to make it stick." : "That knowledge is getting stronger.")
                        .font(.system(size: 12)).foregroundStyle(Theme.secondaryText)
                }
            }
            if lesson.isReview {
                Button("Back to Review", systemImage: "arrow.right") { app.showReview() }
                    .prominentGlassButton()
            } else if app.nextLesson != nil {
                Button("Next mission", systemImage: "arrow.right") { app.openNext() }
                    .prominentGlassButton()
            } else {
                Button("Explore your progress", systemImage: "square.grid.2x2") { app.showDashboard() }
                    .prominentGlassButton()
            }
        }
        .padding(18)
        .frame(maxWidth: .infinity, alignment: .leading)
        .missionCard(tint: Theme.success)
        .accessibilityElement(children: .contain)
    }
}

/// The readme's blocks.
struct MarkdownBlocks: View {
    let blocks: [MDBlock]

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            ForEach(Array(blocks.enumerated()), id: \.offset) { _, block in
                BlockView(block: block)
            }
        }
    }
}

struct BlockView: View {
    let block: MDBlock
    @Environment(\.linkPolicy) private var links

    var body: some View {
        switch block {
        case .heading(let level, let tags, let inlines):
            HStack(alignment: .firstTextBaseline, spacing: 6) {
                ForEach(tags, id: \.self) { tag in
                    Text(tag)
                        .font(.system(size: 11, weight: .semibold))
                        .padding(.horizontal, 7)
                        .padding(.vertical, 2)
                        .background(.quaternary, in: Capsule())
                        .foregroundStyle(Theme.secondaryText)
                }
                Text(Inline.attributed(inlines, size: Theme.headingSize(level), weight: .semibold, links: links))
                    .fixedSize(horizontal: false, vertical: true)
                    .accessibilityAddTraits(.isHeader)
            }
            .padding(.top, level <= 2 ? 14 : 6)
        case .paragraph(let inlines):
            Text(Inline.attributed(inlines, size: Theme.bodySize, links: links))
                .lineSpacing(Theme.bodySize * 0.45)
                .fixedSize(horizontal: false, vertical: true)
                .frame(maxWidth: Theme.proseWidth, alignment: .leading)
        case .code(let language, let text):
            CodeBlockView(language: language, code: text)
        case .quote(let callout, let blocks):
            if let callout {
                CalloutView(kind: callout, blocks: blocks)
            } else {
                QuoteView(blocks: blocks)
            }
        case .list(let ordered, let start, let tight, let items):
            ListBlockView(ordered: ordered, start: start, tight: tight, items: items)
        case .table(let alignments, let header, let rows):
            TableBlockView(alignments: alignments, header: header, rows: rows)
        case .rule:
            Divider().padding(.vertical, 6)
        case .html(let lines):
            VStack(alignment: .leading, spacing: 6) {
                ForEach(Array(lines.enumerated()), id: \.offset) { _, line in
                    Text(line).font(.system(size: Theme.bodySize)).fixedSize(horizontal: false, vertical: true)
                }
            }
        }
    }
}

/// Which links in Markdown act, from the environment: links into the Go documentation, and in
/// a lesson's text links to the lessons of its library (LinkTarget). Any other address is plain
/// text, neither colored nor clickable: Gosmos opens nothing outside itself.
struct LinkPolicy: Sendable {
    /// The library whose lessons links may open.
    var catalog: Catalog = .empty

    /// The URL a link to `destination` carries, or nil when it does nothing: the app's own
    /// (LinkTarget.url), never the web address a link to a lesson gave.
    func url(_ destination: String) -> URL? {
        LinkTarget(destination: destination, isLesson: { catalog.summary($0) != nil }).flatMap { URL(string: $0.url) }
    }
}

private struct LinkPolicyKey: EnvironmentKey {
    static let defaultValue = LinkPolicy()
}

extension EnvironmentValues {
    var linkPolicy: LinkPolicy {
        get { self[LinkPolicyKey.self] }
        set { self[LinkPolicyKey.self] = newValue }
    }
}

/// Inline Markdown as styled text.
enum Inline {
    static func attributed(_ inlines: [MDInline], size: CGFloat, weight: Font.Weight = .regular, links: LinkPolicy = LinkPolicy()) -> AttributedString {
        var out = AttributedString()
        func walk(_ items: [MDInline], bold: Bool, italic: Bool, strike: Bool, link: URL?, raised: Bool = false) {
            for item in items {
                switch item {
                case .text(let s):
                    var a = AttributedString(s)
                    var font = Font.system(size: raised ? (size * 0.72).rounded() : size, weight: bold ? .semibold : weight)
                    if italic { font = font.italic() }
                    a.font = font
                    if raised { a.baselineOffset = (size * 0.38).rounded() }
                    if strike { a.strikethroughStyle = .single }
                    if let link { a.link = link }
                    out += a
                case .code(let s):
                    // A subtle capsule: a narrow no-break space on each side inside a faint
                    // background, so the padding never wraps away from the code.
                    var a = AttributedString("\u{202F}" + s + "\u{202F}")
                    a.font = .system(size: (size * 0.88).rounded(), weight: bold ? .semibold : .regular, design: .monospaced)
                    a.backgroundColor = Color.secondary.opacity(0.13)
                    if let link { a.link = link }
                    out += a
                case .emphasis(let c):
                    walk(c, bold: bold, italic: true, strike: strike, link: link, raised: raised)
                case .strong(let c):
                    walk(c, bold: true, italic: italic, strike: strike, link: link, raised: raised)
                case .strike(let c):
                    walk(c, bold: bold, italic: italic, strike: true, link: link, raised: raised)
                case .superscript(let c):
                    walk(c, bold: bold, italic: italic, strike: strike, link: link, raised: true)
                case .link(let dest, _, let c):
                    walk(c, bold: bold, italic: italic, strike: strike, link: links.url(dest) ?? link, raised: raised)
                case .image:
                    break // no pictures or videos
                case .lineBreak:
                    out += AttributedString("\n")
                case .softBreak:
                    out += AttributedString(" ")
                }
            }
        }
        walk(inlines, bold: false, italic: false, strike: false, link: nil)
        return out
    }
}

/// SF Mono 13 on a rounded panel, Go in the editor's Xcode colors, with the block's
/// language and a Copy button on hover. Long lines wrap, so nothing is cut off at the edge.
struct CodeBlockView: View {
    let language: String
    let code: String
    @ViewState private var hovering = false
    @ViewState private var copied = false

    var body: some View {
        WrappingText(text: CodeText.attributed(code, language: language, size: Theme.codeSize),
                     accessibilityLabel: language.isEmpty ? "Code" : "\(language) code")
            .padding(.horizontal, 14)
            .padding(.vertical, 12)
            .frame(maxWidth: .infinity, alignment: .leading)
            .background(.quaternary.opacity(0.6), in: RoundedRectangle(cornerRadius: 10, style: .continuous))
            .overlay(alignment: .topTrailing) {
                // On hover only, on glass: at rest nothing covers the end of a long first line.
                if hovering {
                    GlassEffectContainer(spacing: 6) {
                        HStack(spacing: 6) {
                            if !language.isEmpty {
                                Text(language)
                                    .font(.system(size: 11, weight: .medium, design: .monospaced))
                                    .foregroundStyle(Theme.secondaryText)
                                    .padding(.horizontal, 8)
                                    .padding(.vertical, 4)
                                    .glassEffect(.regular, in: Capsule())
                            }
                            Button {
                                NSPasteboard.general.clearContents()
                                NSPasteboard.general.setString(code, forType: .string)
                                copied = true
                                Task {
                                    try? await Task.sleep(for: .seconds(1.5))
                                    copied = false
                                }
                            } label: {
                                Label(copied ? "Copied" : "Copy", systemImage: copied ? "checkmark" : "doc.on.doc")
                                    .contentTransition(.symbolEffect(.replace))
                            }
                            .secondaryGlassButton()
                            .controlSize(.small)
                            .accessibilityLabel("Copy code")
                        }
                    }
                    .padding(8)
                }
            }
            .onHover { hovering = $0 }
    }
}

/// Code as attributed text for `WrappingText`: SF Mono in the editor's colors, each line
/// hanging its continuation lines under its own indentation plus two spaces.
enum CodeText {
    /// Highlighted code by its text, language and size. A view's body runs again on every hover
    /// of a code block, and a block of the documentation runs to 48 KB (syscall's constants):
    /// it is highlighted once. Its colors are dynamic and follow the appearance by themselves.
    @MainActor private static let highlighted: NSCache<NSString, NSAttributedString> = {
        let cache = NSCache<NSString, NSAttributedString>()
        cache.totalCostLimit = 8_000_000 // UTF-16 units of code
        return cache
    }()

    /// The code highlighted, the same text object each time for the same code, language and size.
    @MainActor static func attributed(_ code: String, language: String, size: CGFloat) -> NSAttributedString {
        let key = "\(language) \(size)\n\(code)" as NSString
        if let kept = highlighted.object(forKey: key) { return kept }
        let theme = SyntaxTheme.xcode
        let result = NSMutableAttributedString(string: code, attributes: theme.attributes(for: .plain, size: size))
        hangLines(of: result, size: size)
        let length = result.length
        for token in Highlighting.forLanguage(language).tokens(in: code) {
            let range = token.range
            guard range.location != NSNotFound, range.length > 0, NSMaxRange(range) <= length else { continue }
            result.addAttributes(theme.attributes(for: token.kind, size: size), range: range)
        }
        let text = NSAttributedString(attributedString: result)
        highlighted.setObject(text, forKey: key, cost: length)
        return text
    }

    /// Gives every line a paragraph style: tabs four spaces wide, and continuation lines that
    /// start under the line's indentation, past its first `marker` characters, plus two spaces.
    static func hangLines(of text: NSMutableAttributedString, size: CGFloat, marker: Int = 0) {
        let font = SyntaxTheme.xcode.font(for: .plain, size: size)
        let space = (" " as NSString).size(withAttributes: [.font: font]).width
        let tab = space * 4
        let string = text.string as NSString
        var location = 0
        repeat {
            let line = string.paragraphRange(for: NSRange(location: location, length: 0))
            var x = space * CGFloat(min(marker, line.length))
            var index = line.location + min(marker, line.length)
            while index < NSMaxRange(line) {
                let character = string.character(at: index)
                if character == 0x09 {
                    x = (floor(x / tab) + 1) * tab
                } else if character == 0x20 {
                    x += space
                } else {
                    break
                }
                index += 1
            }
            let style = NSMutableParagraphStyle()
            style.tabStops = []
            style.defaultTabInterval = tab
            style.lineSpacing = 3
            // Deep indentation still leaves room for the text beside it.
            style.headIndent = min(x + 2 * space, 24 * space)
            text.addAttribute(.paragraphStyle, value: style, range: line)
            location = NSMaxRange(line)
        } while location < string.length
    }
}

/// Read-only text that wraps to the width it is given and grows to fit, selectable like the
/// rest of the readme. TextKit 1 lays the few lines out at once, so the height is exact.
struct WrappingText: NSViewRepresentable {
    let text: NSAttributedString
    var accessibilityLabel: String?

    /// The text object the view shows, so an update with that very object, such as highlighted
    /// code CodeText kept, isn't compared with it run by run.
    final class Coordinator {
        var shown: NSAttributedString?
    }

    func makeCoordinator() -> Coordinator { Coordinator() }

    func makeNSView(context: Context) -> WrappingTextView {
        let view = WrappingTextView(usingTextLayoutManager: false)
        view.isEditable = false
        view.isSelectable = true
        view.drawsBackground = false
        view.textContainerInset = .zero
        view.textContainer?.lineFragmentPadding = 0
        view.textContainer?.widthTracksTextView = false
        view.isVerticallyResizable = false
        view.isHorizontallyResizable = false
        view.textStorage?.setAttributedString(text)
        context.coordinator.shown = text
        view.setAccessibilityLabel(accessibilityLabel)
        return view
    }

    func updateNSView(_ view: WrappingTextView, context: Context) {
        if context.coordinator.shown !== text {
            if view.textStorage?.isEqual(to: text) != true {
                view.textStorage?.setAttributedString(text)
                view.invalidateIntrinsicContentSize()
            }
            context.coordinator.shown = text
        }
        view.setAccessibilityLabel(accessibilityLabel)
    }

    func sizeThatFits(_ proposal: ProposedViewSize, nsView view: WrappingTextView, context: Context) -> CGSize? {
        let width = proposal.width.flatMap { $0.isFinite && $0 > 0 ? $0 : nil } ?? view.widestLine
        return CGSize(width: width, height: view.height(forWidth: width))
    }
}

final class WrappingTextView: NSTextView {
    /// The height of the text laid out `width` wide.
    func height(forWidth width: CGFloat) -> CGFloat {
        guard let container = textContainer, let layoutManager else { return 0 }
        fit(container, to: width)
        layoutManager.ensureLayout(for: container)
        return ceil(layoutManager.usedRect(for: container).height)
    }

    /// The width of the longest line, unwrapped.
    var widestLine: CGFloat {
        guard let storage = textStorage else { return 0 }
        return ceil(storage.boundingRect(with: NSSize(width: CGFloat.greatestFiniteMagnitude, height: .greatestFiniteMagnitude),
                                         options: [.usesLineFragmentOrigin]).width)
    }

    override func setFrameSize(_ newSize: NSSize) {
        super.setFrameSize(newSize)
        if let container = textContainer { fit(container, to: newSize.width) }
    }

    private func fit(_ container: NSTextContainer, to width: CGFloat) {
        if container.size.width != width {
            container.size = NSSize(width: max(width, 1), height: .greatestFiniteMagnitude)
        }
    }
}

struct CalloutView: View {
    let kind: String
    let blocks: [MDBlock]

    var body: some View {
        let color = Theme.calloutColor(kind)
        HStack(alignment: .top, spacing: 10) {
            Image(systemName: Readme.calloutSymbol(kind))
                .symbolRenderingMode(.hierarchical)
                .font(.system(size: 16, weight: .semibold))
                .foregroundStyle(color)
                .frame(width: 20)
            VStack(alignment: .leading, spacing: 8) {
                Text(Readme.calloutTitle(kind))
                    .font(.system(size: 14, weight: .semibold))
                    .foregroundStyle(color)
                ForEach(Array(blocks.enumerated()), id: \.offset) { _, block in
                    AnyView(BlockView(block: block))
                }
            }
        }
        .padding(14)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(color.opacity(0.1), in: RoundedRectangle(cornerRadius: 10, style: .continuous))
    }
}

struct QuoteView: View {
    let blocks: [MDBlock]

    var body: some View {
        HStack(alignment: .top, spacing: 12) {
            RoundedRectangle(cornerRadius: 1.5)
                .fill(.quaternary)
                .frame(width: 3)
            VStack(alignment: .leading, spacing: 8) {
                ForEach(Array(blocks.enumerated()), id: \.offset) { _, block in
                    AnyView(BlockView(block: block))
                }
            }
            .foregroundStyle(Theme.secondaryText)
        }
        .fixedSize(horizontal: false, vertical: true)
    }
}

struct ListBlockView: View {
    let ordered: Bool
    let start: Int
    let tight: Bool
    let items: [MDListItem]

    var body: some View {
        VStack(alignment: .leading, spacing: tight ? 5 : 10) {
            ForEach(Array(items.enumerated()), id: \.offset) { i, item in
                HStack(alignment: .firstTextBaseline, spacing: 8) {
                    marker(i, item)
                        .frame(minWidth: ordered ? 22 : 12, alignment: .trailing)
                    VStack(alignment: .leading, spacing: 8) {
                        ForEach(Array(item.blocks.enumerated()), id: \.offset) { _, block in
                            AnyView(BlockView(block: block))
                        }
                    }
                }
            }
        }
        .padding(.leading, 2)
    }

    @ViewBuilder
    private func marker(_ i: Int, _ item: MDListItem) -> some View {
        if let ticked = item.task {
            Image(systemName: ticked ? "checkmark.square.fill" : "square")
                .foregroundStyle(ticked ? AnyShapeStyle(.tint) : AnyShapeStyle(.secondary))
                .font(.system(size: Theme.bodySize - 1))
                .accessibilityLabel(ticked ? "Done" : "To do")
        } else if ordered {
            Text("\(start + i).")
                .font(.system(size: Theme.bodySize).monospacedDigit())
                .foregroundStyle(Theme.secondaryText)
        } else {
            Text("•")
                .font(.system(size: Theme.bodySize, weight: .bold))
                .foregroundStyle(Theme.secondaryText)
        }
    }
}

/// Hairline rules and zebra rows. Columns share the readme's width and their text wraps; only
/// when the columns cannot fit even at their narrowest does the table scroll sideways.
struct TableBlockView: View {
    let alignments: [MDAlignment]
    let header: [[MDInline]]
    let rows: [[[MDInline]]]
    @ViewState private var available: CGFloat = 0
    @Environment(\.linkPolicy) private var links

    var body: some View {
        let layout = TableColumns(header: header, rows: rows, available: available)
        Group {
            if layout.fits {
                grid(layout.widths)
            } else {
                ScrollView(.horizontal) {
                    grid(layout.widths)
                        .padding(.bottom, 10)
                }
                .scrollIndicators(.visible)
            }
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .onGeometryChange(for: CGFloat.self) { $0.size.width } action: { available = $0 }
    }

    private func grid(_ widths: [CGFloat]) -> some View {
        Grid(alignment: .leading, horizontalSpacing: 0, verticalSpacing: 0) {
            GridRow {
                ForEach(Array(header.enumerated()), id: \.offset) { c, cell in
                    cellView(cell, column: c, width: widths[safe: c], header: true, zebra: false)
                }
            }
            Divider()
            ForEach(Array(rows.enumerated()), id: \.offset) { r, row in
                GridRow {
                    ForEach(Array(row.enumerated()), id: \.offset) { c, cell in
                        cellView(cell, column: c, width: widths[safe: c], header: false, zebra: r % 2 == 1)
                    }
                }
                if r < rows.count - 1 { Divider().opacity(0.6) }
            }
        }
        .overlay(RoundedRectangle(cornerRadius: 8, style: .continuous).strokeBorder(.separator, lineWidth: 0.5))
        .clipShape(RoundedRectangle(cornerRadius: 8, style: .continuous))
    }

    private func cellView(_ cell: [MDInline], column: Int, width: CGFloat?, header: Bool, zebra: Bool) -> some View {
        let alignment: Alignment = switch column < alignments.count ? alignments[column] : .none {
        case .center: .center
        case .right: .trailing
        default: .leading
        }
        return Text(Inline.attributed(cell, size: TableColumns.fontSize, weight: header ? .semibold : .regular, links: links))
            .multilineTextAlignment(alignment == .trailing ? .trailing : alignment == .center ? .center : .leading)
            .fixedSize(horizontal: false, vertical: true)
            .padding(.horizontal, TableColumns.padding)
            .padding(.vertical, 7)
            .frame(width: width, alignment: alignment)
            .frame(maxHeight: .infinity, alignment: .top)
            .background(header ? AnyShapeStyle(.quaternary.opacity(0.5)) : zebra ? AnyShapeStyle(.quaternary.opacity(0.25)) : AnyShapeStyle(.clear))
    }
}

/// Column widths for a table: each column's text measured at its narrowest (its longest
/// word) and its widest (one line), and the available width shared out between them.
struct TableColumns {
    static let fontSize: CGFloat = 14
    static let padding: CGFloat = 12
    /// A column is never narrower than this, so short headers are not squeezed to a letter.
    static let minimumColumn: CGFloat = 64
    /// A single long word (a URL, a long identifier) wraps rather than widening its column
    /// past this.
    static let longestWord: CGFloat = 200
    /// How wide a column may grow in a table that scrolls sideways anyway.
    static let scrollingColumn: CGFloat = 240

    let widths: [CGFloat]
    let fits: Bool

    init(header: [[MDInline]], rows: [[[MDInline]]], available: CGFloat) {
        let count = max(header.count, rows.map(\.count).max() ?? 0)
        var narrowest = [CGFloat](repeating: Self.minimumColumn, count: count)
        var widest = [CGFloat](repeating: Self.minimumColumn, count: count)
        for (r, row) in ([header] + rows).enumerated() {
            let font = NSFont.systemFont(ofSize: Self.fontSize, weight: r == 0 ? .semibold : .regular)
            for (c, cell) in row.enumerated() where c < count {
                let text = Readme.plainText(cell)
                let line = Self.width(text, font) + 2 * Self.padding
                let word = min(text.split(whereSeparator: \.isWhitespace).map { Self.width(String($0), font) }.max() ?? 0, Self.longestWord)
                widest[c] = max(widest[c], line)
                narrowest[c] = max(narrowest[c], word + 2 * Self.padding)
            }
        }
        let widestTotal = widest.reduce(0, +), narrowestTotal = narrowest.reduce(0, +)
        if available <= 0 || narrowestTotal > available {
            // Too many columns to fit even at their narrowest: each gets room for a few words
            // and the table scrolls sideways.
            widths = zip(narrowest, widest).map { low, high in max(low, min(high, Self.scrollingColumn)) }
            fits = false
        } else if widestTotal <= available {
            widths = widest
            fits = true
        } else {
            // Every column gets its narrowest width plus a share of the rest, in proportion
            // to how much more it would take to show its text on one line.
            let spare = available - narrowestTotal, wanted = widestTotal - narrowestTotal
            widths = zip(narrowest, widest).map { low, high in (low + (high - low) * spare / wanted).rounded(.down) }
            fits = true
        }
    }

    private static func width(_ text: String, _ font: NSFont) -> CGFloat {
        ceil((text as NSString).size(withAttributes: [.font: font]).width)
    }
}

extension Array {
    subscript(safe index: Int) -> Element? { indices.contains(index) ? self[index] : nil }
}
