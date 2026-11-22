The test's first message says the stations must run at the same time. Which of the three functions decides that, and what do the stations' goroutines share?

`runDay` needs a `sync.WaitGroup` and a `sync.Mutex`; `runPass` needs `context.WithDeadline` and a deferred cancel; `runStation` needs `errors.Is` and the index of the pass that failed.

In `runStation`, range over `day` with its index `i`; after `rep.Failed++`, the branch for `errors.Is(err, ErrStowed)` sets `rep.Skipped` to ___, the number of passes after index `i`, then ends the loop.
