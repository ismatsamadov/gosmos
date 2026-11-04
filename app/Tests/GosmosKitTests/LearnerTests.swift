import Foundation
import GRDB
import Testing
@testable import GosmosKit

/// The learner database: made from db/learner.sql when missing, refused when it is not one or
/// is newer, brought up to date when older, and holding everything of yours by library and
/// lesson.
@Suite struct LearnerTests {
    @Test func aNewLearnerStartsWithoutPersonalState() throws {
        let dir = try Fixtures.tempDir("fresh-learner")
        defer { try? FileManager.default.removeItem(at: dir) }
        let url = dir.appendingPathComponent("learner.sqlite")
        let learner = try LearnerDatabase(url: url, schema: Fixtures.learnerSchemaURL)
        try learner.close()

        let queue = try DatabaseQueue(path: url.path)
        defer { try? queue.close() }
        try queue.read { db in
            let tables = try String.fetchAll(db, sql: "SELECT name FROM sqlite_schema WHERE type = 'table' AND name NOT LIKE 'sqlite_%'")
            for table in tables {
                let identifier = "\"" + table.replacingOccurrences(of: "\"", with: "\"\"") + "\""
                let count = try Int.fetchOne(db, sql: "SELECT count(*) FROM \(identifier)")
                #expect(count == 0, "A fresh learner's \(table) table must be empty")
            }
        }
    }

    func header(_ url: URL) throws -> (id: Int, version: Int, tables: [String]) {
        let queue = try DatabaseQueue(path: url.path, configuration: { var c = Configuration(); c.readonly = true; return c }())
        defer { try? queue.close() }
        return try queue.read { db in
            (try Int.fetchOne(db, sql: "PRAGMA application_id") ?? 0, try Int.fetchOne(db, sql: "PRAGMA user_version") ?? 0,
             try String.fetchAll(db, sql: "SELECT name FROM sqlite_schema WHERE type = 'table' ORDER BY name"))
        }
    }

