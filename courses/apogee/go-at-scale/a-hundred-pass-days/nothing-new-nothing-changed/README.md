---
slug: nothing-new-nothing-changed
title: Nothing New, Nothing Changed
type: code_tests
position: 4
difficulty: 9
topics: context, slices, time
created: 2026-11-12T03:18:38Z
modified: 2026-11-12T03:18:38Z
---
Tuesday 26 June. In Monday's rehearsal every ten-minute replan planned the day again from midnight and moved contacts nobody had asked to move, so the test Board blinked with every tick.

A replan should carry the schedule on: contacts that rose before the replan's instant stay as they are, flown or flying, and the rest is planned from the dishes they leave, through yesterday's hand-over. With no news, that gives the schedule you had, contact for contact, so there is nothing to publish.

A published schedule never changes, and `append` onto the start of one can write over the rest of it:

```go
kept := old[:2]
both := append(kept, add...) // overwrites old[2:]
safe := slices.Concat(kept, add)
```

## Your task

`replanFrom` plans `cur`'s whole window, the stretch of time it covers, again with no hand-over. Plan only from `at` to the end of `cur`'s window, with `handover(cur, at, cfg.Turnaround)` and `stows`, and return a new schedule over `cur`'s window: `cur`'s contacts that rose before `at`, then the new ones.

> [!NOTE]
> Version 1 of take-over day, replanned at 09:47 with nothing new, gives the same 2,016 contacts, 1,235 of them from 09:47 on.
