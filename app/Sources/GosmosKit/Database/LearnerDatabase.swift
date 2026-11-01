import Foundation
import GRDB

/// Why the learner database could not be opened.
public enum LearnerDatabaseError: Error, LocalizedError, Equatable, Sendable {
    /// The app has no learner.sql to make a new database from.
    case noSchema
    case notLearner(path: String, applicationID: Int)
    case newerFormat(path: String, version: Int)
    /// A file of an older version for which this app has no migration.
    case noMigration(path: String, version: Int)
    case unreadable(path: String, detail: String)

    public var errorDescription: String? {
        switch self {
        case .noSchema:
            "Gosmos can't make a place for your progress: learner.sql is missing from the app."
        case .notLearner(let path, let id):
            "\(path) is not a Gosmos learner database: application_id is \(id), not \(LearnerDatabase.applicationID)."
        case .newerFormat(let path, let version):
            "\(path) is a learner database of a newer format: user_version is \(version), not \(LearnerDatabase.schemaVersion)."
        case .noMigration(let path, let version):
            "\(path) is a learner database of format \(version), which this version of Gosmos can't bring up to date."
        case .unreadable(let path, let detail):
            "\(path) could not be read: \(detail)"
        }
    }

    public var recoverySuggestion: String? {
        switch self {
        case .noSchema: "Install Gosmos again."
        case .notLearner, .unreadable: "Move the file somewhere else, and Gosmos makes a new one in its place."
        case .newerFormat: "Update Gosmos to use it."
        case .noMigration: "Update Gosmos, or move the file somewhere else and start over."
        }
    }

    public var title: String { "Your Progress Can't Be Opened" }

    /// The file it is about, if any.
    public var path: String? {
        switch self {
        case .noSchema: nil
        case .notLearner(let path, _), .newerFormat(let path, _), .noMigration(let path, _), .unreadable(let path, _): path
        }
    }
}

/// Everything that is yours, in one SQLite file per Mac made from db/learner.sql: the
/// libraries you study in, progress, attempts, notes, code, the hints and solutions you looked
/// at, and your review schedule. A lesson is (library_uuid, lesson_uuid). The course libraries
/// stay read-only; nothing of yours is ever written to them.
public final class LearnerDatabase: Sendable {
    /// 'GSML'.
    public static let applicationID = 1_196_641_612
    public static let schemaVersion = 2
    /// How to bring a file up one version, by the version it is at: each new learner.sql comes
    /// with the step from the one before.
    public typealias Migration = @Sendable (Database) throws -> Void
    static let migrations: [Int: Migration] = [1: dropCommandLineData]

    /// Version 1 to 2 drops what only the lesson types library format 8 removed used: the course
    /// project folders (project), your server addresses (server_address) and the attempts of
    /// command-line checks; attempt takes the other kinds only. Every other row stays as it was,
    /// each attempt with its id. The new attempt table is learner.sql's.
    static let dropCommandLineData: Migration = { db in
        try db.execute(sql: """
            DROP TABLE project;
            DROP TABLE server_address;
            CREATE TABLE new_attempt (
              id           INTEGER PRIMARY KEY,
              library_uuid TEXT NOT NULL REFERENCES library(uuid),
              lesson_uuid  TEXT NOT NULL CHECK (lesson_uuid <> ''),
              kind         TEXT NOT NULL CHECK (kind IN ('run', 'submit', 'answer', 'review')),
              passed       INTEGER NOT NULL CHECK (passed IN (0, 1)),
              output       TEXT NOT NULL,
              duration_ms  INTEGER CHECK (duration_ms >= 0),
              at           TEXT NOT NULL CHECK (at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(at) IS NOT NULL)
            ) STRICT;
            INSERT INTO new_attempt (id, library_uuid, lesson_uuid, kind, passed, output, duration_ms, at)
              SELECT id, library_uuid, lesson_uuid, kind, passed, output, duration_ms, at FROM attempt WHERE kind <> 'check';
            DROP TABLE attempt;
            ALTER TABLE new_attempt RENAME TO attempt;
            CREATE INDEX attempt_lesson ON attempt(library_uuid, lesson_uuid, at);
            """)
    }

