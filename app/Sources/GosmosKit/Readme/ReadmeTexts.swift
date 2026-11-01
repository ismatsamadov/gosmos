import Foundation

/// Each lesson's readme as plain text, as the app renders it (no Markdown marks, callout or
/// system markers), for search. Made once per lesson, off the main actor. Libraries are
/// read-only, so a text never goes stale.
public actor ReadmeTexts {
    private var texts: [LessonRef: String] = [:]

    public init() {}

    /// Makes the text of every lesson of a library ahead of time, so search is quick from the
    /// start.
    public func warm(_ database: AppDatabase) async throws {
        let library = database.info.uuid
        for lesson in try await database.readmes() {
            let ref = LessonRef(library: library, lesson: lesson.id)
            if texts[ref] == nil { texts[ref] = ReadmeTexts.text(lesson.readme, title: lesson.title) }
        }
    }

    /// Searches the libraries, in order: each one's best hits, up to `limit` in all, kept only
    /// when the reader can see a match, in the title or in the readme as the app renders it;
    /// each comes with an excerpt around its first match. The index also matches words that
    /// never show, such as a code block's language, a callout's marker or a link's address.
    public func search(_ query: String, in databases: [AppDatabase], limit: Int = 60) async throws
        -> (hits: [SearchHit], snippets: [LessonRef: [Snippet.Part]]) {
        var hits: [SearchHit] = [], snippets: [LessonRef: [Snippet.Part]] = [:]
        for database in databases where hits.count < limit {
            let library = database.info.uuid
            // Twice as many as shown, so the hits dropped leave room for the next best.
            let candidates = try await database.search(query, limit: 2 * (limit - hits.count))
            let missing = candidates.map(\.ref.lesson).filter { texts[LessonRef(library: library, lesson: $0)] == nil }
            var unread: [String: (id: String, title: String, readme: String)] = [:]
            if !missing.isEmpty {
                for lesson in try await database.readmes(missing) { unread[lesson.id] = lesson }
            }
            for hit in candidates where hits.count < limit {
                if let lesson = unread.removeValue(forKey: hit.ref.lesson) { texts[hit.ref] = ReadmeTexts.text(lesson.readme, title: lesson.title) }
                let text = texts[hit.ref] ?? ""
                guard Snippet.shows(query, in: hit.title) || Snippet.shows(query, in: text) else { continue }
                hits.append(hit)
                snippets[hit.ref] = Snippet.make(text, query: query)
            }
        }
        return (hits, snippets)
    }

    static func text(_ readme: String, title: String) -> String {
        Readme.plainText(Readme.parse(readme, droppingTitle: title))
    }
}
