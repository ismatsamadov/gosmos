import Foundation
import GRDB
import Testing
@testable import GosmosKit

/// The time budgets of docs/architecture.md, on a library the size of a real catalog.
@Suite(.serialized) struct BudgetTests {
    static func bigDatabase() throws -> (AppDatabase, URL) {
        let dir = try Fixtures.tempDir("big")
        let url = dir.appendingPathComponent(Fixtures.libraryName)
        try Fixtures.makeDatabase(at: url)
        let queue = try DatabaseQueue(path: url.path)
        let words = ["goroutine", "channel", "interface", "struct", "pointer", "slice", "map", "closure", "mutex", "context",
                     "handler", "middleware", "json", "token", "cookie", "server", "client", "request", "response", "error"]
        try queue.write { db in
            for c in 0..<13 {
                let course = "bc\(c)"
                try db.execute(sql: "INSERT INTO course(uuid, slug, title, kind, short_description, description, position) VALUES (?,?,?,?,?,?,?)",
                               arguments: [course, "course-\(c)", "Course \(c)", "Course", "", "", c + 10])
                for h in 0..<8 {
                    let chapter = "\(course)h\(h)"
                    try db.execute(sql: "INSERT INTO chapter(uuid, course_uuid, position, slug, title, description) VALUES (?,?,?,?,?,?)",
                                   arguments: [chapter, course, h, "chapter-\(h)", "Chapter \(h)", ""])
                    for l in 0..<7 {
                        var readme = "# Lesson \(c)-\(h)-\(l)\n\n"
                        for p in 0..<12 {
                            readme += (0..<60).map { words[($0 * 7 + p * 3 + l + h + c) % words.count] }.joined(separator: " ") + ".\n\n"
                            if p % 4 == 0 { readme += "```go\nfunc main() {\n\tfmt.Println(\"hi\")\n}\n```\n\n" }
                        }
                        try db.execute(sql: """
                            INSERT INTO lesson(uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
                            VALUES (?,?,?,?,?,?,?,?,?,?,?)
                            """, arguments: ["\(chapter)l\(l)", chapter, l, "lesson-\(l)", "Lesson \(words[l + h]) \(c)-\(h)-\(l)", "type_choice", 1, 2,
                                             readme, Fixtures.ts, Fixtures.ts])
                        try db.execute(sql: "INSERT INTO quiz_lesson(lesson_uuid, contains_complete_dir) VALUES (?,?)",
                                       arguments: ["\(chapter)l\(l)", 0])
                    }
                }
            }
        }
        try queue.close()
        return (try AppDatabase(url: url), dir)
    }

    @Test func searchCatalogLessonAndReadmeAreFast() async throws {
        let (db, dir) = try Self.bigDatabase()
        defer { try? FileManager.default.removeItem(at: dir) }
        let clock = ContinuousClock()
        #expect(try await db.catalog().lessonCount > 700)

        // Search while typing, as the app does it (the readmes' text made when the library
        // opens): each prefix under 50 ms, hits kept only where a match shows, with excerpts.
        let texts = ReadmeTexts()
        try await texts.warm(db)
        var slowest = Duration.zero
        for prefix in ["g", "go", "gor", "goro", "gorou", "goroutine", "goroutine ch", "goroutine channel"] {
            let t = clock.now
            let found = try await texts.search(prefix, in: [db])
            slowest = max(slowest, clock.now - t)
            #expect(!found.hits.isEmpty && found.hits.count <= 60, "\(prefix)")
        }
        #expect(slowest < .milliseconds(50), "slowest search \(slowest)")

        // A lesson switch reads the lesson and parses its readme: under 100 ms.
        let t = clock.now
        let lesson = try #require(try await db.lesson("bc5h3l2"))
        _ = Readme.parse(lesson.readme, droppingTitle: lesson.title)
        #expect(clock.now - t < .milliseconds(100), "lesson switch \(clock.now - t)")

        let c = clock.now
        _ = try await db.catalog()
        #expect(clock.now - c < .milliseconds(200), "catalog \(clock.now - c)")
    }

    /// The documentation's budgets, on the full Go reference the app ships: a package's page
    /// opens in under 100 ms, the biggest ones too, and each keystroke of a search takes under
    /// 30 ms.
    @Test(.enabled(if: FileManager.default.fileExists(atPath: ReferenceFixture.fullURL.path)))
    func documentationPagesAndSearchAreFast() async throws {
        let reference = try ReferenceDatabase(url: ReferenceFixture.fullURL)
        let clock = ContinuousClock()
        // As at launch, when the names for links are read.
        #expect(!(try reference.names()).isEmpty)

        // Each the best of three, so the other tests running meanwhile don't count.
        func best(_ work: () async throws -> Void) async rethrows -> Duration {
            var best = Duration.seconds(60)
            for _ in 0..<3 {
                let t = clock.now
                try await work()
                best = min(best, clock.now - t)
            }
            return best
        }

        var slowest = Duration.zero
        for path in ["net/http", "syscall", "debug/elf", "reflect", "go/types", "strings", "fmt", "os"] {
            slowest = max(slowest, try await best {
                let page = try await reference.page(.package(path))
                #expect(page != nil, "\(path)")
            })
        }
        #expect(slowest < .milliseconds(100), "slowest package page \(slowest)")

        var slowestSearch = Duration.zero
        for query in ["p", "pr", "pri", "prin", "print", "printl", "println", "B", "Bu", "Builder", "strings.B", "strings.Builder.Gr",
                      "g", "go", "gor", "goroutine", "slice of", "http.Cl", "time.Dur"] {
            slowestSearch = max(slowestSearch, try await best {
                let hits = try await reference.search(query)
                #expect(!hits.isEmpty, "\(query)")
            })
        }
        #expect(slowestSearch < .milliseconds(30), "slowest search \(slowestSearch)")
        #expect(try await reference.search("Println").first?.title == "fmt.Println")
        #expect(try await reference.search("Builder").first?.title == "strings.Builder")
    }
}

/// Quitting stops lesson programs, and only the kind asked for.
@Suite struct QuitTests {
    @Test func quittingStopsOnlyTheKindAskedFor() async throws {
        let env = ["PATH": "/usr/bin:/bin"]
        let other = Task { await ProcessRunner.run(ProcessRequest(executable: "/bin/sleep", arguments: ["1"], directory: "/", environment: env, timeout: .seconds(10), kind: ProcessKind("other-test"))) }
        let lesson = Task { await ProcessRunner.run(ProcessRequest(executable: "/bin/sleep", arguments: ["30"], directory: "/", environment: env, timeout: .seconds(60), kind: ProcessKind("quit-test"))) }
        try await Task.sleep(for: .milliseconds(300))
        ProcessRunner.terminateAll(ProcessKind("quit-test"))
        let stopped = await lesson.value
        #expect(stopped.signal == SIGKILL && stopped.duration < .seconds(10))
        let finished = await other.value
        #expect(finished.exitCode == 0 && finished.signal == 0)
    }
}
