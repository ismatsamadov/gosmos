import AppKit
import CodeEditor
import GosmosKit
import SwiftUI
import SyntaxCore

/// Your code beside the lesson's reference solution, file by file: the solution's note on
/// top, then the two side by side, Go highlighted, with the lines that differ marked. It only
/// shows; nothing in it goes into the editor. For a sheet or a split pane.
struct SolutionCompareView: View {
    /// A file you edit: your version and the solution's.
    struct File: Identifiable, Hashable, Sendable {
        let name: String
        let yours: String
        let solution: String
        var id: String { name }

        init(name: String, yours: String, solution: String) {
            self.name = name
            self.yours = yours
            self.solution = solution
        }

        /// The reference solution's file, with your version of it.
        init(_ solution: SolutionFile, yours: String) {
            self.init(name: solution.name, yours: yours, solution: solution.content)
        }
    }

    let files: [File]
    let note: String?
    let onClose: () -> Void
    @ViewState private var selected = 0

    init(files: [File], note: String?, onClose: @escaping () -> Void) {
        self.files = files
        self.note = note
        self.onClose = onClose
    }

    var body: some View {
        let index = files.indices.contains(selected) ? selected : 0
        VStack(alignment: .leading, spacing: 0) {
            header
            if let note = note?.trimmingCharacters(in: .whitespacesAndNewlines), !note.isEmpty {
                SolutionNote(text: note)
                    .padding(.horizontal, 20)
                    .padding(.bottom, 14)
            }
            if files.indices.contains(index) {
                ComparePane(file: files[index])
                    .id(files[index].id)
            } else {
                ContentUnavailableView("No Solution Files", systemImage: "doc.text.magnifyingglass",
                                       description: Text("This lesson has no reference solution to compare with.").foregroundStyle(Theme.secondaryText))
            }
        }
        .frame(minWidth: 760, idealWidth: 1040, maxWidth: .infinity, minHeight: 480, idealHeight: 700, maxHeight: .infinity)
        .background(Theme.canvas)
    }

    private var header: some View {
        HStack(alignment: .center, spacing: 12) {
            MissionIcon(symbol: "doc.text.magnifyingglass", color: Theme.violet)
            VStack(alignment: .leading, spacing: 2) {
                Text("A different way to get there.")
                    .font(.system(size: 19, weight: .bold, design: .rounded))
                    .accessibilityAddTraits(.isHeader)
                Text("Compare it with your code. Nothing here goes into the editor.")
                    .font(.system(size: 12))
                    .foregroundStyle(Theme.secondaryText)
            }
            Spacer(minLength: 12)
            if files.count > 1 {
                Picker("File", selection: $selected) {
                    ForEach(Array(files.enumerated()), id: \.offset) { i, file in
                        Text(file.name).tag(i)
                    }
                }
                .pickerStyle(.segmented)
                .labelsHidden()
                .fixedSize()
                // The selected file's white label needs the deeper fill in dark mode.
                .tint(Theme.prominent)
                .accessibilityLabel("File")
            }
            Button("Done", action: onClose)
                .secondaryGlassButton()
                .keyboardShortcut(.cancelAction)
                .help("Close the solution (Esc)")
        }
        .padding(.horizontal, 20)
        .padding(.top, 16)
        .padding(.bottom, 12)
    }
}

/// Why the solution works, in a quiet panel.
private struct SolutionNote: View {
    let text: String

    var body: some View {
        HStack(alignment: .firstTextBaseline, spacing: 8) {
            Image(systemName: "text.bubble")
                .symbolRenderingMode(.hierarchical)
                .foregroundStyle(.secondary)
                .accessibilityHidden(true)
            Text(Inline.attributed(Readme.parseInline(text), size: 14))
                .fixedSize(horizontal: false, vertical: true)
                .frame(maxWidth: .infinity, alignment: .leading)
                .textSelection(.enabled)
        }
        .padding(.horizontal, 14)
        .padding(.vertical, 10)
        .missionCard(tint: Theme.violet, radius: 12)
        .accessibilityElement(children: .combine)
        .accessibilityLabel("Why the solution works: \(Readme.plainText(Readme.parseInline(text)))")
    }
}

/// The two versions of one file, row by row. The rows are worked out off the main actor, as a
/// long file takes a while, with a small progress indicator in their place meanwhile; when
/// the sheet closes first, the work stops.
private struct ComparePane: View {
    let file: SolutionCompareView.File
    /// Nil until worked out for this file.
    @ViewState private var rows: [SolutionCompareView.Row]?

