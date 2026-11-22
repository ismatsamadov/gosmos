---
slug: catch-the-mutants
title: Catch the Mutants
type: code_tests
position: 1
difficulty: 6
topics: testing
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
Mirela plants three broken copies of `accept`, the core of the replay guard, and asks whether your cases would notice. She calls them mutants.

A case is worth keeping when a likely bug makes it fail. Bugs gather at boundaries: equal values, one more, one less. A case far from every edge passes on a broken copy too:

```go
stows := func(kmh int) bool { return kmh >= 72 }
mutant := func(kmh int) bool { return kmh > 72 }
// 95 can't tell them apart:
fmt.Println(stows(95) == mutant(95)) // true
// 72 catches the mutant:
fmt.Println(stows(72) == mutant(72)) // false
```

## Your task

The real `accept` is `seq > last`. The one case in `cases.go`, far ahead, catches none of Mirela's mutants: `seq >= last`, `seq > last+1` and `seq != last`. Add cases to `acceptCases` until each mutant gets at least one of them wrong.
