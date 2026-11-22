---
slug: quick-check-packages
title: Quick Check: Packages
type: quiz
position: 2
difficulty: 4
topics: modules, packages
created: 2026-11-01T07:36:46Z
modified: 2026-11-22T10:08:26Z
---
Nkem sketches groundctl's new layout on the whiteboard, then caps the marker: three quick questions before the next package moves out of `main`.

```text
groundctl/
  go.mod          module apogee/groundctl
  main.go         package main
  telemetry/
    parse.go      package telemetry
```

`parse.go` declares two functions, `Parse` and `splitFields`, and `main.go` will call one of them. A teammate also wants a second file in the folder, `telemetry/limits.go`, for the battery limits.

## Your task

Answer Nkem's three questions.
