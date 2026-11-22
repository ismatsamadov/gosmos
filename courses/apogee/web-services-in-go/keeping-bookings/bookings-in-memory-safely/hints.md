Start each method by returning the error of `ctx.Err()` if there is one, then call `m.mu.Lock()` and defer `m.mu.Unlock()`, so the rest of the method runs under the lock until it returns.

In `Create`, look through `m.bookings` for one with the same `SatID` and `Station` and an `AOS` that is `Equal` to `b.AOS` before adding anything. The ID is `bk-` and `m.next` in four digits, after which `m.next` goes up by one. Wrap the errors with `fmt.Errorf` and `%w`, the AOS written with `time.RFC3339`.

`List` builds a slice of its own with `make` and `append`, then sorts it with `slices.SortFunc`, comparing AOS with its `Compare` method and, when that gives 0, the IDs with `strings.Compare`. If one pass gets more than one winner, the check and the insert aren't under one lock.
