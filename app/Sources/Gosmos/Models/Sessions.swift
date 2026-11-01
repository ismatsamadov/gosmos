import AppKit
import Foundation
import GosmosKit
import Observation

/// A quiz: pick an answer per question (click, or keys 1 to 5), then Submit checks them all.
@Observable @MainActor
final class QuizSession {
    let lesson: LessonDetail
    /// The picks and verdicts (GosmosKit's QuizState, which the tests cover).
    private(set) var state: QuizState
    var focusedQuestion = 0
    private(set) var successCount = 0
    @ObservationIgnored private weak var app: AppModel?
    @ObservationIgnored private let learner: LearnerDatabase

    init(lesson: LessonDetail, quiz: Quiz, learner: LearnerDatabase, app: AppModel) {
        self.lesson = lesson
        self.state = QuizState(quiz: quiz)
        self.learner = learner
        self.app = app
    }

    var quiz: Quiz { state.quiz }
    var picks: [Int: Int] { state.picks }
    var verdicts: [Int: Bool] { state.verdicts }
    var allCorrect: Bool { state.allCorrect }
    var canSubmit: Bool { state.canSubmit }

    func pick(_ answer: Int, for question: Int) {
        if state.pick(answer, for: question) { focusedQuestion = question }
    }

    /// Key 1 to 5 picks an answer of the focused question.
    func pick(number: Int) {
        pick(number - 1, for: focusedQuestion)
    }

    /// Why each answer of a graded question is right or wrong; nothing before it is graded.
    func explanations(for question: Int) -> [String?] { state.explanations(for: question) }

    /// Grades the picks, records them (a question right for the first time comes back for
    /// review), and finishes the lesson when every answer is right.
    func submit() {
        guard canSubmit, let app else { return }
        let answers = state.submit()
        let learner = learner, ref = lesson.ref
        learner.pendingWrites.start {
            do {
                try await learner.recordAnswers(ref, answers)
                await app.refreshDueCount()
            } catch {
                app.report(error)
            }
        }
        let owner = app.lesson.flatMap { $0.quiz === self ? $0 : nil }
        if answers.contains(where: { !$0.right }) { owner?.noteFailedSubmit() }
        if allCorrect {
            successCount += 1
            NSHapticFeedbackManager.defaultPerformer.perform(.generic, performanceTime: .now)
            owner?.markPassed(review: false)
        }
    }
}
