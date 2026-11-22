---
slug: cutting-the-timeline
title: Cutting the Timeline
type: code_tests
position: 4
difficulty: 4
topics: slices
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:28Z
---
The Board has room for the next three passes, and the alert panel for a satellite's last two battery readings. Both need part of a slice.

`s[a:b]` is the part of `s` from index `a` up to, but not including, `b`. Leave out `a` to start at 0, or `b` to run to the end:

```go
s := []string{"YST", "PEL", "SQT", "LAN", "CWR"}
fmt.Println(s[1:3]) // [PEL SQT]
fmt.Println(s[3:])  // [LAN CWR]
```

## Your task

Make `nextN` return the first `n` pass IDs and `lastN` the last `n` readings. When `n` is more than the slice holds, return all of it.

> [!WARNING]
> On the five IDs in the tests, `ids[:9]` panics. Keep `n` within the slice's length before you slice.
