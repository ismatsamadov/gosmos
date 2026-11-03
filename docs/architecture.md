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
| `app/Packages/` | The app's own package: `CodeEditor` (the editor, plus `SyntaxCore`) |
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
  can be replaced by a newer version without touching anyone's progress. One copy of the app
  uses it at a time: a launch while an earlier copy with the same bundle id runs, such as
  `~/Applications/Gosmos.app` while `build/Gosmos.app` opens, brings that one forward and
  quits before it opens anything (`SingleInstance`). A QA copy under another bundle id runs
  beside it.
- **The Go reference**, `reference/go.sqlite` (`docs/reference-format.md`), is the Go
  documentation built into the app. The app opens it read-only and immutable from
  `Contents/Resources/Reference`, after checking its `application_id` and `user_version`.
- Every connection turns on foreign keys and keeps the rollback journal (no WAL), so each
  file is always whole. Timestamps are RFC 3339 in UTC with a trailing Z.
- Settings that belong to one Mac, like window frames, live in UserDefaults.

## The app

Swift 6.4 (swift-tools-version 6.4) in the Swift 6 language mode with complete concurrency
checking; macOS 27 or later, the current release, with Liquid Glass (see `docs/design.md`
for how it looks). It builds with
SwiftPM and the Command Line Tools, no Xcode needed. `make app` assembles
`build/Gosmos.app` (Info.plist, icon, the Go toolchain, ad-hoc signature); the package also
opens in Xcode as is.

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
  ranges. Readme code blocks, the documentation's declarations and examples, and the editor
  all use it.
- **Go documentation** (`Reference/`):
  - `ReferenceDatabase` reads the Go reference through one read-only `DatabaseQueue`: the
    list (articles, and packages grouped by their first path element), a page with its
    Markdown parsed off the main actor (a package laid out as go doc orders it, a const or
    var block once for all its names; an article with its sections), Quick Help's summary,
    the credits, and search. Search puts the names the text is first (packages by path or
    name, then names, funcs and types before methods, shorter package paths first), then
    names it starts, then FTS5 matches ranked by bm25 with names weighing ten times the text.
  - `DocLink` is a place in it: `godoc:<path>`, `godoc:<path>.<Name>` or
    `godoc:<path>.<Type>.<Name>`, `godoc-article:<slug>#<anchor>`. `DocNames` holds every
    package and name in memory, read once at launch, to turn code into links: inline code in
    a readme (`DocLinker`), and in the editor a selector resolved through the file's imports
    (`GoSource` reads the imports, renamed and dot ones too, and the identifier at a place,
    never in a comment or a string).
  - `LinkTarget` is the one rule for links: into the documentation, or to a lesson of the
    readme's library. Anything else is no link at all. A link carries the app's own address
    (`LinkTarget.url`): `godoc:`, `godoc-article:`, or `gosmos-lesson:<uuid>` for a lesson,
    never the web address a readme gave for one.
- **Process running**:
  - Programs start from a resolved path with an argument array, never a shell string.
  - Each run gets an explicit environment and working folder.
  - stdout and stderr are read concurrently, and stdin, when a run has input, is written
    between the reads.
  - Each run starts in its own process group. A timeout or a cancelled task kills the
    whole group.
- **Go toolchain**: the Go built into the app, unless Settings chose another; `gofmt` is the
  one in that Go's GOROOT. Every Go runs offline and isolated (see The built-in Go below).
- **Formatting**: Format sends the file in the editor to `gofmt` on stdin and reads it back
  from stdout, with a 10-second timeout, so nothing is written to the lesson folder. A syntax
  error comes back as gofmt's messages, naming the file.
- **Code runs**: your saved code is written to
  `~/Library/Caches/Gosmos/work/<library>/<lesson>/`. `go run` or `go test` runs there, and
  the output is compared or the tests are read. A review redo works from the starter code in
  `work/review/<library>/<lesson>/` and is never saved.

**Gosmos** is the SwiftUI app.

- `@main` App with a single main window. It is a `NavigationSplitView`: the sidebar holds
  Mission Control, Review, Go Documentation, then courses, chapters and lessons with done
  marks and search, and the detail shows Mission Control (where the app opens), the lesson
  in a view made for its type, the review, or the documentation browser. XP, rank and badges
  (`LearningProgress`) are worked out from the finished lessons, never stored.
