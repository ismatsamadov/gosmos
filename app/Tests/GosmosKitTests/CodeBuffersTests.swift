import Foundation
import SQLite3
import Testing
@testable import GosmosKit

/// What the buffers asked to remember across launches: whether a file is still in conflict.
@MainActor private final class Remembered {
    var files: [String: Bool] = [:]
}

/// Whether something happened yet, such as the lesson closing.
@MainActor private final class Flag {
    var isSet = false
}

/// Holds a database's write lock on a connection of its own, as another process writing would,
/// so the app's next write waits (its busy timeout is 10 s) until `release`. SQLite's own API:
/// GRDB ends every access with no transaction left open.
private final class WriteLock {
    struct Failed: Error { let code: Int32 }
    private var db: OpaquePointer?

    init(_ url: URL) throws {
        var code = sqlite3_open_v2(url.path, &db, SQLITE_OPEN_READWRITE, nil)
        if code == SQLITE_OK { code = sqlite3_exec(db, "BEGIN IMMEDIATE", nil, nil, nil) }
        guard code == SQLITE_OK else { sqlite3_close(db); throw Failed(code: code) }
    }

    func release() throws {
        let code = sqlite3_exec(db, "COMMIT", nil, nil, nil)
        sqlite3_close(db)
        guard code == SQLITE_OK else { throw Failed(code: code) }
    }
}

/// The editor's buffers before, during and after the load: typing is never lost to a load
/// that lands late, nor overwritten by it.
@Suite struct CodeBuffersTests {
    struct Bench {
        let db: AppDatabase
        let learner: LearnerDatabase
        let dir: URL
        let lesson: LessonDetail
        let code: CodeExercise
        let workspace: Workspace

        init() async throws {
            let (db, dir) = try Fixtures.database()
            self.db = db
            self.dir = dir
            learner = try await Fixtures.learner(in: dir)
            lesson = try #require(try await db.lesson(Fixtures.codeLesson))
            code = try #require(lesson.code)
            workspace = Workspace(root: dir.appendingPathComponent("work"))
        }

        var file: URL { workspace.folder(for: lesson).appendingPathComponent("main.go") }
        func stored() async throws -> String? { try await learner.codeFiles(for: lesson.ref)["main.go"]?.content }
        func onDisk() throws -> String { try String(contentsOf: file, encoding: .utf8) }

        /// Buffers opened with `saved` as code_file had it, whose load waits at `gate`.
        @MainActor
        fileprivate func buffers(saved: [String: CodeFileRecord] = [:], gate: Gate, remembered: @escaping @MainActor (String, Bool) -> Void = { _, _ in }) -> CodeBuffers {
            CodeBuffers(lesson: lesson, exercise: code, saved: saved, learner: learner, workspace: workspace, unsettled: [],
                        goVersion: { await gate.pass(); return "1.27" }, remember: remembered)
        }
    }

    /// The load is slow (Go is still being found); you type meanwhile and leave the lesson.
    /// Leaving waits for the load, then saves what you typed.
    @Test @MainActor func typingDuringASlowLoadIsSavedOnLeaving() async throws {
        let bench = try await Bench()
        defer { try? FileManager.default.removeItem(at: bench.dir) }
        let gate = Gate()
        let buffers = bench.buffers(gate: gate)
        let main = try #require(buffers.files.first { $0.name == "main.go" })
        let typed = "package main\n\nfunc main() {} // typed while it loaded\n"

        let loading = Task { await buffers.load() }
        await gate.somethingWaits()
        #expect(!buffers.isLoaded)
        main.text = typed
        let leaving = Task { try await buffers.settle(within: nil) }
        for _ in 0..<5 { await Task.yield() }
        #expect(!buffers.isLoaded) // the load is still held, and leaving waits for it
        #expect(try await bench.stored() == nil)

        await gate.open()
        #expect(try await leaving.value.isEmpty)
        await loading.value
        #expect(buffers.isLoaded && main.text == typed && buffers.conflict == nil)
        #expect(try await bench.stored() == typed && bench.onDisk() == typed)
    }

    /// Another editor changed the file in the lesson folder; you typed here before the load
    /// read it. The load keeps your text and offers theirs as the other side of a conflict.
    @Test @MainActor func aLoadLandingAfterTypingNeverOverwritesIt() async throws {
        let bench = try await Bench()
        defer { try? FileManager.default.removeItem(at: bench.dir) }
        let earlier = CodeFileRecord(name: "main.go", content: "package main // saved\n", updatedAt: Fixtures.ts)
        try await bench.learner.saveCode(bench.lesson.ref, files: [earlier])
        try FileManager.default.createDirectory(at: bench.file.deletingLastPathComponent(), withIntermediateDirectories: true)
        try "package main // theirs\n".write(to: bench.file, atomically: true, encoding: .utf8) // newer than code_file
        let gate = Gate(), remembered = Remembered()
        let buffers = bench.buffers(saved: ["main.go": earlier], gate: gate) { remembered.files[$0] = $1 }
        let main = try #require(buffers.files.first { $0.name == "main.go" })
        #expect(main.text == earlier.content)

        let loading = Task { await buffers.load() }
        await gate.somethingWaits()
        main.text = "package main // mine\n"
        await gate.open()
        await loading.value
        #expect(buffers.isLoaded)
        #expect(main.text == "package main // mine\n") // not overwritten
        #expect(try await bench.stored() == earlier.content && bench.onDisk() == "package main // theirs\n") // neither side written

        // The first save asks which to keep.
        #expect(try await buffers.flush() == false)
        #expect(buffers.conflict == "main.go")
        #expect(try await bench.stored() == earlier.content && bench.onDisk() == "package main // theirs\n")
        #expect(try await buffers.resolveConflict(keepMine: true)?.value == true)
        #expect(buffers.conflict == nil && remembered.files == ["main.go": false])
        #expect(try await bench.stored() == "package main // mine\n" && bench.onDisk() == "package main // mine\n")
    }

