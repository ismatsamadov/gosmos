Create the file first with `os.Create`, and return its error before anything else happens: with no file, neither the trace nor `work` should start.

`trace.Start(f)` returns an error when it can't start, for one when another trace is already running; close the file and return it then. Otherwise run `work`, call `trace.Stop()`, which returns once the whole trace is written, and return what `f.Close()` returns.

If the test says a later `trace.Start` failed, `trace.Stop()` was never called. If the file doesn't start with `go 1.`, the trace went somewhere other than the file you created.