    let dbQueue: DatabaseQueue
    /// The writes nothing else waits for, such as a review answer; quitting waits for them.
    public let pendingWrites = PendingWrites()

    /// ~/Library/Application Support/Gosmos/learner.sqlite.
    public static var defaultURL: URL {
        URL.applicationSupportDirectory.appendingPathComponent("Gosmos/learner.sqlite", isDirectory: false)
    }

    /// The learner.sql the app carries, which makes a new database.
    public static var bundledSchema: URL? { Bundle.main.url(forResource: "learner", withExtension: "sql") }

    /// Opens the learner database at `url`, making it from `schema` (learner.sql) when the
    /// file is missing, and bringing an older one up to date. A file that is not a learner
    /// database, or one of a newer format, is refused and left as it is.
    public convenience init(url: URL, schema: URL?) throws(LearnerDatabaseError) {
        try self.init(url: url, schema: schema, migrations: LearnerDatabase.migrations)
    }

    init(url: URL, schema: URL?, migrations: [Int: Migration]) throws(LearnerDatabaseError) {
        let path = url.path
        Self.removeStaleDrafts(beside: url)
        if Self.isMissing(path) { try Self.make(url, from: schema) }
        var config = Configuration()
        config.foreignKeysEnabled = true
        config.busyMode = .timeout(10)
        config.label = "Gosmos learner"
        let queue: DatabaseQueue
        let header: (id: Int, version: Int)
        do {
            queue = try DatabaseQueue(path: path, configuration: config)
            header = try queue.read { db in
                let row = try Row.fetchOne(db, sql: "SELECT (SELECT application_id FROM pragma_application_id) AS a, (SELECT user_version FROM pragma_user_version) AS v")
                return (row?["a"] ?? 0, row?["v"] ?? 0)
            }
        } catch {
            throw .unreadable(path: path, detail: AppDatabase.describe(error))
        }
        guard header.id == Self.applicationID else { throw .notLearner(path: path, applicationID: header.id) }
        guard header.version <= Self.schemaVersion else { throw .newerFormat(path: path, version: header.version) }
        if header.version < Self.schemaVersion {
            let steps = (header.version..<Self.schemaVersion).compactMap { migrations[$0] }
            guard steps.count == Self.schemaVersion - header.version else {
                throw .noMigration(path: path, version: header.version)
            }
            do {
                try queue.write { db in
                    for step in steps { try step(db) }
                    try db.execute(sql: "PRAGMA user_version = \(Self.schemaVersion)")
                }
            } catch {
                throw .unreadable(path: path, detail: AppDatabase.describe(error))
            }
        }
        dbQueue = queue
    }

    /// No file, or an empty one (nothing was ever written to it).
    private static func isMissing(_ path: String) -> Bool {
        let size = (try? FileManager.default.attributesOfItem(atPath: path)[.size] as? NSNumber)?.intValue
        return size == nil || size == 0
    }

    /// Makes a new database from learner.sql: in a file of its own next to the place, moved
    /// into it when whole, so a crash never leaves half a database behind (the draft it leaves
    /// goes at a later opening: removeStaleDrafts). It moves in only
    /// while the place is free or holds the empty file: another copy of the app, or the
    /// library tool, may have made the database meanwhile, and that one is kept.
    static func make(_ url: URL, from schema: URL?) throws(LearnerDatabaseError) {
        guard let schema, let sql = try? String(contentsOf: schema, encoding: .utf8) else { throw .noSchema }
        let fm = FileManager.default
        let folder = url.deletingLastPathComponent()
        let draft = folder.appendingPathComponent(".\(url.lastPathComponent).\(UUID().uuidString)", isDirectory: false)
        do {
            try fm.createDirectory(at: folder, withIntermediateDirectories: true)
            let queue = try DatabaseQueue(path: draft.path)
            try queue.writeWithoutTransaction { db in try db.execute(sql: sql) }
            try queue.close()
            // Never over a file that appeared meanwhile; over the empty one, in one step.
            if renamex_np(draft.path, url.path, UInt32(RENAME_EXCL)) != 0 {
                guard errno == EEXIST else { throw POSIXError(POSIXErrorCode(rawValue: errno) ?? .EIO) }
                if isMissing(url.path) {
                    guard rename(draft.path, url.path) == 0 else { throw POSIXError(POSIXErrorCode(rawValue: errno) ?? .EIO) }
                } else {
                    try fm.removeItem(at: draft)
                }
            }
        } catch {
            try? fm.removeItem(at: draft)
            throw .unreadable(path: url.path, detail: AppDatabase.describe(error))
        }
    }

