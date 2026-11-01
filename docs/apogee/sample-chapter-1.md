# Sample: chapter 1, lessons 1 to 5

The first five lessons of Go Foundations, written in full to calibrate writers: tone, length,
the size of a step, how a starter differs from its solution, how a test reads, and how
short a quiz is. They follow `story.md` and the syllabus in `go-foundations.md`.

Each lesson is laid out as a lesson package for `go run ./cmd/author lesson`: the bold path
above each block is the file to create, relative to the lesson's package folder (see
`go run ./cmd/author help`). Every block is the file's whole content, ending in a newline.

| # | Lesson | Type | Idea |
|---|---|---|---|
| 1.1 | Hello, Pipit | code | A program runs `main`; `fmt.Println` prints a line |
| 1.2 | The Pre-Pass Checklist | code | Statements run top to bottom |
| 1.3 | Notes for the Night Shift | code | `//` comments, and switching a line off |
| 1.4 | Checked by Tests | code_tests | A read-only test says what a function must return |
| 1.5 | Quick Check: First Programs | quiz | Checks 1.1 to 1.4 |

## 1.1 Hello, Pipit

**`hello-pipit/README.md`**

````markdown
---
slug: hello-pipit
title: Hello, Pipit
type: code
position: 0
difficulty: 1
topics: printing, programs
---
Welcome to Apogee. By tradition, your first program writes the message the ops room sends up to Pipit-1, which echoes it back on its next pass.

A Go program starts in `func main()`. There, `fmt.Println` prints the text between the quotes and ends the line:

```go
package main

import "fmt"

func main() {
	fmt.Println("Port Ellery standing by")
}
```

Leave `package main` and `import "fmt"` as they are; every program in this chapter starts with them.

## Your task

Change the message so the program prints exactly:

```text
Apogee ground calling Pipit-1
```

Press **Run** to see the output, then **Submit** to check it.
````

**`hello-pipit/starter/main.go`**

```go
package main

import "fmt"

func main() {
	fmt.Println("your message here")
}
```

**`hello-pipit/solution/main.go`**

```go
package main

import "fmt"

func main() {
	fmt.Println("Apogee ground calling Pipit-1")
}
```

**`hello-pipit/expected_output.txt`**

```text
Apogee ground calling Pipit-1
```

Notes for writers: the example prints a different line, so it shows the shape without
being the answer. The readme names Run and Submit once, here, and never again.

## 1.2 The Pre-Pass Checklist

**`pre-pass-checklist/README.md`**

````markdown
---
slug: pre-pass-checklist
title: The Pre-Pass Checklist
type: code
position: 1
difficulty: 1
topics: printing, statements
---
Pipit-1 rises over Port Ellery in ten minutes, and the checklist printer has shuffled Ingrid's steps.

Go runs the statements in `main` from top to bottom, one at a time. Each `fmt.Println` prints one line, so the order of the calls is the order of the lines:

```go
fmt.Println("Antenna parked")
fmt.Println("Antenna moving to the rise azimuth")
```

## Your task

Reorder the `fmt.Println` calls so the title prints first, then steps 1 to 4 in order. Move whole lines; don't change their text.

> [!NOTE]
> AOS, acquisition of signal, is the moment a satellite rises into a station's view.
````

**`pre-pass-checklist/starter/main.go`**

```go
package main

import "fmt"

func main() {
	fmt.Println("3. Point the antenna at the rise azimuth")
	fmt.Println("Pre-pass checklist: Pipit-1 over Port Ellery")
	fmt.Println("1. Confirm the pass time in UTC")
	fmt.Println("4. Send the echo message after AOS")
	fmt.Println("2. Tune the radio to 437.250 MHz")
}
```

**`pre-pass-checklist/solution/main.go`**

```go
package main

import "fmt"

func main() {
	fmt.Println("Pre-pass checklist: Pipit-1 over Port Ellery")
	fmt.Println("1. Confirm the pass time in UTC")
	fmt.Println("2. Tune the radio to 437.250 MHz")
	fmt.Println("3. Point the antenna at the rise azimuth")
	fmt.Println("4. Send the echo message after AOS")
}
```

