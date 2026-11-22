Start one goroutine per contact and wait for them all before you print the count: a `sync.WaitGroup`'s `Go` method starts a function and counts it in one call, and `Wait` returns once every one has returned. Each goroutine waits for its own AOS, so passes at the same moment overlap by themselves.

A contact sleeps for `c.AOS.Sub(w.Now())`, prints its AOS line, then makes its pass's context with `context.WithTimeout(ctx, c.LOS.Sub(w.Now()))`, calls the `cancel` it gets with `defer`, and waits on the context's `Done` channel, which closes at LOS. Write each line with one `fmt.Fprintf` on `w.Out`.

If the tests say `rehearse still running`, a wait or a deadline was a 2029 instant on the `time` package's clock: use durations from `w.Now()` only. If they say `rehearse()` returned before the last LOS, it didn't wait for every contact to end. Many contacts end in the same instant, so guard what they count together with a `sync.Mutex`.
