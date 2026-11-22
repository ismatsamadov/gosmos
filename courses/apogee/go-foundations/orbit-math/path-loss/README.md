---
slug: path-loss
title: Path Loss
type: code_tests
position: 6
difficulty: 5
topics: functions
created: 2026-11-01T07:36:45Z
modified: 2026-11-04T17:07:33Z
---
Tomás wants to know how much the signal weakens on its way to a station, its path loss. He measures the change in decibels (dB), a way to compare signal strengths. The calculation starts from the slant range, which your toolbox already computes.

A function can call another of yours, so each one stays small. This one builds on `periodMinutes` from the onboarding sheet:

```go
// periodHours returns an orbit's period in hours.
func periodHours(altKm float64) float64 {
	return periodMinutes(altKm) / 60
}
```

## Your task

Make `pathLossDB` return 20·log₁₀(d) + 20·log₁₀(f) + 32.44, where d is the slant range in km from `slantRangeKm` and f is `freqMHz`. `math.Log10` computes log₁₀.
