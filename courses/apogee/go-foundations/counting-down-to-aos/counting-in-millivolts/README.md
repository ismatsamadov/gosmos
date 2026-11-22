---
slug: counting-in-millivolts
title: Counting in Millivolts
type: code
position: 1
difficulty: 3
topics: loops
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:19:30Z
---
Sunniva sends a message over the slow link carrying Ystmark's data to headquarters, its backhaul: Pipit-3's battery is at 6700 mV and gains 150 per orbit. A millivolt (mV) is one thousandth of a volt. She needs to know when it reaches 8200 mV, or full.

A `for` with only a condition repeats as long as the condition holds, which suits loops whose number of rounds you don't know in advance. With no condition at all, `for { ... }` repeats forever, unless something inside it leaves the loop. Operators such as `*=` and `+=` update a variable in place: `x += 3` is short for `x = x + 3`.

```go
windowHz := 500
for windowHz < 12000 {
	windowHz *= 2
}
fmt.Println(windowHz) // 16000
```

## Your task

Replace the `TODO` in `main` with a loop that runs while `battMV` is under `fullMV`, adding `chargeMV` and counting one orbit each round.

> [!TIP]
> The count uses whole millivolts because a `float64` can't hold most decimals exactly: in volts, ten charges of 0.15 take 6.7 to 8.200000000000003.
