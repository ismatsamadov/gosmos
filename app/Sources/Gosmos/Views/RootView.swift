import AppKit
import GosmosKit
import SwiftUI

struct RootView: View {
    @Environment(AppModel.self) private var app
    @Environment(\.openSettings) private var openSettings

    var body: some View {
        @Bindable var app = app
        content
            // Declared once, here, for the window's whole life: the items follow the app's
            // state, so changing lessons, lesson types or libraries never adds a second set.
            .toolbar { LessonToolbar() }
            // Links act in the app or not at all: Gosmos opens nothing outside itself.
            .environment(\.openURL, OpenURLAction { url in app.follow(url) ? .handled : .discarded })
            .onChange(of: app.settingsRequests) { openSettings() }
            .alert(app.libraryError?.title ?? "", isPresented: Binding(get: { app.libraryError != nil }, set: { if !$0 { app.libraryError = nil } }),
                   presenting: app.libraryError) { _ in
                Button("OK") { app.libraryError = nil }
            } message: { error in
                Text([error.errorDescription, error.recoverySuggestion].compactMap { $0 }.joined(separator: "\n\n"))
            }
            .overlay {
                if app.isQuitting {
                    SavingOverlay()
                }
            }
    }

    @ViewBuilder private var content: some View {
        switch app.phase {
        case .starting:
            VStack(spacing: 20) {
                OrbitEmblem()
                Text("A little Go. A long way.")
                    .font(.system(size: 25, weight: .bold, design: .rounded))
                ProgressView("Preparing your learning space…")
                    .controlSize(.small)
                    .foregroundStyle(Theme.secondaryText)
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .background(Theme.canvas)
        case .failed(let error):
            UnavailableView(error: error)
        case .ready:
            MainView()
        }
    }
}

/// While the app saves before it quits.
struct SavingOverlay: View {
    var body: some View {
        ZStack {
            Color.black.opacity(0.12).ignoresSafeArea()
            HStack(spacing: 10) {
                ProgressView().controlSize(.small)
                Text("Saving…")
                    .font(.system(size: 14, weight: .medium))
            }
            .padding(.horizontal, 18)
            .padding(.vertical, 12)
            .glassEffect(.regular, in: Capsule())
        }
        .transition(.opacity)
        .accessibilityElement(children: .combine)
        .accessibilityLabel("Saving your study data before Gosmos quits")
    }
}

/// When your study data can't be opened: why, and what to do. Nothing else opens without it,
/// since nothing of what you do could be kept.
struct UnavailableView: View {
    let error: AppError

