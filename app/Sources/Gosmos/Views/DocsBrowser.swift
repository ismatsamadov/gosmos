import AppKit
import GosmosKit
import SwiftUI

/// The Go documentation, in the detail area: its list (the language, the standard library,
/// or what its search found) beside the page.
struct DocsBrowser: View {
    let docs: DocsModel

    var body: some View {
        PersistentSplit(.horizontal, key: "DocsSplitList", defaultFraction: 0.3, minFirst: 200, minSecond: 380) {
            DocsList(docs: docs)
        } second: {
            DocPageView(docs: docs)
                .background(Theme.canvas)
        }
    }
}

/// The documentation's contents, or its search results while the field above them has text.
struct DocsList: View {
    let docs: DocsModel
    @Environment(AppModel.self) private var app
    /// Open groups of packages, by first path element.
    @ViewState private var expanded: Set<String> = []

    var body: some View {
        @Bindable var docs = docs
        VStack(spacing: 0) {
            HStack(spacing: 9) {
                MissionIcon(symbol: "text.book.closed", color: Theme.violet, size: 30)
                VStack(alignment: .leading, spacing: 2) {
                    Text("Field guide").font(.system(size: 14, weight: .semibold))
                    Text("GO DOCUMENTATION").font(.system(size: 8, weight: .bold)).tracking(1.2)
                        .foregroundStyle(Theme.secondaryText)
                }
                Spacer(minLength: 0)
            }
            .padding(.horizontal, 14)
            .padding(.top, 16)
            .padding(.bottom, 4)
            DocsSearchField(text: $docs.searchText, prompt: "Search Documentation") {
                if let first = docs.results.first { app.openDocumentation(first.link) }
            }
            .padding(10)
            // ⌘F searches the documentation: nothing on screen here takes Find.
            .background {
                Button("Search Documentation") {
                    NotificationCenter.default.post(name: DocsSearchField.focusRequest, object: nil)
                }
                .keyboardShortcut("f", modifiers: .command)
                .opacity(0)
                .focusable(false)
                .accessibilityHidden(true)
            }
            Divider()
            if !docs.searchText.trimmingCharacters(in: .whitespaces).isEmpty {
                results
            } else {
                contents
            }
        }
        .background(Theme.surface)
        .onChange(of: docs.location, initial: true) { openGroup(of: docs.location) }
    }

    /// What a row chosen in a list opens, on the next turn of the main actor, after the list
    /// has finished its selection callback (as in the sidebar).
    private func selection(_ current: DocLink?) -> Binding<DocLink?> {
        Binding(
            get: { current },
            set: { link in
                guard let link, link != current else { return }
                Task { @MainActor in app.openDocumentation(link) }
            })
    }

    private var contents: some View {
        List(selection: selection(docs.location?.page)) {
            Section("Language") {
                ForEach(docs.outline.articles) { article in
                    row(article.title, symbol: "doc.text")
                        .tag(article.link)
                }
            }
            Section("Standard Library") {
                ForEach(docs.outline.groups) { group in
                    if group.packages.count == 1, let package = group.packages.first {
                        packageRow(package, in: group)
                    } else {
                        DisclosureGroup(isExpanded: Binding(get: { expanded.contains(group.top) },
                                                            set: { if $0 { expanded.insert(group.top) } else { expanded.remove(group.top) } })) {
                            ForEach(group.packages) { packageRow($0, in: group) }
                        } label: {
                            Text(group.top)
                        }
                    }
                }
            }
            Section {
                row("Acknowledgements", symbol: "heart.text.square")
                    .tag(DocLink.acknowledgements)
            }
        }
        .listStyle(.sidebar)
        .scrollContentBackground(.hidden)
    }

    /// An article or the acknowledgements: a symbol and a title, as the sidebar's rows.
    private func row(_ title: String, symbol: String) -> some View {
        Label {
            Text(title)
                .foregroundStyle(.primary)
                .lineLimit(1)
        } icon: {
            Image(systemName: symbol)
                .symbolRenderingMode(.hierarchical)
                .foregroundStyle(.secondary)
        }
    }

