Each variable is one `expvar.Publish` call with its name and an `expvar.Func`: a `func() any` that returns the value to show. Put all three calls inside the function given to `publishOnce.Do`, in place of the `_, _ =` line.

The three functions return `m.Passes.Snapshot()`, `m.AckWait.Snapshot()` and `runtime.NumGoroutine()`: call them inside each `func() any`, so each request reads them afresh.

If one more pass doesn't show in `passes_total`, the value was worked out once, when `publish` ran: return the call itself from the function, not a snapshot taken before it.
