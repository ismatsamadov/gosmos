# Go Foundations: syllabus

The first course of the Apogee series. The world, the characters, the code and the writing
rules are in `story.md`; lessons 1.1 to 1.5 are written in full in `sample-chapter-1.md`.

This is the plan the lessons were written from. The lessons themselves, in
`courses/apogee.sqlite`, are the source of truth: writing and review changed details such as
numbers, times and extra test cases. `go run ./cmd/author export` writes them out to read.

- **Title:** Go Foundations
- **Slug:** `go-foundations` · **Kind:** Course
- **Short description:** Learn Go from your first line to concurrent code as the newest
  backend engineer at Apogee, a ground-station startup.
- **Description:** You've just joined Apogee, whose five ground stations talk to small
  satellites from the Arctic to the far south. Over one season you go from printing your
  first message to Pipit-1 to running a whole pass day concurrently through a storm. Along
  the way you schedule passes, parse and validate telemetry, model the fleet and the
  stations, sign command uplinks, reassemble downlinked files and handle the anomalies that
  wake the on-call engineer. Every lesson is one idea and one small, real task, checked by
  your own Go toolchain with nothing but the standard library.

## At a glance

| # | Chapter | Slug | Go topics | Lessons | code | tests | quiz |
|---|---|---|---|---|---|---|---|
| 1 | First Contact | `first-contact` | programs, printing, comments, compiler errors, reading tests | 8 | 6 | 1 | 1 |
| 2 | Numbers from Orbit | `numbers-from-orbit` | variables, types, zero values, conversions, constants | 9 | 5 | 2 | 2 |
| 3 | The Pass Board | `the-pass-board` | formatting with fmt | 7 | 5 | 1 | 1 |
| 4 | Orbit Math | `orbit-math` | functions, multiple returns, math, durations | 9 | 4 | 4 | 1 |
| 5 | Go or No-Go | `go-or-no-go` | if, else, logical operators, switch | 8 | 3 | 4 | 1 |
| 6 | Counting Down to AOS | `counting-down-to-aos` | loops, break, continue | 8 | 4 | 3 | 1 |
| 7 | Modeling the Network | `modeling-the-network` | structs, time.Time, embedding | 7 | 3 | 3 | 1 |
| 8 | The Day's Schedule | `the-days-schedule` | slices, sorting, variadic functions | 10 | 4 | 4 | 2 |
| 9 | Who's Who in the Fleet | `whos-who-in-the-fleet` | maps, sets | 8 | 3 | 4 | 1 |
| 10 | Changing Things in Place | `changing-things-in-place` | pointers, methods, iota, Stringer | 9 | 4 | 3 | 2 |
| 11 | Every Radio Speaks | `every-radio-speaks` | interfaces, io.Writer, type switches, assertions | 8 | 3 | 4 | 1 |
| 12 | When the Frame Is Wrong | `when-the-frame-is-wrong` | errors, wrapping, custom errors, defer | 10 | 2 | 6 | 2 |
| 13 | Reading the Raw Feed | `reading-the-raw-feed` | strings, strconv, runes, bytes, HMAC | 8 | 2 | 5 | 1 |
| 14 | Trust, but Verify | `trust-but-verify` | testing basics | 8 | 0 | 7 | 1 |
| 15 | Splitting groundctl | `splitting-groundctl` | packages, modules | 7 | 4 | 1 | 2 |
| 16 | One Helper, Every Type | `one-helper-every-type` | generics | 7 | 1 | 5 | 1 |
| 17 | All Stations at Once | `all-stations-at-once` | goroutines, channels, select | 10 | 4 | 4 | 2 |
| 18 | Locks and Deadlines | `locks-and-deadlines` | mutexes, context | 9 | 1 | 6 | 2 |
| | **Total** | | | **150** | **58** | **67** | **25** |

Mix: 38.7% code, 44.7% code_tests, 16.7% quizzes (83.3% hands-on). Every quiz is a quick
check of one to three questions that comes right after the idea it checks.

## How to read a lesson entry

Each lesson is one idea. The entry gives its slug, type and difficulty (1 to 10), then:

- **Objective:** what the learner can do afterwards.
- **Scenario:** the story beat, in one sentence.
- **Exercise:** what the starter has and what the learner writes.
- **Check:** the expected output (`code`), the test cases (`code_tests`) or the right
  answers (quiz). Numbers come from the formulas and constants in `story.md`; floats are
  checked to within 0.01 unless an entry says otherwise.

Starters fail for the lesson's reason and build, except where an entry names the compiler
message the lesson is about. Readmes stay within 60 to 200 words (`story.md`, Style guide).

---

## Chapter 1 · First Contact

`first-contact` · 8 lessons: 6 code, 1 code_tests, 1 quiz

Your first Go programs: printing a line, the order statements run in, comments, reading a
compiler error, and how a lesson checked by tests works.

**Story arc.** Day one. Ingrid's tradition says a new engineer's first program writes the
message the ops room sends up to Pipit-1's echo buffer. You print it, put the pre-pass
checklist in order, switch the debug noise out of the customer log and meet the tests the
scheduler relies on. By Friday your words have gone up and come back on Pipit-1's beacon,
and the Fishbowl has a roll call of all five stations.

#### 1.1 Hello, Pipit
`hello-pipit` · code · difficulty 1
- **Objective:** Run and submit a program that prints one line with `fmt.Println`.
- **Scenario:** Your first program writes the message the ops room sends up to Pipit-1.
- **Exercise:** The starter prints a placeholder; you change the text between the quotes.
- **Check:** `Apogee ground calling Pipit-1`

#### 1.2 The Pre-Pass Checklist
`pre-pass-checklist` · code · difficulty 1
- **Objective:** See that statements run from top to bottom, one line per `Println`.
- **Scenario:** Ten minutes before AOS, the checklist printer has shuffled Ingrid's steps.
- **Exercise:** The starter prints a title and four numbered steps out of order; you reorder
  the calls.
- **Check:** The title, then steps 1 to 4 (five lines; full text in the sample).

#### 1.3 Notes for the Night Shift
`notes-for-the-night-shift` · code · difficulty 1
- **Objective:** Use `//` comments for notes and to switch a line off.
- **Scenario:** Nkem wants the DEBUG lines out of the customer-facing pass log, but kept.
- **Exercise:** The starter prints two DEBUG lines among three log lines; you comment them
  out.
- **Check:** The three log lines only.

#### 1.4 Checked by Tests
`checked-by-tests` · code_tests · difficulty 2
- **Objective:** Read a test to find out what a function must return.
- **Scenario:** The scheduler calls your code instead of reading what it prints.
- **Exercise:** `stationCode()` returns `"TODO"`; you read `main_test.go` and return the
  value it wants.
- **Check:** `TestStationCode`: `stationCode() == "PEL"`.

#### 1.5 Quick Check: First Programs
`quick-check-first-programs` · quiz · difficulty 1
- **Objective:** Confirm where a program starts, what a comment does and what a test
  controls.
- **Scenario:** Ingrid's check before you sit at a live console.
- **Exercise:** Three questions about a six-line program and a test lesson.
- **Check:** It starts at `func main()`; only `AOS in 2 minutes` prints; you change the
  function in `main.go`, not the test.

#### 1.6 It Won't Build
`it-wont-build` · code · difficulty 2
- **Objective:** Read a compiler error and fix what it names.
- **Scenario:** The banner tool won't build, and AOS is in six minutes.
- **Exercise:** The starter imports `"os"` without using it; you remove the import.
- **Check:** `Banner: Pipit-1 over Port Ellery, AOS 04:12 UTC`. The starter fails with
  `"os" imported and not used`.

#### 1.7 Several Values, One Line
`several-values-one-line` · code · difficulty 2
- **Objective:** Pass several values, text and numbers, to one `Println`, which puts spaces
  between them.
- **Scenario:** The beacon decoder prints each field on its own line; ops wants one line per
  fact.
- **Exercise:** The starter prints `Pipit-1`, `battery`, `7.91` and `V` with four calls, and
  the mode with three; you use one call per line.
- **Check:** `Pipit-1 battery 7.91 V` / `Pipit-1 mode NOMINAL`.

#### 1.8 Station Roll Call
`station-roll-call` · code · difficulty 2
- **Objective:** Combine `Println` calls with several values into a small report.
- **Scenario:** Friday stand-up: Ingrid wants all five stations and their state on the Board.
- **Exercise:** The starter prints the title and Ystmark; you add the other four stations.
- **Check:**
  ```text
  Station roll call
  YST Ystmark up
  PEL Port Ellery up
  SQT Salar Quieto up
  LAN Lantern Atoll up
  CWR Cape Wren stowed
  ```

---

## Chapter 2 · Numbers from Orbit

`numbers-from-orbit` · 9 lessons: 5 code, 2 code_tests, 2 quizzes

Variables and Go's basic types (`int`, `float64`, `string`, `bool`), zero values, integer
and floating-point arithmetic, explicit conversions and constants.

**Story arc.** Ingrid hands you the onboarding sheet: altitudes, frequencies, voltages and
orbit periods, all typed in by hand and slightly wrong. You move each number into a
variable or a constant, find out why Go won't mix an `int` and a `float64` without being
asked, and watch Tomás derive an orbit's period from its altitude and a single constant. By
the end the sheet computes itself.

#### 2.1 Numbers to Know
`numbers-to-know` · code · difficulty 2
- **Objective:** Declare a variable with `:=` and print it.
- **Scenario:** Pipit-1's altitude appears on the sheet three times, typed three ways.
- **Exercise:** The starter declares `name := "Pipit-1"` and prints it; you declare
  `altitudeKm := 525` and print it with a label.
- **Check:** `Pipit-1` / `Altitude km: 525`.

#### 2.2 Every Value Has a Type
`every-value-has-a-type` · code · difficulty 2
- **Objective:** Declare variables with `var`, a type and a value.
- **Scenario:** The sheet's four facts about Pipit-1 each need the right type.
- **Exercise:** You add `var altitudeKm int = 525`, `var beaconMHz float64 = 437.25` and
  `var deployed bool = true` beside the given name, and print all four.
- **Check:** `Name: Pipit-1` / `Altitude km: 525` / `Beacon MHz: 437.25` / `Deployed: true`.

#### 2.3 Nothing Yet Is Still Something
`zero-values` · code · difficulty 2
- **Objective:** Predict the zero value of a declared variable, and assign with `=`.
- **Scenario:** A pass record created before the scheduler fills it shows 0, false and
  empty text.
- **Exercise:** The starter declares `station`, `maxElevDeg` and `confirmed` with `var` and
  prints them; you assign only `station = "PEL"`.
- **Check:** `Station: PEL` / `Max elevation: 0` / `Confirmed: false`.

#### 2.4 Quick Check: Types
`quick-check-types` · quiz · difficulty 2
- **Objective:** Check types and zero values.
- **Scenario:** Ingrid checks you before you edit the live sheet.
- **Exercise:** Three questions.
- **Check:** `437.25` is a `float64`; a `string`'s zero value is `""`; `:=` declares a new
  variable and gives it a value.

#### 2.5 Whole Orbits and Leftovers
`whole-orbits` · code · difficulty 2
- **Objective:** Use integer division `/` and the remainder operator `%`.
- **Scenario:** Ingrid asks how many full orbits Pipit-1 makes in a day, and what's left.
- **Exercise:** The starter has `minutesPerDay := 1440` and `periodMin := 95` and prints
  zeros; you compute both numbers.
- **Check:** `Full orbits: 15` / `Minutes left: 15`.

#### 2.6 Floats Meet Ints
`floats-meet-ints` · code_tests · difficulty 3
- **Objective:** Convert an `int` to `float64` before dividing.
- **Scenario:** The dashboard says Furrowlight averages 333 MB per pass; it's really 333.33.
- **Exercise:** `averageMB(totalMB, passes int) float64` divides, then converts; you
  convert first.
