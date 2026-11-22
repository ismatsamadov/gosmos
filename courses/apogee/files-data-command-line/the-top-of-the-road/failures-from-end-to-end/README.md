---
slug: failures-from-end-to-end
title: Failures from End to End
type: code_tests
position: 7
difficulty: 7
topics: errors, programs, testing
created: 2026-11-06T16:34:05Z
modified: 2026-11-22T10:19:20Z
---
A copy of `bookings.json` on the stick was cut off halfway, and gctl printed an empty plan, as if Thursday had no booked passes. Nkem: "Works. Now make it say what it means."

An error path returns as soon as a step fails, before anything uses that step's result, and says what was being read. Wrapping with `%w` adds those words and keeps the original error for `errors.Is`. A failed step must never turn into a zero value that looks like real data:

```go
entries, err := readLog(path)
if err != nil {
	return Summary{}, fmt.Errorf("log: %w", err)
}
```

## Your task

Complete `planDay`: return each step's error at once, wrapped as `settings: ...` for the settings and `bookings: ...` for the bookings file. Call `validateBookings` before `newPlan`, and wrap its error as `bookings: ...` too.
