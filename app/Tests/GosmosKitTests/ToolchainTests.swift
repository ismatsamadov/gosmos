import Foundation
import Testing
@testable import GosmosKit

/// Which Go runs lessons: the one built into the app, unless Settings chose another; and how
/// every Go runs: offline, isolated, with its caches in the app's.
@Suite struct ToolchainTests {
    /// A Go in `dir`, whose go command answers `go env GOVERSION GOROOT GOTOOLDIR` as a real
    /// one does: with `goroot` as its GOROOT if given, and with `tools` in its tool folder.
    static func fakeGo(in dir: URL, version: String = "go1.27.3", goroot: String? = nil,
                       tools: [String] = ["compile", "link", "asm", "vet"]) throws -> URL {
        let fm = FileManager.default
        let bin = dir.appendingPathComponent("bin", isDirectory: true)
        let toolDir = dir.appendingPathComponent("pkg/tool/darwin_arm64", isDirectory: true)
        try fm.createDirectory(at: bin, withIntermediateDirectories: true)
        try fm.createDirectory(at: toolDir, withIntermediateDirectories: true)
        let answer = "#!/bin/sh\nprintf '%s\\n%s\\n%s\\n' '\(version)' '\(goroot ?? dir.path)' '\(toolDir.path)'\n"
        for (folder, name, script) in [(bin, "go", answer), (bin, "gofmt", answer)] + tools.map({ (toolDir, $0, "#!/bin/sh\n") }) {
            let path = folder.appendingPathComponent(name).path
            try script.write(toFile: path, atomically: true, encoding: .utf8)
            try fm.setAttributes([.posixPermissions: 0o755], ofItemAtPath: path)
        }
        return dir
    }

    @Test func prefersTheBuiltInGo() async throws {
        let dir = try Fixtures.tempDir("toolchain")
        defer { try? FileManager.default.removeItem(at: dir) }
        let builtIn = try Self.fakeGo(in: dir.appendingPathComponent("Gosmos.app/Contents/Resources/Go"))
        let installed = try Self.fakeGo(in: dir.appendingPathComponent("usr/local/go"), version: "go1.28.1")
        let caches = dir.appendingPathComponent("caches")
        // A Go on PATH changes nothing.
        let env = ["PATH": installed.appendingPathComponent("bin").path + ":/usr/bin:/bin", "HOME": dir.path]
        let t = await Toolchain.discover(goOverride: nil, builtIn: builtIn, caches: caches, environment: env)
        #expect(t.go == builtIn.appendingPathComponent("bin/go").path && t.isBuiltIn && t.problem == nil, "\(t.problem ?? "")")
        #expect(t.gofmt == builtIn.appendingPathComponent("bin/gofmt").path)
        #expect(t.goVersion == "go1.27.3" && t.name == "Built-in Go 1.27.3" && t.languageVersion == "1.27")
        #expect(t.caches == caches)
        // Go's telemetry, kept in the caches, is off.
        #expect(try String(contentsOf: caches.appendingPathComponent("telemetry/mode"), encoding: .utf8) == "off")
    }

    @Test func honorsAChosenGo() async throws {
        let dir = try Fixtures.tempDir("toolchain")
        defer { try? FileManager.default.removeItem(at: dir) }
        let builtIn = try Self.fakeGo(in: dir.appendingPathComponent("Gosmos.app/Contents/Resources/Go"))
        let chosen = try Self.fakeGo(in: dir.appendingPathComponent("sdk/go1.28.1"), version: "go1.28.1").appendingPathComponent("bin/go").path
        let t = await Toolchain.discover(goOverride: chosen, builtIn: builtIn, caches: dir.appendingPathComponent("caches"),
                                         environment: ["PATH": "/usr/bin:/bin", "HOME": dir.path])
        #expect(t.go == chosen && !t.isBuiltIn && t.problem == nil, "\(t.problem ?? "")")
        #expect(t.name == "Go 1.28.1" && t.languageVersion == "1.28")
    }

