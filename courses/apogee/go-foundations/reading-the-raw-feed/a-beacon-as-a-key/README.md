---
slug: a-beacon-as-a-key
title: A Beacon as a Key
type: code_tests
position: 7
difficulty: 6
topics: arrays, bytes, maps, slices
created: 2026-11-04T17:00:51Z
modified: 2026-11-05T08:05:37Z
---
Port Ellery's receiver sometimes passes on the same Pipit-1 beacon twice. It needs to remember complete messages already seen, then hand the reading bytes to the existing decoder without another copy.

An array of comparable elements can be a map key; a slice cannot. A key keeps its own array value. Slicing an array gives a view of its backing storage, so changes through that slice change the array:

```go
codes := [3]int{1, 2, 3}
view := codes[1:]
view[0] = 8 // codes is now [1 8 3]
```

## Your task

`firstSeen` gets an initialized map: return `false` for a stored beacon; otherwise record it and return `true`. `readingBytes` gets a pointer to an array: return its slice from byte 3 up to, but excluding, byte 8.
