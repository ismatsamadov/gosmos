---
slug: quick-check-embedded-files
title: Quick Check: Embedded Files
type: quiz
position: 3
difficulty: 7
topics: files
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Nkem checks what ships inside gctl before the van climbs on: the station table, four templates, and nothing that has to be found on the stick at run time.

`go build` copies an embedded file into the program. Inside an `embed.FS`, as in any `fs.FS`, a name is a slash-separated path relative to the root:

```go
//go:embed templates
var bundled embed.FS

// inside a function:
data, err := bundled.ReadFile("templates/day.txt")
```

## Your task

Answer three questions about when an embedded file is read and how a name inside a filesystem is written.
