---
slug: quick-check-fair-comparisons
title: Quick Check: Fair Comparisons
type: quiz
position: 6
difficulty: 8
topics: performance, testing
created: 2026-11-11T22:45:16Z
modified: 2026-11-11T22:45:16Z
---
Thursday 17 May. Tomorrow your table goes up where Tomás's forty was, and before the glass is wiped, Nkem has three last questions about fair comparisons.

They ask when and where to time two versions, what to make of a version that disagrees, and what a benchmark measures once its work changes its own input, as this one from the weather tools does:

```go
func BenchmarkMedian(b *testing.B) {
	gusts := windLog(5000) // as recorded
	for b.Loop() {
		median(gusts) // sorts gusts in place
	}
}
```

## Your task

Answer Nkem's three questions.
