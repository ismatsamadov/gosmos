---
slug: the-missing-chunks
title: The Missing Chunks
type: code_tests
position: 7
difficulty: 6
topics: loops, maps
created: 2026-11-01T07:36:45Z
modified: 2026-11-06T16:51:28Z
---
Yusra, 08:20 UTC: "TL;DR: Furrowlight wants every missing chunk of last night's Furrow-A image listed by number before nine."

A file sent down from a satellite arrives in numbered pieces, or chunks, starting at 0. The ones that arrived sit in a map keyed by number. The map can't list what never came. Walk through everything that should be there instead, and look each one up:

```go
heard := map[string]bool{"YST": true, "SQT": true}
for _, c := range []string{"YST", "PEL", "SQT"} {
	if !heard[c] {
		fmt.Println(c, "is silent")
	}
}
// PEL is silent
```

## Your task

Complete `missing` so it returns, in increasing order, every number from 0 to `total-1` that isn't a key in `chunks`. A chunk that came down empty still came down. When nothing is missing, return no numbers.
