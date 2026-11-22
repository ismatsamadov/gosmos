---
slug: quick-check-fuzz-tests
title: Quick Check: Fuzz Tests
type: quiz
position: 2
difficulty: 8
topics: testing
created: 2026-11-11T12:51:44Z
modified: 2026-11-12T08:43:35Z
---
Before anyone fuzzes overnight again, Mirela has three questions about this test, its seeds in `lineSeeds` and the folder where Saturday's line is saved, `testdata/fuzz/FuzzParseLine/`.

```go
func FuzzParseLine(f *testing.F) {
	for _, line := range lineSeeds {
		f.Add(line)
	}
	f.Fuzz(func(t *testing.T, line string) {
		parseLine(line) // must never panic
	})
}
```

## Your task

Answer Mirela's questions: what plain `go test` runs `FuzzParseLine` on, what `-fuzz` adds to that, and what happens to a line the fuzzer finds.
