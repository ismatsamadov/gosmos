---
slug: quick-check-modules
title: Quick Check: Modules
type: quiz
position: 6
difficulty: 5
topics: modules
created: 2026-11-01T07:36:46Z
modified: 2026-11-22T10:08:26Z
---
groundctl's first tagged release, v1.0.0, goes out on Friday. Before the tag is pushed, Nkem asks three last questions about the file every lesson in this chapter kept read-only:

```text
module apogee/groundctl

go 1.27
```

groundctl uses only the standard library today. If it ever needs another module, `go get` adds a `require` line to this file and creates a second file beside it, `go.sum`.

## Your task

Answer Nkem's three questions.
