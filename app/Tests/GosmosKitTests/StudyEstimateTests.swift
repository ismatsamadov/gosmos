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

    @Test func roundsLongEstimatesToWholeHours() {
        // 31 and 32 lessons of 20 minutes: 10 hr 20 min and 10 hr 40 min.
        #expect(StudyEstimate(lessons: Array(repeating: lessons[2], count: 31), completed: []).description == "About 10 hr left")
        #expect(StudyEstimate(lessons: Array(repeating: lessons[2], count: 32), completed: []).description == "About 11 hr left")
        // Just under ten hours keeps its minutes.
        #expect(StudyEstimate(lessons: Array(repeating: lessons[1], count: 39), completed: []).description == "About 9 hr 45 min left")
    }
}
