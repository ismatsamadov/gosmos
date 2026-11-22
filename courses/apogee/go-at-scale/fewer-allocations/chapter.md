---
slug: fewer-allocations
title: Fewer Allocations
position: 9
---
Measure allocations with `testing.AllocsPerRun` and remove them without changing results: room made ahead with `make`, `strings.Builder`, appending numbers and times into a byte slice, one buffer reused, `sync.Pool` for scratch space, and struct keys instead of formatted ones; then find a struct's padding with `unsafe.Offsetof`, `Sizeof` and `Alignof`, and order its fields so a slice of it carries none.
