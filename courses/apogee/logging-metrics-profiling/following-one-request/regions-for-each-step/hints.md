The `select` can return from `runPass` in one of its cases, so it can't move into a function for `trace.WithRegion` as it is. Start the `wait for worker` region with `trace.StartRegion` right before the `select`, and call its `End` in each case, before anything else happens there.

The call to `fly` fits in a function: give `trace.WithRegion` a function literal that calls `u.fly` and stores the error in a variable declared before it, then check that variable after the region as before.

If the test says a region is never ended, its `End` is missing, or called in a goroutine of its own. A deferred `End` for the wait would only run when `runPass` returns, so the region would hold the whole flight too.
