import GosmosKit
import Observation

/// Where you are in the story (docs/storytelling.md): the opening on screen, the name on your
/// badge, and the openings you have been through. The shelf holds the courses' stories; the app
/// opens the lessons they lead to.
@Observable @MainActor
final class StoryModel {
    /// The opening on screen.
    struct Opening: Equatable {
        let library: String
        let course: String
        var page: Int
        /// Continue Learning returns here after an unseen opening, rather than restarting.
        var destination: LessonRef? = nil
    }

    /// The learner database's settings for the story: the name on your badge, and the openings
    /// you have been through, one "<library>/<course>" a line.
    static let badgeNameKey = "badge_name"
    static let openingsSeenKey = "openings_seen"

    /// The opening on screen, over the whole window: the course it opens, and its page.
    private(set) var opening: Opening?
    /// The name on your badge, the one the team calls you; empty until you give one.
    private(set) var badgeName = ""
    /// The openings you have been through, as "<library>/<course>".
    private(set) var seenOpenings: Set<String> = []
    /// The app whose shelf holds the stories, and which shows what they lead to; it sets itself
    /// here when it is made.
    @ObservationIgnored weak var app: AppModel?

    /// Reads the name on your badge and the openings you have been through.
    func load(from learner: LearnerDatabase) async throws {
        badgeName = try await learner.setting(StoryModel.badgeNameKey) ?? ""
        seenOpenings = Set((try await learner.setting(StoryModel.openingsSeenKey) ?? "").split(separator: "\n").map(String.init))
    }

    /// The opening a new learner sees at launch, one who has seen no opening and opened or
    /// finished no lesson: that of the course Mission Control offers. Every other course's
    /// opening waits until you begin that course.
    var openingAtLaunch: Opening? {
        guard let app, seenOpenings.isEmpty, app.resumeLesson == nil, app.done.values.allSatisfy(\.isEmpty),
              let ref = app.recommendedLesson, let course = app.library(ref.library)?.catalog.course(of: ref.lesson),
              app.story(of: course.id, in: ref.library)?.hasOpening == true else { return nil }
        return Opening(library: ref.library, course: course.id, page: 0)
    }

    /// At launch, a new learner meets the story before anything else.
    func showOpeningAtLaunch() {
        opening = openingAtLaunch
    }

    /// Begins or continues a course: its opening the first time, else its next lesson.
    func begin(_ course: Course, in library: String) {
        guard let app else { return }
        if app.story(of: course.id, in: library)?.hasOpening == true, !seenOpenings.contains(library + "/" + course.id) {
            showOpening(course.id, in: library)
        } else if let next = app.firstUnfinished(course.lessons, in: library) ?? course.lessons.first {
            app.open(LessonRef(library: library, lesson: next.id))
        }
    }

    /// Open the assignment Mission Control actually offers, even before the first pass.
    func continueLearning(_ ref: LessonRef) {
        guard let app, let course = app.library(ref.library)?.catalog.course(of: ref.lesson) else { return }
        if app.story(of: course.id, in: ref.library)?.hasOpening == true,
           !seenOpenings.contains(ref.library + "/" + course.id) {
            showOpening(course.id, in: ref.library)
            opening?.destination = ref
        } else {
            app.open(ref)
        }
    }

    /// The opening View › Replay the Opening shows: that of the course you are in, the lesson's on
    /// screen, done or not; with no lesson on screen, that of the course Mission Control offers;
    /// else the first on the shelf.
    var openingToReplay: (library: String, course: String)? {
        guard let app else { return nil }
        if let ref = app.lesson?.detail.ref ?? app.recommendedLesson ?? app.resumeLesson, let course = app.library(ref.library)?.catalog.course(of: ref.lesson),
           app.story(of: course.id, in: ref.library)?.hasOpening == true {
            return (ref.library, course.id)
        }
        for library in app.shelf.libraries {
            for course in library.catalog.courses where library.stories[course.id]?.hasOpening == true { return (library.id, course.id) }
        }
        return nil
    }

    /// Shows the opening again; a lesson on screen is saved and left first, as when Mission
    /// Control opens.
    func replayOpening() {
        guard let o = openingToReplay else { return }
        app?.showDashboard()
        showOpening(o.course, in: o.library)
    }

    /// Shows a course's opening, from its first page or `page`.
    func showOpening(_ course: String, in library: String, page: Int = 0) {
        guard let app, let pages = app.story(of: course, in: library)?.pages, !pages.isEmpty else { return }
        app.shownTerm = nil
        opening = Opening(library: library, course: course, page: min(max(0, page), pages.count - 1))
    }

    /// The opening's next page; after its last one, the course.
    func advanceOpening() {
        guard let o = opening, let pages = app?.story(of: o.course, in: o.library)?.pages else { return }
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
        guard let app, let course = app.library(o.library)?.catalog.courses.first(where: { $0.id == o.course }) else { return }
        if let destination = o.destination, app.summary(destination) != nil, !app.isDone(destination) {
            app.open(destination)
        } else if let next = app.firstUnfinished(course.lessons, in: o.library) {
            app.open(LessonRef(library: o.library, lesson: next.id))
        } else {
            app.showDashboard()
        }
    }

    private func markSeen(_ o: Opening) {
        guard seenOpenings.insert(o.library + "/" + o.course).inserted, let learner = app?.learner else { return }
        let value = seenOpenings.sorted().joined(separator: "\n")
        learner.pendingWrites.start { [weak app] in
            do { try await learner.setSetting(StoryModel.openingsSeenKey, value) } catch { app?.report(error) }
        }
    }

    /// Puts a name on your badge; an empty one takes it off.
    func setBadgeName(_ name: String) {
        let name = String(name.trimmingCharacters(in: .whitespacesAndNewlines).prefix(40))
        guard name != badgeName else { return }
        badgeName = name
        guard let learner = app?.learner else { return }
        learner.pendingWrites.start { [weak app] in
            do { try await learner.setSetting(StoryModel.badgeNameKey, name.isEmpty ? nil : name) } catch { app?.report(error) }
        }
    }
}

// MARK: - The story

extension AppModel {
    typealias Opening = StoryModel.Opening

    var opening: Opening? { story.opening }
    var badgeName: String { story.badgeName }
    var openingToReplay: (library: String, course: String)? { story.openingToReplay }
    func begin(_ course: Course, in library: String) { story.begin(course, in: library) }
    func continueLearning(_ ref: LessonRef) { story.continueLearning(ref) }
    func replayOpening() { story.replayOpening() }
    func showOpening(_ course: String, in library: String, page: Int = 0) { story.showOpening(course, in: library, page: page) }
    func advanceOpening() { story.advanceOpening() }
    func openingPageBack() { story.openingPageBack() }
    func finishOpening() { story.finishOpening() }
    func setBadgeName(_ name: String) { story.setBadgeName(name) }
}
