---
slug: read-the-first-sheet
title: Read the First Sheet
type: code_tests
position: 0
difficulty: 5
topics: errors, loops, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Hollis Vey, Furrowlight's ops manager, sent September's bookings from his spreadsheet as a CSV file: plain text, one row per line, cells separated by commas. His subject line, like everything he writes, is in capitals: SEPTEMBER BOOKINGS, PLEASE CONFIRM.

A `csv.Reader` reads such a file one record at a time: `Read` returns the next row's cells as a `[]string`, then `io.EOF` once the sheet ends. A misplaced quote is an error.

```go
cr := csv.NewReader(strings.NewReader("CWR,64\n"))
record, err := cr.Read()
fmt.Println(record, err) // [CWR 64] <nil>
_, err = cr.Read()
fmt.Println(err == io.EOF) // true
```

## Your task

Implement `readSheet`: read every record of `r`, the header included, and return them with a nil error at `io.EOF`. On any other error, return no records and that error.
