import Foundation
import os

/// Loggers for the app's subsystems.
public enum Log {
    public static let subsystem = "app.gosmos.Gosmos"
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

/// Which errors reach the person. Work is cancelled when something newer replaces it: another
/// lesson, the next keystroke of a search, another library. That work ends silently, whatever
/// it threw on the way out; the newer work meets any real problem again.
public enum Reporting {
    /// Whether `error` stays unsaid: a cancellation, or anything thrown in a cancelled task.
    public static func isSilent(_ error: any Error) -> Bool {
        error is CancellationError || Task.isCancelled
    }
}

/// One copy of the app at a time. Two copies of the same app, such as ~/Applications/Gosmos.app
/// and a build's, would share the learner database and compete for it, so a launch while another
/// copy runs brings that one forward and quits before it opens anything. A copy under another
/// bundle id, such as a QA build, is another app and runs beside it.
public enum SingleInstance {
    /// A running app as the system lists it: its process, its bundle id and when it launched.
    public struct App: Hashable, Sendable {
        public let pid: Int32
        public let bundleIdentifier: String?
        public let launchDate: Date?

        public init(pid: Int32, bundleIdentifier: String?, launchDate: Date? = nil) {
            self.pid = pid
            self.bundleIdentifier = bundleIdentifier
            self.launchDate = launchDate
        }
    }

    /// The copy that `me`, a launch of the app, gives way to: another running process with the
    /// same bundle id that launched before it (by launch date, else by process id), the first
    /// of them when there are several. Of two copies launched at once, only the later one quits.
    /// Nil when there is none, or when the app has no bundle id, as when `swift run` starts the
    /// bare binary.
    public static func runningCopy(for me: App, in running: [App]) -> App? {
        guard let id = me.bundleIdentifier, !id.isEmpty else { return nil }
        func earlier(_ a: App, _ b: App) -> Bool {
            if let x = a.launchDate, let y = b.launchDate, x != y { return x < y }
            return a.pid < b.pid
        }
        return running.filter { $0.pid != me.pid && $0.bundleIdentifier == id && earlier($0, me) }.min(by: earlier)
    }
}

/// The result of work run again and again, such as checking the Go that Settings › Go names each
/// time it changes: only the latest work's result counts. Work that newer work replaced while it
/// ran, such as the check of a Go that hangs until its timeout, is dropped when it ends,
/// whenever that is.
@MainActor
public final class LatestWork<Value: Sendable> {
    private var task: Task<Value, Never>?

    public init() {}

    /// Starts `work` off the main actor in place of the work before it, and calls `done` with
    /// its result once it ends, unless newer work started meanwhile.
    public func start(priority: TaskPriority? = nil, _ work: @escaping @Sendable () async -> Value,
                      done: @escaping @MainActor (Value) -> Void) {
        let task = Task.detached(priority: priority, operation: work)
        self.task = task
        Task {
            let value = await task.value
            if self.task == task { done(value) }
        }
    }

    /// The latest work's result, waiting for it, and past any work that newer work replaces
    /// meanwhile; nil when no work was started.
    public func value() async -> Value? {
        while let task {
            let value = await task.value
            if self.task == task { return value }
        }
        return nil
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
