import CodeEditor
import GosmosKit
import Observation
import os

/// The sidebar's search: the lessons of every library on the shelf whose title or readme shows
/// what you type, each with an excerpt, then the documentation's best matches. Each change of
/// the text starts a new search, and only the latest one lands.
@Observable @MainActor
final class SearchModel {
    var text = "" {
        didSet {
            // Text copied from a readme carries the characters that lay out its inline code
            // (PastedText); search for what it says.
            if PastedText.needsCleaning(text) { text = PastedText.clean(text) }
            changed()
        }
    }
    private(set) var results: [SearchHit] = []
    /// Excerpts for the results, by lesson.
    private(set) var snippets: [LessonRef: [Snippet.Part]] = [:]
    /// What search found in the documentation, below the lessons.
    private(set) var docResults: [DocSearchHit] = []
    private(set) var searching = false
    /// Readmes as plain text, for search excerpts; made ahead of time for every library.
    let texts = ReadmeTexts()
    /// The app whose shelf and documentation are searched; it sets itself here when it is made.
    @ObservationIgnored weak var app: AppModel?
    @ObservationIgnored private var task: Task<Void, Never>?
    private let signposter = OSSignposter(subsystem: Log.subsystem, category: .pointsOfInterest)

    private func changed() {
        task?.cancel()
        let text = text
        let databases = app?.shelf.libraries.compactMap(\.database) ?? []
        let reference = app?.docs?.reference
        guard !databases.isEmpty || reference != nil, !text.trimmingCharacters(in: .whitespaces).isEmpty else {
            results = []
            docResults = []
            searching = false
            // The place search found stays on screen, as the documentation.
            if let app, case .doc? = app.selection { app.selection = .documentation }
            return
        }
        searching = true
        let texts = texts
        task = Task { [weak self] in
            let interval = self?.signposter.beginInterval("Search")
            let started = ContinuousClock.now
            // The index's hits that show a match, with their excerpts, from the readmes as
            // the app renders them, through every library in the sidebar's order.
            let found: (hits: [SearchHit], snippets: [LessonRef: [Snippet.Part]])
            do {
                found = try await texts.search(text, in: databases)
            } catch {
                found = ([], [:])
                self?.app?.report(error)
            }
            // The documentation's best matches, below the lessons.
            var docHits: [DocSearchHit] = []
            if let reference {
                do { docHits = try await reference.search(text, limit: 8) } catch { self?.app?.report(error) }
            }
            guard let self, !Task.isCancelled else { return }
            self.results = found.hits
            self.docResults = docHits
            self.snippets = found.snippets
            self.searching = false
            if let interval { self.signposter.endInterval("Search", interval) }
            let ms = (ContinuousClock.now - started).milliseconds
            Log.app.notice("Timing search \(String(format: "%.1f", ms), privacy: .public) ms for \(text.count, privacy: .public) characters")
        }
    }
}

// MARK: - Search

extension AppModel {
    /// The sidebar's search field, and what it found.
    var searchText: String {
        get { search.text }
        set { search.text = newValue }
    }
    var searchResults: [SearchHit] { search.results }
    var searchSnippets: [LessonRef: [Snippet.Part]] { search.snippets }
    var docSearchResults: [DocSearchHit] { search.docResults }
    var searching: Bool { search.searching }
}
