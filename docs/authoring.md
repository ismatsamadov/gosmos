# Writing and fixing lessons

Everything needed to add a lesson, a chapter or a course, to work out a number the story
needs, or to fix what is there, is in this repository: the tools in `cmd/` and `scripts/`, the
plans and the world in `docs/`. Every lesson and story is text in the series' library folder,
`courses/apogee/` for Apogee, the only stored copy of each; `make library` builds the library
file the app carries from it, and the docs hold every decision. The rules the steps rely on
live in the story bible (`docs/apogee/story/style.md`, Style guide and Naming rules) and the
library format (`docs/library-format.md`); where they say more, they win.

## What every course is

- **Original.** Written for Gosmos from nothing, in an invented world. Writers never read or
  mirror another course's lessons.
- **Hands-on.** About five lessons in six are code: `code`, graded by what the program
  prints, or `code_tests`, graded by `go test`. The rest are short quizzes. There is nothing
  else: library format 9 has only these three lesson types.
- **Offline.** Standard library only, for the built-in Go (1.27.3). No lesson needs the
  network, a server, GitHub or anything outside the app. No pictures, videos, embeds or links
  that leave the app.
- **Built in.** A course ships inside a library built from `courses/`, and is checked like
  code before it merges.

## The tools

All run from the repository's folder.

| Command | What it does |
|---|---|
| `go run ./cmd/author help` | Prints the formats: a lesson package, a course's story package (opening, briefings and debriefs, field guide), and a library folder |
| `make library` | Builds `courses/apogee.sqlite` from `courses/apogee/` (`go run ./cmd/author build`): the same text makes the same file, and every uuid comes from the slugs, so learners keep their progress |
| `go run ./cmd/author export <library> [--course <c> [--chapter <ch> \| --story]] <folder>` | Writes a library file out as text: the whole library as a library folder, or a course's lessons, one chapter's, or its story |
| `go run ./cmd/author course\|chapter\|lesson\|story <library> …` | Adds or replaces a course, a chapter, a lesson or a course's story in a library file, in place; for a scratch library, never the repository's, which is built |
| `scripts/pkgcheck.sh <folder>` | Checks lesson packages (below); `QUICK=1` in front for a fast first pass |
| `go run ./cmd/library check <library>` | The checks the app makes before it opens a library |
| `go run ./cmd/library verify [--course <c>] --teaching [--slow] <library>` | Runs every code lesson: it must pass with its solution and fail with its starter; also wants the hints, solution notes and quiz explanations. `--slow`, for a slower machine, reports only runs stopped at 10 s |
| `go run ./cmd/library drift [--course <c>] [--accept <list>] [--diff] <library>` | Lists given code (read-only and hidden files) that changes between lessons without a lesson changing it, but for the changes the list accepts (`courses/apogee-drift.txt`) |
| `go run ./cmd/passes --from <day> [--to <day>] [--sat <ids>] [--station <codes>]` | The passes of the Apogee world: pass ID, AOS, LOS, maximum elevation |
| `app/scripts/snapshot.sh <out.png> <light\|dark> [launch arguments]` | Draws one screen of a copy of the app to a PNG |
| `make test`, `make check`, `make verify`, `make install` | The whole repository's tests, the libraries' checks and drift, every lesson run, the app installed; `check`, `verify` and `install` build the library first |

## Where things go

| What | Where |
|---|---|
| The world, the cast, the code the learner grows, the writing rules | the series' story bible: `docs/apogee/story.md` for Apogee, `docs/<series>/story.md` for a new world |
| The plan of every chapter and lesson | `docs/<series>/<course slug>.md`, in the format of `docs/apogee/go-foundations.md`, ending with "Changes in writing and review" |
| The first chapter written out in full | `docs/<series>/sample-chapter-1.md`, for a new series; a new Apogee course adds its own calibration chapter beside it |
| The lessons and each course's story | the series' library folder, as `go run ./cmd/author help` describes it: `courses/apogee/<course>/` holds `course.md`, `story/` and a folder for each chapter, with its `chapter.md` and a lesson package for each lesson; a new series gets `courses/<series>/` |
| The library file the app carries | `courses/apogee.sqlite`, built by `make library` and never committed |
| Notes, screenshots and other scratch work | a scratch folder, `$TMPDIR/<course slug>/` (`$W` below), never committed |

