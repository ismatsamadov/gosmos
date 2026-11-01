import Foundation

// The teaching rules: hints, the reference solution, and spaced review. Pure values and
// functions; the learner database stores their state.

/// A lesson's hint (lesson_hint).
public struct Hint: Hashable, Sendable {
    public let position: Int
    public let text: String

    public init(position: Int, text: String) {
        self.position = position
        self.text = text
    }
}

/// Which of a lesson's hints show, and whether the next may. Hints come one at a time, in
/// order, and one you revealed stays revealed. The last one waits for a failed Submit: you
/// try first, then get the strongest nudge.
public struct HintState: Equatable, Sendable {
    public let hints: [Hint]
    public private(set) var revealed: Set<Int>
    public var failedSubmits: Int

    public init(hints: [Hint], revealed: Set<Int> = [], failedSubmits: Int = 0) {
        self.hints = hints.sorted { $0.position < $1.position }
        self.revealed = revealed
        self.failedSubmits = failedSubmits
    }

    /// The hints you revealed, in order.
    public var shown: [Hint] { hints.filter { revealed.contains($0.position) } }

    /// The first hint not revealed yet.
    public var next: Hint? { hints.first { !revealed.contains($0.position) } }

    /// The next hint is the last one, and no Submit has failed yet.
    public var nextIsLocked: Bool { next != nil && next == hints.last && failedSubmits == 0 }

    public var canReveal: Bool { next != nil && !nextIsLocked }

    /// Reveals the next hint, if it may come now, and returns it.
    public mutating func reveal() -> Hint? {
        guard canReveal, let hint = next else { return nil }
        revealed.insert(hint.position)
        return hint
    }
}

/// The reference solution's version of a file you edit (solution_file).
public struct SolutionFile: Hashable, Sendable {
    public let name: String
    public let content: String

    public init(name: String, content: String) {
        self.name = name
        self.content = content
    }
}

/// What looking at a lesson's reference solution does. Before the lesson is done, the look is
/// recorded, and the pass that follows still counts but puts the lesson in review. After a
/// pass, comparing is free and changes nothing. Struggling counts the same: a pass after three
/// failed Submits or more also brings the lesson back for review.
public enum SolutionPolicy {
    public static let struggleThreshold = 3

    public static func recordsView(lessonDone: Bool) -> Bool { !lessonDone }

    public static func passNeedsReview(viewedSolution: Bool, failedSubmits: Int) -> Bool {
        viewedSolution || failedSubmits >= struggleThreshold
    }
}

/// Spaced review. An item first comes back a day after it is learned, then 3, 7, 21 and 60
/// days after each review it passes; a lapse (a failed review) starts it over at a day.
public enum ReviewSchedule {
    public static let intervals = [1, 3, 7, 21, 60]

    public static var first: Int { intervals[0] }

    /// The interval after a review.
    public static func next(after interval: Int, passed: Bool) -> Int {
        guard passed else { return first }
        return intervals.first { $0 > interval } ?? intervals[intervals.count - 1]
    }

    /// When an item `days` days out is due.
    public static func due(in days: Int, from now: Date) -> Date {
        now.addingTimeInterval(TimeInterval(days) * 86_400)
    }
}

/// A quiz question due for review (review_question).
public struct DueQuestion: Hashable, Sendable {
    public let ref: LessonRef
    public let position: Int
    public let dueAt: String

    public init(ref: LessonRef, position: Int, dueAt: String) {
        self.ref = ref
        self.position = position
        self.dueAt = dueAt
    }
}

/// A lesson due to redo (review_lesson).
public struct DueLesson: Hashable, Sendable {
    public let ref: LessonRef
    public let dueAt: String

    public init(ref: LessonRef, dueAt: String) {
        self.ref = ref
        self.dueAt = dueAt
    }
}

public enum ReviewQueue {
    /// Due questions, interleaved across lessons: lessons take turns, one question each, so
    /// two questions of one lesson are apart whenever others are due. Lessons go in the order
    /// their first question fell due; within a lesson, the longest overdue comes first.
    public static func interleave(_ due: [DueQuestion]) -> [DueQuestion] {
        func earlier(_ a: DueQuestion, _ b: DueQuestion) -> Bool {
            (a.dueAt, a.ref.library, a.ref.lesson, a.position) < (b.dueAt, b.ref.library, b.ref.lesson, b.position)
        }
        // One queue per lesson, made in the order each lesson's first question fell due.
        var queues: [[DueQuestion]] = []
        var queueOf: [LessonRef: Int] = [:]
        for q in due.sorted(by: earlier) {
            if let i = queueOf[q.ref] {
                queues[i].append(q)
            } else {
                queueOf[q.ref] = queues.count
                queues.append([q])
            }
        }
        var out: [DueQuestion] = []
        var round = 0
        while out.count < due.count {
            for queue in queues where round < queue.count {
                out.append(queue[round])
            }
            round += 1
        }
        return out
    }
}
