Drop the slice and the goroutine: `feed` itself loops over `contacts` and sends each one into `q`, so it returns only once the loader's queue has taken the last.

Make each send a `select` with two cases: `q <- c`, which goes on to the next contact, and `<-ctx.Done()`, which returns `ctx.Err()`. After the loop, return `nil`.

If the test says `feed` returned before the test canceled, something still sends from a goroutine of its own. If it is still waiting after the cancel, the send has no `ctx.Done()` case beside it.