    var body: some View {
        VStack(spacing: 0) {
            HStack(spacing: 0) {
                columnTitle("Your code")
                Divider()
                columnTitle("Solution")
            }
            .frame(height: 30)
            .background(.background.secondary)
            Divider()
            Group {
                if let rows {
                    ScrollView(.vertical) {
                        LazyVStack(alignment: .leading, spacing: 0) {
                            ForEach(rows) { row in
                                CompareRowView(row: row)
                            }
                        }
                        .padding(.vertical, 6)
                    }
                } else {
                    ComparingIndicator()
                }
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .background(Color(nsColor: SyntaxTheme.xcode.background))
            Divider()
            footer
                .font(.system(size: 12))
                .foregroundStyle(Theme.secondaryText)
                .padding(.horizontal, 20)
                .padding(.vertical, 8)
                .accessibilityElement(children: .combine)
        }
        // Again only when the file changes. SwiftUI cancels the task when the pane goes, as
        // when the sheet closes, and the comparison gives up then.
        .task(id: file) {
            rows = nil
            let compared = try? await SolutionCompareView.Row.rows(yours: file.yours, solution: file.solution, name: file.name)
            if !Task.isCancelled { rows = compared }
        }
    }

    private var footer: some View {
        HStack(spacing: 8) {
            if let rows {
                let differences = rows.filter(\.differs).count
                if differences == 0 {
                    Label("Your code matches the solution.", systemImage: "equal.circle")
                } else {
                    RoundedRectangle(cornerRadius: 2, style: .continuous)
                        .fill(CompareStyle.marker)
                        .frame(width: 3, height: 13)
                        .accessibilityHidden(true)
                    Text(differences == 1 ? "1 line differs" : "\(differences) lines differ")
                }
            } else {
                // Keeps the footer's height until the rows are ready.
                Text("Comparing…").hidden()
            }
            Spacer()
        }
    }

    private func columnTitle(_ title: String) -> some View {
        Text(title)
            .font(.system(size: 12, weight: .semibold))
            .foregroundStyle(Theme.secondaryText)
            .padding(.leading, 52)
            .frame(maxWidth: .infinity, alignment: .leading)
            .accessibilityAddTraits(.isHeader)
    }
}

/// In the rows' place while they are worked out: nothing for a moment, so a quick comparison
/// never flashes it, then a small spinner.
private struct ComparingIndicator: View {
    @ViewState private var shown = false

    var body: some View {
        HStack(spacing: 8) {
            ProgressView()
                .controlSize(.small)
            Text("Comparing…")
                .font(.system(size: 12))
                .foregroundStyle(Theme.secondaryText)
        }
        .opacity(shown ? 1 : 0)
        .accessibilityElement(children: .ignore)
        .accessibilityLabel("Comparing your code with the solution")
        .accessibilityHidden(!shown)
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .task {
            guard (try? await Task.sleep(for: .milliseconds(250))) != nil else { return }
            shown = true
        }
    }
}

/// How differing lines are marked (a faint wash of the accent and a thin bar at the edge), and
/// the two colors that differ from the editor's: line numbers and dark-mode comments, which in
/// Xcode's palette read under 4.5:1 on the code background.
private enum CompareStyle {
    static let wash = Theme.accent.opacity(0.09)
    static let marker = Theme.accent.opacity(0.75)
    static let missing = Color.secondary.opacity(0.06)
    static let codeSize: CGFloat = 12
    static let lineNumber = color(light: 0x686868, dark: 0x9A9AA0, lightContrast: 0x4D4D4D, darkContrast: 0xB4B4BA)
    /// Xcode's comment gray in light mode; in dark mode, its Increase Contrast gray.
    static let comment = color(light: 0x5D6C79, dark: 0x95A2AE, lightContrast: 0x45515C, darkContrast: 0x95A2AE)

    private static func color(light: UInt32, dark: UInt32, lightContrast: UInt32, darkContrast: UInt32) -> Color {
        func rgb(_ v: UInt32) -> NSColor {
            NSColor(srgbRed: CGFloat(v >> 16 & 0xFF) / 255, green: CGFloat(v >> 8 & 0xFF) / 255, blue: CGFloat(v & 0xFF) / 255, alpha: 1)
        }
        let colors = [rgb(light), rgb(dark), rgb(lightContrast), rgb(darkContrast)]
        return Color(nsColor: NSColor(name: nil) { appearance in
            switch appearance.bestMatch(from: [.aqua, .darkAqua, .accessibilityHighContrastAqua, .accessibilityHighContrastDarkAqua]) {
            case .darkAqua?: colors[1]
            case .accessibilityHighContrastAqua?: colors[2]
            case .accessibilityHighContrastDarkAqua?: colors[3]
            default: colors[0]
            }
        })
    }
}

private struct CompareRowView: View {
    let row: SolutionCompareView.Row

    var body: some View {
        HStack(alignment: .top, spacing: 0) {
            cell(row.left)
            Divider()
            cell(row.right)
        }
        // Both halves take the taller one's height, so a wrapped line keeps its row even.
        .fixedSize(horizontal: false, vertical: true)
        .accessibilityElement(children: .ignore)
        .accessibilityLabel(accessibilityText)
    }

