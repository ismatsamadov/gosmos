Loop over `results` with a count that starts at 0, and add one for each result that passes four checks: its `Station`, its `Outcome`, and its `LOS` against each end of the window.

`now.Add(-d)` is the window's start. `After` is false for a time equal to it, which leaves the start out, as the window wants; at the other end, a LOS that is not after `now` is in, `now` itself included.

If 01:41:28 counts 3, `!LOS.Before(...)` keeps a LOS equal to the window's start, and `After` doesn't. If 08:45 counts 2, flown passes are counted too.