    /// A package by its path below its group's first element: `aes` under crypto.
    private func packageRow(_ package: DocOutline.Package, in group: DocOutline.Group) -> some View {
        let name = package.path == group.top ? package.path : String(package.path.dropFirst(group.top.count + 1))
        return HStack(spacing: 6) {
            Text(name)
                .lineLimit(1)
                .foregroundStyle(package.deprecated ? AnyShapeStyle(Theme.secondaryText) : AnyShapeStyle(.primary))
            if package.deprecated { DeprecatedTag() }
        }
        .help(package.synopsis ?? package.path)
        .accessibilityLabel(package.path + (package.deprecated ? ", deprecated" : ""))
        .tag(package.link)
    }

    private var results: some View {
        Group {
            if docs.results.isEmpty && !docs.searching {
                // Filling the column, so the search field above it stays at the top.
                ContentUnavailableView.search(text: docs.searchText)
                    .frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                List(selection: selection(docs.location)) {
                    ForEach(docs.results) { hit in
                        DocHitRow(hit: hit)
                            .tag(hit.link)
                    }
                }
                .listStyle(.sidebar)
                .scrollContentBackground(.hidden)
                // A new list for each set of results, as in the sidebar's search.
                .id(docs.results.map(\.id))
            }
        }
    }

    /// Opens the group the page on screen is in.
    private func openGroup(of link: DocLink?) {
        guard case .package(let path)? = link?.page else { return }
        expanded.insert(String(path.prefix { $0 != "/" }))
    }
}

/// The system's search field: Escape clears it, Return opens the first result.
struct DocsSearchField: NSViewRepresentable {
    @Binding var text: String
    let prompt: String
    let onSubmit: () -> Void

    /// Posted to put the insertion point in the field (⌘F in the browser, Edit › Search).
    static let focusRequest = Notification.Name("GosmosDocsSearchFocusRequest")

    func makeNSView(context: Context) -> NSSearchField {
        let field = NSSearchField()
        field.placeholderString = prompt
        field.delegate = context.coordinator
        field.sendsSearchStringImmediately = true
        field.setAccessibilityLabel(prompt)
        context.coordinator.field = field
        NotificationCenter.default.addObserver(context.coordinator, selector: #selector(Coordinator.focus),
                                               name: Self.focusRequest, object: nil)
        return field
    }

    func updateNSView(_ field: NSSearchField, context: Context) {
        context.coordinator.parent = self
        if field.stringValue != text { field.stringValue = text }
    }

    func makeCoordinator() -> Coordinator { Coordinator(self) }

    @MainActor
    final class Coordinator: NSObject, NSSearchFieldDelegate {
        var parent: DocsSearchField
        weak var field: NSSearchField?

        init(_ parent: DocsSearchField) { self.parent = parent }

        @objc func focus() {
            guard let field, let window = field.window else { return }
            window.makeFirstResponder(field)
            field.currentEditor()?.selectAll(nil)
        }

        func controlTextDidChange(_ notification: Notification) {
            guard let field = notification.object as? NSSearchField else { return }
            parent.text = field.stringValue
        }

        func control(_ control: NSControl, textView: NSTextView, doCommandBy selector: Selector) -> Bool {
            guard selector == #selector(NSResponder.insertNewline(_:)) else { return false }
            parent.onSubmit()
            return true
        }
    }
}

/// A place search found in the documentation: a package, a name, or an article's section.
struct DocHitRow: View {
    let hit: DocSearchHit
    @Environment(\.backgroundProminence) private var prominence

    var body: some View {
        HStack(alignment: .firstTextBaseline, spacing: 7) {
            Image(systemName: symbol)
                .symbolRenderingMode(.hierarchical)
                .foregroundStyle(.secondary)
                .frame(width: 18)
            VStack(alignment: .leading, spacing: 2) {
                HStack(spacing: 6) {
                    Text(hit.title)
                        .font(hit.kind == .section ? .body : .system(.body, design: .monospaced))
                        .foregroundStyle(.primary)
                        .lineLimit(1)
                        .truncationMode(.middle)
                    if hit.deprecated { DeprecatedTag() }
                }
                if let detail = hit.detail {
                    Text(detail)
                        .font(.caption)
                        .foregroundStyle(prominence == .increased ? AnyShapeStyle(.secondary) : AnyShapeStyle(Theme.secondaryText))
                        .lineLimit(1)
                }
            }
        }
        .padding(.vertical, 2)
        .accessibilityElement(children: .combine)
    }