- A Settings scene, and menu commands: Run ⌘R, Submit ⌘↩, Next ⌘], Previous ⌘[, Mission
  Control ⌘1, and the Editor menu: Format Code ⇧⌥F, Toggle Comment ⌘/, Indent ⌥⌘], Outdent
  ⌥⌘[ and Show Documentation ⌃⌘?.
- The documentation browser (`DocsModel`, `DocsBrowser`): a list and its own search beside
  the page, which is a lazy stack of rows (header, doc blocks, index lines, declarations,
  sections), so only what is in view is made; a place scrolls to its declaration or section.
  The last eight pages stay parsed. Back and Forward keep one history of places: Mission
  Control, lessons and documentation pages, and leaving a lesson for the documentation saves
  it as switching lessons does. Every link in the window goes through `AppModel.follow`, which
  opens the documentation or a lesson and does nothing with any other address.
- Quick Help: the editor reports an Option-click's offset, the app resolves the identifier
  there, and an `NSPopover` beside it shows the summary.
- State lives in `@Observable @MainActor` models that get their services through the
  environment. Views stay thin, and there are no singletons.
- The code editor is an `NSTextView` (TextKit 2) behind `NSViewRepresentable`. It is
  monospaced, keeps tabs, and has line numbers, undo, find and highlighting. It saves your
  code to the learner database after a second of idle typing. The Editor menu and the
  buttons above the editor change the code through the package's `CodeEditorProxy`, each
  change one undo step.
- Errors are typed and shown to you, never swallowed. Logging uses `os.Logger`.

## Self-contained, fast, reliable

- **Self-contained**: the app is one process, and its bundle holds everything it runs: the
  Go toolchain that builds, runs, tests and formats your code, the Go documentation, the
  database layer, Markdown and highlighting. Nothing needs installing. Gosmos never uses the network: it has no
  account, sends nothing anywhere, syncs nothing, fetches nothing and opens no web page, and
  Go never reaches the network either. Every lesson is checked on the Mac alone: code by
  running it, a quiz by its answers. `IsolationTests` keeps it so: it fails on any networking
  module or API (URLSession, Network.framework, SwiftNIO, sockets), any web address written
  into the app's code, and any link handed to the system to open.
- **Fast**: these budgets are measured with `os_signpost` intervals and tested where they
  can be.
  - The window is usable less than a second after launch.
  - Switching lessons takes under 100 ms.
  - Search results arrive in under 50 ms while you type.
  - A documentation page opens in under 100 ms, the biggest packages too, and each keystroke
    of a documentation search takes under 30 ms, on the full Go reference.
  - No work on the main thread takes longer than a frame.
- **Reliable**:
  - No force unwraps, `try!` or `fatalError` except for programmer errors.
  - Every failure is a typed error that you see.
  - Database writes happen in transactions, and `PRAGMA quick_check` runs on open.
  - Every outside process can be cancelled and has a timeout.
  - The window is restored on relaunch, and it opens to Mission Control, which offers the
    lesson you opened last while it is unfinished.

## The built-in Go

Gosmos carries its own Go toolchain, so a learner installs nothing and never needs the
internet. macOS 27 runs only on Apple silicon, so one darwin/arm64 toolchain is enough.

- **In the bundle**: `app/scripts/bundle.sh` copies the Go in `GO_ROOT` (by default the one
  `go env GOROOT` names: Go 1.27.3 from Homebrew, 270 MB) to `Contents/Resources/Go`, with its
  LICENSE, and leaves out what building, running, testing and formatting standard-library
  code never uses. That makes it 146 MB, and the app 171 MB (a 55 MB zip).

  | Left out | Size | Why it can go |
  |---|---|---|
  | `src/cmd` | 64 MB | The sources of the go command, the compiler and the other tools. Their programs are in `bin` and `pkg/tool`; only `go tool` for a tool that isn't there would build one from them |
  | `test/` | 19 MB | Go's own tests of its compiler and runtime |
  | `src/**/testdata` | 18 MB | Test inputs of the standard library's own tests |
  | `pkg/tool/darwin_arm64/fix` | 8.7 MB | `go fix` |
  | `api/` | 8.6 MB | The API lists that check Go's own compatibility |
  | `pkg/tool/darwin_arm64/cover` | 5.3 MB | `go test -cover`, which the app never asks for |
  | `doc/`, `misc/` | 0.5 MB | Go's release notes and spec as HTML, and editor and platform extras |
  | `lib/wasm` | 32 KB | WebAssembly support (`wasm_exec.js`) |
  | `lib/time/mkzip.go` | 2 KB | The program that makes `zoneinfo.zip`, which stays |

  Everything else stays: `bin/go` and `bin/gofmt`; `compile`, `link`, `asm`, `cgo`, `vet`
  (which `go test` runs) and `preprofile` (for a `default.pgo` profile); the standard
  library's sources, which `go build` compiles into the build cache; `pkg/include`;
  `lib/time/zoneinfo.zip` and `lib/fips140`; `go.env` and `VERSION`. Go 1.27 ships no
  `buildid` or `pack` program any more: the go command does their work itself.
- **Out of this repository's module**: `build/` is inside the gosmos Go module, so the app is
  put together under `$TMPDIR` and moved into `build/` whole, the old one out first; `go vet
  ./...` or `go test ./...` running meanwhile never meets a half-copied Go. Every `.go` file
  left is in Go's own module, `src/go.mod`, which `./...` skips: `bundle.sh` stops if one is
  anywhere else, and `BuiltInGoTests` checks the built app.
- **Checking it**: with a new Go, run every lesson with the bundled one, from an empty build
  cache and without the network:

  ```sh
  PATH="$PWD/build/Gosmos.app/Contents/Resources/Go/bin:/usr/bin:/bin" GOCACHE="$(mktemp -d)" GOPROXY=off go run ./cmd/library verify --teaching courses/apogee.sqlite
  ```

  `BuiltInGoTests` runs a program, tests and Format through `GoRunner` with the built app's
  Go, PATH set to `/usr/bin:/bin` and HOME an empty folder, as on a Mac with no Go installed.
- **Signing**: each program in the toolchain (`bin/go`, `bin/gofmt` and the tools in
  `pkg/tool`) is signed with the app's identity before the app itself; today that is ad hoc. A
  Developer ID build for notarization must sign each of them, and then the app, with the
  Developer ID identity, the hardened runtime (`--options runtime`) and a secure timestamp
  (`--timestamp`), then notarize and staple the app. The race detector's `.syso` files in
  `src/runtime/race` are object files, not programs, and are not signed.
- **Which Go runs**: the built-in Go, unless Settings › Go chose another with Choose…, an
  option for advanced use; Use Built-in Go goes back. Nothing else is looked for: not PATH,
  the login shell or the usual install locations. A Go that doesn't run is reported, never
  swapped for another. The built-in one is missing, or damaged when its go command doesn't
  run, it would use another GOROOT, or its compiler, linker, assembler or vet is gone: then
  the app says to reinstall Gosmos. A chosen one that doesn't run is named.
- **Its environment**: every Go the app runs, built in or chosen, gets the app's own
  environment with these changes (`Toolchain.goEnvironment`):

  | Variable | Value | Why |
  |---|---|---|
  | `GOTOOLCHAIN` | `local` | Never another Go, downloaded or not, whatever a go.mod asks for |
  | `GOPROXY`, `GOSUMDB` | `off` | Never a module or a checksum from the network: lessons use the standard library |
  | `GOCACHE` | `~/Library/Caches/Gosmos/go-build` | The build cache, beside the lesson folders, not `~/Library/Caches/go-build` |
  | `GOPATH` | `~/Library/Caches/Gosmos/go` | The module cache, so nothing is written to `~/go` |
  | `TEST_TELEMETRY_DIR` | `~/Library/Caches/Gosmos/telemetry` | Where the go command and its tools keep telemetry; its `mode` file says `off`, so nothing is counted or uploaded, whatever `go telemetry` says for your own Go. Go has no other way to move it: `GOTELEMETRY` and `GOTELEMETRYDIR` can't be set |
  | `GOENV` | `off` | The settings `go env -w` keeps for your own Go don't apply |
  | `CGO_ENABLED` | `0` | Lessons are pure Go, and a Mac without the Command Line Tools has no C compiler |
  | `GOWORK`, `GO111MODULE` | `off`, `on` | A lesson folder is one module |
  | `GOTMPDIR` | the run's folder in `work/.runs` (`GoRunner`) | The go command's work folder, a test's binary included, not in `$TMPDIR`: it goes with the run, even when Stop kills go before it cleans up |
  | Any other `GO…` or `CGO_…` | removed | GOROOT, GOFLAGS, GOOS and the rest would change what is built |

  `GOFLAGS=-mod=mod` isn't set: lessons need no module beyond the standard library, so go.mod
  never needs changing, and an import from outside it fails at once with a clear message. A
  scripted run's `-GosmosDataFolder` holds these caches too.

## Dependencies

Everything is kept at its latest release and pinned by `app/Package.resolved` and `go.sum`,
both committed. These were the latest on 2026-10-31:

| Part | Version |
|---|---|
| Toolchain | Command Line Tools 27.2 (Swift 6.4.2) and Go 1.27.3 |
| App, built in | Go 1.27.3 darwin/arm64, the Go that builds the app (`GO_ROOT`), trimmed |
| App | GRDB.swift 7.11.1, swift-markdown 0.9.0 |
| Library tools | modernc.org/sqlite v1.61.1 (with modernc.org/libc held at the exact v1.79.0 its go.mod requires), goldmark v1.8.7 |

A new dependency needs a reason written here.

## Tests

`make test` runs `go vet ./...`, `go test ./...` and `swift test`. Tests build their own
databases in `$TMPDIR`, from `db/schema.sql` with fixture lessons of every type, from
`db/learner.sql`, and from `db/reference.sql` with a few packages, names and articles. No test
touches `courses/` or a real learner database; the documentation's time budgets read
`reference/go.sqlite`, the full reference, which is only read.
