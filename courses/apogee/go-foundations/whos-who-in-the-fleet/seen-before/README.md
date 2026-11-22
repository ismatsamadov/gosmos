---
slug: seen-before
title: Seen Before
type: code_tests
position: 6
difficulty: 5
topics: maps, slices
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:28Z
---
Some of Thistle's frames arrive twice, once through each of two stations. Joss, a student on its team at Kessick Polytechnic, finds readings doubled in his radiation plot. Each frame has a sequence number, a counter used to spot gaps and repeats.

A map with `bool` values works as a set: store `true` for each key you meet, and any key you haven't met reads as `false`.

```go
met := map[string]bool{}
met["PEL"] = true
fmt.Println(met["PEL"], met["CWR"]) // true false
```

## Your task

Complete `dedupe` so it returns the sequence numbers in `seqs` without repeats, each kept where it first appeared: `[4182 4183 4182 4184]` becomes `[4182 4183 4184]`.