- **Check:** `averageMB(1000, 3)` ≈ 333.333 (within 0.001); `averageMB(900, 3) == 300`;
  `averageMB(10, 4) == 2.5`.

#### 2.7 Constants of Nature
`constants-of-nature` · code · difficulty 3
- **Objective:** Name a fixed value once with `const`.
- **Scenario:** Tomás finds Earth's radius typed as 6371, 6378 and 6370 in three tools.
- **Exercise:** The starter adds three different radii to three altitudes; you declare
  `const earthRadiusKm = 6371` and use it everywhere.
- **Check:** `Pipit-1 orbit radius: 6896 km` / `Thistle orbit radius: 6786 km` /
  `Cirrine-1 orbit radius: 6971 km`.

#### 2.8 Quick Check: Constants
`quick-check-constants` · quiz · difficulty 3
- **Objective:** Check what a constant may hold and how an untyped constant mixes.
- **Scenario:** Tomás asks two questions before he trusts the new constants.
- **Exercise:** Two questions about a three-line `const` block in the readme.
- **Check:** An untyped constant like `6371` works as an `int` or a `float64`;
  `const now = time.Now()` doesn't compile, because a constant is fixed when compiling.

#### 2.9 The Period of an Orbit
`period-from-altitude` · code_tests · difficulty 4
- **Objective:** Write a formula with `float64` arithmetic, constants and the `math`
  package.
- **Scenario:** Tomás shows that an orbit's period follows from its altitude and one
  constant, μ.
- **Exercise:** `periodMinutes(altKm float64) float64` returns 0; you return
  2π·√(a³/μ)/60 with a = 6371 + altKm, using `math.Sqrt` and `math.Pi`.
- **Check:** `periodMinutes(525)` ≈ 94.985; `(415)` ≈ 92.722; `(600)` ≈ 96.539.

---

## Chapter 3 · The Pass Board

`the-pass-board` · 7 lessons: 5 code, 1 code_tests, 1 quiz

Formatted output with `fmt`: `Printf` and its verbs, decimal places, widths and alignment,
zero padding, `Sprintf`, and the verbs that help you debug.

**Story arc.** The Board, the big screen in the Fishbowl, lists the next passes in text
nobody can read: times without leading zeros, frequencies with random decimals, columns
that wander. Furrowlight visits on Friday. You fix the Board one verb at a time, and on
Friday morning it lines up.

#### 3.1 First Verbs
`first-verbs` · code · difficulty 2
- **Objective:** Print a line with `Printf`, `%s`, `%d` and `\n`.
- **Scenario:** The Board's header needs a satellite's name and pass count in one sentence.
- **Exercise:** The starter prints the pieces on separate lines; you write one `Printf`.
- **Check:** `Next passes for Pipit-1: 4`.

#### 3.2 Three Decimals, Always
`three-decimals` · code · difficulty 2
- **Objective:** Set decimal places with `%.3f` and `%.1f`.
- **Scenario:** Operators read frequencies to three decimals; the Board shows 437.25 and 2245.
- **Exercise:** The starter prints with `%v`; you change the verbs.
- **Check:** `Pipit-1 437.250 MHz` / `Pipit-2 2245.000 MHz` / `Max elevation 61.4 deg`
  (the value is 61.37).

#### 3.3 Columns That Line Up
`columns-that-line-up` · code · difficulty 3
- **Objective:** Pad values to a width with `%-13s` and `%5.1f`.
- **Scenario:** The Board's columns wander because station names differ in length.
- **Exercise:** The starter prints rows with single spaces; you add widths to the header and
  the rows.
- **Check:**
  ```text
  STATION        MASK
  Ystmark         5.0
  Salar Quieto    3.0
  Lantern Atoll   8.0
  ```

#### 3.4 Leading Zeros
`leading-zeros` · code_tests · difficulty 3
- **Objective:** Build a string with `Sprintf` and pad numbers with `%02d`.
- **Scenario:** The Board shows AOS as 4:7 instead of 04:07.
- **Exercise:** `clock(h, m int) string` returns `fmt.Sprint(h, ":", m)`; you return a
  zero-padded `hh:mm`.
- **Check:** `clock(4, 7) == "04:07"`; `clock(23, 59) == "23:59"`; `clock(0, 0) == "00:00"`.

#### 3.5 Quick Check: Verbs
`quick-check-verbs` · quiz · difficulty 2
- **Objective:** Read format strings at a glance.
- **Scenario:** Ingrid points at three lines of the Board's code.
- **Exercise:** Three questions: what each of three `Printf` calls prints.
- **Check:** `%5.1f` of 7.91 is `  7.9`; `%-4s|` of `"YST"` is `YST |`; `%03d` of 7 is `007`.

#### 3.6 Seeing Hidden Spaces
`seeing-hidden-spaces` · code · difficulty 3
- **Objective:** Show a string's exact contents with `%q` and a value's type with `%T`.
- **Scenario:** Lookups for Salar Quieto fail: someone typed its code as `"SQT "`.
- **Exercise:** The starter prints three codes with `%s`; you switch to `%q` and add a `%T`
  line for the mask.
- **Check:** `"YST"` / `"SQT "` / `"LAN"` / `mask type: float64`.

#### 3.7 The Board
`the-board` · code · difficulty 4
- **Objective:** Combine widths, decimals and zero padding in one table.
- **Scenario:** Friday, 08:00: the Board must be readable before Furrowlight arrives.
- **Exercise:** The starter holds four passes in variables and prints them with `Println`;
  you print the header with `"%-20s %-8s %-6s %6s  %s\n"` and each row with
  `"%-20s %-8s %02d:%02d  %6.1f  %s\n"`.
- **Check:**
  ```text
  PASS                 STATION  AOS    MAX EL  BAND
  pipit-1/PEL/04182    PEL      04:12    61.4  UHF
  furrow-a/SQT/11873   SQT      05:03    23.0  X
  thistle/LAN/02218    LAN      06:47     8.9  UHF
  cirrine-1/YST/07455  YST      07:30    44.2  S
  ```

---

## Chapter 4 · Orbit Math

`orbit-math` · 9 lessons: 4 code, 4 code_tests, 1 quiz

Functions: parameters and results, several results at once, the blank identifier, calling
the `math` package, building functions on functions, and `time.Duration`.

**Story arc.** Tomás keeps Apogee's orbit formulas in a battered notebook, and every tool
has its own copy of them, each slightly different. You turn the formulas into small
functions: orbits per day, Doppler shift, slant range, path loss. Halfway through you meet
functions that return two answers at once. By the end the notebook is a toolbox, and Tomás
prints a card from it before Pipit-2's pass over Salar Quieto.

#### 4.1 Your First Function
`your-first-function` · code · difficulty 3
- **Objective:** Declare a function with a parameter and a result, and call it.
- **Scenario:** Three tools compute orbits per day three different ways.
- **Exercise:** The starter prints three hard-coded, wrong figures; you write
  `orbitsPerDay(periodMin float64) float64` and call it for each satellite.
- **Check:** `Pipit-1 15.16` / `Thistle 15.53` / `Cirrine-1 14.92` (periods 94.985, 92.722
  and 96.539, printed with `%.2f`).

#### 4.2 Doppler Shift
`doppler-shift` · code_tests · difficulty 3
- **Objective:** Write a function of two parameters that returns a computed value.
- **Scenario:** Pipit-1's beacon arrives about 10 kHz high as it rises; Tomás wants the
  shift computed, not guessed.
- **Exercise:** `dopplerHz(freqHz, rangeRateMps float64) float64` returns 0; you return
  `-freqHz * rangeRateMps / speedOfLight` (the constant is given).
- **Check:** `(437.25e6, -7000)` ≈ +10209.56; `(437.25e6, 7000)` ≈ −10209.56;
  `(2245e6, -6500)` ≈ +48675.34; `(2245e6, 0) == 0`.

#### 4.3 Two Answers at Once
`two-answers-at-once` · code_tests · difficulty 3
- **Objective:** Return two values from one function.
- **Scenario:** The scheduler stores times as minutes after midnight; the Board wants hours
  and minutes.
- **Exercise:** `hoursMinutes(total int) (int, int)` returns `0, 0`; you return both parts.
- **Check:** 252 → (4, 12); 0 → (0, 0); 1439 → (23, 59).

#### 4.4 Only the One You Need
`only-the-one-you-need` · code · difficulty 3
- **Objective:** Discard a result with the blank identifier `_`.
- **Scenario:** The pass banner needs only the minutes, and the build fails.
- **Exercise:** `h, m := hoursMinutes(252)` leaves `h` unused; you replace `h` with `_`.
- **Check:** `Minutes past the hour: 12`. The starter fails with `declared and not used: h`.

#### 4.5 Quick Check: Functions
`quick-check-functions` · quiz · difficulty 3
- **Objective:** Check how arguments and results behave.
- **Scenario:** Nkem's three questions about a function that changes its parameter.
- **Exercise:** Three questions about a five-line snippet in the readme.
- **Check:** Arguments are copies, so the caller's variable is unchanged; `a, b := f()`
  needs `f` to return exactly two values; result types come after the parameters.

#### 4.6 Slant Range
`slant-range` · code_tests · difficulty 4
- **Objective:** Call `math.Sqrt`, `math.Sin` and `math.Cos` from your own function.
- **Scenario:** Link budgets need the distance to the satellite, which grows as it sinks
  toward the horizon.
- **Exercise:** `radians(deg)` is given; `slantRangeKm(altKm, elevDeg float64) float64`
  returns 0; you return √((R+h)² − (R·cos e)²) − R·sin e.
- **Check:** `(525, 90)` ≈ 525.00; `(525, 0)` ≈ 2639.16; `(525, 30)` ≈ 951.24;
  `(415, 10)` ≈ 1479.05.

#### 4.7 Path Loss
`path-loss` · code_tests · difficulty 5
- **Objective:** Build one function on top of another.
- **Scenario:** Tomás wants the free-space path loss of a pass, which starts from the slant
  range.
- **Exercise:** `slantRangeKm` is solved; `pathLossDB(altKm, elevDeg, freqMHz float64)
  float64` returns 0; you return 20·log₁₀(d) + 20·log₁₀(f) + 32.44 using `slantRangeKm` and
  `math.Log10`.
- **Check:** `(525, 0, 437.25)` ≈ 153.68; `(525, 90, 437.25)` ≈ 139.66;
  `(550, 10, 2245)` ≈ 164.64.

#### 4.8 Durations, Not Numbers
`durations-not-numbers` · code · difficulty 4
- **Objective:** Use `time.Duration` and its constants instead of bare numbers.
- **Scenario:** A tool mixed up seconds and minutes and booked a nine-hour pass.
- **Exercise:** `passLength(minutes, seconds int) time.Duration` returns
  `time.Duration(minutes + seconds)`, which counts nanoseconds; you multiply by
  `time.Minute` and `time.Second`. `main` prints three lengths with `%v`.
- **Check:** `9m12s` / `45s` / `12m0s`.

#### 4.9 The Orbit Card
`the-orbit-card` · code · difficulty 5
- **Objective:** Combine your functions in one program.
- **Scenario:** Before Pipit-2's pass over Salar Quieto, Tomás wants every number on one card.
- **Exercise:** The four functions are solved; `main` prints only the title; you print the
  other four lines.
- **Check:**
  ```text
  Pipit-2 over Salar Quieto
  Period: 95.50 min
  Slant range at 10 deg: 1815.08 km
  Doppler at -6.5 km/s: +48.68 kHz
  Path loss at 10 deg: 164.64 dB
  ```

---

## Chapter 5 · Go or No-Go

`go-or-no-go` · 8 lessons: 3 code, 4 code_tests, 1 quiz

Conditionals: `if`, `else` and `else if`, comparison and logical operators, and `switch`
with and without a value.

