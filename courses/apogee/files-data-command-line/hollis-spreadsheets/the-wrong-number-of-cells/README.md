---
slug: the-wrong-number-of-cells
title: The Wrong Number of Cells
type: code_tests
position: 4
difficulty: 5
topics: errors, loops, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Hollis deleted one booking's cell for the moment its satellite sets, its LOS, instead of clearing it, and the row came out one cell short. He answers once a day, so Paz wants every bad row named in one message.

Set a `csv.Reader`'s `FieldsPerRecord` to 6, and a record of 5 or 7 cells comes with an error that names its line and wraps `csv.ErrFieldCount`. That error doesn't stop the reader: the next `Read` goes on.

```go
cr := csv.NewReader(strings.NewReader("CWR\n"))
cr.FieldsPerRecord = 2
_, err := cr.Read()
fmt.Println(err)
// record on line 1: wrong number of fields
fmt.Println(errors.Is(err, csv.ErrFieldCount))
// true
```

## Your task

Implement `checkedRows` with `FieldsPerRecord` set to 6: keep the records of the right width, header included, and collect the others' errors. At `io.EOF`, return the records and `errors.Join` of those errors; on any other error, return no records and that error.