    @Test func madeWhenMissingFromLearnerSQL() async throws {
        let dir = try Fixtures.tempDir("learner")
        defer { try? FileManager.default.removeItem(at: dir) }
        // In a folder that does not exist yet, as Application Support/Gosmos on a new Mac.
        let url = dir.appendingPathComponent("Application Support/Gosmos/learner.sqlite")
        let learner = try LearnerDatabase(url: url, schema: Fixtures.learnerSchemaURL)
        try learner.close()
        let made = try header(url)
        #expect(made.id == LearnerDatabase.applicationID && made.id == 1_196_641_612 && made.version == LearnerDatabase.schemaVersion)
        #expect(made.tables == ["attempt", "code_file", "hint_view", "library", "note", "progress", "review_lesson",
                                "review_question", "setting", "solution_view"])
        // Nothing but the database is left in the folder: it was made aside and moved in whole.
        #expect(try FileManager.default.contentsOfDirectory(atPath: url.deletingLastPathComponent().path) == ["learner.sqlite"])

        // An empty file (nothing ever written) counts as missing.
        let empty = dir.appendingPathComponent("empty.sqlite")
        FileManager.default.createFile(atPath: empty.path, contents: Data())
        try LearnerDatabase(url: empty, schema: Fixtures.learnerSchemaURL).close()
        #expect(try header(empty).id == LearnerDatabase.applicationID)

        // Another copy of the app, or the library tool, made the database while this one was
        // making its own: theirs is in use and stays, and the one made aside goes.
        let theirs = dir.appendingPathComponent("theirs.sqlite")
        let other = try LearnerDatabase(url: theirs, schema: Fixtures.learnerSchemaURL)
        try await other.setSetting("editor_font_size", "17")
        try other.close()
        try LearnerDatabase.make(theirs, from: Fixtures.learnerSchemaURL)
        #expect(try await LearnerDatabase(url: theirs, schema: Fixtures.learnerSchemaURL).setting("editor_font_size") == "17")
        #expect(try FileManager.default.contentsOfDirectory(atPath: dir.path).filter { $0.hasPrefix(".") }.isEmpty)

        // Without learner.sql nothing can be made.
        let nowhere = dir.appendingPathComponent("nowhere.sqlite")
        #expect(throws: LearnerDatabaseError.noSchema) { try LearnerDatabase(url: nowhere, schema: nil) }
        #expect(throws: LearnerDatabaseError.noSchema) { try LearnerDatabase(url: nowhere, schema: dir.appendingPathComponent("missing.sql")) }
        #expect(!FileManager.default.fileExists(atPath: nowhere.path))
    }

    /// What a crash while the database was being made left beside it (the draft, and SQLite's
    /// files for it) goes when the database next opens, once old enough that no copy of the
    /// app can still be making it. Nothing else in the folder is touched, however it is named.
    @Test func draftsACrashLeftBehindGoAtTheNextOpening() throws {
        let fm = FileManager.default
        let dir = try Fixtures.tempDir("learner")
        defer { try? fm.removeItem(at: dir) }
        let folder = dir.appendingPathComponent("Gosmos", isDirectory: true)
        try fm.createDirectory(at: folder, withIntermediateDirectories: true)
        let url = folder.appendingPathComponent("learner.sqlite")
        let anHourAgo = Date().addingTimeInterval(-3600)
        func put(_ name: String, old: Bool = true) throws -> String {
            let path = folder.appendingPathComponent(name).path
            #expect(fm.createFile(atPath: path, contents: Data("x".utf8)))
            if old { try fm.setAttributes([.modificationDate: anHourAgo], ofItemAtPath: path) }
            return name
        }
        let a = UUID().uuidString, b = UUID().uuidString
        // A crash during the first launch: drafts, a journal and WAL files, and no database yet.
        let leftovers = try [".learner.sqlite.\(a)", ".learner.sqlite.\(a)-journal", ".learner.sqlite.\(b)-wal", ".learner.sqlite.\(b)-shm"]
            .map { try put($0) }
        // Another copy of the app making the database right now: its draft is new and stays.
        let making = try put(".learner.sqlite.\(UUID().uuidString)", old: false)
        // Names the app never gives, however close. (The lower-case uuid is another one: this
        // volume may ignore case.)
        let others = try ["learner.sqlite.\(a)", ".learner.sqlite.\(UUID().uuidString.lowercased())", ".learner.sqlite.\(a).bak",
                          ".learner.sqlite.\(a)-journal2", ".learner.sqlite.\(a)-wal-shm", ".learner.sqlite.draft", ".learner.sqlite.",
                          ".notes.sqlite.\(a)", "notes.txt"].map { try put($0) }
        // A folder and a link named as drafts are not drafts either, and the linked file stays.
        let folderDraft = ".learner.sqlite.\(UUID().uuidString)"
        try fm.createDirectory(at: folder.appendingPathComponent(folderDraft), withIntermediateDirectories: false)
        try fm.setAttributes([.modificationDate: anHourAgo], ofItemAtPath: folder.appendingPathComponent(folderDraft).path)
        let linked = dir.appendingPathComponent("elsewhere.txt")
        try "mine".write(to: linked, atomically: true, encoding: .utf8)
        let linkDraft = ".learner.sqlite.\(UUID().uuidString)"
        try fm.createSymbolicLink(at: folder.appendingPathComponent(linkDraft), withDestinationURL: linked)
        let old = timeval(tv_sec: Int(anHourAgo.timeIntervalSince1970), tv_usec: 0)
        #expect(lutimes(folder.appendingPathComponent(linkDraft).path, [old, old]) == 0)

        try LearnerDatabase(url: url, schema: Fixtures.learnerSchemaURL).close()
        let left = Set(try fm.contentsOfDirectory(atPath: folder.path))
        #expect(left.isDisjoint(with: leftovers))
        #expect(left == Set(["learner.sqlite", making, folderDraft, linkDraft] + others))
        #expect(try String(contentsOf: linked, encoding: .utf8) == "mine")

        // One left beside a database in use goes the same way when it next opens.
        let later = try put(".learner.sqlite.\(UUID().uuidString)-journal")
        try LearnerDatabase(url: url, schema: Fixtures.learnerSchemaURL).close()
        #expect(!fm.fileExists(atPath: folder.appendingPathComponent(later).path))
    }

    @Test func refusesAnotherFileAndANewerOneAndLeavesThemAlone() async throws {
        let dir = try Fixtures.tempDir("learner")
        defer { try? FileManager.default.removeItem(at: dir) }
        // A course library is not a learner database.
        let library = dir.appendingPathComponent(Fixtures.libraryName)
        try Fixtures.makeDatabase(at: library)
        let before = try Data(contentsOf: library)
        #expect(throws: LearnerDatabaseError.notLearner(path: library.path, applicationID: AppDatabase.applicationID)) {
            try LearnerDatabase(url: library, schema: Fixtures.learnerSchemaURL)
        }
        #expect(try Data(contentsOf: library) == before)

        // One of a newer format.
        let newer = dir.appendingPathComponent("newer.sqlite")
        try LearnerDatabase(url: newer, schema: Fixtures.learnerSchemaURL).close()
        try await DatabaseQueue(path: newer.path).writeWithoutTransaction { try $0.execute(sql: "PRAGMA user_version = 3") }
        let newerBytes = try Data(contentsOf: newer)
        #expect(throws: LearnerDatabaseError.newerFormat(path: newer.path, version: 3)) { try LearnerDatabase(url: newer, schema: Fixtures.learnerSchemaURL) }
        #expect(try Data(contentsOf: newer) == newerBytes)
        #expect(LearnerDatabaseError.newerFormat(path: newer.path, version: 3).recoverySuggestion == "Update Gosmos to use it.")

        // Not SQLite at all.
        let text = dir.appendingPathComponent("notes.sqlite")
        try String(repeating: "not a database\n", count: 400).write(to: text, atomically: true, encoding: .utf8)
        #expect(throws: LearnerDatabaseError.self) { try LearnerDatabase(url: text, schema: Fixtures.learnerSchemaURL) }
    }

    /// An older file is brought up to date one version at a time, in one transaction; one with
    /// no way up is refused.
    @Test func olderFilesAreMigrated() async throws {
        let dir = try Fixtures.tempDir("learner")
        defer { try? FileManager.default.removeItem(at: dir) }
        let older = dir.appendingPathComponent("older.sqlite")
        try LearnerDatabase(url: older, schema: Fixtures.learnerSchemaURL).close()
        try await DatabaseQueue(path: older.path).writeWithoutTransaction { try $0.execute(sql: "PRAGMA user_version = 0") }

        #expect(throws: LearnerDatabaseError.noMigration(path: older.path, version: 0)) {
            try LearnerDatabase(url: older, schema: Fixtures.learnerSchemaURL, migrations: [:])
        }
        let migrated = try LearnerDatabase(url: older, schema: Fixtures.learnerSchemaURL, migrations: [
            0: { db in try db.execute(sql: "CREATE TABLE migrated (x INTEGER) STRICT") },
            1: { db in try db.execute(sql: "INSERT INTO migrated VALUES (1)") },
        ])
        try migrated.close()
        let after = try header(older)
        #expect(after.version == LearnerDatabase.schemaVersion && after.tables.contains("migrated"))
    }

    /// Version 1 to 2: what only the lesson types library format 8 removed used goes (project
    /// folders, server addresses, the attempts of command-line checks); every other row stays,
    /// each attempt with its id, and the file has a new database's schema.
    @Test func version1LosesWhatOnlyCommandLineLessonsUsed() async throws {
        let dir = try Fixtures.tempDir("learner")
        defer { try? FileManager.default.removeItem(at: dir) }
        let v1 = dir.appendingPathComponent("v1.sqlite")
        let queue = try DatabaseQueue(path: v1.path)
        try await queue.writeWithoutTransaction { db in
            try db.execute(sql: String(contentsOf: Fixtures.repository.appendingPathComponent("db/history/learner_v1.sql"), encoding: .utf8))
            try db.execute(sql: """
                INSERT INTO library VALUES ('lib', 'Apogee', 'builtin:apogee.sqlite', '2026-10-31T10:00:00Z');
                INSERT INTO progress VALUES ('lib', 'l1', '2026-10-31T10:05:00Z', NULL);
                INSERT INTO attempt VALUES (3, 'lib', 'l1', 'run', 0, 'boom', 120, '2026-10-31T10:01:00Z');
                INSERT INTO attempt VALUES (7, 'lib', 'l2', 'check', 1, 'every step passed', 80, '2026-10-31T10:02:00Z');
                INSERT INTO attempt VALUES (9, 'lib', 'l1', 'submit', 1, 'ok', 90, '2026-10-31T10:03:00Z');
                INSERT INTO project VALUES ('lib', 'c1', '/Users/me/shipping', 'https://my-app.example.com');
                INSERT INTO server_address VALUES ('http://localhost:8080', 'http://localhost:8181');
                INSERT INTO setting VALUES ('editor_font_size', '15');
                """)
        }
        try queue.close()

        let learner = try LearnerDatabase(url: v1, schema: Fixtures.learnerSchemaURL)
        #expect(try await learner.setting("editor_font_size") == "15")
        #expect(try await learner.done() == ["lib": ["l1"]])
        #expect(try await learner.failedAttempts(LessonRef(library: "lib", lesson: "l1"), kind: .run) == 1)
        await #expect(throws: DatabaseError.self) { // a check is no kind of attempt now
            try await learner.dbQueue.write { try $0.execute(sql: "INSERT INTO attempt(library_uuid, lesson_uuid, kind, passed, output, at) VALUES ('lib', 'l2', 'check', 1, '', '2026-11-01T10:00:00Z')") }
        }
        try learner.close()
        #expect(try header(v1).version == 2)
        let attempts = try await DatabaseQueue(path: v1.path).read { db in
            try Row.fetchAll(db, sql: "SELECT id, kind FROM attempt ORDER BY id").map { "\($0["id"] as Int) \($0["kind"] as String)" }
        }
        #expect(attempts == ["3 run", "9 submit"])

        let fresh = dir.appendingPathComponent("fresh.sqlite")
        try LearnerDatabase(url: fresh, schema: Fixtures.learnerSchemaURL).close()
        #expect(try schema(of: v1) == schema(of: fresh))
    }

    /// Every object of a database's schema with its SQL, comments left out and spaces collapsed;
    /// a table renamed into place has its name quoted, which counts for nothing.
    func schema(of url: URL) throws -> [String: String] {
        let queue = try DatabaseQueue(path: url.path, configuration: { var c = Configuration(); c.readonly = true; return c }())
        defer { try? queue.close() }
        return try queue.read { db in
            var out: [String: String] = [:]
            for row in try Row.fetchAll(db, sql: "SELECT type || ' ' || name AS object, coalesce(sql, '') AS text FROM sqlite_schema") {
                let text: String = row["text"]
                let lines = text.split(separator: "\n", omittingEmptySubsequences: false).map { $0.split(separator: "--", maxSplits: 1, omittingEmptySubsequences: false).first ?? "" }
                out[row["object"]] = lines.joined(separator: " ").split(whereSeparator: \.isWhitespace).joined(separator: " ")
                    .replacing(/^CREATE TABLE "(\w+)"/) { "CREATE TABLE \($0.1)" }
            }
            return out
        }
    }

    @Test func librariesAreRecordedOnceAndKeptUpToDate() async throws {
        let dir = try Fixtures.tempDir("learner")
        defer { try? FileManager.default.removeItem(at: dir) }
        let learner = try LearnerDatabase(url: dir.appendingPathComponent("learner.sqlite"), schema: Fixtures.learnerSchemaURL)
        let apogee = LibraryInfo(uuid: "lib-apogee", title: "Apogee")
        try await learner.registerLibrary(apogee, location: "builtin:apogee.sqlite", now: Date(timeIntervalSince1970: 1_800_000_000))
        try await learner.registerLibrary(LibraryInfo(uuid: "lib-mine", title: "Mine"), location: "/Users/me/mine.library",
                                          now: Date(timeIntervalSince1970: 1_800_000_100))
        try await learner.registerLibrary(LibraryInfo(uuid: "lib-apogee", title: "Apogee 2"), location: "builtin:apogee.sqlite",
                                          now: Date(timeIntervalSince1970: 1_800_000_200))
        let libraries = try await learner.libraries()
        #expect(libraries.map(\.title) == ["Apogee 2", "Mine"])
        #expect(libraries[0].addedAt == Clock.string(Date(timeIntervalSince1970: 1_800_000_000))) // first seen stays
        #expect(libraries[0].isBuiltIn && !libraries[1].isBuiltIn)
        // The sidebar's opened libraries: those not built in, as you added them.
        #expect(BuiltInLibraries.opened(from: libraries).map(\.uuid) == ["lib-mine"])
        #expect(BuiltInLibraries.opened(from: [
            LearnerLibrary(uuid: "b", title: "B", location: "/b", addedAt: "2027-03-02T00:00:00Z"),
            LearnerLibrary(uuid: "a", title: "A", location: "/a", addedAt: "2027-03-01T00:00:00Z"),
        ]).map(\.uuid) == ["a", "b"])
    }

    /// Everything of yours by (library, lesson): the same lesson uuid in two libraries is two
    /// lessons. The library files are never written.
    @Test func studyDataByLibraryAndLesson() async throws {
        let (_, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        let libraryFile = dir.appendingPathComponent(Fixtures.libraryName)
        let libraryBytes = try Data(contentsOf: libraryFile)
        let learner = try await Fixtures.learner(in: dir)
        let other = LibraryInfo(uuid: "b0000000-0000-0000-0000-000000000002", title: "Other")
        try await learner.registerLibrary(other, location: "/tmp/other.library")
        let code = Fixtures.ref(Fixtures.codeLesson), quiz = Fixtures.ref(Fixtures.quizLesson)
        let otherCode = LessonRef(library: other.uuid, lesson: Fixtures.codeLesson)
        func at(_ s: String) -> Date { Clock.date(s)! }

        try await learner.setDone(quiz, true, now: at("2027-03-01T10:00:00Z"))
        try await learner.setDone(quiz, true, now: at("2027-03-01T11:00:00Z")) // the first time stays
        try await learner.markOpened(code, now: at("2027-03-01T10:05:00Z"))
        try await learner.markOpened(otherCode, now: at("2027-03-01T10:06:00Z"))
        try await learner.setDone(otherCode, true)
        #expect(try await learner.done() == [Fixtures.libraryID: [Fixtures.quizLesson], other.uuid: [Fixtures.codeLesson]])
        try await learner.setDone(otherCode, false) // keeps the row: it was opened
        try await learner.setDone(Fixtures.ref(Fixtures.loopsLesson), true)
        try await learner.setDone(Fixtures.ref(Fixtures.loopsLesson), false) // a row with neither time goes
        #expect(try await learner.done() == [Fixtures.libraryID: [Fixtures.quizLesson]])
        // The last lesson you opened, across libraries.
        #expect(try await learner.lastOpened() == otherCode)
        #expect(try Fixtures.rows(in: dir, "SELECT done_at FROM progress WHERE lesson_uuid = ? AND library_uuid = ?",
                                  [Fixtures.quizLesson, Fixtures.libraryID]).first?["done_at"] == "2027-03-01T10:00:00Z")

        try await learner.addAttempt(code, kind: .run, passed: false, output: "boom", durationMs: 120, now: at("2027-03-01T10:00:00Z"))
        try await learner.addAttempt(code, kind: .submit, passed: false, output: "nope", durationMs: 90, now: at("2027-03-01T10:01:00Z"))
        try await learner.addAttempt(otherCode, kind: .submit, passed: false, output: "elsewhere", durationMs: 90)
        #expect(try await learner.failedAttempts(code, kind: .submit) == 1)
        #expect(try await learner.failedAttempts(code, kind: .run) == 1)
        let attempts = try Fixtures.rows(in: dir, "SELECT kind, output, duration_ms FROM attempt WHERE library_uuid = ? AND lesson_uuid = ? ORDER BY at DESC",
                                         [Fixtures.libraryID, Fixtures.codeLesson])
        #expect(attempts.map { $0["kind"] as String } == ["submit", "run"] && attempts[1]["duration_ms"] as Int? == 120)

        try await learner.saveNote(code, body: "Remember Println.", now: at("2027-03-01T10:00:00Z"))
        try await learner.saveNote(code, body: "Remember Println.", now: at("2027-03-01T11:00:00Z")) // unchanged: nothing written
        #expect(try Fixtures.rows(in: dir, "SELECT updated_at FROM note").first?["updated_at"] == "2027-03-01T10:00:00Z")
        #expect(try await learner.note(for: code) == "Remember Println.")
        #expect(try await learner.note(for: otherCode) == nil)
        try await learner.saveNote(code, body: "  \n")
        #expect(try await learner.note(for: code) == nil)

        try await learner.setSetting("editor_font_size", "15")
        #expect(try await learner.setting("editor_font_size") == "15")

        let v1 = CodeFileRecord(name: "main.go", content: "v1", updatedAt: "2027-03-01T10:00:00Z")
        #expect(try await learner.saveCode(code, files: [v1]))
        #expect(try await !learner.saveCode(code, files: [v1])) // unchanged
        #expect(try await learner.codeFiles(for: code)["main.go"]?.content == "v1")
        #expect(try await learner.codeFiles(for: otherCode).isEmpty)

        // The schema's CHECKs accept only UTC with a Z, and study data names a known library.
        await #expect(throws: DatabaseError.self) {
            try await learner.saveCode(code, files: [CodeFileRecord(name: "main.go", content: "v2", updatedAt: "2027-03-01T10:00:00+04:00")])
        }
        await #expect(throws: DatabaseError.self) { try await learner.saveNote(LessonRef(library: "unknown", lesson: "x"), body: "x") }

        // The course library was only read.
        #expect(try Data(contentsOf: libraryFile) == libraryBytes)
        #expect(!FileManager.default.fileExists(atPath: libraryFile.path + "-journal"))
    }

    @Test func doneMarksAreObserved() async throws {
        let dir = try Fixtures.tempDir("learner")
        defer { try? FileManager.default.removeItem(at: dir) }
        let learner = try await Fixtures.learner(in: dir)
        var iterator = learner.doneUpdates().makeAsyncIterator()
        #expect(try await iterator.next() == [:])
        try await learner.setDone(Fixtures.ref(Fixtures.codeLesson), true)
        #expect(try await iterator.next() == [Fixtures.libraryID: [Fixtures.codeLesson]])
        try await learner.setDone(Fixtures.ref(Fixtures.quizLesson), true)
        #expect(try await iterator.next() == [Fixtures.libraryID: [Fixtures.codeLesson, Fixtures.quizLesson]])
    }
}

