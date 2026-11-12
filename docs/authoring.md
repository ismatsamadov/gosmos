# Writing and fixing lessons

Everything needed to add a lesson, a chapter or a course, to work out a number the story
needs, or to fix what is there, is in this repository: the tools in `cmd/` and `scripts/`, the
plans and the world in `docs/`. Work in progress lives in a scratch folder under `$TMPDIR` and
is never needed afterwards: the library is the only stored copy of every lesson and story, and
the docs hold every decision. The rules the steps rely on live in the story bible
(`docs/apogee/story.md`, Style guide and Naming rules) and the library format
(`docs/library-format.md`); where they say more, they win.

## What every course is

- **Original.** Written for Gosmos from nothing, in an invented world. Writers never read or
  mirror another course's lessons.
- **Hands-on.** About five lessons in six are code: `code`, graded by what the program
  prints, or `code_tests`, graded by `go test`. The rest are short quizzes. There is nothing
  else: library format 9 has only these three lesson types.
- **Offline.** Standard library only, for the built-in Go (1.27.3). No lesson needs the
  network, a server, GitHub or anything outside the app. No pictures, videos, embeds or links
  that leave the app.
- **Built in.** A course ships inside a library in `courses/` and is checked like code
  before it merges.

## The tools

All run from the repository's folder.

| Command | What it does |
|---|---|
| `go run ./cmd/author help` | Prints the two package formats: a lesson package, and a course's story package (opening, briefings and debriefs, field guide) |
| `go run ./cmd/author export <library> --course <c> [--chapter <ch> \| --story] <folder>` | Writes a course's lessons, one chapter's, or its story out as packages to work on |
| `go run ./cmd/author course\|chapter\|lesson\|story <library> …` | Adds or replaces a course, a chapter, a lesson or a course's story, in place; `go run ./cmd/author` lists the flags |
| `scripts/pkgcheck.sh <folder>` | Checks lesson packages before import (below); `QUICK=1` in front for a fast first pass |
| `go run ./cmd/library check <library>` | The checks the app makes before it opens a library |
| `go run ./cmd/library verify [--course <c>] --teaching <library>` | Runs every code lesson: it must pass with its solution and fail with its starter; also wants the hints, solution notes and quiz explanations |
| `go run ./cmd/library drift [--course <c>] [--diff] <library>` | Lists given code (read-only and hidden files) that changes between lessons without a lesson changing it |
| `go run ./cmd/passes --from <day> [--to <day>] [--sat <ids>] [--station <codes>]` | The passes of the Apogee world: pass ID, AOS, LOS, maximum elevation |
| `app/scripts/snapshot.sh <out.png> <light\|dark> [launch arguments]` | Draws one screen of a copy of the app to a PNG |
| `make test`, `make check`, `make verify`, `make install` | The whole repository's tests, the libraries' checks, every lesson run, the app installed |

## Where things go

| What | Where |
|---|---|
| The world, the cast, the code the learner grows, the writing rules | the series' story bible: `docs/apogee/story.md` for Apogee, `docs/<series>/story.md` for a new world |
| The plan of every chapter and lesson | `docs/<series>/<course slug>.md`, in the format of `docs/apogee/go-foundations.md`, ending with "Changes in writing and review" |
| The first chapter written out in full | `docs/<series>/sample-chapter-1.md`, for a new series; a new Apogee course adds its own calibration chapter beside it |
| The lessons and each course's story | the series' library: Apogee's courses are in `courses/apogee.sqlite`; a new series gets `courses/<series>.sqlite` |
| Work in progress: packages, story files, a copy of the library | a scratch folder, `$TMPDIR/<course slug>/`, never committed |

Never import into `courses/apogee.sqlite` while working. Work on a copy, and copy it back only
when every check passes:

```sh
W=$TMPDIR/<course slug>
mkdir -p $W && cp courses/apogee.sqlite $W/lib.sqlite
```

Never edit a library by hand either: the tools keep its tables, keys and search index
consistent.

## Adding a lesson to a course

1. **Plan it.** Add its entry to the course's syllabus: objective, scenario, exercise, check,
   difficulty, slug and type, in its place in the chapter. Anything new in the world (a
   person, a station, a number) goes into the bible first.
