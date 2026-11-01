import Foundation

/// Writes of your study data that nothing else waits for, such as a review answer, a revealed
/// hint, a done mark or a setting. Each is counted from the moment it is asked for, on the main
/// actor, until it lands or fails, so quitting can wait for every one (`finish`), as it waits
/// for a lesson's saves; an answer given just before quitting is kept. Each runs in a task of
/// its own, which no lesson switch cancels: in a cancelled task the database refuses to write.
@MainActor
public final class PendingWrites {
    private var running = 0
    private var waiting: [CheckedContinuation<Void, Never>] = []

    nonisolated public init() {}

    /// Starts `write` and counts it until it ends, however it ends. The task's value is what
    /// `write` returns, for a caller that needs it.
    @discardableResult
    public func start<T: Sendable>(_ write: @escaping @MainActor () async throws -> T) -> Task<T, any Error> {
        running += 1
        return Task {
            defer { ended() }
            return try await write()
        }
    }

    /// Waits until every write started so far has ended, and any started meanwhile.
    public func finish() async {
        guard running > 0 else { return }
        await withCheckedContinuation { waiting.append($0) }
    }

    private func ended() {
        running -= 1
        guard running == 0 else { return }
        let waiters = waiting
        waiting = []
        for waiter in waiters { waiter.resume() }
    }
}
