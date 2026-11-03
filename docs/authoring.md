# Adding a course

Every Gosmos course is made the way Go Foundations was. This page is the procedure. The rules
it relies on live in the story bible (`docs/apogee/story.md`, Style guide and Naming rules)
and the library format (`docs/library-format.md`); where they say more, they win.

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

## Where things go

| What | Where |
|---|---|
| The world, the cast, the code the learner grows, the writing rules | the series' story bible: `docs/apogee/story.md` for Apogee, `docs/<series>/story.md` for a new world |
| The plan of every chapter and lesson | `docs/<series>/<course slug>.md`, in the format of `docs/apogee/go-foundations.md` |
| The first chapter written out in full | `docs/<series>/sample-chapter-1.md`, for a new series; a new Apogee course adds its own calibration chapter beside it |
| The lessons themselves | the series' library: a new Apogee course goes into `courses/apogee.sqlite` beside Go Foundations; a new series gets `courses/<series>.sqlite` |
| Lesson packages while writing | `$TMPDIR/<course slug>/packages/<chapter slug>/<lesson slug>/`, never committed |
| The story package while writing | `$TMPDIR/<course slug>/story/`, holding `opening.md`, `chapters.md` and `terms.md`, never committed |

The library is the only stored copy of a lesson and of a course's story. Once a package is
imported it is deleted; `go run ./cmd/author export` writes lessons out again when one needs
work, and `export --story` the story.

## Steps

1. **Plan.** Add to the story bible first: any new people, systems, data and names, each
   name checked against the web as its Naming rules say. Then write the syllabus: the
   course's title, slug, kind, short description and description; the at-a-glance table; and
   an entry for every lesson (objective, scenario, exercise, check, difficulty). Keep Go
   Foundations' rhythm: chapters of 7 to 10 lessons, a quiz right after the idea it checks,
   a slightly larger lesson closing most chapters, difficulty rising gently from 1 to 10.
2. **Calibrate.** Write the first chapter in full, import it and read it in the app before
   anything else is written. What it teaches goes into the story bible or the syllabus, so
   every later lesson starts right.
3. **Write.** One writer per two chapters, in parallel. Each lesson is a package in the
   format `go run ./cmd/author help` prints: front matter and readme, `starter/`,
   `solution/`, `hints.md` (2 or 3, the last opening after a failed Submit, none giving the
   answer away), `solution_note.md` (1 to 3 sentences), and `expected_output.txt` or
   `quiz.md` (every answer saying why it is right or wrong). The story bible's style guide is
   the law: one idea per lesson, a readme of 60 to 200 words ending in `## Your task`, a
   starter that builds and fails for the lesson's reason, a solution that changes only what
   the task asks, read-only tests whose messages read `call = got, want want`.
4. **Review.** An independent reviewer per two chapters checks every lesson and fixes the
   package:
   - it teaches what the syllabus says, and the Go in it is right;
   - each solution passes 30 runs in a row and under `-race`;
   - each starter fails, for the lesson's reason, within a second or two;
   - hints and quiz explanations hold up, and the style guide is kept.
5. **Tell the story.** Write the course's story as `docs/storytelling.md` says, in a story
   package in the format `go run ./cmd/author help` prints: `opening.md`, the pages the app
   shows before the first lesson; `chapters.md`, a briefing and a debrief for each chapter;
   and `terms.md`, the field guide, the world's words in plain language. It is imported with
   the lessons, by `author story`.
6. **Import.** Into a copy of the library: `author course`, then `author chapter` for each
   chapter (title, description and position from the syllabus), then `author lesson` for
   every package, then `author story` for the story package. Run
   `go run ./cmd/library check` and `go run ./cmd/library verify --teaching` on the copy,
   and copy it over the real file only when both pass. Importing a lesson again replaces it
   in place and keeps its uuid, so learners keep their progress; importing the story again
   replaces the course's whole story.
7. **See it.** Build the app (`make app`) and open representative lessons in a copy of it, in
   light and dark: a code lesson, a tests lesson with several files, a quiz after Submit,
   hints, the solution sheet, a failed Submit, and each callout kind; then the opening, a
   chapter's briefing and debrief, and a readme's links into the field guide. Check every
   readme for Markdown the app doesn't draw.
8. **Ship.** `make test`, `make check` and `make verify` pass. The story bible, the syllabus
   and the library go into one branch and one pull request; merge it, then `make install`.

## Changing a lesson later

Export its chapter with `author export`, edit the package, import it with `author lesson`,
then run `check` and `verify --teaching` again. The story changes the same way:
`author export --story` writes it out, and `author story` imports it again. Never edit the
database by hand: the tools keep its tables, keys and search index consistent.