    private var symbol: String {
        switch hit.kind {
        case .package: "shippingbox"
        case .section: "doc.text"
        case .symbol(let kind): DocDeclarationView.symbol(kind)
        }
    }
}

/// A quiet mark for a deprecated package or name.
struct DeprecatedTag: View {
    var body: some View {
        Text("Deprecated")
            .font(.system(size: 11, weight: .medium))
            .padding(.horizontal, 7)
            .padding(.vertical, 2)
            .background(.quaternary, in: Capsule())
            .foregroundStyle(Theme.secondaryText)
    }
}

// MARK: - Pages

/// The page on screen, or why there is none.
struct DocPageView: View {
    let docs: DocsModel

    var body: some View {
        if let page = docs.page {
            DocPageContent(page: page, docs: docs)
                // A new page starts at its top, with nothing of the last one's scrolling.
                .id(page.link)
        } else if let error = docs.loadError {
            ContentUnavailableView {
                Label("This Page Can't Be Shown", systemImage: "exclamationmark.triangle")
            } description: {
                Text(error).foregroundStyle(Theme.secondaryText)
            }
        } else if docs.location == nil {
            DocsWelcome(docs: docs)
        } else {
            ProgressView().appearsAfterDelay().frame(maxWidth: .infinity, maxHeight: .infinity)
        }
    }
}

private struct DocsWelcome: View {
    let docs: DocsModel
    @Environment(AppModel.self) private var app

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 24) {
                MissionHeading(eyebrow: "YOUR OFFLINE FIELD GUIDE", title: "Curiosity, meet clarity.", subtitle: "The language, the guides, and the Go standard library. A good answer is always within reach.", symbol: "text.book.closed", color: Theme.violet)
                    .padding(24)
                    .frame(maxWidth: .infinity, alignment: .leading)
                    .missionCard(tint: Theme.violet, radius: 22)
                Text("Explore the language").font(.system(size: 18, weight: .semibold))
                ForEach(docs.outline.articles) { article in
                    Button { app.openDocumentation(article.link) } label: {
                        HStack(spacing: 14) {
                            MissionIcon(symbol: "doc.text", color: Theme.violet, size: 34)
                            Text(article.title).font(.system(size: 14, weight: .medium))
                                .frame(maxWidth: .infinity, alignment: .leading)
                            Image(systemName: "arrow.up.right").font(.system(size: 12)).foregroundStyle(Theme.violet)
                        }
                        .padding(16)
                        .missionCard(radius: 14)
                        .contentShape(RoundedRectangle(cornerRadius: 14))
                    }
                    .buttonStyle(.plain)
                }
                Label("Looking for a package or function? Search the field guide.", systemImage: "magnifyingglass")
                    .font(.system(size: 12)).foregroundStyle(Theme.secondaryText)
            }
            .padding(30)
            .frame(maxWidth: Theme.readableWidth)
            .frame(maxWidth: .infinity)
        }
    }
}

/// One row of a page, as the lazy stack lays them out: pages are long (the specification, or
/// a package like net/http), so only the rows in view are made. A declaration's or a
/// section's row has its anchor as its id; the page's own rows start with "#".
enum DocRow: Identifiable {
    case header(eyebrow: String, symbol: String, title: String, subtitle: String?, deprecated: Bool)
    case block(id: String, MDBlock)
    case heading(id: String, String)
    case index(package: String, DocIndexEntry)
    case declaration(DocDeclaration)
    case example(DocExample)
    case contents(slug: String, [DocArticle.Section])
    case section(DocArticle.Section)
    case license(String)

    var id: String {
        switch self {
        case .header: "#top"
        case .block(let id, _), .heading(let id, _): id
        case .index(_, let entry): "#index/" + entry.anchor
        case .declaration(let d): d.id
        case .example(let e): "#example/\(e.id)"
        case .contents: "#contents"
        case .section(let s): s.anchor
        case .license: "#license"
        }
    }

