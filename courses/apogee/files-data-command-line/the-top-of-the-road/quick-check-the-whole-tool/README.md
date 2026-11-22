---
slug: quick-check-the-whole-tool
title: Quick Check: The Whole Tool
type: quiz
position: 5
difficulty: 7
topics: errors, files, programs
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Wednesday evening, Nkem reviews gctl's wiring from Port Ellery: where it checks, where it writes and where it reports trouble.

Each boundary has one job. Input is checked whole before any output starts, results go to standard output and problems to standard error, and a report reaches its real name only through a rename:

```go
err := validateBookings(bookings, stations)
if err != nil {
	return err // nothing written yet
}
```

## Your task

Answer three questions about the order in which gctl checks, writes and reports.
