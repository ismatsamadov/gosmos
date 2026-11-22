---
slug: predictable-by-the-clock
title: Predictable by the Clock
type: code_tests
position: 0
difficulty: 6
topics: crypto, security
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Monday 10 July. Before Cirrine-1's commands are locked so only it can read them, each needs 12 bytes that are never used twice, called a nonce. Mirela reads the draft helper's first line: "It starts from the clock. Tell me the second, and I'll tell you the nonce."

A `math/rand/v2` generator works out its numbers from the one it starts with, its seed: the same seed gives the same numbers, and a clock leaves only a few seeds to try:

```go
sec := uint64(1846800000) // 10 July 2028, 00:00
a := rand.New(rand.NewPCG(sec, 0))
b := rand.New(rand.NewPCG(sec, 0))
fmt.Println(a.Uint32() == b.Uint32()) // true
```

## Your task

`rand.Read` from `crypto/rand` takes no seed: it fills a slice from your Mac's unpredictable source. Make `newNonce` fill its 12 bytes with it, importing `crypto/rand` in place of `math/rand/v2`, and drop the clock and the generator. The test plays Mirela, trying every second around your call.

> [!NOTE]
> `rand.Read` never returns an error and always fills the slice; if the Mac has no randomness to give, the program stops.
