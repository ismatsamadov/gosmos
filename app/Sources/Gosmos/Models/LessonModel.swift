import CodeEditor
import Foundation
import GosmosKit
import Observation
import os

/// The lesson on screen: its content and parsed readme, plus the state its type needs, its
/// hints and its reference solution. A review's redo is a lesson too (`isReview`): the
/// starter code in a folder of its own, graded for the review only.
@Observable @MainActor
final class LessonModel {
    let detail: LessonDetail
    /// The readme and quiz questions as blocks.
    let parsed: ParsedLesson
    let isReview: Bool
    private(set) var code: CodeSession?
    private(set) var quiz: QuizSession?
    var note = "" { didSet { if noteLoaded && note != oldValue { scheduleNoteSave() } } }
    private var noteLoaded = false
    private(set) var noteLoadError: String?
    private var isLoadingNote = false
    var canEditNote: Bool { noteLoaded && !isReview }
    /// The note as the database holds it, so an unchanged note is never written.
    @ObservationIgnored private var savedNote = ""
    private var noteTask: Task<Void, Never>?
    /// Set when the last Submit passed, so the toolbar offers the next lesson and the readme
    /// shows the completion panel; a Submit that fails clears it.
    var justPassed = false
    /// Set when a pass put the lesson in review, due tomorrow.
    private(set) var comesBackForReview = false
    /// Set when the last pass finished the lesson's chapter: it was the one required lesson of
    /// the chapter still to do, so the completion panel shows the chapter's debrief.
    private(set) var finishedChapter = false
    /// Bumped to bring the part below the readme (the quiz) into view.
    private(set) var footerRequests = 0
    /// Which hints show, and whether the next may.
    private(set) var hints: HintState
    /// The reference solution beside your code, while you look at it.
    var showingSolution = false
    /// You looked at the solution before passing: the pass brings the lesson back for review.
    private(set) var solutionWasViewed = false
    @ObservationIgnored weak var app: AppModel?
    @ObservationIgnored let learner: LearnerDatabase

    /// `saved` is your code from code_file, so a code lesson opens with it rather than the
    /// starter code.
    init(detail: LessonDetail, parsed: ParsedLesson, saved: [String: CodeFileRecord], learner: LearnerDatabase, review: Bool, app: AppModel) {
        self.detail = detail
        self.parsed = parsed
        self.learner = learner
        self.isReview = review
        self.app = app
        hints = HintState(hints: review ? [] : detail.hints)
        switch detail.content {
        case .code(let c): code = CodeSession(lesson: detail, exercise: c, tests: false, saved: saved, learner: learner, review: review, app: app)
        case .codeTests(let c): code = CodeSession(lesson: detail, exercise: c, tests: true, saved: saved, learner: learner, review: review, app: app)
        case .choice(let q): quiz = QuizSession(lesson: detail, quiz: q, learner: learner, app: app)
        }
    }

    var ref: LessonRef { detail.ref }
    var isDone: Bool { app?.isDone(ref) ?? false }

    /// The attempts whose failure unlocks the last hint: Submits, or a quiz's answers.
    var submitKind: AttemptKind { quiz != nil ? .answer : .submit }

    /// Reads what this lesson needs from the learner database and the work folder.
    func load() async {
        await loadNote()
        do {
            if !hints.hints.isEmpty {
                hints = HintState(hints: detail.hints, revealed: try await learner.revealedHints(ref),
                                  failedSubmits: try await learner.failedAttempts(ref, kind: submitKind))
            }
            if hasSolution { solutionWasViewed = try await learner.viewedSolution(ref) }
        } catch {
            app?.report(error)
        }
        await code?.load()
    }

    /// Notes wait for their saved text before taking edits, just as code does. A failed read
    /// leaves editing off and can be retried without reloading the code or replacing drafts.
    func loadNote() async {
        guard !isReview, !noteLoaded, !isLoadingNote else { return }
        isLoadingNote = true
        defer { isLoadingNote = false }
        do {
            note = try await learner.note(for: ref) ?? ""
            savedNote = note
            noteLoaded = true
            noteLoadError = nil
        } catch {
            noteLoadError = error.localizedDescription
            app?.report(error)
        }
    }

    /// Saves pending edits; called before another lesson opens and when the app quits. A code
    /// conflict is asked about first; `limit` bounds the wait when quitting. Saving runs to
    /// the end even when the caller was cancelled (a newer lesson replaced it): in a cancelled
    /// task the database refuses to write.
    func leave(within limit: Duration? = nil) async {
        await Task {
            noteTask?.cancel()
            await saveNote()
            await code?.settle(within: limit)
            code?.stop()
        }.value
    }

    private func scheduleNoteSave() {
        noteTask?.cancel()
        noteTask = Task { [weak self] in
            try? await Task.sleep(for: .seconds(1))
            guard !Task.isCancelled else { return }
            await self?.saveNote()
        }
    }

