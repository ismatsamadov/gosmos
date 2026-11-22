---
slug: the-countdown
title: The Countdown
type: code
position: 0
difficulty: 2
topics: loops
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Launch day: Pipit-3 is in orbit among forty other satellites, and its first predicted pass is five minutes out. Ingrid wants a countdown on the console.

A `for` loop repeats its block. The three-part form sets a counter, checks a condition before each round, and updates the counter after it. `++` adds one and `--` subtracts one:

```go
for orbit := 1; orbit <= 3; orbit++ {
	fmt.Println("orbit", orbit)
}
```

## Your task

In `main`, turn `minutes := 5` and the `Printf` below it into a loop that prints `T-5` down to `T-1`, so `AOS` prints once, after them.
