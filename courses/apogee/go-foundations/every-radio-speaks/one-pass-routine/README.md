---
slug: one-pass-routine
title: One Pass Routine
type: code_tests
position: 0
difficulty: 5
topics: functions, interfaces
created: 2026-11-01T07:36:46Z
modified: 2026-11-04T17:19:30Z
---
The radios use different frequency ranges: UHF for Pipit-1, S-band for Pipit-2 and X-band for Furrow-A. The estimate for bringing data down to Earth, the downlink, has been copied once for each radio. Tomás wants one routine for all three.

An interface type lists methods. A value of any type that has those methods fits in a variable of the interface type, so a function that takes the interface accepts them all:

```go
type Thermometer interface {
	TempC() float64
}

func tooHot(th Thermometer) bool {
	return th.TempC() > 45
}
```

## Your task

Add the method all three radios share to `Radio`. Then make `downlinkSeconds` return how many seconds `r` needs to bring `megabytes` down: megabytes times 8 is megabits, and the rate is in megabits per second (Mbps).

> [!NOTE]
> An interface that lists no methods, like `Radio` now, is satisfied by every type. Go has a built-in name for the empty interface: `any`.
