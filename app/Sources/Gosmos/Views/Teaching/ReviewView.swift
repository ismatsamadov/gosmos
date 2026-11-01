import GosmosKit
import SwiftUI

/// The review: quiz questions due again, asked one at a time with lessons taking turns, each
/// explained once answered; then the lessons due to redo. Calm on purpose: no streaks, no
/// scores, nothing about what you missed. Keys 1 to 5 pick an answer, Return checks it, and
/// Return again moves on.
struct ReviewView: View {
    /// A quiz question due for review.
    struct Question: Identifiable, Hashable, Sendable {
        let ref: LessonRef
        /// The question's position in its lesson (quiz_question.position).
        let position: Int
        let lessonTitle: String
        /// The question, in Markdown.
        let text: String
        let answers: [String]
        /// The right answer's text.
        let rightAnswer: String
        /// Why each answer is right or wrong, in the answers' order; nil where the lesson gives none.
        let explanations: [String?]

        var id: String { "\(ref.library)/\(ref.lesson)/\(position)" }

        init(ref: LessonRef, position: Int, lessonTitle: String, text: String, answers: [String], rightAnswer: String,
             explanations: [String?]) {
            self.ref = ref
            self.position = position
            self.lessonTitle = lessonTitle
            self.text = text
            self.answers = answers
            self.rightAnswer = rightAnswer
            self.explanations = explanations
        }

        /// A lesson's quiz question.
        init(ref: LessonRef, lessonTitle: String, question: Quiz.Question) {
            self.init(ref: ref, position: question.position, lessonTitle: lessonTitle, text: question.question,
                      answers: question.answers, rightAnswer: question.answer, explanations: question.explanations)
        }

        func explanation(of answer: Int) -> String? {
            explanations.indices.contains(answer) ? explanations[answer] : nil
        }
    }

    /// A lesson due to redo, and why it came back.
    struct Lesson: Identifiable, Hashable, Sendable {
        let ref: LessonRef
        let title: String
        let type: LessonType?
        /// Why it is due, such as "You viewed the solution" or "3 failed Submits".
        let reason: String

        var id: LessonRef { ref }

        init(ref: LessonRef, title: String, type: LessonType? = nil, reason: String) {
            self.ref = ref
            self.title = title
            self.type = type
            self.reason = reason
        }

        /// Why a lesson came back: its solution was viewed before the pass, or the pass took
        /// three failed Submits or more.
        static func reason(viewedSolution: Bool, failedSubmits: Int) -> String {
            var parts: [String] = []
            if viewedSolution { parts.append("You viewed the solution") }
            if failedSubmits > 0 { parts.append(failedSubmits == 1 ? "1 failed Submit" : "\(failedSubmits) failed Submits") }
            return parts.isEmpty ? "Due for review" : parts.joined(separator: " · ")
        }
    }

    /// The due questions in the order they are asked: interleaved across lessons, as the
    /// learner database gives them (ReviewQueue).
    let questions: [Question]
    let lessons: [Lesson]
    /// A question was answered: the question, the answer picked, and whether it was right.
    let onAnswer: (Question, String, Bool) -> Void
    let onOpenLesson: (Lesson) -> Void

    /// The question picked or answered on screen; kept while its explanation shows, even if
    /// the caller drops it from the list.
    @ViewState private var current: Question?
    @ViewState private var picked: Int?
    @ViewState private var checked = false
    /// The questions answered here, so the rest can be counted and none comes twice.
    @ViewState private var answered: Set<String> = []
    @FocusState private var focused: Bool

    init(questions: [Question], lessons: [Lesson], onAnswer: @escaping (Question, String, Bool) -> Void,
         onOpenLesson: @escaping (Lesson) -> Void) {
        self.questions = questions
        self.lessons = lessons
        self.onAnswer = onAnswer
        self.onOpenLesson = onOpenLesson
    }

    /// Questions not answered yet, in the order they are asked.
    private var waiting: [Question] { questions.filter { !answered.contains($0.id) } }

    /// The questions left, the one on screen included until it is answered.
    private var left: Int { waiting.count }

    /// The question on screen: the one in hand, or else the next waiting.
    private var shown: Question? { current ?? waiting.first }

