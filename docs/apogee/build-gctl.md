# Build gctl: syllabus

The guided project after Season 2 of the Apogee series, story time 25–29 October 2027. The
world, the characters, the code and the writing rules are in `story.md`, whose section "The
gctl project" fixes the program's behavior.

This is the plan the lessons were written from. All 19 are implemented: chapter 1 was read in
the app before the rest was written, and an independent review checked every lesson. The
lessons in `courses/apogee.sqlite` are the source of truth, and the larger differences from
this plan are listed under "Changes in writing and review" at the end.
`go run ./cmd/author export` writes them out to read.

- **Title:** Build gctl
- **Slug:** `build-gctl` · **Kind:** Guided Project
- **Short description:** Build gctl's interactive shell from an empty loop to the prompt the
  night shift runs a pass from.
- **Description:** The test van is home, and the night shift wants gctl without starting it
  again for every command. Over one week at Port Ellery you build one program, a lesson at a
  time: a prompt that reads commands and splits them into words, a table of commands that
  explains itself, the night's passes and stations, a stow that `next` respects, a cache
  that knows when its passes are stale and stays safe beside a background refresher, settings
  from a file and a history that survives a restart. On Friday night the whole session is
  checked against a transcript, and Ingrid runs a pass from your prompt. Standard library only.
- **Story:** a short opening, each chapter's briefing and debrief, and the field guide, told
  as `docs/storytelling.md` says.

## How a guided project differs

The learner builds one program from an empty loop to a finished tool. Each lesson adds one
feature, its starter is the previous lesson's solution, and its read-only tests check the
program as a whole, so every earlier feature stays checked. Readmes say what the feature
must do and show less code than a course's lessons do; hints carry more of the help.
Everything else follows the bible's style guide.

## At a glance

| # | Chapter | Slug | What gets built | Go | Lessons | tests | quiz |
|---|---|---|---|---|---|---|---|
| 1 | The Prompt | `the-prompt` | read, evaluate, print, loop; words; `help`; `exit` | `bufio.Scanner`, `strings.Fields`, maps of functions, sentinel errors | 5 | 4 | 1 |
| 2 | Commands | `commands` | `passes`, `stations`, `stow`, argument checks, never crashing | function values, errors, `fmt`, `recover` | 5 | 5 | 0 |
| 3 | A Cache That Knows Its Age | `a-cache-that-knows-its-age` | the night's passes loaded once, kept for a time, refreshed when stale, `next`, safe beside a refresher | closures, `time`, injected clocks, `sync.Mutex` | 5 | 4 | 1 |
| 4 | The Night Shift | `the-night-shift` | settings, history, a history file, whole sessions checked against transcripts | `encoding/json`, `os`, `testdata`, table tests | 4 | 4 | 0 |
| | **Total** | | | | **19** | **17** | **2** |

Difficulty rises from 4 to 8. Every lesson is `code_tests`: the test drives `Shell.Run` with
a `strings.Reader` of typed lines and compares what it wrote, with the clock and the files
supplied by the test (`t.TempDir()`); nothing reads the real standard input or home folder.

## Chapter 1 · The Prompt

`the-prompt` · 5 lessons · Monday 25 October

Read lines, split them into words, run commands from a table, and end cleanly.

#### 1.1 A Prompt That Waits
`a-prompt-that-waits` · code_tests · difficulty 4
- **Objective:** Read lines with `bufio.Scanner` in a loop that writes a prompt first.
- **Scenario:** Ingrid wants to type during a pass without starting gctl again.
- **Exercise:** `Shell.Run(in, out)` returns at once. Write `gctl> ` before each line, read until the input ends, and answer every line with `unknown command "<line>"; type help`.
- **Check:** Prompts before every line and one more at the end, a last line without a newline, empty input.

#### 1.2 Words, Not Lines
`words-not-lines` · code_tests · difficulty 4
- **Objective:** Split a command line into a command and its arguments with `strings.Fields`.
- **Scenario:** The night shift types two spaces as often as one.
- **Exercise:** Name only the first word in the unknown-command message; a blank line does nothing.
- **Check:** Tabs and repeated spaces, blank and space-only lines, arguments ignored in the message.

#### 1.3 A Table of Commands
`a-table-of-commands` · code_tests · difficulty 5
- **Objective:** Dispatch through a map from names to command values holding functions.
- **Scenario:** Nkem wants each new command to be one entry, not another `if`.
- **Exercise:** Look commands up in `s.commands` with comma-ok; add `help`, which lists every command and its help text sorted by name.
- **Check:** `help` output in sorted order, unknown commands still answered, a command receiving its arguments.

