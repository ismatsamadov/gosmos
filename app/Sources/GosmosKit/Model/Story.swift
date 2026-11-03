import Foundation

/// The story a course is told in (docs/storytelling.md): the opening shown before its first
/// lesson, each chapter's briefing and debrief, and the field guide of its world's words. Every
/// part is optional; a course without a story has none of them.
public struct CourseStory: Sendable, Hashable {
    public enum PageKind: String, Sendable, Hashable {
        case text
        case badge
    }

    /// One of the opening's pages. A badge page shows the learner's badge, made of its cards.
    public struct Page: Identifiable, Sendable, Hashable {
        public let position: Int
        public let kind: PageKind
        public let eyebrow: String?
        public let title: String
        /// Markdown paragraphs, one blank line between them.
        public let body: String
        /// The button that moves on.
        public let action: String
        public let cards: [Card]

        public init(position: Int, kind: PageKind, eyebrow: String?, title: String, body: String, action: String, cards: [Card] = []) {
            self.position = position
            self.kind = kind
            self.eyebrow = eyebrow
            self.title = title
            self.body = body
            self.action = action
            self.cards = cards
        }

        public var id: Int { position }
    }

    /// A person, a place or a fact on a page; on a badge page, one of the badge's lines.
    public struct Card: Identifiable, Sendable, Hashable {
        public let position: Int
        /// A short mark: initials, a station code, a line's name.
        public let label: String
        public let title: String
        public let subtitle: String?
        public let body: String?

        public init(position: Int, label: String, title: String, subtitle: String? = nil, body: String? = nil) {
            self.position = position
            self.label = label
            self.title = title
            self.subtitle = subtitle
            self.body = body
        }

        public var id: Int { position }
    }

    /// A chapter's part in the story.
    public struct ChapterStory: Sendable, Hashable {
        /// Opens the chapter.
        public let briefing: String
        /// Shown once its required lessons are done.
        public let debrief: String

        public init(briefing: String, debrief: String) {
            self.briefing = briefing
            self.debrief = debrief
        }
    }

    /// A word of the course's world with its meaning in plain words. A linked term links from a
    /// readme's first use of it.
    public struct Term: Identifiable, Sendable, Hashable {
        public let term: String
        public let meaning: String
        public let linked: Bool

        public init(term: String, meaning: String, linked: Bool = true) {
            self.term = term
            self.meaning = meaning
            self.linked = linked
        }

        public var id: String { term }
    }

    public let courseID: String
    /// The opening, in order; empty for a course that opens at its first lesson.
    public var pages: [Page]
    /// By chapter uuid.
    public var chapters: [String: ChapterStory]
    /// In the field guide's order: alphabetical, ignoring case.
    public var terms: [Term]

    public init(courseID: String, pages: [Page] = [], chapters: [String: ChapterStory] = [:], terms: [Term] = []) {
        self.courseID = courseID
        self.pages = pages
        self.chapters = chapters
        self.terms = terms
    }

    public var hasOpening: Bool { !pages.isEmpty }

    /// The page whose cards make the learner's badge.
    public var badge: Page? { pages.first { $0.kind == .badge } }

    /// A term by its exact spelling.
    public func term(_ name: String) -> Term? { terms.first { $0.term == name } }

    /// The paragraphs of a story text: its runs of lines between blank lines.
    public static func paragraphs(_ text: String) -> [String] {
        text.components(separatedBy: "\n\n").map { $0.trimmingCharacters(in: .newlines) }.filter { !$0.isEmpty }
    }
}
