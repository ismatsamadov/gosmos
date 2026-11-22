---
slug: the-standard-toolbox
title: The Standard Toolbox
type: code
position: 3
difficulty: 5
topics: generics, slices
created: 2026-11-01T07:36:46Z
modified: 2026-11-01T07:36:46Z
---
Nkem's rule: before you write a helper, look in `slices` and `maps`. The standard library already has generic versions of the loops everyone writes, and they come tested.

```go
temps := []float64{-4.5, 12, 3.25}
fmt.Println(slices.Min(temps)) // -4.5
slices.Sort(temps)             // sorts in place
fmt.Println(temps)             // [-4.5 3.25 12]
```

## Your task

Each of the three loops in `main` has a bug. Replace each loop with a single call: `slices.Contains`, `slices.Index` and `slices.Max`.

> [!WARNING]
> `slices.Max` panics on an empty slice, where your own `Max` returned `false`.
