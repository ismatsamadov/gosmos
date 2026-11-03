# Gosmos, the app

A native macOS app (SwiftUI, with AppKit where SwiftUI falls short) for studying Go courses
offline. It reads **libraries**: SQLite files of courses and their lessons
(`../docs/library-format.md` is the format), the ones built into the app and any you open. Any
library works the same, and the app only reads it. The app checks your code with the Go
toolchain built into it, offline, and keeps everything you do in its **learner database** on
this Mac; it never runs git and syncs nothing. Gosmos never uses the network: every lesson, a
program, its tests or a quiz, is checked on this Mac. How it is built is in
`../docs/architecture.md`, how it looks in `../docs/design.md`.

## Build, install, run

From the repository root, with the Command Line Tools (no Xcode needed) on macOS 27:

```sh
make install
```

That builds `build/Gosmos.app` (release build, Info.plist, icon, the Go toolchain, ad-hoc
signature) and copies it to `~/Applications`. Open it from Spotlight or Launchpad. `make app`
only builds. The Go built in is the one `go env GOROOT` names, without what lessons never use
(`../docs/architecture.md` lists it); `make app GO_ROOT=<folder>` builds in another darwin/arm64
Go. Gosmos then needs nothing installed, and Settings > Go can choose another Go on this Mac.

The sidebar shows the courses built into the app (the libraries in `../courses/`), then the
libraries you added with File > Open Library… (⌘O) or in Settings > General. Before any of it is
shown, a library passes the format's checks, the same as `go run ./cmd/library check <file>`;
a file that fails is refused with the problems the checks found.

Mission Control is the starting screen: continue your next lesson, browse course progress,
and see your explorer rank and achievement badges. Each completed lesson contributes 100 XP,
with a new level every five lessons. Rewards come from the same saved progress as the
sidebar, so earlier completions count and resubmitting never earns duplicate XP. ⌘1 returns
to Mission Control from any lesson or documentation page.

## Test

```sh
make test                      # go vet ./..., go test ./... and swift test
swift test --package-path app  # the app's tests only
```

The tests build their own libraries from `db/schema.sql` and learner databases from
`db/learner.sql` in `$TMPDIR`, with made-up lessons; no test reads a real library or the real
learner database. Some run a real `go build` / `go test` with the Go on PATH, and are skipped
when there is none; `BuiltInGoTests` runs the Go in `build/Gosmos.app`, as on a Mac with no Go
installed, and is skipped until `make app` has built it.

## Where things are

| What | Where |
|---|---|
| Lessons | The libraries, read only: built into the app (`Contents/Resources/Courses`) or files you opened |
| Your study data | `~/Library/Application Support/Gosmos/learner.sqlite` (the libraries you use, progress, attempts, notes, code, hints and solutions looked at, the review schedule, editor settings) |
| Lesson folders for building and running | `~/Library/Caches/Gosmos/work/<library uuid>/<lesson uuid>/`, disposable; a redo for review gets its own under `work/review/`, made afresh each time; each run keeps its build files in `work/.runs` while it lasts |
| Go's caches | `~/Library/Caches/Gosmos/go-build` (builds), `go` (modules, GOPATH) and `telemetry` (set to off), disposable; nothing goes to `~/go` |
| Go | `Contents/Resources/Go`, built into the app |
| Go documentation | `Contents/Resources/Reference/go.sqlite`, built into the app and read only: the Go reference of `../reference/` |
| Settings for this Mac | UserDefaults `app.gosmos.Gosmos`: another Go you chose, window layout, files left in conflict |

Your code lives in the learner database. To run it, the app writes a lesson's files to its work
folder (your files, the read-only tests, `go.mod`, the readme) and builds there; hidden
helper files are added with `go build -overlay`. Edits you make in another editor (Lesson >
Open Lesson Folder in Editor) are taken in when you come back to the app and before every
save and run. If you changed the same file in the app too, it asks which version to keep.

A tests lesson passes when `go test` passes and the tests report no failed case themselves:
some lesson tests print Fail and a summary such as "0 passed, 5 failed" without failing the
test, so the exit status alone is not enough. A code lesson in another language than Go
opens and can be read and edited, but the console says Gosmos runs Go lessons.

## Keyboard

| Keys | Action |
|---|---|
| ⌘R | Run |
| ⌘↩ | Submit |
| ⌘. | Stop |
| ⌘O | Open a library |
| ⌘1 | Mission Control |
| ⌘] / ⌘[ | Next / previous lesson |
| ⌥⌘← / ⌥⌘→ | Back / forward through the lessons and documentation pages you visited |
| 1 to 5 | Pick a quiz answer |
| ⇧⌘D | Mark as done |
| ⇧⌘H / ⌥⌘S | Show the next hint / the solution |
| ⇧⌘Y | Review |
| ⇧⌘R / ⇧⌘O | Reveal the lesson folder in Finder / open it in your editor |
| ⌥⌘I | Show the inspector (note and lesson facts) |
| ⌘= / ⌘- / ⌘0 | Bigger / smaller / actual-size code |
| ⇧⌥F | Format the file with gofmt |
| ⌘/ | Comment or uncomment the selected lines |
| ⌥⌘] / ⌥⌘[ | Indent / outdent the selected lines |
| ⌃⌘? or ⌥-click | Quick Help: the documentation of the name at the insertion point, or clicked |

## Layout

| Path | What it is |
|---|---|
| `Sources/GosmosKit` | Models and services, no SwiftUI or AppKit: the libraries and the learner database (GRDB), the teaching rules, readmes (swift-markdown), the Go lexer, process running, the Go toolchain, code runs, gofmt |
| `Sources/Gosmos` | The SwiftUI app: models, views, rendering |
| `Tests/GosmosKitTests` | Swift Testing tests |
| `Packages/` | The local package of the code editor |
| `Resources/` | `Info.plist` and `AppIcon.icns` |
| `scripts/bundle.sh` | Assembles and signs the .app, with `../db/learner.sql`, the libraries in `../courses/`, the Go documentation in `../reference/` and the Go toolchain (`make app` runs it) |
| `scripts/make-icon.sh` | Redraws `Resources/AppIcon.icns` from `scripts/make-icon.swift` |

## Scripted runs

Launch arguments help with screenshots and debugging; they are read from the command line
only and never stored:

```sh
open -n build/Gosmos.app --args -GosmosDataFolder /tmp/gosmos-shots -GosmosOpenLesson <library uuid>/<lesson uuid> -GosmosLaunchAction run
```

`-GosmosDataFolder <folder>` keeps the learner database, the lesson folders and Go's caches
there instead of your own. `-GosmosAppearance light` or `dark` selects an appearance for that
run, and `-GosmosWindowSize 1100x700` sizes the main window for a layout check. Neither is saved.
`-GosmosOpenLesson` takes `<library uuid>/<lesson uuid>`, or a lesson uuid alone. The
action is one of `run`, `submit`, `answer` (answers a quiz), `footer` (scrolls to the
part below the readme), `hints` (shows every hint that may show), `solution` (shows the
solution), `format` (formats the file in the editor with gofmt), `comment` (comments out the
editor's first line, or uncomments it), `review` (shows the review), `redo-submit` (opens the
lesson as a redo for review and submits it), `settings` (opens Settings), `search:<text>`
(types a search), `doc:<link>` (opens the Go documentation at a link such as
`godoc:fmt.Println` or `godoc-article:spec`, or at `acknowledgements`), `doc-search:<text>`
(types a search in the documentation browser) or `quick-help:<text>` (shows Quick Help for
the first `<text>` in the editor's file).
