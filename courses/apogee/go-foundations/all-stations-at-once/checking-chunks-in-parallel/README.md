---
slug: checking-chunks-in-parallel
title: Checking Chunks in Parallel
type: code_tests
position: 5
difficulty: 7
topics: channels, goroutines
created: 2026-11-01T07:36:46Z
modified: 2026-11-04T17:07:34Z
---
Furrow-A's latest image came down in 2,000 chunks. Checking their CRCs, short checksums that reveal damaged data, one at a time misses the delivery window Yusra promised Furrowlight.

A worker pool shares the work among a fixed number of goroutines that all receive jobs from one channel, until it is closed. Each value goes to exactly one worker. Send indexes rather than data, and each worker can write its result into its own slot of a slice.

```go
for range 3 {
	wg.Go(func() {
		for path := range uploads {
			deliver(path)
		}
	})
}
```

## Your task

Write `checkAll`: start `workers` goroutines that receive indexes from one channel and set `ok[i]` with `chunkOK`. Send every index, close the channel, and wait for the workers before returning `ok`: the last chunks are still being checked after the final send.
