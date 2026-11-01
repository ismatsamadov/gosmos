-- A Gosmos course library: one SQLite file of read-only content, such as courses/apogee.sqlite
-- in this repository. It holds the catalog (courses, chapters, lessons) and everything a lesson
-- teaches with: its readme, hints, code, quizzes and reference solutions. No JSON or other data
-- files. A content producer, such as cmd/author, writes it; the app only reads it. Everything
-- that is yours (progress, attempts, notes, your code, settings) lives in the learner database
-- instead (learner.sql), which refers to a lesson by the library's uuid (library_info.uuid) and
-- the lesson's (lesson.uuid). The file uses the rollback journal (journal_mode DELETE, never
-- WAL), so it is whole whenever no write is under way.
--
-- Conventions:
--   * Every table is STRICT, so a column only ever holds its declared type.
--   * Booleans are INTEGER 0 or 1. Timestamps are RFC 3339 in UTC with a trailing Z, such as
--     2026-10-31T12:00:00Z (fractional seconds allowed). An optional value that is missing is
--     NULL, not ''. Lists keep their order in a `position` column counted from 0.
--   * A row that belongs to one kind of parent carries that kind in a column pinned by CHECK
--     and references the parent's (key, kind), so a quiz_lesson row can only belong to a
--     type_choice lesson, and a code_lesson row only to a code lesson.
--   * Every lesson is one the app checks on the learner's Mac alone, with their Go toolchain:
--     code lessons by running them, quizzes by their answers. Gosmos never uses the network, so
--     format 8 removed the lesson types that reached outside the Mac or were only marked done by
--     hand: command-line lessons (shell commands, HTTP requests, jq checks), GitHub checks and
--     interviews.
--   * Every connection sets PRAGMA foreign_keys = ON; SQLite does not remember it.
--   * The uuids of the library, its courses and its lessons never change: the learner database
--     refers to them. A producer rewrites a lesson in place and never deletes one.
--   * This file makes a new, empty library; it is never run against an existing one. A change
--     raises user_version and comes with a migration in the db package (db/migrate.go), which
--     `library migrate` runs. The app opens only the version it was built for.
--   * No pictures, videos or interactive widgets: producers leave them out (cmd/author refuses
--     them), so none lands here.

PRAGMA foreign_keys = ON;

------------------------------------------------------------------------------------------------
-- The library
------------------------------------------------------------------------------------------------

