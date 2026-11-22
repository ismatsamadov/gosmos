---
slug: the-pass-across-midnight
title: The Pass Across Midnight
type: code_tests
position: 6
difficulty: 6
topics: errors, testing, time
created: 2026-11-06T16:34:05Z
modified: 2026-11-06T16:34:05Z
---
Ingrid's pass log lists how long each pass lasted. Last night's final pass rose at 23:54 UTC and set at 00:05, and the old script, which subtracted clock readings, logged it as minus 23 hours 49 minutes.

A complete timestamp carries its date and offset, so `end.Sub(start)` measures the real time between two instants, across midnight and across offsets. `end.Before(start)` reports an end that comes first: two swapped readings, or one that is wrong.

```go
heard := time.Date(2027, time.September, 22,
	23, 59, 58, 0, time.UTC)
again := time.Date(2027, time.September, 23,
	0, 0, 28, 0, time.UTC)
fmt.Println(again.Sub(heard)) // 30s
```

## Your task

Implement `passDuration`: parse the AOS and the moment the satellite sets out of view, called LOS, both written in RFC 3339, and return the time between them. Return 0 and an error when either does not parse or when the LOS comes before the AOS.
