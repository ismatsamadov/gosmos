The read that waits for the banner can't watch a context, but the connection can be closed under it: then the read returns an error at once. What has to happen is a `Close` at the moment the pass's context ends.

`context.AfterFunc(ctx, f)` runs `f` in its own goroutine when `ctx` ends, and returns a function that cancels that. Right after the dial, register one that closes the connection, and call the function it returns with `defer`.

A call closed this way fails the banner read, so the pass is `failed no_answer` as before, and its LOS line comes at LOS. If the tests still say `rehearse still running`, check that the function you registered closes the call itself, and that the context you gave `AfterFunc` is the pass's, which ends at LOS.
