import Foundation

/// Why gofmt gave back no formatted code.
public enum GoFormatError: LocalizedError, Equatable {
    /// The code doesn't parse: gofmt's messages, naming the file, such as
    /// "main.go:4:1: expected operand, found '}'".
    case syntax(String)
    /// gofmt could not run, took too long, or printed more than the code could be.
    case failed(String)

    public var errorDescription: String? {
        switch self {
        case .syntax(let messages): messages
        case .failed(let why): why
        }
    }
}

/// Formats Go code with the gofmt of the Go that runs it. The code goes in through standard
/// input and comes back on standard output, so nothing is written to the lesson folder.
public struct GoFormatter: Sendable {
    public static let timeout: Duration = .seconds(10)

    let gofmt: String
    let environment: [String: String]

    /// Nil when the toolchain has no gofmt.
    public init?(toolchain: Toolchain) {
        guard let gofmt = toolchain.gofmt else { return nil }
        self.gofmt = gofmt
        environment = toolchain.goEnvironment(submit: false)
    }

    /// The code as gofmt lays it out; code that is formatted already comes back the same.
    /// Cancelling the task stops gofmt and throws `CancellationError`.
    @concurrent
    public func format(_ source: String, fileName: String) async throws -> String {
        let input = Data(source.utf8)
        // gofmt's output is about as long as its input; much longer means something is wrong,
        // and a cut-off copy must never replace the code.
        let out = await ProcessRunner.run(ProcessRequest(
            executable: gofmt, arguments: [], directory: NSHomeDirectory(), environment: environment,
            timeout: GoFormatter.timeout, outputLimit: 2 * input.count + (1 << 20), input: input))
        if out.cancelled { throw CancellationError() }
        if let launch = out.launchError { throw GoFormatError.failed("gofmt could not run: \(launch)") }
        if out.timedOut { throw GoFormatError.failed("gofmt took longer than \(GoFormatter.timeout.components.seconds) seconds.") }
        if out.truncated { throw GoFormatError.failed("gofmt's output was too long.") }
        if out.signal != 0 { throw GoFormatError.failed("gofmt stopped: signal: \(String(cString: strsignal(out.signal)))") }
        guard out.exitCode == 0 else {
            let messages = GoFormatter.naming(out.stderr, fileName).trimmingCharacters(in: .whitespacesAndNewlines)
            throw messages.isEmpty ? GoFormatError.failed("gofmt failed with exit status \(out.exitCode).") : GoFormatError.syntax(messages)
        }
        return out.stdout
    }

    /// gofmt's messages with the file's name where they say "<standard input>".
    static func naming(_ messages: String, _ fileName: String) -> String {
        messages.replacing(/(?m)^<standard input>:/) { _ in fileName + ":" }
    }
}
