import AppKit
import GosmosKit
import Observation
import os
import SwiftUI

/// A failure the app shows: as an inline banner, an alert, or as the window's content when
/// your study data can't be opened.
enum AppError: LocalizedError, Identifiable, Equatable {
    case database(AppDatabaseError)
    case learner(LearnerDatabaseError)
    case lesson(String)
    case files(String)

    var id: String { errorDescription ?? "error" }

    var errorDescription: String? {
        switch self {
        case .database(let e): e.errorDescription
        case .learner(let e): e.errorDescription
        case .lesson(let m): m
        case .files(let m): m
        }
    }

    var recoverySuggestion: String? {
        switch self {
        case .database(let e): e.recoverySuggestion
        case .learner(let e): e.recoverySuggestion
        case .lesson, .files: nil
        }
    }

    /// The headline for the error.
    var title: String {
        switch self {
        case .database(let e): e.title
        case .learner(let e): e.title
        case .lesson: "The Lesson Cannot Open"
        case .files: "Files Error"
        }
    }
}

/// Settings that belong to this Mac, in UserDefaults; your progress lives in the learner
/// database.
struct LocalPreferences {
    static let goPathKey = "GoPath"
    static let unsettledKey = "UnsettledCode"
    /// Scripted runs (screenshots, tests of the app) keep their study data apart:
    /// `-GosmosDataFolder <folder>` holds the learner database and the lesson folders. It lives
    /// in the argument domain only, never stored.
    static let dataFolderKey = "GosmosDataFolder"

    let defaults: UserDefaults

    var goPath: String? {
        let p = defaults.string(forKey: LocalPreferences.goPathKey)?.trimmingCharacters(in: .whitespaces)
        return (p?.isEmpty ?? true) ? nil : p
    }

    func setGoPath(_ path: String?) { defaults.set(path ?? "", forKey: LocalPreferences.goPathKey) }

    var dataFolder: URL? {
        defaults.string(forKey: LocalPreferences.dataFolderKey).map { URL(fileURLWithPath: ($0 as NSString).expandingTildeInPath, isDirectory: true) }
    }

    /// Code files left in conflict when the app quit, as "<library>/<lesson>/<file>": the
    /// lesson folder on this Mac holds the other editor's version, so this list belongs to the
    /// Mac too.
    func unsettledFiles(in ref: LessonRef) -> Set<String> {
        let prefix = ref.library + "/" + ref.lesson + "/"
        return Set((defaults.stringArray(forKey: LocalPreferences.unsettledKey) ?? []).compactMap { entry in
            entry.hasPrefix(prefix) ? String(entry.dropFirst(prefix.count)) : nil
        })
    }

    func setUnsettled(_ file: String, in ref: LessonRef, _ unsettled: Bool) {
        let entry = ref.library + "/" + ref.lesson + "/" + file
        var list = (defaults.stringArray(forKey: LocalPreferences.unsettledKey) ?? []).filter { $0 != entry }
        if unsettled { list.append(entry) }
        if list.isEmpty { defaults.removeObject(forKey: LocalPreferences.unsettledKey) } else { defaults.set(list, forKey: LocalPreferences.unsettledKey) }
    }
}

/// A place Back and Forward return to: a lesson, or a place in the Go documentation.
enum Place: Hashable {
    case dashboard
    case lesson(LessonRef)
    case doc(DocLink)
}

/// Places visited, for the toolbar's back and forward buttons.
struct VisitHistory {
    private(set) var back: [Place] = []
    private(set) var forward: [Place] = []

    mutating func visit(from current: Place?) {
        if let current, back.last != current { back.append(current) }
        if back.count > 200 { back.removeFirst() }
        forward.removeAll()
    }

    mutating func goBack(from current: Place?) -> Place? {
        guard let target = back.popLast() else { return nil }
        if let current { forward.append(current) }
        return target
    }

    mutating func goForward(from current: Place?) -> Place? {
        guard let target = forward.popLast() else { return nil }
        if let current { back.append(current) }
        return target
    }
}

/// A library in the sidebar: built into the app, or one you opened.
struct ShelfLibrary: Identifiable {
    /// The library's uuid; the file's location while it can't be read.
    let id: String
    let title: String
    let url: URL
    let isBuiltIn: Bool
    var database: AppDatabase?
    var catalog: Catalog = .empty
    var error: AppDatabaseError?
}

