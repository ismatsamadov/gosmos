---
slug: a-bucket-that-refills
title: A Bucket That Refills
type: code_tests
position: 0
difficulty: 8
topics: methods, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Brinewatch's client for the beta lost its network on Sunday night and has retried ever since without a pause: a request every millisecond. Nkem: "Ten at once is fine. A thousand a second isn't."

A token bucket draws the line between them. Each customer's bucket holds at most 10 tokens; a request takes one, and tokens flow back at 1 a second: 10 requests at once, then 1 a second. These tokens are counters, not the signed tokens of `POST /tokens`. Time comes from `now`, which tests set by hand:

```go
gap := 1500 * time.Millisecond
fmt.Println(gap.Seconds())      // 1.5
fmt.Println(min(8.5+1.5, 10.0)) // 10
```

## Your task

`Take` lets every request through. When `now` is after `b.Last`, add the seconds since `b.Last` times `refillPerSec` to `b.Tokens`, at most `bucketSize`, and set `b.Last` to `now`. Then take a token if a whole one is left, and return whether you did.

> [!NOTE]
> A zero `Bucket`'s `Last` is in the year 1, so its first `Take` fills it.
