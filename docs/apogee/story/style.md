# Apogee story bible: style and naming

## Style guide

### Voice

Second person, present tense: "You open the pass log." Concise, warm, precise. Explain like
a good colleague next to the learner: no hype, no jokes at the learner's expense, no "simply"
or "just". Prefer concrete numbers to adjectives. American spelling.

### A lesson is short and hands-on

Lessons should feel like doing, not reading. Each one teaches **one idea** and has the
learner use it right away.

- **Readme: 60 to 200 words, most about 100**, counting code. A concept that needs more than
  about 200 words becomes two lessons.
- **Shape**, always in this order:
  1. **Story beat:** one or two sentences that say why this matters today. Never a page.
  2. **The idea:** one to three sentences.
  3. **One tiny example:** at most eight lines, about something other than the task, and
     never the solution.
  4. **`## Your task`:** one to three sentences, or up to three numbered steps. Name the
     function or line to change.
  5. Optionally, one callout of one or two sentences (`> [!TIP]`, `> [!NOTE]` or
     `> [!WARNING]`) for a trap or a term.
- No H1 heading: the app shows the title. `## Your task` is usually the only heading.
- Leave out what this lesson doesn't need, even if it's true. When a learner might wonder,
  point ahead in a sentence ("Why the `&`? Next lesson.") instead of explaining now.
- No lists of rules, no history, no tangents, no "in this lesson you will learn".

Template:

````markdown
<Story beat: one or two sentences.>

<The idea: one to three sentences.>

```go
<one tiny example>
```

## Your task

<What to change or write, in one to three sentences.>
````

### Quizzes are quick checks

- A quiz follows the lesson or two whose idea it checks. It is never where an idea is taught.
- **1 to 3 questions**, each with 2 to 5 answers, exactly one right. Most have 3 answers.
- A question is one line (the format allows nothing else). Put any code the questions refer
  to in the quiz's readme, eight lines at most, and keep that readme to 60 to 120 words.
- Wrong answers are plausible mistakes, not jokes. No "all of the above", no trick wording,
  no answer that is right under some reading. Vary the position of the right answer.
- Title them "Quick Check: <topic>".

### Mix and rhythm

- About 83% of lessons are code (`code`, checked by output) or tests (`code_tests`), about
  17% quizzes.
- Use `code` when the output shows the idea at a glance; use `code_tests` when edge cases
  matter or the learner writes a function.
- Most chapters end with a slightly larger lesson that combines the chapter's ideas, still
  within the word limit.

### Code style

- gofmt, tabs, standard library only, Go 1.27. Idiomatic names (Effective Go): MixedCaps,
  initialisms in one case (`ID`, `AOS`, `LOS`, `CRC`, `HMAC`, `UTC`).
- Units in names: `AltKm`, `MaskDeg`, `FreqHz`, `BattV`, `BattMV`, `TempC`, `RateMbps`,
  `MarginDB`. Times are `time.Time` in UTC, durations `time.Duration`.
- Small functions: a solution's function is rarely longer than 15 lines.
- Every function in starter code except `main` has a one-line doc comment saying what it does.
- Errors: lowercase, no final punctuation, wrapped with `%w` when they carry context.
- No global variables except constants and the lesson's fixed data.
- Never print the wall clock, random values or map iteration order.

### Starter code and solutions

- The starter builds, runs and fails the check for the lesson's reason: wrong output or a
  failing test, never a panic. The exception is a lesson about something the compiler
  checks (an unused import, a missing method, an unexported name, an import cycle): then
  the starter fails to build with exactly the message the lesson explains, and the
  syllabus entry names it.
- The place to change is obvious: a function returning a zero value, a placeholder string,
  or at most one `// TODO:` line naming the action. The solution removes the TODO.
- The solution changes only what the task asks, usually 1 to 10 lines. It never adds a
  file: every file a solution touches is in the starter, and the solution includes every
  editable starter file, changed or not.
- Lessons in the same chapter may reuse the previous lesson's solution as the next starter.

### Tests

- Read-only `main_test.go` in `package main`, standard library only. Every `*_test.go` file
  is read-only, so a learner never edits a test file.
- Each `TestXxx` checks one function. More than three cases go in a table with names from
  the story, run with `t.Run`.