/// What the sidebar can select: the review, the Go documentation, a lesson, or a place in the
/// documentation that search found.
enum SidebarItem: Hashable {
    case dashboard
    case review
    case documentation
    case lesson(LessonRef)
    case doc(DocLink)
}

/// The app's state: the learner database, the libraries in the sidebar with your progress,
/// the lesson on screen, search, review, and the services lessons use. Views read it from the
/// environment.
@Observable @MainActor
final class AppModel {
    enum Phase: Equatable {
        case starting
        /// Your study data can't be opened, so nothing can be kept: the window says why.
        case failed(AppError)
        case ready
    }

    private(set) var phase: Phase = .starting
    private(set) var learner: LearnerDatabase?
    /// The sidebar's libraries: the ones built into the app first, then those you opened.
    private(set) var shelf: [ShelfLibrary] = []
    private(set) var toolchain: Toolchain?
    /// Finished lessons, by library.
    private(set) var done: [String: Set<String>] = [:]
    private(set) var lesson: LessonModel?
    private(set) var resumeLesson: LessonRef?
    private(set) var isLoadingLesson = false
    var selection: SidebarItem?
    /// Open chapters and guided projects, as "<library>/<chapter or course>".
    var expanded: Set<String> = []
    var searchText = "" { didSet { searchChanged() } }
    private(set) var searchResults: [SearchHit] = []
    /// Excerpts for the results, by lesson.
    private(set) var searchSnippets: [LessonRef: [Snippet.Part]] = [:]
    private(set) var searching = false
    var banner: AppError?
    /// A library you chose that can't be added: the alert names the problems.
    var libraryError: AppError?
    private(set) var isAddingLibrary = false
    /// Quiz questions and lessons due for review now.
    private(set) var dueCount = 0
    /// The review, while the sidebar's Review is selected.
    private(set) var review: ReviewSession?
    /// The Go documentation built into the app; nil until it has opened, or when it can't.
    private(set) var docs: DocsModel?
    /// The standard library's names, for links in readmes and Quick Help.
    private(set) var docNames = DocNames.empty
    /// What search found in the documentation, below the lessons.
    private(set) var docSearchResults: [DocSearchHit] = []
    /// The documentation is the place on screen, or about to be, for the history.
    @ObservationIgnored private var docsActive = false
    @ObservationIgnored private var dashboardActive = false
    /// The lesson that is the place on screen, or about to be while it loads, for the history;
    /// nil for the review and a redo.
    @ObservationIgnored private var lessonTarget: LessonRef?

    /// The review is on screen: selected, with no lesson still leaving.
    var showsReview: Bool { selection == .review && lesson == nil }
    var showsDashboard: Bool { selection == .dashboard && lesson == nil }

    var learningProgress: LearningProgress {
        LearningProgress(libraries: shelf.map { .init(id: $0.id, catalog: $0.catalog) }, done: done)
    }

    /// Resume an unfinished lesson, or pick the first unfinished mission across libraries.
    var recommendedLesson: LessonRef? {
        if let resumeLesson, summary(resumeLesson) != nil, !isDone(resumeLesson) { return resumeLesson }
        for library in shelf {
            if let lesson = library.catalog.courses.flatMap(\.lessons).first(where: { !isDone(LessonRef(library: library.id, lesson: $0.id)) }) {
                return LessonRef(library: library.id, lesson: lesson.id)
            }
        }
        return nil
    }

    /// The documentation is on screen: selected, with no lesson still leaving.
    var showsDocumentation: Bool {
        switch selection {
        case .documentation?, .doc?: lesson == nil && docs != nil
        default: false
        }
    }

    /// Shows an error from work that has no place of its own on screen, such as a save. Work
    /// that was cancelled (another lesson or search replaced it) ends silently.
    func report(_ error: any Error) {
        guard !Reporting.isSilent(error) else { return }
        banner = (error as? AppError) ?? (error as? AppDatabaseError).map(AppError.database) ?? .lesson(error.localizedDescription)
    }

