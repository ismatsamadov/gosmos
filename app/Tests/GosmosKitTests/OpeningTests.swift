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

/// One copy of the app at a time: a launch gives way to a copy of the same bundle id that
/// launched before it, and never to another app such as a QA build.
@Suite struct SingleInstanceTests {
    @Test func aLaunchGivesWayToAnEarlierCopyOfTheSameApp() {
        let id = "app.gosmos.Gosmos"
        let t0 = Date(timeIntervalSince1970: 1_800_000_000)
        let me = SingleInstance.App(pid: 500, bundleIdentifier: id, launchDate: t0.addingTimeInterval(60))
        let installed = SingleInstance.App(pid: 400, bundleIdentifier: id, launchDate: t0)
        let qa = SingleInstance.App(pid: 300, bundleIdentifier: "app.gosmos.GosmosQA", launchDate: t0)

        // Alone, or beside another app: it runs. The system lists this process too.
        #expect(SingleInstance.runningCopy(for: me, in: []) == nil)
        #expect(SingleInstance.runningCopy(for: me, in: [me, qa]) == nil)
        // Another copy, such as ~/Applications/Gosmos.app while build/Gosmos.app launches.
        #expect(SingleInstance.runningCopy(for: me, in: [me, qa, installed]) == installed)
        // That copy, launched after this one, gives way instead: of two launched at once, one stays.
        #expect(SingleInstance.runningCopy(for: installed, in: [me, installed]) == nil)
        // Several earlier copies: the first.
        let older = SingleInstance.App(pid: 600, bundleIdentifier: id, launchDate: t0.addingTimeInterval(-60))
        #expect(SingleInstance.runningCopy(for: me, in: [installed, me, older]) == older)
        // The same launch date, or none known: the process id decides.
        let twin = SingleInstance.App(pid: 450, bundleIdentifier: id, launchDate: me.launchDate)
        #expect(SingleInstance.runningCopy(for: me, in: [me, twin]) == twin)
        #expect(SingleInstance.runningCopy(for: twin, in: [me, twin]) == nil)
        let undated = SingleInstance.App(pid: 100, bundleIdentifier: id)
        #expect(SingleInstance.runningCopy(for: SingleInstance.App(pid: 200, bundleIdentifier: id), in: [undated]) == undated)
        // A bare binary has no bundle id: nothing to compare.
        #expect(SingleInstance.runningCopy(for: SingleInstance.App(pid: 500, bundleIdentifier: nil), in: [SingleInstance.App(pid: 1, bundleIdentifier: nil)]) == nil)
    }
}
