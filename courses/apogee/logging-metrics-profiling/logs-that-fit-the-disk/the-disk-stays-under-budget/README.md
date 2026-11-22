---
slug: the-disk-stays-under-budget
title: The Disk Stays Under Budget
type: code_tests
position: 7
difficulty: 8
topics: files, logging
created: 2026-11-10T15:36:38Z
modified: 2026-11-10T15:36:38Z
---
Friday 8 December. At 14:00 the new logging goes to Lantern Atoll, and first Nkem gives it a night on the lab bench worse than any real one: six minutes of lost signal in every 14, for twelve hours, 180,000 `carrier lost` warnings. Files of 64 KiB and a 1 MiB budget stand in for 64 MiB and 1 GiB.

A station's logger is the week's pieces in a row: `stamp` gives each record the bench clock's time, the `Sampler` thins the flood, the JSON handler writes each record as one line to the `Rotator`, and the rotator's `OnOpen` prunes the folder each time it opens a file. Each handler wraps the next:

```go
h := slog.NewJSONHandler(w, nil)
l := slog.New(stamp(h, clock.Now))
```

## Your task

Make `newStationLog` build a `Rotator` for `uplinkd` in `dir`, with `fileBytes` and `now`, whose `OnOpen` prunes with `keepDays` and `totalBytes`. Open it with `open`, and return a logger on `stamp`, around the `Sampler`, around the JSON handler, with the rotator as the `io.Closer`.

> [!WARNING]
> `stamp` goes outermost, so the sampler counts minutes by the bench clock and not by the Mac's.