    var body: some View {
        Group {
            if shown == nil && lessons.isEmpty {
                ContentUnavailableView {
                    Label {
                        Text("Nothing Due").foregroundStyle(.primary)
                    } icon: {
                        Image(systemName: "checkmark.circle").foregroundStyle(.secondary)
                    }
                } description: {
                    Text("Come back tomorrow.")
                        .foregroundStyle(Theme.secondaryText)
                }
            } else {
                ScrollView {
                    VStack(alignment: .leading, spacing: 28) {
                        heading
                        if let shown {
                            card(shown)
                        } else if !questions.isEmpty {
                            Label("No more questions today.", systemImage: "checkmark.circle")
                                .font(.system(size: 14))
                                .foregroundStyle(Theme.secondaryText)
                        }
                        if !lessons.isEmpty {
                            lessonList
                        }
                    }
                    .frame(maxWidth: Theme.readableWidth, alignment: .leading)
                    .padding(.horizontal, 32)
                    .padding(.top, 26)
                    .padding(.bottom, 40)
                    .frame(maxWidth: .infinity)
                }
                // As the readme: text scrolled under the toolbar stops at a hard edge.
                .scrollEdgeEffectStyle(.hard, for: .top)
            }
        }
    }

    // MARK: Heading

    private var heading: some View {
        VStack(alignment: .leading, spacing: 6) {
            Text("Review")
                .font(.system(size: Theme.headingSize(1), weight: .semibold))
                .accessibilityAddTraits(.isHeader)
            Text(summary)
                .font(.system(size: 14))
                .foregroundStyle(Theme.secondaryText)
        }
    }

    private var summary: String {
        var parts: [String] = []
        if left > 0 { parts.append(left == 1 ? "1 question" : "\(left) questions") }
        if !lessons.isEmpty { parts.append(lessons.count == 1 ? "1 lesson to redo" : "\(lessons.count) lessons to redo") }
        return parts.isEmpty ? "All done for today." : "Due now: " + parts.joined(separator: " and ") + "."
    }

    // MARK: The question

    private func card(_ question: Question) -> some View {
        VStack(alignment: .leading, spacing: 14) {
            HStack(alignment: .firstTextBaseline) {
                Label(question.lessonTitle, systemImage: "questionmark.bubble")
                    .font(.system(size: 12, weight: .medium))
                    .foregroundStyle(Theme.secondaryText)
                    .lineLimit(1)
                Spacer(minLength: 12)
                Text("\(left) left")
                    .font(.system(size: 12, weight: .medium).monospacedDigit())
                    .foregroundStyle(Theme.secondaryText)
                    .accessibilityLabel(left == 1 ? "1 question left" : "\(left) questions left")
            }
            MarkdownBlocks(blocks: Readme.parse(question.text))
            GlassEffectContainer(spacing: 8) {
                VStack(spacing: 8) {
                    ForEach(Array(question.answers.prefix(9).enumerated()), id: \.offset) { a, answer in
                        QuizAnswerRow(number: a + 1, text: answer, isPicked: picked == a, isRight: answer == question.rightAnswer,
                                      explanation: question.explanation(of: a), revealed: checked) {
                            pick(a)
                            focused = true
                        }
                    }
                }
            }
            footer(question)
        }
        .padding(18)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(.background.secondary, in: RoundedRectangle(cornerRadius: 16, style: .continuous))
        .overlay(RoundedRectangle(cornerRadius: 16, style: .continuous)
            .strokeBorder(focused ? AnyShapeStyle(.tint.opacity(0.6)) : AnyShapeStyle(.separator.opacity(0.5)), lineWidth: focused ? 1.5 : 0.5))
        .focusable()
        .focusEffectDisabled()
        .focused($focused)
        .onKeyPress(characters: .decimalDigits) { press in
            guard !checked, let n = Int(press.characters), n >= 1, n <= min(5, question.answers.count) else { return .ignored }
            pick(n - 1)
            return .handled
        }
        .onKeyPress(.return) {
            if checked { next() } else if picked != nil { check() } else { return .ignored }
            return .handled
        }
        .onAppear { focused = true }
        .accessibilityElement(children: .contain)
        .accessibilityLabel("Question from \(question.lessonTitle)")
    }

