import Foundation

/// A lesson's Markdown as blocks, parsed once off the main actor: the readme, and a quiz's
/// questions.
public struct ParsedLesson: Sendable {
    public var readme: [MDBlock]
    public var questions: [[MDBlock]] = []

    public init(_ lesson: LessonDetail) {
        readme = Readme.parse(lesson.readme, droppingTitle: lesson.title)
        if case .choice(let quiz) = lesson.content {
            questions = quiz.questions.map { Readme.parse($0.question) }
        }
    }
}

/// Parsed lessons kept for the recently opened ones, by library and lesson.
public actor ReadmeCache {
    private var parsed: [LessonRef: ParsedLesson] = [:]
    private var order: [LessonRef] = []

    public init() {}

    public func parsed(for lesson: LessonDetail) -> ParsedLesson {
        if let cached = parsed[lesson.ref] { return cached }
        let result = ParsedLesson(lesson)
        parsed[lesson.ref] = result
        order.append(lesson.ref)
        if order.count > 120 { parsed[order.removeFirst()] = nil }
        return result
    }
}
