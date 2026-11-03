import Testing
@testable import GosmosKit

@Suite struct LearningProgressTests {
    private func catalog(_ count: Int = 6, optional: Bool = false) -> Catalog {
        let lessons = (0..<count).map {
            LessonSummary(id: "lesson-\($0)", chapterID: "chapter", title: "Mission \($0)",
                          type: $0 == 0 ? .choice : .codeTests, required: !optional || $0 != count - 1)
        }
        return Catalog(courses: [Course(id: "course", title: "Go", kind: "Course", chapters: [
            Chapter(id: "chapter", courseID: "course", title: "First steps", lessons: lessons)
        ])])
    }

    @Test func aNewLearnerStartsAtLevelOneWithLockedBadges() {
        let progress = LearningProgress(libraries: [.init(id: "library", catalog: catalog())], done: [:])
        #expect(progress.xp == 0 && progress.level == 1 && progress.rank == "Cadet")
        #expect(progress.xpToNextLevel == 500 && progress.levelFraction == 0)
        #expect(progress.completedCourses == 0 && progress.completedChapters == 0)
        #expect(progress.achievements.allSatisfy { !$0.isUnlocked })
    }

    @Test func rewardsUseUniqueKnownCompletionsAndSurviveReconstruction() {
        let library = LearningProgress.Library(id: "library", catalog: catalog())
        let saved: [String: Set<String>] = ["library": ["lesson-0", "lesson-1", "removed-lesson"], "missing-library": ["lesson-2"]]
        let progress = LearningProgress(libraries: [library, library], done: saved)
        #expect(progress.xp == 200 && progress.completedLessons == 2)
        #expect(progress.codeLessons == 1 && progress.quizzes == 1)
        #expect(progress.achievements.filter(\.isUnlocked).map(\.id) == ["liftoff", "builder", "thinker"])
        let reopened = LearningProgress(libraries: [library], done: saved)
        #expect(reopened.xp == progress.xp && reopened.level == progress.level)
        let undone = LearningProgress(libraries: [library], done: ["library": ["lesson-0"]])
        #expect(undone.xp == 100 && undone.codeLessons == 0)
    }

    @Test func levelChangesAtFiveLessonsAndOptionalLessonsStillEarnXP() {
        let library = LearningProgress.Library(id: "library", catalog: catalog(optional: true))
        let four = LearningProgress(libraries: [library], done: ["library": Set((0..<4).map { "lesson-\($0)" })])
        #expect(four.level == 1 && four.xpInLevel == 400 && four.xpToNextLevel == 100)
        let five = LearningProgress(libraries: [library], done: ["library": Set((0..<5).map { "lesson-\($0)" })])
        #expect(five.level == 2 && five.rank == "Explorer" && five.xpInLevel == 0)
        #expect(five.xpToNextLevel == 500 && five.completedChapters == 1 && five.completedCourses == 1)
        let six = LearningProgress(libraries: [library], done: ["library": Set((0..<6).map { "lesson-\($0)" })])
        #expect(six.xp == 600 && six.levelFraction == 0.2)
    }

    @Test func identicalLessonIDsInDifferentLibrariesCountSeparately() {
        let libraries = ["a", "b"].map { LearningProgress.Library(id: $0, catalog: catalog(1)) }
        let progress = LearningProgress(libraries: libraries, done: ["a": ["lesson-0"], "b": ["lesson-0"]])
        #expect(progress.xp == 200 && progress.completedChapters == 2 && progress.completedCourses == 2)
    }

    @Test func emptyCoursesAndOptionalOnlyChaptersDoNotUnlockCompletion() {
        let optional = LessonSummary(id: "bonus", chapterID: "optional", title: "Bonus", type: .code, required: false)
        let catalog = Catalog(courses: [
            Course(id: "empty", title: "Empty", kind: "Course"),
            Course(id: "bonus", title: "Bonus", kind: "Course", chapters: [Chapter(id: "optional", courseID: "bonus", title: "Optional", lessons: [optional])])
        ])
        let progress = LearningProgress(libraries: [.init(id: "a", catalog: catalog)], done: ["a": ["bonus"]])
        #expect(progress.xp == 100 && progress.completedChapters == 0 && progress.completedCourses == 0)
    }

    @Test func laterLevelsAndBadgesStayWithinTheirBounds() {
        let catalog = catalog(40)
        let progress = LearningProgress(libraries: [.init(id: "a", catalog: catalog)], done: ["a": Set(catalog.courses.flatMap(\.lessons).map(\.id))])
        #expect(progress.level == 9 && progress.rank == "Captain")
        #expect(progress.achievements.allSatisfy { $0.isUnlocked && $0.fraction == 1 })
    }
}
