import Foundation
import GRDB
import Testing
@testable import GosmosKit

/// Several libraries side by side: a lesson uuid can repeat across libraries, so every cache
/// keys lessons by library too, and each library keeps its own text.
@Suite struct OpeningTests {
    @Test func cachesKeepEachLibrarysOwnText() async throws {
        let (first, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        // The same lesson uuid in another library, with other text.
        let otherURL = dir.appendingPathComponent("other.library")
        try Fixtures.makeDatabase(at: otherURL, id: "b0000000-0000-0000-0000-000000000002", title: "Other")
        try await DatabaseQueue(path: otherURL.path).write {
            try $0.execute(sql: "UPDATE lesson SET readme = ? WHERE uuid = ?", arguments: ["# Hello Gopher\n\nA new kettle.", Fixtures.codeLesson])
        }
        let second = try AppDatabase(url: otherURL)
        let texts = ReadmeTexts(), cache = ReadmeCache()
        let lesson = try #require(try await first.lesson(Fixtures.codeLesson))
        let renewed = try #require(try await second.lesson(Fixtures.codeLesson))
        #expect(lesson.ref != renewed.ref && lesson.id == renewed.id)

        #expect(try await texts.search("greeting", in: [first, second]).hits.map(\.ref) == [lesson.ref])
        let kettle = try await texts.search("kettle", in: [first, second])
        #expect(kettle.hits.map(\.ref) == [renewed.ref])
        #expect(kettle.snippets[renewed.ref]?.contains { $0.matched && $0.text == "kettle" } == true)

        let (old, new) = (await cache.parsed(for: lesson), await cache.parsed(for: renewed))
        #expect(old.readme.first != new.readme.first)
        #expect(await cache.parsed(for: renewed).readme == [.paragraph([.text("A new kettle.")])])
    }
}
