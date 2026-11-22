Make an unbuffered `jobs` channel and start `workers` goroutines that range over it: each calls `parse` on the job's line and sends the `result` into the job's `res`, which never waits because it has room for one.

The dispatching goroutine loops on `receive(ctx, in)`: for each line it makes `res` with room for one, sends it on the queue, then sends the job, and it closes `jobs` and the queue when it stops.

The sending goroutine receives each `res` from the queue in turn, then its result, and sends the frame unless `err` is set. If frames still come out shuffled, results go to `out` in the order they finish rather than in the queue's order.