**Story arc.** Before every pass the Fishbowl runs a go/no-go poll. Is the station up? Will
the satellite clear the mask? Is the battery healthy? Is the wind calm at Cape Wren? Ingrid
wants the poll automated, one condition at a time. In the last lesson your poll makes the
call for Furrow-A's evening pass over Salar Quieto.

#### 5.1 If It's Too Low
`if-its-too-low` · code · difficulty 2
- **Objective:** Run a statement only when a condition holds, with `if`.
- **Scenario:** The pass log should warn about low passes, and only about those.
- **Exercise:** The starter prints a warning after every pass; you wrap each in
  `if maxElev < 10`.
- **Check:** `Pass 1: max 42.0 deg` / `Pass 2: max 7.5 deg` /
  `  low pass: expect a weak signal` / `Pass 3: max 15.2 deg`.

#### 5.2 Above the Mask
`above-the-mask` · code_tests · difficulty 3
- **Objective:** Choose between two results with `if` and `else`.
- **Scenario:** A pass is worth scheduling only if its highest point clears the mask.
- **Exercise:** `elevationCall(maxElevDeg, maskDeg float64) string` returns `""`; you return
  `"GO"` or `"NO-GO"`.
- **Check:** (42, 5) GO; (4.5, 5) NO-GO; (5, 5) GO, because clearing means at or above.

#### 5.3 Battery Bands
`battery-bands` · code_tests · difficulty 3
- **Objective:** Chain conditions with `else if`.
- **Scenario:** Before every pass, ops sorts the battery into FULL, OK, LOW or CRITICAL.
- **Exercise:** `batteryBand(v float64) string` returns `""`; you return the band: FULL from
  8.2 V, OK from 7.2 V, LOW from 6.6 V.
- **Check:** 8.3 FULL; 8.2 FULL; 7.4 OK; 6.6 LOW; 6.59 CRITICAL.

#### 5.4 All Must Agree
`all-must-agree` · code · difficulty 3
- **Objective:** Combine conditions with `&&`, `||` and `!`.
- **Scenario:** A station can track only if it's up, its dish isn't stowed and the wind is
  under 72 km/h.
- **Exercise:** `canTrack(up, stowed bool, windKmh float64) bool` returns `true`; you write
  the combined condition. `main` prints it for three stations.
- **Check:** `YST true` / `SQT false` / `CWR false` (Salar Quieto is down; Cape Wren has
  80 km/h of wind).

#### 5.5 Quick Check: Conditions
`quick-check-conditions` · quiz · difficulty 2
- **Objective:** Check how conditions are written and evaluated.
- **Scenario:** Ingrid's quick check before the poll goes live.
- **Exercise:** Three questions.
- **Check:** `if` needs braces even around one statement; `x > 5 && x < 10` is true for 7;
  in `false && f()`, `f` is never called.

#### 5.6 Switch on Mode
`switch-on-mode` · code · difficulty 3
- **Objective:** Pick one of many cases with `switch`.
- **Scenario:** What ops does during a pass depends on the satellite's mode.
- **Exercise:** `modeAction(mode string) string` returns `"hold"`; you write a switch where
  `DETUMBLE` and `LEOP` share a case. `main` prints five modes.
- **Check:** `NOMINAL: proceed` / `SAFE: recovery pass` / `DETUMBLE: beacon only` /
  `LEOP: beacon only` / `REBOOTING: hold`.

#### 5.7 The Wind at Cape Wren
`the-wind-at-cape-wren` · code_tests · difficulty 4
- **Objective:** Use a `switch` with no value as a tidy chain of conditions.
- **Scenario:** Cape Wren tracks normally under 50 km/h, with a warning from 50, and stows
  from 72.
- **Exercise:** `windCall(kmh float64) string` returns `""`; you write `switch { case ... }`.
- **Check:** 30 → `track`; 50 → `track with warning`; 71.9 → `track with warning`;
  72 → `stow`; 95 → `stow`.

#### 5.8 The Go/No-Go Poll
`the-go-no-go-poll` · code_tests · difficulty 5
- **Objective:** Order checks so that the first failing reason wins.
- **Scenario:** Furrow-A's evening pass over Salar Quieto: your poll makes the call.
- **Exercise:** `poll(up bool, windKmh, maxElevDeg, battV float64) string` returns `"GO"`;
  you return the first of `NO-GO: station down`, `NO-GO: wind`, `NO-GO: below mask` (3°)
  and `NO-GO: battery` (under 6.6 V) that applies.
- **Check:** All good → `GO`; each single failure → its reason; down and windy →
  `NO-GO: station down`.

---

## Chapter 6 · Counting Down to AOS

`counting-down-to-aos` · 8 lessons: 4 code, 3 code_tests, 1 quiz

Loops: the three-part `for`, `for` as a while loop, `range` over an integer, `break`,
`continue`, nested loops, and keeping the best value seen.

**Story arc.** Launch week. Pipit-3 rides to orbit with forty other satellites, and nobody
knows yet which dot is theirs. Ops counts down to its first predicted pass, steps through
the prediction minute by minute and sweeps the receiver across Doppler offsets, hunting for
a faint beacon. You write each loop. The chapter ends when Pipit-3 answers.

Lessons 6.3 to 6.8 use given functions for their data, so they need no slices yet:
`elevationAt(minute int) float64` follows 64 − 1.6·(m − 7)² (−14.4, 6.4, 24.0, 38.4, 49.6,
57.6, 62.4, 64.0, 62.4, … for minutes 0 to 14).

#### 6.1 The Countdown
`the-countdown` · code · difficulty 2
- **Objective:** Repeat a statement with a three-part `for` loop.
- **Scenario:** Ops wants a countdown on the console before Pipit-3's first pass.
- **Exercise:** The starter prints `T-5` once, then `AOS`; you loop from 5 down to 1.
- **Check:** `T-5` / `T-4` / `T-3` / `T-2` / `T-1` / `AOS`.

#### 6.2 Counting in Millivolts
`counting-in-millivolts` · code · difficulty 3
- **Objective:** Loop while a condition holds, with a condition-only `for`.
- **Scenario:** Pipit-3 charges 150 mV per orbit from 6700 mV; Sunniva asks when it reaches
  8200.
- **Exercise:** The starter prints 0 orbits; you loop `for mv < 8200`, counting orbits. A
  tip says why the count uses whole millivolts, not volts.
- **Check:** `Orbits to full: 10` / `Battery: 8200 mV`.

#### 6.3 Ten Minutes of Predictions
`ten-minutes-of-predictions` · code · difficulty 3
- **Objective:** Loop a fixed number of times with `for i := range n`.
- **Scenario:** The first prediction lists Pipit-3's elevation for each minute of the window.
- **Exercise:** `elevationAt` is given; the starter prints minute 0 only; you print minutes
  0 to 9.
- **Check:** Ten lines from `minute 0: -14.4` to `minute 9: 57.6` (`%.1f`).

#### 6.4 Quick Check: Loops
`quick-check-loops` · quiz · difficulty 2
- **Objective:** Check how many times loops run.
- **Scenario:** Sunniva's two questions over the slow backhaul.
- **Exercise:** Two questions.
- **Check:** `for i := 0; i < 3; i++` runs three times; `for { }` without `break` runs
  forever.

#### 6.5 Finding AOS
`finding-aos` · code_tests · difficulty 4
- **Objective:** Leave a loop early with `break`.
- **Scenario:** AOS is the first minute Pipit-3 clears the mask; no need to look further.
- **Exercise:** `aosMinute(maskDeg float64) int` returns 0; you scan minutes 0 to 14 for the
  first `elevationAt(m) >= maskDeg`, returning −1 if there is none.
- **Check:** mask 5 → 1; mask 30 → 3; mask 64 → 7; mask 70 → −1.

#### 6.6 Skipping the Gaps
`skipping-the-gaps` · code_tests · difficulty 4
- **Objective:** Skip an iteration with `continue`.
- **Scenario:** Pipit-3's first battery readings have gaps, marked −1, where frames were lost.
- **Exercise:** `batteryAt(minute int) float64` is given (7.92, 7.90, −1, 7.86, −1, 7.80 for
  minutes 0 to 5); `averageBattery(from, to int) float64` returns 0; you average the minutes
  from `from` up to but not including `to`, skipping gaps, and return 0 when none is left.
- **Check:** `(0, 2)` ≈ 7.91; `(0, 6)` ≈ 7.87; `(2, 3) == 0`; `(4, 5) == 0`.

#### 6.7 Every Station, Every Block
`every-station-every-block` · code · difficulty 5
- **Objective:** Put one loop inside another.
- **Scenario:** The planner wants a grid of which stations see Pipit-3 in each two-hour
  block.
- **Exercise:** `stationCode(i int) string` and `visible(station, block int) bool` are
  given; you print one row per station, `#` when visible and `.` when not, for 12 blocks.
- **Check:** Five rows of the form `YST ##.###.##.##`, exactly as `visible` dictates.

#### 6.8 The Doppler Sweep
`the-doppler-sweep` · code_tests · difficulty 5
- **Objective:** Keep the best value seen while looping.
- **Scenario:** Pipit-3's beacon is somewhere within ±12 kHz of its frequency; the receiver
  steps through offsets and keeps the strongest. It answers at +3500 Hz.
- **Exercise:** `signalStrength(offsetHz int) float64` is given (it peaks at +3500 Hz and
  falls off evenly); `bestOffset(fromHz, toHz, stepHz int) int` returns 0; you scan and
  return the strongest offset, keeping the first one on a tie.
- **Check:** `(-12000, 12000, 500)` → 3500; `(-12000, 0, 500)` → 0;
  `(-12000, 12000, 1000)` → 3000 (3000 and 4000 tie); `(5000, 9000, 1000)` → 5000.

---

## Chapter 7 · Modeling the Network

`modeling-the-network` · 7 lessons: 3 code, 3 code_tests, 1 quiz

Structs: defining types with fields, struct literals, structs as parameters, nested
structs, `time.Time` fields, copies, and embedding.

**Story arc.** The network's facts live in a spreadsheet nobody trusts: Lantern Atoll appears
twice and Cape Wren's mask is wrong. Nkem asks you to model the network in code instead. You
define `Station`, `Satellite` with its `Orbit`, and `Pass`: the types every later chapter
builds on.

#### 7.1 A Station Type
`a-station-type` · code · difficulty 3
- **Objective:** Define a struct type and create a value with named fields.
- **Scenario:** Port Ellery's facts become one value instead of four loose variables.
- **Exercise:** The starter prints four loose variables; you define `Station` with `Code`,
  `Name`, `LatDeg` and `MaskDeg`, build Port Ellery and print it with `%+v`.
- **Check:** `{Code:PEL Name:Port Ellery LatDeg:46.9 MaskDeg:5}`.

#### 7.2 Structs as Parameters
`structs-as-parameters` · code_tests · difficulty 3
- **Objective:** Pass a struct to a function and read its fields.
- **Scenario:** The Board prints each station as one tidy line.
- **Exercise:** `stationLine(s Station) string` returns `""`; you format code, name,
  latitude with N or S, and mask.
- **Check:** Port Ellery → `PEL Port Ellery (46.9 N, mask 5.0)`; Salar Quieto →
  `SQT Salar Quieto (23.4 S, mask 3.0)`; Lantern Atoll → `LAN Lantern Atoll (1.8 S, mask 8.0)`.

#### 7.3 An Orbit Inside a Satellite
`an-orbit-inside` · code_tests · difficulty 4
- **Objective:** Nest one struct inside another and reach through with dots.
- **Scenario:** A satellite's orbit is a little record of its own.
- **Exercise:** `Orbit` and `Satellite` are given; `isPolar(s Satellite) bool` returns
  `false`; you return whether `s.Orbit.IncDeg` is from 80 to 100.
