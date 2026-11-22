---
slug: the-quarter-in-numbers
title: The Quarter in Numbers
type: code
position: 6
difficulty: 9
topics: formatting, metrics, printing
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:38Z
---
Friday 12 January. This afternoon Yusra sends customers Apogee's first report in numbers: "TL;DR: the last quarter, station by station."

A table lines up when every value has a width: `%-7s` pads a string to 7 characters on the right, `%7d` right-aligns a number in 7, and `%6.2f` writes two decimals in 6. `main` prints the header with the same widths, as strings:

```go
fmt.Printf("%-5s|%5d|%6.2f|\n", "SQT", 3, 99.7671)
// SQT  |    3| 99.77|
```

## Your task

Write `row` with your `availability` and `budget`: the station's code, its booked and lost passes, its availability in percent, its budget and what is left of it, as `%-7s %7d %5d %6.2f %7d %5d`. When `availability` refuses the counts, the row is the code and the error. The first row reads:

```text
YST        2208     2  99.91      11     9
```

> [!NOTE]
> Lantern Atoll's row ends in -11: eleven passes more than its budget allowed.