#### 1.4 Saying Goodbye
`saying-goodbye` · code_tests · difficulty 5
- **Objective:** End a loop with a sentinel error and return only real failures.
- **Scenario:** `exit` must leave the prompt, and a broken input must say so.
- **Exercise:** Add `exit`, returning `errExit`; `Run` prints `bye` and returns `nil` on `exit` or the end of the input, and returns the scanner's error.
- **Check:** Lines after `exit` never run; end of input says bye; a failing reader's error is returned.

#### 1.5 Quick Check: The Loop
`quick-check-the-loop` · quiz · difficulty 5
- **Objective:** Recognize how the loop reads, splits and dispatches.
- **Scenario:** Nkem reviews the prompt before the night shift sees it.
- **Exercise:** Three questions on `Scan` at the end of input, `strings.Fields` on a blank line, and a comma-ok lookup of a missing command.
- **Check:** Scan returns false at the end; Fields returns no words; a missing key gives the zero value and false.

## Chapter 2 · Commands

`commands` · 5 lessons · Tuesday 26 October

The night's passes and stations, stowing a station, arguments checked, and a prompt that
survives a broken command.

#### 2.1 Tonight's Passes
`tonights-passes` · code_tests · difficulty 5
- **Objective:** Write a command as a function over the shell's data.
- **Scenario:** The first thing the night shift asks is what flies over tonight.
- **Exercise:** Add `passes`, printing each pass of `s.passes` as `04:12-04:21 PEL pipit-1 max 61.4`, ordered by AOS, without changing the shell's slice.
- **Check:** Exact lines and order, the slice left unsorted, help listing `passes`.

#### 2.2 One Station or All
`one-station-or-all` · code_tests · difficulty 5
- **Objective:** Check a command's arguments and return errors the loop prints.
- **Scenario:** Ingrid wants only Port Ellery's passes on her screen.
- **Exercise:** `passes CODE` shows one station; an unknown code or a second argument returns an error, printed as `error: <message>`.
- **Check:** Filtering, `error:` lines for unknown codes and extra arguments, the session going on.

#### 2.3 Every Station
`every-station` · code_tests · difficulty 5
- **Objective:** Print a table from the shell's stations, in a fixed order.
- **Scenario:** The night shift checks which dishes are up before the first pass.
- **Exercise:** Add `stations`: one line per station in the table's order, `stowed: <reason>` after a stowed one; any argument is an error.
- **Check:** Exact lines, a stowed station, arguments refused.

#### 2.4 Stow a Station
`stow-a-station` · code_tests · difficulty 6
- **Objective:** Change shared state through a command, with its arguments checked.
- **Scenario:** Cape Wren's dish stows itself in wind, and gctl must know it.
- **Exercise:** Add `stow CODE [reason...]`, joining the reason words; unknown, missing and already-stowed stations are errors.
- **Check:** `stations` after a stow, the reason kept as typed, each error case, the station unchanged after an error.

#### 2.5 A Command That Panics
`a-command-that-panics` · code_tests · difficulty 6
- **Objective:** Recover a command's panic at the loop's boundary.
- **Scenario:** One bug in one command must never take the prompt away from the night shift.
- **Exercise:** In dispatch, recover a panic and print `error: internal error in <name>`, then go on.
- **Check:** A test-only panicking command, the next command still answered, errors unchanged.

## Chapter 3 · A Cache That Knows Its Age

`a-cache-that-knows-its-age` · 5 lessons · Wednesday 27 and Thursday 28 October

Load the night's passes when first needed, keep them for a while, and stay correct while a
refresher runs beside the commands.

#### 3.1 Load Once
`load-once` · code_tests · difficulty 6
- **Objective:** Load data lazily through a function value and keep it.
- **Scenario:** Reading the night's schedule file for every command is slow over Ystmark's link.
- **Exercise:** Give the shell a `passCache` whose `passes()` calls its loader on first use only; a loader error is returned and the next call tries again.
- **Check:** Loader calls counted, no load at start, errors printed and retried.

#### 3.2 Stale After Ten Minutes
`stale-after-ten-minutes` · code_tests · difficulty 7
- **Objective:** Expire cached data with an injected clock.
- **Scenario:** A schedule changed at 04:00 must reach the prompt within ten minutes.
- **Exercise:** Record when the passes were loaded with `now()`, and load again once `ttl` has passed.
- **Check:** Fresh, exactly at the time to live and stale, with a fake clock; no sleeping.

#### 3.3 What Comes Next
`what-comes-next` · code_tests · difficulty 7
- **Objective:** Search times relative to an injected now.
- **Scenario:** At 04:00 the only question is which pass is next.
- **Exercise:** Add `next`: `next: pipit-1 over PEL at 04:12, in 12m0s` for the first pass whose AOS is after now, or `no more passes tonight`.
- **Check:** Before the first pass, between passes, exactly at an AOS, after the last.

