import GosmosKit
import SwiftUI

struct SidebarView: View {
    @Environment(AppModel.self) private var app
    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    var body: some View {
        Group {
            if !app.searchText.trimmingCharacters(in: .whitespaces).isEmpty {
                SearchResultsList()
            } else if app.shelf.isEmpty {
                ContentUnavailableView {
                    Label("No Courses Yet", systemImage: "books.vertical")
                } description: {
                    Text("Add a library: an SQLite file of courses. Your progress in it stays on this Mac.")
                        .foregroundStyle(Theme.secondaryText)
                } actions: {
                    Button("Open Library…") { app.chooseLibrary() }
                        .disabled(app.isAddingLibrary)
                }
            } else {
                ScrollViewReader { proxy in
                    List(selection: app.sidebarSelection) {
                        ReviewRow()
                            .tag(SidebarItem.review)
                        ForEach(app.shelf) { library in
                            LibrarySection(library: library)
                        }
                    }
                    .listStyle(.sidebar)
                    .task(id: app.lesson?.detail.ref) {
                        // Bring the lesson on screen into view once its chapter has opened.
                        guard let ref = app.lesson?.detail.ref else { return }
                        try? await Task.sleep(for: .milliseconds(120))
                        withAnimation(reduceMotion ? nil : .snappy) { proxy.scrollTo(SidebarItem.lesson(ref), anchor: .center) }
                    }
                }
            }
        }
    }
}

extension AppModel {
    /// The lists' selection. A row chosen in a list is only recorded here; it opens on the next
    /// turn of the main actor, after the list has finished its selection callback. The model's
    /// own selections (Back, Next Lesson, launch) come back unchanged and open nothing, so
    /// nothing opens twice.
    var sidebarSelection: Binding<SidebarItem?> {
        Binding(
            get: { self.selection },
            set: { item in
                guard item != self.selection else { return }
                self.selection = item
                guard let item else { return }
                Task { @MainActor in
                    guard self.selection == item else { return }
                    switch item {
                    case .review: self.showReview()
                    case .lesson(let ref): self.open(ref)
                    }
                }
            })
    }
}

/// Review: what is due now, with its count.
struct ReviewRow: View {
    @Environment(AppModel.self) private var app

    var body: some View {
        HStack(spacing: 7) {
            Image(systemName: "arrow.trianglehead.2.clockwise.rotate.90")
                .symbolRenderingMode(.hierarchical)
                .foregroundStyle(.secondary)
                .frame(width: 18)
            Text("Review")
            Spacer(minLength: 4)
            ReviewBadge(count: app.dueCount)
        }
        .frame(minHeight: 24)
        .accessibilityElement(children: .combine)
        .accessibilityLabel(app.dueCount == 0 ? "Review, nothing due" : "Review, \(app.dueCount) due")
    }
}

/// One library: its courses, then its guided projects. A library that can't be read says why.
struct LibrarySection: View {
    let library: ShelfLibrary

    var body: some View {
        Section {
            if let error = library.error {
                Label(error.title, systemImage: "exclamationmark.triangle")
                    .foregroundStyle(Theme.secondaryText)
                    .help([error.problems.joined(separator: "\n"), error.recoverySuggestion, library.url.path].compactMap { $0 }.joined(separator: "\n\n"))
                    .accessibilityLabel("\(library.title) can't be opened: \(error.title)")
            } else if library.catalog.courses.isEmpty {
                Text("No courses yet")
                    .foregroundStyle(Theme.secondaryText)
            } else {
                ForEach(library.catalog.courses) { course in
                    CourseGroup(course: course, library: library.id)
                }
            }
        } header: {
            HStack(spacing: 6) {
                Text(library.title)
                    .lineLimit(1)
                if !library.isBuiltIn {
                    Image(systemName: "doc")
                        .imageScale(.small)
                        .foregroundStyle(.tertiary)
                        .help(library.url.path)
                }
            }
            .accessibilityElement(children: .combine)
            .accessibilityLabel(library.isBuiltIn ? library.title : "\(library.title), opened from \(library.url.path)")
        }
    }
}

