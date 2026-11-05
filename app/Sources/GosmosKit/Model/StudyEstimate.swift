import Foundation

/// A planning estimate for unfinished required lessons, rather than a measure of time spent.
public struct StudyEstimate: Sendable {
    public let remainingLessons: Int
    public let minutes: Int

    public init(lessons: [LessonSummary], completed: Set<String>) {
        let remaining = lessons.filter { $0.required && !completed.contains($0.id) }
        remainingLessons = remaining.count
        minutes = remaining.reduce(0) { $0 + Self.minutes(for: $1.type) }
    }

    /// A rough starting point, not a deadline or an estimate of the learner's pace.
    public static func minutes(for type: LessonType) -> Int {
        switch type {
        case .choice: 5
        case .code: 15
        case .codeTests: 20
        }
    }

    public var description: String {
        if minutes == 0 { return "No required lessons left" }
        if minutes < 60 { return "About \(minutes) min left" }
        // From ten hours on, minutes would claim a precision a rough estimate doesn't have.
        if minutes >= 600 { return "About \((minutes + 30) / 60) hr left" }
        let hours = minutes / 60, rest = minutes % 60
        if rest == 0 { return "About \(hours) hr left" }
        return "About \(hours) hr \(rest) min left"
    }
}