- **Check:** Pipit-1 (97.5°) true; Brine-3 (53°) false; Thistle (51.6°) false; 80 and 100
  true.

#### 7.4 Copies, Not References
`copies-not-references` · code · difficulty 4
- **Objective:** See that assigning a struct copies it.
- **Scenario:** You try a 12° mask for Cape Wren on a draft, without touching the live record.
- **Exercise:** The starter copies `live` into `draft` and prints both masks; you set
  `draft.MaskDeg = 12`.
- **Check:** `live mask: 5` / `draft mask: 12`.

#### 7.5 Quick Check: Structs
`quick-check-structs` · quiz · difficulty 3
- **Objective:** Check struct zero values, copies and literals.
- **Scenario:** Nkem's check before your types go into review.
- **Exercise:** Three questions.
- **Check:** A struct's zero value has every field at its zero value; after `b := a` and
  `b.MaskDeg = 9`, `a.MaskDeg` is unchanged; `Station{Code: "PEL"}` sets only `Code`.

#### 7.6 The Pass Type
`the-pass-type` · code_tests · difficulty 4
- **Objective:** Keep `time.Time` values in a struct and subtract them.
- **Scenario:** A pass is a satellite, a station, an AOS and an LOS; the rest follows.
- **Exercise:** `Pass` is given; `duration(p Pass) time.Duration` returns 0; you return
  `p.LOS.Sub(p.AOS)`.
- **Check:** 04:12:00 to 04:21:12 → 9m12s; AOS equal to LOS → 0; 23:58 to 00:07 the next
  day → 9m.

#### 7.7 Shared Equipment
`shared-equipment` · code · difficulty 5
- **Objective:** Embed one struct in another to reuse its fields.
- **Scenario:** Stations and the mobile test van share the same `Dish` record.
- **Exercise:** `Station` holds `Equipment Dish`, so `yst.DiameterM` doesn't compile; you
  embed `Dish` in `Station` and in `Van`.
- **Check:** `YST dish 3.7 m, X-band true` / `van dish 1.2 m, X-band false`. The starter
  fails with `yst.DiameterM undefined`.

---

## Chapter 8 · The Day's Schedule

`the-days-schedule` · 10 lessons: 4 code, 4 code_tests, 2 quizzes

Slices: literals, indexing and `len`, `range`, `append`, slicing, shared backing arrays and
`slices.Clone`, sorting with a comparison function, and variadic functions.

**Story arc.** Three Pipits, two Furrows and the other customers add up to dozens of passes a
day. The schedule becomes a slice: you list it, loop over it, grow it, cut it and sort it.
Along the way, editing a draft moves a Furrow pass on the live Board without anyone touching
it, and in the last lesson you catch two passes booked on the same dish at Ystmark.

#### 8.1 A List of Stations
`a-list-of-stations` · code · difficulty 3
- **Objective:** Create a slice literal, index it and take its length.
- **Scenario:** The five station codes belong in one list, not five variables.
- **Exercise:** The starter has five variables and prints placeholders; you build
  `codes := []string{...}` and print the first, the last and the count.
- **Check:** `First: YST` / `Last: CWR` / `Count: 5`.

#### 8.2 Every Pass in the Log
`every-pass-in-the-log` · code · difficulty 3
- **Objective:** Loop over a slice with `for i, v := range`.
- **Scenario:** The pass log numbers each pass from 1.
- **Exercise:** The starter prints only the first pass ID; you range over the slice.
- **Check:** `1 pipit-1/PEL/04182` / `2 furrow-a/SQT/11873` / `3 thistle/LAN/02218`.

#### 8.3 A Growing Schedule
`a-growing-schedule` · code · difficulty 3
- **Objective:** Grow a slice with `append` and keep its result.
- **Scenario:** Predictions arrive one by one, and the schedule grows.
- **Exercise:** The starter's loop receives three pass IDs but adds none; you write
  `passes = append(passes, id)`.
- **Check:** `passes: 3` / `last: thistle/LAN/02218`.

#### 8.4 Quick Check: Slices
`quick-check-slices` · quiz · difficulty 3
- **Objective:** Check length, indexing and append.
- **Scenario:** Nkem's three questions before you touch the scheduler.
- **Exercise:** Three questions.
- **Check:** `len([]int{4, 8})` is 2; `s[5]` on a slice of length 3 panics at run time;
  `s = append(s, x)` needs the assignment because `append` returns the grown slice.

#### 8.5 Cutting the Timeline
`cutting-the-timeline` · code_tests · difficulty 4
- **Objective:** Take part of a slice with `s[a:b]`.
- **Scenario:** The Board shows the next three passes; the alert panel the last two readings.
- **Exercise:** `nextN(ids []string, n int) []string` and `lastN(readings []float64, n int)
  []float64` return `nil`; you slice, clamping `n` to the length.
- **Check:** `nextN` of five IDs with 3 → the first three, with 9 → all five; `lastN` of
  `[7.9 7.8 7.7]` with 2 → `[7.8 7.7]`, with 0 → empty.

#### 8.6 The Pass That Moved Itself
`the-pass-that-moved` · code · difficulty 5
- **Objective:** See that slices can share one array, and copy with `slices.Clone`.
- **Scenario:** Editing a draft of the schedule moved a Furrow pass on the live Board.
- **Exercise:** The starter makes `draft := live[:2]` and edits `draft[1]`, which changes
  `live` too; you build the draft with `slices.Clone(live[:2])`.
- **Check:** `live: [04:12 05:03 06:47]` / `draft: [04:12 05:30]`.

#### 8.7 Quick Check: Sharing
`quick-check-sharing` · quiz · difficulty 4
- **Objective:** Predict which changes show through a shared array.
- **Scenario:** Ingrid asks how the Board's pass moved.
- **Exercise:** Two questions about a four-line snippet.
- **Check:** After `b := a[:2]` and `b[0] = "X"`, `a[0]` is `"X"`; after
  `c := slices.Clone(a)` and `c[0] = "Y"`, `a[0]` is unchanged.

#### 8.8 Sorting the Day
`sorting-the-day` · code_tests · difficulty 5
- **Objective:** Sort a slice of structs with `slices.SortFunc` and a comparison function.
- **Scenario:** Predictions arrive in any order; the Board needs them by AOS.
- **Exercise:** `sortByAOS(passes []Pass)` is empty; you call `slices.SortFunc` with a
  function literal comparing `a.AOS.Compare(b.AOS)`, breaking ties by station code.
- **Check:** Four passes come out in AOS order; two with the same AOS come out `LAN` before
  `YST`.

#### 8.9 Any Number of Readings
`any-number-of-readings` · code_tests · difficulty 4
- **Objective:** Write a variadic function and pass a slice to it with `...`.
- **Scenario:** The temperature panel averages three readings one day, three hundred the next.
- **Exercise:** `average(readings ...float64) float64` returns 0; you sum and divide,
  returning 0 for no readings.
- **Check:** `average() == 0`; `average(7.9) == 7.9`; `average(1, 2, 3, 4) == 2.5`;
  `average(temps...)` with three values.

#### 8.10 Double-Booked at Ystmark
`double-booked` · code_tests · difficulty 6
- **Objective:** Compare every pair in a slice with nested index loops.
- **Scenario:** Two passes tonight need Ystmark's dish at the same time.
- **Exercise:** `overlaps(passes []Pass) []string` returns `nil`; for each pair at the same
  station whose times overlap, you add `"<sat>@<station> overlaps <sat>@<station>"`, in
  input order.
- **Check:** No overlap → empty; Furrow-B and Pipit-3 at YST overlapping → one message;
  touching passes (LOS equal to the next AOS) don't overlap; the same times at different
  stations don't overlap.

---

## Chapter 9 · Who's Who in the Fleet

`whos-who-in-the-fleet` · 8 lessons: 3 code, 4 code_tests, 1 quiz

Maps: literals and lookups, the comma-ok form, counting, deleting, iterating in a stable
order, maps as sets, and maps keyed by chunk number.

**Story arc.** Yusra forwards an angry email: Furrowlight's dashboard showed a pass for
Brine-2, which never left its deployer. The fleet list is a slice someone searches by hand.
You rebuild it as maps: satellites by ID, frame counts per satellite, the set of sequence
numbers already seen. The chapter ends with the chunk map that tells Furrowlight exactly
which pieces of an image never came down.

#### 9.1 The Fleet Registry
`the-fleet-registry` · code · difficulty 3
- **Objective:** Create a map literal, add a key and look values up.
- **Scenario:** Ops wants each satellite's operator by ID.
- **Exercise:** The starter's map lacks Thistle; you add it and print three lookups.
- **Check:** `pipit-1: Apogee` / `furrow-b: Furrowlight` / `thistle: Kessick Polytechnic`.

#### 9.2 Is It in the Registry?
`is-it-in-the-registry` · code_tests · difficulty 4
- **Objective:** Tell a missing key from a zero value with the comma-ok form.
- **Scenario:** Asking for Brine-2 must say "unknown satellite", not print an empty operator.
- **Exercise:** `operatorOf(reg map[string]string, id string) string` returns `reg[id]`; you
  return `"unknown satellite"` when `id` is missing.
- **Check:** `pipit-1` → `Apogee`; `brine-2` → `unknown satellite`; `pipit-3`, present with
  an empty operator → `""`.

#### 9.3 Counting Frames
`counting-frames` · code · difficulty 4
- **Objective:** Count with `m[key]++`, starting from a missing key's zero value.
- **Scenario:** The hourly report shows how many frames each satellite sent.
- **Exercise:** `countBySat(ids []string) map[string]int` returns an empty map; you count.
  `main` prints the map, which `fmt` shows in key order.
- **Check:** `map[brine-3:1 pipit-1:3 thistle:2]`.

#### 9.4 Quick Check: Maps
`quick-check-maps` · quiz · difficulty 3
- **Objective:** Check how maps behave at their edges.
- **Scenario:** Nkem's quick check before the registry ships.
- **Exercise:** Three questions.
- **Check:** Reading a missing key gives the value type's zero value; `range` over a map has
  no fixed order; writing to a nil map panics.

#### 9.5 Removing Brine-2
`removing-brine-2` · code · difficulty 3
- **Objective:** Remove a key with `delete`.
- **Scenario:** Brine-2 never left its deployer; it must leave the registry.
- **Exercise:** The starter's registry of nine still lists `brine-2`; you delete it.
- **Check:** `satellites: 8` / `brine-2 listed: false`.

#### 9.6 A Report in Stable Order
`a-stable-report` · code_tests · difficulty 5
- **Objective:** Visit a map's keys in sorted order.
- **Scenario:** The hourly report jumps around, because map order changes from run to run.
- **Exercise:** `report(counts map[string]int) []string` returns `nil`; you return
  `"id: n"` lines sorted by ID with `slices.Sorted(maps.Keys(counts))`.
- **Check:** `{thistle: 2, pipit-1: 3, brine-3: 1}` → `["brine-3: 1", "pipit-1: 3",
  "thistle: 2"]`; an empty map → empty.

#### 9.7 Seen Before
`seen-before` · code_tests · difficulty 5
- **Objective:** Use a map as a set.
- **Scenario:** Thistle's frames sometimes arrive twice, once through each of two stations.
- **Exercise:** `dedupe(seqs []int) []int` returns `seqs`; you keep each number's first
  appearance, in order.
- **Check:** `[4182 4183 4182 4184 4183]` → `[4182 4183 4184]`; empty → empty; no
  duplicates → unchanged.

#### 9.8 The Missing Chunks
`the-missing-chunks` · code_tests · difficulty 6
- **Objective:** Combine a map with a counting loop to find what's absent.
- **Scenario:** A Furrow-A image came down in numbered chunks; Furrowlight wants the missing
  numbers before nine.
- **Exercise:** `missing(chunks map[int]string, total int) []int` returns `nil`; you return
  the numbers from 0 to `total-1` that have no chunk, in order.