    /// You choose the file on disk while saving that choice is slow (another connection holds
    /// the learner database), and the text changes as soon as the question goes away. The editor
    /// takes no typing until the choice lands, and should the text change anyway, the file on
    /// disk doesn't go over it: the change was made on your version, so the lesson asks again.
    @Test @MainActor func aSlowChoiceOfTheFileOnDiskNeverOverwritesTextChangedMeanwhile() async throws {
        let bench = try await Bench()
        defer { try? FileManager.default.removeItem(at: bench.dir) }
        let gate = Gate()
        await gate.open()
        let buffers = bench.buffers(gate: gate)
        await buffers.load()
        let main = try #require(buffers.files.first { $0.name == "main.go" })
        main.text = "package main // mine\n"
        try "package main // theirs\n".write(to: bench.file, atomically: true, encoding: .utf8)
        #expect(try await buffers.flush() == false && buffers.conflict == "main.go")

        let lock = try WriteLock(bench.dir.appendingPathComponent("learner.sqlite"))
        let resolving = try #require(buffers.resolveConflict(keepMine: false))
        #expect(buffers.conflict == nil && buffers.isResolving)
        // A save asked for now waits for the choice, rather than build on the text it replaces.
        let saving = Task { try await buffers.flush() }
        try await Task.sleep(for: .milliseconds(100))
        let typed = "package main // typed after choosing\n"
        main.text = typed
        try lock.release()
        #expect(try await resolving.value == false)
        #expect(main.text == typed && buffers.conflict == "main.go" && !buffers.isResolving)
        #expect(try await saving.value == false)
        #expect(try await bench.stored() == "package main // theirs\n" && bench.onDisk() == "package main // theirs\n")

        #expect(try await buffers.resolveConflict(keepMine: true)?.value == true)
        #expect(try await bench.stored() == typed && bench.onDisk() == typed)
    }

    /// The same slow choice with nothing typed: the editor gets the file on disk once it lands.
    @Test @MainActor func aSlowChoiceOfTheFileOnDiskFillsTheEditorWhenItLands() async throws {
        let bench = try await Bench()
        defer { try? FileManager.default.removeItem(at: bench.dir) }
        let gate = Gate()
        await gate.open()
        let buffers = bench.buffers(gate: gate)
        await buffers.load()
        let main = try #require(buffers.files.first { $0.name == "main.go" })
        main.text = "package main // mine\n"
        try "package main // theirs\n".write(to: bench.file, atomically: true, encoding: .utf8)
        #expect(try await buffers.flush() == false)

        let lock = try WriteLock(bench.dir.appendingPathComponent("learner.sqlite"))
        let resolving = try #require(buffers.resolveConflict(keepMine: false))
        let saving = Task { try await buffers.flush() }
        try await Task.sleep(for: .milliseconds(100))
        try lock.release()
        #expect(try await resolving.value == true)
        #expect(try await saving.value == true)
        #expect(main.text == "package main // theirs\n" && buffers.conflict == nil && !buffers.isResolving)
        #expect(try await bench.stored() == "package main // theirs\n" && bench.onDisk() == "package main // theirs\n")
    }

    /// Leaving the lesson while a conflict is asked about: a choice can leave another to make
    /// (here the text changed while the first landed), and the lesson stays until every choice
    /// is made, so nothing typed goes with it.
    @Test @MainActor func leavingWaitsForEveryChoice() async throws {
        let bench = try await Bench()
        defer { try? FileManager.default.removeItem(at: bench.dir) }
        let gate = Gate()
        await gate.open()
        let buffers = bench.buffers(gate: gate)
        await buffers.load()
        let main = try #require(buffers.files.first { $0.name == "main.go" })
        main.text = "package main // mine\n"
        try "package main // theirs\n".write(to: bench.file, atomically: true, encoding: .utf8)
        let left = Flag()
        let leaving = Task { () throws -> [String] in
            let kept = try await buffers.settle(within: nil)
            left.isSet = true
            return kept
        }
        while buffers.conflict == nil { await Task.yield() }

        let lock = try WriteLock(bench.dir.appendingPathComponent("learner.sqlite"))
        let resolving = try #require(buffers.resolveConflict(keepMine: false))
        try await Task.sleep(for: .milliseconds(100))
        let typed = "package main // typed after choosing\n"
        main.text = typed
        try lock.release()
        #expect(try await resolving.value == false && buffers.conflict == "main.go")
        try await Task.sleep(for: .milliseconds(100))
        #expect(!left.isSet) // still there, asking

        #expect(try await buffers.resolveConflict(keepMine: true)?.value == true)
        #expect(try await leaving.value.isEmpty)
        #expect(try await bench.stored() == typed && bench.onDisk() == typed)
    }

