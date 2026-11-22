---
slug: a-filesystem-in-the-binary
title: A Filesystem in the Binary
type: code_tests
position: 1
difficulty: 7
topics: files, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T17:07:41Z
---
Paz's reports start from four fixed templates, two of them in a `notices` folder, and they went missing with the table. She wants the whole `templates` folder built into gctl.

A `//go:embed` directive that names a folder embeds the files below it, except names starting with `.` or `_`, into an `embed.FS`: a small read-only filesystem inside the program. Names inside it start with that folder and always use `/`, so build them with `path.Join`, not `filepath.Join`. Its `ReadFile` method returns a file's bytes, or an error for a name it does not hold:

```go
//go:embed docs
var docs embed.FS

// inside a function:
data, err := docs.ReadFile("docs/dish/wind.txt")
```

## Your task

1. Embed the `templates` folder into `bundled` with a directive like the one above.
2. Implement `loadTemplate`: return the text of `templates/<name>`, or the error from `ReadFile`, so a missing template reports `fs.ErrNotExist`.
