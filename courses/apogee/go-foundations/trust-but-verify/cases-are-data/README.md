---
slug: cases-are-data
title: Cases Are Data
type: code_tests
position: 0
difficulty: 5
topics: structs, testing
created: 2026-11-01T07:36:46Z
modified: 2026-11-04T17:19:31Z
---
A one-character bug nearly let an old command change Cirrine-1's operating state. The customer's weather satellite needs the code that refuses repeated commands, the replay guard, to work every time. Mirela's new rule: nothing touches the uplink without tests its author wrote, starting with `canonical`, which lays a command out for signing.

Go tests usually keep their cases as data: a slice of structs, one per case, each with a name, the input and the result you want. A new case is one more line:

```go
var periodCases = []struct {
	name           string
	altKm, wantMin float64
}{
	{"pipit-1", 525, 94.99},
	{"thistle", 415, 92.72},
}
```

## Your task

Read `canonical` in `main.go`. Then add at least three cases to `canonicalCases` in `cases.go`, each with its own name, a different `Command` and the exact text you expect `canonical` to make of it.