    /// The same, for a problem with a lesson's files.
    func reportFiles(_ error: any Error) {
        guard !Reporting.isSilent(error) else { return }
        banner = .files(error.localizedDescription)
    }

    /// Whether the inspector is open, remembered on this Mac.
    var showInspector: Bool {
        get { access(keyPath: \.showInspector); return preferences.defaults.bool(forKey: "ShowInspector") }
        set { withMutation(keyPath: \.showInspector) { preferences.defaults.set(newValue, forKey: "ShowInspector") } }
    }
    private(set) var history = VisitHistory()
    var editorFontSize: Double = 13
    var showLineNumbers = true
    /// Set while the app saves before it quits.
    private(set) var isQuitting = false
    @ObservationIgnored private var started = false
    var confirmReset = false
    /// Bumped to open the Settings window (views own the openSettings action).
    private(set) var settingsRequests = 0

    let preferences: LocalPreferences
    let workspace: Workspace
    private let readmes = ReadmeCache()
    /// Readmes as plain text, for search excerpts; made ahead of time for every library.
    private let readmeTexts = ReadmeTexts()
    private var observeTask: Task<Void, Never>?
    private var searchTask: Task<Void, Never>?
    private var loadTask: Task<Void, Never>?
    /// Finds the Go to run code with; only the latest check counts (Settings › Go).
    private let toolchainWork = LatestWork<Toolchain>()
    private let signposter = OSSignposter(subsystem: Log.subsystem, category: .pointsOfInterest)
    private var launchInterval: OSSignpostIntervalState?

    init(defaults: UserDefaults = .standard) {
        preferences = LocalPreferences(defaults: defaults)
        workspace = Workspace(root: preferences.dataFolder?.appendingPathComponent("work", isDirectory: true) ?? Workspace.defaultRoot)
    }

    // MARK: Launch

    /// Opens the learner database, the libraries built into the app and the ones you opened,
    /// and the lesson you had open. The app delegate calls it once the app has launched,
    /// before or without a visible window; only the first call does anything.
    func start() async {
        guard !started else { return }
        started = true
        launchInterval = signposter.beginInterval("Launch")
        ObsoletePreferences.remove(from: preferences.defaults)
        openDocumentationDatabase()
        // What runs left behind after a crash; not the lesson folders.
        let workspace = workspace
        Task.detached(priority: .background) { workspace.removeStaleRuns() }
        let goOverride = preferences.goPath, goCaches = goCaches
        toolchainWork.start(priority: .utility, { await Toolchain.discover(goOverride: goOverride, caches: goCaches) }) { [weak self] in
            self?.toolchain = $0
        }

        let opened = ContinuousClock.now
        let url = preferences.dataFolder?.appendingPathComponent("learner.sqlite") ?? LearnerDatabase.defaultURL
        let schema = LearnerDatabase.bundledSchema
        let opening = await Task.detached(priority: .userInitiated) { () -> Result<LearnerDatabase, LearnerDatabaseError> in
            Result { () throws(LearnerDatabaseError) in try LearnerDatabase(url: url, schema: schema) }
        }.value
        let learner: LearnerDatabase
        switch opening {
        case .success(let opened): learner = opened
        case .failure(let e):
            phase = .failed(.learner(e))
            finishLaunch()
            return
        }
        self.learner = learner
        await loadShelf(learner)
        AppModel.logPhase("open (learner database and libraries)", since: opened)
        do {
            done = try await learner.done()
            editorFontSize = Double(try await learner.setting("editor_font_size") ?? "") ?? 13
            showLineNumbers = (try await learner.setting("editor_line_numbers")) != "0"
        } catch {
            report(error)
        }
        observeDone(learner)
        phase = .ready
        let texts = readmeTexts, databases = shelf.compactMap(\.database)
        Task.detached(priority: .utility) {
            for db in databases {
                do { try await texts.warm(db) } catch { Log.app.notice("Search excerpts: \(error.localizedDescription, privacy: .public)") }
            }
        }
        await refreshDueCount()
        resumeLesson = try? await learner.lastOpened()
        if preferences.defaults.string(forKey: LaunchArguments.openLesson) == nil,
           preferences.defaults.string(forKey: LaunchArguments.action) == nil {
            selection = .dashboard
            dashboardActive = true
            finishLaunch()
            return
        }
        if let ref = await launchLesson(learner) {
            // The "redo" launch actions open the lesson as its redo for review.
            if preferences.defaults.string(forKey: LaunchArguments.action)?.hasPrefix("redo") == true {
                redo(ref)
            } else {
                open(ref, remember: false)
            }
        } else {
            finishLaunch()
        }
    }

