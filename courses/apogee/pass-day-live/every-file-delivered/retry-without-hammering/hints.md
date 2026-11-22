Retry inside the event's own goroutine: a loop of at most 8 attempts, each one a new request with its timestamp and signature made from `w.Now()` then. A transport error has no status and is retried, like 408, 429 and 500 to 599; any other status ends the loop at once.

The wait after failed attempt n is 2 s doubled n−1 times, at most 5 minutes, times 0.5 plus half of `rand.Float64()` from `math/rand/v2`. Wait for it with a timer in a `select` on the rehearsal's context, as a pass waits for its AOS, so a stop ends the wait without an `UNDELIVERED` line.

A wait at the very top of its range means the factor wasn't drawn from [0.5, 1.0). Keep each undelivered event's file, customer ID and last answer under a lock, and write them sorted by file after the week, the file there even when nothing failed.
