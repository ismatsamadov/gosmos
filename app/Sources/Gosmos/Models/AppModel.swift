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

/// A place Back and Forward return to: Mission Control, the field guide, a lesson, or a place
/// in the Go documentation.
enum Place: Hashable {
    case dashboard
    case fieldGuide
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

    /// An entry for the place on screen is no step: it is left when you go to a place the
    /// history doesn't keep (the review, a redo) and come straight back.
    mutating func goBack(from current: Place?) -> Place? {
        while let last = back.last, last == current { back.removeLast() }
        guard let target = back.popLast() else { return nil }
        if let current { forward.append(current) }
        return target
    }

    mutating func goForward(from current: Place?) -> Place? {
        while let next = forward.last, next == current { forward.removeLast() }
        guard let target = forward.popLast() else { return nil }
        if let current { back.append(current) }
        return target
    }

    func canGoBack(from current: Place?) -> Bool { back.contains { $0 != current } }
    func canGoForward(from current: Place?) -> Bool { forward.contains { $0 != current } }
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
    /// The stories of its courses, by course uuid.
    var stories: [String: CourseStory] = [:]
    var error: AppDatabaseError?
}

/// What the sidebar can select: Mission Control, the review, the field guide, the Go
/// documentation, a lesson, or a place in the documentation that search found.
enum SidebarItem: Hashable {
    case dashboard
    case review
    case fieldGuide
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
    private(set) var shelf: [ShelfLibrary] = [] { didSet { updateLearningProgress() } }
    private(set) var toolchain: Toolchain?
    /// Finished lessons, by library.
    private(set) var done: [String: Set<String>] = [:] { didSet { updateLearningProgress() } }
    /// XP, rank and badges, from `done` over the shelf's catalogs: made again when either
    /// changes, so the dashboard, the sidebar, the inspector and Settings all read one value.
    private(set) var learningProgress = LearningProgress(libraries: [], done: [:])
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
    /// Mission Control is the place on screen, or about to be, for the history.
    @ObservationIgnored private var dashboardActive = false
    /// The field guide is the place on screen, for the history.
    @ObservationIgnored private var fieldGuideActive = false
    /// The opening on screen, over the whole window: the course it opens, and its page.
    private(set) var opening: Opening?
    /// The name on your badge, the one the team calls you; empty until you give one.
    private(set) var badgeName = ""
    /// The openings you have been through, as "<library>/<course>".
    private(set) var seenOpenings: Set<String> = []
    /// The field guide's term shown over the lesson on screen, after a click on its link.
    var shownTerm: CourseStory.Term?
    /// The term the field guide brings into view when it opens.
    private(set) var fieldGuideFocus: String?
    /// Lessons' story beats, the first paragraph of each readme, for the teaser after a pass.
    @ObservationIgnored private var beats: [LessonRef: [MDInline]] = [:]
    /// The lesson that is the place on screen, or about to be while it loads, for the history;
    /// nil for Mission Control, the documentation, the review and a redo.
    @ObservationIgnored private var lessonTarget: LessonRef?

    /// The review is on screen: selected, with no lesson still leaving.
    var showsReview: Bool { selection == .review && lesson == nil }

    /// The field guide is on screen: selected, with no lesson still leaving.
    var showsFieldGuide: Bool { selection == .fieldGuide && lesson == nil }

    private func updateLearningProgress() {
        learningProgress = LearningProgress(libraries: shelf.map { .init(id: $0.id, catalog: $0.catalog) }, done: done)
    }

