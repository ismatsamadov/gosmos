import Foundation
import GRDB
import Testing
@testable import GosmosKit

/// A library unlike the test fixture, built here from db/schema.sql with made-up lessons: two
/// courses (one a guided project), three chapters, lessons of every type.
enum KitchenLibrary {
    static let boil = "kitchen-boil-water", salt = "kitchen-measure-salt", acid = "kitchen-which-acid"
    static let timer = "kitchen-set-timer", plan = "garden-plan-beds", rows = "garden-count-rows"
    static let kitchen = "course-kitchen", garden = "course-garden"
    static let library = "library-kitchen"

    static let saltTests = """
    package main

    import "testing"

    func TestPinches(t *testing.T) {
    \tfor grams, want := range map[int]int{2: 1, 10: 5} {
    \t\tif got := pinches(grams); got != want {
    \t\t\tt.Errorf("pinches(%d) = %d, want %d", grams, got, want)
    \t\t}
    \t}
    }
    """

    static func build(at url: URL) throws {
        let queue = try DatabaseQueue(path: url.path)
        try queue.writeWithoutTransaction { db in
            try db.execute(sql: String(contentsOf: Fixtures.schemaURL, encoding: .utf8))
            try db.execute(sql: "PRAGMA journal_mode = DELETE")
        }
        try queue.write { db in
            let t = Fixtures.ts
            try db.execute(sql: "INSERT INTO library_info(singleton, uuid, title, version) VALUES (1, ?, ?, ?)", arguments: [library, "Kitchen", "2"])
            for (uuid, title, kind, position) in [(kitchen, "Kitchen Chemistry", "Course", 0), (garden, "Garden Planner", "Guided Project", 1)] {
                try db.execute(sql: "INSERT INTO course(uuid, slug, title, kind, short_description, description, position) VALUES (?,?,?,?,?,?,?)",
                               arguments: [uuid, uuid, title, kind, "", "", position])
            }
            for (uuid, course, position, title) in [("kitchen-basics", kitchen, 0, "Basics"), ("kitchen-practice", kitchen, 1, "Practice"),
                                                    ("garden-start", garden, 0, "Start")] {
                try db.execute(sql: "INSERT INTO chapter(uuid, course_uuid, position, slug, title, description) VALUES (?,?,?,?,?,?)",
                               arguments: [uuid, course, position, uuid, title, ""])
            }
            func lesson(_ uuid: String, _ chapter: String, _ position: Int, _ title: String, _ type: String, _ readme: String) throws {
                try db.execute(sql: """
                    INSERT INTO lesson(uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
                    VALUES (?,?,?,?,?,?,1,4,?,?,?)
                    """, arguments: [uuid, chapter, position, uuid, title, type, readme, t, t])
            }

            try lesson(boil, "kitchen-basics", 0, "Boil Water", "type_code", "# Boil Water\n\nPrint `bubbles`.")
            try db.execute(sql: "INSERT INTO code_lesson(lesson_uuid, lesson_type, prog_lang, expected_output) VALUES (?,?,?,?)",
                           arguments: [boil, "type_code", "go", "bubbles\n"])
            try db.execute(sql: "INSERT INTO starter_file(lesson_uuid, name, position, content, is_hidden, is_readonly) VALUES (?,?,?,?,0,0)",
                           arguments: [boil, "main.go", 0, "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"bubbles\")\n}\n"])

            try lesson(salt, "kitchen-basics", 1, "Measure Salt", "type_code_tests", "# Measure Salt\n\nA pinch is two grams.")
            try db.execute(sql: "INSERT INTO code_lesson(lesson_uuid, lesson_type, prog_lang) VALUES (?,?,?)", arguments: [salt, "type_code_tests", "go"])
            try db.execute(sql: "INSERT INTO starter_file(lesson_uuid, name, position, content, is_hidden, is_readonly) VALUES (?,?,?,?,0,0)",
                           arguments: [salt, "main.go", 0, "package main\n\nfunc pinches(grams int) int {\n\treturn grams / 2\n}\n"])
            try db.execute(sql: "INSERT INTO starter_file(lesson_uuid, name, position, content, is_hidden, is_readonly) VALUES (?,?,?,?,0,1)",
                           arguments: [salt, "main_test.go", 1, saltTests])

            try lesson(acid, "kitchen-basics", 2, "Which Is an Acid?", "type_choice", "# Which Is an Acid?")
            try db.execute(sql: "INSERT INTO quiz_lesson(lesson_uuid, contains_complete_dir) VALUES (?, 0)", arguments: [acid])
            for (q, question, answers) in [(0, "Which is an acid?", [("Vinegar", 1), ("Baking soda", 0)]), (1, "Does salt dissolve in water?", [("No", 0), ("Yes", 1)])] {
                try db.execute(sql: "INSERT INTO quiz_question(lesson_uuid, position, question) VALUES (?,?,?)", arguments: [acid, q, question])
                for (a, (answer, right)) in answers.enumerated() {
                    try db.execute(sql: "INSERT INTO quiz_answer(lesson_uuid, question_position, position, answer, is_correct) VALUES (?,?,?,?,?)",
                                   arguments: [acid, q, a, answer, right])
                }
            }

            try lesson(timer, "kitchen-practice", 0, "Set a Timer", "type_code", "# Set a Timer\n\nPrint `ding` when the time is up.")
            try db.execute(sql: "INSERT INTO code_lesson(lesson_uuid, lesson_type, prog_lang, expected_output) VALUES (?,?,?,?)",
                           arguments: [timer, "type_code", "go", "ding\n"])
            try db.execute(sql: "INSERT INTO starter_file(lesson_uuid, name, position, content, is_hidden, is_readonly) VALUES (?,?,?,?,0,0)",
                           arguments: [timer, "main.go", 0, "package main\n\nfunc main() {\n}\n"])

            try lesson(plan, "garden-start", 0, "Plan the Beds", "type_choice", "# Plan the Beds")
            try db.execute(sql: "INSERT INTO quiz_lesson(lesson_uuid, contains_complete_dir) VALUES (?, 0)", arguments: [plan])
            try db.execute(sql: "INSERT INTO quiz_question(lesson_uuid, position, question) VALUES (?,0,?)", arguments: [plan, "Where do tomatoes go?"])
            for (a, (answer, right)) in [("In the sun", 1), ("In the shade", 0)].enumerated() {
                try db.execute(sql: "INSERT INTO quiz_answer(lesson_uuid, question_position, position, answer, is_correct) VALUES (?,0,?,?,?)",
                               arguments: [plan, a, answer, right])
            }

            try lesson(rows, "garden-start", 1, "Count the Rows", "type_code_tests", "# Count the Rows")
            try db.execute(sql: "INSERT INTO code_lesson(lesson_uuid, lesson_type, prog_lang) VALUES (?,?,?)", arguments: [rows, "type_code_tests", "go"])
            try db.execute(sql: "INSERT INTO starter_file(lesson_uuid, name, position, content, is_hidden, is_readonly) VALUES (?,?,?,?,0,0)",
                           arguments: [rows, "rows.go", 0, "package main\n\nfunc rows(beds int) int {\n\treturn 0\n}\n"])
        }
        try queue.close()
    }
}

