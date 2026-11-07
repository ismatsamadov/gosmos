import Foundation
import GRDB
import Testing
@testable import GosmosKit

/// A course's story, as the app reads it from a library of format 9 (docs/storytelling.md):
/// its opening's pages and their cards, its chapters' briefings and debriefs, and its field
/// guide, whose terms link from readmes.
@Suite struct StoryTests {
    /// The fixture library with a story for its first course: two pages, the second a badge,
    /// a briefing and a debrief for its first chapter, and three terms.
    static func storied() throws -> (AppDatabase, URL) {
        try Fixtures.database { db in
            let course = Fixtures.courseID
            try db.execute(sql: "INSERT INTO story_page(course_uuid, position, kind, eyebrow, title, body, action) VALUES (?,?,?,?,?,?,?)",
                           arguments: [course, 0, "text", "PORT ELLERY · DAY ONE", "Welcome", "Seven satellites.\n\nToday you join us.", "Start my first day"])
            try db.execute(sql: "INSERT INTO story_page(course_uuid, position, kind, eyebrow, title, body, action) VALUES (?,?,?,?,?,?,?)",
                           arguments: [course, 1, "badge", "APOGEE", "Your Badge", "Everyone wears one.", "Print my badge"])
            try db.execute(sql: "INSERT INTO story_card(course_uuid, page_position, position, label, title, subtitle, body) VALUES (?,?,?,?,?,?,?)",
                           arguments: [course, 0, 0, "NH", "Nkem Halloran", "Your buddy", "Reviews your code."])
            for (i, line) in [("Role", "Ground Software Engineer"), ("Badge", "No. 039")].enumerated() {
                try db.execute(sql: "INSERT INTO story_card(course_uuid, page_position, position, label, title) VALUES (?,?,?,?,?)",
                               arguments: [course, 1, i, line.0, line.1])
            }
            try db.execute(sql: "INSERT INTO chapter_story(chapter_uuid, briefing, debrief) VALUES (?,?,?)",
                           arguments: [Fixtures.chapter1, "Day one.", "Your words came back."])
            for (term, meaning, linked) in [("pass", "One flyover.", 0), ("AOS", "The start of a pass.", 1), ("slant range", "The straight-line distance.", 1)] {
                try db.execute(sql: "INSERT INTO story_term(course_uuid, term, meaning, linked) VALUES (?,?,?,?)", arguments: [course, term, meaning, linked])
            }
        }
    }

    @Test func readsEveryPartOfAStory() async throws {
        let (db, dir) = try Self.storied()
        defer { try? FileManager.default.removeItem(at: dir) }
        let stories = try await db.stories()
        #expect(Array(stories.keys) == [Fixtures.courseID]) // the guided project has none
        let story = try #require(stories[Fixtures.courseID])
        #expect(story.hasOpening && story.pages.map(\.title) == ["Welcome", "Your Badge"])
        #expect(story.pages[0].kind == .text && story.pages[0].eyebrow == "PORT ELLERY · DAY ONE" && story.pages[0].action == "Start my first day")
        #expect(CourseStory.paragraphs(story.pages[0].body) == ["Seven satellites.", "Today you join us."])
        #expect(story.pages[0].cards == [CourseStory.Card(position: 0, label: "NH", title: "Nkem Halloran", subtitle: "Your buddy", body: "Reviews your code.")])
        #expect(story.badge?.position == 1 && story.badge?.cards.map(\.title) == ["Ground Software Engineer", "No. 039"])
        #expect(story.badge?.cards.allSatisfy { $0.subtitle == nil && $0.body == nil } == true)
        #expect(story.chapters == [Fixtures.chapter1: CourseStory.ChapterStory(briefing: "Day one.", debrief: "Your words came back.")])
        // The field guide's order: alphabetical, ignoring case.
        #expect(story.terms.map(\.term) == ["AOS", "pass", "slant range"])
        #expect(story.term("pass")?.linked == false && story.term("AOS")?.linked == true && story.term("aos") == nil)
    }

    @Test func aLibraryWithoutAStoryHasNone() async throws {
        let (db, dir) = try Fixtures.database()
        defer { try? FileManager.default.removeItem(at: dir) }
        #expect(try await db.stories().isEmpty)
    }

    static let terms = [
        CourseStory.Term(term: "AOS", meaning: "The start of a pass."),
        CourseStory.Term(term: "pass", meaning: "One flyover.", linked: false),
        CourseStory.Term(term: "range", meaning: "Distance."),
        CourseStory.Term(term: "slant range", meaning: "The straight-line distance."),
        CourseStory.Term(term: "beacon", meaning: "A short status message."),
    ]

