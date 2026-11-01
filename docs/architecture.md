# Architecture

Gosmos is a native macOS app (SwiftUI, with AppKit where SwiftUI falls short) for learning
Go by doing. It reads course **libraries**, SQLite files in the format of
`docs/library-format.md`, and keeps everything the learner does in its own **learner
database**. Nothing in the app depends on a particular library. Go tools create, check and
fill libraries; the courses built into the app are the libraries in `courses/`, committed as
the files themselves. Everything lives in this one repository.

## Repository

| Path | What it is |
|---|---|
| `courses/` | The course libraries built into the app: `apogee.sqlite`, the Apogee courses planned in `docs/apogee/` |
| `db/schema.sql` | A library's schema (format in `user_version`); `db/` also holds the Go package that creates, opens, checks and migrates libraries, with the schema embedded, and `db/history/` the schema of every older format |
| `db/learner.sql` | The learner database's schema |
| `cmd/library/` | The library tool, for any library: `create`, `check`, `migrate`, `verify` |
| `cmd/author/`, `internal/` | The authoring tool: writes lessons into a library from lesson packages and exports them back; `internal/store` writes a lesson's rows, `internal/media` refuses pictures, videos and embeds |
| `app/` | The Swift package: `GosmosKit` (library) and `Gosmos` (the app) |
| `app/Packages/` | The app's own packages: `JQ` (jq compiled in), `LessonHTTP` (the HTTP client for lesson checks) and `CodeEditor` (the editor, plus `SyntaxCore`) |
| `docs/library-format.md` | The library format: what a library holds and what the app checks on open |
| `Makefile` | `make app`, `make install`, `make test`, `make check`, `make verify` |

## The databases

- **Libraries** are content. The `db` package owns their format: the schema, the checks and
  the migrations, used by `cmd/library` and `cmd/author`. A content producer (here
  `cmd/author`) is the only writer. The app opens any library that passes the checks in
  `docs/library-format.md`, read-only; one built into the app is opened immutable, so SQLite
  does not even lock it. The app never creates, changes or migrates a library;
  `cmd/library migrate` does, with the migrations in `db/migrate.go`.
- **The learner database** is the learner's: one file per Mac,
  `~/Library/Application Support/Gosmos/learner.sqlite`, made from `db/learner.sql` (which
  the app carries) on first launch. It lists the libraries in use (the built-in ones and
  each opened file) and keys everything else by `(library_uuid, lesson_uuid)`, so a library
  can be replaced by a newer version without touching anyone's progress.
- Every connection turns on foreign keys and keeps the rollback journal (no WAL), so each
  file is always whole. Timestamps are RFC 3339 in UTC with a trailing Z.
- Settings that belong to one Mac, like window frames, live in UserDefaults.

## The app

Swift 6.4 (swift-tools-version 6.4) in the Swift 6 language mode with complete concurrency
checking; macOS 27 or later, the current release, with Liquid Glass (see `docs/design.md`
for how it looks). It builds with
SwiftPM and the Command Line Tools, no Xcode needed. `make app` assembles
`build/Gosmos.app` (Info.plist, icon, ad-hoc signature); the package also opens in Xcode
as is.

**GosmosKit** imports neither SwiftUI nor AppKit. It holds the models and services and is
covered by Swift Testing tests.

- **AppDatabase**: a library, through one read-only GRDB `DatabaseQueue` (a
  `DatabasePool` would need WAL). It validates the file on open and serves typed reads: the
  catalog, full-text search, and a lesson's content as one enum case per lesson type built
  from its tables.
- **LearnerDatabase**: the learner database, through one GRDB `DatabaseQueue`. It does every
  write: progress, attempts, notes, code, hints and solutions looked at, review. The UI
  follows your progress through `ValueObservation`. A write nothing else waits for, such as
  a review answer, goes through `PendingWrites`, and quitting waits for those as it waits
  for a lesson's saves. A new database is made aside and moved in whole; a draft that a
  crash left beside it goes at a later launch.
- **Teaching**: the rules, as pure values and functions.
  - Hints show one at a time and stay shown; the last waits for a failed Submit.
  - The reference solution is there whenever you want it. Looking before you pass is
    recorded, and the lesson then comes back for review, as it does after a pass that took
    three failed Submits or more.
  - Spaced review brings a lesson or a quiz question back a day after it is learned, then
    3, 7, 21 and 60 days after each review it passes; a failed review starts it over at a
    day. Due quiz questions are interleaved across lessons.
- **Readme**: Apple's swift-markdown parses the Markdown. A visitor turns the tree into the
  render model, including `[!NOTE]`-style callouts and OS-tagged headings. Images and
  videos are never rendered.
