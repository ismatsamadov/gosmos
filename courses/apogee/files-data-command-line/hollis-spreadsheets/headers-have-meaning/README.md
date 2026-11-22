---
slug: headers-have-meaning
title: Headers Have Meaning
type: code_tests
position: 2
difficulty: 5
topics: errors, loops, slices, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Hollis moved the STATION column to the front of his sheet this morning. The van's reader, which took each row's third cell, filed every booking under a satellite instead of a station.

Find a column by its header, not its position, and refuse a sheet where the name is missing or appears twice: with two STATION columns, either could be the right one.

```go
header := []string{"CODE", "WIND", "GUST"}
for i, name := range header {
	if name == "GUST" {
		fmt.Println("GUST is column", i) // 2
	}
}
```

## Your task

Implement `columnIndex`, which `stationColumn` uses: return the position of the header cell that says `name`. When no cell says it, or more than one does, return an error that names it.