    /// A Go built from source says "devel go1.29-abcdef …" for its version: it runs, and new
    /// lesson folders ask for its language version.
    @Test func aDevelopmentGoRuns() async throws {
        let dir = try Fixtures.tempDir("toolchain")
        defer { try? FileManager.default.removeItem(at: dir) }
        let chosen = try Self.fakeGo(in: dir.appendingPathComponent("goroot"), version: "devel go1.29-6f3f1d8 Tue Oct 28 12:00:00 2026 +0000")
            .appendingPathComponent("bin/go").path
        let t = await Toolchain.discover(goOverride: chosen, builtIn: nil, caches: dir.appendingPathComponent("caches"),
                                         environment: ["PATH": "/usr/bin:/bin", "HOME": dir.path])
        #expect(t.go == chosen && t.problem == nil, "\(t.problem ?? "")")
        #expect(t.name == "Go devel go1.29-6f3f1d8" && t.languageVersion == "1.29")
        #expect(Toolchain.languageVersion("go1.27.3") == "1.27" && Toolchain.languageVersion("go1.28rc1") == "1.28")
    }

    @Test func aBrokenGoSaysWhyAndIsNotReplaced() async throws {
        let dir = try Fixtures.tempDir("toolchain")
        defer { try? FileManager.default.removeItem(at: dir) }
        let caches = dir.appendingPathComponent("caches")
        func discover(_ chosen: String?, builtIn: URL?) async -> Toolchain {
            await Toolchain.discover(goOverride: chosen, builtIn: builtIn, caches: caches, environment: ["PATH": "/usr/bin:/bin", "HOME": dir.path])
        }

        // No Go in the app at all.
        let missing = await discover(nil, builtIn: dir.appendingPathComponent("nowhere"))
        #expect(missing.go == nil && missing.isBuiltIn && missing.name == nil)
        #expect(missing.problem == "The Go built into Gosmos is missing. Reinstall Gosmos, or choose another Go in Settings.")

        // Its go command gone, its compiler gone, or it would use the Go installed elsewhere.
        let noGo = try Self.fakeGo(in: dir.appendingPathComponent("no-go"))
        try FileManager.default.removeItem(at: noGo.appendingPathComponent("bin/go"))
        let noCompiler = try Self.fakeGo(in: dir.appendingPathComponent("no-compiler"), tools: ["link", "asm", "vet"])
        let elsewhere = try Self.fakeGo(in: dir.appendingPathComponent("elsewhere"), goroot: "/opt/homebrew/Cellar/go/1.27.3/libexec")
        for (root, why) in [(noGo, "there is no go command at"), (noCompiler, "its compile tool is missing"),
                            (elsewhere, "it uses the Go in /opt/homebrew/Cellar/go/1.27.3/libexec")] {
            let t = await discover(nil, builtIn: root)
            #expect(t.go == nil && t.isBuiltIn, "\(root.lastPathComponent)")
            #expect(t.problem?.hasPrefix("The Go built into Gosmos is damaged: \(why)") == true, "\(t.problem ?? "")")
            #expect(t.problem?.hasSuffix("Reinstall Gosmos, or choose another Go in Settings.") == true)
        }

        // A chosen Go that does not run is named, and the built-in one does not stand in for it.
        let builtIn = try Self.fakeGo(in: dir.appendingPathComponent("built-in"))
        let gone = await discover("/nowhere/bin/go", builtIn: builtIn)
        #expect(gone.go == nil && !gone.isBuiltIn)
        #expect(gone.problem == "The Go chosen in Settings, /nowhere/bin/go, does not run: there is no go command at /nowhere/bin/go. Choose another, or use the built-in Go.")
    }