    /// The libraries built into the app, then the ones you opened, each opened read-only and
    /// checked off the main actor. One that can't be read stays in the sidebar with why.
    private func loadShelf(_ learner: LearnerDatabase) async {
        var libraries: [ShelfLibrary] = []
        for file in BuiltInLibraries.files(in: BuiltInLibraries.folder) {
            var entry = await AppModel.openLibrary(file, builtIn: true)
            if let info = entry.database?.info {
                do {
                    try await learner.registerLibrary(info, location: BuiltInLibraries.location(of: file))
                } catch {
                    report(error)
                }
                if libraries.contains(where: { $0.id == info.uuid }) { entry = ShelfLibrary(id: file.path, title: file.lastPathComponent, url: file, isBuiltIn: true, error: .sqlite("another built-in library has the same uuid")) }
            }
            libraries.append(entry)
        }
        do {
            for row in BuiltInLibraries.opened(from: try await learner.libraries()) {
                let url = URL(fileURLWithPath: row.location)
                var entry = await AppModel.openLibrary(url, builtIn: false)
                if let info = entry.database?.info, info.uuid != row.uuid {
                    // Another library now sits at that path.
                    entry = ShelfLibrary(id: row.uuid, title: row.title, url: url, isBuiltIn: false,
                                         error: .sqlite("\(url.path) now holds another library, \(info.title)."))
                } else if entry.database == nil {
                    entry = ShelfLibrary(id: row.uuid, title: row.title, url: url, isBuiltIn: false, error: entry.error)
                }
                libraries.append(entry)
            }
        } catch {
            report(error)
        }
        shelf = libraries
    }

    /// Opens and checks a library off the main actor, with its catalog.
    private static func openLibrary(_ url: URL, builtIn: Bool) async -> ShelfLibrary {
        await Task.detached(priority: .userInitiated) { () -> ShelfLibrary in
            do {
                let db = try AppDatabase(url: url, immutable: builtIn)
                let catalog = try await db.catalog()
                return ShelfLibrary(id: db.info.uuid, title: db.info.title, url: url, isBuiltIn: builtIn, database: db, catalog: catalog)
            } catch let e as AppDatabaseError {
                return ShelfLibrary(id: url.path, title: url.deletingPathExtension().lastPathComponent, url: url, isBuiltIn: builtIn, error: e)
            } catch {
                return ShelfLibrary(id: url.path, title: url.deletingPathExtension().lastPathComponent, url: url, isBuiltIn: builtIn,
                                    error: .sqlite(error.localizedDescription))
            }
        }.value
    }

    /// The lesson to open at launch: the one a scripted run asks for, else the one you opened
    /// last in any library, else the first one.
    private func launchLesson(_ learner: LearnerDatabase) async -> LessonRef? {
        if let requested = preferences.defaults.string(forKey: LaunchArguments.openLesson) {
            let parts = requested.split(separator: "/").map(String.init)
            let ref = parts.count == 2 ? LessonRef(library: parts[0], lesson: parts[1])
                : shelf.lazy.compactMap { lib in lib.catalog.summary(requested).map { _ in LessonRef(library: lib.id, lesson: requested) } }.first
            if let ref, summary(ref) != nil { return ref }
        }
        if let last = try? await learner.lastOpened(), summary(last) != nil { return last }
        for library in shelf {
            if let first = library.catalog.firstLesson { return LessonRef(library: library.id, lesson: first.id) }
        }
        return nil
    }

    private func finishLaunch() {
        guard let interval = launchInterval else { return }
        signposter.endInterval("Launch", interval)
        launchInterval = nil
        // From the process's start, so the dock click to a usable window is what counts.
        let sinceStart = AppModel.processStart.map { Date().timeIntervalSince($0) * 1000 } ?? -1
        Log.app.notice("Timing launch \(Int(sinceStart), privacy: .public) ms from process start")
    }

