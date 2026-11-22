---
slug: a-span-by-hand
title: A Span by Hand
type: code_tests
position: 4
difficulty: 9
topics: context, mutexes, tracing
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:39Z
---
Friday 16 February. A trace file stays on the bench, and the next failure won't happen there. Nkem wants every uplinkd to write down what each pass did, step by step, as records the night shift can search.

Each record is a span, one timed step of work. The spans of one pass share a trace ID; each has an ID of its own and names its parent, the span it belongs to. `Start` finds the parent in the context and returns a new context carrying the new span, so the steps inside it find it in turn:

```go
ctx, pass := tr.Start(ctx, "pass")
_, worker := tr.Start(ctx, "worker")
fmt.Println(worker.ParentID == pass.SpanID)
// true
```

## Your task

Finish `Start`:
1. Copy the trace ID and parent from the span in `ctx`, or draw a trace ID, `fmt.Sprintf("%016x", t.Rand.Uint64())`.
2. `fmt.Sprintf("%08x", t.Rand.Uint32())` is the span's own ID, drawn next; hold `t.mu` for both draws.
3. Set `Start` from `t.Now()` and `PassID` from `passFrom(ctx)`; return `withSpan(ctx, s)` and the span.

> [!NOTE]
> The test seeds `Rand` with `rand.NewPCG(1, 2)`, so the first trace is `c4f5a58656eef510`.
