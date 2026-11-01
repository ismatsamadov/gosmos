import Foundation
import os

/// Loggers for the app's subsystems.
public enum Log {
    public static let subsystem = "io.github.jobloom.Gosmos"
    public static let app = Logger(subsystem: subsystem, category: "app")
}

/// ISO 8601 timestamps in UTC with a trailing Z, the only form the schema's CHECKs accept,
/// such as 2026-10-31T12:00:00Z.
public enum Clock {
    public static func now() -> String { string(Date()) }

    public static func string(_ date: Date) -> String {
        date.ISO8601Format() // UTC, whole seconds, trailing Z
    }

    public static func date(_ text: String) -> Date? {
        if let d = try? Date(text, strategy: .iso8601) { return d }
        return try? Date(text, strategy: Date.ISO8601FormatStyle(includingFractionalSeconds: true))
    }
}

/// Paths of the system programs the app relies on.
public enum SystemTools {
    public static let zsh = "/bin/zsh"
}

/// Which errors reach the person. Work is cancelled when something newer replaces it: another
/// lesson, the next keystroke of a search, another library. That work ends silently, whatever
/// it threw on the way out; the newer work meets any real problem again.
public enum Reporting {
    /// Whether `error` stays unsaid: a cancellation, or anything thrown in a cancelled task.
    public static func isSilent(_ error: any Error) -> Bool {
        error is CancellationError || Task.isCancelled
    }
}

/// Settings this Mac kept for what the app no longer does: one open library chosen by path,
/// its recent list, git sync, and code conflicts named without their library. Removed when
/// the app starts, so nothing reads them by mistake.
public enum ObsoletePreferences {
    public static let keys = ["LibraryPath", "RecentLibraries", "LastPush", "SyncAutomatically", "UnsettledFiles"]

    public static func remove(from defaults: UserDefaults) {
        for key in keys { defaults.removeObject(forKey: key) }
    }
}
