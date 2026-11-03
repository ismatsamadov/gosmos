import GosmosKit
import SwiftUI

/// A lesson's hints: the ones you revealed, in order, and a button for the next. The last
/// hint can wait for a failed Submit; then the panel says why instead of offering it. With
/// no hints the panel is not there at all.
struct HintsPanel: View {
    let hints: [String]
    /// How many hints are revealed, counted from the first.
    let revealed: Int
    /// Why the next hint cannot show yet, or nil when it can.
    let nextLockedReason: String?
    let onReveal: () -> Void

    init(hints: [String], revealed: Int, nextLockedReason: String?, onReveal: @escaping () -> Void) {
        self.hints = hints
        self.revealed = min(max(revealed, 0), hints.count)
        self.nextLockedReason = nextLockedReason
        self.onReveal = onReveal
    }

    var body: some View {
        if !hints.isEmpty {
            VStack(alignment: .leading, spacing: 10) {
                HStack(alignment: .center, spacing: 8) {
                    Label("Hints", systemImage: "lightbulb")
                        .labelStyle(HintsTitleStyle())
                        .accessibilityAddTraits(.isHeader)
                    Spacer(minLength: 8)
                    if revealed < hints.count && nextLockedReason == nil {
                        Button {
                            onReveal()
                        } label: {
                            Text("Show Hint (\(revealed + 1) of \(hints.count))")
                        }
                        .secondaryGlassButton()
                        .controlSize(.small)
                        .onKeyPress(.return) {
                            onReveal()
                            return .handled
                        }
                        .accessibilityLabel("Show hint \(revealed + 1) of \(hints.count)")
                    }
                }
                ForEach(Array(hints.prefix(revealed).enumerated()), id: \.offset) { i, hint in
                    HStack(alignment: .firstTextBaseline, spacing: 10) {
                        Text("\(i + 1)")
                            .font(.system(size: 11, weight: .semibold, design: .rounded).monospacedDigit())
                            .foregroundStyle(TeachingStyle.digitOnFill)
                            .frame(width: 18, height: 18)
                            .background(.quaternary, in: Circle())
                            .accessibilityHidden(true)
                        Text(Inline.attributed(Readme.parseInline(hint), size: 14))
                            .fixedSize(horizontal: false, vertical: true)
                            .frame(maxWidth: .infinity, alignment: .leading)
                            .textSelection(.enabled)
                    }
                    .accessibilityElement(children: .combine)
                    .accessibilityLabel("Hint \(i + 1): \(Readme.plainText(Readme.parseInline(hint)))")
                }
                if revealed < hints.count, let reason = nextLockedReason {
                    // The lock where a hint's number is, and the reason where its text is.
                    HStack(alignment: .firstTextBaseline, spacing: 10) {
                        Image(systemName: "lock")
                            .frame(width: 18)
                        Text(reason)
                            .fixedSize(horizontal: false, vertical: true)
                    }
                    .font(.system(size: 13))
                    .foregroundStyle(Theme.secondaryText)
                    .accessibilityElement(children: .ignore)
                    .accessibilityLabel("Hint \(revealed + 1) of \(hints.count) is locked. \(reason)")
                }
            }
            .padding(.horizontal, 14)
            .padding(.vertical, 12)
            .frame(maxWidth: .infinity, alignment: .leading)
            .missionCard(tint: Theme.gold, radius: 14)
            .accessibilityElement(children: .contain)
            .accessibilityLabel("Hints, \(revealed) of \(hints.count) shown")
        }
    }
}

extension HintsPanel {
    /// The panel for a lesson's hint state: the hints revealed so far, in order, and the
    /// lock on the last one until a failed Submit.
    init(state: HintState, lockedReason: String = "The last hint opens after a Submit that fails.", onReveal: @escaping () -> Void) {
        let shown = state.shown
        let rest = state.hints.filter { hint in !shown.contains(hint) }
        self.init(hints: (shown + rest).map(\.text), revealed: shown.count,
                  nextLockedReason: state.nextIsLocked ? lockedReason : nil, onReveal: onReveal)
    }
}

/// The panel's title: a quiet symbol and a small semibold word.
private struct HintsTitleStyle: LabelStyle {
    func makeBody(configuration: Configuration) -> some View {
        HStack(spacing: 6) {
            configuration.icon
                .symbolRenderingMode(.hierarchical)
                .foregroundStyle(Theme.gold)
            configuration.title
                .font(.system(size: 13, weight: .semibold))
        }
    }
}