**`pre-pass-checklist/expected_output.txt`**

```text
Pre-pass checklist: Pipit-1 over Port Ellery
1. Confirm the pass time in UTC
2. Tune the radio to 437.250 MHz
3. Point the antenna at the rise azimuth
4. Send the echo message after AOS
```

Notes for writers: a jargon word gets its gloss in a one-sentence callout the first time it
appears, then never again. The task says what not to change, because a learner who
retypes the lines can introduce a typo the check will catch.

## 1.3 Notes for the Night Shift

**`notes-for-the-night-shift/README.md`**

````markdown
---
slug: notes-for-the-night-shift
title: Notes for the Night Shift
type: code
position: 2
difficulty: 1
topics: comments
---
Nkem reviews your first pass log: "Customers read this. Switch the DEBUG lines off, but keep them for the night shift."

Go ignores everything from `//` to the end of the line. Use comments for notes to people, or to switch a line off without deleting it:

```go
// Pipit-1 beacons every 30 seconds.
// fmt.Println("DEBUG: beacon timer started")
```

## Your task

Turn the two `DEBUG` lines into comments. The other three lines should still print, in the same order.

> [!NOTE]
> LOS, loss of signal, is the moment the satellite sets and the link ends.
````

**`notes-for-the-night-shift/starter/main.go`**

```go
package main

import "fmt"

func main() {
	fmt.Println("Pass log: Pipit-1 over Port Ellery")
	fmt.Println("DEBUG: radio warm-up took 41 s")
	fmt.Println("AOS 04:12 UTC, LOS 04:21 UTC")
	fmt.Println("DEBUG: echo buffer holds 29 bytes")
	fmt.Println("Echo received: Apogee ground calling Pipit-1")
}
```

**`notes-for-the-night-shift/solution/main.go`**

```go
package main

import "fmt"

func main() {
	// Customers read this log, so the DEBUG lines stay off.
	fmt.Println("Pass log: Pipit-1 over Port Ellery")
	// fmt.Println("DEBUG: radio warm-up took 41 s")
	fmt.Println("AOS 04:12 UTC, LOS 04:21 UTC")
	// fmt.Println("DEBUG: echo buffer holds 29 bytes")
	fmt.Println("Echo received: Apogee ground calling Pipit-1")
}
```

**`notes-for-the-night-shift/expected_output.txt`**

```text
Pass log: Pipit-1 over Port Ellery
AOS 04:12 UTC, LOS 04:21 UTC
Echo received: Apogee ground calling Pipit-1
```

Notes for writers: the check can't see comments, only their effect, so the task is about
the effect. The 29 bytes are the length of the learner's message from 1.1; small
continuities like this make the world feel real without costing words.

## 1.4 Checked by Tests

**`checked-by-tests/README.md`**

````markdown
---
slug: checked-by-tests
title: Checked by Tests
type: code_tests
position: 3
difficulty: 2
topics: functions, testing
---
The scheduler never reads what your program prints. It calls your functions, and a test says what it expects back.

This lesson comes with `main_test.go`, which you can read but not edit. A test calls a function, compares the result with what it wants, and reports a failure with `t.Errorf` when they differ:

```go
func TestBand(t *testing.T) {
	got := band()
	want := "UHF"
	if got != want {
		t.Errorf("band() = %q, want %q", got, want)
	}
}
```

Here, **Run** and **Submit** both run the tests.

## Your task

Open `main_test.go` and find what `TestStationCode` wants. Make `stationCode` in `main.go` return exactly that, by changing the text after `return`.
````

**`checked-by-tests/starter/main.go`**

```go
package main

import "fmt"

// stationCode returns the code the scheduler uses for Port Ellery.
func stationCode() string {
	return "TODO"
}

func main() {
	fmt.Println(stationCode())
}
```

**`checked-by-tests/starter/main_test.go`** (read-only, as every `_test.go` file is)