    @Test func goRunsOfflineIsolatedAndWithTheAppsCaches() {
        let go = "/Applications/Gosmos.app/Contents/Resources/Go/bin/go"
        let t = Toolchain(environment: ["PATH": "/usr/bin:/bin", "HOME": "/Users/learner",
                                        // What a shell may hand the app; none of it reaches Go.
                                        "GOROOT": "/opt/homebrew/Cellar/go/1.27.3/libexec", "GOFLAGS": "-race", "GOOS": "linux",
                                        "GOPROXY": "https://proxy.golang.org", "GOTOOLCHAIN": "auto", "GODEBUG": "panicnil=1",
                                        "CGO_ENABLED": "1", "CGO_CFLAGS": "-O0", "GOSMOS_SUBMIT": "1"],
                          go: go, goVersion: "go1.27.3", isBuiltIn: true, caches: URL(fileURLWithPath: "/tmp/gosmos-caches"))
        let env = t.goEnvironment(submit: false)
        // Nothing from the network: no other toolchain, no module, no checksum, no telemetry.
        #expect(env["GOTOOLCHAIN"] == "local" && env["GOPROXY"] == "off" && env["GOSUMDB"] == "off")
        #expect(env["TEST_TELEMETRY_DIR"] == "/tmp/gosmos-caches/telemetry")
        // Caches in the app's: nothing in ~/go or ~/Library/Caches/go-build.
        #expect(env["GOCACHE"] == "/tmp/gosmos-caches/go-build" && env["GOPATH"] == "/tmp/gosmos-caches/go")
        // None of your own Go's settings, no C compiler, modules always.
        #expect(env["GOENV"] == "off" && env["CGO_ENABLED"] == "0" && env["GOWORK"] == "off" && env["GO111MODULE"] == "on")
        #expect(env["GOROOT"] == nil && env["GOFLAGS"] == nil && env["GOOS"] == nil && env["GODEBUG"] == nil && env["CGO_CFLAGS"] == nil)
        #expect(env["GOSMOS_SUBMIT"] == "0" && t.goEnvironment(submit: true)["GOSMOS_SUBMIT"] == "1")
        #expect(env["PATH"] == "/Applications/Gosmos.app/Contents/Resources/Go/bin:/usr/bin:/bin" && env["HOME"] == "/Users/learner")
    }
}

