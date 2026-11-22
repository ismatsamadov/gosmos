---
slug: largest-fields-first
title: Largest Fields First
type: code_tests
position: 8
difficulty: 8
topics: performance, structs
created: 2026-11-13T11:32:36Z
modified: 2026-11-13T16:25:47Z
---
Thursday, 17:40. Nkem looks at your picture of `glance` and asks one question: "Same seven fields. Which order needs no dots?"

Order the fields by `unsafe.Alignof`, largest first, and each one starts right where the one before it ends; only the end may still need padding. One value saves a few bytes. A slice, though, stores `Sizeof` bytes for every element, padding included, so the padding repeats once per element: 613,496 times on Tomás's strip.

```go
type reading struct {
	At    int64   // at 0
	BattV float32 // at 8
	Ok    bool    // at 12, then 3 bytes
}
var r reading
fmt.Println(unsafe.Sizeof(r)) // 16, not 24
```

## Your task

Reorder `glance`'s fields so a value takes 24 bytes, exactly what its fields add up to. Keep each field's name, type and comment.

> [!NOTE]
> Go lays fields out in the order you declare them, and neither the compiler nor `go vet` will reorder them or warn you about the gaps.