    /// Writes the note if it changed since it was loaded or last saved.
    private func saveNote() async {
        guard noteLoaded, note != savedNote else { return }
        let body = note
        do {
            try await learner.saveNote(ref, body: body)
            savedNote = body
        } catch {
            app?.report(error)
        }
    }

    // MARK: Hints and the solution

    /// Shows the next hint, when it may come, and keeps it shown from now on.
    func revealHint() {
        guard let hint = hints.reveal() else { return }
        let learner = learner, ref = ref, position = hint.position
        learner.pendingWrites.start { [weak app] in
            do { try await learner.revealHint(ref, position: position) } catch { app?.report(error) }
        }
    }

    /// A graded attempt failed: the last hint may come now.
    func noteFailedSubmit() {
        hints.failedSubmits += 1
        justPassed = false
    }

    /// Whether the lesson has a reference solution to compare with.
    var hasSolution: Bool { !isReview && code?.exercise.solutionFiles.isEmpty == false }

    /// Shows the reference solution beside your code. Before the lesson is done, the look is
    /// recorded, and the pass that follows brings the lesson back for review; after, it is free.
    func showSolution() {
        guard hasSolution else { return }
        showingSolution = true
        let learner = learner, ref = ref
        learner.pendingWrites.start { [weak self, weak app] in
            do {
                if try await learner.recordSolutionView(ref) { self?.solutionWasViewed = true }
            } catch {
                app?.report(error)
            }
        }
    }

    /// Run ⌘R: runs code.
    func run() {
        code?.run(submit: false)
    }

    /// Submit ⌘↩: submits code, or checks a quiz.
    func submit() {
        if let code {
            code.run(submit: true)
        } else if let quiz {
            footerRequests += 1
            quiz.submit()
        }
    }

    /// Brings the part below the readme into view.
    func showFooter() {
        footerRequests += 1
    }

    /// Code in another language opens and reads as usual, but does not run.
    var canRun: Bool { code?.exercise.runsHere == true }
    var canSubmit: Bool { code?.exercise.runsHere == true || quiz?.canSubmit == true }
    var isBusy: Bool { code?.isRunning == true }

    /// The lesson passed: done (the learner database recorded it with the attempt), and the
    /// Next Lesson button. A redo for review leaves the done mark alone.
    func markPassed(review scheduled: Bool) {
        justPassed = true
        guard !isReview else { return }
        if scheduled { comesBackForReview = true }
        finishedChapter = !isDone && (app?.isLastToDo(ref) ?? false)
        app?.setDone(ref, true)
    }

    /// Whether the Editor menu's `command` can act now: on the file in the editor, with no
    /// sheet or question over it.
    func canPerform(_ command: EditorCommand) -> Bool {
        guard let code, !showingSolution, code.conflict == nil else { return false }
        return code.whyUnavailable(command) == nil
    }
}

/// The Editor menu's commands, which the buttons at the end of the editor's file bar give too.
enum EditorCommand {
    case format, toggleComment, indent, outdent
}

/// A code lesson's files, saves and runs. The files, their loading and saving live in
/// `buffers` (GosmosKit), where they are tested.
@Observable @MainActor
final class CodeSession {
    let buffers: CodeBuffers
    let tests: Bool
    private(set) var isRunning = false
    private(set) var result: RunResult?
    private(set) var resultWasSubmit = false
    /// Counts passes, so the done mark can bounce.
    private(set) var passCount = 0
    /// Why the last Format changed nothing, such as gofmt's messages for a syntax error. The
    /// console shows it in place of the last run until the next Format or run.
    private(set) var formatFailure: RunResult?
    /// The active file's editor, which the Editor menu and the buttons above it edit through.
    let editor = CodeEditorProxy()
    @ObservationIgnored private weak var app: AppModel?
    @ObservationIgnored private var runTask: Task<Void, Never>?
    @ObservationIgnored private var formatTask: Task<Void, Never>?
    /// A Run (false) or Submit (true) waiting for a conflict to be settled.
    @ObservationIgnored private var pendingRun: Bool?
    @ObservationIgnored private let learner: LearnerDatabase
    /// A redo for review: the starter code in a folder of its own, saved nowhere.
    let isReview: Bool
    /// A redo's first Submit grades the review; later ones are practice.
    @ObservationIgnored private var reviewGraded = false
    /// Whether a redo's graded Submit passed; nil until it has run.
    private(set) var reviewPassed: Bool?

