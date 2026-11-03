# Library format

Gosmos is not tied to any particular course. It opens **libraries**: SQLite files that each
hold a syllabus (courses, chapters, lessons). Any file in this format plugs in, and the app
behaves the same with it. The courses built into the app are libraries in `courses/`; File ›
Open Library… opens any other. A library is content only: what the learner does lives in
Gosmos's learner database, apart from every library (see the end of this page).

Gosmos is self-contained: it never uses the network. Every lesson a library holds is checked
on the learner's Mac alone, with their Go toolchain, so the format has only the lesson types
that need nothing else: code, code with tests, and quizzes.

## The file

- SQLite 3, rollback journal (`journal_mode` DELETE).
- `PRAGMA application_id = 1196641619` ('GSMS'). The app refuses any other file.
- `PRAGMA user_version` is the format version, now 9. The app opens only the version it was
  built for; `go run ./cmd/library migrate <file>` brings an older file up to date.
- `library_info` holds the library's one row: a uuid that never changes (the learner's data
  is filed under it), a title and a version.
- The tables are exactly those in `db/schema.sql`, which documents every column. Make new
  files with `go run ./cmd/library create <file> --title T [--version V]`, not by hand.
- Any file name and folder work. The app opens a library read-only and never writes to it.

## Content

A producer fills the tables; nothing in a library belongs to a learner.

- **Syllabus**:
  - `course` (kind 'Course' or 'Guided Project'), then `chapter`, then `lesson`.
  - Each list is ordered by `position`, counted from 0, and unique within its parent.
  - The sidebar shows this order exactly.
- **Lessons**: every lesson has exactly one row in the table for its type. `difficulty` runs
  from 1 (easiest) to 10. `created_at` is when the lesson was first written and `last_mod`
  when it last changed; the author tool sets both.

  | `lesson.type` | Table | What the app does |
  |---|---|---|
  | `type_code` | `code_lesson` + `starter_file` | Builds and runs the program; it passes when stdout equals `expected_output`, ignoring `\r` before a newline, spaces and tabs at the ends of lines, and newlines at the end |
  | `type_code_tests` | `code_lesson` + `starter_file` | Runs `go test`; it passes when every test passes and the output reports no failed case |
  | `type_choice` | `quiz_lesson`, `quiz_question`, `quiz_answer` | Grades each question against its one right answer |

- **Code**: `prog_lang` is `go`, since the app runs code with the learner's Go toolchain. A
  lesson in another language still opens, and the app says it can't run it.
  Starter files with `is_hidden` are added at build time and never shown; `is_readonly`
  files are shown but can't be edited. `solution_file` holds the reference solution of each
  file the learner edits (neither hidden nor read-only).
- **Teaching**, so a learner who is stuck can move on:
  - `lesson_hint`: 2 or 3 hints per code lesson, each closer to the answer and none giving
    it; the app shows them one at a time.
  - `code_lesson.solution_note`: why the solution works, in 1 to 3 sentences, shown next to
    the solution when the learner asks for it.
  - `quiz_answer.explanation`: why each answer is right or wrong, shown after the learner
    answers.

  `go run ./cmd/library verify <file>` runs every code lesson with the Go toolchain: it must
  pass with its solution and fail with its starter, every run within 10 s. The solution must
  print `expected_output` exactly; the starter must not match it even the way the app
  compares. With `--teaching` it also requires the hints, solution notes and explanations
  above, as the built-in courses do.
- **Readmes** are Markdown (CommonMark with tables, task lists and strikethrough),
  plus these extensions:
  - Callouts: `> [!NOTE]`, `[!TIP]`, `[!IMPORTANT]`, `[!WARNING]`, `[!CAUTION]`.
    `[!INFO]` shows as a note and `[!DANGER]` as a caution. Any other word shows as a note
    titled with that word.
  - Headings tagged for one system: `## [!windows] …`.
  - Superscript in prose, never in code: `m^e^`.

  Images, videos and HTML embeds are not shown, so leave them out. Other inline HTML shows
  as literal text.
- **Story**, all of it optional: the story a course is told in. A course without one opens
  at its first lesson, and its chapters and readmes show as they are.
  - `story_page` and `story_card`: the opening, a few short pages the app shows before the
    course's first lesson, once, and again whenever the learner asks. A `text` page shows
    its cards, the people and places it introduces, below its text. A `badge` page shows
    the learner's badge: the page's `eyebrow` heads it, naming the organization the learner
    works for (a badge page without one shows no header), its cards are its lines, and a
    field takes the name the team calls the learner by. The last page's `action` opens the
    course.
  - `chapter_story`: a chapter's briefing, which opens it, and its debrief, which the app
    shows once the chapter's required lessons are done.
  - `story_term`: the field guide, the words of the course's world, each with its meaning in
    a sentence or two of plain language. The first time a readme uses a linked term, outside
    code, headings and links, the app links it to the meaning; the term matches as the guide
    spells it, as a whole word, also with a plural "s" or "es". A term lessons also use as an everyday word, such
    as pass, has `linked` 0 and is never linked.

  Their texts are short paragraphs, a blank line between each two, of inline Markdown only:
  code spans, emphasis and strong emphasis, with no headings, lists, quotes, code blocks,
  links or pictures.

