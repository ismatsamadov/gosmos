---
slug: helpers-point-at-the-caller
title: Helpers That Point at the Caller
type: code_tests
position: 5
difficulty: 6
topics: testing
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
Three Doppler shift cases failed last night, and all three failures pointed at the same line: inside `assertNear` in `checks.go`, not at the calls that failed.

`t.Helper()` marks the function that calls it as a test helper. When a helper reports a failure, `go test` names the line that called the helper, which is where the failing case lives:

```go
func assertBand(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("band = %q, want %q", got, want)
	}
}
```

A `testing.TB` parameter accepts a test's `*testing.T`.

## Your task

Make `assertNear` in `checks.go` a helper, so each failure points at the line in `main_test.go` that called it.
