---
slug: doppler-shift
title: Doppler Shift
type: code_tests
position: 1
difficulty: 3
topics: functions
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:27Z
---
Pipit-1's beacon arrives ten thousand cycles per second above its expected radio frequency as it rises: a shift of 10 kHz. One cycle per second is one hertz, written Hz; a thousand is one kHz. That change is called the Doppler shift. Tomás wants the shift computed, not guessed.

A function can take several parameters, separated by commas, and neighbors of the same type can share one type name:

```go
// marginDB returns how far a received signal
// is above what the radio needs.
func marginDB(signalDBm, minDBm float64) float64 {
	return signalDBm - minDBm
}
```

## Your task

Make `dopplerHz` return `-freqHz * rangeRateMps / speedOfLight`. The range rate, how fast the distance to the satellite changes, is negative while it approaches, so the shift is positive as it rises.
