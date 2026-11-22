import Foundation
import GRDB
import GosmosKit
@testable import Gosmos

/// A test library of two courses, with openings, code lessons and a quiz, and apps made on
/// it with study data of their own, in a temporary folder. With `reference`, the apps have the
/// Go documentation the app ships, reference/go.sqlite.
@MainActor struct FlowFixture {
    let directory: URL
    let resources: AppResources

    init(reference: Bool = false) throws {
        let repository = URL(fileURLWithPath: #filePath)
            .deletingLastPathComponent().deletingLastPathComponent().deletingLastPathComponent().deletingLastPathComponent()
        directory = FileManager.default.temporaryDirectory.appendingPathComponent("gosmos-flow-\(UUID().uuidString)")
        let libraries = directory.appendingPathComponent("courses")
        try FileManager.default.createDirectory(at: libraries, withIntermediateDirectories: true)
        resources = AppResources(libraries: libraries, learnerSchema: repository.appendingPathComponent("db/learner.sql"),
                                 reference: reference ? repository.appendingPathComponent("reference/go.sqlite") : nil)
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
        // this fixture's folder, never in ~/Library/Preferences. The study data's folder is one
        // of them, not a registered default: those are the whole process's, so suites running
        // side by side would open each other's.
        let defaults = UserDefaults(suiteName: directory.appendingPathComponent("settings-\(UUID().uuidString)").path)!
        defaults.set(directory.appendingPathComponent("study").path, forKey: LocalPreferences.dataFolderKey)
        return AppModel(defaults: defaults, resources: resources)
    }

    func ref(_ lesson: String) -> LessonRef { LessonRef(library: "flow-library", lesson: lesson) }

    /// Changes the library, as an update brings a new version of it.
    func changeLibrary(_ sql: String) throws {
        let queue = try DatabaseQueue(path: resources.libraries!.appendingPathComponent("practice.sqlite").path)
        try queue.write { try $0.execute(sql: sql) }
        try queue.close()
    }

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