/// A course or guided project, with its progress ring; it opens to its chapters.
struct CourseGroup: View {
    let course: Course
    let library: String
    @Environment(AppModel.self) private var app

    var body: some View {
        let key = library + "/" + course.id
        let expanded = Binding(
            get: { app.expanded.contains(key) },
            set: { open in
                if open { app.expanded.insert(key) } else {
                    app.expanded.remove(key)
                    for ch in course.chapters { app.expanded.remove(library + "/" + ch.id) }
                }
            })
        DisclosureGroup(isExpanded: expanded) {
            ForEach(course.chapters) { chapter in
                ChapterGroup(chapter: chapter, library: library)
            }
        } label: {
            let progress = app.progress(of: course, in: library)
            HStack(spacing: 8) {
                if course.isGuidedProject {
                    Image(systemName: "hammer")
                        .symbolRenderingMode(.hierarchical)
                        .foregroundStyle(.secondary)
                }
                Text(course.title).lineLimit(1)
                Spacer(minLength: 6)
                ProgressRing(done: progress.done, total: progress.total)
            }
            .accessibilityElement(children: .combine)
            .accessibilityLabel("\(course.title)\(course.isGuidedProject ? ", guided project" : ""), \(progress.done) of \(progress.total) required lessons done")
        }
    }
}

/// Finished lessons out of required ones.
struct ProgressRing: View {
    let done: Int
    let total: Int
    var size: CGFloat = 14
    @Environment(\.accessibilityReduceMotion) private var reduceMotion

    var body: some View {
        let fraction = total == 0 ? 0 : Double(done) / Double(total)
        ZStack {
            Circle().stroke(.quaternary, lineWidth: 2.5)
            Circle()
                .trim(from: 0, to: fraction)
                .stroke(.tint, style: StrokeStyle(lineWidth: 2.5, lineCap: .round))
                .rotationEffect(.degrees(-90))
        }
        .frame(width: size, height: size)
        .animation(reduceMotion ? nil : .smooth, value: fraction)
        .help("\(done) of \(total) required lessons done")
    }
}

struct ChapterGroup: View {
    let chapter: Chapter
    let library: String
    @Environment(AppModel.self) private var app

    var body: some View {
        let key = library + "/" + chapter.id
        let expanded = Binding(
            get: { app.expanded.contains(key) },
            set: { open in
                if open { app.expanded.insert(key) } else { app.expanded.remove(key) }
            })
        DisclosureGroup(isExpanded: expanded) {
            ForEach(chapter.lessons) { lesson in
                let ref = LessonRef(library: library, lesson: lesson.id)
                LessonRow(lesson: lesson, ref: ref)
                    .tag(SidebarItem.lesson(ref))
                    .id(SidebarItem.lesson(ref))
            }
        } label: {
            let done = app.doneCount(chapter, in: library)
            HStack(spacing: 6) {
                Text(chapter.title)
                    .lineLimit(1)
                Spacer(minLength: 6)
                Text("\(done)/\(chapter.lessons.count)")
                    .font(.caption.monospacedDigit())
                    .foregroundStyle(done == chapter.lessons.count && done > 0 ? AnyShapeStyle(.tint) : AnyShapeStyle(Theme.secondaryText))
            }
            .accessibilityElement(children: .combine)
            .accessibilityLabel("\(chapter.title), \(done) of \(chapter.lessons.count) done")
        }
    }
}

struct LessonRow: View {
    let lesson: LessonSummary
    let ref: LessonRef
    @Environment(AppModel.self) private var app
    @Environment(\.backgroundProminence) private var prominence

    var body: some View {
        let done = app.isDone(ref)
        HStack(spacing: 7) {
            Image(systemName: lesson.type.symbol)
                .symbolRenderingMode(.hierarchical)
                .foregroundStyle(.secondary)
                .frame(width: 18)
            Text(lesson.title)
                .lineLimit(1)
                // Selected, the row is white on the accent, so an optional title takes the list's own secondary.
                .foregroundStyle(lesson.required ? AnyShapeStyle(.primary) : prominence == .increased ? AnyShapeStyle(.secondary) : AnyShapeStyle(Theme.secondaryText))
            Spacer(minLength: 4)
            // Always there, hidden until done, so the bounce plays when the lesson is finished.
            Image(systemName: "checkmark.circle.fill")
                .foregroundStyle(.tint)
                .opacity(done ? 1 : 0)
                .symbolEffect(.bounce, value: done)
        }
        .frame(minHeight: 24)
        .accessibilityElement(children: .combine)
        .accessibilityLabel("\(lesson.title), \(lesson.type.label)\(done ? ", done" : "")\(lesson.required ? "" : ", optional")")
    }
}

