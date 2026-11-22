import Foundation
import Testing

/// Gosmos never uses the network: it is self-contained, and everything a lesson needs is on the
/// Mac. These tests read every Swift file the app is built from, of app/Sources, of its own
/// packages and of the third-party ones kept in app/Packages/Vendor, and fail on anything that
/// would reach outside it:
///   - a networking framework or library imported: Network, SwiftNIO, FoundationNetworking,
///     CFNetwork and the like;
///   - a networking API used: URLSession, Network.framework's connections, listeners and path
///     monitors, CFNetwork's streams, sockets and name lookups;
///   - a web address written into the code: a string literal that starts with http:// or
///     https://, which is how a link, an NSWorkspace open or a request would come by one;
///   - a link handed to the system (OpenURLAction's systemAction), which opens a web address in
///     the browser. Links act in the app or not at all.
/// Comments may say anything. Opening a local file or folder, such as the lesson folder in your
/// editor with NSWorkspace, stays allowed.
@Suite struct IsolationTests {
    struct Rule {
        let what: String
        let pattern: Regex<AnyRegexOutput>

        init(_ what: String, _ pattern: String) {
            self.what = what
            // \b between a word character and another: Unicode's word boundaries would take
            // URLSession.shared for one word.
            self.pattern = try! Regex(pattern).wordBoundaryKind(.simple)
        }
    }

    static func rules() -> [Rule] { [
        Rule("imports a networking module", #"\bimport\s+(?:(?:struct|class|enum|protocol|func|var|let|typealias)\s+)?(?:Network|NIO\w*|FoundationNetworking|CFNetwork|AsyncHTTPClient|WebSocketKit|MultipeerConnectivity)\b"#),
        Rule("uses a networking API", #"\b(?:URLSession\w*|NSURLSession\w*|NSURLConnection|URLSessionWebSocketTask|NWConnection\w*|NWListener|NWBrowser|NWPathMonitor|NWEndpoint|NWParameters|NWProtocol\w*|CFSocket\w*|CFStreamCreatePairWithSocket\w*|CFReadStreamCreateForHTTPRequest|CFHTTPMessage\w*|getaddrinfo|gethostbyname)\b|\bsocket\("#),
        Rule("writes a web address into the code", #""https?://"#),
        Rule("hands a link to the system, which opens a web address in the browser", #"\.systemAction\b"#),
    ] }

    /// The Swift files of the app: app/Sources, the sources of app/Packages, and the third-party
    /// packages in app/Packages/Vendor, whole.
    static func appSources() throws -> [URL] {
        let app = Fixtures.repository.appendingPathComponent("app")
        var roots = [app.appendingPathComponent("Sources")]
        for package in try FileManager.default.contentsOfDirectory(at: app.appendingPathComponent("Packages"), includingPropertiesForKeys: nil) {
            roots.append(package.appendingPathComponent("Sources"))
        }
        roots += try vendored()
        var files: [URL] = []
        for root in roots {
            guard let walk = FileManager.default.enumerator(at: root, includingPropertiesForKeys: [.isRegularFileKey]) else { continue }
            for case let file as URL in walk where file.pathExtension == "swift" {
                // GRDB.swift is a folder.
                if (try? file.resourceValues(forKeys: [.isRegularFileKey]))?.isRegularFile == true { files.append(file) }
            }
        }
        return files
    }

    /// The third-party packages kept in the repository: app/Packages/Vendor's folders.
    static func vendored() throws -> [URL] {
        let vendor = Fixtures.repository.appendingPathComponent("app/Packages/Vendor")
        return try FileManager.default.contentsOfDirectory(at: vendor, includingPropertiesForKeys: nil).filter(\.hasDirectoryPath)
    }

    /// A line of Swift without its comment: the text before a `//` that is outside a string.
    static func code(_ line: Substring) -> Substring {
        var inString = false, escaped = false
        var i = line.startIndex
        while i < line.endIndex {
            let c = line[i]
            if escaped {
                escaped = false
            } else if c == "\\" && inString {
                escaped = true
            } else if c == "\"" {
                inString.toggle()
            } else if c == "/" && !inString && line[line.index(after: i)...].hasPrefix("/") {
                return line[..<i]
            }
            i = line.index(after: i)
        }
        return line
    }

    @Test func theAppNeverUsesTheNetwork() throws {
        let files = try Self.appSources(), rules = Self.rules()
        #expect(files.count > 40, "found only \(files.count) Swift files of the app")
        var found: [String] = []
        for file in files {
            let text = try String(contentsOf: file, encoding: .utf8)
            for (n, line) in text.split(separator: "\n", omittingEmptySubsequences: false).enumerated() {
                let code = Self.code(line)
                for rule in rules where code.contains(rule.pattern) {
                    let path = file.path.replacingOccurrences(of: Fixtures.repository.path + "/", with: "")
                    found.append("\(path):\(n + 1) \(rule.what): \(line.trimmingCharacters(in: .whitespaces))")
                }
            }
        }
        #expect(found.isEmpty, "Gosmos never uses the network:\n\(found.joined(separator: "\n"))")
    }

    /// The app fetches no package, a networking library least of all: every package it builds
    /// with is a folder of the repository, whose code the test above reads too.
    @Test func everyPackageIsInTheRepository() throws {
        let app = Fixtures.repository.appendingPathComponent("app")
        var manifests = [app.appendingPathComponent("Package.swift")]
        for package in try FileManager.default.contentsOfDirectory(at: app.appendingPathComponent("Packages"), includingPropertiesForKeys: nil) + Self.vendored() {
            let manifest = package.appendingPathComponent("Package.swift")
            if FileManager.default.fileExists(atPath: manifest.path) { manifests.append(manifest) }
        }
        #expect(manifests.count == 5, "found the manifests \(manifests.map(\.path))")
        for manifest in manifests {
            let text = try String(contentsOf: manifest, encoding: .utf8)
            #expect(!text.contains(".package(url:"), "\(manifest.path) fetches a package")
        }
        #expect(!FileManager.default.fileExists(atPath: app.appendingPathComponent("Package.resolved").path), "app/Package.resolved pins a fetched package")
    }

    /// The rules catch what they are for, and let comments and local files be.
    @Test func theRulesCatchWhatTheyAreFor() {
        let rules = Self.rules()
        func flagged(_ line: Substring) -> Bool { rules.contains { Self.code(line).contains($0.pattern) } }
        for line: Substring in [
            "import Network",
            "@preconcurrency import NIOCore",
            "let task = URLSession.shared.dataTask(with: url)",
            "let monitor = NWPathMonitor()",
            #"NSWorkspace.shared.open(URL(string: "https://go.dev/dl/")!)"#,
            ##"Link("Download Go", destination: URL(string: #"http://example.com"#)!)"##,
            "OpenURLAction { _ in .systemAction }",
        ] {
            #expect(flagged(line), "\(line)")
        }
        for line: Substring in [
            "// Install it from https://go.dev/dl/, or use URLSession.",
            #"let message = "See the docs" // https://example.com"#,
            "NSWorkspace.shared.activateFileViewerSelecting([folder])",
            "_ = try await NSWorkspace.shared.open(targets, withApplicationAt: editor, configuration: config)",
            #"let scheme = url.scheme == "https""#,
            "OpenURLAction { url in app.follow(url) ? .handled : .discarded }",
        ] {
            #expect(!flagged(line), "\(line)")
        }
    }
}
