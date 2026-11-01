# Gosmos

A quiet, native macOS app for learning Go by doing: courses on the left, the lesson in the
middle, your own Go toolchain checking your code. Short lessons, hints when you are stuck,
and spaced review. No browser, no local web server, no account, no pictures or videos.

The app is not tied to any course. It reads **libraries**: SQLite files holding a syllabus,
in the format described in `docs/library-format.md`. Any file in that format plugs in and
works the same way. The courses that ship with the app are the libraries in `courses/`;
File › Open Library… opens any other.

`courses/apogee.sqlite` holds the Apogee courses: original lessons set at an invented
ground-station startup, planned in `docs/apogee/`. Everything in them is written for Gosmos.

## Two databases

- A **library** is content only: courses, chapters and lessons, starter and solution files,
  hints, quizzes with explanations. The app opens it read-only and never writes to it. Every
  part of a lesson has its own table, and keys, foreign keys and checks keep bad data out: a
  quiz question has one right answer, a test sets only its own kind's fields. The one piece
  of JSON is the body of an HTTP request a lesson sends to your server, kept whole like a
  code file. `db/schema.sql` explains each table.
- The **learner database**, `~/Library/Application Support/Gosmos/learner.sqlite`, holds
  everything that is yours: progress, every Run and Submit, notes, your code, the hints and
  solutions you looked at, and your review schedule. `db/learner.sql` explains each table.

Nothing is kept as JSON or loose files. To build and run your code the app writes a lesson's
files to a scratch folder and reads them back; the database is where they live. The app
never runs git and syncs nothing.

## A new laptop

`gh repo clone jobloom/gosmos`, then `make install` in it, and open Gosmos.

## Parts

| Path | What it is |
|---|---|
| `courses/` | The course libraries built into the app |
| `db/schema.sql` | A library's tables, constraints and conventions (format 7) |
| `db/learner.sql` | The learner database's tables |
| `db/` | The Go package that creates, opens, checks and migrates libraries; `db/history/` keeps the schema of every older format |
| `cmd/library/` | Go tool for any library: `create`, `check`, `migrate`, `verify` |
| `cmd/author/`, `internal/` | Go tool: writes lessons into a library from lesson packages, and exports them back |
| `app/` | The app: a Swift package with GosmosKit and the SwiftUI app |
| `docs/library-format.md` | The library format: what any library must hold, and what the app checks |
| `docs/architecture.md` | How the pieces fit together, and why |
| `docs/design.md` | How the app looks and behaves |
| `docs/apogee/` | The Apogee courses: their world, their Go syllabus, a sample chapter |
| `Makefile` | `make install` builds and installs the app; `make test`, `make check`, `make verify` |

## Requirements

- macOS 27 or later, Xcode Command Line Tools (`xcode-select --install`)
- Go 1.27 or later, for the lessons and the library tools