    /// The lesson Mission Control offers: the one you opened last while it is unfinished, else
    /// the next unfinished one in its course, so finishing a lesson moves you on from it; else
    /// the first unfinished one on the shelf.
    var recommendedLesson: LessonRef? {
        if let ref = resumeLesson, let catalog = library(ref.library)?.catalog, catalog.summary(ref.lesson) != nil {
            if !isDone(ref) { return ref }
            let finished = done[ref.library] ?? []
            if let next = catalog.nextUnfinished(after: ref.lesson, isDone: { finished.contains($0) }) {
                return LessonRef(library: ref.library, lesson: next.id)
            }
        }
        for library in shelf {
            for course in library.catalog.courses {
                if let lesson = firstUnfinished(course.lessons, in: library.id) {
                    return LessonRef(library: library.id, lesson: lesson.id)
                }
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
    /// then shows Mission Control, or the lesson a scripted run asks for. The app delegate calls
    /// it once the app has launched, before or without a visible window; only the first call
    /// does anything.
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
            badgeName = try await learner.setting(AppModel.badgeNameKey) ?? ""
            seenOpenings = Set((try await learner.setting(AppModel.openingsSeenKey) ?? "").split(separator: "\n").map(String.init))
            resumeLesson = try await learner.lastOpened()
        } catch {
            report(error)
        }
        observeDone(learner)
        // Mission Control is the place on screen before the window can take a click, so nothing
        // chosen while the launch finishes is undone; a scripted run opens its lesson below.
        let action = preferences.defaults.string(forKey: LaunchArguments.action)
        // The story's launch actions show a screen of their own, with Mission Control behind it.
        let storyAction = action.flatMap { $0 == "field-guide" || $0 == "opening" || $0.hasPrefix("opening:") ? $0 : nil }
        let scripted = storyAction == nil && (preferences.defaults.string(forKey: LaunchArguments.openLesson) != nil || action != nil)
        if !scripted {
            selection = .dashboard
            dashboardActive = true
            // A new learner meets the story before anything else.
            if storyAction == nil { self.opening = openingAtLaunch }
        }
        phase = .ready
        let texts = readmeTexts, databases = shelf.compactMap(\.database)
        Task.detached(priority: .utility) {
            for db in databases {
                do { try await texts.warm(db) } catch { Log.app.notice("Search excerpts: \(error.localizedDescription, privacy: .public)") }
            }
        }
        await refreshDueCount()
        if let storyAction { performStoryAction(storyAction) }
        guard scripted else {
            finishLaunch()
            return
        }
        if let ref = launchLesson() {
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
                let stories = try await db.stories()
                return ShelfLibrary(id: db.info.uuid, title: db.info.title, url: url, isBuiltIn: builtIn, database: db, catalog: catalog, stories: stories)
            } catch let e as AppDatabaseError {
                return ShelfLibrary(id: url.path, title: url.deletingPathExtension().lastPathComponent, url: url, isBuiltIn: builtIn, error: e)
            } catch {
                return ShelfLibrary(id: url.path, title: url.deletingPathExtension().lastPathComponent, url: url, isBuiltIn: builtIn,
                                    error: .sqlite(error.localizedDescription))
            }
        }.value
    }

    /// The lesson a scripted run opens at launch: the one it asks for, else the one you opened
    /// last in any library, else the first one. A launch without one shows Mission Control.
    private func launchLesson() -> LessonRef? {
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
        dashboardActive = false
        fieldGuideActive = false
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
        fieldGuideActive = false
        review = nil
        docsActive = false
        lessonTarget = nil
        startLoading(ref, db: db, learner: learner, review: true)
    }

