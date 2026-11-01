import Foundation
import GRDB
import Testing
@testable import GosmosKit

/// A course library, as the app reads it: read-only, from the tables of format 7.
@Suite struct DatabaseTests {
    @Test func opensAnyFileNameAndKnowsItsLibrary() async throws {
        let (db, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        #expect(db.info == LibraryInfo(uuid: Fixtures.libraryID, title: Fixtures.libraryTitle))
        #expect(try await db.catalog().lessonCount == 8)
        let missing = dir.appendingPathComponent("missing.library")
        #expect(throws: AppDatabaseError.missing(path: missing.path)) { try AppDatabase(url: missing) }
        #expect(!FileManager.default.fileExists(atPath: missing.path)) // never created
    }

    @Test func catalogAndNavigation() async throws {
        let (db, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        let catalog = try await db.catalog()
        #expect(catalog.courses.map(\.title) == ["Learn Testing", "Build a Thing"])
        #expect(catalog.courses[1].isGuidedProject && !catalog.courses[0].isGuidedProject)
        #expect(catalog.courses[0].chapters.map(\.title) == ["Basics", "More"])
        #expect(catalog.courses[0].chapters[0].lessons.map(\.title) == ["Hello Gopher", "Adding Numbers", "Functions Quiz"])
        #expect(catalog.lessonCount == 8)
        #expect(catalog.courses[1].requiredLessons.isEmpty)
        #expect(catalog.next(after: Fixtures.quizLesson)?.id == Fixtures.cliLesson) // across chapters
        #expect(catalog.previous(before: Fixtures.cliLesson)?.id == Fixtures.quizLesson)
        #expect(catalog.previous(before: Fixtures.codeLesson) == nil)
        #expect(catalog.next(after: Fixtures.interviewLesson)?.id == Fixtures.askingLesson)
        #expect(catalog.next(after: Fixtures.askingLesson) == nil) // not into another course
        #expect(catalog.firstLesson?.id == Fixtures.codeLesson)
    }

    @Test func lessonContentForEveryType() async throws {
        let (db, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }

        let code = try #require(try await db.lesson(Fixtures.codeLesson))
        #expect(code.ref == Fixtures.ref(Fixtures.codeLesson))
        #expect(code.courseTitle == "Learn Testing" && code.chapterTitle == "Basics" && code.difficulty == 3)
        #expect(code.topics == ["fmt", "printing"])
        #expect(code.hints.map(\.position) == [0, 1, 2] && code.hints[1].text == "Println takes the text to print.")
        guard case .code(let exercise) = code.content else { Issue.record("not a code lesson"); return }
        #expect(exercise.expectedOutput == "Hello, gopher\n2 + 2 = 4\n")
        #expect(exercise.starterFiles.map(\.name) == ["main.go"])
        // Format 7: the reference solution and why it works.
        #expect(exercise.solutionFiles == [SolutionFile(name: "main.go", content: Fixtures.helloSolution)])
        #expect(exercise.solutionNote == "Println ends each line; %d formats the sum.")

        guard case .codeTests(let tests)? = try await db.lesson(Fixtures.testsLesson)?.content else { Issue.record("tests"); return }
        #expect(tests.expectedOutput == nil && tests.solutionFiles.isEmpty && tests.solutionNote == nil)
        #expect(tests.visibleFiles.map(\.name) == ["main.go", "main_test.go"])
        #expect(tests.visibleFiles.filter { !$0.isReadonly }.map(\.name) == ["main.go"])

        let quizLesson = try #require(try await db.lesson(Fixtures.quizLesson))
        #expect(quizLesson.hints.isEmpty)
        guard case .choice(let quiz) = quizLesson.content else { Issue.record("quiz"); return }
        #expect(quiz.questions.count == 2 && quiz.questions.map(\.position) == [0, 1])
        #expect(quiz.questions[0].answers == ["`fn`", "`func`", "`def`"])
        #expect(quiz.questions[0].isCorrect("`func`") && !quiz.questions[0].isCorrect("`fn`"))
        #expect(quiz.questions[1].answer == "Yes")
        // Format 7: why each answer is right or wrong, where the lesson says.
        #expect(quiz.questions[0].explanation(of: 0) == nil)
        #expect(quiz.questions[0].explanation(of: 1) == "Go declares every function with func.")
        #expect(quiz.questions[0].explanation(of: 2) == "def is Python's keyword.")
        #expect(quiz.questions[1].explanation(of: 0) == nil && quiz.questions[1].explanation(of: 7) == nil)
        #expect(try await db.quiz(Fixtures.quizLesson)?.questions.count == 2)
        #expect(try await db.quiz(Fixtures.codeLesson) == nil)

        guard case .cli(let cli)? = try await db.lesson(Fixtures.cliLesson)?.content else { Issue.record("cli"); return }
        #expect(cli.defaultBaseURL == "http://localhost:8080" && !cli.asksForBaseURL)
        #expect(cli.steps.count == 2 && cli.usesBaseURL && cli.hasCommands)
        guard case .command(let c) = cli.steps[0].action, case .http(let h) = cli.steps[1].action else { Issue.record("steps"); return }
        #expect(c.command == "echo built" && c.tests.count == 3 && c.variables == [NameValue("word", "(b\\w+)")])
        if case .stdoutJq(let q, let jsonl, let expected) = c.tests[2] {
            #expect(q == ".ok" && !jsonl && expected == [JqExpectation("==", .bool(true))])
        } else { Issue.record("jq test") }
        #expect(h.method == "POST" && !h.followRedirects && h.headers == [NameValue("X-Trace", "abc")])
        if case .json(let body)? = h.body { #expect(body.contains("ada@example.com")) } else { Issue.record("json body") }
        #expect(h.bodyVariables == [.init(name: "token", jsonPath: ".token")])
        #expect(h.headerVariables == [.init(name: "session", header: "Set-Cookie", regex: "^(session_id=[^;]+)")])
        #expect(h.tests.count == 3)
        if case .jsonValue(let path, let op, let value) = h.tests[1] { #expect(path == ".user.id" && op == "gte" && value == .int(7)) }
        if case .trailersEqual(let k, let v) = h.tests[2] { #expect(k == "X-Sum" && v == "42") }
        #expect(cli.steps[1].sleepAfterMs == 10)

        guard case .githubChecks(let gh)? = try await db.lesson(Fixtures.githubLesson)?.content else { Issue.record("github"); return }
        #expect(gh.steps.map(\.kind) == ["workflow_status", "workflow_logs_contains"])
        #expect(gh.steps[1].summary.contains("staticcheck") && gh.steps[1].summary.contains("Style"))

        guard case .interviewChat(let interview)? = try await db.lesson(Fixtures.interviewLesson)?.content else { Issue.record("interview"); return }
        #expect(interview.studentPrompt == "Explain tracing." && interview.acceptanceCriteria == "Mentions spans.")
        // No criteria is NULL, and the answer section shows none.
        guard case .interviewChat(let open)? = try await db.lesson(Fixtures.optionalLesson)?.content else { Issue.record("optional"); return }
        #expect(open.acceptanceCriteria == nil)
        #expect(try await db.lesson("no-such-lesson") == nil)
    }

    @Test func fullTextSearch() async throws {
        let (db, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        let hits = try await db.search("greeting")
        #expect(hits.map(\.ref) == [Fixtures.ref(Fixtures.codeLesson)])
        let found = try await ReadmeTexts().search("greeting", in: [db])
        #expect(found.hits == hits)
        #expect(found.snippets[Fixtures.ref(Fixtures.codeLesson)]?.contains { $0.matched && $0.text.lowercased() == "greeting" } == true)
        #expect(try await db.search("gre").map(\.ref.lesson) == [Fixtures.codeLesson]) // the last word is a prefix
        #expect(try await db.search("Gopher").first?.ref.lesson == Fixtures.codeLesson) // title matches rank first
        #expect(try await db.search("\"quoted\" OR (stuff").isEmpty) // FTS syntax stays literal
        #expect(try await db.search("   ").isEmpty)
    }

    /// Search goes through every library, in the sidebar's order, and keeps which one a hit
    /// is from: lesson uuids only need to be unique within their library.
    @Test func searchSpansLibraries() async throws {
        let (first, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        let other = dir.appendingPathComponent("other.library")
        try Fixtures.makeDatabase(at: other, id: "b0000000-0000-0000-0000-000000000002", title: "Other")
        let second = try AppDatabase(url: other)
        let found = try await ReadmeTexts().search("greeting", in: [first, second])
        #expect(found.hits.map(\.ref) == [Fixtures.ref(Fixtures.codeLesson), LessonRef(library: second.info.uuid, lesson: Fixtures.codeLesson)])
        #expect(found.snippets.count == 2)
        #expect(try await ReadmeTexts().search("greeting", in: [second, first], limit: 1).hits.map(\.ref.library) == [second.info.uuid])
    }

    /// The index also holds words a reader never sees: a code block's language, a callout's
    /// marker, a link's address. Search keeps a hit only when its title or its readme as
    /// rendered shows a match, so every hit has a word to point at.
    @Test func searchKeepsOnlyMatchesYouCanSee() async throws {
        let readmes = [
            Fixtures.testsLesson: "# Adding Numbers\n\nComplete `add` for the café.\n\n```zig\nconst sum = 1 + 2;\n```",
            // The address names warning three times, so the index ranks this lesson first for it.
            Fixtures.quizLesson: "# Functions Quiz\n\n> [!info]\n> Each answer counts once.\n\nSee [the guide](https://kettle.example.org/warning/warning/warning).",
            Fixtures.cliLesson: "# Your Server\n\n> [!warning]\n> Stop it when you are done.",
        ]
        let (db, dir) = try Fixtures.database { db in
            for (id, readme) in readmes { try db.execute(sql: "UPDATE lesson SET readme = ? WHERE uuid = ?", arguments: [readme, id]) }
        }
        defer { try? FileManager.default.removeItem(at: dir) }
        func ref(_ lesson: String) -> LessonRef { Fixtures.ref(lesson) }
        let texts = ReadmeTexts()
        // A fence's language, a callout's marker (info is drawn as Note), a link's address.
        for word in ["zig", "info", "kettle", "example", "https"] {
            #expect(try await db.search(word).isEmpty == false, "the index holds \(word)")
            #expect(try await texts.search(word, in: [db]).hits.isEmpty, "\(word) shows nowhere")
        }
        // A callout's title is drawn, so its marker's word is a match you can see; the hit the
        // index ranks first only had it in an address.
        #expect(try await db.search("warning").map(\.ref) == [ref(Fixtures.quizLesson), ref(Fixtures.cliLesson)])
        let warning = try await texts.search("warning", in: [db])
        #expect(warning.hits.map(\.ref) == [ref(Fixtures.cliLesson)])
        #expect(warning.snippets[ref(Fixtures.cliLesson)]?.filter(\.matched).map(\.text) == ["Warning"])
        // Dropped hits leave room for the next best.
        #expect(try await texts.search("warning", in: [db], limit: 1).hits.map(\.ref) == [ref(Fixtures.cliLesson)])
        // The title counts, code counts, and diacritics fold as the index folds them.
        #expect(try await texts.search("quiz", in: [db]).hits.map(\.ref) == [ref(Fixtures.quizLesson)])
        #expect(try await texts.search("sum", in: [db]).hits.map(\.ref) == [ref(Fixtures.testsLesson)])
        let cafe = try await texts.search("cafe", in: [db])
        #expect(cafe.hits.map(\.ref) == [ref(Fixtures.testsLesson)])
        #expect(cafe.snippets[ref(Fixtures.testsLesson)]?.filter(\.matched).map(\.text) == ["café"])
    }
}

/// v4: where a command-line lesson's ${baseURL} comes from, from the library's suggestion and
/// the learner database's addresses and projects.
@Suite struct ServerAddressTests {
    @Test func baseURLFromTheSuggestionYourAddressOrYourServer() async throws {
        let (db, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        let learner = try await Fixtures.learner(in: dir)
        let library = Fixtures.libraryID
        guard case .cli(let suggests)? = try await db.lesson(Fixtures.cliLesson)?.content,
              case .cli(let asks)? = try await db.lesson(Fixtures.askingLesson)?.content else { Issue.record("cli lessons"); return }
        #expect(asks.asksForBaseURL && asks.baseURLDefault == nil && asks.defaultBaseURL == nil)

        // A suggested address, until you set yours; "Use Suggested" removes your row.
        #expect(try await learner.baseURL(for: suggests, library: library, course: Fixtures.courseID) == .suggested("http://localhost:8080", yours: nil))
        try await learner.setServerAddress("http://192.168.1.20:8080", for: "http://localhost:8080")
        let replaced = try await learner.baseURL(for: suggests, library: library, course: Fixtures.courseID)
        #expect(replaced == .suggested("http://localhost:8080", yours: "http://192.168.1.20:8080") && replaced.url == "http://192.168.1.20:8080")
        try await learner.setServerAddress("http://localhost:8080", for: "http://localhost:8080") // the suggestion itself
        #expect(try await learner.baseURL(for: suggests, library: library, course: Fixtures.courseID) == .suggested("http://localhost:8080", yours: nil))
        try await learner.setServerAddress("http://localhost:9090", for: "http://localhost:8080")
        try await learner.setServerAddress(nil, for: "http://localhost:8080")
        #expect(try await learner.baseURL(for: suggests, library: library, course: Fixtures.courseID).url == "http://localhost:8080")

        // Your own server: missing until the course's project has it; then the project wins.
        let missing = try await learner.baseURL(for: asks, library: library, course: Fixtures.courseID)
        #expect(missing == .yourServer(nil) && missing.isMissing && missing.url == nil)
        try await learner.saveProject(library: library, course: Fixtures.courseID, Project(dir: nil, baseURL: "https://quill.example.com"))
        #expect(try await learner.baseURL(for: asks, library: library, course: Fixtures.courseID).url == "https://quill.example.com")

        // Neither: no ${baseURL}.
        #expect(BaseURL(lesson: CLILesson(steps: []), projectBaseURL: "https://x.example.com", serverAddress: nil) == .none)
    }
}

/// A quiz: graded answers go into the attempt history, and all right finishes the lesson.
@Suite struct QuizTests {
    @Test func aQuizIsDoneWhenEveryAnswerIsRight() async throws {
        let (db, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        let learner = try await Fixtures.learner(in: dir)
        let ref = Fixtures.ref(Fixtures.quizLesson)
        guard case .choice(let quiz)? = try await db.lesson(Fixtures.quizLesson)?.content else { Issue.record("quiz"); return }
        var state = QuizState(quiz: quiz)
        let early = state.submit()
        #expect(!state.canSubmit && early.isEmpty)

        // One wrong, one right.
        let picked = [state.pick(0, for: 0), state.pick(0, for: 1), state.pick(9, for: 0)]
        #expect(picked == [true, true, false]) // no answer 9
        let first = state.submit()
        #expect(first == [QuizState.Answer(question: 0, text: "`fn`", right: false), QuizState.Answer(question: 1, text: "Yes", right: true)])
        #expect(state.verdicts == [0: false, 1: true] && !state.allCorrect && !state.canSubmit)
        let changedRight = state.pick(1, for: 1)
        #expect(!changedRight) // a right answer stays
        try await learner.recordAnswers(ref, first, now: Date(timeIntervalSince1970: 1_800_000_000))

        // The wrong one again, now right: the quiz is passed and the lesson done.
        let repicked = state.pick(1, for: 0)
        #expect(repicked && state.canSubmit)
        let second = state.submit()
        #expect(second == [QuizState.Answer(question: 0, text: "`func`", right: true)] && state.allCorrect && !state.canSubmit)
        try await learner.recordAnswers(ref, second, now: Date(timeIntervalSince1970: 1_800_003_600))
        if state.allCorrect { try await learner.setDone(ref, true) }
        #expect(try await learner.done() == [Fixtures.libraryID: [Fixtures.quizLesson]])

        let attempts = try Fixtures.rows(in: dir, "SELECT kind, passed, output FROM attempt WHERE lesson_uuid = ? ORDER BY at DESC, id DESC", [Fixtures.quizLesson])
        #expect(attempts.map { $0["output"] as String } == ["`func`", "Yes", "`fn`"])
        #expect(attempts.map { $0["passed"] as Int } == [1, 1, 0] && attempts.allSatisfy { $0["kind"] as String == "answer" })
    }
}
