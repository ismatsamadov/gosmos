---
slug: quick-check-filesystem-views
title: Quick Check: Filesystem Views
type: quiz
position: 7
difficulty: 7
topics: files, interfaces
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Nkem reviews where gctl's files come from before the last day on the road: built into the program, from the stick, or from a test's memory.

An `fs.FS` opens and reads files; it has no method that creates or changes one. A view made by `fs.Sub` and a test's `fstest.MapFS` are both `fs.FS` values:

```go
view, err := fs.Sub(resources, "templates")
fixture := fstest.MapFS{
	"stations.csv": {Data: table},
}
```

## Your task

Answer three questions about reading through an `fs.FS`, the views `fs.Sub` makes, and `fstest.MapFS`.