Work on a branch, in the packages themselves. After a change, `make library` builds the file
the checks, `verify` and the app read. Never edit a library file by hand, or import into the
built one: the next build replaces it, and only the text is kept.

## Adding a lesson to a course

1. **Plan it.** Add its entry to the course's syllabus: objective, scenario, exercise, check,
   difficulty, slug and type, in its place in the chapter. Anything new in the world (a
   person, a station, a number) goes into the bible first.
2. **Read its neighbors** in its chapter's folder, `courses/apogee/<c>/<ch>/`. A lesson whose
   code continues another's starts from that lesson's solution; in a guided project every
   lesson does.
3. **Write the package** in `courses/apogee/<c>/<ch>/<lesson>/`, in the format `author help`
   prints: front matter and readme, `starter/`, `solution/`, `hints.md`, `solution_note.md`,
   and `expected_output.txt` or `quiz.md`. Its front matter ends with `created:` and
   `modified:`, both the time you write it, in UTC. The style guide is the law. A pass in it
   comes from `cmd/passes`, any other number from the bible or from running the code (see
   Calculating the world).
4. **Make room.** Positions count from 0 and each is used once in a chapter. To put a lesson
   between two others, raise the `position:` of every lesson after it, and set their
   `modified:`.
5. **Check it:** `scripts/pkgcheck.sh courses/apogee/<c>/<ch>/<lesson>`, until it prints `ok`.
6. **Build and check the library:**

   ```sh
   make library
   go run ./cmd/library check courses/apogee.sqlite
   go run ./cmd/library verify --course <c> --teaching courses/apogee.sqlite
   go run ./cmd/library drift --course <c> --accept courses/apogee-drift.txt courses/apogee.sqlite
   ```

7. **Tell the story.** A world word the readme uses for the first time needs a card in the
   field guide; a lesson that now ends its chapter may need the debrief changed. Edit
   `terms.md` or `chapters.md` in `courses/apogee/<c>/story/`, then build again. See
   `docs/storytelling.md`.
8. **See it in the app** (below).
9. **Ship it.** Note in the syllabus what changed from its plan, then `make test`, `make
   check` and `make verify`. The syllabus, any bible change and the packages go into one
   branch and one pull request, which shows every change as text.

## Adding a chapter

A folder `courses/apogee/<c>/<ch>/`, named for the chapter's slug, with a `chapter.md`: its
slug, title and position as front matter, then the syllabus's description (`author help`
shows one). Then each lesson as above, and a briefing and a debrief for the chapter in the
story's `chapters.md` (at most 70 and 50 words; `docs/storytelling.md` says what they hold).
To put a chapter between two others, raise the `position:` of every chapter after it.

## Adding a course

Every Gosmos course is made the way Go Foundations was:

1. **Plan.** Add to the story bible first: any new people, systems, data and names, each
   name checked against the web as its Naming rules say. Then write the syllabus: the
   course's title, slug, kind, short description and description; the at-a-glance table; and
   an entry for every lesson. Keep Go Foundations' rhythm: chapters of 7 to 10 lessons, a quiz
   right after the idea it checks, a slightly larger lesson closing most chapters, difficulty
   rising gently from 1 to 10.