    var body: some View {
        ContentUnavailableView {
            Label(error.title, systemImage: "externaldrive.badge.exclamationmark")
        } description: {
            VStack(spacing: 8) {
                if let text = error.errorDescription { Text(text) }
                if let suggestion = error.recoverySuggestion { Text(suggestion) }
            }
            .foregroundStyle(Theme.secondaryText)
            .textSelection(.enabled)
        } actions: {
            if case .learner(let e) = error, let path = e.path {
                Button("Show in Finder") { NSWorkspace.shared.activateFileViewerSelecting([URL(fileURLWithPath: path)]) }
                    .secondaryGlassButton()
            }
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
}

struct MainView: View {
    @Environment(AppModel.self) private var app
    @ViewState private var columns: NavigationSplitViewVisibility = .all
    @ViewState private var windowWidth: CGFloat = 0
    @ViewState private var sidebarWidth: CGFloat = 0
    @FocusState private var searchFocused: Bool

    var body: some View {
        @Bindable var app = app
        NavigationSplitView(columnVisibility: $columns) {
            SidebarView()
                .onGeometryChange(for: CGFloat.self) { $0.size.width } action: { sidebarWidth = $0 }
                // On the column's outermost view, where the split view reads it. The split view
                // keeps the sidebar and the inspector at their ideal widths, so those ideals plus
                // the lesson's 600 fit the 1100 pt minimum window: 260 + 600 + 240. 260 leaves a
                // lesson title about 20 characters a line; longer ones take a second line.
                .navigationSplitViewColumnWidth(min: 220, ideal: 260, max: 400)
        } detail: {
            DetailView(minWidth: lessonMinimum)
        }
        .searchable(text: $app.searchText, placement: .toolbar, prompt: "Search")
        .searchFocused($searchFocused)
        // Edit › Search (⇧⌘F): the documentation's own field while it is on screen, else the
        // toolbar's.
        .focusedSceneValue(\.searchFocus, SearchFocusAction { [app] in
            if app.showsDocumentation {
                NotificationCenter.default.post(name: DocsSearchField.focusRequest, object: nil)
            } else {
                searchFocused = true
            }
        })
        .onGeometryChange(for: CGFloat.self) { $0.size.width } action: { windowWidth = $0 }
    }

    /// 600, or less when a sidebar dragged wide and the inspector leave less room in the
    /// window: the split view keeps the sidebar's width, so a fixed minimum would push the
    /// inspector out of the window.
    private var lessonMinimum: CGFloat {
        guard windowWidth > 0 else { return 600 }
        let sidebar = columns == .detailOnly ? 0 : sidebarWidth
        let room = windowWidth - sidebar - (app.showInspector ? 240 : 0) - 2
        return min(600, max(440, room.rounded(.down)))
    }
}

/// What Edit › Search does in the focused window: puts the insertion point in its search field.
struct SearchFocusAction {
    let focus: @MainActor () -> Void
}

private struct SearchFocusKey: FocusedValueKey {
    typealias Value = SearchFocusAction
}

extension FocusedValues {
    var searchFocus: SearchFocusAction? {
        get { self[SearchFocusKey.self] }
        set { self[SearchFocusKey.self] = newValue }
    }
}

struct DetailView: View {
    var minWidth: CGFloat = 600
    @Environment(AppModel.self) private var app

    var body: some View {
        @Bindable var app = app
        Group {
            if app.showsDashboard {
                DashboardView()
            } else if app.showsDocumentation, let docs = app.docs {
                DocsBrowser(docs: docs)
            } else if app.showsReview {
                if let review = app.review {
                    ReviewScreen(session: review)
                } else {
                    ProgressView().appearsAfterDelay().frame(maxWidth: .infinity, maxHeight: .infinity)
                }
            } else if let lesson = app.lesson {
                LessonView(lesson: lesson)
                    .id(lesson.detail.ref)
            } else if app.isLoadingLesson {
                ProgressView().appearsAfterDelay().frame(maxWidth: .infinity, maxHeight: .infinity)
            } else {
                DashboardView()
            }
        }
        // Sidebar 260 + lesson 600 + inspector 240 fit the 1100 pt minimum window.
        .frame(minWidth: minWidth)
        .background(Theme.canvas)
        .safeAreaInset(edge: .top, spacing: 0) { ErrorBanner() }
        .navigationTitle(app.showsDashboard ? "Mission Control" : app.showsDocumentation ? app.docs?.title ?? "Go Documentation" : app.showsReview ? "Review" : app.lesson?.detail.title ?? "Gosmos")
        // With no page open the title already says Go Documentation, so the subtitle doesn't repeat it.
        .navigationSubtitle(app.showsDashboard ? "Your Go learning journey" : app.showsDocumentation ? (app.docs?.page == nil ? "" : "Go Documentation") : app.showsReview ? "" : app.lesson.map { ($0.isReview ? "Redo for review · " : "") + "\($0.detail.courseTitle) › \($0.detail.chapterTitle)" } ?? "")
        .inspector(isPresented: $app.showInspector) {
            InspectorView()
                .inspectorColumnWidth(min: 240, ideal: 240, max: 380)
        }
        .confirmationDialog("Reset to the starter code?", isPresented: $app.confirmReset) {
            Button("Reset to Starter Code", role: .destructive) { app.lesson?.code?.resetToStarter() }
        } message: {
            Text("Your code for this lesson is replaced by the lesson's starter code.")
        }
    }
}

struct LessonView: View {
    let lesson: LessonModel

    var body: some View {
        if let code = lesson.code {
            CodeLessonView(lesson: lesson, session: code)
        } else if let quiz = lesson.quiz {
            QuizLessonView(lesson: lesson, session: quiz)
        }
    }
}

struct ErrorBanner: View {
    @Environment(AppModel.self) private var app

    var body: some View {
        if let error = app.banner {
            HStack(spacing: 10) {
                Image(systemName: "exclamationmark.triangle.fill")
                    .symbolRenderingMode(.multicolor)
                VStack(alignment: .leading, spacing: 2) {
                    Text(error.errorDescription ?? "Something went wrong").fontWeight(.medium)
                    if let s = error.recoverySuggestion { Text(s).foregroundStyle(Theme.secondaryText) }
                }
                .font(.system(size: 12.5))
                Spacer()
                Button("Dismiss") { app.banner = nil }
                    .secondaryGlassButton()
                    .controlSize(.small)
            }
            .padding(.horizontal, 14)
            .padding(.vertical, 10)
            .glassEffect(.regular.tint(.orange.opacity(0.15)), in: RoundedRectangle(cornerRadius: 12, style: .continuous))
            .padding(.horizontal, 12)
            .padding(.top, 8)
            .transition(.move(edge: .top).combined(with: .opacity))
        }
    }
}

/// The window's toolbar items, all of them, as a function of the app's state: navigation,
/// the lesson's actions for its type and state, and the inspector toggle. Only the root view
/// declares it, so no lesson view adds items of its own.
struct LessonToolbar: ToolbarContent {
    @Environment(AppModel.self) private var app

    var body: some ToolbarContent {
        if app.phase == .ready {
            ToolbarItemGroup(placement: .navigation) {
                Button { app.goBack() } label: { Label("Back", systemImage: "chevron.backward") }
                    .disabled(!app.canGoBack)
                    .help("Back (⌥⌘←)")
                Button { app.goForward() } label: { Label("Forward", systemImage: "chevron.forward") }
                    .disabled(!app.canGoForward)
                    .help("Forward (⌥⌘→)")
            }
            ToolbarItemGroup(placement: .primaryAction) {
                if let lesson = app.lesson {
                    LessonActions(lesson: lesson)
                }
            }
            ToolbarItem(placement: .primaryAction) {
                @Bindable var app = app
                Button { app.showInspector.toggle() } label: { Label("Inspector", systemImage: "sidebar.trailing") }
                    .help("Show the note and lesson facts (⌥⌘I)")
            }
        }
    }
}

/// The open lesson's actions: Next Lesson after a pass, Run or Stop, and Submit, by its type
/// and state.
struct LessonActions: View {
    let lesson: LessonModel
    @Environment(AppModel.self) private var app

    var body: some View {
        if lesson.justPassed, let next = app.nextLesson {
            Button { app.openNext() } label: { Label("Next Lesson", systemImage: "arrow.right") }
                .labelStyle(.titleAndIcon)
                .help("Next lesson: \(next.title) (⌘])")
        }
        if lesson.canRun {
            Button {
                if lesson.isBusy { lesson.code?.stop() } else { lesson.run() }
            } label: {
                Label(lesson.isBusy ? "Stop" : "Run", systemImage: lesson.isBusy ? "stop.fill" : "play.fill")
                    .contentTransition(.symbolEffect(.replace))
            }
            .labelStyle(.titleAndIcon)
            .help(lesson.isBusy ? "Stop (⌘.)" : "Run (⌘R)")
        }
        if lesson.code != nil || lesson.quiz != nil {
            Button { lesson.submit() } label: { Label("Submit", systemImage: "paperplane.fill") }
                .labelStyle(.titleAndIcon)
                .prominentGlassButton()
                .disabled(lesson.isBusy || !lesson.canSubmit)
                .help("Submit (⌘↩)")
        }
    }
}

/// The Review screen's place in the window. It hosts the review views (Views/Teaching) with the
/// session's values and actions; until a session has loaded, it waits in place.
struct ReviewScreen: View {
    let session: ReviewSession

    var body: some View {
        if session.isLoaded {
            ReviewView(questions: session.questions.map { ReviewView.Question(ref: $0.ref, lessonTitle: $0.lessonTitle, question: $0.question) },
                       lessons: session.lessons.map {
                           ReviewView.Lesson(ref: $0.ref, title: $0.title, type: $0.type,
                                             reason: ReviewView.Lesson.reason(viewedSolution: $0.viewedSolution, failedSubmits: $0.failedSubmits))
                       },
                       onAnswer: { question, answer, right in
                           session.answer(question.ref, position: question.position, answer: answer, correct: right)
                       },
                       onOpenLesson: { session.redo($0.ref) })
                // A new session (Review chosen again) starts the view over.
                .id(ObjectIdentifier(session))
        } else {
            ProgressView().appearsAfterDelay().frame(maxWidth: .infinity, maxHeight: .infinity)
        }
    }
}
