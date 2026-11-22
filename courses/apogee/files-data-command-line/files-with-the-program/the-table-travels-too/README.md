---
slug: the-table-travels-too
title: The Table Travels Too
type: code_tests
position: 0
difficulty: 7
topics: files, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T17:07:41Z
---
Paz copied gctl onto a fresh memory stick at 3,700 meters and left `stations.csv` behind. This morning gctl cannot say where Salar Quieto is, and nobody up here can send the file again.

A `//go:embed` comment right above a package-level variable tells `go build` to copy a file into the program. The file is read once, when the program is built; when it runs, the variable already holds the file's contents, whatever folder it runs in. The directive needs the `embed` package imported, even if only as `_`:

```go
import _ "embed"

//go:embed checklist.txt
var checklist string
```

## Your task

Embed `stations.csv` in a package-level string with a directive like the one above, and make `stationTable` return it. The test calls `stationTable` from an empty folder, so a table read from the disk at run time is not there.

> [!NOTE]
> No space after the slashes: `// go:embed` is an ordinary comment, and the variable stays empty.
