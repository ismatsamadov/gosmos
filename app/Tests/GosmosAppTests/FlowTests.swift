import Foundation
import GRDB
import GosmosKit
import Testing
@testable import Gosmos

/// Real app navigation and persistence, with a small test library and temporary study data.
@Suite(.serialized) @MainActor struct FlowTests {
    @Test func theBadgeNameSurvivesNavigationAndRelaunch() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        #expect(app.opening?.page == 0)
        app.advanceOpening()
        #expect(app.opening?.page == 1)
        app.setBadgeName("  Ada Lovelace  ")
        app.advanceOpening()
        #expect(app.opening?.page == 2)
        app.openingPageBack()
        #expect(app.badgeName == "Ada Lovelace")
        app.finishOpening()
        try await fixture.wait { app.lesson?.ref == fixture.ref("first") && app.lesson?.code?.isLoaded == true }
        await app.prepareToQuit()

        let reopened = fixture.app()
        await reopened.start()
        #expect(reopened.badgeName == "Ada Lovelace")
        #expect(reopened.opening == nil)
        #expect(reopened.recommendedLesson == fixture.ref("first"))
        await reopened.prepareToQuit()
    }

    /// The first course's opening, finished, then a relaunch before the course is done: Mission
    /// Control, with the lesson to continue. The second course's opening waits until you begin
    /// that course.
    @Test func aRelaunchShowsNoOpeningOfACourseNotBegun() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        #expect(app.opening?.course == "course")
        app.finishOpening()
        try await fixture.wait { app.resumeLesson == fixture.ref("first") }
        await app.prepareToQuit()

        let reopened = fixture.app()
        await reopened.start()
        #expect(reopened.opening == nil)
        #expect(reopened.selection == .dashboard)
        #expect(reopened.recommendedLesson == fixture.ref("first"))
        let second = try #require(reopened.library("flow-library")?.catalog.courses.first { $0.id == "second-course" })
        reopened.begin(second, in: "flow-library")
        #expect(reopened.opening?.course == "second-course")
        await reopened.prepareToQuit()
    }

    /// Quitting halfway through the opening: you are still new, so the next launch shows it again.
    @Test func anOpeningLeftUnfinishedShowsAgainAtTheNextLaunch() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.advanceOpening()
        #expect(app.opening?.page == 1)
        await app.prepareToQuit()

        let reopened = fixture.app()
        await reopened.start()
        #expect(reopened.opening == AppModel.Opening(library: "flow-library", course: "course", page: 0))
        await reopened.prepareToQuit()
    }

    /// View › Replay the Opening replays the course of the lesson on screen, also once that course
    /// is done and Mission Control offers another; on Mission Control, the course it offers.
    @Test func replayingTheOpeningKeepsToTheCourseOnScreen() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        app.open(fixture.ref("second"))
        try await fixture.wait { app.lesson?.ref == fixture.ref("second") && app.resumeLesson == fixture.ref("second") }
        let lesson = try #require(app.lesson)
        lesson.markPassed(review: false) // the second course's only lesson: that course is done
        #expect(app.recommendedLesson == fixture.ref("first"))
        #expect(app.openingToReplay?.course == "second-course")
        app.replayOpening()
        #expect(app.opening?.course == "second-course")
        app.finishOpening()
        try await fixture.wait { app.lesson == nil }
        #expect(app.selection == .dashboard)
        #expect(app.openingToReplay?.course == "course")
        await app.prepareToQuit()
    }

    /// While the opening covers the window, no command acts on what it covers: here a lesson
    /// still on screen behind it, as when Replay the Opening leaves one that waits to be saved.
    /// They come back once the opening ends.
    @Test func theOpeningTurnsOffTheCommandsBehindIt() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        app.open(fixture.ref("optional"))
        try await fixture.wait { app.lesson?.ref == fixture.ref("optional") && app.lesson?.code?.isLoaded == true }
        app.open(fixture.ref("first"))
        try await fixture.wait { app.lesson?.ref == fixture.ref("first") && app.lesson?.code?.isLoaded == true }
        app.goBack()
        try await fixture.wait { app.lesson?.ref == fixture.ref("optional") && app.lesson?.code?.isLoaded == true }
        // Format needs a Go with gofmt, which a test may not have; the others need nothing.
        let editing: [EditorCommand] = [.toggleComment, .indent, .outdent]
        var commands = app.commands
        #expect(commands.run && commands.submit && commands.nextLesson && commands.previousLesson && commands.back && commands.forward)
        #expect(commands.review && commands.markDone && commands.resetToStarter && editing.allSatisfy(commands.editor))
        #expect(commands.missionControl && commands.replayOpening)

        app.showOpening("course", in: "flow-library")
        commands = app.commands
        #expect(commands.covered && app.lesson?.ref == fixture.ref("optional"))
        #expect(!commands.run && !commands.submit && !commands.nextLesson && !commands.previousLesson && !commands.back && !commands.forward)
        #expect(!commands.review && !commands.markDone && !commands.resetToStarter && !editing.contains(where: commands.editor))
        #expect(!commands.missionControl && !commands.fieldGuide && !commands.replayOpening)

        app.finishOpening()
        try await fixture.wait { app.lesson?.ref == fixture.ref("first") && app.lesson?.code?.isLoaded == true }
        #expect(app.commands.run && app.commands.back)
        await app.prepareToQuit()
    }

    /// A redo for review, then Back to the lesson it redoes: the lesson's views and its editor
    /// are made anew, never the redo's kept with its undo history, whose Undo would walk back
    /// through the redo's typing and save it as the lesson's code.
    @Test func aRedoThenBackGivesTheLessonViewsOfItsOwn() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        try await fixture.wait { app.lesson?.ref == fixture.ref("first") && app.lesson?.code?.isLoaded == true }
        app.showReview()
        try await fixture.wait { app.showsReview }
        app.redo(fixture.ref("first"))
        try await fixture.wait { app.lesson?.isReview == true && app.lesson?.code?.isLoaded == true }
        let redo = try #require(app.lesson?.code)
        let redoFile = try #require(redo.activeFile)
        let redoViews = try #require(app.lesson?.viewIdentity), redoEditor = redo.editorIdentity(of: redoFile)
        redoFile.text = "package main // typed in the redo\n"

        app.goBack()
        try await fixture.wait { app.lesson?.isReview == false && app.lesson?.code?.isLoaded == true }
        let lesson = try #require(app.lesson)
        let code = try #require(lesson.code)
        let file = try #require(code.activeFile)
        #expect(lesson.ref == fixture.ref("first"))
        #expect(lesson.viewIdentity != redoViews && code.editorIdentity(of: file) != redoEditor)
        #expect(file !== redoFile && file.text == "package main\n\nfunc main() {}\n")
        await app.prepareToQuit()
    }

    /// Text copied from a readme carries the word joiners that keep its inline code on one line,
    /// unseen. Pasted into Search, it searches for what it says.
    @Test func aSearchForTextCopiedFromAReadmeFindsWhatItSays() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.searchText = "Print\u{2060} hel\u{2060}lo"
        #expect(app.searchText == "Print hello")
        try await fixture.wait { !app.searching && !app.searchResults.isEmpty }
        #expect(app.searchResults.map(\.ref).contains(fixture.ref("first")))
        await app.prepareToQuit()
    }

    /// An update changes a lesson between two launches: a new starter, and a test file it no
    /// longer gives. Its folder, where you never typed, takes the new starter, and the dropped
    /// file goes, where go would still build it.
    @Test func aLessonChangedBetweenLaunchesGivesItsNewStarter() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        try fixture.changeLibrary("""
            INSERT INTO starter_file (lesson_uuid, name, position, content, is_hidden, is_readonly)
              VALUES ('first', 'extra_test.go', 1, 'package main\n', 0, 1);
            """)
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        try await fixture.wait { app.lesson?.ref == fixture.ref("first") && app.lesson?.code?.isLoaded == true }
        let folder = try #require(app.lesson?.code?.folder)
        let extra = folder.appendingPathComponent("extra_test.go")
        #expect(FileManager.default.fileExists(atPath: extra.path))
        await app.prepareToQuit()

        let starter = "package main\n\nfunc main() { /* the new starter */ }\n"
        try fixture.changeLibrary("""
            UPDATE starter_file SET content = '\(starter)' WHERE lesson_uuid = 'first' AND name = 'main.go';
            DELETE FROM starter_file WHERE lesson_uuid = 'first' AND name = 'extra_test.go';
            """)
        let reopened = fixture.app()
        await reopened.start()
        reopened.open(fixture.ref("first"))
        try await fixture.wait { reopened.lesson?.ref == fixture.ref("first") && reopened.lesson?.code?.isLoaded == true }
        #expect(reopened.lesson?.code?.files.map(\.text) == [starter])
        #expect(try String(contentsOf: folder.appendingPathComponent("main.go"), encoding: .utf8) == starter)
        #expect(!FileManager.default.fileExists(atPath: extra.path))
        #expect(try await reopened.learner?.codeFiles(for: fixture.ref("first")).isEmpty == true)
        await reopened.prepareToQuit()
    }

    /// A lesson whose code can't load is on screen, but the next launch doesn't resume it, nor
    /// does its folder count it as opened: the time of the last open is what a folder from
    /// before the folders' records is measured against. Once a retry loads it, it counts.
    @Test func aLessonWhoseCodeFailedToLoadIsNotTheLastOpened() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        try await fixture.wait { app.resumeLesson == fixture.ref("first") }
        // The work folder can't hold lesson folders: a file sits in its place.
        let work = fixture.directory.appendingPathComponent("study/work")
        try FileManager.default.removeItem(at: work)
        try "not a folder".write(to: work, atomically: true, encoding: .utf8)
        app.open(fixture.ref("optional"))
        try await fixture.wait { app.lesson?.ref == fixture.ref("optional") && app.lesson?.code?.loadError != nil && !app.isLoadingLesson }
        try await Task.sleep(for: .milliseconds(100))
        #expect(app.resumeLesson == fixture.ref("first"))
        #expect(try await app.learner?.lastOpened() == fixture.ref("first"))
        #expect(try await app.learner?.lastOpenedAt(fixture.ref("optional")) == nil)

        // Fixed, and Try Again: now it loads, and it is the lesson opened last, after a relaunch too.
        try FileManager.default.removeItem(at: work)
        let session = try #require(app.lesson?.code)
        await session.load()
        #expect(session.isLoaded && app.resumeLesson == fixture.ref("optional"))
        await app.prepareToQuit()
        let reopened = fixture.app()
        await reopened.start()
        #expect(reopened.resumeLesson == fixture.ref("optional") && reopened.recommendedLesson == fixture.ref("optional"))
        await reopened.prepareToQuit()
    }

    /// A Run asked for while the database is slow, just before you answer a conflict: its save
    /// meets the file your answer is settling, and the Run waits for the answer to land rather
    /// than being dropped.
    @Test func aRunMeetingAnAnswerStillLandingRunsAfterIt() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        try fixture.changeLibrary("""
            INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
              VALUES ('pair', 'chapter-two', 1, 'pair', 'Two files', 'type_code', 0, 1, '# Pair\n\nPrint hello.', '2026-10-01T00:00:00Z', '2026-10-01T00:00:00Z');
            INSERT INTO code_lesson (lesson_uuid, lesson_type, prog_lang, expected_output) VALUES ('pair', 'type_code', 'go', 'Hello');
            INSERT INTO starter_file (lesson_uuid, name, position, content, is_hidden, is_readonly)
              VALUES ('pair', 'helper.go', 0, 'package main\n', 0, 0), ('pair', 'main.go', 1, 'package main\n\nfunc main() {}\n', 0, 0);
            """)
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        app.open(fixture.ref("pair"))
        try await fixture.wait { app.lesson?.ref == fixture.ref("pair") && app.lesson?.code?.isLoaded == true }
        let code = try #require(app.lesson?.code)
        let folder = try #require(code.folder)
        let main = try #require(code.files.first { $0.name == "main.go" })
        let helper = try #require(code.files.first { $0.name == "helper.go" })
        main.text = "package main // mine\n"
        try "package main // theirs\n".write(to: folder.appendingPathComponent("main.go"), atomically: true, encoding: .utf8)
        #expect(await code.flush() == false && code.conflict == "main.go")

        let lock = try WriteLock(fixture.directory.appendingPathComponent("study/learner.sqlite"))
        helper.text = "package main // helper, edited\n"
        code.run(submit: false) // its save writes helper.go, and waits for the database
        try await Task.sleep(for: .milliseconds(100))
        code.resolveConflict(keepMine: false)
        try lock.release()
        try await fixture.wait { code.result != nil && !code.isRunning }
        #expect(main.text == "package main // theirs\n" && code.conflict == nil)
        await app.prepareToQuit()
    }

    @Test func continueResumesBeforeAnyLessonHasBeenCompleted() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        try await fixture.wait { app.lesson?.code?.isLoaded == true }
        app.open(fixture.ref("optional"))
        try await fixture.wait { app.resumeLesson == fixture.ref("optional") }
        app.showDashboard()
        try await fixture.wait { app.lesson == nil }
        #expect(app.learningProgress.completedLessons == 0)
        let recommendation = try #require(app.recommendedLesson)
        #expect(recommendation == fixture.ref("optional"))
        app.continueLearning(recommendation)
        try await fixture.wait { app.lesson?.ref == recommendation && app.lesson?.code?.isLoaded == true }
        #expect(app.selection == .lesson(recommendation))
        await app.prepareToQuit()
    }

    @Test func anUnseenOpeningReturnsToTheRequestedLesson() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.continueLearning(fixture.ref("optional"))
        #expect(app.opening?.destination == fixture.ref("optional"))
        app.advanceOpening()
        app.advanceOpening()
        app.advanceOpening()
        try await fixture.wait { app.lesson?.ref == fixture.ref("optional") && app.lesson?.code?.isLoaded == true }
        #expect(app.opening == nil)
        await app.prepareToQuit()
    }

    @Test func theLastLessonReturnsToMissionControl() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        app.open(fixture.ref("quiz"))
        try await fixture.wait { app.lesson?.ref == fixture.ref("quiz") && app.resumeLesson == fixture.ref("quiz") }
        let lesson = try #require(app.lesson)
        lesson.markPassed(review: false)
        #expect(app.nextLesson == nil)
        app.continueAfterLesson()
        try await fixture.wait { app.lesson == nil }
        #expect(app.selection == .dashboard)
        #expect(app.isDone(fixture.ref("quiz")))
        await app.prepareToQuit()
    }

    @Test func aQuizOffersOnlyTheKeysItsAnswersUse() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        app.open(fixture.ref("quiz"))
        try await fixture.wait { app.lesson?.ref == fixture.ref("quiz") && app.lesson?.quiz != nil }
        // Its one question has two answers, so keys 3 to 5 pick nothing.
        #expect(app.lesson?.quiz?.answerKeys == 2)
        await app.prepareToQuit()
    }

    @Test func aChapterPassStillAllowsItsOptionalLesson() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        try await fixture.wait { app.lesson?.ref == fixture.ref("first") && app.lesson?.code?.isLoaded == true }
        let lesson = try #require(app.lesson)
        lesson.markPassed(review: false)
        #expect(lesson.finishedChapter)
        #expect(app.nextLesson?.id == "optional")
        app.continueAfterLesson()
        try await fixture.wait { app.lesson?.ref == fixture.ref("optional") && app.lesson?.code?.isLoaded == true }
        await app.prepareToQuit()
    }

    @Test func aPassRecordedBeforeItIsMarkedStillFinishesTheChapter() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        try await fixture.wait { app.lesson?.ref == fixture.ref("first") && app.lesson?.code?.isLoaded == true }
        let lesson = try #require(app.lesson)
        let finishing = lesson.finishesChapter
        // A code Submit records the pass before marking it, and the learner database reports
        // the lesson done as soon as the record is written.
        let learner = try #require(app.learner)
        _ = try await learner.recordSubmit(lesson.ref, passed: true, output: "Hello", durationMs: 1, reviewable: true)
        try await fixture.wait { app.isDone(lesson.ref) }
        lesson.markPassed(review: false, finishing: finishing)
        #expect(lesson.finishedChapter)
        await app.prepareToQuit()
    }

    @Test func aQuizPassThatFinishesTheChapterShowsItsDebriefNotTheQuiz() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        app.open(fixture.ref("quiz"))
        try await fixture.wait { app.lesson?.ref == fixture.ref("quiz") && app.lesson?.quiz != nil }
        let lesson = try #require(app.lesson)
        let quiz = try #require(lesson.quiz)
        // A wrong answer brings the quiz's verdicts into view.
        quiz.pick(1, for: 0)
        lesson.submit()
        #expect(lesson.footerRequests == 1)
        #expect(lesson.topRequests == 0)
        // The right one passes the chapter's only lesson: its debrief, above the readme, comes
        // into view instead of the quiz.
        quiz.pick(0, for: 0)
        lesson.submit()
        #expect(lesson.finishedChapter)
        #expect(lesson.topRequests == 1)
        #expect(lesson.footerRequests == 1)
        await app.prepareToQuit()
    }

    @Test func aRedoReturnsToReviewWithoutAdvancingTheCourse() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        app.redo(fixture.ref("first"))
        try await fixture.wait { app.lesson?.isReview == true && app.lesson?.code?.isLoaded == true }
        #expect(app.nextLesson == nil)
        app.continueAfterLesson()
        try await fixture.wait { app.showsReview && app.review?.isLoaded == true }
        #expect(!app.isDone(fixture.ref("first")))
        await app.prepareToQuit()
    }

    @Test func leavingALessonKeepsItsCodeAndNote() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        try await fixture.wait { app.lesson?.code?.isLoaded == true }
        let lesson = try #require(app.lesson)
        let file = try #require(lesson.code?.files.first)
        file.text = "package main\n\nfunc main() { println(42) }\n"
        lesson.note = "Remember the next assignment."
        app.showDashboard()
        try await fixture.wait { app.lesson == nil }
        let learner = try #require(app.learner)
        #expect(try await learner.note(for: lesson.ref) == "Remember the next assignment.")
        #expect(try await learner.codeFiles(for: lesson.ref)["main.go"]?.content == file.text)
        app.continueLearning(lesson.ref)
        try await fixture.wait { app.lesson?.ref == lesson.ref && app.lesson?.code?.isLoaded == true }
        #expect(app.lesson?.note == "Remember the next assignment.")
        #expect(app.lesson?.code?.files.first?.text == file.text)
        await app.prepareToQuit()
    }

    @Test func rapidNavigationKeepsOnlyTheLatestDestination() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        app.finishOpening()
        app.open(fixture.ref("optional"))
        app.showDashboard()
        app.open(fixture.ref("quiz"))
        try await fixture.wait { app.resumeLesson == fixture.ref("quiz") && !app.isLoadingLesson }
        #expect(app.lesson?.ref == fixture.ref("quiz"))
        #expect(app.selection == .lesson(fixture.ref("quiz")))
        #expect(app.banner == nil)
        await app.prepareToQuit()
    }

    @Test func notesWaitForTheirSavedTextBeforeAllowingEdits() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let learner = try #require(app.learner)
        let ref = fixture.ref("first")
        try await learner.saveNote(ref, body: "A note from an earlier session.")
        let database = try #require(app.library(ref.library)?.database)
        let detail = try #require(try await database.lesson(ref.lesson))
        let parsed = await ReadmeCache().parsed(for: detail)
        let lesson = LessonModel(detail: detail, parsed: parsed, saved: [:], learner: learner, review: false, app: app)
        #expect(!lesson.canEditNote)
        await lesson.loadNote()
        #expect(lesson.canEditNote)
        #expect(lesson.note == "A note from an earlier session.")
        lesson.note = "A draft still being edited."
        await lesson.loadNote()
        #expect(lesson.note == "A draft still being edited.")
        await lesson.leave()
        #expect(try await learner.note(for: ref) == "A draft still being edited.")
        await app.prepareToQuit()
    }
}
