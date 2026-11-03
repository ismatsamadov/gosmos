import Foundation
import GRDB

/// Why a library could not be opened or used. The checks and their words are those of the
/// library tool's validator (`go run ./cmd/library check <file>`), so the app and the tool say
/// the same thing about a file.
public enum AppDatabaseError: Error, LocalizedError, Equatable, Sendable {
    case missing(path: String)
    /// macOS refused access: the file sits in a folder it protects (Documents, Desktop, an
    /// external drive) and the learner said no, or file permissions forbid reading it.
    case notPermitted(path: String)
    case unreadable(path: String, detail: String)
    case notLibrary(path: String, applicationID: Int)
    case olderFormat(path: String, version: Int)
    case newerFormat(path: String, version: Int)
    /// quick_check's findings, one line each.
    case damaged(path: String, problems: [String])
    /// The foreign key, per-type row, quiz answer, solution and library_info checks' findings.
    case failedChecks(path: String, problems: [String])
    /// A write to the file was cut short; only a writer, the library tool, can roll it back.
    case interruptedWrite(path: String)
    /// Another program held the file locked for longer than the open waits.
    case busy(path: String)
    case sqlite(String)

    public var errorDescription: String? {
        switch self {
        case .missing(let path): "There is no file at \(path)."
        case .notPermitted, .unreadable, .notLibrary, .olderFormat, .newerFormat, .damaged, .failedChecks, .interruptedWrite, .busy: problems.joined(separator: "\n")
        case .sqlite(let message): message
        }
    }

    /// What is wrong, in the validator's words, one problem per line.
    public var problems: [String] {
        let id = AppDatabase.applicationID, version = AppDatabase.schemaVersion
        return switch self {
        case .missing(let path): ["There is no file at \(path)."]
        case .notPermitted(let path): ["Gosmos isn't allowed to read \(path)."]
        case .unreadable(_, let detail): ["the file is not an SQLite database SQLite can read: \(detail)"]
        case .notLibrary(_, let found): ["application_id is \(found), not \(id): the file is not a Gosmos library"]
        case .olderFormat(_, let found): ["user_version is \(found), not \(version): the library is of an older format; migrating brings it up to date"]
        case .newerFormat(_, let found): ["user_version is \(found), not \(version): the library is of a newer format than this program knows"]
        case .damaged(_, let lines), .failedChecks(_, let lines): lines
        case .interruptedWrite: ["a write to the library was cut short and has to be rolled back first"]
        case .busy: ["another program is writing to this library; try again"]
        case .sqlite(let message): [message]
        }
    }

    /// A headline for the error.
    public var title: String {
        switch self {
        case .missing: "Library Not Found"
        case .notPermitted: "No Access to the Library"
        case .unreadable: "Not an SQLite File"
        case .notLibrary: "Not a Gosmos Library"
        case .olderFormat: "Library Needs Migrating"
        case .newerFormat: "Library Too New"
        case .damaged: "Damaged Library"
        case .failedChecks: "Library Failed Its Checks"
        case .interruptedWrite: "Library Needs Recovery"
        case .busy: "Library Busy"
        case .sqlite: "Library Error"
        }
    }

    public var recoverySuggestion: String? {
        switch self {
        case .missing: "Open another library, or put the file back."
        case .notPermitted: "Allow Gosmos in System Settings › Privacy & Security › Files and Folders, then try again."
        case .unreadable, .notLibrary: "Open a Gosmos library."
        case .olderFormat(let path, _): "Run go run ./cmd/library migrate \(path), then open it again."
        case .newerFormat: "Update Gosmos to open it."
        case .damaged: "Restore the file from a backup, or open another library."
        case .failedChecks(let path, _): "Whoever made the library can list the same problems with go run ./cmd/library check \(path)."
        case .interruptedWrite(let path): "Run go run ./cmd/library check \(path), which rolls the write back, then open it again."
        case .busy: "Try again when the other program is done with the library."
        case .sqlite: nil
        }
    }
}

