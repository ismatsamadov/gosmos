`merge` has to return `out` before anything is sent on it, since its caller is the only receiver. So which parts of the work must run in goroutines?

Use a `sync.WaitGroup`: start one forwarding goroutine per feed with `wg.Go`, and close `out` from one more goroutine that calls `wg.Wait()` first.

Each forwarder is `for seq := range ___ { out <- seq }`; the closer is `go func() { wg.Wait(); ___ }()`.
