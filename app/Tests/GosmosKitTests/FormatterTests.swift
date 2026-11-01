import Foundation
import Testing
@testable import GosmosKit

/// gofmt from the Go on this machine, with the code through standard input and output.
@Suite struct GoFormatterTests {
    /// gofmt beside the Go this machine has, or nil: then the tests that run it are skipped.
    static func formatter() async -> GoFormatter? {
        guard let toolchain = await Fixtures.toolchain() else { return nil }
        return GoFormatter(toolchain: toolchain)
    }

    static let formatted = "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tx := 1\n\tfmt.Println(x)\n}\n"

    @Test func formatsTheCode() async throws {
        guard let formatter = await Self.formatter() else { return } // no gofmt on this machine
        let messy = "package main\nimport \"fmt\"\nfunc main(){\nx:=1\n  fmt.Println( x )\n}"
        #expect(try await formatter.format(messy, fileName: "main.go") == Self.formatted)
    }

    @Test func aSyntaxErrorNamesTheFile() async throws {
        guard let formatter = await Self.formatter() else { return }
        let broken = "package main\n\nfunc main() {\n\tx := \n}\n"
        let error = await #expect(throws: GoFormatError.self) { try await formatter.format(broken, fileName: "orbit/orbit.go") }
        guard case .syntax(let messages)? = error else {
            Issue.record("not a syntax error: \(String(describing: error))")
            return
        }
        #expect(messages.hasPrefix("orbit/orbit.go:5:1: "), "\(messages)")
        #expect(!messages.contains("<standard input>") && !messages.hasSuffix("\n"))
    }

    @Test func formattedCodeComesBackTheSame() async throws {
        guard let formatter = await Self.formatter() else { return }
        #expect(try await formatter.format(Self.formatted, fileName: "main.go") == Self.formatted)
        // Longer than a pipe holds, both ways, so gofmt reads and writes while the runner
        // still writes and reads.
        let long = "package main\n" + (0..<6000).map { "\nfunc f\($0)(a, b int) int {\n\treturn a*b + \($0)\n}\n" }.joined()
        #expect(long.utf8.count > 256 << 10)
        #expect(try await formatter.format(long, fileName: "main.go") == long)
    }

    @Test func messagesNameTheFile() {
        let messages = "<standard input>:5:4: expected ';', found ':='\n<standard input>:6:3: expected '}', found 'EOF'\n"
        #expect(GoFormatter.naming(messages, "main.go") == "main.go:5:4: expected ';', found ':='\nmain.go:6:3: expected '}', found 'EOF'\n")
        // Only where a line names the input: a message quoting the same words keeps them.
        #expect(GoFormatter.naming("x.go:1:1: found '<standard input>:'", "main.go") == "x.go:1:1: found '<standard input>:'")
    }

    @Test func noFormatterWithoutGofmt() {
        #expect(GoFormatter(toolchain: Toolchain(environment: [:], go: "/usr/local/go/bin/go", goVersion: "go1.27.3")) == nil)
        #expect(GoFormatter(toolchain: Toolchain(environment: [:], go: "/usr/local/go/bin/go", goVersion: "go1.27.3", gofmt: "/usr/local/go/bin/gofmt")) != nil)
    }
}
