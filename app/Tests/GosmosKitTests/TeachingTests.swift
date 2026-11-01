import Foundation
import GRDB
import Testing
@testable import GosmosKit

/// The teaching rules and what the learner database keeps of them: hints, the reference
/// solution, spaced review of quiz questions and lessons, and failures that say what failed.
@Suite struct TeachingTests {
    let hints = [Hint(position: 0, text: "Read the error."), Hint(position: 1, text: "Look at the loop."), Hint(position: 2, text: "i <= len is one too far.")]

    @Test func hintsComeOneAtATimeAndTheLastAfterAFailedSubmit() {
        var state = HintState(hints: hints.reversed()) // in position order whatever the input order
        #expect(state.shown.isEmpty && state.next?.position == 0 && state.canReveal)
        #expect(state.reveal()?.text == "Read the error.")
        #expect(state.reveal()?.text == "Look at the loop.")
        // The last one waits until a Submit has failed.
        #expect(state.nextIsLocked && !state.canReveal && state.reveal() == nil)
        state.failedSubmits = 1
        #expect(!state.nextIsLocked && state.reveal()?.position == 2)
        #expect(state.shown.map(\.position) == [0, 1, 2] && state.next == nil && !state.canReveal)

        // Hints you revealed before stay revealed, the last one too.
        let again = HintState(hints: hints, revealed: [0, 1, 2], failedSubmits: 0)
        #expect(again.shown.count == 3 && !again.nextIsLocked)
        // A lesson with one hint: it is the last, so it waits too.
        var one = HintState(hints: [hints[0]])
        #expect(one.reveal() == nil && one.nextIsLocked)
        #expect(HintState(hints: []).next == nil)
    }

    @Test func revealedHintsAndFailedSubmitsAreKept() async throws {
        let dir = try Fixtures.tempDir("hints")
        defer { try? FileManager.default.removeItem(at: dir) }
        let learner = try await Fixtures.learner(in: dir)
        let ref = Fixtures.ref(Fixtures.codeLesson)
        try await learner.revealHint(ref, position: 0)
        try await learner.revealHint(ref, position: 0) // once
        try await learner.revealHint(ref, position: 1)
        #expect(try await learner.revealedHints(ref) == [0, 1])
        #expect(try await learner.revealedHints(Fixtures.ref(Fixtures.testsLesson)).isEmpty)
        try await learner.recordSubmit(ref, passed: false, output: "no", durationMs: 5, reviewable: true)
        let state = HintState(hints: hints, revealed: try await learner.revealedHints(ref),
                              failedSubmits: try await learner.failedAttempts(ref, kind: .submit))
        #expect(state.canReveal && state.next?.position == 2)
    }

    /// Looking at the solution before passing is recorded, and the pass that follows counts
    /// but brings the lesson back for review a day later. After a pass, comparing is free.
    @Test func lookingAtTheSolution() async throws {
        let dir = try Fixtures.tempDir("solution")
        defer { try? FileManager.default.removeItem(at: dir) }
        let learner = try await Fixtures.learner(in: dir)
        let ref = Fixtures.ref(Fixtures.codeLesson)
        let now = Date(timeIntervalSince1970: 1_803_000_000)
        #expect(SolutionPolicy.recordsView(lessonDone: false) && !SolutionPolicy.recordsView(lessonDone: true))

        #expect(try await learner.recordSolutionView(ref, now: now))
        #expect(try await !learner.recordSolutionView(ref, now: now.addingTimeInterval(60))) // the first look counts
        #expect(try await learner.viewedSolution(ref))
        // The pass still finishes the lesson, and puts it in review, due in a day.
        #expect(try await learner.recordSubmit(ref, passed: true, output: "ok", durationMs: 10, reviewable: true, now: now.addingTimeInterval(120)))
        #expect(try await learner.done()[Fixtures.libraryID] == [Fixtures.codeLesson])
        let review = try Fixtures.rows(in: dir, "SELECT due_at, interval_days, reps, lapses FROM review_lesson")
        #expect(review.count == 1 && review[0]["interval_days"] as Int == 1 && review[0]["reps"] as Int == 0)
        #expect(review[0]["due_at"] as String == Clock.string(now.addingTimeInterval(120 + 86_400)))
        #expect(try await learner.dueLessons(now: now.addingTimeInterval(3_600)).isEmpty)
        #expect(try await learner.dueLessons(now: now.addingTimeInterval(120 + 86_400)).map(\.ref) == [ref])

        // Another lesson, passed first: looking afterwards records nothing and changes nothing.
        let other = Fixtures.ref(Fixtures.testsLesson)
        try await learner.recordSubmit(other, passed: true, output: "ok", durationMs: 10, reviewable: true, now: now)
        #expect(try await !learner.recordSolutionView(other, now: now.addingTimeInterval(60)))
        #expect(try await !learner.viewedSolution(other))
        #expect(try await !learner.recordSubmit(other, passed: true, output: "ok", durationMs: 10, reviewable: true, now: now.addingTimeInterval(120)))
        #expect(try Fixtures.rows(in: dir, "SELECT lesson_uuid FROM review_lesson").map { $0["lesson_uuid"] as String } == [Fixtures.codeLesson])
    }