    private func footer(_ question: Question) -> some View {
        HStack(spacing: 12) {
            if checked {
                let right = picked.map { question.answers[$0] == question.rightAnswer } ?? false
                Label(right ? "Right." : "Not quite. The right answer is marked.",
                      systemImage: right ? "checkmark.circle" : "info.circle")
                    .font(.system(size: 13, weight: right ? .semibold : .regular))
                    .foregroundStyle(right ? AnyShapeStyle(Theme.success) : AnyShapeStyle(Theme.secondaryText))
                Spacer()
                Button(waiting.isEmpty ? "Done" : "Next Question", action: next)
                    .prominentGlassButton()
                    .help("Return")
            } else {
                Text(picked == nil ? "Pick an answer (keys 1 to \(min(5, question.answers.count))), then press Return." : "Press Return to check your answer.")
                    .font(.system(size: 13))
                    .foregroundStyle(Theme.secondaryText)
                Spacer()
                Button("Check Answer", action: check)
                    .prominentGlassButton()
                    .disabled(picked == nil)
                    .help("Return")
            }
        }
    }

    // MARK: Lessons to redo

    private var lessonList: some View {
        VStack(alignment: .leading, spacing: 10) {
            Text("Lessons to Redo")
                .font(.system(size: Theme.headingSize(3), weight: .semibold))
                .accessibilityAddTraits(.isHeader)
            Text("A redo starts from the starter code in a folder of its own and isn't kept; your own code stays as it is.")
                .font(.system(size: 13))
                .foregroundStyle(Theme.secondaryText)
                .fixedSize(horizontal: false, vertical: true)
            VStack(spacing: 8) {
                ForEach(lessons) { lesson in
                    HStack(alignment: .center, spacing: 12) {
                        Image(systemName: lesson.type?.symbol ?? "arrow.counterclockwise")
                            .symbolRenderingMode(.hierarchical)
                            .foregroundStyle(.secondary)
                            .frame(width: 20)
                            .accessibilityHidden(true)
                        VStack(alignment: .leading, spacing: 2) {
                            Text(lesson.title)
                                .font(.system(size: 14, weight: .medium))
                                .lineLimit(1)
                            Text(lesson.reason)
                                .font(.system(size: 12))
                                .foregroundStyle(Theme.secondaryText)
                                .lineLimit(1)
                        }
                        .accessibilityElement(children: .combine)
                        Spacer(minLength: 12)
                        Button("Redo") { onOpenLesson(lesson) }
                            .secondaryGlassButton()
                            .accessibilityLabel("Redo \(lesson.title)")
                            .help("Open \(lesson.title) to redo it")
                    }
                    .padding(.horizontal, 14)
                    .padding(.vertical, 10)
                    .frame(minHeight: 44)
                    .background(.background.secondary, in: RoundedRectangle(cornerRadius: 12, style: .continuous))
                }
            }
        }
    }

    // MARK: Steps

    private func pick(_ answer: Int) {
        guard !checked, let question = shown, question.answers.indices.contains(answer) else { return }
        current = question
        picked = answer
    }

    private func check() {
        guard !checked, let question = shown, let picked, question.answers.indices.contains(picked) else { return }
        current = question
        checked = true
        answered.insert(question.id)
        let answer = question.answers[picked]
        onAnswer(question, answer, answer == question.rightAnswer)
    }

    private func next() {
        guard checked else { return }
        current = nil
        picked = nil
        checked = false
        if shown != nil { focused = true }
    }
}

/// The number of things due for review, for the sidebar's Review row. Nothing shows at zero.
struct ReviewBadge: View {
    let count: Int

    init(count: Int) {
        self.count = count
    }

    var body: some View {
        if count > 0 {
            Text(count > 99 ? "99+" : "\(count)")
                .font(.system(size: 11, weight: .semibold).monospacedDigit())
                .foregroundStyle(TeachingStyle.digitOnFill)
                .padding(.horizontal, 7)
                .padding(.vertical, 1.5)
                .background(.quaternary, in: Capsule())
                .accessibilityLabel(count == 1 ? "1 due for review" : "\(count) due for review")
        }
    }
}
