-- The quiz goes with its questions and answers, so only the lesson is left without its row.
PRAGMA foreign_keys = ON;
DELETE FROM quiz_lesson WHERE lesson_uuid = 'quiz';