    /// Removes what a crash while making the database left beside it: drafts `make` names
    /// (".learner.sqlite.<UUID>") and SQLite's files for them ("-journal", "-wal", "-shm"),
    /// once older than `age`. Making one takes milliseconds, so a younger draft may be another
    /// copy of the app's, still at work, and stays. Only regular files named so go; nothing
    /// else in the folder is touched. Returns how many went.
    @discardableResult
    static func removeStaleDrafts(beside url: URL, olderThan age: Duration = .seconds(600), now: Date = Date()) -> Int {
        let fm = FileManager.default
        let folder = url.deletingLastPathComponent()
        guard let names = try? fm.contentsOfDirectory(atPath: folder.path) else { return 0 }
        var removed = 0
        for name in names where isDraft(name, of: url.lastPathComponent) {
            let path = folder.appendingPathComponent(name).path
            // The item itself, never what a link points to.
            guard let attributes = try? fm.attributesOfItem(atPath: path), attributes[.type] as? FileAttributeType == .typeRegular,
                  let modified = attributes[.modificationDate] as? Date,
                  now.timeIntervalSince(modified) > Double(age.components.seconds) else { continue }
            // unlink, never a folder, even one put in its place since the look above.
            if unlink(path) == 0 { removed += 1 }
        }
        return removed
    }

    /// Whether `name` is a draft `make` names for the database called `base`, or one of
    /// SQLite's files for such a draft.
    static func isDraft(_ name: String, of base: String) -> Bool {
        let prefix = "." + base + "."
        guard name.hasPrefix(prefix) else { return false }
        let rest = name.dropFirst(prefix.count)
        let id = String(["-journal", "-wal", "-shm"].first { rest.hasSuffix($0) }.map { rest.dropLast($0.count) } ?? rest)
        // UUID().uuidString, as make writes it: upper case, with its dashes.
        return UUID(uuidString: id)?.uuidString == id
    }

    public func close() throws {
        try dbQueue.close()
    }

    // MARK: Libraries

    /// The libraries you study in, as you added them.
    public func libraries() async throws -> [LearnerLibrary] {
        try await dbQueue.read { db in
            try Row.fetchAll(db, sql: "SELECT uuid, title, location, added_at FROM library ORDER BY added_at, uuid").map {
                LearnerLibrary(uuid: $0["uuid"], title: $0["title"], location: $0["location"], addedAt: $0["added_at"])
            }
        }
    }

    /// Records a library you study in, or brings its title and location up to date.
    public func registerLibrary(_ info: LibraryInfo, location: String, now: Date = Date()) async throws {
        try await dbQueue.write { db in
            try db.execute(sql: """
                INSERT INTO library(uuid, title, location, added_at) VALUES (?, ?, ?, ?)
                ON CONFLICT(uuid) DO UPDATE SET title = excluded.title, location = excluded.location
                WHERE title IS NOT excluded.title OR location IS NOT excluded.location
                """, arguments: [info.uuid, info.title, location, Clock.string(now)])
        }
    }

    // MARK: Progress

    /// The lessons you finished, by library.
    public func done() async throws -> [String: Set<String>] {
        try await dbQueue.read { db in try LearnerDatabase.fetchDone(db) }
    }

    /// The finished lessons, now and after every change.
    public func doneUpdates() -> AsyncThrowingStream<[String: Set<String>], Error> {
        let observation = ValueObservation.tracking { db in try LearnerDatabase.fetchDone(db) }.removeDuplicates()
        let queue = dbQueue
        return AsyncThrowingStream { continuation in
            let cancellable = observation.start(in: queue, scheduling: .async(onQueue: .main),
                                                onError: { continuation.finish(throwing: $0) },
                                                onChange: { continuation.yield($0) })
            continuation.onTermination = { _ in cancellable.cancel() }
        }
    }

