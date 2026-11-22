---
slug: a-day-at-a-time
title: A Day at a Time
type: code_tests
position: 4
difficulty: 9
topics: iterators, sorting
created: 2026-11-11T21:35:30Z
modified: 2026-11-11T21:35:30Z
---
Wednesday 28 March. The third quarter's catalog holds 672,292 passes, and the old passplan hands a slice of all of them to every service that asks: about eighty-four megabytes a copy, made before the first pass is read.

An iterator can make its values only as the loop asks for them. A catalog runs in AOS order and different satellites' passes interleave, so a day must be predicted whole and sorted before its first pass goes out. Days don't interleave: sorted one at a time, they stay in order, and a loop that stops on 1 July costs one day, not ninety-two.

```go
for h := range 24 {
	for _, r := range readHour(log, h) {
		if !yield(r) {
			return // later hours are never read
		}
	}
}
```

## Your task

`Catalog` predicts every satellite at every station that can serve it (`canServe`) over all of `w`, sorts the lot with `byAOS` and yields it. Predict and sort one day at a time instead, as `Passes` walks its days, and yield each day before predicting the next.
