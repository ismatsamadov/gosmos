import Foundation

/// The Go toolchain and the environment programs run with. Lessons run with the Go built into
/// the app (Contents/Resources/Go), so nothing needs installing; Settings can choose another Go
/// on this Mac instead. Either runs isolated and offline: it never downloads a toolchain or a
/// module, never uploads telemetry, ignores the settings `go env -w` keeps for your own Go, never
/// calls a C compiler, and keeps its caches in the app's (~/Library/Caches/Gosmos).
public struct Toolchain: Sendable {
    /// The environment programs run with: the app's own.
    public var environment: [String: String]
    public var go: String?
    /// Like "go1.27.3".
    public var goVersion: String?
    /// The same Go's gofmt, in the bin folder of its GOROOT; nil when it isn't there.
    public var gofmt: String?
    /// Whether this is the Go built into the app, rather than one chosen in Settings.
    public var isBuiltIn: Bool
    /// Why there is no Go to run: the built-in one is missing or damaged, or the chosen one
    /// does not run.
    public var problem: String?
    /// Where Go keeps its build cache (go-build), its module cache (go) and its telemetry
    /// setting (telemetry).
    public var caches: URL

    public init(environment: [String: String], go: String? = nil, goVersion: String? = nil, gofmt: String? = nil,
                isBuiltIn: Bool = false, problem: String? = nil, caches: URL = Toolchain.defaultCaches) {
        self.environment = environment
        self.go = go
        self.goVersion = goVersion
        self.gofmt = gofmt
        self.isBuiltIn = isBuiltIn
        self.problem = problem
        self.caches = caches
    }

    /// The Go built into the app.
    public static var builtInRoot: URL? {
        Bundle.main.resourceURL?.appendingPathComponent("Go", isDirectory: true)
    }

    /// ~/Library/Caches/Gosmos, which the lesson folders share.
    public static var defaultCaches: URL {
        FileManager.default.urls(for: .cachesDirectory, in: .userDomainMask)[0].appendingPathComponent("Gosmos", isDirectory: true)
    }

    /// The go.mod language version for new lesson folders, like "1.27".
    public var languageVersion: String { Toolchain.languageVersion(goVersion ?? "") }

    /// The language version of a GOVERSION: "1.27" for go1.27.3 or go1.27rc1, and for a Go
    /// built from source, "devel go1.29-6f3f1d8 …", the one it is developing.
    public static func languageVersion(_ v: String) -> String {
        guard let m = v.firstMatch(of: /go(\d+)\.(\d+)/) else { return "1.22" }
        return "\(m.1).\(m.2)"
    }

    /// The Go as Settings names it: "Built-in Go 1.27.3", "Go 1.28.1" for a chosen one, or "Go
    /// devel go1.29-6f3f1d8" for one built from source; nil without a Go to run.
    public var name: String? {
        guard go != nil, let words = goVersion?.split(separator: " "), let version = words.first else { return nil }
        if version == "devel", words.count > 1 { return (isBuiltIn ? "Built-in Go devel " : "Go devel ") + words[1] }
        return (isBuiltIn ? "Built-in Go " : "Go ") + version.dropFirst(version.hasPrefix("go") ? 2 : 0)
    }

    /// The Go chosen in Settings (`goOverride`), else the one built into the app. A Go that does
    /// not run is reported in `problem`, never swapped for another. Off the main thread; it runs
    /// `go env` (with a timeout). `environment` is the app's own; tests pass theirs.
    public static func discover(goOverride: String?, builtIn: URL? = Toolchain.builtInRoot, caches: URL = Toolchain.defaultCaches,
                                environment: [String: String] = ProcessInfo.processInfo.environment) async -> Toolchain {
        var env = environment
        env["HOME"] = env["HOME"] ?? NSHomeDirectory()
        var t = Toolchain(environment: env, isBuiltIn: goOverride == nil, caches: caches)
        t.turnTelemetryOff()
        if let chosen = goOverride {
            let go = (chosen as NSString).expandingTildeInPath
            switch await t.probe(go, root: nil) {
            case .found(let version, let gofmt): t.use(go, version, gofmt)
            case .failed(let why): t.problem = "The Go chosen in Settings, \(go), does not run: \(why). Choose another, or use the built-in Go."
            }
        } else if let root = builtIn, FileManager.default.fileExists(atPath: root.path) {
            let go = root.appendingPathComponent("bin/go").path
            switch await t.probe(go, root: root.path) {
            case .found(let version, let gofmt): t.use(go, version, gofmt)
            case .failed(let why): t.problem = "The Go built into Gosmos is damaged: \(why). Reinstall Gosmos, or choose another Go in Settings."
            }
        } else {
            t.problem = "The Go built into Gosmos is missing. Reinstall Gosmos, or choose another Go in Settings."
        }
        Log.app.notice("Go: \(t.go ?? "none", privacy: .public) \(t.goVersion ?? "", privacy: .public) \(t.problem ?? "", privacy: .public)")
        return t
    }

