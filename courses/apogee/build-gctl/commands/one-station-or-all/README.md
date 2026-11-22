---
slug: one-station-or-all
title: One Station or All
type: code_tests
position: 1
difficulty: 5
topics: errors, slices
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
Ingrid wants only Port Ellery's passes on her screen, and at 04:00 someone will mistype a station code. gctl must say so rather than print nothing.

A command checks its arguments before it does any work, and returns an error when they are wrong. `Run` already prints that error as `error: <message>`, and the session goes on.

## Your task

1. Give `Shell` a `stations []Station` field, set to `stationTable()` in `newShell`.
2. Make `passes CODE` print only that station's passes. `stationIndex` in `data.go` finds a code in the table, or returns -1.
3. Return an error for an unknown code or a second argument, so that these mistakes are answered with:

```text
error: unknown station "XYZ"
error: passes takes at most one station code
```