/// The Go in a built Gosmos.app (`make app`), run as the app runs it on a Mac with no Go
/// installed: PATH is /usr/bin:/bin and HOME an empty folder. `swift test` alone has no app
/// and skips it.
@Suite struct BuiltInGoTests {
    /// This repository's build/Gosmos.app/Contents/Resources/Go, when it is there.
    static let root: URL? = {
        let repo = URL(fileURLWithPath: #filePath).deletingLastPathComponent().deletingLastPathComponent()
            .deletingLastPathComponent().deletingLastPathComponent()
        let go = repo.appendingPathComponent("build/Gosmos.app/Contents/Resources/Go", isDirectory: true)
        return FileManager.default.fileExists(atPath: go.appendingPathComponent("bin/go").path) ? go : nil
    }()

    /// build/ is inside this repository's Go module, so every .go file of the built-in Go must be
    /// in Go's own module, src/go.mod, which ./... (go vet, go test) skips.
    @Test(.enabled(if: root != nil, "no build/Gosmos.app: make app builds it"))
    func staysOutOfThisRepositorysModule() throws {
        let root = try #require(Self.root)
        let fm = FileManager.default
        #expect(fm.fileExists(atPath: root.appendingPathComponent("src/go.mod").path))
        let outside = try fm.contentsOfDirectory(atPath: root.path).filter { $0 != "src" }.flatMap { top in
            [top] + (fm.subpaths(atPath: root.appendingPathComponent(top).path) ?? []).map { top + "/" + $0 }
        }.filter { $0.hasSuffix(".go") }
        #expect(outside.isEmpty, "\(outside)")
    }

    @Test(.enabled(if: root != nil, "no build/Gosmos.app: make app builds it"))
    func runsTestsAndFormatsWithNoGoInstalled() async throws {
        let root = try #require(Self.root)
        let dir = try Fixtures.tempDir("built-in-go")
        defer { try? FileManager.default.removeItem(at: dir) }
        let home = dir.appendingPathComponent("home", isDirectory: true)
        try FileManager.default.createDirectory(at: home, withIntermediateDirectories: true)
        let t = await Toolchain.discover(goOverride: nil, builtIn: root, caches: Fixtures.goCaches,
                                         environment: ["PATH": "/usr/bin:/bin", "HOME": home.path])
        #expect(t.problem == nil && t.isBuiltIn && t.go == root.appendingPathComponent("bin/go").path, "\(t.problem ?? "")")
        let version = try #require(t.goVersion)
        let ws = Workspace(root: dir.appendingPathComponent("work"))
        let runner = GoRunner(toolchain: t, scratch: ws.runsFolder)

        // A program, built and run: it was compiled by this Go.
        let program = CodeExercise(starterFiles: [StarterFile(name: "main.go", content: "package main\n\nimport (\n\t\"fmt\"\n\t\"runtime\"\n)\n\nfunc main() {\n\tfmt.Println(runtime.Version())\n}\n")],
                                   expectedOutput: version + "\n")
        let p = Fixtures.codeLesson("built-in-program", .code, program)
        let ran = await runner.run(program, tests: false, folder: try ws.prepare(p, code: program, saved: [:], goVersion: t.languageVersion).folder, submit: true)
        #expect(ran.status == .passed, "\(ran.output)")

        // Tests, through go test and its vet checks.
        let tests = CodeExercise(starterFiles: [
            StarterFile(name: "main.go", content: "package main\n\nfunc add(a, b int) int {\n\treturn a + b\n}\n"),
            StarterFile(name: "main_test.go", content: Fixtures.addTests, isReadonly: true),
        ])
        let l = Fixtures.codeLesson("built-in-tests", .codeTests, tests)
        let tested = await runner.run(tests, tests: true, folder: try ws.prepare(l, code: tests, saved: [:], goVersion: t.languageVersion).folder, submit: true)
        #expect(tested.status == .passed, "\(tested.output)")

        // Format, with its gofmt.
        let formatter = try #require(GoFormatter(toolchain: t))
        #expect(try await formatter.format("package main\nfunc main(){}", fileName: "main.go") == "package main\n\nfunc main() {}\n")

        // Nothing went to HOME (no ~/go, no ~/Library), and telemetry counted nothing.
        #expect(try FileManager.default.contentsOfDirectory(atPath: home.path).isEmpty)
        #expect(!FileManager.default.fileExists(atPath: Fixtures.goCaches.appendingPathComponent("telemetry/local").path))
    }
}

/// Settings › Go checks the Go it names each time it changes, and only the latest check counts:
/// a Go that hangs until its timeout, chosen and then replaced by the built-in one, never comes
/// back when its check ends at last.
@Suite @MainActor struct LatestWorkTests {
    @MainActor final class Results {
        var applied: [String] = []
    }

    @Test func anEarlierCheckThatEndsLaterNeverReplacesTheLatest() async throws {
        let latest = LatestWork<String>()
        #expect(await latest.value() == nil)
        let results = Results()
        let (gate, open) = AsyncStream<Void>.makeStream()
        latest.start({ for await _ in gate {}; return "hung" }) { results.applied.append($0) }
        // Something waits for a Go to run code with meanwhile, such as a Run.
        async let waited = latest.value()
        try await Task.sleep(for: .milliseconds(50))
        latest.start({ "built-in" }) { results.applied.append($0) }
        try await Task.sleep(for: .milliseconds(50))
        #expect(results.applied == ["built-in"])
        open.finish()
        #expect(await waited == "built-in")
        #expect(await latest.value() == "built-in")
        try await Task.sleep(for: .milliseconds(100))
        #expect(results.applied == ["built-in"])
    }
}
