import Foundation

/// The cases of db/testdata/checks, which the library tool's tests read too (db/checks_test.go):
/// library.sql's library, changed by a case's .sql (run with foreign keys off, outside a
/// transaction) or replaced by its .file, and the problems its .want says the checks report.
/// The libraries are made by macOS's own sqlite3 tool: a case may write around the schema
/// (writable_schema), which the app's SQLite, in its defensive mode, refuses.
enum SharedChecks {
    struct Failed: Error, CustomStringConvertible {
        let description: String
    }

    static let folder = Fixtures.repository.appendingPathComponent("db/testdata/checks")

    static func cases() -> [String] {
        let names = (try? FileManager.default.contentsOfDirectory(atPath: folder.path)) ?? []
        return names.filter { $0.hasSuffix(".want") }.map { String($0.dropLast(5)) }.sorted()
    }

    static func want(_ name: String) throws -> String {
        try String(contentsOf: folder.appendingPathComponent(name + ".want"), encoding: .utf8)
    }

    static func make(_ name: String, at url: URL) throws {
        if let raw = try? Data(contentsOf: folder.appendingPathComponent(name + ".file")) {
            try raw.write(to: url)
            return
        }
        let script = [
            try String(contentsOf: Fixtures.schemaURL, encoding: .utf8),
            try String(contentsOf: folder.appendingPathComponent("library.sql"), encoding: .utf8),
            "PRAGMA foreign_keys = OFF;\n",
            try String(contentsOf: folder.appendingPathComponent(name + ".sql"), encoding: .utf8),
        ].joined(separator: "\n")
        let sqlite = Process()
        sqlite.executableURL = URL(fileURLWithPath: "/usr/bin/sqlite3")
        sqlite.arguments = ["-bail", url.path]
        let input = Pipe(), errors = Pipe()
        sqlite.standardInput = input
        sqlite.standardError = errors
        sqlite.standardOutput = FileHandle.nullDevice
        try sqlite.run()
        input.fileHandleForWriting.write(Data(script.utf8))
        try input.fileHandleForWriting.close()
        sqlite.waitUntilExit()
        let said = String(decoding: errors.fileHandleForReading.readDataToEndOfFile(), as: UTF8.self)
        guard sqlite.terminationStatus == 0, said.isEmpty else {
            throw Failed(description: "making case \(name): sqlite3 exited \(sqlite.terminationStatus): \(said)")
        }
    }
}