    init(lesson: LessonDetail, exercise: CodeExercise, tests: Bool, saved: [String: CodeFileRecord], learner: LearnerDatabase, review: Bool, app: AppModel) {
        self.tests = tests
        self.app = app
        self.learner = learner
        self.isReview = review
        let ref = lesson.ref
        buffers = CodeBuffers(lesson: lesson, exercise: exercise, saved: saved, learner: review ? nil : learner,
                              workspace: review ? app.workspace.review : app.workspace,
                              unsettled: review ? [] : app.preferences.unsettledFiles(in: ref),
                              goVersion: { [weak app] in await app?.awaitToolchain()?.languageVersion ?? "1.27" },
                              remember: { [weak app] name, unsettled in app?.preferences.setUnsettled(name, in: ref, unsettled) })
    }

    var lesson: LessonDetail { buffers.lesson }
    var exercise: CodeExercise { buffers.exercise }
    var files: [FileBuffer] { buffers.files }
    var activeIndex: Int {
        get { buffers.activeIndex }
        set { buffers.activeIndex = newValue }
    }
    var activeFile: FileBuffer? { files.indices.contains(activeIndex) ? files[activeIndex] : nil }
    /// True once the editor holds your code as code_file and the lesson folder have it, and
    /// saving is ready: the editor takes typing from then on. Until then it shows a load
    /// under way, or `loadError` when the load failed.
    var isLoaded: Bool { buffers.isLoaded }
    var loadError: String? { buffers.loadError }
    var folder: URL? { buffers.folder }
    /// A file that changed both here and in another editor; the lesson asks which to keep.
    var conflict: String? { buffers.conflict }

    /// Your code from the database, reconciled with the work folder (an edit made there in
    /// another editor wins when it is newer). Then a save, which asks about a file left in
    /// conflict when the app last quit.
    func load() async {
        await buffers.load()
        if buffers.isLoaded { await flush() }
    }

    /// Saves your edits soon; the editor calls this a second after typing stops, and the app
    /// when it comes to the front (to pick up edits made in another editor).
    func save() {
        Task { await flush() }
    }

    /// Saves now (see `CodeBuffers.flush`). Runs, lesson switches and quitting call this
    /// first. Returns false while a conflict waits for your answer or the save failed.
    @discardableResult
    func flush() async -> Bool {
        do {
            return try await buffers.flush()
        } catch {
            app?.reportFiles(error)
            return false
        }
    }

    /// Settles a conflict with your version or the one on disk (read again now), then asks
    /// about the next one, or runs what was waiting.
    func resolveConflict(keepMine: Bool) {
        guard let resolving = buffers.resolveConflict(keepMine: keepMine) else { return }
        Task {
            do {
                guard try await resolving.value else { return }
            } catch {
                app?.reportFiles(error)
                return
            }
            if let submit = pendingRun {
                pendingRun = nil
                run(submit: submit)
            }
        }
    }

    /// Saves before the lesson closes, after a load under way (see `CodeBuffers.settle`).
    func settle(within limit: Duration?) async {
        do {
            let kept = try await buffers.settle(within: limit)
            if !kept.isEmpty {
                app?.report(AppError.files("Kept both versions of \(kept.joined(separator: ", ")): yours is saved in Gosmos, the other editor's stays in the lesson folder. Gosmos asks which to keep when you open the lesson."))
            }
        } catch {
            app?.reportFiles(error)
        }
    }

    func run(submit: Bool) {
        guard !isRunning, exercise.runsHere, let app else { return }
        isRunning = true
        result = nil
        formatFailure = nil
        resultWasSubmit = submit
        let lesson = lesson, exercise = exercise, tests = tests, scratch = app.workspace.runsFolder
        runTask = Task { [weak self] in
            guard let self else { return }
            if self.folder == nil { await self.load() }
            // Saving waits for a save under way and takes in other editors' edits, so the
            // build sees the code on screen.
            guard await self.flush() else {
                self.isRunning = false
                if self.conflict != nil { self.pendingRun = submit }
                return
            }
            guard let folder = self.folder else {
                self.isRunning = false
                self.result = RunResult(status: .error, output: self.loadError ?? "The lesson folder could not be prepared.")
                return
            }
            guard let toolchain = await app.awaitToolchain() else { return }
            let result = await GoRunner(toolchain: toolchain, scratch: scratch).run(exercise, tests: tests, folder: folder, submit: submit)
            self.result = result
            // The run is newer than a Format that failed while it ran.
            self.formatFailure = nil
            self.isRunning = false
            let passed = submit ? result.status == .passed : result.status.isSuccess
            var output = result.output
            if let expected = result.expected { output += "\n\n--- expected output ---\n" + expected }
            // Stopped, or it could not run at all (no Go, say): no attempt of yours, so it
            // counts as no failed Submit, and a redo's grading waits for one that runs.
            guard result.status != .cancelled, result.status != .error else { return }
            // The record is written to the end, even when this run's task is cancelled meanwhile
            // (another lesson opened), and quitting waits for it.
            let learner = self.learner, ref = lesson.ref, recorded = output, ms = result.durationMs
            var scheduled = false
            do {
                if self.isReview {
                    // The redo's first Submit grades the review; it never touches your code or done mark.
                    if submit && !self.reviewGraded {
                        self.reviewGraded = true
                        self.reviewPassed = passed
                        try await learner.pendingWrites.start {
                            try await learner.recordLessonReview(ref, passed: passed, output: recorded, durationMs: ms)
                        }.value
                        await app.refreshDueCount()
                    }
                } else if submit {
                    scheduled = try await learner.pendingWrites.start {
                        try await learner.recordSubmit(ref, passed: passed, output: recorded, durationMs: ms, reviewable: true)
                    }.value
                } else {
                    try await learner.pendingWrites.start {
                        try await learner.addAttempt(ref, kind: .run, passed: passed, output: recorded, durationMs: ms)
                    }.value
                }
            } catch {
                app.report(error)
            }
            // Only a passing Submit finishes the lesson; a Run shows the output (or the tests'
            // results) and leaves the lesson as it was.
            if submit && result.status == .passed {
                self.passCount += 1
                self.ownLesson?.markPassed(review: scheduled)
                if scheduled { await app.refreshDueCount() }
            } else if submit {
                self.ownLesson?.noteFailedSubmit()
            }
        }
    }

