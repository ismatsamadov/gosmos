import Foundation
import GosmosKit
import Testing
@testable import Gosmos

/// The shelf: the lesson Mission Control offers, what your finished lessons complete, and the
/// stories its libraries hold.
@Suite(.serialized) @MainActor struct ShelfModelTests {
    /// The lesson opened last while it is unfinished; once it is done, the next unfinished one
    /// in its course, after it and then before it; then the first unfinished one on the shelf.
    @Test func missionControlOffersTheLessonToGoOnWith() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let shelf = app.shelf
        #expect(shelf.recommendedLesson == fixture.ref("first"))
        shelf.resumeLesson = fixture.ref("optional")
        #expect(shelf.recommendedLesson == fixture.ref("optional"))
        shelf.setDone(fixture.ref("optional"), true)
        #expect(shelf.recommendedLesson == fixture.ref("quiz"))
        shelf.setDone(fixture.ref("quiz"), true)
        #expect(shelf.recommendedLesson == fixture.ref("first"))
        shelf.setDone(fixture.ref("first"), true)
        #expect(shelf.recommendedLesson == fixture.ref("second"))
        shelf.setDone(fixture.ref("second"), true)
        #expect(shelf.recommendedLesson == nil)
        // A lesson opened last that its library no longer has is passed over.
        shelf.resumeLesson = LessonRef(library: "flow-library", lesson: "gone")
        shelf.setDone(fixture.ref("quiz"), false)
        #expect(shelf.recommendedLesson == fixture.ref("quiz"))
        await app.prepareToQuit()
    }

    /// A chapter is complete once its required lessons are done; an optional one doesn't count.
    @Test func finishedLessonsCompleteTheirChapters() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let shelf = app.shelf
        let course = try #require(shelf.library("flow-library")?.catalog.courses.first { $0.id == "course" })
        let chapter = course.chapters[0]
        #expect(!shelf.isComplete(chapter, in: "flow-library"))
        #expect(shelf.isLastToDo(fixture.ref("first")) && !shelf.isLastToDo(fixture.ref("optional")))
        shelf.setDone(fixture.ref("first"), true)
        #expect(shelf.isComplete(chapter, in: "flow-library") && shelf.doneCount(chapter, in: "flow-library") == 1)
        let progress = shelf.progress(of: course, in: "flow-library")
        #expect(progress.done == 1 && progress.total == 2)
        #expect(shelf.learningProgress.completedLessons == 1)
        await app.prepareToQuit()
    }

    /// The courses' stories and the field guide's words come from the libraries on the shelf.
    @Test func theStoriesAreTheLibrariesOwn() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        try fixture.changeLibrary("""
            INSERT INTO chapter_story VALUES ('chapter-one', 'The first week.', 'A week done.');
            INSERT INTO story_term VALUES ('second-course', 'Uplink', 'The link up to a satellite.', 1);
            """)
        let app = fixture.app()
        await app.start()
        let shelf = app.shelf
        #expect(shelf.story(for: fixture.ref("first"))?.pages.count == 3)
        #expect(shelf.story(for: fixture.ref("second"))?.term("Uplink") != nil)
        let chapters = try #require(shelf.library("flow-library")?.catalog.courses.first?.chapters)
        #expect(shelf.chapterStory(chapters[0], in: "flow-library")?.briefing == "The first week.")
        #expect(shelf.chapterStory(chapters[1], in: "flow-library") == nil)
        #expect(shelf.hasFieldGuide)
        await app.prepareToQuit()
    }

    /// A library added in place of one with its uuid that can't be read takes that one's place;
    /// any other goes after the others.
    @Test func anAddedLibraryReplacesTheOneThatCantBeRead() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let shelf = app.shelf
        let url = fixture.directory.appendingPathComponent("moved.sqlite")
        shelf.add(ShelfLibrary(id: "moved", title: "Moved", url: url, isBuiltIn: false, error: .missing(path: url.path)))
        shelf.add(ShelfLibrary(id: "other", title: "Other", url: url, isBuiltIn: false, error: .missing(path: url.path)))
        shelf.add(ShelfLibrary(id: "moved", title: "Found again", url: url, isBuiltIn: false))
        #expect(shelf.libraries.map(\.id) == ["flow-library", "moved", "other"])
        #expect(shelf.library("moved")?.title == "Found again")
        await app.prepareToQuit()
    }
}
