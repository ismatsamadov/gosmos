import Foundation
import GosmosKit

extension AppModel {
    /// Launch arguments for scripted runs (screenshots, debugging): `-GosmosOpenLesson
    /// [<library>/]<lesson>` opens a lesson, and `-GosmosLaunchAction <action>` acts on it once,
    /// after it loads. They live in the argument domain only, never stored.
    enum LaunchArguments {
        static let openLesson = "GosmosOpenLesson"
        static let action = "GosmosLaunchAction"
    }

    /// Acts once on the first lesson that loads, as `-GosmosLaunchAction` asks.
    func performLaunchAction(on model: LessonModel) {
        guard !launchActionDone, let action = preferences.defaults.string(forKey: LaunchArguments.action) else { return }
        launchActionDone = true
        if action.hasPrefix("search:") {
            // Types the query a letter at a time, as someone would.
            let query = String(action.dropFirst(7))
            Task { [weak self] in
                for end in query.indices {
                    try? await Task.sleep(for: .milliseconds(120))
                    self?.search.text = String(query[...end])
                }
            }
            return
        }
        if action.hasPrefix("doc") || action.hasPrefix("quick-help:") {
            Task { [weak self] in await self?.performDocumentationAction(action, on: model) }
            return
        }
        switch action {
        case "settings": settingsRequests += 1
        case "run": model.run()
        case "submit": model.submit()
        case "footer": model.showFooter()
        case "hints":
            // Every hint that may show: all but a last one that waits for a failed Submit.
            while model.hints.canReveal { model.revealHint() }
        case "solution": model.showSolution()
        case "format", "comment":
            guard let code = model.code else { return }
            Task {
                // The editor attaches itself to the session once it is on screen.
                for _ in 0..<50 {
                    if code.editor.text != nil { break }
                    try? await Task.sleep(for: .milliseconds(100))
                }
                code.perform(action == "format" ? .format : .toggleComment)
            }
        case "review": showReview()
        case "redo-submit": model.submit()
        case _ where action.hasPrefix("term:"):
            shownTerm = story(for: model.ref)?.term(String(action.dropFirst(5)))
        case "answer":
            // The first question right, the others wrong: both outcomes on screen.
            guard let quiz = model.quiz else { return }
            for (i, q) in quiz.quiz.questions.enumerated() {
                let right = q.answers.firstIndex(of: q.answer) ?? 0
                let wrong = q.answers.indices.first { $0 != right } ?? right
                quiz.pick(i == 0 ? right : wrong, for: i)
            }
            model.submit()
        default: break
        }
    }

    /// The documentation's launch actions, once it has opened: `doc:<link>` shows a place in
    /// it (`doc:godoc:fmt.Println`, `doc:acknowledgements`), `doc-search:<text>` types a search
    /// in the browser, and `quick-help:<text>` shows Quick Help for the first `<text>` in the
    /// editor's file.
    private func performDocumentationAction(_ action: String, on model: LessonModel) async {
        for _ in 0..<50 where docs == nil || (action.hasPrefix("quick-help:") && model.code?.editor.text == nil) {
            try? await Task.sleep(for: .milliseconds(100))
        }
        guard let docs else { return }
        let argument = String(action.drop { $0 != ":" }.dropFirst())
        if action.hasPrefix("doc:") {
            openDocumentation(argument == "acknowledgements" ? .acknowledgements : DocLink(string: argument))
        } else if action.hasPrefix("doc-search:") {
            openDocumentation()
            for end in argument.indices {
                try? await Task.sleep(for: .milliseconds(120))
                docs.searchText = String(argument[...end])
            }
        } else if let code = model.code, let text = code.editor.text {
            let found = (text as NSString).range(of: argument)
            if found.location != NSNotFound { showQuickHelp(code, at: found.location) }
        }
    }

    /// The lesson a scripted run opens at launch: the one it asks for, else the one you opened
    /// last in any library, else the first one. A launch without one shows Mission Control.
    func launchLesson() -> LessonRef? {
        if let requested = preferences.defaults.string(forKey: LaunchArguments.openLesson) {
            let parts = requested.split(separator: "/").map(String.init)
            let ref = parts.count == 2 ? LessonRef(library: parts[0], lesson: parts[1])
                : shelf.lazy.compactMap { lib in lib.catalog.summary(requested).map { _ in LessonRef(library: lib.id, lesson: requested) } }.first
            if let ref, summary(ref) != nil { return ref }
        }
        if let last = resumeLesson, summary(last) != nil { return last }
        for library in shelf {
            if let first = library.catalog.firstLesson { return LessonRef(library: library.id, lesson: first.id) }
        }
        return nil
    }

    /// The story's launch actions: `opening` or `opening:<page>` shows the first opening on the
    /// shelf, and `field-guide` the field guide.
    func performStoryAction(_ action: String) {
        if action == "field-guide" {
            showFieldGuide(remember: false)
            return
        }
        for library in shelf {
            for course in library.catalog.courses where library.stories[course.id]?.hasOpening == true {
                story.showOpening(course.id, in: library.id, page: Int(action.dropFirst("opening:".count)) ?? 0)
                return
            }
        }
    }
}