    static func logPhase(_ name: String, since start: ContinuousClock.Instant) {
        Log.app.notice("Timing \(name, privacy: .public) \(Int((ContinuousClock.now - start).milliseconds), privacy: .public) ms")
    }

    /// When this process started, from the kernel.
    static let processStart: Date? = {
        var info = kinfo_proc()
        var size = MemoryLayout<kinfo_proc>.stride
        var mib: [Int32] = [CTL_KERN, KERN_PROC, KERN_PROC_PID, getpid()]
        guard sysctl(&mib, 4, &info, &size, nil, 0) == 0 else { return nil }
        let t = info.kp_proc.p_un.__p_starttime
        return Date(timeIntervalSince1970: Double(t.tv_sec) + Double(t.tv_usec) / 1_000_000)
    }()

    /// Follows finished lessons as the app records them.
    private func observeDone(_ learner: LearnerDatabase) {
        observeTask?.cancel()
        observeTask = Task { [weak self] in
            do {
                for try await done in learner.doneUpdates() {
                    self?.done = done
                }
            } catch {
                self?.report(error)
            }
        }
    }

    // MARK: Libraries

    /// File › Open Library…: adds a library file to the sidebar, after its checks.
    func chooseLibrary() {
        let panel = NSOpenPanel()
        panel.canChooseDirectories = false
        panel.canChooseFiles = true
        panel.prompt = "Add"
        panel.message = "Choose a Gosmos library: an SQLite file of courses. Your progress in it stays on this Mac."
        guard panel.runModal() == .OK, let url = panel.url else { return }
        Task { await addLibrary(url) }
    }

    /// Adds a library you chose: it is checked first, recorded in the learner database with
    /// where it is, and shown after the others. One that is already there opens instead; one
    /// there that can't be read (its file moved, say) is this file from now on.
    func addLibrary(_ url: URL) async {
        guard let learner, !isAddingLibrary else { return }
        isAddingLibrary = true
        defer { isAddingLibrary = false }
        let entry = await AppModel.openLibrary(url, builtIn: false)
        guard let db = entry.database else {
            libraryError = .database(entry.error ?? .missing(path: url.path))
            return
        }
        if let existing = shelf.first(where: { $0.id == db.info.uuid }), existing.database != nil {
            if let first = existing.catalog.firstLesson { open(LessonRef(library: existing.id, lesson: first.id)) }
            return
        }
        do {
            let info = db.info, location = url.path
            try await learner.pendingWrites.start { try await learner.registerLibrary(info, location: location) }.value
        } catch {
            report(error)
            return
        }
        if let unreadable = shelf.firstIndex(where: { $0.id == db.info.uuid }) { shelf[unreadable] = entry } else { shelf.append(entry) }
        let texts = readmeTexts
        Task.detached(priority: .utility) { try? await texts.warm(db) }
        if let first = entry.catalog.firstLesson { open(LessonRef(library: entry.id, lesson: first.id)) }
    }

    func library(_ id: String) -> ShelfLibrary? { shelf.first { $0.id == id } }

    // MARK: Lessons

    func summary(_ ref: LessonRef) -> LessonSummary? { library(ref.library)?.catalog.summary(ref.lesson) }

    /// Shows a lesson; `remember` adds the one on screen to the back history.
    func open(_ ref: LessonRef, remember: Bool = true) {
        guard let learner, let db = library(ref.library)?.database, let catalog = library(ref.library)?.catalog,
              catalog.summary(ref.lesson) != nil else { return }
        review = nil
        // The lesson on screen, or the one still loading to take its place; none while the
        // documentation is the place, even with a lesson still leaving for it.
        let current = docsActive ? nil : lessonTarget
        if remember && ref != current { history.visit(from: visited) }
        resumeLesson = ref
        dashboardActive = false
        docsActive = false
        lessonTarget = ref
        if selection != .lesson(ref) { selection = .lesson(ref) }
        if let loc = catalog.location(of: ref.lesson) {
            let course = catalog.courses[loc.course]
            expanded.insert(ref.library + "/" + course.chapters[loc.chapter].id)
            expanded.insert(ref.library + "/" + course.id)
        }
        // Not again while it is on screen or loading; a load that failed may be tried again.
        let shown = lesson.map { $0.detail.ref == ref && !$0.isReview } ?? false
        guard current != ref || !(shown || isLoadingLesson) else { return }
        startLoading(ref, db: db, learner: learner, review: false)
    }

