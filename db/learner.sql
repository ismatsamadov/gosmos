-- Gosmos's learner database: one SQLite file per Mac, ~/Library/Application Support/Gosmos/
-- learner.sqlite. It holds everything that is yours: the libraries you use, your progress,
-- every attempt, notes, your code, the hints and solutions you looked at, and your review
-- schedule. Course libraries (schema.sql) stay read-only; a lesson is referenced here by
-- (library_uuid, lesson_uuid), since a foreign key cannot reach into another file.
--
-- Conventions, as in schema.sql:
--   * Every table is STRICT. Booleans are INTEGER 0 or 1. Timestamps are RFC 3339 in UTC
--     with a trailing Z. An optional value that is missing is NULL, not ''.
--   * Every connection sets PRAGMA foreign_keys = ON.
--   * This file makes a new, empty learner database. A change raises user_version and comes
--     with a migration; the app applies them when it opens an older file (LearnerDatabase).
--     Version 2 dropped what only the lesson types format 8 removed used: course project
--     folders, server addresses and command-line check attempts.

PRAGMA foreign_keys = ON;

-- Every library you have studied in: built into the app, or opened from a file.
CREATE TABLE library (
  uuid     TEXT PRIMARY KEY CHECK (uuid <> ''),         -- the library's library_info.uuid
  title    TEXT NOT NULL CHECK (title <> ''),
  location TEXT NOT NULL CHECK (location <> ''),        -- 'builtin:<file name>' or an absolute path
  added_at TEXT NOT NULL CHECK (added_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(added_at) IS NOT NULL)
) STRICT;

CREATE TABLE progress (
  library_uuid   TEXT NOT NULL REFERENCES library(uuid),
  lesson_uuid    TEXT NOT NULL CHECK (lesson_uuid <> ''),
  done_at        TEXT CHECK (done_at IS NULL OR (done_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(done_at) IS NOT NULL)),  -- NULL until done
  last_opened_at TEXT CHECK (last_opened_at IS NULL OR (last_opened_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(last_opened_at) IS NOT NULL)),
  CHECK (done_at IS NOT NULL OR last_opened_at IS NOT NULL),
  PRIMARY KEY (library_uuid, lesson_uuid)
) STRICT;
CREATE INDEX progress_last_opened ON progress(last_opened_at);

-- Every Run, Submit, quiz answer and review, passed or not.
CREATE TABLE attempt (
  id           INTEGER PRIMARY KEY,
  library_uuid TEXT NOT NULL REFERENCES library(uuid),
  lesson_uuid  TEXT NOT NULL CHECK (lesson_uuid <> ''),
  kind         TEXT NOT NULL CHECK (kind IN ('run', 'submit', 'answer', 'review')),
  passed       INTEGER NOT NULL CHECK (passed IN (0, 1)),
  output       TEXT NOT NULL,                           -- what the run printed, or the answer you picked
  duration_ms  INTEGER CHECK (duration_ms >= 0),
  at           TEXT NOT NULL CHECK (at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(at) IS NOT NULL)
) STRICT;
CREATE INDEX attempt_lesson ON attempt(library_uuid, lesson_uuid, at);

CREATE TABLE note (
  library_uuid TEXT NOT NULL REFERENCES library(uuid),
  lesson_uuid  TEXT NOT NULL CHECK (lesson_uuid <> ''),
  body         TEXT NOT NULL,
  updated_at   TEXT NOT NULL CHECK (updated_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(updated_at) IS NOT NULL),
  PRIMARY KEY (library_uuid, lesson_uuid)
) STRICT;

-- Your code, one row per file. The app writes a lesson's files to a scratch folder only to
-- build, run and test them, and reads them back from there.
CREATE TABLE code_file (
  library_uuid TEXT NOT NULL REFERENCES library(uuid),
  lesson_uuid  TEXT NOT NULL CHECK (lesson_uuid <> ''),
  name         TEXT NOT NULL CHECK (name <> ''),        -- path inside the lesson, e.g. main.go
  content      TEXT NOT NULL,
  updated_at   TEXT NOT NULL CHECK (updated_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(updated_at) IS NOT NULL),
  PRIMARY KEY (library_uuid, lesson_uuid, name)
) STRICT;

-- Each hint you revealed (lesson_hint.position), so the lesson shows it again.
CREATE TABLE hint_view (
  library_uuid TEXT NOT NULL REFERENCES library(uuid),
  lesson_uuid  TEXT NOT NULL CHECK (lesson_uuid <> ''),
  position     INTEGER NOT NULL CHECK (position >= 0),
  at           TEXT NOT NULL CHECK (at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(at) IS NOT NULL),
  PRIMARY KEY (library_uuid, lesson_uuid, position)
) STRICT;

-- The first time you looked at a lesson's reference solution.
CREATE TABLE solution_view (
  library_uuid TEXT NOT NULL REFERENCES library(uuid),
  lesson_uuid  TEXT NOT NULL CHECK (lesson_uuid <> ''),
  at           TEXT NOT NULL CHECK (at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(at) IS NOT NULL),
  PRIMARY KEY (library_uuid, lesson_uuid)
) STRICT;

-- Spaced review. A lesson comes back as a whole (review_lesson), a quiz question on its own
-- (review_question). due_at is when it is next due; interval_days grows with each pass and
-- starts over after a lapse (a failed review).
CREATE TABLE review_lesson (
  library_uuid     TEXT NOT NULL REFERENCES library(uuid),
  lesson_uuid      TEXT NOT NULL CHECK (lesson_uuid <> ''),
  due_at           TEXT NOT NULL CHECK (due_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(due_at) IS NOT NULL),
  interval_days    INTEGER NOT NULL CHECK (interval_days >= 1),
  reps             INTEGER NOT NULL DEFAULT 0 CHECK (reps >= 0),
  lapses           INTEGER NOT NULL DEFAULT 0 CHECK (lapses >= 0),
  last_reviewed_at TEXT CHECK (last_reviewed_at IS NULL OR (last_reviewed_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(last_reviewed_at) IS NOT NULL)),
  PRIMARY KEY (library_uuid, lesson_uuid)
) STRICT;
CREATE INDEX review_lesson_due ON review_lesson(due_at);

CREATE TABLE review_question (
  library_uuid      TEXT NOT NULL REFERENCES library(uuid),
  lesson_uuid       TEXT NOT NULL CHECK (lesson_uuid <> ''),
  question_position INTEGER NOT NULL CHECK (question_position >= 0),
  due_at            TEXT NOT NULL CHECK (due_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(due_at) IS NOT NULL),
  interval_days     INTEGER NOT NULL CHECK (interval_days >= 1),
  reps              INTEGER NOT NULL DEFAULT 0 CHECK (reps >= 0),
  lapses            INTEGER NOT NULL DEFAULT 0 CHECK (lapses >= 0),
  last_reviewed_at  TEXT CHECK (last_reviewed_at IS NULL OR (last_reviewed_at GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]*Z' AND julianday(last_reviewed_at) IS NOT NULL)),
  PRIMARY KEY (library_uuid, lesson_uuid, question_position)
) STRICT;
CREATE INDEX review_question_due ON review_question(due_at);

-- App preferences that follow you on this Mac, such as the editor font size.
CREATE TABLE setting (
  key   TEXT PRIMARY KEY CHECK (key <> ''),
  value TEXT NOT NULL
) STRICT;

-- 'GSML' marks the file as a Gosmos learner database; user_version is its schema version.
PRAGMA application_id = 1196641612;
PRAGMA user_version = 2;
