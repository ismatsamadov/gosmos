import Foundation
import SQLite3

/// Holds a database's write lock on a connection of its own, as another process writing would,
/// so the app's next write waits (its busy timeout is 10 s) until `release`. SQLite's own API:
/// GRDB ends every access with no transaction left open.
final class WriteLock {
    struct Failed: Error { let code: Int32 }
    private var db: OpaquePointer?

    init(_ url: URL) throws {
        var code = sqlite3_open_v2(url.path, &db, SQLITE_OPEN_READWRITE, nil)
        // A write of the app's still landing, such as the one that records a lesson just loaded
        // as opened, is waited for, as the app waits for another's.
        if code == SQLITE_OK { code = sqlite3_busy_timeout(db, 10_000) }
        if code == SQLITE_OK { code = sqlite3_exec(db, "BEGIN IMMEDIATE", nil, nil, nil) }
        guard code == SQLITE_OK else { sqlite3_close(db); throw Failed(code: code) }
    }

    func release() throws {
        let code = sqlite3_exec(db, "COMMIT", nil, nil, nil)
        sqlite3_close(db)
        guard code == SQLITE_OK else { throw Failed(code: code) }
    }
}
