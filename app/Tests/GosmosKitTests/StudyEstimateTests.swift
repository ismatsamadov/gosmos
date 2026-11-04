import Testing
@testable import GosmosKit

@Suite struct StudyEstimateTests {
    private let lessons = [
        LessonSummary(id: "quiz", chapterID: "c", title: "Quiz", type: .choice, required: true),
        LessonSummary(id: "code", chapterID: "c", title: "Code", type: .code, required: true),
        LessonSummary(id: "tests", chapterID: "c", title: "Tests", type: .codeTests, required: true),
        LessonSummary(id: "extra", chapterID: "c", title: "Extra", type: .codeTests, required: false),
    ]

    @Test func estimatesOnlyUnfinishedRequiredLessons() {
        let fresh = StudyEstimate(lessons: lessons, completed: [])
        #expect(fresh.remainingLessons == 3 && fresh.minutes == 40)
        let partial = StudyEstimate(lessons: lessons, completed: ["quiz", "extra", "unknown"])
        #expect(partial.remainingLessons == 2 && partial.minutes == 35)
        let complete = StudyEstimate(lessons: lessons, completed: ["quiz", "code", "tests"])
        #expect(complete.remainingLessons == 0 && complete.description == "No required lessons left")
    }

    @Test func formatsHoursAndEmptyCourses() {
        let many = StudyEstimate(lessons: Array(repeating: lessons[2], count: 4), completed: [])
        #expect(many.description == "About 1 hr 20 min left")
        #expect(StudyEstimate(lessons: [], completed: []).minutes == 0)
    }
}
