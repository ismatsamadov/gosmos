---
slug: breaking-the-cycle
title: Breaking the Cycle
type: code
position: 5
difficulty: 6
topics: compiler-errors, packages
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
On Nkem's whiteboard, `alert` points at `telemetry` for its frames, and `telemetry` points back at `alert` for its severities. In code those arrows are imports, and the build stops: `import cycle not allowed`.

Go refuses import cycles of any length: if A imports B, B can't import A, directly or through other packages. The usual fix moves what both need into a third package that imports neither:

```text
schedule -> board -> schedule
  refused: a cycle
schedule -> pass <- board
  fine: pass imports neither
```

## Your task

1. Move `Severity`, its constants and its `String` method from `alert/alert.go` to `severity/severity.go`.
2. In `telemetry` and `alert`, import `apogee/groundctl/severity` and write `severity.Severity`, `severity.Warning` and so on.
3. Remove `telemetry`'s import of `alert`.
