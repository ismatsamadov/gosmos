---
slug: commas-inside-names
title: Commas Inside Names
type: code_tests
position: 1
difficulty: 5
topics: slices, strings, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Paz ran Hollis's sheet through the van's quick reader, which cuts each line at its commas. The note `PRIORITY, SPRAYING FROM 06:00` came out as two cells, and a note with a line break became two rows.

In CSV, a cell in double quotes may hold commas and line breaks, and a quote inside it is written twice: `""`. A `csv.Reader` undoes all of it; `ReadAll` returns every record at once.

```go
line := `CWR,"STOWED, WIND ""72"""` + "\n"
cr := csv.NewReader(strings.NewReader(line))
record, _ := cr.Read()
fmt.Println(len(record)) // 2
fmt.Println(record[1])   // STOWED, WIND "72"
```

## Your task

Implement `readNotes`: read the sheet with a `csv.Reader` and return the NOTE cell of every booking, in order, without the header's. NOTE is the sixth cell, at index 5. Return any reading error.
