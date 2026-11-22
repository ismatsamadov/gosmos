---
slug: capital-letters-export
title: Capital Letters Export
type: code
position: 1
difficulty: 5
topics: compiler-errors, packages
created: 2026-11-01T07:36:46Z
modified: 2026-11-06T16:51:28Z
---
Joss, on Thistle's team at Kessick Polytechnic, asks for its orbital period. Last night's tidy-up of `orbit` gave it a `cube` helper but also renamed `PeriodMinutes` to `periodMinutes`, and now `main` can't reach it: `undefined: orbit.periodMinutes`.

A name is visible outside its package only if it starts with a capital letter: it's *exported*. A lowercase name stays private to its package, which is how a package keeps its helpers to itself:

```go
package board

const rows = 12 // unexported: only board sees it

// Rows is exported: main can call board.Rows().
func Rows() int { return rows }
```

## Your task

In `orbit/orbit.go`, rename `periodMinutes` to `PeriodMinutes`, doc comment included, and call it by its new name in `main.go`. Leave `cube` lowercase: only `orbit` needs it.
