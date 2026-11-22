---
slug: read-any-filesystem
title: Read Any Filesystem
type: code_tests
position: 2
difficulty: 7
topics: files, interfaces, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-22T10:21:59Z
---
Nkem's review asks one question: how would a test hand the table reader a small table of its own, without rebuilding gctl?

`fs.FS` is the interface for any read-only tree of files: an `embed.FS`, a folder on the disk through `os.DirFS`, or a filesystem a test builds in memory. `fs.ReadFile(fsys, name)` reads one file from whichever one it is given. Names are slash-separated and relative, with no leading `/`:

```go
func checklist(fsys fs.FS) (string, error) {
	b, err := fs.ReadFile(fsys, "ops/checks.txt")
	return string(b), err
}
```

## Your task

Implement `readStationsFS`: read `stations.csv` from `fsys` with `fs.ReadFile`, then return what the supplied `parseStations` makes of the bytes. Return a read error as it is, so a missing table reports `fs.ErrNotExist`.