    /// Redoes a lesson due for review: its starter code in a folder of its own; your saved
    /// code and its done mark stay as they are. The redo is off the sidebar's selection, so
    /// Review brings the review back when it is done.
    func redo(_ ref: LessonRef) {
        guard let learner, let db = library(ref.library)?.database, summary(ref) != nil else { return }
        selection = nil
        dashboardActive = false
        review = nil
        docsActive = false
        lessonTarget = nil
        startLoading(ref, db: db, learner: learner, review: true)
    }

    /// The place on screen as the back and forward history knows it: the documentation's page,
    /// or the lesson, also while it still loads, so Back or Forward meanwhile keeps it; a redo is
    /// not in it.
    private var visited: Place? {
        if dashboardActive { return .dashboard }
        if docsActive { return docs?.location.map(Place.doc) }
        return lessonTarget.map(Place.lesson)
    }

    private func startLoading(_ ref: LessonRef, db: AppDatabase, learner: LearnerDatabase, review: Bool) {
        let interval = signposter.beginInterval("Lesson switch")
        let started = ContinuousClock.now
        isLoadingLesson = true
        loadTask?.cancel()
        let previous = lesson
        loadTask = Task { [weak self] in
            await previous?.leave()
            guard let self else { return }
            do {
                guard let detail = try await db.lesson(ref.lesson) else { throw AppError.lesson("This lesson is not in its library.") }
                try Task.checkCancellation()
                let parsed = await self.readmes.parsed(for: detail)
                var saved: [String: CodeFileRecord] = [:]
                if detail.code != nil && !review { saved = try await learner.codeFiles(for: ref) }
                try Task.checkCancellation()
                let model = LessonModel(detail: detail, parsed: parsed, saved: saved, learner: learner, review: review, app: self)
                self.lesson = model
                self.isLoadingLesson = false
                self.signposter.endInterval("Lesson switch", interval)
                let ms = (ContinuousClock.now - started).milliseconds
                Log.app.notice("Timing lesson switch \(Int(ms), privacy: .public) ms (\(detail.type.rawValue, privacy: .public))")
                self.finishLaunch()
                await model.load()
                // Not once a newer lesson replaced this one: that one is the last opened.
                if !review && !Task.isCancelled {
                    do { try await learner.pendingWrites.start { try await learner.markOpened(ref) }.value } catch { self.report(error) }
                }
                self.performLaunchAction(on: model)
            } catch let error where Reporting.isSilent(error) {
                self.signposter.endInterval("Lesson switch", interval)
            } catch let e as AppError {
                self.isLoadingLesson = false
                self.banner = e
            } catch {
                self.isLoadingLesson = false
                self.banner = .lesson(error.localizedDescription)
            }
        }
    }

    /// Launch arguments for scripted runs (screenshots, debugging): `-GosmosOpenLesson
    /// [<library>/]<lesson>` opens a lesson, and `-GosmosLaunchAction <action>` acts on it once,
    /// after it loads. They live in the argument domain only, never stored.
    enum LaunchArguments {
        static let openLesson = "GosmosOpenLesson"
        static let action = "GosmosLaunchAction"
    }

    @ObservationIgnored private var launchActionDone = false