    static func rows(_ page: DocPage) -> [DocRow] {
        switch page {
        case .package(let p):
            var rows: [DocRow] = [.header(eyebrow: "Package", symbol: "shippingbox", title: p.path, subtitle: "import \"\(p.path)\"", deprecated: p.deprecated)]
            rows += p.doc.enumerated().map { .block(id: "#doc/\($0.offset)", $0.element) }
            rows += p.examples.map(DocRow.example)
            if !p.index.isEmpty {
                rows.append(.heading(id: "#index", "Index"))
                rows += p.index.map { .index(package: p.path, $0) }
            }
            for section in p.sections {
                rows.append(.heading(id: "#section/" + section.title, section.title))
                rows += section.declarations.map(DocRow.declaration)
            }
            return rows
        case .article(let a):
            var rows: [DocRow] = [.header(eyebrow: "Language", symbol: "doc.text", title: a.title, subtitle: a.subtitle, deprecated: false)]
            if a.contents.count > 1 { rows.append(.contents(slug: a.slug, a.contents)) }
            rows += a.lead.enumerated().map { .block(id: "#lead/\($0.offset)", $0.element) }
            rows += a.sections.map(DocRow.section)
            return rows
        case .acknowledgements(let credits):
            var rows: [DocRow] = [.header(eyebrow: "Go Documentation", symbol: "heart.text.square", title: "Acknowledgements",
                                          subtitle: "From \(credits.goVersion)", deprecated: false)]
            rows += credits.credit.enumerated().map { .block(id: "#credit/\($0.offset)", $0.element) }
            rows.append(.heading(id: "#license-heading", "Go's License"))
            rows.append(.license(credits.license))
            return rows
        }
    }
}

/// A page in the readable column, scrolled to the place on screen.
struct DocPageContent: View {
    let page: DocPage
    let docs: DocsModel

    var body: some View {
        let rows = DocRow.rows(page)
        ScrollViewReader { proxy in
            ScrollView {
                LazyVStack(alignment: .leading, spacing: 12) {
                    ForEach(rows) { row in
                        DocRowView(row: row)
                    }
                }
                .frame(maxWidth: Theme.readableWidth, alignment: .leading)
                .padding(.horizontal, 32)
                // The space above the page is the header row's own, so a scroll to the top
                // keeps it.
                .padding(.bottom, 40)
                .frame(maxWidth: .infinity)
            }
            .scrollEdgeEffectStyle(.hard, for: .top)
            .task(id: docs.scrollRequests) {
                // Rows far down the page are laid out with estimated heights until they come into
                // view, so the place moves as they take their real ones: the scroll is repeated
                // while they settle.
                for delay in [30, 60, 120, 240] {
                    try? await Task.sleep(for: .milliseconds(delay))
                    guard !Task.isCancelled else { return }
                    proxy.scrollTo(target, anchor: .top)
                }
            }
        }
        .textSelection(.enabled)
    }

    /// The row of the place on screen: a declaration, a section, or the top.
    private var target: String {
        guard let anchor = docs.location?.anchor, docs.location?.page == page.link else { return "#top" }
        switch page {
        case .package(let p): return p.declaration(for: anchor) ?? "#top"
        case .article(let a): return a.sections.contains { $0.anchor == anchor } ? anchor : "#top"
        case .acknowledgements: return "#top"
        }
    }
}

struct DocRowView: View {
    let row: DocRow
    @Environment(AppModel.self) private var app

