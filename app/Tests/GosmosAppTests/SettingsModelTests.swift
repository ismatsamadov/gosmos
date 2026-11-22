import Foundation
import GosmosKit
import Testing
@testable import Gosmos

/// Settings › Editor and Settings › Go: the editor's settings, kept with your study data, and
/// the Go chosen on this Mac.
@Suite(.serialized) @MainActor struct SettingsModelTests {
    @Test func theEditorsSettingsStayInRangeAndAreKept() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let settings = app.settings
        #expect(settings.editorFontSize == 13 && settings.showLineNumbers)
        settings.setEditorFontSize(40)
        #expect(settings.editorFontSize == 28)
        settings.setEditorFontSize(4)
        #expect(settings.editorFontSize == 9)
        // Each change is a write of its own, and they may land in any order.
        await app.learner?.pendingWrites.finish()
        settings.setEditorFontSize(14.7)
        settings.setShowLineNumbers(false)
        await app.prepareToQuit()

        let reopened = fixture.app()
        await reopened.start()
        // Kept in whole points.
        #expect(reopened.settings.editorFontSize == 14 && !reopened.settings.showLineNumbers)
        await reopened.prepareToQuit()
    }

    /// The Go chosen is kept on this Mac and checked anew, with no Go to run until the check
    /// lands; its caches are the study data's folder of a scripted run.
    @Test func aChosenGoIsKeptAndChecked() async throws {
        let folder = FileManager.default.temporaryDirectory.appendingPathComponent("gosmos-settings-\(UUID().uuidString)")
        try FileManager.default.createDirectory(at: folder, withIntermediateDirectories: true)
        defer { try? FileManager.default.removeItem(at: folder) }
        let defaults = try #require(UserDefaults(suiteName: folder.appendingPathComponent("settings").path))
        let study = folder.appendingPathComponent("study")
        defaults.set(study.path, forKey: LocalPreferences.dataFolderKey)
        let settings = SettingsModel(preferences: LocalPreferences(defaults: defaults))

        let missing = folder.appendingPathComponent("go/bin/go").path
        settings.rediscoverToolchain(goPath: missing)
        #expect(settings.toolchain == nil && settings.preferences.goPath == missing)
        let chosen = try #require(await settings.awaitToolchain())
        #expect(chosen.go == nil && !chosen.isBuiltIn && chosen.problem?.contains(missing) == true)
        #expect(chosen.caches.path == study.path)

        settings.rediscoverToolchain(goPath: nil)
        #expect(settings.preferences.goPath == nil)
        #expect(await settings.awaitToolchain()?.isBuiltIn == true)
    }
}
