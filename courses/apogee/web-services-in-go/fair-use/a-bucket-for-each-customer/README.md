---
slug: a-bucket-for-each-customer
title: A Bucket for Each Customer
type: code_tests
position: 1
difficulty: 8
topics: maps, mutexes, time
created: 2026-11-07T06:18:32Z
modified: 2026-11-07T06:18:32Z
---
Yusra's TL;DR: "Brinewatch's flood must not spend Cirrine's share." One bucket for the whole service would let it.

So each customer gets a bucket of its own, kept in a map and made on first use. Requests arrive on many goroutines at once, so a mutex guards the map and every bucket in it. `Allow` also says what is left: the whole tokens, and, when no whole token is, how long until one is back, rounded up to a second, the unit customers' programs read:

```go
wait := 400 * time.Millisecond
secs := math.Ceil(wait.Seconds())
fmt.Println(time.Duration(secs) * time.Second)
// 1s
```

## Your task

`Allow` says yes to everyone, with 10 left. Under `l.mu`, find the customer's bucket in `l.buckets`, making it on first use, and `Take` from it at `l.now()`. Return whether it allowed the request, the whole tokens left and, when fewer than one is left, the time until a whole one is back, rounded up to a second; otherwise 0.