    static func fetchDone(_ db: Database) throws -> [String: Set<String>] {
        var out: [String: Set<String>] = [:]
        for row in try Row.fetchAll(db, sql: "SELECT library_uuid, lesson_uuid FROM progress WHERE done_at IS NOT NULL") {
            out[row["library_uuid"], default: []].insert(row["lesson_uuid"])
        }
        return out
    }

    public func setDone(_ ref: LessonRef, _ done: Bool, now: Date = Date()) async throws {
        try await dbQueue.write { db in try LearnerDatabase.writeDone(db, ref, done, at: Clock.string(now)) }
    }

    static func writeDone(_ db: Database, _ ref: LessonRef, _ done: Bool, at: String) throws {
        if done {
            try db.execute(sql: """
                INSERT INTO progress(library_uuid, lesson_uuid, done_at) VALUES (?, ?, ?)
                ON CONFLICT(library_uuid, lesson_uuid) DO UPDATE SET done_at = coalesce(progress.done_at, excluded.done_at)
                """, arguments: [ref.library, ref.lesson, at])
        } else {
            // A progress row needs done_at or last_opened_at; one that would have neither goes.
            try db.execute(sql: "DELETE FROM progress WHERE library_uuid = ? AND lesson_uuid = ? AND last_opened_at IS NULL",
                           arguments: [ref.library, ref.lesson])
            try db.execute(sql: "UPDATE progress SET done_at = NULL WHERE library_uuid = ? AND lesson_uuid = ?",
                           arguments: [ref.library, ref.lesson])
        }
    }

    public func markOpened(_ ref: LessonRef, now: Date = Date()) async throws {
        try await dbQueue.write { db in
            try db.execute(sql: """
                INSERT INTO progress(library_uuid, lesson_uuid, last_opened_at) VALUES (?, ?, ?)
                ON CONFLICT(library_uuid, lesson_uuid) DO UPDATE SET last_opened_at = excluded.last_opened_at
                """, arguments: [ref.library, ref.lesson, Clock.string(now)])
        }
    }

    /// The lesson you opened last, in any library.
    public func lastOpened() async throws -> LessonRef? {
        try await dbQueue.read { db in
            try Row.fetchOne(db, sql: """
                SELECT library_uuid, lesson_uuid FROM progress WHERE last_opened_at IS NOT NULL
                ORDER BY last_opened_at DESC LIMIT 1
                """).map { LessonRef(library: $0["library_uuid"], lesson: $0["lesson_uuid"]) }
        }
    }

    // MARK: Attempts

    public func addAttempt(_ ref: LessonRef, kind: AttemptKind, passed: Bool, output: String, durationMs: Int?, now: Date = Date()) async throws {
        try await dbQueue.write { db in
            try LearnerDatabase.insertAttempt(db, ref, kind: kind, passed: passed, output: output, durationMs: durationMs, at: Clock.string(now))
        }
    }

    static func insertAttempt(_ db: Database, _ ref: LessonRef, kind: AttemptKind, passed: Bool, output: String, durationMs: Int?, at: String) throws {
        try db.execute(sql: """
            INSERT INTO attempt(library_uuid, lesson_uuid, kind, passed, output, duration_ms, at) VALUES (?, ?, ?, ?, ?, ?, ?)
            """, arguments: [ref.library, ref.lesson, kind.rawValue, passed ? 1 : 0, output, durationMs.map { max(0, $0) }, at])
    }

    /// Failed attempts of one kind, such as Submits.
    public func failedAttempts(_ ref: LessonRef, kind: AttemptKind) async throws -> Int {
        try await dbQueue.read { db in try LearnerDatabase.countFailed(db, ref, kind: kind) }
    }