    private func performLaunchAction(on model: LessonModel) {
        guard !launchActionDone, let action = preferences.defaults.string(forKey: LaunchArguments.action) else { return }
        launchActionDone = true
        if action.hasPrefix("search:") {
            // Types the query a letter at a time, as someone would.
            let query = String(action.dropFirst(7))
            Task { [weak self] in
                for end in query.indices {
                    try? await Task.sleep(for: .milliseconds(120))
                    self?.searchText = String(query[...end])
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

    var canGoBack: Bool { !history.back.isEmpty }
    var canGoForward: Bool { !history.forward.isEmpty }

    func goBack() {
        if let place = history.goBack(from: visited) { go(to: place) }
    }

    func goForward() {
        if let place = history.goForward(from: visited) { go(to: place) }
    }

    private func go(to place: Place) {
        switch place {
        case .dashboard: showDashboard(remember: false)
        case .lesson(let ref): open(ref, remember: false)
        case .doc(let link): openDocumentation(link, remember: false)
        }
    }

    var nextLesson: LessonSummary? {
        guard let ref = lesson?.detail.ref, lesson?.isReview == false else { return nil }
        return library(ref.library)?.catalog.next(after: ref.lesson)
    }
    var previousLesson: LessonSummary? {
        guard let ref = lesson?.detail.ref, lesson?.isReview == false else { return nil }
        return library(ref.library)?.catalog.previous(before: ref.lesson)
    }

    func openNext() { if let n = nextLesson, let library = lesson?.detail.library { open(LessonRef(library: library, lesson: n.id)) } }
    func openPrevious() { if let p = previousLesson, let library = lesson?.detail.library { open(LessonRef(library: library, lesson: p.id)) } }

    func isDone(_ ref: LessonRef) -> Bool { done[ref.library]?.contains(ref.lesson) ?? false }

    func setDone(_ ref: LessonRef, _ value: Bool) {
        guard let learner else { return }
        if value { done[ref.library, default: []].insert(ref.lesson) } else { done[ref.library]?.remove(ref.lesson) }
        learner.pendingWrites.start { [weak self] in
            do { try await learner.setDone(ref, value) } catch { self?.report(error) }
        }
    }

    /// Finished lessons out of required ones, for a course's ring.
    func progress(of course: Course, in library: String) -> (done: Int, total: Int) {
        let required = course.requiredLessons
        let finished = done[library] ?? []
        return (required.filter { finished.contains($0.id) }.count, required.count)
    }

    func doneCount(_ chapter: Chapter, in library: String) -> Int {
        let finished = done[library] ?? []
        return chapter.lessons.filter { finished.contains($0.id) }.count
    }

    // MARK: Review

    func showDashboard(remember: Bool = true) {
        guard !dashboardActive else { return }
        if remember { history.visit(from: visited) }
        dashboardActive = true
        docsActive = false
        lessonTarget = nil
        selection = .dashboard
        isLoadingLesson = false
        loadTask?.cancel()
        let previous = lesson
        loadTask = Task { [weak self] in
            await previous?.leave()
            guard let self, !Task.isCancelled else { return }
            self.lesson = nil
            self.review = nil
        }
    }

    /// Shows the review: due quiz questions, interleaved across lessons, and lessons to redo.
    /// The lesson on screen stays until it is saved, as when another lesson opens, so a code
    /// conflict can still be asked about; Back returns to it.
    func showReview() {
        guard let learner else { return }
        if let ref = visited { history.visit(from: ref) }
        dashboardActive = false
        docsActive = false
        lessonTarget = nil
        selection = .review
        isLoadingLesson = false
        loadTask?.cancel()
        let previous = lesson
        let session = ReviewSession(learner: learner, app: self)
        loadTask = Task { [weak self] in
            await previous?.leave()
            guard let self, !Task.isCancelled else { return }
            self.lesson = nil
            self.review = session
            await session.load()
        }
    }

    /// Counts what is due now, for the sidebar's Review: what the review would show.
    func refreshDueCount() async {
        guard let learner else { return }
        do {
            let due = try await ReviewSession.due(learner, in: self)
            dueCount = due.questions.count + due.lessons.count
        } catch {
            report(error)
        }
    }

    // MARK: Documentation

    /// Opens the Go documentation built into the app, off the main actor, and reads the
    /// standard library's names for links. Without it the app works as before, with no
    /// documentation.
    private func openDocumentationDatabase() {
        guard let url = ReferenceDatabase.bundledURL else { return }
        Task { [weak self] in
            let opened = await Task.detached(priority: .userInitiated) { () -> Result<(ReferenceDatabase, DocNames), any Error> in
                Result {
                    let reference = try ReferenceDatabase(url: url)
                    return (reference, try reference.names())
                }
            }.value
            guard let self else { return }
            switch opened {
            case .success(let (reference, names)):
                self.docs = DocsModel(reference: reference, app: self)
                self.docNames = names
            case .failure(let error):
                Log.app.error("Go documentation: \(error.localizedDescription, privacy: .public)")
                self.report(error)
            }
        }
    }

    /// Shows a place in the documentation, or the page shown last; `remember` adds the place on
    /// screen to the back history. A lesson on screen is saved and left, as when another opens.
    func openDocumentation(_ link: DocLink? = nil, remember: Bool = true, selecting item: SidebarItem = .documentation) {
        guard let docs else { return }
        let target = link ?? docs.location
        if remember, let place = visited, place != target.map(Place.doc) { history.visit(from: place) }
        dashboardActive = false
        docsActive = true
        lessonTarget = nil
        if selection != item { selection = item }
        if let target { docs.show(target) }
        review = nil
        isLoadingLesson = false
        loadTask?.cancel()
        guard let previous = lesson else { return }
        loadTask = Task { [weak self] in
            await previous.leave()
            guard let self, !Task.isCancelled else { return }
            self.lesson = nil
        }
    }

    /// Follows a link clicked in the app: into the documentation, or to a lesson of the
    /// library on screen. Any other address does nothing; Gosmos opens nothing outside itself.
    func follow(_ url: URL) -> Bool {
        let library = lesson?.detail.library
        let catalog = library.flatMap { self.library($0)?.catalog }
        switch LinkTarget(destination: url.absoluteString, isLesson: { catalog?.summary($0) != nil }) {
        case .doc(let link)?:
            openDocumentation(link)
            return true
        case .lesson(let id)?:
            if let library { open(LessonRef(library: library, lesson: id)) }
            return true
        case nil:
            return false
        }
    }

    // MARK: Search

    private func searchChanged() {
        searchTask?.cancel()
        let text = searchText
        let databases = shelf.compactMap(\.database)
        let reference = docs?.reference
        guard !databases.isEmpty || reference != nil, !text.trimmingCharacters(in: .whitespaces).isEmpty else {
            searchResults = []
            docSearchResults = []
            searching = false
            // The place search found stays on screen, as the documentation.
            if case .doc? = selection { selection = .documentation }
            return
        }
        searching = true
        let texts = readmeTexts
        searchTask = Task { [weak self] in
            let interval = self?.signposter.beginInterval("Search")
            let started = ContinuousClock.now
            // The index's hits that show a match, with their excerpts, from the readmes as
            // the app renders them, through every library in the sidebar's order.
            let found: (hits: [SearchHit], snippets: [LessonRef: [Snippet.Part]])
            do {
                found = try await texts.search(text, in: databases)
            } catch {
                found = ([], [:])
                self?.report(error)
            }
            // The documentation's best matches, below the lessons.
            var docHits: [DocSearchHit] = []
            if let reference {
                do { docHits = try await reference.search(text, limit: 8) } catch { self?.report(error) }
            }
            guard let self, !Task.isCancelled else { return }
            self.searchResults = found.hits
            self.docSearchResults = docHits
            self.searchSnippets = found.snippets
            self.searching = false
            if let interval { self.signposter.endInterval("Search", interval) }
            let ms = (ContinuousClock.now - started).milliseconds
            Log.app.notice("Timing search \(String(format: "%.1f", ms), privacy: .public) ms for \(text.count, privacy: .public) characters")
        }
    }

    // MARK: Settings

    func setEditorFontSize(_ size: Double) {
        editorFontSize = min(max(size, 9), 28)
        guard let learner else { return }
        let value = String(Int(editorFontSize))
        learner.pendingWrites.start { [weak self] in
            do { try await learner.setSetting("editor_font_size", value) } catch { self?.report(error) }
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
        guard let learner else { return }
        learner.pendingWrites.start { [weak self] in
            do { try await learner.setSetting("editor_line_numbers", on ? "1" : "0") } catch { self?.report(error) }
        }
    }

    // MARK: Lesson folders

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

    /// Saves what is pending before the app quits. Lesson programs stop at once; the app
    /// delegate bounds the whole wait.
    func prepareToQuit() async {
        isQuitting = true
        // A code conflict asks first; without a choice in 15 seconds both versions are kept.
        await lesson?.leave(within: .seconds(15))
        ProcessRunner.terminateAll(.lesson)
        // The writes nothing else waits for (a review answer, a hint, a done mark) land before
        // the app goes.
        await learner?.pendingWrites.finish()
    }
}