/// A course library, read-only, through one GRDB DatabaseQueue.
///
/// It opens a file only when it passes the format's checks (docs/library-format.md), and it
/// never writes to it: the connection is read-only, and a library built into the app is opened
/// immutable, so SQLite does not even lock it. Everything of yours lives in the learner
/// database instead.
public final class AppDatabase: Sendable {
    public static let schemaVersion = 9
    /// 'GSMS'.
    public static let applicationID = 1_196_641_619

    let dbQueue: DatabaseQueue
    /// The library's identity, from its library_info row.
    public let info: LibraryInfo

    /// How long opening waits for another program's lock, as the library tool does.
    static let busyTimeout: TimeInterval = 10

    /// Opens a library read-only after its checks. `immutable` is for a file nothing changes
    /// while the app runs, one built into the app: SQLite then reads it without any locking.
    public convenience init(url: URL, immutable: Bool = false) throws(AppDatabaseError) {
        try self.init(url: url, immutable: immutable, busyTimeout: AppDatabase.busyTimeout)
    }

    init(url: URL, immutable: Bool, busyTimeout: TimeInterval) throws(AppDatabaseError) {
        let path = url.path
        if let refused = AppDatabase.access(url) { throw refused }
        var config = Configuration()
        config.readonly = true
        config.foreignKeysEnabled = false // nothing is written, so there is nothing to enforce
        config.busyMode = .timeout(busyTimeout)
        config.label = "Gosmos library"
        let queue: DatabaseQueue
        do {
            // A URI, so SQLite takes the immutable flag; the system SQLite reads URIs.
            queue = try DatabaseQueue(path: immutable ? URL(fileURLWithPath: path).absoluteString + "?immutable=1" : path, configuration: config)
        } catch {
            throw AppDatabase.refusal(error, path: path)
        }
        // Outside a transaction: on a damaged file, ending one would fail too and hide what
        // the checks found.
        var outcome: Result<LibraryInfo, AppDatabaseError> = .failure(.sqlite("the library's checks did not run"))
        queue.inDatabase { db in
            outcome = AppDatabase.validate(db, path: path).flatMap { AppDatabase.readInfo(db, path: path) }
        }
        switch outcome {
        case .success(let info):
            self.info = info
            dbQueue = queue
        case .failure(let error):
            try? queue.close()
            throw error
        }
    }

    /// Opens the file for reading and closes it again. In a folder macOS protects, this is where
    /// macOS asks the learner for access; the call waits for the answer, however long it takes.
    /// Nil when the file can be read.
    public static func access(_ url: URL) -> AppDatabaseError? {
        let fd = Darwin.open(url.path, O_RDONLY)
        if fd >= 0 {
            Darwin.close(fd)
            return nil
        }
        switch errno {
        case ENOENT, ENOTDIR: return .missing(path: url.path)
        case EPERM, EACCES: return .notPermitted(path: url.path)
        default: return nil // the checks that follow report it in their words
        }
    }

    /// Why SQLite would not read the file: another program's lock, a write cut short whose
    /// hot journal only a writer can roll back, or the file itself.
    static func refusal(_ error: any Error, path: String) -> AppDatabaseError {
        if isBusy(error) { return .busy(path: path) }
        if (error as? DatabaseError)?.extendedResultCode == .SQLITE_READONLY_ROLLBACK { return .interruptedWrite(path: path) }
        return .unreadable(path: path, detail: describe(error))
    }

    static func isBusy(_ error: any Error) -> Bool {
        guard let e = error as? DatabaseError else { return false }
        return e.resultCode == .SQLITE_BUSY || e.resultCode == .SQLITE_LOCKED
    }

    static func readInfo(_ db: Database, path: String) -> Result<LibraryInfo, AppDatabaseError> {
        do {
            guard let row = try Row.fetchOne(db, sql: "SELECT uuid, title FROM library_info") else {
                return .failure(.failedChecks(path: path, problems: ["library_info holds 0 rows, not exactly 1"]))
            }
            return .success(LibraryInfo(uuid: row["uuid"], title: row["title"]))
        } catch {
            return .failure(.sqlite("the library's checks could not run: \(describe(error))"))
        }
    }

