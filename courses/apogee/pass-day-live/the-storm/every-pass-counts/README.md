---
slug: every-pass-counts
title: Every Pass Counts
type: code_tests
position: 4
difficulty: 10
topics: files, formatting, testing
created: 2026-11-12T09:13:52Z
modified: 2026-11-12T11:24:16Z
---
Friday 20 July, 15:00, the Fishbowl: the investors' week flies in full, Ingrid's script. Her go/no-go comes from its report.

A golden file holds output known to be right, in `testdata`, which Go builds skip; tests compare a run with it byte for byte.

## Your task

After the week, write `reports/week.txt`, titled with the window's first day's ISO week:

```text
2029-W30  booked  lost  avail  budget  left
YST          210     9  95.71       1    -8
```

Rows are `%-8s %7d %5d %6.2f %7d %5d`, one per rehearsed station in `newNetwork`'s order, then `all`, the columns summed. Booked counts a station's bookings in the window, matched or not; lost is booked minus flown; avail, flown over booked in percent, or `-` (`%6s`) with none; budget, booked × (10000 − target) / 10000 in integers, target 9900 at `LAN` and `CAL`, else 9950; left, budget minus lost.

Then each lost booking, by AOS then ID, with its outcome or `unplanned` and its reason:

```text
bk-8743 Thu 11:05:54 furrow-h/YST skipped stowed
```

Last, one of these verdicts:

```text
GO: every booked pass flown
NO-GO: 1 booked pass lost
NO-GO: 9 booked passes lost
```

The tests check the investors' week against `testdata`, then Yusra's moves.