    /// Three failed Submits or more before a pass bring a code lesson back too.
    @Test func strugglingBringsALessonBack() async throws {
        let dir = try Fixtures.tempDir("struggle")
        defer { try? FileManager.default.removeItem(at: dir) }
        let learner = try await Fixtures.learner(in: dir)
        let ref = Fixtures.ref(Fixtures.testsLesson)
        for _ in 0..<2 { try await learner.recordSubmit(ref, passed: false, output: "FAIL", durationMs: 5, reviewable: true) }
        #expect(SolutionPolicy.passNeedsReview(viewedSolution: false, failedSubmits: 2) == false)
        try await learner.recordSubmit(ref, passed: false, output: "FAIL", durationMs: 5, reviewable: true)
        #expect(SolutionPolicy.passNeedsReview(viewedSolution: false, failedSubmits: 3))
        #expect(try await learner.recordSubmit(ref, passed: true, output: "ok", durationMs: 5, reviewable: true))
        #expect(try Fixtures.rows(in: dir, "SELECT lesson_uuid FROM review_lesson").map { $0["lesson_uuid"] as String } == [Fixtures.testsLesson])
    }

    @Test func theReviewScheduleGrowsAndStartsOverAfterALapse() {
        #expect(ReviewSchedule.intervals == [1, 3, 7, 21, 60] && ReviewSchedule.first == 1)
        var interval = ReviewSchedule.first
        var seen = [interval]
        for _ in 0..<5 {
            interval = ReviewSchedule.next(after: interval, passed: true)
            seen.append(interval)
        }
        #expect(seen == [1, 3, 7, 21, 60, 60]) // 60 days is the longest
        #expect(ReviewSchedule.next(after: 21, passed: false) == 1)
        #expect(ReviewSchedule.next(after: 2, passed: true) == 3) // an interval off the list moves to the next one on it
        let now = Date(timeIntervalSince1970: 1_803_000_000)
        #expect(ReviewSchedule.due(in: 7, from: now) == now.addingTimeInterval(7 * 86_400))
    }