    /// The library tool's validator (db.Validate), check for check and word for word. After a
    /// wrong header or a damaged file it stops: the other checks would only fail to read it.
    static func validate(_ db: Database, path: String) -> Result<Void, AppDatabaseError> {
        let id: Int, version: Int
        do {
            let row = try Row.fetchOne(db, sql: "SELECT (SELECT application_id FROM pragma_application_id) AS a, (SELECT user_version FROM pragma_user_version) AS v")
            id = row?["a"] ?? 0
            version = row?["v"] ?? 0
        } catch {
            return .failure(refusal(error, path: path))
        }
        if id != applicationID { return .failure(.notLibrary(path: path, applicationID: id)) }
        if version < schemaVersion { return .failure(.olderFormat(path: path, version: version)) }
        if version > schemaVersion { return .failure(.newerFormat(path: path, version: version)) }

        var damage: [String] = []
        do {
            // A cursor, so the lines found before a failure are kept.
            let messages = try String.fetchCursor(db, sql: "SELECT quick_check FROM pragma_quick_check")
            while let message = try messages.next() {
                guard message != "ok" else { continue }
                for line in message.split(separator: "\n").map({ $0.trimmingCharacters(in: .whitespaces) })
                where !line.isEmpty && !line.hasPrefix("*** in database") {
                    damage.append("quick_check: \(line)")
                }
            }
        } catch {
            // A lock is not damage: waiting ran out, so the file is busy, not broken.
            if isBusy(error) { return .failure(.busy(path: path)) }
            damage.append("quick_check: \(describe(error))") // too damaged for the check to finish
        }
        if !damage.isEmpty { return .failure(.damaged(path: path, problems: damage)) }

        var problems: [String] = []
        do {
            for row in try Row.fetchAll(db, sql: #"SELECT "table" AS t, rowid AS r, parent AS p FROM pragma_foreign_key_check ORDER BY 1, 2"#) {
                problems.append("foreign key: \(row["t"] as String) row \(row["r"] as Int64? ?? 0) points at a missing \(row["p"] as String) row")
            }
            let typeTables: [(table: String, types: [String])] = [
                ("code_lesson", ["type_code", "type_code_tests"]), ("quiz_lesson", ["type_choice"]),
            ]
            let counts = typeTables.enumerated().map { "(SELECT count(*) FROM \($1.table) x WHERE x.lesson_uuid = l.uuid) AS n\($0)" }
            for row in try Row.fetchAll(db, sql: "SELECT l.uuid AS uuid, l.type AS type, \(counts.joined(separator: ", ")) FROM lesson l ORDER BY l.uuid") {
                let uuid: String = row["uuid"], type: String = row["type"]
                for (i, t) in typeTables.enumerated() {
                    let want = t.types.contains(type) ? 1 : 0
                    let n: Int = row["n\(i)"]
                    if n != want { problems.append("lesson \(uuid) (\(type)) has \(n) \(t.table) rows, want \(want)") }
                }
                // Only a write past lesson.type's CHECK leaves one, and quick_check lets it pass.
                if !typeTables.contains(where: { $0.types.contains(type) }) {
                    problems.append("lesson \(uuid) (\(type)) is of a type format \(schemaVersion) does not have")
                }
            }
            for row in try Row.fetchAll(db, sql: """
                SELECT lesson_uuid, position, n FROM (
                    SELECT q.lesson_uuid, q.position, (SELECT count(*) FROM quiz_answer a
                        WHERE a.lesson_uuid = q.lesson_uuid AND a.question_position = q.position AND a.is_correct = 1) AS n
                    FROM quiz_question q)
                WHERE n <> 1 ORDER BY 1, 2
                """) {
                problems.append("quiz question \(row["position"] as Int) of lesson \(row["lesson_uuid"] as String) has \(row["n"] as Int) right answers, want 1")
            }
            for row in try Row.fetchAll(db, sql: """
                SELECT s.lesson_uuid, s.name FROM solution_file s
                WHERE NOT EXISTS (SELECT 1 FROM starter_file f
                    WHERE f.lesson_uuid = s.lesson_uuid AND f.name = s.name AND f.is_hidden = 0 AND f.is_readonly = 0)
                ORDER BY 1, 2
                """) {
                problems.append("solution_file \(row["name"] as String) of lesson \(row["lesson_uuid"] as String) solves no editable starter file")
            }
            let infoRows = try Int.fetchOne(db, sql: "SELECT count(*) FROM library_info") ?? 0
            if infoRows != 1 { problems.append("library_info holds \(infoRows) rows, not exactly 1") }
        } catch {
            if isBusy(error) { return .failure(.busy(path: path)) }
            return .failure(.sqlite("the library's checks could not run: \(describe(error))"))
        }
        return problems.isEmpty ? .success(()) : .failure(.failedChecks(path: path, problems: problems))
    }

