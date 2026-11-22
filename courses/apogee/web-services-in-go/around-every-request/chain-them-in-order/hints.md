Loop over `mws` by index, from its last element down to the first, and each time replace `h` with the handler that the middleware at that index returns when it wraps `h`.

Each element of `mws` is a function: call it with the handler built so far, and it returns that handler wrapped. After the loop, `h` is the whole chain.

If the trace reads `c b a handler a b c`, your loop runs from the first middleware to the last: run it the other way. With no middlewares the loop never runs, and `h` comes back as it was.
