---
slug: what-comes-next
title: What Comes Next
type: code_tests
position: 2
difficulty: 7
topics: sorting, time
created: 2026-11-06T16:59:09Z
modified: 2026-11-06T16:59:09Z
---
Thursday. At 04:00 the night shift asks one question over and over: which pass is next, and how long until it? Ingrid wants one command that answers it.

`After` tells whether one time is later than another, and `Sub` gives the time between two moments as a `time.Duration`, which `%v` prints the way ops reads it:

```go
wait := aos.Sub(now)
fmt.Printf("in %v\n", wait) // in 12m0s
```

## Your task

Write `cmdNext`: from the cache's passes, find the first one, by AOS, that begins after `s.cache.now()`, and print it with the wait. When none is left, say so. Any argument is an error. Add `next` to the table. At 04:00 and later, among other lines:

```text
next: pipit-1 over PEL at 04:12, in 12m0s
no more passes tonight
error: next takes no arguments
next      show the next pass
```
