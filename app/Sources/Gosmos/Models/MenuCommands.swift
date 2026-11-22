import GosmosKit

/// When the menus' commands can act, worked out from the app's state, so the menus and the
/// flow tests read the same rules. While the opening covers the window, nothing acts on what
/// it covers: the lesson behind it, Back and Forward, the review, Mission Control.
@MainActor
struct MenuCommands {
    let app: AppModel

    /// The opening covers the window.
    var covered: Bool { app.opening != nil }
    /// The lesson the commands act on: none while the opening covers it.
    private var lesson: LessonModel? { covered ? nil : app.lesson }

    var run: Bool { lesson.map { $0.canRun && !$0.isBusy } ?? false }
    var submit: Bool { lesson.map { $0.canSubmit && !$0.isBusy } ?? false }
    /// Stop ends a run under way, whatever covers it.
    var stop: Bool { app.lesson?.isBusy ?? false }
    var nextLesson: Bool { !covered && app.nextLesson != nil }
    var previousLesson: Bool { !covered && app.previousLesson != nil }
    var back: Bool { !covered && app.canGoBack }
    var forward: Bool { !covered && app.canGoForward }
    var showHint: Bool { lesson?.hints.canReveal ?? false }
    var showSolution: Bool { lesson?.hasSolution ?? false }
    var review: Bool { !covered && app.learner != nil }
    /// Mark as Done, or as Not Done: a lesson's, never a redo's.
    var markDone: Bool { lesson.map { !$0.isReview } ?? false }
    var resetToStarter: Bool { lesson.map { $0.code != nil && !$0.isReview } ?? false }
    func editor(_ command: EditorCommand) -> Bool { lesson?.canPerform(command) ?? false }
    var showDocumentation: Bool { !covered && app.canShowQuickHelp }
    var missionControl: Bool { !covered && app.learner != nil }
    var fieldGuide: Bool { !covered && app.learner != nil && app.hasFieldGuide }
    var replayOpening: Bool { !covered && app.learner != nil && app.openingToReplay != nil }
}

extension AppModel {
    /// When the menus' commands can act now.
    var commands: MenuCommands { MenuCommands(app: self) }
}
