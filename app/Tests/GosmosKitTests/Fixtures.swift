import Foundation
import GRDB
@testable import GosmosKit

/// Test libraries made from the repository's db/schema.sql and learner databases made from
/// db/learner.sql, in $TMPDIR, filled with lessons written for these tests (they only mimic
/// the shape of real lessons). Nothing here touches a real library or the real learner
/// database.
enum Fixtures {
    static let repository = URL(fileURLWithPath: #filePath)
        .deletingLastPathComponent().deletingLastPathComponent().deletingLastPathComponent().deletingLastPathComponent()
    static let schemaURL = repository.appendingPathComponent("db/schema.sql")
    static let learnerSchemaURL = repository.appendingPathComponent("db/learner.sql")

    /// The fixture library's identity (library_info).
    static let libraryID = "b0000000-0000-0000-0000-000000000001"
    static let libraryTitle = "Practice"

    /// A code lesson outside any database, for runs and folders.
    static func codeLesson(_ id: String, _ type: LessonType, _ code: CodeExercise) -> LessonDetail {
        LessonDetail(library: libraryID, summary: LessonSummary(id: id, chapterID: "h", title: "Test", type: type, required: true),
                     courseID: "c", courseTitle: "Test Course", chapterTitle: "Chapter",
                     readme: "# Test lesson", difficulty: 1, topics: [], content: type == .code ? .code(code) : .codeTests(code))
    }

    /// A fresh directory under $TMPDIR; the caller removes it.
    static func tempDir(_ name: String = "test") throws -> URL {
        let url = URL(fileURLWithPath: NSTemporaryDirectory())
            .appendingPathComponent("gosmos-tests-\(name)-\(UUID().uuidString)", isDirectory: true)
        try FileManager.default.createDirectory(at: url, withIntermediateDirectories: true)
        return URL(fileURLWithPath: GoRunner.realPath(url.path), isDirectory: true)
    }

    static let ts = "2026-10-01T09:00:00Z"
    /// Test libraries' file name: any name works.
    static let libraryName = "practice.library"
    static let courseID = "c0000000-0000-0000-0000-000000000001"
    static let projectCourseID = "c0000000-0000-0000-0000-000000000002"
    static let chapter1 = "h0000000-0000-0000-0000-000000000001"
    static let chapter2 = "h0000000-0000-0000-0000-000000000002"
    static let chapter3 = "h0000000-0000-0000-0000-000000000003"
    static let codeLesson = "l0000000-0000-0000-0000-000000000001"
    static let testsLesson = "l0000000-0000-0000-0000-000000000002"
    static let quizLesson = "l0000000-0000-0000-0000-000000000003"
    static let planetsLesson = "l0000000-0000-0000-0000-000000000004"
    static let loopsLesson = "l0000000-0000-0000-0000-000000000005"
    static let optionalLesson = "l0000000-0000-0000-0000-000000000007"

    static let addTests = """
    package main

    import (
    \t"fmt"
    \t"testing"
    )

    func TestAdd(t *testing.T) {
    \tsums := [][3]int{{1, 2, 3}, {2, 2, 4}}
    \tif withSubmit {
    \t\tsums = append(sums, [3]int{1000, 1, 1001})
    \t}
    \tfor _, s := range sums {
    \t\tgot := add(s[0], s[1])
    \t\tif got != s[2] {
    \t\t\tt.Errorf("add(%d, %d) = %d, want %d", s[0], s[1], got, s[2])
    \t\t\tcontinue
    \t\t}
    \t\tfmt.Printf("add(%d, %d) = %d ok\\n", s[0], s[1], got)
    \t}
    }

    var withSubmit = true

    """

    static let wasmOnly = """
    //go:build js && wasm

    package main

    func init() { withSubmit = false }

    """

    /// Creates a library from db/schema.sql with lessons of every type, in two courses; `id` and
    /// `title` go into its library_info.
    static func makeDatabase(at url: URL, id: String = libraryID, title: String = libraryTitle) throws {
        let queue = try DatabaseQueue(path: url.path)
        try queue.writeWithoutTransaction { db in
            try db.execute(sql: String(contentsOf: schemaURL, encoding: .utf8))
            try db.execute(sql: "PRAGMA journal_mode = DELETE")
        }
        try queue.write { db in
            try db.execute(sql: "INSERT INTO library_info(singleton, uuid, title, version) VALUES (1, ?, ?, '1')", arguments: [id, title])
            try db.execute(sql: "INSERT INTO course(uuid, slug, title, kind, short_description, description, position) VALUES (?,?,?,?,?,?,?)",
                           arguments: [courseID, "learn-testing", "Learn Testing", "Course", "Short", "Long", 0])
            try db.execute(sql: "INSERT INTO course(uuid, slug, title, kind, short_description, description, position) VALUES (?,?,?,?,?,?,?)",
                           arguments: [projectCourseID, "build-a-thing", "Build a Thing", "Guided Project", "Short", "Long", 1])
            for (id, course, pos, slug, title) in [(chapter1, courseID, 0, "1-basics", "Basics"), (chapter2, courseID, 1, "2-more", "More"),
                                                   (chapter3, projectCourseID, 0, "1-start", "Start")] {
                try db.execute(sql: "INSERT INTO chapter(uuid, course_uuid, position, slug, title, description) VALUES (?,?,?,?,?,?)",
                               arguments: [id, course, pos, slug, title, ""])
            }
            func lesson(_ id: String, _ chapter: String, _ pos: Int, _ slug: String, _ title: String, _ type: String, _ readme: String, required: Int = 1) throws {
                try db.execute(sql: """
                    INSERT INTO lesson(uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
                    VALUES (?,?,?,?,?,?,?,?,?,?,?)
                    """, arguments: [id, chapter, pos, slug, title, type, required, 3, readme, ts, ts])
            }
            try lesson(codeLesson, chapter1, 0, "1-hello", "Hello Gopher", "type_code", "# Hello Gopher\n\nPrint a greeting with `fmt.Println`.")
            try db.execute(sql: "INSERT INTO code_lesson(lesson_uuid, lesson_type, prog_lang, expected_output, solution_note) VALUES (?,?,?,?,?)",
                           arguments: [codeLesson, "type_code", "go", "Hello, gopher\n2 + 2 = 4\n", "Println ends each line; %d formats the sum."])
            try db.execute(sql: "INSERT INTO starter_file(lesson_uuid, name, position, content, is_hidden, is_readonly) VALUES (?,?,?,?,?,?)",
                           arguments: [codeLesson, "main.go", 0, "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, world\")\n}\n", 0, 0])
            try db.execute(sql: "INSERT INTO solution_file(lesson_uuid, name, content) VALUES (?,?,?)",
                           arguments: [codeLesson, "main.go", helloSolution])
            for (i, hint) in ["Look at what the program prints now.", "Println takes the text to print.", "Print \"Hello, gopher\", then the sum with Printf."].enumerated() {
                try db.execute(sql: "INSERT INTO lesson_hint(lesson_uuid, position, text) VALUES (?,?,?)", arguments: [codeLesson, i, hint])
            }
            try db.execute(sql: "INSERT INTO lesson_topic(lesson_uuid, topic) VALUES (?,?)", arguments: [codeLesson, "printing"])
            try db.execute(sql: "INSERT INTO lesson_topic(lesson_uuid, topic) VALUES (?,?)", arguments: [codeLesson, "fmt"])

            try lesson(testsLesson, chapter1, 1, "2-add", "Adding Numbers", "type_code_tests", "# Adding Numbers\n\nComplete `add`.")
            try db.execute(sql: "INSERT INTO code_lesson(lesson_uuid, lesson_type, prog_lang) VALUES (?,?,?)", arguments: [testsLesson, "type_code_tests", "go"])
            for (i, f) in [("main.go", "package main\n\nfunc add(a, b int) int {\n\treturn 0\n}\n", 0, 0),
                           ("main_test.go", addTests, 0, 1), ("main_js_test.go", wasmOnly, 1, 1)].enumerated() {
                try db.execute(sql: "INSERT INTO starter_file(lesson_uuid, name, position, content, is_hidden, is_readonly) VALUES (?,?,?,?,?,?)",
                               arguments: [testsLesson, f.0, i, f.1, f.2, f.3])
            }

            try lesson(quizLesson, chapter1, 2, "3-quiz", "Functions Quiz", "type_choice", "# Functions Quiz\n\nWhich keyword declares a function?")
            try db.execute(sql: "INSERT INTO quiz_lesson(lesson_uuid, contains_complete_dir) VALUES (?,?)", arguments: [quizLesson, 0])
            try db.execute(sql: "INSERT INTO quiz_question(lesson_uuid, position, question) VALUES (?,?,?)", arguments: [quizLesson, 0, "Which keyword declares a function?"])
            for (i, a) in ["`fn`", "`func`", "`def`"].enumerated() {
                let why: String? = a == "`func`" ? "Go declares every function with func." : a == "`def`" ? "def is Python's keyword." : nil
                try db.execute(sql: "INSERT INTO quiz_answer(lesson_uuid, question_position, position, answer, is_correct, explanation) VALUES (?,?,?,?,?,?)",
                               arguments: [quizLesson, 0, i, a, a == "`func`" ? 1 : 0, why])
            }
            try db.execute(sql: "INSERT INTO quiz_question(lesson_uuid, position, question) VALUES (?,?,?)", arguments: [quizLesson, 1, "Is Go compiled?"])
            for (i, a) in ["Yes", "No"].enumerated() {
                try db.execute(sql: "INSERT INTO quiz_answer(lesson_uuid, question_position, position, answer, is_correct) VALUES (?,?,?,?,?)",
                               arguments: [quizLesson, 1, i, a, a == "Yes" ? 1 : 0])
            }

            try lesson(planetsLesson, chapter2, 0, "1-planets", "Planets", "type_code", "# Planets\n\nPrint the inner planets, one a line.")
            try db.execute(sql: "INSERT INTO code_lesson(lesson_uuid, lesson_type, prog_lang, expected_output) VALUES (?,?,?,?)",
                           arguments: [planetsLesson, "type_code", "go", "Mercury\nVenus\nEarth\nMars\n"])
            try db.execute(sql: "INSERT INTO starter_file(lesson_uuid, name, position, content, is_hidden, is_readonly) VALUES (?,?,?,?,?,?)",
                           arguments: [planetsLesson, "main.go", 0, "package main\n\nfunc main() {\n}\n", 0, 0])

            try lesson(loopsLesson, chapter2, 1, "2-loops", "Loops", "type_choice", "# Loops\n\nHow does Go loop?")
            try db.execute(sql: "INSERT INTO quiz_lesson(lesson_uuid, contains_complete_dir) VALUES (?,?)", arguments: [loopsLesson, 0])
            try db.execute(sql: "INSERT INTO quiz_question(lesson_uuid, position, question) VALUES (?,?,?)", arguments: [loopsLesson, 0, "Which keyword loops?"])
            for (i, a) in ["`for`", "`while`"].enumerated() {
                try db.execute(sql: "INSERT INTO quiz_answer(lesson_uuid, question_position, position, answer, is_correct) VALUES (?,?,?,?,?)",
                               arguments: [loopsLesson, 0, i, a, i == 0 ? 1 : 0])
            }

            try lesson(optionalLesson, chapter3, 0, "1-extra", "Extra Credit", "type_choice", "# Extra Credit", required: 0)
            try db.execute(sql: "INSERT INTO quiz_lesson(lesson_uuid, contains_complete_dir) VALUES (?,?)", arguments: [optionalLesson, 0])
        }
        try queue.close()
    }

    /// Rows of the learner database in `dir`, read on a connection of their own (the app's
    /// API has no reader for some of these tables).
    static func rows(in dir: URL, _ sql: String, _ arguments: StatementArguments = []) throws -> [Row] {
        var config = Configuration()
        config.readonly = true
        let queue = try DatabaseQueue(path: dir.appendingPathComponent("learner.sqlite").path, configuration: config)
        return try queue.read { try Row.fetchAll($0, sql: sql, arguments: arguments) }
    }

    /// The fixture library in a directory of its own, opened as the app opens it; `edit`
    /// changes the file first, on a connection of its own (the app's is read-only).
    static func database(_ edit: ((Database) throws -> Void)? = nil) throws -> (AppDatabase, URL) {
        let dir = try tempDir("db")
        let url = dir.appendingPathComponent(Fixtures.libraryName)
        try makeDatabase(at: url)
        if let edit {
            let queue = try DatabaseQueue(path: url.path)
            try queue.write(edit)
            try queue.close()
        }
        return (try AppDatabase(url: url), dir)
    }

    /// A new learner database in `dir`, made from db/learner.sql, with the fixture library in
    /// its library table.
    static func learner(in dir: URL) async throws -> LearnerDatabase {
        let learner = try LearnerDatabase(url: dir.appendingPathComponent("learner.sqlite"), schema: learnerSchemaURL)
        try await learner.registerLibrary(LibraryInfo(uuid: libraryID, title: libraryTitle), location: dir.appendingPathComponent(libraryName).path)
        return learner
    }

    /// A lesson of the fixture library, for the learner database.
    static func ref(_ lesson: String) -> LessonRef { LessonRef(library: libraryID, lesson: lesson) }

    static let helloSolution = "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, gopher\")\n\tfmt.Printf(\"2 + 2 = %d\\n\", 2+2)\n}\n"

    /// The Go on this machine's PATH, as if chosen in Settings (`swift test` has no app with a
    /// Go built in), or nil: then the tests that run it are skipped.
    static func toolchain() async -> Toolchain? {
        guard let go = ProcessRunner.resolve("go", environment: ProcessInfo.processInfo.environment) else { return nil }
        let t = await Toolchain.discover(goOverride: go, builtIn: nil, caches: goCaches)
        return t.go == nil ? nil : t
    }

    /// Go's caches for the tests, kept between runs so builds after the first are quick; never
    /// the app's.
    static let goCaches = FileManager.default.temporaryDirectory.appendingPathComponent("gosmos-tests-go", isDirectory: true)
}

/// Holds work where it calls `pass`, until opened: a load where it asks for the Go version, so
/// a test can type while the load is under way, or a write, so a test can quit meanwhile.
actor Gate {
    private var isOpen = false
    private var held: [CheckedContinuation<Void, Never>] = []
    private var watchers: [CheckedContinuation<Void, Never>] = []

    /// Waits here until the gate opens.
    func pass() async {
        guard !isOpen else { return }
        for watcher in watchers { watcher.resume() }
        watchers = []
        await withCheckedContinuation { held.append($0) }
    }

    /// Returns once something waits at the gate.
    func somethingWaits() async {
        guard held.isEmpty else { return }
        await withCheckedContinuation { watchers.append($0) }
    }

    func open() {
        isOpen = true
        for waiter in held { waiter.resume() }
        held = []
    }
}
