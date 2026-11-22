---
slug: a-different-separator
title: A Different Separator
type: code_tests
position: 7
difficulty: 5
topics: runes, streams, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
At the last signal, Salar Quieto sent the record of where La Grande pointed through each pass, and Paz wants it read before the next stop. The computer that steers the dish, its controller, puts a semicolon between cells and `#` before comment lines.

A `csv.Reader` reads other dialects through two of its fields: `Comma` is the separator, and `Comment` is the character that starts a comment line. Quoting works as before, so a quoted cell may hold the separator.

```go
in := strings.NewReader("CODE|WIND\nCWR|64\n")
cr := csv.NewReader(in)
cr.Comma = '|'
records, err := cr.ReadAll()
fmt.Println(records, err)
// [[CODE WIND] [CWR 64]] <nil>
```

## Your task

Implement `readControllerCSV`: set the reader's `Comma` to `;` and its `Comment` to `#`, then return all its records and its error.