- **Check:** total 5 holding 0, 1 and 3 → `[2 4]`; all present → empty; none of 3 →
  `[0 1 2]`.

---

## Chapter 10 · Changing Things in Place

`changing-things-in-place` · 9 lessons: 4 code, 3 code_tests, 2 quizzes

Pointers and methods: addresses and dereferencing, `nil`, value and pointer receivers,
named types with `iota`, and the `String` method.

**Story arc.** Sunniva radios from Ystmark: a storm is coming and both dishes must stow. Your
stow function runs, and Ystmark still shows as tracking, because the function changed a
copy. Pointers explain why. Then Nkem shows you how to give types behavior of their own: a
Pass that knows its length, a Station that stows itself, a severity that prints as a word.
The chapter ends with the storm itself, hour by hour.

#### 10.1 The Stow That Didn't
`the-stow-that-didnt` · code · difficulty 4
- **Objective:** Pass a pointer so a function can change the caller's value.
- **Scenario:** `stow(yst)` ran before the storm, but Ystmark still shows as tracking.
- **Exercise:** `stow(s Station)` sets `s.Stowed`; you make it take `*Station` and call
  `stow(&yst)`.
- **Check:** `YST stowed: true`.

#### 10.2 Quick Check: Pointers
`quick-check-pointers` · quiz · difficulty 4
- **Objective:** Check `&`, `*` and `nil`.
- **Scenario:** Sunniva's three questions, one word each in her reply.
- **Exercise:** Three questions about a four-line snippet.
- **Check:** `&x` is the address of `x`; `*p = 5` changes the variable `p` points to; a
  pointer's zero value is `nil`.

#### 10.3 Nobody Home
`nobody-home` · code_tests · difficulty 5
- **Objective:** Return a pointer, or `nil` when there is nothing to point to.
- **Scenario:** Looking up `XYZ` must say there's no such station, not invent an empty one.
- **Exercise:** `find(stations []Station, code string) *Station` returns `nil`; you return
  `&stations[i]` for the match.
- **Check:** `find(..., "SQT")` points into the slice, so setting `Stowed` through it shows
  in the slice; `find(..., "XYZ")` is `nil`.

#### 10.4 A Pass Knows Its Length
`a-pass-knows-its-length` · code · difficulty 4
- **Objective:** Declare a method with a value receiver.
- **Scenario:** Everyone writes `p.LOS.Sub(p.AOS)`; Nkem wants `p.Duration()`.
- **Exercise:** `main` calls `p.Duration()` on three passes, but the method doesn't exist
  yet; you add `func (p Pass) Duration() time.Duration`.
- **Check:** `pipit-1@PEL 9m12s` / `furrow-a@SQT 11m4s` / `thistle@LAN 6m30s`. The starter
  fails with `p.Duration undefined`.

#### 10.5 Stow Yourself
`stow-yourself` · code_tests · difficulty 5
- **Objective:** Use a pointer receiver for a method that changes its value.
- **Scenario:** Sunniva wants `st.Stow("wind 81 km/h")` to record why a dish is parked.
- **Exercise:** `Stow(reason string)` and `Unstow()` have value receivers, so their changes
  vanish; you switch both to pointer receivers.
- **Check:** After `Stow("wind 81 km/h")`: `Stowed` true, `StowReason` set; after `Unstow()`:
  false and `""`.

#### 10.6 Quick Check: Receivers
`quick-check-receivers` · quiz · difficulty 4
- **Objective:** Check when a method needs a pointer receiver.
- **Scenario:** Nkem's review comment comes as two questions.
- **Exercise:** Two questions.
- **Check:** A method that changes its struct needs a pointer receiver; `st.Stow(...)` on a
  variable works without `&`, because Go takes the address for you.

#### 10.7 Severity Levels
`severity-levels` · code_tests · difficulty 5
- **Objective:** Declare a named type with a block of `iota` constants.
- **Scenario:** The pager ranks alerts Info, Warning and Critical, and must pick the worst.
- **Exercise:** `type Severity int` has three constants all set to 0, and
  `worst(a, b Severity) Severity` returns `a`; you number them with `iota` and write `worst`.
- **Check:** `Info < Warning < Critical`; `worst(Info, Critical) == Critical`;
  `worst(Warning, Info) == Warning`.

#### 10.8 A Word, Not a Number
`a-word-not-a-number` · code · difficulty 5
- **Objective:** Add a `String` method so `fmt` prints a type as text.
- **Scenario:** Pages say "severity 2", and the on-call engineer has to look it up.
- **Exercise:** `main` prints three severities as `0`, `1`, `2`; you add
  `func (s Severity) String() string`.
- **Check:** `INFO` / `WARNING` / `CRITICAL`.

#### 10.9 Storm at Ystmark
`storm-at-ystmark` · code · difficulty 6
- **Objective:** Call pointer methods on the elements of a slice, not on copies.
- **Scenario:** Hour by hour the wind rises and falls at Ystmark, and your code stows and
  unstows both dishes.
- **Exercise:** The starter loops with `for _, d := range dishes` and stows copies; you loop
  by index so the dishes in the slice change.
- **Check:**
  ```text
  02:00 wind 64 km/h: tracking
  03:00 wind 75 km/h: stowed
  04:00 wind 58 km/h: tracking
  05:00 wind 81 km/h: stowed
  dish 1 stowed 2 times
  dish 2 stowed 2 times
  ```
  With the starter's copies, every hour reads `tracking` and both counts are 0.

---

## Chapter 11 · Every Radio Speaks

`every-radio-speaks` · 8 lessons: 3 code, 4 code_tests, 1 quiz

Interfaces: declaring one, satisfying it implicitly, `fmt.Stringer` and `io.Writer`, type
switches, type assertions, and slices of interface values.

**Story arc.** Apogee's stations run three kinds of radio: UHF for small cubesats, S-band for
housekeeping and commands, X-band for the big imagery dumps. Every new radio has meant
another copy of the pass code. Tomás wants one pass routine for any radio, and Mirela wants
anomaly rules she can add without touching the pager. Interfaces give them both.

#### 11.1 One Pass Routine
`one-pass-routine` · code_tests · difficulty 5
- **Objective:** Declare an interface and accept it as a parameter.
- **Scenario:** The downlink-time estimate exists three times, once per radio.
- **Exercise:** `UHFRadio`, `SBandRadio` and `XBandRadio` each have `RateMbps()`; `Radio` is
  an empty interface and `downlinkSeconds(r Radio, megabytes float64) float64` returns 0;
  you add the method to `Radio` and compute megabytes × 8 / rate.
- **Check:** X-band (150 Mbps), 9000 MB → 480 s; S-band (4 Mbps), 30 MB → 60 s; UHF
  (0.0096 Mbps), 0.0072 MB → 6 s.

#### 11.2 No implements Keyword
`no-implements-keyword` · code · difficulty 5
- **Objective:** Satisfy an interface just by having its methods.
- **Scenario:** Salar Quieto's new Ka-band test radio must plug into the routine unchanged.
- **Exercise:** `main` passes a `KaBandRadio` to `downlinkSeconds`, but it has no `RateMbps`
  method; you add one returning 300.
- **Check:** `Ka-band: 9000 MB in 240 s`. The starter fails with
  `KaBandRadio does not implement Radio (missing method RateMbps)`.

#### 11.3 Quick Check: Interfaces
`quick-check-interfaces` · quiz · difficulty 4
- **Objective:** Check what satisfying an interface means.
- **Scenario:** Tomás asks three questions while the Ka-band radio warms up.
- **Exercise:** Three questions.
- **Check:** A type satisfies an interface by having its methods, with no declaration; one
  interface variable can hold different types over time; every type satisfies `any`.

#### 11.4 Satellites That Introduce Themselves
`satellites-introduce-themselves` · code · difficulty 4
- **Objective:** Implement `fmt.Stringer`, the interface `fmt` looks for.
- **Scenario:** Logs print satellites as `{pipit-1 Pipit-1 {525 97.5}}`; ops wants
  `Pipit-1 (525 km)`.
- **Exercise:** `main` prints three satellites with `Println`; you add
  `func (s Satellite) String() string`.
- **Check:** `Pipit-1 (525 km)` / `Thistle (415 km)` / `Cirrine-1 (600 km)`.

#### 11.5 Write Anywhere
`write-anywhere` · code_tests · difficulty 5
- **Objective:** Accept an `io.Writer` so output can go anywhere.
- **Scenario:** The Board must draw on the Fishbowl's screen and, in tests, into a string.
- **Exercise:** `writeBoard(w io.Writer, ids []string)` ignores `w` and prints with
  `fmt.Println`; you write with `fmt.Fprintln(w, ...)`.
- **Check:** Given a `strings.Builder`, the result is exactly `NEXT PASSES\n` followed by one
  ID per line; with no IDs, just the header.

#### 11.6 Which Payload Is This?
`which-payload` · code · difficulty 5
- **Objective:** Branch on a value's concrete type with a type switch.
- **Scenario:** Downlinked items are images, housekeeping frames or file chunks, and each
  goes somewhere different.
- **Exercise:** `route(item any) string` returns `"quarantine"`; you write a type switch.
  `main` prints `%T -> route` for four items.
- **Check:** `main.Image -> imagery` / `main.Housekeeping -> telemetry` /
  `main.Chunk -> reassembly` / `string -> quarantine`.

#### 11.7 Can It Track?
`can-it-track` · code_tests · difficulty 6
- **Objective:** Check for an optional method with a type assertion's comma-ok form.
- **Scenario:** Some radios correct Doppler on their own; the routine must pre-compensate
  the rest.
- **Exercise:** `DopplerTracker` (one method, `TracksDoppler()`) is given;
  `needsPrecompensation(r Radio) bool` returns `true`; you return `false` when `r` is also a
  `DopplerTracker`.
