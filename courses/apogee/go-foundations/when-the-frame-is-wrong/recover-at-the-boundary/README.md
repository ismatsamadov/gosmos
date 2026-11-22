---
slug: recover-at-the-boundary
title: Recover at the Boundary
type: code_tests
position: 10
difficulty: 6
topics: defer, errors, functions
created: 2026-11-04T17:00:51Z
modified: 2026-11-06T16:51:28Z
---
One broken decoder must not take the Fishbowl's whole reading service down. Nkem puts a boundary around each decoder call so its failure can be logged and the service can keep working.

`recover()` stops a panic only when called directly inside a deferred function. It returns the panic value, or `nil` when there is no panic. Execution does not resume at the panicking line. Named results let the deferred function set the boundary's return values:

```go
func total() (n int) {
	defer func() { n++ }()
	return 4 // returns 5
}
```

## Your task

Add a deferred function to `runDecoder`. In it, keep the value `recover()` returns as `problem`; if that isn't `nil`, set `err` with `fmt.Errorf("decoder panic: %v", problem)`. Keep a successful decoder's value and nil error. Remove the import placeholder; recovered calls return a zero reading.
