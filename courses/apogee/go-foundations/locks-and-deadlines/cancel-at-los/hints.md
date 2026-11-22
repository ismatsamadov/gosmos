`for range chunks` stops only when the channel closes. While `receive` waits for the next chunk, what else should it be watching?

Loop over a `select` with two cases: `<-ctx.Done()`, and a receive from `chunks`. The receive's two-value form, `_, ok := <-chunks`, sets `ok` to `false` once the channel is closed and empty.