struct SearchResultsList: View {
    @Environment(AppModel.self) private var app

    var body: some View {
        if app.searchResults.isEmpty && !app.searching {
            ContentUnavailableView.search(text: app.searchText)
        } else {
            List(selection: app.sidebarSelection) {
                Section("\(app.searchResults.count) Lesson\(app.searchResults.count == 1 ? "" : "s")") {
                    ForEach(app.searchResults) { hit in
                        SearchHitRow(hit: hit)
                            .tag(SidebarItem.lesson(hit.ref))
                    }
                }
            }
            .listStyle(.sidebar)
            // A new list for each set of results. Updated in place, the table measures rows of
            // differing heights while it keeps its top row still, resizes in the middle of that
            // and re-enters its own row-height cache, which AppKit logs as a reentrant
            // operation in its delegate and will turn into an assertion.
            .id(app.searchResults.map(\.ref))
        }
    }
}

struct SearchHitRow: View {
    let hit: SearchHit
    @Environment(AppModel.self) private var app
    @Environment(\.backgroundProminence) private var prominence

    var body: some View {
        VStack(alignment: .leading, spacing: 3) {
            HStack(spacing: 6) {
                Image(systemName: hit.type.symbol)
                    .symbolRenderingMode(.hierarchical)
                    .foregroundStyle(.secondary)
                Text(title)
                    .fontWeight(.medium)
                    .lineLimit(1)
                Spacer(minLength: 4)
                if app.isDone(hit.ref) {
                    Image(systemName: "checkmark.circle.fill").foregroundStyle(.tint)
                }
            }
            Text("\(hit.courseTitle) › \(hit.chapterTitle)")
                .font(.caption)
                .foregroundStyle(secondary)
                .lineLimit(1)
            // A lesson found by its title alone may have no readme text to quote.
            if !(app.searchSnippets[hit.ref] ?? []).isEmpty {
                Text(snippet)
                    .font(.caption)
                    .foregroundStyle(secondary)
                    .lineLimit(2)
            }
        }
        .padding(.vertical, 3)
        .accessibilityElement(children: .combine)
    }

    /// A focused list draws its selected row in white on the accent, where colors of our own
    /// would sink in.
    private var emphasized: Bool { prominence == .increased }
    private var secondary: AnyShapeStyle { emphasized ? AnyShapeStyle(.secondary) : AnyShapeStyle(Theme.secondaryText) }
    /// On the selected row, gray or accent, the matched words keep their weight but drop the
    /// accent, which reads poorly on either fill.
    private var plainMarks: Bool { emphasized || app.selection == .lesson(hit.ref) }

    /// The title with the words the query matches marked by the excerpt's rule: words that
    /// start with a word of the query.
    private var title: AttributedString {
        // Enough words that the excerpt is the whole title.
        let parts = Snippet.make(hit.title, query: app.searchText, words: 4 * (hit.title.count + 1))
        guard parts.contains(where: \.matched) else { return AttributedString(hit.title) }
        var out = AttributedString()
        for part in parts {
            var a = AttributedString(part.text)
            if part.matched {
                a.inlinePresentationIntent = .stronglyEmphasized
                if !plainMarks { a.foregroundColor = Theme.searchMark }
            }
            out += a
        }
        return out
    }

    private var snippet: AttributedString {
        var out = AttributedString()
        for part in app.searchSnippets[hit.ref] ?? [] {
            var a = AttributedString(part.text)
            if part.matched {
                a.font = .caption.weight(.semibold)
                if !plainMarks { a.foregroundColor = Theme.searchMark }
            }
            out += a
        }
        return out
    }
}