    private func cell(_ line: SolutionCompareView.Row.Line?) -> some View {
        HStack(alignment: .firstTextBaseline, spacing: 10) {
            Text(line.map { "\($0.number)" } ?? "")
                .font(.system(size: 11, design: .monospaced).monospacedDigit())
                .foregroundStyle(CompareStyle.lineNumber)
                .frame(width: 30, alignment: .trailing)
            Text(line.map { $0.text.characters.isEmpty ? AttributedString(" ") : $0.text } ?? AttributedString(" "))
                .font(.system(size: CompareStyle.codeSize, design: .monospaced))
                .fixedSize(horizontal: false, vertical: true)
                .frame(maxWidth: .infinity, alignment: .leading)
                .textSelection(.enabled)
        }
        .padding(.leading, 8)
        .padding(.trailing, 12)
        .padding(.vertical, 1.5)
        .frame(maxWidth: .infinity, maxHeight: .infinity, alignment: .topLeading)
        .background {
            if line == nil {
                CompareStyle.missing
            } else if row.differs {
                CompareStyle.wash
            }
        }
        .overlay(alignment: .leading) {
            if row.differs && line != nil {
                Rectangle().fill(CompareStyle.marker).frame(width: 2)
            }
        }
    }

    private var accessibilityText: String {
        func side(_ name: String, _ line: SolutionCompareView.Row.Line?) -> String {
            guard let line else { return "\(name): no line" }
            let text = String(line.text.characters).trimmingCharacters(in: .whitespaces)
            return "\(name), line \(line.number): \(text.isEmpty ? "blank" : text)"
        }
        let both = "\(side("Yours", row.left)). \(side("Solution", row.right))"
        return row.differs ? "Differs. \(both)" : both
    }
}

extension SolutionCompareView {
    /// One row of the side-by-side view: a line of yours, a line of the solution, or both.
    struct Row: Identifiable, Hashable {
        struct Line: Hashable {
            let number: Int
            let text: AttributedString
        }

        let id: Int
        let left: Line?
        let right: Line?
        let differs: Bool

        /// The rows for two versions of a file, from LineDiff: lines that match side by side,
        /// the lines that differ paired between them, each in the editor's colors. Worked out
        /// off the main actor, since a long file takes a while; it gives up once its task is
        /// cancelled.
        @concurrent static func rows(yours: String, solution: String, name: String) async throws(CancellationError) -> [Row] {
            let highlighter: any SyntaxHighlighting = name.lowercased().hasSuffix(".go") ? GoLexer() : PlainTextHighlighter()
            let left = try highlightedLines(yours, highlighter: highlighter)
            let right = try highlightedLines(solution, highlighter: highlighter)
            return try LineDiff.rows(left.map { String($0.characters) }, right.map { String($0.characters) }).enumerated().map { i, row in
                Row(id: i, left: row.left.map { Line(number: $0 + 1, text: left[$0]) },
                    right: row.right.map { Line(number: $0 + 1, text: right[$0]) }, differs: row.differs)
            }
        }

        /// The code's lines, highlighted as the editor highlights them. Tabs become four spaces,
        /// and a final newline adds no empty line. Cancelled, it stops within a line.
        static func highlightedLines(_ code: String, highlighter: any SyntaxHighlighting) throws(CancellationError) -> [AttributedString] {
            var source = code.replacingOccurrences(of: "\r\n", with: "\n").replacingOccurrences(of: "\t", with: "    ")
            if source.hasSuffix("\n") { source.removeLast() }
            let full = highlighted(source, highlighter: highlighter)
            var lines: [AttributedString] = []
            var start = full.startIndex
            var i = full.startIndex
            while i < full.endIndex {
                let next = full.characters.index(after: i)
                if full.characters[i] == "\n" {
                    if Task.isCancelled { throw CancellationError() }
                    lines.append(AttributedString(full[start..<i]))
                    start = next
                }
                i = next
            }
            lines.append(AttributedString(full[start..<full.endIndex]))
            return lines
        }

        /// The code in the editor's colors and weights (SyntaxTheme), comments excepted.
        static func highlighted(_ source: String, highlighter: any SyntaxHighlighting) -> AttributedString {
            let theme = SyntaxTheme.xcode
            var result = AttributedString(source)
            result.font = .system(size: CompareStyle.codeSize, design: .monospaced)
            result.foregroundColor = Color(nsColor: theme.color(for: .plain))
            for token in highlighter.tokens(in: source) where token.kind != .plain {
                guard let range = Range(token.range, in: result) else { continue }
                result[range].font = .system(size: CompareStyle.codeSize, weight: theme.isBold(token.kind) ? .bold : .regular, design: .monospaced)
                result[range].foregroundColor = token.kind == .comment ? CompareStyle.comment : Color(nsColor: theme.color(for: token.kind))
            }
            return result
        }
    }
}
