---
slug: where-the-padding-goes
title: Where the Padding Goes
type: code
position: 7
difficulty: 8
topics: performance, structs
created: 2026-11-13T11:32:35Z
modified: 2026-11-13T16:25:47Z
---
Thursday, 17:10. Tomás's preview strip of take-over day on the test Board draws the window's 613,496 candidates from one `glance` value each, and his battered notebook adds its seven fields up to 24 bytes, 14.7 MB in all. The strip's slice takes 19.6 MB.

A field starts at a multiple of its alignment, so the compiler may leave unused bytes, padding, between a small field and a larger one, and it rounds a struct's size up to a multiple of its largest alignment. `unsafe.Offsetof`, `unsafe.Sizeof` and `unsafe.Alignof` report the layout it chose:

```go
type reading struct {
	Ok    bool    // at 0, then 7 bytes of padding
	At    int64   // at 8
	BattV float32 // at 16, then 4 more
}
var r reading
fmt.Println(unsafe.Offsetof(r.BattV)) // 16
fmt.Println(unsafe.Sizeof(r))         // 24
```

## Your task

`main` draws `glance` byte by byte: each field's number, 1 to 7 as they're declared, on every byte it takes, and a dot for padding. Mark the six fields after `Sat` the way `Sat` is marked.