    /// SQLite's message and result code, as the library tool shows them: "file is not a
    /// database (26)".
    static func describe(_ error: any Error) -> String {
        guard let e = error as? DatabaseError else { return error.localizedDescription }
        return "\(e.message ?? "SQLite error") (\(e.extendedResultCode.rawValue))"
    }

    public func close() throws {
        try dbQueue.close()
    }

    // MARK: Reads

    public func catalog() async throws -> Catalog {
        try await dbQueue.read { db in try AppDatabase.fetchCatalog(db) }
    }

    public func lesson(_ uuid: String) async throws -> LessonDetail? {
        let library = info.uuid
        return try await dbQueue.read { db in try AppDatabase.fetchLesson(db, uuid, library: library) }
    }

    /// The quiz of a lesson, for review: its questions with their answers and explanations.
    public func quiz(_ uuid: String) async throws -> Quiz? {
        try await dbQueue.read { db in
            try Bool.fetchOne(db, sql: "SELECT EXISTS (SELECT 1 FROM quiz_lesson WHERE lesson_uuid = ?)", arguments: [uuid]) == true
                ? try AppDatabase.fetchQuiz(db, uuid) : nil
        }
    }

    public func search(_ text: String, limit: Int = 60) async throws -> [SearchHit] {
        let library = info.uuid
        return try await dbQueue.read { db in try AppDatabase.fetchSearch(db, text, limit: limit, library: library) }
    }

    /// The stories of the library's courses, by course uuid; a course without one is left out.
    public func stories() async throws -> [String: CourseStory] {
        try await dbQueue.read { db in try AppDatabase.fetchStories(db) }
    }

    /// Lessons' titles and readmes, for search excerpts.
    public func readmes(_ ids: [String]? = nil) async throws -> [(id: String, title: String, readme: String)] {
        try await dbQueue.read { db in
            let rows: [Row]
            if let ids {
                guard !ids.isEmpty else { return [] }
                rows = try Row.fetchAll(db, sql: "SELECT uuid, title, readme FROM lesson WHERE uuid IN (\(ids.map { _ in "?" }.joined(separator: ",")))",
                                        arguments: StatementArguments(ids))
            } else {
                rows = try Row.fetchAll(db, sql: "SELECT uuid, title, readme FROM lesson")
            }
            return rows.map { (id: $0["uuid"], title: $0["title"], readme: $0["readme"]) }
        }
    }

    // MARK: Queries shared by the async methods and the tests

    static func fetchCatalog(_ db: Database) throws -> Catalog {
        var courses = try Row.fetchAll(db, sql: "SELECT uuid, title, kind, short_description FROM course ORDER BY position").map {
            Course(id: $0["uuid"], title: $0["title"], kind: $0["kind"], shortDescription: $0["short_description"])
        }
        var lessonsByChapter: [String: [LessonSummary]] = [:]
        for row in try Row.fetchAll(db, sql: "SELECT uuid, chapter_uuid, title, type, required FROM lesson ORDER BY chapter_uuid, position") {
            guard let type = LessonType(rawValue: row["type"]) else { continue }
            let l = LessonSummary(id: row["uuid"], chapterID: row["chapter_uuid"], title: row["title"], type: type, required: row["required"] as Int != 0)
            lessonsByChapter[l.chapterID, default: []].append(l)
        }
        var chaptersByCourse: [String: [Chapter]] = [:]
        for row in try Row.fetchAll(db, sql: "SELECT uuid, course_uuid, title FROM chapter ORDER BY course_uuid, position") {
            var ch = Chapter(id: row["uuid"], courseID: row["course_uuid"], title: row["title"])
            ch.lessons = lessonsByChapter[ch.id] ?? []
            chaptersByCourse[ch.courseID, default: []].append(ch)
        }
        for i in courses.indices { courses[i].chapters = chaptersByCourse[courses[i].id] ?? [] }
        return Catalog(courses: courses)
    }

