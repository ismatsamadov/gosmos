import Foundation
import GosmosKit
import Observation

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

    /// Opens and checks a library off the main actor, with its catalog.
    static func open(_ url: URL, builtIn: Bool) async -> ShelfLibrary {
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
}

/// The libraries in the sidebar and your progress through them: the lessons you finished, the
/// one you opened last, the lesson Mission Control offers, and what the courses' stories and
/// chapters make of them.
@Observable @MainActor
final class ShelfModel {
    /// The sidebar's libraries: the ones built into the app first, then those you opened.
    private(set) var libraries: [ShelfLibrary] = [] { didSet { updateLearningProgress() } }
    /// Finished lessons, by library.
    private(set) var done: [String: Set<String>] = [:] { didSet { updateLearningProgress() } }
    /// XP, rank and badges, from `done` over the shelf's catalogs: made again when either
    /// changes, so the dashboard, the sidebar, the inspector and Settings all read one value.
    private(set) var learningProgress = LearningProgress(libraries: [], done: [:])
    /// The lesson you opened last; the app sets it as a lesson loads.
    var resumeLesson: LessonRef?
    /// The app whose learner database keeps your progress; it sets itself here when it is made.
    @ObservationIgnored weak var app: AppModel?
    /// Lessons' story beats, the first paragraph of each readme, for the teaser after a pass.
    @ObservationIgnored private var beats: [LessonRef: [MDInline]] = [:]
    private var observeTask: Task<Void, Never>?

    private func updateLearningProgress() {
        learningProgress = LearningProgress(libraries: libraries.map { .init(id: $0.id, catalog: $0.catalog) }, done: done)
    }

    /// The libraries built into the app, then the ones you opened, each opened read-only and
    /// checked off the main actor. One that can't be read stays in the sidebar with why.
    func load(builtIn folder: URL?, learner: LearnerDatabase) async {
        var libraries: [ShelfLibrary] = []
        for file in BuiltInLibraries.files(in: folder) {
            var entry = await ShelfLibrary.open(file, builtIn: true)
            if let info = entry.database?.info {
                do {
                    try await learner.registerLibrary(info, location: BuiltInLibraries.location(of: file))
                } catch {
                    app?.report(error)
                }
                if libraries.contains(where: { $0.id == info.uuid }) { entry = ShelfLibrary(id: file.path, title: file.lastPathComponent, url: file, isBuiltIn: true, error: .sqlite("another built-in library has the same uuid")) }
            }
            libraries.append(entry)
        }
        do {
            for row in BuiltInLibraries.opened(from: try await learner.libraries()) {
                let url = URL(fileURLWithPath: row.location)
                var entry = await ShelfLibrary.open(url, builtIn: false)
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
            app?.report(error)
        }
        self.libraries = libraries
    }

    /// Puts a library you added on the shelf: in place of the one with its uuid that can't be
    /// read, else after the others.
    func add(_ entry: ShelfLibrary) {
        if let unreadable = libraries.firstIndex(where: { $0.id == entry.id }) { libraries[unreadable] = entry } else { libraries.append(entry) }
    }

    /// Reads the lessons you finished.
    func readDone(from learner: LearnerDatabase) async throws {
        done = try await learner.done()
    }

    /// Follows finished lessons as the app records them.
    func observeDone(_ learner: LearnerDatabase) {
        observeTask?.cancel()
        observeTask = Task { [weak self] in
            do {
                for try await done in learner.doneUpdates() {
                    self?.done = done
                }
            } catch {
                self?.app?.report(error)
            }
        }
    }

    func library(_ id: String) -> ShelfLibrary? { libraries.first { $0.id == id } }

    func summary(_ ref: LessonRef) -> LessonSummary? { library(ref.library)?.catalog.summary(ref.lesson) }

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
        for library in libraries {
            for course in library.catalog.courses {
                if let lesson = firstUnfinished(course.lessons, in: library.id) {
                    return LessonRef(library: library.id, lesson: lesson.id)
                }
            }
        }
        return nil
    }

    func isDone(_ ref: LessonRef) -> Bool { done[ref.library]?.contains(ref.lesson) ?? false }

    func setDone(_ ref: LessonRef, _ value: Bool) {
        guard let learner = app?.learner else { return }
        if value { done[ref.library, default: []].insert(ref.lesson) } else { done[ref.library]?.remove(ref.lesson) }
        learner.pendingWrites.start { [weak app] in
            do { try await learner.setDone(ref, value) } catch { app?.report(error) }
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
    var hasFieldGuide: Bool { libraries.contains { $0.stories.values.contains { !$0.terms.isEmpty } } }

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

    /// A lesson's story beat, the first paragraph of its readme, for the teaser after a pass.
    func beat(of ref: LessonRef) async -> [MDInline]? {
        if let beat = beats[ref] { return beat }
        guard let db = library(ref.library)?.database,
              let row = try? await db.readmes([ref.lesson]).first else { return nil }
        guard case .paragraph(let inlines)? = Readme.parse(row.readme, droppingTitle: row.title).first else { return nil }
        beats[ref] = inlines
        return inlines
    }
}

// MARK: - The shelf

extension AppModel {
    var done: [String: Set<String>] { shelf.done }
    var learningProgress: LearningProgress { shelf.learningProgress }
    var resumeLesson: LessonRef? { shelf.resumeLesson }
    var recommendedLesson: LessonRef? { shelf.recommendedLesson }
    func library(_ id: String) -> ShelfLibrary? { shelf.library(id) }
    func summary(_ ref: LessonRef) -> LessonSummary? { shelf.summary(ref) }
    func isDone(_ ref: LessonRef) -> Bool { shelf.isDone(ref) }
    func setDone(_ ref: LessonRef, _ value: Bool) { shelf.setDone(ref, value) }
    func progress(of course: Course, in library: String) -> (done: Int, total: Int) { shelf.progress(of: course, in: library) }
    func estimate(of course: Course, in library: String) -> StudyEstimate { shelf.estimate(of: course, in: library) }
    func doneCount(_ chapter: Chapter, in library: String) -> Int { shelf.doneCount(chapter, in: library) }
    func firstUnfinished(_ lessons: [LessonSummary], in library: String) -> LessonSummary? { shelf.firstUnfinished(lessons, in: library) }
    func story(of course: String, in library: String) -> CourseStory? { shelf.story(of: course, in: library) }
    func story(for ref: LessonRef) -> CourseStory? { shelf.story(for: ref) }
    func chapterStory(_ chapter: Chapter, in library: String) -> CourseStory.ChapterStory? { shelf.chapterStory(chapter, in: library) }
    var hasFieldGuide: Bool { shelf.hasFieldGuide }
    func isComplete(_ chapter: Chapter, in library: String) -> Bool { shelf.isComplete(chapter, in: library) }
    func isLastToDo(_ ref: LessonRef) -> Bool { shelf.isLastToDo(ref) }
    func beat(of ref: LessonRef) async -> [MDInline]? { await shelf.beat(of: ref) }
}