/// The libraries built into the app: Resources/Courses/*.sqlite, by file name.
@Suite struct BuiltInLibraryTests {
    @Test func listsTheCoursesFolder() throws {
        let dir = try Fixtures.tempDir("courses")
        defer { try? FileManager.default.removeItem(at: dir) }
        for name in ["b-orbits.sqlite", "a-apogee.sqlite", "notes.txt", ".hidden.sqlite", "c.sqlite-journal"] {
            FileManager.default.createFile(atPath: dir.appendingPathComponent(name).path, contents: Data())
        }
        let files = BuiltInLibraries.files(in: dir)
        #expect(files.map(\.lastPathComponent) == ["a-apogee.sqlite", "b-orbits.sqlite"])
        #expect(BuiltInLibraries.location(of: files[0]) == "builtin:a-apogee.sqlite")
        #expect(BuiltInLibraries.files(in: nil).isEmpty)
        #expect(BuiltInLibraries.files(in: dir.appendingPathComponent("none")).isEmpty)
        // A test run has no app bundle, so no built-in courses either.
        #expect(BuiltInLibraries.files(in: BuiltInLibraries.folder).allSatisfy { $0.pathExtension == "sqlite" })
    }
}

/// Nothing is left of git sync or of the single open library: their settings go at launch.
@Suite struct LeftoverTests {
    @Test func obsoleteSettingsAreRemoved() throws {
        // A suite named by a path keeps its file there, in $TMPDIR, not in ~/Library/Preferences.
        let dir = try Fixtures.tempDir("leftovers")
        defer { try? FileManager.default.removeItem(at: dir) }
        let suite = dir.appendingPathComponent("settings").path
        let defaults = try #require(UserDefaults(suiteName: suite))
        defer { defaults.removePersistentDomain(forName: suite) }
        for key in ObsoletePreferences.keys { defaults.set("old", forKey: key) }
        defaults.set("/usr/local/go/bin/go", forKey: "GoPath")
        ObsoletePreferences.remove(from: defaults)
        #expect(ObsoletePreferences.keys.allSatisfy { defaults.object(forKey: $0) == nil })
        #expect(defaults.string(forKey: "GoPath") == "/usr/local/go/bin/go") // what still means something stays
    }
}

