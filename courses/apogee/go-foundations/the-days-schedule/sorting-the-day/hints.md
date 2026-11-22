Your comparison gets two passes. When does it need to look at their stations at all?

`a.AOS.Compare(b.AOS)` already returns negative, 0 or positive, as `SortFunc` wants. Only when it returns 0 do the stations decide, through `cmp.Compare`.

Skeleton: `slices.SortFunc(passes, func(a, b Pass) int { ___ })`, where the body returns the time comparison unless it is 0.
