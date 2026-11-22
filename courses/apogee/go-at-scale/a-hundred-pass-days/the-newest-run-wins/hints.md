Keep the run in progress's cancel function in a variable outside the loop, starting as a function that does nothing. For each event, call it, then make the new run's context with `context.WithCancel(ctx)` and keep its cancel in that variable.

Move `r.runs.Add(1)` out of the goroutine, beside `r.now()` and `r.pub.Current()`, and hand the goroutine the number. Start the goroutine with a `sync.WaitGroup`'s `Go`, and defer a function that waits for the group, then cancels the last run.

If the 00:00:21 check finds the tick's run published, its context never ended; if run 30 isn't the event at 00:00:30, runs are still numbered when they end; if the last test finds nothing published, `Run` returned before its last run did.
