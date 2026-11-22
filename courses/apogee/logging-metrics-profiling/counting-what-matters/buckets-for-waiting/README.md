---
slug: buckets-for-waiting
title: Buckets for Waiting
type: code_tests
position: 4
difficulty: 8
topics: metrics, slices, time
created: 2026-11-10T15:36:37Z
modified: 2026-11-11T17:33:38Z
---
Tomás wants to know how long satellites take to acknowledge the commands a pass sends up: "I don't need every wait, only their shape! How many answer within a second, and how many keep us waiting a minute?"

A histogram sorts observations into buckets by size. `ack_wait_seconds` has buckets up to 0.5, 1, 2, 5, 10, 30, 60, 120 and 300 seconds, then one above them all. A wait goes in the first bucket whose bound is at or above it, and in no other, so exactly 2 seconds counts as "up to 2". The histogram also keeps how many waits it has seen and their sum, in whole milliseconds:

```go
d := 1500 * time.Millisecond
fmt.Println(d.Seconds())      // 1.5
fmt.Println(d.Milliseconds()) // 1500
```

## Your task

`Observe` does nothing. Add one to the first of `h.counts` whose bound in `h.Bounds` is at or above `d` in seconds, or to the last of `h.counts` when `d` is above every bound. Then add one to `h.count`, and `d` in whole milliseconds to `h.sumMs`.

> [!NOTE]
> `h.counts` has one element more than `h.Bounds`: the bucket above them all.
