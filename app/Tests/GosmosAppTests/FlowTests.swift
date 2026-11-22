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
        // Format needs a Go with gofmt, which a test may not have; the others need nothing.
        let editing: [EditorCommand] = [.toggleComment, .indent, .outdent]
        var commands = app.commands
        #expect(commands.run && commands.submit && commands.nextLesson && commands.previousLesson && commands.back)
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

@MainActor private struct FlowFixture {
    let directory: URL
    let resources: AppResources

    init() throws {
        let repository = URL(fileURLWithPath: #filePath)
            .deletingLastPathComponent().deletingLastPathComponent().deletingLastPathComponent().deletingLastPathComponent()
        directory = FileManager.default.temporaryDirectory.appendingPathComponent("gosmos-flow-\(UUID().uuidString)")
        let libraries = directory.appendingPathComponent("courses")
        try FileManager.default.createDirectory(at: libraries, withIntermediateDirectories: true)
        resources = AppResources(libraries: libraries, learnerSchema: repository.appendingPathComponent("db/learner.sql"), reference: nil)
        let queue = try DatabaseQueue(path: libraries.appendingPathComponent("practice.sqlite").path)
        try queue.writeWithoutTransaction { db in
            try db.execute(sql: String(contentsOf: repository.appendingPathComponent("db/schema.sql"), encoding: .utf8))
        }
        try queue.write { db in
            try db.execute(sql: """
                INSERT INTO library_info VALUES (1, 'flow-library', 'Flow tests', '1');
                INSERT INTO course (uuid, slug, title, kind, short_description, description, position)
                  VALUES ('course', 'course', 'Practice', 'Course', '', '', 0),
                         ('second-course', 'second-course', 'More practice', 'Course', '', '', 1);
                INSERT INTO chapter (uuid, course_uuid, position, slug, title, description)
                  VALUES ('chapter-one', 'course', 0, 'one', 'First chapter', ''),
                         ('chapter-two', 'course', 1, 'two', 'Second chapter', ''),
                         ('chapter-three', 'second-course', 0, 'three', 'The second course''s chapter', '');
                INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
                  VALUES ('first', 'chapter-one', 0, 'first', 'First assignment', 'type_code', 1, 1, '# First\n\nPrint hello.', '2026-10-01T00:00:00Z', '2026-10-01T00:00:00Z'),
                         ('optional', 'chapter-one', 1, 'optional', 'Optional practice', 'type_code', 0, 1, '# Optional\n\nTry another greeting.', '2026-10-01T00:00:00Z', '2026-10-01T00:00:00Z'),
                         ('quiz', 'chapter-two', 0, 'quiz', 'Last assignment', 'type_choice', 1, 1, '# Quiz\n\nCheck your understanding.', '2026-10-01T00:00:00Z', '2026-10-01T00:00:00Z'),
                         ('second', 'chapter-three', 0, 'second', 'The second course''s assignment', 'type_code', 1, 1, '# Second\n\nPrint hello again.', '2026-10-01T00:00:00Z', '2026-10-01T00:00:00Z');
                INSERT INTO code_lesson (lesson_uuid, lesson_type, prog_lang, expected_output)
                  VALUES ('first', 'type_code', 'go', 'Hello'), ('optional', 'type_code', 'go', 'Hello'), ('second', 'type_code', 'go', 'Hello');
                INSERT INTO starter_file (lesson_uuid, name, position, content, is_hidden, is_readonly)
                  VALUES ('first', 'main.go', 0, 'package main\n\nfunc main() {}\n', 0, 0),
                         ('optional', 'main.go', 0, 'package main\n\nfunc main() {}\n', 0, 0),
                         ('second', 'main.go', 0, 'package main\n\nfunc main() {}\n', 0, 0);
                INSERT INTO quiz_lesson (lesson_uuid, contains_complete_dir) VALUES ('quiz', 0);
                INSERT INTO quiz_question VALUES ('quiz', 0, 'Is Go compiled?');
                INSERT INTO quiz_answer VALUES ('quiz', 0, 0, 'Yes', 1, 'Go builds programs.'), ('quiz', 0, 1, 'No', 0, 'Go compiles its programs.');
                INSERT INTO story_page VALUES ('course', 0, 'text', NULL, 'Welcome', 'Welcome to the course.', 'Continue'),
                                              ('course', 1, 'badge', 'Practice', 'Your badge', 'Choose your name.', 'Continue'),
                                              ('course', 2, 'text', NULL, 'Your assignment', 'Start when you are ready.', 'Start'),
                                              ('second-course', 0, 'text', NULL, 'Further out', 'A second course begins.', 'Start');
                """)
        }
        try queue.close()
    }

    func app() -> AppModel {
        // A suite named by an absolute path is stored at that path: each app's settings stay in
        // this fixture's folder, never in ~/Library/Preferences.
        let defaults = UserDefaults(suiteName: directory.appendingPathComponent("settings-\(UUID().uuidString)").path)!
        defaults.register(defaults: [LocalPreferences.dataFolderKey: directory.appendingPathComponent("study").path])
        return AppModel(defaults: defaults, resources: resources)
    }

    func ref(_ lesson: String) -> LessonRef { LessonRef(library: "flow-library", lesson: lesson) }

    func wait(_ condition: () -> Bool) async throws {
        for _ in 0..<500 {
            if condition() { return }
            try await Task.sleep(for: .milliseconds(10))
        }
        struct TimedOut: Error {}
        throw TimedOut()
    }

    func remove() { try? FileManager.default.removeItem(at: directory) }
}