    static func countFailed(_ db: Database, _ ref: LessonRef, kind: AttemptKind) throws -> Int {
        try Int.fetchOne(db, sql: "SELECT count(*) FROM attempt WHERE library_uuid = ? AND lesson_uuid = ? AND kind = ? AND passed = 0",
                         arguments: [ref.library, ref.lesson, kind.rawValue]) ?? 0
    }

    /// Records a Submit. A pass finishes the lesson. For a lesson you can redo (`reviewable`: a
    /// code lesson), a pass after a look at its solution, or after three failed Submits or more,
    /// also puts it in review, due in a day; returns whether it did.
    @discardableResult
    public func recordSubmit(_ ref: LessonRef, passed: Bool, output: String, durationMs: Int?,
                             reviewable: Bool, now: Date = Date()) async throws -> Bool {
        try await dbQueue.write { db in
            let at = Clock.string(now)
            try LearnerDatabase.insertAttempt(db, ref, kind: .submit, passed: passed, output: output, durationMs: durationMs, at: at)
            guard passed else { return false }
            try LearnerDatabase.writeDone(db, ref, true, at: at)
            guard reviewable else { return false }
            let viewed = try Bool.fetchOne(db, sql: "SELECT EXISTS (SELECT 1 FROM solution_view WHERE library_uuid = ? AND lesson_uuid = ?)",
                                           arguments: [ref.library, ref.lesson]) ?? false
            let failed = try LearnerDatabase.countFailed(db, ref, kind: .submit)
            guard SolutionPolicy.passNeedsReview(viewedSolution: viewed, failedSubmits: failed) else { return false }
            try db.execute(sql: """
                INSERT INTO review_lesson(library_uuid, lesson_uuid, due_at, interval_days) VALUES (?, ?, ?, ?)
                ON CONFLICT(library_uuid, lesson_uuid) DO NOTHING
                """, arguments: [ref.library, ref.lesson, Clock.string(ReviewSchedule.due(in: ReviewSchedule.first, from: now)), ReviewSchedule.first])
            return db.changesCount > 0
        }
    }

    /// Records a quiz Submit's graded answers. A question answered right for the first time
    /// goes into review, due in a day.
    public func recordAnswers(_ ref: LessonRef, _ answers: [QuizState.Answer], now: Date = Date()) async throws {
        guard !answers.isEmpty else { return }
        try await dbQueue.write { db in
            let at = Clock.string(now), due = Clock.string(ReviewSchedule.due(in: ReviewSchedule.first, from: now))
            for a in answers {
                try LearnerDatabase.insertAttempt(db, ref, kind: .answer, passed: a.right, output: a.text, durationMs: nil, at: at)
                if a.right {
                    try db.execute(sql: """
                        INSERT INTO review_question(library_uuid, lesson_uuid, question_position, due_at, interval_days) VALUES (?, ?, ?, ?, ?)
                        ON CONFLICT(library_uuid, lesson_uuid, question_position) DO NOTHING
                        """, arguments: [ref.library, ref.lesson, a.question, due, ReviewSchedule.first])
                }
            }
        }
    }

    // MARK: Notes, code, settings

    public func note(for ref: LessonRef) async throws -> String? {
        try await dbQueue.read { db in
            try String.fetchOne(db, sql: "SELECT body FROM note WHERE library_uuid = ? AND lesson_uuid = ?", arguments: [ref.library, ref.lesson])
        }
    }

    /// Saves a note; an empty one is removed. Saving the same text again changes nothing.
    public func saveNote(_ ref: LessonRef, body: String, now: Date = Date()) async throws {
        try await dbQueue.write { db in
            if body.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty {
                try db.execute(sql: "DELETE FROM note WHERE library_uuid = ? AND lesson_uuid = ?", arguments: [ref.library, ref.lesson])
            } else {
                try db.execute(sql: """
                    INSERT INTO note(library_uuid, lesson_uuid, body, updated_at) VALUES (?, ?, ?, ?)
                    ON CONFLICT(library_uuid, lesson_uuid) DO UPDATE SET body = excluded.body, updated_at = excluded.updated_at
                    WHERE note.body IS NOT excluded.body
                    """, arguments: [ref.library, ref.lesson, body, Clock.string(now)])
            }
        }
    }

