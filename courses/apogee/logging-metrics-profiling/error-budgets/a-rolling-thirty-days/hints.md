Work out where the window starts once, before the loop, from `now` and `30 * 24 * time.Hour`; then skip every result of another station, and every result whose LOS is outside the range.

The range is half-open: a LOS exactly at the start is out, so it must be after the start, and a LOS exactly at `now` is in, so it must not be after `now`. `time.Time` has `After` and `Before` for both.

If Lantern Atoll's count is too high, check the start and end: 30 days back is excluded and a pass after `now` hasn't happened yet. If Cape Wren's lost count is 0, a skipped pass isn't counted as lost.
