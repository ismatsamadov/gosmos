---
slug: hollis-gets-the-report
title: Hollis Gets the Report
type: code_tests
position: 8
difficulty: 6
topics: maps, sorting, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
The van sent Hollis its list of passes that delivered nothing. He sent it back from his spreadsheet with notes added, commas and line breaks included, and STATION moved again. His question, in capitals: HOW MANY AT EACH STATION?

The read-only `stationColumn` reads the sheet by its header, and `writeSheet` writes one with real quoting. Between them, count each station's passes, sort the stations and add a total.

```go
counts := map[string]int{}
for _, s := range []string{"LG", "CH", "LG"} {
	counts[s]++
}
keys := slices.Sorted(maps.Keys(counts))
fmt.Println(keys, counts["LG"]) // [CH LG] 2
```

## Your task

Implement `summarizeSheet`: read the stations from `r`, then write to `w` the header `STATION,MISSED`, one row per station in alphabetical order with its count, and a last row `TOTAL` with the sum. If the sheet can't be read, write nothing and return the error.
