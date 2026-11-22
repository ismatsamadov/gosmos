---
slug: numbers-from-cells
title: Numbers from Cells
type: code_tests
position: 5
difficulty: 5
topics: conversions, errors, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Hollis's pass sheet gives the highest each pass climbs above the horizon, its maximum elevation, in degrees under ELEV. In one row he typed `HIGH`. Paz wants the error to point at that cell the way Hollis's spreadsheet would: by row and column.

Parse each cell with `strconv`, and wrap a failure with where it happened. A spreadsheet numbers rows by record, the header as row 1, so a note with a line break still takes one row.

```go
_, err := strconv.ParseFloat("GUSTY", 64)
fmt.Println(fmt.Errorf("row 7, WIND: %w", err))
// row 7, WIND: strconv.ParseFloat:
//   parsing "GUSTY": invalid syntax
```

## Your task

Complete the loop in `readPassRows`: parse each record's ELEV cell, at index 2, into `MaxElevDeg`. If a cell doesn't parse, return no rows and an error that starts `row N, ELEV: ` and wraps the `strconv` error.
