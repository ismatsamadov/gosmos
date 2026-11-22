`Do` is two calls you already have: `c.send(f)` for the frame, then `recvCtx(ctx, c)` for its answer.

Check `send`'s error before reading: if it isn't `nil`, return it with an empty `Frame`.

If the controller that stops reading fails after 10 s with a read error, `Do` read after a failed send. If the pass ending after 2 s fails after 5 s instead, the answer is read with `c.recv()` rather than `recvCtx`.
