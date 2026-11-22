Each counter is an `atomic.Int64`: call its `Add` method with an `int64`, so convert `planned` and `dropped` with `int64(...)` first.

`Add` makes three calls: one to each counter's `Add`, the last with 1 for the run. `Snapshot` returns each counter's `Load`, in the order planned, dropped, runs.

If the totals are right but the run count is 0, `Add` leaves `Runs` alone. If they come back in the wrong places, check the order of `Snapshot`'s results.
