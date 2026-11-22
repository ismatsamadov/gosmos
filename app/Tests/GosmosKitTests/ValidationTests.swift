import Foundation
import GRDB
import Testing
@testable import GosmosKit

/// Each check docs/library-format.md lists refuses a file that breaks it, with the same words
/// as the library tool, and writes nothing to the file. A library is only ever read.
@Suite struct ValidationTests {
    /// A good library, broken by `change` (run with foreign keys off), then opened.
    func open(_ change: (Database) throws -> Void) throws -> (error: AppDatabaseError?, untouched: Bool) {
        let dir = try Fixtures.tempDir("broken")
        defer { try? FileManager.default.removeItem(at: dir) }
        let url = dir.appendingPathComponent(Fixtures.libraryName)
        try Fixtures.makeDatabase(at: url)
        let queue = try DatabaseQueue(path: url.path)
        try queue.writeWithoutTransaction { db in
            try db.execute(sql: "PRAGMA foreign_keys = OFF")
            try change(db)
        }
        try queue.close()
        let before = try Data(contentsOf: url)
        let opened = Result { () throws(AppDatabaseError) in _ = try AppDatabase(url: url) }
        let untouched = try Data(contentsOf: url) == before
        if case .failure(let error) = opened { return (error, untouched) }
        return (nil, untouched)
    }

    /// The cases in db/testdata/checks, which the library tool's tests read too
    /// (db/checks_test.go): each a library broken for one check, and the problems the checks
    /// report for it. The app refuses each in the tool's words, problem for problem, and writes
    /// nothing to the file; the case that breaks nothing opens.
    @Test(arguments: SharedChecks.cases()) func sharedCase(_ name: String) throws {
        let dir = try Fixtures.tempDir("checks")
        defer { try? FileManager.default.removeItem(at: dir) }
        let url = dir.appendingPathComponent(Fixtures.libraryName)
        try SharedChecks.make(name, at: url)
        let before = try Data(contentsOf: url)
        var problems: [String] = []
        do {
            _ = try AppDatabase(url: url)
        } catch {
            problems = error.problems
        }
        #expect(problems.map { $0 + "\n" }.joined() == (try SharedChecks.want(name)))
        #expect(try Data(contentsOf: url) == before, "a refused file must not be written")
    }

    /// An older library is made again by whoever made it; nothing brings it up to date.
    @Test func anOlderLibraryIsMadeAgain() throws {
        let older = try open { try $0.execute(sql: "PRAGMA user_version = 8") }.error
        #expect(older?.title == "Library Too Old" && older?.recoverySuggestion?.contains("make it again") == true)
    }

    @Test func damage() throws {
        let (big, dir) = try BudgetTests.bigDatabase()
        defer { try? FileManager.default.removeItem(at: dir) }
        try big.close()
        let url = dir.appendingPathComponent(Fixtures.libraryName)
        // Overwrite readme pages near the end, so the header reads but quick_check finds damage.
        let handle = try FileHandle(forWritingTo: url)
        let size = try handle.seekToEnd()
        try handle.seek(toOffset: size / 4 * 3)
        try handle.write(contentsOf: Data(repeating: 0xA5, count: 16384))
        try handle.close()
        do {
            _ = try AppDatabase(url: url)
            Issue.record("opened a damaged file")
        } catch {
            guard case .damaged(_, let lines) = error else { Issue.record("\(error)"); return }
            #expect(!lines.isEmpty && lines.allSatisfy { $0.hasPrefix("quick_check: ") })
        }
    }

    /// While another program holds the file locked, opening waits; still locked when waiting
    /// runs out, the library is busy, never damaged or not a database.
    @Test func aLockedLibraryIsBusyNotBroken() async throws {
        let dir = try Fixtures.tempDir("locked")
        defer { try? FileManager.default.removeItem(at: dir) }
        let url = dir.appendingPathComponent(Fixtures.libraryName)
        try Fixtures.makeDatabase(at: url)
        // Another program holds the file locked, readers out too, until it is done.
        let lock = try WriteLock(url, exclusive: true)
        do {
            _ = try AppDatabase(url: url, immutable: false, busyTimeout: 0.2)
            Issue.record("opened a locked file")
        } catch {
            #expect(error == .busy(path: url.path) && error.problems == ["another program is writing to this library; try again"])
            #expect(error.title == "Library Busy")
        }
        // Waiting longer than the lock lasts: the library opens.
        let releasing = Task.detached {
            try await Task.sleep(for: .milliseconds(300))
            try lock.release()
        }
        _ = try AppDatabase(url: url, immutable: false, busyTimeout: 5)
        try await releasing.value
    }

