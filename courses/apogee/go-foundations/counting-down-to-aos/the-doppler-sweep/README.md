---
slug: the-doppler-sweep
title: The Doppler Sweep
type: code_tests
position: 7
difficulty: 5
topics: loops
created: 2026-11-01T07:36:45Z
modified: 2026-11-01T07:36:45Z
---
Pipit-3's orbit is still a rough guess, and so is its Doppler shift: the faint beacon could be anywhere within 12 kHz of the predicted frequency. Tomás steps the receiver through the offsets and keeps the strongest.

To find the best value in a loop, keep the best seen so far and replace it only when something beats it. Start from the first value, not from 0, which could beat every real one:

```go
coldestC := tempAt(0)
for minute := 1; minute < 60; minute++ {
	tempC := tempAt(minute)
	if tempC < coldestC {
		coldestC = tempC
	}
}
```

## Your task

Make `bestOffset` step from `fromHz` to `toHz`, both included, by `stepHz`, and return the offset with the highest `signalStrength`. On a tie, keep the first.
