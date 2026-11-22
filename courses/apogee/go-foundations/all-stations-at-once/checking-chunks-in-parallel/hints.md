Each worker needs to know which chunk to check next. What could all the workers receive from, so that every index goes to exactly one of them?

You need a `chan int` of indexes, `for range workers` with `wg.Go` to start the workers, then `close` and `wg.Wait()` once every index has been sent.

Each worker's body is `for i := range ___ { ok[i] = chunkOK(___, ___) }`.
