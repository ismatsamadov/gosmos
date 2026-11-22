Keep the two loops over satellites and stations, and inside them ask `findPasses` instead of `scanEverySecond`; `passViews` turns its contacts into passes as before. Add each station's passes to one slice declared before the loops, with `append` and `...`.

After the loops, nothing is per station any more: one `attach` gives the whole day's passes their bookings, and one `buildDay` orders them. It takes groups, so give it one group, the slice `attach` returns.

If the test reports too many `elevationDeg` calls, a call to `scanEverySecond` is still in a loop. If the plan differs, compare your loops with the old ones: the same satellites, stations, window and mask.
