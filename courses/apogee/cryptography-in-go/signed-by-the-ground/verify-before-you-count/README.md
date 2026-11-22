---
slug: verify-before-you-count
title: Verify Before You Count
type: code_tests
position: 3
difficulty: 7
topics: crypto, errors, maps
created: 2026-11-07T07:46:48Z
modified: 2026-11-22T10:06:32Z
---
Wednesday. Mirela sends one forged command, numbered 18,446,744,073,709,551,615, the largest a `uint64` holds, to the code that refuses old commands sent again, the replay guard. The guard writes that number down before checking the signature, so every real command after it looks like a replay, for good.

Check first, change state last: a refused message must leave everything as it was. The guard's refusals are sentinel errors, a replay's wrapped with its numbers:

```go
e := fmt.Errorf("%w: %d, last %d",
	ErrReplay, 7, 9)
fmt.Println(errors.Is(e, ErrReplay)) // true
fmt.Println(e)
// replayed sequence number: 7, last 9
```

## Your task

`Accept` records `c.Seq` first and verifies after. Reorder it: build the canonical form and `verify` the signature, else `ErrForged`; then compare `c.Seq` with the satellite's last number, else `ErrReplay` as above; and only then record it.