    static func fetchStories(_ db: Database) throws -> [String: CourseStory] {
        var stories: [String: CourseStory] = [:]
        func story(_ course: String) -> CourseStory { stories[course] ?? CourseStory(courseID: course) }
        var cards: [String: [CourseStory.Card]] = [:] // by "<course>/<page position>"
        for row in try Row.fetchAll(db, sql: "SELECT course_uuid, page_position, position, label, title, subtitle, body FROM story_card ORDER BY course_uuid, page_position, position") {
            let key = "\(row["course_uuid"] as String)/\(row["page_position"] as Int)"
            cards[key, default: []].append(CourseStory.Card(position: row["position"], label: row["label"], title: row["title"],
                                                            subtitle: row["subtitle"], body: row["body"]))
        }
        for row in try Row.fetchAll(db, sql: "SELECT course_uuid, position, kind, eyebrow, title, body, action FROM story_page ORDER BY course_uuid, position") {
            let course: String = row["course_uuid"], position: Int = row["position"]
            guard let kind = CourseStory.PageKind(rawValue: row["kind"]) else { continue }
            var s = story(course)
            s.pages.append(CourseStory.Page(position: position, kind: kind, eyebrow: row["eyebrow"], title: row["title"], body: row["body"],
                                            action: row["action"], cards: cards["\(course)/\(position)"] ?? []))
            stories[course] = s
        }
        for row in try Row.fetchAll(db, sql: "SELECT ch.course_uuid, s.chapter_uuid, s.briefing, s.debrief FROM chapter_story s JOIN chapter ch ON ch.uuid = s.chapter_uuid") {
            let course: String = row["course_uuid"]
            var s = story(course)
            s.chapters[row["chapter_uuid"]] = CourseStory.ChapterStory(briefing: row["briefing"], debrief: row["debrief"])
            stories[course] = s
        }
        for row in try Row.fetchAll(db, sql: "SELECT course_uuid, term, meaning, linked FROM story_term ORDER BY course_uuid, lower(term), term") {
            let course: String = row["course_uuid"]
            var s = story(course)
            s.terms.append(CourseStory.Term(term: row["term"], meaning: row["meaning"], linked: row["linked"] as Int != 0))
            stories[course] = s
        }
        return stories
    }

    static func fetchSearch(_ db: Database, _ text: String, limit: Int, library: String) throws -> [SearchHit] {
        guard let match = ftsQuery(text) else { return [] }
        return try Row.fetchAll(db, sql: """
            SELECT l.uuid, l.title, l.type, c.title AS course, ch.title AS chapter
            FROM lesson_search
            JOIN lesson l ON l.id = lesson_search.rowid
            JOIN chapter ch ON ch.uuid = l.chapter_uuid
            JOIN course c ON c.uuid = ch.course_uuid
            WHERE lesson_search MATCH ?
            ORDER BY bm25(lesson_search, 5.0, 1.0), c.position, ch.position, l.position
            LIMIT ?
            """, arguments: [match, limit]).compactMap { row in
            guard let type = LessonType(rawValue: row["type"]) else { return nil }
            return SearchHit(ref: LessonRef(library: library, lesson: row["uuid"]), title: row["title"], courseTitle: row["course"],
                             chapterTitle: row["chapter"], type: type)
        }
    }

    /// Turns what someone typed into an FTS5 query: every word must appear, the last one as
    /// a prefix. Quotes keep FTS5 syntax characters literal.
    static func ftsQuery(_ text: String) -> String? {
        let words = text.split(whereSeparator: { !$0.isLetter && !$0.isNumber && $0 != "_" }).map(String.init)
        guard !words.isEmpty else { return nil }
        return words.enumerated().map { i, w in
            "\"" + w.replacingOccurrences(of: "\"", with: "\"\"") + "\"" + (i == words.count - 1 ? "*" : "")
        }.joined(separator: " ")
    }
}