    /// A quiz question answered right the first time comes back a day later; reviews move it
    /// along the schedule, and a wrong answer is a lapse.
    @Test func quizQuestionsComeBack() async throws {
        let dir = try Fixtures.tempDir("quizreview")
        defer { try? FileManager.default.removeItem(at: dir) }
        let learner = try await Fixtures.learner(in: dir)
        let ref = Fixtures.ref(Fixtures.quizLesson)
        let day = 86_400.0, start = Date(timeIntervalSince1970: 1_803_000_000)
        try await learner.recordAnswers(ref, [QuizState.Answer(question: 0, text: "`fn`", right: false),
                                              QuizState.Answer(question: 1, text: "Yes", right: true)], now: start)
        try await learner.recordAnswers(ref, [QuizState.Answer(question: 0, text: "`func`", right: true)], now: start.addingTimeInterval(3_600))
        // Answering it right again later does not move it back to the start.
        try await learner.recordAnswers(ref, [QuizState.Answer(question: 1, text: "Yes", right: true)], now: start.addingTimeInterval(2 * 3_600))
        #expect(try await learner.dueQuestions(now: start.addingTimeInterval(day - 1)).isEmpty)
        let due = try await learner.dueQuestions(now: start.addingTimeInterval(day + 3_600))
        #expect(due.map(\.position) == [1, 0] && due.allSatisfy { $0.ref == ref })

        func schedule(_ position: Int) throws -> (interval: Int, reps: Int, lapses: Int, due: String) {
            let row = try #require(try Fixtures.rows(in: dir, "SELECT interval_days, reps, lapses, due_at FROM review_question WHERE question_position = ?", [position]).first)
            return (row["interval_days"], row["reps"], row["lapses"], row["due_at"])
        }
        let reviewed = start.addingTimeInterval(day + 3_600)
        try await learner.recordQuestionReview(ref, position: 1, answer: "Yes", correct: true, now: reviewed)
        #expect(try schedule(1).interval == 3 && schedule(1).reps == 1 && schedule(1).due == Clock.string(reviewed.addingTimeInterval(3 * day)))
        try await learner.recordQuestionReview(ref, position: 0, answer: "`def`", correct: false, now: reviewed)
        #expect(try schedule(0).interval == 1 && schedule(0).lapses == 1 && schedule(0).due == Clock.string(reviewed.addingTimeInterval(day)))
        // Reviews are attempts too.
        let attempts = try Fixtures.rows(in: dir, "SELECT kind, passed FROM attempt WHERE kind = 'review' ORDER BY id")
        #expect(attempts.map { $0["passed"] as Int } == [1, 0])
    }

    @Test func lessonsDueToRedoMoveAlongTheSchedule() async throws {
        let dir = try Fixtures.tempDir("lessonreview")
        defer { try? FileManager.default.removeItem(at: dir) }
        let learner = try await Fixtures.learner(in: dir)
        let ref = Fixtures.ref(Fixtures.codeLesson)
        let day = 86_400.0, start = Date(timeIntervalSince1970: 1_803_000_000)
        try await learner.recordSolutionView(ref, now: start)
        try await learner.recordSubmit(ref, passed: true, output: "ok", durationMs: 5, reviewable: true, now: start)
        try await learner.recordLessonReview(ref, passed: true, output: "ok", durationMs: 5, now: start.addingTimeInterval(day))
        let after = try #require(try Fixtures.rows(in: dir, "SELECT interval_days, reps, due_at FROM review_lesson").first)
        #expect(after["interval_days"] as Int == 3 && after["reps"] as Int == 1)
        #expect(after["due_at"] as String == Clock.string(start.addingTimeInterval(4 * day)))
        try await learner.recordLessonReview(ref, passed: false, output: "FAIL", durationMs: 5, now: start.addingTimeInterval(4 * day))
        let lapsed = try #require(try Fixtures.rows(in: dir, "SELECT interval_days, lapses FROM review_lesson").first)
        #expect(lapsed["interval_days"] as Int == 1 && lapsed["lapses"] as Int == 1)
        // A redo never touches your saved code or the lesson's done mark.
        #expect(try await learner.codeFiles(for: ref).isEmpty)
        #expect(try await learner.done()[Fixtures.libraryID] == [Fixtures.codeLesson])
    }

    /// Due questions take turns across lessons, so two of one lesson are apart while others
    /// are due.
    @Test func dueQuestionsAreInterleaved() {
        func q(_ lesson: String, _ position: Int, _ due: String) -> DueQuestion {
            DueQuestion(ref: LessonRef(library: "lib", lesson: lesson), position: position, dueAt: due)
        }
        let due = [q("a", 0, "2027-03-01T08:00:00Z"), q("a", 1, "2027-03-01T08:00:00Z"), q("a", 2, "2027-03-02T08:00:00Z"),
                   q("b", 0, "2027-03-01T09:00:00Z"), q("c", 3, "2027-02-28T09:00:00Z"), q("c", 4, "2027-03-03T09:00:00Z")]
        let order = ReviewQueue.interleave(due.shuffled())
        #expect(order.map { "\($0.ref.lesson)\($0.position)" } == ["c3", "a0", "b0", "c4", "a1", "a2"])
        #expect(ReviewQueue.interleave([]).isEmpty)
    }

