Keep the `next` functions in a slice, one per sequence: for each `seq`, `next, stop := iter.Pull(seq)`, then `defer` its `stop`, so every source stops however the merge ends.

Push `head[T]{v, i}` for each first value `next()` gives you. Then loop: `Pop` the smallest head, `yield` its `v`, and call `nexts[h.i]()` for that sequence's next value.

If a stream was never stopped after the loop broke, a `stop` isn't deferred. If values come out of order, a sequence's next value isn't pushed once its head has been yielded.
