import Foundation
import GRDB

// A lesson and its content, read from the per-type tables of format 8, in position order.
// The app never writes a library.

extension AppDatabase {
    static func fetchLesson(_ db: Database, _ uuid: String, library: String) throws -> LessonDetail? {
        guard let row = try Row.fetchOne(db, sql: """
            SELECT l.uuid, l.chapter_uuid, l.title, l.type, l.required, l.readme, l.difficulty,
                   c.uuid AS course_uuid, c.title AS course_title, ch.title AS chapter_title
            FROM lesson l JOIN chapter ch ON ch.uuid = l.chapter_uuid JOIN course c ON c.uuid = ch.course_uuid
            WHERE l.uuid = ?
            """, arguments: [uuid]), let type = LessonType(rawValue: row["type"]) else { return nil }
        let summary = LessonSummary(id: row["uuid"], chapterID: row["chapter_uuid"], title: row["title"], type: type, required: row["required"] as Int != 0)
        let topics = try String.fetchAll(db, sql: "SELECT topic FROM lesson_topic WHERE lesson_uuid = ? ORDER BY topic", arguments: [uuid])
        let hints = try Row.fetchAll(db, sql: "SELECT position, text FROM lesson_hint WHERE lesson_uuid = ? ORDER BY position", arguments: [uuid]).map {
            Hint(position: $0["position"], text: $0["text"])
        }
        let content: LessonContent
        switch type {
        case .code: content = .code(try fetchCode(db, lesson: uuid))
        case .codeTests: content = .codeTests(try fetchCode(db, lesson: uuid))
        case .choice: content = .choice(try fetchQuiz(db, uuid))
        }
        return LessonDetail(library: library, summary: summary, courseID: row["course_uuid"], courseTitle: row["course_title"],
                            chapterTitle: row["chapter_title"], readme: row["readme"], difficulty: row["difficulty"], topics: topics,
                            hints: hints, content: content)
    }

    static func fetchCode(_ db: Database, lesson uuid: String) throws -> CodeExercise {
        let head = try Row.fetchOne(db, sql: "SELECT prog_lang, expected_output, solution_note FROM code_lesson WHERE lesson_uuid = ?", arguments: [uuid])
        let files = try Row.fetchAll(db, sql: """
            SELECT name, content, is_hidden, is_readonly FROM starter_file WHERE lesson_uuid = ? ORDER BY position
            """, arguments: [uuid]).map {
            StarterFile(name: $0["name"], content: $0["content"], isHidden: $0["is_hidden"] as Int != 0, isReadonly: $0["is_readonly"] as Int != 0)
        }
        let solution = try Row.fetchAll(db, sql: """
            SELECT s.name, s.content FROM solution_file s JOIN starter_file f ON f.lesson_uuid = s.lesson_uuid AND f.name = s.name
            WHERE s.lesson_uuid = ? ORDER BY f.position
            """, arguments: [uuid]).map { SolutionFile(name: $0["name"], content: $0["content"]) }
        return CodeExercise(progLang: head?["prog_lang"] ?? "go", starterFiles: files, expectedOutput: head?["expected_output"],
                            solutionFiles: solution, solutionNote: head?["solution_note"])
    }

    static func fetchQuiz(_ db: Database, _ uuid: String) throws -> Quiz {
        let answers = try Row.fetchAll(db, sql: """
            SELECT question_position, answer, is_correct, explanation FROM quiz_answer WHERE lesson_uuid = ? ORDER BY question_position, position
            """, arguments: [uuid])
        let questions = try Row.fetchAll(db, sql: "SELECT position, question FROM quiz_question WHERE lesson_uuid = ? ORDER BY position", arguments: [uuid])
        return Quiz(questions: questions.map { q in
            let position: Int = q["position"]
            let own = answers.filter { ($0["question_position"] as Int) == position }
            return Quiz.Question(position: position, question: q["question"], answers: own.map { $0["answer"] },
                                 answer: own.first(where: { ($0["is_correct"] as Int) != 0 }).map { $0["answer"] } ?? "",
                                 explanations: own.map { $0["explanation"] })
        })
    }
}
