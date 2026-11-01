import Foundation

// What each kind of lesson holds besides its readme, as read from the per-type tables of
// db/schema.sql (format 8). Lists keep the tables' position order. Every kind is checked on
// this Mac alone: code by running it with your Go, a quiz by its answers.

/// A lesson's content: one case per lesson type, built from that type's tables.
public enum LessonContent: Sendable {
    case code(CodeExercise)
    case codeTests(CodeExercise)
    case choice(Quiz)
}

/// A file of a code lesson as the lesson ships it.
public struct StarterFile: Hashable, Sendable {
    public let name: String
    public let content: String
    public let isHidden: Bool
    public let isReadonly: Bool

    public init(name: String, content: String, isHidden: Bool = false, isReadonly: Bool = false) {
        self.name = name
        self.content = content
        self.isHidden = isHidden
        self.isReadonly = isReadonly
    }

    /// Shown in the editor; hidden files never are.
    public var isVisible: Bool { !isHidden }
}

/// type_code (checked by what the program prints) and type_code_tests (checked by its tests).
public struct CodeExercise: Sendable {
    public let progLang: String
    public let starterFiles: [StarterFile]
    /// What the program must print; set for type_code only.
    public let expectedOutput: String?
    /// The reference solution's version of each file you edit, in the files' order. Shown
    /// beside your code on request, never put in the editor.
    public let solutionFiles: [SolutionFile]
    /// A few sentences on why the solution works.
    public let solutionNote: String?

    /// Whether the app can run the code: it runs Go with your toolchain.
    public var runsHere: Bool { progLang == "go" }

    /// What the console says about a lesson in another language.
    public var otherLanguageNote: String { "This lesson uses \(progLang); Gosmos runs Go lessons." }

    public init(progLang: String = "go", starterFiles: [StarterFile], expectedOutput: String? = nil,
                solutionFiles: [SolutionFile] = [], solutionNote: String? = nil) {
        self.progLang = progLang
        self.starterFiles = starterFiles
        self.expectedOutput = expectedOutput
        self.solutionFiles = solutionFiles
        self.solutionNote = solutionNote
    }

    public var visibleFiles: [StarterFile] { starterFiles.filter(\.isVisible) }
}

/// type_choice.
public struct Quiz: Sendable {
    public struct Question: Hashable, Sendable {
        /// quiz_question.position: the question's place in its lesson, and its key for review.
        public let position: Int
        public let question: String
        public let answers: [String]
        /// The one right answer.
        public let answer: String
        /// Why each answer is right or wrong, in the answers' order; nil where the lesson
        /// gives no reason.
        public let explanations: [String?]

        public init(position: Int, question: String, answers: [String], answer: String, explanations: [String?] = []) {
            self.position = position
            self.question = question
            self.answers = answers
            self.answer = answer
            self.explanations = explanations
        }

        /// The explanation of an answer, by its index.
        public func explanation(of answer: Int) -> String? {
            explanations.indices.contains(answer) ? explanations[answer] : nil
        }

        /// Checked on this Mac: the chosen text must be the right answer.
        public func isCorrect(_ choice: String) -> Bool { choice == answer }
    }

    public let questions: [Question]

    public init(questions: [Question]) { self.questions = questions }
}
