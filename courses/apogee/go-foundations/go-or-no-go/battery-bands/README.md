---
slug: battery-bands
title: Battery Bands
type: code_tests
position: 2
difficulty: 3
topics: conditionals
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
Before every pass, the operations team, or ops, sorts the satellite's battery into a group, or band. The limits hang on the Fishbowl wall: FULL from 8.2 V, OK from 7.2 V, LOW from 6.6 V, and CRITICAL below that.

`else if` adds another condition to an `if`. Go tries the conditions from the top, runs the block of the first true one and skips the rest; a final `else` catches whatever is left:

```go
if marginDB >= 6 {
	fmt.Println("strong link")
} else if marginDB >= 0 {
	fmt.Println("weak link")
} else {
	fmt.Println("no link")
}
```

## Your task

In `batteryLine`, replace the TODO with an `if`, `else if` and `else` chain that sets `band` to `"FULL"`, `"OK"`, `"LOW"` or `"CRITICAL"` for `battV`, using the limits above. The `Sprintf` after it adds the voltage.
