---
slug: cirrines-first-sealed-command
title: Cirrine's First Sealed Command
type: code_tests
position: 7
difficulty: 9
topics: crypto, security
created: 2026-11-07T07:46:48Z
modified: 2026-11-07T19:06:42Z
---
Friday 25 August. Cirrine-1 rises over Ystmark at 07:30, and Ingrid holds the pass for its first sealed command: "AOS 07:30. Sealed PING, number 4471. Line on my mark."

`sealCommand` puts the chapter together, from the master keys uplinkd holds in `ring` to the line Ystmark queues. The sealed bytes hold the operation and its arguments joined by `|`, and the header carries the rest. `append` puts one element in front of a slice's others:

```go
codes := []string{"YST", "PEL"}
all := append([]string{"SQT"}, codes...)
fmt.Println(strings.Join(all, "|")) // SQT|YST|PEL
```

## Your task

1. Check `c` with `canonical`, then get its satellite's seal key from `ring.SatKey(keyID, "seal", c.SatID)`.
2. Seal `c.Op` and `c.Args` joined by `|` with `newAEAD`, a fresh `newNonce()` and the header's `AD()`.
3. Return `formatSealed`'s line, or `""` and the first error.

> [!NOTE]
> The test hands your line to yesterday's satellite stand-in, which must open it as the same command.
