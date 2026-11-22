---
slug: a-benchmark-inside-a-test
title: A Benchmark Inside a Test
type: code_tests
position: 2
difficulty: 8
topics: profiling, testing
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:38Z
---
Wednesday. Before anyone changes passplan, the program that plans the passes, Nkem wants a number to beat: "How long does one plan take now?" The app runs tests, never benchmarks, so the measuring happens inside a test.

A benchmark runs some work over and over and reports the time per run. Inside its body, `for b.Loop()` repeats the loop's body as often as measuring it takes, while whatever comes before the loop, such as building the inputs, runs once and isn't timed. `testing.Benchmark` runs such a body from an ordinary test:

```go
r := testing.Benchmark(func(b *testing.B) {
	words := strings.Fields(notes)
	for b.Loop() {
		slices.Sort(slices.Clone(words))
	}
})
fmt.Println(r.N > 1, r.NsPerOp() > 0)
// true true
```

## Your task

Write `benchPlan`, in `bench.go`: build the inputs once with `benchInputs`, then plan them inside `for b.Loop()` with `planDay`.

> [!NOTE]
> `r.N` is how many times the loop's body ran, and `r.NsPerOp()` its average time per run, in nanoseconds.
