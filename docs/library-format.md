# Library format

Gosmos is not tied to any particular course. It opens **libraries**: SQLite files that each
hold a syllabus (courses, chapters, lessons). Any file in this format plugs in, and the app
behaves the same with it. The courses built into the app are libraries in `courses/`; File ›
Open Library… opens any other. A library is content only: what the learner does lives in
Gosmos's learner database, apart from every library (see the end of this page).

## The file

- SQLite 3, rollback journal (`journal_mode` DELETE).
- `PRAGMA application_id = 1196641619` ('GSMS'). The app refuses any other file.
- `PRAGMA user_version` is the format version, now 7. The app opens only the version it was
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
  | `type_cli` | `cli_lesson`, `cli_step` and its parts | Runs the steps in order: commands under `/bin/zsh -c` in the course's folder, HTTP requests to the learner's server, each checked by its tests |
  | `type_github_checks` | `github_lesson`, `github_step` and texts | Shows the checks as a list; the learner marks the lesson done |
  | `type_interview_chat` | `interview_lesson` | Shows the prompt and criteria; the learner marks the lesson done |

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
- **Command-line steps**:
  - `${name}` in commands, URLs, headers and bodies is replaced with values captured by
    earlier steps. `${baseURL}` is the learner's server: the address the lesson suggests
    (`base_url_default`), the learner's own address for it, or, for a lesson that
    `asks_for_base_url`, the learner's server for the course.
  - jq programs follow jq 1.8.
  - HTTP behaves like Go's `net/http` client.
- **Readmes** are Markdown (CommonMark with GitHub tables, task lists and strikethrough),
  plus these extensions:
  - Callouts: `> [!NOTE]`, `[!TIP]`, `[!IMPORTANT]`, `[!WARNING]`, `[!CAUTION]`.
    `[!INFO]` shows as a note and `[!DANGER]` as a caution. Any other word shows as a note
    titled with that word.
  - Headings tagged for one system: `## [!windows] …`.
  - Superscript in prose, never in code: `m^e^`.

  Images, videos and HTML embeds are not shown, so leave them out. Other inline HTML shows
  as literal text.

## What the app checks on open

A file that fails any of these is refused with a message that names the problem:

1. The `application_id` and `user_version` above.
2. `PRAGMA quick_check` is `ok`, and `PRAGMA foreign_key_check` finds nothing.
3. Every lesson has exactly one row in its type's table, of the matching type.
4. Every quiz question has exactly one right answer.
5. Every `cli_step` has its command or request row, matching its kind.
6. Every `solution_file` solves an editable starter file of its lesson: one neither hidden
   nor read-only.
7. `library_info` holds exactly one row.

`go run ./cmd/library check <file>` (or `make check` for this repository's library) runs the
same checks, written once in the `db` package, and prints the same messages the app shows.

## Authoring

`go run ./cmd/author` writes original content into any library: `course` and `chapter` add
or update one by slug, `lesson` imports a lesson package (a folder with the lesson's readme,
starter and solution files, hints and solution note, or its quiz with explanations), and
`export` writes lessons back out as packages, byte for byte, to edit and import again.
`go run ./cmd/author help` describes the package layout exactly.

## The learner database

Everything a learner does is in one SQLite file on their Mac,
`~/Library/Application Support/Gosmos/learner.sqlite`, with the schema in
`db/learner.sql` (`application_id` 'GSML'): the libraries they use, progress, attempts,
notes, their code, the hints and solutions they looked at, and their review schedule. Each
row names its library by `library_info.uuid` and its lesson by uuid, so updating a library
never touches the learner's data. Libraries of format 6 and older kept the learner's data
inside the library; `go run ./cmd/library migrate <file> --learner <learner.sqlite>` moves it
into the learner database and leaves the library content only. The same migration puts
`difficulty` on the 1 to 10 scale (the old score plus 6, divided by 7) and drops the lesson
columns `code_last_mod` and `synced_at`.
