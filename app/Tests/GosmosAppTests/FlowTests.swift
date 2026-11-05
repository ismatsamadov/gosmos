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
    /// The settings suite of every app made here. Each becomes a file in ~/Library/Preferences
    /// once written, and a test run must leave nothing there.
    private let suites = Suites()

    @MainActor private final class Suites { var names: [String] = [] }

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
                  VALUES ('course', 'course', 'Practice', 'Course', '', '', 0);
                INSERT INTO chapter (uuid, course_uuid, position, slug, title, description)
                  VALUES ('chapter-one', 'course', 0, 'one', 'First chapter', ''),
                         ('chapter-two', 'course', 1, 'two', 'Second chapter', '');
                INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
                  VALUES ('first', 'chapter-one', 0, 'first', 'First assignment', 'type_code', 1, 1, '# First\n\nPrint hello.', '2026-10-01T00:00:00Z', '2026-10-01T00:00:00Z'),
                         ('optional', 'chapter-one', 1, 'optional', 'Optional practice', 'type_code', 0, 1, '# Optional\n\nTry another greeting.', '2026-10-01T00:00:00Z', '2026-10-01T00:00:00Z'),
                         ('quiz', 'chapter-two', 0, 'quiz', 'Last assignment', 'type_choice', 1, 1, '# Quiz\n\nCheck your understanding.', '2026-10-01T00:00:00Z', '2026-10-01T00:00:00Z');
                INSERT INTO code_lesson (lesson_uuid, lesson_type, prog_lang, expected_output)
                  VALUES ('first', 'type_code', 'go', 'Hello'), ('optional', 'type_code', 'go', 'Hello');
                INSERT INTO starter_file (lesson_uuid, name, position, content, is_hidden, is_readonly)
                  VALUES ('first', 'main.go', 0, 'package main\n\nfunc main() {}\n', 0, 0),
                         ('optional', 'main.go', 0, 'package main\n\nfunc main() {}\n', 0, 0);
                INSERT INTO quiz_lesson (lesson_uuid, contains_complete_dir) VALUES ('quiz', 0);
                INSERT INTO quiz_question VALUES ('quiz', 0, 'Is Go compiled?');
                INSERT INTO quiz_answer VALUES ('quiz', 0, 0, 'Yes', 1, 'Go builds programs.'), ('quiz', 0, 1, 'No', 0, 'Go compiles its programs.');
                INSERT INTO story_page VALUES ('course', 0, 'text', NULL, 'Welcome', 'Welcome to the course.', 'Continue'),
                                              ('course', 1, 'badge', 'Practice', 'Your badge', 'Choose your name.', 'Continue'),
                                              ('course', 2, 'text', NULL, 'Your assignment', 'Start when you are ready.', 'Start');
                """)
        }
        try queue.close()
    }

    func app() -> AppModel {
        let name = "app.gosmos.flow.\(UUID().uuidString)"
        suites.names.append(name)
        let defaults = UserDefaults(suiteName: name)!
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

    func remove() {
        try? FileManager.default.removeItem(at: directory)
        // Removing a domain leaves an empty plist behind, so flush it, then delete the file.
        let preferences = FileManager.default.homeDirectoryForCurrentUser.appendingPathComponent("Library/Preferences")
        for name in suites.names {
            UserDefaults.standard.removePersistentDomain(forName: name)
            CFPreferencesAppSynchronize(name as CFString)
            try? FileManager.default.removeItem(at: preferences.appendingPathComponent("\(name).plist"))
        }
    }
}