    /// The links a readme's blocks hold, in order, with the text each one links.
    static func links(_ blocks: [MDBlock]) -> [(String, String)] {
        func inlines(_ items: [MDInline]) -> [(String, String)] {
            items.flatMap { item -> [(String, String)] in
                switch item {
                case .link(let destination, _, let children):
                    [(destination, children.map { if case .text(let t) = $0 { t } else { "" } }.joined())]
                case .emphasis(let c), .strong(let c), .strike(let c), .superscript(let c): inlines(c)
                default: []
                }
            }
        }
        return blocks.flatMap { block -> [(String, String)] in
            switch block {
            case .paragraph(let i), .heading(_, _, let i): inlines(i)
            case .quote(_, let b): links(b)
            case .list(_, _, _, let items): items.flatMap { links($0.blocks) }
            default: []
            }
        }
    }

    /// The first use of a linked term in prose links to it, once a readme; code, headings and
    /// links keep their text, an everyday word is not linked, and the longer term wins.
    @Test func linksTheFirstUseOfEachTerm() {
        let readme = """
        ## The AOS heading

        Pipit-1 rises, its pass begins at `AOS`, and the beacons come in. AOS again, and AOS.

        > [!NOTE]
        > The slant range shrinks to the range at the top. [AOS](godoc:fmt.Println) is a link.

        - Tests pass when the beacon is heard.
        """
        let linked = TermLinker.link(Readme.parse(readme), terms: Self.terms)
        let found = Self.links(linked)
        #expect(found.map(\.0) == ["gosmos-term:beacon", "gosmos-term:AOS", "gosmos-term:slant%20range", "gosmos-term:range", "godoc:fmt.Println"])
        #expect(found.map(\.1) == ["beacons", "AOS", "slant range", "range", "AOS"])
        // The heading, the code span and the unlinked word stay as they were.
        #expect(linked[0] == Readme.parse(readme)[0])
        // Without terms, nothing changes.
        #expect(TermLinker.link(Readme.parse(readme), terms: []) == Readme.parse(readme))
    }

    /// A word that is a term of its own is never a shorter term with an ending, even once its
    /// own term has been linked, or when it is never linked: the second "uplinkd" links to
    /// nothing, not to "uplink".
    @Test func aTermWithAnEndingIsNotAnotherTerm() {
        let terms = [CourseStory.Term(term: "uplink", meaning: "Ground to space."),
                     CourseStory.Term(term: "uplinkd", meaning: "The service that sends commands up.")]
        let readme = "uplinkd signs each command. Then uplinkd sends it up the uplink."
        let found = Self.links(TermLinker.link(Readme.parse(readme), terms: terms))
        #expect(found.map(\.0) == ["gosmos-term:uplinkd", "gosmos-term:uplink"])
        #expect(found.map(\.1) == ["uplinkd", "uplink"])
        let unlinked = [terms[0], CourseStory.Term(term: "uplinkd", meaning: "The service that sends commands up.", linked: false)]
        let alone = Self.links(TermLinker.link(Readme.parse(readme), terms: unlinked))
        #expect(alone.map(\.0) == ["gosmos-term:uplink"])
        #expect(alone.map(\.1) == ["uplink"])
    }

    @Test func termsMatchWholeWordsOnly() {
        #expect(TermLinker.wholeWord("AOS", in: "CHAOS reigns") == nil)
        #expect(TermLinker.wholeWord("pass", in: "a password") == nil)
        let text = "two passes, then a pass"
        #expect(TermLinker.wholeWord("pass", in: text).map { String(text[$0]) } == "passes")
        let possessive = "the beacon's echo"
        #expect(TermLinker.wholeWord("beacon", in: possessive).map { String(possessive[$0]) } == "beacon")
        #expect(TermLinker.wholeWord("AOS", in: "aos") == nil) // as the guide spells it
        let verbs = "Both dishes stowed, then tracking resumed"
        #expect(TermLinker.wholeWord("stow", in: verbs).map { String(verbs[$0]) } == "stowed")
        #expect(TermLinker.wholeWord("track", in: verbs).map { String(verbs[$0]) } == "tracking")
        #expect(TermLinker.wholeWord("stow", in: "a stowaway") == nil)
        // A lowercase term also starts a sentence; an acronym is only ever itself.
        #expect(TermLinker.spellings(of: "the Board") == ["the Board", "The Board"])
        #expect(TermLinker.spellings(of: "AOS") == ["AOS"])
        let linked = TermLinker.link(Readme.parse("The Board lists the next passes."), terms: [CourseStory.Term(term: "the Board", meaning: "The big screen.")])
        #expect(Self.links(linked).map(\.0) == ["gosmos-term:the%20Board"] && Self.links(linked).map(\.1) == ["The Board"])
    }

    /// A term's link carries the app's own address and leads back to the term, and only to a
    /// term of the course's field guide.
    @Test func termLinksLeadBackToTheTerm() throws {
        let target = LinkTarget.term("slant range")
        let url = try #require(URL(string: target.url))
        #expect(url.scheme == LinkTarget.termScheme)
        #expect(LinkTarget(destination: url.absoluteString, isTerm: { $0 == "slant range" }) == target)
        #expect(LinkTarget(destination: url.absoluteString) == nil)
    }
}