- **Go syntax**: a lexer that follows the Go specification and turns code into highlight
  ranges. Readme code blocks and the editor both use it.
- **Process running**:
  - Programs start from a resolved path with an argument array, never a shell string,
    except for a lesson's own command line, which runs under `/bin/zsh -c`.
  - Each run gets an explicit environment and working folder.
  - stdout and stderr are read concurrently.
  - Each run starts in its own process group. A timeout or a cancelled task kills the
    whole group.
- **Go toolchain**: `go` comes from the path set in Settings, else from the login shell's
  PATH, else from the usual install locations.
- **Code runs**: your saved code is written to
  `~/Library/Caches/Gosmos/work/<library>/<lesson>/`. `go run` or `go test` runs there, and
  the output is compared or the tests are read. A review redo works from the starter code in
  `work/review/<library>/<lesson>/` and is never saved.
- **Command-line checks**:
  - Steps run in order in the course's project folder, with captured `${names}` passed on
    to later steps.
  - HTTP requests go through the app's own HTTP/1.1 client on SwiftNIO over
    Network.framework (NIOTransportServices). Each request goes out exactly as the step
    defines it, with no added cookies, cache or headers. Redirects are followed only when
    the step allows it, and response trailers are read, which URLSession cannot do.
  - jq programs, both tests and JSON-path captures, run in process on jq 1.8.2. It is
    compiled into the app from its C sources, with Oniguruma for regular expressions,
    because the lessons use the real language (`scan`, `any`, `all`, slices).

**Gosmos** is the SwiftUI app.

- `@main` App with a single main window. It is a `NavigationSplitView`: the sidebar holds
  courses, chapters and lessons with done marks and search, and the detail shows the lesson
  in a view made for its type.
- A Settings scene, and menu commands: Run ⌘R, Submit ⌘↩, Next ⌘], Previous ⌘[.
- State lives in `@Observable @MainActor` models that get their services through the
  environment. Views stay thin, and there are no singletons.
- The code editor is an `NSTextView` (TextKit 2) behind `NSViewRepresentable`. It is
  monospaced, keeps tabs, and has line numbers, undo, find and highlighting. It saves your
  code to the learner database after a second of idle typing.
- Errors are typed and shown to you, never swallowed. Logging uses `os.Logger`.

## Self-contained, fast, reliable

- **Self-contained**: the app is one process, and its bundle holds everything it runs: the
  database layer, Markdown, highlighting, jq and the HTTP client. Outside it, the app uses
  only what a Mac for learning Go has anyway:
  - the Go toolchain, to run your code (the app finds it, or tells you how to install it);
  - `/bin/zsh`, for a lesson's own command lines.

  It has no account, sends nothing anywhere and syncs nothing; the network is used only by
  a lesson that sends requests to your own server.
- **Fast**: these budgets are measured with `os_signpost` intervals and tested where they
  can be.
  - The window is usable less than a second after launch.
  - Switching lessons takes under 100 ms.
  - Search results arrive in under 50 ms while you type.
  - No work on the main thread takes longer than a frame.
- **Reliable**:
  - No force unwraps, `try!` or `fatalError` except for programmer errors.
  - Every failure is a typed error that you see.
  - Database writes happen in transactions, and `PRAGMA quick_check` runs on open.
  - Every outside process can be cancelled and has a timeout.
  - The window and your place are restored on relaunch.

## Dependencies

Everything is kept at its latest release and pinned by `app/Package.resolved` and `go.sum`,
both committed. These were the latest on 2026-10-31:

| Part | Version |
|---|---|
| Toolchain | Command Line Tools 27.2 (Swift 6.4.2) and Go 1.27.3 |
| App | GRDB.swift 7.11.1, swift-markdown 0.9.0, SwiftNIO 2.104.0, swift-nio-transport-services 1.28.0 |
| App, compiled in | jq 1.8.2 and Oniguruma 6.9.10, vendored as C targets by `app/Packages/JQ/scripts/vendor.sh` from the release tarballs, with SHA-256 checked; one line is changed so jq never reads `~/.jq` |
| Library tools | modernc.org/sqlite v1.61.1 (with modernc.org/libc held at the exact v1.79.0 its go.mod requires), goldmark v1.8.7 |

A new dependency needs a reason written here.

## Tests

`make test` runs `go vet ./...`, `go test ./...` and `swift test`. Tests build their own
databases in `$TMPDIR`, from `db/schema.sql` with fixture lessons of every type and from
`db/learner.sql`. No test touches `courses/` or a real learner database.
