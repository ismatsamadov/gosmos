---
slug: units-with-tilde
title: Units with ~
type: code_tests
position: 4
difficulty: 7
topics: generics, types
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:27Z
---
Ingrid wants every gauge on the Board to clamp its reading to its scale. `Clamp` works on a plain `float64`, but the battery gauge passes `Volts`, and the build stops: `Volts does not satisfy Real`.

A constraint can list types, joined with `|`. A bare `float64` admits only `float64` itself; `~float64` also admits every type whose underlying type is `float64`, such as `Volts`:

```go
// Whole admits int, int64 and every type
// built on them.
type Whole interface {
	~int | ~int64
}

type Orbits int // fits Whole, thanks to ~int
```

## Your task

Change `Real` in `main.go` so that `Volts`, `Kilometers`, `KmPerHour` and every other type built on `float32` or `float64` fit it.