    var body: some View {
        switch row {
        case .header(let eyebrow, let symbol, let title, let subtitle, let deprecated):
            VStack(alignment: .leading, spacing: 6) {
                Label(eyebrow, systemImage: symbol)
                    .symbolRenderingMode(.hierarchical)
                    .font(.system(size: 12, weight: .medium))
                    .foregroundStyle(Theme.violet)
                HStack(alignment: .firstTextBaseline, spacing: 10) {
                    Text(title)
                        .font(.system(size: Theme.headingSize(1), weight: .semibold))
                        .fixedSize(horizontal: false, vertical: true)
                        .accessibilityAddTraits(.isHeader)
                    if deprecated { DeprecatedTag() }
                }
                if let subtitle {
                    Text(subtitle)
                        .font(.system(size: 13, design: .monospaced))
                        .foregroundStyle(Theme.secondaryText)
                }
            }
            .padding(20)
            .frame(maxWidth: .infinity, alignment: .leading)
            .missionCard(tint: Theme.violet)
            .padding(.top, 26)
            .padding(.bottom, 6)
        case .block(_, let block):
            BlockView(block: block)
        case .heading(_, let title):
            Text(title)
                .font(.system(size: Theme.headingSize(2), weight: .semibold))
                .accessibilityAddTraits(.isHeader)
                .padding(.top, 22)
        case .index(let package, let entry):
            Button { app.openDocumentation(.symbol(package: package, name: entry.anchor)) } label: {
                HStack(spacing: 6) {
                    Text(entry.text)
                        .font(.system(size: 12.5, design: .monospaced))
                        .lineLimit(1)
                        .truncationMode(.tail)
                    if entry.deprecated { DeprecatedTag() }
                }
            }
            .buttonStyle(.link)
            .help(entry.text)
            .padding(.leading, entry.isMember ? 20 : 0)
            .padding(.vertical, -4)
        case .declaration(let declaration):
            DocDeclarationView(declaration: declaration)
        case .example(let example):
            DocExampleView(example: example)
        case .contents(let slug, let sections):
            ArticleContents(slug: slug, sections: sections)
        case .section(let section):
            VStack(alignment: .leading, spacing: 12) {
                Text(section.title)
                    .font(.system(size: Theme.headingSize(section.level), weight: .semibold))
                    .fixedSize(horizontal: false, vertical: true)
                    .accessibilityAddTraits(.isHeader)
                MarkdownBlocks(blocks: section.body)
            }
            .padding(.top, section.level <= 2 ? 22 : 10)
        case .license(let text):
            Text(text)
                .font(.system(size: 12, design: .monospaced))
                .fixedSize(horizontal: false, vertical: true)
                .padding(14)
                .frame(maxWidth: .infinity, alignment: .leading)
                .background(.quaternary.opacity(0.6), in: RoundedRectangle(cornerRadius: 10, style: .continuous))
        }
    }
}

/// A declaration: its heading, its Go source highlighted, its doc and its examples.
struct DocDeclarationView: View {
    let declaration: DocDeclaration

    var body: some View {
        VStack(alignment: .leading, spacing: 10) {
            if let title = declaration.title {
                HStack(alignment: .firstTextBaseline, spacing: 8) {
                    Text(title)
                        .font(.system(size: declaration.isMember ? 15 : 17, weight: .semibold, design: .monospaced))
                        .fixedSize(horizontal: false, vertical: true)
                        .accessibilityAddTraits(.isHeader)
                    if declaration.deprecated { DeprecatedTag() }
                }
            } else if declaration.deprecated {
                DeprecatedTag()
            }
            CodeBlockView(language: "go", code: declaration.source)
            if !declaration.doc.isEmpty {
                MarkdownBlocks(blocks: declaration.doc)
            }
            ForEach(declaration.examples) { example in
                DocExampleView(example: example)
            }
        }
        .padding(.top, declaration.title == nil ? 4 : 16)
    }

    /// A kind's symbol, as Xcode marks names: a letter in a square.
    static func symbol(_ kind: DocKind) -> String {
        switch kind {
        case .const: "c.square"
        case .var: "v.square"
        case .func: "f.square"
        case .type: "t.square"
        case .method: "m.square"
        }
    }
}

/// An example from the package's tests: what it shows, its code, and what it prints.
struct DocExampleView: View {
    let example: DocExample

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            Label(example.title, systemImage: "play.rectangle")
                .symbolRenderingMode(.hierarchical)
                .font(.system(size: 13, weight: .semibold))
                .foregroundStyle(Theme.secondaryText)
            if !example.doc.isEmpty {
                MarkdownBlocks(blocks: example.doc)
            }
            CodeBlockView(language: "go", code: example.code)
            if let output = example.output {
                Text(example.unordered ? "Output, in any order" : "Output")
                    .font(.system(size: 12, weight: .medium))
                    .foregroundStyle(Theme.secondaryText)
                if output.isEmpty {
                    Text("Prints nothing.")
                        .font(.system(size: 13))
                        .foregroundStyle(Theme.secondaryText)
                } else {
                    // Without the newline that ends the output, which would show as an empty line.
                    CodeBlockView(language: "", code: output.hasSuffix("\n") ? String(output.dropLast()) : output)
                }
            }
        }
        .padding(.top, 6)
    }
}

