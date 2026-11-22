import AppKit
import GosmosKit

extension AppModel {
    /// The folder of the lesson on screen: a code lesson's work folder.
    var lessonFolder: URL? { lesson?.code?.folder }

    func revealLessonFolder() {
        Task {
            await lesson?.code?.flush()
            guard let folder = lessonFolder else { return }
            let main = lesson?.code?.activeFile.map { folder.appendingPathComponent($0.name) }
            NSWorkspace.shared.activateFileViewerSelecting([main.flatMap { FileManager.default.fileExists(atPath: $0.path) ? $0 : nil } ?? folder])
        }
    }

    /// Opens the folder in the app that opens your Go files; an editor that only opens
    /// files gets the file you are editing.
    func openLessonFolderInEditor() {
        Task {
            await lesson?.code?.flush()
            guard let folder = lessonFolder else { return }
            let file = lesson?.code?.activeFile.map { folder.appendingPathComponent($0.name) }
            let probe = file ?? folder.appendingPathComponent("main.go")
            let config = NSWorkspace.OpenConfiguration()
            guard let editor = NSWorkspace.shared.urlForApplication(toOpen: probe) else {
                NSWorkspace.shared.activateFileViewerSelecting([folder])
                return
            }
            let fileOnly = ["com.apple.TextEdit", "com.apple.dt.Xcode"].contains(Bundle(url: editor)?.bundleIdentifier ?? "")
            let targets = fileOnly ? [file ?? probe] : [folder]
            do {
                _ = try await NSWorkspace.shared.open(targets, withApplicationAt: editor, configuration: config)
            } catch {
                banner = .files(error.localizedDescription)
            }
        }
    }
}
