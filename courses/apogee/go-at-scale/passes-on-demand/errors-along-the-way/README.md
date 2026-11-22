---
slug: errors-along-the-way
title: Errors Along the Way
type: code_tests
position: 5
difficulty: 8
topics: errors, iterators, streams
created: 2026-11-11T21:35:30Z
modified: 2026-11-11T21:35:30Z
---
Thursday 29 March. Hollis's spreadsheet lists one more Furrow-A pass over Salar Quieto in February than bookingd does. A line of the pass archive, the file of every past pass, was damaged in a copy, and bookingd's reader skipped it without a word.

An `iter.Seq2[V, error]` yields a value and an error together, so a sequence can report a failure where it happens: yield the zero value with the error, then stop. The loop checks each error as it comes:

```go
for r, err := range readings(windLog) {
	if err != nil {
		return fmt.Errorf("wind log: %w", err)
	}
	fmt.Println(r.Kmh, "km/h")
}
```

## Your task

`ReadPasses` skips a line `parsePass` can't read. Yield the zero `PassView` and the error instead, wrapped with the line's number, and stop:

```text
archive line 3: 3 fields, want 4
```

If reading itself fails, yield the scanner's error with the zero `PassView` once the loop has ended.

> [!NOTE]
> Each archive line is a pass ID, then its AOS, LOS and max elevation.
