import Foundation

/// Rewards are a projection of saved completions, so reopening or resubmitting a lesson
/// never earns it twice. Library content and the learner schema stay unchanged.
public struct LearningProgress: Sendable {
    public static let lessonXP = 100
    public static let rankXP = 500

    public struct Library: Sendable {
        public let id: String
        public let catalog: Catalog
        public init(id: String, catalog: Catalog) { self.id = id; self.catalog = catalog }
    }

    public struct Achievement: Identifiable, Sendable {
        public let id: String
        public let title: String
        public let detail: String
        public let symbol: String
        public let value: Int
        public let target: Int
        public var isUnlocked: Bool { value >= target }
        public var fraction: Double { min(1, Double(value) / Double(target)) }
    }

    public let completedLessons: Int
    public let completedChapters: Int
    public let completedCourses: Int
    public let codeLessons: Int
    public let quizzes: Int

    public init(libraries: [Library], done: [String: Set<String>]) {
        var lessons: Set<LessonRef> = [], code: Set<LessonRef> = [], quizzes: Set<LessonRef> = []
        var chapters: Set<String> = [], courses: Set<String> = []
        for library in libraries {
            let finished = done[library.id] ?? []
            for course in library.catalog.courses {
                for lesson in course.lessons where finished.contains(lesson.id) {
                    let ref = LessonRef(library: library.id, lesson: lesson.id)
                    lessons.insert(ref)
                    if lesson.type == .choice { quizzes.insert(ref) } else { code.insert(ref) }
                }
                for chapter in course.chapters {
                    let required = chapter.lessons.filter(\.required)
                    if !required.isEmpty && required.allSatisfy({ finished.contains($0.id) }) {
                        chapters.insert(library.id + "/" + chapter.id)
                    }
                }
                if !course.requiredLessons.isEmpty && course.requiredLessons.allSatisfy({ finished.contains($0.id) }) {
                    courses.insert(library.id + "/" + course.id)
                }
            }
        }
        completedLessons = lessons.count
        completedChapters = chapters.count
        completedCourses = courses.count
        codeLessons = code.count
        self.quizzes = quizzes.count
    }

    public var xp: Int { completedLessons * Self.lessonXP }
    public var level: Int { xp / Self.rankXP + 1 }
    public var xpInLevel: Int { xp % Self.rankXP }
    public var xpToNextLevel: Int { Self.rankXP - xpInLevel }
    public var levelFraction: Double { Double(xpInLevel) / Double(Self.rankXP) }
    public var rank: String {
        let ranks = ["Cadet", "Explorer", "Navigator", "Pilot", "Commander", "Captain"]
        return ranks[min(level - 1, ranks.count - 1)]
    }

    public var achievements: [Achievement] {
        [
            Achievement(id: "liftoff", title: "Liftoff", detail: "Complete your first lesson", symbol: "paperplane.fill", value: completedLessons, target: 1),
            Achievement(id: "builder", title: "Go Builder", detail: "Complete a coding lesson", symbol: "curlybraces", value: codeLessons, target: 1),
            Achievement(id: "thinker", title: "Sharp Mind", detail: "Complete a quiz", symbol: "brain", value: quizzes, target: 1),
            Achievement(id: "explorer", title: "Explorer", detail: "Complete 5 lessons", symbol: "safari.fill", value: completedLessons, target: 5),
            Achievement(id: "chapter", title: "New Horizons", detail: "Finish a chapter’s required lessons", symbol: "flag.checkered", value: completedChapters, target: 1),
            Achievement(id: "voyager", title: "Voyager", detail: "Complete 25 lessons", symbol: "sparkles", value: completedLessons, target: 25),
        ]
    }
}
