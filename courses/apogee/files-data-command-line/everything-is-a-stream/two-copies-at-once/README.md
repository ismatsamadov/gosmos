---
slug: two-copies-at-once
title: Two Copies at Once
type: code_tests
position: 4
difficulty: 3
topics: interfaces, streams, testing
created: 2026-11-04T17:36:06Z
modified: 2026-11-04T17:36:06Z
---
Paz needs a live view of the report and a saved copy. Reading the source twice would fail once it is a one-way recording stream.

`io.MultiWriter` makes one writer that sends each write to its destinations in order. If one destination fails, it returns that error without writing the same block to later destinations. It does not roll back earlier writes:

```go
both := io.MultiWriter(&first, &second)
_, err := fmt.Fprintln(both, "ready")
```

## Your task

Implement `mirrorReport` by copying `src` to a `MultiWriter` built from `screen` and `archive`, in that order. Return the copy count and error. Both successful destinations must receive the same bytes from a single consumption of the source.