2. **Export its chapter** to see its neighbors:
   `go run ./cmd/author export $W/lib.sqlite --course <c> --chapter <ch> $W/packages/<ch>`.
   A lesson whose code continues another's starts from that lesson's solution; in a guided
   project every lesson does.
3. **Write the package** in `$W/packages/<ch>/<lesson>/`, in the format `author help` prints:
   front matter and readme, `starter/`, `solution/`, `hints.md`, `solution_note.md`, and
   `expected_output.txt` or `quiz.md`. The style guide is the law. A pass in it comes from
   `cmd/passes`, any other number from the bible or from running the code (see Calculating
   the world).
4. **Make room.** Positions count from 0 and each is used once in a chapter: `author lesson`
   refuses a taken one. To put a lesson between two others, raise the `position:` of every
   lesson after it in its package and import those again, the last one first.
5. **Check it:** `scripts/pkgcheck.sh $W/packages/<ch>/<lesson>`, until it prints `ok`.
6. **Import and check the copy:**

   ```sh
   go run ./cmd/author lesson $W/lib.sqlite --course <c> --chapter <ch> $W/packages/<ch>/<lesson>
   go run ./cmd/library check $W/lib.sqlite
   go run ./cmd/library verify --course <c> --teaching $W/lib.sqlite
   go run ./cmd/library drift --course <c> $W/lib.sqlite
   ```

7. **Tell the story.** A world word the readme uses for the first time needs a card in the
   field guide; a lesson that now ends its chapter may need the debrief changed. Export the
   story (`author export … --story $W/story`), edit `terms.md` or `chapters.md`, and import it
   (`go run ./cmd/author story $W/lib.sqlite --course <c> $W/story`). See
   `docs/storytelling.md`.
8. **See it in the app** (below).
9. **Ship it.** `cp $W/lib.sqlite courses/apogee.sqlite`, note in the syllabus what changed
   from its plan, then `make test`, `make check` and `make verify`. The syllabus, any bible
   change and the library go into one branch and one pull request.

## Adding a chapter

`go run ./cmd/author chapter $W/lib.sqlite --course <c> --slug <ch> --title "<title>"
--description "<description>" --position <n>` with the syllabus's title and description, then
each lesson as above, and a briefing and a debrief for the chapter in the story's
`chapters.md` (at most 70 and 50 words; `docs/storytelling.md` says what they hold).

## Adding a course

Every Gosmos course is made the way Go Foundations was:

1. **Plan.** Add to the story bible first: any new people, systems, data and names, each
   name checked against the web as its Naming rules say. Then write the syllabus: the
   course's title, slug, kind, short description and description; the at-a-glance table; and
   an entry for every lesson. Keep Go Foundations' rhythm: chapters of 7 to 10 lessons, a quiz
   right after the idea it checks, a slightly larger lesson closing most chapters, difficulty
   rising gently from 1 to 10.