#### 3.4 The Background Refresher
`the-background-refresher` · code_tests · difficulty 8
- **Objective:** Guard shared state with `sync.Mutex`.
- **Scenario:** A refresher reloads the passes every minute while the night shift types.
- **Exercise:** Add `refresh()` and lock the cache's fields in every method, returning a copy of the slice.
- **Check:** Concurrent `passes()` and `refresh()` calls from many goroutines, correct by construction; results never share the cache's slice.

#### 3.5 Quick Check: Caches
`quick-check-caches` · quiz · difficulty 7
- **Objective:** Reason about time to live, injected clocks and locking.
- **Scenario:** Nkem asks three questions before the cache meets a real night.
- **Exercise:** Questions on an expiry boundary, why the clock is injected, and why `passes()` returns a copy.
- **Check:** Expiry at exactly the time to live; tests without sleeping; callers never hold the cache's slice.

## Chapter 4 · The Night Shift

`the-night-shift` · 4 lessons · Friday 29 October

Settings and history, and the whole night's session checked against a transcript.

#### 4.1 Settings from a File
`settings-from-a-file` · code_tests · difficulty 7
- **Objective:** Read optional settings strictly, with defaults.
- **Scenario:** The night shift wants its own prompt and Port Ellery as the default station.
- **Exercise:** Implement `loadSettings(path)`: a missing file gives the defaults; unknown fields, trailing data and bad durations are errors; `passes` with no argument uses the default station.
- **Check:** Missing, partial, full and broken files under `t.TempDir()`.

#### 4.2 Remember What You Typed
`remember-what-you-typed` · code_tests · difficulty 7
- **Objective:** Keep a session's history in a slice and print it numbered.
- **Scenario:** At 05:00 nobody remembers what was typed at 04:41.
- **Exercise:** Record every non-blank line before it runs; `history` prints them numbered from 1.
- **Check:** Numbering and padding, blank lines left out, `history` listing itself last.

#### 4.3 History That Survives
`history-that-survives` · code_tests · difficulty 7
- **Objective:** Append to a file and read it back at the next start.
- **Scenario:** The relief shift wants to see what the last shift typed.
- **Exercise:** Load the history file at start, append each line as it is typed, and report a file error once without ending the session.
- **Check:** Two sessions in a row under `t.TempDir()`, a missing file at first, an unwritable history reported.

#### 4.4 The Night Shift Runs a Pass
`the-night-shift-runs-a-pass` · code_tests · difficulty 8
- **Objective:** Check whole sessions against golden files in `testdata`.
- **Scenario:** Friday night at 04:00: Ingrid runs the night's passes from your prompt, and Cape Wren stows at 04:40.
- **Exercise:** Make `next` skip passes at stowed stations and `passes` end their lines with ` (stowed)`: the last two things the night's transcripts need. Then every session in `testdata` matches.
- **Check:** Every `testdata/*.in` session against its `.golden` transcript, with the clock moving between commands as the session file says.

## Changes in writing and review

A lesson's read-only tests can use only names its starter already has, or the starter would
not build; most of these changes follow from that.

- **1.4** `Run` starts printing a command's error as `error: <message>` here, not in 2.2, since
  1.4 is where `Run` first decides which errors end a session.
- **2.1–2.2** `passes []Pass` arrives in 2.1 and `stations []Station` in 2.2; the check that
  `passes` leaves `s.passes` unsorted runs in 2.2's tests. `passes SQT` with no pass tonight
  prints nothing.
- **2.3–2.5** `stowNote` ends a stowed station's line; 2.5's recovered panic becomes the error
  `internal error in <name>`, which `Run` prints, so errors are printed in one place.
- **3.1–3.4** 3.1 tests the cache directly and 3.2 checks the shell uses it; `refresh` arrives
  in 3.2 as the loading step, and 3.4 adds the lock, a shared `reload` and copies from
  `passes`. `main.go` starts the refresher from 3.4. 3.3's "between passes" case is 04:25.
- **4.1** The skeleton gives `decodeStrict` and `ttlText` (Season 2's strict reading); an
  unknown default station is reported by `passes` when used.
- **4.3** `add` returns an error, `load` runs when `Run` starts, and the first file error is
  printed once as `error: history: <err>`; an empty path keeps history in memory only.
- **4.4** Two sessions, `friday` and `mistakes`, replayed with a moving clock against their
  `.golden` transcripts, which read like the screen.
- **Story time** The night's passes are on Saturday 30 October UTC, with end times that fit
  each orbit and mask (the story bible's table); `help` lines for `next` and `stow` are shorter
  so they fit the readme's code blocks.