2. **Calibrate.** Write the first chapter in full, build it and read it in the app before
   anything else is written. What it teaches goes into the bible or the syllabus, so every
   later lesson starts right. A guided project builds its whole program and any hidden test
   harness first: each lesson's solution is then a cut of that program, with that lesson's
   features and no later ones (the capstone's section of the bible shows how).
3. **Write.** Chapters can be written in parallel, one writer per chapter (see Working in
   parallel). Each lesson is written as in Adding a lesson.
4. **Review.** Someone who didn't write a chapter checks every lesson of it against the
   checklist below and fixes the packages.
5. **Tell the story** as `docs/storytelling.md` says, in the course's `story/`: `opening.md`,
   the pages before the first lesson; `chapters.md`, a briefing and a debrief for each
   chapter; and `terms.md`, the field guide. The course itself is `course.md` in its folder,
   `courses/apogee/<c>/`: slug, title, kind, position and short description as front matter,
   then the description.
6. **Build** with `make library`, and run `check`, `verify --teaching` and `drift`. Every uuid
   comes from the slugs, so a lesson changed and built again keeps its uuid, and learners
   keep their progress; never rename a slug that has shipped.
7. **See it** in the app, every chapter.
8. **Ship.** As for a lesson. Then write the syllabus's "Authoring status" and "Changes in
   writing and review", and bring the bible up to date with what writing and review settled.

## Fixing a lesson

1. Edit its package in `courses/apogee/<c>/<ch>/<lesson>/`, set its `modified:` to now, run
   `pkgcheck`, then `make library`, and `check`, `verify --teaching` and `drift` for the
   course. The lesson keeps its uuid, so learners keep their progress.
2. **Code that carries forward.** When later lessons give this lesson's code as read-only or
   hidden files, or start from its solution (every lesson of a guided project does), the same
   fix goes into each of them: their starters and their solutions. `library drift --course
   <c> --diff` lists every place the given code no longer matches what the learner last saw.
3. **Shared hidden files.** A harness hidden in many lessons, such as the capstone's
   `sim_*_test.go` files, stays byte for byte the same in every lesson that has it: change
   every copy, then verify the whole course.
4. **The story.** Edit the course's `story/` files, then build.
5. **Say why.** A fix that changes what a lesson teaches or tests goes into the syllabus's
   "Changes in writing and review".

## Given code changes only for a reason

The read-only and hidden files are the team's codebase as the learner reads it. A
declaration that is reworded, renamed or rewritten between lessons, with no lesson or story
beat behind it, reads as a mistake, and a different error text for the same failure is a real
inconsistency.

- A declaration in a given file keeps the text it had when the learner last saw it: its
  earlier given copy, or the solution of the last lesson that changed it. Byte for byte, doc
  comments included.
- A lesson may give only the declarations it needs; a file holding part of an earlier file is
  fine if every declaration in it is unchanged.
- New behavior goes in new declarations. When a lesson's task changes an old declaration, the
  new text appears first in that lesson's solution and stays the same afterwards.
- A struct may gain a field, or a route list a route, when a chapter's feature needs it and a
  readme or the story introduces it.
- A fix the story names (a debrief, a briefing or a readme saying who fixed what) may appear
  in the next given copy.

`library drift` reports each change: `its code changed` means behavior or names, `only its
comments or layout changed` the wording; `back to an older text` that an earlier version
returned. A change a lesson or the story asks for is fine: add it to `courses/apogee-drift.txt`
as drift names it (`<course>/<chapter>/<lesson>: <file> gives <declaration>`), with an
indented line or more saying where the readme, the story or the syllabus asks for it, and
drift stops reporting it. An entry that no longer names a change is reported in its place, so
the list goes when its changes go. A slip is fixed in the lessons instead. `make check` runs
drift with the list over the whole library.

## Calculating the world

- **Passes.** Every pass in a lesson or the story (pass ID, AOS, LOS, maximum elevation) comes
  from the pass model; none is invented, and none is longer than its orbit allows:

  ```sh
  go run ./cmd/passes --from 2029-07-23 --sat pipit-2 --station YST
  pipit-2/YST/22814 2029-07-23T00:41:10Z 2029-07-23T00:46:33Z 9.0
  ```

  The model and the elements of all eighty satellites are written out in the bible under
  "The pass model", and `cmd/passes` holds the same code, copied from the lessons. The
  elements describe each satellite from 1 January 2029, or from its deployment; the six fitted
  to passes the stations recorded (Pipit-1, Pipit-2, Furrow-A, Furrow-B, Brine-3, Cirrine-1)
  also run back to 2027, and there they give the passes Seasons 3 to 6 used, to the second.
  Two things differ from those seasons' lists: a pass shorter than two minutes no longer
  counts, and Season 3's catalog keeps two Furrow passes a few seconds shorter than the model
  (its lessons print the catalog's). Season 1's passes came before the model. A new satellite
  gets its elements in the bible and in `cmd/passes/network.go` first, with a test.
- **The world's facts.** The bible holds the stations (place, mask, dishes), the fleet
  (orbits, radios), the people, the timeline of every season, and useful numbers such as
  orbital periods and slant ranges. Take them from there; add a new one there first.
- **Booking IDs only grow.** Each season's range is in its section of the bible; the
  highest used is `bk-9951` (the capstone's test bookings). Four digits leave `bk-9952` to
  `bk-9999`, so a story after the capstone needs the bible to settle a longer form first.
- **What a program produces.** Run times, allocation counts, a week's totals and anything else
  a program measures are found by running the lesson's own code with the built-in Go, then
  written into the bible; never guessed. The capstone's whole program is its last lesson's
  solution (`author export … --chapter the-storm`).

## Checks

- **`scripts/pkgcheck.sh`**, on packages before a build: the front matter; a readme of 60 to
  200 words (a quiz's 60 to 120) with `## Your task` and no H1, picture or embed, its code
  lines at most 50 characters (the app's readme column wraps longer ones); 2 or 3 hints and a
  solution note of 1 to 3 sentences; gofmt; the solution passing 30 runs in a row (a tests
  lesson's tests also 3 times under `-race`); the starter building and failing for the
  lesson's reason, never by a panic.
- **`library check`**: what the app checks when it opens a library.
- **`library verify --teaching`**: every code lesson passes with its solution and fails with
  its starter, as the app runs it; a run over 2 s (5 s in a guided project) is reported slow
  and one over 10 s stopped. Run it on a quiet Mac: under heavy load runs slow down.
- **`library drift`**: given code changes only for a reason (above).
- **The review checklist**, for someone who didn't write the lesson:
  1. It teaches what its syllabus entry says, one idea, and the Go is right and idiomatic
     (check APIs with `go doc`); every claim in the readme, hints, note and quiz holds.
  2. The task is fair: a newcomer who has done the earlier lessons can do it from the readme
     and hints alone, and nothing is tested that the readme doesn't ask for.
  3. The tests check the normal case, a boundary and the mistake learners make; they accept
     every correct solution, not only the official one; messages read `call = got, want
     want`; nothing depends on the wall clock, randomness, map order, the time zone, the
     network or the home folder.
  4. The starter fails for the lesson's reason and the place to change is obvious; the
     solution changes only what the task asks.
  5. The words keep the style guide; world terms are explained in plain words first and are
     in the field guide; every number is right.
  6. The story keeps the bible's time, places and voices, and no lesson contradicts an
     earlier one.
  7. Shared types, formats, file names and output texts match the bible exactly.

## Seeing it in the app

Tests can't see how a lesson reads. Build the app once (`make app`), then draw any screen of a
copy of it with the library as you last built it:

```sh
L=$(sqlite3 courses/apogee.sqlite "SELECT l.uuid FROM lesson l JOIN chapter ch ON ch.uuid = l.chapter_uuid
  WHERE ch.slug = '<ch>' AND l.slug = '<lesson>'")
QA_LIBRARY=courses/apogee.sqlite app/scripts/snapshot.sh $W/shots/readme.png dark -GosmosOpenLesson $L
QA_DELAY=8 app/scripts/snapshot.sh $W/shots/submit.png dark -GosmosOpenLesson $L -GosmosLaunchAction submit
app/scripts/snapshot.sh --clean
```

The copy runs under its own bundle identifier with its own learner data, so your app, its
settings and your progress stay untouched; `app/scripts/snapshot.sh` says how to set the data
folder (`QA_DATA`), the window's size (`QA_SIZE`) and the delay (`QA_DELAY`).
`-GosmosLaunchAction` acts once the lesson has loaded:

| Action | Shows |
|---|---|
| `run`, `submit` | a Run, or a Submit (on the starter: the failure the learner meets) |
| `answer` | a quiz submitted with its first question right and the rest wrong, so both kinds of explanation show |
| `hints` | every hint that can show before a failed Submit |
| `solution`, `footer` | the solution sheet; the lesson's footer |
| `term:<term>` | a field-guide card |
| `field-guide`, `opening`, `opening:<n>` | the field guide; the first course's opening, or its page n |
| `doc:<link>`, `doc-search:<text>`, `quick-help:<text>` | the Go documentation |
| `format`, `comment`, `review`, `settings` | the editor's Format and Toggle Comment, the review, Settings |

What to look at: every readme in dark (light too for each chapter's first lesson and any
lesson with a callout, a table or a long code block), all of it (`QA_SIZE=1400x1600` for a
long one); a failed Submit on every code lesson's starter (is the failure readable in a few
lines?); every quiz after `answer`; the hints of a few lessons a chapter; each chapter's
briefing, a teaser after a pass, and the debrief; and the field guide's cards and the links
into them. Look for Markdown the app doesn't draw, readme code that wraps, inline code that
breaks badly at a line's end, a term link that opens the wrong card, and anything that
contradicts the story or an earlier lesson.

- A new data folder starts in the first course's opening. To see another course's, give the
  copy a throwaway library holding only that course (`PRAGMA foreign_keys=ON; DELETE FROM
  course WHERE slug <> '<c>'` on a copy of the built library, `$W/one.sqlite`). Mission
  Control with a chapter's briefing needs the opening marked seen: the `openings_seen` setting
  (`<library uuid>/<course uuid>`, one a line) in the data folder's `learner.sqlite`. The teaser after a pass, and the
  debrief on the pass that finishes a chapter, show after `submit` once that `learner.sqlite`
  holds the lesson's solution files in `code_file` (and, for the debrief, the chapter's other
  lessons as done in `progress`).
- A drawn screen leaves out sheets (the solution comparison), Liquid Glass buttons (blank,
  with dark labels) and the sidebar's list: judge those in the app itself. A light shot can
  come out with text missing: trust the dark one, or take it again.

## Branches

The lessons are text, so content branches diff, review and merge like code: two branches that
change different lessons merge on their own, and a conflict is in the lines both changed.
After a merge, `make library` builds the library from the merged text; run `check` and
`verify` on it as on any change.

## Working in parallel

Writers, reviewers and readers can work at the same time, one chapter each. Each gets:

- this page, `docs/storytelling.md`, the bible and the course's syllabus to read first, and
  the shipped lessons of the course (exported) as the quality bar;
- their chapter's syllabus entries and a branch or worktree of their own, whose
  `courses/apogee/` they write in and build;
- for a reviewer or a reader, the packages to check, fixed in place.

Each hands back the packages (and story files: a chapter's `chapters.md` entries and its new
`terms.md` cards), plus a short report: per lesson what was written or changed and why,
anything the bible lacked and what was assumed, and the last lines of `pkgcheck`, `check`,
`verify` and `drift`. Whoever puts the chapters together merges the story files into the one
story package (chapters in course order; terms alphabetical, one card per term) and brings the
bible and the syllabus up to date before the next chapter is written. Shared numbers, such as
booking IDs, are split between writers in advance.