/// An article's table of contents: its sections two levels below the title, each a link.
struct ArticleContents: View {
    let slug: String
    let sections: [DocArticle.Section]
    @Environment(AppModel.self) private var app

    var body: some View {
        VStack(alignment: .leading, spacing: 5) {
            Text("Contents")
                .font(.system(size: 13, weight: .semibold))
                .foregroundStyle(Theme.secondaryText)
                .padding(.bottom, 2)
            ForEach(sections) { section in
                Button(section.title) { app.openDocumentation(.article(slug: slug, anchor: section.anchor)) }
                    .buttonStyle(.link)
                    .font(.system(size: section.level <= 2 ? 14 : 13))
                    .padding(.leading, CGFloat(max(0, section.level - 2)) * 18)
            }
        }
        .padding(14)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(.background.secondary, in: RoundedRectangle(cornerRadius: 12, style: .continuous))
    }
}

// MARK: - Quick Help

/// Quick Help in the editor: the name, its declaration, the first paragraph of its doc, and
/// the way to its page.
struct QuickHelpView: View {
    let summary: DocSummary
    let open: () -> Void
    /// Lines of a long declaration shown before it is cut; its page shows it all.
    static let lines = 14

    var body: some View {
        VStack(alignment: .leading, spacing: 10) {
            HStack(spacing: 8) {
                Text(summary.title)
                    .font(.system(size: 13, weight: .semibold, design: .monospaced))
                    .lineLimit(1)
                    .truncationMode(.middle)
                if summary.deprecated { DeprecatedTag() }
            }
            WrappingText(text: CodeText.attributed(declaration, language: "go", size: 12), accessibilityLabel: "Declaration")
                .padding(10)
                .frame(maxWidth: .infinity, alignment: .leading)
                .background(.quaternary.opacity(0.6), in: RoundedRectangle(cornerRadius: 8, style: .continuous))
            if case .paragraph(let inlines)? = summary.paragraph.first {
                Text(Inline.attributed(inlines, size: 13))
                    .lineSpacing(3)
                    .fixedSize(horizontal: false, vertical: true)
                    .textSelection(.enabled)
            }
            HStack {
                Spacer()
                Button("Open in Documentation", action: open)
                    .keyboardShortcut(.defaultAction)
                    .controlSize(.small)
            }
        }
        .padding(14)
        .frame(width: 440)
    }

    private var declaration: String {
        let lines = summary.declaration.split(separator: "\n", omittingEmptySubsequences: false)
        return lines.count > Self.lines ? lines.prefix(Self.lines).joined(separator: "\n") + "\n\t…" : summary.declaration
    }
}

// MARK: - Acknowledgements

/// Where the documentation comes from and Go's license, as the reference database says, for
/// Settings. The browser shows the same as a page.
struct AcknowledgementsSheet: View {
    let reference: ReferenceDatabase?
    @Environment(\.dismiss) private var dismiss
    @ViewState private var credits: DocCredits?
    @ViewState private var failure: String?

    var body: some View {
        VStack(spacing: 0) {
            ScrollView {
                VStack(alignment: .leading, spacing: 12) {
                    Text("Acknowledgements")
                        .font(.system(size: Theme.headingSize(2), weight: .semibold))
                    if let credits {
                        MarkdownBlocks(blocks: credits.credit)
                        Text("Go's License")
                            .font(.system(size: Theme.headingSize(3), weight: .semibold))
                            .padding(.top, 8)
                        DocRowView(row: .license(credits.license))
                    } else if let failure {
                        Text(failure).foregroundStyle(Theme.secondaryText)
                    } else {
                        ProgressView()
                    }
                }
                .padding(24)
                .frame(maxWidth: .infinity, alignment: .leading)
                .textSelection(.enabled)
            }
            Divider()
            HStack {
                Spacer()
                Button("Done") { dismiss() }
                    .keyboardShortcut(.defaultAction)
            }
            .padding(14)
        }
        .frame(width: 620, height: 560)
        .task {
            guard let reference else {
                failure = "The Go documentation isn't part of this copy of Gosmos."
                return
            }
            do { credits = try await reference.credits() } catch { failure = error.localizedDescription }
        }
    }
}
