Start with `next, stop := iter.Pull(b)` and `defer stop()` on the line after, then range over `a` as you would any sequence.

Inside the loop, call `next()`. If it returns false, `b` has ended: return. Otherwise yield the value of `a` and the value of `b`, and return if `yield` gives back false.

A test saying b's loop had not ended means `stop` never ran: without it, `b` waits inside its loop after `Zip` returns. A pair in the wrong order means `b`'s value came first. If `b` was read to its end for two rows, it was collected first, not pulled one value at a time.
