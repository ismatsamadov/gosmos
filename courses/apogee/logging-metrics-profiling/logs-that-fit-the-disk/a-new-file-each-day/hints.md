Work out today's day once, before the size check, from `r.Now()` turned to UTC and written with `time.DateOnly`, the same form as `r.day`. Two such days compare with `>` in date order.

There are now three cases: a new day and an empty file, where the file only takes the new day and a count of 0; a new day and a file with bytes, which rotates; and the size rule as before. A `switch` with no value, its cases tried in order, keeps them apart. `rotate` already names the old file by `r.day`, and `r.open` gives the fresh one today's day.

If a clock that steps back rotates the file, you compared the days with `!=`: only a later day rotates. If nothing rotates at 00:00:00, the day was taken without `UTC()`: the test's clock writes its times eight hours ahead of UTC, where it was already the 5th.
