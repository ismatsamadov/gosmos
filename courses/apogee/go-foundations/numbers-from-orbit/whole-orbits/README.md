---
slug: whole-orbits
title: Whole Orbits and Leftovers
type: code
position: 4
difficulty: 2
topics: types, variables
created: 2026-11-01T07:36:45Z
modified: 2026-11-05T08:05:37Z
---
Ingrid asks how many full orbits Thistle, a small satellite that students at Kessick Polytechnic built from two 10 cm cubes, makes in a day. This kind of satellite is called a cubesat. She also needs the minutes left after its last orbit.

Between two `int`s, `/` divides and drops the fraction, so the result is always whole. `%` gives the remainder, what that division leaves over:

```go
conduits := 23
fmt.Println(conduits / 4) // 5 full crates
fmt.Println(conduits % 4) // 3 conduits left over
```

## Your task

Replace the two zeros: set `fullOrbits` to the number of whole orbits in `minutesPerDay`, and `minutesLeft` to the minutes that remain, so the program ends with:

```text
Full orbits: 15
Minutes left: 45
```