/// A library migrated by the tool: its learner rows sit in learner.sqlite under its uuid, with
/// its path as the location. The app lists it with the libraries you opened and shows that
/// progress as it is.
@Suite struct MigratedLibraryTests {
    @Test func migratedProgressShowsAsIs() async throws {
        let (_, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        let libraryFile = dir.appendingPathComponent(Fixtures.libraryName)
        let url = dir.appendingPathComponent("learner.sqlite")
        try LearnerDatabase(url: url, schema: Fixtures.learnerSchemaURL).close()
        // What the tool writes, row by row, as it moves the progress out of the library.
        try await DatabaseQueue(path: url.path).write { db in
            let lib = Fixtures.libraryID
            try db.execute(sql: "INSERT INTO library(uuid, title, location, added_at) VALUES (?, ?, ?, ?)",
                           arguments: [lib, Fixtures.libraryTitle, libraryFile.path, "2026-10-01T09:00:00Z"])
            try db.execute(sql: "INSERT INTO progress(library_uuid, lesson_uuid, done_at, last_opened_at) VALUES (?, ?, ?, ?)",
                           arguments: [lib, Fixtures.codeLesson, "2026-10-02T10:00:00Z", "2026-10-02T10:00:00Z"])
            try db.execute(sql: "INSERT INTO progress(library_uuid, lesson_uuid, done_at, last_opened_at) VALUES (?, ?, NULL, ?)",
                           arguments: [lib, Fixtures.quizLesson, "2026-10-03T10:00:00Z"])
            try db.execute(sql: "INSERT INTO code_file(library_uuid, lesson_uuid, name, content, updated_at) VALUES (?, ?, ?, ?, ?)",
                           arguments: [lib, Fixtures.codeLesson, "main.go", "package main // mine\n", "2026-10-02T10:00:00Z"])
            try db.execute(sql: "INSERT INTO note(library_uuid, lesson_uuid, body, updated_at) VALUES (?, ?, ?, ?)",
                           arguments: [lib, Fixtures.codeLesson, "Println adds the newline.", "2026-10-02T10:00:00Z"])
            try db.execute(sql: "INSERT INTO setting(key, value) VALUES ('editor_font_size', '15')")
        }
        let learner = try LearnerDatabase(url: url, schema: Fixtures.learnerSchemaURL)
        let opened = BuiltInLibraries.opened(from: try await learner.libraries())
        #expect(opened.map(\.uuid) == [Fixtures.libraryID] && opened[0].location == libraryFile.path)
        let library = try AppDatabase(url: URL(fileURLWithPath: opened[0].location))
        #expect(library.info.uuid == opened[0].uuid)
        #expect(try await learner.done() == [Fixtures.libraryID: [Fixtures.codeLesson]])
        #expect(try await learner.lastOpened() == Fixtures.ref(Fixtures.quizLesson))
        #expect(try await learner.codeFiles(for: Fixtures.ref(Fixtures.codeLesson))["main.go"]?.content == "package main // mine\n")
        #expect(try await learner.note(for: Fixtures.ref(Fixtures.codeLesson)) == "Println adds the newline.")
        #expect(try await learner.setting("editor_font_size") == "15")
    }
}

/// Writes nothing else waits for, such as a review answer, and quitting, which waits for them.
@Suite struct PendingWriteTests {
    /// A review answer given a moment before the app quits is there when it opens again; a
    /// write that fails holds nothing up.
    @Test @MainActor func quittingWaitsForAnAnswerJustGiven() async throws {
        let dir = try Fixtures.tempDir("pending")
        defer { try? FileManager.default.removeItem(at: dir) }
        let learner = try await Fixtures.learner(in: dir)
        let ref = Fixtures.ref("quiz"), day: TimeInterval = 86_400
        let start = Date(timeIntervalSince1970: 1_800_000_000)
        // A question answered right in its lesson: due again a day later.
        try await learner.recordAnswers(ref, [QuizState.Answer(question: 0, text: "4", right: true)], now: start)
        // With nothing pending, quitting goes on at once.
        await learner.pendingWrites.finish()

        // On day two the review asks it again and it is answered right, as the review view
        // sends answers: started, and not waited for. The write is held until the app quits.
        let gate = Gate()
        learner.pendingWrites.start {
            await gate.pass()
            try await learner.recordQuestionReview(ref, position: 0, answer: "4", correct: true, now: start + 2 * day)
        }
        learner.pendingWrites.start { throw CancellationError() }
        Task {
            await gate.somethingWaits()
            await gate.open()
        }
        // The app quits.
        await learner.pendingWrites.finish()
        try learner.close()

        // Answered right on day two, it comes back three days on, on day five.
        let again = try LearnerDatabase(url: dir.appendingPathComponent("learner.sqlite"), schema: Fixtures.learnerSchemaURL)
        defer { try? again.close() }
        #expect(try await again.dueQuestions(now: start + 3 * day).isEmpty)
        #expect(try await again.dueQuestions(now: start + 5 * day).map(\.position) == [0])
    }
}
