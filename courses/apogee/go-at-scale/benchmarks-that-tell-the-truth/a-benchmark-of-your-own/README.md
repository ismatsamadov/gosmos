---
slug: a-benchmark-of-your-own
title: A Benchmark of Your Own
type: code_tests
position: 0
difficulty: 8
topics: performance, testing
created: 2026-11-11T22:45:16Z
modified: 2026-11-12T08:43:36Z
---
Nkem won't argue about the forty. She asks you for one thing: "Show me a benchmark anyone can run."

In January you timed passplan inside a test with `testing.Benchmark`. A benchmark anyone can run lives in a test file as a function named `BenchmarkXxx`: `go test -bench=Parse` runs every one whose name matches the pattern `Parse`, and prints a line for each. The setup comes before `for b.Loop()`, and only the loop's body is measured:

```go
// in parse_test.go
func BenchmarkParse(b *testing.B) {
	line := formatLine(drillFrame)
	for b.Loop() {
		parseLine(line)
	}
}
```

## Your task

The test file wraps your `benchPlanDay` as `BenchmarkPlanDay`. `benchPlanDay` predicts Cape Wren's candidates for take-over day with `stationDay`, then plans them once. Keep `stationDay` where it is, and call `planDay` inside `for b.Loop()`.

> [!TIP]
> In a terminal, `go test -bench=PlanDay -benchtime=20x` runs the loop exactly 20 times, as the tests do.
