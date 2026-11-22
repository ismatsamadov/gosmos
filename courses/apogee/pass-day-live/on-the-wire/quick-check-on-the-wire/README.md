---
slug: quick-check-on-the-wire
title: Quick Check: On the Wire
type: quiz
position: 4
difficulty: 9
topics: bytes, networking, time
created: 2026-11-11T22:09:38Z
modified: 2026-11-12T10:46:09Z
---
Monday evening. Sunniva has three questions before the backhaul meets the rehearsal for real.

This is how a pass reads its call, with `idle` a reader that moves the call's read deadline 90 s past `time.Now()` before every read:

```go
br := bufio.NewReader(idle{conn})
if err := readBanner(br, code); err != nil {
	return callFailed
}
for {
	f, err := readFrame(br)
	// ...
}
```

## Your task

Answer Sunniva's questions: why one reader, why a deadline moved before every read, and what a bad CRC and a bad length each mean.
