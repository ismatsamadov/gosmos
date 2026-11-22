Wrap the dial in a `for` loop with no condition. When `l.Dial` returns a connection, return it with a nil error; that is the only way out besides the context.

After a failed dial, call `l.Wait(ctx, time.Second)` and look at what it returns. Nil means the second has passed, so the loop goes round and dials again.

When the wait returns an error, return nil and that error, not the dial's. If the test counts a third dial after the stop, the loop went round once more after the wait failed.