    /// A write cut short leaves a hot journal. Rolling it back is a write, which the app never
    /// makes: the library is refused with what to do, and the file and its journal stay.
    @Test func anInterruptedWriteIsReportedNotRolledBack() throws {
        let dir = try Fixtures.tempDir("interrupted")
        defer { try? FileManager.default.removeItem(at: dir) }
        let url = dir.appendingPathComponent(Fixtures.libraryName), copy = dir.appendingPathComponent("interrupted.library")
        try Fixtures.makeDatabase(at: url)
        // Copy the file and its journal in the middle of a write too big for the cache, so part of
        // it has reached the file, as a crash leaves them.
        let writer = try DatabaseQueue(path: url.path)
        try writer.writeWithoutTransaction { db in
            // The system SQLite spills only past 20,000 cached pages unless told otherwise.
            try db.execute(sql: "PRAGMA cache_size = 10")
            try db.execute(sql: "PRAGMA cache_spill = 1")
            try db.execute(sql: "BEGIN IMMEDIATE")
            try db.execute(sql: "CREATE TABLE spill (x BLOB)")
            try db.execute(sql: """
                WITH RECURSIVE n(i) AS (SELECT 1 UNION ALL SELECT i + 1 FROM n WHERE i < 6000)
                INSERT INTO spill SELECT randomblob(500) FROM n
                """)
            try FileManager.default.copyItem(at: url, to: copy)
            try FileManager.default.copyItem(atPath: url.path + "-journal", toPath: copy.path + "-journal")
            try db.execute(sql: "ROLLBACK")
        }
        try writer.close()
        #expect(try Data(contentsOf: copy) != Data(contentsOf: url)) // the copy holds part of the write
        let bytes = try Data(contentsOf: copy), journal = try Data(contentsOf: URL(fileURLWithPath: copy.path + "-journal"))
        #expect(throws: AppDatabaseError.interruptedWrite(path: copy.path)) { try AppDatabase(url: copy) }
        #expect(try Data(contentsOf: copy) == bytes && Data(contentsOf: URL(fileURLWithPath: copy.path + "-journal")) == journal)
        let refused = AppDatabaseError.interruptedWrite(path: copy.path)
        #expect(refused.title == "Library Needs Recovery" && refused.recoverySuggestion?.contains("go run ./cmd/library check \(copy.path)") == true)
    }

    /// Opening and reading everything writes nothing: not a byte, not the modification time, no
    /// journal. A library built into the app opens immutable, even from a folder no one may
    /// write to, and nothing can be written through the app's connection.
    @Test func librariesAreOnlyRead() async throws {
        let dir = try Fixtures.tempDir("readonly")
        defer {
            try? FileManager.default.setAttributes([.posixPermissions: 0o755], ofItemAtPath: dir.appendingPathComponent("Courses").path)
            try? FileManager.default.removeItem(at: dir)
        }
        let url = dir.appendingPathComponent(Fixtures.libraryName)
        try Fixtures.makeDatabase(at: url)
        let before = try Data(contentsOf: url)
        let modified = try FileManager.default.attributesOfItem(atPath: url.path)[.modificationDate] as? Date
        func readEverything(_ db: AppDatabase) async throws {
            let catalog = try await db.catalog()
            for course in catalog.courses { for lesson in course.lessons { _ = try await db.lesson(lesson.id) } }
            _ = try await db.search("greeting")
            _ = try await db.readmes()
        }
        let db = try AppDatabase(url: url)
        try await readEverything(db)
        await #expect(throws: DatabaseError.self) {
            try await db.dbQueue.writeWithoutTransaction { try $0.execute(sql: "UPDATE lesson SET title = 'x'") }
        }
        try db.close()
        #expect(try Data(contentsOf: url) == before)
        #expect(try FileManager.default.attributesOfItem(atPath: url.path)[.modificationDate] as? Date == modified)
        for suffix in ["-journal", "-wal", "-shm"] { #expect(!FileManager.default.fileExists(atPath: url.path + suffix)) }

        // As built into the app: a read-only file in a read-only folder.
        let courses = dir.appendingPathComponent("Courses")
        try FileManager.default.createDirectory(at: courses, withIntermediateDirectories: true)
        let bundled = courses.appendingPathComponent("practice.sqlite")
        try FileManager.default.copyItem(at: url, to: bundled)
        try FileManager.default.setAttributes([.posixPermissions: 0o444], ofItemAtPath: bundled.path)
        try FileManager.default.setAttributes([.posixPermissions: 0o555], ofItemAtPath: courses.path)
        let builtIn = try AppDatabase(url: bundled, immutable: true)
        #expect(builtIn.info.uuid == Fixtures.libraryID)
        try await readEverything(builtIn)
        try builtIn.close()
        #expect(try Data(contentsOf: bundled) == before)
        #expect(try FileManager.default.contentsOfDirectory(atPath: courses.path) == ["practice.sqlite"])
    }

    @Test func accessIsAskedForFirstAndARefusalIsNamed() throws {
        let dir = FileManager.default.temporaryDirectory.appendingPathComponent("gosmos-tests-access-\(UUID().uuidString)")
        try FileManager.default.createDirectory(at: dir, withIntermediateDirectories: true)
        defer { try? FileManager.default.removeItem(at: dir) }
        let file = dir.appendingPathComponent("library.sqlite")
        try Data("x".utf8).write(to: file)
        #expect(AppDatabase.access(file) == nil)
        #expect(AppDatabase.access(dir.appendingPathComponent("gone.sqlite")) == .missing(path: dir.appendingPathComponent("gone.sqlite").path))

        // File permissions forbid reading: the same refusal macOS gives for a protected folder.
        try FileManager.default.setAttributes([.posixPermissions: 0o000], ofItemAtPath: file.path)
        defer { try? FileManager.default.setAttributes([.posixPermissions: 0o600], ofItemAtPath: file.path) }
        #expect(AppDatabase.access(file) == .notPermitted(path: file.path))
        #expect(throws: AppDatabaseError.notPermitted(path: file.path)) { try AppDatabase(url: file) }
        #expect(AppDatabaseError.notPermitted(path: file.path).title == "No Access to the Library")
    }
}
