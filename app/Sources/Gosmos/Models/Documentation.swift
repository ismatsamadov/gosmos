import AppKit
import GosmosKit
import Observation
import SwiftUI

/// The Go documentation browser: its list, the page on screen and the place on it, and its own
/// search. Pages are read and parsed off the main actor; the last few stay at hand, so Back is
/// instant.
@Observable @MainActor
final class DocsModel {
    let reference: ReferenceDatabase
    private(set) var outline = DocOutline.empty
    /// The place on screen: its page, and where on it.
    private(set) var location: DocLink?
    private(set) var page: DocPage?
    /// Why the page can't be shown.
    private(set) var loadError: String?
    /// Bumped to bring `location` into view, even when it is the same place again.
    private(set) var scrollRequests = 0
    var searchText = "" { didSet { searchChanged() } }
    private(set) var results: [DocSearchHit] = []
    private(set) var searching = false
    @ObservationIgnored private weak var app: AppModel?
    @ObservationIgnored private var pages: [DocLink: DocPage] = [:]
    @ObservationIgnored private var recent: [DocLink] = []
    @ObservationIgnored private var loadTask: Task<Void, Never>?
    @ObservationIgnored private var searchTask: Task<Void, Never>?

    init(reference: ReferenceDatabase, app: AppModel) {
        self.reference = reference
        self.app = app
        Task { [weak self] in
            do {
                let outline = try await reference.outline()
                self?.outline = outline
            } catch {
                self?.app?.report(error)
            }
        }
    }

    /// The title of the page on screen.
    var title: String { page?.title ?? "Go Documentation" }

    /// Shows a place: its page, read unless it is on screen or at hand, scrolled to the place.
    /// AppModel.openDocumentation calls it, keeping the history.
    func show(_ link: DocLink) {
        location = link
        scrollRequests += 1
        let target = link.page
        // A page still loading is not wanted any more.
        loadTask?.cancel()
        guard page?.link != target else { return }
        if let kept = pages[target] {
            page = kept
            loadError = nil
            return
        }
        let reference = reference, started = ContinuousClock.now
        loadTask = Task { [weak self] in
            do {
                let found = try await reference.page(target)
                guard let self, !Task.isCancelled else { return }
                guard let found else {
                    self.loadError = "\(target.url) is not in the documentation."
                    self.page = nil
                    return
                }
                self.keep(found)
                self.page = found
                self.loadError = nil
                Log.app.notice("Timing doc page \(Int((ContinuousClock.now - started).milliseconds), privacy: .public) ms (\(found.title, privacy: .public))")
            } catch {
                guard let self, !Reporting.isSilent(error) else { return }
                self.loadError = error.localizedDescription
                self.page = nil
            }
        }
    }

    private func keep(_ page: DocPage) {
        pages[page.link] = page
        recent.removeAll { $0 == page.link }
        recent.append(page.link)
        if recent.count > 8 { pages[recent.removeFirst()] = nil }
    }

    private func searchChanged() {
        searchTask?.cancel()
        let text = searchText
        guard !text.trimmingCharacters(in: .whitespaces).isEmpty else {
            results = []
            searching = false
            return
        }
        searching = true
        let reference = reference
        searchTask = Task { [weak self] in
            let started = ContinuousClock.now
            let found: [DocSearchHit]
            do {
                found = try await reference.search(text)
            } catch {
                guard !Reporting.isSilent(error) else { return }
                found = []
                self?.app?.report(error)
            }
            guard let self, !Task.isCancelled else { return }
            self.results = found
            self.searching = false
            Log.app.notice("Timing doc search \(String(format: "%.1f", (ContinuousClock.now - started).milliseconds), privacy: .public) ms")
        }
    }
}

// MARK: - Quick Help

extension AppModel {
    /// Whether Show Documentation can act: a Go file in the editor, with the documentation at
    /// hand and nothing over the editor.
    var canShowQuickHelp: Bool {
        guard docs != nil, let lesson, let code = lesson.code, !lesson.showingSolution, code.conflict == nil else { return false }
        return code.activeFile?.name.hasSuffix(".go") == true
    }

    /// Quick Help: the declaration of the package or name at `offset` in the editor's file, or
    /// at the insertion point, with the first paragraph of its doc, in a popover beside it.
    /// The file's imports say which package a name is in. A name the standard library does not
    /// have gets the system's beep.
    func showQuickHelp(_ session: CodeSession, at offset: Int? = nil) {
        guard let docs, let text = session.editor.text, let at = offset ?? session.editor.selectedRange?.location,
              let selector = GoSource.selector(at: at, in: text),
              let link = docNames.link(for: selector, imports: GoSource.imports(in: text)) else {
            NSSound.beep()
            return
        }
        Task { [weak self] in
            let summary: DocSummary?
            do {
                summary = try await docs.reference.summary(link)
            } catch {
                self?.report(error)
                return
            }
            guard let summary else {
                NSSound.beep()
                return
            }
            let popover = NSPopover()
            popover.behavior = .transient
            let open = { [weak self, weak popover] (link: DocLink) in
                popover?.performClose(nil)
                self?.openDocumentation(link)
            }
            let content = QuickHelpView(summary: summary) { open(link) }
                .environment(\.openURL, OpenURLAction { url in
                    guard let target = DocLink(url) else { return .discarded }
                    open(target)
                    return .handled
                })
            let host = NSHostingController(rootView: content)
            host.sizingOptions = .preferredContentSize
            popover.contentViewController = host
            if !session.editor.show(popover, at: selector.range) { NSSound.beep() }
        }
    }
}
