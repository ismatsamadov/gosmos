---
slug: only-a-newer-plan
title: Only a Newer Plan
type: code_tests
position: 8
difficulty: 9
topics: goroutines, mutexes, pointers
created: 2026-11-11T12:51:42Z
modified: 2026-11-11T12:51:42Z
---
Friday 23 March. Wednesday's real trouble: two planning runs finished out of order, the older plan replaced the newer one, and the test Board showed a plan twenty minutes old. Nkem: "What happens when two planners touch the same plan?"

`CompareAndSwap(old, new)` stores `new` only if the pointer still holds `old`, and reports whether it did. So load the current schedule, decide, then swap: if another offer got in between, the swap fails, and you load again and decide again. That loop stays right in any order the runs finish.

```go
for c := peak.Load(); kmh > c; c = peak.Load() {
	if peak.CompareAndSwap(c, kmh) {
		break // a new peak gust
	}
}
```

## Your task

`Offer` stores `s` whatever its version. In a loop, load the current schedule (`nil` while there is none) and return `false` unless `s.Version` is greater; then `CompareAndSwap` it for `s`: return `true` if it swapped, or go round again.

> [!NOTE]
> A check and a plain `Store` lose a plan only now and then, so the test also reads `publish.go` for the swap.
