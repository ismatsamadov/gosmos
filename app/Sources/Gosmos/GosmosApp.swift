import AppKit
import GosmosKit
import SwiftUI

@main
struct GosmosApp: App {
    @NSApplicationDelegateAdaptor(AppDelegate.self) private var delegate

    private var app: AppModel { delegate.app }

    var body: some Scene {
        Window("Gosmos", id: "main") {
            RootView()
                .environment(app)
                .frame(minWidth: 1100, minHeight: 700)
                .tint(Theme.accent)
        }
        .defaultSize(width: 1440, height: 900)
        .windowToolbarStyle(.unified)
        .commands {
            GosmosCommands(app: app)
            // Find (⌘F, ⌘G, ⌘E) for the editor's find bar, and the other text commands.
            TextEditingCommands()
        }

        Settings {
            SettingsView()
                .environment(app)
                .tint(Theme.accent)
        }
    }
}

/// Owns the app's model and starts it when the app has launched, whether or not its window is
/// visible yet (a window that can't be seen, such as on a locked screen, runs no view tasks).
/// Saves pending work before the app quits.
@MainActor
final class AppDelegate: NSObject, NSApplicationDelegate {
    let app = AppModel()

    func applicationDidFinishLaunching(_ notification: Notification) {
        Log.app.info("Launched")
        // Another copy of this app runs and has the learner database: it comes forward, and
        // this one quits before it opens anything (SingleInstance).
        if let other = runningCopy() {
            Log.app.notice("Another copy of Gosmos runs, process \(other.processIdentifier, privacy: .public): quitting")
            NSApp.yieldActivation(to: other)
            other.activate()
            NSApp.terminate(nil)
            return
        }
        // Launch-only overrides for repeatable visual checks. They never change preferences.
        if let appearance = UserDefaults.standard.string(forKey: "GosmosAppearance") {
            switch appearance {
            case "light": NSApp.appearance = NSAppearance(named: .aqua)
            case "dark": NSApp.appearance = NSAppearance(named: .darkAqua)
            default: break
            }
        }
        Task {
            await app.start()
            if let size = UserDefaults.standard.string(forKey: "GosmosWindowSize") {
                // Exactly "<width>x<height>", two finite numbers; anything else is ignored.
                let parts = size.split(separator: "x", omittingEmptySubsequences: false).map { Double($0) }
                if parts.count == 2, let width = parts[0], let height = parts[1], width.isFinite, height.isFinite,
                   let window = NSApp.windows.first(where: { $0.canBecomeMain }) {
                    // The window's frame is autosaved; this run saves none, so a size for a
                    // layout check never becomes the size of the next launch.
                    _ = window.setFrameAutosaveName("")
                    window.setContentSize(NSSize(width: max(1100, width), height: max(700, height)))
                }
            }
        }
    }

    /// The copy of this app, by bundle id, that this launch gives way to: one launched before it.
    private func runningCopy() -> NSRunningApplication? {
        guard let id = Bundle.main.bundleIdentifier else { return nil }
        func described(_ r: NSRunningApplication) -> SingleInstance.App {
            SingleInstance.App(pid: r.processIdentifier, bundleIdentifier: r.bundleIdentifier, launchDate: r.launchDate)
        }
        let running = NSRunningApplication.runningApplications(withBundleIdentifier: id).filter { !$0.isTerminated }
        guard let copy = SingleInstance.runningCopy(for: described(.current), in: running.map(described)) else { return nil }
        return running.first { $0.processIdentifier == copy.pid }
    }

    func applicationShouldTerminate(_ sender: NSApplication) -> NSApplication.TerminateReply {
        // Saves, but quits within 20 seconds whatever happens.
        let reply = ReplyOnce { sender.reply(toApplicationShouldTerminate: true) }
        Task {
            await app.prepareToQuit()
            reply.send()
        }
        Task {
            try? await Task.sleep(for: .seconds(20))
            reply.send()
        }
        return .terminateLater
    }

    /// Back from another app, such as the editor the lesson folder is open in: take in edits
    /// made there, and count what has come due meanwhile.
    func applicationDidBecomeActive(_ notification: Notification) {
        app.lesson?.code?.save()
        Task { await app.refreshDueCount() }
    }

    func applicationShouldTerminateAfterLastWindowClosed(_ sender: NSApplication) -> Bool { true }
}

/// Calls its action the first time `send` is called.
@MainActor
final class ReplyOnce {
    private var action: (() -> Void)?
    init(_ action: @escaping () -> Void) { self.action = action }
    func send() {
        action?()
        action = nil
    }
}

struct GosmosCommands: Commands {
    let app: AppModel
    @FocusedValue(\.searchFocus) private var searchFocus

