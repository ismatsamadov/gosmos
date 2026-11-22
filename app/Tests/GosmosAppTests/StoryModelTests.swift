import GosmosKit
import Testing
@testable import Gosmos

/// The story as you go through it: who meets an opening at launch, which opening Replay shows,
/// where an opening leads, its pages, and the name on your badge.
@Suite(.serialized) @MainActor struct StoryModelTests {
    /// Only a new learner meets the opening at launch: one who has seen no opening, and opened
    /// or finished no lesson.
    @Test func onlyANewLearnerMeetsTheOpeningAtLaunch() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let story = app.story
        #expect(story.opening == StoryModel.Opening(library: "flow-library", course: "course", page: 0))
        #expect(story.openingAtLaunch == story.opening)
        app.setDone(fixture.ref("quiz"), true)
        #expect(story.openingAtLaunch == nil)
        app.setDone(fixture.ref("quiz"), false)
        #expect(story.openingAtLaunch != nil)
        app.open(fixture.ref("optional"))
        try await fixture.wait { app.resumeLesson == fixture.ref("optional") }
        #expect(story.openingAtLaunch == nil)
        await app.prepareToQuit()
    }

    @Test func anOpeningEndedIsSeenForGood() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.story.finishOpening()
        // Seen, before the lesson it opens has loaded.
        #expect(app.story.seenOpenings == ["flow-library/course"])
        #expect(app.resumeLesson == nil && app.story.openingAtLaunch == nil)
        try await fixture.wait { app.lesson?.ref == fixture.ref("first") }
        await app.prepareToQuit()

        let reopened = fixture.app()
        await reopened.start()
        #expect(reopened.story.seenOpenings == ["flow-library/course"])
        await reopened.prepareToQuit()
    }

    /// A course with no opening starts at its lesson, at launch and from Begin; Replay shows
    /// the first opening on the shelf instead.
    @Test func aCourseWithoutAnOpeningStartsAtItsLesson() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        try fixture.changeLibrary("DELETE FROM story_page WHERE course_uuid = 'course';")
        let app = fixture.app()
        await app.start()
        let story = app.story
        #expect(story.opening == nil && app.recommendedLesson == fixture.ref("first"))
        #expect(story.openingToReplay?.course == "second-course")
        let course = try #require(app.library("flow-library")?.catalog.courses.first { $0.id == "course" })
        story.begin(course, in: "flow-library")
        #expect(story.opening == nil)
        try await fixture.wait { app.lesson?.ref == fixture.ref("first") }
        await app.prepareToQuit()
    }

    /// After the opening it waited for, Continue Learning's lesson opens, unless it was done
    /// meanwhile: then the course's next unfinished one, and Mission Control once none is left.
    @Test func anOpeningLeadsToTheLessonAskedForWhileItIsUnfinished() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let story = app.story
        story.continueLearning(fixture.ref("optional"))
        #expect(story.opening?.destination == fixture.ref("optional"))
        app.setDone(fixture.ref("optional"), true)
        story.finishOpening()
        try await fixture.wait { app.lesson?.ref == fixture.ref("first") && app.lesson?.code?.isLoaded == true }

        for lesson in ["first", "quiz"] { app.setDone(fixture.ref(lesson), true) }
        story.showOpening("course", in: "flow-library")
        story.finishOpening()
        try await fixture.wait { app.lesson == nil }
        #expect(app.selection == .dashboard && story.opening == nil)
        await app.prepareToQuit()
    }

    @Test func theOpeningsPagesStayInRange() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let story = app.story
        app.shownTerm = CourseStory.Term(term: "Uplink", meaning: "The link up to a satellite.")
        story.showOpening("course", in: "flow-library", page: 9)
        #expect(story.opening?.page == 2 && app.shownTerm == nil)
        story.showOpening("course", in: "flow-library", page: -1)
        story.openingPageBack()
        #expect(story.opening?.page == 0)
        // A course with no opening shows none, and leaves the one on screen.
        story.showOpening("no-such-course", in: "flow-library")
        #expect(story.opening?.course == "course")
        await app.prepareToQuit()
    }

    @Test func theBadgeNameIsTrimmedCutShortAndTakenOffWhenEmpty() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let story = app.story
        let learner = try #require(app.learner)
        story.setBadgeName("  " + String(repeating: "x", count: 50) + "\n")
        #expect(story.badgeName == String(repeating: "x", count: 40))
        await learner.pendingWrites.finish()
        #expect(try await learner.setting(StoryModel.badgeNameKey) == String(repeating: "x", count: 40))
        story.setBadgeName(" ")
        #expect(story.badgeName.isEmpty)
        await learner.pendingWrites.finish()
        #expect(try await learner.setting(StoryModel.badgeNameKey) == nil)
        await app.prepareToQuit()
    }
}