    /// Without typing, the load takes in the other editor's newer edit, as before.
    @Test @MainActor func aLoadWithoutTypingTakesInTheFolder() async throws {
        let bench = try await Bench()
        defer { try? FileManager.default.removeItem(at: bench.dir) }
        let earlier = CodeFileRecord(name: "main.go", content: "package main // saved\n", updatedAt: Fixtures.ts)
        try await bench.learner.saveCode(bench.lesson.ref, files: [earlier])
        try FileManager.default.createDirectory(at: bench.file.deletingLastPathComponent(), withIntermediateDirectories: true)
        try "package main // theirs\n".write(to: bench.file, atomically: true, encoding: .utf8)
        let gate = Gate()
        await gate.open()
        let buffers = bench.buffers(saved: ["main.go": earlier], gate: gate)
        await buffers.load()
        #expect(buffers.isLoaded && buffers.files.first { $0.name == "main.go" }?.text == "package main // theirs\n")
        #expect(try await bench.stored() == "package main // theirs\n")
        #expect(try await buffers.flush() && buffers.conflict == nil)
    }
}

extension CodeBuffersTests {
    /// A review redoes a lesson from its starter code, in a folder of its own: your saved code
    /// and the lesson's folder stay as they are, and the next review starts over.
    @Test @MainActor func aReviewStartsFromTheStarterAndKeepsNothing() async throws {
        let bench = try await Bench()
        defer { try? FileManager.default.removeItem(at: bench.dir) }
        let starter = try #require(bench.code.starterFiles.first { $0.name == "main.go" }?.content)
        let mine = CodeFileRecord(name: "main.go", content: "package main // mine\n", updatedAt: Fixtures.ts)
        try await bench.learner.saveCode(bench.lesson.ref, files: [mine])
        let review = bench.workspace.review
        let gate = Gate()
        await gate.open()
        func reviewBuffers() -> CodeBuffers {
            CodeBuffers(lesson: bench.lesson, exercise: bench.code, saved: [:], learner: nil, workspace: review, unsettled: [],
                        goVersion: { await gate.pass(); return "1.27" }, remember: { _, _ in })
        }
        let buffers = reviewBuffers()
        await buffers.load()
        let main = try #require(buffers.files.first { $0.name == "main.go" })
        #expect(buffers.isLoaded && main.text == starter)
        #expect(buffers.folder == review.folder(for: bench.lesson) && buffers.folder != bench.workspace.folder(for: bench.lesson))
        main.text = "package main // review try\n"
        #expect(try await buffers.flush())
        let tried = review.folder(for: bench.lesson).appendingPathComponent("main.go")
        #expect(try String(contentsOf: tried, encoding: .utf8) == "package main // review try\n")
        #expect(try await bench.stored() == mine.content) // your saved code is untouched
        #expect(!FileManager.default.fileExists(atPath: bench.file.path)) // and the lesson's folder too

        let again = reviewBuffers()
        await again.load()
        #expect(again.files.first { $0.name == "main.go" }?.text == starter)
    }

    /// The conflicts a review meets in its own folder stay there: the lesson's list of files
    /// left in conflict (which decides what its next opening keeps) is never written, and
    /// leaving keeps nothing.
    @Test @MainActor func aReviewsConflictsAreNotRemembered() async throws {
        let bench = try await Bench()
        defer { try? FileManager.default.removeItem(at: bench.dir) }
        let review = bench.workspace.review
        let gate = Gate(), remembered = Remembered()
        await gate.open()
        let buffers = CodeBuffers(lesson: bench.lesson, exercise: bench.code, saved: [:], learner: nil, workspace: review, unsettled: [],
                                  goVersion: { await gate.pass(); return "1.27" }, remember: { remembered.files[$0] = $1 })
        await buffers.load()
        let main = try #require(buffers.files.first { $0.name == "main.go" })
        try "package main // another editor\n".write(to: review.folder(for: bench.lesson).appendingPathComponent("main.go"),
                                                     atomically: true, encoding: .utf8)
        main.text = "package main // redo\n"
        #expect(try await buffers.flush() == false && buffers.conflict == "main.go")
        // Quitting with no choice made: a review has nothing to keep.
        #expect(try await buffers.settle(within: .milliseconds(50)).isEmpty)
        #expect(try await buffers.resolveConflict(keepMine: true)?.value == true)
        #expect(remembered.files.isEmpty)
    }
}
