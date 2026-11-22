---
slug: between-two-passes
title: Between Two Passes
type: code_tests
position: 6
difficulty: 9
topics: slices, sorting, time
created: 2026-11-07T07:46:47Z
modified: 2026-11-07T07:46:47Z
---
Ingrid's log for Rotation Day has a column she worked out by hand last week: when each station changed certificates. "Worked by hand. Before I sign, the code works it out too."

Her rule: a switch takes 10 minutes and keeps 5 minutes clear of every pass at that station on either side. So it must not touch a pass's time from 5 minutes before AOS, the moment the satellite comes into view, to 5 minutes after LOS, the moment it sets. The earliest such start is `from` itself, or 5 minutes after some pass's LOS. Two windows collide when each starts before the other ends:

```go
collide := aStart.Before(bEnd) &&
	bStart.Before(aEnd)
```

## Your task

Make `switchTime` sort `station`'s passes by AOS and return the earliest fitting start from `from` on. When no switch fits within 24 hours of `from`, return the error below, ending with `from` plus 24 hours in RFC 3339:

```text
no 10-minute gap for <station> before <end>
```
