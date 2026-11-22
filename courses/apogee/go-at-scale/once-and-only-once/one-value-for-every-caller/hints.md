`sync.OnceValues` takes a function that returns two values and gives back a function of the same type, which is what `networkOnce` returns.

Pass `load` to `sync.OnceValues` and return what it gives you, with `"sync"` imported.

If the second table's server is never called, both tables share one function. If a failed load is tried again, the function `networkOnce` returns calls `load` itself.