    var body: some Commands {
        CommandGroup(replacing: .newItem) {
            Button("Open Library…") { app.chooseLibrary() }
                .keyboardShortcut("o", modifiers: .command)
                .disabled(app.learner == nil || app.isAddingLibrary)
        }
        CommandGroup(after: .textEditing) {
            // ⇧⌘F, as Xcode's Find in Project, since ⌘F is the editor's Find.
            Button("Search") { searchFocus?.focus() }
                .keyboardShortcut("f", modifiers: [.command, .shift])
                .disabled(searchFocus == nil)
        }
        CommandMenu("Lesson") {
            Button("Run") { app.lesson?.run() }
                .keyboardShortcut("r", modifiers: .command)
                .disabled(!(app.lesson?.canRun ?? false) || (app.lesson?.isBusy ?? false))
            Button("Submit") { app.lesson?.submit() }
                .keyboardShortcut(.return, modifiers: .command)
                .disabled(!(app.lesson?.canSubmit ?? false) || (app.lesson?.isBusy ?? false))
            Button("Stop") { app.lesson?.code?.stop() }
                .keyboardShortcut(".", modifiers: .command)
                .disabled(!(app.lesson?.isBusy ?? false))
            Divider()
            Button("Next Lesson") { app.openNext() }
                .keyboardShortcut("]", modifiers: .command)
                .disabled(app.nextLesson == nil)
            Button("Previous Lesson") { app.openPrevious() }
                .keyboardShortcut("[", modifiers: .command)
                .disabled(app.previousLesson == nil)
            Button("Back") { app.goBack() }
                .keyboardShortcut(.leftArrow, modifiers: [.command, .option])
                .disabled(!app.canGoBack)
            Button("Forward") { app.goForward() }
                .keyboardShortcut(.rightArrow, modifiers: [.command, .option])
                .disabled(!app.canGoForward)
            Divider()
            // ⇧⌘H: the system keeps ⌥⌘H for Hide Others.
            Button("Show Hint") { app.lesson?.revealHint() }
                .keyboardShortcut("h", modifiers: [.command, .shift])
                .disabled(!(app.lesson?.hints.canReveal ?? false))
            Button("Show Solution") { app.lesson?.showSolution() }
                .keyboardShortcut("s", modifiers: [.command, .option])
                .disabled(!(app.lesson?.hasSolution ?? false))
            Button("Review") { app.showReview() }
                .keyboardShortcut("y", modifiers: [.command, .shift])
                .disabled(app.learner == nil)
            Divider()
            if let lesson = app.lesson, !lesson.isReview {
                let done = app.isDone(lesson.ref)
                Button(done ? "Mark as Not Done" : "Mark as Done") { app.setDone(lesson.ref, !done) }
                    .keyboardShortcut("d", modifiers: [.command, .shift])
            }
            Button("Reset to Starter Code…") { app.confirmReset = true }
                .disabled(app.lesson?.code == nil || app.lesson?.isReview == true)
            Divider()
            Button("Reveal Lesson Folder in Finder") { app.revealLessonFolder() }
                .keyboardShortcut("r", modifiers: [.command, .shift])
                .disabled(app.lessonFolder == nil)
            Button("Open Lesson Folder in Editor") { app.openLessonFolderInEditor() }
                .keyboardShortcut("o", modifiers: [.command, .shift])
                .disabled(app.lessonFolder == nil)
        }
        // The file in the editor; the buttons above it do the same.
        CommandMenu("Editor") {
            Button("Format Code") { app.lesson?.code?.perform(.format) }
                .keyboardShortcut("f", modifiers: [.shift, .option])
                .disabled(!(app.lesson?.canPerform(.format) ?? false))
            Button("Toggle Comment") { app.lesson?.code?.perform(.toggleComment) }
                .keyboardShortcut("/", modifiers: .command)
                .disabled(!(app.lesson?.canPerform(.toggleComment) ?? false))
            Divider()
            // With Option: ⌘] and ⌘[ are the next and previous lesson.
            Button("Indent") { app.lesson?.code?.perform(.indent) }
                .keyboardShortcut("]", modifiers: [.command, .option])
                .disabled(!(app.lesson?.canPerform(.indent) ?? false))
            Button("Outdent") { app.lesson?.code?.perform(.outdent) }
                .keyboardShortcut("[", modifiers: [.command, .option])
                .disabled(!(app.lesson?.canPerform(.outdent) ?? false))
            Divider()
            // Quick Help for the name at the insertion point, as Option-click gives it; ⌃⌘? as
            // in Xcode.
            Button("Show Documentation") { if let code = app.lesson?.code { app.showQuickHelp(code) } }
                .keyboardShortcut("?", modifiers: [.command, .control])
                .disabled(!app.canShowQuickHelp)
        }
        CommandGroup(after: .sidebar) {
            Button("Mission Control") { app.showDashboard() }
                .keyboardShortcut("1", modifiers: .command)
                .disabled(app.learner == nil)
            Divider()
            Button("Show Inspector") { app.showInspector.toggle() }
                .keyboardShortcut("i", modifiers: [.command, .option])
            Divider()
            Button("Bigger Code") { app.setEditorFontSize(app.editorFontSize + 1) }
                .keyboardShortcut("=", modifiers: .command)
            Button("Smaller Code") { app.setEditorFontSize(app.editorFontSize - 1) }
                .keyboardShortcut("-", modifiers: .command)
            Button("Actual Size") { app.setEditorFontSize(13) }
                .keyboardShortcut("0", modifiers: .command)
        }
    }
}
