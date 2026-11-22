---
slug: one-rule-per-station
title: One Rule per Station
type: code_tests
position: 6
difficulty: 4
topics: functions, loops
created: 2026-11-04T16:56:59Z
modified: 2026-11-06T16:51:28Z
---
Nkem is checking the rule builder before it reaches five stations. Three test stations use minimum angles of 1, 2 and 3 degrees; every rule must remember its own minimum after the loop ends.

Since Go 1.22, a loop variable declared with `:=` gets a separate variable for each iteration. A closure keeps that iteration's variable, even after the loop ends. A variable declared outside the loop and assigned with `=` is shared instead:

```go
var first func()
m := 0 // declared outside the loop
for n := 0; n < 3; n, m = n+1, m+1 {
	if n == 0 {
		first = func() { fmt.Println(n, m) }
	}
}
first() // 0 3: its own n, but the shared m
```

## Your task

Change only the function literal in `stationRules` to compare `elevationDeg` with `minimum`, including equality. The supplied loop keeps all three closures so tests can call them after it has finished.
