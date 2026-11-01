# Gosmos, the app

A native macOS app (SwiftUI, with AppKit where SwiftUI falls short) for studying Go courses
offline. It reads **libraries**: SQLite files of courses and their lessons
(`../docs/library-format.md` is the format), the ones built into the app and any you open. Any
library works the same, and the app only reads it. The app checks your code with your own Go
toolchain and keeps everything you do in its **learner database** on this Mac; it never runs
git and syncs nothing. How it is built is in `../docs/architecture.md`, how it looks in
`../docs/design.md`.

## Build, install, run

From the repository root, with the Command Line Tools (no Xcode needed) on macOS 27:

```sh
make install
```

That builds `build/Gosmos.app` (release build, Info.plist, icon, ad-hoc signature) and
copies it to `~/Applications`. Open it from Spotlight or Launchpad. `make app` only builds.

The sidebar shows the courses built into the app (the libraries in `../courses/`), then the
libraries you added with File > Open Library… (⌘O) or in Settings > General. Before any of it is
shown, a library passes the format's checks, the same as `go run ./cmd/library check <file>`;
a file that fails is refused with the problems the checks found.

## Test

```sh
make test                      # go vet ./..., go test ./... and swift test
swift test --package-path app  # the app's tests only
```

The tests build their own libraries from `db/schema.sql` and learner databases from
`db/learner.sql` in `$TMPDIR`, with made-up lessons; no test reads a real library or the real
learner database. Some run a real `go build` / `go test` and are skipped when Go is not
installed.

## Where things are

| What | Where |
|---|---|
| Lessons | The libraries, read only: built into the app (`Contents/Resources/Courses`) or files you opened |
| Your study data | `~/Library/Application Support/Gosmos/learner.sqlite` (the libraries you use, progress, attempts, notes, code, projects, server addresses, hints and solutions looked at, the review schedule, editor settings) |
| Lesson folders for building and running | `~/Library/Caches/Gosmos/work/<library uuid>/<lesson uuid>/`, disposable; a redo for review gets its own under `work/review/`, made afresh each time; each run keeps its build files in `work/.runs` while it lasts |
| Settings for this Mac | UserDefaults `io.github.jobloom.Gosmos`: the Go path, window layout, files left in conflict |

Your code lives in the learner database. To run it, the app writes a lesson's files to its work
folder (your files, the read-only tests, `go.mod`, the readme) and builds there; hidden
helper files are added with `go build -overlay`. Edits you make in another editor (Lesson >
Open Lesson Folder in Editor) are taken in when you come back to the app and before every
save and run. If you changed the same file in the app too, it asks which version to keep.

A tests lesson passes when `go test` passes and the tests report no failed case themselves:
some lesson tests print Fail and a summary such as "0 passed, 5 failed" without failing the
test, so the exit status alone is not enough. A code lesson in another language than Go
opens and can be read and edited, but the console says Gosmos runs Go lessons.

## Command-line lessons

Commands run in the course's project folder, which you choose once per course. A lesson's
`${baseURL}` is the server address it suggests, or yours for it: an address you enter
replaces that suggestion for every lesson that makes it, and Use Suggested goes back. A
lesson that checks a server you deployed suggests none and asks for your server's address,
which the course keeps.

## Keyboard

| Keys | Action |
|---|---|
| ⌘R | Run (or Run Checks) |
| ⌘↩ | Submit |
| ⌘. | Stop |
| ⌘O | Open a library |
| ⌘] / ⌘[ | Next / previous lesson |
| ⌥⌘← / ⌥⌘→ | Back / forward through lessons you visited |
| 1 to 5 | Pick a quiz answer |
| ⇧⌘D | Mark as done |
| ⌥⌘H / ⌥⌘S | Show the next hint / the solution |
| ⇧⌘Y | Review |
| ⇧⌘R / ⇧⌘O | Reveal the lesson folder in Finder / open it in your editor |
| ⌥⌘I | Show the inspector (note and lesson facts) |
| ⌘= / ⌘- / ⌘0 | Bigger / smaller / actual-size code |

## Layout

| Path | What it is |
|---|---|
| `Sources/GosmosKit` | Models and services, no SwiftUI or AppKit: the libraries and the learner database (GRDB), the teaching rules, readmes (swift-markdown), the Go lexer, process running, the Go toolchain, code runs, command-line checks |
| `Sources/Gosmos` | The SwiftUI app: models, views, rendering |
| `Tests/GosmosKitTests` | Swift Testing tests |
| `Packages/` | Local packages: the HTTP client for checks, jq, the code editor |
| `Resources/` | `Info.plist` and `AppIcon.icns` |
| `scripts/bundle.sh` | Assembles and signs the .app, with `../db/learner.sql` and the libraries in `../courses/` (`make app` runs it) |
| `scripts/make-icon.sh` | Redraws `Resources/AppIcon.icns` from `scripts/make-icon.swift` |

## Scripted runs

Launch arguments help with screenshots and debugging; they are read from the command line
only and never stored:

```sh
open -n build/Gosmos.app --args -GosmosDataFolder /tmp/gosmos-shots -GosmosOpenLesson <library uuid>/<lesson uuid> -GosmosLaunchAction run
```

`-GosmosDataFolder <folder>` keeps the learner database and the lesson folders there instead of
your own. `-GosmosOpenLesson` takes `<library uuid>/<lesson uuid>`, or a lesson uuid alone. The
action is one of `run`, `submit`, `check`, `answer` (answers a quiz), `footer` (scrolls to the
part below the readme), `hints` (shows every hint that may show), `solution` (shows the
solution), `review` (shows the review), `redo-submit` (opens the lesson as a redo for review
and submits it), `settings` (opens Settings) or `search:<text>` (types a search).
