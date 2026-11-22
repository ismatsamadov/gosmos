import GosmosKit
import Testing
@testable import Gosmos

/// The sidebar's search: what pasted text searches for, only the latest search landing, and
/// what clearing it leaves.
@Suite(.serialized) @MainActor struct SearchModelTests {
    /// The characters that lay out a readme's inline code go, even with nothing to search.
    @Test func pastedTextLosesItsLayoutCharacters() {
        let search = SearchModel()
        search.text = "Print\u{2060} hel\u{200B}lo\u{FEFF}"
        #expect(search.text == "Print hello")
        #expect(!search.searching && search.results.isEmpty)
    }

    @Test func onlyTheLatestSearchLands() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let search = app.search
        search.text = "hello"
        search.text = "greeting"
        #expect(search.searching)
        try await fixture.wait { !search.searching }
        try await Task.sleep(for: .milliseconds(100))
        #expect(search.results.map(\.ref) == [fixture.ref("optional")])
        #expect(search.snippets[fixture.ref("optional")]?.contains { $0.matched && $0.text == "greeting" } == true)
        await app.prepareToQuit()
    }

    /// Cleared, the search shows nothing, and a place in the documentation it found stays
    /// selected as the documentation, the sidebar's row for it.
    @Test func clearingLeavesNoResultsAndTheDocumentationSelected() async throws {
        let fixture = try FlowFixture()
        defer { fixture.remove() }
        let app = fixture.app()
        await app.start()
        let search = app.search
        search.text = "hello"
        try await fixture.wait { !search.searching }
        #expect(Set(search.results.map(\.ref)) == [fixture.ref("first"), fixture.ref("second")])
        app.selection = .doc(.package("fmt"))
        search.text = "  "
        #expect(search.results.isEmpty && search.docResults.isEmpty && !search.searching)
        #expect(app.selection == .documentation)
        // A lesson found stays selected.
        search.text = "hello"
        app.selection = .lesson(fixture.ref("first"))
        search.text = ""
        #expect(!search.searching && app.selection == .lesson(fixture.ref("first")))
        await app.prepareToQuit()
    }
}
