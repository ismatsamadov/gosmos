Replace `for h := range in` with a plain `for` loop around a `select` with two cases: `<-ctx.Done()`, which returns, and a receive from `in`.

When the receive gives `ok` false, `in` is closed: return, and the deferred `close(out)` still runs. Otherwise parse the line as before, and send the frame with a second `select` beside `<-ctx.Done()`.

A goroutine left after cancel means a send still has no `ctx.Done()` case. If only the case canceled before the first line fails, the receive from `in` is still outside a `select`.
