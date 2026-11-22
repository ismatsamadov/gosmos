-- A small library of every lesson type, which every check passes. Each case in this folder
-- changes it (<case>.sql, run with foreign keys off, outside a transaction) or replaces the
-- file (<case>.file), and <case>.want holds what the checks report, one problem a line.
INSERT INTO library_info (singleton, uuid, title, version) VALUES (1, 'b0000000-0000-0000-0000-000000000001', 'Practice', '1');
INSERT INTO course (uuid, slug, title, kind, short_description, description, position)
  VALUES ('course', 'learn-testing', 'Learn Testing', 'Course', 'Short.', 'Long.', 0);
INSERT INTO chapter (uuid, course_uuid, position, slug, title, description)
  VALUES ('chapter', 'course', 0, 'basics', 'Basics', 'The basics.');
INSERT INTO lesson (uuid, chapter_uuid, position, slug, title, type, required, difficulty, readme, created_at, last_mod)
  VALUES ('hello', 'chapter', 0, 'hello', 'Hello', 'type_code', 1, 1, '# Hello', '2026-10-01T09:00:00Z', '2026-10-01T09:00:00Z'),
         ('add', 'chapter', 1, 'add', 'Adding Numbers', 'type_code_tests', 1, 2, '# Adding', '2026-10-01T09:00:00Z', '2026-10-01T09:00:00Z'),
         ('quiz', 'chapter', 2, 'quiz', 'A Quiz', 'type_choice', 1, 1, '# A Quiz', '2026-10-01T09:00:00Z', '2026-10-01T09:00:00Z');
INSERT INTO lesson_topic (lesson_uuid, topic) VALUES ('hello', 'printing'), ('add', 'functions');
INSERT INTO code_lesson (lesson_uuid, lesson_type, prog_lang, expected_output)
  VALUES ('hello', 'type_code', 'go', 'Hello'), ('add', 'type_code_tests', 'go', NULL);
INSERT INTO starter_file (lesson_uuid, name, position, content, is_hidden, is_readonly)
  VALUES ('hello', 'main.go', 0, 'package main', 0, 0),
         ('add', 'main.go', 0, 'package main', 0, 0),
         ('add', 'main_test.go', 1, 'package main', 0, 1);
INSERT INTO solution_file (lesson_uuid, name, content) VALUES ('hello', 'main.go', 'package main'), ('add', 'main.go', 'package main');
INSERT INTO quiz_lesson (lesson_uuid, lesson_type, contains_complete_dir) VALUES ('quiz', 'type_choice', 0);
INSERT INTO quiz_question (lesson_uuid, position, question) VALUES ('quiz', 0, 'Is Go compiled?'), ('quiz', 1, 'Which keyword declares a function?');
INSERT INTO quiz_answer (lesson_uuid, question_position, position, answer, is_correct, explanation)
  VALUES ('quiz', 0, 0, 'Yes', 1, NULL), ('quiz', 0, 1, 'No', 0, NULL),
         ('quiz', 1, 0, 'func', 1, NULL), ('quiz', 1, 1, 'def', 0, NULL);
