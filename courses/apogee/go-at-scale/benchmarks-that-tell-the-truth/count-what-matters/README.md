---
slug: count-what-matters
title: Count What Matters
type: code_tests
position: 4
difficulty: 9
topics: iterators, performance, testing
created: 2026-11-11T22:45:16Z
modified: 2026-11-11T22:45:16Z
---
Friday 11 May. Another of passplan 2's benchmarks read Cape Wren's day from a stream made once, before its loop: from the second turn on the stream was empty, and `assign` planned nothing, very fast.

A benchmark can report its own unit of work beside the time: `b.ReportMetric(v, "contacts/op")` adds a column to its line, so a wrong count shows at once. After the loop, `b.N` holds how many turns ran:

```go
frames := 0
for b.Loop() {
	frames += len(parseAll(lines))
}
perTurn := float64(frames) / float64(b.N)
b.ReportMetric(perTurn, "frames/op")
```

## Your task

1. In `benchStream`, add up the contacts `assign` returns, and after the loop report them per turn as `contacts/op`. With the stream, that's 3.8.
2. Then pass `slices.Values(day)` to `assign` instead of `stream`: a sequence from `slices.Values` starts from the first pass each time it's ranged.

> [!NOTE]
> A function built on `iter.Pull`'s `next` hands out each value once: after the last, `next` returns false for good.
