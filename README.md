# Gosmos

A quiet, native macOS app for learning Go by doing: courses on the left, the lesson in the
middle, the Go toolchain built into the app checking your code. Short lessons, hints when you
are stuck, and spaced review. No browser, no local web server, no account, no network, no
pictures or videos.

The app is not tied to any course. It reads **libraries**: SQLite files holding a syllabus,
in the format described in `docs/library-format.md`. Any file in that format plugs in and
works the same way. The courses that ship with the app are the libraries in `courses/`;
File › Open Library… opens any other.

`courses/apogee.sqlite` holds the Apogee courses: original lessons set at Apogee, an
invented small space company whose radio dishes talk to satellites, planned in
`docs/apogee/`. They are told as a story (`docs/storytelling.md`): an opening that gives you
your badge, a briefing and a debrief for every chapter, a mission patch for each one you
finish, and a field guide that explains every word of that world in plain language. Everything
in them is written for Gosmos. The path from here to production Go is jobloom/gosmos#18.

## Two databases

- A **library** is content only: courses, chapters and lessons, starter and solution files,
  hints, quizzes with explanations. The app opens it read-only and never writes to it. Every
  part of a lesson has its own table, and keys, foreign keys and checks keep bad data out: a
  quiz question has one right answer, and only a lesson checked by what its program prints
  has expected output. Every lesson is checked on the Mac: a program or its tests with the
  built-in Go, a quiz by its answers. `db/schema.sql` explains each table.
- The **learner database**, `~/Library/Application Support/Gosmos/learner.sqlite`, holds
  everything that is yours: progress, every Run and Submit, notes, your code, the hints and
  solutions you looked at, and your review schedule. `db/learner.sql` explains each table.
  Each user's learner database is created empty; no personal study data ships in the
  repository or app.

Nothing is kept as JSON or loose files. To build and run your code the app writes a lesson's
files to a scratch folder and reads them back; the database is where they live. The app
never runs git and syncs nothing.

## A new laptop

`gh repo clone jobloom/gosmos`, then `make install` in it, and open Gosmos.

## Parts

| Path | What it is |
|---|---|
| `courses/` | The course libraries built into the app |
| `reference/go.sqlite`, `cmd/reference/` | The Go reference built into the app, read offline: the standard library's docs and examples, the spec, the memory model, Effective Go, the FAQ and How to Write Go Code; `make reference` remakes it, `docs/reference-format.md` says how |
| `db/schema.sql` | A library's tables, constraints and conventions (format 9) |
| `db/learner.sql` | The learner database's tables |
| `db/` | The Go package that creates, opens, checks and migrates libraries; `db/history/` keeps the schema of every older format |
| `cmd/library/` | Go tool for any library: `create`, `check`, `migrate`, `verify` |
| `cmd/author/`, `internal/` | Go tool: writes lessons, and a course's story, into a library from packages, and exports them back |
| `app/` | The app: a Swift package with GosmosKit and the SwiftUI app |
| `docs/library-format.md` | The library format: what any library must hold, and what the app checks |
| `docs/architecture.md` | How the pieces fit together, and why |
| `docs/design.md` | How the app looks and behaves |
| `docs/apogee/` | The Apogee courses: their world, their Go syllabus, a sample chapter |
| `docs/authoring.md` | How every new course is made: the steps, where things go, the checks |
| `Makefile` | `make install` builds and installs the app; `make test`, `make check`, `make verify` |

## Requirements

- To use Gosmos: macOS 27 or later, nothing else. The app carries its own Go toolchain and
  works offline; you install nothing, and it never needs the internet.
- To build it: the Xcode Command Line Tools (`xcode-select --install`) and Go 1.27 or later.
  `make app` builds the Go it runs with into the app; `GO_ROOT=<folder>` picks another Go.