    @Test func theSolutionSitsBesideYourCodeLineByLine() throws {
        let yours = ["package main", "", "func main() {", "\tfmt.Println(\"hi\")", "}"]
        let solution = ["package main", "", "import \"fmt\"", "", "func main() {", "\tfmt.Println(\"Hello, gopher\")", "\tfmt.Printf(\"%d\\n\", 2+2)", "}"]
        // Matching lines side by side; between them, the lines that differ paired in order and
        // the rest alone.
        #expect(try LineDiff.rows(yours, solution) == [
            LineDiff.Row(left: 0, right: 0, differs: false),
            LineDiff.Row(left: 1, right: 1, differs: false),
            LineDiff.Row(left: nil, right: 2, differs: true),
            LineDiff.Row(left: nil, right: 3, differs: true),
            LineDiff.Row(left: 2, right: 4, differs: false),
            LineDiff.Row(left: 3, right: 5, differs: true),
            LineDiff.Row(left: nil, right: 6, differs: true),
            LineDiff.Row(left: 4, right: 7, differs: false),
        ])
        #expect(try LineDiff.rows(["a", "b", "c"], ["a", "c"]) == [
            LineDiff.Row(left: 0, right: 0, differs: false), LineDiff.Row(left: 1, right: nil, differs: true), LineDiff.Row(left: 2, right: 1, differs: false),
        ])
        // Spacing alone is no difference.
        let respaced = yours.map { $0.replacingOccurrences(of: "\t", with: "    ") + "  " }
        #expect(try LineDiff.rows(yours, respaced).allSatisfy { !$0.differs })
        #expect(try LineDiff.rows([], ["x"]) == [LineDiff.Row(left: nil, right: 0, differs: true)] && LineDiff.rows([], []).isEmpty)

        // Past four million line pairs, lines are compared by position.
        let long = (0..<2001).map { "line \($0)" }
        var changed = long
        changed[1000] = "changed"
        let rows = try LineDiff.rows(long, changed)
        #expect(rows.count == 2001 && rows[1000] == LineDiff.Row(left: 1000, right: 1000, differs: true))
        #expect(rows.filter(\.differs).count == 1)
    }

    /// A long comparison gives up once its task is cancelled, as when the sheet showing it
    /// closes first; left to run, the same comparison ends.
    @Test func aCancelledComparisonGivesUp() async throws {
        let yours = (0..<600).map { "line \($0)" }, solution = yours.map { $0 + " changed" }
        let cancelled = Task.detached {
            withUnsafeCurrentTask { $0?.cancel() }
            return try LineDiff.rows(yours, solution)
        }
        await #expect(throws: CancellationError.self) { try await cancelled.value }
        let rows = try await Task.detached { try LineDiff.rows(yours, solution) }.value
        #expect(rows.count == 600 && rows.allSatisfy(\.differs))
    }

    @Test func failingTestsAreNamedWithTheirMessages() {
        let output = """
        --- FAIL: TestAdd (0.00s)
            main_test.go:16: add(1000, 1) = 0, want 1001
            main_test.go:16: add(2, 3) = 0, want 5
        --- FAIL: TestSplit (0.00s)
            --- FAIL: TestSplit/empty (0.00s)
                split_test.go:9: got [""], want []
                    because an empty string has no parts
        --- PASS: TestJoin (0.00s)
        FAIL
        exit status 1
        FAIL\tlesson\t0.005s
        """
        let failed = FailedTest.parse(output)
        #expect(failed.map(\.name) == ["TestAdd", "TestSplit/empty"])
        #expect(failed[0].messages == ["main_test.go:16: add(1000, 1) = 0, want 1001", "main_test.go:16: add(2, 3) = 0, want 5"])
        #expect(failed[1].messages == ["split_test.go:9: got [\"\"], want []\nbecause an empty string has no parts"])
        #expect(FailedTest.parse("ok  \tlesson\t0.003s\n").isEmpty)
        // A test that failed without a message of its own still counts.
        #expect(FailedTest.parse("--- FAIL: TestQuiet (0.00s)\nFAIL\n") == [FailedTest(name: "TestQuiet", messages: [])])
    }
}
