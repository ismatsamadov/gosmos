import GosmosKit
import SwiftUI

/// One answer of a quiz question, as QuizView draws it: the key that picks it, the answer,
/// and its state. Once the question is graded (`revealed`), the right answer turns green and
/// a wrong pick red, and every answer shows its explanation under it.
struct QuizAnswerRow: View {
    let number: Int
    let text: String
    let isPicked: Bool
    /// This is the question's right answer.
    let isRight: Bool
    let explanation: String?
    /// The question is graded: verdicts and explanations show.
    let revealed: Bool
    let action: () -> Void

    init(number: Int, text: String, isPicked: Bool, isRight: Bool, explanation: String?, revealed: Bool,
         action: @escaping () -> Void) {
        self.number = number
        self.text = text
        self.isPicked = isPicked
        self.isRight = isRight
        self.explanation = explanation
        self.revealed = revealed
        self.action = action
    }

    private enum Mark { case idle, picked, right, wrong }


    private var mark: Mark {
        guard revealed else { return isPicked ? .picked : .idle }
        if isRight { return .right }
        return isPicked ? .wrong : .idle
    }

    /// The explanation, when the question is graded and the lesson gives one.
    private var shownExplanation: String? {
        guard revealed, let explanation = explanation?.trimmingCharacters(in: .whitespacesAndNewlines), !explanation.isEmpty else { return nil }
        return explanation
    }

    var body: some View {
        Button(action: action) {
            HStack(alignment: .firstTextBaseline, spacing: 12) {
                Text("\(number)")
                    .font(.system(size: 11, weight: .medium, design: .monospaced))
                    .foregroundStyle(mark == .idle ? TeachingStyle.digitOnFill : tint)
                    .frame(width: 26, height: 26)
                    .background(mark == .idle ? Theme.secondaryText.opacity(0.08) : tint.opacity(0.12), in: RoundedRectangle(cornerRadius: 7))
                VStack(alignment: .leading, spacing: 7) {
                    Text(Inline.attributed(Readme.parseInline(text), size: 14))
                        .multilineTextAlignment(.leading)
                        .fixedSize(horizontal: false, vertical: true)
                    if let shownExplanation {
                        ExplanationLine(text: shownExplanation, isRight: isRight, isPicked: isPicked)
                            .transition(.opacity)
                    }
                }
                .frame(maxWidth: .infinity, alignment: .leading)
                Image(systemName: symbol)
                    .symbolRenderingMode(.hierarchical)
                    .foregroundStyle(mark == .idle ? Theme.secondaryText.opacity(0.55) : tint)
                    .contentTransition(.symbolEffect(.replace))
            }
            .padding(.horizontal, 14)
            .padding(.vertical, 11)
            .frame(maxWidth: .infinity, minHeight: 44, alignment: .leading)
            .contentShape(RoundedRectangle(cornerRadius: 12, style: .continuous))
        }
        .buttonStyle(MissionRowButtonStyle())
        .background(mark == .idle ? Theme.inset.opacity(0.65) : tint.opacity(0.07),
                    in: RoundedRectangle(cornerRadius: 10, style: .continuous))
        .overlay {
            RoundedRectangle(cornerRadius: 10).strokeBorder(mark == .idle ? Theme.border : tint.opacity(0.6), lineWidth: 1)
                .allowsHitTesting(false)
        }
        .accessibilityLabel("Answer \(number): \(Readme.plainText(Readme.parseInline(text)))")
        .accessibilityValue(accessibilityState)
        .accessibilityHint("Press \(number) to pick it")
    }

    private var symbol: String {
        switch mark {
        case .idle: "circle"
        case .picked: "circle.inset.filled"
        case .right: "checkmark.circle.fill"
        case .wrong: "xmark.circle.fill"
        }
    }

    private var tint: Color {
        switch mark {
        case .idle, .picked: Theme.accent
        case .right: Theme.success
        case .wrong: Theme.failure
        }
    }

    private var accessibilityState: String {
        let state = switch mark {
        case .idle: revealed ? "Not the right answer" : ""
        case .picked: "Picked"
        case .right: isPicked ? "Right answer, your pick" : "Right answer"
        case .wrong: "Wrong answer, your pick"
        }
        guard let shownExplanation else { return state }
        let why = Readme.plainText(Readme.parseInline(shownExplanation))
        return state.isEmpty ? why : "\(state). \(why)"
    }
}

/// Why an answer is right or wrong, under the answer: the reason in the body's color, with a
/// small mark in green for the right answer, red for your wrong pick, gray for the others.
private struct ExplanationLine: View {
    let text: String
    let isRight: Bool
    let isPicked: Bool

    var body: some View {
        HStack(alignment: .firstTextBaseline, spacing: 6) {
            Image(systemName: isRight ? "checkmark.circle" : "xmark.circle")
                .font(.system(size: 12, weight: .semibold))
                .foregroundStyle(isRight ? Theme.success : isPicked ? Theme.failure : Theme.secondaryText)
            Text(Inline.attributed(Readme.parseInline(text), size: 13))
                .foregroundStyle(.primary)
                .multilineTextAlignment(.leading)
                .fixedSize(horizontal: false, vertical: true)
        }
        .accessibilityHidden(true)
    }
}