- **Check:** S-band (tracks) → false; UHF (doesn't) → true; X-band (tracks) → false.

#### 11.8 Anomaly Rules
`anomaly-rules` · code_tests · difficulty 6
- **Objective:** Keep different types in one slice through an interface.
- **Scenario:** Mirela wants to add anomaly rules without editing the pager.
- **Exercise:** `Rule` (`Check(f TelemetryFrame) (string, bool)`) and
  `evaluate(rules []Rule, f TelemetryFrame) []string` are given; `BatteryLow{MinV}` and
  `TooHot{MaxC}` return `"", false`; you write their `Check` methods.
- **Check:** 6.4 V, 20 °C → `["battery 6.40 V below 6.60 V"]`; 7.9 V, 51 °C →
  `["temperature 51.0 C above 45.0 C"]`; both → both messages in rule order; healthy →
  none.

---

## Chapter 12 · When the Frame Is Wrong

`when-the-frame-is-wrong` · 10 lessons: 2 code, 6 code_tests, 2 quizzes

Errors: returning and checking `error` values, messages with `fmt.Errorf`, sentinel errors
and `errors.Is`, wrapping with `%w`, custom error types with `errors.AsType`, `defer` for
cleanup, and `errors.Join`.

**Story arc.** Thistle starts sending garbage: temperatures of 900 °C and frames that arrive
twice. The telemetry service drops the bad ones without a word, and nobody notices for a
week. Ingrid's new rule: every failure gets a name, a reason and a line in the log. You
learn Go's error values one step at a time, and finish with a validator that says exactly
what is wrong with a frame.

#### 12.1 Errors Are Values
`errors-are-values` · code_tests · difficulty 4
- **Objective:** Return an error, or `nil` when all is well.
- **Scenario:** A battery reading under 6.6 V must be reported, not stored quietly.
- **Exercise:** `checkBattery(v float64) error` returns `nil`; you return
  `errors.New("battery too low")` under 6.6 V.
- **Check:** 7.4 → nil; 6.6 → nil; 6.2 → an error reading `battery too low`.

#### 12.2 Check It Right Away
`check-it-right-away` · code · difficulty 4
- **Objective:** Handle an error with `if err != nil` right after the call.
- **Scenario:** The ingest loop stores every reading, even the rejected ones.
- **Exercise:** The starter ignores `checkBattery`'s error; you print the rejection and
  skip the store.
- **Check:** `7.91: stored` / `6.20: rejected: battery too low` / `7.40: stored` /
  `stored 2 of 3`.

#### 12.3 Say What Went Wrong
`say-what-went-wrong` · code_tests · difficulty 5
- **Objective:** Put values into an error message with `fmt.Errorf`.
- **Scenario:** "battery too low" doesn't say which satellite, or how low.
- **Exercise:** `checkBattery(sat string, v float64) error` returns the old fixed message;
  you include the satellite and both voltages.
- **Check:** `("thistle", 6.2)` → `thistle: battery 6.20 V below 6.60 V`;
  `("pipit-1", 7.9)` → nil.

#### 12.4 Quick Check: Errors
`quick-check-errors` · quiz · difficulty 4
- **Objective:** Check what an error is and how messages are written.
- **Scenario:** Ingrid's quick check before the new log goes live.
- **Exercise:** Three questions.
- **Check:** `error` is an interface with one method, `Error() string`; a `nil` error means
  success; Go style is `battery too low`: lowercase, no final period.

#### 12.5 Known Failures
`known-failures` · code_tests · difficulty 5
- **Objective:** Declare a sentinel error and recognize it with `errors.Is`.
- **Scenario:** Duplicate frames are normal; the pager must ignore them and page on anything
  else.
- **Exercise:** `ErrDuplicate` is declared; `accept(seq, last uint32) error` returns a new
  error each time, and `shouldPage(err error) bool` returns `err != nil`; you return
  `ErrDuplicate` itself and skip it in `shouldPage`.
- **Check:** `accept(4182, 4182)` is `ErrDuplicate`; `accept(4183, 4182)` is nil;
  `shouldPage(ErrDuplicate)` false; `shouldPage(errors.New("bad crc"))` true;
  `shouldPage(nil)` false.

#### 12.6 Wrap It, Keep It
`wrap-it-keep-it` · code_tests · difficulty 5
- **Objective:** Add context with `%w` and keep the original error recognizable.
- **Scenario:** The log wants `thistle seq 4182: duplicate frame`, and the pager must still
  spot the duplicate.
- **Exercise:** `ingest(sat string, seq, last uint32) error` formats with `%v`; you switch
  to `%w`.
- **Check:** The message is `thistle seq 4182: duplicate frame`;
  `errors.Is(err, ErrDuplicate)` is true.

#### 12.7 An Error with Fields
`an-error-with-fields` · code_tests · difficulty 6
- **Objective:** Define an error type with fields and extract it with `errors.AsType`.
- **Scenario:** The pager wants the field and the value that broke the limit, not a sentence.
- **Exercise:** `LimitError{Field string; Value, Min, Max float64}` has its `Error` method;
  `checkTemp(c float64) error` returns `errors.New(...)` and `pagerDetail(err error) string`
  returns `""`; you return a `*LimitError` and read it back with
  `errors.AsType[*LimitError](err)`. A note mentions `errors.As`, the older form.
- **Check:** `checkTemp(51)` gives `Field "temp_c"`, `Value 51`, `Min -10`, `Max 45`;
  `pagerDetail` of it wrapped with `%w` → `temp_c=51.0 outside [-10.0, 45.0]`;
  `pagerDetail` of another error → `""`.

#### 12.8 Quick Check: Is or As?
`quick-check-is-or-as` · quiz · difficulty 5
- **Objective:** Choose between `errors.Is` and `errors.AsType`.
- **Scenario:** Mirela asks which tool the pager should use, and when.
- **Exercise:** Two questions.
- **Check:** To recognize `ErrDuplicate` inside a wrapped error, `errors.Is`; to get a
  `*LimitError`'s fields, `errors.AsType`.

#### 12.9 Close the Session, Always
`close-the-session` · code · difficulty 5
- **Objective:** Use `defer` to run cleanup however a function returns.
- **Scenario:** A pass session left open after an error blocked the dish for the next pass.
- **Exercise:** `runPass` closes its session at the end, but an early return skips it; you
  `defer` the close right after opening.
- **Check:** `open YST` / `pass ok` / `close YST` / `open CWR` /
  `pass failed: wind 81 km/h` / `close CWR`.

#### 12.10 The Frame Validator
`the-frame-validator` · code_tests · difficulty 7
- **Objective:** Return every problem at once with `errors.Join`.
- **Scenario:** One of Thistle's frames is both a duplicate and too hot, and the log should
  say both.
- **Exercise:** `validate(f TelemetryFrame, lastSeq uint32) error` returns the first problem
  only; you collect them all and join them.
- **Check:** Healthy → nil; duplicate → `errors.Is(err, ErrDuplicate)`; 900 °C →
  `errors.AsType[*LimitError]` with `Field "temp_c"`; duplicate and too hot → both hold and
  the message has two lines.

---

## Chapter 13 · Reading the Raw Feed

`reading-the-raw-feed` · 8 lessons: 2 code, 5 code_tests, 1 quiz

Strings, bytes and runes: splitting and cutting lines, parsing numbers with `strconv`,
UTF-8 and rune counts, iterating over runes, hex and big-endian bytes, and HMAC signatures.

**Story arc.** Thistle's frames come down as text lines, Pipit beacons as raw bytes, and
place names arrive from Salar Quieto's console in UTF-8 that the Board mangles. You take
them all apart. Then you build the most careful bytes Apogee sends: a signed command for
Pipit-2, which Mirela reviews line by line.

#### 13.1 Taking a Line Apart
`taking-a-line-apart` · code · difficulty 4
- **Objective:** Split a line into fields with `strings.Fields`.
- **Scenario:** A telemetry line is satellite, sequence and time, then `key=value` pairs, and
  this one has a doubled space.
- **Exercise:** The starter splits with `strings.Split(line, " ")`, which yields an empty
  field; you use `strings.Fields`.
- **Check:** `fields: 6` / `satellite: thistle` / `seq: 4182`.

#### 13.2 Key Equals Value
`key-equals-value` · code_tests · difficulty 5
- **Objective:** Split a string once with `strings.Cut`.
- **Scenario:** `batt_v=7.91` has to become a key and a value.
- **Exercise:** `field(line, key string) (string, bool)` returns `"", false`; you cut each
  field at the first `=`.
- **Check:** `batt_v` → `"7.91", true`; `mode` → `"NOMINAL", true`; a missing key →
  `"", false`; `note=a=b` → `"a=b"`.

#### 13.3 Text to Numbers
`text-to-numbers` · code_tests · difficulty 5
- **Objective:** Parse a number with `strconv.ParseFloat` and handle its error.
- **Scenario:** A frame said `batt_v=7.9l`, with a letter l, and the parser stored 0 V.
- **Exercise:** `battery(line string) (float64, error)` ignores the parse error; you return
  an error that names `batt_v`.
- **Check:** `batt_v=7.91` → 7.91, nil; `batt_v=7.9l` → an error containing `batt_v`; no
  `batt_v` field → an error.

#### 13.4 Bytes Are Not Letters
`bytes-are-not-letters` · code · difficulty 5
- **Objective:** Count characters with `utf8.RuneCountInString`, not `len`.
- **Scenario:** The Board cuts "Quebrada Añil", the road to Salar Quieto, a letter short,
  because it counts bytes.
- **Exercise:** The starter prints `len` of three names; you print their rune counts.
- **Check:** `Ystmark: 7` / `Quebrada Añil: 13` / `Tomás: 5` (the starter prints 7, 14, 6).

#### 13.5 Walking Runes
`walking-runes` · code_tests · difficulty 5
- **Objective:** Loop over a string's runes with `for range`.
- **Scenario:** Radio call signs accept only A to Z, digits and hyphens.
- **Exercise:** `callSign(name string) string` returns `name`; you keep ASCII letters
  (uppercased), digits and hyphens, turn spaces into hyphens and everything else into `_`.
- **Check:** `Salar Quieto` → `SALAR-QUIETO`; `Añil` → `A_IL`; `Pipit-3` → `PIPIT-3`;
  `""` → `""`.

#### 13.6 Quick Check: Text
`quick-check-text` · quiz · difficulty 4
- **Objective:** Check bytes, runes and string immutability.
- **Scenario:** Paz asks three questions in a voice memo from 3,840 m.
- **Exercise:** Three questions.
- **Check:** `len("°C")` is 3; `s[0]` is a byte; `s[0] = 'x'` doesn't compile, because
  strings can't be changed.

#### 13.7 Raw Beacon Bytes
`raw-beacon-bytes` · code_tests · difficulty 6
- **Objective:** Decode hex into bytes and read big-endian numbers with `encoding/binary`.
- **Scenario:** Pipit beacons come down as 12 raw bytes; ops needs the sequence number and
  battery voltage.
- **Exercise:** `decodeBeacon(hexText string) (seq, battMV uint16, err error)` returns
  zeros; you decode with `hex.DecodeString`, check the length and the `"AP"` sync bytes, and
  read bytes 3–4 and 6–7 with `binary.BigEndian.Uint16`.
- **Check:** A valid beacon (format in `story.md`) → 4182 and 7910 mV; bad hex → error;
  8 bytes → error; wrong sync bytes → error.

#### 13.8 Signing the Uplink
`signing-the-uplink` · code_tests · difficulty 7
- **Objective:** Sign and verify bytes with HMAC-SHA256, comparing with `hmac.Equal`.
- **Scenario:** Mirela won't let a command reach Pipit-2 unless it carries a signature only
  Apogee can make.
- **Exercise:** `sign(canonical string, key []byte) string` and
  `verify(canonical, sigHex string, key []byte) bool` return `""` and `false`; you compute
  the hex HMAC and compare in constant time.
- **Check:** `sign("pipit-2|118|SET_MODE|SAFE", key)` equals the HMAC the test computes with
  `crypto/hmac`; `verify` accepts it and rejects a changed command, a changed signature, the
  wrong key and a signature that isn't hex.

---

## Chapter 14 · Trust, but Verify

`trust-but-verify` · 8 lessons: 7 code_tests, 1 quiz

Testing basics: cases as data, choosing cases that catch bugs, `*testing.T` with `Errorf`
and `Fatalf`, subtests with `t.Run`, helpers with `t.Helper`, and fixtures in `testdata`.

**Story arc.** A one-character bug in the replay guard nearly lets an old SET_MODE command
through to Cirrine-1. Mirela's new rule: nothing touches the uplink path without tests its
author wrote. Until now the tests came with each lesson; here you write them, and Mirela
proves they work by planting broken copies of the code, mutants, that your tests must catch.

How these lessons check the learner's tests: `_test.go` files are always read-only, so the
learner writes cases and checking functions in `cases.go` or `checks.go` (plain files that
import `testing`). The read-only `main_test.go` runs them against the real code, which must
pass, and runs each mutant in a re-executed test binary
(`exec.Command(os.Args[0], "-test.run=^TestMutant$")` with `APOGEE_MUTANT=<name>`), which
must fail. Subtest names and helper line numbers are checked the same way, from a
re-executed run's output.

#### 14.1 Cases Are Data
`cases-are-data` · code_tests · difficulty 5
- **Objective:** Write test cases as a table of inputs and expected results.
- **Scenario:** Mirela asks for cases for `canonical()`, which lays a command out for signing.
- **Exercise:** `cases.go` holds `canonicalCases`, an empty `[]canonicalCase` (name,
  command, want); you add at least three.
- **Check:** `TestEnoughCases`: three or more, names unique and non-empty;
  `TestCasesHold`: every case passes against `canonical`.

#### 14.2 Catch the Mutants
`catch-the-mutants` · code_tests · difficulty 6
- **Objective:** Pick cases that would catch real bugs, boundaries first.
- **Scenario:** Mirela plants three broken copies of `accept(seq, last uint64) bool`, the
  replay guard's core (`seq > last`).
- **Exercise:** You fill `acceptCases` in `cases.go`.
- **Check:** `TestCasesHold` against the real `accept`; `TestCasesCatchMutants`: each of
  `seq >= last` (accepts a replay), `seq > last+1` (rejects the next command) and
  `seq != last` (accepts an old command) fails at least one case.

#### 14.3 Your First Check
`your-first-check` · code_tests · difficulty 5
- **Objective:** Report a failure with `t.Errorf` from your own checking code.
- **Scenario:** Now Mirela wants the checking code written by you as well.
- **Exercise:** `checks.go` has `checkBand(t *testing.T, band func(float64) string)`, empty;
  you call `band` on values around each threshold and `t.Errorf` on a wrong result.
- **Check:** `TestRealCodePasses` (with `batteryBand`); `TestChecksCatchMutants` (two
  mutants, each must fail in a re-executed run).

#### 14.4 Quick Check: Errorf or Fatalf
`quick-check-errorf-or-fatalf` · quiz · difficulty 4
- **Objective:** Check when a test stops.
- **Scenario:** Mirela's two questions after your first review.
- **Exercise:** Two questions.
- **Check:** After `t.Errorf` the test keeps running; `t.Fatalf` fits when later checks make
  no sense, such as after a parse error.

#### 14.5 One Subtest per Case
`one-subtest-per-case` · code_tests · difficulty 6
- **Objective:** Run each case as a named subtest with `t.Run`.
- **Scenario:** When one of twenty cases fails, the log should name it.
- **Exercise:** `checkWind` in `checks.go` loops over cases without `t.Run`; you wrap each
  case in `t.Run(c.name, ...)`.
- **Check:** `TestSubtestsNamed` (a re-executed `-test.v` run lists
  `TestWind/calm-30` and the other names); `TestChecksCatchMutants`.

#### 14.6 Helpers That Point at the Caller
`helpers-point-at-the-caller` · code_tests · difficulty 6
- **Objective:** Mark a helper with `t.Helper` so failures point at the line that called it.
- **Scenario:** Every failure in the Doppler tests points at the same line inside the helper.
- **Exercise:** `assertNear(t testing.TB, got, want, tol float64)` in `checks.go` works but
  lacks `t.Helper()`; you add it.
- **Check:** `TestHelperPointsAtCaller`: a re-executed probe that fails on purpose reports
  `main_test.go:<line>`, not `checks.go`; `TestAssertNear` passes within the tolerance.

#### 14.7 Fixtures in testdata
`fixtures-in-testdata` · code_tests · difficulty 6
- **Objective:** Load test inputs from `testdata/`, stopping with `t.Fatalf` if they can't
  be read.
- **Scenario:** Last week's real Thistle frames become fixtures for the parser.
- **Exercise:** `testdata/thistle-0511.tlm` is a read-only starter file;
  `loadFixture(t testing.TB, name string) []string` returns `nil`; you read the file with
  `os.ReadFile`, call `t.Fatalf` on error and return its non-empty lines.
- **Check:** `TestLoadsFixture` gets the file's 12 lines; `TestMissingFixtureIsFatal`: a
  re-executed probe asking for a missing file fails with your message.

#### 14.8 Guard the Uplink
`guard-the-uplink` · code_tests · difficulty 7
- **Objective:** Combine a case table, subtests and a helper into a full test of one
  function.
- **Scenario:** Mirela's sign-off: your checks for `verify` must catch all five planted bugs.
- **Exercise:** `checkVerify(t *testing.T, verify func(cmd, sigHex string, key []byte) bool)`
  is empty; you write cases as subtests: valid, changed command, changed signature, wrong
  key, empty signature.
- **Check:** `TestRealCodePasses`; `TestChecksCatchMutants` with five mutants: always true,
  accepts an empty signature, compares only the first 8 hex digits, ignores the key, ignores
  the command.

---

## Chapter 15 · Splitting groundctl

`splitting-groundctl` · 7 lessons: 4 code, 1 code_tests, 2 quizzes

Packages and modules: `go.mod` and import paths, exported and unexported names, a small
package API, `internal` packages, import cycles, and versions.

**Story arc.** groundctl has grown into one huge main package that three teams edit at once.
Nkem draws boxes on the whiteboard: orbit, telemetry, uplink, alert. You move code into
packages, decide what each one shows the world, and find out why Go refuses to let
telemetry and alert import each other.

Every lesson here ships `go.mod` (`module apogee/groundctl`, listed in `readonly.txt`) and
packages in subfolders, with tests, when there are any, in the root `main_test.go`.

#### 15.1 The Module Line
`the-module-line` · code · difficulty 5
- **Objective:** Import a package from your own module by its path.
- **Scenario:** The orbit math has moved into its own folder, `orbit/`.
- **Exercise:** `main.go` calls `orbit.PeriodMinutes` without importing it; you add
  `import "apogee/groundctl/orbit"`.
- **Check:** `Pipit-1 period: 94.99 min`. The starter fails with `undefined: orbit`.

#### 15.2 Capital Letters Export
`capital-letters-export` · code · difficulty 5
- **Objective:** Export a name by capitalizing it, and keep helpers unexported.
- **Scenario:** `main` can't call `orbit.periodMinutes`: outside its package, a lowercase
  name is invisible.
- **Exercise:** You rename `periodMinutes` to `PeriodMinutes` in `orbit/orbit.go` and in
  `main.go`, leaving the helper `cube` unexported.
- **Check:** `Thistle period: 92.72 min`. The starter fails with
  `undefined: orbit.periodMinutes`.

#### 15.3 Quick Check: Packages
`quick-check-packages` · quiz · difficulty 4
- **Objective:** Check import paths, folders and exported names.
- **Scenario:** Nkem's whiteboard quiz.
- **Exercise:** Three questions.
- **Check:** A package's import path is the module path plus its folder; one folder holds
  one package; other packages can use only names that start with a capital letter.

#### 15.4 A Small API
`a-small-api` · code_tests · difficulty 6
- **Objective:** Export exactly what other packages need: a type and a function.
- **Scenario:** The telemetry parser moves to `telemetry/`, and other teams will call it.
- **Exercise:** `telemetry/telemetry.go` has `frame` and `parse`, unexported; the root test
  calls `telemetry.Parse` and reads `Frame` fields; you export `Frame`, its fields and
  `Parse`, and nothing else.
- **Check:** `TestParse` (a good line gives `SatID`, `Seq`, `BattV`); `TestParseError` (a bad
  line gives an error). The starter fails with `undefined: telemetry.Parse`.

#### 15.5 internal Means Internal
`internal-means-internal` · code · difficulty 6
- **Objective:** Keep a package private to its parent with an `internal` folder.
- **Scenario:** The signing keys must be reachable only through `uplink`, never directly.
- **Exercise:** `main.go` imports `apogee/groundctl/uplink/internal/keys`; you call
  `uplink.KeyID` instead.
- **Check:** `signing with key k-2027-03`. The starter fails with
  `use of internal package apogee/groundctl/uplink/internal/keys not allowed`.

#### 15.6 Breaking the Cycle
`breaking-the-cycle` · code · difficulty 6
- **Objective:** Fix an import cycle by moving a shared type into its own package.
- **Scenario:** `alert` imports `telemetry` for frames, `telemetry` imports `alert` for
  severities, and Go refuses both.
- **Exercise:** `severity/severity.go` holds only its package line; you move `Severity`
  there and point both packages at it.
- **Check:** `thistle 4182: WARNING`. The starter fails with `import cycle not allowed`.

#### 15.7 Quick Check: Modules
`quick-check-modules` · quiz · difficulty 5
- **Objective:** Check what `go.mod` and `go.sum` record.
- **Scenario:** Nkem's last question before groundctl's first tagged release.
- **Exercise:** Three questions.
- **Check:** The `module` line sets the path every package's import path starts with;
  `go.sum` holds checksums of dependencies; a v2 module's import path ends in `/v2`.

---

## Chapter 16 · One Helper, Every Type

`one-helper-every-type` · 7 lessons: 1 code, 5 code_tests, 1 quiz

Generics: type parameters, `cmp.Ordered` and constraints of your own, the `~` operator,
generic types with methods, `comparable`, and the generic helpers in `slices`.

**Story arc.** telemetry, downlink and schedule each keep their own copy of filter, max and
last-N, for passes, frames and chunks. Nkem asks for one of each. Generics let you write a
helper once for every type, and a ring buffer that keeps the last readings of anything.

#### 16.1 One Filter for Everything
`one-filter` · code_tests · difficulty 6
- **Objective:** Write a function with a type parameter.
- **Scenario:** `filterPasses` and `filterFrames` are the same eight lines, twice.
- **Exercise:** `Filter[T any](items []T, keep func(T) bool) []T` returns `nil`; you write it.
- **Check:** passes above 10° → 2 of 3; even numbers of 1 to 6 → `[2 4 6]`; nothing kept →
  empty.

#### 16.2 The Highest Reading
`the-highest-reading` · code_tests · difficulty 6
- **Objective:** Constrain a type parameter with `cmp.Ordered`.
- **Scenario:** Highest elevation, highest temperature, last station code: one helper.
- **Exercise:** `Max[T cmp.Ordered](items []T) (T, bool)` returns the zero value and
  `false`; you return the largest and `true`.
- **Check:** `[3 9 4]` → 9, true; `[61.4 23 44.2]` → 61.4; `["YST" "CWR" "SQT"]` → `"YST"`;
  empty → 0, false.

#### 16.3 Quick Check: Generics
`quick-check-generics` · quiz · difficulty 5
- **Objective:** Check constraints and type inference.
- **Scenario:** Nkem's two questions before generics reach the scheduler.
- **Exercise:** Two questions.
- **Check:** `Max[T any]` can't use `>`, because `any` allows no ordering; `Filter(passes, f)`
  needs no type argument, because Go infers it.

#### 16.4 The Standard Toolbox
`the-standard-toolbox` · code · difficulty 5
- **Objective:** Use the generic helpers the standard library already has.
- **Scenario:** Nkem's rule: before writing a helper, look in `slices` and `maps`.
- **Exercise:** The starter has three hand-written loops, each with a bug; you replace them
  with `slices.Contains`, `slices.Index` and `slices.Max`.
- **Check:** `has SQT: true` / `index of LAN: 3` / `highest mask: 8`.

#### 16.5 Units with ~
`units-with-tilde` · code_tests · difficulty 7
- **Objective:** Write a constraint with `~` so named types fit it.
- **Scenario:** `Clamp` works for `float64` but rejects `Volts` and `Kilometers`.
- **Exercise:** `type Real interface{ float64 }`; the test calls `Clamp` with `Volts`, so it
  doesn't compile; you change the constraint to `~float32 | ~float64`.
- **Check:** `Clamp(Volts(9.1), 6.6, 8.4) == 8.4`; `Clamp(Kilometers(300), 400, 600) == 400`;
  `Clamp(7.5, 6.6, 8.4) == 7.5`. The starter fails with
  `Volts does not satisfy Real (possibly missing ~ for float64 in Real)`.

#### 16.6 A Ring of Readings
`a-ring-of-readings` · code_tests · difficulty 7
- **Objective:** Define a generic type with methods.
- **Scenario:** The alert panel keeps the last N of anything: voltages, frames, pass IDs.
- **Exercise:** `Ring[T any]` with `NewRing`, `Push` and `Items` is stubbed; you implement it
  with a slice and the next write position.
- **Check:** capacity 3, push 1 to 5 → `Items()` is `[3 4 5]`; push 1, 2 → `[1 2]`; works for
  strings.

#### 16.7 A Set of Anything
`a-set-of-anything` · code_tests · difficulty 6
- **Objective:** Build a generic type on a map with the `comparable` constraint.
- **Scenario:** Which stations heard both Thistle and Pipit-1 tonight?
- **Exercise:** `type Set[T comparable] map[T]struct{}` with `Add`, `Has` and `Intersect` is
  stubbed; you implement the three methods.
- **Check:** `Has` after `Add`; Thistle's `{PEL LAN SQT}` ∩ Pipit-1's `{LAN CWR PEL}` →
  `{LAN PEL}`; works for ints.

---

## Chapter 17 · All Stations at Once

`all-stations-at-once` · 10 lessons: 4 code, 4 code_tests, 2 quizzes

Goroutines and channels: `go` and `sync.WaitGroup`, sending and receiving, closing and
ranging, fan-in, worker pools, `select`, and timeouts.

**Story arc.** Thistle's orbit is decaying, and Kessick Polytechnic has days left with their
satellite. Professor Marr asks for every station that can hear it to listen on every pass
(all but Ystmark, which Thistle's orbit never reaches). The telemetry service polls
stations one at a time, too slowly: by the time it reaches Lantern Atoll, the pass is over.
You make it concurrent. The chapter ends with Thistle's last beacons arriving at four
stations at once.

Programs here print in a fixed order however goroutines are scheduled: results go into a
slice by index or a map, and are printed after the work ends.

#### 17.1 Five Stations, One Moment
`five-stations-one-moment` · code · difficulty 6
- **Objective:** Start goroutines with `go` and wait for them with `sync.WaitGroup`.
- **Scenario:** Polling five stations one after another takes five times too long.
- **Exercise:** The starter starts a goroutine per station, each poll taking 20 ms, but
  doesn't wait, so `main` prints zeros before any result is in; you wait with a
  `WaitGroup` (`wg.Go`).
- **Check:** `YST: 3 frames` / `PEL: 2 frames` / `SQT: 4 frames` / `LAN: 2 frames` /
  `CWR: 5 frames`.

#### 17.2 A Channel for Frames
`a-channel-for-frames` · code · difficulty 6
- **Objective:** Send and receive values on a channel.
- **Scenario:** The Ystmark poller hands frames to the archiver through a channel.
- **Exercise:** A goroutine sends three sequence numbers; `main` receives one; you receive
  all three.
- **Check:** `archived 4182` / `archived 4183` / `archived 4184`.

#### 17.3 Quick Check: Who Blocks?
`quick-check-who-blocks` · quiz · difficulty 6
- **Objective:** Check when sends and receives wait.
- **Scenario:** Nkem sketches three channels on the whiteboard.
- **Exercise:** Three questions.
- **Check:** A send on an unbuffered channel waits for a receiver; a send on a buffered
  channel with room doesn't wait; receiving from a closed, empty channel returns the zero
  value at once.

#### 17.4 The Last Frame
`the-last-frame` · code · difficulty 6
- **Objective:** Close a channel when done, and `range` over it.
- **Scenario:** The archiver has to know when Ystmark's pass is over.
- **Exercise:** The producer never closes its channel, so `main`'s `range` waits forever
  and the program dies with a deadlock; you close it after the last send.
- **Check:** `archived 4182` … `archived 4185` / `pass over: 4 frames`.

#### 17.5 Many Feeds, One Stream
`many-feeds-one-stream` · code_tests · difficulty 7
- **Objective:** Merge several channels into one, closing it after every input is done.
- **Scenario:** Four stations hear Thistle at once; the archive wants one stream.
- **Exercise:** `merge(feeds ...<-chan int) <-chan int` returns a channel nobody writes to;
  you start a goroutine per feed, wait for all of them, then close.
- **Check:** feeds of 2, 3 and 0 values → all five values (compared sorted); the output
  closes (the test gives up after 1 s).

#### 17.6 Checking Chunks in Parallel
`checking-chunks-in-parallel` · code_tests · difficulty 7
- **Objective:** Share work among a fixed number of goroutines reading one channel.
- **Scenario:** A Furrow-A image has 2,000 chunks, and checking their CRCs one by one misses
  the delivery window.
- **Exercise:** `checkAll(chunks [][]byte, sums []uint32, workers int) []bool` returns all
  `false`; you send indexes to `workers` goroutines that write results by index.
- **Check:** 2,000 chunks with 3 corrupted → exactly those 3 false, with 1 worker and with 8.

#### 17.7 Whichever Comes First
`whichever-comes-first` · code_tests · difficulty 7
- **Objective:** Wait on several channels at once with `select`.
- **Scenario:** Thistle may be heard first at Lantern Atoll or at Cape Wren.
- **Exercise:** `first(lan, cwr <-chan string) string` receives only from `lan`; you
  `select` on both.
- **Check:** only `cwr` has a beacon → it's returned; only `lan` → it's returned.

#### 17.8 Give Up at LOS
`give-up-at-los` · code_tests · difficulty 7
- **Objective:** Add a timeout to `select` with `time.After`.
- **Scenario:** If no acknowledgment arrives before LOS, the uplink must stop waiting.
- **Exercise:** `waitAck(acks <-chan string, limit time.Duration) (string, error)` waits
  forever; you add a `time.After(limit)` case that returns `ErrNoAck`.
- **Check:** an ack already waiting → returned with nil; no ack and a 20 ms limit →
  `ErrNoAck` within 1 s.

#### 17.9 Quick Check: Deadlocks
`quick-check-deadlocks` · quiz · difficulty 6
- **Objective:** Spot a deadlock and a misuse of `close`.
- **Scenario:** Two snippets from last night's incident review.
- **Exercise:** Two questions about two short snippets in the readme.
- **Check:** The program that sends on an unbuffered channel in `main` with no receiver
  deadlocks; closing an already closed channel panics.

#### 17.10 Thistle's Last Pass
`thistles-last-pass` · code · difficulty 8
- **Objective:** Combine goroutines, a merged channel and a map to collect beacons.
- **Scenario:** Thistle's final pass before re-entry: four stations listen, and Professor
  Marr wants to know which beacons were heard, and by whom.
- **Exercise:** Four simulated station feeds are given; you merge them, record each beacon's
  stations in a map, and print in sequence order with sorted station lists.
- **Check:**
  ```text
  9001 heard by CWR LAN
  9002 heard by CWR LAN PEL
  9003 heard by PEL SQT
  9004 heard by SQT
  last beacon: 9004
  ```

---

## Chapter 18 · Locks and Deadlines

`locks-and-deadlines` · 9 lessons: 1 code, 6 code_tests, 2 quizzes

Shared state and cancellation: `sync.Mutex` and `sync.RWMutex`, contexts that cancel,
deadlines, passing `ctx` down a call chain, and stopping sibling goroutines after the first
failure.

**Story arc.** Your first on-call week. The status cache every service reads goes wrong under
load, a downlink keeps writing after LOS, and a stuck uplink holds a dish for an hour. You
fix the cache with a mutex, put a deadline on every pass with `context`, and end the week
running a whole pass day, concurrent and cancellable, through a storm at Ystmark.

#### 18.1 Lost Updates
`lost-updates` · code_tests · difficulty 7
- **Objective:** Protect shared data with `sync.Mutex`.
- **Scenario:** The frame counter undercounts when all five feeds update it at once.
- **Exercise:** `Counter` holds a mutex and a map; `Add(sat string)` does nothing yet; you
  lock, count and unlock with `defer`.
- **Check:** 50 goroutines × 200 `Add`s over 5 satellites → exactly 2,000 each; `Get` of an
  unknown ID → 0.

#### 18.2 Many Readers, One Writer
`many-readers-one-writer` · code_tests · difficulty 7
- **Objective:** Use `sync.RWMutex`, and hand out copies, not your internal map.
- **Scenario:** Every service reads station status, only the poller writes, and someone
  edited the cache through a map it returned.
- **Exercise:** `StatusCache.Snapshot()` returns its internal map; you return a copy
  (`maps.Clone`) under `RLock`, and use `RLock` in `Get`.
- **Check:** changing a snapshot leaves the cache unchanged; concurrent `Set`s and
  `Snapshot`s end in the right final state.

#### 18.3 Quick Check: Locks
`quick-check-locks` · quiz · difficulty 6
- **Objective:** Check how mutexes are used safely.
- **Scenario:** Mirela's two questions during your on-call handover.
- **Exercise:** Two questions.
- **Check:** A zero-value `sync.Mutex` is ready to use; `defer mu.Unlock()` right after
  `Lock` unlocks on every return path.

#### 18.4 Cancel at LOS
`cancel-at-los` · code_tests · difficulty 7
- **Objective:** Stop work when a context is cancelled.
- **Scenario:** A downlink kept writing after LOS and filled the archive with noise.
- **Exercise:** `receive(ctx context.Context, chunks <-chan []byte) (int, error)` ignores
  `ctx`; you `select` on `ctx.Done()` and return `ctx.Err()`.
- **Check:** a closed channel of 3 chunks → 3, nil; an already cancelled context and an
  open, empty channel → 0, `context.Canceled`.

#### 18.5 Every Pass Has a Deadline
`every-pass-has-a-deadline` · code_tests · difficulty 7
- **Objective:** Make a context that ends at a deadline with `context.WithDeadline`.
- **Scenario:** A pass ends at LOS whether the code is finished or not.
- **Exercise:** `passContext(parent context.Context, los time.Time) (context.Context,
  context.CancelFunc)` returns `context.WithCancel(parent)`; you use `WithDeadline`.
- **Check:** `Deadline()` equals LOS; with LOS in the past, `Err()` is
  `context.DeadlineExceeded`.

#### 18.6 Pass It Down
`pass-it-down` · code · difficulty 7
- **Objective:** Pass `ctx` as the first parameter through every call.
- **Scenario:** The stuck uplink ignored its deadline because one function in the chain
  started a fresh context.
- **Exercise:** `sendCommand` calls `transmit(context.Background(), ...)`; you add a `ctx`
  parameter and pass it down.
- **Check:** `uplink pipit-2 118: context deadline exceeded` (transmitting takes 1 s; the
  deadline is 30 ms).

#### 18.7 Quick Check: Context
`quick-check-context` · quiz · difficulty 6
- **Objective:** Check the conventions around `context`.
- **Scenario:** Nkem's three questions before you take the pager alone.
- **Exercise:** Three questions.
- **Check:** `ctx` goes first in a parameter list; call the cancel function even if the
  context ends on its own, to release it; pass a context to each call instead of storing it
  in a struct.

#### 18.8 First Failure Cancels the Rest
`first-failure-cancels` · code_tests · difficulty 8
- **Objective:** Cancel sibling goroutines when one fails, and return the first error.
- **Scenario:** If Ystmark's pre-pass check fails, the other stations' checks for that pass
  are wasted.
- **Exercise:** `runAll(ctx context.Context, checks []func(context.Context) error) error`
  runs every check and waits; you cancel a shared context on the first failure
  (`context.WithCancelCause`) and return that error.
- **Check:** all succeed → nil; one fails at once while the others wait on `ctx.Done()` →
  that error, with every check returned.

#### 18.9 Pass Day
`pass-day` · code_tests · difficulty 9
- **Objective:** Combine goroutines, a mutex, deadlines and errors in one small service.
- **Scenario:** Friday, your last on-call shift: a full day of passes through a storm at
  Ystmark.
- **Exercise:** `runDay(ctx context.Context, passes []Pass, run func(context.Context, Pass)
  error) Report` is stubbed; you run each station's passes in AOS order and the stations
  concurrently, give each pass a context ending at its LOS, record results under a mutex,
  and skip a station's remaining passes after `ErrStowed`.
- **Check:** a fake `run` and a fixed schedule of 9 passes → `Report{Done: 6, Failed: 1,
  Skipped: 2}`, the skipped ones being Ystmark's last two; each station's passes run in
  order; the test finishes well under 1 s.