    /// The place on screen as the back and forward history knows it: Mission Control, the
    /// documentation's page, or the lesson, also while it still loads, so Back or Forward
    /// meanwhile keeps it; the review and a redo are not in it.
    private var visited: Place? {
        if dashboardActive { return .dashboard }
        if fieldGuideActive { return .fieldGuide }
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
                // Not once a newer lesson replaced this one: that one is the last opened. Mission
                // Control resumes the same lesson the next launch will, never one that failed to load.
                if !review && !Task.isCancelled {
                    self.resumeLesson = ref
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

    var canGoBack: Bool { history.canGoBack(from: visited) }
    var canGoForward: Bool { history.canGoForward(from: visited) }

    func goBack() {
        if let place = history.goBack(from: visited) { go(to: place) }
    }

    func goForward() {
        if let place = history.goForward(from: visited) { go(to: place) }
    }

    private func go(to place: Place) {
        switch place {
        case .dashboard: showDashboard(remember: false)
        case .fieldGuide: showFieldGuide(remember: false)
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

    func estimate(of course: Course, in library: String) -> StudyEstimate {
        StudyEstimate(lessons: course.lessons, completed: done[library] ?? [])
    }

    func doneCount(_ chapter: Chapter, in library: String) -> Int {
        let finished = done[library] ?? []
        return chapter.lessons.filter { finished.contains($0.id) }.count
    }

    /// The first of `lessons` not done yet, in their order.
    func firstUnfinished(_ lessons: [LessonSummary], in library: String) -> LessonSummary? {
        let finished = done[library] ?? []
        return lessons.first { !finished.contains($0.id) }
    }

    // MARK: Mission Control

    /// Shows Mission Control; `remember` adds the place on screen to the back history. A lesson
    /// on screen is saved and left, as when another opens.
    func showDashboard(remember: Bool = true) {
        guard !dashboardActive else {
            // Already the place: the list may have cleared its selection (a Command-click on
            // the selected row), so the sidebar shows it selected again.
            if selection != .dashboard { selection = .dashboard }
            return
        }
        if remember { history.visit(from: visited) }
        dashboardActive = true
        fieldGuideActive = false
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

    // MARK: Story

    /// The opening on screen.
    struct Opening: Equatable {
        let library: String
        let course: String
        var page: Int
    }

    /// The learner database's settings for the story: the name on your badge, and the openings
    /// you have been through, one "<library>/<course>" a line.
    static let badgeNameKey = "badge_name"
    static let openingsSeenKey = "openings_seen"

    func story(of course: String, in library: String) -> CourseStory? { self.library(library)?.stories[course] }

    /// The story of the course a lesson is in.
    func story(for ref: LessonRef) -> CourseStory? {
        library(ref.library)?.catalog.course(of: ref.lesson).flatMap { story(of: $0.id, in: ref.library) }
    }

    /// A chapter's part in its course's story.
    func chapterStory(_ chapter: Chapter, in library: String) -> CourseStory.ChapterStory? {
        story(of: chapter.courseID, in: library)?.chapters[chapter.id]
    }

    /// Whether the field guide has a word to show: a course on the shelf has terms.
    var hasFieldGuide: Bool { shelf.contains { $0.stories.values.contains { !$0.terms.isEmpty } } }

    /// A chapter whose required lessons are all done, as the chapter milestone counts it.
    func isComplete(_ chapter: Chapter, in library: String) -> Bool {
        let required = chapter.lessons.filter(\.required)
        return !required.isEmpty && required.allSatisfy { isDone(LessonRef(library: library, lesson: $0.id)) }
    }

    /// Whether `ref` is the one required lesson of its chapter not done yet, so passing it
    /// finishes the chapter.
    func isLastToDo(_ ref: LessonRef) -> Bool {
        guard let chapter = library(ref.library)?.catalog.chapter(of: ref.lesson) else { return false }
        let required = chapter.lessons.filter(\.required)
        return required.contains { $0.id == ref.lesson }
            && required.allSatisfy { $0.id == ref.lesson || isDone(LessonRef(library: ref.library, lesson: $0.id)) }
    }

    /// The opening a new learner sees at launch: the first course with an opening not seen yet,
    /// and no lesson of it done.
    private var openingAtLaunch: Opening? {
        for library in shelf {
            for course in library.catalog.courses where library.stories[course.id]?.hasOpening == true {
                guard !seenOpenings.contains(library.id + "/" + course.id),
                      !course.lessons.contains(where: { isDone(LessonRef(library: library.id, lesson: $0.id)) }) else { continue }
                return Opening(library: library.id, course: course.id, page: 0)
            }
        }
        return nil
    }

    /// Begins or continues a course: its opening the first time, else its next lesson.
    func begin(_ course: Course, in library: String) {
        if story(of: course.id, in: library)?.hasOpening == true, !seenOpenings.contains(library + "/" + course.id) {
            showOpening(course.id, in: library)
        } else if let next = firstUnfinished(course.lessons, in: library) ?? course.lessons.first {
            open(LessonRef(library: library, lesson: next.id))
        }
    }

    /// The opening View › Replay the Opening shows: that of the course you are in, else the first
    /// on the shelf.
    var openingToReplay: (library: String, course: String)? {
        if let ref = recommendedLesson ?? resumeLesson, let course = library(ref.library)?.catalog.course(of: ref.lesson),
           story(of: course.id, in: ref.library)?.hasOpening == true {
            return (ref.library, course.id)
        }
        for library in shelf {
            for course in library.catalog.courses where library.stories[course.id]?.hasOpening == true { return (library.id, course.id) }
        }
        return nil
    }

    /// Shows the opening again; a lesson on screen is saved and left first, as when Mission
    /// Control opens.
    func replayOpening() {
        guard let o = openingToReplay else { return }
        showDashboard()
        showOpening(o.course, in: o.library)
    }

    /// Shows a course's opening, from its first page or `page`.
    func showOpening(_ course: String, in library: String, page: Int = 0) {
        guard let pages = story(of: course, in: library)?.pages, !pages.isEmpty else { return }
        shownTerm = nil
        opening = Opening(library: library, course: course, page: min(max(0, page), pages.count - 1))
    }

    /// The opening's next page; after its last one, the course.
    func advanceOpening() {
        guard let o = opening, let pages = story(of: o.course, in: o.library)?.pages else { return }
        if o.page + 1 < pages.count { opening?.page = o.page + 1 } else { finishOpening() }
    }

    /// The opening's page before the one on screen.
    func openingPageBack() {
        if let o = opening, o.page > 0 { opening?.page = o.page - 1 }
    }

    /// Ends the opening, read or skipped, and opens the course where you are in it: its next
    /// unfinished lesson, or Mission Control once every lesson is done.
    func finishOpening() {
        guard let o = opening else { return }
        opening = nil
        markSeen(o)
        guard let course = library(o.library)?.catalog.courses.first(where: { $0.id == o.course }) else { return }
        if let next = firstUnfinished(course.lessons, in: o.library) {
            open(LessonRef(library: o.library, lesson: next.id))
        } else {
            showDashboard()
        }
    }

    private func markSeen(_ o: Opening) {
        guard seenOpenings.insert(o.library + "/" + o.course).inserted, let learner else { return }
        let value = seenOpenings.sorted().joined(separator: "\n")
        learner.pendingWrites.start { [weak self] in
            do { try await learner.setSetting(AppModel.openingsSeenKey, value) } catch { self?.report(error) }
        }
    }

    /// Puts a name on your badge; an empty one takes it off.
    func setBadgeName(_ name: String) {
        let name = String(name.trimmingCharacters(in: .whitespacesAndNewlines).prefix(40))
        guard name != badgeName else { return }
        badgeName = name
        guard let learner else { return }
        learner.pendingWrites.start { [weak self] in
            do { try await learner.setSetting(AppModel.badgeNameKey, name.isEmpty ? nil : name) } catch { self?.report(error) }
        }
    }

    /// Shows the field guide, with `term` in view; `remember` adds the place on screen to the
    /// back history. A lesson on screen is saved and left, as when another opens.
    func showFieldGuide(focus term: String? = nil, remember: Bool = true) {
        shownTerm = nil
        fieldGuideFocus = term
        guard !fieldGuideActive else {
            if selection != .fieldGuide { selection = .fieldGuide }
            return
        }
        if remember { history.visit(from: visited) }
        dashboardActive = false
        fieldGuideActive = true
        docsActive = false
        lessonTarget = nil
        selection = .fieldGuide
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

    /// A lesson's story beat, the first paragraph of its readme, for the teaser after a pass.
    func beat(of ref: LessonRef) async -> [MDInline]? {
        if let beat = beats[ref] { return beat }
        guard let db = library(ref.library)?.database,
              let row = try? await db.readmes([ref.lesson]).first else { return nil }
        guard case .paragraph(let inlines)? = Readme.parse(row.readme, droppingTitle: row.title).first else { return nil }
        beats[ref] = inlines
        return inlines
    }

    /// The story's launch actions: `opening` or `opening:<page>` shows the first opening on the
    /// shelf, and `field-guide` the field guide.
    private func performStoryAction(_ action: String) {
        if action == "field-guide" {
            showFieldGuide(remember: false)
            return
        }
        for library in shelf {
            for course in library.catalog.courses where library.stories[course.id]?.hasOpening == true {
                showOpening(course.id, in: library.id, page: Int(action.dropFirst("opening:".count)) ?? 0)
                return
            }
        }
    }

    // MARK: Review

    /// Shows the review: due quiz questions, interleaved across lessons, and lessons to redo.
    /// The lesson on screen stays until it is saved, as when another lesson opens, so a code
    /// conflict can still be asked about; Back returns to it.
    func showReview() {
        guard let learner else { return }
        if let ref = visited { history.visit(from: ref) }
        dashboardActive = false
        fieldGuideActive = false
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
        fieldGuideActive = false
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
        let story = lesson.flatMap { story(for: $0.detail.ref) }
        switch LinkTarget(destination: url.absoluteString, isLesson: { catalog?.summary($0) != nil }, isTerm: { story?.term($0) != nil }) {
        case .doc(let link)?:
            openDocumentation(link)
            return true
        case .term(let name)?:
            shownTerm = story?.term(name)
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
