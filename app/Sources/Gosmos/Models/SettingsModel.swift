import Foundation
import GosmosKit
import Observation

/// Settings › Editor and Settings › Go: the editor's font size and line numbers, kept in the
/// learner database, and the Go that runs your code, chosen on this Mac.
@Observable @MainActor
final class SettingsModel {
    private(set) var editorFontSize: Double = 13
    private(set) var showLineNumbers = true
    private(set) var toolchain: Toolchain?
    let preferences: LocalPreferences
    /// The app whose learner database keeps the editor's settings; it sets itself here when it
    /// is made.
    @ObservationIgnored weak var app: AppModel?
    /// Finds the Go to run code with; only the latest check counts (Settings › Go).
    private let toolchainWork = LatestWork<Toolchain>()

    init(preferences: LocalPreferences) {
        self.preferences = preferences
    }

    /// Reads the editor's settings from the learner database.
    func load(from learner: LearnerDatabase) async throws {
        editorFontSize = Double(try await learner.setting("editor_font_size") ?? "") ?? 13
        showLineNumbers = (try await learner.setting("editor_line_numbers")) != "0"
    }

    /// Starts the first check of the Go to run code with, as the app launches.
    func discoverToolchain() {
        let goOverride = preferences.goPath, goCaches = goCaches
        toolchainWork.start(priority: .utility, { await Toolchain.discover(goOverride: goOverride, caches: goCaches) }) { [weak self] in
            self?.toolchain = $0
        }
    }

    func setEditorFontSize(_ size: Double) {
        editorFontSize = min(max(size, 9), 28)
        guard let learner = app?.learner else { return }
        let value = String(Int(editorFontSize))
        learner.pendingWrites.start { [weak app] in
            do { try await learner.setSetting("editor_font_size", value) } catch { app?.report(error) }
        }
    }

    /// The Go to run code with, waiting for the latest check of it.
    func awaitToolchain() async -> Toolchain? {
        if let toolchain { return toolchain }
        guard let t = await toolchainWork.value() else { return nil }
        toolchain = t
        return t
    }

    /// Where Go keeps its caches: ~/Library/Caches/Gosmos, beside the lesson folders, or the
    /// data folder of a scripted run.
    private var goCaches: URL { preferences.dataFolder ?? Toolchain.defaultCaches }

    /// Runs your code with the Go at `goPath`, or with the built-in Go when it is nil.
    func rediscoverToolchain(goPath: String?) {
        preferences.setGoPath(goPath)
        toolchain = nil
        let override = preferences.goPath, goCaches = goCaches
        // A check still under way, such as of a Go that hangs, never lands after this one.
        toolchainWork.start(priority: .userInitiated, { await Toolchain.discover(goOverride: override, caches: goCaches) }) { [weak self] in
            self?.toolchain = $0
        }
    }

    func setShowLineNumbers(_ on: Bool) {
        showLineNumbers = on
        guard let learner = app?.learner else { return }
        learner.pendingWrites.start { [weak app] in
            do { try await learner.setSetting("editor_line_numbers", on ? "1" : "0") } catch { app?.report(error) }
        }
    }
}

// MARK: - Settings

extension AppModel {
    var editorFontSize: Double { settings.editorFontSize }
    var showLineNumbers: Bool { settings.showLineNumbers }
    var toolchain: Toolchain? { settings.toolchain }
    func setEditorFontSize(_ size: Double) { settings.setEditorFontSize(size) }
    func setShowLineNumbers(_ on: Bool) { settings.setShowLineNumbers(on) }
    /// The Go to run code with, waiting for the latest check of it.
    func awaitToolchain() async -> Toolchain? { await settings.awaitToolchain() }
    func rediscoverToolchain(goPath: String?) { settings.rediscoverToolchain(goPath: goPath) }
}