```go
package main

import "testing"

func TestStationCode(t *testing.T) {
	got := stationCode()
	want := "PEL"
	if got != want {
		t.Errorf("stationCode() = %q, want %q", got, want)
	}
}
```

**`checked-by-tests/solution/main.go`**

```go
package main

import "fmt"

// stationCode returns the code the scheduler uses for Port Ellery.
func stationCode() string {
	return "PEL"
}

func main() {
	fmt.Println(stationCode())
}
```

What the learner sees when the starter runs (the time varies):

```text
--- FAIL: TestStationCode (0.00s)
    main_test.go:9: stationCode() = "TODO", want "PEL"
FAIL
exit status 1
FAIL	lesson	0.263s
```

Notes for writers: the readme's example has the same shape as the real test but a
different function and value, so the learner has to open the test to finish. Functions get
their own chapter later; here the learner only edits a returned string.

## 1.5 Quick Check: First Programs

**`quick-check-first-programs/README.md`**

````markdown
---
slug: quick-check-first-programs
title: Quick Check: First Programs
type: quiz
position: 4
difficulty: 1
topics: comments, programs, testing
---
Before you sit at a console during a live pass, Ingrid runs a quick check: three questions, a few seconds each. The first two are about this program:

```go
package main

import "fmt"

func main() {
	// fmt.Println("DEBUG: radio on")
	fmt.Println("AOS in 2 minutes")
}
```

Pick one answer for each question, then press **Submit**. A wrong pick shows you the right one.
````

**`quick-check-first-programs/quiz.md`**

```markdown
## Where does this program start running?

- [ ] At the first line of the file
- [ ] At `import "fmt"`
- [x] At the start of `func main()`

## What does this program print?

- [ ] `DEBUG: radio on`, then `AOS in 2 minutes`
- [x] `AOS in 2 minutes`
- [ ] Nothing, because the comment stops the program

## You can't edit `main_test.go`. What do you change to make its test pass?

- [ ] The `want` value inside the test
- [x] The function in `main.go` that the test calls
- [ ] The name of the test
```

Right answers: 3, 2, 2, marked `[x]` above.

Notes for writers: one question per idea from lessons 1.1 to 1.4, each answerable in a few
seconds. Every wrong answer is a mistake a beginner might make, and none is right under a
different reading: the program starts in `main`, not at `package` or `import`; a comment
does nothing at all; and `main_test.go` can't be edited, so the `want` value can't be the
thing to change.

## Verification

Checked on 2026-11-01 with Go 1.27.3 (darwin/arm64).

| Lesson | Starter | Solution | Run time | Readme |
|---|---|---|---|---|
| 1.1 `hello-pipit` | fails: prints `your message here` | passes | 0.26 s | 104 words |
| 1.2 `pre-pass-checklist` | fails: lines out of order | passes | 0.25 s | 100 words |
| 1.3 `notes-for-the-night-shift` | fails: prints the DEBUG lines | passes | 0.25 s | 98 words |
| 1.4 `checked-by-tests` | fails: `stationCode() = "TODO", want "PEL"` | passes | 0.39 s | 110 words |
| 1.5 `quick-check-first-programs` | | | | 63 words |

- Every block above was extracted into its lesson package and run the way the app runs it:
  in a fresh folder with `go.mod` reading `module lesson` and `go 1.27`, `go build` and the
  binary for `code` (stdout compared after the app's normalization), and
  `go test -count=1 -timeout 10s` for `code_tests`. gofmt and `go vet` are clean on every
  starter and solution.
- The five packages were imported into a scratch library with `go run ./cmd/author lesson`.
  `go run ./cmd/library check` and `go run ./cmd/library verify` both report ok: the four
  code lessons pass with their solutions and fail with their starters. `author export`
  wrote the packages back byte for byte.
- The quiz's program was run: it prints only `AOS in 2 minutes`. Each question has one
  defensible answer: nothing runs at `package` or `import`, a comment has no effect, and
  `main_test.go` is read-only, so its `want` value can't be the thing to change.
- Word counts include code and exclude the front matter.
