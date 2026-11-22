---
slug: quick-check-measuring-speed
title: Quick Check: Measuring Speed
type: quiz
position: 5
difficulty: 8
topics: profiling, testing
created: 2026-11-10T15:36:39Z
modified: 2026-11-11T17:33:38Z
---
Thursday evening, Mirela reads the week's speed checks before they merge: "A number proves nothing until I know how it was taken."

Last, she points at passplan's pass-ID check, which runs once for every planned pass:

```go
// validID reports whether id is a pass ID.
func validID(id string) bool {
	re := regexp.MustCompile(
		`^[a-z0-9-]+/[A-Z]{3}/\d{5}$`)
	return re.MatchString(id)
}
```

## Your task

Answer Mirela's three questions: whether `b.Loop` times the inputs built before it, which check of a speed-up holds on any Mac, and what each call of `validID` costs.
