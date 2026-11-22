---
slug: check-it-right-away
title: Check It Right Away
type: code
position: 1
difficulty: 4
topics: errors, loops
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
Your `checkBattery` works, but the ingest loop calls it and ignores the answer, so Thistle's 6.20 V reading still lands in the store.

Check an error right after the call that returns it. When it isn't `nil`, handle it, here by printing it, and skip the work that needed success:

```go
err := checkWind(kmh)
if err != nil {
	fmt.Println("no pass:", err)
	return
}
fmt.Println("tracking")
```

## Your task

In the loop, keep `checkBattery`'s result in `err`. When it isn't `nil`, print a line like `6.20: rejected: battery too low` and go on to the next reading without storing it.

> [!TIP]
> Go often puts the call and the check on one line: `if err := checkWind(kmh); err != nil {`.
