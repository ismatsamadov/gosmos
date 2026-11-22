---
slug: quick-check-reading-the-numbers
title: Quick Check: Reading the Numbers
type: quiz
position: 2
difficulty: 8
topics: performance, testing
created: 2026-11-11T22:45:16Z
modified: 2026-11-12T08:43:36Z
---
Wednesday 9 May. Beside the forty, Nkem tapes up a run of the benchmarks from passplan's machine and asks what its numbers say.

`go test -bench=. -benchmem` printed these two lines, each wrapped here to fit:

```text
BenchmarkPlanDay-12     110263
     11078 ns/op   42000 B/op   9 allocs/op
BenchmarkHigh-12    1000000000
    0.3589 ns/op       0 B/op   0 allocs/op
```

`BenchmarkHigh` is from before `b.Loop`:

```go
func BenchmarkHigh(b *testing.B) {
	p := PassView{Pass: Pass{MaxElevDeg: 47.9}}
	for i := 0; i < b.N; i++ {
		high(p) // MaxElevDeg >= 10
	}
}
```

## Your task

Answer Nkem's three questions about the two lines.