    public func codeFiles(for ref: LessonRef) async throws -> [String: CodeFileRecord] {
        try await dbQueue.read { db in
            var out: [String: CodeFileRecord] = [:]
            for row in try Row.fetchAll(db, sql: "SELECT name, content, updated_at FROM code_file WHERE library_uuid = ? AND lesson_uuid = ?",
                                        arguments: [ref.library, ref.lesson]) {
                let r = CodeFileRecord(name: row["name"], content: row["content"], updatedAt: row["updated_at"])
                out[r.name] = r
            }
            return out
        }
    }

    /// Stores your files; rows that already hold the same content are left alone. Returns
    /// whether anything changed.
    @discardableResult
    public func saveCode(_ ref: LessonRef, files: [CodeFileRecord]) async throws -> Bool {
        try await dbQueue.write { db in
            var changed = false
            for f in files {
                try db.execute(sql: """
                    INSERT INTO code_file(library_uuid, lesson_uuid, name, content, updated_at) VALUES (?, ?, ?, ?, ?)
                    ON CONFLICT(library_uuid, lesson_uuid, name) DO UPDATE SET content = excluded.content, updated_at = excluded.updated_at
                    WHERE code_file.content IS NOT excluded.content
                    """, arguments: [ref.library, ref.lesson, f.name, f.content, f.updatedAt])
                changed = changed || db.changesCount > 0
            }
            return changed
        }
    }

    public func setting(_ key: String) async throws -> String? {
        try await dbQueue.read { db in try String.fetchOne(db, sql: "SELECT value FROM setting WHERE key = ?", arguments: [key]) }
    }

    public func setSetting(_ key: String, _ value: String?) async throws {
        try await dbQueue.write { db in
            if let value {
                try db.execute(sql: "INSERT INTO setting(key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value",
                               arguments: [key, value])
            } else {
                try db.execute(sql: "DELETE FROM setting WHERE key = ?", arguments: [key])
            }
        }
    }

    // MARK: Hints and solutions

    /// The positions of the hints you revealed.
    public func revealedHints(_ ref: LessonRef) async throws -> Set<Int> {
        try await dbQueue.read { db in
            Set(try Int.fetchAll(db, sql: "SELECT position FROM hint_view WHERE library_uuid = ? AND lesson_uuid = ?", arguments: [ref.library, ref.lesson]))
        }
    }

    public func revealHint(_ ref: LessonRef, position: Int, now: Date = Date()) async throws {
        try await dbQueue.write { db in
            try db.execute(sql: """
                INSERT INTO hint_view(library_uuid, lesson_uuid, position, at) VALUES (?, ?, ?, ?)
                ON CONFLICT(library_uuid, lesson_uuid, position) DO NOTHING
                """, arguments: [ref.library, ref.lesson, position, Clock.string(now)])
        }
    }

    /// Whether you looked at the lesson's solution before passing it.
    public func viewedSolution(_ ref: LessonRef) async throws -> Bool {
        try await dbQueue.read { db in
            try Bool.fetchOne(db, sql: "SELECT EXISTS (SELECT 1 FROM solution_view WHERE library_uuid = ? AND lesson_uuid = ?)",
                              arguments: [ref.library, ref.lesson]) ?? false
        }
    }

    /// Records a look at the lesson's solution, unless the lesson is done: after a pass,
    /// comparing is free. Returns whether it was recorded now.
    @discardableResult
    public func recordSolutionView(_ ref: LessonRef, now: Date = Date()) async throws -> Bool {
        try await dbQueue.write { db in
            let done = try Bool.fetchOne(db, sql: "SELECT EXISTS (SELECT 1 FROM progress WHERE library_uuid = ? AND lesson_uuid = ? AND done_at IS NOT NULL)",
                                         arguments: [ref.library, ref.lesson]) ?? false
            guard SolutionPolicy.recordsView(lessonDone: done) else { return false }
            try db.execute(sql: """
                INSERT INTO solution_view(library_uuid, lesson_uuid, at) VALUES (?, ?, ?)
                ON CONFLICT(library_uuid, lesson_uuid) DO NOTHING
                """, arguments: [ref.library, ref.lesson, Clock.string(now)])
            return db.changesCount > 0
        }
    }