-- What this file is: exactly one row (the library's checks see to it).
CREATE TABLE library_info (
  singleton INTEGER PRIMARY KEY CHECK (singleton = 1),
  uuid      TEXT NOT NULL UNIQUE CHECK (uuid <> ''),   -- the learner database's library.uuid
  title     TEXT NOT NULL CHECK (title <> ''),         -- such as Apogee
  version   TEXT NOT NULL CHECK (version <> '')        -- the content's version, such as 1
) STRICT;

------------------------------------------------------------------------------------------------
-- The catalog
------------------------------------------------------------------------------------------------

CREATE TABLE course (
  uuid              TEXT PRIMARY KEY,
  slug              TEXT NOT NULL UNIQUE,
  title             TEXT NOT NULL,
  kind              TEXT NOT NULL CHECK (kind IN ('Course', 'Guided Project')),
  short_description TEXT NOT NULL,
  description       TEXT NOT NULL,
  position          INTEGER NOT NULL UNIQUE CHECK (position >= 0)  -- order in the catalog
) STRICT;

CREATE TABLE chapter (
  uuid        TEXT PRIMARY KEY,
  course_uuid TEXT NOT NULL REFERENCES course(uuid) ON DELETE CASCADE,
  position    INTEGER NOT NULL CHECK (position >= 0),
  slug        TEXT NOT NULL,
  title       TEXT NOT NULL,
  description TEXT NOT NULL,
  UNIQUE (course_uuid, position),
  UNIQUE (course_uuid, slug)
) STRICT;

CREATE TABLE lesson (
  id            INTEGER PRIMARY KEY,           -- stable row id for lesson_search; VACUUM keeps it
  uuid          TEXT NOT NULL UNIQUE,
  chapter_uuid  TEXT NOT NULL REFERENCES chapter(uuid) ON DELETE CASCADE,
  position      INTEGER NOT NULL CHECK (position >= 0),
  slug          TEXT NOT NULL,
  title         TEXT NOT NULL,
  type          TEXT NOT NULL CHECK (type IN ('type_code', 'type_code_tests', 'type_choice')),
  required      INTEGER NOT NULL CHECK (required IN (0, 1)),
  difficulty    INTEGER NOT NULL CHECK (difficulty BETWEEN 1 AND 10),  -- 1 (easiest) to 10
  readme        TEXT NOT NULL,                                         -- Markdown, no pictures or videos
  created_at    TEXT NOT NULL CHECK (created_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(created_at) IS NOT NULL),
  last_mod      TEXT NOT NULL CHECK (last_mod GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(last_mod) IS NOT NULL),
  UNIQUE (chapter_uuid, position),
  UNIQUE (chapter_uuid, slug),
  UNIQUE (uuid, type)                          -- what the per-type tables reference
) STRICT;
-- created_at is when the lesson was first written, and last_mod when it last changed; the
-- author tool sets both.

CREATE TABLE lesson_topic (
  lesson_uuid TEXT NOT NULL REFERENCES lesson(uuid) ON DELETE CASCADE,
  topic       TEXT NOT NULL CHECK (topic <> ''),
  PRIMARY KEY (lesson_uuid, topic)
) STRICT;
CREATE INDEX lesson_topic_topic ON lesson_topic(topic);

-- Graduated hints, shown one at a time in order: each a short paragraph (inline Markdown) that
-- moves the learner closer without ever giving the answer away.
CREATE TABLE lesson_hint (
  lesson_uuid TEXT NOT NULL REFERENCES lesson(uuid) ON DELETE CASCADE,
  position    INTEGER NOT NULL CHECK (position >= 0),
  text        TEXT NOT NULL CHECK (text <> ''),
  PRIMARY KEY (lesson_uuid, position)
) STRICT;

-- Full-text search over lesson titles and readmes, kept in step by the triggers below.
CREATE VIRTUAL TABLE lesson_search USING fts5(
  title, readme, content = 'lesson', content_rowid = 'id', tokenize = 'unicode61 remove_diacritics 2'
);
CREATE TRIGGER lesson_search_insert AFTER INSERT ON lesson BEGIN
  INSERT INTO lesson_search(rowid, title, readme) VALUES (new.id, new.title, new.readme);
END;
CREATE TRIGGER lesson_search_delete AFTER DELETE ON lesson BEGIN
  INSERT INTO lesson_search(lesson_search, rowid, title, readme) VALUES ('delete', old.id, old.title, old.readme);
END;
CREATE TRIGGER lesson_search_update AFTER UPDATE OF title, readme ON lesson BEGIN
  INSERT INTO lesson_search(lesson_search, rowid, title, readme) VALUES ('delete', old.id, old.title, old.readme);
  INSERT INTO lesson_search(rowid, title, readme) VALUES (new.id, new.title, new.readme);
END;

------------------------------------------------------------------------------------------------
-- Code lessons: type_code (checked by what the program prints) and type_code_tests (checked by
-- its tests)
------------------------------------------------------------------------------------------------

CREATE TABLE code_lesson (
  lesson_uuid     TEXT PRIMARY KEY,
  lesson_type     TEXT NOT NULL CHECK (lesson_type IN ('type_code', 'type_code_tests')),
  prog_lang       TEXT NOT NULL CHECK (prog_lang <> ''),
  expected_output TEXT,                        -- type_code: what the program must print
  solution_note   TEXT CHECK (solution_note <> ''),  -- 1-3 sentences on why the solution works
  CHECK ((lesson_type = 'type_code') = (expected_output IS NOT NULL)),
  FOREIGN KEY (lesson_uuid, lesson_type) REFERENCES lesson(uuid, type) ON DELETE CASCADE
) STRICT;

CREATE TABLE starter_file (
  lesson_uuid TEXT NOT NULL REFERENCES code_lesson(lesson_uuid) ON DELETE CASCADE,
  name        TEXT NOT NULL CHECK (name <> ''),  -- path inside the lesson, e.g. main.go
  position    INTEGER NOT NULL CHECK (position >= 0),
  content     TEXT NOT NULL,
  is_hidden   INTEGER NOT NULL CHECK (is_hidden IN (0, 1)),    -- added at build time, never shown
  is_readonly INTEGER NOT NULL CHECK (is_readonly IN (0, 1)),  -- shown but not editable, like the tests
  PRIMARY KEY (lesson_uuid, name),
  UNIQUE (lesson_uuid, position)
) STRICT;

-- The reference solution's version of each file the learner edits (a starter file that is
-- neither hidden nor read-only). The app shows it, with code_lesson.solution_note, when the
-- learner asks to see the solution. Producers check their lessons with it (library verify):
-- the solution must pass and the starter must not.
CREATE TABLE solution_file (
  lesson_uuid TEXT NOT NULL REFERENCES code_lesson(lesson_uuid) ON DELETE CASCADE,
  name        TEXT NOT NULL CHECK (name <> ''),  -- the starter file it solves, e.g. main.go
  content     TEXT NOT NULL,
  PRIMARY KEY (lesson_uuid, name)
) STRICT;

------------------------------------------------------------------------------------------------
-- Quizzes: type_choice
------------------------------------------------------------------------------------------------

CREATE TABLE quiz_lesson (
  lesson_uuid           TEXT PRIMARY KEY,
  lesson_type           TEXT NOT NULL DEFAULT 'type_choice' CHECK (lesson_type = 'type_choice'),
  contains_complete_dir INTEGER NOT NULL CHECK (contains_complete_dir IN (0, 1)),
  FOREIGN KEY (lesson_uuid, lesson_type) REFERENCES lesson(uuid, type) ON DELETE CASCADE
) STRICT;

CREATE TABLE quiz_question (
  lesson_uuid TEXT NOT NULL REFERENCES quiz_lesson(lesson_uuid) ON DELETE CASCADE,
  position    INTEGER NOT NULL CHECK (position >= 0),
  question    TEXT NOT NULL,
  PRIMARY KEY (lesson_uuid, position)
) STRICT;

CREATE TABLE quiz_answer (
  lesson_uuid       TEXT NOT NULL,
  question_position INTEGER NOT NULL,
  position          INTEGER NOT NULL CHECK (position >= 0),
  answer            TEXT NOT NULL,
  is_correct        INTEGER NOT NULL CHECK (is_correct IN (0, 1)),
  explanation       TEXT CHECK (explanation <> ''),  -- why this answer is right or wrong
  PRIMARY KEY (lesson_uuid, question_position, position),
  UNIQUE (lesson_uuid, question_position, answer),
  FOREIGN KEY (lesson_uuid, question_position)
    REFERENCES quiz_question(lesson_uuid, position) ON DELETE CASCADE
) STRICT;
-- Exactly one right answer per question: this index allows at most one, and the library's
-- checks (db.Validate) make sure there is one.
CREATE UNIQUE INDEX quiz_answer_one_correct
  ON quiz_answer(lesson_uuid, question_position) WHERE is_correct = 1;

-- 'GSMS' marks the file as a Gosmos course library; user_version is the format's version.
PRAGMA application_id = 1196641619;
PRAGMA user_version = 8;