    /// The lesson this session belongs to, while it is on screen. A run can end after another
    /// lesson opened, whose done mark and hints are not its to change.
    private var ownLesson: LessonModel? { app?.lesson.flatMap { $0.code === self ? $0 : nil } }

    /// Stops a run under way; its process group is killed. A Format under way stops too.
    func stop() {
        runTask?.cancel()
        formatTask?.cancel()
    }

    /// Puts the starter code back into every file of yours.
    func resetToStarter() {
        for buffer in files where !buffer.readOnly {
            if let starter = exercise.starterFiles.first(where: { $0.name == buffer.name }) {
                buffer.text = starter.content
            }
        }
        save()
    }

    // MARK: Editor commands

    /// Why `command` can't act on the active file now, as its button's help says; nil when it
    /// can. Format waits for Go while Go is still being looked for, as a run does.
    func whyUnavailable(_ command: EditorCommand) -> String? {
        guard let file = activeFile else { return "There is no file to edit" }
        if file.readOnly { return "This file is read-only" }
        if !isLoaded { return loadError == nil ? "Your code is still loading" : "Your code could not load" }
        switch command {
        case .format:
            if !file.name.hasSuffix(".go") { return "Format works on .go files" }
            if let toolchain = app?.toolchain {
                if toolchain.go == nil { return "Go can't run: see Settings" }
                if toolchain.gofmt == nil { return "This Go has no gofmt: choose another in Settings" }
            }
        case .toggleComment:
            if LineSyntax(fileName: file.name).commentMarker == nil { return "This kind of file has no comments" }
        case .indent, .outdent:
            break
        }
        return nil
    }

    /// Carries out an editor command on the active file, if it can act now. Each change is one
    /// step on the editor's undo stack.
    func perform(_ command: EditorCommand) {
        guard whyUnavailable(command) == nil else { return }
        switch command {
        case .format: format()
        case .toggleComment: editor.toggleComment()
        case .indent: editor.indent()
        case .outdent: editor.outdent()
        }
    }

    /// Puts the active file through gofmt, of the Go your code runs with, and its output into
    /// the editor; formatted code stays as it is. When gofmt refuses the code, such as for a
    /// syntax error, the code stays too and the console shows why.
    private func format() {
        guard let app, let file = activeFile, let source = editor.text else { return }
        formatTask?.cancel()
        formatTask = Task { [weak self] in
            guard let toolchain = await app.awaitToolchain() else { return }
            let started = ContinuousClock.now
            let outcome: Result<String, GoFormatError>
            do {
                guard let formatter = GoFormatter(toolchain: toolchain) else {
                    throw GoFormatError.failed(toolchain.go == nil ? toolchain.problem ?? "Go can't run: see Settings."
                                                                   : "This Go has no gofmt: choose another in Settings.")
                }
                outcome = .success(try await formatter.format(source, fileName: file.name))
            } catch let error as GoFormatError {
                outcome = .failure(error)
            } catch {
                return // stopped: another lesson opened, or another Format began
            }
            // Only for the code it was made from: not once you typed, or another file is open.
            guard let self, !Task.isCancelled, self.activeFile === file, self.editor.text == source else { return }
            switch outcome {
            case .success(let formatted):
                self.formatFailure = nil
                self.editor.replaceText(with: formatted, actionName: "Format")
            case .failure(let error):
                let status: RunStatus = if case .syntax = error { .failed } else { .error }
                self.formatFailure = RunResult(status: status, output: error.errorDescription ?? "", duration: ContinuousClock.now - started)
            }
        }
    }
}
