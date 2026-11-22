---
slug: too-big-to-keep
title: Too Big to Keep
type: code_tests
position: 6
difficulty: 9
topics: http, networking, security
created: 2026-11-07T07:44:28Z
modified: 2026-11-07T19:06:41Z
---
Thursday 8 June. Overnight a test recorder set up wrong sent a whole day of passes as one body, and Sunniva found the bench computer out of memory reading it.

A limit protects the server only if checked before the bytes are read. A chunk's size line says how much is coming, so the server can refuse as soon as that size would take the body past `max`. Compare the size with the room left: a 16-digit size can overflow a sum.

```go
used, room := int64(65000), int64(536)
size := int64(math.MaxInt64)
fmt.Println(used+size > 65536) // false
fmt.Println(size > room)       // true
```

## Your task

`readChunked` accepts chunks of any size. Before reading a chunk, refuse it with 413 when its size would take the body past `max`, as in `body over 65536 bytes`. A `max` of 0 means no limit.

> [!NOTE]
> The refused chunk is still on the line, so the server hangs up after its answer.
