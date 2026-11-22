`sync.OnceFunc` takes a function with no arguments and no results and returns one of the same type, so its result can go straight into `r.stop`.

Pass the function that closes `r.done` to `sync.OnceFunc` and store what it returns in `r.stop`, with `"sync"` imported.

If the test still reports a panic from a second `Stop`, `r.stop` is still the plain function, or `OnceFunc` is called inside `Stop` itself, which makes a new once on every call.
