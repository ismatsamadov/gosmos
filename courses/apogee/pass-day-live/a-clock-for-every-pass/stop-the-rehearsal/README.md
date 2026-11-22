---
slug: stop-the-rehearsal
title: Stop the Rehearsal
type: code_tests
position: 4
difficulty: 9
topics: context, select, time
created: 2026-11-11T16:41:21Z
modified: 2026-11-12T02:35:46Z
---
Thursday afternoon. Ingrid stops a run halfway when she spots a wrong setting, and the next run must not meet a pass the last one left behind.

When `ctx` ends, the whole week stops at once, wherever each pass is: open calls hang up, and a pass still waiting for its AOS never dials. `time.Sleep` can't be woken early; a timer in a `select` beside `ctx.Done()` can.

## Your task

When `ctx` ends, dial nothing more and hang up every open call, its LOS line saying `stopped`. Once every pass has ended, say when the run stopped, and return `context.Cause(ctx)`:

```text
Tue 12:00:00 LOS furrow-f/YST/00406 stopped
rehearsal stopped at Tue 12:00:00
```

The tests stop a two-day run at Ystmark and Tirrin Rock at Tuesday noon, with a cause, and check that no call is dialed after it and every open one is hung up at 12:00:00. A stop without a cause, before the first AOS, dials nothing and returns `context.Canceled`.
