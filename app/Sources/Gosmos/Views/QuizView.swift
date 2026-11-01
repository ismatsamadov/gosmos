import GosmosKit
import SwiftUI

struct QuizLessonView: View {
    let lesson: LessonModel
    let session: QuizSession

    var body: some View {
        ReadmeColumn(lesson: lesson) {
            VStack(alignment: .leading, spacing: 18) {
                ForEach(Array(session.quiz.questions.enumerated()), id: \.offset) { i, question in
                    QuestionCard(session: session, index: i, question: question,
                                 blocks: lesson.parsed.questions.indices.contains(i) ? lesson.parsed.questions[i] : [])
                }
                QuizFooter(session: session, lesson: lesson)
            }
        }
    }
}

struct QuestionCard: View {
    let session: QuizSession
    let index: Int
    let question: Quiz.Question
    /// The question's Markdown, parsed when the lesson opened.
    let blocks: [MDBlock]
    @FocusState private var focused: Bool

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            if session.quiz.questions.count > 1 {
                Text("Question \(index + 1) of \(session.quiz.questions.count)")
                    .font(.system(size: 12, weight: .medium))
                    .foregroundStyle(Theme.secondaryText)
            }
            MarkdownBlocks(blocks: blocks)
            // Empty until the question is graded: the reasons would give the answer away.
            let explanations = session.explanations(for: index)
            GlassEffectContainer(spacing: 8) {
                VStack(spacing: 8) {
                    ForEach(Array(question.answers.enumerated()), id: \.offset) { a, answer in
                        QuizAnswerRow(number: a + 1, text: answer, isPicked: session.picks[index] == a, isRight: question.isCorrect(answer),
                                      explanation: explanations.indices.contains(a) ? explanations[a] : nil,
                                      revealed: session.verdicts[index] != nil) {
                            session.pick(a, for: index)
                            focused = true
                        }
                    }
                }
            }
            if session.verdicts[index] == false {
                Label("Not quite. The right answer is marked.", systemImage: "info.circle")
                    .font(.system(size: 13))
                    .foregroundStyle(Theme.secondaryText)
            }
        }
        .padding(18)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(.background.secondary, in: RoundedRectangle(cornerRadius: 16, style: .continuous))
        .overlay(RoundedRectangle(cornerRadius: 16, style: .continuous).strokeBorder(focused ? AnyShapeStyle(.tint.opacity(0.6)) : AnyShapeStyle(.separator.opacity(0.5)), lineWidth: focused ? 1.5 : 0.5))
        .focusable()
        .focusEffectDisabled()
        .focused($focused)
        .onChange(of: focused) { _, isFocused in if isFocused { session.focusedQuestion = index } }
        .onKeyPress(characters: .decimalDigits) { press in
            guard let n = Int(press.characters), (1...5).contains(n) else { return .ignored }
            session.pick(n - 1, for: index)
            return .handled
        }
        .onAppear { if index == 0 { focused = true } }
        .accessibilityElement(children: .contain)
        .accessibilityLabel("Question \(index + 1)")
    }
}

struct QuizFooter: View {
    let session: QuizSession
    let lesson: LessonModel
    @Environment(AppModel.self) private var app

    var body: some View {
        let correct = session.allCorrect
        let wrong = session.verdicts.values.contains(false)
        HStack(spacing: 12) {
            // One icon that stays and changes, so its bounce plays when the quiz is passed.
            Image(systemName: correct ? "checkmark.seal.fill" : wrong ? "arrow.uturn.backward.circle" : "keyboard")
                .contentTransition(.symbolEffect(.replace))
                .symbolEffect(.bounce, value: session.successCount)
                .foregroundStyle(correct ? AnyShapeStyle(Theme.success) : AnyShapeStyle(.secondary))
                .font(.system(size: correct ? 14 : 13, weight: correct ? .semibold : .regular))
                .accessibilityHidden(true)
            if correct {
                Text("All correct")
                    .foregroundStyle(Theme.success)
                    .font(.system(size: 14, weight: .semibold))
                Spacer()
                if let next = app.nextLesson {
                    Button("Next: \(next.title)", systemImage: "arrow.right") { app.openNext() }
                        .prominentGlassButton()
                        .lineLimit(1)
                }
            } else {
                Text(wrong ? "Pick another answer where one is marked wrong, then submit again (⌘↩)."
                           : session.canSubmit ? "Ready: submit your answers (⌘↩)." : "Pick an answer for each question (keys 1 to 5), then submit.")
                    .font(.system(size: 13))
                    .foregroundStyle(Theme.secondaryText)
                Spacer()
            }
        }
    }
}
