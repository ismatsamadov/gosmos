import AppKit
import CodeEditor
import GosmosKit
import Observation
import os
import SwiftUI

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
    private(set) var lesson: LessonModel?
    private(set) var isLoadingLesson = false
    var selection: SidebarItem?
    /// Open chapters and guided projects, as "<library>/<chapter or course>".
    var expanded: Set<String> = []
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
    /// The field guide's term shown over the lesson on screen, after a click on its link.
    var shownTerm: CourseStory.Term?
    /// The term the field guide brings into view when it opens.
    private(set) var fieldGuideFocus: String?

    /// The review is on screen: selected, with no lesson still leaving.
    var showsReview: Bool { selection == .review && lesson == nil }

    /// The field guide is on screen: selected, with no lesson still leaving.
    var showsFieldGuide: Bool { selection == .fieldGuide && lesson == nil }

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
    /// The sidebar's libraries and your progress through them.
    let shelf = ShelfModel()
    /// The place on screen and the places visited, for Back and Forward.
    let navigation = NavigationModel()
    /// The sidebar's search, through the shelf and the documentation.
    let search = SearchModel()
    /// The opening on screen, the name on your badge and the openings seen.
    let story = StoryModel()
    /// Settings › Editor and Settings › Go.
    let settings: SettingsModel
    /// Set while the app saves before it quits.
    private(set) var isQuitting = false
    @ObservationIgnored private var started = false
    var confirmReset = false
    /// Bumped to open the Settings window (views own the openSettings action).
    var settingsRequests = 0
    /// Whether a scripted run's launch action has run; it runs once.
    @ObservationIgnored var launchActionDone = false

    let preferences: LocalPreferences
    let workspace: Workspace
    private let resources: AppResources
    private let readmes = ReadmeCache()
    private var loadTask: Task<Void, Never>?
    private let signposter = OSSignposter(subsystem: Log.subsystem, category: .pointsOfInterest)
    private var launchInterval: OSSignpostIntervalState?

    init(defaults: UserDefaults = .standard, resources: AppResources = .bundled) {
        preferences = LocalPreferences(defaults: defaults)
        settings = SettingsModel(preferences: preferences)
        self.resources = resources
        workspace = Workspace(root: preferences.dataFolder?.appendingPathComponent("work", isDirectory: true) ?? Workspace.defaultRoot)
        navigation.docLocation = { [weak self] in self?.docs?.location }
        shelf.app = self
        search.app = self
        story.app = self
        settings.app = self
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
        settings.discoverToolchain()

        let opened = ContinuousClock.now
        let url = preferences.dataFolder?.appendingPathComponent("learner.sqlite") ?? LearnerDatabase.defaultURL
        let schema = resources.learnerSchema
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
        await shelf.load(builtIn: resources.libraries, learner: learner)
        AppModel.logPhase("open (learner database and libraries)", since: opened)
        do {
            try await shelf.readDone(from: learner)
            try await settings.load(from: learner)
            try await story.load(from: learner)
            shelf.resumeLesson = try await learner.lastOpened()
        } catch {
            report(error)
        }
        shelf.observeDone(learner)
        // Mission Control is the place on screen before the window can take a click, so nothing
        // chosen while the launch finishes is undone; a scripted run opens its lesson below.
        let action = preferences.defaults.string(forKey: LaunchArguments.action)
        // The story's launch actions show a screen of their own, with Mission Control behind it.
        let storyAction = action.flatMap { $0 == "field-guide" || $0 == "opening" || $0.hasPrefix("opening:") ? $0 : nil }
        let scripted = storyAction == nil && (preferences.defaults.string(forKey: LaunchArguments.openLesson) != nil || action != nil)
        if !scripted {
            selection = .dashboard
            navigation.screen = .dashboard
            // A new learner meets the story before anything else.
            if storyAction == nil { story.showOpeningAtLaunch() }
        }
        phase = .ready
        let texts = search.texts, databases = shelf.libraries.compactMap(\.database)
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
        let entry = await ShelfLibrary.open(url, builtIn: false)
        guard let db = entry.database else {
            libraryError = .database(entry.error ?? .missing(path: url.path))
            return
        }
        if let existing = shelf.library(db.info.uuid), existing.database != nil {
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
        shelf.add(entry)
        let texts = search.texts
        Task.detached(priority: .utility) { try? await texts.warm(db) }
        if let first = entry.catalog.firstLesson { open(LessonRef(library: entry.id, lesson: first.id)) }
    }

    // MARK: Lessons

    /// Shows a lesson; `remember` adds the one on screen to the back history.
    func open(_ ref: LessonRef, remember: Bool = true) {
        guard let learner, let db = library(ref.library)?.database, let catalog = library(ref.library)?.catalog,
              catalog.summary(ref.lesson) != nil else { return }
        review = nil
        // The lesson on screen, or the one still loading to take its place; none while the
        // documentation is the place, even with a lesson still leaving for it.
        let current = navigation.lessonTarget
        if remember && ref != current { navigation.visit() }
        navigation.screen = .lesson(ref)
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
        navigation.screen = nil
        review = nil
        startLoading(ref, db: db, learner: learner, review: true)
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
                // Control resumes the same lesson the next launch will, never one that failed to load,
                // its code included: a lesson folder from before the folders' records is measured
                // against this time when it is next prepared (Workspace.recordAppsCopies).
                if !review && !Task.isCancelled && model.code?.loadError == nil {
                    self.shelf.resumeLesson = ref
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

    /// A lesson on screen whose code failed to load, then loaded on a retry: it is the lesson you
    /// opened last from now on, as it would have been had it loaded at once.
    func codeLoaded(_ ref: LessonRef) {
        guard let learner, lesson?.ref == ref, lesson?.isReview == false else { return }
        shelf.resumeLesson = ref
        learner.pendingWrites.start { [weak self] in
            do { try await learner.markOpened(ref) } catch { self?.report(error) }
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

    /// A clear way on after a pass, including the course's last lesson and a review redo.
    func continueAfterLesson() {
        guard let lesson else { return }
        if lesson.isReview { showReview() }
        else if nextLesson != nil { openNext() }
        else { showDashboard() }
    }

    // MARK: Mission Control

    /// Shows Mission Control; `remember` adds the place on screen to the back history. A lesson
    /// on screen is saved and left, as when another opens.
    func showDashboard(remember: Bool = true) {
        guard navigation.screen != .dashboard else {
            // Already the place: the list may have cleared its selection (a Command-click on
            // the selected row), so the sidebar shows it selected again.
            if selection != .dashboard { selection = .dashboard }
            return
        }
        if remember { navigation.visit() }
        navigation.screen = .dashboard
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

    /// Shows the field guide, with `term` in view; `remember` adds the place on screen to the
    /// back history. A lesson on screen is saved and left, as when another opens.
    func showFieldGuide(focus term: String? = nil, remember: Bool = true) {
        shownTerm = nil
        fieldGuideFocus = term
        guard navigation.screen != .fieldGuide else {
            if selection != .fieldGuide { selection = .fieldGuide }
            return
        }
        if remember { navigation.visit() }
        navigation.screen = .fieldGuide
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

    // MARK: Review

    /// Shows the review: due quiz questions, interleaved across lessons, and lessons to redo.
    /// The lesson on screen stays until it is saved, as when another lesson opens, so a code
    /// conflict can still be asked about; Back returns to it.
    func showReview() {
        guard let learner else { return }
        if navigation.visited != nil { navigation.visit() }
        navigation.screen = nil
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
        guard let url = resources.reference else { return }
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
        if remember, let place = navigation.visited, place != target.map(Place.doc) { navigation.visit() }
        navigation.screen = .documentation
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