## What the app checks on open

A file that fails any of these is refused with a message that names the problem:

1. The `application_id` and `user_version` above.
2. `PRAGMA quick_check` is `ok`, and `PRAGMA foreign_key_check` finds nothing.
3. Every lesson is of one of the three types above, with exactly one row in its type's table,
   of the matching type.
4. Every quiz question has exactly one right answer.
5. Every `solution_file` solves an editable starter file of its lesson: one neither hidden
   nor read-only.
6. `library_info` holds exactly one row.

`go run ./cmd/library check <file>` (or `make check` for this repository's library) runs the
same checks, written once in the `db` package, and prints the same messages the app shows.

## Authoring

`go run ./cmd/author` writes original content into any library: `course` and `chapter` add
or update one by slug, `lesson` imports a lesson package (a folder with the lesson's readme,
starter and solution files, hints and solution note, or its quiz with explanations), `story`
imports a course's story package (its opening, its chapters' briefings and debriefs, and its
field guide) in place of the story the course had, and `export` writes lessons, or with
`--story` a course's story, back out as packages, byte for byte, to edit and import again.
`go run ./cmd/author help` describes both package layouts exactly.

## The learner database

Everything a learner does is in one SQLite file on their Mac,
`~/Library/Application Support/Gosmos/learner.sqlite`, with the schema in
`db/learner.sql` (`application_id` 'GSML', now version 2): the libraries they use, progress,
attempts, notes, their code, the hints and solutions they looked at, and their review
schedule. Each row names its library by `library_info.uuid` and its lesson by uuid, so
updating a library never touches the learner's data. The app brings an older learner
database up to date when it opens it. Libraries of format 6 and older kept the learner's data
inside the library; `go run ./cmd/library migrate <file> --learner <learner.sqlite>` moves it
into the learner database and leaves the library content only. The same migration puts
`difficulty` on the 1 to 10 scale (the old score plus 6, divided by 7) and drops the lesson
columns `code_last_mod` and `synced_at`.

## Format 8: what it removed, and why

Gosmos is an isolated app: no account, no server, nothing it fetches or sends. Format 7 still
had three lesson types from the app Gosmos grew out of, and each reached outside the Mac or
could only be marked done by hand:

- `type_cli`, command-line lessons: shell commands, HTTP requests to a server the learner ran
  or deployed, and jq checks of the output (the `cli_*` tables).
- `type_github_checks`: checks of the learner's GitHub repository, its branches, pull
  requests and workflow runs (the `github_*` tables).
- `type_interview_chat`: a question to answer in one's own words (`interview_lesson`).

Format 8 drops those tables and takes only `type_code`, `type_code_tests` and `type_choice`
in `lesson.type`. `go run ./cmd/library migrate <file>` takes a library of format 7 (or older)
to 8 and keeps every row of the kept types, its `library_info` row included, so the
library's uuid, and the learner's progress filed under it, stay as they were. A library that
holds a lesson of a removed type is refused, with the count of each, and left unchanged.

The learner database's version 2 dropped what only those lessons used: the course project
folders (`project`), the learner's server addresses (`server_address`) and the attempts of
command-line checks. Migrating a library of format 6 or older drops such rows instead of
moving them, and says how many.

## Format 9: what it added, and why

Gosmos's courses are set in invented worlds, such as Apogee, a small space company. Until
format 9 that world lived only in the readmes, so each lesson set its own scene and explained
the world's words again. Format 9 gives the story tables of its own (`story_page`,
`story_card`, `chapter_story` and `story_term`), so the app can open a course with a few
pages that set the scene and hand the learner their badge, frame each chapter with a briefing
and a debrief, and explain each of the world's words once, in the field guide, linking readmes
to it. All of it is optional: a library without a story opens as before.

`go run ./cmd/library migrate <file>` takes a library of format 8 to 9 by creating the four
tables, empty. Every other row stays as it is, its `library_info` row included, so the
library keeps its uuid and the learner's progress filed under it. A library of format 7 or
older goes through format 8 on the way. `go run ./cmd/author story` then gives a course its
story.
