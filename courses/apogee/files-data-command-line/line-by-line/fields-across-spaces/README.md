---
slug: fields-across-spaces
title: Fields Across Spaces
type: code_tests
position: 4
difficulty: 4
topics: slices, strings, testing
created: 2026-11-06T16:34:04Z
modified: 2026-11-06T16:34:04Z
---
Paz wants each reading split into its four columns: the time; the dish, `LG` for La Grande or `CH` for the smaller dish, Chiquita; the compass direction it points, called azimuth; and its angle above the horizon, called elevation. The controller separates them with tabs on some lines, runs of spaces on others, and both on a few.

`strings.Fields` splits a string around every run of spaces, tabs and other white space, and ignores white space at either end. Splitting on one space instead gives an empty field for every extra space.

```go
cols := strings.Fields("  PEL \t 46.9\t-38.1 ")
fmt.Println(len(cols), cols) // 3 [PEL 46.9 -38.1]
```

## Your task

Implement `parseColumns`: return a reading's four columns in order. For any other number of columns, return nil and an error that reads like this, with the count found:

```
3 columns instead of 4
```
