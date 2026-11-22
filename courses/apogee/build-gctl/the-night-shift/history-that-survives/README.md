---
slug: history-that-survives
title: History That Survives
type: code_tests
position: 2
difficulty: 7
topics: errors, files
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
The relief shift starts gctl at 06:00 and wants to see what the night shift typed. Today every session forgets its history when it ends.

`os.OpenFile` with these flags opens a file for adding to its end, and makes it if it is missing:

```go
flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
f, err := os.OpenFile("log.txt", flags, 0o600)
```

## Your task

`main.go` sets `s.history.path` to `gctl_history`; an empty path, as `newShell` leaves it, keeps the history in memory only.

1. Write a `load` method that reads the file at `h.path` into `h.lines`; a missing file is an empty history. Call it at the start of `Run`.
2. Make `add` append each line it keeps to the file too, and give `add` an `error` result.
3. The first time the file fails, in `load` or `add`, return an error starting `history:`, and leave the file alone from then on. `Run` prints it as it prints a command's error, and the session goes on:

```text
error: history: read gctl_history: is a directory
```
