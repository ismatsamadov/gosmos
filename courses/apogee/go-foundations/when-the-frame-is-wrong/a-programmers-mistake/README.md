---
slug: a-programmers-mistake
title: A Programmer's Mistake
type: code_tests
position: 9
difficulty: 5
topics: errors, functions
created: 2026-11-04T17:00:51Z
modified: 2026-11-05T08:05:38Z
---
Nkem finds a decoder configured to hold zero readings. That is a mistake in Apogee's code, not a bad radio message; continuing would hide the broken configuration.

An `error` reports a failure the caller should handle, such as damaged input. `panic` stops normal execution for a broken internal assumption. Deferred cleanup still runs while calls unwind. Do not panic for malformed satellite data: return an error for that.

```go
if slots < 0 {
	panic("negative internal slot count")
}
```

## Your task

In `requireLimit`, panic with `"invalid decoder limit"` when `limit` is zero or negative. Positive limits return normally. The read-only tests catch the panic to check its value; no real satellite input reaches this helper.
