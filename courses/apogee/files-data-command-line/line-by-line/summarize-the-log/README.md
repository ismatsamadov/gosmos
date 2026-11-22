---
slug: summarize-the-log
title: Summarize the Log
type: code_tests
position: 8
difficulty: 5
topics: maps, streams, testing
created: 2026-11-06T16:34:04Z
modified: 2026-11-06T17:07:42Z
---
Before the van climbs on, Paz wants the controller log in a few numbers: readings per dish and the highest elevation of all. Next year's log will be longer still, so the tool must not keep the lines it has read.

Keep running totals instead of records: each line updates the summary and is gone by the next `Scan`. A map counts readings per dish. Start the highest elevation from the first reading, not from 0, since readings can sit below zero.

```go
coldest := 0.0
for i, tempC := range []float64{4.5, -3.0, 1.5} {
    if i == 0 || tempC < coldest {
        coldest = tempC
    }
} // coldest is -3
```

## Your task

Implement `summarize` with a scanner and the read-only helpers in `log.go`. Count lines from 1, skip notes and blank lines with `isReading`, parse readings with `parseReading`, and return the summary with the scanner's error. A parsing error ends the scan: return it with the zero `Summary`, wrapped with `%w` and its line number:

```
line 2: 3 columns instead of 4
```
