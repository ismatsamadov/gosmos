---
slug: the-incident-review
title: The Incident Review
type: code
position: 7
difficulty: 10
topics: formatting, sorting, tracing
created: 2026-11-10T15:36:37Z
modified: 2026-11-10T15:36:37Z
---
Friday 23 February, 15:00, the Fishbowl: the incident review, where the team works out what went wrong and what to change, blaming no one. Nkem puts the bench replay of 13 February on the Board as one picture.

A timeline turns spans back into steps: each root, a pass, in order of its start, followed by the spans whose parent it is, in the order they started. A stable sort by start keeps spans of the same second in the order given. Padding lines up the columns: `%-24s` fills one from the left, `%8s` from the right, and a duration prints as `11m5s`:

```go
wait := 9*time.Minute + 30*time.Second
fmt.Printf("[%-8s][%7s]\n", "prep", wait)
// [prep    ][  9m30s]
```

## Your task

Write `timeline`: for each root, by start, the line `%s %-24s %8s` (its start as `15:04:05`, `pass ` and its pass ID, its duration); after it, `%s   %-22s %8s` for each of its children, in the order they started, with their names.

> [!NOTE]
> `benchReplay`, in `replay.go`, holds the spans in the order they ended, as `Spans` returns them.