    // MARK: Review

    /// The quiz questions due by `now`, interleaved across lessons.
    public func dueQuestions(now: Date = Date()) async throws -> [DueQuestion] {
        try await dbQueue.read { db in
            ReviewQueue.interleave(try Row.fetchAll(db, sql: """
                SELECT library_uuid, lesson_uuid, question_position, due_at FROM review_question WHERE due_at <= ?
                """, arguments: [Clock.string(now)]).map {
                DueQuestion(ref: LessonRef(library: $0["library_uuid"], lesson: $0["lesson_uuid"]), position: $0["question_position"], dueAt: $0["due_at"])
            })
        }
    }

    /// The lessons due to redo by `now`, longest overdue first.
    public func dueLessons(now: Date = Date()) async throws -> [DueLesson] {
        try await dbQueue.read { db in
            try Row.fetchAll(db, sql: """
                SELECT library_uuid, lesson_uuid, due_at FROM review_lesson WHERE due_at <= ? ORDER BY due_at, library_uuid, lesson_uuid
                """, arguments: [Clock.string(now)]).map {
                DueLesson(ref: LessonRef(library: $0["library_uuid"], lesson: $0["lesson_uuid"]), dueAt: $0["due_at"])
            }
        }
    }

    /// Records a review of a quiz question: right moves it to its next interval, wrong is a
    /// lapse and starts it over at a day.
    public func recordQuestionReview(_ ref: LessonRef, position: Int, answer: String, correct: Bool, now: Date = Date()) async throws {
        try await dbQueue.write { db in
            try LearnerDatabase.insertAttempt(db, ref, kind: .review, passed: correct, output: answer, durationMs: nil, at: Clock.string(now))
            try LearnerDatabase.reschedule(db, table: "review_question", keys: ["library_uuid", "lesson_uuid", "question_position"],
                                           values: [ref.library, ref.lesson, position], passed: correct, now: now)
        }
    }

    /// Records a redo of a lesson due for review: a passing Submit moves it to its next
    /// interval; a failing one is a lapse and starts it over at a day.
    public func recordLessonReview(_ ref: LessonRef, passed: Bool, output: String, durationMs: Int?, now: Date = Date()) async throws {
        try await dbQueue.write { db in
            try LearnerDatabase.insertAttempt(db, ref, kind: .review, passed: passed, output: output, durationMs: durationMs, at: Clock.string(now))
            try LearnerDatabase.reschedule(db, table: "review_lesson", keys: ["library_uuid", "lesson_uuid"],
                                           values: [ref.library, ref.lesson], passed: passed, now: now)
        }
    }

    private static func reschedule(_ db: Database, table: String, keys: [String], values: [any DatabaseValueConvertible],
                                   passed: Bool, now: Date) throws {
        let match = keys.map { "\($0) = ?" }.joined(separator: " AND ")
        guard let interval = try Int.fetchOne(db, sql: "SELECT interval_days FROM \(table) WHERE \(match)", arguments: StatementArguments(values)) else { return }
        let next = ReviewSchedule.next(after: interval, passed: passed)
        try db.execute(sql: """
            UPDATE \(table) SET interval_days = ?, due_at = ?, last_reviewed_at = ?, reps = reps + ?, lapses = lapses + ? WHERE \(match)
            """, arguments: StatementArguments([next, Clock.string(ReviewSchedule.due(in: next, from: now)), Clock.string(now),
                                                passed ? 1 : 0, passed ? 0 : 1] as [any DatabaseValueConvertible] + values))
    }
}

// MARK: - Study records

/// What an attempt row records: a Run, a Submit, a quiz answer, or a review.
public enum AttemptKind: String, Sendable, CaseIterable {
    case run, submit, answer, review
}

public struct CodeFileRecord: Hashable, Sendable {
    public let name: String
    public let content: String
    public let updatedAt: String

    public init(name: String, content: String, updatedAt: String) {
        self.name = name
        self.content = content
        self.updatedAt = updatedAt
    }
}
