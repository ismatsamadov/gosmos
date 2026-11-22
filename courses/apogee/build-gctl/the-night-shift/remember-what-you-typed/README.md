---
slug: remember-what-you-typed
title: Remember What You Typed
type: code_tests
position: 1
difficulty: 7
topics: formatting, slices
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
At 05:00 nobody remembers what was typed at 04:41, and Ingrid needs it for the log before the morning team, the relief shift, takes over.

A line recorded before it runs is in the history even when its command fails, or when it is `history` itself. `%3d` pads a number to three columns, so the numbers line up:

```go
fmt.Printf("%3d  %s\n", 7, "stations")
//   7  stations
```

## Your task

1. Give `Shell` a `history *history` field, `&history{}` in `newShell`, and make `Run` pass each line to `s.history.add` before it runs.
2. Write `add`: keep every line that isn't blank, as typed.
3. Write `cmdHistory`, listing the kept lines numbered from 1, and add it to the table. An argument is an error.

```text
  1  passes PEL
  2  history
error: history takes no arguments
history   list the commands typed so far
```