/// The whole path with a library the app has never seen: open and check it, read every lesson
/// type, grade the quiz, and run the code lessons when Go is here.
@Suite(.serialized) struct LibraryTests {
    @Test func anyLibraryInTheFormatWorks() async throws {
        let dir = try Fixtures.tempDir("kitchen")
        defer { try? FileManager.default.removeItem(at: dir) }
        let url = dir.appendingPathComponent("kitchen.gosmos")
        try KitchenLibrary.build(at: url)

        // Opened through the same open-and-check path as the app's; your study data goes into a
        // learner database of its own.
        let db = try AppDatabase(url: url)
        #expect(db.info == LibraryInfo(uuid: KitchenLibrary.library, title: "Kitchen"))
        let learner = try LearnerDatabase(url: dir.appendingPathComponent("learner.sqlite"), schema: Fixtures.learnerSchemaURL)
        try await learner.registerLibrary(db.info, location: url.path)
        let catalog = try await db.catalog()
        #expect(catalog.courses.map(\.title) == ["Kitchen Chemistry", "Garden Planner"])
        #expect(catalog.courses.map(\.isGuidedProject) == [false, true])
        #expect(catalog.courses.flatMap(\.chapters).map(\.title) == ["Basics", "Practice", "Start"])
        #expect(catalog.lessonCount == 6 && catalog.firstLesson?.id == KitchenLibrary.boil)
        #expect(catalog.next(after: KitchenLibrary.acid)?.id == KitchenLibrary.timer)

        // Every type reads back as itself.
        var types: [LessonType] = []
        for summary in catalog.courses.flatMap(\.lessons) {
            let lesson = try #require(try await db.lesson(summary.id))
            types.append(lesson.type)
        }
        #expect(types == [.code, .codeTests, .choice, .code, .choice, .codeTests])
        guard case .choice(let beds)? = try await db.lesson(KitchenLibrary.plan)?.content,
              case .codeTests(let count)? = try await db.lesson(KitchenLibrary.rows)?.content else { Issue.record("garden lessons"); return }
        #expect(beds.questions.map(\.answer) == ["In the sun"] && count.visibleFiles.map(\.name) == ["rows.go"])

        // The quiz: one wrong, then right; done.
        guard case .choice(let quiz)? = try await db.lesson(KitchenLibrary.acid)?.content else { Issue.record("quiz"); return }
        var state = QuizState(quiz: quiz)
        _ = state.pick(1, for: 0)
        _ = state.pick(1, for: 1)
        let first = state.submit()
        #expect(first.map(\.right) == [false, true] && !state.allCorrect)
        _ = state.pick(0, for: 0)
        let second = state.submit()
        #expect(second.map(\.right) == [true] && state.allCorrect)
        let acid = LessonRef(library: KitchenLibrary.library, lesson: KitchenLibrary.acid)
        try await learner.recordAnswers(acid, first + second)
        try await learner.setDone(acid, true)
        #expect(try await learner.done() == [KitchenLibrary.library: [KitchenLibrary.acid]])

        // The code lessons, with Go when it is here.
        guard let toolchain = await Fixtures.toolchain() else { return }
        let workspace = Workspace(root: dir.appendingPathComponent("work"))
        for id in [KitchenLibrary.boil, KitchenLibrary.salt] {
            let lesson = try #require(try await db.lesson(id))
            let code = try #require(lesson.code)
            let folder = try workspace.prepare(lesson, code: code, saved: [:], goVersion: toolchain.languageVersion).folder
            let result = await GoRunner(toolchain: toolchain, scratch: workspace.runsFolder).run(code, tests: lesson.type == .codeTests, folder: folder, submit: true)
            #expect(result.status == .passed, "\(lesson.title): \(result.output)")
        }
    }

    @Test func aLessonInAnotherLanguageOpensButDoesNotRun() async throws {
        let dir = try Fixtures.tempDir("python")
        defer { try? FileManager.default.removeItem(at: dir) }
        let url = dir.appendingPathComponent("kitchen.gosmos")
        try KitchenLibrary.build(at: url)
        let queue = try DatabaseQueue(path: url.path)
        try await queue.write { try $0.execute(sql: "UPDATE code_lesson SET prog_lang = 'python' WHERE lesson_uuid = ?", arguments: [KitchenLibrary.boil]) }
        try queue.close()
        let db = try AppDatabase(url: url)
        let code = try #require(try await db.lesson(KitchenLibrary.boil)?.code)
        #expect(!code.runsHere && code.otherLanguageNote == "This lesson uses python; Gosmos runs Go lessons.")
        guard let toolchain = await Fixtures.toolchain() else { return }
        let result = await GoRunner(toolchain: toolchain, scratch: dir.appendingPathComponent(".runs")).run(code, tests: false, folder: dir, submit: true)
        #expect(result.status == .error && result.output == code.otherLanguageNote)
    }
}
