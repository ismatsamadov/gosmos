import Foundation
import GosmosKit
import Observation

/// The review on screen: the quiz questions due now, interleaved across lessons, and the
/// lessons due to redo. The review view asks the questions; the session records each answer
/// and opens a redo. Calm on purpose: no streaks, no counts of what you missed, no reminders.
@Observable @MainActor
final class ReviewSession {
    /// A due question with what it needs on screen.
    struct Question {
        let ref: LessonRef
        let lessonTitle: String
        let question: Quiz.Question
    }

    /// A lesson due to redo, and what brought it back.
    struct Lesson {
        let ref: LessonRef
        let title: String
        let type: LessonType
        let viewedSolution: Bool
        let failedSubmits: Int
    }

    private(set) var questions: [Question] = []
    private(set) var lessons: [Lesson] = []
    private(set) var isLoaded = false
    @ObservationIgnored private let learner: LearnerDatabase
    @ObservationIgnored private weak var app: AppModel?

    init(learner: LearnerDatabase, app: AppModel) {
        self.learner = learner
        self.app = app
    }

    /// Reads what is due from the learner database, and the questions and titles from the
    /// libraries in the sidebar. Something whose library or lesson is gone is left out.
    func load() async {
        guard let app else { return }
        do {
            let due = try await ReviewSession.due(learner, in: app)
            var redo: [Lesson] = []
            for (ref, summary) in due.lessons {
                redo.append(Lesson(ref: ref, title: summary.title, type: summary.type,
                                   viewedSolution: try await learner.viewedSolution(ref),
                                   failedSubmits: try await learner.failedAttempts(ref, kind: .submit)))
            }
            questions = due.questions
            lessons = redo
        } catch {
            app.report(error)
        }
        isLoaded = true
    }

    /// What is due now that the libraries in the sidebar still hold: the questions their
    /// lessons still ask, interleaved, and the lessons. The review shows exactly these, and the
    /// sidebar counts them, so a question a library update dropped is never counted.
    static func due(_ learner: LearnerDatabase, in app: AppModel) async throws
        -> (questions: [Question], lessons: [(ref: LessonRef, summary: LessonSummary)]) {
        var asked: [Question] = []
        var quizzes: [LessonRef: Quiz] = [:]
        for due in try await learner.dueQuestions() {
            guard let library = app.library(due.ref.library), let db = library.database,
                  let summary = library.catalog.summary(due.ref.lesson) else { continue }
            if quizzes[due.ref] == nil { quizzes[due.ref] = try await db.quiz(due.ref.lesson) }
            guard let question = quizzes[due.ref]?.questions.first(where: { $0.position == due.position }) else { continue }
            asked.append(Question(ref: due.ref, lessonTitle: summary.title, question: question))
        }
        var lessons: [(ref: LessonRef, summary: LessonSummary)] = []
        for due in try await learner.dueLessons() {
            if let summary = app.summary(due.ref) { lessons.append((due.ref, summary)) }
        }
        return (asked, lessons)
    }

    /// Records the answer to a due question: right moves it along its schedule, wrong starts
    /// it over at a day. Quitting waits for the write, so an answer given just before is kept.
    func answer(_ ref: LessonRef, position: Int, answer: String, correct: Bool) {
        let learner = learner
        learner.pendingWrites.start { [weak app] in
            do {
                try await learner.recordQuestionReview(ref, position: position, answer: answer, correct: correct)
                await app?.refreshDueCount()
            } catch {
                app?.report(error)
            }
        }
    }

    /// Redoes a lesson due for review, from its starter code.
    func redo(_ ref: LessonRef) {
        app?.redo(ref)
    }
}