2. **Calibrate.** Write the first chapter in full, import it and read it in the app before
   anything else is written. What it teaches goes into the bible or the syllabus, so every
   later lesson starts right. A guided project builds its whole program and any hidden test
   harness first: each lesson's solution is then a cut of that program, with that lesson's
   features and no later ones (the capstone's section of the bible shows how).
3. **Write.** Chapters can be written in parallel, one writer per chapter (see Working in
   parallel). Each lesson is written as in Adding a lesson.
4. **Review.** Someone who didn't write a chapter checks every lesson of it against the
   checklist below and fixes the packages.
5. **Tell the story** as `docs/storytelling.md` says: `opening.md`, the pages before the
   first lesson; `chapters.md`, a briefing and a debrief for each chapter; and `terms.md`, the
   field guide. `author course` makes the course, `author story` imports its story.
6. **Import** into the copy: `author course`, then `author chapter` for each chapter, then
   `author lesson` for every package, then `author story`. Run `check`, `verify --teaching`
   and `drift` on the copy. Importing a lesson again replaces it in place and keeps its uuid,
   so learners keep their progress; importing the story again replaces the whole story.
7. **See it** in the app, every chapter.
8. **Ship.** As for a lesson. Then write the syllabus's "Authoring status" and "Changes in
   writing and review", and bring the bible up to date with what writing and review settled.

## Fixing a lesson

1. Export its chapter (and its story, if the fix touches it), edit the package, run
   `pkgcheck`, import it into the copy, and run `check`, `verify --teaching` and `drift` for
   the course. Importing replaces the lesson in place and keeps its uuid, so learners keep
   their progress.
2. **Code that carries forward.** When later lessons give this lesson's code as read-only or
   hidden files, or start from its solution (every lesson of a guided project does), the same
   fix goes into each of them: their starters and their solutions. `library drift --course
   <c> --diff` lists every place the given code no longer matches what the learner last saw.
3. **Shared hidden files.** A harness hidden in many lessons, such as the capstone's
   `sim_*_test.go` files, stays byte for byte the same in every lesson that has it: change
   every copy, then verify the whole course.
4. **The story.** `author export … --story`, edit, `author story`.
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
returned. A change a lesson or the story asks for is fine: say why in the pull request. The
courses before Season 6 were written before this check existed, and it still reports changes
in them.

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

- **`scripts/pkgcheck.sh`**, on packages before import: the front matter; a readme of 60 to
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
copy of it with your library:

```sh
L=$(sqlite3 $W/lib.sqlite "SELECT l.uuid FROM lesson l JOIN chapter ch ON ch.uuid = l.chapter_uuid
  WHERE ch.slug = '<ch>' AND l.slug = '<lesson>'")
QA_LIBRARY=$W/lib.sqlite app/scripts/snapshot.sh $W/shots/readme.png dark -GosmosOpenLesson $L
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
  course WHERE slug <> '<c>'` on a copy of your copy). Mission Control with a chapter's
  briefing needs the opening marked seen: the `openings_seen` setting (`<library uuid>/<course
  uuid>`, one a line) in the data folder's `learner.sqlite`. The teaser after a pass, and the
  debrief on the pass that finishes a chapter, show after `submit` once that `learner.sqlite`
  holds the lesson's solution files in `code_file` (and, for the debrief, the chapter's other
  lessons as done in `progress`).
- A drawn screen leaves out sheets (the solution comparison), Liquid Glass buttons (blank,
  with dark labels) and the sidebar's list: judge those in the app itself. A light shot can
  come out with text missing: trust the dark one, or take it again.

## Branches: the library is one binary file

Git can't merge two branches that both change `courses/apogee.sqlite`. Keep one content branch
open at a time. When two must overlap, base the second on the first and merge them in that
order (retarget the second to `main` before the first's branch is deleted). To carry a fix
from the lower branch into the higher one:

1. `git merge --no-commit <lower>`, then `git checkout --ours courses/apogee.sqlite`.
2. Import the same fixed packages and story into it (`author lesson`, `author story`, `author
   chapter` or `author course` for a changed row).
3. `git add courses/apogee.sqlite` again after importing (adding it before keeps the old
   file), then commit.
4. Check the committed file, not the working one: export it from `git show
   <branch>:courses/apogee.sqlite` and compare with the fixed packages, then run `check` and
   `verify`.

A docs-only branch whose library is identical to its base's merges as usual.

## Working in parallel

Writers, reviewers and readers can work at the same time, one chapter each. Each gets:

- this page, `docs/storytelling.md`, the bible and the course's syllabus to read first, and
  the shipped lessons of the course (exported) as the quality bar;
- their chapter's syllabus entries and a fresh copy of the library to import into;
- for a reviewer or a reader, the packages to check, fixed in place.

Each hands back the packages (and story files: a chapter's `chapters.md` entries and its new
`terms.md` cards), plus a short report: per lesson what was written or changed and why,
anything the bible lacked and what was assumed, and the last lines of `pkgcheck`, `check`,
`verify` and `drift`. Whoever puts the chapters together merges the story files into the one
story package (chapters in course order; terms alphabetical, one card per term) and brings the
bible and the syllabus up to date before the next chapter is written. Shared numbers, such as
booking IDs, are split between writers in advance.