- Failure messages read `call = got, want want`: `stationCode() = "TODO", want "PEL"`. Use
  `%q` for strings, `%v` otherwise, and a tolerance for floats:
  `math.Abs(got-want) > 0.01` with `%.2f` in the message.
- `t.Errorf` when the other checks still make sense; `t.Fatalf` when they don't.
- Cover the normal case, a boundary, and the mistake learners actually make.
- Concurrency tests are correct by construction for a right solution: they never rely on
  timing luck or the race detector (the app doesn't use `-race`). Sleeps are 50 ms at most.
- In the testing chapter, where learners write checks themselves, their code goes in an
  editable non-test file (`checks.go`, `cases.go`) and the read-only test calls it. The
  test proves the checks work by running them against planted mutants in a re-executed test
  binary (`os.Args[0] -test.run=...` with an environment variable naming the mutant).

### Expected output

- Only stdout is compared, and the program must exit with status 0.
- `expected_output.txt` is exactly what the solution prints, its final newline included:
  `library verify` compares it byte for byte. The app, grading a learner, ignores line endings,
  trailing spaces and final newlines, and verify checks that the starter fails even then.
- ASCII unless the lesson is about text encoding; at most about 15 lines.

### What the app does that writers must know

- Code runs in a folder whose `go.mod` reads `module lesson` and the learner's Go version
  (`go 1.27`), unless the starter brings its own `go.mod`. A `code` lesson is built and
  run, and Submit compares its stdout; a `code_tests` lesson runs
  `go test -count=1 -timeout 10s` on Run and on Submit.
- `go test` runs part of `go vet`: starter code must pass it. Never write
  `fmt.Printf(msg)` with a variable format string.
- The app reads `go test` output for cases between lines of ten or more dashes and for lines
  like "3 passed, 1 failed". Tests must not print either.
- Every run takes well under 2 seconds on a laptop; the app stops runs at 10.

## Naming rules

- **Everything invented.** No real companies, people, products or places. Before a new name
  ships, search the web for its exact spelling. If it belongs to a real company, product,
  place, flown or announced space mission, or a person with a public profile, pick another.
  Generic words (birds, plants) may turn up in tiny or cancelled projects; record such
  overlaps below and keep the name only if nobody would mistake one for the other.
- **Name checks** (web search, 2026-11-01): nothing real was found for Furrowlight,
  Brinewatch, Cirrine Atmospherics, Kessick Polytechnic, Ystmark, Port Ellery, Salar
  Quieto, Lantern Atoll, Cape Wren, Thistle as a satellite, or the cast's full names.
  Known overlaps: a small satellite project called "Pipit" was cancelled before launch and
  never had a numbered series, so Pipit-1 to Pipit-3 stay. (A rename should keep five
  letters, so lesson 1.3's "29 bytes" stays true.) The Salar Quieto engineer was renamed
  from an earlier draft because that name belongs to a public figure.
- **Satellites:** IDs lowercase with hyphens (`pipit-1`, `furrow-a`); display names as in
  the fleet table (`Pipit-1`, `Furrow-A`).
- **Stations:** three capital letters (`YST`) and the name (`Ystmark`). Never abbreviate a
  name otherwise.
- **Characters:** first name in story beats; spellings exactly as above, accents included
  (Tomás).
- **Course, chapter and lesson slugs:** lowercase letters, digits and single hyphens, two
  to five words. Lesson slugs are unique across the whole course, not just the chapter.
- **Titles:** Title Case, two to five words, specific: "The Pre-Pass Checklist", not
  "Statements". Quizzes are "Quick Check: <topic>".
- **Files:** `main.go`, `main_test.go`; other files named for what they hold
  (`station.go`, `checks.go`); package folders are one lowercase word.
- **Topics** (front matter), from this list, alphabetical in each lesson: arrays, bytes, channels,
  comments, compiler-errors, conditionals, constants, context, conversions, defer, errors,
  files, formatting, functions, generics, goroutines, interfaces, loops, maps, methods, modules,
  multiple-returns, mutexes, packages, pointers, printing, programs, runes, select, slices,
  sorting, statements, streams, strings, structs, switch, testing, time, types, variables.