    private mutating func use(_ go: String, _ version: String, _ gofmt: String?) {
        self.go = go
        goVersion = version
        self.gofmt = gofmt
    }

    private enum Probe {
        case found(version: String, gofmt: String?)
        case failed(String)
    }

    /// Asks the go command at `go` what it is. It fails when it doesn't run, when its compiler,
    /// linker, assembler or vet are missing, or, for the built-in Go, when it would use another
    /// GOROOT than `root`.
    private func probe(_ go: String, root: String?) async -> Probe {
        guard let go = Toolchain.runnable(go) else { return .failed("there is no go command at \(go)") }
        var probing = self
        probing.go = go
        // From the root folder, where no go.mod can ask for another version.
        let out = await ProcessRunner.run(ProcessRequest(executable: go, arguments: ["env", "GOVERSION", "GOROOT", "GOTOOLDIR"], directory: "/",
                                                         environment: probing.goEnvironment(submit: false), timeout: .seconds(20)))
        // One line each: the version, the folder Go is installed in, the folder of its tools. A
        // release says go1.27.3; a Go built from source, devel go1.29-6f3f1d8 and a date.
        let lines = out.stdout.components(separatedBy: "\n").map { $0.trimmingCharacters(in: .whitespaces) }
        if let launch = out.launchError { return .failed(launch) }
        if out.timedOut { return .failed("it did not answer within 20 seconds") }
        guard out.succeeded, lines.count >= 3, lines[0].hasPrefix("go") || lines[0].hasPrefix("devel") else {
            let said = out.combined.trimmingCharacters(in: .whitespacesAndNewlines).components(separatedBy: "\n").first ?? ""
            return .failed(said.isEmpty ? "go env failed with exit status \(out.exitCode)" : said)
        }
        if let root, GoRunner.realPath(lines[1]) != GoRunner.realPath(root) {
            return .failed("it uses the Go in \(lines[1])")
        }
        for tool in ["compile", "link", "asm", "vet"] where Toolchain.runnable(lines[2] + "/" + tool) == nil {
            return .failed("its \(tool) tool is missing")
        }
        return .found(version: lines[0], gofmt: Toolchain.runnable(lines[1] + "/bin/gofmt"))
    }

    private static func runnable(_ path: String) -> String? {
        var isDir: ObjCBool = false
        guard FileManager.default.fileExists(atPath: path, isDirectory: &isDir), !isDir.boolValue,
              FileManager.default.isExecutableFile(atPath: path) else { return nil }
        return path
    }

    /// Go's telemetry stays off: the go command and its tools keep it in TEST_TELEMETRY_DIR,
    /// the caches' telemetry folder, so nothing is counted or uploaded, whatever `go telemetry`
    /// says for your own Go. Without the mode file Go only counts, in that folder, and uploads
    /// nothing either.
    private func turnTelemetryOff() {
        let folder = caches.appendingPathComponent("telemetry", isDirectory: true)
        let mode = folder.appendingPathComponent("mode")
        guard (try? String(contentsOf: mode, encoding: .utf8)) != "off" else { return }
        do {
            try FileManager.default.createDirectory(at: folder, withIntermediateDirectories: true)
            try Data("off".utf8).write(to: mode, options: .atomic)
        } catch {
            Log.app.error("Go telemetry: \(error.localizedDescription, privacy: .public)")
        }
    }

    /// The environment for running Go: the go command's folder first on PATH, so tools it
    /// starts find the same Go; none of the Go settings the app was started with (GOROOT,
    /// GOFLAGS, GOOS and the rest); and nothing that would reach the network, switch toolchains
    /// or modules modes, or call a C compiler.
    public func goEnvironment(submit: Bool) -> [String: String] {
        var env = environment.filter { !$0.key.hasPrefix("GO") && !$0.key.hasPrefix("CGO_") }
        env["GOTOOLCHAIN"] = "local"
        env["GOPROXY"] = "off"
        env["GOSUMDB"] = "off"
        env["GOENV"] = "off"
        env["GOWORK"] = "off"
        env["GO111MODULE"] = "on"
        // Lessons are pure Go; a Mac without the Command Line Tools has no C compiler.
        env["CGO_ENABLED"] = "0"
        env["GOCACHE"] = caches.appendingPathComponent("go-build", isDirectory: true).path
        env["GOPATH"] = caches.appendingPathComponent("go", isDirectory: true).path
        env["TEST_TELEMETRY_DIR"] = caches.appendingPathComponent("telemetry", isDirectory: true).path
        env["GOSMOS_SUBMIT"] = submit ? "1" : "0"
        if let go {
            env["PATH"] = (go as NSString).deletingLastPathComponent + ":" + (env["PATH"] ?? "/usr/bin:/bin")
        }
        return env
    }
}
