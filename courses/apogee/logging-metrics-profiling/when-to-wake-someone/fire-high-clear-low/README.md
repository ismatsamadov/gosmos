---
slug: fire-high-clear-low
title: Fire High, Clear Low
type: code_tests
position: 3
difficulty: 8
topics: methods, metrics
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
A disk that bobs between 88% and 91% as recordings land and are collected would fire `disk-full`, resolve it and fire it again all night. Ingrid calls a rule like that flapping.

The cure is two lines instead of one, which engineers call hysteresis. A rule fires at `FireAt` and is over only when its value falls below a lower `ClearBelow`, so `disk-full` fires at 0.90 and resolves below 0.85, and a disk at 0.88 that has fired stays firing. Once a rule fires, `Observe` hands each sample to `whileFiring`. A heater works the same way:

```go
on := false // a heater: on below -8, off from -4
for _, c := range []float64{-9, -6, -3, -6} {
	on = c < -8 || on && c < -4
	fmt.Print(on, " ")
}
// true true false false
```

## Your task

`whileFiring` resolves as soon as a sample falls below `FireAt`. Resolve only below `r.ClearBelow`, so a value between the two lines keeps the rule firing.

> [!NOTE]
> Until the next lesson a firing rule repeats FIRING at every sample, so the test counts only when it fires and resolves.
